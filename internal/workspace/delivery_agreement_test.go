package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func fixtureAgreement(mode, target string) DeliveryAgreement {
	a := DeliveryAgreement{SchemaVersion: 1, Mode: mode, ProjectID: "ply", RepoID: "ply", EpicID: "epic", TargetRef: target}
	if mode == DeliveryPullRequest {
		a.Remote, a.GitHubRepository = "origin", "Fixture/Repository"
	}
	return a
}

func TestDeliveryAgreementSurvivesPublicationQueueAndFrozenExecution(t *testing.T) {
	for _, mode := range []string{DeliveryPullRequest, DeliveryLocalEpic, DeliveryLocalBranch} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			a := fixtureAgreement(mode, "refs/heads/epic")
			if mode != DeliveryPullRequest {
				a.TargetWorktree = filepath.Join(f.wrapper, "epic")
			}
			goal := goalFixture(t, f, "task", "goal", a)
			q := goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			if len(q.Pending) != 1 || !contentTypedEqual(q.Pending[0].Delivery, &a) {
				t.Fatal("queue lost exact agreement")
			}
			first, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
			if err != nil {
				t.Fatal(err)
			}
			if !contentTypedEqual(first.Goal.Delivery, &a) || first.Delivery == nil || first.Delivery.SourceRef != "refs/heads/"+first.Branch {
				t.Fatal("preview did not preserve goal and bind actual source")
			}
			prepared, err := PrepareTaskGoalExecution(f.dependencies, first.Input, first.Confirmation, goalHumanFixture())
			if err != nil {
				t.Fatal(err)
			}
			if !contentTypedEqual(prepared.Delivery, first.Delivery) {
				t.Fatal("preparation changed delivery")
			}
			r := mustQueueRegistry(t, f)
			spec, err := readRegisteredTaskManifest(f.dependencies.TaskContent, filepath.Dir(f.wrapper), r, prepared.Basis.Spec.ManifestSHA256)
			if err != nil {
				t.Fatal(err)
			}
			got, err := DeliveryAgreementFromSpec(spec)
			if err != nil || !contentTypedEqual(got, first.Delivery) {
				t.Fatalf("frozen Spec lost agreement: %v", err)
			}
			var frozen TaskGoalExecutePreview
			if err = contentDecode(contentFields(spec)["execution_request"], &frozen); err != nil {
				t.Fatal(err)
			}
			changed := *frozen.Delivery
			changed.TargetRef = "refs/heads/other"
			frozen.Delivery = &changed
			if validateGoalDeliveryBinding(frozen) == nil {
				t.Fatal("execution accepted target switch without goal revision")
			}
			// A subsequent published goal changes no bytes of the already-started run.
			a.TargetRef = "refs/heads/other"
			goalFixture(t, f, "task", "revised-goal", a)
			old, _, err := storedGoalExecution(f.dependencies, filepath.Dir(f.wrapper), mustQueueRegistry(t, f), first.Input, first.Confirmation)
			if err != nil || old == nil || !contentTypedEqual(old.Delivery, first.Delivery) {
				t.Fatalf("later goal revision changed frozen execution: %v", err)
			}
		})
	}
}

func deliveryIntegrationInput(f integrationJourneyFixture) TaskIntegrationInput {
	return TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid, DeliveryAuthorization: f.evidence.evidence.DeliveryAuthorization}
}

