package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func integrationDigest(t *testing.T, readback WorkspaceTaskIntegrationReadback) string {
	t.Helper()
	for _, member := range readback.Value {
		if member.Name != "plan" {
			continue
		}
		plan, ok := member.Value.(canonicaljson.Object)
		if !ok {
			t.Fatalf("plan is not an object: %#v", member.Value)
		}
		for _, field := range plan {
			if field.Name == "sha256" {
				digest, ok := field.Value.(string)
				if !ok {
					t.Fatalf("plan sha256 is not text: %#v", field.Value)
				}
				return digest
			}
		}
	}
	t.Fatal("readback has no plan digest")
	return ""
}

type driftAfterPlanGit struct {
	TaskIntegrationGit
	drift, merged bool
	worktreeCalls int
}

func (g *driftAfterPlanGit) ObserveIntegrationWorktree(path, ref string) (IntegrationWorktreeObservation, error) {
	observed, err := g.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
	if err != nil {
		return observed, err
	}
	g.worktreeCalls++
	if g.drift && g.worktreeCalls == 4 {
		observed.OID = strings.Repeat("f", len(observed.OID))
	}
	return observed, nil
}
func (g *driftAfterPlanGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	g.merged = true
	return g.TaskIntegrationGit.MergeFastForward(input)
}

type noEffectIntegrationGit struct {
	TaskIntegrationGit
	merges int
}

func (g *noEffectIntegrationGit) MergeFastForward(IntegrationMergeInput) GitCommandOutcome {
	g.merges++
	return GitCommandOutcome{Exit: 1, Err: errors.New("injected no effect")}
}

type countingIntegrationGit struct {
	TaskIntegrationGit
	mu     sync.Mutex
	merges int
	lost   bool
}

func (g *countingIntegrationGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	g.mu.Lock()
	g.merges++
	g.mu.Unlock()
	outcome := g.TaskIntegrationGit.MergeFastForward(input)
	if g.lost && outcome.Err == nil {
		outcome.Exit = 1
		outcome.Err = errors.New("injected lost response")
	}
	return outcome
}

func (g *countingIntegrationGit) count() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.merges
}

type partialAfterCommandGit struct {
	TaskIntegrationGit
	commanded bool
}

type unsupportedIntegrationGit struct{ TaskIntegrationGit }

func (g unsupportedIntegrationGit) ObserveIntegrationRepository(RepoRecord) (IntegrationRepositoryObservation, error) {
	return IntegrationRepositoryObservation{}, errors.New("injected unsupported repository")
}

type inventoryMismatchIntegrationGit struct{ TaskIntegrationGit }

func (g inventoryMismatchIntegrationGit) ObserveIntegrationRepository(repo RepoRecord) (IntegrationRepositoryObservation, error) {
	observed, err := g.TaskIntegrationGit.ObserveIntegrationRepository(repo)
	if err == nil && len(observed.Inventory) > 0 {
		observed.Inventory = append([]IntegrationInventoryEntry(nil), observed.Inventory...)
		observed.Inventory[0].OID = strings.Repeat("f", len(observed.Inventory[0].OID))
	}
	return observed, err
}

type repositoryFailureAfterCommandGit struct {
	TaskIntegrationGit
	commanded bool
}

type sourceDriftAfterCommandGit struct {
	TaskIntegrationGit
	commanded bool
	oid       string
	tree      string
}

func (git *sourceDriftAfterCommandGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	git.commanded = true
	return git.TaskIntegrationGit.MergeFastForward(input)
}

func (git *sourceDriftAfterCommandGit) ObserveIntegrationWorktree(path, ref string) (IntegrationWorktreeObservation, error) {
	observed, err := git.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
	if err == nil && git.commanded && strings.HasSuffix(ref, "/task") {
		observed.OID = git.oid
		observed.Tree = git.tree
	}
	return observed, err
}

type sourceFailureAfterCommandGit struct {
	TaskIntegrationGit
	commanded bool
}

func (git *sourceFailureAfterCommandGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	git.commanded = true
	return git.TaskIntegrationGit.MergeFastForward(input)
}

func (git *sourceFailureAfterCommandGit) ObserveIntegrationWorktree(path, ref string) (IntegrationWorktreeObservation, error) {
	if git.commanded && strings.HasSuffix(ref, "/task") {
		return IntegrationWorktreeObservation{}, errors.New("injected post-attempt source observation failure")
	}
	return git.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
}

func (git *repositoryFailureAfterCommandGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	git.commanded = true
	return git.TaskIntegrationGit.MergeFastForward(input)
}

