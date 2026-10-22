package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type TaskGoalExecuteInput struct {
	Target         QueueTargetInput `json:"target"`
	Next           bool             `json:"next"`
	SpecID         string           `json:"spec_id"`
	AcceptancePath string           `json:"acceptance_path"`
}

type TaskGoalExecutePreview struct {
	SourceMode        string                  `yaml:"source_mode,omitempty" json:"source_mode,omitempty"`
	AdoptionID        WorktreeOperationID     `yaml:"adoption_id,omitempty" json:"adoption_id,omitempty"`
	StartOID          string                  `yaml:"start_oid,omitempty" json:"start_oid,omitempty"`
	StartTree         string                  `yaml:"start_tree,omitempty" json:"start_tree,omitempty"`
	Delivery          *DeliveryAgreement      `json:"delivery,omitempty"`
	DeliveryGates     []string                `json:"delivery_gates,omitempty"`
	DeliveryEffects   []string                `json:"delivery_effects,omitempty"`
	DeliveryStopAfter string                  `json:"delivery_stop_after,omitempty"`
	Kind              string                  `json:"kind"`
	SchemaVersion     int                     `json:"schema_version"`
	Workspace         string                  `json:"workspace"`
	Input             TaskGoalExecuteInput    `json:"input"`
	Goal              TaskGoalContract        `json:"goal"`
	QueueID           string                  `json:"queue_id"`
	QueueRevision     int                     `json:"queue_revision"`
	Target            QueueTarget             `json:"target"`
	BaseRevision      int                     `json:"base_revision"`
	ParentOID         string                  `json:"parent_oid"`
	ParentTree        string                  `json:"parent_tree"`
	Branch            string                  `json:"branch"`
	WorktreePath      string                  `json:"worktree_path"`
	AcceptancePath    string                  `json:"acceptance_path"`
	ObservedParent    PlanWorktreeObservation `json:"observed_parent"`
	PreviousSpec      *TaskRevisionRef        `json:"previous_spec"`
	PreviousSelection *TaskDecisionRef        `json:"previous_selection"`
	RegistryUpgrade   *TaskRegistryUpgrade    `json:"registry_upgrade"`
	Confirmation      string                  `json:"confirmation,omitempty"`
}

type TaskGoalPreparedExecution struct {
	Delivery       *DeliveryAgreement      `json:"delivery,omitempty"`
	Goal           TaskGoalRef             `json:"goal"`
	Preparation    TaskPreparationReadback `json:"preparation"`
	Basis          TaskSpecBasis           `json:"basis"`
	Executor       TaskExecutorAssignment  `json:"executor"`
	AcceptancePath string                  `json:"acceptance_path"`
}

func goalExecutionDigest(p TaskGoalExecutePreview) string {
	p.Confirmation = ""
	return queueDigest(p)
}

func validateGoalExecuteInput(in TaskGoalExecuteInput) error {
	if in.Next == (in.SpecID != "") {
		return WorkInvalidArguments("choose exactly one of next or a goal Spec ID")
	}
	if in.SpecID != "" && contentSlug(in.SpecID) != nil {
		return WorkInvalidArguments("invalid goal Spec ID")
	}
	if in.AcceptancePath == "" {
		return nil
	}
	if e := contentPath(in.AcceptancePath); e != nil {
		return WorkInvalidArguments("acceptance path must be absolute and clean")
	}
	// Preview permits a missing workarea, but never a symlink in its existing
	// ancestry. The launcher owns safe directory creation and script writes.
	for p := in.AcceptancePath; ; p = filepath.Dir(p) {
		st, e := os.Lstat(p)
		if e == nil {
			if st.Mode()&os.ModeSymlink != 0 {
				return WorkInvalidArguments("acceptance path has a symlink component")
			}
			physical, e := filepath.EvalSymlinks(p)
			if e != nil || physical != p {
				return WorkInvalidArguments("acceptance path must have physical ancestry")
			}
			if p != in.AcceptancePath && !st.IsDir() {
				return WorkInvalidArguments("acceptance parent is not a directory")
			}
			if p == in.AcceptancePath && !st.Mode().IsRegular() {
				return WorkInvalidArguments("acceptance path is not a regular file")
			}
			return nil
		}
		if !os.IsNotExist(e) || filepath.Dir(p) == p {
			return WorkInvalidArguments("cannot inspect acceptance ancestry")
		}
	}
}

func storedGoalExecution(d Dependencies, root string, r WorkItemRegistry, in TaskGoalExecuteInput, digest string) (*TaskGoalExecutePreview, *TaskRevisionRef, error) {
	projects, repos, err := d.Projects.Snapshot(root)
	if err != nil {
		return nil, nil, err
	}
	target, _, err := resolveGoalQueueTarget(d, root, r, ProjectSnapshot{Projects: projects, Repos: repos}, in.Target, in.SpecID)
	if err != nil {
		return nil, nil, err
	}
	return storedGoalExecutionInTarget(d, root, r, in, digest, target)
}

