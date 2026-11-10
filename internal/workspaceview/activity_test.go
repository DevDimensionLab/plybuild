package workspaceview

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestActivitySinceUsesInclusiveUTCAndPositiveDuration(t *testing.T) {
	now := time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC)
	for _, v := range []string{"1h", "2026-10-06T09:00:00Z"} {
		got, err := ParseActivitySince(v, now)
		if err != nil || got == nil || !got.Equal(now.Add(-time.Hour)) {
			t.Fatalf("%q: %v %v", v, got, err)
		}
	}
	for _, v := range []string{"0s", "-1h", "bad", "2026-10-06T11:00:00+02:00"} {
		if _, err := ParseActivitySince(v, now); err == nil {
			t.Fatalf("accepted %q", v)
		}
	}
}

func TestProcessWaitingRequiresOpenUnsupersededExplicitHuman(t *testing.T) {
	s := taskjournal.Snapshot{Events: []taskjournal.Event{{EventID: "start", Origin: "contribution", Type: "step_started", Title: "Need a decision", RecordedAt: stringValue("2026-10-06T09:00:00Z"), SourceIDs: []string{"s"}, SupersededBy: []string{}}}, Sources: []taskjournal.Source{{SourceID: "s", Status: "valid"}}, Steps: []taskjournal.Step{{StepID: "waiting", StartEventID: "start", State: "no_end_recorded", Waiting: &taskjournal.Waiting{Reason: "Choose behavior", NextActor: "human"}}}, Coverage: taskjournal.Coverage{State: "complete", Reasons: []taskjournal.Reason{}}}
	f := processFactsFromJournal(s)
	if !f.HasProgress || len(f.Waiting) != 1 || f.Waiting[0].Actor != "human" {
		t.Fatalf("missing open human wait: %+v", f)
	}
	s.Steps[0].EndEventID = stringValue("end")
	if f = processFactsFromJournal(s); len(f.Waiting) != 0 {
		t.Fatalf("closed wait survived: %+v", f)
	}
	s.Steps[0].EndEventID = nil
	s.Steps[0].Waiting.NextActor = "some developer"
	if f = processFactsFromJournal(s); len(f.Waiting) != 0 {
		t.Fatalf("guessed actor from free text: %+v", f)
	}
	s.Steps[0].Waiting.NextActor = "human"
	s.Events[0].SupersededBy = []string{"replacement"}
	if f = processFactsFromJournal(s); len(f.Waiting) != 0 {
		t.Fatalf("superseded wait survived: %+v", f)
	}
}

func TestPlanningAndLifecycleAloneAreNotProcessProgress(t *testing.T) {
	s := taskjournal.Snapshot{Events: []taskjournal.Event{{Type: "task_state", Origin: "native"}, {Type: "record_problem", Origin: "native"}, {Type: "record_spec", Origin: "native"}, {Type: "lifecycle_changed", Origin: "native"}}, Coverage: taskjournal.Coverage{State: "complete"}}
	if f := processFactsFromJournal(s); f.HasProgress {
		t.Fatalf("planning became work progress: %+v", f)
	}
}

func TestUnknownEventClockDoesNotMakeFreshSourcesUnavailable(t *testing.T) {
	s := taskjournal.Snapshot{Sources: []taskjournal.Source{{SourceID: "s", Status: "valid"}}, Coverage: taskjournal.Coverage{State: "partial", Reasons: []taskjournal.Reason{{Code: "time_conflict", Detail: "Occurrence order is uncertain"}}}}
	if f := processFactsFromJournal(s); f.Freshness != "fresh" {
		t.Fatalf("temporal coverage became source drift: %+v", f)
	}
	s.Sources[0].Status = "missing"
	if f := processFactsFromJournal(s); f.Freshness != "unknown" {
		t.Fatalf("missing source was fresh: %+v", f)
	}
}

type activityClock struct{ now time.Time }

func (c activityClock) Now() time.Time { return c.now }

