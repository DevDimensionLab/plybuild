package taskrun

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
)

func TestParallelRetirementUnknownAbsencePreservesBothMembers(t *testing.T) {
	f, a, b, second := parallelOwnerPair(t)
	proof := writeAny(t, f.R.WorkspaceRoot, "retire-denial.json", map[string]string{"observation": "synthetic A-only refusal"})
	if _, e := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, a.RunID, DeliveryStopInput{EventID: "a-retire-stop", Provider: "codex", NativeSession: a.Transport.AgentSessionID, Stage: "report", Observer: "coordinator", Reason: "simulated refused report", EvidencePath: proof}); e != nil {
		t.Fatal(e)
	}
	before, e := workflowRead(b.R.WorkspaceRoot, second.RunID)
	if e != nil {
		t.Fatal(e)
	}
	coordinator := reviewRoundAsCoordinator(f)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "inventory": []any{}, "retire_inventory_denied": true})
	input := DeliveryRetireInput{Observer: "synthetic coordinator", Reason: "retire A only after known writer absence"}
	if _, e = WorkflowRetireDelivery(coordinator, f.R.WorkspaceRoot, a.RunID, input); e == nil {
		t.Fatal("unknown writer absence retired A")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if !workspace.PreparationIsActive(r, f.R.PreparationID) || !workspace.PreparationIsActive(r, b.R.PreparationID) {
		t.Fatal("unknown retirement removed membership")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"retire_inventory_denied": false})
	if _, e = WorkflowRetireDelivery(coordinator, f.R.WorkspaceRoot, a.RunID, input); e != nil {
		t.Fatal(e)
	}
	r, e = f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if workspace.PreparationIsActive(r, f.R.PreparationID) || !workspace.PreparationIsActive(r, b.R.PreparationID) {
		t.Fatal("retirement crossed member identity")
	}
	after, e := workflowRead(b.R.WorkspaceRoot, second.RunID)
	if e != nil || !equal(before, after) {
		t.Fatal("A retirement changed B Run")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Plan.WorktreePath); e != nil {
		t.Fatal("retirement removed retained source")
	}
}

func parallelOwnerPair(t *testing.T) (deliveryFixture, WorkflowRun, deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, Parallel: true})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o, e := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	b := parallelSecondOwner(t, f)
	p := workflowTestStart(t, b.workflowFixture)
	p, e = WorkflowDeliveryReceipt(b.D, b.R.WorkspaceRoot, p.RunID, p.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: p.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	return f, o, b, p
}