func (git *repositoryFailureAfterCommandGit) ObserveIntegrationRepository(repo RepoRecord) (IntegrationRepositoryObservation, error) {
	if git.commanded {
		return IntegrationRepositoryObservation{}, errors.New("injected post-attempt repository observation failure")
	}
	return git.TaskIntegrationGit.ObserveIntegrationRepository(repo)
}

func integrationReadbackMember(t *testing.T, object canonicaljson.Object, name string) canonicaljson.Value {
	t.Helper()
	for _, member := range object {
		if member.Name == name {
			return member.Value
		}
	}
	t.Fatalf("readback member %q is missing", name)
	return nil
}

func integrationReadbackObject(t *testing.T, value canonicaljson.Value, name string) canonicaljson.Object {
	t.Helper()
	object, ok := value.(canonicaljson.Object)
	if !ok {
		t.Fatalf("readback member %q is not an object: %#v", name, value)
	}
	return object
}

func (g *partialAfterCommandGit) MergeFastForward(IntegrationMergeInput) GitCommandOutcome {
	g.commanded = true
	return GitCommandOutcome{Exit: 1, Err: errors.New("injected ambiguous command")}
}

func (g *partialAfterCommandGit) ObserveIntegrationWorktree(path, ref string) (IntegrationWorktreeObservation, error) {
	observed, err := g.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
	if err == nil && g.commanded && strings.HasSuffix(ref, "/epic") {
		observed.OID = strings.Repeat("f", len(observed.OID))
		observed.Tree = strings.Repeat("e", len(observed.Tree))
	}
	return observed, err
}

type sequenceTaskLifecycleIDs struct{ next int }

func (s *sequenceTaskLifecycleIDs) id(prefix string) string {
	s.next++
	return prefix + strings.Repeat("0", 31) + string("0123456789abcdef"[s.next])
}
func (s *sequenceTaskLifecycleIDs) NewTaskResultID() (TaskResultID, error) {
	return TaskResultID(s.id("trs_")), nil
}
func (s *sequenceTaskLifecycleIDs) NewHumanQARecordID() (HumanQARecordID, error) {
	return HumanQARecordID(s.id("hqa_")), nil
}
func (s *sequenceTaskLifecycleIDs) NewIntegrationAuthorityID() (IntegrationAuthorityID, error) {
	return IntegrationAuthorityID(s.id("iauth_")), nil
}
func (s *sequenceTaskLifecycleIDs) NewIntegrationIntentID() (IntegrationIntentID, error) {
	return IntegrationIntentID(s.id("iint_")), nil
}
func (s *sequenceTaskLifecycleIDs) NewIntegrationAttemptID() (IntegrationAttemptID, error) {
	return IntegrationAttemptID(s.id("iat_")), nil
}
func (s *sequenceTaskLifecycleIDs) NewIntegrationResultID() (IntegrationResultID, error) {
	return IntegrationResultID(s.id("ires_")), nil
}

type fixedWorkClock struct{ value time.Time }

func (c fixedWorkClock) Now() time.Time { return c.value }

type failFirstAttemptIDs struct {
	TaskLifecycleIDSource
	failed bool
}

func (source *failFirstAttemptIDs) NewIntegrationAttemptID() (IntegrationAttemptID, error) {
	if !source.failed {
		source.failed = true
		return "", errors.New("injected attempt ID failure")
	}
	return source.TaskLifecycleIDSource.NewIntegrationAttemptID()
}

type fixedProjectSnapshotLocker struct{ snapshot ProjectSnapshot }

func (locker fixedProjectSnapshotLocker) WithSnapshotLock(_ string, operation func(ProjectSnapshot) error) error {
	return operation(locker.snapshot)
}

type integrationJourneyFixture struct {
	workItemJourneyFixture
	root, taskPath, resultOID, resultTree string
	resultID                              TaskResultID
	qaID                                  HumanQARecordID
	evidence                              *fixedTaskEvidenceReader
}

