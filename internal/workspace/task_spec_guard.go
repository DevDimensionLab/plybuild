package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

type TaskContentGit interface {
	ObserveContentCommit(RepoRecord, string) (string, error)
	ObserveContentDocument(canonicaljson.Value, []byte) error
}

func observeTaskContentGit(d Dependencies, provenance canonicaljson.Value, b []byte) error {
	g, ok := d.WorkGit.(TaskContentGit)
	if !ok {
		return contentError("task_content_observation_unknown", "document Git observer is unavailable", nil)
	}
	return g.ObserveContentDocument(provenance, b)
}
func validateSpecImplementationBasis(d Dependencies, projects ProjectSnapshot, r WorkItemRegistry, t TaskRecord, value canonicaljson.Value, fresh bool) error {
	if fresh && t.WorktreeState == WorkItemRetired {
		return contentError("task_spec_binding_conflict", "retired Task worktree is historical and cannot start new work", nil)
	}
	m := contentFields(value)
	_, repo, e := projectAndRepo(projects, t.ProjectID, t.RepoID)
	if e != nil {
		return e
	}
	epic, _ := findEpic(r, t.ParentEpicID)
	if epic == nil {
		return contentError("task_spec_binding_conflict", "Task Epic is missing", nil)
	}
	parent, _ := findEpicRepo(*epic, t.RepoID)
	if parent != nil {
		historical, ok := historicalEpicBinding(r, *epic, *parent, contentString(m, "parent_oid"), contentString(m, "parent_tree"))
		if !ok {
			return contentError("task_spec_binding_conflict", "implementation basis is not a preserved Epic base", nil)
		}
		if fresh {
			current := currentEpicBase(r, *epic, *parent)
			if current.OID != historical.Worktree.OID || current.Tree != historical.Worktree.Tree {
				return contentError("epic_base_changed", "implementation basis is older than the current Epic base", nil)
			}
			if pendingBaseUpdate(r, t.ParentEpicID, t.RepoID) {
				return contentError("epic_base_pending", "resolve the pending base request first", nil)
			}
		}
		parent = &historical
	}

	if parent == nil || repo.GitCommonDir != t.GitCommonDir || contentString(m, "project_id") != string(t.ProjectID) || contentString(m, "repo_id") != string(t.RepoID) || contentString(m, "git_common_dir") != t.GitCommonDir || contentString(m, "epic_id") != string(t.ParentEpicID) || contentString(m, "parent_worktree_id") != string(parent.Worktree.ID) || contentString(m, "parent_ref") != parent.Worktree.Ref || contentString(m, "parent_oid") != parent.Worktree.OID || contentString(m, "parent_tree") != parent.Worktree.Tree {
		return contentError("task_spec_binding_conflict", "implementation basis differs from registered Task and Epic", nil)
	}
	startOID, startTree := contentString(m, "start_oid"), contentString(m, "start_tree")
	if t.Worktree == nil {
		if startOID != parent.Worktree.OID || startTree != parent.Worktree.Tree {
			return contentError("task_spec_binding_conflict", "start basis before worktree creation must equal parent", nil)
		}
	} else if t.Worktree.ParentOID != parent.Worktree.OID || t.Worktree.ParentTree != parent.Worktree.Tree || t.Worktree.ParentRef != parent.Worktree.Ref || t.Worktree.ParentWorktreeID != parent.Worktree.ID {
		return contentError("task_spec_binding_conflict", "Task worktree parent binding differs", nil)
	}

	g, ok := d.WorkGit.(TaskContentGit)
	if !ok {
		return contentError("task_content_observation_unknown", "commit observer is unavailable", nil)
	}
	tree, e := g.ObserveContentCommit(repo, startOID)
	if e != nil {
		return e
	}
	if tree != startTree {
		return contentError("task_spec_binding_conflict", "start commit tree differs", nil)
	}
	if d.IntegrationGit == nil {
		return contentError("task_content_observation_unknown", "ancestry observer is unavailable", nil)
	}
	ancestor, e := d.IntegrationGit.CheckAncestor(repo, parent.Worktree.OID, startOID)
	if e != nil {
		return e
	}
	if !ancestor {
		return contentError("task_spec_binding_conflict", "start commit is not descended from the recorded parent", nil)
	}
	if !fresh {
		return nil
	}
	p, e := d.WorkGit.ObserveWorktree(parent.Worktree.Locator)
	if e != nil {
		return contentError("task_content_observation_unknown", "cannot observe Epic basis", e)
	}
	if p.Ref != parent.Worktree.Ref || p.OID != parent.Worktree.OID || p.Tree != parent.Worktree.Tree || p.GitCommonDir != t.GitCommonDir || !p.Clean || !p.InventoryMatch {
		return contentError("task_spec_basis_stale", "Epic no longer matches the recorded clean basis", nil)
	}
	if t.Worktree != nil {
		w, e := d.WorkGit.ObserveWorktree(t.Worktree.Locator)
		if e != nil {
			return contentError("task_content_observation_unknown", "cannot observe Task basis", e)
		}
		if w.Ref != t.Worktree.Ref || w.OID != startOID || w.Tree != startTree || w.GitCommonDir != t.GitCommonDir || !w.Clean || !w.InventoryMatch {
			return contentError("task_spec_basis_stale", "Task no longer matches the clean start basis", nil)
		}
	}
	return nil
}
func taskDeliveryDigest(t TaskRecord) (string, error) {
	b, e := contentCanonical(map[string]any{"task_id": t.ID, "project_id": t.ProjectID, "repo_id": t.RepoID, "git_common_dir": t.GitCommonDir, "parent_epic_id": t.ParentEpicID, "worktree": t.Worktree})
	if e != nil {
		return "", e
	}
	return digestTaskBytes(b), nil
}
func validateTaskSpecBasisRegistry(r WorkItemRegistry, b TaskSpecBasis) error {
	t, _ := findTask(r, b.TaskID)
	if t == nil || t.Worktree == nil || t.Worktree.ID != b.TaskWorktreeID || (t.WorktreeState != WorkItemReady && t.WorktreeState != WorkItemRetired) || b.Dependencies == nil || len(b.Dependencies) != 0 {
		return contentError("task_spec_binding_conflict", "basis Task worktree binding is missing or differs", nil)
	}
	d, e := taskDeliveryDigest(*t)
	if e != nil || d != b.DeliveryBindingSHA256 {
		return contentError("task_spec_binding_conflict", "delivery binding digest differs", e)
	}
	p, s, a, c := false, false, false, false
	for _, v := range r.TaskProblemRevisions {
		if v.TaskID == b.TaskID && v.Revision == b.Problem.Revision && v.ManifestSHA256 == b.Problem.ManifestSHA256 {
			p = true
		}
	}
	for _, v := range r.TaskSpecRevisions {
		if v.TaskID == b.TaskID && v.SpecID == b.SpecID && v.Revision == b.Spec.Revision && v.ManifestSHA256 == b.Spec.ManifestSHA256 {
			s = true
		}
	}
	for _, v := range r.TaskSpecAssessments {
		if v.TaskID == b.TaskID && v.SpecID == b.SpecID && v.SpecRevision == b.Spec.Revision && v.ID == b.Assessment.ID && v.ManifestSHA256 == b.Assessment.ManifestSHA256 {
			a = true
		}
	}
	for _, v := range r.TaskSolutionSelections {
		if v.TaskID == b.TaskID && v.ID == b.Selection.ID && v.ManifestSHA256 == b.Selection.ManifestSHA256 {
			c = true
		}
	}
	if !p || !s || !a || !c {
		return contentError("task_spec_binding_conflict", "basis contains unregistered references", nil)
	}
	return nil
}

