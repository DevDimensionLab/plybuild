package taskwork

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// TaskContentStorage owns bounded, physical reads and publish-once content. Its
// fault boundary is injected per dependency instance, never through environment.
type TaskContentStorage struct {
	scope        *taskContentScope
	fault        func(string) error
	readCaptured func(work workspace.ProjectWork, kind, digest string) ([]byte, error)
}

func (s *TaskContentStorage) check(stage string) error {
	if s == nil {
		return contentError("task_content_observation_unknown", "content storage dependency is missing", nil)
	}
	if s.fault != nil {
		return s.fault(stage)
	}
	return nil
}
func taskContentPath(work workspace.ProjectWork, kind, digest string) string {
	suffix := ""
	if kind == "manifests" || kind == "requests" {
		suffix = ".json"
	}
	return filepath.Join(work.Content(), kind, "sha256", strings.TrimPrefix(digest, "sha256:")+suffix)
}
func (s *TaskContentStorage) ReadSource(path string, limit int) ([]byte, error) {
	if e := s.check("source-read"); e != nil {
		return nil, e
	}
	return contentReadBounded(path, limit, false)
}
func contentDirectories(work workspace.ProjectWork, kind string) []string {
	return []string{filepath.Dir(work.Content()), work.Content(), filepath.Join(work.Content(), kind), filepath.Join(work.Content(), kind, "sha256")}
}
func (s *TaskContentStorage) directories(work workspace.ProjectWork, kind string, create bool) error {
	if create {
		if e := work.Create(); e != nil {
			return e
		}
	}
	if e := contentPhysical(work.Dir); e != nil {
		return e
	}
	for _, dir := range contentDirectories(work, kind) {
		info, e := os.Lstat(dir)
		if errors.Is(e, fs.ErrNotExist) && create {
			if e = s.check("mkdir"); e != nil {
				return e
			}
			if e = os.Mkdir(dir, 0o700); e != nil {
				return e
			}
			if e = workspace.PlatformSyncDirectory(filepath.Dir(dir)); e != nil {
				return e
			}
			info, e = os.Lstat(dir)
		}
		if e != nil {
			return e
		}
		if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 || info.Mode().Perm() != 0o700 {
			return contentError("task_content_integrity_conflict", "managed directory must be physical and mode 0700: "+dir, nil)
		}
	}
	return nil
}
func (s *TaskContentStorage) Read(work workspace.ProjectWork, kind, digest string) (out []byte, err error) {
	defer func() {
		if err != nil {
			err = contentArtifactError(work, kind, digest, "read content", err)
		}
	}()
	if s != nil && s.readCaptured != nil {
		return s.readCaptured(work, kind, digest)
	}
	if !digestPattern.MatchString(digest) {
		return nil, contentError("task_content_integrity_conflict", "invalid managed digest", nil)
	}
	if e := s.check("managed-read"); e != nil {
		return nil, e
	}
	if e := s.directories(work, kind, false); e != nil {
		return nil, e
	}
	limit := taskManifestLimit
	if kind == "objects" {
		limit = taskDocumentLimit
	}
	if kind == integrationEvidenceKind {
		limit = integrationEvidenceLimit
	}
	b, e := contentReadBounded(taskContentPath(work, kind, digest), limit, true)
	if e != nil {
		return nil, e
	}
	if digestTaskBytes(b) != digest {
		return nil, contentError("task_content_integrity_conflict", "managed bytes do not match their digest", nil)
	}
	return b, nil
}
func (s *TaskContentStorage) Publish(work workspace.ProjectWork, kind string, b []byte) (string, error) {
	if s != nil && s.readCaptured != nil {
		return "", contentError("task_content_read_only", "captured content is read-only", nil)
	}
	digest := digestTaskBytes(b)
	if e := s.directories(work, kind, true); e != nil {
		return "", e
	}
	dest := taskContentPath(work, kind, digest)
	if _, e := os.Lstat(dest); e == nil {
		old, e := s.Read(work, kind, digest)
		if e != nil {
			return "", e
		}
		if string(old) != string(b) {
			return "", contentError("task_content_integrity_conflict", "immutable destination differs", nil)
		}
		return digest, s.Sync(work, kind, digest)
	} else if !errors.Is(e, fs.ErrNotExist) {
		return "", e
	}
	if e := s.check("object-temp-open"); e != nil {
		return "", e
	}
	f, e := os.CreateTemp(filepath.Dir(dest), ".task-content-*")
	if e != nil {
		return "", e
	}
	temp := f.Name()
	renamed := false
	defer func() {
		_ = f.Close()
		if !renamed {
			_ = os.Remove(temp)
		}
	}()
	if e = s.check("object-write"); e != nil {
		return "", e
	}
	n, e := f.Write(b)
	if e != nil {
		return "", e
	}
	if n != len(b) {
		return "", io.ErrShortWrite
	}
	if e = f.Chmod(0o400); e != nil {
		return "", e
	}
	if e = s.check("object-sync"); e != nil {
		return "", e
	}
	if e = f.Sync(); e != nil {
		return "", e
	}
	if e = f.Close(); e != nil {
		return "", e
	}
	if e = s.check("object-rename"); e != nil {
		return "", e
	}
	// The work-item lock serializes the destination check and rename.
	if _, e = os.Lstat(dest); !errors.Is(e, fs.ErrNotExist) {
		return "", contentError("task_content_integrity_conflict", "immutable destination appeared during publication", e)
	}
	if e = os.Rename(temp, dest); e != nil {
		return "", e
	}
	renamed = true
	if e = s.check("object-directory-sync"); e != nil {
		return "", e
	}
	if e = workspace.PlatformSyncDirectory(filepath.Dir(dest)); e != nil {
		return "", e
	}
	return digest, nil
}
func (s *TaskContentStorage) Sync(work workspace.ProjectWork, kind, digest string) error {
	if s != nil && s.readCaptured != nil {
		return contentError("task_content_read_only", "captured content is read-only", nil)
	}
	if _, e := s.Read(work, kind, digest); e != nil {
		return e
	}
	f, e := os.Open(taskContentPath(work, kind, digest))
	if e != nil {
		return e
	}
	syncErr := f.Sync()
	closeErr := f.Close()
	if e = firstError(syncErr, closeErr, s.check("recovery-sync")); e != nil {
		return e
	}
	for _, d := range contentDirectories(work, kind) {
		if e = workspace.PlatformSyncDirectory(d); e != nil {
			return e
		}
	}
	return workspace.PlatformSyncDirectory(work.Dir)
}

