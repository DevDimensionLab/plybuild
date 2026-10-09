package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Exercise the real historical writer, installed repair command, and unchanged
// old reader. All provider authority and human answers belong to an isolated
// synthetic fixture; this test grants no human approval to the repair itself.
func TestDeliveryQARecoveryNativeOldControlCompatibility(t *testing.T) {
	current := deliveryCLIBinary(t)
	old := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalEpic}})
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	control := filepath.Join(f.R.WorkspaceRoot, "frozen-original-ply-control")
	raw, err := os.ReadFile(old)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(control, raw, 0500); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.PlyExecutable = Executable{control, hashFileTest(t, control)}
	f.D.Executable = func() (string, error) { return control, nil }
	f.File = writeAny(t, f.R.WorkspaceRoot, "historical-qa-request.json", f.R)
	var start WorkflowPreview
	if err = json.Unmarshal(continuityCLIOK(t, control, cwd, "workflow", "run", "start", "--file", f.File, "--check"), &start); err != nil || start.Confirmation == nil {
		t.Fatalf("original control did not preview its own run: %+v %v", start, err)
	}
	o := continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "run", "start", "--file", f.File, "--apply", "--confirm", *start.Confirmation))
	acceptance := writeAny(t, f.R.WorkspaceRoot, "synthetic-historical-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance))
	counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-executions.txt")
	oid := continuityCandidate(t, f, "synthetic historical QA candidate", counter)
	oreview := deliveryTestReview(t, f, "historical-qa")
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", oreview))
	if len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != oid || o.Delivery.Candidates[0].HumanQA != nil {
		t.Fatalf("original control did not qualify exactly one candidate: %+v", o.Delivery)
	}
	session := o.Transport.AgentSessionID
	corrected := deliveryTestHuman(t, f, o, "pass")
	var original workflowhandoff.DeliveryHumanAttestation
	if err = readValue(corrected, 1<<20, &original); err != nil {
		t.Fatal(err)
	}
	original.StartedAtUTC = strings.TrimSuffix(original.StartedAtUTC, "Z") + "+00:00"
	original.CompletedAtUTC = strings.TrimSuffix(original.CompletedAtUTC, "Z") + "+00:00"
	rejected := writeAny(t, f.R.WorkspaceRoot, "rejected-synthetic-human-answer.json", original)
	raw, err = continuityCLI(control, cwd, "workflow", "execute", "qa", o.RunID, "pass", "--context", o.Paths.Context, "--evidence", rejected)
	if err == nil || !strings.Contains(err.Error(), "invalid human QA actor") {
		t.Fatalf("historical CLI did not reproduce the prepublication timestamp rejection: %v\n%s", err, raw)
	}
	o = continuityRun(t, raw)
	if o.Delivery.Attempt == nil || o.Delivery.Attempt.Kind != "human_qa" || o.Delivery.Attempt.State != "attempted" || o.Delivery.Candidates[0].HumanQA != nil {
		t.Fatalf("historical rejection did not preserve its unresolved reservation: %+v", o.Delivery)
	}
	attempt := *o.Delivery.Attempt
	if _, err = continuityCLI(control, cwd, "workflow", "execute", "qa", o.RunID, "pass", "--context", o.Paths.Context, "--evidence", corrected); err == nil || !strings.Contains(err.Error(), "another delivery effect is unresolved") {
		t.Fatalf("historical control did not block the normalized same answer before repair: %v", err)
	}
	handoff, err := f.D.Workflow.Store.ReadByLocator(o.Delivery.Candidates[0].TaskResult.HandoffLocator)
	if err != nil {
		t.Fatal(err)
	}
	staging := filepath.Join(handoff.Handoff.ReplyRoot, "staging", "human-qa-"+hashFileTest(t, rejected)[7:])
	protected := map[string]string{}
	for _, path := range []string{control, f.File, workflowIndex(f.R.WorkspaceRoot, o.RunID), o.Paths.Context, o.Handoff.Locator, filepath.Join(o.Paths.RunRoot, "acceptance.json"), rejected, corrected, filepath.Join(staging, "attestation.json"), filepath.Join(staging, "qa-draft.json"), filepath.Join(attempt.Path, "attempt.json")} {
		protected[path] = hashFileTest(t, path)
	}
	observation := deliveryTestRuntimeObservation(t, f, o)
	args := []string{"workflow", "execute", "recover-qa", o.RunID, "--attempt", attempt.ID, "--context", o.Paths.Context, "--runtime-evidence", observation}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	stateBefore := hashFileTest(t, statePath)
	var preview DeliveryQARecoveryPreview
	if err = json.Unmarshal(continuityCLIOK(t, current, cwd, append(args, "--check")...), &preview); err != nil || preview.State != "ready" || preview.Confirmation == "" || preview.Receipt == nil {
		t.Fatalf("real repair CLI did not preview known rejection: %+v %v", preview, err)
	}
	if hashFileTest(t, statePath) != stateBefore {
		t.Fatal("read-only recovery preview changed original state")
	}
	if _, err = os.Stat(preview.Receipt.Locator); !os.IsNotExist(err) {
		t.Fatalf("preview published its recovery receipt: %v", err)
	}
	var recovered DeliveryQARecoveryPreview
	if err = json.Unmarshal(continuityCLIOK(t, current, cwd, append(args, "--confirm", preview.Confirmation)...), &recovered); err != nil || recovered.State != "recovered" || recovered.Receipt == nil || *recovered.Receipt != *preview.Receipt {
		t.Fatalf("real repair CLI did not apply its exact preview: %+v %v", recovered, err)
	}
	var receipt deliveryQARejectionReceipt
	if err = readValue(recovered.Receipt.Locator, 8<<20, &receipt); err != nil || receipt.OriginalControl != f.R.Runtime.PlyExecutable || receipt.RepairExecutable.SHA256 != hashFileTest(t, current) || receipt.RepairExecutable.SHA256 == receipt.OriginalControl.SHA256 {
		t.Fatalf("repair changed the frozen callback identity: %+v %v", receipt, err)
	}
	protected[recovered.Receipt.Locator] = recovered.Receipt.SHA256

	// This decode is performed by the real old binary, not the new Go types.
	// It must accept the standard report event and retired (nil) reservation.
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "show", o.RunID))
	if o.Delivery.Attempt != nil || o.Delivery.Candidates[0].HumanQA != nil || o.Delivery.Phase != "awaiting_human_qa" {
		t.Fatalf("old reader confused reservation recovery with human QA: %+v", o.Delivery)
	}
	if _, err = continuityCLI(control, cwd, "workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil {
		t.Fatal("recovery supplied a human pass to the frozen control")
	}
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "qa", o.RunID, "pass", "--context", o.Paths.Context, "--evidence", corrected))
	qa := o.Delivery.Candidates[0].HumanQA
	if qa == nil || qa.Outcome != "pass" || len(qa.Evidence) != 1 || qa.Evidence[0].SHA256 != protected[corrected] {
		t.Fatalf("frozen control failed to publish the normalized synthetic answer: %+v", qa)
	}
	var published workflowhandoff.DeliveryHumanAttestation
	if err = readValue(qa.Evidence[0].Locator, 1<<20, &published); err != nil {
		t.Fatal(err)
	}
	normalized := original
	normalized.StartedAtUTC = strings.TrimSuffix(original.StartedAtUTC, "+00:00") + "Z"
	normalized.CompletedAtUTC = strings.TrimSuffix(original.CompletedAtUTC, "+00:00") + "Z"
	if published != normalized {
		t.Fatal("repair or old control altered the original synthetic answer or candidate provenance")
	}
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context))
	if o.Delivery.Phase != "completed" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != oid {
		t.Fatalf("old control could not finish its synthetic local integration: %+v", o.Delivery)
	}
	stateBefore = hashFileTest(t, statePath)
	var retry DeliveryQARecoveryPreview
	if err = json.Unmarshal(continuityCLIOK(t, current, cwd, append(args, "--confirm", preview.Confirmation)...), &retry); err != nil || retry.State != "existing" || retry.Receipt == nil || *retry.Receipt != *recovered.Receipt || hashFileTest(t, statePath) != stateBefore {
		t.Fatalf("late recovery retry changed completed original-control state: %+v %v", retry, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 1 || len(registry.HumanQARecords) != 1 {
		t.Fatalf("repair duplicated candidate or human QA publications: %v", err)
	}
	if continuityCounter(t, counter) != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || o.Transport.AgentSessionID != session {
		t.Fatal("repair restarted the provider, changed its session, or repeated candidate verification")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Continuation != nil || s.Request.Runtime.PlyExecutable != f.R.Runtime.PlyExecutable {
		t.Fatalf("repair replaced the original callback control: %v", err)
	}
	for path, want := range protected {
		if hashFileTest(t, path) != want {
			t.Fatalf("repair or later callbacks changed preserved bytes: %s", path)
		}
	}
	t.Log("Synthetic CLI fixture: real 6472bb87 QA rejection reproduced; installed repair check/apply preserved all rejected bytes; unchanged historical control recorded only the normalized same answer and completed local integration, with one provider start and one verifier execution.")
}
