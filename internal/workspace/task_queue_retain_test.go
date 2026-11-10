package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func retainQueueDraft(t *testing.T, f workItemJourneyFixture, key string, revision int, entries []TaskGoalQueueEntry) string {
	t.Helper()
	draft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: key, ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "fixture", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}, RegistryUpgrade: queueUpgrade(mustQueueRegistry(t, f))}
	raw, err := contentCanonical(draft)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(filepath.Dir(f.wrapper), key+".json")
	if err = os.WriteFile(file, raw, 0600); err != nil {
		t.Fatal(err)
	}
	return file
}

func TestTaskQueueRetainConcurrentRevision(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	entry := TaskGoalQueueEntry{TaskID: "task", Goal: &goal}
	q := goalQueueFixture(t, f, []TaskGoalQueueEntry{entry})
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.wrapper, "task"), ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	files := []string{retainQueueDraft(t, f, "retain-a", q.Revision, []TaskGoalQueueEntry{entry}), retainQueueDraft(t, f, "remove-b", q.Revision, []TaskGoalQueueEntry{})}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, file := range files {
		wg.Add(1)
		go func(file string) { defer wg.Done(); _, err := SetTaskQueue(f.dependencies, file); errs <- err }(file)
	}
	wg.Wait()
	close(errs)
	passed, conflicted := 0, 0
	for err := range errs {
		var qe *TaskContentError
		if err == nil {
			passed++
		} else if errors.As(err, &qe) && qe.Code == "task_queue_revision_conflict" {
			conflicted++
		} else {
			t.Fatal(err)
		}
	}
	if passed != 1 || conflicted != 1 {
		t.Fatalf("passed %d, conflicts %d", passed, conflicted)
	}
}

func TestTaskQueueRetainStillValidatesReferencesAndTarget(t *testing.T) {
	for _, variant := range []string{"goal"} {
		t.Run(variant, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			entry := TaskGoalQueueEntry{TaskID: "task"}
			var digest string
			if variant == "goal" {
				goal := goalFixture(t, f, "task", "goal")
				entry.Goal, digest = &goal, goal.Spec.ManifestSHA256
			} else {
				sel := queueSelect(t, f, "task", "selection", nil)
				entry.Selection, digest = &sel, sel.ManifestSHA256
			}
			q := goalQueueFixture(t, f, []TaskGoalQueueEntry{entry})
			file := retainQueueDraft(t, f, "invalid", q.Revision, []TaskGoalQueueEntry{entry})
			raw, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(file, []byte(strings.Replace(string(raw), `"epic_id":"epic"`, `"epic_id":"missing"`, 1)), 0600); err != nil {
				t.Fatal(err)
			}
			state := queueStateBytes(t, f)
			if _, err = SetTaskQueue(f.dependencies, file); err == nil || state != queueStateBytes(t, f) {
				t.Fatalf("invalid target changed state: %v", err)
			}
			if err = os.WriteFile(file, raw, 0600); err != nil {
				t.Fatal(err)
			}
			path := taskContentPath(filepath.Dir(f.wrapper), "manifests", digest)
			if err = os.Chmod(path, 0600); err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(path, []byte("corrupt"), 0600); err != nil {
				t.Fatal(err)
			}
			state = queueStateBytes(t, f)
			if _, err = SetTaskQueue(f.dependencies, file); err == nil || state != queueStateBytes(t, f) {
				t.Fatalf("corrupt retained reference changed state: %v", err)
			}
		})
	}
}
