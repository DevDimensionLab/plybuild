package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestDeliveryAgreementModesFreezeAuthorityAndComplete(t *testing.T) {
	for _, tc := range []struct{ mode, parent string }{{workspace.DeliveryLocalEpic, "epic"}, {workspace.DeliveryLocalBranch, "main"}, {workspace.DeliveryLocalBranch, "master"}, {workspace.DeliveryPullRequest, "main"}, {workspace.DeliveryPullRequest, "master"}} {
		t.Run(tc.mode+"/"+tc.parent, func(t *testing.T) {
			agreement := workspace.DeliveryAgreement{Mode: tc.mode}
			if tc.mode == workspace.DeliveryPullRequest {
				agreement.GitHubRepository, agreement.Remote = "fixture/repository", "origin"
			}
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement, ParentRef: tc.parent})
			bound := f.R.Delivery.Agreement
			if bound == nil || bound.SourceRef != "refs/heads/"+f.Prepared.Preparation.Preparation.Plan.Branch || !equal(bound, f.Prepared.Delivery) {
				t.Fatal("agreement lost its prepared Task source")
			}
			if err := workflowhandoff.ValidateDeliveryMandate(f.R.HandoffDraft, bound); err != nil {
				t.Fatal(err)
			}
			changed := *bound
			changed.TargetRef = "refs/heads/another-target"
			if err := workflowhandoff.ValidateDeliveryMandate(f.R.HandoffDraft, &changed); err == nil {
				t.Fatal("changed target acquired frozen authority")
			}
			o := workflowTestStart(t, f.workflowFixture)
			a := deliveryTestAcceptance(t, f, o)
			wanted := append([]string(nil), a.DeliveryPermission.AllowedEffects...)
			a.DeliveryPermission.AllowedEffects = []string{"merge"}
			if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "wrong-effects.json", a)); err == nil {
				t.Fatal("mode-mismatched actual permission was accepted")
			}
			a.DeliveryPermission.AllowedEffects = wanted
			a.DeliveryPermission.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(changed)
			if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "wrong-target.json", a)); err == nil {
				t.Fatal("target-mismatched actual permission was accepted")
			}
			o = deliveryTestAccept(t, f, o)
			before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			registryBefore, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			deliveryTestCandidate(t, f, "selected delivery candidate")
			if tc.mode == workspace.DeliveryPullRequest {
				// Normal PR-base advancement must not force a rebase or change the
				// selected target, including a main/master return context.
				if err = os.WriteFile(filepath.Join(f.Parent, "base-progress.txt"), []byte("unrelated parent progress\n"), 0600); err != nil {
					t.Fatal(err)
				}
				runGit(t, f.Parent, "add", "base-progress.txt")
				runGit(t, f.Parent, "commit", "-m", "normal PR base progress")
				before = gitOutput(t, f.Parent, "rev-parse", "HEAD")
			}
			o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "mode"))
			if err != nil {
				t.Fatal(err)
			}
			candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			authority, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult, *bound)
			if err != nil {
				t.Fatal(err)
			}
			if authority.Candidate.HumanQA != nil || authority.Authorization.CandidateRunID != candidate.TaskResult.RunID {
				t.Fatal("technical qualification fabricated human QA or lost candidate identity")
			}
			if _, err = ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult, changed); err == nil {
				t.Fatal("registration accepted a changed delivery target")
			}
			if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult.ID, ""); err == nil {
				t.Fatal("delivery without human QA had an effect")
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
				t.Fatal("missing QA changed the parent")
			}
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
			if err != nil {
				t.Fatal(err)
			}
			if tc.mode == workspace.DeliveryPullRequest {
				if !strings.Contains(o.NextAction.Message, "agreed PR") || strings.Contains(o.NextAction.Message, "local integration") {
					t.Fatalf("human pass directs PR owner to wrong delivery effect: %s", o.NextAction.Message)
				}
				if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
					t.Fatal("PR authority allowed local integration")
				}
				observed := PullRequestDeliveryObservation{DeliveryID: "fixture-delivery", Repository: bound.GitHubRepository, HeadRef: bound.SourceRef, HeadOID: candidate.OID, BaseRef: bound.TargetRef, URL: "https://github.com/fixture/repository/pull/17", Number: 17, MetadataApplied: true}
				receiptPath := writeAny(t, f.R.WorkspaceRoot, "observed-pr.json", observed)
				receipt := FileBinding{receiptPath, hashFileTest(t, receiptPath)}
				o, err = CompletePullRequestDelivery(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult.ID, receipt, observed)
				if err != nil {
					t.Fatal(err)
				}
				if o.Delivery.Phase != "completed" || o.Delivery.Candidates[0].Integration != nil || o.Delivery.Candidates[0].PullRequest == nil {
					t.Fatal("PR completion was not truthful")
				}
				if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
					t.Fatal("PR completion altered the local target")
				}
				registryAfter, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				if !equal(registryBefore.EpicBaseVersions, registryAfter.EpicBaseVersions) || len(registryAfter.IntegrationResults) != 0 {
					t.Fatal("PR completion claimed an Epic base update or integration")
				}
				events := len(o.Delivery.Events)
				o, err = CompletePullRequestDelivery(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult.ID, receipt, observed)
				if err != nil || len(o.Delivery.Events) != events {
					t.Fatalf("PR completion retry duplicated event: %v", err)
				}
			} else {
				qa := o.Delivery.Candidates[0].HumanQA
				completed, err := CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult.ID, qa.ID)
				if err != nil || !completed.Completed {
					t.Fatalf("local native delivery failed: %+v %v", completed, err)
				}
				if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate.OID {
					t.Fatal("local exact target did not receive the accepted candidate")
				}
				if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, candidate.TaskResult.ID, qa.ID); err != nil {
					t.Fatalf("native local retry: %v", err)
				}
			}
		})
	}
}

