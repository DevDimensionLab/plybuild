package taskrun

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Inspect the actual generated inputs, not a separately maintained prompt.
func al15ReviewInstructions(t *testing.T, f deliveryFixture, round DeliveryReviewRoundStatus) {
	t.Helper()
	brief, err := os.ReadFile(round.Brief.Locator)
	if err != nil {
		t.Fatal(err)
	}
	prompts := reviewRoundArgv(t, f, "agent prompt")
	prompt := prompts[len(prompts)-1][3]
	for name, text := range map[string]string{"brief": string(brief), "launch prompt": prompt} {
		for _, want := range []string{"relevant ordinary tests", "source", "cache", "temporary", "runtime", "workflow execute verify", "workflow execute integrate", "owner-only", "-mod=readonly", "assigned product Task/Run", "isolated fixtures", "synthetic owner callbacks", "do not grant", "actual refusal", "unknown prior effect"} {
			if !strings.Contains(text, want) {
				t.Errorf("%s lacks %q: %s", name, want, text)
			}
		}
		for _, ambiguous := range []string{"commit, verify or integrate", "must not be invoked by the reviewer"} {
			if strings.Contains(text, ambiguous) {
				t.Errorf("%s retains ambiguous prohibition: %s", name, text)
			}
		}
	}
	for _, want := range []string{round.CandidateOID, round.CandidateTree, "Never reset", "blocked", "evidence"} {
		if !strings.Contains(string(brief), want) {
			t.Errorf("brief lost candidate/evidence boundary %q", want)
		}
	}
}

func TestAL15GeneratedReviewTestMandate(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed mandate")
	reviewRoundIn(t, "w-fixture:root")
	run, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	al15ReviewInstructions(t, f, reviewRoundLatest(t, run))
}

// The executable and native Git/review storage are real. Herdr/provider identities
// and judgments are modeled; this does not establish actual runtime approval.
func TestAL15ReviewTestMandateThroughCandidateCLI(t *testing.T) {
	if dir := os.Getenv("PLY_AL15_EVIDENCE_DIR"); dir != "" {
		t.Setenv("PLY_AL13_EVIDENCE_DIR", dir)
	}
	binary := deliveryCLIBinary(t)
	f, o, _ := reviewRoundFixture(t)
	oid, tree := reviewRoundCandidate(t, f, "fixed CLI mandate")
	run, detail, exit := al13CLI(t, binary, f, "w-fixture:root", "request", "review", o.RunID, "--coordinator", "modeled coordinator")
	if exit != 0 {
		t.Fatal(detail)
	}
	round := reviewRoundLatest(t, run)
	if round.CandidateOID != oid || round.CandidateTree != tree || round.ReviewerStart != "started" {
		t.Fatalf("candidate/start drift: %+v", round)
	}
	al15ReviewInstructions(t, f, round)
	briefBefore, err := os.ReadFile(round.Brief.Locator)
	if err != nil {
		t.Fatal(err)
	}
	// A harmless ordinary command reads the candidate without owner callbacks.
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if got := al15CandidateStatus(t, target); got != "" {
		t.Fatal(got)
	}
	file := reviewRoundReturnFile(t, f, "review.json", `{"decision":"passed","summary":"Modeled review completed; candidate inspected","findings":[],"evidence_gaps":["Provider and Herdr observations are modeled; no live runtime approval claimed."]}`)
	returned, detail, exit := al13CLI(t, binary, f, round.ReviewerPane, "return", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
	if exit != 0 {
		t.Fatal(detail)
	}
	sealed := reviewRoundLatest(t, returned)
	if sealed.Decision != "passed" || sealed.CandidateOID != oid || sealed.CandidateTree != tree || sealed.ReviewerSession == o.Transport.AgentSessionID {
		t.Fatalf("return binding: %+v", sealed)
	}
	// Safe native counterexamples in the isolated fixture: reviewer context cannot
	// authorize owner verification, and an unqualified review cannot integrate.
	// No verifier process, candidate registration or delivery effect may start.
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	stateBefore, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	callbackFixture := f
	callbackFixture.Parent = target
	_, detail, exit = al13CLI(t, binary, callbackFixture, round.ReviewerPane, "reviewer-verify-refused", "verify", o.RunID, "--context", round.Context.Locator, "--candidate-binary", binary)
	if exit != 4 || !strings.Contains(detail, "exact current private context") {
		t.Fatalf("reviewer context authorized owner verification: exit %d %s", exit, detail)
	}
	_, detail, exit = al13CLI(t, binary, callbackFixture, round.ReviewerPane, "reviewer-integrate-refused", "integrate", o.RunID, "--context", round.Context.Locator)
	if exit != 4 || !strings.Contains(detail, "qualified candidate") {
		t.Fatalf("unqualified reviewer integrated: exit %d %s", exit, detail)
	}
	stateAfter, err := os.ReadFile(statePath)
	if err != nil || string(stateAfter) != string(stateBefore) {
		t.Fatal("refused owner callbacks changed Run state", err)
	}
	// A repeat reads the sealed result; it cannot replace a judgment or brief.
	_, detail, exit = al13CLI(t, binary, f, round.ReviewerPane, "repeat", "review-return", o.RunID, "--context", round.Context.Locator, "--file", file)
	if exit != 0 {
		t.Fatal(detail)
	}
	briefAfter, err := os.ReadFile(round.Brief.Locator)
	if err != nil || string(briefAfter) != string(briefBefore) {
		t.Fatal("sealed brief changed", err)
	}
	if gitOutput(t, target, "rev-parse", "HEAD") != oid || al15CandidateStatus(t, target) != "" {
		t.Fatal("candidate changed")
	}
	// Keep raw generated instructions with the process evidence if requested.
	if dir := os.Getenv("PLY_AL15_EVIDENCE_DIR"); dir != "" {
		if err := os.WriteFile(filepath.Join(dir, "generated-review-brief.md"), briefAfter, 0600); err != nil {
			t.Fatal(err)
		}
		prompts := reviewRoundArgv(t, f, "agent prompt")
		if err := os.WriteFile(filepath.Join(dir, "generated-review-launch.txt"), []byte(prompts[len(prompts)-1][3]), 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func al15CandidateStatus(t *testing.T, target string) string {
	t.Helper()
	c := exec.Command("git", "-C", target, "status", "--porcelain")
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("candidate status: %v %s", err, out)
	}
	return string(out)
}
