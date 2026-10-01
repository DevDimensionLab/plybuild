package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// This is explicitly synthetic human-observation evidence, never product QA.
func syntheticStatus(t *testing.T, d Dependencies, r Request, state string) TaskStatus {
	t.Helper()
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	return TaskStatus{Envelope: env("task-status"), RunID: j.Result.RunID, RequestSHA256: digest(r), SessionID: "ply:" + j.Result.RunID, BasisEventSHA256: *j.Result.LastEventSHA256, Source: "human_observation", ActorClaim: "synthetic test observer (not human QA)", ObservedAtUTC: d.Now().UTC().Format(time.RFC3339Nano), State: state, WorktreePath: j.Binding.Preparation.Plan.WorktreePath, VisibleGroupPath: ptr(j.Binding.Preparation.Plan.WorktreePath), TaskLabel: ptr("synthetic task row"), MatchingTasks: ptr(1)}
}
func collectSyntheticStatus(t *testing.T, d Dependencies, r Request, s TaskStatus) (Result, error) {
	t.Helper()
	d.TaskStatusPath = writeAny(t, r.WorkspaceRoot, "synthetic-status.json", s)
	p, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		return Result{}, e
	}
	return Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *p.Confirmation)
}
func startAndObserve(t *testing.T, d Dependencies, r Request, file, confirm string) (Result, error) {
	t.Helper()
	out, e := Start(d, file, confirm)
	var typed *Error
	if !errors.As(e, &typed) || typed.Code != "task_run_task_status_unknown" {
		return out, e
	}
	return collectSyntheticStatus(t, d, r, syntheticStatus(t, d, r, "inactive"))
}

func TestReturnRepairStatusValidationAndQualification(t *testing.T) {
	for _, mode := range []string{"inactive", "active", "unknown", "ambiguous", "future", "stale", "unbound", "wrong-native", "hidden-native", "live-client", "unknown-client", "missing-report"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			if mode != "missing-report" {
				completeInside(t, &d, r, nil)
			}
			if mode == "live-client" || mode == "unknown-client" {
				state := "exited"
				q := ptr(false)
				if mode == "unknown-client" {
					state = "unknown"
					q = nil
				}
				d.Runner.(*fakeRunner).process = Process{State: state, PID: ptr(731), ProcessGroup: ptr(731), StartIdentity: ptr("synthetic-boot:birth-123"), ExitCode: ptr(0), Quiescence: q}
			}
			if mode == "wrong-native" || mode == "hidden-native" {
				d.Runner.(*fakeRunner).inside = func(LaunchSpec) error {
					a := acceptance(t, d, r)
					a.RuntimeClaim.NativeSessionID = ptr("synthetic-native")
					if _, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "claim.json", a)); e != nil {
						return e
					}
					_, e := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "return.json", semanticReport(t, d, r)))
					return e
				}
			}
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			Start(d, f, *p.(Preview).Confirmation)
			if e = d.Runner.(*fakeRunner).lastErr; e != nil {
				t.Fatal(e)
			}
			s := syntheticStatus(t, d, r, "inactive")
			switch mode {
			case "active", "unknown":
				s.State = mode
			case "ambiguous":
				s.MatchingTasks = ptr(2)
			case "future":
				s.ObservedAtUTC = d.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano)
			case "stale":
				s.ObservedAtUTC = "2020-01-01T00:00:00Z"
			case "unbound":
				s.RequestSHA256 = hash([]byte("other"))
			case "wrong-native":
				s.NativeSessionID = ptr("different")
			}
			d.TaskStatusPath = writeAny(t, r.WorkspaceRoot, "status.json", s)
			before := treeState(t, r.WorkspaceRoot)
			preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if before != treeState(t, r.WorkspaceRoot) {
				t.Fatal("check wrote")
			}
			invalid := mode == "ambiguous" || mode == "future" || mode == "stale" || mode == "unbound" || mode == "wrong-native"
			if invalid {
				if e == nil {
					t.Fatal("invalid status accepted")
				}
				return
			}
			if e != nil {
				t.Fatal(e)
			}
			out, e := Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			wantQualified := mode == "inactive" || mode == "hidden-native"
			if (out.Collection.State == "qualified") != wantQualified {
				t.Fatalf("qualification: %+v", out)
			}
			reg, e := d.Workspace.WorkItems.Snapshot(r.WorkspaceRoot)
			if e != nil {
				t.Fatal(e)
			}
			expected := 0
			if wantQualified {
				expected = 1
			}
			if len(reg.TaskResults) != expected {
				t.Fatal("wrong publication count")
			}
			if out.TaskExecution.State != s.State || out.TaskExecution.Source != "human_observation" {
				t.Fatal("task provenance lost")
			}
		})
	}
}

