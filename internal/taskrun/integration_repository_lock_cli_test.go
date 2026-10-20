package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type heldIntegrationEvidence struct{ evidence workspace.TaskHandoffEvidence }

func (r heldIntegrationEvidence) ReadTaskEvidence(workspace.TaskHandoffEvidenceRequest) (workspace.TaskHandoffEvidence, error) {
	return r.evidence, nil
}

type attemptPublicationFailureStore struct {
	workspace.WorkItemStore
	publishes *int
}

func (s attemptPublicationFailureStore) WithLock(root string, action func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		return action(attemptPublicationFailureSession{session, s.publishes})
	})
}

type attemptPublicationFailureSession struct {
	workspace.WorkItemStoreSession
	publishes *int
}

func (s attemptPublicationFailureSession) Publish(registry workspace.WorkItemRegistry) error {
	*s.publishes++
	if *s.publishes == 2 {
		return errors.New("synthetic interruption after durable integration authority")
	}
	return s.WorkItemStoreSession.Publish(registry)
}

// heldIntegration registers a Task under epic with a synthetic result and QA.
// The returned start begins its integration in process and interrupts it after
// the durable authority: the production state of an integration in flight.
func heldIntegration(t *testing.T, w workspace.Dependencies, root string, project workspace.ProjectID, repo workspace.RepoID, epic workspace.EpicID, parentOID string) (start func()) {
	t.Helper()
	if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "held", Title: "Held integration", Description: "Synthetic Task whose integration stays in flight.", ParentEpicID: epic, ProjectID: project, RepoID: repo}); err != nil {
		t.Fatal(err)
	}
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	var problem workspace.TaskRevisionRef
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == "held" {
			problem = workspace.TaskRevisionRef{Revision: p.Revision, ManifestSHA256: p.ManifestSHA256}
		}
	}
	doc := filepath.Join(root, "held-goal.md")
	contents := []byte("# Held goal\nSynthetic integration kept in flight.\n")
	if err = os.WriteFile(doc, contents, 0600); err != nil {
		t.Fatal(err)
	}
	draft := map[string]any{
		"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "format": "json", "format_version": 1, "canonicalization": "RFC8785",
		"publication_key": "fixture/held-goal", "task_id": "held", "spec_id": "solution", "registry_upgrade": nil, "expected_previous": nil, "problem": problem,
		"recorder":      map[string]any{"actor_claim": "fixture planner", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"},
		"contract_kind": "goal", "title": "Held integration", "objective": "Hold an integration lock without attributing synthetic evidence to a real human.", "change_reason": "Synthetic fixture goal",
		"documents": []any{map[string]any{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(contents), "size_bytes": len(contents), "media_type": "text/markdown", "git_provenance": nil}}},
		"design":    []any{map[string]any{"document_id": "design", "section": nil}}, "constraints": []string{"Local isolated fixture only."},
		"requirements": []any{map[string]any{"id": "behavior", "acceptance": "The synthetic result exists."}},
		"executor":     map[string]any{"provider": "codex", "model": "fixture-model", "effort": "medium"},
	}
	publication, err := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "held", File: writeAny(t, root, "held-goal-draft.json", draft)})
	if err != nil {
		t.Fatal(err)
	}
	goal := workspace.TaskGoalRef{SpecID: "solution", Spec: workspace.TaskRevisionRef{Revision: *publication.OutcomeRef.Revision, ManifestSHA256: publication.OutcomeRef.ManifestSHA256}}
	queue := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/held-goals", ProjectID: project, RepoID: repo, EpicID: epic, Entries: []workspace.TaskGoalQueueEntry{{TaskID: "held", Goal: &goal}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}}
	if _, err = workspace.SetTaskQueue(w, writeAny(t, root, "held-queue.json", queue)); err != nil {
		t.Fatal(err)
	}
	preview, err := workspace.PreviewTaskGoalExecution(w, workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: project, RepoID: repo, EpicID: epic}, Next: true, AcceptancePath: filepath.Join(root, "held-acceptance.sh")})
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := workspace.PrepareTaskGoalExecution(w, preview.Input, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "synthetic human", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Synthetic execution choice; not actual product QA."})
	if err != nil {
		t.Fatal(err)
	}
	path := prepared.Preparation.Preparation.Plan.WorktreePath
	if err = os.WriteFile(filepath.Join(path, "held.txt"), []byte("held\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, path, "add", "held.txt")
	runGit(t, path, "commit", "-m", "held result")
	resultOID, resultTree := gitOutput(t, path, "rev-parse", "HEAD"), gitOutput(t, path, "rev-parse", "HEAD^{tree}")
	r, err = w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	var task workspace.TaskRecord
	for _, candidate := range r.Tasks {
		if candidate.ID == "held" {
			task = candidate
		}
	}
	digest := "sha256:" + strings.Repeat("1", 64)
	resultID, qaID := workspace.TaskResultID("trs_"+strings.Repeat("1", 32)), workspace.HumanQARecordID("hqa_"+strings.Repeat("1", 32))
	report := []byte("Synthetic QA report; not actual product QA.\n")
	reportPath := filepath.Join(root, "held-qa.txt")
	if err = os.WriteFile(reportPath, report, 0600); err != nil {
		t.Fatal(err)
	}
	evidence := workspace.TaskHandoffEvidence{
		TaskSpecBasis: &prepared.Basis, TaskSpecValid: true, TaskRequirementsValid: true, AcceptedStartOutcome: "started",
		ActivityID: "act_" + strings.Repeat("1", 32), RunID: "run_" + strings.Repeat("1", 32), HandoffID: "hnd_" + strings.Repeat("1", 32), HandoffLocator: filepath.Join(root, "held-handoff"), HandoffSHA256: digest,
		StartReceiptID: "rcp_" + strings.Repeat("1", 32), StartReceiptLocator: filepath.Join(root, "held-start"), StartReceiptSHA256: digest,
		TerminalResultID: "res_" + strings.Repeat("1", 32), TerminalResultLocator: filepath.Join(root, "held-terminal"), TerminalResultSHA256: digest, InspectionSHA256: digest,
		ReportedOutcome: "complete", TargetWorktree: path, TargetRef: task.Worktree.Ref, ResultOID: resultOID, ResultTree: resultTree, GitCommonDir: task.GitCommonDir,
		SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true,
		VerifierResults: []workspace.TaskVerifierResultRecord{}, ExpectedVerifierIDs: []string{}, Review: workspace.TaskReviewRecord{Findings: []workspace.TaskReviewEntry{}, Fixes: []workspace.TaskReviewEntry{}, OpenActionableFindings: []workspace.TaskReviewEntry{}}, Artifacts: []workspace.TaskArtifactRecord{}, EvidenceGaps: []string{},
	}
	err = w.WorkItems.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		r.TaskResults = append(r.TaskResults, workspace.TaskResultRecord{ID: resultID, PublicationKey: "held/result", DraftSHA256: digest, StoreTransition: "none", TaskID: "held", TaskWorktreeID: task.Worktree.ID, ProjectID: project, RepoID: repo, GitCommonDir: task.GitCommonDir, SourceLocator: path, SourceRef: task.Worktree.Ref, ResultOID: resultOID, ResultTree: resultTree, ActivityID: evidence.ActivityID, RunID: evidence.RunID, HandoffID: evidence.HandoffID, HandoffLocator: evidence.HandoffLocator, HandoffSHA256: digest, StartReceiptID: evidence.StartReceiptID, StartReceiptLocator: evidence.StartReceiptLocator, StartReceiptSHA256: digest, TerminalResultID: evidence.TerminalResultID, TerminalResultLocator: evidence.TerminalResultLocator, TerminalResultSHA256: digest, InspectionSHA256: digest, ReportedOutcome: "complete", TechnicalGate: "passed", VerifierResults: []workspace.TaskVerifierResultRecord{}, Review: evidence.Review, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{}, Artifacts: []workspace.TaskArtifactRecord{}, Recorder: workspace.TaskRecorderRecord{ActorClaim: "fixture", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}})
		r.TaskResultSpecBindings = append(r.TaskResultSpecBindings, workspace.TaskResultSpecBinding{TaskID: "held", TaskResultID: resultID, Basis: prepared.Basis, HandoffSHA256: digest, StartReceiptSHA256: digest})
		r.HumanQARecords = append(r.HumanQARecords, workspace.TaskHumanQARecord{ID: qaID, PublicationKey: "held/qa", DraftSHA256: digest, TaskID: "held", TaskResultID: resultID, ResultOID: resultOID, ResultTree: resultTree, Outcome: "pass", Actor: workspace.TaskHumanActorRecord{ActorClaim: "synthetic human, not actual approval", StartSurface: "synthetic test", StartedAtUTC: "2026-10-05T00:00:00Z", CompletedAtUTC: "2026-10-05T00:01:00Z"}, Evidence: []workspace.TaskHumanQAEvidenceRecord{{ID: "report", Role: "report", Locator: reportPath, SHA256: hash(report), SizeBytes: int64(len(report))}}, Observation: "synthetic", AcceptedResidualRisks: []workspace.TaskResidualRiskRecord{}})
		return session.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		held := w
		held.HandoffEvidence = heldIntegrationEvidence{evidence}
		input := workspace.TaskIntegrationInput{TaskID: "held", TaskResultID: resultID, HumanQARecordID: qaID, ExpectedResultOID: resultOID, ExpectedParentOID: parentOID}
		checked, err := workspace.CheckTaskIntegration(held, input)
		if err != nil || checked.Readback.Classification != "ready" {
			t.Fatalf("held integration not ready: %+v %v", checked.Readback.Classification, err)
		}
		input.Apply, input.Confirmation = true, checked.Readback.NextAction.Argv[len(checked.Readback.NextAction.Argv)-1]
		publishes := 0
		held.WorkItems = attemptPublicationFailureStore{w.WorkItems, &publishes}
		if _, err = workspace.ApplyTaskIntegration(held, input); err == nil {
			t.Fatal("interrupted held integration completed")
		}
		r, err = w.WorkItems.Snapshot(root)
		if err != nil || len(r.IntegrationAuthorities) != 1 || len(r.IntegrationResults) != 0 {
			t.Fatalf("held integration is not in flight: authorities=%d results=%d %v", len(r.IntegrationAuthorities), len(r.IntegrationResults), err)
		}
	}
}

// IL01 public CLI journey: an integration in flight for one repository does not
// block another repository's same target ref, while the same repository and
// target stays serialized. The held integration is synthetic; the exact CLI
// candidate's acceptance, Git integration and native closure are real.
func TestIntegrationLockCLIScopesInFlightIntegrationToRepositoryAndTarget(t *testing.T) {
	for _, scenario := range []struct {
		name           string
		sameRepository bool
	}{{"independent repository same ref", false}, {"same repository and ref", true}} {
		t.Run(scenario.name, func(t *testing.T) {
			binary := automaticCLIBinary(t)
			heldParent := ""
			var hold func()
			a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic repository lock orchestrator", TimeoutSeconds: 30}}
			f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: a, BeforePrepare: func(w workspace.Dependencies, root, epic string, _ workspace.TaskGoalRef) {
				oid := gitOutput(t, epic, "rev-parse", "HEAD")
				if scenario.sameRepository {
					heldParent = epic
					if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "sibling", Title: "Sibling Epic", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: oid}); err != nil {
						t.Fatal(err)
					}
					hold = heldIntegration(t, w, root, "ply", "ply", "sibling", oid)
					return
				}
				wrapper := filepath.Join(root, "other")
				main := filepath.Join(wrapper, "main")
				if err := os.MkdirAll(main, 0o755); err != nil {
					t.Fatal(err)
				}
				runGit(t, main, "init", "-b", "main")
				runGit(t, main, "config", "user.name", "IL01 Test")
				runGit(t, main, "config", "user.email", "il01@example.invalid")
				if err := os.WriteFile(filepath.Join(main, "README.md"), []byte("other\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				runGit(t, main, "add", "README.md")
				runGit(t, main, "commit", "-m", "other base")
				otherOID := gitOutput(t, main, "rev-parse", "HEAD")
				heldParent = filepath.Join(wrapper, "epic")
				runGit(t, main, "worktree", "add", "-b", "epic", heldParent, otherOID)
				if _, err := workspace.AddProject(w, workspace.ProjectAddInput{ProjectID: "other", Name: "Other", Wrapper: wrapper, Repos: []workspace.RepoInput{{RepoID: "other", Path: main}}}); err != nil {
					t.Fatal(err)
				}
				if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "other-epic", Title: "Other Epic", ProjectID: "other", RepoID: "other", Worktree: heldParent, Ref: "refs/heads/epic", ExpectedOID: otherOID}); err != nil {
					t.Fatal(err)
				}
				hold = heldIntegration(t, w, root, "other", "other", "other-epic", otherOID)
			}})
			hold()
			f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
			f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
			f.D.Executable = func() (string, error) { return binary, nil }
			ownerPaneStandin(t, &f)
			o := workflowTestStart(t, f.workflowFixture)
			o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o)))
			deliveryTestCandidate(t, f, "repository lock")
			if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'repository lock'\n"), 0600); err != nil {
				t.Fatal(err)
			}
			o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, deliveryTestReview(t, f, "repository-lock"), binary)...)
			heldHead := gitOutput(t, heldParent, "rev-parse", "HEAD")
			base := gitOutput(t, f.Parent, "rev-parse", "HEAD")

			cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "repository-lock"))
			id := deliveryCLIString(t, registered, "id")
			if !scenario.sameRepository {
				ready := deliveryCLIOK(t, binary, cwd, "check", id)
				if deliveryCLIString(t, ready, "state") != "ready" {
					t.Fatalf("independent repository was blocked: %s", ready["reasons"])
				}
				integrated, output, _ := deliveryCLIRun(binary, cwd, "execute", id)
				if string(integrated["local_integration"]) == "null" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != o.Delivery.Candidates[len(o.Delivery.Candidates)-1].OID {
					t.Fatalf("independent repository did not integrate the exact candidate:\n%s", output)
				}
			} else {
				reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
				// Delivery summarizes the plan reason parent_integration_already_active.
				blocked, output, _ := deliveryCLIRun(binary, cwd, "check", id)
				if deliveryCLIString(t, blocked, "state") != "blocked" || !strings.Contains(string(blocked["reasons"]), "local integration: conflict") {
					t.Fatalf("same repository and target was not serialized:\n%s", output)
				}
				if _, output, err := deliveryCLIRun(binary, cwd, "execute", id); err == nil {
					t.Fatalf("same repository and target integrated while another integration is in flight:\n%s", output)
				}
				if gitOutput(t, f.Parent, "rev-parse", "HEAD") != base || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
					t.Fatal("blocked integration changed the shared target")
				}
			}
			if gitOutput(t, heldParent, "rev-parse", "HEAD") != heldHead {
				t.Fatal("the held integration's target moved")
			}
			r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			for _, authority := range r.IntegrationAuthorities {
				if authority.TaskID != "held" {
					continue
				}
				for _, result := range r.IntegrationResults {
					if result.AuthorityID == authority.ID {
						t.Fatal("the held integration was resolved by another Task")
					}
				}
			}
		})
	}
}
