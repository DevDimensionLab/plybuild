package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func releaseReportCandidate(t *testing.T) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 2, Mode: workspace.DeliveryLocalBranch, IntegrationOwner: workspace.IntegrationOwnerHuman}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "qualified source handover")
	counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-count.txt")
	script, err := os.ReadFile(f.AcceptancePath)
	if err != nil {
		t.Fatal(err)
	}
	script = append(script, []byte("printf 'executed\\n' >> "+ShellQuote(counter)+"\n")...)
	if err = os.WriteFile(f.AcceptancePath, script, 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "source-handover")
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	return f, o, review, counter
}

func releaseProgressReport(t *testing.T, f deliveryFixture, o WorkflowRun, phase, id string) WorkflowRun {
	t.Helper()
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID,
		EventID: id, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: phase, Summary: "The exact qualified candidate is ready for source handover.",
		Meaning: "Progress only; the candidate, acceptance and review are unchanged. Human QA remains separate.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if phase == "needs_input" {
		r.Question = ptr("Which supported source handover should this unchanged candidate use?")
	}
	out, err := WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, id+".json", r))
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func assertReleaseCandidatePreserved(t *testing.T, f deliveryFixture, o WorkflowRun, candidate DeliveryCandidate, counter, parent string) {
	t.Helper()
	if len(o.Delivery.Candidates) != 1 || !equal(o.Delivery.Candidates[0], candidate) || o.Delivery.Attempt.ID != candidate.Key || o.Delivery.Attempt.State != "recorded" {
		t.Fatalf("handover replaced candidate, TaskResult, receipt or attempt: %+v", o.Delivery)
	}
	raw, err := os.ReadFile(counter)
	if err != nil || string(raw) != "executed\n" {
		t.Fatalf("acceptance must execute exactly once: %q %v", raw, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 1 || len(registry.HumanQARecords) != 0 {
		t.Fatalf("handover created a result or human judgment: results=%d QA=%d error=%v", len(registry.TaskResults), len(registry.HumanQARecords), err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parent {
		t.Fatal("source release integrated the candidate")
	}
}
