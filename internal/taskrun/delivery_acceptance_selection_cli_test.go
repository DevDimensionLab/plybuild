package taskrun

import (
	"os"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func acceptanceChoiceFile(t *testing.T, f deliveryFixture, o WorkflowRun, name string) string {
	t.Helper()
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	x, err := deliveryTarget(f.D, s)
	if err != nil {
		t.Fatal(err)
	}
	return writeAny(t, f.R.WorkspaceRoot, name, map[string]any{
		"kind": "DeliveryAcceptanceChoice@1", "schema_version": 1,
		"run_id": o.RunID, "request_sha256": o.RequestSHA256,
		"previous_agreement_sha256": workspace.DeliveryAgreementDigest(*f.R.Delivery.Agreement),
		"candidate_oid":             x.OID, "candidate_tree": x.Tree,
		"actor_claim": "fixture human via orchestrator", "answer": "use automatic acceptance for this delivery",
		"requested_at_utc": time.Now().UTC().Format(time.RFC3339Nano),
		"policy":           workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30},
	})
}

func TestExistingRunAutomaticAcceptanceCLI(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, false, 30)
	deliveryTestCandidate(t, f, "explicit automatic choice")
	review := deliveryTestReview(t, f, "before-choice")
	o = automaticCLIOK(t, binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review)
	old := o.Delivery.Candidates[0]
	original, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	choice := acceptanceChoiceFile(t, f, o, "automatic-choice.json")
	args := []string{"execute", "acceptance", o.RunID, "--context", o.Paths.Context, "--file", choice}
	o = automaticCLIOK(t, binary, f, args...)
	o = automaticCLIOK(t, binary, f, args...)
	if o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Mode != "automatic" || o.DeliveryStatus.Acceptance.Outcome == "pass" || o.DeliveryStatus.HumanJudgment != nil {
		t.Fatalf("selection is not a separate policy choice: %+v", o.DeliveryStatus)
	}
	if len(o.Delivery.Candidates) != 1 || !equal(old, o.Delivery.Candidates[0]) {
		t.Fatal("selection rewrote candidate history")
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || !equal(original.Request, after.Request) || original.Result.Handoff != after.Result.Handoff || !equal(original.Acceptance, after.Acceptance) || original.ContextSHA256 != after.ContextSHA256 {
		t.Fatalf("frozen authority changed: %v", err)
	}
	if _, _, err := automaticCLI(binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--reuse", old.Key); err == nil {
		t.Fatal("selection alone upgraded a historical receipt to machine pass")
	}
	if _, _, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil {
		t.Fatal("selection delivered without automatic test")
	}
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'explicit automatic choice'\n\"$PLY_CANDIDATE_BINARY\" capabilities --format json\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
	if len(o.Delivery.Candidates) != 2 || !equal(old, o.Delivery.Candidates[0]) {
		t.Fatal("new acceptance did not preserve prior candidate")
	}
	assertAutomaticCLIClosure(t, binary, f, o)
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if run.Coverage != "complete" || len(traceKind(run, "verification")) != 2 || len(traceKind(run, "acceptance_policy_choice")) != 1 || len(traceKind(run, "human_qa")) != 0 || len(traceKind(run, "task_closeout")) != 1 {
		t.Fatalf("mixed old/new history lost in trace: %+v", run.Reasons)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.IntegrationAuthorities) != 1 {
		t.Fatal("duplicate native integration authority")
	}
	auth := registry.IntegrationAuthorities[0].Plan.DeliveryAuthorization
	if auth == nil || auth.AcceptanceAmendment == nil || auth.AcceptanceAmendment.PreviousAgreementSHA256 != workspace.DeliveryAgreementDigest(*f.R.Delivery.Agreement) || auth.Agreement.SchemaVersion != 3 {
		t.Fatal("native integration lost original/effective policy distinction")
	}
}
