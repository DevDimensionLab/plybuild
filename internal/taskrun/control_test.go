//go:build !windows

package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func mr05Worker(t *testing.T) (*aw01Fixture, *atomic.Int32) {
	f := newAW01Fixture(t)
	wake := &atomic.Int32{}
	f.d.ControlQueue = func(string, string) (string, error) { return "current-queue-revision-1", nil }
	f.d.ControlWake = func(string, string) error { wake.Add(1); return nil }
	f.s.Request.ControlContract = NewControlContract("standalone", "", "")
	f.s.Result.RequestSHA256 = digest(f.s.Request)
	f.s.Result.Worker.Job = nil
	f.save()
	return f, wake
}
func mr05Preview(t *testing.T, f *aw01Fixture, action, id string) ControlRequest {
	t.Helper()
	p, e := WorkerControls(f.d).Preview(f.s.Request.WorkspaceRoot, f.s.Result.RunID, action, id)
	if e != nil {
		t.Fatal(e)
	}
	return *p.Request
}
func mr05Consume(t *testing.T, f *aw01Fixture) {
	t.Helper()
	f.reload()
	if f.d.WorkerOwner == nil {
		f.own()
	}
	yes, e := workerConsumeBoundControl(f.d, f.s)
	if !yes || e != nil {
		t.Fatalf("native consume=%v %v", yes, e)
	}
	f.reload()
}
func mr05Status(t *testing.T, f *aw01Fixture, id, effect string) ControlStatus {
	t.Helper()
	s, e := WorkerControls(f.d).Status(f.s.Request.WorkspaceRoot, id)
	if e != nil || s.ControlEffect != effect {
		t.Fatalf("status %s: %+v %v", id, s, e)
	}
	return s
}

