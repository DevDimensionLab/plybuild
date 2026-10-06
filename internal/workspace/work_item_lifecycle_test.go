package workspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestWorkItemLifecycleAbsentReadAndTransitionsPreserveRegistry(t *testing.T) {
	for _, legacy := range []bool{true, false} {
		f := newWorkItemJourneyFixtureVersion(t, legacy)
		root := filepath.Dir(f.wrapper)
		registryBefore, err := os.ReadFile(workItemsPath(root))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(root, MarkerDirectory, workItemsLock)); err != nil {
			t.Fatal(err)
		}
		snapshot, err := ReadWorkItemLifecycle(f.dependencies, root)
		if err != nil || snapshot.Exists || snapshot.Revision != 0 || snapshot.Task("task") != LifecycleActive || snapshot.Epic("epic") != LifecycleActive {
			t.Fatalf("absent read = %#v, %v", snapshot, err)
		}
		if _, err := os.Stat(filepath.Join(root, MarkerDirectory, workItemsLock)); !os.IsNotExist(err) {
			t.Fatalf("read created lock: %v", err)
		}
		for _, subject := range []struct{ kind, id string }{{"task", "task"}, {"epic", "epic"}} {
			for _, state := range []LifecycleState{LifecycleParked, LifecycleFrozen, LifecycleArchived, LifecycleActive} {
				result, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: state, ActorClaim: "fixture-agent"})
				if err != nil || !result.Changed || result.Lifecycle != state || result.Event == nil || result.Event.ActorClaim != "fixture-agent" {
					t.Fatalf("set %s %s = %#v, %v", subject.kind, state, result, err)
				}
				before, err := os.ReadFile(result.Source.Locator)
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(result.Source.Locator)
				if err != nil {
					t.Fatal(err)
				}
				repeat, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: state, ActorClaim: "another-agent"})
				if err != nil || repeat.Changed || repeat.Revision != result.Revision || repeat.Event != nil {
					t.Fatalf("repeat = %#v, %v", repeat, err)
				}
				after, err := os.ReadFile(result.Source.Locator)
				if err != nil {
					t.Fatal(err)
				}
				afterInfo, err := os.Stat(result.Source.Locator)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(before, after) || !info.ModTime().Equal(afterInfo.ModTime()) {
					t.Fatal("no-op rewrote metadata")
				}
			}
		}
		after, err := os.ReadFile(workItemsPath(root))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(registryBefore, after) {
			t.Fatal("lifecycle changed historical registry bytes")
		}
	}
}

