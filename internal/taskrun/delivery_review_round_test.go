package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// The synthetic Herdr gives each reviewer tab its own pane and native session.
// It is a stand-in: no real Claude or Codex reviewer is started by these tests.
const reviewRoundHerdrStandin = `
# Each pane has its own shell. pane_shells binds panes to shells explicitly;
# without it the caller's HERDR_PANE_ID is fresh and holds the caller shell.
if cmd in (["pane", "list"], ["pane", "process-info"]):
    shells = model.get("pane_shells")
    if shells is None:
        caller = os.environ.get("HERDR_PANE_ID", "")
        shells = {caller: model["caller_shell"]} if caller else {}
    shells = {"w-fixture:p1": 99999, **shells}
    if cmd == ["pane", "list"]:
        print(json.dumps({"id":"fixture", "result":{"panes":[{"workspace_id":"w-fixture", "pane_id":p, "focused":p == model.get("focused_pane")} for p in sorted(shells)]}}))
        raise SystemExit(0)
    pane = args[args.index("--pane") + 1]
    if pane not in shells:
        print(json.dumps({"id":"fixture", "error":{"code":"pane_not_found", "message":"pane not found"}}))
        raise SystemExit(1)
    print(json.dumps({"id":"fixture", "result":{"process_info":{"pane_id":pane, "shell_pid":shells[pane]}}}))
    raise SystemExit(0)
if cmd == ["tab", "create"] and (args[args.index("--label") + 1].startswith("review ") or args[args.index("--label") + 1].split(" ")[1:2] == ["review"]):
    n = model.get("review_tabs", 0) + 1
    model["review_tabs"] = n
    pane = "w-fixture:r%d" % n
    path.write_text(json.dumps(model))
    print(json.dumps({"id":"fixture", "result":{"root_pane":{"workspace_id":"w-fixture", "tab_id":"tab-review-%d" % n, "pane_id":pane, "terminal_id":"terminal-review-%d" % n}}}))
    raise SystemExit(0)
review_pane = args[args.index("--pane") + 1] if cmd == ["agent", "start"] and "--pane" in args else (args[2] if len(args) > 2 and cmd in (["agent", "get"], ["agent", "prompt"]) else (json.loads(Path(args[4]).read_text())["target"] if cmd == ["agent", "startup"] else ""))
if review_pane.startswith("w-fixture:r"):
    agents = model.setdefault("review_agents", {})
    if cmd == ["agent", "start"]:
        agents[review_pane] = args[2]
    if cmd == ["agent", "prompt"]:
        model.setdefault("review_prompts", {})[review_pane] = args[3]
        if model.get("review_prompt_lost"):
            path.write_text(json.dumps(model))
            raise SystemExit(7)
    session = model["agent_session_id"] if model.get("review_session_is_owner") else "reviewer-session-" + model.get("review_session_prefix", "") + review_pane.split(":r")[1]
    agent = {"workspace_id":"w-fixture", "tab_id":"tab-review-" + review_pane.split(":r")[1], "pane_id":review_pane, "terminal_id":"terminal-review-" + review_pane.split(":r")[1],
             "name":agents.get(review_pane, "unstarted"), "agent":model.get("provider", "codex"), "agent_status":"idle",
             "agent_session":{"agent":model.get("provider", "codex"), "kind":"id", "value":session}}
    if model.get("review_sessionless") and not model.get("review_bootstrapped"):
        agent["agent_session"] = None
        agent["interactive_ready"] = True
    if model.get("al33_bootstrap"):
        agent.update(name=model["al33_name"],agent_session=None,interactive_ready=True)
    elif model.get("al33_live"):
        agent.update(name=model["al33_name"],agent_status="working",agent_session={"agent":"codex","kind":"id","value":"synthetic-live-session"})
    startup_fixture(agent)
    path.write_text(json.dumps(model))
    print(json.dumps({"id":"fixture", "result":{"agent":agent}}))
    raise SystemExit(0)
`

