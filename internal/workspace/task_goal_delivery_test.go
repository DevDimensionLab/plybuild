package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type goalDeliveryFixture struct {
	workItemJourneyFixture
	root, taskPath, epicPath, resultOID string
	input                               TaskIntegrationInput
}

// automaticGoalDelivery is the delivery member a Goal registered with
// --acceptance automatic stored before DP01.
func automaticGoalDelivery(f workItemJourneyFixture) map[string]any {
	return map[string]any{"schema_version": 3, "mode": "local_epic_integration", "project_id": "ply", "repo_id": "ply", "epic_id": "epic", "target_ref": "refs/heads/epic", "target_worktree": filepath.Join(f.wrapper, "epic"), "acceptance": map[string]any{"schema_version": 1, "mode": "automatic", "responsible_actor": "ply-admin", "require_human_qa": false, "timeout_seconds": 30}}
}

// The fixture has only a registered Goal, its explicit Task worktree and an
// ordinary commit.
func newGoalDeliveryFixture(t *testing.T, delivery func(map[string]any)) goalDeliveryFixture {
	t.Helper()
	base := newWorkItemJourneyFixture(t)
	a := automaticGoalDelivery(base)
	if delivery != nil {
		delivery(a)
	}
	goalFixture(t, base, "task", "goal-delivery", a)
	taskPath := filepath.Join(base.wrapper, "task")
	input, _ := ParseTaskWorktreeCreateInput("task", "task", taskPath, base.oid)
	if _, err := CreateTaskWorktree(base.dependencies, input); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "result.txt"), []byte("result\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, taskPath, "add", "result.txt")
	runLocalGit(t, taskPath, "commit", "-m", "goal result")
	result := runLocalGit(t, taskPath, "rev-parse", "HEAD")
	base.dependencies.TaskLifecycleIDs = &sequenceTaskLifecycleIDs{}
	return goalDeliveryFixture{workItemJourneyFixture: base, root: filepath.Dir(base.wrapper), taskPath: taskPath, epicPath: filepath.Join(base.wrapper, "epic"), resultOID: result, input: TaskIntegrationInput{TaskID: "task", ExpectedResultOID: result, ExpectedParentOID: base.oid}}
}

func (f goalDeliveryFixture) apply(t *testing.T) WorkspaceTaskIntegrationReadback {
	t.Helper()
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	if checked.Readback.Classification != "ready" {
		t.Fatalf("check = %s %v", checked.Readback.Classification, checked.Readback.NextAction)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	return applied.Readback
}

func TestGoalDeliveryFastForwardsEpicWithoutProcessDocuments(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	before, _ := os.ReadFile(workItemsPath(f.root))
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(workItemsPath(f.root))
	if string(before) != string(after) || git.count() != 0 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.oid {
		t.Fatal("check had an effect")
	}
	if plan := integrationPlanFromReadback(t, checked.Readback); strings.Contains(plan, "task_result") || strings.Contains(plan, "human_qa") {
		t.Fatalf("integration presented process records: %s", plan)
	}
	applied := f.apply(t)
	if applied.Classification != "exact_effect" || git.count() != 1 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.resultOID {
		t.Fatalf("applied=%s merges=%d", applied.Classification, git.count())
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	if len(r.TaskResults) != 0 || len(r.HumanQARecords) != 0 || len(r.IntegrationAuthorities) != 1 {
		t.Fatalf("records: results=%d qa=%d authorities=%d", len(r.TaskResults), len(r.HumanQARecords), len(r.IntegrationAuthorities))
	}
	a := r.IntegrationAuthorities[0]
	if a.historical != nil || a.Plan.Kind != "WorkspaceTaskIntegrationPlan@6" {
		t.Fatalf("authority claims more than Git facts: %#v", a)
	}

	// Replay never merges again, whether by check, by apply or through the Task read.
	replay, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || replay.Readback.Classification != "exact_effect" {
		t.Fatalf("replay check = %v %v", replay.Readback.Classification, err)
	}
	in := f.input
	in.Confirmation = a.PlanSHA256
	again, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || again.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("replay apply = %s merges=%d err=%v", again.Readback.Classification, git.count(), err)
	}
	shown, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if shown.Integration == nil || shown.Integration.Classification != "exact_effect" || shown.Integration.ResultOID != f.resultOID {
		t.Fatalf("show = %s %s", shown.Integration.Classification, shown.Integration.NextAction.Kind)
	}
}

func TestGoalDeliveryTaskShowNamesDeliveryInsteadOfTaskResult(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	shown, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if shown.Integration == nil || shown.Integration.Classification != "ready" || shown.Integration.NextAction.Kind == "record_task_result" {
		t.Fatalf("show = %s %#v", shown.Integration.Classification, shown.Integration.NextAction)
	}
}

func TestGoalDeliveryReplayAfterLaterEpicWorkIsHistorical(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	f.apply(t)
	runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "later member")
	replay, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || replay.Readback.Classification != "exact_effect" || !replay.Readback.HistoricalCompletion || git.count() != 1 {
		t.Fatalf("replay=%s historical=%t merges=%d err=%v", replay.Readback.Classification, replay.Readback.HistoricalCompletion, git.count(), err)
	}
}