func storedGoalExecutionInTarget(d Dependencies, root string, r WorkItemRegistry, in TaskGoalExecuteInput, digest string, target QueueTarget) (*TaskGoalExecutePreview, *TaskRevisionRef, error) {
	// An adopted execution intent precedes the first content publication.
	for _, o := range r.WorktreeOperations {
		if o.Execution != nil && o.Release == nil {
			p := *o.Execution
			comparison := in
			if comparison.AcceptancePath == "" {
				comparison.AcceptancePath = p.AcceptancePath
			}
			if p.Target == target && contentTypedEqual(p.Input, comparison) && (digest == "" || digest == p.Confirmation) {
				return &p, nil, nil
			}
		}
	}
	candidates := goalExecutionCandidates(root, r, in, target)
	// A known confirmation can name a completed execution outside today's
	// queue. Its native publication key is authoritative recovery metadata.
	if digest != "" {
		if pub := taskPublication(r, "goal-execute/"+strings.TrimPrefix(digest, "sha256:")+"/spec"); pub != nil {
			task, _ := findTask(r, pub.TaskID)
			if task != nil && task.ProjectID == target.ProjectID && task.RepoID == target.RepoID && task.ParentEpicID == target.EpicID {
				candidates[task.ID] = true
			}
		}
	}

	for _, s := range r.TaskSpecRevisions {
		if !candidates[s.TaskID] || (in.SpecID != "" && s.SpecID != in.SpecID) {
			continue
		}
		v, e := readRegisteredTaskManifest(d.TaskContent, root, r, s.ManifestSHA256)
		if e != nil {
			return nil, nil, e
		}
		m := contentFields(v)
		if contentString(m, "contract_kind") != "execution" {
			continue
		}
		var p TaskGoalExecutePreview
		if e = contentDecode(m["execution_request"], &p); e != nil {
			return nil, nil, e
		}
		if digest != "" && contentString(m, "execution_request_sha256") != digest {
			continue
		}
		comparison := in
		if comparison.AcceptancePath == "" {
			// An omitted path is derived at preview time. Reusing the same
			// public input must recover that preserved choice, not require the
			// caller to reconstruct a field it never supplied.
			comparison.AcceptancePath = p.AcceptancePath
		}
		if digest == "" && !contentTypedEqual(p.Input, comparison) {
			continue
		}
		if !contentTypedEqual(p.Input, comparison) {
			return nil, nil, queueError("task_goal_execution_conflict", "confirmation belongs to another selector or acceptance path")
		}
		p.Confirmation = goalExecutionDigest(p)
		if p.Confirmation != contentString(m, "execution_request_sha256") {
			return nil, nil, queueError("task_goal_execution_integrity", "preserved execution request differs")
		}
		ref := TaskRevisionRef{s.Revision, s.ManifestSHA256}
		return &p, &ref, nil
	}
	return nil, nil, nil
}

func buildGoalExecutionPreview(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in TaskGoalExecuteInput) (TaskGoalExecutePreview, error) {
	target, repo, e := resolveGoalQueueTarget(d, root, r, projects, in.Target, in.SpecID)
	if e != nil {
		return TaskGoalExecutePreview{}, e
	}
	return buildGoalExecutionPreviewInTarget(d, root, r, projects, in, target, repo)
}