func reviewRoundPatchHerdr(t *testing.T, f *deliveryFixture) {
	t.Helper()
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin+reviewRoundHerdrStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "review-round-request.json", f.R)
}

func reviewRoundFixture(t *testing.T, provider ...string) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}
	var f deliveryFixture
	var o WorkflowRun
	if len(provider) > 0 && provider[0] == "claude" {
		f = receiptFixture(t, "claude", deliveryFixtureOptions{Agreement: &agreement, AllowLocalInstall: true})
		reviewRoundPatchHerdr(t, &f)
		started := workflowTestStart(t, f.workflowFixture)
		var err error
		if o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, started.RunID, started.Paths.Context, DeliveryBeginInput{Provider: "claude", NativeSession: started.Transport.AgentSessionID}, true); err != nil {
			t.Fatal(err)
		}
	} else {
		f = newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement, AllowLocalInstall: true, AdoptExisting: len(provider) > 0 && provider[0] == "adopted"})
		reviewRoundPatchHerdr(t, &f)
		o = deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	}
	f.D.HerdrTimeout = 3 * time.Second
	// A shell of the caller's own process ancestry, as a Herdr pane holds it.
	workflowTestModel(t, f.workflowFixture, map[string]any{"caller_shell": os.Getppid()})
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	return f, o, binary
}

// reviewRoundIn runs fn as a process inside the given Herdr pane.
func reviewRoundIn(t *testing.T, pane string) {
	t.Helper()
	t.Setenv("HERDR_PANE_ID", pane)
}

func reviewRoundArgv(t *testing.T, f deliveryFixture, command string) [][]string {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if err != nil {
		t.Fatal(err)
	}
	out := [][]string{}
	for _, line := range bytes.Split(b, []byte("\n")) {
		var x struct {
			Argv []string `json:"argv"`
		}
		if len(line) == 0 {
			continue
		}
		if err = json.Unmarshal(line, &x); err != nil {
			t.Fatal(err)
		}
		x.Argv = workflowInputArgv(t, x.Argv)
		if strings.Join(x.Argv[:2], " ") == command {
			out = append(out, x.Argv)
		}
	}
	return out
}

// reviewRoundCandidate commits a candidate whose acceptance script stays fixed:
// the oracle is never rewritten to follow a later candidate.
func reviewRoundCandidate(t *testing.T, f deliveryFixture, text string) (string, string) {
	t.Helper()
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte(text+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, target, "add", "README.md")
	runGit(t, target, "commit", "-m", "candidate "+text)
	script := "#!/bin/sh\nset -eu\ngrep -q fixed README.md\ntest \"$(\"$PLY_CANDIDATE_BINARY\")\" = 'fixture candidate'\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	return gitOutput(t, target, "rev-parse", "HEAD"), gitOutput(t, target, "rev-parse", "HEAD^{tree}")
}

func reviewRoundAsCoordinator(f deliveryFixture) Dependencies {
	d := f.D
	d.CWD = func() (string, error) { return f.Parent, nil }
	d.Executable = func() (string, error) { return "/installed/ply", nil }
	return d
}

// The reviewer writes an ordinary file: ordinary mode, capitalised words and
// unsorted findings are adapted by Ply rather than rejected.
func reviewRoundReturnFile(t *testing.T, f deliveryFixture, name, body string) string {
	t.Helper()
	path := filepath.Join(f.R.WorkspaceRoot, name)
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func reviewRoundLatest(t *testing.T, o WorkflowRun) DeliveryReviewRoundStatus {
	t.Helper()
	if o.DeliveryStatus == nil || o.DeliveryStatus.Review == nil || len(o.DeliveryStatus.Review.Rounds) == 0 {
		t.Fatalf("native status shows no review: %+v", o.DeliveryStatus)
	}
	r := o.DeliveryStatus.Review.Rounds
	return r[len(r)-1]
}

func reviewRoundQueue(t *testing.T, f deliveryFixture) string {
	t.Helper()
	a := f.R.Delivery.Agreement
	q, err := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: a.ProjectID, RepoID: a.RepoID, EpicID: a.EpicID}, false)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Canonical(q)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// FS01 AC01-AC05/AC07: one existing Task goes through a Ply-started