func TestGoalDeliveryPartialOrUnknownEffectBlocksNewAttempt(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	partial := &partialAfterCommandGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = partial
	out := f.apply(t)
	if out.Classification != "partial" {
		t.Fatalf("classification = %s", out.Classification)
	}
	f.dependencies.IntegrationGit = partial.TaskIntegrationGit
	again, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || again.Readback.Classification != "partial" || again.Readback.NextAction.Kind != "read_only_recovery_control" {
		t.Fatalf("again = %s %v", again.Readback.Classification, err)
	}
}

func TestGoalDeliveryResumesAfterCrashBetweenMergeAndResult(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	f.apply(t)
	err := f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		r.IntegrationResults = []IntegrationResult{}
		r.IntegrationAttempts[0].State, r.IntegrationAttempts[0].ResultID = "prepared", nil
		return session.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = mustQueueRegistry(t, f.workItemJourneyFixture).IntegrationAuthorities[0].PlanSHA256
	out, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || out.Readback.Classification != "exact_effect" || !strings.Contains(integrationPlanFromReadback(t, out.Readback), `"recovered":true`) || git.count() != 1 {
		t.Fatalf("resume=%s merges=%d err=%v", out.Readback.Classification, git.count(), err)
	}
}

func TestGoalDeliveryLetsTheEpicBaseAdvanceForTheNextTask(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	f.apply(t)
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	if unresolved := queueUnresolved(r, QueueTarget{ProjectID: "ply", RepoID: "ply", EpicID: "epic", GitCommonDir: r.Tasks[0].GitCommonDir, ParentRef: "refs/heads/epic"}); len(unresolved) != 0 {
		t.Fatalf("delivered Goal left unresolved records: %v", unresolved)
	}
	facts := RecordedTaskProgressFacts(r)
	if fact := facts["task"]; fact.IntegrationResult == nil || fact.IntegrationResult.Outcome != "exact_effect" || fact.IntegrationPending {
		t.Fatalf("progress does not show the delivery: %#v", fact)
	}
	preview, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || preview.Confirmation == nil {
		t.Fatalf("base preview = %#v err=%v", preview, err)
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Upgrade: r.RawSHA256, Apply: true, Confirmation: *preview.Confirmation}); err != nil {
		t.Fatal(err)
	}
	r = mustQueueRegistry(t, f.workItemJourneyFixture)
	epic, _ := findEpic(r, "epic")
	binding, _ := findEpicRepo(*epic, "ply")
	if base := currentEpicBase(r, *epic, *binding); base.OID != f.resultOID {
		t.Fatalf("Epic base = %s, want %s", base.OID, f.resultOID)
	}
	if _, err = CreateTask(f.dependencies, TaskCreateInput{TaskID: "next", Title: "Next", Description: "Next Task", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "next", Branch: "next", Path: filepath.Join(f.wrapper, "next"), ExpectedParentOID: f.resultOID}); err != nil {
		t.Fatal(err)
	}
}

