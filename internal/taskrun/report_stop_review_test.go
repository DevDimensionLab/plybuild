package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAL01FollowConsumesPreservedStopBeforeOwnerPolling(t *testing.T) {
	for _, exited := range []bool{false, true} {
		t.Run(map[bool]string{false: "owner_idle", true: "owner_exited"}[exited], func(t *testing.T) {
			f, o := al01AbsentReportFixture(t)
			source := mustReportWorktree(t, f, o)
			f.D.CWD = func() (string, error) { return source, nil }
			r := al01StopReport(f)
			r.RequestSHA256 = o.RequestSHA256
			if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "follow-stop.json", r)); err != nil {
				t.Fatal(err)
			}
			workflowTestModel(t, f, map[string]any{"closed_transport": exited, "agent_status": "idle"})
			f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
			calls := workflowTestCalls(t, f, "agent get")
			got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 100*time.Millisecond)
			if err != nil {
				t.Fatalf("preserved stop must return without polling exited/idle owner: %v", err)
			}
			if got.Round.State != "stopped" || got.Round.ReportSHA256 != nil || got.FinalReturn.TerminalSHA256 != nil || got.Review.Decision != nil {
				t.Fatal("stop became successful working-round review")
			}
			if workflowTestCalls(t, f, "agent get") != calls {
				t.Fatal("preserved stop polled live transport")
			}
		})
	}
}

func TestAL01CLIUnknownReportEffectsCanBeReconciledThenRetired(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := al01AbsentReportFixture(t)
	evidence := writeAny(t, f.R.WorkspaceRoot, "initial-unknown-denial.json", map[string]any{"kind": "synthetic_runtime_observation", "known": "callback result lost; prior effect unknown"})
	in := DeliveryStopInput{EventID: "unknown-report", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "synthetic coordinator", Reason: "Initial callback result unknown", Action: "round-report callback", EffectStatus: DeliveryEffectUnknown, EvidencePath: evidence}
	if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	original := mustReadAL01(t, s.ReportProviderStop.Locator)
	source, err := workflowTarget(f.D, s)
	if err != nil {
		t.Fatal(err)
	}
	// This is a deliberately synthetic action-result proof, not writer absence.
	observed := writeAny(t, f.R.WorkspaceRoot, "native-action-result.json", map[string]any{"kind": "synthetic_native_callback_result", "request_accepted": false, "effect_started": false, "result": "rejected before report publication"})
	args := []string{"run", "resolve-stop", o.RunID, "--stop", s.ReportProviderStop.SHA256, "--source", source, "--effect-status", DeliveryEffectNotStarted, "--basis", "action_result", "--observer", "synthetic coordinator", "--reason", "Actual fixture callback result confirms pre-effect rejection", "--evidence", observed}
	al01CLIOK(t, binary, f, args...)
	al01CLIOK(t, binary, f, args...)
	if string(mustReadAL01(t, s.ReportProviderStop.Locator)) != string(original) {
		t.Fatal("resolution rewrote original refusal")
	}
	retired := al01CLIOK(t, binary, f, "run", "retire", o.RunID, "--observer", "synthetic coordinator", "--reason", "Replace after sourced no-effect reconciliation")
	if retired.FinalReturn.State != "retired" || retired.ReportStartup.Retirement == nil {
		t.Fatal("reconciled absent owner could not retire")
	}
	// The same exact transition remains idempotent after its reply is lost.
	al01CLIOK(t, binary, f, "run", "retire", o.RunID, "--observer", "synthetic coordinator", "--reason", "Replace after sourced no-effect reconciliation")
}

func TestAL01ReportEffectReconciliationPreservesUncertaintyAndSource(t *testing.T) {
	for _, which := range []string{"absence_only", "wrong_stop", "wrong_source", "original_partial", "observed_partial", "source_changed", "durable_retry"} {
		t.Run(which, func(t *testing.T) {
			f, o := al01AbsentReportFixture(t)
			initialStatus := DeliveryEffectUnknown
			if which == "original_partial" {
				initialStatus = DeliveryEffectPartial
			}
			initial := writeAny(t, f.R.WorkspaceRoot, "effect-unknown.json", map[string]any{"source": "synthetic native action transport", "effect_status": initialStatus})
			if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "effect-stop", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "fixture coordinator", Reason: "Preserve actual known uncertainty", Action: "report publication", EffectStatus: initialStatus, EvidencePath: initial}); err != nil {
				t.Fatal(err)
			}
			s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
			source, _ := workflowTarget(f.D, s)
			proof := writeAny(t, f.R.WorkspaceRoot, "effect-observation.json", map[string]any{"source": "synthetic native result", "accepted": false, "effect_started": false})
			in := ReportStopResolutionInput{StopSHA256: s.ReportProviderStop.SHA256, SourceSHA256: source, EffectStatus: DeliveryEffectNotStarted, Basis: "action_result", Observer: "fixture coordinator", Reason: "Observed actual pre-effect rejection", EvidencePath: proof}
			switch which {
			case "absence_only":
				in.Basis = "writer_absence"
			case "wrong_stop":
				in.StopSHA256 = "sha256:" + strings.Repeat("0", 64)
			case "wrong_source":
				in.SourceSHA256 = "sha256:" + strings.Repeat("0", 64)
			case "observed_partial":
				in.EffectStatus = DeliveryEffectPartial
			case "durable_retry":
				f.D.Fault = func(at string) error {
					if at == "report_after_stop_resolution" {
						return fmt.Errorf("synthetic lost publication reply")
					}
					return nil
				}
			}
			_, err := WorkflowResolveReportStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if which == "absence_only" || which == "wrong_stop" || which == "wrong_source" || which == "original_partial" {
				if err == nil {
					t.Fatal("invalid proof erased uncertainty")
				}
				return
			}
			if which == "durable_retry" {
				if err == nil {
					t.Fatal("fault not reached")
				}
				f.D.Fault = nil
				if _, err = WorkflowResolveReportStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
					t.Fatal(err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if which == "source_changed" {
				if err = os.WriteFile(filepath.Join(s.Observed.Target.WorktreeLocator, "unobserved-change"), []byte("changed after proof"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			retirement := DeliveryRetireInput{Observer: "fixture coordinator", Reason: "Selected source-bound replacement"}
			_, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, retirement)
			if which == "observed_partial" || which == "source_changed" {
				if err == nil {
					t.Fatal("partial/stale effect proof released the owner")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			// A lost resolution reply is a readback after retirement, not a new effect.
			if _, err = WorkflowResolveReportStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAL01FollowConsumesCoordinatorStopWithoutOwner(t *testing.T) {
	f, o := al01AbsentReportFixture(t)
	evidence := writeAny(t, f.R.WorkspaceRoot, "follow-provider-stop.json", map[string]any{"source": "synthetic denial before callback"})
	if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "follow-stop", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "fixture coordinator", Reason: "Known blocked callback", Action: "report", EffectStatus: DeliveryEffectNotStarted, EvidencePath: evidence}); err != nil {
		t.Fatal(err)
	}
	before := workflowTestCalls(t, f, "agent get")
	got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 100*time.Millisecond)
	if err != nil || got.Round.State != "stopped" || got.TaskResultState != "not_published" || got.Round.ReportSHA256 != nil {
		t.Fatalf("coordinator stop was lost or qualified: %+v %v", got, err)
	}
	if workflowTestCalls(t, f, "agent get") != before {
		t.Fatal("follow queried exited provider after its preserved stop")
	}
}