func buildGoalExecutionPreviewInTarget(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in TaskGoalExecuteInput, target QueueTarget, repo RepoRecord) (TaskGoalExecutePreview, error) {
	var e error
	if in.AcceptancePath != "" {
		if old, _, e := storedGoalExecutionInTarget(d, root, r, in, "", target); e != nil {
			return TaskGoalExecutePreview{}, e
		} else if old != nil {
			return *old, nil
		}
	}
	out := TaskGoalExecutePreview{Kind: "WorkspaceTaskGoalExecutePlan@1", SchemaVersion: 1, Workspace: root, Input: in, AcceptancePath: in.AcceptancePath}
	out.Target, out.QueueID = target, queueID(root, target)
	q := foldQueue(r, out.QueueID)
	if q.Current != nil {
		prep := findPreparation(r, *q.Current)
		if old, e := resumeGoalPreview(d, root, r, in, prep.Plan.TaskID, target); e != nil {
			return out, e
		} else if old != nil {
			return *old, nil
		}
		return out, queueError("task_queue_current_conflict", "complete the current preparation before choosing another goal")
	}
	out.QueueRevision = q.Revision
	// Recover a partially published execution before evaluating the planner's
	// current order. No worktree reservation exists yet at this boundary, and a
	// newly inserted goal must not make an interrupted --next choose new work.
	var interrupted *TaskGoalExecutePreview
	for _, entry := range q.Pending {
		if preparationForTask(r, entry.TaskID) != nil {
			continue
		}
		old, e := resumeGoalPreview(d, root, r, in, entry.TaskID, target)
		if e != nil {
			return out, e
		}
		if old != nil {
			if interrupted != nil && interrupted.Confirmation != old.Confirmation {
				return out, queueError("task_goal_execution_conflict", "multiple unfinished goal executions match; inspect their preserved requests before selecting work")
			}
			interrupted = old
		}
	}
	if interrupted != nil {
		return *interrupted, nil
	}
	var chosen *TaskGoalQueueEntry
	if in.SpecID != "" {
		for _, entry := range q.Pending {
			if entry.Goal != nil && entry.Goal.SpecID == in.SpecID {
				if chosen != nil {
					return out, queueError("task_goal_ambiguous", "goal Spec ID belongs to multiple pending Tasks in this target")
				}
				v := entry
				chosen = &v
			}
		}
		if chosen == nil {
			return out, queueError("task_goal_not_pending", "goal Spec ID is not pending in this target")
		}
	} else {
		for i, entry := range q.Pending {
			if entry.Goal == nil {
				continue
			}
			if queueGoalEvaluation(d, root, r, entry, i+1).State == "ready" {
				v := entry
				chosen = &v
				break
			}
		}
		if chosen == nil {
			return out, queueError("task_goal_no_ready", "no eligible pending goal; inspect the target queue")
		}
	}
	row := queueGoalEvaluation(d, root, r, *chosen, 1)
	if row.State != "ready" {
		return out, queueError(row.Reasons[0].Code, row.Reasons[0].Message)
	}
	out.Goal, e = loadTaskGoal(d, root, r, chosen.TaskID, *chosen.Goal)
	if e != nil {
		return out, e
	}
	if in.AcceptancePath == "" {
		key := strings.TrimPrefix(queueDigest(map[string]any{"queue_id": out.QueueID, "task_id": chosen.TaskID, "goal": chosen.Goal}), "sha256:")
		in.AcceptancePath = filepath.Join(root, MarkerDirectory, "task-executions", key, "acceptance.sh")
		out.Input, out.AcceptancePath = in, in.AcceptancePath
	}
	if e = validateGoalExecuteInput(in); e != nil {
		return out, e
	}
	if old, _, e := storedGoalExecutionInTarget(d, root, r, in, "", target); e != nil {
		return out, e
	} else if old != nil {
		return *old, nil
	}
	if len(queueUnresolved(r, target)) != 0 {
		return out, queueError("task_queue_effect_pending", "resolve outstanding target effects first")
	}
	epic, binding := queueTargetBinding(r, target)
	base := currentEpicBase(r, epic, binding)
	out.BaseRevision, out.ParentOID, out.ParentTree = base.Revision, base.OID, base.Tree
	task, _ := findTask(r, chosen.TaskID)
	title, e := queueTaskTitle(d, root, r, *task)
	if e != nil {
		return out, e
	}
	out.Branch, out.WorktreePath = queueTaskNames(out.QueueID, chosen.TaskID, title, target.ParentLocator)
	if op, _ := findOperation(r, task.ID); op != nil && op.Adoption != nil && op.State == "ready" {
		if !contentTypedEqual(op.Adoption.Plan.Goal, chosen.Goal) && op.Adoption.Plan.Goal != *chosen.Goal {
			return out, queueError("adoption_goal_stale", "adoption pins another Goal")
		}
		if _, e = observeAdoption(d, root, r, projects, op.Adoption.Plan.Input, op.ID); e != nil {
			return out, e
		}
		out.SourceMode = "adopted"
		out.AdoptionID = op.ID
		out.StartOID = operationStartOID(*op)
		out.StartTree = operationStartTree(*op)
		out.Branch = strings.TrimPrefix(op.SourceRef, "refs/heads/")
		out.WorktreePath = op.TargetLocator
	}
	if out.Goal.Delivery != nil {
		bound, err := BindDeliveryAgreement(*out.Goal.Delivery, target, "refs/heads/"+out.Branch)
		if err != nil {
			return out, queueError("task_delivery_binding_conflict", err.Error())
		}
		out.Delivery = &bound
		out.DeliveryGates, out.DeliveryEffects, out.DeliveryStopAfter = bound.RequiredGates(), bound.AllowedEffects(), bound.StopAfter()
	}
	if out.SourceMode != "adopted" {
		if e = precheckParentAndResources(d, repo, currentEpicBinding(r, epic, binding), "refs/heads/"+out.Branch, out.WorktreePath); e != nil {
			return out, e
		}
	}
	out.ObservedParent, e = queueParentObservation(d, target)
	if e != nil {
		return out, e
	}
	state := taskContentState(r, chosen.TaskID)
	out.PreviousSpec, out.PreviousSelection, out.RegistryUpgrade = state.SpecHead, state.SelectionEvent, queueUpgrade(r)
	out.Confirmation = goalExecutionDigest(out)
	return out, nil
}

func PreviewTaskGoalExecution(d Dependencies, in TaskGoalExecuteInput) (TaskGoalExecutePreview, error) {
	if e := validateGoalExecuteInput(in); e != nil {
		return TaskGoalExecutePreview{}, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return TaskGoalExecutePreview{}, e
	}
	var out TaskGoalExecutePreview
	e = WithTaskSpecRegistrationSnapshot(d, root, func(s *TaskSpecSession) error {
		var e error
		out, e = buildGoalExecutionPreview(d, root, s.Registry, s.Projects, in)
		if e != nil {
			return e
		}
		return validateGoalOperationContent(d, root, s.Registry, out.Goal.TaskID)
	})
	return out, e
}

func goalDraftBase(task TaskID, kind, key, at string) map[string]canonicaljson.Value {
	return map[string]canonicaljson.Value{"kind": kind, "schema_version": int64(1), "format": "json", "format_version": int64(1), "canonicalization": "RFC8785", "publication_key": key, "task_id": string(task), "registry_upgrade": nil,
		"recorder": contentObject(map[string]canonicaljson.Value{"actor_claim": "Ply goal execution", "control_surface": "ply workflow execute", "recorded_at_utc": at})}
}