func TestDeliveryLegacyCannotAcquireStructuredAuthority(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "legacy candidate")
	var err error
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "legacy"))
	if err != nil {
		t.Fatal(err)
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryPullRequest, ProjectID: "ply", RepoID: "ply", EpicID: "epic", SourceRef: o.Delivery.Candidates[0].TaskResult.SourceRef, TargetRef: "refs/heads/main", GitHubRepository: "fixture/repository", Remote: "origin"}
	if _, err = ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, o.Delivery.Candidates[0].TaskResult, a); err == nil || !strings.Contains(err.Error(), "historical") {
		t.Fatalf("old run acquired PR authority: %v", err)
	}
	var draft map[string]any
	if err = json.Unmarshal(f.R.HandoffDraft, &draft); err != nil {
		t.Fatal(err)
	}
	if _, exists := draft["delivery_binding"].(map[string]any)["agreement"]; exists {
		t.Fatal("legacy immutable mandate gained an agreement")
	}
}

func TestDeliveryFailedHumanCheckRevokesDirectIntegrationAuthority(t *testing.T) {
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}, ParentRef: "main"})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "candidate with subsequent human fail")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "revoked"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "fail", deliveryTestHuman(t, f, o, "fail"))
	if err != nil {
		t.Fatal(err)
	}
	// A caller may replay the old native pass and old authority object, but the
	// evidence reader must withhold current authority after the actual fail.
	evidence, err := workflowhandoff.NewTaskHandoffEvidenceReader(f.D.Workflow).ReadTaskEvidence(taskResultEvidenceRequest(c.TaskResult))
	if err != nil || !evidence.TaskSpecValid || evidence.DeliveryAuthorization != nil {
		t.Fatalf("historical evidence acquired live authority: %+v %v", evidence.DeliveryAuthorization, err)
	}
	in, err := workspace.ParseTaskIntegrationInput("task", string(c.TaskResult.ID), string(c.HumanQA.ID), c.OID, a.ExpectedParentOID, "", false, "")
	if err != nil {
		t.Fatal(err)
	}
	in.DeliveryAuthorization = a.Authorization
	in.DeliveryOwner = &workspace.DeliveryIntegrationOwner{RunID: a.RunID, RequestSHA256: a.RequestSHA256, ActorClaim: a.OwnerClaim, PreparationID: a.PreparationID}
	checked, err := workspace.CheckTaskIntegration(f.D.Workspace, in)
	if err == nil && (checked.Readback.Classification == "ready" || checked.Readback.RecoveryStatus == "complete") {
		t.Fatal("old pass bypassed failed current candidate through direct native integration")
	}
	if got := gitOutput(t, f.Parent, "rev-parse", "HEAD"); got != a.ExpectedParentOID {
		t.Fatal("revoked authority changed the parent")
	}
}

