package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func resolvedHistoryStop(t *testing.T, f deliveryFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	proof := writeAny(t, f.R.WorkspaceRoot, "history-refusal.json", map[string]string{"source": "synthetic provider refusal before execution"})
	stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "history-refusal", Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "fixture coordinator", Reason: "Synthetic refusal", EvidencePath: proof, Action: "refused fixture action", EffectStatus: DeliveryEffectNotStarted})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stopped.Delivery.ActiveStop.SHA256, "Cause clarified by fixture coordinator", "Continue a different fixture action")
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

// Reproduce PE01's already integrated/released state using native writers in a
// disposable workspace, then recover only closeout across the public CLI.
func TestAutomaticDeliveryCLIClosesResolvedStopHistoryWithoutReplay(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	o = resolvedHistoryStop(t, f, o)
	history := append([]DeliveryEventRecord(nil), o.Delivery.Events...)
	deliveryTestCandidate(t, f, "resolved stop candidate")
	marker := filepath.Join(f.R.WorkspaceRoot, "acceptance-count")
	script := "#!/bin/sh\nset -eu\nprintf 'once\\n' >> " + ShellQuote(marker) + "\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, deliveryTestReview(t, f, "resolved-stop"), binary)...)
	c := o.Delivery.Candidates[0]
	d := f.D
	d.Fault = func(point string) error {
		if point == "automatic_closeout_after_release" {
			return errors.New("fixture interruption after recorded integration and release")
		}
		return nil
	}
	interrupted, err := WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err == nil || interrupted.Delivery.Phase != "closing" || interrupted.Delivery.Candidates[0].Integration == nil || interrupted.Delivery.OwnershipRelease == nil {
		t.Fatalf("fixture did not reach recorded integration/release: %+v %v", interrupted.Delivery, err)
	}
	integration := *interrupted.Delivery.Candidates[0].Integration
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	// This is the ordinary same-owner continuation, not a new delivery attempt.
	done := automaticCLIOK(t, binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if done.Delivery.Phase != "completed" || *done.Delivery.Candidates[0].Integration != integration {
		t.Fatalf("closeout lost original integration: %+v", done.Delivery)
	}
	if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != before {
		t.Fatal("closeout replayed integration")
	}
	if count, err := os.ReadFile(marker); err != nil || string(count) != "once\n" {
		t.Fatalf("closeout replayed verification: %q %v", count, err)
	}
	for i, event := range history {
		if done.Delivery.Events[i] != event {
			t.Fatal("closeout rewrote refusal/resolution")
		}
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || !closeout.LifecycleCompleted || closeout.ResourceState != "kept" {
		t.Fatalf("Task not closed: %+v %v", closeout, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.HumanQARecords) != 0 || len(registry.IntegrationResults) != 1 {
		t.Fatalf("closeout added QA or integration: %+v %v", registry, err)
	}
	if current := retireQueueCurrent(t, f); current != nil {
		t.Fatalf("queue not empty: %+v", current)
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatalf("base not advanced: %+v %v", base, err)
	}
	row := deliveryInventoryRows(t, f)[o.RunID]
	if row.Freshness != "fresh" || row.State != "completed" {
		t.Fatalf("completed history unreadable: %+v", row.Reasons)
	}
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	for i, kind := range []string{"coordinator_stop", "stop_resolution"} {
		entries := traceKind(run, kind)
		if len(entries) != 1 || entries[0].Source != history[i].Binding || len(entries[0].Data) == 0 {
			t.Fatalf("trace lost original stop history: %+v", entries)
		}
	}
}

func TestReadInventoryNativeRetirementNeedsNoAgentReport(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Selected replacement"})
	if err != nil {
		t.Fatal(err)
	}
	row := deliveryInventoryRows(t, f)[o.RunID]
	if row.Freshness != "fresh" || len(row.Reasons) != 0 || row.Delivery.Report != nil || row.Delivery.Phase != "stopped" || out.Delivery.PermissionState != o.Delivery.PermissionState {
		t.Fatalf("native retirement requires fabricated report or grants permission: %+v", row)
	}
}

// Historical fixtures exercise frozen legacy and current native field shapes.
// They carry no live permission or delivery authority.
func TestReadInventoryStopHistoryIdentityContracts(t *testing.T) {
	cases := []string{"current", "legacy", "startup-before-session", "missing-work-session", "missing-work-session-before-binding", "wrong-stop-run", "wrong-provider", "wrong-stop-session", "wrong-event", "wrong-resolution-run", "wrong-request", "wrong-resolution-session", "wrong-stop", "wrong-stop-locator", "unknown-effect", "partial-effect", "legacy-unknown-effect", "same-refused-action", "legacy-bool-with-current-fields", "corrupt", "unknown-kind"}
	for _, condition := range cases {
		t.Run(condition, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			s := deliveryInventoryState(t, f, condition)
			if condition == "missing-work-session-before-binding" {
				s.Result.Transport.AgentSessionID = ""
			}
			proof := writeAny(t, f.R.WorkspaceRoot, "history-proof.json", map[string]string{"observed": "synthetic refusal"})
			evidence, _, err := deliveryReadBinding(proof, 1<<20)
			if err != nil {
				t.Fatal(err)
			}
			stop := DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: s.Result.RunID, Input: DeliveryStopInput{EventID: "stop", Provider: "codex", NativeSession: s.Result.Transport.AgentSessionID, Stage: "work", Observer: "fixture coordinator", Reason: "Synthetic refusal", EvidencePath: proof, Action: "refused fixture action", EffectStatus: DeliveryEffectNotStarted}, Source: evidence, Evidence: evidence, ObservedAtUTC: "2026-10-05T10:00:00Z"}
			switch condition {
			case "startup-before-session":
				stop.Input.Stage = "startup"
				stop.Input.NativeSession = ""
			case "missing-work-session":
				stop.Input.NativeSession = ""
			case "wrong-stop-run":
				stop.RunID = "another-run"
			case "wrong-provider":
				stop.Input.Provider = "claude"
			case "wrong-stop-session":
				stop.Input.NativeSession = "another-session"
			case "wrong-event":
				stop.Input.EventID = "another-event"
			case "unknown-effect", "legacy-unknown-effect":
				stop.Input.EffectStatus = DeliveryEffectUnknown
			case "partial-effect":
				stop.Input.EffectStatus = DeliveryEffectPartial
			}
			binding, err := deliveryAppendEvent(f.D, &s, "stop", "coordinator_stop", stop)
			if err != nil {
				t.Fatal(err)
			}
			resolution := DeliveryStopResolution{Kind: "ply.workflow.stop-resolution", RunID: s.Result.RunID, Stop: binding, NativeSession: s.Result.Transport.AgentSessionID, Resolution: "Clarified cause", ContinuationAction: "Different fixture action", RequestSHA256: s.Result.RequestSHA256}
			switch condition {
			case "legacy", "legacy-unknown-effect":
				resolution.PermissionConfirmed = true
				resolution.ContinuationAction = ""
				resolution.RequestSHA256 = ""
			case "wrong-resolution-run":
				resolution.RunID = "another-run"
			case "wrong-request":
				resolution.RequestSHA256 = hash([]byte("another-request"))
			case "wrong-resolution-session":
				resolution.NativeSession = "another-session"
			case "wrong-stop-locator":
				resolution.Stop.Locator += ".other"
			case "wrong-stop":
				resolution.Stop.SHA256 = hash([]byte("another-stop"))
			case "same-refused-action":
				resolution.ContinuationAction = " refused   fixture action "
			case "legacy-bool-with-current-fields":
				resolution.PermissionConfirmed = true
				resolution.RequestSHA256 = hash([]byte("another-request"))
			}
			if _, err = deliveryAppendEvent(f.D, &s, "resolution", "stop_resolution", resolution); err != nil {
				t.Fatal(err)
			}
			if condition == "corrupt" {
				if err = os.WriteFile(binding.Locator, []byte("broken"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if condition == "unknown-kind" {
				s.Result.Delivery.Events[0].Kind = "future_stop"
			}
			s.Result.Delivery.Phase = "working"
			if err = workflowSave(f.D, s); err != nil {
				t.Fatal(err)
			}
			row := deliveryInventoryRows(t, f)[s.Result.RunID]
			valid := condition == "current" || condition == "legacy" || condition == "startup-before-session"
			if valid && (row.Freshness != "fresh" || len(row.Reasons) != 0) {
				t.Fatalf("supported history rejected: %+v", row.Reasons)
			}
			if !valid && row.Freshness != "unknown" {
				t.Fatalf("invalid history accepted: %+v", row)
			}
			if row.Delivery.Report != nil || row.State == "completed" {
				t.Fatal("history conferred an outcome")
			}
		})
	}
}

func TestReadInventoryRejectsRetirementIdentityMismatch(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Selected replacement"})
	if err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	event := out.Delivery.Events[len(out.Delivery.Events)-1]
	raw, err := workflowBound(event.Binding, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	for _, condition := range []string{"run", "request", "session", "stop", "schema", "kind"} {
		t.Run(condition, func(t *testing.T) {
			var record DeliveryRetirement
			if err := decode(raw, 1<<20, &record); err != nil {
				t.Fatal(err)
			}
			switch condition {
			case "run":
				record.RunID = "another-run"
			case "request":
				record.RequestSHA256 = hash([]byte("another-request"))
			case "session":
				record.NativeSession = "another-session"
			case "stop":
				record.ActiveStop = nil
			case "schema":
				record.SchemaVersion++
			case "kind":
				record.Kind = "future-retirement"
			}
			changed, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(event.Binding.Locator, changed, 0600); err != nil {
				t.Fatal(err)
			}
			s.Result.Delivery.Events[len(s.Result.Delivery.Events)-1].Binding.SHA256 = hash(changed)
			s.Result.Delivery.LastEventSHA256 = ptr(hash(changed))
			if err = workflowSave(f.D, s); err != nil {
				t.Fatal(err)
			}
			row := deliveryInventoryRows(t, f)[o.RunID]
			if row.Freshness != "unknown" || !hasInventoryReason(row, "delivery_event_invalid") {
				t.Fatalf("mismatched retirement accepted: %+v", row)
			}
		})
	}
}
