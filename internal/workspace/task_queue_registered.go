package workspace

import (
	"fmt"
	"sort"
)

// RegisteredTaskQueues is a request-local view of registered queue facts. A
// ready pending goal still requires execute preflight; this view observes no
// live checkout, process, provider, or repository state.
type RegisteredTaskQueues struct {
	Queues      []RegisteredTaskQueue
	Freshness   string
	Diagnostics []RegisteredQueueDiagnostic
}

type RegisteredTaskQueue struct {
	Target                 QueueTarget
	QueueID                string
	Revision               int
	Current                *RegisteredQueueCurrent
	Pending                []RegisteredQueuePending
	UnresolvedOperationIDs []string
	Reasons                []QueueReason
	Freshness              string
	SinceUTC               *string
	EvidenceIDs            []string
}

type RegisteredQueuePending struct {
	QueuePending
	Spec        *TaskRevisionRef
	Freshness   string
	SinceUTC    *string
	EvidenceIDs []string
}

type RegisteredQueueCurrent struct {
	QueueCurrent
	Goal        *TaskGoalRef
	Selection   *TaskDecisionRef
	SpecID      *string
	Spec        *TaskRevisionRef
	Freshness   string
	SinceUTC    *string
	EvidenceIDs []string
}

// Diagnostics retain safely observed source identities even when no Task can
// be bound. They must remain visible independently of Task lifecycle filters.
type RegisteredQueueDiagnostic struct {
	Code        string
	Message     string
	Source      string
	QueueID     string
	Target      *QueueTarget
	TaskID      *TaskID
	Freshness   string
	SinceUTC    *string
	EvidenceIDs []string
}

type registeredQueueIndex struct {
	tasks        map[TaskID]TaskRecord
	preparations map[string]TaskPreparation
	preparation  map[TaskID]TaskPreparation
	operations   map[TaskID]WorktreeOperationRecord
	events       map[string][]TaskQueueEvent
	folds        map[string]queueFold
	effects      map[string][]string
}