type TaskSpecRequiredInput struct {
	ID         string  `yaml:"id" json:"id"`
	Role       string  `yaml:"role" json:"role"`
	Locator    string  `yaml:"locator" json:"locator"`
	SHA256     string  `yaml:"sha256" json:"sha256"`
	SizeBytes  int64   `yaml:"size_bytes" json:"size_bytes"`
	MediaType  string  `yaml:"media_type" json:"media_type"`
	GitBinding *string `yaml:"git_binding" json:"git_binding"`
}
type TaskSpecEvaluation struct {
	Basis                                                 *TaskSpecBasis
	RequiredInputs                                        []TaskSpecRequiredInput
	Problem, Spec, Assessment, Selection                  canonicaljson.Object
	SelectionFreshness, ContentIntegrity, TargetFreshness string
	Reasons                                               []string
}

func historicalTaskSpec(d Dependencies, root string, r WorkItemRegistry, b TaskSpecBasis) (TaskSpecEvaluation, error) {
	if e := validateTaskSpecBasisRegistry(r, b); e != nil {
		return TaskSpecEvaluation{}, e
	}
	return loadTaskSpecContent(d, root, r, b)
}
func loadTaskSpecContent(d Dependencies, root string, r WorkItemRegistry, b TaskSpecBasis) (TaskSpecEvaluation, error) {
	out := TaskSpecEvaluation{Basis: &b, RequiredInputs: []TaskSpecRequiredInput{}, SelectionFreshness: "unknown", ContentIntegrity: "unknown", TargetFreshness: "unknown", Reasons: []string{}}
	refs := []struct {
		name, digest string
		target       *canonicaljson.Object
	}{{"problem", b.Problem.ManifestSHA256, &out.Problem}, {"spec", b.Spec.ManifestSHA256, &out.Spec}, {"assessment", b.Assessment.ManifestSHA256, &out.Assessment}, {"selection", b.Selection.ManifestSHA256, &out.Selection}}
	for _, ref := range refs {
		v, e := readRegisteredTaskManifest(d.TaskContent, root, r, ref.digest)
		if e != nil {
			return out, e
		}
		*ref.target = v
		if contentString(contentFields(v), "task_id") != string(b.TaskID) {
			return out, contentError("task_spec_binding_conflict", "manifest Task owner differs", nil)
		}
		bytes, e := canonicaljson.Marshal(v)
		if e != nil {
			return out, e
		}
		out.RequiredInputs = append(out.RequiredInputs, TaskSpecRequiredInput{ID: "task-basis/" + ref.name, Role: "other", Locator: taskContentPath(root, "manifests", ref.digest), SHA256: ref.digest, SizeBytes: int64(len(bytes)), MediaType: "application/json"})
	}
	pm, sm, am, cm := contentFields(out.Problem), contentFields(out.Spec), contentFields(out.Assessment), contentFields(out.Selection)
	if contentInt(pm, "revision") != b.Problem.Revision || contentInt(sm, "revision") != b.Spec.Revision || contentString(sm, "spec_id") != b.SpecID || !contentTypedEqual(valueRevision(sm["problem"]), &b.Problem) || !contentTypedEqual(valueRevision(am["spec"]), &b.Spec) || contentString(am, "spec_id") != b.SpecID || contentString(am, "id") != b.Assessment.ID || contentString(am, "outcome") != "ready" || contentString(cm, "id") != b.Selection.ID || contentString(cm, "action") != "select" {
		return out, contentError("task_spec_binding_conflict", "historical manifest chain differs from basis", nil)
	}
	solution := contentFields(cm["solution"])
	if contentString(solution, "spec_id") != b.SpecID || !contentTypedEqual(valueRevision(solution["spec"]), &b.Spec) || !contentTypedEqual(valueRevision(solution["problem"]), &b.Problem) || !contentTypedEqual(valueDecision(solution["assessment"]), &b.Assessment) {
		return out, contentError("task_spec_binding_conflict", "historical selection differs from basis", nil)
	}
	documentRoles := map[string]string{}
	parts := contentFields(sm["parts"])
	for _, name := range []string{"abstract", "functional", "technical"} {
		for _, v := range contentArray(contentFields(parts[name]), "documents") {
			documentRoles[contentString(contentFields(v), "document_id")] = "spec"
		}
	}
	for _, v := range contentArray(sm, "supporting") {
		p := contentFields(v)
		if contentString(p, "usage") == "normative" {
			documentRoles[contentString(p, "document_id")] = "design"
		}
	}
	objects := map[string]TaskSpecRequiredInput{}
	total := int64(0)
	priority := map[string]int{"other": 0, "design": 1, "spec": 2}
	for index, manifest := range []canonicaljson.Object{out.Problem, out.Spec, out.Assessment} {
		for _, v := range contentArray(contentFields(manifest), "documents") {
			doc := contentFields(v)
			digest := contentString(doc, "sha256")
			role := "other"
			if index == 1 && documentRoles[contentString(doc, "id")] != "" {
				role = documentRoles[contentString(doc, "id")]
			}
			bytes, e := d.TaskContent.Read(root, "objects", digest)
			if e != nil {
				return out, e
			}
			if len(bytes) != contentInt(doc, "size_bytes") || contentString(doc, "locator") != taskContentPath(root, "objects", digest) {
				return out, contentError("task_content_integrity_conflict", "closure document descriptor differs", nil)
			}
			if old, ok := objects[digest]; ok {
				if priority[role] > priority[old.Role] {
					old.Role = role
					objects[digest] = old
				}
				continue
			}
			total += int64(len(bytes))
			objects[digest] = TaskSpecRequiredInput{ID: "task-doc/" + strings.TrimPrefix(digest, "sha256:"), Role: role, Locator: taskContentPath(root, "objects", digest), SHA256: digest, SizeBytes: int64(len(bytes)), MediaType: contentString(doc, "media_type")}
		}
	}
	if len(objects) > 96 || total > 64<<20 {
		return out, contentError("task_spec_input_coverage", "handoff closure exceeds 96 documents or 64 MiB", nil)
	}
	for _, v := range objects {
		out.RequiredInputs = append(out.RequiredInputs, v)
	}
	sort.Slice(out.RequiredInputs, func(i, j int) bool { return out.RequiredInputs[i].ID < out.RequiredInputs[j].ID })
	out.ContentIntegrity = "valid"
	return out, nil
}
func currentTaskSpec(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, t TaskRecord, fresh bool) (TaskSpecEvaluation, error) {
	out := TaskSpecEvaluation{RequiredInputs: []TaskSpecRequiredInput{}, SelectionFreshness: "none", ContentIntegrity: "not_recorded", TargetFreshness: "unknown", Reasons: []string{}}
	state := taskContentState(r, t.ID)
	if state.SelectionEvent == nil {
		return out, nil
	}
	selected, e := readContentManifest(d.TaskContent, root, state.SelectionEvent.ManifestSHA256)
	if e != nil {
		return out, e
	}
	out.Selection = selected
	m := contentFields(selected)
	if contentString(m, "action") == "withdraw" {
		out.SelectionFreshness = "withdrawn"
		return out, nil
	}
	sol := contentFields(m["solution"])
	b := TaskSpecBasis{TaskID: t.ID, Problem: *valueRevision(sol["problem"]), SpecID: contentString(sol, "spec_id"), Spec: *valueRevision(sol["spec"]), Assessment: *valueDecision(sol["assessment"]), Selection: *state.SelectionEvent, Dependencies: []string{}}
	if t.Worktree == nil {
		out, e = loadTaskSpecContent(d, root, r, b)
		out.Basis = nil
		out.TargetFreshness = "stale"
	} else {
		digest, err := taskDeliveryDigest(t)
		if err != nil {
			return out, err
		}
		b.TaskWorktreeID = t.Worktree.ID
		b.DeliveryBindingSHA256 = digest
		out, e = historicalTaskSpec(d, root, r, b)
		out.TargetFreshness = "fresh"
	}
	if e != nil {
		return out, e
	}
	out.SelectionFreshness = "current"

	last := lastTaskAssessment(r, t.ID, b.SpecID, b.Spec.Revision)
	if !contentTypedEqual(state.ProblemHead, &b.Problem) || last == nil || last.ID != b.Assessment.ID || last.ManifestSHA256 != b.Assessment.ManifestSHA256 {
		out.SelectionFreshness = "stale"
		out.Reasons = append(out.Reasons, "task_spec_selection_stale")
	}
	if ep, _ := findEpic(r, t.ParentEpicID); ep != nil {
		if parent, _ := findEpicRepo(*ep, t.RepoID); parent != nil {
			base := currentEpicBase(r, *ep, *parent)
			basis := contentFields(contentFields(out.Spec)["implementation_basis"])
			parentOID, parentTree := contentString(basis, "parent_oid"), contentString(basis, "parent_tree")
			if parentOID != base.OID || parentTree != base.Tree {
				out.TargetFreshness = "stale"
				out.Reasons = append(out.Reasons, "epic_base_changed")
			}
		}
	}
	if fresh {
		if e = validateSpecImplementationBasis(d, projects, r, t, contentFields(out.Spec)["implementation_basis"], true); e != nil {
			out.TargetFreshness = "unknown"
			if r.FormatVersion < 4 {
				out.SelectionFreshness = "unknown"
			}
			reason := "task_content_observation_unknown"
			var contentErr *TaskContentError
			if errors.As(e, &contentErr) && (contentErr.Code == "task_spec_basis_stale" || contentErr.Code == "task_spec_binding_conflict" || contentErr.Code == "epic_base_changed") {
				out.TargetFreshness = "stale"
				if r.FormatVersion < 4 {
					out.SelectionFreshness = "stale"
				}
				reason = contentErr.Code
			}
			out.Reasons = append(out.Reasons, reason)
		}
	}
	return out, nil
}

