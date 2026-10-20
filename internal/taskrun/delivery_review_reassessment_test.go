package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the public commands with real native records and a synthetic Herdr.
// No commit or source edit is needed when the missing review evidence arrives.
func TestRV02EvidenceBlockedReassessmentThroughCandidateCLI(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, acceptanceBinary := reviewRoundFixture(t)
	oid, tree := reviewRoundCandidate(t, f, "fixed candidate")
	call := func(pane string, args ...string) WorkflowRun {
		t.Helper()
		out, err := reviewRoundCLI(t, binary, f.Parent, pane, append(args, "--format", "json")...)
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
		var run WorkflowRun
		if err = json.Unmarshal([]byte(out), &run); err != nil {
			t.Fatalf("decode %v: %v\n%s", args, err, out)
		}
		return run
	}
	request := []string{"review", o.RunID, "--coordinator", "fixture root coordinator"}
	first := reviewRoundLatest(t, call("w-fixture:root", request...))
	blocked := reviewRoundReturnFile(t, f, "blocked.json", `{"decision":"blocked","summary":"Cannot judge without the runtime transcript","findings":[],"evidence_gaps":["runtime-transcript.txt is missing"]}`)
	call(first.ReviewerPane, "review-return", o.RunID, "--context", first.Context.Locator, "--file", blocked)
	before := map[string][]byte{}
	for _, path := range []string{first.Request.Locator, first.Brief.Locator, first.Context.Locator, first.Template.Locator, filepath.Join(filepath.Dir(first.Context.Locator), "return.json"), filepath.Join(filepath.Dir(first.Context.Locator), "return-input.json"), filepath.Join(filepath.Dir(first.Context.Locator), "agent.json"), filepath.Join(filepath.Dir(first.Context.Locator), "launch.json"), filepath.Join(o.Paths.RunRoot, "state.json")} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		before[path] = raw
	}
	starts := len(reviewRoundArgv(t, f, "agent start"))
	shown := call("w-fixture:root", request...)
	if len(shown.DeliveryStatus.Review.Rounds) != 1 || len(reviewRoundArgv(t, f, "agent start")) != starts {
		t.Fatal("ordinary repeat restarted the sealed evidence-blocked review")
	}
	// Supply actual additional fixture evidence outside the unchanged source.
	if err := os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "runtime-transcript.txt"), []byte("Observed exact candidate runtime check passed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	secondRun := call("w-fixture:root", append(request, "--new-round")...)
	if !strings.Contains(shown.DeliveryStatus.Review.Next.Message, "--new-round") {
		t.Fatalf("blocked status has no concrete reassessment action: %+v", shown.DeliveryStatus.Review)
	}
	second := reviewRoundLatest(t, secondRun)
	if second.Round != 2 || second.CandidateOID != oid || second.CandidateTree != tree || second.Supersedes == nil || *second.Supersedes != first.Request || second.ReviewerPane == first.ReviewerPane || secondRun.DeliveryStatus.Review.CurrentCandidate != "review_requested" {
		t.Fatalf("fresh independent same-candidate review not current: %+v", secondRun.DeliveryStatus.Review)
	}
	old := secondRun.DeliveryStatus.Review.Rounds[0]
	if old.State != "superseded" || old.Decision != "blocked" || len(old.EvidenceGaps) != 1 || old.ReviewerSession != "reviewer-session-1" {
		t.Fatalf("original blocked judgment is not preserved as history: %+v", old)
	}
	// An active replacement is awaited; it does not inherit any judgment.
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", acceptanceBinary); err == nil || !strings.Contains(err.Error(), "has not returned") {
		t.Fatalf("open replacement qualified: %v", err)
	}
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"New runtime transcript resolves the evidence gap","findings":[],"evidence_gaps":[]}`)
	current := call(second.ReviewerPane, "review-return", o.RunID, "--context", second.Context.Locator, "--file", passed)
	if current.DeliveryStatus.Review.CurrentCandidate != "review_passed" || reviewRoundLatest(t, current).ReviewerSession != "reviewer-session-2" {
		t.Fatalf("new reviewer pass not current: %+v", current.DeliveryStatus.Review)
	}
	// The old identical callback is idempotent; a changed late return is refused.
	current = call(first.ReviewerPane, "review-return", o.RunID, "--context", first.Context.Locator, "--file", blocked)
	if current.DeliveryStatus.Review.CurrentCandidate != "review_passed" {
		t.Fatal("old callback overrode the current pass")
	}
	if out, err := reviewRoundCLI(t, binary, f.Parent, first.ReviewerPane, "review-return", o.RunID, "--context", first.Context.Locator, "--file", passed); err == nil || !strings.Contains(out, "sealed return") {
		t.Fatalf("old reviewer changed its sealed judgment: %v %s", err, out)
	}
	for path, raw := range before {
		after, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(raw, after) {
			t.Fatalf("reassessment rewrote preserved bytes at %s: %v", path, err)
		}
	}
	qualified, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", acceptanceBinary)
	if err != nil || qualified.DeliveryStatus.Acceptance == nil || qualified.DeliveryStatus.Acceptance.Outcome != "pass" {
		t.Fatalf("unchanged candidate cannot qualify using current review: %v %+v", err, qualified.DeliveryStatus)
	}
	qualification, err := os.ReadFile(filepath.Join(filepath.Dir(second.Context.Locator), "candidate-review.json"))
	if err != nil || !bytes.Contains(qualification, []byte(`"fixes":[]`)) {
		t.Fatalf("evidence gap was invented as a product fix: %s %v", qualification, err)
	}
	if got := gitOutput(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "rev-parse", "HEAD"); got != oid {
		t.Fatal("reassessment changed the candidate")
	}
}

func TestRV02OnlyEvidenceBlockedReviewCanBeReassessed(t *testing.T) {
	for _, body := range []string{
		`{"decision":"blocked","summary":"Product issue and missing evidence","findings":[{"severity":"high","summary":"Known product defect"}],"evidence_gaps":["Missing transcript"]}`,
		`{"decision":"changes_requested","summary":"Product issue","findings":[{"severity":"high","summary":"Known product defect"}],"evidence_gaps":["Missing transcript"]}`,
		`{"decision":"passed","summary":"Complete","findings":[],"evidence_gaps":[]}`,
	} {
		t.Run(body, func(t *testing.T) {
			f, o, binary := reviewRoundFixture(t)
			reviewRoundCandidate(t, f, "fixed candidate")
			d := reviewRoundAsCoordinator(f)
			reviewRoundIn(t, "w-fixture:root")
			in := DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}
			requested, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, in)
			if err != nil {
				t.Fatal(err)
			}
			round := reviewRoundLatest(t, requested)
			reviewRoundIn(t, round.ReviewerPane)
			file := reviewRoundReturnFile(t, f, "judgment.json", body)
			if _, err = WorkflowDeliveryReviewReturn(d, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, file); err != nil {
				t.Fatal(err)
			}
			reviewRoundIn(t, "w-fixture:root")
			in.NewRound = true
			if _, err = WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, in); err == nil {
				t.Fatal("same-candidate reassessment bypassed a sealed judgment outside the evidence-only path")
			}
			shown := mustShow(t, f, o.RunID)
			if len(shown.DeliveryStatus.Review.Rounds) != 1 {
				t.Fatal("refused reassessment wrote a round")
			}
			if strings.Contains(body, "Known product defect") {
				if _, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil || !strings.Contains(err.Error(), "r001-f01") {
					t.Fatalf("known finding stopped governing verification: %v", err)
				}
			}
		})
	}
}

func TestRV02ReassessmentKeepsNewProductFindingsOpen(t *testing.T) {
	f, o, binary := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed candidate")
	d := reviewRoundAsCoordinator(f)
	in := DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}
	for n, body := range []string{
		`{"decision":"blocked","summary":"Missing evidence","findings":[],"evidence_gaps":["Runtime transcript missing"]}`,
		`{"decision":"blocked","summary":"Still missing evidence","findings":[],"evidence_gaps":["Runtime output incomplete"]}`,
		`{"decision":"blocked","summary":"New evidence exposes a defect","findings":[{"severity":"high","summary":"Runtime accepts invalid input"}],"evidence_gaps":["Additional output missing"]}`,
	} {
		reviewRoundIn(t, "w-fixture:root")
		in.NewRound = n > 0
		requested, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, in)
		if err != nil {
			t.Fatal(err)
		}
		round := reviewRoundLatest(t, requested)
		reviewRoundIn(t, round.ReviewerPane)
		file := reviewRoundReturnFile(t, f, "judgment.json", body)
		if _, err = WorkflowDeliveryReviewReturn(d, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, file); err != nil {
			t.Fatal(err)
		}
	}
	reviewRoundIn(t, "w-fixture:root")
	if _, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("new product finding disappeared into another reassessment")
	}
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil || !strings.Contains(err.Error(), "r003-f01") {
		t.Fatalf("new finding no longer governs verification: %v", err)
	}
	shown := mustShow(t, f, o.RunID)
	if len(shown.DeliveryStatus.Review.Rounds) != 3 || shown.DeliveryStatus.Review.Next.Actor != "owner" || strings.Contains(shown.DeliveryStatus.Review.Next.Message, "--new-round") {
		t.Fatalf("status offers a bypass of the known finding: %+v", shown.DeliveryStatus.Review)
	}
}

// Legacy history can contain multiple judgments for an identical candidate.
// Even a supersedes link must not turn a known product defect into a fix.
func TestRV02KnownNegativeHistoryGovernsLaterEvidenceRounds(t *testing.T) {
	rounds := []deliveryReviewRecord{
		{Superseded: true, Request: deliveryReviewRound{Round: 1, CandidateOID: "oid", CandidateTree: "tree"}, Result: &deliveryReviewResult{Decision: "changes_requested", Findings: []DeliveryReviewFinding{{ID: "r001-f01", Summary: "Known defect"}}}},
		{Request: deliveryReviewRound{Round: 2, CandidateOID: "oid", CandidateTree: "tree"}, Result: &deliveryReviewResult{Decision: "blocked", EvidenceGaps: []string{"Missing transcript"}}},
		{Request: deliveryReviewRound{Round: 3, CandidateOID: "oid", CandidateTree: "tree"}, Result: &deliveryReviewResult{Decision: "passed"}},
	}
	verdict, current := deliveryReviewJudge(rounds, "oid", "tree")
	if verdict != "negative" || current.Request.Round != 1 || current.Result.Findings[0].ID != "r001-f01" {
		t.Fatalf("earlier known finding no longer governs: %s %+v", verdict, current)
	}
}
