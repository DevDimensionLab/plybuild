package weekrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Only the queue is synthetic. Child readback, executable validation, retained
// attempt evidence and Project/WorkItem locking use the production adapters.
type mr05NativeQueue struct {
	NativeAdapter
	queue   *fakeNative
	effects *int
}

func (n mr05NativeQueue) Queue(c Contract) (workspace.WorkspaceTaskQueueReadback, error) {
	return n.queue.Queue(c)
}
func (n mr05NativeQueue) Start(Contract, Selection, taskrun.Executable, bool) (Run, error) {
	*n.effects++
	return Run{}, fmt.Errorf("unexpected native start in readback fixture")
}
func (n mr05NativeQueue) Resume(Contract, string, bool) (Run, error) {
	*n.effects++
	return Run{}, fmt.Errorf("unexpected native resume in readback fixture")
}
func (n mr05NativeQueue) Control(Contract, string, string) (Run, error) {
	*n.effects++
	return Run{}, fmt.Errorf("unexpected native control in readback fixture")
}
func (n mr05NativeQueue) Preview(Contract, Selection, taskrun.Executable) error {
	*n.effects++
	return fmt.Errorf("unexpected next Goal preview in readback fixture")
}
func (n mr05NativeQueue) Release(Contract, string, workspace.TaskID) error {
	*n.effects++
	return fmt.Errorf("unexpected release in readback fixture")
}

type mr05ChildWorkspace struct {
	workflowhandoff.WorkspaceObserver
	root string
}

func (w mr05ChildWorkspace) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	return workflowhandoff.WorkspaceSnapshot{Observation: workspace.WorkspaceObservation{Root: w.root}}, nil
}

type mr05ChildHandoff struct{ workflowhandoff.Store }

func (mr05ChildHandoff) ReadByLocator(string) (workflowhandoff.Snapshot, error) {
	return workflowhandoff.Snapshot{}, nil
}

type mr05ChildGit struct {
	workspace.TaskIntegrationGit
	observation workspace.IntegrationWorktreeObservation
}

func (g mr05ChildGit) ObserveIntegrationWorktree(string, string) (workspace.IntegrationWorktreeObservation, error) {
	return g.observation, nil
}

