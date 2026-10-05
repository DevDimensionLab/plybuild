package workspace

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// TaskContentStorage owns bounded, physical reads and publish-once content. Its
// fault boundary is injected per dependency instance, never through environment.
type TaskContentStorage struct{ fault func(string) error }

func (s *TaskContentStorage) check(stage string) error {
	if s == nil {
		return contentError("task_content_observation_unknown", "content storage dependency is missing", nil)
	}
	if s.fault != nil {
		return s.fault(stage)
	}
	return nil
}
func taskContentRoot(root string) string {
	return filepath.Join(root, MarkerDirectory, "task-content", "v1")
}
func taskContentPath(root, kind, digest string) string {
	suffix := ""
	if kind == "manifests" || kind == "requests" {
		suffix = ".json"
	}
	if kind == "backups" {
		suffix = ".yaml"
	}
	return filepath.Join(taskContentRoot(root), kind, "sha256", strings.TrimPrefix(digest, "sha256:")+suffix)
}
func contentPhysical(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return contentError("task_content_invalid_input", "path must be absolute and clean", nil)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	if resolved != path {
		return contentError("task_content_integrity_conflict", "path contains a symlink: "+path, nil)
	}
	return nil
}
func contentReadBounded(path string, limit int, immutable bool) (data []byte, err error) {
	if err = contentPhysical(path); err != nil {
		return nil, err
	}
	before, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !before.Mode().IsRegular() || immutable && before.Mode().Perm() != 0o400 {
		return nil, contentError("task_content_integrity_conflict", "unexpected file type or immutable mode: "+path, nil)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()
	opened, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(before, opened) {
		return nil, contentError("task_content_integrity_conflict", "file changed before read", nil)
	}
	data, err = io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > limit {
		return nil, contentError("task_content_invalid_input", fmt.Sprintf("file exceeds %d bytes", limit), nil)
	}
	after, e := f.Stat()
	current, p := os.Lstat(path)
	if e != nil || p != nil {
		return nil, firstError(e, p)
	}
	if !current.Mode().IsRegular() || !os.SameFile(opened, current) || opened.Size() != int64(len(data)) || after.Size() != opened.Size() || after.ModTime() != opened.ModTime() {
		return nil, contentError("task_content_integrity_conflict", "file changed or was short during read", nil)
	}
	return data, nil
}
func (s *TaskContentStorage) ReadSource(path string, limit int) ([]byte, error) {
	if e := s.check("source-read"); e != nil {
		return nil, e
	}
	return contentReadBounded(path, limit, false)
}
func contentDirectories(root, kind string) []string {
	return []string{filepath.Join(root, MarkerDirectory, "task-content"), taskContentRoot(root), filepath.Join(taskContentRoot(root), kind), filepath.Join(taskContentRoot(root), kind, "sha256")}
}
func (s *TaskContentStorage) directories(root, kind string, create bool) error {
	if e := contentPhysical(filepath.Join(root, MarkerDirectory)); e != nil {
		return e
	}
	for _, dir := range contentDirectories(root, kind) {
		info, e := os.Lstat(dir)
		if errors.Is(e, fs.ErrNotExist) && create {
			if e = s.check("mkdir"); e != nil {
				return e
			}
			if e = os.Mkdir(dir, 0o700); e != nil {
				return e
			}
			if e = platformSyncDirectory(filepath.Dir(dir)); e != nil {
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
func (s *TaskContentStorage) Read(root, kind, digest string) ([]byte, error) {
	if !digestPattern.MatchString(digest) {
		return nil, contentError("task_content_integrity_conflict", "invalid managed digest", nil)
	}
	if e := s.check("managed-read"); e != nil {
		return nil, e
	}
	if e := s.directories(root, kind, false); e != nil {
		return nil, e
	}
	limit := taskManifestLimit
	if kind == "objects" {
		limit = taskDocumentLimit
	}
	if kind == "backups" {
		limit = int(^uint(0)>>1) - 1
	}
	b, e := contentReadBounded(taskContentPath(root, kind, digest), limit, true)
	if e != nil {
		return nil, e
	}
	if digestTaskBytes(b) != digest {
		return nil, contentError("task_content_integrity_conflict", "managed bytes do not match their digest", nil)
	}
	return b, nil
}
func (s *TaskContentStorage) Publish(root, kind string, b []byte) (string, error) {
	digest := digestTaskBytes(b)
	if e := s.directories(root, kind, true); e != nil {
		return "", e
	}
	dest := taskContentPath(root, kind, digest)
	if _, e := os.Lstat(dest); e == nil {
		old, e := s.Read(root, kind, digest)
		if e != nil {
			return "", e
		}
		if string(old) != string(b) {
			return "", contentError("task_content_integrity_conflict", "immutable destination differs", nil)
		}
		return digest, s.Sync(root, kind, digest)
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
	if e = platformSyncDirectory(filepath.Dir(dest)); e != nil {
		return "", e
	}
	return digest, nil
}
func (s *TaskContentStorage) Sync(root, kind, digest string) error {
	if _, e := s.Read(root, kind, digest); e != nil {
		return e
	}
	f, e := os.Open(taskContentPath(root, kind, digest))
	if e != nil {
		return e
	}
	syncErr := f.Sync()
	closeErr := f.Close()
	if e = firstError(syncErr, closeErr, s.check("recovery-sync")); e != nil {
		return e
	}
	for _, d := range contentDirectories(root, kind) {
		if e = platformSyncDirectory(d); e != nil {
			return e
		}
	}
	return platformSyncDirectory(filepath.Join(root, MarkerDirectory))
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
func readContentManifest(s *TaskContentStorage, root, digest string) (canonicaljson.Object, error) {
	b, e := s.Read(root, "manifests", digest)
	if e != nil {
		return nil, e
	}
	v, e := canonicaljson.DecodeStrict(b)
	if e != nil {
		return nil, e
	}
	m := contentFields(v)
	op := contentManifestOperation(contentString(m, "kind"))
	if op == "spec_select" && contentString(m, "action") == "withdraw" {
		op = "spec_withdraw"
	}
	if op == "" {
		return nil, contentError("task_content_integrity_conflict", "unknown manifest kind", nil)
	}
	o, e := decodeTaskContent(b, op, true, true)
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
func contentResolveDocuments(d Dependencies, root string, r WorkItemRegistry, id TaskID, inputs []canonicaljson.Value) ([]canonicaljson.Value, map[string][]byte, error) {
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
			if !contentManifestOwned(r, id, digest) {
				return nil, nil, contentError("task_content_integrity_conflict", "snapshot manifest is not registered to this Task", nil)
			}
			manifest, e := readContentManifest(d.TaskContent, root, digest)
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
			bytes, err = d.TaskContent.Read(root, "objects", contentString(descriptor, "sha256"))
			if err != nil {
				return nil, nil, err
			}
			if contentInt(descriptor, "size_bytes") != len(bytes) || contentString(descriptor, "locator") != taskContentPath(root, "objects", contentString(descriptor, "sha256")) {
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
			descriptor = map[string]canonicaljson.Value{"id": m["id"], "locator": taskContentPath(root, "objects", digestTaskBytes(bytes)), "sha256": digestTaskBytes(bytes), "size_bytes": int64(len(bytes)), "media_type": src["media_type"], "provenance": contentObject(map[string]canonicaljson.Value{"origin_locator": src["locator"], "git": src["git_provenance"]})}
		}
		if err = taskDocumentBytesValid(bytes); err != nil {
			return nil, nil, err
		}
		objects[digestTaskBytes(bytes)] = bytes
		out = append(out, contentObject(descriptor))
	}
	return out, objects, nil
}

func validateTaskContentClosure(s *TaskContentStorage, root string, r WorkItemRegistry, sync bool) error {
	if r.FormatVersion < 3 {
		return nil
	}
	for _, pub := range r.TaskContentPublications {
		manifest, e := readContentManifest(s, root, pub.OutcomeRef.ManifestSHA256)
		if e != nil {
			return e
		}
		m := contentFields(manifest)
		if e := validateExecutionGoalOrigin(Dependencies{TaskContent: s}, root, r, pub.TaskID, manifest); e != nil {
			return e
		}
		if contentString(m, "task_id") != string(pub.TaskID) || contentString(m, "publication_key") != pub.PublicationKey || contentString(m, "source_draft_sha256") != pub.IntentSHA256 {
			return contentError("task_content_integrity_conflict", "publication and manifest binding differ", nil)
		}
		request, e := s.Read(root, "requests", pub.RequestSHA256)
		if e != nil {
			return e
		}
		if e = validateContentPublicationChain(s, root, r, pub, m, request); e != nil {
			return contentError("task_content_integrity_conflict", "preserved publication chain differs", e)
		}
		if pub.BackupSHA256 != nil {
			if _, e = s.Read(root, "backups", *pub.BackupSHA256); e != nil {
				return e
			}
		}
		for _, v := range contentArray(m, "documents") {
			doc := contentFields(v)
			digest := contentString(doc, "sha256")
			b, e := s.Read(root, "objects", digest)
			if e != nil {
				return e
			}
			if e = taskDocumentBytesValid(b); e != nil {
				return e
			}
			if len(b) != contentInt(doc, "size_bytes") || contentString(doc, "locator") != taskContentPath(root, "objects", digest) {
				return contentError("task_content_integrity_conflict", "document descriptor differs from managed object", nil)
			}
			if sync {
				if e = s.Sync(root, "objects", digest); e != nil {
					return e
				}
			}
		}
		if sync {
			for _, pair := range [][2]string{{"manifests", pub.OutcomeRef.ManifestSHA256}, {"requests", pub.RequestSHA256}} {
				if e = s.Sync(root, pair[0], pair[1]); e != nil {
					return e
				}
			}
			if pub.BackupSHA256 != nil {
				if e = s.Sync(root, "backups", *pub.BackupSHA256); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

func validateContentPublicationChain(s *TaskContentStorage, root string, r WorkItemRegistry, pub TaskContentPublication, manifest map[string]canonicaljson.Value, request []byte) error {
	op := pub.Operation
	if op == "task_create" {
		op = "problem_record"
	}
	draft, err := decodeTaskContent(request, op, false, pub.Operation == "task_create")
	if err != nil {
		return err
	}
	canonical, err := canonicaljson.Marshal(draft)
	if err != nil || string(canonical) != string(request) {
		return fmt.Errorf("request is not canonical")
	}
	dm := contentFields(draft)
	if pub.Operation == "task_create" {
		t, _ := findTask(r, pub.TaskID)
		if t == nil {
			return fmt.Errorf("missing Task")
		}
		var upgrade *TaskRegistryUpgrade
		if err := contentDecode(dm["registry_upgrade"], &upgrade); err != nil {
			return err
		}
		intent, err := contentCanonical(map[string]any{"task_id": t.ID, "title": t.Title, "description": t.Description, "parent_epic_id": t.ParentEpicID, "project_id": t.ProjectID, "repo_id": t.RepoID, "registry_upgrade": upgrade})
		if err != nil || digestTaskBytes(intent) != pub.IntentSHA256 {
			return fmt.Errorf("create intent differs")
		}
	} else if pub.RequestSHA256 != pub.IntentSHA256 {
		return fmt.Errorf("request and intent digests differ")
	}
	for k, v := range dm {
		switch k {
		case "kind", "registry_upgrade", "documents":
			continue
		case "expected_previous", "expected_previous_assessment", "expected_previous_selection":
			k = "previous"
		}
		if !contentEqual(v, manifest[k]) {
			return fmt.Errorf("request field %s differs", k)
		}
	}
	if contentString(manifest, "recorded_at_utc") != pub.RecordedAtUTC {
		return fmt.Errorf("publication time differs")
	}
	var previous canonicaljson.Value
	switch pub.OutcomeRef.Kind {
	case "problem":
		if pub.OutcomeRef.Revision == nil || *pub.OutcomeRef.Revision != contentInt(manifest, "revision") {
			return fmt.Errorf("problem revision differs")
		}
		for _, v := range r.TaskProblemRevisions {
			if v.TaskID == pub.TaskID && v.Revision == *pub.OutcomeRef.Revision-1 {
				previous = contentRefValue(TaskRevisionRef{v.Revision, v.ManifestSHA256})
			}
		}
	case "spec":
		if pub.OutcomeRef.SpecID == nil || *pub.OutcomeRef.SpecID != contentString(manifest, "spec_id") || pub.OutcomeRef.Revision == nil || *pub.OutcomeRef.Revision != contentInt(manifest, "revision") {
			return fmt.Errorf("Spec identity differs")
		}
		for _, v := range r.TaskSpecRevisions {
			if v.TaskID == pub.TaskID && v.SpecID == *pub.OutcomeRef.SpecID && v.Revision == *pub.OutcomeRef.Revision-1 {
				previous = contentRefValue(TaskRevisionRef{v.Revision, v.ManifestSHA256})
			}
		}
	case "assessment":
		if pub.OutcomeRef.ID == nil || *pub.OutcomeRef.ID != contentString(manifest, "id") || pub.OutcomeRef.SpecID == nil || *pub.OutcomeRef.SpecID != contentString(manifest, "spec_id") || pub.OutcomeRef.Revision == nil || *pub.OutcomeRef.Revision != contentInt(contentFields(manifest["spec"]), "revision") {
			return fmt.Errorf("assessment identity differs")
		}
		for _, v := range r.TaskSpecAssessments {
			if v.TaskID == pub.TaskID && v.SpecID == *pub.OutcomeRef.SpecID && v.SpecRevision == *pub.OutcomeRef.Revision && v.Ordinal == contentInt(manifest, "ordinal")-1 {
				previous = contentRefValue(TaskDecisionRef{v.ID, v.ManifestSHA256})
			}
		}
	case "selection":
		if pub.OutcomeRef.ID == nil || *pub.OutcomeRef.ID != contentString(manifest, "id") {
			return fmt.Errorf("selection identity differs")
		}
		for _, v := range r.TaskSolutionSelections {
			if v.TaskID == pub.TaskID && v.Ordinal == contentInt(manifest, "ordinal")-1 {
				previous = contentRefValue(TaskDecisionRef{v.ID, v.ManifestSHA256})
			}
		}
	default:
		return fmt.Errorf("unknown publication outcome")
	}
	if !contentEqual(previous, manifest["previous"]) {
		return fmt.Errorf("predecessor differs from registered chain")
	}
	if p := valueRevision(manifest["problem"]); p != nil {
		found := false
		for _, v := range r.TaskProblemRevisions {
			if v.TaskID == pub.TaskID && v.Revision == p.Revision && v.ManifestSHA256 == p.ManifestSHA256 {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("problem reference is not owned by Task")
		}
	}
	if spec := valueRevision(manifest["spec"]); spec != nil {
		found := false
		for _, v := range r.TaskSpecRevisions {
			if v.TaskID == pub.TaskID && v.SpecID == contentString(manifest, "spec_id") && v.Revision == spec.Revision && v.ManifestSHA256 == spec.ManifestSHA256 {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Spec reference is not owned by Task")
		}
	}
	if sol := contentFields(manifest["solution"]); len(sol) > 0 {
		p, sp, a := valueRevision(sol["problem"]), valueRevision(sol["spec"]), valueDecision(sol["assessment"])
		found := false
		for _, v := range r.TaskSpecAssessments {
			if v.TaskID == pub.TaskID && v.SpecID == contentString(sol, "spec_id") && v.SpecRevision == sp.Revision && v.ID == a.ID && v.ManifestSHA256 == a.ManifestSHA256 {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("selection assessment is not registered")
		}
		assessment, err := readContentManifest(s, root, a.ManifestSHA256)
		if err != nil {
			return err
		}
		am := contentFields(assessment)
		if contentString(am, "outcome") != "ready" || !contentTypedEqual(valueRevision(am["spec"]), sp) {
			return fmt.Errorf("selection does not bind a ready assessment")
		}
		spec, err := readContentManifest(s, root, sp.ManifestSHA256)
		if err != nil {
			return err
		}
		sm := contentFields(spec)
		if contentString(sm, "task_id") != string(pub.TaskID) || contentString(sm, "spec_id") != contentString(sol, "spec_id") || !contentTypedEqual(valueRevision(sm["problem"]), p) {
			return fmt.Errorf("selection problem and Spec differ")
		}
	}
	if contentString(contentFields(dm["origin"]), "kind") == "task_create_summary" || contentString(contentFields(dm["origin"]), "kind") == "legacy_summary_import" {
		t, _ := findTask(r, pub.TaskID)
		if t == nil || contentString(dm, "title") != t.Title || contentString(dm, "summary") != t.Description {
			return fmt.Errorf("generated problem differs from the original Task summary")
		}
		body := []byte("# " + t.Title + "\n\n" + t.Description + "\n")
		digest := digestTaskBytes(body)
		expected := contentObject(map[string]canonicaljson.Value{"id": "problem", "locator": taskContentPath(root, "objects", digest), "sha256": digest, "size_bytes": int64(len(body)), "media_type": "text/markdown", "provenance": contentObject(map[string]canonicaljson.Value{"origin_locator": nil, "git": nil})})
		docs := contentArray(manifest, "documents")
		if len(docs) != 1 || !contentEqual(docs[0], expected) {
			return fmt.Errorf("generated problem document differs from the original Task summary")
		}
		if contentString(contentFields(dm["origin"]), "kind") == "legacy_summary_import" {
			original, err := contentCanonical(map[string]any{"task_id": t.ID, "title": t.Title, "description": t.Description})
			if err != nil || contentString(contentFields(dm["origin"]), "legacy_summary_sha256") != digestTaskBytes(original) {
				return fmt.Errorf("legacy summary digest differs")
			}
		}
		return nil
	}
	docs := contentArray(manifest, "documents")
	sourceDocs := contentArray(dm, "documents")
	if len(docs) != len(sourceDocs) {
		return fmt.Errorf("document set differs from request")
	}
	for i, raw := range sourceDocs {
		in := contentFields(raw)
		source := contentFields(in["source"])
		doc := contentFields(docs[i])
		if contentString(in, "id") != contentString(doc, "id") {
			return fmt.Errorf("document ID differs")
		}
		if contentString(source, "kind") == "file" {
			for _, k := range []string{"sha256", "size_bytes", "media_type"} {
				if !contentEqual(source[k], doc[k]) {
					return fmt.Errorf("document metadata differs")
				}
			}
			provenance := contentFields(doc["provenance"])
			if !contentEqual(source["locator"], provenance["origin_locator"]) || !contentEqual(source["git_provenance"], provenance["git"]) {
				return fmt.Errorf("document provenance differs")
			}
		} else {
			digest := contentString(source, "manifest_sha256")
			if !contentManifestOwned(r, pub.TaskID, digest) {
				return fmt.Errorf("snapshot source is not owned by Task")
			}
			origin, err := readContentManifest(s, root, digest)
			if err != nil {
				return err
			}
			found := false
			for _, raw := range contentArray(contentFields(origin), "documents") {
				old := contentFields(raw)
				if contentString(old, "id") == contentString(source, "document_id") {
					old["id"] = doc["id"]
					if !contentEqual(contentObject(old), docs[i]) {
						return fmt.Errorf("snapshot descriptor differs")
					}
					found = true
				}
			}
			if !found {
				return fmt.Errorf("snapshot source document is missing")
			}
		}
	}
	return nil
}

func readRegisteredTaskManifest(s *TaskContentStorage, root string, r WorkItemRegistry, digest string) (canonicaljson.Object, error) {
	m, err := readContentManifest(s, root, digest)
	if err != nil {
		return nil, err
	}
	for _, p := range r.TaskContentPublications {
		if p.OutcomeRef.ManifestSHA256 == digest {
			request, err := s.Read(root, "requests", p.RequestSHA256)
			if err != nil {
				return nil, err
			}
			if err = validateContentPublicationChain(s, root, r, p, contentFields(m), request); err != nil {
				return nil, contentError("task_content_integrity_conflict", "manifest publication chain differs", err)
			}
			return m, nil
		}
	}
	return nil, contentError("task_content_integrity_conflict", "manifest has no registered publication", nil)
}