func TestPullRequestCompletionRejectsChangedHumanEvidence(t *testing.T) {
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryPullRequest, GitHubRepository: "fixture/repository", Remote: "origin"}, ParentRef: "main"})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "candidate with exact human evidence")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "evidence"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	evidence := c.HumanQA.Evidence[0]
	if err = os.WriteFile(evidence.Locator, []byte("changed candidate evidence"), 0600); err != nil {
		t.Fatal(err)
	}
	a := *f.R.Delivery.Agreement
	observed := PullRequestDeliveryObservation{DeliveryID: "fixture", Repository: a.GitHubRepository, HeadRef: a.SourceRef, HeadOID: c.OID, BaseRef: a.TargetRef, URL: "https://github.com/fixture/repository/pull/17", Number: 17, MetadataApplied: true}
	path := writeAny(t, f.R.WorkspaceRoot, "pr-evidence.json", observed)
	if _, err = CompletePullRequestDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, FileBinding{path, hashFileTest(t, path)}, observed); err == nil {
		t.Fatal("changed human evidence completed PR delivery")
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.Delivery.Phase == "completed" || shown.Delivery.Candidates[0].PullRequest != nil {
		t.Fatalf("invalid PR completion was preserved: %+v %v", shown.Delivery, err)
	}
}

type reversedDeliveryQAIDs struct {
	workspace.TaskLifecycleIDSource
	count int
}

func (s *reversedDeliveryQAIDs) NewHumanQARecordID() (workspace.HumanQARecordID, error) {
	s.count++
	if s.count == 1 {
		return workspace.HumanQARecordID("hqa_" + strings.Repeat("f", 32)), nil
	}
	return workspace.HumanQARecordID("hqa_" + strings.Repeat("0", 32)), nil
}

func TestDeliveryExternalLaterHumanFailOverridesEarlierPassRegardlessRecordID(t *testing.T) {
	for _, mode := range []string{workspace.DeliveryLocalBranch, workspace.DeliveryPullRequest} {
		t.Run(mode, func(t *testing.T) {
			agreement := workspace.DeliveryAgreement{Mode: mode}
			if mode == workspace.DeliveryPullRequest {
				agreement.GitHubRepository, agreement.Remote = "fixture/repository", "origin"
			}
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement, ParentRef: "main"})
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			deliveryTestCandidate(t, f, "candidate with externally recorded newer fail")
			o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "qa-order"))
			if err != nil {
				t.Fatal(err)
			}
			c := o.Delivery.Candidates[0]
			wd := f.D.Workspace
			wd.TaskLifecycleIDs = &reversedDeliveryQAIDs{TaskLifecycleIDSource: wd.TaskLifecycleIDs}
			wf := f.D.Workflow
			wf.TaskWorkspace = &wd
			pass, err := workflowhandoff.RecordDeliveryHumanQA(wf, "task", c.TaskResult, "pass", deliveryTestHuman(t, f, o, "pass"))
			if err != nil {
				t.Fatal(err)
			}
			failurePath := deliveryTestHuman(t, f, o, "fail")
			failure := workflowProviderDocument(t, failurePath)
			failure["completed_at_utc"] = "2026-10-05T00:00:02Z"
			failurePath = writeAny(t, f.R.WorkspaceRoot, "external-later-fail.json", failure)
			failed, err := workflowhandoff.RecordDeliveryHumanQA(wf, "task", c.TaskResult, "fail", failurePath)
			if err != nil {
				t.Fatal(err)
			}
			if failed.ID >= pass.ID {
				t.Fatal("fixture did not reverse ID ordering")
			}
			parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			if mode == workspace.DeliveryLocalBranch {
				if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, pass.ID); err == nil {
					t.Error("earlier human pass overrode later native fail")
				}
			} else {
				a := *f.R.Delivery.Agreement
				observed := PullRequestDeliveryObservation{DeliveryID: "fixture", Repository: a.GitHubRepository, HeadRef: a.SourceRef, HeadOID: c.OID, BaseRef: a.TargetRef, URL: "https://github.com/fixture/repository/pull/17", Number: 17, MetadataApplied: true}
				path := writeAny(t, f.R.WorkspaceRoot, "old-pass-pr.json", observed)
				if _, err = CompletePullRequestDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, FileBinding{path, hashFileTest(t, path)}, observed); err == nil {
					t.Error("earlier human pass overrode later native fail")
				}
			}
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || shown.Delivery.Phase == "completed" {
				t.Errorf("later fail closed native workflow: %v", err)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
				t.Error("later fail allowed a local target effect")
			}
		})
	}
}

