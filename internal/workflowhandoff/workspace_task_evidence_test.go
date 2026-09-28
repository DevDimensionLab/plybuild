package workflowhandoff

import (
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkspaceTaskEvidenceReaderProjectsExactImmutableEvidence(t *testing.T) {
	dependencies, created, _ := completeRoundTrip(t)
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	request := workspace.TaskHandoffEvidenceRequest{
		ActivityID: snapshot.Handoff.Identity.ActivityID, RunID: snapshot.Handoff.Identity.RunID, HandoffID: snapshot.Handoff.Identity.HandoffID,
		HandoffLocator: snapshot.Handoff.Locator, HandoffSHA256: snapshot.Handoff.SHA256,
		StartReceiptID: snapshot.Start.DocumentID, StartReceiptLocator: snapshot.Start.Locator, StartReceiptSHA256: snapshot.Start.SHA256,
		TerminalResultID: snapshot.Terminal.DocumentID, TerminalResultLocator: snapshot.Terminal.Locator, TerminalResultSHA256: snapshot.Terminal.SHA256,
		InspectionSHA256: digestBytes(inspection.Bytes),
	}
	reader := NewTaskHandoffEvidenceReader(dependencies)
	evidence, err := reader.ReadTaskEvidence(request)
	if err != nil {
		t.Fatal(err)
	}
	if evidence.ReportedOutcome != "complete" || !evidence.SchemaValid || !evidence.DigestValid || !evidence.LifecycleValid || !evidence.BindingValid || !evidence.CapabilityValid || !evidence.PrincipalSessionValid || !evidence.PolicyValid || !evidence.EvidenceCoverageValid {
		t.Fatalf("evidence = %#v", evidence)
	}
	if evidence.ActivityID != request.ActivityID || evidence.TerminalResultID != request.TerminalResultID || evidence.InspectionSHA256 != request.InspectionSHA256 || evidence.TargetWorktree == "" || evidence.TargetRef == "" || evidence.ResultOID == "" || evidence.ResultTree == "" || evidence.GitCommonDir == "" {
		t.Fatalf("projection = %#v", evidence)
	}
	if strings.Contains(evidence.HandoffLocator, snapshot.Handoff.ReplySecret) {
		t.Fatal("projection leaked reply secret")
	}

	request.TerminalResultSHA256 = "sha256:" + strings.Repeat("0", 64)
	if _, err := reader.ReadTaskEvidence(request); err == nil {
		t.Fatal("mismatched immutable evidence was accepted")
	}
}
