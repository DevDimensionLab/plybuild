//go:build !windows

package taskrun

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const aw01ResumeMessage = "Work is preserved; use worker resume for this same Run."

// These fixtures preserve real request/state/attempt files and use the public
// readers. Native observations are isolated adapters, never a registry, Git
// checkout, workspace.Init, provider, or another Run.
type aw01Fixture struct {
	t   *testing.T
	d   Dependencies
	s   workflowState
	now time.Time
}

type aw01Workspace struct {
	workflowhandoff.WorkspaceObserver
	root string
}

func (w aw01Workspace) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	return workflowhandoff.WorkspaceSnapshot{Observation: workspace.WorkspaceObservation{Root: w.root}}, nil
}

type aw01HandoffStore struct{ workflowhandoff.Store }

func (aw01HandoffStore) ReadByLocator(string) (workflowhandoff.Snapshot, error) {
	return workflowhandoff.Snapshot{}, nil
}

type aw01Git struct {
	workspace.TaskIntegrationGit
	x workspace.IntegrationWorktreeObservation
}

func (g aw01Git) ObserveIntegrationWorktree(string, string) (workspace.IntegrationWorktreeObservation, error) {
	return g.x, nil
}

func newAW01Fixture(t *testing.T) *aw01Fixture {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := &aw01Fixture{t: t, now: time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC)}
	f.d.Now = func() time.Time { return f.now }
	f.d.Workflow.Workspace = aw01Workspace{root: root}
	f.d.Workflow.Store = aw01HandoffStore{}
	provider := f.file(filepath.Join(root, "provider"), []byte("#!/bin/sh\nexit 99\n"))
	control := f.file(filepath.Join(root, "ply-control"), []byte("#!/bin/sh\nexit 98\n"))
	for _, path := range []string{provider.Locator, control.Locator} {
		if err := os.Chmod(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	goal := f.file(filepath.Join(root, "goal.json"), []byte(`{"fixture":"AW01"}`))
	handoff := f.file(filepath.Join(root, "handoff.json"), []byte(`{"fixture":"isolated readback"}`))
	r := WorkflowRequest{
		Envelope: deliveryEnv("headless-run-request"), RequestKey: "aw01-guidance", WorkspaceRoot: root,
		HandoffDraft: []byte(`{"inputs":[]}`), ReturnMode: "delivery_owner",
		Runtime:  Runtime{Provider: "codex", Mode: "headless", Model: "fixture-model", Executable: Executable{provider.Locator, provider.SHA256}, PlyExecutable: Executable{control.Locator, control.SHA256}, PermissionBinding: Permission{Evidence: []Evidence{}}},
		Worker:   &WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 120},
		Delivery: &DeliveryContract{Goal: goal, OwnerClaim: "fixture owner", ReasoningEffort: "max"},
	}
	id := workflowID(r)
	x := workspace.PlanWorktreeObservation{WorktreeLocator: root, Ref: "refs/heads/fixture", GitCommonDir: filepath.Join(root, "fixture.git"), OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), Clean: true, StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
	f.d.Workspace.IntegrationGit = aw01Git{x: workspace.IntegrationWorktreeObservation{Locator: root, Ref: x.Ref, GitCommonDir: x.GitCommonDir, OID: x.OID, Tree: x.Tree, Clean: true}}
	f.s = workflowState{Request: r, Observed: Observed{Target: x, Epic: x, RuntimeBindings: []FileBinding{}}, Records: []workflowRecord{}, Result: WorkflowRun{
		Envelope: deliveryEnv("run"), RunID: id, SessionID: "ply:" + id, RequestSHA256: digest(r), Paths: workflowPaths(r, 1),
		Handoff: WorkflowHandoff{ID: "fixture", Locator: handoff.Locator, SHA256: handoff.SHA256},
		Round:   WorkflowRound{Number: 1, State: "working"}, Reasons: []Reason{}, Delivery: &DeliveryState{Phase: "working", ActualPolicyEvidence: []Evidence{}, Events: []DeliveryEventRecord{}, Candidates: []DeliveryCandidate{}},
		Worker:     &WorkerStatus{SessionID: "ply-worker:" + id, Phase: "running", Stage: "reviewer", Generation: 2, Iteration: 1, OwnerPID: 12345, DeveloperSessions: []string{}},
		NextAction: WorkflowAction{"human", aw01ResumeMessage},
	}}
	f.job("reviewer")
	f.save()
	return f
}

