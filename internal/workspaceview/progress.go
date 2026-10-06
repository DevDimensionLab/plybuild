package workspaceview

import (
	"sort"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const (
	TaskListKind = "WorkspaceTaskProgressListReadback@1"
	TaskKind     = "WorkspaceTaskProgressReadback@1"
)

type Action struct {
	Source      string   `json:"source"`
	Kind        string   `json:"kind"`
	Actor       string   `json:"actor"`
	Reason      string   `json:"reason"`
	Severity    string   `json:"severity"`
	SinceUTC    *string  `json:"since_utc"`
	EvidenceIDs []string `json:"evidence_ids"`
}

type Progress struct {
	State                     string   `json:"state"`
	HasProgress               bool     `json:"has_progress"`
	IntegrationClassification string   `json:"integration_classification"`
	TechnicalGate             *string  `json:"technical_gate"`
	HumanQAOutcome            *string  `json:"human_qa_outcome"`
	ResultCount               int      `json:"result_count"`
	TaskResultID              *string  `json:"task_result_id"`
	IntegrationResultID       *string  `json:"integration_result_id"`
	NextActions               []Action `json:"next_actions"`
	LastActivityUTC           *string  `json:"last_activity_utc"`
	Freshness                 string   `json:"freshness"`
	SourceBasis               string   `json:"source_basis"`
	Reasons                   []string `json:"reasons"`
}

type TaskRow struct {
	TaskID        workspace.TaskID         `json:"task_id"`
	ProjectID     workspace.ProjectID      `json:"project_id"`
	RepoID        workspace.RepoID         `json:"repo_id"`
	ParentEpicID  workspace.EpicID         `json:"parent_epic_id"`
	WorktreeState workspace.WorkItemState  `json:"worktree_state"`
	Lifecycle     workspace.LifecycleState `json:"lifecycle"`
	EpicLifecycle workspace.LifecycleState `json:"epic_lifecycle"`
	workspace.TaskListTitle
	Progress Progress `json:"progress"`
}

type TaskList struct {
	Kind          string                    `json:"kind"`
	SchemaVersion int                       `json:"schema_version"`
	Workspace     WorkspaceRef              `json:"workspace"`
	Scope         workspace.TaskListFilters `json:"scope"`
	Tasks         []TaskRow                 `json:"tasks"`
	Freshness     string                    `json:"freshness"`
	SourceBasis   string                    `json:"source_basis"`
	Reasons       []string                  `json:"reasons"`
}

type TaskReadback struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	Workspace     WorkspaceRef `json:"workspace"`
	Task          TaskRow      `json:"task"`
	Freshness     string       `json:"freshness"`
	SourceBasis   string       `json:"source_basis"`
	Reasons       []string     `json:"reasons"`
}