// ReadRegisteredTaskQueues reuses the caller's validated registration snapshot.
// Only referenced queue content is read, once per unique managed source plus a
// final integrity check. The caller remains responsible for checking drift in
// the supplied registration snapshot (CheckReadViewBasis).
func ReadRegisteredTaskQueues(d Dependencies, ws WorkspaceObservation, projects ProjectSnapshot, r WorkItemRegistry) RegisteredTaskQueues {
	out := RegisteredTaskQueues{Queues: []RegisteredTaskQueue{}, Freshness: "fresh", Diagnostics: []RegisteredQueueDiagnostic{}}
	index := indexRegisteredQueues(ws.Root, r)
	content := captureTaskContent(d.TaskContent)
	d.TaskContent = content.store
	knownQueues := map[string]bool{}
	for _, epic := range r.Epics {
		for _, repo := range epic.RepoBindings {
			target := QueueTarget{ProjectID: epic.ProjectID, RepoID: repo.RepoID, EpicID: epic.ID, GitCommonDir: repo.GitCommonDir, ParentWorktreeID: repo.Worktree.ID, ParentLocator: repo.Worktree.Locator, ParentRef: repo.Worktree.Ref}
			id := queueID(ws.Root, target)
			knownQueues[id] = true
			events := index.events[id]
			if len(events) == 0 && len(index.effects[id]) == 0 {
				continue
			}
			fold := index.folds[id]
			q := RegisteredTaskQueue{Target: target, QueueID: id, Revision: fold.Revision, Pending: []RegisteredQueuePending{}, UnresolvedOperationIDs: append([]string{}, index.effects[id]...), Reasons: []QueueReason{}, Freshness: "fresh", EvidenceIDs: []string{id}}
			var pendingSince *string
			pendingEvidence := []string{id}
			validTarget := true
			for _, event := range events {
				q.SinceUTC = recordedLatestTime(q.SinceUTC, event.RecordedAtUTC)
				q.EvidenceIDs = append(q.EvidenceIDs, event.RequestSHA256)
				if event.Kind == "set" {
					draft := decodeQueueDraft(queueRequestValue(event.Request))
					if draft.ProjectID != target.ProjectID || draft.RepoID != target.RepoID || draft.EpicID != target.EpicID {
						validTarget = false
					}
					pendingSince = recordedLatestTime(nil, event.RecordedAtUTC)
					pendingEvidence = []string{id, event.RequestSHA256}
				}
			}
			if _, registeredRepo, err := projectAndRepo(projects, target.ProjectID, target.RepoID); err != nil || registeredRepo.GitCommonDir != target.GitCommonDir {
				validTarget = false
			}
			if !validTarget {
				q.Freshness = "unknown"
				q.Reasons = queueReasons("task_queue_target_conflict", "queue target differs from the registered Project, repository, or Epic")
				out.addDiagnostic(q, nil, q.Reasons[0], "unknown", q.EvidenceIDs)
				out.Queues = append(out.Queues, q)
				continue
			}
			if fold.Current != nil {
				q.Reasons = append(q.Reasons, QueueReason{"task_queue_current", "the target has a current preparation; resolve it before starting another goal"})
				prep, exists := index.preparations[*fold.Current]
				if !exists || prep.Plan.Target != target || prep.Plan.Workspace.Root != ws.Root {
					q.Freshness = "unknown"
					reason := QueueReason{"task_queue_preparation_unknown", "current preparation is missing or belongs to another registered target"}
					q.Reasons = append(q.Reasons, reason)
					out.addDiagnostic(q, nil, reason, "unknown", []string{id, *fold.Current})
				} else if task, exists := index.tasks[prep.Plan.TaskID]; !exists || !registeredQueueOwns(target, task) {
					q.Freshness = "unknown"
					reason := QueueReason{"task_queue_task_missing", "current preparation has no safe registered Task binding"}
					q.Reasons = append(q.Reasons, reason)
					out.addDiagnostic(q, nil, reason, "unknown", []string{id, prep.ID})
				} else {
					content.owner = registeredQueueOwner(id, task.ID)
					current, err := registeredCurrentPreparation(d, ws.Root, r, prep)
					q.Current = &current
					if err != nil {
						q.Freshness = "unknown"
						for _, reason := range queueErrorReasons(err) {
							out.addDiagnostic(q, &task.ID, reason, "unknown", current.EvidenceIDs)
						}
					}
				}
			}
			if len(q.UnresolvedOperationIDs) != 0 {
				q.Reasons = append(q.Reasons, QueueReason{"task_queue_effect_pending", "the target has unresolved registered effects"})
				q.EvidenceIDs = append(q.EvidenceIDs, q.UnresolvedOperationIDs...)
				out.addDiagnostic(q, nil, q.Reasons[len(q.Reasons)-1], "fresh", q.UnresolvedOperationIDs)
			}
			for rank, entry := range fold.Pending {
				task, exists := index.tasks[entry.TaskID]
				if !exists || !registeredQueueOwns(target, task) {
					q.Freshness = "unknown"
					reason := QueueReason{"task_queue_task_missing", "pending queue entry has no safe registered Task binding"}
					out.addDiagnostic(q, nil, reason, "unknown", append(append([]string{}, pendingEvidence...), string(entry.TaskID)))
					continue
				}
				content.owner = registeredQueueOwner(id, task.ID)
				pending, err := registeredQueuePending(d, ws.Root, r, index, target, task, entry, rank+1)
				pending.SinceUTC = pendingSince
				pending.EvidenceIDs = registeredQueueEvidence(append(pending.EvidenceIDs, pendingEvidence...))
				if err != nil {
					for _, reason := range queueErrorReasons(err) {
						out.addDiagnostic(q, &task.ID, reason, pending.Freshness, pending.EvidenceIDs)
					}
				}
				if len(q.Reasons) != 0 {
					pending.Reasons = append(pending.Reasons, q.Reasons...)
					pending.EvidenceIDs = registeredQueueEvidence(append(pending.EvidenceIDs, q.EvidenceIDs...))
					if pending.State == "ready" || pending.State == "unknown" && pending.Freshness == "fresh" {
						pending.State = "blocked"
					}
				}
				if q.Freshness != "fresh" && pending.State == "ready" {
					pending.State, pending.Freshness = "unknown", "unknown"
				}
				pending.Reasons = sortQueueReasons(pending.Reasons)
				q.Pending = append(q.Pending, pending)
			}
			q.EvidenceIDs = registeredQueueEvidence(q.EvidenceIDs)
			q.Reasons = sortQueueReasons(q.Reasons)
			out.Queues = append(out.Queues, q)
		}
	}
	for id, events := range index.events {
		if knownQueues[id] {
			continue
		}
		// A syntactically valid registry can still carry a queue ID hashed for
		// another physical workspace. Do not invent a Task or target for it.
		diagnostic := RegisteredQueueDiagnostic{Code: "task_queue_orphaned", Message: "queue identity does not match a registered target in this workspace", Source: "registered_queue", QueueID: id, Freshness: "unknown", EvidenceIDs: []string{id}}
		for _, event := range events {
			diagnostic.EvidenceIDs = append(diagnostic.EvidenceIDs, event.RequestSHA256)
			diagnostic.SinceUTC = recordedLatestTime(diagnostic.SinceUTC, event.RecordedAtUTC)
		}
		diagnostic.EvidenceIDs = registeredQueueEvidence(diagnostic.EvidenceIDs)
		out.Diagnostics = append(out.Diagnostics, diagnostic)
	}
	out.applyContentChanges(content.changedOwners())
	out.finish()
	return out
}