func publishGoalExecutionContent(d Dependencies, root string, session WorkItemStoreSession, r WorkItemRegistry, projects ProjectSnapshot, task TaskRecord, operation string, m map[string]canonicaljson.Value) (TaskContentOutcomeRef, error) {
	key := contentString(m, "publication_key")
	if old := taskPublication(r, key); old != nil {
		if old.TaskID != task.ID || old.Operation != operation {
			return TaskContentOutcomeRef{}, queueError("task_goal_execution_conflict", "execution publication key belongs to other work")
		}
		return old.OutcomeRef, nil
	}
	raw, e := canonicaljson.Marshal(contentObject(m))
	if e != nil {
		return TaskContentOutcomeRef{}, e
	}
	draft, e := decodeTaskContent(raw, operation, false, false)
	if e != nil {
		return TaskContentOutcomeRef{}, e
	}
	result := newTaskContentMutation(root, task.ID, key, digestTaskBytes(raw))
	result.PreState = taskContentState(r, task.ID)
	if e = publishTaskContent(d, root, session, r, task, projects, draft, operation, raw, &result, false); e != nil {
		return TaskContentOutcomeRef{}, e
	}
	return *result.OutcomeRef, nil
}

// materializeGoalExecution retains a full immutable request before the first Git
// effect. Partial metadata publication is recovered by its publication keys.
func materializeGoalExecution(d Dependencies, root string, session WorkItemStoreSession, r WorkItemRegistry, projects ProjectSnapshot, p TaskGoalExecutePreview, human QueueHumanDecision) (TaskDecisionRef, error) {
	if r.FormatVersion < 4 {
		if e := upgradeQueueRegistry(d, root, &r, p.RegistryUpgrade); e != nil {
			return TaskDecisionRef{}, e
		}
	}
	task, _ := findTask(r, p.Goal.TaskID)
	if task == nil {
		return TaskDecisionRef{}, queueError("task_goal_missing", "Task is no longer registered")
	}
	prefix := "goal-execute/" + strings.TrimPrefix(p.Confirmation, "sha256:")
	at := queueNow(d)
	specDraft, e := buildGoalExecutionSpec(d, root, r, p, prefix+"/spec", at)
	if e != nil {
		return TaskDecisionRef{}, e
	}
	spec, e := publishGoalExecutionContent(d, root, session, r, projects, *task, "spec_record", specDraft)
	if e != nil {
		return TaskDecisionRef{}, e
	}
	r, e = session.Snapshot()
	if e != nil {
		return TaskDecisionRef{}, e
	}
	specRef := TaskRevisionRef{*spec.Revision, spec.ManifestSHA256}
	checks := []canonicaljson.Value{}
	for _, name := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, contentObject(map[string]canonicaljson.Value{"id": name, "outcome": "pass", "reason": "The execution contract preserves the chosen outcomes and declares the owner's work and acceptance entrypoint; implementation and test results are not yet claimed.", "evidence_document_ids": []canonicaljson.Value{}}))
	}
	assessmentDraft := goalDraftBase(task.ID, "WorkspaceTaskSpecAssessmentDraft@1", prefix+"/assessment", at)
	for k, v := range map[string]canonicaljson.Value{"spec_id": p.Goal.Goal.SpecID, "spec": contentRefValue(specRef), "expected_previous_assessment": nil, "outcome": "ready", "reason": "Ready to begin the authorized implementation; detailed design, tests and product judgment remain to be performed.", "open_questions": []canonicaljson.Value{}, "checks": checks, "documents": []canonicaljson.Value{}} {
		assessmentDraft[k] = v
	}
	assessment, e := publishGoalExecutionContent(d, root, session, r, projects, *task, "spec_assess", assessmentDraft)
	if e != nil {
		return TaskDecisionRef{}, e
	}
	r, e = session.Snapshot()
	if e != nil {
		return TaskDecisionRef{}, e
	}
	selectionDraft := goalDraftBase(task.ID, "WorkspaceTaskSolutionSelectionDraft@1", prefix+"/selection", at)
	for k, v := range map[string]canonicaljson.Value{"expected_previous_selection": contentRefValue(p.PreviousSelection), "action": "select", "reason": "The human execute command chose this exact goal and late-bound execution contract.", "human_decision": contentRefValue(human), "solution": contentObject(map[string]canonicaljson.Value{"spec_id": p.Goal.Goal.SpecID, "spec": contentRefValue(specRef), "problem": contentRefValue(p.Goal.Problem), "assessment": contentRefValue(TaskDecisionRef{*assessment.ID, assessment.ManifestSHA256})})} {
		selectionDraft[k] = v
	}
	selection, e := publishGoalExecutionContent(d, root, session, r, projects, *task, "spec_select", selectionDraft)
	if e != nil {
		return TaskDecisionRef{}, e
	}
	return TaskDecisionRef{*selection.ID, selection.ManifestSHA256}, nil
}

func resumeGoalPreview(d Dependencies, root string, r WorkItemRegistry, in TaskGoalExecuteInput, task TaskID, target QueueTarget) (*TaskGoalExecutePreview, error) {
	for i := len(r.TaskSpecRevisions) - 1; i >= 0; i-- {
		ref := r.TaskSpecRevisions[i]
		if ref.TaskID != task {
			continue
		}
		v, e := readRegisteredTaskManifest(d.TaskContent, root, r, ref.ManifestSHA256)
		if e != nil {
			return nil, e
		}
		m := contentFields(v)
		if contentString(m, "contract_kind") != "execution" {
			continue
		}
		var p TaskGoalExecutePreview
		if e = contentDecode(m["execution_request"], &p); e != nil {
			return nil, e
		}
		if in.SpecID != "" && in.SpecID != p.Goal.Goal.SpecID {
			continue
		}
		if in.AcceptancePath == "" {
			in.AcceptancePath = p.AcceptancePath
		}
		old, _, e := storedGoalExecutionInTarget(d, root, r, in, contentString(m, "execution_request_sha256"), target)
		return old, e
	}
	return nil, nil
}

