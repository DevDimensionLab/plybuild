package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func workflowCandidateReport(t *testing.T, f workflowFixture, round int) WorkflowReport {
	t.Helper()
	r := f.Report
	cwd, _ := f.D.CWD()
	oid := gitOutput(t, cwd, "rev-parse", "HEAD")
	tree := gitOutput(t, cwd, "rev-parse", "HEAD^{tree}")
	var artifacts []map[string]any
	if e := json.Unmarshal(r.Artifacts, &artifacts); e != nil {
		t.Fatal(e)
	}
	p := artifacts[0]["locator"].(string)
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	var req map[string]any
	if e = json.Unmarshal(b, &req); e != nil {
		t.Fatal(e)
	}
	req["result_oid"], req["result_tree"] = oid, tree
	p = writeAny(t, filepath.Dir(p), fmt.Sprintf("requirements-candidate-%d.json", round), req)
	b, _ = os.ReadFile(p)
	artifacts[0]["locator"], artifacts[0]["sha256"], artifacts[0]["size_bytes"] = p, hash(b), len(b)
	r.Artifacts, _ = Canonical(artifacts)
	var verifiers []map[string]any
	json.Unmarshal(r.VerifierResults, &verifiers)
	verifiers[0]["bound_oid_or_sha256"] = oid
	r.VerifierResults, _ = Canonical(verifiers)
	r.BudgetUsage = &BudgetUsage{true, round, 10 + round*10, 0}
	return r
}
func TestWorkflowCommittedCandidateAndDirtyCorrection(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	cwd, _ := f.D.CWD()
	os.WriteFile(filepath.Join(cwd, "README.md"), []byte("candidate change\n"), 0600)
	runGit(t, cwd, "add", "README.md")
	runGit(t, cwd, "commit", "-m", "fixture candidate")
	f.Report = workflowCandidateReport(t, f, 0)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(cwd, "README.md"), []byte("dirty correction\n"), 0600)
	f.Report = workflowCandidateReport(t, f, 1)
	f.Report.Outcome = "blocked"
	f.Report.TechnicalAssessment.Gate = "failed"
	f.Report.StopReasons = json.RawMessage(`[{"type":"repository_state_unexpected","detail":"Candidate needs its final commit.","primary":true}]`)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, e = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
	if e != nil {
		t.Fatal("dirty correction was not reviewable:", e)
	}
	runGit(t, cwd, "add", "README.md")
	runGit(t, cwd, "commit", "-m", "fixture final correction")
	f.Report = workflowCandidateReport(t, f, 2)
	f.Report.Outcome = "complete"
	f.Report.TechnicalAssessment.Gate = "passed"
	f.Report.StopReasons = json.RawMessage("[]")
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, e = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if e != nil {
		t.Fatal(e)
	}
	if o.FinalReturn.State != "accepted" || o.Budget.Used.CorrectionRounds != 2 {
		t.Fatal("changed candidate journey did not seal")
	}
}
func TestWorkflowDirtyContentDriftAndSymlinkAreNotTrusted(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	cwd, _ := f.D.CWD()
	os.WriteFile(filepath.Join(cwd, "README.md"), []byte("dirty before report\n"), 0600)
	os.WriteFile(filepath.Join(cwd, "untracked.txt"), []byte("first\n"), 0600)
	f.Report.Outcome = "blocked"
	f.Report.TechnicalAssessment.Gate = "failed"
	f.Report.StopReasons = json.RawMessage(`[{"type":"repository_state_unexpected","detail":"Dirty candidate needs correction.","primary":true}]`)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	file := workflowTestReview(t, f, o, "changes_requested")
	os.WriteFile(filepath.Join(cwd, "untracked.txt"), []byte("other\n"), 0600)
	if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); e == nil {
		t.Fatal("untracked content drift was not detected")
	}
	os.WriteFile(filepath.Join(cwd, "untracked.txt"), []byte("first\n"), 0600)
	os.WriteFile(filepath.Join(cwd, "README.md"), []byte("dirty after report!\n"), 0600)
	if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); e == nil {
		t.Fatal("same-status tracked content drift was not detected")
	}
	if e := os.Symlink(f.File, filepath.Join(f.R.WorkspaceRoot, "request-link.json")); e != nil {
		t.Fatal(e)
	}
	if _, e := WorkflowPreviewStart(f.D, filepath.Join(f.R.WorkspaceRoot, "request-link.json")); e == nil {
		t.Fatal("symlink request accepted")
	}
	outside := filepath.Join(f.R.WorkspaceRoot, "outside.txt")
	os.WriteFile(outside, []byte("unchanged\n"), 0600)
	os.Symlink(outside, filepath.Join(cwd, "outside-link"))
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	one, e := workflowTarget(f.D, s)
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(outside, []byte("changed outside target\n"), 0600)
	two, e := workflowTarget(f.D, s)
	if e != nil {
		t.Fatal(e)
	}
	if one != two {
		t.Fatal("candidate content hashing followed a symlink outside the target")
	}
}
func TestWorkflowContextAndEvidenceDriftAreRejected(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestStart(t, f)
	file := writeAny(t, f.R.WorkspaceRoot, "accept.json", f.Acceptance)
	d := f.D
	d.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	if _, e := WorkflowAccept(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); e == nil {
		t.Fatal("wrong callback cwd accepted")
	}
	if _, e := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context+"-other", file); e == nil {
		t.Fatal("wrong context accepted")
	}
	o = workflowTestAccept(t, f, o)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	review := workflowTestReview(t, f, o, "accepted")
	var artifacts []map[string]any
	json.Unmarshal(f.Report.Artifacts, &artifacts)
	p := artifacts[0]["locator"].(string)
	os.WriteFile(p, []byte("changed"), 0600)
	if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, review); e == nil {
		t.Fatal("artifact drift accepted")
	}
	shown, e := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if len(shown.Reasons) == 0 || shown.Round.State != "unknown" {
		t.Fatal("show hid artifact drift")
	}
}