func registeredQueueOwns(target QueueTarget, task TaskRecord) bool {
	return task.ProjectID == target.ProjectID && task.RepoID == target.RepoID && task.ParentEpicID == target.EpicID && task.GitCommonDir == target.GitCommonDir
}

func registeredQueueOwner(id string, task TaskID) string { return id + "/" + string(task) }

func indexRegisteredQueues(root string, r WorkItemRegistry) registeredQueueIndex {
	x := registeredQueueIndex{tasks: map[TaskID]TaskRecord{}, preparations: map[string]TaskPreparation{}, preparation: map[TaskID]TaskPreparation{}, operations: map[TaskID]WorktreeOperationRecord{}, events: map[string][]TaskQueueEvent{}, folds: map[string]queueFold{}, effects: map[string][]string{}}
	for _, task := range r.Tasks {
		x.tasks[task.ID] = task
	}
	for _, prep := range r.TaskPreparations {
		x.preparations[prep.ID], x.preparation[prep.Plan.TaskID] = prep, prep
	}
	for _, event := range r.TaskQueueEvents {
		x.events[event.QueueID] = append(x.events[event.QueueID], event)
		fold, exists := x.folds[event.QueueID]
		if !exists {
			fold = emptyQueueFold()
		}
		applyQueueEvent(&fold, event)
		if event.Kind == "reserve" && fold.Current != nil {
			if prep, exists := x.preparations[*fold.Current]; exists {
				pending := []TaskGoalQueueEntry{}
				for _, entry := range fold.Pending {
					if entry.TaskID != prep.Plan.TaskID {
						pending = append(pending, entry)
					}
				}
				fold.Pending = pending
			}
		}
		x.folds[event.QueueID] = fold
	}
	for _, op := range r.WorktreeOperations {
		x.operations[op.TaskID] = op
		if op.State == "ready" {
			continue
		}
		if prep, exists := x.preparation[op.TaskID]; exists && x.folds[prep.Plan.QueueID].Terminal[prep.ID] == "released" {
			continue
		}
		id := queueID(root, QueueTarget{ProjectID: op.ProjectID, RepoID: op.RepoID, EpicID: op.ParentEpicID})
		x.effects[id] = append(x.effects[id], string(op.ID))
	}
	complete := map[IntegrationAuthorityID]bool{}
	for _, result := range r.IntegrationResults {
		switch result.Outcome {
		case "exact_effect", "already_integrated", "no_effect", "conflict":
			complete[result.AuthorityID] = true
		}
	}
	for _, authority := range r.IntegrationAuthorities {
		if !complete[authority.ID] {
			id := queueID(root, QueueTarget{ProjectID: authority.Plan.Project.ProjectID, RepoID: authority.Plan.Repository.RepoID, EpicID: authority.Plan.Epic.EpicID})
			x.effects[id] = append(x.effects[id], string(authority.ID))
		}
	}
	for _, update := range r.EpicBaseUpdates {
		if update.Phase == "intent" {
			id := queueID(root, update.Plan.Target)
			x.effects[id] = append(x.effects[id], update.ID)
		}
	}
	for id, effects := range x.effects {
		x.effects[id] = registeredQueueEvidence(effects)
	}
	return x
}

