package taskrun

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkflowAcceptedRequiresNativeEvidencePolicy(t *testing.T) {
	f := workflowTestFixture(t)
	var verifiers []map[string]any
	if e := json.Unmarshal(f.Report.VerifierResults, &verifiers); e != nil {
		t.Fatal(e)
	}
	// The verifier did not authorize stdout capture. Shape-valid evidence must
	// still satisfy the native policy, not just the reported passed gate.
	verifiers[0]["stdout_artifact_id"] = "task-requirements"
	f.Report.VerifierResults, _ = Canonical(verifiers)
	o := workflowTestReady(t, f)
	if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted")); e == nil {
		t.Fatal("accepted bypassed native evidence coverage policy")
	}
}

type workflowRejectReportFile struct{}

func (workflowRejectReportFile) WriteOnce(p string, b []byte) error {
	if filepath.Base(p) == "report.json" {
		return errors.New("interrupted before report projection")
	}
	return writeOnce(p, b)
}
func (workflowRejectReportFile) Replace(p string, b []byte) error { return publishFile(p, b, true) }

func TestWorkflowInterruptedReportCannotReplaceItsFrozenSemantics(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	d := f.D
	d.Files = workflowRejectReportFile{}
	file := writeAny(t, f.R.WorkspaceRoot, "offered-report.json", f.Report)
	if _, e := WorkflowReportRound(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); e == nil {
		t.Fatal("fault did not interrupt")
	}
	different := f.Report
	different.Summary = "Different semantics after interrupted publication."
	if _, e := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "changed-report.json", different)); e == nil {
		t.Fatal("interrupted report accepted different semantics against the old terminal draft")
	}
	after, e := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
	if e != nil {
		t.Fatal(e)
	}
	if after.Round.State != "report_received" {
		t.Fatal("exact report recovery failed")
	}
}

func TestWorkflowSettledObservationsMustBeNoMoreThanOneSecondApart(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	o = workflowTestReport(t, f, o)
	workflowTestModel(t, f, map[string]any{"agent_status": "idle", "delay_command": "agent get", "delay_seconds": 2})
	after, e := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 5*time.Second)
	if e == nil || after.Round.State == "ready_for_review" {
		t.Fatal("widely separated observations incorrectly completed the round")
	}
}

func TestWorkflowPassedAssessmentCannotOmitRequiredVerifiers(t *testing.T) {
	f := workflowTestFixture(t)
	f.Report.TechnicalAssessment.RequiredVerifierIDs = []string{}
	o := workflowTestReady(t, f)
	file := workflowTestReview(t, f, o, "accepted")
	if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); e == nil {
		t.Fatal("passed assessment ignored the native required verifier set")
	}
}
