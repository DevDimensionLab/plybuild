package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestParallelForeignTargetAndFirstGoalRegistrationIsolation(t *testing.T) {
	f, goals := parallelPendingFixture(t)
	parent := filepath.Join(f.wrapper, "foreign-epic")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "foreign-epic", parent, f.oid)
	if _, e := AdoptEpic(f.dependencies, EpicAdoptInput{EpicID: "foreign", Title: "Foreign target", ProjectID: "ply", RepoID: "ply", Worktree: parent, Ref: "refs/heads/foreign-epic", ExpectedOID: f.oid}); e != nil {
		t.Fatal(e)
	}
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "foreign-c", Title: "C", Description: "Independent foreign target", ParentEpicID: "foreign", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	agreement := DeliveryAgreement{SchemaVersion: 1, Mode: DeliveryLocalEpic, ProjectID: "ply", RepoID: "ply", EpicID: "foreign", TargetRef: "refs/heads/foreign-epic", TargetWorktree: parent}
	cGoal := parallelGoalFixture(t, f, "foreign-c", "foreign-c-goal", agreement)
	target := QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "foreign"}
	draft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "foreign/queue", ProjectID: "ply", RepoID: "ply", EpicID: "foreign", ExpectedRevision: 0, Entries: []TaskGoalQueueEntry{{TaskID: "foreign-c", Goal: &cGoal}}, Recorder: TaskGoalQueueRecorder{ActorClaim: "synthetic planner", ControlSurface: "isolated fixture", RecordedAtUTC: "2026-10-23T12:00:00Z"}}
	raw, e := contentCanonical(draft)
	if e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(filepath.Dir(f.wrapper), "foreign-queue.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = SetTaskQueue(f.dependencies, file); e != nil {
		t.Fatal(e)
	}
	in := TaskGoalExecuteInput{Target: target, SpecID: cGoal.SpecID, Parallel: true}
	p, e := PreviewTaskGoalExecution(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	c, e := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatal(e)
	}
	r := mustQueueRegistry(t, f)
	foreignDigest := func(r WorkItemRegistry) string {
		task, _ := findTask(r, "foreign-c")
		op, _ := findOperation(r, "foreign-c")
		prep := findPreparation(r, c.Preparation.Preparation.ID)
		content := taskContentState(r, "foreign-c")
		content.RegistrySHA256 = nil // Shared container identity is not a C-owned record.
		return queueDigest([]any{task, op, prep, foldQueue(r, c.Preparation.Plan.QueueID), content})
	}
	before := foreignDigest(r)
	a := parallelPrepare(t, f, goals["task"])
	b := parallelPrepare(t, f, goals["task-b"])
	registration := firstGoalFixture(t, f)
	preview, e := RegisterFirstGoal(f.dependencies, registration, false, "")
	if e != nil {
		t.Fatal(e)
	}
	if len(preview.Queue.Active) != 2 {
		t.Fatal("registration preview lost A/B")
	}
	registered, e := RegisterFirstGoal(f.dependencies, registration, true, preview.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if registered.Queue.SchemaVersion != 2 || len(registered.Queue.Active) != 2 || len(registered.Queue.Pending) != 3 {
		t.Fatalf("registration nested queue unusable: %+v", registered.Queue)
	}
	state := queueStateBytes(t, f)
	if _, e = RegisterFirstGoal(f.dependencies, registration, true, preview.Confirmation); e != nil || queueStateBytes(t, f) != state {
		t.Fatal("firstGoal retry rewrote shared membership")
	}
	q := foldQueue(mustQueueRegistry(t, f), a.Preparation.Plan.QueueID)
	if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), a.Preparation.Preparation.ID, q.Revision, "explicit fixture release A only", "release"); e != nil {
		t.Fatal(e)
	}
	r = mustQueueRegistry(t, f)
	if foreignDigest(r) != before || !PreparationIsActive(r, b.Preparation.Preparation.ID) || PreparationIsActive(r, a.Preparation.Preparation.ID) {
		t.Fatal("registration/release changed C or B")
	}
	if got := runLocalGit(t, c.Preparation.Plan.WorktreePath, "rev-parse", "HEAD"); got != f.oid {
		t.Fatal("C source changed")
	}
}

