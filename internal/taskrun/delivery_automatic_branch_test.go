package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Public preview uses a real isolated registry and Git checkout. No provider,
// process inventory, execution, or integration is needed to register a Goal.
func TestPWADep01AutomaticBranchRegistrationPreview(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, branch := range []string{"main", "master"} {
		t.Run(branch, func(t *testing.T) {
			root, parent, _, oid := prepareServiceWorkspace(t)
			if branch != "main" {
				runGit(t, parent, "branch", "-m", branch)
			}
			w := workspace.SystemDependencies()
			if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Synthetic branch return", ProjectID: "ply", RepoID: "ply", Worktree: parent, Ref: "refs/heads/" + branch, ExpectedOID: oid}); err != nil {
				t.Fatal(err)
			}
			doc := filepath.Join(root, "goal.md")
			if err := os.WriteFile(doc, []byte("# Explicit local return\n"), 0600); err != nil {
				t.Fatal(err)
			}
			before := branchFixtureState(t, root)
			args := []string{"workspace", "task", "goal", "register", "branch-task", "--spec", "branch-v1", "--document", doc, "--title", "Explicit branch", "--objective", "Return the exact candidate", "--criterion", "The selected branch receives the qualified candidate", "--project", "ply", "--repo", "ply", "--epic", "epic", "--delivery", workspace.DeliveryLocalBranch, "--acceptance", "automatic", "--responsible-actor", "synthetic fixture owner", "--acceptance-timeout", "30", "--format", "json"}
			cmd := exec.Command(binary, args...)
			cmd.Dir = parent
			raw, err := cmd.CombinedOutput()
			t.Logf("argv=%q cwd=%s exit_error=%v\n%s", args, parent, err, raw)
			if err != nil {
				t.Fatal(err)
			}
			var out workspace.FirstGoalResult
			if err := json.Unmarshal(raw, &out); err != nil {
				t.Fatal(err)
			}
			a := out.Delivery
			if out.Applied || out.Phase != "preview" || out.Run != nil || a.Mode != workspace.DeliveryLocalBranch || a.TargetRef != "refs/heads/"+branch || a.TargetWorktree != parent || a.SourceRef != "" || !a.AutomaticAcceptance() || a.Acceptance.RequireHumanQA || a.Acceptance.ResponsibleActor != "synthetic fixture owner" || a.Acceptance.TimeoutSeconds != 30 {
				t.Fatalf("preview changed explicit target/policy or started work: %+v", out)
			}
			if !reflect.DeepEqual(before, branchFixtureState(t, root)) || gitOutput(t, parent, "rev-parse", "HEAD") != oid {
				t.Fatal("registration preview wrote native state or Git")
			}
		})
	}
}

func branchFixtureState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(filepath.Join(root, ".ply"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			out[path] = hashFileTest(t, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// All OS/provider identity observations in this journey are injected. Git,
// native review sealing, acceptance execution, result/queue/base and return
// are real disposable fixture operations, not approval of this implementation.
func automaticBranchFixture(t *testing.T, branch string, requireHuman bool, owner string) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalBranch, IntegrationOwner: owner, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic branch owner", RequireHumanQA: requireHuman, TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: a, ParentRef: branch})
	reviewRoundPatchHerdr(t, &f)
	// Never fall through to the real lsof/ps observers, including cleanup.
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: 70001, CWD: f.Parent}, {PID: 70001, ParentPID: 1, CWD: f.Parent}}, nil
	}
	f.D.StartupProcessTree = func(pid int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"caller_shell": 70001})
	started := workflowTestStart(t, f.workflowFixture)
	o, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, started.RunID, started.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: started.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	oid, tree := reviewRoundCandidate(t, f, "fixed explicit branch candidate")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic branch coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	reviewRoundIn(t, round.ReviewerPane)
	returned, err := WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, reviewRoundReturnFile(t, f, "branch-review.json", `{"decision":"passed","summary":"Synthetic exact candidate review completed","findings":[],"evidence_gaps":[]}`))
	if err != nil || reviewRoundLatest(t, returned).ReviewerSession == started.Transport.AgentSessionID {
		t.Fatalf("separate native review failed: %v", err)
	}
	reviewRoundIn(t, started.Transport.PaneID)
	// The bound owner's context is required; a review context cannot qualify.
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, "", binary); err == nil {
		t.Fatal("reviewer context authorized owner verification")
	}
	o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary)
	if err != nil || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != oid || o.Delivery.Candidates[0].Tree != tree || !o.DeliveryStatus.Verification.InputsMatch || !o.DeliveryStatus.QualifiedCandidate.Current {
		t.Fatalf("exact branch qualification failed: %+v %v", o.DeliveryStatus, err)
	}
	guide, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "delivery", "callback-guide.md"))
	if err != nil || !strings.Contains(string(guide), "Mode: "+workspace.DeliveryLocalBranch) || !requireHuman && owner == "" && !strings.Contains(string(guide), "## Automatic local branch delivery") {
		t.Fatalf("branch guide lost selected mode: %v", err)
	}
	return f, o, binary
}