func TestMR05ControlDiscoveryAndOwnership(t *testing.T) {
	t.Run("unowned_running_requires_result", mr05UnownedRunningEvidence)
	f := newAW01Fixture(t)
	f.d.ControlQueue = func(string, string) (string, error) { return "current", nil }
	f.s.Request.ControlContract = NewControlContract("standalone", "", "")
	f.s.Result.RequestSHA256 = digest(f.s.Request)
	f.save()
	f.own()
	b, err := ShowControl(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if b.Scope.Kind != "run" || len(b.Actions) != 2 || b.Actions[0].Action != "pause" || b.Actions[1].Action != "stop" {
		t.Fatalf("new explicitly bound worker must offer pause and stop: %+v", b)
	}
	for _, tc := range []struct {
		name, phase, reason string
		contract            *ControlContract
		released            bool
		actions             int
	}{
		{"old", "running", "control_unsupported", nil, false, 0},
		{"parent", "running", "parent_control_unsupported", NewControlContract("week", "week_"+strings.Repeat("a", 32), "sha256:"+strings.Repeat("a", 64)), false, 0},
		{"released", "paused", "queue_released", NewControlContract("standalone", "", ""), true, 0},
		{"completed", "completed", "completed", NewControlContract("standalone", "", ""), false, 0},
		{"blocked", "blocked", "control_unavailable", NewControlContract("standalone", "", ""), false, 0},
		{"paused", "paused", "", NewControlContract("standalone", "", ""), false, 2},
		{"stopped", "stopped", "", NewControlContract("standalone", "", ""), false, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, _ := mr05Worker(t)
			f.s.Request.ControlContract = tc.contract
			f.s.Result.RequestSHA256 = digest(f.s.Request)
			f.s.Result.Worker.Phase = tc.phase
			f.s.Result.Worker.QueueReleased = tc.released
			f.s.Result.Worker.OwnerPID = 0
			f.save()
			b, e := ShowControl(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
			if e != nil || len(b.Actions) != tc.actions {
				t.Fatalf("discovery %+v %v", b, e)
			}
			if tc.reason != "" && (b.UnavailableReason == nil || *b.UnavailableReason != tc.reason) {
				t.Fatalf("reason %+v", b)
			}
		})
	}
}

func TestMR05ControlPreviewStatusReadOnlyStrictRequest(t *testing.T) {
	f, wakes := mr05Worker(t)
	r := mr05Preview(t, f, "pause", "readonly-0001")
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	f.d.Files = aw01NoEffects{t}
	for i := 0; i < 3; i++ {
		b, e := ShowControl(f.d, r.Workspace, f.s.Result.RunID)
		if e != nil || b.RunID == nil {
			t.Fatal(b, e)
		}
		p, e := WorkerControls(f.d).Preview(r.Workspace, f.s.Result.RunID, "pause", r.RequestID)
		if e != nil || p.Request.Binding != r.Binding {
			t.Fatal(p, e)
		}
		s, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
		if e != nil || s.Scope != nil || s.Action != nil || s.RequestDigest != nil || s.ControlEffect != "not_confirmed" {
			t.Fatal(s, e)
		}
		raw, _ := json.Marshal([]any{b, p, s})
		for _, secret := range []string{"callback_secret", "provider-stream", f.s.Request.Runtime.Executable.Path, f.s.Result.Paths.RunRoot} {
			if strings.Contains(string(raw), secret) {
				t.Fatal("private native information leaked")
			}
		}
	}
	if !reflect.DeepEqual(before, aw01Inventory(t, r.Workspace)) || wakes.Load() != 0 {
		t.Fatal("read changed store or woke provider")
	}
	raw, _ := Canonical(r)
	if _, e := ParseControlRequest(raw); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []string{strings.Replace(string(raw), `"schema_version":1`, `"schema_version":2`, 1), strings.Replace(string(raw), `"schema_version":1`, `"schema_version":"1"`, 1), strings.Replace(string(raw), `"action":"pause"`, `"action":"kill"`, 1), strings.Replace(string(raw), `"kind":"run"`, `"kind":"RUN"`, 1), strings.Replace(string(raw), `"request_id":`, `"unknown":true,"request_id":`, 1), strings.Replace(string(raw), `"request_id":`, `"request_id":"duplicate","request_id":`, 1), strings.Replace(string(raw), `"binding":"`+r.Binding+`"`, `"binding":null`, 1), `null`} {
		if _, e := ParseControlRequest([]byte(bad)); e == nil {
			t.Fatalf("accepted malformed request: %s", bad)
		}
	}
	tooLong := r
	tooLong.Validity.ExpiresAt = ControlTimestamp(f.now.Add(6 * time.Minute))
	if ValidateControlRequest(tooLong) == nil {
		t.Fatal("accepted validity > 5 minutes")
	}
}

func TestMR05ControlStaleBasisAndABA(t *testing.T) {
	t.Run("observation_timestamps_do_not_change_revision", func(t *testing.T) {
		f, _ := mr05Worker(t)
		r := mr05Preview(t, f, "pause", "telemetry-0001")
		before, e := ShowControl(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
		if e != nil {
			t.Fatal(e)
		}
		f.now = f.now.Add(time.Minute)
		if e = workflowUpdate(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID, func(s *workflowState) error {
			s.Result.Worker.UpdatedAt = ControlTimestamp(f.now)
			s.Result.Worker.OwnerState = "observed heartbeat"
			s.Result.Worker.ControlRequestedAt = ControlTimestamp(f.now)
			return nil
		}); e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(workerDirectory(f.s), "controller.log"), []byte("additional observed log output\n"), 0600); e != nil {
			t.Fatal(e)
		}
		after, e := ShowControl(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
		if e != nil || after.StateRevision != before.StateRevision {
			t.Fatal("observation clock invalidated semantic revision", e)
		}
		if _, e = WorkerControls(f.d).Apply(r); e != nil {
			t.Fatal("observational log/heartbeat growth invalidated approval", e)
		}
	})
	mutations := map[string]func(*workflowState){"owner": func(s *workflowState) { s.Result.Worker.Generation++ }, "stage": func(s *workflowState) { s.Result.Worker.Stage = "integrate" }, "job": func(s *workflowState) {
		s.Result.Worker.Job = &WorkerJob{ID: "other", IntentHash: hash([]byte("other"))}
	}, "phase": func(s *workflowState) { s.Result.Worker.Phase = "waiting" }, "queue_release": func(s *workflowState) { s.Result.Worker.QueueReleased = true }, "provider_wait": func(s *workflowState) { s.Result.Worker.RetryNotBefore = "2027-01-01T00:00:00Z" }, "policy": func(s *workflowState) { s.Result.Delivery.PermissionState = "revoked" }}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			f, w := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "stale-0001")
			if e := workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error { mutate(s); return nil }); e != nil {
				t.Fatal(e)
			}
			s, e := WorkerControls(f.d).Apply(r)
			if e == nil || s.ControlEffect != "none" || w.Load() != 0 {
				t.Fatal("stale request had effect", s, e)
			}
			if _, e = os.Stat(controlPath(r.Workspace, r.RequestID, "request")); !os.IsNotExist(e) {
				t.Fatal("stale reserved")
			}
		})
	}
	for _, name := range []string{"runtime", "frozen", "scope", "parent", "request", "rights"} {
		t.Run(name, func(t *testing.T) {
			f, w := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "drift-0001")
			switch name {
			case "runtime":
				os.WriteFile(f.s.Request.Runtime.Executable.Path, []byte("changed"), 0700)
			case "frozen":
				os.WriteFile(f.s.Request.Runtime.PlyExecutable.Path, []byte("changed"), 0700)
			case "rights":
				f.d.ControlQueue = func(string, string) (string, error) { return "released-revision-2", errors.New("released") }
			case "scope":
				f.s.Request.WorkspaceRoot += "/other"
				f.save()
			case "parent":
				f.s.Request.ControlContract = NewControlContract("week", "week_"+strings.Repeat("b", 32), hash([]byte("week")))
				f.save()
			case "request":
				f.s.Request.Runtime.Model = "changed"
				f.save()
			}
			_, e := WorkerControls(f.d).Apply(r)
			if name == "frozen" {
				if e != nil || w.Load() != 1 {
					t.Fatal("Ply upgrade blocked unchanged control", e)
				}
			} else if e == nil || w.Load() != 0 {
				t.Fatal("binding drift accepted", e)
			}
		})
	}
	f, w := mr05Worker(t)
	r := mr05Preview(t, f, "pause", "aba-test-0001")
	for _, stage := range []string{"developer", "reviewer"} {
		if e := workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error { s.Result.Worker.Stage = stage; return nil }); e != nil {
			t.Fatal(e)
		}
	}
	if _, e := WorkerControls(f.d).Apply(r); e == nil || w.Load() != 0 {
		t.Fatal("ABA reused revision")
	}
	r = mr05Preview(t, f, "pause", "expiry-0001")
	f.now = f.now.Add(6 * time.Minute)
	if _, e := WorkerControls(f.d).Apply(r); ControlFailureOf(e).Code != "control_expired" {
		t.Fatal(e)
	}
	r = mr05Preview(t, f, "pause", "historical-0001")
	if _, e := WorkerControls(f.d).Apply(r); e != nil {
		t.Fatal(e)
	}
	f.now = f.now.Add(time.Hour)
	if _, e := WorkerControls(f.d).Apply(r); e != nil || w.Load() != 1 {
		t.Fatal("historical replay changed effect", e)
	}
}