func (f *aw01Fixture) file(path string, raw []byte) FileBinding {
	f.t.Helper()
	if err := (Dependencies{}).writeOnce(path, raw); err != nil {
		f.t.Fatal(err)
	}
	return FileBinding{Locator: path, SHA256: hash(raw)}
}

func (f *aw01Fixture) job(kind string) {
	f.t.Helper()
	id := kind + "-0001"
	dir := filepath.Join(workerDirectory(f.s), "attempts", id)
	in, err := workflowKeep(Dependencies{}, filepath.Join(dir, "intent.json"), workerAttemptIntent{RunID: f.s.Result.RunID, JobID: id, Kind: kind, Generation: f.s.Result.Worker.Generation, Argv: []string{}, Environment: []string{}})
	if err != nil {
		f.t.Fatal(err)
	}
	f.s.Result.Worker.Stage = kind
	f.s.Result.Worker.Job = &WorkerJob{ID: id, Kind: kind, Directory: dir, IntentHash: in.SHA256}
}

func (f *aw01Fixture) save() {
	f.t.Helper()
	if err := (Dependencies{}).replaceValue(workflowIndex(f.s.Request.WorkspaceRoot, f.s.Result.RunID), f.s); err != nil {
		f.t.Fatal(err)
	}
	if err := workflowSave(Dependencies{}, f.s); err != nil {
		f.t.Fatal(err)
	}
}

func (f *aw01Fixture) reads(check func(*testing.T, WorkflowRun)) {
	f.t.Helper()
	for _, read := range []struct {
		name string
		fn   func(Dependencies, string, string) (WorkflowRun, error)
	}{{"WorkflowShow", WorkflowShow}, {"ShowWorker", ShowWorker}} {
		f.t.Run(read.name, func(t *testing.T) {
			o, err := read.fn(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if o.RunID != f.s.Result.RunID || o.RequestSHA256 != f.s.Result.RequestSHA256 || o.SessionID != f.s.Result.SessionID || !equal(o.Worker.Job, f.s.Result.Worker.Job) || o.Worker.Generation != f.s.Result.Worker.Generation {
				t.Fatalf("read changed Run, binding or attempt: %+v", o)
			}
			check(t, o)
		})
	}
}

func aw01AssertObserving(t *testing.T, o WorkflowRun) {
	t.Helper()
	if len(o.Reasons) != 0 || o.Worker.Reason != "" {
		t.Fatalf("fixture has unrelated obstacle: %+v / %+v", o.Reasons, o.Worker)
	}
	if o.NextAction.Actor != "worker" || !strings.Contains(o.NextAction.Message, "Observe") || !strings.Contains(o.NextAction.Message, "same Run") || strings.Contains(o.NextAction.Message, "use worker resume") || strings.Contains(o.NextAction.Message, "waits until") {
		t.Fatalf("stale active worker guidance: %+v", o.NextAction)
	}
	if !strings.Contains(o.NextAction.Message, "not") || !strings.Contains(o.NextAction.Message, "liveness") {
		t.Fatalf("cached worker phase must not attest liveness: %+v", o.NextAction)
	}
}

func TestAW01GuidanceLegacyRunningReadback(t *testing.T) {
	for _, stage := range []string{"developer", "reviewer"} {
		for _, action := range []WorkflowAction{{"human", aw01ResumeMessage}, {"worker_wait", "usage_limit: same Run waits until 2026-10-10T05:00:00Z"}} {
			t.Run(stage+"/"+action.Actor, func(t *testing.T) {
				f := newAW01Fixture(t)
				f.job(stage)
				f.s.Result.NextAction = action
				f.save()
				f.reads(aw01AssertObserving)
			})
		}
	}
}

func (f *aw01Fixture) reload() {
	f.t.Helper()
	var err error
	f.s, err = workflowRead(f.s.Request.WorkspaceRoot, f.s.Result.RunID)
	if err != nil {
		f.t.Fatal(err)
	}
}

func (f *aw01Fixture) own() {
	f.t.Helper()
	lock, err := workerLock(filepath.Join(workerDirectory(f.s), "owner.lock"), true)
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { lock.Close() })
	f.d.WorkerOwner = &workerOwner{root: f.s.Request.WorkspaceRoot, runID: f.s.Result.RunID, generation: f.s.Result.Worker.Generation, lock: lock}
}