func TestLocalBranchDeliveryIntegratesMainAndMasterWithoutRemoteAndPreservesRemote(t *testing.T) {
	for _, branch := range []string{"main", "master"} {
		for _, remote := range []bool{false, true} {
			t.Run(branch+map[bool]string{false: "/offline", true: "/with-remote"}[remote], func(t *testing.T) {
				f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryLocalBranch, "refs/heads/"+branch))
				parent := f.evidence.evidence.DeliveryAuthorization.Agreement.TargetWorktree
				if remote {
					runLocalGit(t, parent, "remote", "add", "origin", filepath.Join(f.root, "unreachable-remote.git"))
					runLocalGit(t, parent, "update-ref", "refs/remotes/origin/"+branch, f.oid)
				}
				git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
				f.dependencies.IntegrationGit = git
				input := deliveryIntegrationInput(f)
				check, err := CheckTaskIntegration(f.dependencies, input)
				if err != nil || check.Readback.Classification != "ready" {
					t.Fatalf("explicit local product agreement not ready: %#v %v", check.Readback, err)
				}
				input.Apply, input.Confirmation = true, integrationDigest(t, check.Readback)
				applied, err := ApplyTaskIntegration(f.dependencies, input)
				if err != nil || applied.Readback.Classification != "exact_effect" || git.count() != 1 {
					t.Fatalf("local product delivery: %#v %v", applied.Readback, err)
				}
				if got := runLocalGit(t, parent, "rev-parse", "HEAD"); got != f.resultOID {
					t.Fatal("actual target was not integrated")
				}
				if remote && runLocalGit(t, parent, "rev-parse", "refs/remotes/origin/"+branch) != f.oid {
					t.Fatal("local delivery changed remote refs")
				}
				q, err := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
				if err != nil {
					t.Fatal(err)
				}
				_, err = CloseTaskQueue(f.dependencies, queueTargetFixture(), q.Current.PreparationID, q.Revision, "Fixture observed local Delivery", "advance")
				if err != nil {
					t.Fatal(err)
				}
				base, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
				if err != nil {
					t.Fatal(err)
				}
				_, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation})
				if err != nil {
					t.Fatal(err)
				}
				r := mustQueueRegistry(t, f.workItemJourneyFixture)
				ep, p := queueTargetBinding(r, q.Target)
				if currentEpicBase(r, ep, p).OID != f.resultOID || foldQueue(r, q.QueueID).Current != nil {
					t.Fatal("native completion did not preserve base and queue facts")
				}
				if _, err = CheckTaskIntegration(f.dependencies, input); err != nil || git.count() != 1 {
					t.Fatalf("recovery did not reuse effect: %v", err)
				}
			})
		}
	}
}

func TestLocalBranchDeliveryRejectsUnauthenticatedOrChangedAgreement(t *testing.T) {
	for _, defect := range []string{"missing", "forged-evidence", "wrong-run", "wrong-target", "wrong-worktree", "wrong-mode", "dirty-parent", "moved-parent", "changed-candidate"} {
		t.Run(defect, func(t *testing.T) {
			f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryLocalBranch, "refs/heads/main"))
			input := deliveryIntegrationInput(f)
			copy := *input.DeliveryAuthorization
			input.DeliveryAuthorization = &copy
			parent := copy.Agreement.TargetWorktree
			switch defect {
			case "missing":
				input.DeliveryAuthorization = nil
			case "forged-evidence":
				f.evidence.evidence.DeliveryAuthorization = nil
			case "wrong-run":
				copy.CandidateRunID = "run_other"
			case "wrong-target":
				copy.Agreement.TargetRef = "refs/heads/other"
			case "wrong-worktree":
				copy.Agreement.TargetWorktree = f.taskPath
			case "wrong-mode":
				copy.Agreement.Mode = DeliveryLocalEpic
			case "dirty-parent":
				if err := os.WriteFile(filepath.Join(parent, "dirty.txt"), []byte("dirty"), 0600); err != nil {
					t.Fatal(err)
				}
			case "moved-parent":
				runLocalGit(t, parent, "commit", "--allow-empty", "-m", "parent moved")
			case "changed-candidate":
				runLocalGit(t, f.taskPath, "commit", "--allow-empty", "-m", "candidate moved")
			}
			copy.AgreementSHA256 = DeliveryAgreementDigest(copy.Agreement)
			before := runLocalGit(t, parent, "rev-parse", "HEAD")
			checked, err := CheckTaskIntegration(f.dependencies, input)
			if err == nil && (checked.Readback.Classification == "ready" || checked.Readback.Classification == "already_integrated") {
				t.Fatalf("unsafe %s was ready", defect)
			}
			if runLocalGit(t, parent, "rev-parse", "HEAD") != before {
				t.Fatal("failed check changed target")
			}
		})
	}
}

