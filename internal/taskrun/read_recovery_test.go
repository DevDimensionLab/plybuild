package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Historical metadata matches the native startup-recovery contract. The
// synthetic executable identities are deliberately absent, because reading a
// preserved report must not inspect or authorize current executable bytes.
func deliveryInventoryRecovery(t *testing.T, f deliveryFixture) (workflowState, map[string]any) {
	t.Helper()
	s := deliveryInventoryState(t, f, "recovered")
	// Recovery records the interrupted pre-Task generation, not an unused
	// reservation. This fixture is readable by both native and history readers.
	s.Phase = "bootstrap_attempted"
	dir := filepath.Join(s.Result.Paths.RunRoot, "startup-recovery", "001")
	before, err := workflowKeep(f.D, filepath.Join(dir, "before-state.json"), s)
	if err != nil {
		t.Fatal(err)
	}
	previousTransport := s.Result.Transport
	provider := Executable{Path: filepath.Join(f.R.WorkspaceRoot, "removed-provider"), SHA256: hash([]byte("historical provider"))}
	control := Executable{Path: filepath.Join(dir, "ply-control"), SHA256: hash([]byte("historical control"))}
	context := workflowContext{Envelope: workflowEnv("run-context"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, HandoffSHA256: s.Result.Handoff.SHA256, PlyExecutable: control}
	contextBinding, err := workflowKeep(f.D, filepath.Join(dir, "context.json"), context)
	if err != nil {
		t.Fatal(err)
	}
	s.Result.Transport.TabID, s.Result.Transport.PaneID, s.Result.Transport.TerminalID = "replacement-tab", "replacement-pane", "replacement-terminal"
	transport, err := workflowKeep(f.D, filepath.Join(dir, "transport.json"), map[string]any{"workspace_id": s.Result.Transport.WorkspaceID, "tab_id": s.Result.Transport.TabID, "pane_id": s.Result.Transport.PaneID, "terminal_id": s.Result.Transport.TerminalID, "cwd": s.Observed.Target.WorktreeLocator})
	if err != nil {
		t.Fatal(err)
	}
	record, err := workflowKeep(f.D, filepath.Join(dir, "recovery.json"), map[string]any{
		"kind": "PlyDeliveryStartupRecovery@1", "schema_version": 1, "generation": 1, "run_id": s.Result.RunID, "request_sha256": s.Result.RequestSHA256,
		"before_state": before, "source_control": control, "provider_executable": provider, "control_executable": control, "context": contextBinding, "confirmation": hash([]byte("historical confirmation")), "recorded_at_utc": "2026-10-06T10:00:00Z", "transport_mode": "replace_missing_terminal", "destination_workspace_id": s.Result.Transport.WorkspaceID,
		"absent_transport": map[string]any{"transport": previousTransport, "pane_state": "not_found", "terminal_state": "absent_from_all_panes", "inventory_sha256": hash([]byte("historical transport inventory")), "process_liveness": "unknown", "task_authority": "not_issued"},
	})
	if err != nil {
		t.Fatal(err)
	}
	s.Result.Paths.Context, s.ContextSHA256 = contextBinding.Locator, contextBinding.SHA256
	deliveryInventoryReport(t, f, &s, "recovered-question", "needs_input")
	value := map[string]any{}
	raw, err := Canonical(s)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &value); err != nil {
		t.Fatal(err)
	}
	value["startup_recovery"], value["startup_recovery_transport"] = record, transport
	value["result"].(map[string]any)["startup_recovery"] = map[string]any{"generation": 1, "binding": record, "original_provider_executable": s.Request.Runtime.Executable, "original_control_executable": s.Request.Runtime.PlyExecutable, "provider_executable": provider, "control_executable": control, "transport_mode": "replace_missing_terminal", "original_transport": previousTransport, "replacement_transport": transport}
	writeInventoryRecoveryState(t, s, value)
	return s, value
}