// Read-only handles take the same advisory locks as ordinary metadata writers.
// No file or directory is created, truncated or chmodded by this guard.
func acquireExistingContentLock(path string) (*filelock.Lock, error) {
	if e := contentPhysical(path); e != nil {
		return nil, contentError("task_spec_lock_unavailable", "existing lock is unavailable: "+path, e)
	}
	lock, e := filelock.AcquireExisting(path)
	var lockErr *filelock.Error
	if errors.As(e, &lockErr) {
		return nil, contentLockError(lockErr)
	}
	return lock, e
}

func contentLockError(e *filelock.Error) error {
	if e.Op == "lock" {
		return contentError("task_spec_lock_unavailable", "cannot acquire existing lock", e.Err)
	}
	return contentError("task_spec_lock_unavailable", "existing lock identity differs", e.Err)
}

func releaseContentLock(lock *filelock.Lock) error {
	if e := lock.Release(); e != nil {
		return contentReleaseError(e)
	}
	return nil
}

func contentReleaseError(e *filelock.ReleaseError) error { return firstError(e.Unlock, e.Close) }

type TaskSpecSession struct {
	Dependencies Dependencies
	Root         string
	Registry     WorkItemRegistry
	Projects     ProjectSnapshot
}

func WithTaskSpecSnapshot(d Dependencies, root string, operation func(*TaskSpecSession) error) (err error) {
	return withTaskSpecSnapshot(d, root, false, operation)
}

