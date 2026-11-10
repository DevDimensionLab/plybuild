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
	// Preserve omission in historical canonical execution requests. Only an
	// explicit true in a newly selected Goal grants the owner installation.
	AllowLocalInstall *bool `yaml:"allow_local_install,omitempty" json:"allow_local_install,omitempty"`
}

func (a TaskExecutorAssignment) LocalInstallAllowed() bool {
	return a.AllowLocalInstall != nil && *a.AllowLocalInstall
}

func taskGoalExecutorRule(v canonicaljson.Value) error {
	fields := map[string]contentRule{"provider": contentNullable(contentEnum("codex", "claude")), "model": contentNullable(contentText(256)), "effort": contentNullable(contentText(128))}
	if _, present := contentFields(v)["allow_local_install"]; present {
		fields["allow_local_install"] = func(v canonicaljson.Value) error {
			if _, ok := v.(bool); !ok {
				return fmt.Errorf("must be a boolean")
			}
			return nil
		}
	}
	return contentExact(fields)(v)
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
	s["executor"] = taskGoalExecutorRule
	doc := contentDocumentInputRule()
	if published {
		doc = contentDescriptorRule()
	}
	s["documents"] = contentList(31, doc, contentIDKey, false)
	return s
}

func validateTaskGoalSemantics(m map[string]canonicaljson.Value) error {
	if len(contentArray(m, "requirements")) == 0 {
		return fmt.Errorf("a goal needs at least one observable requirement")
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