func TestMR05ControlIdempotentConcurrentApply(t *testing.T) {
	t.Run("concurrent_first_store_guard", func(t *testing.T) {
		f, _ := mr05Worker(t)
		var wg sync.WaitGroup
		errs := make(chan error, 12)
		for i := 0; i < 12; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); errs <- withStore(f.s.Request.WorkspaceRoot, func() error { return nil }) }()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Errorf("concurrent native store guard: %T: %v", err, err)
			}
		}
	})
	f, w := mr05Worker(t)
	r := mr05Preview(t, f, "pause", "parallel-0001")
	var wg sync.WaitGroup
	type outcome struct {
		status ControlStatus
		err    error
	}
	errs := make(chan outcome, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); s, e := WorkerControls(f.d).Apply(r); errs <- outcome{s, e} }()
	}
	wg.Wait()
	close(errs)
	for result := range errs {
		if result.err != nil {
			raw, _ := json.Marshal(result.status)
			t.Fatalf("concurrent apply: %v; status=%s", result.err, raw)
		}
	}
	if w.Load() != 1 {
		t.Fatal("duplicate wake", w.Load())
	}
	mr05Consume(t, f)
	mr05Status(t, f, r.RequestID, "paused")
	fresh := f.d
	fresh.ControlWake = func(string, string) error { t.Fatal("replay woke controller"); return nil }
	if _, e := WorkerControls(fresh).Apply(r); e != nil {
		t.Fatal(e)
	}
	conflict := r
	conflict.Input.Action = "stop"
	if out, e := WorkerControls(fresh).Apply(conflict); ControlFailureOf(e).Code != "control_request_conflict" || out.Error == nil || out.Error.Code != "control_request_conflict" {
		t.Fatal("conflict must carry a public JSON error", out, e)
	}
	f2, w2 := mr05Worker(t)
	a := mr05Preview(t, f2, "pause", "winner-0001")
	b := mr05Preview(t, f2, "stop", "winner-0002")
	wins := atomic.Int32{}
	for _, r := range []ControlRequest{a, b} {
		wg.Add(1)
		go func(r ControlRequest) {
			defer wg.Done()
			if _, e := WorkerControls(f2.d).Apply(r); e == nil {
				wins.Add(1)
			}
		}(r)
	}
	wg.Wait()
	if wins.Load() != 1 || w2.Load() != 1 {
		t.Fatal("competing controls both won")
	}
}

