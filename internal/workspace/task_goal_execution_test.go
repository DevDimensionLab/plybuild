package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func goalFixture(t *testing.T, f workItemJourneyFixture, task TaskID, key string, agreement ...DeliveryAgreement) TaskGoalRef {
	t.Helper()
	in := contentFixtureSpec(t, f, key, nil)
	raw, e := os.ReadFile(in.File)
	if e != nil {
		t.Fatal(e)
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		t.Fatal(e)
	}
	m := contentFields(v)
	r := mustQueueRegistry(t, f)
	m["task_id"] = string(task)
	m["problem"] = contentRefValue(taskContentState(r, task).ProblemHead)
	m["expected_previous"] = contentRefValue(taskContentState(r, task).SpecHead)
	for _, k := range []string{"parts", "supporting", "removed_requirement_ids", "phases", "implementation_basis", "dependencies"} {
		delete(m, k)
	}
	first := contentString(contentFields(contentArray(m, "documents")[0]), "id")
	m["kind"] = "WorkspaceTaskSpecDraft@2"
	m["schema_version"] = int64(2)
	m["contract_kind"] = "goal"
	m["objective"] = "Make the chosen behavior observable without choosing the detailed implementation."
	m["requirements"] = []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "behavior", "acceptance": "The meaningful acceptance test observes the chosen behavior."})}
	m["design"] = []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"document_id": first, "section": nil})}
	m["constraints"] = []canonicaljson.Value{"Preserve existing behavior outside this goal."}
	m["executor"] = contentObject(map[string]canonicaljson.Value{"provider": "codex", "model": "fixture-model", "effort": "high"})
	if len(agreement) != 0 {
		m["delivery"] = contentRefValue(agreement[0])
	}
	raw, e = canonicaljson.Marshal(contentObject(m))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(in.File, raw, 0600); e != nil {
		t.Fatal(e)
	}
	in.TaskID = task
	out, e := RecordTaskSpec(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	return TaskGoalRef{SpecID: *out.OutcomeRef.SpecID, Spec: TaskRevisionRef{Revision: *out.OutcomeRef.Revision, ManifestSHA256: out.OutcomeRef.ManifestSHA256}}
}

func goalQueueFixture(t *testing.T, f workItemJourneyFixture, entries []TaskGoalQueueEntry) WorkspaceTaskQueueReadback {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r := mustQueueRegistry(t, f)
	q := foldQueue(r, queueID(root, QueueTarget{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}))
	draft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: fmt.Sprintf("goals/queue/%d", q.Revision), ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: q.Revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}, RegistryUpgrade: queueUpgrade(r)}
	raw, e := contentCanonical(draft)
	if e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(root, "goal-queue.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	out, e := SetTaskQueue(f.dependencies, file)
	if e != nil {
		t.Fatal(e)
	}
	return out
}

func goalHumanFixture() QueueHumanDecision {
	return QueueHumanDecision{ActorClaim: "fixture human", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "human_cli", Statement: "Execute this exact synthetic goal; this is not product QA."}
}

func goalPreparedFixture(t *testing.T, f workItemJourneyFixture) (TaskGoalExecutePreview, TaskGoalPreparedExecution) {
	t.Helper()
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	preview, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
	if e != nil {
		t.Fatal(e)
	}
	prepared, e := PrepareTaskGoalExecution(f.dependencies, preview.Input, preview.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatal(e)
	}
	return preview, prepared
}

