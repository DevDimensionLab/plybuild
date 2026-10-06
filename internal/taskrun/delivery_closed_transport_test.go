package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const closedTransportInspectStandin = `
if model.get("closed_transport") and cmd == ["pane", "get"] and args[2] == model.get("closed_pane", "w-fixture:p1"):
    print(json.dumps({"id":"fixture", "error":{"code":"pane_not_found", "message":"pane missing"}}), file=sys.stderr)
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["pane", "list"]:
    print(json.dumps({"id":"fixture", "result":{"panes":model.get("inventory", [])}}))
    raise SystemExit(0)
if model.get("closed_transport") and cmd == ["workspace", "get"]:
    print(json.dumps({"id":"fixture", "result":{"workspace":{"workspace_id":model.get("destination", "w-fixture")}}}))
    raise SystemExit(0)
`

const closedTransportAgentStandin = `
if model.get("closed_transport") and cmd == ["tab", "create"]:
    model["replacement_count"] = model.get("replacement_count", 0) + 1
if model.get("replacement_count"):
    suffix = str(model["replacement_count"])
    agent.update(workspace_id=model.get("destination", "w-fixture"), tab_id="recovery-tab-"+suffix, pane_id="recovery-pane-"+suffix, terminal_id="recovery-terminal-"+suffix)
`

func closedTransportFixture(t *testing.T) (deliveryFixture, WorkflowRun, DeliveryStartRecoveryInput) {
	t.Helper()
	f := deliveryRecoveryNewFixture(t)
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+closedTransportInspectStandin), 1)
	herdr = bytes.Replace(herdr, []byte("if \"foreground_cwd\" in model:\n"), []byte(closedTransportAgentStandin+"if \"foreground_cwd\" in model:\n"), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil || o.RunID == "" {
		t.Fatalf("fixture did not stop before Task: %v", err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "recovery_armed": true})
	f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
		return nil, errors.New("closed transport must not pretend to inspect an unavailable shell")
	}
	f.D.HerdrTimeout = 3 * time.Second
	return f, o, DeliveryStartRecoveryInput{RunID: o.RunID, ProviderExecutable: f.R.Runtime.Executable, ControlExecutable: f.R.Runtime.PlyExecutable, Timeout: 3 * time.Second}
}

func TestDeliveryClosedTransportPreservesRunAndStartsOneNewTab(t *testing.T) {
	f, original, in := closedTransportFixture(t)
	before := deliveryRecoveryFiles(t, original.Paths.RunRoot)
	p := deliveryRecoveryTestPreview(t, f, in)
	if p.TransportMode != deliveryRecoveryReplaceTerminal || p.Generation != 1 || p.ExitEvidence != nil || p.AbsentTransport == nil || p.AbsentTransport.ProcessLiveness != "unknown" || p.AbsentTransport.TaskAuthority != "not_issued" {
		t.Fatalf("absence was not distinguished from exit: %+v", p)
	}
	deliveryRecoveryUnchanged(t, before, "")
	deliveryRecoveryPreserveControl(t, in, p)
	o, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || o.Transport.TerminalID == original.Transport.TerminalID || o.StartupRecovery == nil || o.StartupRecovery.ReplacementTransport == nil {
		t.Fatalf("new transport was not bound: %+v, %v", o.StartupRecovery, err)
	}
	if o.RunID != original.RunID || o.RequestSHA256 != original.RequestSHA256 || o.Handoff != original.Handoff {
		t.Fatal("replacement changed Task ownership")
	}
	deliveryRecoveryUnchanged(t, before, filepath.Join(original.Paths.RunRoot, "state.json"))
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 {
		t.Fatal("replacement did not create exactly one new tab and agent")
	}
	if _, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 {
		t.Fatal("same confirmation repeated a native effect")
	}
}