func contentManifestOperation(kind string) string {
	switch kind {
	case "WorkspaceTaskProblemRevision@1":
		return "problem_record"
	case "WorkspaceTaskSpecRevision@1", "WorkspaceTaskSpecRevision@2":
		return "spec_record"
	case "WorkspaceTaskSpecAssessment@1":
		return "spec_assess"
	case "WorkspaceTaskSolutionSelection@1":
		return "spec_select"
	}
	return ""
}
func readContentManifest(s *TaskContentStorage, work workspace.ProjectWork, digest string) (out canonicaljson.Object, err error) {
	defer func() {
		if err != nil {
			err = contentArtifactError(work, "manifests", digest, "read publication", err)
		}
	}()
	b, e := s.Read(work, "manifests", digest)
	if e != nil {
		return nil, e
	}
	v, e := canonicaljson.DecodeStrict(b)
	if e != nil {
		return nil, contentError("task_content_invalid_input", "malformed publication: "+e.Error(), e)
	}
	m := contentFields(v)
	op := contentManifestOperation(contentString(m, "kind"))
	if op == "spec_select" && contentString(m, "action") == "withdraw" {
		op = "spec_withdraw"
	}
	if op == "" {
		return nil, contentError("task_content_reader_incompatible", "unsupported manifest kind for reader "+taskContentReaderCapability, nil)
	}
	o, e := decodeStoredTaskContent(b, op, true, true)
	if e != nil {
		return nil, e
	}
	canonical, e := canonicaljson.Marshal(o)
	if e != nil || string(canonical) != string(b) {
		return nil, contentError("task_content_integrity_conflict", "manifest is not canonical", e)
	}
	return o, nil
}
func contentManifestOwned(r WorkItemRegistry, id TaskID, digest string) bool {
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == id && p.ManifestSHA256 == digest {
			return true
		}
	}
	for _, p := range r.TaskSpecRevisions {
		if p.TaskID == id && p.ManifestSHA256 == digest {
			return true
		}
	}
	for _, p := range r.TaskSpecAssessments {
		if p.TaskID == id && p.ManifestSHA256 == digest {
			return true
		}
	}
	return false
}
func contentResolveDocuments(d Dependencies, work workspace.ProjectWork, r WorkItemRegistry, id TaskID, inputs []canonicaljson.Value) ([]canonicaljson.Value, map[string][]byte, error) {
	out := []canonicaljson.Value{}
	objects := map[string][]byte{}
	for _, v := range inputs {
		m := contentFields(v)
		src := contentFields(m["source"])
		var descriptor map[string]canonicaljson.Value
		var bytes []byte
		var err error
		if contentString(src, "kind") == "snapshot" {
			digest := contentString(src, "manifest_sha256")
			owner := id
			if explicit := contentString(src, "task_id"); explicit != "" {
				owner = TaskID(explicit)
			}
			if !contentManifestOwned(r, owner, digest) {
				return nil, nil, contentError("task_content_integrity_conflict", "snapshot manifest is not registered to this Task", nil)
			}
			manifest, e := readContentManifest(d.TaskContent, work, digest)
			if e != nil {
				return nil, nil, e
			}
			for _, doc := range contentArray(contentFields(manifest), "documents") {
				candidate := contentFields(doc)
				if contentString(candidate, "id") == contentString(src, "document_id") {
					descriptor = candidate
					break
				}
			}
			if descriptor == nil {
				return nil, nil, contentError("task_content_missing", "snapshot document does not exist", nil)
			}
			bytes, err = d.TaskContent.Read(work, "objects", contentString(descriptor, "sha256"))
			if err != nil {
				return nil, nil, err
			}
			if contentInt(descriptor, "size_bytes") != len(bytes) || contentString(descriptor, "locator") != taskContentPath(work, "objects", contentString(descriptor, "sha256")) {
				return nil, nil, contentError("task_content_integrity_conflict", "snapshot descriptor differs", nil)
			}
			descriptor["id"] = m["id"]
		} else {
			bytes, err = d.TaskContent.ReadSource(contentString(src, "locator"), taskDocumentLimit)
			if err != nil {
				return nil, nil, err
			}
			if digestTaskBytes(bytes) != contentString(src, "sha256") || len(bytes) != contentInt(src, "size_bytes") {
				return nil, nil, contentError("task_content_integrity_conflict", "source hash or size differs", nil)
			}
			if src["git_provenance"] != nil {
				if err = observeTaskContentGit(d, src["git_provenance"], bytes); err != nil {
					return nil, nil, err
				}
			}
			descriptor = map[string]canonicaljson.Value{"id": m["id"], "locator": taskContentPath(work, "objects", digestTaskBytes(bytes)), "sha256": digestTaskBytes(bytes), "size_bytes": int64(len(bytes)), "media_type": src["media_type"], "provenance": contentObject(map[string]canonicaljson.Value{"origin_locator": src["locator"], "git": src["git_provenance"]})}
		}
		if err = taskDocumentBytesValid(bytes); err != nil {
			return nil, nil, err
		}
		objects[digestTaskBytes(bytes)] = bytes
		out = append(out, contentObject(descriptor))
	}
	return out, objects, nil
}