func TestReturnRepairStatusCASRecoveryAndLateReport(t *testing.T) {
	for _, mode := range []string{"file-changed", "tip-changed", "after_task_status_file", "after_event_task_status_observed", "after_task_result", "late-report"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			if mode != "late-report" {
				completeInside(t, &d, r, nil)
			}
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			Start(d, f, *p.(Preview).Confirmation)
			s := syntheticStatus(t, d, r, "inactive")
			d.TaskStatusPath = writeAny(t, r.WorkspaceRoot, "status.json", s)
			preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			if mode == "file-changed" {
				s.ActorClaim = "changed synthetic observer"
				writeAny(t, r.WorkspaceRoot, "status.json", s)
			}
			if mode == "tip-changed" {
				// Another genuine synthetic observation advances the journal.
				if _, e = collectSyntheticStatus(t, d, r, syntheticStatus(t, d, r, "active")); e != nil {
					t.Fatal(e)
				}
			}
			if strings.HasPrefix(mode, "after_") {
				d.Fault = func(point string) error {
					if point == mode {
						return errors.New("synthetic crash")
					}
					return nil
				}
			}
			before := treeState(t, r.WorkspaceRoot)
			out, e := Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
			if mode == "file-changed" || mode == "tip-changed" {
				if e == nil || before != treeState(t, r.WorkspaceRoot) {
					t.Fatal("stale preview wrote")
				}
				return
			}
			if mode == "late-report" {
				if e != nil || out.Collection.State == "qualified" {
					t.Fatalf("missing report: %v", e)
				}
				a := acceptance(t, d, r)
				if _, e = Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "late-claim.json", a)); e != nil {
					t.Fatal(e)
				}
				report := semanticReport(t, d, r)
				if _, e = SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "late-report.json", report)); e != nil {
					t.Fatal(e)
				}
				d.TaskStatusPath = ""
				preview, e = PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
				if e != nil {
					t.Fatal(e)
				}
				out, e = Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
				if e != nil || out.Collection.State == "qualified" {
					t.Fatal("old observation qualified later report")
				}
				out, e = collectSyntheticStatus(t, d, r, syntheticStatus(t, d, r, "inactive"))
				if e != nil || out.Collection.State != "qualified" {
					t.Fatalf("new observation did not qualify: %v", e)
				}
				return
			}
			if e == nil {
				t.Fatal("fault not reached")
			}
			d.Fault = nil
			preview, e = PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			out, e = Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
			if e != nil || out.Collection.State != "qualified" {
				t.Fatalf("recovery: %v %+v", e, out)
			}
			before = treeState(t, r.WorkspaceRoot)
			preview, e = PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			if _, e = Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation); e != nil {
				t.Fatal(e)
			}
			if before != treeState(t, r.WorkspaceRoot) {
				t.Fatal("identical retry wrote")
			}
			j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			count := 0
			for _, ev := range j.Events {
				if ev.Type == "task_status_observed" {
					count++
				}
			}
			if count != 1 {
				t.Fatal("duplicate status event")
			}
			info, e := os.Stat(taskStatusPath(r, digest(s)))
			if e != nil || info.Mode().Perm() != 0600 {
				t.Fatal("status file is not private")
			}
		})
	}
}

func TestReturnRepairUnknownAndNegativeClaims(t *testing.T) {
	for _, mode := range []string{"unknown-v2", "unknown-v1", "negative-v2", "negative-v1"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			runner := d.Runner.(*fakeRunner)
			runner.inside = func(s LaunchSpec) error {
				a := acceptance(t, d, r)
				b, _ := Canonical(a)
				var m map[string]any
				json.Unmarshal(b, &m)
				c := m["runtime_claim"].(map[string]any)
				if strings.HasSuffix(mode, "v2") {
					m["schema_version"] = 2
					c["model_id"] = nil
				} else {
					c["model_id"] = "unknown"
				}
				if strings.HasPrefix(mode, "negative") {
					m["acceptance"] = "unknown"
					m["issues"] = []any{map[string]any{"type": "unknown", "detail": "Synthetic missing authority observation."}}
					if strings.HasSuffix(mode, "v2") {
						c["runtime_id"] = nil
						c["profile_id"] = nil
						c["effective_policy_sha256"] = nil
						m["sandbox"] = nil
					} else {
						c["runtime_id"] = "different"
					}
				}
				path := writeAny(t, r.WorkspaceRoot, "claim-repair.json", m)
				out, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), path)
				if e != nil {
					return e
				}
				if out.Acceptance.State != m["acceptance"] {
					return fmt.Errorf("state lost: %+v", out.Acceptance)
				}
				return nil
			}
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			Start(d, f, *p.(Preview).Confirmation)
			if runner.lastErr != nil {
				t.Fatal(runner.lastErr)
			}
		})
	}
}

func TestReturnRepairExitAloneRetainsReservation(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, nil)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, e := Start(d, f, *p.(Preview).Confirmation)
	if out.Collection.State == "qualified" || e == nil {
		t.Fatalf("client exit qualified task: %v %+v", e, out)
	}
	r.RequestKey = "fixture/second-run"
	if _, e = PreviewStart(d, writeAny(t, r.WorkspaceRoot, "second-request.json", r)); e == nil {
		t.Fatal("client exit released target")
	}
}

func TestReturnRepairGeneratedContextAndMissingEvidence(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	instructions := recipientInstructions(r, *j.Binding)
	contextArg := "--context " + ShellQuote(filepath.Join(j.Binding.TempRoot, "context.json"))
	if strings.Count(instructions, contextArg) != 2 {
		t.Error("both generated callbacks need exact explicit context")
	}
	preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	for _, reason := range preview.Reasons {
		if reason.Code == "task_run_candidate_drift" {
			t.Error("missing report was described as candidate drift")
		}
	}
}

func TestReturnRepairStatusChangeInsideApplyCannotEscapeConfirmation(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, nil)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	Start(d, f, *p.(Preview).Confirmation)
	s := syntheticStatus(t, d, r, "inactive")
	d.TaskStatusPath = writeAny(t, r.WorkspaceRoot, "status.json", s)
	preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	now := d.Now()
	calls := 0
	d.Now = func() time.Time {
		calls++
		if calls == 2 {
			// The in-lock preview has read inactive, but an external writer replaces
			// the input before the final publication read. It must need a new preview.
			s.State = "active"
			writeAny(t, r.WorkspaceRoot, "status.json", s)
		}
		return now
	}
	before := treeState(t, runPaths(r).RunRoot)
	if _, e = Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation); e == nil {
		t.Fatal("unconfirmed status published")
	}
	if before != treeState(t, runPaths(r).RunRoot) {
		t.Fatal("CAS rejection wrote run evidence")
	}
}