func ReadTaskList(d workspace.Dependencies, filters workspace.TaskListFilters) (TaskList, error) {
	if err := filters.Validate(); err != nil {
		return TaskList{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return TaskList{}, err
	}
	if err := ValidateTaskFilters(s, filters); err != nil {
		return TaskList{}, err
	}
	process := readProgressProcess(d, s)
	s.RefreshFreshness(d)
	return BuildTaskList(s, filters, process)
}

func ReadTask(d workspace.Dependencies, id workspace.TaskID) (TaskReadback, error) {
	if _, err := workspace.ParseTaskID(string(id)); err != nil {
		return TaskReadback{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return TaskReadback{}, err
	}
	found := false
	for _, task := range s.Registry.Tasks {
		found = found || task.ID == id
	}
	if !found {
		return TaskReadback{}, notFound("Task", string(id))
	}
	process := readProgressProcess(d, s)
	s.RefreshFreshness(d)
	list, err := BuildTaskList(s, workspace.TaskListFilters{}, process)
	if err != nil {
		return TaskReadback{}, err
	}
	for _, row := range list.Tasks {
		if row.TaskID == id {
			return TaskReadback{TaskKind, 1, list.Workspace, row, list.Freshness, list.SourceBasis, list.Reasons}, nil
		}
	}
	return TaskReadback{}, notFound("Task", string(id))
}

func ReadTaskProgress(d workspace.Dependencies, id workspace.TaskID) (Progress, error) {
	result, err := ReadTask(d, id)
	return result.Task.Progress, err
}

func readProgressProcess(d workspace.Dependencies, s *Snapshot) map[workspace.TaskID]ProcessFacts {
	process, err := ReadProcessFacts(d, s)
	if err == nil {
		return process
	}
	// A damaged optional process history must not hide independently validated
	// registrations/results. Unavailability remains explicit in the response.
	s.Freshness = "unknown"
	s.Reasons = sortedUnique(append(s.Reasons, "process_history_unavailable"))
	process = map[workspace.TaskID]ProcessFacts{}
	for _, task := range s.Registry.Tasks {
		process[task.ID] = ProcessFacts{Freshness: "unknown", Reasons: []string{"process_history_unavailable"}}
	}
	return process
}

func BuildTaskList(s *Snapshot, filters workspace.TaskListFilters, process map[workspace.TaskID]ProcessFacts) (TaskList, error) {
	if err := ValidateTaskFilters(s, filters); err != nil {
		return TaskList{}, err
	}
	result := TaskList{TaskListKind, 1, WorkspaceRef{s.Workspace.Root}, filters, []TaskRow{}, s.Freshness, "registered", sortedUnique(s.Reasons)}
	facts := workspace.RecordedTaskProgressFacts(s.Registry)
	for _, task := range s.Registry.Tasks {
		if !MatchesTaskFilters(task, filters) {
			continue
		}
		title, ok := s.Titles[task.ID]
		if !ok {
			registeredTitle := task.Title
			title = workspace.TaskListTitle{Title: &registeredTitle, Source: "registration", Status: "available"}
		}
		p := buildProgress(task, facts[task.ID], process[task.ID], s.Freshness, s.Reasons)
		for _, event := range s.Lifecycle.Events {
			if event.SubjectKind == "task" && event.SubjectID == string(task.ID) {
				p.LastActivityUTC = latestTime(p.LastActivityUTC, &event.RecordedAtUTC)
			}
		}
		result.Tasks = append(result.Tasks, TaskRow{task.ID, task.ProjectID, task.RepoID, task.ParentEpicID, task.WorktreeState, s.Lifecycle.Task(task.ID), s.Lifecycle.Epic(task.ParentEpicID), title, p})
	}
	sort.Slice(result.Tasks, func(i, j int) bool { return result.Tasks[i].TaskID < result.Tasks[j].TaskID })
	return result, nil
}

func BuildTaskProgress(s *Snapshot, id workspace.TaskID, process map[workspace.TaskID]ProcessFacts) (Progress, error) {
	if _, err := workspace.ParseTaskID(string(id)); err != nil {
		return Progress{}, err
	}
	list, err := BuildTaskList(s, workspace.TaskListFilters{}, process)
	if err != nil {
		return Progress{}, err
	}
	for _, row := range list.Tasks {
		if row.TaskID == id {
			return row.Progress, nil
		}
	}
	return Progress{}, notFound("Task", string(id))
}

func buildProgress(task workspace.TaskRecord, fact workspace.TaskRecordedProgressFacts, process ProcessFacts, freshness string, reasons []string) Progress {
	p := Progress{State: "nothing", IntegrationClassification: "unknown", ResultCount: fact.ResultCount, HasProgress: fact.HasProgress || process.HasProgress, NextActions: []Action{}, LastActivityUTC: latestTime(fact.LastActivityUTC, process.LastActivityUTC), Freshness: freshness, SourceBasis: "registered", Reasons: sortedUnique(append(append([]string{}, reasons...), process.Reasons...))}
	if process.Freshness != "" && process.Freshness != "fresh" {
		p.Freshness = "unknown"
	}
	if p.HasProgress {
		p.State = "in_progress"
	}
	if fact.TaskResult != nil && !fact.Ambiguous {
		p.TaskResultID = stringValue(string(fact.TaskResult.ID))
		p.TechnicalGate = stringValue(fact.TaskResult.TechnicalGate)
	}
	if fact.HumanQA != nil && !fact.Ambiguous {
		p.HumanQAOutcome = stringValue(fact.HumanQA.Outcome)
	}
	add := func(source, kind, actor, reason, severity string, since *string, evidence string) {
		ids := []string{}
		if evidence != "" {
			ids = append(ids, evidence)
		}
		p.NextActions = append(p.NextActions, Action{source, kind, actor, reason, severity, since, ids})
	}
	switch {
	case fact.Ambiguous:
		p.State, p.IntegrationClassification = "attention", "conflict"
		p.Reasons = append(p.Reasons, "recorded_result_selection_conflict")
		add("integration", "resolve_integration_conflict", "human", "Choose the intended result and QA binding among conflicting recorded facts.", "attention", fact.LastActivityUTC, "")
	case fact.ResultBasisStale:
		p.Reasons = append(p.Reasons, "recorded_result_basis_stale")
		add("content", "inspect_task_basis", "agent", "Historical results do not cover the current registered Task basis; inspect the selected work before continuing.", "next_step", fact.LastActivityUTC, "")
	case fact.IntegrationResult != nil:
		integration := fact.IntegrationResult
		p.IntegrationClassification = integration.Outcome
		p.IntegrationResultID = stringValue(string(integration.ID))
		switch integration.Outcome {
		case "exact_effect", "already_integrated":
			p.State = "integrated"
		case "conflict", "partial", "unknown":
			p.State = "attention"
			add("integration", "resolve_integration_conflict", "human", "Inspect the recorded integration outcome and choose the recovery action.", "attention", stringValue(integration.RecordedAtUTC), string(integration.ID))
		case "no_effect", "blocked":
			p.State = "attention"
			add("integration", "inspect_integration", "agent", "Inspect the recorded integration prerequisite before another attempt.", "attention", stringValue(integration.RecordedAtUTC), string(integration.ID))
		}
	case fact.IntegrationPending:
		add("integration", "inspect_integration", "agent", "Follow the existing integration before considering another attempt.", "next_step", fact.LastActivityUTC, "")
	case p.TechnicalGate != nil && *p.TechnicalGate == "failed":
		p.State = "attention"
		add("content", "correct_task_result", "agent", "Correct the failed technical result before requesting product verification.", "attention", fact.LastActivityUTC, string(fact.TaskResult.ID))
	case p.HumanQAOutcome != nil && (*p.HumanQAOutcome == "fail" || *p.HumanQAOutcome == "blocked"):
		p.State = "attention"
		add("content", "follow_human_qa", "agent", "Follow the recorded human QA outcome and its observations.", "attention", stringValue(fact.HumanQA.Actor.CompletedAtUTC), string(fact.HumanQA.ID))
	case p.TechnicalGate != nil && (*p.TechnicalGate == "passed" || *p.TechnicalGate == "good_enough_with_known_debt") && p.HumanQAOutcome == nil:
		add("content", "record_human_qa", "human", "Try the technically controlled result and record your product judgment.", "next_step", stringValue(fact.TaskResult.Recorder.RecordedAtUTC), string(fact.TaskResult.ID))
	case p.HumanQAOutcome != nil && *p.HumanQAOutcome == "pass":
		add("integration", "check_integration", "agent", "Check the current local integration basis; no live Git readiness is asserted by this overview.", "next_step", stringValue(fact.HumanQA.Actor.CompletedAtUTC), string(fact.HumanQA.ID))
	case task.WorktreeState == workspace.WorkItemReconciliationRequired:
		p.State = "attention"
		add("readiness", "inspect_worktree", "agent", "Inspect the recorded worktree reconciliation state.", "attention", nil, "")
	}
	for _, waiting := range process.Waiting {
		if !p.HasProgress || waiting.Actor != "human" && waiting.Actor != "agent" && waiting.Actor != "ply" {
			continue
		}
		kind := waiting.Kind
		if kind == "" {
			kind = "resolve_wait"
		}
		add("journal", kind, waiting.Actor, waiting.Reason, "next_step", waiting.SinceUTC, waiting.EventID)
	}
	p.Reasons = sortedUnique(p.Reasons)
	sort.Slice(p.NextActions, func(i, j int) bool {
		a, b := p.NextActions[i], p.NextActions[j]
		if a.Source != b.Source {
			return a.Source < b.Source
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Actor != b.Actor {
			return a.Actor < b.Actor
		}
		return a.Reason < b.Reason
	})
	return p
}
