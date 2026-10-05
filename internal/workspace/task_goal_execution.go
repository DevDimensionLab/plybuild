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
	for _, s := range r.TaskSpecRevisions {
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
	if in.AcceptancePath != "" {
		if old, _, e := storedGoalExecution(d, root, r, in, ""); e != nil {
			return TaskGoalExecutePreview{}, e
		} else if old != nil {
			return *old, nil
		}
	}
	out := TaskGoalExecutePreview{Kind: "WorkspaceTaskGoalExecutePlan@1", SchemaVersion: 1, Workspace: root, Input: in, AcceptancePath: in.AcceptancePath}
	target, repo, e := resolveQueueTarget(d, root, r, projects, in.Target)
	if e != nil {
		return out, e
	}
	out.Target, out.QueueID = target, queueID(root, target)
	q := foldQueue(r, out.QueueID)
	if q.Current != nil {
		prep := findPreparation(r, *q.Current)
		if old, e := resumeGoalPreview(d, root, r, in, prep.Plan.TaskID); e != nil {
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
		old, e := resumeGoalPreview(d, root, r, in, entry.TaskID)
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
	if old, _, e := storedGoalExecution(d, root, r, in, ""); e != nil {
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
	if e = precheckParentAndResources(d, repo, currentEpicBinding(r, epic, binding), "refs/heads/"+out.Branch, out.WorktreePath); e != nil {
		return out, e
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
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		var e error
		out, e = buildGoalExecutionPreview(d, root, s.Registry, s.Projects, in)
		return e
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

func resumeGoalPreview(d Dependencies, root string, r WorkItemRegistry, in TaskGoalExecuteInput, task TaskID) (*TaskGoalExecutePreview, error) {
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
		old, _, e := storedGoalExecution(d, root, r, in, contentString(m, "execution_request_sha256"))
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
	if p.Input.SpecID != "" && p.Input.SpecID != p.Goal.Goal.SpecID {
		return fmt.Errorf("goal selector differs from request")
	}
	if e := validateQueuePlanObservation(p.ObservedParent, p.Target, p.ParentOID, p.ParentTree); e != nil {
		return e
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
	text += "\n## Implementation ownership\n\nThe implementing owner investigates the code, chooses and revises the detailed design, implements meaningful acceptance tests, and owns verification and correction within this goal and actual runtime permissions. The planner has not chosen these technical details.\n\nThe declared acceptance entrypoint is /bin/sh " + p.AcceptancePath + " with the Task worktree as cwd. The owner must implement this script and meaningful tests, execute it, and preserve the actual script and output as artifacts. A declaration is not an executed test or a passing result.\n\nPreserve work and evidence on interruption. Do not reset, remove worktrees, expand scope, or claim human judgment. Human product judgment and local integration are later distinct facts.\n"
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
		"implementation_basis": contentObject(map[string]canonicaljson.Value{"project_id": string(p.Target.ProjectID), "repo_id": string(p.Target.RepoID), "git_common_dir": p.Target.GitCommonDir, "epic_id": string(p.Target.EpicID), "parent_worktree_id": string(p.Target.ParentWorktreeID), "parent_ref": p.Target.ParentRef, "parent_oid": p.ParentOID, "parent_tree": p.ParentTree, "start_oid": p.ParentOID, "start_tree": p.ParentTree}), "dependencies": []canonicaljson.Value{},
		"goal_origin": contentRefValue(p.Goal.Goal), "execution_request_sha256": p.Confirmation, "execution_request": contentRefValue(frozen), "acceptance_path": p.AcceptancePath,
	} {
		m[k] = v
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
	var p TaskGoalExecutePreview
	var selected TaskDecisionRef
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			old, _, e := storedGoalExecution(d, root, r, in, confirmation)
			if e != nil {
				return e
			}
			if old != nil {
				p = *old
			} else {
				p, e = buildGoalExecutionPreview(d, root, r, projects, in)
				if e != nil {
					return e
				}
				if p.Confirmation != confirmation {
					return queueError("task_goal_execution_conflict", "goal, queue, base or destination changed before execution")
				}
			}
			if prep := preparationForTask(r, p.Goal.TaskID); prep != nil {
				if prep.Plan.ParentOID != p.ParentOID || prep.Plan.WorktreePath != p.WorktreePath {
					return queueError("task_goal_execution_conflict", "Task preparation belongs to another execution")
				}
				selected = prep.Plan.Selection
				return nil
			}
			target, _, e := resolveQueueTarget(d, root, r, projects, in.Target)
			if e != nil {
				return e
			}
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
	prepareInput := TaskPrepareInput{Target: in.Target, Selector: QueueSelector{Kind: "task", TaskID: &id}}
	prepared, e := PrepareTask(d, prepareInput)
	if e != nil {
		return out, e
	}
	if prepared.Plan == nil || prepared.Plan.TaskID != id || prepared.Plan.ParentOID != p.ParentOID || prepared.Plan.ParentTree != p.ParentTree || prepared.Plan.WorktreePath != p.WorktreePath || prepared.Plan.Selection != selected {
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
		out = TaskGoalPreparedExecution{Goal: p.Goal.Goal, Preparation: prepared, Basis: TaskSpecBasis{TaskID: id, Problem: plan.Problem, SpecID: plan.SpecID, Spec: plan.Spec, Assessment: plan.Assessment, Selection: plan.Selection, TaskWorktreeID: task.Worktree.ID, DeliveryBindingSHA256: digest, Dependencies: []string{}}, Executor: p.Goal.Executor, AcceptancePath: p.AcceptancePath}
		return nil
	})
	return out, e
}