type mr05CommittedWriteError struct{ systemFileSystem }

func (fs mr05CommittedWriteError) WriteOnce(path string, raw []byte) error {
	if filepath.Base(path) != "request.json" {
		return fs.systemFileSystem.WriteOnce(path, raw)
	}
	if err := fs.systemFileSystem.WriteOnce(path, raw); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
		return err
	}
	return errors.New("directory fsync failed after publication; readback is corrupt")
}

func TestMR05ControlFaultsPreserveUnknown(t *testing.T) {
	t.Run("resume_drift_after_reservation", mr05ResumeConsumptionDrift)
	t.Run("terminal_consumption_exits_supervisor", mr05TerminalConsumption)
	t.Run("receipt_publication_recovers_after_successor", mr05WorkerReceiptRecovery)
	t.Run("publication_error_after_commit", func(t *testing.T) {
		f, wakes := mr05Worker(t)
		r := mr05Preview(t, f, "pause", "write-error-0001")
		f.d.Files = mr05CommittedWriteError{}
		status, err := WorkerControls(f.d).Apply(r)
		if err == nil || status.RequestState != "outcome_unknown" || status.ControlEffect != "unknown" || status.RequestDigest == nil || *status.RequestDigest != digest(r) || wakes.Load() != 0 {
			t.Fatal("ambiguous publication error was reported as no reservation", status, err)
		}
	})
	t.Run("missing_record_replay_cannot_create_new_reservation", func(t *testing.T) {
		f, wakes := mr05Worker(t)
		r := mr05Preview(t, f, "pause", "missing-record-0001")
		if _, err := WorkerControls(f.d).Apply(r); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(controlPath(r.Workspace, r.RequestID, "request")); err != nil {
			t.Fatal(err)
		}
		before := aw01Inventory(t, r.Workspace)
		status, err := WorkerControls(f.d).Apply(r)
		if err == nil || status.RequestState != "outcome_unknown" || status.ControlEffect != "unknown" || wakes.Load() != 1 || !reflect.DeepEqual(before, aw01Inventory(t, r.Workspace)) {
			t.Fatal("partial identity was treated as a fresh request", status, err)
		}
	})
	for _, point := range []string{"control_after_reserve", "control_after_wake"} {
		t.Run("durable_reservation_corrupt_readback/"+point, func(t *testing.T) {
			f, wakes := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "reserved-corrupt-0001")
			var original []byte
			f.d.Fault = func(p string) error {
				if p != point {
					return nil
				}
				path := controlPath(r.Workspace, r.RequestID, "request")
				var err error
				original, err = os.ReadFile(path)
				if err != nil {
					t.Fatal("reservation was not durably published", err)
				}
				var record ControlRecord
				if err = json.Unmarshal(original, &record); err != nil || record.Digest != digest(r) {
					t.Fatal("reservation did not preserve the original request", err)
				}
				if err = os.WriteFile(path, []byte("{"), 0600); err != nil {
					t.Fatal(err)
				}
				return errors.New("failure after durable reservation and corrupt readback")
			}
			s, err := WorkerControls(f.d).Apply(r)
			if len(original) == 0 {
				t.Fatal("post-reservation fault was not exercised")
			}
			wantWakes := int32(0)
			if point == "control_after_wake" {
				wantWakes = 1
			}
			if wakes.Load() != wantWakes {
				t.Fatalf("wakes=%d, want %d", wakes.Load(), wantWakes)
			}
			if err == nil || s.RequestState != "outcome_unknown" || s.ControlEffect != "unknown" {
				t.Errorf("durable reservation was falsely rejected: state=%s effect=%s err=%v", s.RequestState, s.ControlEffect, err)
			}
			if s.RequestDigest == nil || *s.RequestDigest != digest(r) || s.ApprovalBinding == nil || *s.ApprovalBinding != r.Binding || s.Scope == nil || *s.Scope != r.Input.Target {
				t.Error("post-reservation uncertainty lost the invocation's known identity")
			}
			f.d.Fault = nil
			before := aw01Inventory(t, r.Workspace)
			for i := 0; i < 2; i++ {
				status, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
				if e == nil || status.RequestState != "outcome_unknown" || status.ControlEffect != "unknown" {
					t.Fatal("corrupt status concealed uncertainty", status, e)
				}
				replay, e := WorkerControls(f.d).Apply(r)
				if e == nil || replay.RequestState != "outcome_unknown" || replay.ControlEffect != "unknown" {
					t.Fatal("corrupt replay concealed uncertainty", replay, e)
				}
			}
			if wakes.Load() != wantWakes || !reflect.DeepEqual(before, aw01Inventory(t, r.Workspace)) {
				t.Fatal("status/replay repeated the effect")
			}
		})
	}
	t.Run("ply_upgrade_after_reservation", func(t *testing.T) {
		f, _ := mr05Worker(t)
		r := mr05Preview(t, f, "pause", "consume-drift-0001")
		if _, e := WorkerControls(f.d).Apply(r); e != nil {
			t.Fatal(e)
		}
		f.own()
		if e := os.WriteFile(f.s.Request.Runtime.PlyExecutable.Path, []byte("different controller"), 0700); e != nil {
			t.Fatal(e)
		}
		if consumed, e := workerConsumeBoundControl(f.d, f.s); e != nil || !consumed {
			t.Fatal("Ply upgrade blocked the reserved control", e)
		}
		f.reload()
		if f.s.Result.Worker.Phase != "paused" {
			t.Fatal("reserved pause was not confirmed")
		}
	})
	for _, point := range []string{"control_before_reserve", "control_after_reserve", "control_before_wake", "control_after_wake", "control_before_consume", "control_after_consume", "control_before_receipt", "control_after_receipt"} {
		t.Run(point, func(t *testing.T) {
			f, w := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "fault-0001")
			f.d.Fault = func(p string) error {
				if p == point {
					return errors.New("injected crash")
				}
				return nil
			}
			_, e := WorkerControls(f.d).Apply(r)
			if point == "control_before_reserve" {
				if e == nil || w.Load() != 0 {
					t.Fatal("pre-reservation fault had effect")
				}
				return
			}
			if strings.Contains(point, "consume") || strings.Contains(point, "receipt") {
				if e != nil {
					t.Fatal(e)
				}
				f.own()
				_, e = workerConsumeBoundControl(f.d, f.s)
				if e == nil {
					t.Fatal("fault not reached")
				}
			}
			before := aw01Inventory(t, r.Workspace)
			count := w.Load()
			s, _ := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
			if s.RequestDigest == nil || s.Sequence == nil {
				t.Fatal("lost durable identity")
			}
			if point != "control_after_receipt" && point != "control_before_consume" && s.ControlEffect != "unknown" {
				t.Fatalf("uncertainty became success %+v", s)
			}
			if !reflect.DeepEqual(before, aw01Inventory(t, r.Workspace)) || w.Load() != count {
				t.Fatal("status recovered effect")
			}
			f.d.Fault = nil
			_, _ = WorkerControls(f.d).Apply(r)
			if w.Load() != count {
				t.Fatal("replay woke after uncertain effect")
			}
		})
	}
	f, w := mr05Worker(t)
	r := mr05Preview(t, f, "pause", "corrupt-0001")
	if _, e := WorkerControls(f.d).Apply(r); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(controlPath(r.Workspace, r.RequestID, "request"), []byte("{"), 0600); e != nil {
		t.Fatal(e)
	}
	s, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
	if e == nil || s.ControlEffect != "unknown" || w.Load() != 1 {
		t.Fatal("truncation hidden", s, e)
	}
}

