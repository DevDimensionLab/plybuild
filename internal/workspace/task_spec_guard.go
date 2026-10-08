package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
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
			if contentString(basis, "parent_oid") != base.OID || contentString(basis, "parent_tree") != base.Tree {
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
func acquireExistingContentLock(path string) (*os.File, error) {
	if e := contentPhysical(path); e != nil {
		return nil, contentError("task_spec_lock_unavailable", "existing lock is unavailable: "+path, e)
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	opened, e := f.Stat()
	current, p := os.Lstat(path)
	if e != nil || p != nil || !opened.Mode().IsRegular() || !current.Mode().IsRegular() || !os.SameFile(opened, current) {
		_ = f.Close()
		return nil, contentError("task_spec_lock_unavailable", "existing lock identity differs", firstError(e, p))
	}
	if e = platformLock(f); e != nil {
		_ = f.Close()
		return nil, contentError("task_spec_lock_unavailable", "cannot acquire existing lock", e)
	}
	return f, nil
}
func releaseContentLock(f *os.File) error { return firstError(platformUnlock(f), f.Close()) }

type TaskSpecSession struct {
	Dependencies Dependencies
	Root         string
	Registry     WorkItemRegistry
	Projects     ProjectSnapshot
}

// TaskSpecReadPaths includes Git sources outside the Task or workspace roots.
// The caller already holds the Project and work-item snapshot locks.
func (s *TaskSpecSession) TaskSpecReadPaths(b *TaskSpecBasis) ([]string, error) {
	if b == nil {
		return nil, nil
	}
	t, _ := findTask(s.Registry, b.TaskID)
	if t != nil {
		epic, _ := findEpic(s.Registry, t.ParentEpicID)
		if epic != nil {
			parent, _ := findEpicRepo(*epic, t.RepoID)
			if parent != nil {
				return []string{parent.Worktree.Locator}, nil
			}
		}
	}
	return nil, contentError("task_spec_binding_conflict", "Task Epic read binding is unavailable", nil)
}

func WithTaskSpecSnapshot(d Dependencies, root string, operation func(*TaskSpecSession) error) (err error) {
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
	registry, e := d.WorkItems.Snapshot(root)
	if e != nil {
		return e
	}
	return operation(&TaskSpecSession{d, root, registry, ProjectSnapshot{Projects: projects, Repos: repos}})
}
func (s *TaskSpecSession) ValidateTarget(locator, ref, common string, basis *TaskSpecBasis) (TaskSpecEvaluation, error) {
	if err := validateCloseoutTargetAvailable(s.Root, locator, ref, common); err != nil {
		return TaskSpecEvaluation{}, err
	}
	physical, e := filepath.EvalSymlinks(locator)
	if e != nil {
		return TaskSpecEvaluation{}, e
	}
	var owner *TaskRecord
	for _, t := range s.Registry.Tasks {
		if t.Worktree == nil || t.WorktreeState == WorkItemRetired {
			continue
		}
		p, e := filepath.EvalSymlinks(t.Worktree.Locator)
		if e != nil {
			return TaskSpecEvaluation{}, contentError("task_content_observation_unknown", "cannot resolve registered Task target", e)
		}
		if p == physical || t.Worktree.Ref == ref && t.GitCommonDir == common {
			if owner != nil {
				return TaskSpecEvaluation{}, contentError("task_spec_binding_conflict", "ambiguous target ownership", nil)
			}
			copy := t
			owner = &copy
		}
	}
	if owner == nil {
		if basis != nil {
			return TaskSpecEvaluation{}, contentError("task_spec_binding_conflict", "target has no registered Task owner", nil)
		}
		return TaskSpecEvaluation{}, nil
	}
	if taskRequiresSpec(s.Registry, owner.ID) && basis == nil {
		return TaskSpecEvaluation{}, contentError("task_spec_required", "target requires an exact selected Task Spec basis", nil)
	}
	if basis == nil {
		return TaskSpecEvaluation{}, nil
	}
	if owner.ID != basis.TaskID || owner.Worktree.Locator != physical || owner.Worktree.Ref != ref || owner.GitCommonDir != common {
		return TaskSpecEvaluation{}, contentError("task_spec_binding_conflict", "target differs from the registered delivery binding", nil)
	}
	out, e := currentTaskSpec(s.Dependencies, s.Root, s.Registry, s.Projects, *owner, true)
	if e != nil {
		return out, e
	}
	if out.SelectionFreshness == "unknown" || out.TargetFreshness == "unknown" && out.SelectionFreshness == "current" {
		return out, contentError("task_content_observation_unknown", "selected Task basis could not be observed", nil)
	}
	if out.SelectionFreshness == "withdrawn" {
		return out, contentError("task_spec_selection_withdrawn", "Task selection has been withdrawn", nil)
	}
	if out.SelectionFreshness == "none" {
		return out, contentError("task_spec_selection_missing", "Task has no selected solution", nil)
	}
	if out.ContentIntegrity != "valid" {
		return out, contentError(contentIntegrityReason(out.ContentIntegrity), "selected Task content could not be verified", nil)
	}
	if out.TargetFreshness == "stale" {
		return out, contentError("task_spec_basis_stale", "Task target differs from the selected implementation basis", nil)
	}
	if !contentTypedEqual(out.Basis, basis) || out.SelectionFreshness != "current" || out.ContentIntegrity != "valid" || out.TargetFreshness != "fresh" {
		return out, contentError("task_spec_selection_stale", "selected Task basis is not current, intact and ready", nil)
	}
	return out, nil
}
func ReadHistoricalTaskSpec(d Dependencies, root string, b TaskSpecBasis) (TaskSpecEvaluation, error) {
	r, e := d.WorkItems.Snapshot(root)
	if e != nil {
		return TaskSpecEvaluation{}, e
	}
	return historicalTaskSpec(d, root, r, b)
}
func (s *TaskSpecSession) StartObservation(expected *TaskSpecBasis, out TaskSpecEvaluation) canonicaljson.Value {
	if expected == nil {
		return nil
	}
	eb, _ := contentCanonical(expected)
	var ob canonicaljson.Value
	if out.Basis != nil {
		bytes, _ := contentCanonical(out.Basis)
		ob = digestTaskBytes(bytes)
	}
	st := taskContentState(s.Registry, expected.TaskID)
	p, _ := contentValue(st.ProblemHead)
	sel, _ := contentValue(st.SelectionEvent)
	var registry canonicaljson.Value
	if s.Registry.RawSHA256 != "" {
		registry = s.Registry.RawSHA256
	}
	return contentObject(map[string]canonicaljson.Value{"expected_basis_sha256": digestTaskBytes(eb), "observed_basis_sha256": ob, "registry_sha256": registry, "current_problem": p, "current_selection": sel, "selection_freshness": out.SelectionFreshness, "content_integrity": out.ContentIntegrity, "target_freshness": out.TargetFreshness, "validated_at_utc": s.Dependencies.WorkClock.Now().UTC().Format(time.RFC3339Nano), "reasons": sortedContentStrings(out.Reasons)})
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

type TaskSpecRelevance struct {
	Basis       *TaskSpecBasis   `yaml:"basis" json:"basis"`
	BasisSHA256 *string          `yaml:"basis_sha256" json:"basis_sha256"`
	Selection   *TaskDecisionRef `yaml:"selection" json:"selection"`
	Problem     *TaskRevisionRef `yaml:"problem" json:"problem"`
	Assessment  *TaskDecisionRef `yaml:"assessment" json:"assessment"`
	Relevance   string           `yaml:"relevance" json:"relevance"`
	Reasons     []string         `yaml:"reasons" json:"reasons"`
}

func taskResultSpecLink(r WorkItemRegistry, id TaskResultID) *TaskResultSpecBinding {
	for _, v := range r.TaskResultSpecBindings {
		if v.TaskResultID == id {
			copy := v
			return &copy
		}
	}
	return nil
}
func taskResultSpecRelevance(d Dependencies, root string, projects ProjectSnapshot, r WorkItemRegistry, t TaskRecord, result TaskResultID) *TaskSpecRelevance {
	if !taskRequiresSpec(r, t.ID) {
		return nil
	}
	out := &TaskSpecRelevance{Relevance: "stale", Reasons: []string{}}
	state := taskContentState(r, t.ID)
	out.Problem = state.ProblemHead
	out.Selection = state.SelectionEvent
	link := taskResultSpecLink(r, result)
	if link == nil {
		out.Reasons = append(out.Reasons, "task_spec_result_basis_missing")
		return out
	}
	out.Basis = &link.Basis
	bytes, e := contentCanonical(link.Basis)
	if e != nil {
		out.Relevance = "unknown"
		out.Reasons = append(out.Reasons, "task_content_observation_unknown")
		return out
	}
	hash := digestTaskBytes(bytes)
	out.BasisSHA256 = &hash
	if last := lastTaskAssessment(r, t.ID, link.Basis.SpecID, link.Basis.Spec.Revision); last != nil {
		out.Assessment = &TaskDecisionRef{last.ID, last.ManifestSHA256}
	}
	if _, e = historicalTaskSpec(d, root, r, link.Basis); e != nil {
		out.Relevance = "unknown"
		out.Reasons = append(out.Reasons, contentIntegrityReason(contentIntegrity(e)))
		return out
	}
	cur, e := currentTaskSpec(d, root, r, projects, t, false)
	if e != nil {
		out.Relevance = "unknown"
		out.Reasons = append(out.Reasons, "task_content_observation_unknown")
		return out
	}
	if cur.SelectionFreshness == "current" && cur.ContentIntegrity == "valid" && cur.TargetFreshness != "stale" && contentTypedEqual(cur.Basis, &link.Basis) {
		out.Relevance = "current"
		if r.FormatVersion == 4 {
			ep, _ := findEpic(r, t.ParentEpicID)
			parent, _ := findEpicRepo(*ep, t.RepoID)
			base := currentEpicBase(r, *ep, *parent)
			target := QueueTarget{ProjectID: t.ProjectID, RepoID: t.RepoID, EpicID: ep.ID, GitCommonDir: parent.GitCommonDir, ParentWorktreeID: parent.Worktree.ID, ParentLocator: parent.Worktree.Locator, ParentRef: parent.Worktree.Ref}
			observed, err := queueParentObservation(d, target)
			if err != nil {
				out.Relevance = "unknown"
				out.Reasons = append(out.Reasons, "task_content_observation_unknown")
				var contentErr *TaskContentError
				if errors.As(err, &contentErr) && contentErr.Code == "epic_base_stale" {
					out.Relevance = "stale"
					out.Reasons = []string{"epic_base_changed"}
				}
			} else if observed.OID != base.OID || observed.Tree != base.Tree {
				out.Relevance = "stale"
				out.Reasons = append(out.Reasons, "epic_base_changed")
			}
		}
	} else {
		out.Reasons = append(out.Reasons, "task_spec_result_not_current")
		out.Reasons = append(out.Reasons, cur.Reasons...)
		out.Reasons = sortedReasons(out.Reasons)
	}
	return out
}
func validateStoredTaskSpecGuard(r WorkItemRegistry, p WorkspaceTaskIntegrationPlan) error {
	if p.SchemaVersion == 1 {
		if p.Kind != "WorkspaceTaskIntegrationPlan@1" || p.TaskSpecGuard != nil || p.DeliveryAuthorization != nil {
			return fmt.Errorf("invalid legacy integration plan")
		}
		return nil
	}
	validVersion := p.SchemaVersion == 2 && p.Kind == "WorkspaceTaskIntegrationPlan@2" && p.DeliveryAuthorization == nil || p.SchemaVersion == 3 && p.Kind == "WorkspaceTaskIntegrationPlan@3" && p.DeliveryAuthorization != nil
	if !validVersion || p.TaskSpecGuard == nil {
		return fmt.Errorf("unsupported integration plan version")
	}
	g := p.TaskSpecGuard
	if g.Basis == nil || g.BasisSHA256 == nil || g.Relevance != "current" || g.Reasons == nil || len(g.Reasons) != 0 {
		return fmt.Errorf("persisted integration authority requires a current exact Task Spec basis")
	}
	if e := validateTaskSpecBasisRegistry(r, *g.Basis); e != nil {
		return e
	}
	bytes, e := contentCanonical(g.Basis)
	if e != nil || *g.BasisSHA256 != digestTaskBytes(bytes) || !contentTypedEqual(g.Selection, &g.Basis.Selection) || !contentTypedEqual(g.Problem, &g.Basis.Problem) || !contentTypedEqual(g.Assessment, &g.Basis.Assessment) {
		return fmt.Errorf("persisted Task Spec guard is inconsistent")
	}
	link := taskResultSpecLink(r, p.TaskResult.ID)
	if link == nil || !contentTypedEqual(link.Basis, *g.Basis) {
		return fmt.Errorf("plan basis differs from Task result link")
	}
	return nil
}

func ValidateTaskRequirementEvidence(raw []byte, b TaskSpecBasis, spec canonicaljson.Object, oid, tree string, verifiers []TaskVerifierResultRecord, artifacts []TaskArtifactRecord) error {
	if len(raw) > taskManifestLimit {
		return fmt.Errorf("task-requirements exceeds 256 KiB")
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return e
	}
	schema := map[string]contentRule{"kind": contentEnum("WorkspaceTaskRequirementEvidence@1"), "schema_version": contentInteger(1, 1), "format": contentEnum("json"), "format_version": contentInteger(1, 1), "canonicalization": contentEnum("RFC8785"), "task_id": contentSlug, "spec_id": contentSlug, "spec": contentRevisionRule(), "result_oid": contentOID, "result_tree": contentOID, "requirements": contentList(128, contentExact(map[string]contentRule{"id": contentKey, "outcome": contentEnum("passed", "failed", "not_run", "unknown"), "verifier_ids": contentStringList(128, contentKey), "artifact_ids": contentStringList(128, contentKey), "reason": contentText(2000)}), contentIDKey, false)}
	if e = contentExact(schema)(v); e != nil {
		return e
	}
	m := contentFields(v)
	if contentString(m, "task_id") != string(b.TaskID) || contentString(m, "spec_id") != b.SpecID || !contentTypedEqual(valueRevision(m["spec"]), &b.Spec) || contentString(m, "result_oid") != oid || contentString(m, "result_tree") != tree {
		return fmt.Errorf("requirement evidence result or Spec binding differs")
	}
	vs := map[string]bool{}
	for _, r := range verifiers {
		vs[r.VerifierID] = true
	}
	as := map[string]bool{}
	for _, a := range artifacts {
		if a.ArtifactID != "task-requirements" {
			as[a.ArtifactID] = true
		}
	}
	reqs := contentArray(contentFields(spec), "requirements")
	coverage := contentArray(m, "requirements")
	if len(reqs) != len(coverage) {
		return fmt.Errorf("requirement evidence ID set differs")
	}
	for i, v := range coverage {
		entry := contentFields(v)
		if contentIDKey(reqs[i]) != contentString(entry, "id") {
			return fmt.Errorf("requirement evidence ID set differs")
		}
		if contentString(entry, "outcome") != "passed" {
			return fmt.Errorf("requirement %s has outcome %s", contentString(entry, "id"), contentString(entry, "outcome"))
		}
		for _, id := range contentArray(entry, "verifier_ids") {
			if !vs[id.(string)] {
				return fmt.Errorf("requirement verifier is not present in terminal")
			}
		}
		for _, id := range contentArray(entry, "artifact_ids") {
			if !as[id.(string)] {
				return fmt.Errorf("requirement artifact is absent or self-referential")
			}
		}
	}
	return nil
}

// ValidateTaskRunTechnicalAssessment keeps the TaskResult gate and debt schema
// authoritative for semantic reports without creating a TaskResult draft.
func ValidateTaskRunTechnicalAssessment(raw []byte) error {
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return e
	}
	f, e := exactTaskObject(v, "technical_assessment", "gate", "required_verifier_ids", "accepted_debt")
	if e != nil {
		return e
	}
	if !setString("passed", "good_enough_with_known_debt", "failed", "unknown")[mustTaskString(f, "gate")] {
		return fmt.Errorf("invalid technical gate")
	}
	if _, e = taskStringArray(f["required_verifier_ids"], true); e != nil {
		return e
	}
	_, e = parseAcceptedDebt(f["accepted_debt"])
	return e
}