// independent reviewer, a negative review, correction, a new exact-candidate
// review and agreed delivery without a review Task, queue move, hand-written
// binding or oracle edit. Herdr and both sessions are synthetic stand-ins.
func TestFS01ReviewRoundsFromExistingTaskToDelivery(t *testing.T) {
	f, o, binary := reviewRoundFixture(t)
	first, firstTree := reviewRoundCandidate(t, f, "first attempt")
	coordinator := reviewRoundAsCoordinator(f)
	queueBefore := reviewRoundQueue(t, f)
	stateFile := filepath.Join(o.Paths.RunRoot, "state.json")
	stateBefore, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{}); err == nil {
		t.Fatal("review started without an actual coordinator claim")
	}
	reviewRoundIn(t, o.Transport.PaneID)
	if _, err = WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "owner"}); err == nil || !strings.Contains(err.Error(), "own review") {
		t.Fatalf("the owner pane coordinated its own review: %v", err)
	}
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatalf("one native operation did not start review of the existing Task: %v", err)
	}
	round := reviewRoundLatest(t, requested)
	actualBrief, e := os.ReadFile(round.Brief.Locator)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(actualBrief), "Review/diff base: "+f.Prepared.Preparation.Preparation.Plan.ParentOID) || strings.Contains(string(actualBrief), "Adopted start:") {
		t.Fatal("ordinary create reviewer scope changed")
	}
	if round.Round != 1 || round.CandidateOID != first || round.CandidateTree != firstTree || round.Coordinator != "fixture root coordinator" || round.State != "requested" || !round.Current || round.ReviewerStart != "started" || round.ReviewerPane != "w-fixture:r1" {
		t.Fatalf("review round does not bind the exact candidate, coordinator and started reviewer: %+v", round)
	}
	if next := requested.DeliveryStatus.Review.Next; next.Actor != "reviewer" || !strings.Contains(next.Message, "w-fixture:r1") {
		t.Fatalf("next responsible action is not the started reviewer: %+v", next)
	}
	tabs, starts, prompts := reviewRoundArgv(t, f, "tab create"), reviewRoundArgv(t, f, "agent start"), reviewRoundArgv(t, f, "agent prompt")
	review := tabs[len(tabs)-1]
	if len(tabs) != 2 || !strings.Contains(strings.Join(review, " "), "--cwd "+f.Prepared.Preparation.Preparation.Plan.WorktreePath+" --label review 1 ") || !strings.Contains(strings.Join(review, " "), "--workspace w-fixture") {
		t.Fatalf("reviewer tab is not a visible Herdr tab in the Task worktree: %v", tabs)
	}
	reviewer := starts[len(starts)-1]
	if len(starts) != 2 || reviewer[2] == starts[0][2] || !strings.HasPrefix(reviewer[2], "ply-review-") || strings.Join(reviewer[3:7], " ") != "--kind codex --pane w-fixture:r1" || !strings.Contains(strings.Join(reviewer, " "), "--model "+f.R.Runtime.Model) {
		t.Fatalf("reviewer was not started as a separate agent with the Run's provider: %v", starts)
	}
	if p := prompts[len(prompts)-1]; p[2] != "w-fixture:r1" || !strings.Contains(p[3], round.Brief.Locator) {
		t.Fatalf("reviewer was not sent its brief: %v", prompts)
	}
	if p := prompts[len(prompts)-1][3]; !strings.Contains(p, "Review only;") || strings.Contains(p, "Local installation allowed:") {
		t.Fatalf("reviewer inherited the installation-capable owner's assignment: %s", p)
	}
	brief, err := os.ReadFile(round.Brief.Locator)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{first, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "review-return " + o.RunID, round.Context.Locator, round.Template.Locator, "/installed/ply"} {
		if !strings.Contains(string(brief), want) {
			t.Fatalf("reviewer brief lacks %q:\n%s", want, brief)
		}
	}
	if reviewRoundQueue(t, f) != queueBefore {
		t.Fatal("starting review changed the native queue")
	}
	again, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil || len(again.DeliveryStatus.Review.Rounds) != 1 || len(reviewRoundArgv(t, f, "tab create")) != 2 || len(reviewRoundArgv(t, f, "agent start")) != 2 {
		t.Fatalf("repeating the request for the same candidate started another reviewer: %v", err)
	}

	negative := reviewRoundReturnFile(t, f, "negative.json", `{"reviewer_claim":"fixture independent reviewer","decision":"Changes_Requested","summary":"Complete review; README lacks the fix.","findings":[{"severity":"High","summary":"README does not contain the fix","evidence":"cat README.md"},{"severity":"low","summary":"Commit message is vague"}],"evidence_gaps":["No real provider journey was observed"]}`)
	// Only the started reviewer pane can return, with its observed session.
	reviewRoundIn(t, o.Transport.PaneID)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, negative); err == nil || !strings.Contains(err.Error(), "bound reviewer pane") {
		t.Fatalf("the owner pane returned the independent review: %v", err)
	}
	reviewRoundIn(t, "w-fixture:r1")
	forged := reviewRoundReturnFile(t, f, "forged.json", `{"reviewer_session_id":"some-other-session","decision":"passed","summary":"Forged","findings":[],"evidence_gaps":[]}`)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, forged); err == nil || !strings.Contains(err.Error(), "observed") {
		t.Fatalf("a self-reported reviewer session replaced the observed one: %v", err)
	}
	contradictory := reviewRoundReturnFile(t, f, "contradictory.json", `{"decision":"passed","summary":"Passed with a finding","findings":[{"severity":"low","summary":"Something"}],"evidence_gaps":[]}`)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, contradictory); err == nil {
		t.Fatal("a passed decision with an open finding was accepted")
	}
	returned, err := WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, negative)
	if err != nil {
		t.Fatalf("complete negative review with truthful gaps was not accepted as a completed review: %v", err)
	}
	round = reviewRoundLatest(t, returned)
	if round.State != "completed" || round.Decision != "changes_requested" || len(round.Findings) != 2 || round.Findings[0].ID != "r001-f01" || round.Findings[0].Severity != "high" || len(round.EvidenceGaps) != 1 || round.ReviewerSession != "reviewer-session-1" {
		t.Fatalf("completed negative review is not readable with findings, gaps and the observed session: %+v", round)
	}
	if next := returned.DeliveryStatus.Review.Next; next.Actor != "owner" || !strings.Contains(next.Message, "finding") {
		t.Fatalf("negative review does not return to the developer: %+v", next)
	}
	if replay, err := WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, negative); err != nil || reviewRoundLatest(t, replay).Decision != "changes_requested" {
		t.Fatalf("identical return was not idempotent: %v", err)
	}
	other := reviewRoundReturnFile(t, f, "other.json", `{"decision":"passed","summary":"Changed my mind","findings":[],"evidence_gaps":[]}`)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, other); err == nil {
		t.Fatal("a sealed review return was overwritten")
	}

	if _, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil || !strings.Contains(err.Error(), "review") {
		t.Fatalf("a candidate with open review findings was verified: %v", err)
	}

	second, secondTree := reviewRoundCandidate(t, f, "fixed attempt")
	if r := reviewRoundLatest(t, mustShow(t, f, o.RunID)); r.Current || mustShow(t, f, o.RunID).DeliveryStatus.Review.CurrentCandidate != "not_reviewed" {
		t.Fatalf("old review still appears to cover the corrected candidate: %+v", r)
	}
	reviewRoundIn(t, "w-fixture:root")
	rereview, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round2 := reviewRoundLatest(t, rereview)
	if round2.Round != 2 || round2.CandidateOID != second || round2.CandidateTree != secondTree || round2.Previous == nil || round2.ReviewerPane != "w-fixture:r2" {
		t.Fatalf("rereview is not an explicit new binding with its own reviewer: %+v", round2)
	}
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"Finding corrected","findings":[],"evidence_gaps":[]}`)
	reviewRoundIn(t, "w-fixture:r1")
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round2.Context.Locator, passed); err == nil {
		t.Fatal("the first round's reviewer returned the second round")
	}
	reviewRoundIn(t, "w-fixture:r2")
	if o, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round2.Context.Locator, passed); err != nil {
		t.Fatal(err)
	}
	if o.DeliveryStatus.Review.CurrentCandidate != "review_passed" || o.DeliveryStatus.Review.Next.Actor != "owner" {
		t.Fatalf("passed rereview not visible: %+v", o.DeliveryStatus.Review)
	}
	stateAfter, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(stateAfter) != string(stateBefore) {
		t.Fatal("review rounds rewrote frozen Run state")
	}

	o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary)
	if err != nil {
		t.Fatalf("owner could not verify with the native passed review: %v", err)
	}
	if len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != second || o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" {
		t.Fatalf("native review did not qualify the exact corrected candidate: %+v", o.DeliveryStatus)
	}
	qualification, err := os.ReadFile(filepath.Join(filepath.Dir(round2.Context.Locator), "candidate-review.json"))
	if err != nil || !strings.Contains(string(qualification), "r001-f01") || !strings.Contains(string(qualification), "reviewer-session-2") {
		t.Fatalf("qualification review does not keep earlier findings as fixes: %s %v", qualification, err)
	}
	if _, err = os.Stat(filepath.Join(filepath.Dir(round.Context.Locator), "return.json")); err != nil {
		t.Fatal("first negative review was not preserved")
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	// Integration finishes while the modeled owner pane is still present.
	// Reuse the ordinary cleanup fixture contract, then explicitly observe its
	// exit before the coordinator follows the same delivered integration.
	cleanupOwnerHandoff(t, f, o, err)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != second {
		t.Fatal("agreed Epic delivery did not integrate the reviewed candidate")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	cleanupOwnerLeaves(t, f)
	reviewRoundIn(t, "w-fixture:root")
	// The owner/provider are modeled, so use the existing process adapter as
	// well. This table contains only the observer, with no Task/session writer;
	// it makes no claim about real OS process inspection.
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}, nil
	}
	assertCleanupUnknownInspectionPreservesDeliveredIntegration(t, f, o, reflog, func() (WorkflowRun, error) {
		return WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 5*time.Second)
	})
	o, err = WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 30*time.Second)
	if err != nil {
		t.Fatalf("coordinator did not finish source cleanup: %v", err)
	}
	assertCleanupRemoved(t, f, o, reflog)
}

func TestFS01ChangedCandidateCannotInheritPass(t *testing.T) {
	f, o, binary := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"Looks right","findings":[],"evidence_gaps":[]}`)
	// The owner moves on while the review is still open. The current source
	// is no longer the round's candidate, so Ply cannot confirm a pass for it
	// (RV01 r001-f02) and nothing is reset; a truthful late negative judgment
	// is still sealed and shows the drift.
	later, _ := reviewRoundCandidate(t, f, "fixed attempt, later edit")
	reviewRoundIn(t, round.ReviewerPane)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, passed); err == nil || !strings.Contains(err.Error(), later) || !strings.Contains(err.Error(), "does not reset") {
		t.Fatalf("a pass was sealed for a candidate other than the tested HEAD: %v", err)
	}
	if head := gitOutput(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "rev-parse", "HEAD"); head != later {
		t.Fatalf("the later candidate was reset: %s", head)
	}
	blocked := reviewRoundReturnFile(t, f, "blocked.json", `{"decision":"blocked","summary":"HEAD moved to a later candidate during review","findings":[],"evidence_gaps":["Current source is a later candidate; which commit the tests exercised is unconfirmed"]}`)
	o, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, blocked)
	if err != nil {
		t.Fatal(err)
	}
	if r := reviewRoundLatest(t, o); r.Decision != "blocked" || r.Current || o.DeliveryStatus.Review.CurrentCandidate != "not_reviewed" || o.DeliveryStatus.Review.Next.Actor != "coordinator" {
		t.Fatalf("drift after review is not visible: %+v", o.DeliveryStatus.Review)
	}
	if _, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil || !strings.Contains(err.Error(), "review") {
		t.Fatalf("changed candidate inherited an earlier pass: %v", err)
	}
	if s := mustShow(t, f, o.RunID); s.Delivery.Attempt != nil {
		t.Fatal("refused verification reserved an attempt")
	}
}