func TestPWADep01AutomaticBranchNativeJourney(t *testing.T) {
	for _, branch := range []string{"main", "master"} {
		t.Run(branch, func(t *testing.T) {
			f, o, _ := automaticBranchFixture(t, branch, false, "")
			c := o.Delivery.Candidates[0]
			if o.Delivery.Phase != "automatic_acceptance_passed" || o.DeliveryStatus.HumanJudgment != nil || !strings.Contains(o.NextAction.Message, "local branch delivery") {
				t.Fatalf("branch acquired an artificial human gate or Epic substitution: %+v", o.DeliveryStatus)
			}
			a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
			if err != nil || a.Authorization.AuthorityKind != workspace.DeliveryMandateReceiptV1 || a.Authorization.PermissionConfirmed || !reflect.DeepEqual(a.Authorization.AllowedEffects, []string{workspace.DeliveryLocalBranch}) {
				t.Fatalf("branch lost original native owner/effect: %+v %v", a.Authorization, err)
			}
			if checked, err := workspace.CheckTaskIntegration(f.D.Workspace, autoGateInput(a)); err != nil || checked.Readback.Classification != "ready" {
				t.Fatalf("branch not ready without HumanQA: %+v %v", checked.Readback, err)
			}
			handoff, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
			cleanupOwnerHandoff(t, f, handoff, err)
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
				t.Fatal("exact local branch target did not receive candidate")
			}
			reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+branch)
			cleanupOwnerLeaves(t, f)
			done, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
			if err != nil || done.Delivery.Phase != "completed" || done.DeliveryStatus.FinalDelivery.State != "delivered" {
				t.Fatalf("branch closeout failed: %+v %v", done.DeliveryStatus, err)
			}
			if _, _, err := qualifiedCLIInstallation(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, filepath.Join(f.R.WorkspaceRoot, "installed-ply"), "synthetic coordinator"); err == nil || !strings.Contains(err.Error(), "local Epic delivery") {
				t.Fatalf("branch expanded self-hosted installation scope: %v", err)
			}
			r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if err != nil || len(r.HumanQARecords) != 0 || len(r.TaskResults) != 1 || len(r.IntegrationAttempts) != 1 || len(r.IntegrationResults) != 1 {
				t.Fatalf("branch return lost history or fabricated QA: %v", err)
			}
			base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
			if err != nil || deliveryCLILatestBase(base).OID != c.OID {
				t.Fatalf("branch return left registered base stale: %v", err)
			}
			q, err := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
			if err != nil || q.Current != nil {
				t.Fatalf("branch queue not closed: %v", err)
			}
			if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err != nil || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/"+branch) != reflog {
				t.Fatalf("branch return repeated Git effect: %v", err)
			}
		})
	}
}

