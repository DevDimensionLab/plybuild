package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type integrationEvidenceByHandoff map[string]TaskHandoffEvidence

func (evidence integrationEvidenceByHandoff) ReadTaskEvidence(request TaskHandoffEvidenceRequest) (TaskHandoffEvidence, error) {
	found, ok := evidence[request.HandoffID]
	if !ok {
		return TaskHandoffEvidence{}, errors.New("fixture has no handoff evidence")
	}
	return found, nil
}

// lockFixture is an integration fixture whose evidence reader serves several
// Tasks, so independent integrations can share one workspace.
type lockFixture struct {
	integrationJourneyFixture
	handoffs integrationEvidenceByHandoff
}

func newLockFixture(t *testing.T) lockFixture {
	t.Helper()
	f := newIntegrationJourneyFixture(t)
	handoffs := integrationEvidenceByHandoff{f.evidence.evidence.HandoffID: f.evidence.evidence}
	f.dependencies.HandoffEvidence = handoffs
	return lockFixture{integrationJourneyFixture: f, handoffs: handoffs}
}

func (f lockFixture) input() TaskIntegrationInput {
	return TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
}

// leaveUnresolved persists the fixture Task's integration authority and then
// loses the attempt publication, the production state of an active integration.
// It returns the confirmed input that resumes that durable authority.
func (f lockFixture) leaveUnresolved(t *testing.T) TaskIntegrationInput {
	t.Helper()
	input := f.input()
	checked, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	publishes := 0
	d := f.dependencies
	d.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
		if operation == "temp-open" {
			publishes++
			if publishes == 2 {
				return errInjected
			}
		}
		return nil
	}}}
	if _, err := ApplyTaskIntegration(d, input); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
		t.Fatalf("interrupted apply error = %v", err)
	}
	registry, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationResults) != 0 {
		t.Fatalf("fixture integration is not unresolved: %#v %v", registry.IntegrationAuthorities, err)
	}
	return input
}