func TestGoalDeliveryRefusesWhatGitOrTheGoalDoesNotAllow(t *testing.T) {
	expectReason := func(t *testing.T, f goalDeliveryFixture, in TaskIntegrationInput, state, reason string) {
		t.Helper()
		out, err := CheckTaskIntegration(f.dependencies, in)
		if err != nil {
			t.Fatal(err)
		}
		plan := integrationPlanFromReadback(t, out.Readback)
		if out.Readback.Classification != state || !strings.Contains(plan, reason) {
			t.Fatalf("classification=%s plan=%s", out.Readback.Classification, plan)
		}
		if runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.oid {
			t.Fatal("refusal moved the Epic")
		}
	}
	t.Run("dirty source", func(t *testing.T) {
		f := newGoalDeliveryFixture(t, nil)
		os.WriteFile(filepath.Join(f.taskPath, "untracked"), []byte("x"), 0o644)
		expectReason(t, f, f.input, "blocked", "source_not_clean")
	})
	t.Run("dirty parent", func(t *testing.T) {
		f := newGoalDeliveryFixture(t, nil)
		os.WriteFile(filepath.Join(f.epicPath, "untracked"), []byte("x"), 0o644)
		expectReason(t, f, f.input, "blocked", "parent_not_clean")
	})
	t.Run("source advanced after expected result", func(t *testing.T) {
		f := newGoalDeliveryFixture(t, nil)
		runLocalGit(t, f.taskPath, "commit", "--allow-empty", "-m", "later")
		expectReason(t, f, f.input, "conflict", "source_binding_mismatch")
	})
	t.Run("wrong expected parent", func(t *testing.T) {
		f := newGoalDeliveryFixture(t, nil)
		in := f.input
		in.ExpectedParentOID = f.resultOID
		expectReason(t, f, in, "conflict", "expected_oid_mismatch")
	})
	t.Run("parent moved before apply is refused by the confirmation", func(t *testing.T) {
		f := newGoalDeliveryFixture(t, nil)
		git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
		f.dependencies.IntegrationGit = git
		checked, err := CheckTaskIntegration(f.dependencies, f.input)
		if err != nil {
			t.Fatal(err)
		}
		runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "moved")
		in := f.input
		in.Confirmation = integrationDigest(t, checked.Readback)
		if _, err := ApplyTaskIntegration(f.dependencies, in); err == nil || git.count() != 0 {
			t.Fatalf("merges=%d err=%v", git.count(), err)
		}
		if r := mustQueueRegistry(t, f.workItemJourneyFixture); len(r.IntegrationAuthorities) != 0 {
			t.Fatal("refused apply published authority")
		}
	})
}