func mr05NativeChild(t *testing.T) (*fixture, string, string, string, *int) {
	t.Helper()
	f := newMR05Week(t)
	root := f.c.Workspace
	file := func(path string, raw []byte, mode os.FileMode) taskrun.FileBinding {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, raw, mode); err != nil {
			t.Fatal(err)
		}
		return taskrun.FileBinding{Locator: path, SHA256: hash(raw)}
	}
	provider := file(filepath.Join(root, "child-provider"), []byte("#!/bin/sh\nexit 99\n"), 0700)
	goal := file(filepath.Join(root, "child-goal.json"), []byte(`{"fixture":"MR05"}`), 0600)
	handoff := file(filepath.Join(root, "child-handoff.json"), []byte(`{"fixture":"isolated child readback"}`), 0600)
	r := taskrun.WorkflowRequest{Envelope: taskrun.Envelope{Kind: "ply.workflow.headless-run-request", SchemaVersion: 2}, RequestKey: "mr05-native-child", WorkspaceRoot: root,
		HandoffDraft: json.RawMessage(`{"inputs":[]}`), ReturnMode: "delivery_owner",
		Runtime:  taskrun.Runtime{Provider: "codex", Mode: "headless", Model: f.c.Model, Executable: taskrun.Executable{Path: provider.Locator, SHA256: provider.SHA256}, PlyExecutable: f.b.Control, PermissionBinding: taskrun.Permission{Evidence: []taskrun.Evidence{}}},
		Worker:   &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 120},
		Delivery: &taskrun.DeliveryContract{Goal: goal, OwnerClaim: "isolated fixture", ReasoningEffort: "max"}, ControlContract: f.b.ControlContract}
	id := "wfr_" + strings.TrimPrefix(taskrun.ControlRevision([]string{root, r.RequestKey}), "sha256:")
	runRoot := filepath.Join(root, ".ply", "workflow-runs", "runs", id)
	jobDir := filepath.Join(runRoot, "worker", "attempts", "reviewer-0001")
	intent := map[string]any{"run_id": id, "job_id": "reviewer-0001", "kind": "reviewer", "generation": 2, "candidate_oid": "", "candidate_tree": "", "argv": []string{}, "environment": []string{}, "timeout_seconds": 0}
	intentBinding := file(filepath.Join(jobDir, "intent.json"), canonical(intent), 0600)
	job := &taskrun.WorkerJob{ID: "reviewer-0001", Kind: "reviewer", Directory: jobDir, IntentHash: intentBinding.SHA256}
	stdout := file(filepath.Join(jobDir, "stdout.jsonl"), []byte("synthetic completed attempt\n"), 0600)
	stderr := file(filepath.Join(jobDir, "stderr.log"), []byte(""), 0600)
	zero, yes := 0, true
	result := map[string]any{"run_id": id, "job_id": job.ID, "intent_sha256": job.IntentHash, "started_at": stamp(*f.now), "finished_at": stamp(*f.now), "exit": 0, "stdout": stdout, "stderr": stderr, "native_completed": false, "process": taskrun.Process{State: "exited", ExitCode: &zero, Quiescence: &yes}}
	file(filepath.Join(jobDir, "result.json"), canonical(result), 0600)
	x := workspace.PlanWorktreeObservation{WorktreeLocator: root, Ref: "refs/heads/fixture", GitCommonDir: filepath.Join(root, "fixture.git"), OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), Clean: true, StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
	run := taskrun.WorkflowRun{Envelope: taskrun.Envelope{Kind: "ply.workflow.run", SchemaVersion: 2}, RunID: id, SessionID: "ply:" + id, RequestSHA256: taskrun.ControlRevision(r),
		Paths:   taskrun.WorkflowPaths{RunRoot: runRoot, Context: filepath.Join(runRoot, "rounds", "001", "context.json")},
		Handoff: taskrun.WorkflowHandoff{ID: "fixture", Locator: handoff.Locator, SHA256: handoff.SHA256}, Round: taskrun.WorkflowRound{Number: 1, State: "working"},
		Delivery: &taskrun.DeliveryState{Phase: "working", ActualPolicyEvidence: []taskrun.Evidence{}, Events: []taskrun.DeliveryEventRecord{}, Candidates: []taskrun.DeliveryCandidate{}},
		Worker:   &taskrun.WorkerStatus{SessionID: "ply-worker:" + id, Phase: "paused", Stage: "reviewer", Generation: 2, Iteration: 1, Job: job, DeveloperSessions: []string{}}, Reasons: []taskrun.Reason{}}
	state := map[string]any{"request": r, "observed": taskrun.Observed{Target: x, Epic: x, RuntimeBindings: []taskrun.FileBinding{}}, "result": run, "records": []any{}, "phase": "", "context_sha256": "", "acceptance": nil, "start_draft": nil, "start_sha256": nil}
	for _, p := range []string{filepath.Join(root, ".ply", "workflow-runs", "requests", id+".json"), filepath.Join(runRoot, "state.json")} {
		if err := publish(p, state, false); err != nil {
			t.Fatal(err)
		}
	}
	file(filepath.Join(root, ".ply", "projects.lock"), nil, 0600)
	d := taskrun.Dependencies{Workspace: workspace.SystemDependencies(), Now: f.m.Now}
	d.Workflow.Workspace = mr05ChildWorkspace{root: root}
	d.Workflow.Store = mr05ChildHandoff{}
	d.Workspace.IntegrationGit = mr05ChildGit{observation: workspace.IntegrationWorktreeObservation{Locator: root, Ref: x.Ref, GitCommonDir: x.GitCommonDir, OID: x.OID, Tree: x.Tree, Clean: true}}
	d.ControlQueue = func(string, string) (string, error) { return "isolated-child-current-queue", nil }
	f.m.Native = mr05NativeQueue{NativeAdapter: NativeAdapter{Dependencies: d}, queue: f.n, effects: new(int)}
	wakes := new(int)
	f.m.Launch = func(Binding) error { *wakes++; return nil }
	f.s.ActiveTaskID = "first"
	f.s.Entries[0].RunID = id
	f.s.Entries[0].State = "running"
	f.s.Phase, f.s.DesiredAction = "paused", "pause"
	f.n.current = "first"
	f.n.runs["first"] = Run{ID: id, Phase: "paused"}
	g, err := taskrun.ReadParentControl(root, f.b.ControlContract)
	if err != nil {
		t.Fatal(err)
	}
	g.ActiveRun, g.DesiredAction = id, "pause"
	if err = f.m.ControlGuard(root, f.b.ID, func() error { return f.m.saveBound(f.b, f.dir, &f.s, &g) }); err != nil {
		t.Fatal(err)
	}
	observed, err := f.m.Native.Observe(f.c, id)
	if err != nil || observed.Phase != "paused" {
		t.Fatal("actual NativeAdapter could not read seeded child", observed, err)
	}
	return f, id, provider.Locator, jobDir, wakes
}

