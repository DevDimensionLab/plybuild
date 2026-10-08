package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeliveryProviderUpgradePreservesNativeLegacyContinuation(t *testing.T) {
	// Both builds must precede the fixture's process cwd change. The legacy
	// proof is published by the real pinned reader, never assembled with the
	// current structs or by substituting its reader-capability string.
	current := deliveryCLIBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "")
	f := newDeliveryFixture(t, "codex")
	oldBytes, err := os.ReadFile(before)
	if err != nil {
		t.Fatal(err)
	}
	originalControl := f.R.Runtime.PlyExecutable.Path
	if err = os.WriteFile(originalControl, oldBytes, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.PlyExecutable.SHA256 = hashFileTest(t, originalControl)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")

	var start WorkflowPreview
	if err = json.Unmarshal(continuityCLIOK(t, originalControl, cwd, "workflow", "run", "start", "--file", f.File, "--check"), &start); err != nil || start.Confirmation == nil {
		t.Fatalf("pinned reader could not preview its native launch: %+v %v", start, err)
	}
	o := continuityRun(t, continuityCLIOK(t, originalControl, cwd, "workflow", "run", "start", "--file", f.File, "--apply", "--confirm", *start.Confirmation))
	acceptance := writeAny(t, f.R.WorkspaceRoot, "legacy-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = continuityRun(t, continuityCLIOK(t, originalControl, cwd, "workflow", "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance))
	legacyControl := continuityControl(t, continuityCLIOK(t, before, cwd, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context))
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Continuation == nil || deliveryContinuationGeneration(s) != 1 {
		t.Fatalf("pinned reader did not publish exactly one native continuation: %+v %v", s.Continuation, err)
	}
	var legacy deliveryContinuationRecord
	if err = readValue(s.Continuation.Locator, 1<<20, &legacy); err != nil {
		t.Fatal(err)
	}
	const expectedLegacyCapability = "delivery-owner-v2/task-spec-v3/task-content-v2/scoped-publications-v1"
	if legacy.ReaderCapability != expectedLegacyCapability || legacy.RuntimeObservation != nil || legacy.ControlExecutable.Path != legacyControl || legacy.PreviousControl != f.R.Runtime.PlyExecutable || legacy.SourceControl.SHA256 != hashFileTest(t, before) {
		t.Fatalf("pinned reader did not supply the actual legacy proof: %+v", legacy)
	}
	protected := map[string]string{}
	for _, path := range []string{originalControl, legacyControl, s.Continuation.Locator, legacy.BeforeState.Locator, workflowIndex(f.R.WorkspaceRoot, o.RunID), o.Paths.Context, o.Handoff.Locator, filepath.Join(o.Paths.RunRoot, "mandate.json"), s.Acceptance.Locator} {
		protected[path] = hashFileTest(t, path)
	}

	probe := providerUpgradeFixture(t, f, "removed")
	if raw, err := continuityCLI(before, cwd, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context, "--check"); err == nil || !strings.Contains(err.Error()+string(raw), "no such file or directory") {
		t.Fatalf("pinned reader did not reproduce its historical launcher dependency: %v\n%s", err, raw)
	}
	if raw, err := continuityCLI(current, cwd, "workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context, "--check"); err == nil || !strings.Contains(err.Error()+string(raw), "delivery_runtime_observation_required") {
		t.Fatalf("legacy continuation was treated as current runtime evidence: %v\n%s", err, raw)
	}
	observation := deliveryTestRuntimeObservation(t, f, o)
	args := []string{"workflow", "execute", "continue", o.RunID, "--context", o.Paths.Context, "--runtime-evidence", observation}
	control := continuityControl(t, continuityCLIOK(t, current, cwd, args...))
	if again := continuityControl(t, continuityCLIOK(t, current, cwd, args...)); again != control {
		t.Fatal("repeated compatible continuation changed the preserved control")
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || deliveryContinuationGeneration(s) != 2 || s.Result.Transport.AgentSessionID != o.Transport.AgentSessionID || s.Result.SessionID != o.SessionID {
		t.Fatalf("compatible reader lost the legacy generation or native owner: generation=%d err=%v", deliveryContinuationGeneration(s), err)
	}
	continued, err := workflowContinuationRecord(s)
	if err != nil || continued == nil || continued.ReaderCapability != DeliveryContinuationReaderCapability || continued.RuntimeObservation == nil || continued.RuntimeObservation.Locator != observation || continued.PreviousControl != legacy.ControlExecutable || continued.ControlExecutable.Path != control {
		t.Fatalf("current proof did not extend the exact legacy proof: %+v %v", continued, err)
	}
	if runtime, err := workflowEffectiveRuntime(s); err != nil || runtime.Executable != f.R.Runtime.Executable || runtime.PlyExecutable.Path != control {
		t.Fatalf("legacy history substituted a provider or lost the new control: %+v %v", runtime, err)
	}

	candidate := deliveryTestCandidate(t, f, "legacy continuation after launcher cleanup")
	review := deliveryTestReview(t, f, "legacy-continuation-upgrade")
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review))
	if o.Delivery.Phase != "awaiting_human_qa" || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != candidate || o.Delivery.Candidates[0].HumanQA != nil {
		t.Fatalf("legacy continuation failed qualification or implied human QA: %+v", o.Delivery)
	}
	if _, err := continuityCLI(control, cwd, "workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil {
		t.Fatal("legacy continuation supplied missing human pass authority")
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 1 || len(registry.HumanQARecords) != 0 || gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatalf("continuation fabricated human QA or integrated its candidate: results=%d qa=%d err=%v", len(registry.TaskResults), len(registry.HumanQARecords), err)
	}
	for path, want := range protected {
		if got := hashFileTest(t, path); got != want {
			t.Fatalf("compatible continuation rewrote preserved legacy artifact %s", path)
		}
	}
	providerUpgradeNoRestart(t, f, probe)
	t.Logf("Pinned reader %s published legacy generation 1; current reader preserved it and qualified generation 2's candidate in the same native session without a human pass", continuityBeforeUpgradeSource)
}