func goalExecutionRequestRule(v canonicaljson.Value) error {
	var p TaskGoalExecutePreview
	if e := contentDecode(v, &p); e != nil {
		return e
	}
	if !contentEqual(v, contentRefValue(p)) || p.Kind != "WorkspaceTaskGoalExecutePlan@1" || p.SchemaVersion != 1 || p.Confirmation != "" || p.Input.Next == (p.Input.SpecID != "") || contentPath(p.Workspace) != nil || contentPath(p.AcceptancePath) != nil || p.Input.AcceptancePath != p.AcceptancePath || p.QueueID != queueID(p.Workspace, p.Target) || p.QueueRevision < 0 || p.BaseRevision < 1 || !validOIDText(p.ParentOID) || !validOIDText(p.ParentTree) || !validShortBranchText(p.Branch) || contentPath(p.WorktreePath) != nil || p.PreviousSpec == nil || p.Goal.TaskID == "" {
		return fmt.Errorf("invalid preserved goal execution request")
	}
	if p.SourceMode == "" {
		if p.AdoptionID != "" || p.StartOID != "" || p.StartTree != "" {
			return fmt.Errorf("create source cannot carry adoption start")
		}
	} else if p.SourceMode != "adopted" || !operationIDPattern.MatchString(string(p.AdoptionID)) || !validOIDText(p.StartOID) || !validOIDText(p.StartTree) || p.StartOID == p.ParentOID || p.StartTree == p.ParentTree {
		return fmt.Errorf("invalid adopted source request")
	}
	if p.Input.SpecID != "" && p.Input.SpecID != p.Goal.Goal.SpecID {
		return fmt.Errorf("goal selector differs from request")
	}
	if e := validateQueuePlanObservation(p.ObservedParent, p.Target, p.ParentOID, p.ParentTree); e != nil {
		return e
	}
	if e := validateGoalDeliveryBinding(p); e != nil {
		return e
	}
	return nil
}

func validateGoalDeliveryBinding(p TaskGoalExecutePreview) error {
	if p.Goal.Delivery == nil {
		if p.Delivery != nil || len(p.DeliveryGates) != 0 || len(p.DeliveryEffects) != 0 || p.DeliveryStopAfter != "" {
			return fmt.Errorf("legacy goal cannot acquire a delivery agreement during execution")
		}
		return nil
	}
	bound, e := BindDeliveryAgreement(*p.Goal.Delivery, p.Target, "refs/heads/"+p.Branch)
	if e != nil {
		return e
	}
	if p.Delivery == nil || !contentTypedEqual(bound, *p.Delivery) {
		return fmt.Errorf("preserved delivery differs from exact goal/source/target")
	}
	if !contentTypedEqual(p.DeliveryGates, bound.RequiredGates()) || !contentTypedEqual(p.DeliveryEffects, bound.AllowedEffects()) || p.DeliveryStopAfter != bound.StopAfter() {
		return fmt.Errorf("preserved delivery gates, effects or stop point differ from agreement")
	}
	return nil
}