func TestActivityMatchesTaskJournalIncludesLifecycleAndFiltersInclusively(t *testing.T) {
	d, root := changesFixture(t)
	baseTime := time.Date(2026, 10, 6, 9, 0, 0, 0, time.UTC)
	d.WorkClock = activityClock{baseTime}
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, repo, "init", "-b", "epic")
	if err := os.WriteFile(filepath.Join(repo, "README"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, repo, "add", "README")
	inventoryGit(t, repo, "commit", "-m", "fixture")
	if _, err := workspace.AddProject(d, workspace.ProjectAddInput{ProjectID: "project", Name: "Activity", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AdoptEpic(d, workspace.EpicAdoptInput{EpicID: "epic", Title: "Epic", ProjectID: "project", RepoID: "repo", Worktree: repo, Ref: "refs/heads/epic", ExpectedOID: inventoryGit(t, repo, "rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	service := taskjournal.New(d)
	recorded := baseTime.Add(time.Hour)
	service.Now = func() time.Time { return recorded }
	source := filepath.Join(root, "evidence.txt")
	if err := os.WriteFile(source, []byte("synthetic process fact"), 0600); err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, id := range []string{"a", "b"} {
		if _, err := workspace.CreateTask(d, workspace.TaskCreateInput{TaskID: workspace.TaskID(id), Title: id, Description: "Activity fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
			t.Fatal(err)
		}
		input := taskjournal.EventInput{Common: taskjournal.Common{Kind: "ply.workspace.task-journal-event-input", SchemaVersion: 1, PublicationKey: "activity:" + id, TaskID: id, Actor: taskjournal.Actor{ID: "fixture-agent", Role: "developer"}, TimeBasis: taskjournal.TimeBasis{Kind: "unknown", Precision: "unknown"}, Sources: []taskjournal.SourceRef{{Locator: source, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte("synthetic process fact")))}}}, Type: "note", Title: "A reported fact " + id, Relations: []taskjournal.Relation{}}
		b, _ := json.Marshal(input)
		file := filepath.Join(root, "input-"+id+".json")
		if err := os.WriteFile(file, b, 0600); err != nil {
			t.Fatal(err)
		}
		result, err := service.Append(id, file, "event")
		if err != nil {
			t.Fatal(err)
		}
		ids[result.EventID] = true
	}
	d.WorkClock = activityClock{recorded.Add(time.Minute)}
	if _, err := workspace.SetWorkItemLifecycle(d, workspace.WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "a", State: workspace.LifecycleFrozen, ActorClaim: "Synthetic organizer"}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.SetWorkItemLifecycle(d, workspace.WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: "epic", State: workspace.LifecycleParked, ActorClaim: "Synthetic organizer"}); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	before := inventoryFingerprint(t, root)
	all, err := ListActivity(d, s, ActivityOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, inventoryFingerprint(t, root)) {
		t.Fatal("activity changed files")
	}
	known := map[string]ActivityEvent{}
	for _, e := range all.Events {
		known[e.EventID] = e
	}
	for _, id := range []string{"a", "b"} {
		single, err := taskjournal.New(d).Show(id, taskjournal.Options{})
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range single.Events {
			v, ok := known[e.EventID]
			if !ok || v.Summary != e.Title || !reflect.DeepEqual(v.RecordedAtUTC, e.RecordedAt) || !reflect.DeepEqual(v.OccurredAtUTC, e.OccurredAt) {
				t.Fatalf("Task event differs: %+v => %+v", e, v)
			}
		}
	}
	lifecycle := 0
	for _, e := range all.Events {
		if e.Type == "lifecycle_changed" {
			lifecycle++
		}
	}
	if lifecycle != 2 || all.UnknownTimeCount < 2 {
		t.Fatalf("lifecycle/time facts lost: %+v", all)
	}
	filtered, err := ListActivity(d, s, ActivityOptions{Since: recorded.Format(time.RFC3339), ProjectID: "project"})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range filtered.Events {
		delete(ids, e.EventID)
		if e.RecordedAtUTC == nil {
			t.Fatal("undated event passed --since")
		}
	}
	if len(ids) != 0 {
		t.Fatalf("inclusive boundary excluded records: %v", ids)
	}
	limited, err := ListActivity(d, s, ActivityOptions{Limit: 1})
	if err != nil || !limited.Truncated || len(limited.Events) != 1 {
		t.Fatalf("limit: %+v %v", limited, err)
	}
	for i := 1; i < len(all.Events); i++ {
		a, b := all.Events[i-1], all.Events[i]
		if a.RecordedAtUTC == nil && b.RecordedAtUTC != nil {
			t.Fatal("unknown times were not last")
		}
		if a.RecordedAtUTC != nil && b.RecordedAtUTC != nil {
			x, _ := time.Parse(time.RFC3339Nano, *a.RecordedAtUTC)
			y, _ := time.Parse(time.RFC3339Nano, *b.RecordedAtUTC)
			if x.Before(y) || x.Equal(y) && a.EventID > b.EventID {
				t.Fatal("unstable activity order")
			}
		}
	}
}