func TestMR05ControlAttemptBoundariesAndSameRunResume(t *testing.T) {
	t.Run("existing_waiting_supervisor", mr05ExistingWaitingSupervisor)
	t.Run("pause_before_job_then_successor_resume", mr05UnreservedIntentResume)
	for _, scenario := range []string{"exact_wake_successor", "owner_observed_before_wake", "second_successor"} {
		t.Run(scenario, func(t *testing.T) {
			f, wakes := mr05Worker(t)
			observed := f.s
			predecessor := f.s.Result.Worker.Generation
			r := mr05Preview(t, f, "pause", "linked-owner-0001")
			if _, err := WorkerControls(f.d).Apply(r); err != nil {
				t.Fatal(err)
			}
			f.reload()
			f.own()
			if scenario != "owner_observed_before_wake" {
				observed = f.s
			}
			if err := workerClaimOwner(&f.d, observed, f.d.WorkerOwner.lock); err != nil {
				t.Fatal(err)
			}
			f.reload()
			if scenario == "second_successor" {
				if err := workerClaimOwner(&f.d, f.s, f.d.WorkerOwner.lock); err != nil {
					t.Fatal(err)
				}
				f.reload()
			}
			_, err := workerConsumeBoundControl(f.d, f.s)
			status, _ := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
			if scenario == "exact_wake_successor" {
				if err != nil || status.ControlEffect != "paused" || status.Receipt == nil || status.Receipt.PredecessorGeneration == nil || *status.Receipt.PredecessorGeneration != predecessor || status.Receipt.OwnerGeneration == nil || *status.Receipt.OwnerGeneration != predecessor+1 {
					t.Fatal("successor did not bind its approved predecessor", status, err)
				}
			} else if err == nil || status.Receipt != nil {
				t.Fatal("unlinked owner consumed approval", status, err)
			}
			if wakes.Load() != 1 {
				t.Fatal("owner claim repeated the wake", wakes.Load())
			}
		})
	}
	t.Run("unlinked_successor_cannot_consume", func(t *testing.T) {
		f, _ := mr05Worker(t)
		r := mr05Preview(t, f, "pause", "unlinked-owner-0001")
		if _, err := WorkerControls(f.d).Apply(r); err != nil {
			t.Fatal(err)
		}
		// An unrelated native mutation is not a controller wake claim. A mere
		// generation increment cannot authorize this owner to consume approval.
		if err := workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error {
			s.Result.Worker.Generation++
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		f.reload()
		f.own()
		_, err := workerConsumeBoundControl(f.d, f.s)
		status, _ := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
		if err == nil || status.ControlEffect == "paused" || status.Receipt != nil {
			t.Fatal("unlinked successor consumed the prior generation approval", status, err)
		}
	})
	t.Run("commit_requires_native_completion", func(t *testing.T) {
		f, _ := mr05Worker(t)
		f.job("commit")
		f.save()
		f.own()
		r := mr05Preview(t, f, "pause", "commit-boundary-0001")
		if _, e := WorkerControls(f.d).Apply(r); e != nil {
			t.Fatal(e)
		}
		f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
		if _, e := workerConsumeBoundControl(f.d, f.s); e == nil {
			t.Fatal("commit without native completion was parked")
		}
	})
	for _, kind := range []string{"developer", "reviewer", "pretest", "verify", "integrate"} {
		t.Run(kind, func(t *testing.T) {
			for _, complete := range []bool{false, true} {
				t.Run(map[bool]string{false: "unproven_stays_unknown", true: "positive_boundary"}[complete], func(t *testing.T) {
					f, _ := mr05Worker(t)
					f.job(kind)
					f.save()
					f.own()
					r := mr05Preview(t, f, "pause", "boundary-0001")
					if _, e := WorkerControls(f.d).Apply(r); e != nil {
						t.Fatal(e)
					}
					mr05Status(t, f, r.RequestID, "pending_boundary")
					if !complete {
						if _, e := workerConsumeBoundControl(f.d, f.s); e == nil {
							t.Fatal("parked without positive completion")
						}
						f.reload()
						if f.s.Result.Worker.Phase != "running" {
							t.Fatal("unconfirmed park")
						}
					}
					result := f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
					if kind == "verify" || kind == "integrate" {
						result.NativeCompleted = true
						if e := f.d.replaceValue(filepath.Join(f.s.Result.Worker.Job.Directory, "result.json"), result); e != nil {
							t.Fatal(e)
						}
					}
					if !complete {
						// Late evidence cannot automatically reverse a terminal rejection.
						if _, e := workerConsumeBoundControl(f.d, f.s); e == nil {
							t.Fatal("late evidence retried a terminal request")
						}
						status, e := WorkerControls(f.d).Status(r.Workspace, r.RequestID)
						if e == nil || status.ControlEffect != "unknown" || status.Receipt != nil {
							t.Fatal("terminal rejection became positive", status, e)
						}
						return
					}
					effect := "paused"
					if kind == "integrate" {
						effect = "completed_before_boundary"
						if e := workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error { s.Result.Delivery.Phase = "completed"; return nil }); e != nil {
							t.Fatal(e)
						}
						if e := workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error {
							s.Result.Worker.Job = nil
							s.Result.Worker.Stage = "completed"
							return nil
						}); e != nil {
							t.Fatal(e)
						}
					}
					mr05Consume(t, f)
					mr05Status(t, f, r.RequestID, effect)
				})
			}
		})
	}
	f, _ := mr05Worker(t)
	original := f.s.Request
	session := f.s.Result.Worker.SessionID
	f.s.Result.Worker.RetryNotBefore = ControlTimestamp(f.now.Add(time.Hour))
	f.s.Result.Worker.Phase = "waiting"
	f.save()
	r := mr05Preview(t, f, "pause", "wait-pause-0001")
	WorkerControls(f.d).Apply(r)
	mr05Consume(t, f)
	r = mr05Preview(t, f, "resume", "wait-resume-0001")
	if _, e := WorkerControls(f.d).Apply(r); e != nil {
		t.Fatal(e)
	}
	mr05Consume(t, f)
	mr05Status(t, f, r.RequestID, "resumed")
	if f.s.Result.Worker.Phase != "waiting" || f.s.Result.Worker.RetryNotBefore != ControlTimestamp(f.now.Add(time.Hour)) || f.s.Result.Worker.SessionID != session || !equal(original, f.s.Request) {
		t.Fatal("resume replaced Run or shortened provider wait")
	}
}