func TestParallelOwnerStopQuestionAndContextIsolation(t *testing.T) {
	f, a, b, second := parallelOwnerPair(t)
	reviewRoundIn(t, a.Transport.PaneID)
	a = releaseProgressReport(t, f, a, "needs_input", "a-question")
	reviewRoundIn(t, second.Transport.PaneID)
	second = releaseProgressReport(t, b, second, "working", "b-progress")
	if a.Delivery.Phase != "needs_input" || second.Delivery.Phase != "working" {
		t.Fatal("reports crossed Run ownership")
	}
	if _, e := WorkflowDeliveryReceipt(b.D, b.R.WorkspaceRoot, second.RunID, a.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: second.Transport.AgentSessionID}, true); e == nil {
		t.Fatal("A context accepted for B")
	}
	proof := writeAny(t, f.R.WorkspaceRoot, "parallel-simulated-denial.json", map[string]string{"observation": "synthetic refusal, not a live provider denial"})
	stopped, e := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, a.RunID, DeliveryStopInput{EventID: "a-refusal", Provider: "codex", NativeSession: a.Transport.AgentSessionID, Stage: "report", Observer: "coordinator", Reason: "simulated A-only refused action", EvidencePath: proof})
	if e != nil {
		t.Fatal(e)
	}
	reviewRoundIn(t, a.Transport.PaneID)
	if _, e = WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, a.RunID, a.Paths.Context, workspace.TaskSourceReconciliationInput{}); e == nil {
		t.Fatal("A known stop did not block dependent reconciliation")
	}
	reviewRoundIn(t, second.Transport.PaneID)
	second = releaseProgressReport(t, b, second, "working", "b-after-a-stop")
	source := b.Prepared.Preparation.Plan.WorktreePath
	if e = os.WriteFile(filepath.Join(source, "independent.txt"), []byte("lawful B source work\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, source, "add", "independent.txt")
	runGit(t, source, "commit", "-m", "B independent source")
	state, e := workflowRead(f.R.WorkspaceRoot, a.RunID)
	if e != nil || state.Result.Delivery.ActiveStop == nil || *state.Result.Delivery.ActiveStop != *stopped.Delivery.ActiveStop {
		t.Fatal("B progress erased A refusal")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if !workspace.PreparationIsActive(r, f.R.PreparationID) || !workspace.PreparationIsActive(r, b.R.PreparationID) {
		t.Fatal("stop/question closed participation")
	}
}

func TestParallelSameOIDNewGenerationNeedsFreshReviewAndVerification(t *testing.T) {
	f, o, _, _ := parallelOwnerPair(t)
	source := f.Prepared.Preparation.Plan.WorktreePath
	runGit(t, source, "commit", "--allow-empty", "-m", "C parent to be advanced")
	parentNext := gitOutput(t, source, "rev-parse", "HEAD")
	runGit(t, source, "commit", "--allow-empty", "-m", "fixed exact source candidate")
	oid := gitOutput(t, source, "rev-parse", "HEAD")
	binary := filepath.Join(f.R.WorkspaceRoot, "modeled-candidate")
	if e := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nexit 0\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	runGit(t, f.Parent, "merge", "--ff-only", parentNext)
	base, e := workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
		t.Fatal(e)
	}
	preview, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	out, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: preview.Confirmation})
	if e != nil || out.State != "committed" {
		t.Fatalf("no-op reconciliation: %+v %v", out, e)
	}
	if gitOutput(t, source, "rev-parse", "HEAD") != oid {
		t.Fatal("already-inclusive source was rewritten")
	}
	shown, e := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if shown.DeliveryStatus.Review != nil && shown.DeliveryStatus.Review.CurrentCandidate == "review_passed" {
		t.Fatal("readback presents old-generation review as current at same OID")
	}
	if _, e = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); e == nil {
		t.Fatal("same OID inherited old-generation review/acceptance")
	}
	o = parallelQualify(t, f, o, binary)
	if len(o.Delivery.Candidates) != 2 || o.Delivery.Candidates[0].OID != o.Delivery.Candidates[1].OID || o.Delivery.Candidates[0].Verification == o.Delivery.Candidates[1].Verification {
		t.Fatal("fresh generation did not execute distinct qualification")
	}
}

func TestParallelReviewerAndOwnerSessionExclusion(t *testing.T) {
	for _, direction := range []string{"peer-reviewer", "owner-after-peer-reviewer"} {
		t.Run(direction, func(t *testing.T) {
			f, a, b, second := parallelOwnerPair(t)
			for _, f := range []deliveryFixture{f, b} {
				raw, e := os.ReadFile(f.Model)
				if e != nil {
					t.Fatal(e)
				}
				var m map[string]any
				if e = json.Unmarshal(raw, &m); e != nil {
					t.Fatal(e)
				}
				m["review_session_prefix"] = "deliberate-collision-"
				writeAny(t, filepath.Dir(f.Model), filepath.Base(f.Model), m)
			}
			reviewRoundIn(t, "w-fixture:root")
			first, e := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, a.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic coordinator"})
			if e != nil {
				t.Fatal(e)
			}
			round := reviewRoundLatest(t, first)
			state, e := workflowRead(f.R.WorkspaceRoot, a.RunID)
			if e != nil {
				t.Fatal(e)
			}
			rounds, e := deliveryReviewRounds(state)
			if e != nil {
				t.Fatal(e)
			}
			started := deliveryReviewStartState(rounds[0].Dir)
			if started.Agent == nil {
				t.Fatal("no simulated reviewer session")
			}
			if direction == "peer-reviewer" {
				// Each model intentionally returns reviewer-session-1: the two
				// visible fake tabs cannot establish distinct native ownership.
				_, e = WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(b), b.R.WorkspaceRoot, second.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic coordinator"})
				if e == nil {
					t.Fatal("peer active reviewer session reused for B")
				}
			} else {
				other, e := workflowRead(b.R.WorkspaceRoot, second.RunID)
				if e != nil {
					t.Fatal(e)
				}
				if e = parallelSessionAvailable(b.D, other, "codex", started.Agent.SessionID, true); e == nil {
					t.Fatal("new owner could bind peer live reviewer session")
				}
			}
			reviewRoundIn(t, round.ReviewerPane)
			if _, e = WorkflowDeliveryReviewReturn(b.D, b.R.WorkspaceRoot, second.RunID, round.Context.Locator, reviewRoundReturnFile(t, f, "foreign-review-control.json", `{"decision":"passed","summary":"synthetic wrong control","findings":[],"evidence_gaps":[]}`)); e == nil {
				t.Fatal("A reviewer context/control returned on B")
			}
		})
	}
}