func TestDeliveryTechnicalEvidencePreservesNativeAcceptedDebt(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*workspace.TaskResultRecord, *workspace.TaskHandoffEvidence)
		wantOK bool
	}{
		{"named evidence debt", nil, true},
		{"missing named debt", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) { r.AcceptedDebt = nil }, false},
		{"uncovered gap", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) {
			e.EvidenceGaps = append(e.EvidenceGaps, "another evidence gap")
		}, false},
		{"unmanaged debt evidence", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) { r.Artifacts[0].Role = "other" }, false},
		{"debt disguised as passed", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) {
			r.TechnicalGate, e.EvidenceCoverageValid = "passed", true
		}, false},
		{"failed functional requirement", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) { e.TaskRequirementsValid = false }, false},
		{"changed verifier", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) { e.VerifierResults[0].Exit = 8 }, false},
		{"high risk debt", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) {
			e.Review.Findings = []workspace.TaskReviewEntry{{Severity: "high"}}
			r.Review = e.Review
		}, false},
		{"complete passed evidence", func(r *workspace.TaskResultRecord, e *workspace.TaskHandoffEvidence) {
			r.TechnicalGate, r.AcceptedDebt, e.EvidenceCoverageValid, e.EvidenceGaps = "passed", nil, true, nil
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// This isolates assessment revalidation from native candidate creation,
			// which currently publishes only the passed technical gate.
			verifier := workspace.TaskVerifierResultRecord{VerifierID: "acceptance", Outcome: "passed", Exit: 0}
			artifact := workspace.TaskArtifactRecord{ArtifactID: "debt-evidence", Role: "debt_control", Locator: "/fixture/debt.json", SHA256: "sha256:" + strings.Repeat("a", 64), SizeBytes: 12}
			authority := &workspace.DeliveryAuthorization{RunID: "native-authority"}
			r := workspace.TaskResultRecord{TechnicalGate: "good_enough_with_known_debt", ReportedOutcome: "complete", VerifierResults: []workspace.TaskVerifierResultRecord{verifier}, Artifacts: []workspace.TaskArtifactRecord{artifact}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{{ID: "optional-evidence", RiskClass: "evidence", Severity: "low", Summary: "Optional evidence unavailable", Control: "Recorded fallback observation", EvidenceArtifactIDs: []string{"debt-evidence"}}}}
			e := workspace.TaskHandoffEvidence{SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, TaskSpecValid: true, TaskRequirementsValid: true, ReportedOutcome: "complete", ExpectedVerifierIDs: []string{"acceptance"}, VerifierResults: []workspace.TaskVerifierResultRecord{verifier}, Artifacts: []workspace.TaskArtifactRecord{artifact}, EvidenceGaps: []string{"optional evidence unavailable"}, DeliveryAuthorization: authority}
			if tc.mutate != nil {
				tc.mutate(&r, &e)
			}
			if err := validateDeliveryTechnicalEvidence(r, e, authority); (err == nil) != tc.wantOK {
				t.Fatalf("native technical gate acceptance = %v, want %v: %v", err == nil, tc.wantOK, err)
			}
		})
	}
}