func registeredQueuePending(d Dependencies, root string, r WorkItemRegistry, index registeredQueueIndex, target QueueTarget, task TaskRecord, entry TaskGoalQueueEntry, rank int) (RegisteredQueuePending, error) {
	row := RegisteredQueuePending{QueuePending: QueuePending{TaskID: task.ID, Title: task.Title, Goal: entry.Goal, Selection: entry.Selection, Rank: rank, State: "ready", Reasons: []QueueReason{}}, Freshness: "fresh", EvidenceIDs: []string{}}
	_, prepared := index.preparation[task.ID]
	_, operated := index.operations[task.ID]
	if task.Worktree != nil || task.WorktreeState != WorkItemUnbound || prepared || operated {
		row.State = "blocked"
		row.Reasons = append(row.Reasons, QueueReason{"task_queue_task_bound", "Task already owns a worktree, preparation, or preserved intent"})
	}
	if entry.Goal != nil {
		row.SpecID, row.SpecRevision, row.Spec = &entry.Goal.SpecID, &entry.Goal.Spec.Revision, &entry.Goal.Spec
		row.EvidenceIDs = append(row.EvidenceIDs, entry.Goal.Spec.ManifestSHA256)
		goal, err := loadTaskGoal(d, root, r, task.ID, *entry.Goal)
		if err != nil {
			return unknownRegisteredPending(row, err)
		}
		row.Title, row.Executor = goal.Title, &goal.Executor
		if !contentTypedEqual(taskContentState(r, task.ID).ProblemHead, &goal.Problem) {
			row.State = "blocked"
			row.Reasons = append(row.Reasons, QueueReason{"task_goal_problem_changed", "queued goal no longer refers to the current registered Problem"})
		}
	} else if entry.Selection == nil {
		row.State = "blocked"
		row.Reasons = append(row.Reasons, QueueReason{"task_queue_selection_required", "queue entry does not bind an exact goal or selected solution"})
	} else {
		row.EvidenceIDs = append(row.EvidenceIDs, entry.Selection.ID, entry.Selection.ManifestSHA256)
		if err := registeredQueueSelection(d, root, r, target, task, &row); err != nil {
			return unknownRegisteredPending(row, err)
		}
	}
	return row, nil
}

func unknownRegisteredPending(row RegisteredQueuePending, err error) (RegisteredQueuePending, error) {
	row.State, row.Freshness = "unknown", "unknown"
	row.Reasons = append(row.Reasons, queueErrorReasons(err)...)
	return row, err
}