// addRepository registers an independent repository in its own Project with an
// Epic worktree on the given ref.
func (f lockFixture) addRepository(t *testing.T, id string, ref string) (EpicID, string, string) {
	t.Helper()
	wrapper := filepath.Join(f.root, id)
	repository := filepath.Join(wrapper, "main")
	if err := os.MkdirAll(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	repository = physicalPath(t, repository)
	runLocalGit(t, repository, "init", "-b", "main")
	runLocalGit(t, repository, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, repository, "config", "user.name", "Ply tests")
	oid := sharedTargetCommit(t, repository, "README.md", id+"\n")
	epicPath := filepath.Join(wrapper, "epic")
	runLocalGit(t, repository, "worktree", "add", "-b", strings.TrimPrefix(ref, "refs/heads/"), epicPath, oid)
	project := ProjectRecord{ID: ProjectID(id), Name: id, Wrapper: physicalPath(t, wrapper), RepoIDs: []RepoID{RepoID(id)}}
	if _, _, _, err := f.dependencies.Projects.(*systemProjectStore).Add(f.root, project, []RepoRecord{{ID: RepoID(id), Locator: repository, GitCommonDir: physicalPath(t, filepath.Join(repository, ".git"))}}); err != nil {
		t.Fatal(err)
	}
	epic := EpicID(id + "-epic")
	adopt, err := ParseEpicAdoptInput(string(epic), "Epic "+id, id, id, physicalPath(t, epicPath), ref, oid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
		t.Fatal(err)
	}
	return epic, physicalPath(t, wrapper), oid
}

// addTask registers a Task with a selected Spec, a committed result and a
// passing QA record, returning the input for its own integration.
func (f lockFixture) addTask(t *testing.T, id TaskID, epic EpicID, project ProjectID, repo RepoID, wrapper, parentOID string, n byte) TaskIntegrationInput {
	t.Helper()
	create, err := ParseTaskCreateInput(string(id), "Task "+string(id), "Independent integration", string(epic), string(project), string(repo))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateTask(f.dependencies, create); err != nil {
		t.Fatal(err)
	}
	worktree, _ := ParseTaskWorktreeCreateInput(string(id), string(id), filepath.Join(wrapper, string(id)), parentOID)
	ready, err := CreateTaskWorktree(f.dependencies, worktree)
	if err != nil {
		t.Fatal(err)
	}
	taskPath := ready.Task.Worktree.Locator
	basis := fixtureSelectTaskSpecFor(t, f.workItemJourneyFixture, id, "fixture/"+string(id))
	resultOID := sharedTargetCommit(t, taskPath, string(id)+".txt", "result\n")
	resultTree := runLocalGit(t, taskPath, "rev-parse", "HEAD^{tree}")
	hex := strings.Repeat(string(n), 32)
	resultID, qaID := TaskResultID("trs_"+hex), HumanQARecordID("hqa_"+hex)
	digest := "sha256:" + strings.Repeat(string(n), 64)
	evidence := f.evidence.evidence
	evidence.TaskSpecBasis = &basis
	evidence.ActivityID, evidence.RunID, evidence.HandoffID = "act_"+hex, "run_"+hex, "hnd_"+hex
	evidence.StartReceiptID, evidence.TerminalResultID = "rcp_"+hex, "res_"+hex
	evidence.HandoffLocator, evidence.StartReceiptLocator, evidence.TerminalResultLocator = filepath.Join(f.root, "handoff-"+string(id)), filepath.Join(f.root, "start-"+string(id)), filepath.Join(f.root, "terminal-"+string(id))
	evidence.HandoffSHA256, evidence.StartReceiptSHA256, evidence.TerminalResultSHA256, evidence.InspectionSHA256 = digest, digest, digest, digest
	evidence.TargetWorktree, evidence.TargetRef, evidence.ResultOID, evidence.ResultTree, evidence.GitCommonDir = taskPath, ready.Task.Worktree.Ref, resultOID, resultTree, ready.Task.GitCommonDir
	f.handoffs[evidence.HandoffID] = evidence
	report := filepath.Join(f.root, "qa-"+string(id)+".txt")
	reportBytes := []byte("QA passed\n")
	if err := os.WriteFile(report, reportBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	err = f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		registry.TaskResults = append(registry.TaskResults, TaskResultRecord{ID: resultID, PublicationKey: "task/result/" + string(id), DraftSHA256: digest, StoreTransition: "none", TaskID: id, TaskWorktreeID: ready.Task.Worktree.ID, ProjectID: project, RepoID: repo, GitCommonDir: ready.Task.GitCommonDir, SourceLocator: taskPath, SourceRef: ready.Task.Worktree.Ref, ResultOID: resultOID, ResultTree: resultTree, ActivityID: evidence.ActivityID, RunID: evidence.RunID, HandoffID: evidence.HandoffID, HandoffLocator: evidence.HandoffLocator, HandoffSHA256: digest, StartReceiptID: evidence.StartReceiptID, StartReceiptLocator: evidence.StartReceiptLocator, StartReceiptSHA256: digest, TerminalResultID: evidence.TerminalResultID, TerminalResultLocator: evidence.TerminalResultLocator, TerminalResultSHA256: digest, InspectionSHA256: digest, ReportedOutcome: "complete", TechnicalGate: "passed", VerifierResults: []TaskVerifierResultRecord{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, AcceptedDebt: []TaskAcceptedDebtRecord{}, Artifacts: []TaskArtifactRecord{}, Recorder: TaskRecorderRecord{ActorClaim: "test", ControlSurface: "test", RecordedAtUTC: "2026-09-28T12:00:00Z"}})
		registry.TaskResultSpecBindings = append(registry.TaskResultSpecBindings, TaskResultSpecBinding{TaskID: id, TaskResultID: resultID, Basis: basis, HandoffSHA256: digest, StartReceiptSHA256: digest})
		registry.HumanQARecords = append(registry.HumanQARecords, TaskHumanQARecord{ID: qaID, PublicationKey: "task/qa/" + string(id), DraftSHA256: digest, TaskID: id, TaskResultID: resultID, ResultOID: resultOID, ResultTree: resultTree, Outcome: "pass", Actor: TaskHumanActorRecord{ActorClaim: "fixture human, not actual approval", StartSurface: "test", StartedAtUTC: "2026-09-28T12:00:00Z", CompletedAtUTC: "2026-09-28T12:01:00Z"}, Evidence: []TaskHumanQAEvidenceRecord{{ID: "report", Role: "report", Locator: report, SHA256: digestTaskBytes(reportBytes), SizeBytes: int64(len(reportBytes))}}, Observation: "looks good", AcceptedResidualRisks: []TaskResidualRiskRecord{}})
		sortWorkRegistry(&registry)
		return session.Publish(registry)
	})
	if err != nil {
		t.Fatal(err)
	}
	return TaskIntegrationInput{TaskID: id, TaskResultID: resultID, HumanQARecordID: qaID, ExpectedResultOID: resultOID, ExpectedParentOID: parentOID}
}

func checkedIntegration(t *testing.T, d Dependencies, input TaskIntegrationInput) (WorkspaceTaskIntegrationReadback, string) {
	t.Helper()
	checked, err := CheckTaskIntegration(d, input)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := MarshalTaskIntegrationReadback(checked.Readback)
	if err != nil {
		t.Fatal(err)
	}
	return checked.Readback, string(encoded)
}

// IL01 ac-01: an unresolved integration of repository A to refs/heads/epic must
// not block repository B from integrating to its own refs/heads/epic.
func TestTaskIntegrationUnresolvedRepositoryDoesNotBlockOtherRepositorySameRef(t *testing.T) {
	f := newLockFixture(t)
	original := f.leaveUnresolved(t)
	epic, wrapper, parentOID := f.addRepository(t, "other", "refs/heads/epic")
	input := f.addTask(t, "other-task", epic, "other", "other", wrapper, parentOID, '2')

	readback, encoded := checkedIntegration(t, f.dependencies, input)
	if readback.Classification != "ready" || strings.Contains(encoded, "parent_integration_already_active") {
		t.Fatalf("independent repository was blocked: %s", encoded)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, readback)
	applied, err := ApplyTaskIntegration(f.dependencies, input)
	if err != nil || applied.Readback.Classification != "exact_effect" {
		t.Fatalf("independent apply = %#v %v", applied.Readback, err)
	}
	if got := runLocalGit(t, filepath.Join(wrapper, "epic"), "rev-parse", "HEAD"); got != input.ExpectedResultOID {
		t.Fatalf("independent parent HEAD = %s", got)
	}
	if got := runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD"); got != f.oid {
		t.Fatalf("unresolved repository parent moved: %s", got)
	}

	// The unresolved integration keeps its own authority and completes once.
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	recovered, err := ApplyTaskIntegration(f.dependencies, original)
	if err != nil || recovered.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("recovered = %#v %v merges=%d", recovered.Readback, err, git.count())
	}
	registry, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(registry.IntegrationAuthorities) != 2 || len(registry.IntegrationResults) != 2 {
		t.Fatalf("registry authorities=%d results=%d %v", len(registry.IntegrationAuthorities), len(registry.IntegrationResults), err)
	}
}