func buildGoalExecutionSpec(d Dependencies, root string, r WorkItemRegistry, p TaskGoalExecutePreview, key, at string) (map[string]canonicaljson.Value, error) {
	goal, e := registeredTaskSpec(d, root, r, p.Goal.TaskID, p.Goal.Goal.SpecID, &p.Goal.Goal.Spec)
	if e != nil {
		return nil, e
	}
	gm := contentFields(goal)
	docs, supporting := []canonicaljson.Value{}, []canonicaljson.Value{}
	ids := map[string]bool{}
	for _, v := range contentArray(gm, "documents") {
		id := contentString(contentFields(v), "id")
		ids[id] = true
		docs = append(docs, contentObject(map[string]canonicaljson.Value{"id": id, "source": contentObject(map[string]canonicaljson.Value{"kind": "snapshot", "manifest_sha256": p.Goal.Goal.Spec.ManifestSHA256, "document_id": id})}))
		supporting = append(supporting, contentObject(map[string]canonicaljson.Value{"document_id": id, "usage": "normative"}))
	}
	id := "execution-contract"
	for ids[id] {
		id += "-generated"
	}
	text := "# " + p.Goal.Title + "\n\n## Goal\n\n" + p.Goal.Objective + "\n\n## Outcomes\n\n"
	for _, req := range p.Goal.Requirements {
		text += "- " + req.ID + ": " + req.Acceptance + "\n"
	}
	if p.SourceMode == "adopted" {
		text += "\n## Adopted source\n\nObserved existing source " + p.WorktreePath + " on refs/heads/" + p.Branch + " starts at " + p.StartOID + " (tree " + p.StartTree + "). Chosen Epic parent: " + p.ParentOID + ". Review, assess, adapt and test the entire inherited delta from parent to source, plus subsequent work. No prior session, mandate, tests or qualification is inherited. Ordinary native delivery includes cleanup of this now-owned pre-existing branch and worktree after owner exit.\n"
	}
	text += "\n## Implementation ownership\n\nThe implementing owner investigates the code, chooses and revises the detailed design, implements meaningful acceptance tests, and owns verification and correction within this goal and actual runtime permissions. The planner has not chosen these technical details.\n\nThe declared acceptance entrypoint is /bin/sh " + p.AcceptancePath + " with the Task worktree as cwd. The owner must implement this script and meaningful tests, execute it, and preserve the actual script and output as artifacts. A declaration is not an executed test or a passing result.\n\nPreserve work and evidence on interruption. Do not reset, remove worktrees, expand scope, or claim human judgment. Human product judgment and local integration are later distinct facts.\n"
	if p.Delivery != nil {
		text += "\n## Delivery agreement\n\nMode: " + p.Delivery.Mode + ". Source: " + p.Delivery.SourceRef + ". Target: " + p.Delivery.TargetRef + ".\n\nGates: " + strings.Join(p.DeliveryGates, ", ") + ". Allowed effects require matching accepted native runtime authority: " + strings.Join(p.DeliveryEffects, ", ") + ". Stop after " + p.DeliveryStopAfter + ".\n"
		if p.Delivery.TargetWorktree != "" {
			text += "\nExact local target worktree: " + p.Delivery.TargetWorktree + ".\n"
		}
		if p.Delivery.GitHubRepository != "" {
			text += "\nExact GitHub repository: " + p.Delivery.GitHubRepository + "; remote: " + p.Delivery.Remote + ".\n"
		}
	}
	if len(p.Goal.Constraints) > 0 {
		text += "\n## Constraints\n\n"
		for _, c := range p.Goal.Constraints {
			text += "- " + c + "\n"
		}
	}
	bytes := []byte(text)
	hash, e := d.TaskContent.Publish(root, "objects", bytes)
	if e != nil {
		return nil, e
	}
	docs = append(docs, contentObject(map[string]canonicaljson.Value{"id": id, "source": contentObject(map[string]canonicaljson.Value{"kind": "file", "locator": taskContentPath(root, "objects", hash), "sha256": hash, "size_bytes": int64(len(bytes)), "media_type": "text/markdown", "git_provenance": nil})}))
	sort.Slice(docs, func(i, j int) bool { return contentIDKey(docs[i]) < contentIDKey(docs[j]) })
	partRefs := func(section string) []canonicaljson.Value {
		return []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"document_id": id, "section": section})}
	}
	part := func(section string) canonicaljson.Value {
		return contentObject(map[string]canonicaljson.Value{"state": "present", "reason": nil, "documents": partRefs(section)})
	}
	reqs, reqIDs := []canonicaljson.Value{}, []canonicaljson.Value{}
	for _, req := range p.Goal.Requirements {
		reqIDs = append(reqIDs, req.ID)
		reqs = append(reqs, contentObject(map[string]canonicaljson.Value{"id": req.ID, "functional_refs": partRefs("Outcomes"), "acceptance": req.Acceptance, "verification_ids": []canonicaljson.Value{"goal-acceptance"}, "technical_refs": partRefs("Implementation ownership")}))
	}
	removed := []string{}
	if p.PreviousSpec != nil {
		previous, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.PreviousSpec.ManifestSHA256)
		if e != nil {
			return nil, e
		}
		current := map[string]bool{}
		for _, req := range p.Goal.Requirements {
			current[req.ID] = true
		}
		for _, req := range contentArray(contentFields(previous), "requirements") {
			if !current[contentIDKey(req)] {
				removed = append(removed, contentIDKey(req))
			}
		}
	}
	frozen := p
	frozen.Confirmation = ""
	m := goalDraftBase(p.Goal.TaskID, "WorkspaceTaskSpecDraft@2", key, at)
	for k, v := range map[string]canonicaljson.Value{
		"schema_version": int64(2), "contract_kind": "execution", "spec_id": p.Goal.Goal.SpecID, "problem": contentRefValue(p.Goal.Problem), "title": p.Goal.Title, "change_reason": "Materialize the human's exact goal at the execution base; implementation choices remain delegated.", "expected_previous": contentRefValue(p.PreviousSpec),
		"documents": docs, "parts": contentObject(map[string]canonicaljson.Value{"abstract": part("Goal"), "functional": part("Outcomes"), "technical": part("Implementation ownership")}), "supporting": supporting, "requirements": reqs, "removed_requirement_ids": sortedContentStrings(removed),
		"phases":               []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "implement-and-verify", "purpose": "Investigate, implement and verify the chosen goal.", "requirement_ids": reqIDs, "entry_criteria": []canonicaljson.Value{"Exact authorized goal and clean execution base are bound."}, "exit_criteria": []canonicaljson.Value{"Meaningful acceptance tests have run and the candidate is ready for actual product judgment."}, "verification_ids": []canonicaljson.Value{"goal-acceptance"}})},
		"implementation_basis": contentObject(map[string]canonicaljson.Value{"project_id": string(p.Target.ProjectID), "repo_id": string(p.Target.RepoID), "git_common_dir": p.Target.GitCommonDir, "epic_id": string(p.Target.EpicID), "parent_worktree_id": string(p.Target.ParentWorktreeID), "parent_ref": p.Target.ParentRef, "parent_oid": p.ParentOID, "parent_tree": p.ParentTree, "start_oid": goalStartOID(p), "start_tree": goalStartTree(p)}), "dependencies": []canonicaljson.Value{},
		"goal_origin": contentRefValue(p.Goal.Goal), "execution_request_sha256": p.Confirmation, "execution_request": contentRefValue(frozen), "acceptance_path": p.AcceptancePath,
	} {
		m[k] = v
	}
	if p.Delivery != nil {
		m["delivery"] = contentRefValue(*p.Delivery)
	}
	return m, nil
}