func TestGoalQueueSurvivesEpicAdvanceAndBindsBaseAtExecute(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	q := goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	if q.Pending[0].State != "ready" || q.Pending[0].ParentOID != nil {
		t.Fatalf("goal was prematurely base-bound: %#v", q.Pending[0])
	}
	if len(mustQueueRegistry(t, f).TaskSolutionSelections) != 0 {
		t.Fatal("planner publication fabricated a human selection")
	}
	before := queueStateBytes(t, f)
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true}
	first, e := PreviewTaskGoalExecution(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("goal preview changed state")
	}
	if _, e = os.Stat(filepath.Dir(first.AcceptancePath)); !os.IsNotExist(e) {
		t.Fatal("preview created an acceptance directory")
	}
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "advance before goal execution")
	update, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	_, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *update.Confirmation})
	if e != nil {
		t.Fatal(e)
	}
	show, e := ShowTaskSpec(f.dependencies, TaskContentQuery{TaskID: "task", SpecID: goal.SpecID, Revision: goal.Spec.Revision})
	if e != nil {
		t.Fatal(e)
	}
	if contentString(contentFields(show.Value), "target_freshness") != "late_bound" {
		t.Fatalf("goal readback stale after unrelated advancement: %#v", show.Value)
	}
	second, e := PreviewTaskGoalExecution(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	if first.ParentOID == second.ParentOID || second.Goal.Goal != goal || second.AcceptancePath != first.AcceptancePath {
		t.Fatal("goal identity was changed or execution base was not refreshed")
	}
	if _, e = PrepareTaskGoalExecution(f.dependencies, first.Input, first.Confirmation, goalHumanFixture()); e == nil {
		t.Fatal("stale preview applied")
	}
	prepared, e := PrepareTaskGoalExecution(f.dependencies, second.Input, second.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatal(e)
	}
	if prepared.Preparation.Plan.ParentOID != second.ParentOID || prepared.Goal != goal {
		t.Fatal("prepared wrong goal or base")
	}
	if prepared.Executor.Provider == nil || *prepared.Executor.Provider != "codex" || *prepared.Executor.Model != "fixture-model" || *prepared.Executor.Effort != "high" {
		t.Fatal("executor assignment lost")
	}
	if _, e = os.Stat(prepared.AcceptancePath); !os.IsNotExist(e) {
		t.Fatal("preparation fabricated acceptance test implementation")
	}
	registry := mustQueueRegistry(t, f)
	if len(registry.TaskSpecRevisions) != 2 || len(registry.TaskSolutionSelections) != 1 {
		t.Fatal("goal was rewritten instead of materialized")
	}
	execution, e := readRegisteredTaskManifest(f.dependencies.TaskContent, filepath.Dir(f.wrapper), registry, prepared.Basis.Spec.ManifestSHA256)
	if e != nil {
		t.Fatal(e)
	}
	if contentString(contentFields(execution), "contract_kind") != "execution" || !contentEqual(contentFields(execution)["goal_origin"], contentRefValue(goal)) {
		t.Fatal("execution lost exact goal origin")
	}
}

func TestGoalExecutionRecoversEachPublishedContentStep(t *testing.T) {
	for _, failAt := range []int{1, 2, 3} {
		t.Run(string(rune('0'+failAt)), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			goal := goalFixture(t, f, "task", "goal")
			goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			p, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
			if e != nil {
				t.Fatal(e)
			}
			count := 0
			f.dependencies.TaskContent.fault = func(stage string) error {
				if stage == "postread" {
					count++
					if count == failAt {
						return errors.New("synthetic lost metadata response")
					}
				}
				return nil
			}
			if _, e = PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); e == nil {
				t.Fatal("expected injected failure")
			}
			f.dependencies.TaskContent.fault = nil
			resume, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
			if e != nil {
				t.Fatal(e)
			}
			if resume.Confirmation != p.Confirmation {
				t.Fatal("retry replaced the frozen request")
			}
			out, e := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture())
			if e != nil {
				t.Fatal(e)
			}
			before := queueStateBytes(t, f)
			again, e := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture())
			if e != nil {
				t.Fatal(e)
			}
			if before != queueStateBytes(t, f) || again.Preparation.Preparation.ID != out.Preparation.Preparation.ID {
				t.Fatal("retry duplicated a worktree or metadata")
			}
			r := mustQueueRegistry(t, f)
			if len(r.TaskSpecRevisions) != 2 || len(r.TaskSpecAssessments) != 1 || len(r.TaskSolutionSelections) != 1 || len(r.TaskPreparations) != 1 {
				t.Fatal("partial retry duplicated content")
			}
		})
	}
}

