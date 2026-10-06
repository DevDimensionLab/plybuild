package workspaceview

import (
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ReadWorkflowStatus composes one registered basis and the existing process
// batch. It performs no live preflight, creates no state, and starts no process.
func ReadWorkflowStatus(d workspace.Dependencies, o WorkflowStatusOptions) (WorkflowStatus, error) {
	if err := o.Filters.Validate(); err != nil {
		return WorkflowStatus{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return WorkflowStatus{}, err
	}
	if err := ValidateTaskFilters(s, o.Filters); err != nil {
		return WorkflowStatus{}, err
	}
	process := readProgressProcess(d, s)
	var runs taskrun.Inventory
	if s.process != nil {
		runs = s.process.runs
	} else {
		// A failed journal read must not hide independently preserved runs.
		runs, err = taskrun.ReadInventory(taskrun.SystemDependencies(d), s.Workspace.Root)
		if err != nil {
			runs = taskrun.Inventory{Runs: []taskrun.InventoryRun{}, Reasons: []taskrun.Reason{{Code: "run_store_unavailable", Detail: err.Error()}}}
		}
	}
	queues := workspace.ReadRegisteredTaskQueues(d, s.Workspace, s.Projects, s.Registry)
	s.RefreshFreshness(d)
	return BuildWorkflowStatus(s, o, process, runs, queues)
}

// BuildWorkflowStatus is a pure projection. The native Progress remains intact;
// presentation actions supplement it without manufacturing result/QA facts.
func BuildWorkflowStatus(s *Snapshot, o WorkflowStatusOptions, process map[workspace.TaskID]ProcessFacts, runs taskrun.Inventory, queues workspace.RegisteredTaskQueues) (WorkflowStatus, error) {
	tasks, err := BuildTaskList(s, o.Filters, process)
	if err != nil {
		return WorkflowStatus{}, err
	}
	out := WorkflowStatus{Kind: WorkflowStatusKind, SchemaVersion: 1, Workspace: tasks.Workspace, Filters: o.Filters, IncludeAll: o.IncludeAll, ObservedAtUTC: s.ObservedAtUTC, SourceBasis: "registered", Freshness: s.Freshness, Projects: []WorkflowStatusProject{}, Epics: []EpicRow{}, Items: []WorkflowStatusItem{}, Diagnostics: []WorkflowDiagnostic{}, Counts: WorkflowStatusCounts{Categories: map[string]int{}}}
	for _, category := range workflowCategories {
		out.Counts.Categories[category] = 0
	}
	known := map[string]workspace.TaskRecord{}
	for _, task := range s.Registry.Tasks {
		known[string(task.ID)] = task
	}
	rows := make(map[workspace.TaskID]*WorkflowStatusItem, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		row := &WorkflowStatusItem{TaskID: task.TaskID, ProjectID: task.ProjectID, RepoID: task.RepoID, EpicID: task.ParentEpicID, TaskListTitle: task.TaskListTitle, WorktreeState: task.WorktreeState, TaskLifecycle: task.Lifecycle, EpicLifecycle: task.EpicLifecycle, Progress: task.Progress, LastActivityUTC: task.Progress.LastActivityUTC, Freshness: task.Progress.Freshness, NextActions: []WorkflowStatusAction{}, Runs: []taskrun.InventoryRun{}, Queues: []WorkflowStatusQueue{}}
		for _, action := range task.Progress.NextActions {
			row.NextActions = append(row.NextActions, WorkflowStatusAction{Action: action, Current: true})
		}
		rows[task.TaskID] = row
		if task.Status != "available" {
			out.addDiagnostic("task_title_unavailable", "content", "The current registered Task title could not be read.", string(task.TaskID), "", "", nil)
			row.Freshness = combineFreshness(row.Freshness, "unknown")
		}
		for _, reason := range task.Progress.Reasons {
			out.addDiagnostic(reason, "progress", "Registered progress: "+reason, string(task.TaskID), "", "", nil)
		}
	}
	for _, reason := range s.Reasons {
		out.addDiagnostic(reason, "workspace", "Workspace source observation: "+reason, "", "", "", nil)
	}
	for _, reason := range runs.Reasons {
		out.addDiagnostic(reason.Code, "run", reason.Detail, "", "", "", nil)
		out.Freshness = combineFreshness(out.Freshness, "unknown")
	}
	basis := workflowBasisIndex(s.Registry)
	for _, run := range runs.Runs {
		var taskID string
		if run.TaskID != nil {
			taskID = *run.TaskID
		}
		if task, exists := known[taskID]; exists && !MatchesTaskFilters(task, o.Filters) {
			continue
		}
		row := rows[workspace.TaskID(taskID)]
		if row == nil {
			out.addDiagnostic("run_task_unbound", "run", "The preserved run cannot be bound to a registered Task in this workspace.", taskID, run.RunID, "", workflowRunEvidence(run))
			out.Freshness = combineFreshness(out.Freshness, "unknown")
		} else {
			row.Runs = append(row.Runs, run)
			row.Progress.HasProgress = true
			if row.Progress.State == "nothing" {
				row.Progress.State = "in_progress"
			}
			row.Freshness = combineFreshness(row.Freshness, run.Freshness)
			if run.Freshness != "fresh" {
				row.addAction("run", "inspect_run_evidence", "unknown", "Inspect the unavailable or changing preserved run evidence.", "attention", nil, workflowRunEvidence(run), run.RunID, "")
			}
			if delivery := run.Delivery; delivery != nil {
				row.Freshness = combineFreshness(row.Freshness, delivery.Freshness)
				if delivery.Freshness != "fresh" {
					row.addAction("delivery_report", "inspect_delivery_evidence", "unknown", "Inspect the incomplete or changing delivery report evidence.", "attention", nil, workflowRunEvidence(run), run.RunID, "")
				}
				currentBasis := basis.matches(workspace.TaskID(taskID), delivery.Basis)
				if !currentBasis {
					out.addDiagnostic("delivery_task_basis_stale", "delivery_report", "The delivery run does not cover the current registered Task basis.", taskID, run.RunID, "", workflowRunEvidence(run))
					row.addAction("delivery_report", "inspect_delivery_basis", "agent", "Inspect the historical delivery and the current registered Task basis.", "attention", nil, workflowRunEvidence(run), run.RunID, "")
				}
				a := delivery.NextAction
				if delivery.Phase == "needs_input" || delivery.Phase == "stopped" || delivery.Phase == "awaiting_acceptance" || delivery.Phase == "unknown" {
					if a.Kind != "" {
						row.addAction("delivery_report", a.Kind, a.Actor, a.Reason, "next_step", a.SinceUTC, a.EvidenceIDs, run.RunID, "")
						row.NextActions[len(row.NextActions)-1].Current = currentBasis && delivery.Freshness == "fresh"
					}
				}
				// All old questions stay available as history, including reports
				// superseded by verification or an explicit newer progress event.
				for _, report := range delivery.Reports {
					if report.Question == nil || report.Current && delivery.Phase == "needs_input" {
						continue
					}
					row.addAction("delivery_report", "answer_question", "human", *report.Question, "next_step", report.RecordedAtUTC, []string{run.RunID, report.EventID, report.Binding.SHA256}, run.RunID, "")
					row.NextActions[len(row.NextActions)-1].Current = false
				}
			} else if run.State == "stopped" {
				row.addAction("run", "follow_stopped_run", "agent", "Inspect the stopped run and its preserved outcome.", "next_step", run.LastSeenUTC, workflowRunEvidence(run), run.RunID, "")
			}
		}
		for _, reason := range run.Reasons {
			out.addDiagnostic(reason.Code, "run", reason.Detail, taskID, run.RunID, "", workflowRunEvidence(run))
		}
		if run.Delivery != nil {
			for _, reason := range run.Delivery.Reasons {
				out.addDiagnostic(reason.Code, "delivery_report", reason.Detail, taskID, run.RunID, "", workflowRunEvidence(run))
			}
		}
		out.Freshness = combineFreshness(out.Freshness, run.Freshness)
	}
	out.addQueues(rows, queues, o.Filters)
	for _, row := range rows {
		if s.Freshness != "fresh" && (row.Progress.HasProgress || len(row.Queues) > 0) {
			for i := range row.NextActions {
				row.NextActions[i].Current = false
			}
			row.addAction("readiness", "inspect_registered_sources", "agent", "Read the changing or unavailable registered sources before acting on this snapshot.", "attention", nil, nil, "", "")
		}
		workflowClassify(row)
		sortWorkflowActions(row.NextActions)
		sort.Slice(row.Runs, func(i, j int) bool { return row.Runs[i].RunID < row.Runs[j].RunID })
		sort.Slice(row.Queues, func(i, j int) bool {
			if row.Queues[i].QueueID != row.Queues[j].QueueID {
				return row.Queues[i].QueueID < row.Queues[j].QueueID
			}
			return row.Queues[i].State < row.Queues[j].State
		})
		out.Counts.Total++
		out.Freshness = combineFreshness(out.Freshness, row.Freshness)
		if !o.IncludeAll {
			switch row.Category {
			case "inactive":
				out.Counts.Hidden.Inactive++
				continue
			case "completed":
				out.Counts.Hidden.Completed++
				continue
			case "backlog":
				out.Counts.Hidden.Backlog++
				continue
			}
		}
		out.Items = append(out.Items, *row)
		out.Counts.Visible++
		out.Counts.Categories[row.Category]++
	}
	out.fillContainers(s, tasks, o.Filters)
	sort.Slice(out.Items, func(i, j int) bool { return workflowItemLess(out.Items[i], out.Items[j]) })
	sort.Slice(out.Diagnostics, func(i, j int) bool {
		a, b := out.Diagnostics[i], out.Diagnostics[j]
		return strings.Join([]string{a.Source, valueOrEmpty(a.TaskID), valueOrEmpty(a.RunID), valueOrEmpty(a.QueueID), a.Code, a.Reason}, "\x00") < strings.Join([]string{b.Source, valueOrEmpty(b.TaskID), valueOrEmpty(b.RunID), valueOrEmpty(b.QueueID), b.Code, b.Reason}, "\x00")
	})
	return out, nil
}

var workflowCategories = []string{"needs_you", "follow_up", "in_progress", "ready_next", "inactive", "completed", "backlog"}

func (out *WorkflowStatus) addQueues(rows map[workspace.TaskID]*WorkflowStatusItem, queues workspace.RegisteredTaskQueues, filters workspace.TaskListFilters) {
	for _, q := range queues.Queues {
		if !workflowTargetMatches(q.Target, filters) {
			continue
		}
		for _, pending := range q.Pending {
			row := rows[pending.TaskID]
			if row == nil {
				continue
			}
			entry := WorkflowStatusQueue{QueueID: q.QueueID, Revision: q.Revision, Target: q.Target, State: "pending", Rank: &pending.Rank, SpecID: pending.SpecID, Spec: pending.Spec, Readiness: pending.State, Freshness: combineFreshness(q.Freshness, pending.Freshness), SinceUTC: pending.SinceUTC, Reasons: append([]workspace.QueueReason{}, pending.Reasons...), EvidenceIDs: sortedUnique(append(append([]string{}, q.EvidenceIDs...), pending.EvidenceIDs...))}
			entry.Goal, entry.Selection = pending.Goal, pending.Selection
			if q.Current != nil {
				current := q.Current.QueueCurrent
				entry.Current = &current
			}
			entry.Freshness = combineFreshness(entry.Freshness, row.Freshness)
			if entry.Freshness != "fresh" && entry.Readiness == "ready" {
				entry.Readiness = "unknown"
				entry.Reasons = append(entry.Reasons, workspace.QueueReason{Code: "queue_source_not_fresh", Message: "Registered queue readiness is unknown because its captured source basis is not fresh."})
			}
			row.Queues = append(row.Queues, entry)
			row.Freshness = combineFreshness(row.Freshness, entry.Freshness)
			row.LastActivityUTC = latestTime(row.LastActivityUTC, entry.SinceUTC)
			if entry.Readiness != "ready" {
				reason := "Inspect the registered queue prerequisites; start readiness is unknown."
				if len(entry.Reasons) > 0 {
					reason = entry.Reasons[0].Message
				}
				row.addAction("queue", "inspect_queue", "agent", reason, "next_step", entry.SinceUTC, entry.EvidenceIDs, "", q.QueueID)
			}
		}
		if current := q.Current; current != nil {
			if row := rows[current.TaskID]; row != nil {
				c := current.QueueCurrent
				entry := WorkflowStatusQueue{QueueID: q.QueueID, Revision: q.Revision, Target: q.Target, State: "current", SpecID: current.SpecID, Spec: current.Spec, Readiness: "blocked", Current: &c, Freshness: current.Freshness, SinceUTC: current.SinceUTC, Reasons: append([]workspace.QueueReason{}, current.Reasons...), EvidenceIDs: sortedUnique(append(append([]string{}, q.EvidenceIDs...), current.EvidenceIDs...))}
				entry.Goal, entry.Selection = current.Goal, current.Selection
				row.Queues = append(row.Queues, entry)
				row.Progress.HasProgress = true
				row.Freshness = combineFreshness(row.Freshness, current.Freshness)
				row.LastActivityUTC = latestTime(row.LastActivityUTC, current.SinceUTC)
				// A reservation with an executing run is expected. An unhandled
				// reservation or unresolved effect needs owner inspection.
				if len(row.Runs) == 0 || len(q.UnresolvedOperationIDs) > 0 || len(current.Reasons) > 0 || current.State == "blocked" || current.State == "unknown" || current.Freshness != "fresh" {
					row.addAction("queue", "follow_current_preparation", "agent", "Follow the current preparation and its recorded effects before selecting another goal.", "next_step", current.SinceUTC, entry.EvidenceIDs, "", q.QueueID)
				}
			}
		}
	}
	for _, d := range queues.Diagnostics {
		if d.Target != nil && !workflowTargetMatches(*d.Target, filters) {
			continue
		}
		var task string
		if d.TaskID != nil {
			task = string(*d.TaskID)
		}
		out.addDiagnostic(d.Code, "queue", d.Message, task, "", d.QueueID, d.EvidenceIDs)
		out.Diagnostics[len(out.Diagnostics)-1].SinceUTC = d.SinceUTC
		out.Freshness = combineFreshness(out.Freshness, d.Freshness)
	}
}

func workflowClassify(row *WorkflowStatusItem) {
	switch {
	case row.TaskLifecycle != workspace.LifecycleActive || row.EpicLifecycle != workspace.LifecycleActive:
		row.Category = "inactive"
	case row.Progress.State == "integrated":
		row.Category = "completed"
	default:
		row.Category = "backlog"
		if row.Progress.HasProgress {
			row.Category = "in_progress"
		}
		for _, q := range row.Queues {
			if q.State == "pending" && q.Readiness == "ready" && row.Category == "backlog" {
				row.Category = "ready_next"
			}
		}
		for _, a := range row.NextActions {
			if !a.Current {
				continue
			}
			if a.Actor == "human" && row.Progress.HasProgress {
				row.Category = "needs_you"
				return
			}
			row.Category = "follow_up"
		}
		if row.Category == "in_progress" && row.Freshness != "fresh" {
			row.Category = "follow_up"
		}
		return
	}
	for i := range row.NextActions {
		row.NextActions[i].Current = false
	}
}

func (out *WorkflowStatus) fillContainers(s *Snapshot, tasks TaskList, filters workspace.TaskListFilters) {
	epics := epicsFromTasks(s, EpicFilters{ProjectID: filters.ProjectID}, tasks)
	projectEpics := map[workspace.ProjectID]bool{}
	for _, epic := range epics.Epics {
		if filters.EpicID != nil && epic.EpicID != *filters.EpicID {
			continue
		}
		if filters.RepoID != nil {
			bases, worktrees := []EpicBase{}, []EpicWorktree{}
			for _, base := range epic.Base {
				if base.RepoID == *filters.RepoID {
					bases = append(bases, base)
				}
			}
			for _, wt := range epic.Worktrees {
				if wt.RepoID == *filters.RepoID {
					worktrees = append(worktrees, wt)
				}
			}
			if len(bases) == 0 {
				continue
			}
			epic.Base, epic.Worktrees = bases, worktrees
		}
		out.Epics = append(out.Epics, epic)
		projectEpics[epic.ProjectID] = true
	}
	for _, project := range s.Projects.Projects {
		if filters.ProjectID != nil && project.ID != *filters.ProjectID {
			continue
		}
		if filters.EpicID != nil && !projectEpics[project.ID] {
			continue
		}
		repos := []workspace.RepoID{}
		for _, id := range project.RepoIDs {
			if filters.RepoID == nil || id == *filters.RepoID {
				repos = append(repos, id)
			}
		}
		if filters.RepoID != nil && len(repos) == 0 {
			continue
		}
		sort.Slice(repos, func(i, j int) bool { return repos[i] < repos[j] })
		out.Projects = append(out.Projects, WorkflowStatusProject{ProjectID: project.ID, Name: project.Name, RepoIDs: repos})
	}
	sort.Slice(out.Projects, func(i, j int) bool { return out.Projects[i].ProjectID < out.Projects[j].ProjectID })
	sort.Slice(out.Epics, func(i, j int) bool {
		if out.Epics[i].ProjectID != out.Epics[j].ProjectID {
			return out.Epics[i].ProjectID < out.Epics[j].ProjectID
		}
		return out.Epics[i].EpicID < out.Epics[j].EpicID
	})
}

func (row *WorkflowStatusItem) addAction(source, kind, actor, reason, severity string, since *string, evidence []string, run, queue string) {
	if actor == "" {
		actor = "unknown"
	}
	row.NextActions = append(row.NextActions, WorkflowStatusAction{Action: Action{Source: source, Kind: kind, Actor: actor, Reason: reason, Severity: severity, SinceUTC: since, EvidenceIDs: sortedUnique(evidence)}, Current: true, RunID: stringValue(run), QueueID: stringValue(queue)})
}

func (out *WorkflowStatus) addDiagnostic(code, source, reason, task, run, queue string, evidence []string) {
	out.Diagnostics = append(out.Diagnostics, WorkflowDiagnostic{Code: code, Source: source, Actor: "unknown", Severity: "attention", Reason: reason, TaskID: stringValue(task), RunID: stringValue(run), QueueID: stringValue(queue), EvidenceIDs: sortedUnique(evidence)})
}

func workflowRunEvidence(run taskrun.InventoryRun) []string {
	ids := []string{run.RunID}
	for _, source := range run.Sources {
		ids = append(ids, source.SHA256)
	}
	return sortedUnique(ids)
}

func workflowTargetMatches(t workspace.QueueTarget, f workspace.TaskListFilters) bool {
	return (f.ProjectID == nil || t.ProjectID == *f.ProjectID) && (f.RepoID == nil || t.RepoID == *f.RepoID) && (f.EpicID == nil || t.EpicID == *f.EpicID)
}

func combineFreshness(a, b string) string {
	if a == "stale" || b == "stale" {
		return "stale"
	}
	if a == "unknown" || b == "unknown" {
		return "unknown"
	}
	return "fresh"
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func sortWorkflowActions(actions []WorkflowStatusAction) {
	sort.Slice(actions, func(i, j int) bool {
		a, b := actions[i], actions[j]
		if a.Current != b.Current {
			return a.Current
		}
		if a.Severity != b.Severity {
			return a.Severity == "attention"
		}
		if c := compareWorkflowTime(a.SinceUTC, b.SinceUTC); c != 0 {
			return c < 0
		}
		return strings.Join([]string{a.Source, a.Kind, a.Actor, valueOrEmpty(a.RunID), valueOrEmpty(a.QueueID), a.Reason}, "\x00") < strings.Join([]string{b.Source, b.Kind, b.Actor, valueOrEmpty(b.RunID), valueOrEmpty(b.QueueID), b.Reason}, "\x00")
	})
}

func compareWorkflowTime(a, b *string) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1
	}
	if b == nil {
		return -1
	}
	at, ae := time.Parse(time.RFC3339Nano, *a)
	bt, be := time.Parse(time.RFC3339Nano, *b)
	if ae != nil || be != nil || at.Equal(bt) {
		return 0
	}
	if at.Before(bt) {
		return -1
	}
	return 1
}

func workflowItemLess(a, b WorkflowStatusItem) bool {
	rank := func(category string) int {
		for i, c := range workflowCategories {
			if c == category {
				return i
			}
		}
		return len(workflowCategories)
	}
	if a.Category != b.Category {
		return rank(a.Category) < rank(b.Category)
	}
	if a.Category == "needs_you" || a.Category == "follow_up" {
		var aa, ba *WorkflowStatusAction
		for i := range a.NextActions {
			if a.NextActions[i].Current {
				aa = &a.NextActions[i]
				break
			}
		}
		for i := range b.NextActions {
			if b.NextActions[i].Current {
				ba = &b.NextActions[i]
				break
			}
		}
		if aa != nil && ba != nil {
			if aa.Severity != ba.Severity {
				return aa.Severity == "attention"
			}
			if c := compareWorkflowTime(aa.SinceUTC, ba.SinceUTC); c != 0 {
				return c < 0
			}
		}
	}
	if a.ProjectID != b.ProjectID {
		return a.ProjectID < b.ProjectID
	}
	if a.RepoID != b.RepoID {
		return a.RepoID < b.RepoID
	}
	if a.EpicID != b.EpicID {
		return a.EpicID < b.EpicID
	}
	if a.Category == "ready_next" && len(a.Queues) > 0 && len(b.Queues) > 0 {
		aq, bq := a.Queues[0], b.Queues[0]
		if aq.QueueID == bq.QueueID && aq.Rank != nil && bq.Rank != nil && *aq.Rank != *bq.Rank {
			return *aq.Rank < *bq.Rank
		}
	}
	return a.TaskID < b.TaskID
}

type workflowTaskBasis struct {
	problems    map[workspace.TaskID]workspace.TaskProblemReference
	selections  map[workspace.TaskID]workspace.TaskSelectionReference
	assessments map[workflowAssessmentKey]workspace.TaskAssessmentReference
	specs       map[workflowAssessmentKey]string
}

type workflowAssessmentKey struct {
	task     workspace.TaskID
	spec     string
	revision int
}

func workflowBasisIndex(r workspace.WorkItemRegistry) workflowTaskBasis {
	b := workflowTaskBasis{problems: map[workspace.TaskID]workspace.TaskProblemReference{}, selections: map[workspace.TaskID]workspace.TaskSelectionReference{}, assessments: map[workflowAssessmentKey]workspace.TaskAssessmentReference{}, specs: map[workflowAssessmentKey]string{}}
	for _, spec := range r.TaskSpecRevisions {
		b.specs[workflowAssessmentKey{spec.TaskID, spec.SpecID, spec.Revision}] = spec.ManifestSHA256
	}
	for _, p := range r.TaskProblemRevisions {
		if p.Revision > b.problems[p.TaskID].Revision {
			b.problems[p.TaskID] = p
		}
	}
	for _, s := range r.TaskSolutionSelections {
		if s.Ordinal > b.selections[s.TaskID].Ordinal {
			b.selections[s.TaskID] = s
		}
	}
	for _, assessment := range r.TaskSpecAssessments {
		key := workflowAssessmentKey{assessment.TaskID, assessment.SpecID, assessment.SpecRevision}
		if assessment.Ordinal > b.assessments[key].Ordinal {
			b.assessments[key] = assessment
		}
	}
	return b
}

func (b workflowTaskBasis) matches(taskID workspace.TaskID, basis *workspace.TaskSpecBasis) bool {
	if basis == nil {
		return true
	}
	if basis.TaskID != taskID {
		return false
	}
	if p, exists := b.problems[taskID]; !exists || p.Revision != basis.Problem.Revision || p.ManifestSHA256 != basis.Problem.ManifestSHA256 {
		return false
	}
	if s, exists := b.selections[taskID]; !exists || s.ID != basis.Selection.ID || s.ManifestSHA256 != basis.Selection.ManifestSHA256 {
		return false
	}
	key := workflowAssessmentKey{taskID, basis.SpecID, basis.Spec.Revision}
	if digest, exists := b.specs[key]; !exists || digest != basis.Spec.ManifestSHA256 {
		return false
	}
	if a, exists := b.assessments[key]; !exists || a.ID != basis.Assessment.ID || a.ManifestSHA256 != basis.Assessment.ManifestSHA256 {
		return false
	}
	return true
}
