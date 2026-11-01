package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestWorkflowReportAcceptancePreservesLimitations(t *testing.T) {
	for _, name := range []string{"evidence_gap", "failed_assessment", "budget_exceeded"} {
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			switch name {
			case "evidence_gap":
				f.Report.EvidenceGaps = json.RawMessage(`[{"type":"known_baseline_failure","detail":"Existing suite fails on unchanged base; targeted evidence is preserved.","artifact_ids":[],"effect_ids":[]}]`)
			case "failed_assessment":
				f.Report.TechnicalAssessment.Gate = "failed"
			case "budget_exceeded":
				f.Report.BudgetUsage.ActiveSeconds = f.R.Agreement.MaxActiveSeconds + 1
			}
			o := workflowTestReady(t, f)
			before, err := os.ReadFile(*o.Round.ReportPath)
			if err != nil {
				t.Fatal(err)
			}
			review := workflowTestReview(t, f, o, "accepted")
			got, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, review)
			if err != nil {
				t.Fatalf("coordinator cannot accept truthful report: %v", err)
			}
			after, err := os.ReadFile(*got.Round.ReportPath)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("acceptance rewrote evidence: %v", err)
			}
			if got.FinalReturn.State != "accepted" || got.FinalReturn.TerminalSHA256 == nil || got.TaskResultState != "not_published" || got.Delivery != nil {
				t.Fatalf("wrong report-only result: %+v", got)
			}
			again, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, review)
			if err != nil || !equal(got.FinalReturn, again.FinalReturn) {
				t.Fatalf("retry: %v", err)
			}
		})
	}
}

func TestWorkflowReportAcceptancePreservesIncompleteDirtyCandidate(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := al10CheckedIn(t, f)
	worktree := mustReportWorktree(t, f, o)
	if err := os.WriteFile(worktree+"/partial.txt", []byte("unfinished\n"), 0600); err != nil {
		t.Fatal(err)
	}
	report := al10Report(t, f, o, "unfinished", "blocked", "failed", nil, nil)
	o, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, report)
	if err != nil {
		t.Fatal(err)
	}
	o = workflowTestSettle(t, f, o)
	got, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || got.FinalReturn.State != "accepted" || got.TaskResultState != "not_published" {
		t.Fatalf("acceptance: %+v %v", got, err)
	}
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	stored, err := workflowReadReport(state.Records[0])
	if err != nil || stored.Outcome != "blocked" || stored.TechnicalAssessment.Gate != "failed" {
		t.Fatalf("unfinished outcome changed: %+v %v", stored, err)
	}
	content, err := os.ReadFile(worktree + "/partial.txt")
	if err != nil || string(content) != "unfinished\n" {
		t.Fatalf("user work changed: %q %v", content, err)
	}
}

func TestWorkflowReportAcceptanceStillBindsDecisionAndSession(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	file := workflowTestReview(t, f, o, "accepted")
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var original WorkflowReview
	if err = json.Unmarshal(b, &original); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"report_hash", "coordinator", "request_hash", "handoff_hash"} {
		r := original
		switch name {
		case "report_hash":
			r.ReportSHA256 = hash([]byte("different report"))
		case "coordinator":
			r.Reviewer = "another-coordinator"
		case "request_hash":
			r.RequestSHA256 = hash([]byte("different request"))
		case "handoff_hash":
			r.HandoffSHA256 = hash([]byte("different handoff"))
		}
		if _, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, writeAny(t, f.R.WorkspaceRoot, name+".json", r)); err == nil {
			t.Fatalf("accepted wrong %s", name)
		}
	}
	workflowTestModel(t, f, map[string]any{"agent_session_id": "replacement-session"})
	if _, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); err == nil {
		t.Fatal("accepted replacement session")
	}
	workflowTestModel(t, f, map[string]any{"agent_session_id": o.Transport.AgentSessionID})
	if _, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); err != nil {
		t.Fatal(err)
	}
	original.Decision = "blocked"
	original.Findings = []WorkflowFinding{{"A later decision", "Original remains immutable", "Reject replacement"}}
	if _, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, writeAny(t, f.R.WorkspaceRoot, "replace.json", original)); err == nil {
		t.Fatal("replaced sealed decision")
	}
}