func registeredQueueSelection(d Dependencies, root string, r WorkItemRegistry, target QueueTarget, task TaskRecord, row *RegisteredQueuePending) error {
	if err := validateQueueSelection(r, task.ID, row.Selection); err != nil {
		return err
	}
	selection, err := readRegisteredTaskManifest(d.TaskContent, root, r, row.Selection.ManifestSHA256)
	if err != nil {
		return err
	}
	m := contentFields(selection)
	solution := contentFields(m["solution"])
	row.SpecID, row.Spec = registeredQueueString(contentString(solution, "spec_id")), valueRevision(solution["spec"])
	if contentString(m, "task_id") != string(task.ID) || contentString(m, "id") != row.Selection.ID || contentString(m, "action") != "select" || row.SpecID == nil || row.Spec == nil {
		return queueError("task_queue_selection_required", "queue selection does not bind an exact selected solution for this Task")
	}
	row.SpecRevision = &row.Spec.Revision
	row.EvidenceIDs = append(row.EvidenceIDs, row.Spec.ManifestSHA256)
	spec, err := readRegisteredTaskManifest(d.TaskContent, root, r, row.Spec.ManifestSHA256)
	if err != nil {
		return err
	}
	if _, err = registeredTaskSpec(d, root, r, task.ID, *row.SpecID, row.Spec); err != nil {
		return err
	}
	if contentString(contentFields(spec), "task_id") != string(task.ID) || contentString(contentFields(spec), "spec_id") != *row.SpecID || contentInt(contentFields(spec), "revision") != row.Spec.Revision {
		return queueError("task_spec_binding_conflict", "selected Spec identity differs from its registered reference")
	}
	if contentString(contentFields(spec), "contract_kind") == "execution" {
		var goal TaskGoalRef
		if err = contentDecode(contentFields(spec)["goal_origin"], &goal); err != nil {
			return err
		}
		if err = validateExecutionGoalOrigin(d, root, r, task.ID, spec); err != nil {
			return err
		}
		row.Goal = &goal
		row.EvidenceIDs = append(row.EvidenceIDs, goal.Spec.ManifestSHA256)
	}
	if err = taskSpecStructurallyReady(spec); err != nil {
		row.State = "blocked"
		row.Reasons = append(row.Reasons, queueErrorReasons(err)...)
	}
	state := taskContentState(r, task.ID)
	assessment := lastTaskAssessment(r, task.ID, *row.SpecID, row.Spec.Revision)
	selectedAssessment := valueDecision(solution["assessment"])
	if !contentTypedEqual(row.Selection, state.SelectionEvent) || !contentTypedEqual(valueRevision(solution["problem"]), state.ProblemHead) || assessment == nil || selectedAssessment == nil || assessment.ID != selectedAssessment.ID || assessment.ManifestSHA256 != selectedAssessment.ManifestSHA256 {
		row.State = "blocked"
		row.Reasons = append(row.Reasons, QueueReason{"task_queue_selection_stale", "registered Problem, selection, or assessment changed after this queue binding"})
	}
	basis := contentFields(contentFields(spec)["implementation_basis"])
	row.ParentOID = registeredQueueString(contentString(basis, "parent_oid"))
	var latest EpicBaseVersion
	for _, base := range r.EpicBaseVersions {
		if base.ProjectID == target.ProjectID && base.RepoID == target.RepoID && base.EpicID == target.EpicID {
			if row.ParentOID != nil && base.OID == *row.ParentOID && base.Tree == contentString(basis, "parent_tree") {
				revision := base.Revision
				row.BaseRevision = &revision
			}
			if base.Revision > latest.Revision {
				latest = base
			}
		}
	}
	if row.BaseRevision == nil || contentString(basis, "project_id") != string(target.ProjectID) || contentString(basis, "repo_id") != string(target.RepoID) || contentString(basis, "epic_id") != string(target.EpicID) || contentString(basis, "git_common_dir") != target.GitCommonDir || contentString(basis, "parent_worktree_id") != string(target.ParentWorktreeID) || contentString(basis, "parent_ref") != target.ParentRef {
		return queueError("task_spec_binding_conflict", "selected implementation basis does not match a preserved registered target")
	}
	if *row.BaseRevision != latest.Revision {
		row.State = "blocked"
		row.Reasons = append(row.Reasons, QueueReason{"epic_base_changed", "selected implementation basis is older than the current registered Epic base"})
	}
	if row.State == "ready" {
		row.State = "unknown"
	}
	row.Reasons = append(row.Reasons, QueueReason{"task_queue_live_preflight_required", "selected solution start readiness requires live preflight, which this registered view does not perform"})
	return nil
}

func registeredCurrentPreparation(d Dependencies, root string, r WorkItemRegistry, prep TaskPreparation) (RegisteredQueueCurrent, error) {
	p := prep.Plan
	row := RegisteredQueueCurrent{QueueCurrent: currentPreparation(d, r, prep), Selection: &p.Selection, SpecID: &p.SpecID, Spec: &p.Spec, Freshness: "fresh", SinceUTC: recordedLatestTime(nil, prep.CreatedAtUTC), EvidenceIDs: registeredQueueEvidence([]string{prep.ID, prep.PlanSHA256, string(prep.OperationID), p.Selection.ID, p.Selection.ManifestSHA256, p.Spec.ManifestSHA256})}
	basis := TaskSpecBasis{TaskID: p.TaskID, Problem: p.Problem, SpecID: p.SpecID, Spec: p.Spec, Assessment: p.Assessment, Selection: p.Selection, Dependencies: []string{}}
	eval, err := loadTaskSpecContent(d, root, r, basis)
	if err == nil {
		implementation := contentFields(contentFields(eval.Spec)["implementation_basis"])
		if contentString(implementation, "parent_oid") != p.ParentOID || contentString(implementation, "parent_tree") != p.ParentTree || contentString(implementation, "start_oid") != p.ParentOID || contentString(implementation, "start_tree") != p.ParentTree {
			err = queueError("task_queue_integrity_conflict", "preparation plan differs from the selected registered implementation basis")
		}
		if err == nil && contentString(contentFields(eval.Spec), "contract_kind") == "execution" {
			var goal TaskGoalRef
			err = contentDecode(contentFields(eval.Spec)["goal_origin"], &goal)
			if err == nil {
				row.Goal = &goal
				row.EvidenceIDs = registeredQueueEvidence(append(row.EvidenceIDs, goal.Spec.ManifestSHA256))
				err = validateExecutionGoalOrigin(d, root, r, p.TaskID, eval.Spec)
			}
		}
		if err == nil && !contentTypedEqual(eval.RequiredInputs, p.RequiredInputs) {
			err = queueError("task_queue_integrity_conflict", "preserved preparation inputs differ from their selected content")
		}
		if !contentTypedEqual(taskContentState(r, p.TaskID).ProblemHead, &p.Problem) {
			row.Reasons = append(row.Reasons, QueueReason{"task_goal_problem_changed", "the current registered Problem differs from this preparation's preserved basis"})
		}
	}
	if err != nil {
		row.State, row.Freshness = "unknown", "unknown"
		row.Reasons = append(row.Reasons, queueErrorReasons(err)...)
	}
	return row, err
}

