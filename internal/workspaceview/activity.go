package workspaceview

import (
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const ActivityKind = "WorkspaceActivityReadback@1"

type ActivityOptions struct {
	Since, ProjectID string
	Limit            int
}
type ActivityEvent struct {
	EventID       string             `json:"event_id"`
	SubjectKind   string             `json:"subject_kind"`
	SubjectID     string             `json:"subject_id"`
	TaskID        *string            `json:"task_id"`
	EpicID        *string            `json:"epic_id"`
	ProjectID     string             `json:"project_id"`
	Type          string             `json:"type"`
	Origin        string             `json:"origin"`
	RecordedAtUTC *string            `json:"recorded_at_utc"`
	OccurredAtUTC *string            `json:"occurred_at_utc"`
	Summary       string             `json:"summary"`
	Actor         *taskjournal.Actor `json:"actor"`
	RunID         *string            `json:"run_id"`
	SourceIDs     []string           `json:"source_ids"`
	Freshness     string             `json:"freshness"`
}
type ActivityResult struct {
	Kind             string               `json:"kind"`
	SchemaVersion    int                  `json:"schema_version"`
	Workspace        WorkspaceRef         `json:"workspace"`
	AsOf             string               `json:"as_of"`
	SinceUTC         *string              `json:"since_utc"`
	Events           []ActivityEvent      `json:"events"`
	Sources          []taskjournal.Source `json:"sources"`
	UnknownTimeCount int                  `json:"unknown_time_count"`
	Truncated        bool                 `json:"truncated"`
	Freshness        string               `json:"freshness"`
	Coverage         string               `json:"coverage"`
	Reasons          []string             `json:"reasons"`
}

type processHistory struct {
	journal taskjournal.EventBatch
	runs    taskrun.Inventory
	facts   map[workspace.TaskID]ProcessFacts
}

func readProcessHistory(d workspace.Dependencies, s *Snapshot) (*processHistory, error) {
	if s.process != nil {
		return s.process, nil
	}
	j, err := taskjournal.ReadEventBatch(d, taskjournal.EventBatchInput{Workspace: s.Workspace, Projects: s.Projects, Registry: s.Registry, Lifecycle: s.Lifecycle})
	if err != nil {
		return nil, err
	}
	r, err := taskrun.ReadInventory(taskrun.SystemDependencies(d), s.Workspace.Root)
	if err != nil {
		return nil, err
	}
	h := &processHistory{journal: j, runs: r, facts: map[workspace.TaskID]ProcessFacts{}}
	for id, snap := range j.Tasks {
		h.facts[id] = processFactsFromJournal(snap)
	}
	for _, run := range r.Runs {
		if run.TaskID == nil {
			continue
		}
		id := workspace.TaskID(*run.TaskID)
		f, ok := h.facts[id]
		if !ok {
			continue
		}
		f.HasProgress = true
		f.LastActivityUTC = latestTime(f.LastActivityUTC, run.LastSeenUTC)
		if run.Freshness != "fresh" {
			f.Freshness = "unknown"
			f.Reasons = append(f.Reasons, "run_evidence_unavailable")
		}
		f.Reasons = sortedUnique(f.Reasons)
		h.facts[id] = f
	}
	if len(r.Reasons) > 0 {
		for id, f := range h.facts {
			f.Freshness = "unknown"
			f.Reasons = append(f.Reasons, "run_store_incomplete")
			h.facts[id] = f
		}
	}
	s.process = h
	return h, nil
}

// ReadProcessFacts is one memoized batch shared by Task progress, Epic counts
// and Attention. It never derives native technical/human QA from agent claims.
func ReadProcessFacts(d workspace.Dependencies, s *Snapshot) (map[workspace.TaskID]ProcessFacts, error) {
	h, err := readProcessHistory(d, s)
	if err != nil {
		return nil, err
	}
	return h.facts, nil
}

func processFactsFromJournal(s taskjournal.Snapshot) ProcessFacts {
	f := ProcessFacts{Waiting: []ProcessWaiting{}, Freshness: "fresh", Reasons: taskjournal.StableBatchReasons(s)}
	for _, source := range s.Sources {
		if source.Status != "valid" {
			f.Freshness = "unknown"
		}
	}
	for _, reason := range f.Reasons {
		switch reason {
		case "source_invalid", "source_missing", "source_changed", "run_unbound", "task_identity_unavailable", "lifecycle_unavailable":
			f.Freshness = "unknown"
		}
	}
	events := map[string]taskjournal.Event{}
	for _, e := range s.Events {
		events[e.EventID] = e
		f.LastActivityUTC = latestTime(f.LastActivityUTC, e.RecordedAt)
		if e.Origin == "contribution" || strings.HasPrefix(e.Type, "run_") || strings.HasPrefix(e.Type, "task_result_") || e.Type == "human_qa" || e.Type == "integration" {
			f.HasProgress = true
		}
	}
	for _, step := range s.Steps {
		if step.EndEventID != nil || step.Waiting == nil {
			continue
		}
		e, ok := events[step.StartEventID]
		if !ok || len(e.SupersededBy) > 0 || eventFreshness(e, s.Sources) != "fresh" {
			continue
		}
		actor := step.Waiting.NextActor
		// Legacy next_actor is free text. Only explicit role values carry a
		// machine-readable actor; do not infer human ownership from a name.
		if actor == "user" {
			actor = "human"
		}
		if actor != "human" && actor != "agent" && actor != "ply" {
			continue
		}
		f.Waiting = append(f.Waiting, ProcessWaiting{Kind: "recorded_wait", Reason: step.Waiting.Reason, Actor: actor, EventID: e.EventID, SinceUTC: e.RecordedAt})
	}
	sort.Slice(f.Waiting, func(i, j int) bool { return f.Waiting[i].EventID < f.Waiting[j].EventID })
	return f
}

// ParseActivitySince accepts an absolute UTC clock or a positive Go duration.
// Callers may validate argv before workspace discovery using this same parser.
func ParseActivitySince(value string, now time.Time) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		_, offset := t.Zone()
		if offset != 0 {
			return nil, workspace.WorkInvalidArguments("--since timestamp must be UTC")
		}
		t = t.UTC()
		return &t, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return nil, workspace.WorkInvalidArguments("--since must be a UTC RFC3339 timestamp or a positive duration")
	}
	t := now.UTC().Add(-duration)
	return &t, nil
}