func TestGoalDeliveryNoEffectThenExactlyOneRetry(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	real := f.dependencies.IntegrationGit
	failing := &noEffectIntegrationGit{TaskIntegrationGit: real}
	f.dependencies.IntegrationGit = failing
	if first := f.apply(t); first.Classification != "no_effect" || failing.merges != 1 {
		t.Fatalf("first = %s", first.Classification)
	}
	f.dependencies.IntegrationGit = real
	// A new attempt names the no-effect result; without it the outcome stands.
	pending, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || pending.Readback.Classification != "no_effect" || !strings.Contains(strings.Join(pending.Readback.NextAction.Argv, " "), "--retry-after ires_") {
		t.Fatalf("without --retry-after = %s %v %v", pending.Readback.Classification, pending.Readback.NextAction, err)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	prior := r.IntegrationResults[0].ID
	f.input.RetryAfterResultID = &prior
	second := f.apply(t)
	if second.Classification != "exact_effect" || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.resultOID {
		t.Fatalf("second = %s", second.Classification)
	}
	replay, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || replay.Readback.Classification != "exact_effect" {
		t.Fatalf("replay = %s %v", replay.Readback.Classification, err)
	}
}

func TestGoalDeliveryDriftAfterConfirmationConflictsWithoutMergeAndIsNotPermanent(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	real := f.dependencies.IntegrationGit
	git := &driftAfterPlanGit{TaskIntegrationGit: real}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	git.drift, git.worktreeCalls = true, 0
	out, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || out.Readback.Classification != "conflict" || git.merged {
		t.Fatalf("drift = %s merged=%t err=%v", out.Readback.Classification, git.merged, err)
	}
	f.dependencies.IntegrationGit = real
	if done := f.apply(t); done.Classification != "exact_effect" {
		t.Fatalf("after resolved drift = %s", done.Classification)
	}
}

func TestGoalDeliveryRetryConflictKeepsPredecessorUsableAndMergesOnce(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	real := f.dependencies.IntegrationGit
	f.dependencies.IntegrationGit = &noEffectIntegrationGit{TaskIntegrationGit: real}
	if first := f.apply(t); first.Classification != "no_effect" {
		t.Fatal(first.Classification)
	}
	prior := mustQueueRegistry(t, f.workItemJourneyFixture).IntegrationResults[0].ID
	f.input.RetryAfterResultID = &prior
	git := &driftAfterPlanGit{TaskIntegrationGit: real}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	git.drift, git.worktreeCalls = true, 0
	out, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || out.Readback.Classification != "conflict" || git.merged {
		t.Fatalf("drift=%s merged=%t err=%v", out.Readback.Classification, git.merged, err)
	}
	counting := &countingIntegrationGit{TaskIntegrationGit: real}
	f.dependencies.IntegrationGit = counting
	if done := f.apply(t); done.Classification != "exact_effect" || counting.count() != 1 {
		t.Fatalf("after resolved = %s merges=%d", done.Classification, counting.count())
	}
	if replay, err := CheckTaskIntegration(f.dependencies, f.input); err != nil || replay.Readback.Classification != "exact_effect" {
		t.Fatalf("replay = %s %v", replay.Readback.Classification, err)
	}
}

type failAttemptIDs struct {
	sequenceTaskLifecycleIDs
	fail bool
}

func (s *failAttemptIDs) NewIntegrationAttemptID() (IntegrationAttemptID, error) {
	if s.fail {
		return "", errors.New("injected crash before attempt")
	}
	return s.sequenceTaskLifecycleIDs.NewIntegrationAttemptID()
}

func TestGoalDeliveryRecoversInterruptedAttemptFromReflogWithoutSecondMerge(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit, lost: true}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	if merged := runLocalGit(t, f.epicPath, "rev-parse", "HEAD"); merged != f.resultOID || git.count() != 1 {
		t.Fatalf("head=%s merges=%d", merged, git.count())
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	if len(r.IntegrationResults) != 1 || r.IntegrationResults[0].Outcome != "exact_effect" {
		t.Fatalf("results = %#v", r.IntegrationResults)
	}
}

func TestGoalDeliveryConcurrentIdenticalApplyMergesOnce(t *testing.T) {
	f := newGoalDeliveryFixture(t, nil)
	f.dependencies.TaskLifecycleIDs = cryptoTaskLifecycleIDSource{}
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	var wg sync.WaitGroup
	errs := make([]error, 4)
	for i := range errs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = ApplyTaskIntegration(f.dependencies, in)
		}(i)
	}
	wg.Wait()
	for _, e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	if git.count() != 1 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.resultOID {
		t.Fatalf("merges=%d", git.count())
	}
	if r := mustQueueRegistry(t, f.workItemJourneyFixture); len(r.IntegrationAttempts) != 1 {
		t.Fatalf("attempts=%d", len(r.IntegrationAttempts))
	}
}

func integrationPlanFromReadback(t *testing.T, r WorkspaceTaskIntegrationReadback) string {
	t.Helper()
	raw, err := MarshalTaskIntegrationReadback(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}