func mr05NativeChildEvidence(t *testing.T) {
	for _, mutation := range []string{"provider_bytes", "intent.json", "result.json", "stdout.jsonl", "stderr.log", "unchanged"} {
		t.Run(mutation, func(t *testing.T) {
			f, id, provider, job, wakes := mr05NativeChild(t)
			before := mr05ChildInventory(t, f.c.Workspace)
			r := weekPreview(t, f, "resume", "child-proof-resume-0001")
			if !reflect.DeepEqual(before, mr05ChildInventory(t, f.c.Workspace)) {
				t.Fatal("native parent preview wrote to its fixture")
			}
			if r.Input.Target.Kind != "week" {
				t.Fatal("child escaped its parent scope")
			}
			switch mutation {
			case "unchanged":
				weekApply(t, f, r)
				if *wakes != 1 {
					t.Fatal("legitimate parent routing was disabled")
				}
				return
			case "provider_bytes":
				if err := os.WriteFile(provider, []byte("changed executable bytes\n"), 0700); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.Remove(filepath.Join(job, mutation)); err != nil {
					t.Fatal(err)
				}
			}
			b, err := f.m.Controls().Show(f.c.Workspace, f.b.ID)
			if err != nil {
				t.Fatal(err)
			}
			if len(b.Actions) != 0 {
				t.Errorf("unsafe native child still offers parent control after %s: %v", mutation, b.Actions)
			}
			reader := f.m.Native.(childControlReader)
			child, err := reader.BoundControlBasis(f.c, id)
			if err != nil {
				t.Fatal(err)
			}
			if child.UnavailableReason == nil || *child.UnavailableReason == "parent_control_unsupported" {
				t.Errorf("child routing concealed missing evidence after %s", mutation)
			}
			status, err := f.m.Controls().Apply(r)
			if err == nil || status.RequestState != "rejected" || status.ControlEffect != "none" || *wakes != 0 {
				t.Errorf("native child drift reached reservation/wake: mutation=%s state=%s effect=%s wakes=%d err=%v", mutation, status.RequestState, status.ControlEffect, *wakes, err)
			}
			if _, err := os.Stat(filepath.Join(f.c.Workspace, ".ply", "runtime-controls", "v1", r.RequestID, "request.json")); !os.IsNotExist(err) {
				t.Error("invalid child proof acquired durable reservation")
			}
			if !reflect.DeepEqual(r.Input.Target, taskrun.ControlTarget{Kind: "week", ID: f.b.ID}) {
				t.Fatal("request target changed")
			}
		})
	}
}