func validateProjectFilter(s *Snapshot, id string) error {
	if id == "" {
		return nil
	}
	p, err := workspace.ParseProjectID(id)
	if err != nil {
		return err
	}
	return ValidateTaskFilters(s, workspace.TaskListFilters{ProjectID: &p})
}

func ListActivity(d workspace.Dependencies, s *Snapshot, o ActivityOptions) (ActivityResult, error) {
	out := ActivityResult{Kind: ActivityKind, SchemaVersion: 1, Workspace: WorkspaceRef{s.Workspace.Root}, AsOf: s.ObservedAtUTC, Events: []ActivityEvent{}, Sources: []taskjournal.Source{}, Freshness: s.Freshness, Coverage: "complete", Reasons: []string{}}
	if o.Limit < 0 {
		return out, workspace.WorkInvalidArguments("--limit must not be negative")
	}
	if err := validateProjectFilter(s, o.ProjectID); err != nil {
		return out, err
	}
	now, err := time.Parse(time.RFC3339Nano, s.ObservedAtUTC)
	if err != nil {
		now = time.Now().UTC()
	}
	since, err := ParseActivitySince(o.Since, now)
	if err != nil {
		return out, err
	}
	if since != nil {
		x := since.Format(time.RFC3339Nano)
		out.SinceUTC = &x
	}
	h, err := readProcessHistory(d, s)
	if err != nil {
		return out, err
	}
	sources := map[string]taskjournal.Source{}
	appendEvents := func(kind, id, project, epic string, events []taskjournal.Event, src []taskjournal.Source) {
		for _, e := range events {
			if e.RecordedAt == nil {
				out.UnknownTimeCount++
				if since != nil {
					continue
				}
			}
			if since != nil && e.RecordedAt != nil {
				t, err := time.Parse(time.RFC3339Nano, *e.RecordedAt)
				if err != nil || t.Before(*since) {
					continue
				}
			}
			var task *string
			if kind == "task" {
				task = stringValue(id)
			}
			out.Events = append(out.Events, ActivityEvent{e.EventID, kind, id, task, stringValue(epic), project, e.Type, e.Origin, e.RecordedAt, e.OccurredAt, e.Title, e.Actor, e.RunID, append([]string{}, e.SourceIDs...), eventFreshness(e, src)})
			for _, source := range src {
				for _, used := range e.SourceIDs {
					if source.SourceID == used {
						sources[used] = source
					}
				}
			}
		}
	}
	for _, task := range s.Registry.Tasks {
		if o.ProjectID != "" && string(task.ProjectID) != o.ProjectID {
			continue
		}
		journal := h.journal.Tasks[task.ID]
		appendEvents("task", string(task.ID), string(task.ProjectID), string(task.ParentEpicID), journal.Events, journal.Sources)
		if journal.Coverage.State == "partial" {
			out.Coverage = "partial"
			out.Reasons = append(out.Reasons, taskjournal.StableBatchReasons(journal)...)
		}
	}
	for _, epic := range s.Registry.Epics {
		if o.ProjectID != "" && string(epic.ProjectID) != o.ProjectID {
			continue
		}
		appendEvents("epic", string(epic.ID), string(epic.ProjectID), string(epic.ID), h.journal.EpicEvents[epic.ID], h.journal.EpicSources[epic.ID])
	}
	sort.Slice(out.Events, func(i, j int) bool {
		a, b := out.Events[i], out.Events[j]
		if a.RecordedAtUTC == nil && b.RecordedAtUTC != nil {
			return false
		}
		if a.RecordedAtUTC != nil && b.RecordedAtUTC == nil {
			return true
		}
		if a.RecordedAtUTC != nil && b.RecordedAtUTC != nil {
			x, _ := time.Parse(time.RFC3339Nano, *a.RecordedAtUTC)
			y, _ := time.Parse(time.RFC3339Nano, *b.RecordedAtUTC)
			if !x.Equal(y) {
				return x.After(y)
			}
		}
		return a.EventID < b.EventID
	})
	if o.Limit > 0 && len(out.Events) > o.Limit {
		out.Events = out.Events[:o.Limit]
		out.Truncated = true
	}
	needed := map[string]bool{}
	for _, e := range out.Events {
		for _, id := range e.SourceIDs {
			needed[id] = true
		}
	}
	for id, source := range sources {
		if needed[id] {
			out.Sources = append(out.Sources, source)
		}
	}
	sort.Slice(out.Sources, func(i, j int) bool { return out.Sources[i].SourceID < out.Sources[j].SourceID })
	s.RefreshFreshness(d)
	out.Freshness = s.Freshness
	out.Reasons = sortedUnique(append(out.Reasons, s.Reasons...))
	if out.Freshness != "fresh" {
		out.Coverage = "partial"
	}
	return out, nil
}

func eventFreshness(e taskjournal.Event, sources []taskjournal.Source) string {
	if len(e.SourceIDs) == 0 {
		return "unknown"
	}
	result := "fresh"
	for _, id := range e.SourceIDs {
		status := "missing"
		for _, s := range sources {
			if s.SourceID == id {
				status = s.Status
				break
			}
		}
		if status == "changed" {
			result = "stale"
		} else if status != "valid" {
			return "unknown"
		}
	}
	return result
}