func TestWorkItemLifecycleRejectsCorruptionWithoutDefaultingToActive(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	input := WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: "actor"}
	if _, err := SetWorkItemLifecycle(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	good, err := os.ReadFile(lifecycleLocator(root))
	if err != nil {
		t.Fatal(err)
	}
	for name, bad := range map[string][]byte{
		"empty":           {},
		"unknown version": []byte(strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 2`, 1)),
		"null history":    []byte(`{"kind":"WorkspaceWorkItemLifecycleStore@1","schema_version":1,"events":null}`),
		"duplicate key":   []byte(strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 1, "schema_version": 1`, 1)),
		"unknown field":   []byte(strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 1, "unknown": true`, 1)),
		"changed event":   []byte(strings.Replace(string(good), `"to": "archived"`, `"to": "active"`, 1)),
	} {
		t.Run(name, func(t *testing.T) {
			if err := os.WriteFile(lifecycleLocator(root), bad, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := ReadWorkItemLifecycle(f.dependencies, root); err == nil {
				t.Fatal("corrupt lifecycle read as active")
			}
			if _, err := SetWorkItemLifecycle(f.dependencies, input); err == nil {
				t.Fatal("corrupt store was overwritten")
			}
			after, err := os.ReadFile(lifecycleLocator(root))
			if err != nil || !bytes.Equal(after, bad) {
				t.Fatalf("corruption changed: %v", err)
			}
		})
	}
}

func TestWorkItemLifecycleRejectsSymlinksIncludingDanglingLinks(t *testing.T) {
	for _, dangling := range []bool{false, true} {
		f := newWorkItemJourneyFixture(t)
		root := filepath.Dir(f.wrapper)
		target := filepath.Join(t.TempDir(), "target")
		if !dangling {
			if err := os.WriteFile(target, []byte("preserve\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.Symlink(target, lifecycleLocator(root)); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadWorkItemLifecycle(f.dependencies, root); err == nil {
			t.Fatal("symlink read accepted")
		}
		if _, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: "actor"}); err == nil {
			t.Fatal("symlink write accepted")
		}
		if _, err := os.Readlink(lifecycleLocator(root)); err != nil {
			t.Fatalf("symlink replaced: %v", err)
		}
		if !dangling {
			if got, err := os.ReadFile(target); err != nil || string(got) != "preserve\n" {
				t.Fatalf("target changed: %q %v", got, err)
			}
		}
	}
}

func TestWorkItemLifecycleFaultsLeaveOldOrCompleteNewHistory(t *testing.T) {
	for _, stage := range []string{"temp-open", "write", "file-sync", "replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			input := WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleParked, ActorClaim: "actor"}
			if _, err := SetWorkItemLifecycle(f.dependencies, input); err != nil {
				t.Fatal(err)
			}
			input.State = LifecycleArchived
			injected := errors.New("injected " + stage)
			_, err := setWorkItemLifecycle(f.dependencies, input, workspaceMetadataPublisher{fault: func(point string) error {
				if point == stage {
					return injected
				}
				return nil
			}})
			if !errors.Is(err, injected) {
				t.Fatalf("fault = %v", err)
			}
			snapshot, err := ReadWorkItemLifecycle(f.dependencies, root)
			if err != nil {
				t.Fatal(err)
			}
			wantState, wantCount := LifecycleParked, 1
			if stage == "directory-sync" {
				wantState, wantCount = LifecycleArchived, 2
			}
			if snapshot.Task("task") != wantState || len(snapshot.Events) != wantCount {
				t.Fatalf("fault left partial history: %#v", snapshot)
			}
			result, err := SetWorkItemLifecycle(f.dependencies, input)
			if err != nil || result.Revision != 2 || result.Changed != (stage != "directory-sync") {
				t.Fatalf("recovery duplicated or lost event: %#v %v", result, err)
			}
			matches, err := filepath.Glob(filepath.Join(root, MarkerDirectory, ".metadata-*.tmp"))
			if err != nil || len(matches) != 0 {
				t.Fatalf("temporary files left: %v %v", matches, err)
			}
		})
	}
}

func TestWorkItemLifecycleConcurrentChangesAndCompareAndSet(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	const count = 8
	for i := 0; i < count; i++ {
		_, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: TaskID(fmt.Sprintf("task-%d", i)), Title: "Parallel Task", Description: "Fixture", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"})
		if err != nil {
			t.Fatal(err)
		}
	}
	errorsOut := make(chan error, count)
	var group sync.WaitGroup
	for i := 0; i < count; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			_, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: fmt.Sprintf("task-%d", i), State: LifecycleArchived, ActorClaim: "parallel-agent"})
			errorsOut <- err
		}(i)
	}
	group.Wait()
	close(errorsOut)
	for err := range errorsOut {
		if err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := ReadWorkItemLifecycle(f.dependencies, root)
	if err != nil || snapshot.Revision != count {
		t.Fatalf("lost concurrent changes: %#v %v", snapshot, err)
	}
	for i := 0; i < count; i++ {
		if snapshot.Task(TaskID(fmt.Sprintf("task-%d", i))) != LifecycleArchived {
			t.Fatalf("lost Task %d", i)
		}
	}
	revision := snapshot.Revision
	changes := make(chan bool, 2)
	for _, state := range []LifecycleState{LifecycleParked, LifecycleFrozen} {
		group.Add(1)
		go func(state LifecycleState) {
			defer group.Done()
			out, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: state, ActorClaim: "parallel-agent", ExpectedRevision: &revision})
			changes <- err == nil && out.Changed
		}(state)
	}
	group.Wait()
	close(changes)
	wins := 0
	for won := range changes {
		if won {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("compare-and-set winners = %d", wins)
	}
}

func TestWorkItemLifecycleDoesNotGateExistingWorktreeWorkflow(t *testing.T) {
	f := newWorkItemJourneyFixtureVersion(t, true)
	for _, subject := range []struct{ kind, id string }{{"task", "task"}, {"epic", "epic"}} {
		if _, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: LifecycleArchived, ActorClaim: "actor"}); err != nil {
			t.Fatal(err)
		}
	}
	input, err := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if err != nil {
		t.Fatal(err)
	}
	result, err := CreateTaskWorktree(f.dependencies, input)
	if err != nil || result.Task.Worktree == nil {
		t.Fatalf("organizing metadata gated existing workflow: %#v %v", result, err)
	}
}

func TestWorkItemLifecycleAppearsInExistingListAndShowJSON(t *testing.T) {
	for _, legacy := range []bool{true, false} {
		f := newWorkItemJourneyFixtureVersion(t, legacy)
		for _, subject := range []struct{ kind, id string }{{"task", "task"}, {"epic", "epic"}} {
			if _, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: LifecycleFrozen, ActorClaim: "actor"}); err != nil {
				t.Fatal(err)
			}
		}
		task, err := ShowTask(f.dependencies, "task")
		if err != nil {
			t.Fatal(err)
		}
		ep, err := ShowEpic(f.dependencies, "epic")
		if err != nil {
			t.Fatal(err)
		}
		list, err := ListTasksWithFilters(f.dependencies, TaskListFilters{})
		if err != nil {
			t.Fatal(err)
		}
		for _, marshal := range []func() ([]byte, error){func() ([]byte, error) { return MarshalTaskReadback(task) }, func() ([]byte, error) { return MarshalEpicReadback(ep) }} {
			b, err := marshal()
			if err != nil {
				t.Fatal(err)
			}
			var value map[string]any
			if err := json.Unmarshal(b, &value); err != nil {
				t.Fatal(err)
			}
			if value["lifecycle"] != "frozen" {
				t.Fatalf("missing lifecycle: %s", b)
			}
		}
		b, err := MarshalTaskList(list)
		if err != nil {
			t.Fatal(err)
		}
		var value struct {
			Tasks []struct {
				Lifecycle string `json:"lifecycle"`
			} `json:"tasks"`
		}
		if err := json.Unmarshal(b, &value); err != nil || len(value.Tasks) != 1 || value.Tasks[0].Lifecycle != "frozen" {
			t.Fatalf("list lifecycle: %s %v", b, err)
		}
	}
}

func TestWorkItemLifecyclePreservesV2AndV4RegistryBindings(t *testing.T) {
	for _, version := range []int{2, 4} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, version == 2)
			root := filepath.Dir(f.wrapper)
			if version == 2 {
				if err := f.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
					r, err := session.Snapshot()
					if err != nil {
						return err
					}
					upgradeRegistryToV2(&r)
					return session.Publish(r)
				}); err != nil {
					t.Fatal(err)
				}
			} else {
				queueSetFixture(t, f, "metadata-fixture", []QueueEntry{})
			}
			before, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil || r.FormatVersion != version {
				t.Fatalf("fixture version %d: %#v %v", version, r, err)
			}
			if out, err := ShowWorkItemLifecycle(f.dependencies, "task", "task"); err != nil || out.Lifecycle != LifecycleActive || out.Source.Exists {
				t.Fatalf("legacy default %#v %v", out, err)
			}
			if _, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: "actor"}); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(workItemsPath(root))
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("v%d bindings changed %v", version, err)
			}
			if _, err := f.dependencies.WorkItems.Snapshot(root); err != nil {
				t.Fatalf("v%d closure no longer valid: %v", version, err)
			}
		})
	}
}

func TestWorkItemLifecyclePreservesHistoricalResultAndIntegrationBinding(t *testing.T) {
	f := newIntegrationJourneyFixtureVersion(t, true)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	before, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	for _, subject := range []struct{ kind, id string }{{"task", "task"}, {"epic", "epic"}} {
		if _, err := SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: LifecycleArchived, ActorClaim: "actor"}); err != nil {
			t.Fatal(err)
		}
	}
	after, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil || integrationDigest(t, before.Readback) != integrationDigest(t, after.Readback) {
		t.Fatalf("metadata changed bound integration plan: %v", err)
	}
	shown, err := ShowTask(f.dependencies, "task")
	if err != nil || shown.Integration == nil || shown.Integration.TaskResult == nil || shown.Integration.TaskResult.ID != f.resultID {
		t.Fatalf("historical result unreadable: %#v %v", shown.Integration, err)
	}
	b, err := MarshalTaskReadback(shown)
	if err != nil || !bytes.Contains(b, []byte(`"lifecycle":"archived"`)) || !bytes.Contains(b, []byte(`"task_result":{`)) {
		t.Fatalf("historical show: %s %v", b, err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, before.Readback)
	if result, err := ApplyTaskIntegration(f.dependencies, input); err != nil || result.Readback.Classification != "exact_effect" {
		t.Fatalf("metadata gated bound integration: %#v %v", result.Readback, err)
	}
}

func TestWorkItemLifecycleInvalidInputsAndExpectedRevisionDoNotPublish(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	for _, input := range []WorkItemLifecycleInput{
		{SubjectKind: "task", SubjectID: "task", State: "deleted", ActorClaim: "actor"},
		{SubjectKind: "project", SubjectID: "ply", State: LifecycleArchived, ActorClaim: "actor"},
		{SubjectKind: "task", SubjectID: "missing", State: LifecycleArchived, ActorClaim: "actor"},
		{SubjectKind: "epic", SubjectID: "missing", State: LifecycleArchived, ActorClaim: "actor"},
		{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: ""},
		{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: "actor\nother"},
	} {
		if _, err := SetWorkItemLifecycle(f.dependencies, input); err == nil {
			t.Fatalf("accepted invalid input %#v", input)
		}
	}
	if snapshot, err := ReadWorkItemLifecycle(f.dependencies, root); err != nil || snapshot.Exists {
		t.Fatalf("invalid input wrote metadata: %#v %v", snapshot, err)
	}
	zero := 0
	input := WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleParked, ActorClaim: "actor", ExpectedRevision: &zero}
	result, err := SetWorkItemLifecycle(f.dependencies, input)
	if err != nil || result.Revision != 1 {
		t.Fatalf("initial write %#v %v", result, err)
	}
	input.State = LifecycleArchived
	if _, err := SetWorkItemLifecycle(f.dependencies, input); err == nil {
		t.Fatal("stale expected revision accepted")
	}
	snapshot, err := ReadWorkItemLifecycle(f.dependencies, root)
	if err != nil || snapshot.Task("task") != LifecycleParked || len(snapshot.Events) != 1 {
		t.Fatalf("conflict altered state: %#v %v", snapshot, err)
	}
}