func parallelPendingFixture(t *testing.T) (workItemJourneyFixture, map[TaskID]TaskGoalRef) {
	t.Helper()
	f := parallelFixture(t)
	goals := map[TaskID]TaskGoalRef{}
	entries := []TaskGoalQueueEntry{}
	for _, id := range []TaskID{"task", "task-b", "ply-consumer-agent-status", "ply-registered-cross-repo-plan"} {
		if id != "task" {
			if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: id, Title: string(id), Description: "isolated native fixture", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
				t.Fatal(err)
			}
		}
		g := parallelGoalFixture(t, f, id, "parallel-"+string(id))
		goals[id] = g
		entries = append(entries, TaskGoalQueueEntry{TaskID: id, Goal: &g})
	}
	goalQueueFixture(t, f, entries)
	return f, goals
}

func parallelPrepare(t *testing.T, f workItemJourneyFixture, goal TaskGoalRef) TaskGoalPreparedExecution {
	t.Helper()
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goal.SpecID, Parallel: true}
	p, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestParallelProtectedPendingAndPlannerCAS(t *testing.T) {
	f, goals := parallelPendingFixture(t)
	b := parallelPrepare(t, f, goals["task-b"])
	r := mustQueueRegistry(t, f)
	foreign := append([]TaskPreparation{}, r.TaskPreparations...)
	before := foldQueue(r, b.Preparation.Plan.QueueID).Pending
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goals["task"].SpecID, Parallel: true}
	p, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	// A planner reorders only the two protected entries. Their complete native
	// references must survive both stale refusal and refreshed named admission.
	goalQueueFixture(t, f, []TaskGoalQueueEntry{before[2], before[0], before[1]})
	state := queueStateBytes(t, f)
	if _, err = PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture()); err == nil || state != queueStateBytes(t, f) {
		t.Fatal("stale planner preview caused effects")
	}
	config := filepath.Join(b.Preparation.Plan.Target.GitCommonDir, "config")
	identity, err := os.Stat(config)
	if err != nil {
		t.Fatal(err)
	}
	configBytes, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	a := parallelPrepare(t, f, goals["task"])
	afterIdentity, err := os.Stat(config)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(identity, afterIdentity) || string(configBytes) != string(afterBytes) {
		t.Fatal("controlled native creation replaced shared Git config while B active")
	}
	t.Log("isolated native creation: shared config inode/bytes unchanged; live process-local trust remains unproven")
	r = mustQueueRegistry(t, f)
	q := foldQueue(r, p.QueueID)
	expected := []TaskGoalQueueEntry{before[2], before[1]}
	if !contentTypedEqual(q.Pending, expected) || !contentTypedEqual(*findPreparation(r, b.Preparation.Preparation.ID), foreign[0]) {
		t.Fatal("reserve lost B or protected pending reference/order")
	}
	// Explicit close removes A alone; serial and old readback consumers still
	// refuse the parallel queue rather than interpreting a nil singleton as idle.
	if _, err = ListTaskQueue(f.dependencies, queueTargetFixture(), false); err == nil {
		t.Fatal("schema1 silently hid multi-active")
	}
	if _, err = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true}); err == nil {
		t.Fatal("next adopted parallel participation")
	}
	if _, err = CloseTaskQueue(f.dependencies, queueTargetFixture(), a.Preparation.Preparation.ID, q.Revision, "simulated explicit named completion", "advance"); err != nil {
		t.Fatal(err)
	}
	r = mustQueueRegistry(t, f)
	q = foldQueue(r, p.QueueID)
	if len(q.Active) != 1 || q.Active[0] != b.Preparation.Preparation.ID || !contentTypedEqual(q.Pending, expected) || !contentTypedEqual(*findPreparation(r, b.Preparation.Preparation.ID), foreign[0]) {
		t.Fatal("close changed B or protected pending")
	}
}