func TestDeliveryClosedTransportRequiresCompleteGlobalAbsence(t *testing.T) {
	for _, tc := range []struct {
		name      string
		inventory any
	}{
		{"moved_terminal", []any{map[string]any{"workspace_id": "other-workspace", "tab_id": "other-tab", "pane_id": "other-pane", "terminal_id": "terminal-fixture"}}},
		{"unknown_inventory", nil},
		{"incomplete_identity", []any{map[string]any{"pane_id": "unresolved"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, o, in := closedTransportFixture(t)
			workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": tc.inventory})
			before := deliveryRecoveryFiles(t, o.Paths.RunRoot)
			if p, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in); err == nil || p.State == "ready" {
				t.Fatal("unsafe inventory qualified a replacement")
			}
			deliveryRecoveryUnchanged(t, before, "")
			if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
				t.Fatal("preview created a tab")
			}
		})
	}
}

func TestDeliveryClosedTransportLostTabReplyPreservesUnknownAttempt(t *testing.T) {
	f, o, in := closedTransportFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": "tab create"})
	if _, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("lost tab response was reported as confirmed")
	}
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Phase != "recovery_tab_create_attempted" || state.RecoveryTransport != nil {
		t.Fatalf("unknown tab identity was not preserved: %v", err)
	}
	for n := 0; n < 2; n++ {
		preview, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in)
		if err != nil || preview.State != "existing" {
			t.Fatalf("unknown create could be replaced: %+v %v", preview, err)
		}
		if _, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
			t.Fatal(err)
		}
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
		t.Fatal("uncertain tab create was replayed or started a provider")
	}
}

func TestDeliveryClosedTransportCanReplaceAnotherClosedReadinessGeneration(t *testing.T) {
	f, original, in := closedTransportFixture(t)
	// The original Herdr workspace can also have closed. The new destination is
	// explicitly selected for this attempt; the frozen request remains unchanged.
	in.DestinationWorkspaceID = "current-workspace"
	workflowTestModel(t, f.workflowFixture, map[string]any{"destination": in.DestinationWorkspaceID})
	f.D.Files = deliveryRecoveryRejectTaskPrompt{}
	first := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, first)
	if _, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, first.Confirmation); err == nil {
		t.Fatal("fixture should stop before Task prompt")
	}
	state, err := workflowRead(f.R.WorkspaceRoot, original.RunID)
	if err != nil || state.Phase != "session_bound" || state.Result.Transport.AgentSessionID == "" {
		t.Fatalf("readiness was not bound: %v", err)
	}
	prior := deliveryRecoveryFiles(t, original.Paths.RunRoot)
	oldContext := state.Result.Paths.Context
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_pane": state.Result.Transport.PaneID})
	second := deliveryRecoveryTestPreview(t, f, in)
	if second.Generation != 2 || second.RecoveryDirectory == first.RecoveryDirectory {
		t.Fatalf("next explicit replacement did not append a generation: %+v", second)
	}
	deliveryRecoveryPreserveControl(t, in, second)
	if _, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, second.Confirmation); err == nil {
		t.Fatal("fixture should again stop before Task prompt")
	}
	state, err = workflowRead(f.R.WorkspaceRoot, original.RunID)
	if err != nil || state.Result.StartupRecovery.Generation != 2 || state.Result.Paths.Context == oldContext || state.Result.Transport.WorkspaceID != in.DestinationWorkspaceID {
		t.Fatalf("replacement chain is not readable: %v", err)
	}
	deliveryRecoveryUnchanged(t, prior, filepath.Join(original.Paths.RunRoot, "state.json"))
	f.D.Files = nil
	o, err := WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, original.RunID, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	f.D.Executable = func() (string, error) { return second.ControlExecutable.Path, nil }
	o = deliveryTestAccept(t, f, o)
	if o.StartupRecovery.Generation != 2 || o.Delivery.PermissionState != "recipient_confirmed_contract" {
		t.Fatal("latest generation could not accept the unchanged Task mandate")
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 3 || workflowTestCalls(t, f.workflowFixture, "agent start") != 3 {
		t.Fatal("replacement chain repeated a transport effect")
	}
}

func TestDeliveryClosedTransportIgnoresUnrelatedPaneInventoryChanges(t *testing.T) {
	f, _, in := closedTransportFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{map[string]any{"workspace_id": "another-workspace", "tab_id": "unrelated-tab", "pane_id": "unrelated-pane", "terminal_id": "unrelated-terminal"}}})
	o, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || o.StartupRecovery == nil || o.StartupRecovery.ReplacementTransport == nil {
		t.Fatalf("unrelated pane change blocked the selected replacement: %v", err)
	}
}