// WithTaskSpecRegistrationSnapshot keeps the same locks and full structural
// registry checks. ValidateTarget validates the resolved Task's content before
// any handoff effect; historical reads independently validate their exact basis.
func WithTaskSpecRegistrationSnapshot(d Dependencies, root string, operation func(*TaskSpecSession) error) error {
	return withTaskSpecSnapshot(d, root, true, operation)
}

func withTaskSpecSnapshot(d Dependencies, root string, registrations bool, operation func(*TaskSpecSession) error) (err error) {
	project, e := acquireExistingContentLock(filepath.Join(root, MarkerDirectory, "projects.lock"))
	if e != nil {
		return e
	}
	defer func() {
		if e := releaseContentLock(project); err == nil {
			err = e
		}
	}()
	lockPath := filepath.Join(root, MarkerDirectory, workItemsLock)
	_, lockErr := os.Lstat(lockPath)
	_, storeErr := os.Lstat(workItemsPath(root))
	if errors.Is(lockErr, fs.ErrNotExist) && errors.Is(storeErr, fs.ErrNotExist) {
	} else {
		work, e := acquireExistingContentLock(lockPath)
		if e != nil {
			return e
		}
		defer func() {
			if e := releaseContentLock(work); err == nil {
				err = e
			}
		}()
	}
	projects, repos, e := d.Projects.Snapshot(root)
	if e != nil {
		return e
	}
	var registry WorkItemRegistry
	if registrations {
		registry, e = d.WorkItems.SnapshotRegistrations(root)
	} else {
		registry, e = d.WorkItems.Snapshot(root)
	}
	if e != nil {
		return e
	}
	return operation(&TaskSpecSession{d, root, registry, ProjectSnapshot{Projects: projects, Repos: repos}})
}

func ValidateTaskSpecBasisValue(v canonicaljson.Value) (*TaskSpecBasis, error) {
	if v == nil {
		return nil, nil
	}
	rule := contentExact(map[string]contentRule{"task_id": contentSlug, "problem": contentRevisionRule(), "spec_id": contentSlug, "spec": contentRevisionRule(), "assessment": contentDecisionRule("asm_"), "selection": contentDecisionRule("sel_"), "task_worktree_id": func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok || !worktreeIDPattern.MatchString(s) {
			return fmt.Errorf("invalid Task worktree ID")
		}
		return nil
	}, "delivery_binding_sha256": contentDigest, "dependencies": contentList(0, contentKey, nil, false)})
	if e := rule(v); e != nil {
		return nil, e
	}
	var b TaskSpecBasis
	if e := contentDecode(v, &b); e != nil {
		return nil, e
	}
	return &b, nil
}
