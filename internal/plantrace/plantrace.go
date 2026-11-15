// Package plantrace reads the recorded workflow of one plan: its revisions,
// the Tasks whose journal links them to the plan, and every event of those
// journals, without interpretation. A drawing tool builds its own model.
package plantrace

import (
	"encoding/json"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Kind is the versioned contract of a plan trace.
const Kind = "PlyPlanTrace@1"

type Trace struct {
	Kind          string              `json:"kind"`
	SchemaVersion int                 `json:"schema_version"`
	Project       workspace.ProjectID `json:"project"`
	Plan          Plan                `json:"plan"`
	Tasks         []Task              `json:"tasks"`
	Events        []Event             `json:"events"`
	Reasons       []Reason            `json:"reasons"`
}

type Plan struct {
	ID        string              `json:"id"`
	Title     string              `json:"title"`
	Status    string              `json:"status"`
	Revisions []plan.HistoryEntry `json:"revisions"`
}

// Task is one Task whose latest plan link names the plan.
type Task struct {
	TaskID   string  `json:"task_id"`
	Title    string  `json:"title"`
	PlanStep *string `json:"plan_step"`
	Coverage string  `json:"coverage"`
}

// Event is one journal event with what a drawing needs; Data is the event as
// the journal holds it.
type Event struct {
	TaskID       string            `json:"task_id"`
	EventID      string            `json:"event_id"`
	Type         string            `json:"type"`
	Title        string            `json:"title"`
	ActorID      *string           `json:"actor_id"`
	Role         *string           `json:"role"`
	StepID       *string           `json:"step_id"`
	StepKind     *string           `json:"step_kind"`
	ParentStepID *string           `json:"parent_step_id"`
	Outcome      *string           `json:"outcome"`
	OccurredAt   *string           `json:"occurred_at"`
	RecordedAt   *string           `json:"recorded_at"`
	PlanStep     *string           `json:"plan_step"`
	PlyOperation *string           `json:"ply_operation"`
	Refs         map[string]string `json:"refs"`
}

// Reason is a coverage gap of one Task journal.
type Reason struct {
	TaskID string `json:"task_id"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// recorded is the part of a journal event's data this trace reads.
type recorded struct {
	Step *taskjournal.StepRef   `json:"step"`
	Ply  *taskjournal.PlyAction `json:"ply"`
}

// Read returns the trace of planID in project, whose plans live in plans.
func Read(d workspace.Dependencies, project workspace.ProjectID, plans plan.Store, planID string) (Trace, error) {
	view, err := plans.Show(planID, 0)
	if err != nil {
		return Trace{}, err
	}
	history, err := plans.History(planID)
	if err != nil {
		return Trace{}, err
	}
	out := Trace{Kind: Kind, SchemaVersion: 1, Project: project, Plan: Plan{planID, view.Plan.Title, view.Plan.Status, history}, Tasks: []Task{}, Events: []Event{}, Reasons: []Reason{}}
	basis, err := workspace.ReadViewBasis(d)
	if err != nil {
		return Trace{}, err
	}
	selected := basis.Registry
	selected.Tasks = nil
	for _, t := range basis.Registry.Tasks {
		if t.ProjectID == project {
			selected.Tasks = append(selected.Tasks, t)
		}
	}
	// History needs the Tasks' content closure, not queue or Epic base readiness.
	selected.TaskQueueEvents, selected.EpicBaseUpdates = nil, nil
	lifecycle, err := workspace.ReadWorkItemLifecycle(d, basis.Workspace.Root)
	if err != nil {
		return Trace{}, err
	}
	batch, err := taskjournal.ReadEventBatch(d, taskjournal.EventBatchInput{Workspace: basis.Workspace, Projects: basis.Projects, Registry: selected, Lifecycle: lifecycle})
	if err != nil {
		return Trace{}, err
	}
	for _, t := range selected.Tasks {
		snap := batch.Tasks[t.ID]
		var link *taskjournal.PlyAction
		events := make([]recorded, len(snap.Events))
		for i, e := range snap.Events {
			_ = json.Unmarshal(e.Data, &events[i]) // native events carry other data
			if p := events[i].Ply; p != nil && p.Operation == "plan_link" {
				link = p
			}
		}
		if link == nil || link.Plan == nil || *link.Plan != planID {
			continue
		}
		out.Tasks = append(out.Tasks, Task{string(t.ID), t.Title, link.PlanStep, snap.Coverage.State})
		for _, r := range snap.Coverage.Reasons {
			out.Reasons = append(out.Reasons, Reason{string(t.ID), r.Code, r.Detail})
		}
		for i, e := range snap.Events {
			ev := Event{TaskID: string(t.ID), EventID: e.EventID, Type: e.Type, Title: e.Title, Outcome: e.Outcome, OccurredAt: e.OccurredAt, RecordedAt: e.RecordedAt, PlanStep: link.PlanStep, StepID: e.StepID}
			if e.Actor != nil {
				ev.ActorID, ev.Role = &e.Actor.ID, &e.Actor.Role
			}
			if s := events[i].Step; s != nil {
				ev.StepKind, ev.ParentStepID = &s.Kind, s.ParentStepID
			}
			if p := events[i].Ply; p != nil {
				ev.PlyOperation, ev.Refs = &p.Operation, p.Refs
			}
			out.Events = append(out.Events, ev)
		}
	}
	sort.SliceStable(out.Events, func(i, j int) bool {
		a, b := at(out.Events[i]), at(out.Events[j])
		if a != b {
			return a < b
		}
		return out.Events[i].EventID < out.Events[j].EventID
	})
	return out, nil
}

// at is the time an event is ordered by: occurred, else recorded, else last.
func at(e Event) string {
	switch {
	case e.OccurredAt != nil:
		return *e.OccurredAt
	case e.RecordedAt != nil:
		return *e.RecordedAt
	}
	return "~"
}