func mr05NativeRegistryContention(t *testing.T) {
	f, _, _, _, wakes := mr05NativeChild(t)
	root := f.c.Workspace
	if err := os.WriteFile(filepath.Join(root, ".ply", "work-items.lock"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	r := weekPreview(t, f, "resume", "native-registry-0001")
	validated, releaseControl := make(chan struct{}), make(chan struct{})
	registryHeld, releaseRegistry := make(chan struct{}), make(chan struct{})
	var controlOnce, registryOnce sync.Once
	defer controlOnce.Do(func() { close(releaseControl) })
	defer registryOnce.Do(func() { close(releaseRegistry) })
	f.m.Fault = func(point string) error {
		if point == "control_before_reserve" {
			close(validated)
			<-releaseControl
		}
		return nil
	}
	applyDone := make(chan error, 1)
	go func() { _, err := f.m.Controls().Apply(r); applyDone <- err }()
	select {
	case <-validated:
	case err := <-applyDone:
		t.Fatal("apply did not validate under its native guards", err)
	case <-time.After(3 * time.Second):
		t.Fatal("apply did not reach its publication barrier")
	}
	nativeDone := make(chan error, 1)
	d := f.m.Native.(mr05NativeQueue).Dependencies
	go func() {
		nativeDone <- workspace.WithTaskSpecRegistrationSnapshot(d.Workspace, root, func(*workspace.TaskSpecSession) error {
			close(registryHeld)
			<-releaseRegistry
			return nil
		})
	}()
	select {
	case <-registryHeld:
	case err := <-nativeDone:
		t.Fatal("native registry guard failed", err)
	case <-time.After(3 * time.Second):
		t.Fatal("native registry could not acquire its own guards")
	}
	for _, path := range []string{filepath.Join(f.dir, "control.lock"), filepath.Join(root, ".ply", "task-runs", "v1", "store.lock"), filepath.Join(root, ".ply", "projects.lock"), filepath.Join(root, ".ply", "work-items.lock")} {
		probe, err := lock(path)
		if probe != nil {
			probe.Close()
		}
		if !errors.Is(err, errOwned) {
			t.Fatal("expected native lock is not held", path, err)
		}
	}
	controlOnce.Do(func() { close(releaseControl) })
	select {
	case err := <-applyDone:
		t.Fatal("apply escaped the contended registry publication guard", err)
	case <-time.After(150 * time.Millisecond):
	}
	if _, err := os.Stat(filepath.Join(root, ".ply", "runtime-controls", "v1", r.RequestID, "request.json")); !os.IsNotExist(err) {
		t.Fatal("control published before native registry revalidation", err)
	}
	registryOnce.Do(func() { close(releaseRegistry) })
	if err := <-nativeDone; err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-applyDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("parent/TaskRun/Project/WorkItem publication deadlocked")
	}
	if *wakes != 1 || *f.m.Native.(mr05NativeQueue).effects != 0 {
		t.Fatal("contention changed the authorized effect count")
	}
}

func mr05ChildInventory(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if !entry.IsDir() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += hash(b)
		}
		out[path] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// The native completed result is a disposable retained-store fixture, read by
// the real NativeAdapter. Parent reservation/consumption and dispatch exclusion
// are production paths; this is not a claim of live integration acceptance.
func mr05CompletedChildPause(t *testing.T) {
	f, id, _, _, wakes := mr05NativeChild(t)
	root := f.c.Workspace
	statePath := filepath.Join(root, ".ply", "workflow-runs", "runs", id, "state.json")
	updateChild := func(completed bool) {
		t.Helper()
		if err := taskrun.WithControlStore(root, func() error {
			var state map[string]any
			if err := readJSON(statePath, &state); err != nil {
				return err
			}
			result := state["result"].(map[string]any)
			phase := "running"
			if completed {
				phase = "completed"
			}
			result["worker"].(map[string]any)["phase"] = phase
			if completed {
				result["delivery"].(map[string]any)["phase"] = "completed"
			}
			return publish(statePath, state, true)
		}); err != nil {
			t.Fatal(err)
		}
	}
	updateChild(false)
	g, err := taskrun.ReadParentControl(root, f.b.ControlContract)
	if err != nil {
		t.Fatal(err)
	}
	f.s.Phase, f.s.DesiredAction = "running", "continue"
	g.DesiredAction = "continue"
	if err = f.m.ControlGuard(root, f.b.ID, func() error { return f.m.saveBound(f.b, f.dir, &f.s, &g) }); err != nil {
		t.Fatal(err)
	}
	r := weekPreview(t, f, "pause", "between-goals-pause-0001")
	weekApply(t, f, r)
	weekEffect(t, f, r.RequestID, "pending_boundary")
	updateChild(true)
	f.tick(t)
	status := weekEffect(t, f, r.RequestID, "paused")
	if status.Basis.RunID == nil || *status.Basis.RunID != id {
		t.Fatal("receipt lost the originally approved child")
	}
	for i := 0; i < 3; i++ {
		basis, err := f.m.Controls().Show(root, f.b.ID)
		if err != nil || basis.Phase != "paused" || basis.RunID != nil || basis.TaskID != nil || basis.Week == nil || basis.Week.DesiredAction != "pause" {
			t.Fatal("winning completion pause did not publish a null-Run parked boundary", basis, err)
		}
		f.tick(t)
	}
	if f.s.Entries[0].State != "completed" || f.s.Entries[1].State != "pending" || f.s.ActiveTaskID != "" || *wakes != 1 || *f.m.Native.(mr05NativeQueue).effects != 0 {
		t.Fatal("next Goal escaped the completion pause", f.s)
	}
}