func validateTaskContentClosure(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, sync bool) error {
	for _, pub := range r.TaskContentPublications {
		if !s.includesTask(pub.TaskID) {
			continue
		}
		manifest, e := readContentManifest(s, work, pub.OutcomeRef.ManifestSHA256)
		if e != nil {
			return e
		}
		m := contentFields(manifest)
		if contentString(m, "task_id") != string(pub.TaskID) || contentString(m, "publication_key") != pub.PublicationKey || contentString(m, "source_draft_sha256") != pub.IntentSHA256 {
			return contentError("task_content_integrity_conflict", "publication and manifest binding differ", nil)
		}
		request, e := s.Read(work, "requests", pub.RequestSHA256)
		if e != nil {
			return e
		}
		if e = validateContentPublicationChain(s, work, r, pub, m, request); e != nil {
			var classified *TaskContentError
			if errors.As(e, &classified) {
				return contentArtifactError(work, "requests", pub.RequestSHA256, pub.Operation, e)
			}
			return contentArtifactError(work, "requests", pub.RequestSHA256, pub.Operation, contentError("task_content_integrity_conflict", "preserved publication chain differs: "+e.Error(), e))
		}
		for _, v := range contentArray(m, "documents") {
			doc := contentFields(v)
			digest := contentString(doc, "sha256")
			b, e := s.Read(work, "objects", digest)
			if e != nil {
				return e
			}
			if e = taskDocumentBytesValid(b); e != nil {
				return e
			}
			if len(b) != contentInt(doc, "size_bytes") || contentString(doc, "locator") != taskContentPath(work, "objects", digest) {
				return contentError("task_content_integrity_conflict", "document descriptor differs from managed object", nil)
			}
			if sync {
				if e = s.Sync(work, "objects", digest); e != nil {
					return e
				}
			}
		}
		if sync {
			for _, pair := range [][2]string{{"manifests", pub.OutcomeRef.ManifestSHA256}, {"requests", pub.RequestSHA256}} {
				if e = s.Sync(work, pair[0], pair[1]); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// validateContentPublicationChain checks a registered publication against its
// preserved request and its manifest. The rules run in a fixed order and the
// first one that fails answers. Classified errors come from decoding the
// request and from reading stored manifests; the callers class the others.
func validateContentPublicationChain(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, pub TaskContentPublication, manifest map[string]canonicaljson.Value, request []byte) error {
	draft, err := publicationRequestDraft(pub, request)
	if err != nil {
		return err
	}
	if err = validatePublicationIntent(r, pub); err != nil {
		return err
	}
	if err = validateRequestFields(draft, manifest); err != nil {
		return err
	}
	if contentString(manifest, "recorded_at_utc") != pub.RecordedAtUTC {
		return fmt.Errorf("publication time differs")
	}
	previous, err := registeredPredecessor(r, pub, manifest)
	if err != nil {
		return err
	}
	if !contentEqual(previous, manifest["previous"]) {
		return fmt.Errorf("predecessor differs from registered chain")
	}
	if err = validateOwnedReferences(r, pub.TaskID, manifest); err != nil {
		return err
	}
	if err = validateSelectedSolution(s, work, r, pub.TaskID, manifest); err != nil {
		return err
	}
	if origin := contentString(contentFields(draft["origin"]), "kind"); origin == "task_create_summary" || origin == "legacy_summary_import" {
		return validateGeneratedProblem(work, r, pub.TaskID, draft, manifest)
	}
	return validatePublishedDocuments(s, work, r, pub.TaskID, draft, manifest)
}

// publicationRequestDraft decodes the preserved request as the canonical draft
// of the publication's operation and returns its fields. A Task's creation
// preserved an internal problem draft.
func publicationRequestDraft(pub TaskContentPublication, request []byte) (map[string]canonicaljson.Value, error) {
	created := pub.Operation == "task_create"
	operation := pub.Operation
	if created {
		operation = "problem_record"
	}
	draft, err := decodeStoredTaskContent(request, operation, false, created)
	if err != nil {
		return nil, err
	}
	canonical, err := canonicaljson.Marshal(draft)
	if err != nil || string(canonical) != string(request) {
		return nil, fmt.Errorf("request is not canonical")
	}
	return contentFields(draft), nil
}

// The intent of a Task's creation is the create intent of the registered Task;
// the intent of every other publication is its request.
func validatePublicationIntent(r WorkItemRegistry, pub TaskContentPublication) error {
	if pub.Operation != "task_create" {
		if pub.RequestSHA256 != pub.IntentSHA256 {
			return fmt.Errorf("request and intent digests differ")
		}
		return nil
	}
	t, _ := findTask(r, pub.TaskID)
	if t == nil {
		return fmt.Errorf("missing Task")
	}
	intent, err := taskCreateIntent(*t)
	if err != nil || digestTaskBytes(intent) != pub.IntentSHA256 {
		return fmt.Errorf("create intent differs")
	}
	return nil
}

// The manifest repeats the fields of its request. Fields are compared by name,
// so of several differing fields the error names the first.
func validateRequestFields(draft, manifest map[string]canonicaljson.Value) error {
	for _, field := range slices.Sorted(maps.Keys(draft)) {
		if member, repeated := manifestMember(field); repeated && !contentEqual(draft[field], manifest[member]) {
			return fmt.Errorf("request field %s differs", member)
		}
	}
	return nil
}

// manifestMember names the manifest member that repeats a request field. The
// predecessor that the request expected is the manifest's "previous"; the
// kind, the registry upgrade and the documents are not repeated.
func manifestMember(field string) (member string, repeated bool) {
	switch field {
	case "kind", "registry_upgrade", "documents":
		return "", false
	case "expected_previous", "expected_previous_assessment", "expected_previous_selection":
		return "previous", true
	}
	return field, true
}

// registeredPredecessor checks that the publication's outcome identifies its
// manifest, and returns the reference to the record registered just before it
// in its chain: nil for the first one, the last match if several are
// registered.
func registeredPredecessor(r WorkItemRegistry, pub TaskContentPublication, manifest map[string]canonicaljson.Value) (canonicaljson.Value, error) {
	switch pub.OutcomeRef.Kind {
	case "problem":
		return problemPredecessor(r, pub.TaskID, pub.OutcomeRef, manifest)
	case "spec":
		return specPredecessor(r, pub.TaskID, pub.OutcomeRef, manifest)
	case "assessment":
		return assessmentPredecessor(r, pub.TaskID, pub.OutcomeRef, manifest)
	case "selection":
		return selectionPredecessor(r, pub.TaskID, pub.OutcomeRef, manifest)
	}
	return nil, fmt.Errorf("unknown publication outcome")
}

// pointsTo reports whether an optional member of an outcome is set to value.
func pointsTo[T comparable](member *T, value T) bool {
	return member != nil && *member == value
}

// A Task's problems form one chain by revision.
func problemPredecessor(r WorkItemRegistry, task TaskID, outcome TaskContentOutcomeRef, manifest map[string]canonicaljson.Value) (canonicaljson.Value, error) {
	if !pointsTo(outcome.Revision, contentInt(manifest, "revision")) {
		return nil, fmt.Errorf("problem revision differs")
	}
	var previous canonicaljson.Value
	for _, v := range r.TaskProblemRevisions {
		if v.TaskID == task && v.Revision == *outcome.Revision-1 {
			previous = contentRefValue(TaskRevisionRef{v.Revision, v.ManifestSHA256})
		}
	}
	return previous, nil
}

// Each Spec of a Task forms a chain by revision.
func specPredecessor(r WorkItemRegistry, task TaskID, outcome TaskContentOutcomeRef, manifest map[string]canonicaljson.Value) (canonicaljson.Value, error) {
	if !pointsTo(outcome.SpecID, contentString(manifest, "spec_id")) || !pointsTo(outcome.Revision, contentInt(manifest, "revision")) {
		return nil, fmt.Errorf("Spec identity differs")
	}
	var previous canonicaljson.Value
	for _, v := range r.TaskSpecRevisions {
		if v.TaskID == task && v.SpecID == *outcome.SpecID && v.Revision == *outcome.Revision-1 {
			previous = contentRefValue(TaskRevisionRef{v.Revision, v.ManifestSHA256})
		}
	}
	return previous, nil
}

// The assessments of one Spec revision form a chain by ordinal. The outcome's
// revision is that of the assessed Spec.
func assessmentPredecessor(r WorkItemRegistry, task TaskID, outcome TaskContentOutcomeRef, manifest map[string]canonicaljson.Value) (canonicaljson.Value, error) {
	assessed := contentInt(contentFields(manifest["spec"]), "revision")
	if !pointsTo(outcome.ID, contentString(manifest, "id")) || !pointsTo(outcome.SpecID, contentString(manifest, "spec_id")) || !pointsTo(outcome.Revision, assessed) {
		return nil, fmt.Errorf("assessment identity differs")
	}
	var previous canonicaljson.Value
	for _, v := range r.TaskSpecAssessments {
		if v.TaskID == task && v.SpecID == *outcome.SpecID && v.SpecRevision == *outcome.Revision && v.Ordinal == contentInt(manifest, "ordinal")-1 {
			previous = contentRefValue(TaskDecisionRef{v.ID, v.ManifestSHA256})
		}
	}
	return previous, nil
}

// A Task's selections and withdrawals form one chain by ordinal.
func selectionPredecessor(r WorkItemRegistry, task TaskID, outcome TaskContentOutcomeRef, manifest map[string]canonicaljson.Value) (canonicaljson.Value, error) {
	if !pointsTo(outcome.ID, contentString(manifest, "id")) {
		return nil, fmt.Errorf("selection identity differs")
	}
	var previous canonicaljson.Value
	for _, v := range r.TaskSolutionSelections {
		if v.TaskID == task && v.Ordinal == contentInt(manifest, "ordinal")-1 {
			previous = contentRefValue(TaskDecisionRef{v.ID, v.ManifestSHA256})
		}
	}
	return previous, nil
}

// The problem and the Spec that a manifest refers to are registered to its
// Task.
func validateOwnedReferences(r WorkItemRegistry, task TaskID, manifest map[string]canonicaljson.Value) error {
	if problem := valueRevision(manifest["problem"]); problem != nil && !problemRegistered(r, task, *problem) {
		return fmt.Errorf("problem reference is not owned by Task")
	}
	if spec := valueRevision(manifest["spec"]); spec != nil && !specRegistered(r, task, contentString(manifest, "spec_id"), *spec) {
		return fmt.Errorf("Spec reference is not owned by Task")
	}
	return nil
}

func problemRegistered(r WorkItemRegistry, task TaskID, ref TaskRevisionRef) bool {
	return slices.ContainsFunc(r.TaskProblemRevisions, func(v TaskProblemReference) bool {
		return v.TaskID == task && v.Revision == ref.Revision && v.ManifestSHA256 == ref.ManifestSHA256
	})
}

func specRegistered(r WorkItemRegistry, task TaskID, specID string, ref TaskRevisionRef) bool {
	return slices.ContainsFunc(r.TaskSpecRevisions, func(v TaskSpecReference) bool {
		return v.TaskID == task && v.SpecID == specID && v.Revision == ref.Revision && v.ManifestSHA256 == ref.ManifestSHA256
	})
}

// A selected solution names a registered assessment that found the selected
// Spec revision ready, and that Spec is the Task's answer to the selected
// problem. A withdrawal has no solution.
func validateSelectedSolution(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, task TaskID, manifest map[string]canonicaljson.Value) error {
	solution := contentFields(manifest["solution"])
	if len(solution) == 0 {
		return nil
	}
	specID := contentString(solution, "spec_id")
	problem, spec, assessment := valueRevision(solution["problem"]), valueRevision(solution["spec"]), valueDecision(solution["assessment"])
	registered := slices.ContainsFunc(r.TaskSpecAssessments, func(v TaskAssessmentReference) bool {
		return v.TaskID == task && v.SpecID == specID && v.SpecRevision == spec.Revision && v.ID == assessment.ID && v.ManifestSHA256 == assessment.ManifestSHA256
	})
	if !registered {
		return fmt.Errorf("selection assessment is not registered")
	}
	stored, err := readContentManifest(s, work, assessment.ManifestSHA256)
	if err != nil {
		return err
	}
	assessed := contentFields(stored)
	if contentString(assessed, "outcome") != "ready" || !contentTypedEqual(valueRevision(assessed["spec"]), spec) {
		return fmt.Errorf("selection does not bind a ready assessment")
	}
	if stored, err = readContentManifest(s, work, spec.ManifestSHA256); err != nil {
		return err
	}
	selected := contentFields(stored)
	if contentString(selected, "task_id") != string(task) || contentString(selected, "spec_id") != specID || !contentTypedEqual(valueRevision(selected["problem"]), problem) {
		return fmt.Errorf("selection problem and Spec differ")
	}
	return nil
}

// A problem that Ply wrote from a Task's summary has the title and the
// description of the registered Task and one document, the summary as
// Markdown. The import of a legacy summary also keeps that summary's digest.
// The request of a generated problem has no documents to compare.
func validateGeneratedProblem(work workspace.ProjectWork, r WorkItemRegistry, task TaskID, draft, manifest map[string]canonicaljson.Value) error {
	t, _ := findTask(r, task)
	if t == nil || contentString(draft, "title") != t.Title || contentString(draft, "summary") != t.Description {
		return fmt.Errorf("generated problem differs from the original Task summary")
	}
	docs := contentArray(manifest, "documents")
	if len(docs) != 1 || !contentEqual(docs[0], generatedProblemDocument(work, *t)) {
		return fmt.Errorf("generated problem document differs from the original Task summary")
	}
	origin := contentFields(draft["origin"])
	if contentString(origin, "kind") != "legacy_summary_import" {
		return nil
	}
	original, err := contentCanonical(map[string]any{"task_id": t.ID, "title": t.Title, "description": t.Description})
	if err != nil || contentString(origin, "legacy_summary_sha256") != digestTaskBytes(original) {
		return fmt.Errorf("legacy summary digest differs")
	}
	return nil
}

// generatedProblemDocument describes the managed Markdown object that holds a
// Task's summary.
func generatedProblemDocument(work workspace.ProjectWork, t TaskRecord) canonicaljson.Object {
	body := []byte("# " + t.Title + "\n\n" + t.Description + "\n")
	digest := digestTaskBytes(body)
	return contentObject(map[string]canonicaljson.Value{"id": "problem", "locator": taskContentPath(work, "objects", digest), "sha256": digest, "size_bytes": int64(len(body)), "media_type": "text/markdown", "provenance": contentObject(map[string]canonicaljson.Value{"origin_locator": nil, "git": nil})})
}

// The manifest has one document for each document of the request, in the same
// order and under the same ID; the first one that differs answers.
func validatePublishedDocuments(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, task TaskID, draft, manifest map[string]canonicaljson.Value) error {
	published := contentArray(manifest, "documents")
	requested := contentArray(draft, "documents")
	if len(published) != len(requested) {
		return fmt.Errorf("document set differs from request")
	}
	for i := range requested {
		if err := validatePublishedDocument(s, work, r, task, requested[i], published[i]); err != nil {
			return err
		}
	}
	return nil
}

// A published document has the ID of its request and follows the request's
// source: a file, or else a snapshot.
func validatePublishedDocument(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, task TaskID, requested, published canonicaljson.Value) error {
	in, doc := contentFields(requested), contentFields(published)
	if contentString(in, "id") != contentString(doc, "id") {
		return fmt.Errorf("document ID differs")
	}
	source := contentFields(in["source"])
	if contentString(source, "kind") == "file" {
		return validateFileDocument(source, doc)
	}
	return validateSnapshotDocument(s, work, r, task, source, published)
}

// A document from a file keeps the content metadata and the provenance that
// its request gave.
func validateFileDocument(source, doc map[string]canonicaljson.Value) error {
	for _, k := range []string{"sha256", "size_bytes", "media_type"} {
		if !contentEqual(source[k], doc[k]) {
			return fmt.Errorf("document metadata differs")
		}
	}
	provenance := contentFields(doc["provenance"])
	if !contentEqual(source["locator"], provenance["origin_locator"]) || !contentEqual(source["git_provenance"], provenance["git"]) {
		return fmt.Errorf("document provenance differs")
	}
	return nil
}

// A snapshot repeats the descriptor of a document in a manifest that the
// Task, or the Task that the source names, has registered. Only the document
// ID may be new.
func validateSnapshotDocument(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, task TaskID, source map[string]canonicaljson.Value, published canonicaljson.Value) error {
	digest := contentString(source, "manifest_sha256")
	owner := task
	if explicit := contentString(source, "task_id"); explicit != "" {
		owner = TaskID(explicit)
	}
	if !contentManifestOwned(r, owner, digest) {
		return fmt.Errorf("snapshot source is not owned by Task")
	}
	origin, err := readContentManifest(s, work, digest)
	if err != nil {
		return err
	}
	originals := documentsWithID(origin, contentString(source, "document_id"))
	if len(originals) == 0 {
		return fmt.Errorf("snapshot source document is missing")
	}
	for _, original := range originals {
		original["id"] = contentFields(published)["id"]
		if !contentEqual(contentObject(original), published) {
			return fmt.Errorf("snapshot descriptor differs")
		}
	}
	return nil
}

// documentsWithID returns the fields of the manifest's documents under an ID.
// A valid manifest has at most one.
func documentsWithID(manifest canonicaljson.Object, id string) []map[string]canonicaljson.Value {
	var found []map[string]canonicaljson.Value
	for _, raw := range contentArray(contentFields(manifest), "documents") {
		if doc := contentFields(raw); contentString(doc, "id") == id {
			found = append(found, doc)
		}
	}
	return found
}

func readRegisteredTaskManifest(s *TaskContentStorage, work workspace.ProjectWork, r WorkItemRegistry, digest string) (canonicaljson.Object, error) {
	m, err := readContentManifest(s, work, digest)
	if err != nil {
		return nil, err
	}
	for _, p := range r.TaskContentPublications {
		if p.OutcomeRef.ManifestSHA256 == digest {
			request, err := s.Read(work, "requests", p.RequestSHA256)
			if err != nil {
				return nil, err
			}
			if err = validateContentPublicationChain(s, work, r, p, contentFields(m), request); err != nil {
				var classified *TaskContentError
				if errors.As(err, &classified) {
					return nil, contentArtifactError(work, "requests", p.RequestSHA256, p.Operation, err)
				}
				return nil, contentArtifactError(work, "requests", p.RequestSHA256, p.Operation, contentError("task_content_integrity_conflict", "manifest publication chain differs: "+err.Error(), err))
			}
			return m, nil
		}
	}
	return nil, contentError("task_content_integrity_conflict", "manifest has no registered publication", nil)
}