func TestLocalBranchStoredPlanCannotAcquireMainAuthorityByRemovingContract(t *testing.T) {
	f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryLocalBranch, "refs/heads/master"))
	input := deliveryIntegrationInput(f)
	check, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Confirmation = integrationDigest(t, check.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	result, qa := r.TaskResults[0], r.HumanQARecords[0]
	plan := r.IntegrationAuthorities[0].Plan
	if err = validateStoredIntegrationPlan(plan, r, result, qa); err != nil {
		t.Fatal(err)
	}
	for _, defect := range []string{"remove", "version", "target", "run", "permission", "effect"} {
		bad := plan
		auth := *plan.DeliveryAuthorization
		bad.DeliveryAuthorization = &auth
		switch defect {
		case "remove":
			bad.DeliveryAuthorization = nil
		case "version":
			bad.Kind, bad.SchemaVersion = "WorkspaceTaskIntegrationPlan@2", 2
		case "target":
			auth.Agreement.TargetRef = "refs/heads/main"
		case "run":
			auth.CandidateRunID = "run_other"
		case "permission":
			auth.PermissionConfirmed = false
		case "effect":
			auth.AllowedEffects = []string{DeliveryLocalEpic}
		}
		auth.AgreementSHA256 = DeliveryAgreementDigest(auth.Agreement)
		if validateStoredIntegrationPlan(bad, r, result, qa) == nil {
			t.Fatalf("stored plan accepted %s bypass", defect)
		}
	}
}

func TestDeliveryAgreementStrictShapeAndHistoricalAbsence(t *testing.T) {
	if a, err := DeliveryAgreementFromSpec(canonicaljson.Object{}); err != nil || a != nil {
		t.Fatal("historical contract acquired authority")
	}
	a := fixtureAgreement(DeliveryPullRequest, "refs/heads/master")
	a.SourceRef = "refs/heads/task"
	v := contentFields(contentRefValue(a))
	v["surprise_authority"] = true
	if deliveryAgreementRule(contentObject(v)) == nil {
		t.Fatal("unknown authority field accepted")
	}
	for _, mode := range []string{DeliveryLocalEpic, "", "local"} {
		a := fixtureAgreement(mode, "refs/heads/main")
		a.TargetWorktree = "/tmp/main"
		if ValidateDeliveryAgreement(a) == nil {
			t.Fatalf("%s gained main permission", mode)
		}
	}
	if !strings.HasPrefix(DeliveryAgreementDigest(fixtureAgreement(DeliveryPullRequest, "refs/heads/main")), "sha256:") {
		t.Fatal("agreement identity missing")
	}
}

func TestStoredDeliveryCannotRetrofitAgreementIntoLegacySpec(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Confirmation = integrationDigest(t, checked.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	a := fixtureAgreement(DeliveryLocalBranch, "refs/heads/epic")
	a.SourceRef, a.TargetWorktree = r.TaskResults[0].SourceRef, r.Epics[0].RepoBindings[0].Worktree.Locator
	plan := &r.IntegrationAuthorities[0].Plan
	plan.Kind, plan.SchemaVersion = "WorkspaceTaskIntegrationPlan@3", 3
	plan.DeliveryAuthorization = &DeliveryAuthorization{Agreement: a, AgreementSHA256: DeliveryAgreementDigest(a), RunID: "wfr_forged", CandidateRunID: r.TaskResults[0].RunID, RequestSHA256: digestTaskBytes([]byte("request")), MandateSHA256: digestTaskBytes([]byte("mandate")), AllowedEffects: a.AllowedEffects(), PermissionConfirmed: true}
	if validateTaskContentClosure(f.dependencies.TaskContent, f.root, r, false) == nil {
		t.Fatal("storage accepted an invented explicit agreement over a historical Spec")
	}
}

func TestPullRequestCandidateAllowsNormalBaseProgressButNoLocalIntegration(t *testing.T) {
	f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryPullRequest, "refs/heads/main"))
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	task, _ := findTask(r, "task")
	parent := r.Epics[0].RepoBindings[0].Worktree.Locator
	runLocalGit(t, parent, "commit", "--allow-empty", "-m", "normal base progress")
	err := WithTaskSpecSnapshot(f.dependencies, f.root, func(s *TaskSpecSession) error {
		_, err := s.ValidateDeliveryCandidateTarget(f.taskPath, task.Worktree.Ref, task.GitCommonDir, f.evidence.evidence.TaskSpecBasis)
		return err
	})
	if err != nil {
		t.Fatalf("normal PR base progress invalidated exact candidate: %v", err)
	}
	before := runLocalGit(t, parent, "rev-parse", "HEAD")
	if _, err := CheckTaskIntegration(f.dependencies, deliveryIntegrationInput(f)); err == nil {
		t.Fatal("PR agreement authorized local integration")
	}
	if runLocalGit(t, parent, "rev-parse", "HEAD") != before {
		t.Fatal("PR touched local parent")
	}
}
