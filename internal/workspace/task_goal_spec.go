package workspace

import (
	"fmt"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// TaskGoalRef identifies the immutable goal chosen by a planner or execution.
// It deliberately contains no Git base or claim of human solution selection.
type TaskGoalRef struct {
	SpecID string          `yaml:"spec_id" json:"spec_id"`
	Spec   TaskRevisionRef `yaml:"spec" json:"spec"`
}

type TaskExecutorAssignment struct {
	Provider *string `yaml:"provider" json:"provider"`
	Model    *string `yaml:"model" json:"model"`
	Effort   *string `yaml:"effort" json:"effort"`
}

type TaskGoalRequirement struct {
	ID         string `json:"id"`
	Acceptance string `json:"acceptance"`
}

type TaskGoalContract struct {
	TaskID         TaskID                  `json:"task_id"`
	Goal           TaskGoalRef             `json:"goal"`
	Problem        TaskRevisionRef         `json:"problem"`
	Title          string                  `json:"title"`
	Objective      string                  `json:"objective"`
	Requirements   []TaskGoalRequirement   `json:"requirements"`
	Constraints    []string                `json:"constraints"`
	RequiredInputs []TaskSpecRequiredInput `json:"required_inputs"`
	Executor       TaskExecutorAssignment  `json:"executor"`
}

func taskGoalRefRule() contentRule {
	return contentExact(map[string]contentRule{"spec_id": contentSlug, "spec": contentRevisionRule()})
}

func isTaskGoalSpec(v canonicaljson.Value) bool {
	m := contentFields(v)
	return contentInt(m, "schema_version") == 2 && contentString(m, "contract_kind") == "goal"
}

func taskGoalSpecSchema(published bool) map[string]contentRule {
	s := taskContentSchema("spec_record", published, false)
	s["schema_version"] = contentInteger(2, 2)
	s["kind"] = contentEnum("WorkspaceTaskSpecDraft@2")
	if published {
		s["kind"] = contentEnum("WorkspaceTaskSpecRevision@2")
	}
	for _, k := range []string{"parts", "supporting", "removed_requirement_ids", "phases", "implementation_basis", "dependencies"} {
		delete(s, k)
	}
	s["contract_kind"] = contentEnum("goal")
	s["objective"] = contentText(8000)
	s["requirements"] = contentList(128, contentExact(map[string]contentRule{"id": contentKey, "acceptance": contentText(2000)}), contentIDKey, false)
	s["design"] = contentPartRefs()
	s["constraints"] = contentStringList(128, contentText(2000))
	s["executor"] = contentExact(map[string]contentRule{"provider": contentNullable(contentEnum("codex", "claude")), "model": contentNullable(contentText(256)), "effort": contentNullable(contentText(128))})
	doc := contentDocumentInputRule()
	if published {
		doc = contentDescriptorRule()
	}
	s["documents"] = contentList(31, doc, contentIDKey, false)
	return s
}

func taskExecutionSpecSchema(published bool) map[string]contentRule {
	s := taskContentSchema("spec_record", published, false)
	s["schema_version"] = contentInteger(2, 2)
	s["kind"] = contentEnum("WorkspaceTaskSpecDraft@2")
	if published {
		s["kind"] = contentEnum("WorkspaceTaskSpecRevision@2")
	}
	s["contract_kind"] = contentEnum("execution")
	s["goal_origin"] = taskGoalRefRule()
	s["execution_request_sha256"] = contentDigest
	s["acceptance_path"] = contentPath
	s["execution_request"] = goalExecutionRequestRule
	return s
}

func validateTaskGoalSemantics(m map[string]canonicaljson.Value) error {
	if len(contentArray(m, "requirements")) == 0 || len(contentArray(m, "design")) == 0 {
		return fmt.Errorf("a goal needs at least one observable requirement and design reference")
	}
	docs := map[string]bool{}
	for _, v := range contentArray(m, "documents") {
		docs[contentString(contentFields(v), "id")] = false
	}
	for _, v := range contentArray(m, "design") {
		id := contentString(contentFields(v), "document_id")
		if _, ok := docs[id]; !ok {
			return fmt.Errorf("goal design document %s is missing", id)
		}
		docs[id] = true
	}
	for id, used := range docs {
		if !used {
			return fmt.Errorf("unused goal document %s", id)
		}
	}
	return nil
}

func validateRegisteredGoalRef(r WorkItemRegistry, task TaskID, ref *TaskGoalRef) error {
	if ref == nil {
		return nil
	}
	if contentSlug(ref.SpecID) != nil || ref.Spec.Revision < 1 || !digestPattern.MatchString(ref.Spec.ManifestSHA256) {
		return queueError("task_goal_invalid_reference", "invalid exact goal reference")
	}
	for _, s := range r.TaskSpecRevisions {
		if s.TaskID == task && s.SpecID == ref.SpecID && s.Revision == ref.Spec.Revision && s.ManifestSHA256 == ref.Spec.ManifestSHA256 {
			return nil
		}
	}
	return queueError("task_goal_invalid_reference", "goal reference is not registered to this Task")
}

func loadTaskGoal(d Dependencies, root string, r WorkItemRegistry, task TaskID, ref TaskGoalRef) (TaskGoalContract, error) {
	out := TaskGoalContract{TaskID: task, Goal: ref, Requirements: []TaskGoalRequirement{}, Constraints: []string{}, RequiredInputs: []TaskSpecRequiredInput{}}
	if e := validateRegisteredGoalRef(r, task, &ref); e != nil {
		return out, e
	}
	v, e := registeredTaskSpec(d, root, r, task, ref.SpecID, &ref.Spec)
	if e != nil {
		return out, e
	}
	if !isTaskGoalSpec(v) {
		return out, queueError("task_goal_required", "queue goal must identify a base-free goal Spec")
	}
	m := contentFields(v)
	out.Problem = *valueRevision(m["problem"])
	out.Title, out.Objective = contentString(m, "title"), contentString(m, "objective")
	if e = contentDecode(m["executor"], &out.Executor); e != nil {
		return out, e
	}
	if e = contentDecode(m["requirements"], &out.Requirements); e != nil {
		return out, e
	}
	if e = contentDecode(m["constraints"], &out.Constraints); e != nil {
		return out, e
	}
	for _, basis := range []struct{ name, hash string }{{"goal", ref.Spec.ManifestSHA256}, {"problem", out.Problem.ManifestSHA256}} {
		manifest, e := readRegisteredTaskManifest(d.TaskContent, root, r, basis.hash)
		if e != nil {
			return out, e
		}
		bytes, e := canonicaljson.Marshal(manifest)
		if e != nil {
			return out, e
		}
		out.RequiredInputs = append(out.RequiredInputs, TaskSpecRequiredInput{ID: "task-basis/" + basis.name, Role: "other", Locator: taskContentPath(root, "manifests", basis.hash), SHA256: basis.hash, SizeBytes: int64(len(bytes)), MediaType: "application/json"})
		for _, doc := range contentArray(contentFields(manifest), "documents") {
			x := contentFields(doc)
			hash := contentString(x, "sha256")
			bytes, e := d.TaskContent.Read(root, "objects", hash)
			if e != nil || len(bytes) != contentInt(x, "size_bytes") {
				return out, contentError("task_content_integrity_conflict", "goal document differs", e)
			}
			role := "other"
			if basis.name == "goal" {
				role = "design"
			}
			out.RequiredInputs = append(out.RequiredInputs, TaskSpecRequiredInput{ID: "task-goal-doc/" + basis.name + "/" + contentString(x, "id"), Role: role, Locator: taskContentPath(root, "objects", hash), SHA256: hash, SizeBytes: int64(len(bytes)), MediaType: contentString(x, "media_type")})
		}
	}
	sort.Slice(out.RequiredInputs, func(i, j int) bool { return out.RequiredInputs[i].ID < out.RequiredInputs[j].ID })
	return out, nil
}

func validateExecutionGoalOrigin(d Dependencies, root string, r WorkItemRegistry, task TaskID, spec canonicaljson.Value) error {
	m := contentFields(spec)
	if contentInt(m, "schema_version") != 2 || contentString(m, "contract_kind") != "execution" {
		return nil
	}
	var origin TaskGoalRef
	if e := contentDecode(m["goal_origin"], &origin); e != nil {
		return e
	}
	goal, e := loadTaskGoal(d, root, r, task, origin)
	if e != nil {
		return e
	}
	var p TaskGoalExecutePreview
	if e = contentDecode(m["execution_request"], &p); e != nil {
		return e
	}
	if goalExecutionDigest(p) != contentString(m, "execution_request_sha256") || !contentTypedEqual(goal, p.Goal) || contentString(m, "acceptance_path") != p.AcceptancePath {
		return queueError("task_goal_binding_conflict", "execution request differs from its exact goal or acceptance path")
	}
	t, _ := findTask(r, task)
	if t == nil || p.Workspace != root {
		return queueError("task_goal_binding_conflict", "execution request belongs to another workspace or Task")
	}
	epic, _ := findEpic(r, t.ParentEpicID)
	if epic == nil {
		return queueError("task_goal_binding_conflict", "execution Epic is not registered")
	}
	parent, _ := findEpicRepo(*epic, t.RepoID)
	if parent == nil || p.Target != (QueueTarget{ProjectID: t.ProjectID, RepoID: t.RepoID, EpicID: t.ParentEpicID, GitCommonDir: t.GitCommonDir, ParentWorktreeID: parent.Worktree.ID, ParentLocator: parent.Worktree.Locator, ParentRef: parent.Worktree.Ref}) {
		return queueError("task_goal_binding_conflict", "execution request differs from the registered return worktree")
	}
	previous := m["previous"]
	if _, draft := m["expected_previous"]; draft {
		previous = m["expected_previous"]
	}
	if !contentTypedEqual(valueRevision(previous), p.PreviousSpec) {
		return queueError("task_goal_binding_conflict", "execution request differs from its preserved predecessor Spec")
	}
	basis := contentFields(m["implementation_basis"])
	if contentString(basis, "parent_oid") != p.ParentOID || contentString(basis, "parent_tree") != p.ParentTree || contentString(basis, "start_oid") != p.ParentOID || contentString(basis, "start_tree") != p.ParentTree || contentString(basis, "project_id") != string(p.Target.ProjectID) || contentString(basis, "repo_id") != string(p.Target.RepoID) || contentString(basis, "epic_id") != string(p.Target.EpicID) || contentString(basis, "git_common_dir") != p.Target.GitCommonDir || contentString(basis, "parent_worktree_id") != string(p.Target.ParentWorktreeID) || contentString(basis, "parent_ref") != p.Target.ParentRef {
		return queueError("task_goal_binding_conflict", "execution base differs from its frozen request")
	}
	if !contentTypedEqual(valueRevision(m["problem"]), &goal.Problem) || contentString(m, "spec_id") != goal.Goal.SpecID {
		return queueError("task_goal_binding_conflict", "execution must retain the exact goal and problem")
	}
	reqs := contentArray(m, "requirements")
	if len(reqs) != len(goal.Requirements) {
		return queueError("task_goal_binding_conflict", "execution requirements differ from goal")
	}
	for i, req := range reqs {
		x := contentFields(req)
		if contentString(x, "id") != goal.Requirements[i].ID || contentString(x, "acceptance") != goal.Requirements[i].Acceptance {
			return queueError("task_goal_binding_conflict", "execution requirements differ from goal")
		}
	}
	return nil
}