// A Claude owner's reviewer uses the same provider, model and native
// permission mode; a reviewer showing the owner's session is not bound, and an
// unknown start is a visible stop that is never replayed.
func TestFS01ReviewerStartUsesOwnerRuntimeAndStopsOnUnknownStart(t *testing.T) {
	f, o, _ := reviewRoundFixture(t, "claude")
	reviewRoundCandidate(t, f, "first attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_session_is_owner": true})
	if _, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}); err == nil || !strings.Contains(err.Error(), "distinct session") {
		t.Fatalf("a reviewer reporting the owner's session was bound: %v", err)
	}
	starts := reviewRoundArgv(t, f, "agent start")
	reviewer := strings.Join(starts[len(starts)-1], " ")
	if !strings.Contains(reviewer, "--kind claude --pane w-fixture:r1") || !strings.Contains(reviewer, "--permission-mode "+f.R.Runtime.PermissionBinding.ProfileID) || !strings.Contains(reviewer, "--model "+f.R.Runtime.Model) {
		t.Fatalf("Claude reviewer does not use the owner's provider, model and mode: %s", reviewer)
	}
	shown := mustShow(t, f, o.RunID)
	r := reviewRoundLatest(t, shown)
	if r.ReviewerStart != "start_unknown" || shown.DeliveryStatus.Review.Next.Actor != "coordinator" || !strings.Contains(shown.DeliveryStatus.Review.Next.Message, "--new-round") {
		t.Fatalf("unknown reviewer start is not a concrete stop with a way forward: %+v %+v", r, shown.DeliveryStatus.Review.Next)
	}
	if _, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}); err != nil || len(reviewRoundArgv(t, f, "agent start")) != len(starts) {
		t.Fatalf("an unknown reviewer start was replayed: %v", err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_session_is_owner": false})
	fresh, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", NewRound: true})
	if err != nil {
		t.Fatal(err)
	}
	if r := reviewRoundLatest(t, fresh); r.Round != 2 || r.ReviewerStart != "started" || r.ReviewerPane != "w-fixture:r2" || reviewRoundLatest(t, fresh).CandidateOID != reviewRoundLatest(t, shown).CandidateOID {
		t.Fatalf("--new-round did not start a fresh reviewer for the same candidate: %+v", r)
	}
}

func mustShow(t *testing.T, f deliveryFixture, id string) WorkflowRun {
	t.Helper()
	o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, id)
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func reviewRoundCLI(t *testing.T, binary, cwd, pane string, args ...string) (string, error) {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow", "execute"}, args...)...)
	c.Dir = cwd
	c.Env = append(os.Environ(), "HERDR_PANE_ID="+pane)
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	exit := -1
	if c.ProcessState != nil {
		exit = c.ProcessState.ExitCode()
	}
	t.Logf("candidate CLI argv=%q cwd=%q exit=%d\nstdout:\n%s\nstderr:\n%s", c.Args, c.Dir, exit, stdout.String(), stderr.String())
	return stdout.String() + stderr.String(), err
}

// FS01 AC01/AC03/AC07 through the exact candidate CLI: the coordinator and the
// reviewer use only the ordinary commands and native status shows the result.
func TestFS01ReviewRoundsThroughCandidateCLI(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	oid, _ := reviewRoundCandidate(t, f, "first attempt")
	binary := deliveryCLIBinary(t)
	out, err := reviewRoundCLI(t, binary, f.Parent, "w-fixture:root", "review", o.RunID, "--coordinator", "fixture root coordinator", "--format", "json")
	if err != nil {
		t.Fatalf("execute review: %v\n%s", err, out)
	}
	var run WorkflowRun
	if err = json.Unmarshal([]byte(out), &run); err != nil {
		t.Fatalf("execute review output: %v\n%s", err, out)
	}
	round := reviewRoundLatest(t, run)
	if round.CandidateOID != oid || round.State != "requested" || round.ReviewerStart != "started" {
		t.Fatalf("CLI review round differs: %+v", round)
	}
	file := reviewRoundReturnFile(t, f, "cli-negative.json", `{"reviewer_claim":"fixture independent reviewer","decision":"changes_requested","summary":"Complete review with one finding","findings":[{"severity":"medium","summary":"README lacks the fix"}],"evidence_gaps":["No real provider journey"]}`)
	if out, err = reviewRoundCLI(t, binary, f.Parent, round.ReviewerPane, "review-return", o.RunID, "--context", round.Context.Locator, "--file", file); err != nil {
		t.Fatalf("execute review-return: %v\n%s", err, out)
	}
	if out, err = reviewRoundCLI(t, binary, f.Parent, "w-fixture:root", "show", o.RunID); err != nil {
		t.Fatalf("execute show: %v\n%s", err, out)
	}
	for _, want := range []string{"Review round 1 for " + oid[:12] + " (current): completed — changes_requested", "Finding r001-f01 [medium]: README lacks the fix", "Evidence gap: No real provider journey", "Review next: owner"} {
		if !strings.Contains(out, want) {
			t.Fatalf("native status lacks %q:\n%s", want, out)
		}
	}
}

// Review r001-f01: after a lost review prompt the coordinator explicitly
// replaces the round for the same candidate. The replacement completes
// without waking the old reviewer; a late return of the superseded round is
// refused, so it can neither hide a finding nor become a "fix". Sealed
// findings of a candidate keep governing it and cannot be replaced.
func TestFS01SupersededRoundCannotReturnAndReplacementCompletes(t *testing.T) {
	f, o, binary := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_prompt_lost": true})
	if _, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}); err == nil {
		t.Fatal("a lost review prompt was reported as a started review")
	}
	first := reviewRoundLatest(t, mustShow(t, f, o.RunID))
	if first.ReviewerStart != "prompt_unknown" {
		t.Fatalf("lost prompt is not visible: %+v", first)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_prompt_lost": false})
	fresh, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", NewRound: true})
	if err != nil {
		t.Fatal(err)
	}
	second := reviewRoundLatest(t, fresh)
	if old := fresh.DeliveryStatus.Review.Rounds[0]; old.State != "superseded" || old.ReviewerStart != "prompt_unknown" || second.Supersedes == nil || second.Supersedes.Locator != old.Request.Locator {
		t.Fatalf("replacement is not an explicit superseded binding with preserved history: %+v %+v", old, second.Supersedes)
	}
	prompted := len(reviewRoundArgv(t, f, "agent prompt"))
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"Looks right","findings":[],"evidence_gaps":[]}`)
	reviewRoundIn(t, second.ReviewerPane)
	if o, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, second.Context.Locator, passed); err != nil {
		t.Fatal(err)
	}
	if c := o.DeliveryStatus.Review.CurrentCandidate; c != "review_passed" || o.DeliveryStatus.Review.Next.Actor != "owner" {
		t.Fatalf("the replacement review does not complete while the old reviewer is silent: %s %+v", c, o.DeliveryStatus.Review.Next)
	}
	negative := reviewRoundReturnFile(t, f, "late-negative.json", `{"decision":"changes_requested","summary":"Late review","findings":[{"severity":"high","summary":"Late finding"}],"evidence_gaps":[]}`)
	reviewRoundIn(t, first.ReviewerPane)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, first.Context.Locator, negative); err == nil || !strings.Contains(err.Error(), "superseded") {
		t.Fatalf("a superseded round returned a judgment: %v", err)
	}
	if _, err = os.Stat(filepath.Join(filepath.Dir(first.Context.Locator), "return.json")); err == nil {
		t.Fatal("a superseded round sealed a return")
	}
	if len(reviewRoundArgv(t, f, "agent prompt")) != prompted {
		t.Fatal("the old reviewer was woken")
	}
	reviewRoundIn(t, "w-fixture:root")
	if o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err != nil {
		t.Fatalf("the replacement pass could not qualify its candidate: %v", err)
	}
	qualification, err := os.ReadFile(filepath.Join(filepath.Dir(second.Context.Locator), "candidate-review.json"))
	if err != nil || !strings.Contains(string(qualification), `"fixes":[]`) {
		t.Fatalf("qualification invented fixes: %s %v", qualification, err)
	}
}

// A sealed negative judgment of a candidate cannot be replaced by a fresh
// round for that same candidate; only a corrected candidate gets a new round.
func TestFS01SealedFindingsCannotBeReplacedForTheSameCandidate(t *testing.T) {
	f, o, binary := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "first attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	negative := reviewRoundReturnFile(t, f, "negative.json", `{"decision":"changes_requested","summary":"Complete review","findings":[{"severity":"high","summary":"Real finding"}],"evidence_gaps":[]}`)
	reviewRoundIn(t, round.ReviewerPane)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, negative); err != nil {
		t.Fatal(err)
	}
	reviewRoundIn(t, "w-fixture:root")
	if _, err = WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", NewRound: true}); err == nil {
		t.Fatal("a fresh round replaced sealed findings of the same candidate")
	}
	if _, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil || !strings.Contains(err.Error(), "r001-f01") {
		t.Fatalf("verify qualified a candidate with its own open finding: %v", err)
	}
	if s := mustShow(t, f, o.RunID); s.Delivery.Attempt != nil || s.DeliveryStatus.Review.CurrentCandidate != "findings_open" {
		t.Fatal("sealed findings no longer govern the candidate")
	}
}

// A round requested by 8066ecb has no supersedes field. Later candidates read
// that preserved request unchanged and can start the next round after it.
func TestFS01ReadsRoundRequestWithoutSupersedesFrom8066ecb(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "first attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	var old map[string]any
	raw, err := os.ReadFile(round.Request.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &old); err != nil {
		t.Fatal(err)
	}
	delete(old, "supersedes")
	legacy, err := Canonical(old)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(round.Request.Locator, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.DeliveryStatus.Review == nil || len(shown.DeliveryStatus.Review.Rounds) != 1 {
		t.Fatalf("8066ecb round request is unreadable: %v %+v", err, shown.DeliveryStatus)
	}
	reviewRoundCandidate(t, f, "fixed attempt")
	next, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatalf("next round cannot start after a 8066ecb round: %v", err)
	}
	if r := reviewRoundLatest(t, next); r.Round != 2 || r.Supersedes != nil {
		t.Fatalf("next round after a legacy round differs: %+v", r)
	}
	after, err := os.ReadFile(round.Request.Locator)
	if err != nil || string(after) != string(legacy) {
		t.Fatal("the preserved 8066ecb request was rewritten")
	}
}