func newIntegrationJourneyFixture(t *testing.T) integrationJourneyFixture {
	return newIntegrationJourneyFixtureVersion(t, false)
}
func newIntegrationJourneyFixtureVersion(t *testing.T, legacy bool) integrationJourneyFixture {
	t.Helper()
	base := newWorkItemJourneyFixtureVersion(t, legacy)
	taskPath := filepath.Join(base.wrapper, "task")
	input, _ := ParseTaskWorktreeCreateInput("task", "task", taskPath, base.oid)
	ready, err := CreateTaskWorktree(base.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	var basis *TaskSpecBasis
	if !legacy {
		b := fixtureSelectTaskSpec(t, base)
		basis = &b
	}
	if err := os.WriteFile(filepath.Join(taskPath, "result.txt"), []byte("result\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, taskPath, "add", "result.txt")
	runLocalGit(t, taskPath, "commit", "-m", "task result")
	resultOID := runLocalGit(t, taskPath, "rev-parse", "HEAD")
	resultTree := runLocalGit(t, taskPath, "rev-parse", "HEAD^{tree}")
	epicPath := filepath.Join(base.wrapper, "epic")
	runLocalGit(t, epicPath, "commit", "--allow-empty", "-m", "seed integration reflog")
	seedOID := runLocalGit(t, epicPath, "rev-parse", "HEAD")
	runLocalGit(t, epicPath, "update-ref", "-m", "restore integration base", "refs/heads/epic", base.oid, seedOID)
	resultID := TaskResultID("trs_11111111111111111111111111111111")
	qaID := HumanQARecordID("hqa_11111111111111111111111111111111")
	digest := "sha256:" + strings.Repeat("1", 64)
	root := filepath.Dir(base.wrapper)
	qaReport := filepath.Join(root, "qa-report.txt")
	qaReportBytes := []byte("QA passed\n")
	if err := os.WriteFile(qaReport, qaReportBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	err = base.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		if legacy {
			upgradeRegistryToV2(&registry)
		}
		registry.TaskResults = append(registry.TaskResults, TaskResultRecord{ID: resultID, PublicationKey: "task/result", DraftSHA256: digest, StoreTransition: "none", TaskID: "task", TaskWorktreeID: ready.Task.Worktree.ID, ProjectID: "ply", RepoID: "ply", GitCommonDir: ready.Task.GitCommonDir, SourceLocator: taskPath, SourceRef: "refs/heads/task", ResultOID: resultOID, ResultTree: resultTree, ActivityID: "act_11111111111111111111111111111111", RunID: "run_11111111111111111111111111111111", HandoffID: "hnd_11111111111111111111111111111111", HandoffLocator: filepath.Join(root, "handoff"), HandoffSHA256: digest, StartReceiptID: "rcp_11111111111111111111111111111111", StartReceiptLocator: filepath.Join(root, "start"), StartReceiptSHA256: digest, TerminalResultID: "res_11111111111111111111111111111111", TerminalResultLocator: filepath.Join(root, "terminal"), TerminalResultSHA256: digest, InspectionSHA256: digest, ReportedOutcome: "complete", TechnicalGate: "passed", VerifierResults: []TaskVerifierResultRecord{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, AcceptedDebt: []TaskAcceptedDebtRecord{}, Artifacts: []TaskArtifactRecord{}, Recorder: TaskRecorderRecord{ActorClaim: "test", ControlSurface: "test", RecordedAtUTC: "2026-09-28T12:00:00Z"}})
		if basis != nil {
			registry.TaskResultSpecBindings = append(registry.TaskResultSpecBindings, TaskResultSpecBinding{TaskID: "task", TaskResultID: resultID, Basis: *basis, HandoffSHA256: digest, StartReceiptSHA256: digest})
		}
		registry.HumanQARecords = append(registry.HumanQARecords, TaskHumanQARecord{ID: qaID, PublicationKey: "task/qa", DraftSHA256: digest, TaskID: "task", TaskResultID: resultID, ResultOID: resultOID, ResultTree: resultTree, Outcome: "pass", Actor: TaskHumanActorRecord{ActorClaim: "fixture human, not actual approval", StartSurface: "test", StartedAtUTC: "2026-09-28T12:00:00Z", CompletedAtUTC: "2026-09-28T12:01:00Z"}, Evidence: []TaskHumanQAEvidenceRecord{{ID: "report", Role: "report", Locator: qaReport, SHA256: digestTaskBytes(qaReportBytes), SizeBytes: int64(len(qaReportBytes))}}, Observation: "looks good", AcceptedResidualRisks: []TaskResidualRiskRecord{}})
		sortWorkRegistry(&registry)
		return session.Publish(registry)
	})
	if err != nil {
		t.Fatal(err)
	}
	ids := &sequenceTaskLifecycleIDs{}
	evidenceReader := &fixedTaskEvidenceReader{evidence: TaskHandoffEvidence{
		TaskSpecBasis: basis, TaskSpecValid: basis != nil, TaskRequirementsValid: basis != nil, AcceptedStartOutcome: "started",
		ActivityID: "act_11111111111111111111111111111111", RunID: "run_11111111111111111111111111111111", HandoffID: "hnd_11111111111111111111111111111111", HandoffLocator: filepath.Join(root, "handoff"), HandoffSHA256: digest,
		StartReceiptID: "rcp_11111111111111111111111111111111", StartReceiptLocator: filepath.Join(root, "start"), StartReceiptSHA256: digest,
		TerminalResultID: "res_11111111111111111111111111111111", TerminalResultLocator: filepath.Join(root, "terminal"), TerminalResultSHA256: digest, InspectionSHA256: digest,
		ReportedOutcome: "complete", TargetWorktree: taskPath, TargetRef: "refs/heads/task", ResultOID: resultOID, ResultTree: resultTree, GitCommonDir: ready.Task.GitCommonDir,
		SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true,
		VerifierResults: []TaskVerifierResultRecord{}, ExpectedVerifierIDs: []string{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, Artifacts: []TaskArtifactRecord{}, EvidenceGaps: []string{},
	}}
	base.dependencies.TaskLifecycleIDs = ids
	base.dependencies.WorkClock = fixedWorkClock{value: time.Date(2026, 9, 28, 12, 2, 0, 0, time.UTC)}
	base.dependencies.HandoffEvidence = evidenceReader
	return integrationJourneyFixture{workItemJourneyFixture: base, root: root, taskPath: taskPath, resultOID: resultOID, resultTree: resultTree, resultID: resultID, qaID: qaID, evidence: evidenceReader}
}

func TestTaskIntegrationCheckAndSingleConfirmedApply(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	before, err := os.ReadFile(workItemsPath(fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if checked.Readback.Classification != "ready" || checked.Readback.NextAction.Kind != "apply_confirmed_plan" || checked.Readback.GitChanged == nil || *checked.Readback.GitChanged {
		t.Fatalf("check = %#v", checked.Readback)
	}
	rechecked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	firstBytes, _ := MarshalTaskIntegrationReadback(checked.Readback)
	secondBytes, _ := MarshalTaskIntegrationReadback(rechecked.Readback)
	if integrationDigest(t, checked.Readback) != integrationDigest(t, rechecked.Readback) || string(firstBytes) != string(secondBytes) {
		t.Fatal("read-only check or canonical plan digest was not stable")
	}
	afterCheck, _ := os.ReadFile(workItemsPath(fixture.root))
	if string(before) != string(afterCheck) || runLocalGit(t, filepath.Join(fixture.wrapper, "epic"), "rev-parse", "HEAD") != fixture.oid {
		t.Fatal("check changed persisted or Git state")
	}
	confirmation := checked.Readback.NextAction.Argv[len(checked.Readback.NextAction.Argv)-1]
	input.Apply, input.Confirmation = true, confirmation
	applied, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Readback.Classification != "exact_effect" || applied.Readback.GitChanged == nil || !*applied.Readback.GitChanged || applied.Readback.RecoveryStatus != "complete" {
		t.Fatalf("apply = %#v", applied.Readback)
	}
	if got := runLocalGit(t, filepath.Join(fixture.wrapper, "epic"), "rev-parse", "HEAD"); got != fixture.resultOID {
		t.Fatalf("parent HEAD = %s", got)
	}
	storedBeforeRetry, _ := os.ReadFile(workItemsPath(fixture.root))
	retry, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	storedAfterRetry, _ := os.ReadFile(workItemsPath(fixture.root))
	if retry.Readback.Classification != "exact_effect" || string(storedBeforeRetry) != string(storedAfterRetry) {
		t.Fatal("identical apply retry was not a read-only replay")
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 1 || len(registry.IntegrationResults) != 1 {
		t.Fatalf("registry = %#v, %v", registry, err)
	}
}

func TestTaskIntegrationConfirmationMismatchStartsNoAuthorityOrGit(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid, Apply: true, Confirmation: "sha256:" + strings.Repeat("f", 64)}
	_, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationConfirmationMismatch) {
		t.Fatalf("error = %v", err)
	}
	registry, _ := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if len(registry.IntegrationAuthorities) != 0 || runLocalGit(t, filepath.Join(fixture.wrapper, "epic"), "rev-parse", "HEAD") != fixture.oid {
		t.Fatal("mismatched confirmation changed state")
	}
}

func TestTaskIntegrationDurableAuthorityLookupRequiresExactRedundantBinding(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err != nil {
		t.Fatal(err)
	}
	input.ExpectedParentOID = strings.Repeat("f", len(fixture.oid))
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err == nil {
		t.Fatal("durable authority was replayed for a different expected parent OID")
	}
}

func TestTaskIntegrationAuthorityRecoveryRejectsProjectBindingDrift(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	fixture.dependencies.TaskLifecycleIDs = &failFirstAttemptIDs{TaskLifecycleIDSource: fixture.dependencies.TaskLifecycleIDs}
	git := &countingIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err == nil {
		t.Fatal("injected attempt ID failure was not returned")
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 0 || git.count() != 0 {
		t.Fatalf("registry=%#v error=%v merges=%d", registry, err, git.count())
	}
	projects, repos, err := fixture.dependencies.Projects.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	repos[0].Locator = filepath.Join(fixture.root, "different-repository")
	fixture.dependencies.ProjectLocks = fixedProjectSnapshotLocker{snapshot: ProjectSnapshot{Projects: projects, Repos: repos}}

	recovered, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Readback.Classification != "conflict" || git.count() != 0 {
		t.Fatalf("recovered=%#v merges=%d", recovered.Readback, git.count())
	}
	registry, err = fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || len(registry.IntegrationAttempts) != 0 || len(registry.IntegrationResults) != 1 {
		t.Fatalf("registry=%#v error=%v", registry, err)
	}
}

func TestTaskIntegrationRechecksConfirmedStateBeforeAttempt(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	git := &driftAfterPlanGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply = true
	input.Confirmation = checked.Readback.NextAction.Argv[len(checked.Readback.NextAction.Argv)-1]
	git.drift, git.worktreeCalls = true, 0
	applied, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Readback.Classification != "conflict" || applied.Readback.GitChanged == nil || *applied.Readback.GitChanged || git.merged {
		t.Fatalf("drift result = %#v, merged=%t", applied.Readback, git.merged)
	}
	registry, _ := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 0 || len(registry.IntegrationResults) != 1 {
		t.Fatalf("registry = %#v", registry)
	}
}

func TestTaskIntegrationNoEffectNamesExactRetryPredecessor(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	git := &noEffectIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply = true
	input.Confirmation = checked.Readback.NextAction.Argv[len(checked.Readback.NextAction.Argv)-1]
	applied, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Readback.Classification != "no_effect" || applied.Readback.NextAction.Kind != "retry_after_no_effect" || git.merges != 1 {
		t.Fatalf("result = %#v, merges=%d", applied.Readback, git.merges)
	}
	registry, _ := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	want := string(registry.IntegrationResults[0].ID)
	argv := applied.Readback.NextAction.Argv
	found := false
	for i := 0; i+1 < len(argv); i++ {
		if argv[i] == "--retry-after" && argv[i+1] == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("retry argv = %#v, want result %s", argv, want)
	}
	input.Apply, input.Confirmation = false, ""
	rechecked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || rechecked.Readback.Classification != "no_effect" || rechecked.Readback.NextAction.Kind != "retry_after_no_effect" {
		t.Fatalf("plain recheck = %#v, %v", rechecked, err)
	}
}

func TestTaskIntegrationRetriesAfterNamedNoEffectWithOneNewAttempt(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	base := fixture.dependencies.IntegrationGit
	noEffect := &noEffectIntegrationGit{TaskIntegrationGit: base}
	fixture.dependencies.IntegrationGit = noEffect
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	first, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil || first.Readback.Classification != "no_effect" {
		t.Fatalf("first=%#v error=%v", first, err)
	}
	registry, _ := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	prior := registry.IntegrationResults[0].ID
	fixture.dependencies.IntegrationGit = base
	input.Apply, input.Confirmation, input.RetryAfterResultID = false, "", &prior
	recheck, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || recheck.Readback.Classification != "ready" {
		t.Fatalf("recheck=%#v error=%v", recheck, err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, recheck.Readback)
	retried, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil || retried.Readback.Classification != "exact_effect" || retried.Readback.GitChanged == nil || !*retried.Readback.GitChanged {
		t.Fatalf("retried=%#v error=%v", retried, err)
	}
	registry, _ = fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if len(registry.IntegrationAuthorities) != 2 || len(registry.IntegrationAttempts) != 2 || len(registry.IntegrationResults) != 2 || registry.IntegrationAuthorities[1].RetryAfterResultID == nil || *registry.IntegrationAuthorities[1].RetryAfterResultID != prior {
		t.Fatalf("retry registry = %#v", registry)
	}
}

func TestTaskIntegrationRecoversLostResponseFromReflogWithoutSecondMerge(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit, lost: true}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil || applied.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("applied=%#v error=%v merges=%d", applied, err, git.count())
	}
	replayed, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err != nil || replayed.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("replayed=%#v error=%v merges=%d", replayed, err, git.count())
	}
}

func TestTaskIntegrationRecoversDurableAuthorityAndAttemptPublicationBoundaries(t *testing.T) {
	for _, scenario := range []struct {
		name        string
		failPublish int
		wantMerges  int
	}{
		{name: "authority persisted but attempt publication failed", failPublish: 2, wantMerges: 0},
		{name: "attempt persisted but result publication failed", failPublish: 3, wantMerges: 1},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			fixture := newIntegrationJourneyFixture(t)
			git := &countingIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
			fixture.dependencies.IntegrationGit = git
			input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
			checked, err := CheckTaskIntegration(fixture.dependencies, input)
			if err != nil {
				t.Fatal(err)
			}
			input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
			publishes := 0
			fixture.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
				if operation == "temp-open" {
					publishes++
					if publishes == scenario.failPublish {
						return errInjected
					}
				}
				return nil
			}}}
			if _, err := ApplyTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
				t.Fatalf("first error = %v", err)
			}
			if git.count() != scenario.wantMerges {
				t.Fatalf("merges after interruption=%d want %d", git.count(), scenario.wantMerges)
			}
			fixture.dependencies.WorkItems = newSystemWorkItemStore()
			recovered, err := ApplyTaskIntegration(fixture.dependencies, input)
			if err != nil || recovered.Readback.Classification != "exact_effect" || git.count() != 1 {
				t.Fatalf("recovered=%#v error=%v merges=%d", recovered, err, git.count())
			}
			registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
			if err != nil || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 1 || len(registry.IntegrationResults) != 1 {
				t.Fatalf("registry=%#v error=%v", registry, err)
			}
		})
	}
}

