// Package plantrace reads the recorded workflow of one plan: its revisions,
// the Tasks whose journal links them to the plan, and every journal record of
// those Tasks, without interpretation. It reads only journals and the plan,
// not the work-item registry: Ply records its own actions in the journals.
package plantrace

import (
	"fmt"
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
}

// Event is one journal record with what a drawing needs. A record of the
// plan's own journal has an empty TaskID.
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

// Reason is a journal that could not be read completely; an empty TaskID is
// the plan's own journal.
type Reason struct {
	TaskID string `json:"task_id"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
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
	w, err := workspace.ObserveContainingRoot(d)
	if err != nil {
		return Trace{}, err
	}
	tasks, err := taskjournal.JournalTasks(w.Root)
	if err != nil {
		return Trace{}, err
	}
	for _, task := range tasks {
		records, readErr := taskjournal.ReadContributions(w, task)
		var link *taskjournal.PlyAction
		for _, r := range records {
			if r.Event != nil && r.Event.Ply != nil && r.Event.Ply.Operation == "plan_link" {
				link = r.Event.Ply
			}
		}
		if link == nil || link.Plan == nil || *link.Plan != planID || link.Refs["project"] != string(project) {
			continue
		}
		if readErr != nil {
			out.Reasons = append(out.Reasons, Reason{task, "journal_damaged", "Only the validated prefix of the journal is read: " + readErr.Error()})
		}
		out.Tasks = append(out.Tasks, Task{task, link.Refs["task_title"], link.PlanStep})
		for _, r := range records {
			out.Events = append(out.Events, event(task, link, r))
		}
	}
	journal, err := plans.Journal(planID)
	if err != nil {
		out.Reasons = append(out.Reasons, Reason{"", "journal_damaged", "Only the readable prefix of the plan journal is read: " + err.Error()})
	}
	for i, e := range journal {
		// go.mod is go 1.18, so e is shared across iterations: copy what is pointed to.
		at, actor, role := e.At, e.ActorID, e.Role
		out.Events = append(out.Events, Event{EventID: fmt.Sprintf("%s-journal-%d", planID, i+1), Type: e.Type, Title: e.Title, ActorID: &actor, Role: &role,
			StepID: e.StepID, StepKind: e.StepKind, ParentStepID: e.ParentStepID, Outcome: e.Outcome, OccurredAt: &at, RecordedAt: &at})
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

func event(task string, link *taskjournal.PlyAction, r taskjournal.Contribution) Event {
	recorded := r.RecordedAt
	e := Event{TaskID: task, EventID: r.EventID, ActorID: &r.Common.Actor.ID, Role: &r.Common.Actor.Role, OccurredAt: r.Common.OccurredAt, RecordedAt: &recorded, PlanStep: link.PlanStep}
	if x := r.Event; x != nil {
		e.Type, e.Title, e.Outcome = x.Type, x.Title, x.Outcome
		if s := x.Step; s != nil {
			e.StepID, e.StepKind, e.ParentStepID = &s.ID, &s.Kind, s.ParentStepID
		}
		if p := x.Ply; p != nil {
			e.PlyOperation, e.Refs = &p.Operation, p.Refs
		}
	} else if x := r.Observation; x != nil {
		e.Type, e.Title = "observation_"+x.Type, x.Finding
	}
	return e
}

// at is the time an event is ordered by: occurred, else recorded.
func at(e Event) string {
	if e.OccurredAt != nil {
		return *e.OccurredAt
	}
	return *e.RecordedAt
}
