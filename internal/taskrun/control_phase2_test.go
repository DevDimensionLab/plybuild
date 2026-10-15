//go:build !windows

package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func mr05PositiveResult(f *aw01Fixture) {
	zero, yes := 0, true
	r := f.result(&Process{State: "exited", ExitCode: &zero, Quiescence: &yes}, nil)
	r.Stderr = f.file(filepath.Join(f.s.Result.Worker.Job.Directory, "stderr.log"), nil)
	if e := f.d.replaceValue(filepath.Join(f.s.Result.Worker.Job.Directory, "result.json"), r); e != nil {
		f.t.Fatal(e)
	}
}

func mr05WorkerReceiptRecovery(t *testing.T) {
	f, wakes := mr05Worker(t)
	r := mr05Preview(t, f, "pause", "worker-receipt-recovery-0001")
	if _, e := WorkerControls(f.d).Apply(r); e != nil {
		t.Fatal(e)
	}
	f.reload()
	f.own()
	f.d.Fault = func(point string) error {
		if point == "control_before_receipt" {
			return errors.New("receipt publication interrupted")
		}
		return nil
	}
	if _, e := workerConsumeBoundControl(f.d, f.s); e == nil {
		t.Fatal("publication fault missed")
	}
	f.reload()
	if f.s.ControlReceipt == nil {
		t.Fatal("no durable owner receipt")
	}
	retained := *f.s.ControlReceipt
	f.d.Fault = nil
	if e := workerClaimOwner(&f.d, f.s, f.d.WorkerOwner.lock); e != nil {
		t.Fatal(e)
	}
	mr05Consume(t, f)
	status := mr05Status(t, f, r.RequestID, "paused")
	if !equal(*status.Receipt, retained) || wakes.Load() != 1 {
		t.Error("receipt recovery changed original evidence or woke twice")
	}
}

// The observer seams only bound and step the production supervisor loop. Its
// actual ServeWorker consumes the request; no controller/provider subprocess is
// needed and the retry clock remains fixed throughout.
func mr05ExistingWaitingSupervisor(t *testing.T) {
	for _, action := range []string{"pause", "stop"} {
		t.Run(action, func(t *testing.T) {
			f, wakes := mr05Worker(t)
			f.s.Result.Worker.Phase = "waiting"
			due := ControlTimestamp(f.now.Add(time.Hour))
			f.s.Result.Worker.RetryNotBefore = due
			f.save()
			r := mr05Preview(t, f, action, "waiting-supervisor-0001")
			waiting, served := 0, 0
			done := errors.New("bounded supervisor observation")
			t.Chdir(f.s.Request.WorkspaceRoot)
			f.d.Fault = func(point string) error {
				switch point {
				case "worker_supervisor_waiting":
					waiting++
					if waiting > 1 {
						return done
					}
					if _, e := WorkerControls(f.d).Apply(r); e != nil {
						t.Fatal(e)
					}
				case "worker_supervisor_before_serve":
					served++
					if e := ServeWorker(f.d, r.Workspace, f.s.Result.RunID); e != nil {
						t.Fatal(e)
					}
					return done
				}
				return nil
			}
			err := SuperviseWorker(f.d, r.Workspace, f.s.Result.RunID)
			if !errors.Is(err, done) {
				t.Fatal(err)
			}
			f.reload()
			want := map[string]string{"pause": "paused", "stop": "stopped"}[action]
			if served != 1 || f.s.Result.Worker.Phase != want || f.s.Result.Worker.RetryNotBefore != due || wakes.Load() != 1 {
				t.Error("existing supervisor ignored bound request until future retry", served, waiting, f.s.Result.Worker)
			}
			if !f.now.Before(timeMustParse(t, due)) {
				t.Fatal("test advanced retry clock")
			}
		})
	}
}
func timeMustParse(t *testing.T, s string) time.Time {
	t.Helper()
	v, e := time.Parse(time.RFC3339Nano, s)
	if e != nil {
		t.Fatal(e)
	}
	return v
}