func TestTaskIntegrationConcurrentIdenticalApplyProducesOneAttempt(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	fixture.dependencies.TaskLifecycleIDs = cryptoTaskLifecycleIDSource{}
	git := &countingIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	start := make(chan struct{})
	results := make(chan TaskIntegrationResult, 2)
	errors := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			result, err := ApplyTaskIntegration(fixture.dependencies, input)
			results <- result
			errors <- err
		}()
	}
	close(start)
	for i := 0; i < 2; i++ {
		if err := <-errors; err != nil {
			t.Fatal(err)
		}
		if result := <-results; result.Readback.Classification != "exact_effect" {
			t.Fatalf("result = %#v", result)
		}
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || git.count() != 1 || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 1 || len(registry.IntegrationResults) != 1 {
		t.Fatalf("registry=%#v error=%v merges=%d", registry, err, git.count())
	}
}

func TestTaskIntegrationClassifiesAlreadyIntegratedBlockedConflictAndPartial(t *testing.T) {
	t.Run("already integrated", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		git := &countingIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		fixture.dependencies.IntegrationGit = git
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, _ := CheckTaskIntegration(fixture.dependencies, input)
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		first, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		if first.Readback.ParentOID != fixture.resultOID || !strings.Contains(RenderTaskIntegrationText(first.Readback), "Parent: epic / refs/heads/epic@"+fixture.resultOID) {
			t.Fatalf("human readback = %q", RenderTaskIntegrationText(first.Readback))
		}
		input.Apply, input.Confirmation = false, ""
		checked, err = CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "already_integrated" {
			t.Fatalf("check=%#v error=%v", checked, err)
		}
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil || applied.Readback.Classification != "already_integrated" || git.count() != 1 {
			t.Fatalf("apply=%#v error=%v merges=%d", applied, err, git.count())
		}
	})

	t.Run("QA pass alone cannot overcome failed technical gate", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		if err := fixture.dependencies.WorkItems.WithLock(fixture.root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			registry.TaskResults[0].TechnicalGate = "failed"
			return session.Publish(registry)
		}); err != nil {
			t.Fatal(err)
		}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "blocked" {
			t.Fatalf("check=%#v error=%v", checked, err)
		}
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		registry, _ := fixture.dependencies.WorkItems.Snapshot(fixture.root)
		if err != nil || applied.Readback.Classification != "blocked" || len(registry.IntegrationAuthorities) != 0 {
			t.Fatalf("apply=%#v error=%v registry=%#v", applied, err, registry)
		}
	})

	t.Run("expected OID mismatch is conflict", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: strings.Repeat("f", len(fixture.resultOID)), ExpectedParentOID: fixture.oid}
		if _, err := CheckTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationConflict) {
			t.Fatalf("error=%v", err)
		}
	})

	t.Run("ambiguous command becomes partial", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.dependencies.IntegrationGit = &partialAfterCommandGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, _ := CheckTaskIntegration(fixture.dependencies, input)
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil || applied.Readback.Classification != "partial" || applied.Readback.GitChanged != nil || applied.Readback.RecoveryStatus != "reconciliation-required" || applied.Readback.ParentOID != strings.Repeat("f", len(fixture.oid)) {
			t.Fatalf("apply=%#v error=%v", applied, err)
		}
	})

	t.Run("post-attempt source drift is the fresh observed source", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		postOID := strings.Repeat("d", len(fixture.resultOID))
		postTree := strings.Repeat("e", len(fixture.resultTree))
		fixture.dependencies.IntegrationGit = &sourceDriftAfterCommandGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit, oid: postOID, tree: postTree}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		observedSource := integrationReadbackObject(t, integrationReadbackMember(t, applied.Readback.Value, "observed_source"), "observed_source")
		values := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "values"), "observed_source.values")
		freshness := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "freshness"), "observed_source.freshness")
		planSource := integrationReadbackObject(t, integrationReadbackMember(t, applied.Readback.Value, "source"), "source")
		if applied.Readback.Classification != "partial" || applied.Readback.GitChanged != nil || applied.Readback.RecoveryStatus != "reconciliation-required" {
			t.Fatalf("readback = %#v", applied.Readback)
		}
		if got := integrationReadbackMember(t, values, "oid"); got != postOID {
			t.Fatalf("observed source OID = %#v, want %s", got, postOID)
		}
		if got := integrationReadbackMember(t, values, "tree"); got != postTree {
			t.Fatalf("observed source tree = %#v, want %s", got, postTree)
		}
		for _, field := range []string{"worktree_locator", "ref", "oid", "tree", "git_common_dir", "object_format", "ref_format", "symbolic", "clean", "status_entries", "in_progress"} {
			if got := integrationReadbackMember(t, freshness, field); got != "fresh" {
				t.Fatalf("observed source freshness %s = %#v", field, got)
			}
		}
		if got := integrationReadbackMember(t, planSource, "oid"); got != fixture.resultOID {
			t.Fatalf("authority-plan source OID = %#v, want %s", got, fixture.resultOID)
		}
		if applied.Readback.ParentOID != fixture.resultOID || !strings.Contains(RenderTaskIntegrationText(applied.Readback), "Parent: epic / refs/heads/epic@"+fixture.resultOID) {
			t.Fatalf("human readback = %q", RenderTaskIntegrationText(applied.Readback))
		}
	})

	t.Run("post-attempt source observation failure is unknown without stale scalars", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.dependencies.IntegrationGit = &sourceFailureAfterCommandGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		observedSource := integrationReadbackObject(t, integrationReadbackMember(t, applied.Readback.Value, "observed_source"), "observed_source")
		values := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "values"), "observed_source.values")
		freshness := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "freshness"), "observed_source.freshness")
		if applied.Readback.Classification != "unknown" || applied.Readback.GitChanged != nil || applied.Readback.RecoveryStatus != "unknown" {
			t.Fatalf("readback = %#v", applied.Readback)
		}
		for _, field := range []string{"worktree_locator", "ref", "oid", "tree", "git_common_dir", "object_format", "ref_format", "symbolic", "clean"} {
			if got := integrationReadbackMember(t, values, field); got != nil {
				t.Fatalf("unknown observed source %s = %#v, want null", field, got)
			}
			if got := integrationReadbackMember(t, freshness, field); got != "unknown" {
				t.Fatalf("unknown observed source freshness %s = %#v", field, got)
			}
		}
		if applied.Readback.ParentOID != fixture.resultOID || !strings.Contains(RenderTaskIntegrationText(applied.Readback), "Parent: epic / refs/heads/epic@"+fixture.resultOID) {
			t.Fatalf("human readback = %q", RenderTaskIntegrationText(applied.Readback))
		}
	})

	t.Run("post-attempt repository failure keeps inventory freshness unknown", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.dependencies.IntegrationGit = &repositoryFailureAfterCommandGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := MarshalTaskIntegrationReadback(applied.Readback)
		if err != nil {
			t.Fatal(err)
		}
		if applied.Readback.Classification != "unknown" || !strings.Contains(string(encoded), `"observed_inventory":{"entries":[],"freshness":"unknown"}`) {
			t.Fatalf("readback = %s", encoded)
		}
		registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
		if err != nil || len(registry.IntegrationResults) != 1 || len(registry.IntegrationResults[0].AfterObservation.InventoryEntries) != 0 {
			t.Fatalf("registry=%#v error=%v", registry, err)
		}
	})
}