func TestGoalExecutionExplicitSpecAmbiguityAndPhysicalAcceptance(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	a := goalFixture(t, f, "task", "goal-a")
	_, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "Other", Description: "Other goal", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	b := goalFixture(t, f, "task-b", "goal-b")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &a}, {TaskID: "task-b", Goal: &b}})
	if _, e = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: a.SpecID}); e == nil || !strings.Contains(e.Error(), "ambiguous") {
		t.Fatalf("ambiguous Spec ID accepted: %v", e)
	}
	link := filepath.Join(filepath.Dir(f.wrapper), "acceptance-link")
	if e = os.Symlink(f.wrapper, link); e != nil {
		t.Fatal(e)
	}
	if _, e = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true, AcceptancePath: filepath.Join(link, "future", "acceptance.sh")}); e == nil {
		t.Fatal("symlink acceptance parent accepted")
	}
}

func TestDeliveryCandidateAcceptsOwnedDescendantButNotDirtyOrMovedParent(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	_, prepared := goalPreparedFixture(t, f)
	path := prepared.Preparation.Plan.WorktreePath
	runLocalGit(t, path, "commit", "--allow-empty", "-m", "implemented candidate")
	check := func() error {
		return WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
			_, e := s.ValidateDeliveryCandidateTarget(path, "refs/heads/"+prepared.Preparation.Plan.Branch, prepared.Preparation.Plan.Target.GitCommonDir, &prepared.Basis)
			return e
		})
	}
	if e := check(); e != nil {
		t.Fatal(e)
	}
	if e := WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
		_, e := s.ValidateTarget(path, "refs/heads/"+prepared.Preparation.Plan.Branch, prepared.Preparation.Plan.Target.GitCommonDir, &prepared.Basis)
		return e
	}); e == nil {
		t.Fatal("legacy start guard accepted advanced candidate")
	}
	file := filepath.Join(path, "dirty.txt")
	if e := os.WriteFile(file, []byte("dirty"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := check(); e == nil {
		t.Fatal("dirty candidate accepted")
	}
	if e := os.Remove(file); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "parent moved")
	if e := check(); e == nil {
		t.Fatal("moved parent accepted")
	}
}

func TestGoalExecutionRecoversAfterPlannerAddsIndependentGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	p, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
	if err != nil {
		t.Fatal(err)
	}
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "postread" {
			return errors.New("synthetic interruption after execution Spec publication")
		}
		return nil
	}
	if _, err = PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
		t.Fatal("expected interruption")
	}
	f.dependencies.TaskContent.fault = nil
	if _, err = CreateTask(f.dependencies, TaskCreateInput{TaskID: "later", Title: "Next independently planned feature", Description: "Planner keeps working", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	later := goalFixture(t, f, "later", "goal-later")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "later", Goal: &later}, {TaskID: "task", Goal: &goal}})
	resumed, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
	if err != nil {
		t.Fatal(err)
	}
	if resumed.Confirmation != p.Confirmation {
		t.Fatal("default recovery reranked the interrupted execution behind new planner work")
	}
	prepared, err := PrepareTaskGoalExecution(f.dependencies, resumed.Input, resumed.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	if prepared.Goal != goal || prepared.Preparation.Plan.TaskID != "task" {
		t.Fatal("recovery reranked the chosen goal")
	}
	queue, err := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(queue.Pending) != 1 || queue.Pending[0].TaskID != "later" || queue.Pending[0].Goal == nil || *queue.Pending[0].Goal != later || queue.Pending[0].State != "ready" {
		t.Fatalf("planner's independent goal was lost: %#v", queue.Pending)
	}
}