func mr05UnreservedIntentResume(t *testing.T) {
	for _, parent := range []bool{false, true} {
		t.Run(map[bool]string{false: "standalone", true: "week_child"}[parent], func(t *testing.T) {
			f, _ := mr05Worker(t)
			root := f.s.Request.WorkspaceRoot
			id := f.s.Result.RunID
			marker := filepath.Join(root, "attempt-started")
			control := f.file(filepath.Join(root, "count-attempt"), []byte("#!/bin/sh\nprintf 'attempt\\n' >> '"+marker+"'\n"))
			if e := os.Chmod(control.Locator, 0700); e != nil {
				t.Fatal(e)
			}
			f.s.Request.Runtime.PlyExecutable = Executable{control.Locator, control.SHA256}
			var gate ParentControlGate
			if parent {
				f.s.Request.ControlContract = NewControlContract("week", "week_"+strings.Repeat("a", 32), "sha256:"+strings.Repeat("b", 64))
				gate = ParentControlGate{Contract: f.s.Request.ControlContract, Deadline: ControlTimestamp(f.now.Add(time.Hour)), DesiredAction: "continue", ActiveRun: id}
				if e := f.d.replaceValue(ParentControlPath(root, *gate.Contract.WeekID), gate); e != nil {
					t.Fatal(e)
				}
			}
			f.s.Result.RequestSHA256 = digest(f.s.Request)
			f.save()
			f.own()
			dir := filepath.Join(workerDirectory(f.s), "attempts", "developer-0002")
			in := workerAttemptIntent{Kind: "developer", Argv: []string{}, Environment: []string{}}
			won := false
			f.d.Fault = func(point string) error {
				if point != "worker_before_attempt_reservation" || won {
					return nil
				}
				won = true
				_, e := os.Stat(filepath.Join(dir, "intent.json"))
				t.Log("intent exists at pause barrier:", e == nil)
				if parent {
					gate.DesiredAction = "pause"
					return WithWorkerParentGuard(root, gate.Contract, func() error { return f.d.replaceValue(ParentControlPath(root, *gate.Contract.WeekID), gate) })
				}
				r := mr05Preview(t, f, "pause", "intent-pause-0001")
				_, e = WorkerControls(f.d).Apply(r)
				return e
			}
			if e := workerLaunchAttempt(f.d, f.s, in, dir); e == nil {
				t.Fatal("pause lost dispatch barrier")
			}
			if !won {
				t.Fatal("barrier not reached")
			}
			f.reload()
			if f.s.Result.Worker.Job != nil {
				t.Fatal("pause reserved job")
			}
			if _, e := os.Stat(marker); !os.IsNotExist(e) {
				t.Fatal("attempt started before pause", e)
			}
			f.d.Fault = nil
			if parent {
				f.d.WorkerOwner.lock.Close()
				f.d.WorkerOwner = nil
				t.Chdir(root)
				if e := ServeWorker(f.d, root, id); e != nil {
					t.Fatal(e)
				}
				f.reload()
				if f.s.Result.Worker.Phase != "paused" {
					t.Fatal("parent child failed to park")
				}
				gate.DesiredAction = "continue"
				gate.ResumeAllowed = true
				if e := WithWorkerParentGuard(root, gate.Contract, func() error { return f.d.replaceValue(ParentControlPath(root, *gate.Contract.WeekID), gate) }); e != nil {
					t.Fatal(e)
				}
				c := f.control("resume")
				f.reload()
				f.own()
				if e := workerClaimOwner(&f.d, f.s, f.d.WorkerOwner.lock); e != nil {
					t.Fatal(e)
				}
				f.reload()
				if e := workerConsumeResume(f.d, f.s, c); e != nil {
					t.Fatal(e)
				}
			} else {
				mr05Consume(t, f)
				f.d.WorkerOwner.lock.Close()
				f.d.WorkerOwner = nil
				r := mr05Preview(t, f, "resume", "intent-resume-0001")
				if _, e := WorkerControls(f.d).Apply(r); e != nil {
					t.Fatal(e)
				}
				f.reload()
				f.own()
				if e := workerClaimOwner(&f.d, f.s, f.d.WorkerOwner.lock); e != nil {
					t.Fatal(e)
				}
				mr05Consume(t, f)
			}
			f.reload()
			if e := workerLaunchAttempt(f.d, f.s, in, dir); e != nil {
				t.Error("won pause poisoned legitimate successor dispatch", e)
			}
			b, e := os.ReadFile(marker)
			if e != nil || string(b) != "attempt\n" {
				t.Error("successor did not launch exactly one attempt", string(b), e)
			}
		})
	}
}
func mr05UnownedRunningEvidence(t *testing.T) {
	for _, owned := range []bool{false, true} {
		t.Run(map[bool]string{false: "unowned", true: "owned"}[owned], func(t *testing.T) {
			f, _ := mr05Worker(t)
			f.job("reviewer")
			f.save()
			if owned {
				f.own()
			}
			if !owned {
				for _, name := range []string{"owner.lock", "supervisor.lock"} {
					lock, e := workerLock(filepath.Join(workerDirectory(f.s), name), true)
					if e != nil {
						t.Fatal(e)
					}
					lock.Close()
				}
			}
			b, e := ShowControl(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
			if e != nil {
				t.Fatal(e)
			}
			if owned {
				if len(b.Actions) != 2 {
					t.Error("owned active attempt cannot request boundary")
				}
			} else if len(b.Actions) != 0 || b.UnavailableReason == nil || *b.UnavailableReason != "attempt_outcome_unknown" {
				t.Error("unowned job without result still offers control", b)
			}
		})
	}
}
func mr05ResumeConsumptionDrift(t *testing.T) {
	for _, mutation := range []string{"provider", "result.json", "stdout.jsonl", "stderr.log", "job_removed", "unchanged"} {
		t.Run(mutation, func(t *testing.T) {
			f, _ := mr05Worker(t)
			f.job("reviewer")
			f.s.Result.Worker.Phase = "paused"
			f.save()
			mr05PositiveResult(f)
			r := mr05Preview(t, f, "resume", "resume-drift-0001")
			if _, e := WorkerControls(f.d).Apply(r); e != nil {
				t.Fatal(e)
			}
			f.reload()
			f.own()
			switch mutation {
			case "unchanged":
			case "provider":
				if e := os.WriteFile(f.s.Request.Runtime.Executable.Path, []byte("changed"), 0700); e != nil {
					t.Fatal(e)
				}
			case "job_removed":
				f.s.Result.Worker.Job = nil
				f.save()
			default:
				if e := os.Remove(filepath.Join(f.s.Result.Worker.Job.Directory, mutation)); e != nil {
					t.Fatal(e)
				}
			}
			_, err := workerConsumeBoundControl(f.d, f.s)
			f.reload()
			status, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
			if mutation == "unchanged" {
				if err != nil || e != nil || status.ControlEffect != "resumed" {
					t.Fatal("legitimate resume refused", status, err, e)
				}
				return
			}
			if err == nil || status.ControlEffect != "unknown" || status.RequestState != "outcome_unknown" || status.Receipt != nil || f.s.Result.Worker.Phase != "paused" {
				t.Error("unproven resume acknowledged", status, err)
			}
		})
	}
}
func mr05TerminalConsumption(t *testing.T) {
	for _, mutation := range []string{"generation", "job"} {
		t.Run(mutation, func(t *testing.T) {
			f, wakes := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "terminal-worker-0001")
			if _, e := WorkerControls(f.d).Apply(r); e != nil {
				t.Fatal(e)
			}
			f.reload()
			if mutation == "generation" {
				f.s.Result.Worker.Generation += 2
			} else {
				f.job("reviewer")
			}
			f.save()
			f.own()
			_, err := workerConsumeBoundControl(f.d, f.s)
			if err == nil {
				t.Fatal("mismatch consumed")
			}
			if e := workerBlocked(f.d, r.Workspace, f.s.Result.RunID, err); e != nil {
				t.Fatal(e)
			}
			supervisor, e := workerLock(filepath.Join(workerDirectory(f.s), "supervisor.lock"), true)
			if e != nil {
				t.Fatal(e)
			}
			defer supervisor.Close()
			exited, e := workerSupervisorExit(f.d, r.Workspace, f.s.Result.RunID, supervisor)
			if !exited && e == nil {
				t.Error("terminal mismatch asks supervisor to spawn another owner")
			}
			supervisor.Close()
			f.d.WorkerOwner.lock.Close()
			f.d.Fault = func(point string) error {
				if point == "worker_supervisor_before_serve" {
					t.Error("terminal supervisor relaunched controller")
					return errors.New("unexpected owner launch")
				}
				return nil
			}
			if e = SuperviseWorker(f.d, r.Workspace, f.s.Result.RunID); e == nil {
				t.Error("terminal supervisor lost unknown outcome")
			}
			f.reload()
			generation := f.s.Result.Worker.Generation
			for i := 0; i < 3; i++ {
				status, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
				if e == nil || status.RequestState != "outcome_unknown" || status.ControlEffect != "unknown" || status.Receipt != nil {
					t.Error("terminal mismatch reported pending", status, e)
				}
				_, _ = WorkerControls(f.d).Apply(r)
			}
			f.reload()
			if f.s.Result.Worker.Generation != generation || wakes.Load() != 1 {
				t.Error("terminal replay retried owner")
			}
		})
	}
}