func TestPWADep01AutomaticBranchHumanGates(t *testing.T) {
	for _, tc := range []struct {
		name, owner string
		requireQA   bool
	}{{"required QA", "", true}, {"human integration owner", workspace.IntegrationOwnerHuman, false}, {"later veto", "", false}} {
		t.Run(tc.name, func(t *testing.T) {
			f, o, _ := automaticBranchFixture(t, "main", tc.requireQA, tc.owner)
			c := o.Delivery.Candidates[0]
			before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			for _, verdict := range []string{"fail", "blocked"} {
				if tc.name == "later veto" {
					answer := workflowProviderDocument(t, deliveryTestHuman(t, f, o, verdict))
					if verdict == "blocked" {
						answer["completed_at_utc"] = "2026-10-05T00:00:03Z"
					}
					if _, err := workflowhandoff.RecordDeliveryHumanQA(f.D.Workflow, "task", c.TaskResult, verdict, writeAny(t, f.R.WorkspaceRoot, verdict+".json", answer)); err != nil {
						t.Fatal(err)
					}
				}
				if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
					t.Fatal("automatic branch bypassed actual human gate")
				}
				if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
					t.Fatal("human-gated branch changed target")
				}
			}
		})
	}
}

func TestPWADep01AutomaticBranchFreshnessAndAuthority(t *testing.T) {
	f, o, binary := automaticBranchFixture(t, "main", false, "")
	c := o.Delivery.Candidates[0]
	a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
	if err != nil {
		t.Fatal(err)
	}
	assertBlocked := func(d Dependencies) {
		t.Helper()
		if _, err := CompleteLocalDelivery(d, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, ""); err == nil {
			t.Fatal("stale candidate or authority allowed local delivery")
		}
		if gitOutput(t, f.Parent, "rev-parse", "HEAD") != a.ExpectedParentOID {
			t.Fatal("blocked delivery changed target")
		}
	}
	var receipt deliveryVerificationReceipt
	raw, err := os.ReadFile(c.Verification.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &receipt); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{binary, f.AcceptancePath, receipt.Review.Locator} {
		original, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(append([]byte{}, original...), '\n'), 0600); err != nil {
			t.Fatal(err)
		}
		assertBlocked(f.D)
		if err := os.WriteFile(path, original, 0700); err != nil {
			t.Fatal(err)
		}
	}
	input := autoGateInput(a)
	input.DeliveryOwner = nil
	checked, err := workspace.CheckTaskIntegration(f.D.Workspace, input)
	if err == nil && checked.Readback.Classification == "ready" {
		t.Fatal("missing native owner acquired automatic authority")
	}
	// Shared automatic acceptance must not expand the separate headless policy.
	r := f.R
	cfg := *r.Delivery
	r.Delivery = &cfg
	cfg.StartupProtocol = ""
	r.Envelope = deliveryEnv("headless-run-request")
	r.Worker = &WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 60}
	r.Herdr = (WorkflowRequest{}).Herdr
	r.Runtime.Mode, r.Runtime.PermissionBinding.ProfileID = "headless", "ply-worker"
	r.HumanAuthority.StartSurface = "human_authorized_worker"
	if err := validateWorkerRequest(r); err == nil || !strings.Contains(err.Error(), "automatic local Epic agreement") {
		t.Fatalf("branch expanded headless scope: %v", err)
	}
	for _, field := range []string{"mode", "ref", "worktree"} {
		changed := *f.R.Delivery.Agreement
		switch field {
		case "mode":
			changed.Mode = workspace.DeliveryLocalEpic
		case "ref":
			changed.TargetRef = "refs/heads/other"
		case "worktree":
			changed.TargetWorktree = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		}
		if _, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, changed); err == nil {
			t.Fatal("changed delivery agreement acquired authority")
		}
	}
	// Live target drift must stop before effects, even with old native authority.
	runGit(t, f.Parent, "commit", "--allow-empty", "-m", "fixture target drift")
	input = autoGateInput(a)
	if checked, err := workspace.CheckTaskIntegration(f.D.Workspace, input); err == nil && checked.Readback.Classification == "ready" {
		t.Fatal("target drift accepted")
	}
	// An unchanged review cannot qualify a later source commit.
	runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "commit", "--allow-empty", "-m", "fixture stale review")
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); err == nil {
		t.Fatal("later candidate inherited sealed review")
	}
}