func registeredQueueString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func registeredQueueEvidence(ids []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, id := range ids {
		if id != "" && !seen[id] {
			out = append(out, id)
			seen[id] = true
		}
	}
	sort.Strings(out)
	return out
}

func (out *RegisteredTaskQueues) addDiagnostic(q RegisteredTaskQueue, task *TaskID, reason QueueReason, freshness string, evidence []string) {
	out.Diagnostics = append(out.Diagnostics, RegisteredQueueDiagnostic{Code: reason.Code, Message: reason.Message, Source: "registered_queue", QueueID: q.QueueID, Target: &q.Target, TaskID: task, Freshness: freshness, EvidenceIDs: registeredQueueEvidence(evidence)})
}

func (out *RegisteredTaskQueues) applyContentChanges(changed map[string][]string) {
	reason := QueueReason{"task_queue_source_changed", "referenced queue content changed or became unavailable during this read"}
	for i := range out.Queues {
		q := &out.Queues[i]
		if q.Current != nil {
			row := q.Current
			if evidence := changed[registeredQueueOwner(q.QueueID, row.TaskID)]; len(evidence) != 0 {
				row.State, row.Freshness, q.Freshness = "unknown", "stale", "stale"
				row.Reasons = append(row.Reasons, reason)
				out.addDiagnostic(*q, &row.TaskID, reason, "stale", evidence)
			}
		}
		for j := range q.Pending {
			row := &q.Pending[j]
			if evidence := changed[registeredQueueOwner(q.QueueID, row.TaskID)]; len(evidence) != 0 {
				row.State, row.Freshness, q.Freshness = "unknown", "stale", "stale"
				row.Reasons = append(row.Reasons, reason)
				out.addDiagnostic(*q, &row.TaskID, reason, "stale", evidence)
			}
		}
	}
}

func (out *RegisteredTaskQueues) finish() {
	for _, q := range out.Queues {
		if q.Freshness != "fresh" {
			out.Freshness = "unknown"
		}
	}
	for _, diagnostic := range out.Diagnostics {
		if diagnostic.Freshness != "fresh" {
			out.Freshness = "unknown"
		}
	}
	sort.Slice(out.Queues, func(i, j int) bool {
		a, b := out.Queues[i].Target, out.Queues[j].Target
		if a.ProjectID != b.ProjectID {
			return a.ProjectID < b.ProjectID
		}
		if a.RepoID != b.RepoID {
			return a.RepoID < b.RepoID
		}
		return a.EpicID < b.EpicID
	})
	sort.Slice(out.Diagnostics, func(i, j int) bool {
		a, b := out.Diagnostics[i], out.Diagnostics[j]
		return fmt.Sprint(a.QueueID, "/", a.Code, "/", registeredQueueTaskID(a.TaskID), "/", a.Message) < fmt.Sprint(b.QueueID, "/", b.Code, "/", registeredQueueTaskID(b.TaskID), "/", b.Message)
	})
}

func registeredQueueTaskID(id *TaskID) string {
	if id == nil {
		return ""
	}
	return string(*id)
}
