package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// FS01 AC01/AC07: review is one coordinator operation on the existing Run and
// the reviewer returns with only its own judgment file.
func TestExecuteReviewIsOneCoordinatorOperation(t *testing.T) {
	root := newWorkflowExecuteCommand(taskrun.Dependencies{})
	review, _, err := root.Find([]string{"review"})
	if err != nil || review.Name() != "review" {
		t.Fatalf("missing coordinator review: %v", err)
	}
	if review.Flags().Lookup("coordinator") == nil || review.Flags().Lookup("context") != nil || review.Flags().Lookup("file") != nil {
		t.Error("review must need only the coordinator claim, no owner context or metadata file")
	}
	ret, _, err := root.Find([]string{"review-return"})
	if err != nil || ret.Name() != "review-return" {
		t.Fatalf("missing reviewer return: %v", err)
	}
	for _, flag := range []string{"context", "file"} {
		if ret.Flags().Lookup(flag) == nil {
			t.Errorf("missing review-return selector %s", flag)
		}
	}
	var help bytes.Buffer
	ret.SetOut(&help)
	if err := ret.Help(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"summary: 1-2000 Unicode codepoints", "findings[].summary: 1-600 Unicode codepoints", "256 KiB (262144 bytes)"} {
		if !strings.Contains(help.String(), want) {
			t.Errorf("review-return help lacks %q", want)
		}
	}
	verify, _, err := root.Find([]string{"verify"})
	if err != nil {
		t.Fatal(err)
	}
	verify.SetArgs(nil)
	if err := verify.Flags().Set("context", "/absolute/context.json"); err != nil {
		t.Fatal(err)
	}
	if err := verify.Args(verify, []string{"wfr_x"}); err != nil {
		t.Errorf("verify still requires a hand-written --review: %v", err)
	}
	found := map[string]string{}
	for _, op := range coreCapabilities(capabilityBuild{}).WorkflowExtensions {
		found[op.ID] = op.Effect
	}
	if found["workflow.execute.review"] != "coordinator_review_round" || found["workflow.execute.review-return"] != "reviewer_review_return" {
		t.Errorf("review capabilities missing: %v %v", found["workflow.execute.review"], found["workflow.execute.review-return"])
	}
}

func TestDeliveryStatusShowsReviewFindingsAndNextActor(t *testing.T) {
	var out bytes.Buffer
	writeDeliveryStatus(&out, &taskrun.DeliveryStatus{Source: taskrun.DeliverySourceStatus{State: "observed", OID: "1234567890abcdef", Clean: true}, Review: &taskrun.DeliveryReviewStatus{CurrentCandidate: "findings_open", Rounds: []taskrun.DeliveryReviewRoundStatus{{Round: 1, CandidateOID: "1234567890abcdef", State: "completed", Decision: "changes_requested", Current: true, Findings: []taskrun.DeliveryReviewFinding{{ID: "r001-f01", Severity: "high", Summary: "README lacks the fix"}}, EvidenceGaps: []string{"No provider journey"}}}, Next: taskrun.WorkflowAction{Actor: "owner", Message: "Correct the 1 finding(s)"}}})
	for _, want := range []string{"Review round 1 for 1234567890ab (current): completed — changes_requested", "Finding r001-f01 [high]: README lacks the fix", "Evidence gap: No provider journey", "Review next: owner — Correct the 1 finding(s)"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("status lacks %q:\n%s", want, out.String())
		}
	}
}

func TestAL33TextShowsReviewerInspection(t *testing.T) {
	var out bytes.Buffer
	writeDeliveryStatus(&out, &taskrun.DeliveryStatus{Review: &taskrun.DeliveryReviewStatus{Rounds: []taskrun.DeliveryReviewRoundStatus{{Round: 1, ReviewerStart: "start_unknown", StartupDiagnostics: &taskrun.DeliveryReviewStartupDiagnostics{StartAttempt: "attempted", SessionBinding: "not_bound", TaskPrompt: "not_attempted", ProviderLiveness: "unknown", InspectionArgv: [][]string{{"/path with space/herdr", "pane", "process-info", "--pane", "w-fixture:r1"}}}}}}})
	if !strings.Contains(out.String(), "'process-info' '--pane' 'w-fixture:r1'") || !strings.Contains(out.String(), "provider_liveness=unknown") {
		t.Fatalf("no actionable inspection: %s", out.String())
	}
}