func bindGoalExecutionQueue(d Dependencies, root string, session WorkItemStoreSession, r WorkItemRegistry, p TaskGoalExecutePreview, selection TaskDecisionRef) error {
	key := "goal-execute/" + strings.TrimPrefix(p.Confirmation, "sha256:") + "/queue"
	for _, ev := range r.TaskQueueEvents {
		if ev.QueueID == p.QueueID && ev.Kind == "set" && contentString(contentFields(queueRequestValue(ev.Request)), "publication_key") == key {
			return nil
		}
	}
	q := foldQueue(r, p.QueueID)
	if q.Current != nil || q.Revision < p.QueueRevision {
		return queueError("task_queue_revision_conflict", "another preparation owns this queue or its preserved history changed")
	}
	// The execution Spec has already frozen the human's chosen goal. A planner
	// may add or reorder independent entries while an interrupted publication
	// is recovered. Keep their current queue, but require the exact chosen goal
	// below; removing or replacing it still stops before any Git effect.
	entries := []canonicaljson.Value{}
	found := false
	for _, entry := range q.Pending {
		m := map[string]canonicaljson.Value{"task_id": string(entry.TaskID)}
		if entry.TaskID == p.Goal.TaskID {
			if !contentTypedEqual(entry.Goal, &p.Goal.Goal) {
				return queueError("task_goal_execution_conflict", "pending goal changed")
			}
			m["selection"] = contentRefValue(selection)
			found = true
		} else if entry.Goal != nil {
			m["goal"] = contentRefValue(entry.Goal)
		} else {
			m["selection"] = contentRefValue(entry.Selection)
		}
		entries = append(entries, contentObject(m))
	}
	if !found {
		return queueError("task_goal_not_pending", "chosen goal is no longer pending")
	}
	request := contentObject(map[string]canonicaljson.Value{"kind": "WorkspaceTaskQueueDraft@2", "schema_version": int64(2), "publication_key": key, "project_id": string(p.Target.ProjectID), "repo_id": string(p.Target.RepoID), "epic_id": string(p.Target.EpicID), "expected_revision": int64(q.Revision), "entries": entries, "recorder": contentObject(map[string]canonicaljson.Value{"actor_claim": "Ply goal execution", "control_surface": "ply workflow execute", "recorded_at_utc": queueNow(d)}), "registry_upgrade": nil})
	requestMap := queueRequest(request)
	r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{p.QueueID, q.Revision + 1, "set", queueDigest(requestMap), requestMap, queueNow(d)})
	sortQueueRegistry(&r)
	return session.Publish(r)
}