func (f *aw01Fixture) result(process *Process, failure *WorkerProviderFailure) workerAttemptResult {
	f.t.Helper()
	job := f.s.Result.Worker.Job
	raw := []byte("synthetic completed attempt\n")
	exit := 0
	if failure != nil {
		exit = 1
		raw = []byte(fmt.Sprintf("{\"type\":\"error\",\"code\":\"usage_limit_reached\",\"resets_at\":%q}\n", failure.ResetAt))
	}
	stream := f.file(filepath.Join(job.Directory, "stdout.jsonl"), raw)
	r := workerAttemptResult{RunID: f.s.Result.RunID, JobID: job.ID, IntentSHA256: job.IntentHash, FinishedAt: f.now.Format(time.RFC3339Nano), Exit: ptr(exit), Stdout: stream, Stderr: stream, Process: process, ProviderFailure: failure}
	if err := (Dependencies{}).replaceValue(filepath.Join(job.Directory, "result.json"), r); err != nil {
		f.t.Fatal(err)
	}
	return r
}

func (f *aw01Fixture) control(action string) workerControl {
	f.t.Helper()
	if err := workerControlWrite(f.d, f.s, action); err != nil {
		f.t.Fatal(err)
	}
	return workerControl{RunID: f.s.Result.RunID, Action: action, RequestedAt: f.now.UTC().Format(time.RFC3339Nano)}
}

func (f *aw01Fixture) park(action string) {
	f.t.Helper()
	if f.d.WorkerOwner != nil {
		f.d.WorkerOwner.lock.Close()
		f.d.WorkerOwner = nil
	}
	f.control(action)
	// ServeWorker consumes the real pause/stop at a positively completed attempt.
	// It returns before any preparation lookup, provider, or native delivery.
	f.t.Chdir(f.s.Request.WorkspaceRoot)
	if err := ServeWorker(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID); err != nil {
		f.t.Fatal(err)
	}
	f.reload()
}

func TestAW01GuidanceResumeTransition(t *testing.T) {
	f := newAW01Fixture(t)
	f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
	originalRequest, originalJob := digest(f.s.Request), *f.s.Result.Worker.Job
	f.park("pause")
	if f.s.Result.Worker.Phase != "paused" || f.s.Result.NextAction.Message != aw01ResumeMessage {
		t.Fatal("actual pause was not consumed")
	}
	f.own()
	control := f.control("resume")
	stale := f.d
	wrongOwner := *f.d.WorkerOwner
	wrongOwner.generation--
	stale.WorkerOwner = &wrongOwner
	if err := workerConsumeResume(stale, f.s, control); err == nil {
		t.Fatal("stale owner consumed resume")
	}
	if err := workerConsumeResume(f.d, f.s, control); err != nil {
		t.Fatal(err)
	}
	f.reload()
	var consumed workerControl
	if err := readValue(filepath.Join(workerDirectory(f.s), "control.json"), 64<<10, &consumed); err != nil || consumed.Action != "continue" || consumed.RunID != control.RunID || consumed.RequestedAt != control.RequestedAt {
		t.Fatalf("resume identity not consumed: %+v / %v", consumed, err)
	}
	if f.s.Result.Worker.Phase != "running" || f.s.Result.Worker.Stage != "reviewer" || *f.s.Result.Worker.Job != originalJob || digest(f.s.Request) != originalRequest {
		t.Fatal("resume changed the binding or reserved reviewer")
	}
	f.reads(aw01AssertObserving)
}