// IL01 ac-02: the same repository and target stays serialized, while another
// target in the same repository is independent.
func TestTaskIntegrationSameRepositoryLockFollowsTargetRef(t *testing.T) {
	t.Run("same repository and ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		input := f.addTask(t, "sibling", "epic", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification == "ready" || !strings.Contains(encoded, "parent_integration_already_active") {
			t.Fatalf("same repository and ref was not serialized: %s", encoded)
		}
	})
	t.Run("worktree alias of the same repository and ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		alias := filepath.Join(f.wrapper, "alias")
		runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "move", filepath.Join(f.wrapper, "epic"), alias)
		adopt, err := ParseEpicAdoptInput("alias", "Alias Epic", "ply", "ply", physicalPath(t, alias), "refs/heads/epic", f.oid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
			t.Fatal(err)
		}
		input := f.addTask(t, "alias-task", "alias", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification == "ready" || !strings.Contains(encoded, "parent_integration_already_active") {
			t.Fatalf("worktree alias of the same repository and ref was not serialized: %s", encoded)
		}
	})
	t.Run("same repository other ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		path := filepath.Join(f.wrapper, "second")
		runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "second", path, f.oid)
		adopt, err := ParseEpicAdoptInput("second", "Second Epic", "ply", "ply", physicalPath(t, path), "refs/heads/second", f.oid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
			t.Fatal(err)
		}
		input := f.addTask(t, "second-task", "second", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification != "ready" {
			t.Fatalf("other ref in the same repository was blocked: %s", encoded)
		}
	})
}