func TestMR05ControlLegacyAndNativeMutationRaces(t *testing.T) {
	for _, first := range []string{"public", "legacy", "native"} {
		t.Run(first, func(t *testing.T) {
			f, w := mr05Worker(t)
			r := mr05Preview(t, f, "pause", "race-0001")
			entered, release := make(chan struct{}), make(chan struct{})
			done := make(chan error, 1)
			go func() { done <- withStore(r.Workspace, func() error { close(entered); <-release; return nil }) }()
			<-entered
			apply := make(chan error, 1)
			go func() {
				if first == "legacy" {
					apply <- workerControlWrite(f.d, f.s, "stop")
				} else if first == "native" {
					apply <- workflowUpdate(f.d, r.Workspace, f.s.Result.RunID, func(s *workflowState) error { s.Result.Worker.Generation++; return nil })
				} else {
					_, e := WorkerControls(f.d).Apply(r)
					apply <- e
				}
			}()
			close(release)
			if e := <-done; e != nil {
				t.Fatal(e)
			}
			if e := <-apply; e != nil {
				t.Fatal(e)
			}
			if first == "public" {
				if e := workerControlWrite(f.d, f.s, "stop"); e == nil {
					t.Fatal("legacy overwrote pending public control")
				}
			} else {
				if _, e := WorkerControls(f.d).Apply(r); e == nil {
					t.Fatal("public overwrote preceding writer")
				}
			}
			if first == "native" && w.Load() != 0 || first != "native" && w.Load() != 1 {
				t.Fatal("unexpected control effects", w.Load())
			}
		})
	}
}
