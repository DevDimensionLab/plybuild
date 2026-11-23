// Package plantrace reads the recorded workflow of one plan: its revisions,
// the Tasks whose journal links them to the plan, and every journal record of
// those Tasks, without interpretation. It reads only journals and the plan,
// not the work-item registry: Ply records its own actions in the journals.
package plantrace

import (
	"fmt"
	"sort"
	"time"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskactor"
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

// Task is one Task whose latest plan link names the plan, with the role its
// goal register gave it.
type Task struct {
	TaskID   string  `json:"task_id"`
	Title    string  `json:"title"`
	PlanStep *string `json:"plan_step"`
	Role     string  `json:"role"`
}

// Event is one journal record with what a drawing needs. A record of the
// plan's own journal has an empty TaskID. ActorBasis is derived when the
// record was stored with actor unknown and the trace found the actor with
// the rule Ply records new actors by.
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
	ActorBasis   string            `json:"actor_basis,omitempty"`
}

// Reason is a journal that could not be read completely; an empty TaskID is
// the plan's own journal.
type Reason struct {
	TaskID string `json:"task_id"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ReasonTimeUnreadable: an event's time is not RFC 3339; the event is ordered first.
const ReasonTimeUnreadable = "time_unreadable"

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
	w, err := workspace.ObserveContaining(d)
	if err != nil {
		return Trace{}, err
	}
	journals, err := taskjournal.Journals(w.Root)
	if err != nil {
		return Trace{}, err
	}
	for _, journal := range journals {
		task := journal.Task
		records, readErr := taskjournal.ReadContributions(w, journal.Work, task)
		link := taskjournal.LatestPlanLink(records)
		if link == nil || link.Plan == nil || *link.Plan != planID || link.Refs["project"] != string(project) {
			continue
		}
		if readErr != nil {
			out.Reasons = append(out.Reasons, Reason{task, "journal_damaged", "Only the validated prefix of the journal is read: " + readErr.Error()})
		}
		out.Tasks = append(out.Tasks, Task{task, link.Refs["task_title"], link.PlanStep, taskactor.Role(records)})
		agent := taskactor.StartedAgent(journal.Work, task)
		for _, r := range records {
			out.Events = append(out.Events, derived(event(task, link, r), r, agent, view.Plan.Owner))
		}
	}
	journal, err := plans.Journal(planID)
	if err != nil {
		out.Reasons = append(out.Reasons, Reason{"", "journal_damaged", "Only the readable prefix of the plan journal is read: " + err.Error()})
	}
	planEvents := make([]Event, 0, len(journal))
	for i, e := range journal {
		planEvents = append(planEvents, planEvent(planID, i, e))
	}
	out.Reasons = append(out.Reasons, timeReasons(append(append([]Event{}, out.Events...), planEvents...))...)
	out.Events = mergeEvents(out.Events, planEvents)
	return out, nil
}

// timeReasons names each event whose ordering time cannot be read.
func timeReasons(events []Event) []Reason {
	var reasons []Reason
	for _, e := range events {
		if stamp := orderStamp(e); stamp != nil {
			if _, err := time.Parse(time.RFC3339Nano, *stamp); err != nil {
				reasons = append(reasons, Reason{e.TaskID, ReasonTimeUnreadable, fmt.Sprintf("Event %s has an unreadable time %q; it is ordered first.", e.EventID, *stamp)})
			}
		}
	}
	return reasons
}

func mergeEvents(taskEvents, planEvents []Event) []Event {
	sort.SliceStable(taskEvents, func(i, j int) bool {
		return eventBefore(taskEvents[i], taskEvents[j])
	})

	// The plan journal already has numeric EventID order. Only compare its
	// next entry with the next task event so clock changes cannot reorder it.
	events := make([]Event, 0, len(taskEvents)+len(planEvents))
	for len(taskEvents) > 0 && len(planEvents) > 0 {
		if eventBefore(taskEvents[0], planEvents[0]) {
			events = append(events, taskEvents[0])
			taskEvents = taskEvents[1:]
		} else {
			events = append(events, planEvents[0])
			planEvents = planEvents[1:]
		}
	}
	events = append(events, taskEvents...)
	return append(events, planEvents...)
}

func eventBefore(a, b Event) bool {
	if order := at(a).Compare(at(b)); order != 0 {
		return order < 0
	}
	return a.EventID < b.EventID
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

// derived gives an event stored with actor unknown the actor Ply records new
// events with: owner for Ply's operations, agent for the agent's records.
func derived(e Event, r taskjournal.Contribution, agent, owner string) Event {
	if *e.ActorID != taskactor.Unknown || r.Event == nil {
		return e
	}
	if actor := taskactor.Of(r.Event.Ply != nil, agent, owner); actor != taskactor.Unknown {
		e.ActorID, e.ActorBasis = &actor, "derived"
	}
	return e
}

// planEvent is record i of the plan's own journal, with an empty TaskID. It
// takes e by value, so the pointers it returns are its own (go.mod is go 1.18).
func planEvent(planID string, i int, e plan.JournalEntry) Event {
	return Event{EventID: fmt.Sprintf("%s-journal-%d", planID, i+1), Type: e.Type, Title: e.Title, ActorID: &e.ActorID, Role: &e.Role,
		StepID: e.StepID, StepKind: e.StepKind, ParentStepID: e.ParentStepID, Outcome: e.Outcome, OccurredAt: &e.At, RecordedAt: &e.At}
}

// orderStamp is the time an event is ordered by: occurred, else recorded.
func orderStamp(e Event) *string {
	if e.OccurredAt != nil {
		return e.OccurredAt
	}
	return e.RecordedAt
}

// at is orderStamp as a time; an unreadable one is the zero time, which
// timeReasons reports.
func at(e Event) time.Time {
	stamp := orderStamp(e)
	if stamp == nil {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339Nano, *stamp)
	if err != nil {
		return time.Time{}
	}
	return parsed
}