func TestTaskIntegrationPrecheckRejectsDirtyStaleMissingAndUnsupportedInputs(t *testing.T) {
	t.Run("dirty source", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		if err := os.WriteFile(filepath.Join(fixture.taskPath, "dirty.txt"), []byte("dirty\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "blocked" {
			t.Fatalf("checked=%#v error=%v", checked, err)
		}
	})

	t.Run("dirty parent", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		if err := os.WriteFile(filepath.Join(fixture.wrapper, "epic", "dirty.txt"), []byte("dirty\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "blocked" {
			t.Fatalf("checked=%#v error=%v", checked, err)
		}
	})

	t.Run("stale source", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		runLocalGit(t, fixture.taskPath, "commit", "--allow-empty", "-m", "source moved")
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "conflict" {
			t.Fatalf("checked=%#v error=%v", checked, err)
		}
	})

	t.Run("stale parent", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		runLocalGit(t, filepath.Join(fixture.wrapper, "epic"), "commit", "--allow-empty", "-m", "parent moved")
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "conflict" {
			t.Fatalf("checked=%#v error=%v", checked, err)
		}
	})

	t.Run("unsupported repository", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.dependencies.IntegrationGit = unsupportedIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		if _, err := CheckTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationUnsupportedRepository) {
			t.Fatalf("error=%v", err)
		}
	})

	t.Run("inventory mismatch", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.dependencies.IntegrationGit = inventoryMismatchIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		checked, err := CheckTaskIntegration(fixture.dependencies, input)
		if err != nil || checked.Readback.Classification != "conflict" {
			t.Fatalf("checked=%#v error=%v", checked, err)
		}
	})

	t.Run("handoff result without QA", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		missing := HumanQARecordID("hqa_ffffffffffffffffffffffffffffffff")
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: missing, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		if _, err := CheckTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationBlocked) {
			t.Fatalf("error=%v", err)
		}
	})

	t.Run("QA evidence drift", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		if err := os.WriteFile(filepath.Join(fixture.root, "qa-report.txt"), []byte("changed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		if _, err := CheckTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationBlocked) {
			t.Fatalf("error=%v", err)
		}
	})

	t.Run("handoff evidence drift", func(t *testing.T) {
		fixture := newIntegrationJourneyFixture(t)
		fixture.evidence.evidence.PolicyValid = false
		input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
		if _, err := CheckTaskIntegration(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorTaskIntegrationBlocked) {
			t.Fatalf("error=%v", err)
		}
	})
}