func TestAW01GuidanceWaitingResume(t *testing.T) {
	for _, completion := range []string{"confirmed", "missing", "unknown"} {
		t.Run(completion, func(t *testing.T) {
			f := newAW01Fixture(t)
			deadline := f.now.Add(time.Hour).Format(time.RFC3339Nano)
			r := f.result(&Process{State: "exited", ExitCode: ptr(1), Quiescence: ptr(true)}, &WorkerProviderFailure{Kind: "usage_limit", ResetAt: deadline})
			f.own()
			if err := workerWaitForProvider(f.d, f.s, r); err != nil {
				t.Fatal(err)
			}
			f.reload()
			f.park("pause")
			f.own()
			if err := workerConsumeResume(f.d, f.s, f.control("resume")); err != nil {
				t.Fatal(err)
			}
			f.reload()
			assertWait := func(t *testing.T, o WorkflowRun) {
				if o.Worker.Phase != "waiting" || o.Worker.RetryNotBefore != deadline || o.Worker.ProviderRetries != 1 || o.NextAction.Actor != "worker_wait" || !strings.Contains(o.NextAction.Message, deadline) || strings.Contains(o.NextAction.Message, "use worker resume") {
					t.Fatalf("saved wait not respected: %+v / %+v", o.NextAction, o.Worker)
				}
			}
			f.reads(assertWait)
			if err := workerContinueAfterWait(f.d, f.s); err == nil {
				t.Fatal("resume shortened the saved reset")
			}
			path := filepath.Join(f.s.Result.Worker.Job.Directory, "result.json")
			switch completion {
			case "missing":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "unknown":
				r.Process.Quiescence = ptr(false)
				if err := (Dependencies{}).replaceValue(path, r); err != nil {
					t.Fatal(err)
				}
			}
			f.now = f.now.Add(time.Hour)
			// Reading the elapsed timestamp alone never consumes the wait.
			f.reads(assertWait)
			err := workerContinueAfterWait(f.d, f.s)
			if completion != "confirmed" {
				if err == nil {
					t.Fatal("uncertain attempt was released")
				}
				if blockErr := workerBlocked(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID, err); blockErr != nil {
					t.Fatal(blockErr)
				}
				f.reload()
				f.reads(func(t *testing.T, o WorkflowRun) {
					if o.Worker.Phase != "blocked" || o.Worker.RetryNotBefore != deadline || !strings.Contains(o.NextAction.Message, err.Error()) {
						t.Fatalf("unknown completion concealed: %+v", o)
					}
				})
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			f.reload()
			if f.s.Result.Worker.RetryNotBefore != "" || f.s.Result.Worker.Job != nil || f.s.Result.Worker.Stage != "reviewer" {
				t.Fatal("confirmed expiry did not preserve stage and release attempt")
			}
			f.reads(aw01AssertObserving)
			if _, err := workflowBound(r.Stdout, providerStreamLimit); err != nil {
				t.Fatal("wait consumption lost attempt history:", err)
			}
		})
	}
}

func TestAW01GuidancePauseStopAndBlocked(t *testing.T) {
	for _, phase := range []string{"paused", "stopped", "blocked", "running", "unknown"} {
		t.Run(phase, func(t *testing.T) {
			f := newAW01Fixture(t)
			f.s.Result.Worker.Phase = phase
			if phase != "paused" && phase != "stopped" {
				f.s.Result.Worker.Reason = "attempt reviewer-0001 has no positive process completion"
				f.s.Result.Reasons = []Reason{{"worker_attempt_unknown", "ownership and process outcome are unknown"}}
				f.s.Result.NextAction = WorkflowAction{"human", "Inspect the preserved unknown attempt and ownership; no replacement writer."}
			}
			f.save()
			f.reads(func(t *testing.T, o WorkflowRun) {
				if o.NextAction != f.s.Result.NextAction || o.Worker.Phase != phase || o.Worker.Reason != f.s.Result.Worker.Reason || !equal(o.Reasons, f.s.Result.Reasons) {
					t.Fatalf("real obstacle hidden: %+v", o)
				}
			})
		})
	}
}

func TestAW01GuidancePendingControls(t *testing.T) {
	for _, action := range []string{"pause", "stop", "resume"} {
		t.Run(action, func(t *testing.T) {
			f := newAW01Fixture(t)
			if action == "resume" {
				f.s.Result.Worker.Phase = "paused"
				f.save()
			}
			control := f.control(action)
			f.reads(func(t *testing.T, o WorkflowRun) {
				if o.Worker.ControlAction != action || o.Worker.ControlRequestedAt != control.RequestedAt || o.Worker.Phase != f.s.Result.Worker.Phase {
					t.Fatalf("pending control identity or unconsumed phase lost: %+v", o.Worker)
				}
				if o.NextAction.Actor != "worker" || !strings.Contains(o.NextAction.Message, action+" requested") || !strings.Contains(o.NextAction.Message, "safe boundary") || strings.Contains(o.NextAction.Message, "use worker resume") {
					t.Fatalf("pending control treated as completed or repeated: %+v", o.NextAction)
				}
			})
			var preserved workerControl
			if err := readValue(filepath.Join(workerDirectory(f.s), "control.json"), 64<<10, &preserved); err != nil || preserved != control {
				t.Fatalf("read consumed pending control: %+v / %v", preserved, err)
			}
		})
	}
	for _, controlState := range []string{"other_run", "malformed", "consumed"} {
		t.Run(controlState, func(t *testing.T) {
			f := newAW01Fixture(t)
			f.s.Result.Worker.ControlAction = "resume"
			f.s.Result.Worker.ControlRequestedAt = "old observation"
			f.save()
			f.control("continue")
			path := filepath.Join(workerDirectory(f.s), "control.json")
			if controlState == "other_run" {
				if err := (Dependencies{}).replaceValue(path, workerControl{RunID: "another-run", Action: "resume", RequestedAt: f.now.Format(time.RFC3339Nano)}); err != nil {
					t.Fatal(err)
				}
			} else if controlState == "malformed" {
				if err := os.WriteFile(path, []byte("{broken"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
			if controlState == "consumed" {
				f.reads(func(t *testing.T, o WorkflowRun) {
					aw01AssertObserving(t, o)
					if o.Worker.ControlAction != "" || o.Worker.ControlRequestedAt != "" {
						t.Fatal("consumed control retained stale pending fields")
					}
				})
			} else {
				for _, read := range []func(Dependencies, string, string) (WorkflowRun, error){WorkflowShow, ShowWorker} {
					if _, err := read(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID); err == nil {
						t.Fatal("unknown or foreign control was accepted")
					}
				}
			}
			if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
				t.Fatal("control read changed preserved files")
			}
		})
	}
}

func (f *aw01Fixture) closeout(state string) workspace.TaskCloseoutReceipt {
	f.t.Helper()
	x := f.s.Observed.Target
	verification := f.file(filepath.Join(f.s.Result.Paths.RunRoot, "verification.json"), []byte(`{"fixture":"verification binding"}`))
	c := DeliveryCandidate{OID: x.OID, Tree: x.Tree, Verification: verification, TaskResult: workspace.TaskResultRecord{
		ID: "result-fixture", TaskID: "aw01-fixture", SourceLocator: x.WorktreeLocator, SourceRef: x.Ref,
		VerifierResults: []workspace.TaskVerifierResultRecord{}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{}, Artifacts: []workspace.TaskArtifactRecord{},
		Review: workspace.TaskReviewRecord{Findings: []workspace.TaskReviewEntry{}, Fixes: []workspace.TaskReviewEntry{}, OpenActionableFindings: []workspace.TaskReviewEntry{}},
	}}
	f.s.Result.Delivery.Candidates = []DeliveryCandidate{c}
	p := workspace.TaskCloseoutPlan{WorkspaceRoot: f.s.Request.WorkspaceRoot, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, Keep: true,
		Source: workspace.TaskWorktreeBinding{Locator: x.WorktreeLocator, Ref: x.Ref, GitCommonDir: x.GitCommonDir}}
	raw, err := json.Marshal(p)
	if err != nil {
		f.t.Fatal(err)
	}
	r := workspace.TaskCloseoutReceipt{Kind: "WorkspaceTaskCloseoutReceipt@1", SchemaVersion: 1, OperationID: "isolated-closeout", Plan: p, PlanSHA256: hash(raw), State: state, ResourceState: "retained", Phase: "cleanup_pending", CreatedAtUTC: f.now.Format(time.RFC3339Nano), UpdatedAtUTC: f.now.Format(time.RFC3339Nano)}
	if state == "complete" {
		r.LifecycleCompleted, r.ResourceState, r.Phase = true, "kept", "complete"
		r.Retention = &workspace.TaskRetentionManifest{Entries: []workspace.TaskRetainedEvidence{}}
	}
	store := map[string]any{"kind": "WorkspaceTaskCloseoutStore@1", "schema_version": 1, "receipts": []workspace.TaskCloseoutReceipt{r}}
	if err := (Dependencies{}).writeValue(filepath.Join(f.s.Request.WorkspaceRoot, ".ply", "task-closeouts.json"), store); err != nil {
		f.t.Fatal(err)
	}
	f.save()
	return r
}

func TestAW01GuidanceCloseoutAndDriftPriority(t *testing.T) {
	for _, action := range []WorkflowAction{{"human", aw01ResumeMessage}, {"worker_wait", "usage_limit: same Run waits until 2026-10-10T05:00:00Z"}} {
		for _, scenario := range []string{"complete", "cleanup_pending", "closeout_unavailable", "binding_drift", "missing_evidence"} {
			t.Run(action.Actor+"/"+scenario, func(t *testing.T) {
				f := newAW01Fixture(t)
				f.s.Result.NextAction = action
				var receipt workspace.TaskCloseoutReceipt
				switch scenario {
				case "complete", "cleanup_pending", "closeout_unavailable":
					state := scenario
					if state == "closeout_unavailable" {
						state = "cleanup_pending"
					}
					receipt = f.closeout(state)
					if scenario == "closeout_unavailable" {
						if err := os.WriteFile(filepath.Join(f.s.Request.WorkspaceRoot, ".ply", "task-closeouts.json"), []byte("{broken"), 0600); err != nil {
							t.Fatal(err)
						}
					}
				case "binding_drift":
					if err := os.WriteFile(f.s.Request.Runtime.Executable.Path, []byte("changed executable"), 0700); err != nil {
						t.Fatal(err)
					}
				case "missing_evidence":
					if err := os.Remove(f.s.Request.Delivery.Goal.Locator); err != nil {
						t.Fatal(err)
					}
				}
				f.save()
				f.reads(func(t *testing.T, o WorkflowRun) {
					switch scenario {
					case "complete":
						if o.Closeout == nil || o.Closeout.State != "complete" || o.Round.State != "completed" || o.FinalReturn.State != "completed" || o.NextAction.Actor != "user" || !strings.Contains(o.NextAction.Message, "exact Task delivery is closed") {
							t.Fatalf("complete closeout lost priority: %+v", o)
						}
					case "cleanup_pending":
						code, _, message := workspace.TaskCloseoutPendingAction(receipt)
						if o.Closeout == nil || o.NextAction != (WorkflowAction{"user", message}) || !aw01HasReason(o, code) {
							t.Fatalf("pending closeout lost priority: %+v", o)
						}
					case "closeout_unavailable":
						if !aw01HasReason(o, "task_closeout_unavailable") || o.NextAction != action {
							t.Fatalf("unavailable closeout replaced with active guidance: %+v / %+v", o.Reasons, o.NextAction)
						}
					default:
						if !aw01HasReason(o, "workflow_run_drift") || o.Round.State != "unknown" || o.NextAction.Actor != "coordinator" || !strings.Contains(o.NextAction.Message, "binding drift") {
							t.Fatalf("required evidence failure lost priority: %+v", o)
						}
					}
				})
			})
		}
	}
}

func aw01HasReason(o WorkflowRun, code string) bool {
	for _, reason := range o.Reasons {
		if reason.Code == code {
			return true
		}
	}
	return false
}

type aw01NoEffects struct{ t *testing.T }

func (n aw01NoEffects) WriteOnce(string, []byte) error { n.t.Fatal("read wrote a file"); return nil }
func (n aw01NoEffects) Replace(string, []byte) error   { n.t.Fatal("read replaced a file"); return nil }
func (n aw01NoEffects) Check() error                   { n.t.Fatal("read checked a process runner"); return nil }
func (n aw01NoEffects) Run(LaunchSpec, func(Process) error) (Process, error) {
	n.t.Fatal("read started a provider or native effect")
	return Process{}, nil
}

func aw01Inventory(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
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
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += "\n" + string(raw)
		}
		files[path] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestAW01GuidanceReadDoesNotWrite(t *testing.T) {
	for _, owner := range []string{"absent", "held", "unowned", "unknown"} {
		t.Run(owner, func(t *testing.T) {
			f := newAW01Fixture(t)
			f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
			f.file(filepath.Join(f.s.Result.Paths.RunRoot, "history.jsonl"), []byte("preserved paused/resumed history\n"))
			f.control("pause")
			path := filepath.Join(workerDirectory(f.s), "owner.lock")
			if owner == "held" || owner == "unowned" {
				f.own()
				if owner == "unowned" {
					f.d.WorkerOwner.lock.Close()
				}
			} else if owner == "unknown" {
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			}
			no := aw01NoEffects{t}
			f.d.Files, f.d.Runner, f.d.ExecRunner = no, no, no
			f.d.Fault = func(string) error { t.Fatal("read entered an effect boundary"); return nil }
			f.d.StartupProcessTree = func(int) ([]StartupProcess, error) { t.Fatal("read inferred provider liveness"); return nil, nil }
			before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
			for i := 0; i < 3; i++ {
				f.reads(func(t *testing.T, o WorkflowRun) {
					if o.Worker.ControlAction != "pause" || o.Worker.Phase != "running" || strings.Contains(o.NextAction.Message, "use worker resume") {
						t.Fatalf("read changed control semantics: %+v", o)
					}
				})
				shown, err := ShowWorker(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
				want := map[string]string{"absent": "not_started", "held": "lock_held", "unowned": "unowned", "unknown": "unknown"}[owner]
				if err != nil || shown.Worker.OwnerState != want {
					t.Fatalf("ownership observation changed: %+v / %v", shown.Worker, err)
				}
			}
			if after := aw01Inventory(t, f.s.Request.WorkspaceRoot); !reflect.DeepEqual(before, after) {
				t.Fatal("read changed actual inventory, bytes or permissions")
			}
			if owner == "absent" {
				if _, err := os.Lstat(path); !os.IsNotExist(err) {
					t.Fatal("reading created an absent owner lock")
				}
			}
		})
	}
}

func TestAW01GuidanceInteractiveUnchanged(t *testing.T) {
	for _, scenario := range []string{"current", "legacy_resume", "legacy_wait", "unbound", "startup", "drift"} {
		t.Run(scenario, func(t *testing.T) {
			f := newAW01Fixture(t)
			f.s.Request.Envelope = deliveryEnv("herdr-run-request")
			f.s.Request.Worker, f.s.Result.Worker = nil, nil
			f.s.Request.Runtime.Mode = "interactive"
			f.s.Request.Herdr.Executable = f.s.Request.Runtime.Executable
			f.s.Result.Transport.AgentSessionID = "fixture-interactive-session"
			f.s.Result.NextAction = WorkflowAction{"recipient", "Continue the same interactive Task."}
			switch scenario {
			case "legacy_resume":
				f.s.Result.NextAction = WorkflowAction{"human", aw01ResumeMessage}
			case "legacy_wait":
				f.s.Result.NextAction = WorkflowAction{"worker_wait", "saved interactive message"}
			case "unbound", "startup":
				f.s.Result.Transport.AgentSessionID = ""
				if scenario == "startup" {
					f.s.Phase = "agent_start_attempted"
					f.s.Result.Round.Number = 0
					f.s.Result.Transport.TabID, f.s.Result.Transport.PaneID, f.s.Result.Transport.TerminalID = "tab", "pane", "terminal"
				}
			case "drift":
				if err := os.Remove(f.s.Request.Delivery.Goal.Locator); err != nil {
					t.Fatal(err)
				}
			}
			f.s.Result.RequestSHA256 = digest(f.s.Request)
			f.save()
			before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
			o, err := WorkflowShow(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if o.Worker != nil || o.RunID != f.s.Result.RunID || o.Transport.AgentSessionID != f.s.Result.Transport.AgentSessionID {
				t.Fatal("interactive identity changed or acquired worker guidance")
			}
			switch scenario {
			case "unbound":
				if o.NextAction.Actor != "coordinator" || !strings.Contains(o.NextAction.Message, "no native session is bound") {
					t.Fatal("unbound startup priority changed:", o.NextAction)
				}
			case "startup":
				if o.NextAction != deliveryStartupAction(f.s) {
					t.Fatal("startup priority changed:", o.NextAction)
				}
			case "drift":
				if !aw01HasReason(o, "workflow_run_drift") || o.NextAction.Actor != "coordinator" {
					t.Fatal("interactive drift priority changed:", o.NextAction)
				}
			default:
				if o.NextAction != f.s.Result.NextAction || len(o.Reasons) != 0 {
					t.Fatal("interactive guidance changed:", o.NextAction, o.Reasons)
				}
			}
			if _, err := ShowWorker(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID); err == nil {
				t.Fatal("interactive Run was accepted as a worker")
			}
			if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
				t.Fatal("interactive read wrote state")
			}
		})
	}
}