func TestGoalExecutionSamePublicInputRetriesWithoutDuplicateEffects(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	input := TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true}
	preview, err := PreviewTaskGoalExecution(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	first, err := PrepareTaskGoalExecution(f.dependencies, input, preview.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	before := queueStateBytes(t, f)
	second, err := PrepareTaskGoalExecution(f.dependencies, input, preview.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	if queueStateBytes(t, f) != before || first.Preparation.Preparation.ID != second.Preparation.Preparation.ID {
		t.Fatal("an identical public input created duplicate state")
	}
}

func TestGoalExecutionSelectionReadinessAndRejectedSelectorsAreReadOnly(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	initial := queueStateBytes(t, f)
	if _, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true}); err == nil || !strings.Contains(err.Error(), "no eligible") {
		t.Fatalf("empty queue did not explain the absence of an eligible goal: %v", err)
	}
	if initial != queueStateBytes(t, f) {
		t.Fatal("empty queue selection changed state")
	}
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	before := queueStateBytes(t, f)
	for _, input := range []TaskGoalExecuteInput{
		{Target: queueTargetFixture()},
		{Target: queueTargetFixture(), Next: true, SpecID: goal.SpecID},
		{Target: queueTargetFixture(), SpecID: "unknown"},
		{Target: QueueTargetInput{ProjectID: "ply"}, Next: true},
	} {
		if _, err := PreviewTaskGoalExecution(f.dependencies, input); err == nil {
			t.Fatalf("invalid or missing selector accepted: %#v", input)
		}
	}
	preview, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goal.SpecID})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Goal.Goal != goal || preview.Input.Next || preview.Input.SpecID != goal.SpecID {
		t.Fatal("explicit selector lost its exact goal")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("a read-only selector changed workspace or Git state")
	}
	if _, err = RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "changed-problem")); err != nil {
		t.Fatal(err)
	}
	before = queueStateBytes(t, f)
	if _, err = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true}); err == nil || !strings.Contains(err.Error(), "no eligible") {
		t.Fatalf("next silently selected a stale goal: %v", err)
	}
	if _, err = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: goal.SpecID}); err == nil || !strings.Contains(err.Error(), "task_goal_problem_changed") {
		t.Fatalf("explicit stale goal did not explain its problem: %v", err)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("rejected stale goal changed state")
	}
}

func TestLegacyPrepareExplainsThatGoalRequiresExecute(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	before := queueStateBytes(t, f)
	id := TaskID("task")
	preview, err := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "task", TaskID: &id}})
	if err != nil {
		t.Fatal(err)
	}
	if preview.State != "blocked" || len(preview.Reasons) != 1 || preview.Reasons[0].Code != "task_goal_requires_execute" || !strings.Contains(preview.Reasons[0].Message, "workflow execute") {
		t.Fatalf("legacy prepare presented an unmaterialized goal as ready: %#v", preview)
	}
	if before != queueStateBytes(t, f) || preview.Confirmation != nil || preview.Plan != nil {
		t.Fatal("legacy prepare should not publish a goal execution choice")
	}
}

func TestExecutionManifestRejectsForeignPreservedScope(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	_, prepared := goalPreparedFixture(t, f)
	registry := mustQueueRegistry(t, f)
	root := filepath.Dir(f.wrapper)
	manifest, err := readRegisteredTaskManifest(f.dependencies.TaskContent, root, registry, prepared.Basis.Spec.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"workspace", "repo", "parent_ref", "parent_worktree_id", "previous_spec"} {
		t.Run(field, func(t *testing.T) {
			m := contentFields(manifest)
			var preview TaskGoalExecutePreview
			if err := contentDecode(m["execution_request"], &preview); err != nil {
				t.Fatal(err)
			}
			switch field {
			case "workspace":
				preview.Workspace = filepath.Join(root, "foreign-workspace")
			case "repo":
				preview.Target.RepoID = "foreign-repo"
			case "parent_ref":
				preview.Target.ParentRef = "refs/heads/foreign-parent"
				preview.ObservedParent.Ref = preview.Target.ParentRef
			case "parent_worktree_id":
				preview.Target.ParentWorktreeID = "wt_00000000000000000000000000000000"
			case "previous_spec":
				preview.PreviousSpec = &prepared.Basis.Spec
			}
			preview.QueueID = queueID(preview.Workspace, preview.Target)
			m["execution_request"] = contentRefValue(preview)
			m["execution_request_sha256"] = goalExecutionDigest(preview)
			if err := validateExecutionGoalOrigin(f.dependencies, root, registry, "task", contentObject(m)); err == nil {
				t.Fatal("foreign preserved execution scope was accepted")
			}
		})
	}
}