func writeInventoryRecoveryState(t *testing.T, s workflowState, value map[string]any) {
	t.Helper()
	raw, err := Canonical(value)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), raw, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestReadInventoryPreservesDeliveryReportAfterBoundStartupRecovery(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s, _ := deliveryInventoryRecovery(t, f)
	row := deliveryInventoryRows(t, f)[s.Result.RunID]
	if row.Freshness != "fresh" || row.Delivery == nil || row.Delivery.Report == nil || !row.Delivery.Report.Current || row.Delivery.NextAction.Kind != "answer_question" || row.Delivery.NextAction.Reason != "Which default should this command use?" {
		t.Fatalf("known startup recovery metadata hid the intact delivery question: %+v delivery=%+v", row, row.Delivery)
	}
	if row.Herdr == nil || row.Herdr.TerminalID == nil || *row.Herdr.TerminalID != "replacement-terminal" || row.StateFreshness != "unknown" || row.State != "unknown" {
		t.Fatalf("recovery lost cached transport binding or claimed runtime authority: %+v", row)
	}
	if row.StartupRecovery == nil || row.StartupRecovery.Generation != 1 || row.StartupRecovery.ProviderExecutable.Path != filepath.Join(f.R.WorkspaceRoot, "removed-provider") || row.StartupRecovery.ReplacementTransport == nil {
		t.Fatalf("bound recovery metadata was not preserved: %+v", row.StartupRecovery)
	}
	for _, name := range []string{"recovery.json", "before-state.json", "context.json", "transport.json"} {
		found := false
		for _, source := range row.Sources {
			found = found || filepath.Base(source.Locator) == name
		}
		if !found {
			t.Errorf("recovery source %s was not preserved", name)
		}
	}
	raw, err := os.ReadFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var mutationState workflowState
	if err := decode(raw, 8<<20, &mutationState); err != nil {
		t.Fatalf("installed native recovery schema must remain readable: %v", err)
	}
	var unknown map[string]any
	if err := json.Unmarshal(raw, &unknown); err != nil {
		t.Fatal(err)
	}
	unknown["unrecognized_authority"] = true
	changed, err := Canonical(unknown)
	if err != nil {
		t.Fatal(err)
	}
	if decode(changed, 8<<20, &mutationState) == nil {
		t.Fatal("native recovery compatibility accepted an unknown authority field")
	}
}

func TestReadInventoryRejectsUnknownOrUnboundRecoveryMetadata(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s, base := deliveryInventoryRecovery(t, f)
	raw, err := Canonical(base)
	if err != nil {
		t.Fatal(err)
	}
	for _, condition := range []string{"unknown_state_field", "unknown_result_field", "unknown_recovery_field", "missing_record", "generation", "record_digest", "transport_digest", "original_provider", "request_identity"} {
		t.Run(condition, func(t *testing.T) {
			var changed map[string]any
			if err := json.Unmarshal(raw, &changed); err != nil {
				t.Fatal(err)
			}
			metadata := changed["result"].(map[string]any)["startup_recovery"].(map[string]any)
			switch condition {
			case "unknown_state_field":
				changed["unrecognized_authority"] = true
			case "unknown_result_field":
				changed["result"].(map[string]any)["unrecognized_authority"] = true
			case "unknown_recovery_field":
				metadata["unrecognized_authority"] = true
			case "missing_record":
				delete(changed, "startup_recovery")
			case "generation":
				metadata["generation"] = 2
			case "record_digest":
				changed["startup_recovery"].(map[string]any)["sha256"] = hash([]byte("different record"))
				metadata["binding"] = changed["startup_recovery"]
			case "transport_digest":
				changed["startup_recovery_transport"].(map[string]any)["sha256"] = hash([]byte("different transport"))
				metadata["replacement_transport"] = changed["startup_recovery_transport"]
			case "original_provider":
				metadata["original_provider_executable"].(map[string]any)["sha256"] = hash([]byte("different provider"))
			case "request_identity":
				changed["request"].(map[string]any)["request_key"] = "another/request"
			}
			writeInventoryRecoveryState(t, s, changed)
			row := readHerdrInventoryRun(f.D, f.R.WorkspaceRoot, s.Result.RunID)
			if row.Freshness != "unknown" || row.StartupRecovery != nil || row.Delivery == nil || row.Delivery.NextAction.Kind != "inspect_delivery" || row.Delivery.Report != nil && row.Delivery.Report.Current || !hasInventoryReason(row, "delivery_state_unavailable") {
				t.Fatalf("%s weakened strict source identity checks: %+v delivery=%+v", condition, row, row.Delivery)
			}
		})
	}
}