func TestParallelConcurrentIdenticalPreparation(t *testing.T) {
	f, goals := parallelPendingFixture(t)
	b := parallelPrepare(t, f, goals["task-b"])
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goals["task"].SpecID, Parallel: true}
	p, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	ids := make([]string, 2)
	errs := make([]error, 2)
	for i := range ids {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out, e := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
			errs[i] = e
			if out.Preparation.Preparation != nil {
				ids[i] = out.Preparation.Preparation.ID
			}
		}(i)
	}
	wg.Wait()
	if errs[0] != nil || errs[1] != nil || ids[0] == "" || ids[0] != ids[1] {
		t.Fatalf("identical admission: %v %v %v", ids, errs[0], errs[1])
	}
	r := mustQueueRegistry(t, f)
	q := foldQueue(r, b.Preparation.Plan.QueueID)
	if len(r.TaskPreparations) != 2 || len(q.Active) != 2 || len(q.Pending) != 2 {
		t.Fatal("concurrent admission duplicated/lost participation")
	}
}

func TestParallelPublicationFaultsPreserveOtherMember(t *testing.T) {
	for _, point := range []string{"spec", "assessment", "selection", "reserve", "worktree", "ready"} {
		t.Run(point, func(t *testing.T) {
			f, goals := parallelPendingFixture(t)
			b := parallelPrepare(t, f, goals["task-b"])
			r := mustQueueRegistry(t, f)
			bRecord := *findPreparation(r, b.Preparation.Preparation.ID)
			pending := foldQueue(r, b.Preparation.Plan.QueueID).Pending
			in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goals["task"].SpecID, Parallel: true}
			p, err := PreviewTaskGoalExecution(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			original := f.dependencies.WorkGit
			contentCount := 0
			if point == "spec" || point == "assessment" || point == "selection" {
				want := map[string]int{"spec": 1, "assessment": 2, "selection": 3}[point]
				f.dependencies.TaskContent = &TaskContentStorage{fault: func(at string) error {
					if at == "before-registry-replace" {
						contentCount++
						if contentCount == want {
							return errors.New("simulated content interruption")
						}
					}
					return nil
				}}
			} else if point == "worktree" {
				f.dependencies.WorkGit = queueFaultGit{original, original.(TaskContentGit), "none"}
			} else {
				// Three content publications precede atomic reservation and ready.
				count := 0
				want := 4
				if point == "ready" {
					want = 5
				}
				f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(at string) error {
					if at == "temp-open" {
						count++
						if count == want {
							return errors.New("simulated registry interruption")
						}
					}
					return nil
				}}}
			}
			_, err = PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
			if err == nil {
				t.Fatal("injected interruption not observed")
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			f.dependencies.TaskContent = &TaskContentStorage{}
			f.dependencies.WorkGit = original
			r = mustQueueRegistry(t, f)
			if !contentTypedEqual(*findPreparation(r, b.Preparation.Preparation.ID), bRecord) {
				t.Fatal("interruption changed B")
			}
			out, err := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
			if err != nil {
				t.Fatal(err)
			}
			r = mustQueueRegistry(t, f)
			q := foldQueue(r, b.Preparation.Plan.QueueID)
			if len(r.TaskPreparations) != 2 || len(q.Active) != 2 || !contentTypedEqual(q.Pending, pending[1:]) || out.Preparation.State != "prepared" {
				t.Fatal("recovery lost protected references or duplicated native identity")
			}
			if _, err = filepath.EvalSymlinks(out.Preparation.Plan.WorktreePath); err != nil {
				t.Fatal(err)
			}
		})
	}
}
