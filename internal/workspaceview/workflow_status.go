package workspaceview

import (
	"sort"
	"strings"
	"time"

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
	queues := workspace.ReadRegisteredTaskQueues(d, s.Workspace, s.Projects, s.Registry)
	s.RefreshFreshness(d)
	return BuildWorkflowStatus(s, o, process, queues)
}

// BuildWorkflowStatus is a pure projection. The native Progress remains intact;
// presentation actions supplement it without manufacturing result/QA facts.
func BuildWorkflowStatus(s *Snapshot, o WorkflowStatusOptions, process map[workspace.TaskID]ProcessFacts, queues workspace.RegisteredTaskQueues) (WorkflowStatus, error) {
	version := o.SchemaVersion
	if version == 0 {
		version = 1
	}
	if version != 1 && version != 2 {
		return WorkflowStatus{}, workspace.WorkInvalidArguments("workflow status schema-version must be 1 or 2")
	}
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
	if version == 2 {
		out.Kind = "WorkflowStatusReadback@2"
		out.SchemaVersion = 2
	}
	rows := make(map[workspace.TaskID]*WorkflowStatusItem, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		row := &WorkflowStatusItem{TaskID: task.TaskID, ProjectID: task.ProjectID, RepoID: task.RepoID, EpicID: task.ParentEpicID, TaskListTitle: task.TaskListTitle, WorktreeState: task.WorktreeState, TaskLifecycle: task.Lifecycle, EpicLifecycle: task.EpicLifecycle, Progress: task.Progress, LastActivityUTC: task.Progress.LastActivityUTC, Freshness: task.Progress.Freshness, NextActions: []WorkflowStatusAction{}, Queues: []WorkflowStatusQueue{}}
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
	case row.Progress.State == "cleanup_pending":
		row.Category = "follow_up"
		return
	case row.TaskLifecycle == workspace.LifecycleCompleted && row.Progress.State == "completed":
		row.Category = "completed"
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