func TestTaskIntegrationAcceptsRevalidatedDebtControlArtifactRole(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	artifact := TaskArtifactRecord{ArtifactID: "debt-control", Role: "debt_control", Locator: filepath.Join(fixture.root, "debt-control.txt"), SHA256: "sha256:" + strings.Repeat("2", 64), SizeBytes: 7}
	if err := fixture.dependencies.WorkItems.WithLock(fixture.root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		registry.TaskResults[0].TechnicalGate = "good_enough_with_known_debt"
		registry.TaskResults[0].Artifacts = []TaskArtifactRecord{artifact}
		registry.TaskResults[0].AcceptedDebt = []TaskAcceptedDebtRecord{{ID: "coverage-gap", RiskClass: "coverage", Severity: "low", Summary: "A bounded coverage gap remains.", Control: "Use the managed control artifact.", EvidenceArtifactIDs: []string{artifact.ArtifactID}}}
		return session.Publish(registry)
	}); err != nil {
		t.Fatal(err)
	}
	terminalArtifact := artifact
	terminalArtifact.Role = "other"
	fixture.evidence.evidence.EvidenceCoverageValid = false
	fixture.evidence.evidence.EvidenceGaps = []string{"bounded coverage gap"}
	fixture.evidence.evidence.Artifacts = []TaskArtifactRecord{terminalArtifact}

	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if checked.Readback.Classification != "ready" {
		t.Fatalf("readback = %#v", checked.Readback)
	}
}

func TestIntegrationIntentDigestOmitsGeneratedIdentityFields(t *testing.T) {
	intent := IntegrationIntent{ID: "iint_11111111111111111111111111111111", IntentSHA256: "sha256:" + strings.Repeat("a", 64), AuthorityID: "iauth_11111111111111111111111111111111", State: "prepared", WorkspaceRoot: "/workspace", WorkspaceMarkerSHA256: "sha256:" + strings.Repeat("b", 64), MaxAttempts: 1}
	digest := integrationIntentDigest(intent)
	intent.ID, intent.IntentSHA256 = "", ""
	withEmptyIdentityBytes, err := canonicaljson.Marshal(integrationIntentCanonical(intent))
	if err != nil {
		t.Fatal(err)
	}
	if digest == digestTaskBytes(withEmptyIdentityBytes) {
		t.Fatal("intent digest included empty id and intent_sha256 fields")
	}
}