func TestParallelParentRacesDuringReviewAndBeforeReturn(t *testing.T) {
	f, o, _, _ := parallelOwnerPair(t)
	source := f.Prepared.Preparation.Plan.WorktreePath
	if e := os.WriteFile(filepath.Join(source, "owned.txt"), []byte("fixed\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, source, "add", "owned.txt")
	runGit(t, source, "commit", "-m", "owned candidate")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, e := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic race coordinator"})
	if e != nil {
		t.Fatal(e)
	}
	round := reviewRoundLatest(t, requested)
	move := func(message string) {
		t.Helper()
		runGit(t, f.Parent, "commit", "--allow-empty", "-m", message)
		base, e := workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
		if e != nil {
			t.Fatal(e)
		}
		if _, e = workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
			t.Fatal(e)
		}
	}
	move("C moved during A review")
	reviewRoundIn(t, round.ReviewerPane)
	pass := reviewRoundReturnFile(t, f, "stale-parent-review.json", `{"decision":"passed","summary":"synthetic outdated pass","findings":[],"evidence_gaps":[]}`)
	if _, e = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, pass); e == nil {
		t.Fatal("parent move inherited review pass")
	}
	blocked := reviewRoundReturnFile(t, f, "changed-parent-blocked.json", `{"decision":"blocked","summary":"parent changed during review","findings":[],"evidence_gaps":["current parent generation changed"]}`)
	if _, e = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, blocked); e != nil {
		t.Fatal(e)
	}
	reviewRoundIn(t, o.Transport.PaneID)
	p, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: p.Confirmation}); e != nil {
		t.Fatal(e)
	}
	binary := filepath.Join(f.R.WorkspaceRoot, "modeled-race-candidate")
	if e = os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\ngrep -q fixed owned.txt\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	move("C moved immediately before return")
	before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	if _, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
		t.Fatal("stale effective parent integrated")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
		t.Fatal("stale return changed target")
	}
}

func TestParallelHumanVetoSurvivesNativeReconciliation(t *testing.T) {
	f, o, _, _ := parallelOwnerPair(t)
	source := f.Prepared.Preparation.Plan.WorktreePath
	if e := os.WriteFile(filepath.Join(source, "own.txt"), []byte("fixed owned behavior\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, source, "add", "own.txt")
	runGit(t, source, "commit", "-m", "fixed behavior")
	binary := filepath.Join(f.R.WorkspaceRoot, "simulated-candidate")
	if e := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\ngrep -q fixed own.txt\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	o, e := WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "fail", deliveryTestHuman(t, f, o, "fail"))
	if e != nil {
		t.Fatal(e)
	}
	original := o.Delivery.Candidates[0].HumanQA
	runGit(t, f.Parent, "commit", "--allow-empty", "-m", "C parent movement after human veto")
	base, e := workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
		t.Fatal(e)
	}
	p, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: p.Confirmation}); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	if o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "fail" || o.DeliveryStatus.Acceptance.HumanQARecordID != original.ID {
		t.Fatal("new source/automatic qualification cleared native human veto")
	}
	if _, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
		t.Fatal("native return bypassed preserved human fail")
	}
	passFile := deliveryTestHuman(t, f, o, "pass")
	raw, e := os.ReadFile(passFile)
	if e != nil {
		t.Fatal(e)
	}
	var pass map[string]any
	if e = json.Unmarshal(raw, &pass); e != nil {
		t.Fatal(e)
	}
	pass["completed_at_utc"] = "2026-10-05T00:00:02Z"
	o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", writeAny(t, f.R.WorkspaceRoot, "later-exact-pass.json", pass))
	if e != nil {
		t.Fatal(e)
	}
	if o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" {
		t.Fatal("later actual modeled exact human pass did not resolve veto")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || len(r.HumanQARecords) != 2 {
		t.Fatal("machine qualification invented/erased human QA")
	}
	found := false
	for _, qa := range r.HumanQARecords {
		if qa.ID == original.ID {
			found = equal(qa, *original)
		}
	}
	if !found {
		t.Fatal("old human verdict was rewritten")
	}
}
