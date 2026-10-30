package taskrun

import (
	"encoding/json"

	"os"

	"testing"
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