func PrepareTaskGoalExecution(d Dependencies, in TaskGoalExecuteInput, confirmation string, human QueueHumanDecision) (TaskGoalPreparedExecution, error) {
	var out TaskGoalPreparedExecution
	if e := validateGoalExecuteInput(in); e != nil {
		return out, e
	}
	if !digestPattern.MatchString(confirmation) {
		return out, queueError("task_goal_confirmation_required", "exact goal execution preview digest is required")
	}
	if e := contentExact(map[string]contentRule{"actor_claim": contentText(256), "decided_at_utc": contentUTC, "source": contentEnum("human_cli", "explicit_human_instruction"), "statement": contentText(2000)})(contentRefValue(human)); e != nil {
		return out, WorkInvalidArguments("an explicit human execution decision is required: " + e.Error())
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	d, bindContent := beginGoalOperationContent(d)
	var p TaskGoalExecutePreview
	var selected TaskDecisionRef
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			target, repo, e := resolveGoalQueueTarget(d, root, r, projects, in.Target, in.SpecID)
			if e != nil {
				return e
			}
			old, _, e := storedGoalExecutionInTarget(d, root, r, in, confirmation, target)
			if e != nil {
				return e
			}
			if old != nil {
				p = *old
			} else {
				p, e = buildGoalExecutionPreviewInTarget(d, root, r, projects, in, target, repo)
				if e != nil {
					return e
				}
				if p.Confirmation != confirmation {
					return queueError("task_goal_execution_conflict", "goal, queue, base or destination changed before execution")
				}
			}
			if e = bindContent(root, r, p.Goal.TaskID); e != nil {
				return e
			}
			if prep := preparationForTask(r, p.Goal.TaskID); prep != nil {
				if prep.Plan.ParentOID != p.ParentOID || prep.Plan.WorktreePath != p.WorktreePath {
					return queueError("task_goal_execution_conflict", "Task preparation belongs to another execution")
				}
				selected = prep.Plan.Selection
				return nil
			}
			// Recheck the chosen identity, including on a retry after the goal
			// queue entry has become a selected execution solution.
			if target != p.Target {
				return queueError("task_goal_execution_conflict", "execution target changed")
			}
			ep, b := queueTargetBinding(r, target)
			base := currentEpicBase(r, ep, b)
			observed, e := queueParentObservation(d, target)
			if e != nil {
				return e
			}
			if base.OID != p.ParentOID || base.Tree != p.ParentTree || observed.OID != p.ParentOID || observed.Tree != p.ParentTree {
				return queueError("epic_base_changed", "execution parent changed; preserved request cannot silently change base")
			}
			if !contentTypedEqual(taskContentState(r, p.Goal.TaskID).ProblemHead, &p.Goal.Problem) {
				return queueError("task_goal_problem_changed", "chosen problem changed before execution")
			}
			if p.SourceMode == "adopted" {
				if e = validateAdoptedExecution(d, root, r, projects, p); e != nil {
					return e
				}
			}
			if p.SourceMode == "adopted" {
				op, oi := adoptionOperation(r, p.AdoptionID)
				if op == nil || op.Release != nil {
					return queueError("adoption_generation_stale", "adoption was released before materialization")
				}
				if op.Execution == nil {
					frozen := p
					r.WorktreeOperations[oi].Execution = &frozen
					if e = session.Publish(r); e != nil {
						return e
					}
					r, e = session.Snapshot()
					if e != nil {
						return e
					}
				} else if !contentTypedEqual(*op.Execution, p) {
					return queueError("adoption_execution_conflict", "another execution request owns this adoption")
				}
			}
			selected, e = materializeGoalExecution(d, root, session, r, projects, p, human)
			if e != nil {
				return e
			}
			r, e = session.Snapshot()
			if e != nil {
				return e
			}
			return bindGoalExecutionQueue(d, root, session, r, p, selected)
		})
	})
	if e != nil {
		return out, e
	}
	id := p.Goal.TaskID
	prepareInput := TaskPrepareInput{Target: QueueTargetInput{p.Target.ProjectID, p.Target.RepoID, p.Target.EpicID}, Selector: QueueSelector{Kind: "task", TaskID: &id}}
	prepared, e := PrepareTask(d, prepareInput)
	if e != nil {
		return out, e
	}
	if prepared.Plan == nil || prepared.Plan.Target != p.Target || prepared.Plan.TaskID != id || prepared.Plan.ParentOID != p.ParentOID || prepared.Plan.ParentTree != p.ParentTree || prepared.Plan.WorktreePath != p.WorktreePath || prepared.Plan.Selection != selected {
		return out, queueError("task_goal_execution_conflict", "prepared binding differs from the chosen execution")
	}
	if prepared.Preparation == nil || prepared.Preparation.Outcome == nil {
		if prepared.Confirmation == nil {
			return out, queueError("task_goal_execution_conflict", "execution has no recoverable preparation")
		}
		prepareInput.Apply = true
		prepareInput.Confirmation = *prepared.Confirmation
		if prepared.Plan.RegistryUpgrade != nil {
			prepareInput.Upgrade = prepared.Plan.RegistryUpgrade.RegistrySHA256
		}
		prepared, e = PrepareTask(d, prepareInput)
		if e != nil {
			return out, e
		}
	}
	if prepared.Preparation == nil || prepared.Preparation.Outcome == nil {
		return out, queueError("task_goal_execution_unknown", "execution worktree is not confirmed prepared")
	}
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		task, _ := findTask(s.Registry, id)
		if task == nil || task.Worktree == nil {
			return queueError("task_goal_execution_unknown", "prepared Task binding is unavailable")
		}
		digest, e := taskDeliveryDigest(*task)
		if e != nil {
			return e
		}
		plan := prepared.Preparation.Plan
		out = TaskGoalPreparedExecution{Delivery: p.Delivery, Goal: p.Goal.Goal, Preparation: prepared, Basis: TaskSpecBasis{TaskID: id, Problem: plan.Problem, SpecID: plan.SpecID, Spec: plan.Spec, Assessment: plan.Assessment, Selection: plan.Selection, TaskWorktreeID: task.Worktree.ID, DeliveryBindingSHA256: digest, Dependencies: []string{}}, Executor: p.Goal.Executor, AcceptancePath: p.AcceptancePath}
		return nil
	})
	return out, e
}

func goalStartOID(p TaskGoalExecutePreview) string {
	if p.SourceMode == "adopted" {
		return p.StartOID
	}
	return p.ParentOID
}
func goalStartTree(p TaskGoalExecutePreview) string {
	if p.SourceMode == "adopted" {
		return p.StartTree
	}
	return p.ParentTree
}
func validateAdoptedExecution(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, p TaskGoalExecutePreview) error {
	o, _ := findOperation(r, p.Goal.TaskID)
	if o == nil || o.Adoption == nil || o.ID != p.AdoptionID || o.State != "ready" || o.Release != nil || o.TargetLocator != p.WorktreePath || o.SourceRef != "refs/heads/"+p.Branch || operationStartOID(*o) != p.StartOID || operationStartTree(*o) != p.StartTree || o.ParentOID != p.ParentOID || o.ParentTree != p.ParentTree || o.Adoption.Plan.Goal != p.Goal.Goal {
		return queueError("adoption_generation_stale", "execution preview no longer owns active adoption; inspect or preview again")
	}
	source, e := observeAdoptionSource(d, p.WorktreePath, o.SourceRef)
	if e != nil {
		return e
	}
	if !source.Clean || !source.InventoryMatch || source.OID != p.StartOID || source.Tree != p.StartTree || source.Ref != o.SourceRef || source.GitCommonDir != o.GitCommonDir || source.Locator != o.TargetLocator {
		return queueError("adoption_source_changed", "source differs from adopted start")
	}
	parent, e := queueParentObservation(d, p.Target)
	if e != nil {
		return e
	}
	if parent.OID != p.ParentOID || parent.Tree != p.ParentTree {
		return queueError("adoption_base_changed", "parent differs from adoption")
	}
	_, e = adoptionWriters(d, p.WorktreePath)
	return e
}
