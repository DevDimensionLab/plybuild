package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func firstGoalFixture(t *testing.T, f workItemJourneyFixture) FirstGoalInput {
	t.Helper()
	doc := filepath.Join(filepath.Dir(f.wrapper), "supplied.md")
	if err := os.WriteFile(doc, []byte("# Supplied design\nPreserve these exact bytes.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return FirstGoalInput{TaskID: "fresh", SpecID: "fresh-v1", Title: "Supplied Goal", Objective: "Register the supplied behavior", Criteria: []string{"The specified behavior is observable."}, Documents: []GoalDocumentInput{{ID: "spec", Locator: doc}}, Target: queueTargetFixture(), Delivery: GoalDeliveryInput{Mode: DeliveryLocalEpic, Acceptance: "human"}}
}

func existingFirstGoalFixture(t *testing.T, f workItemJourneyFixture) FirstGoalInput {
	t.Helper()
	if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "import/problem-2")); err != nil {
		t.Fatal(err)
	}
	in := firstGoalFixture(t, f)
	in.TaskID, in.SpecID = "task", "imported-first-goal"
	return in
}

func TestExistingFirstGoalPreservesProblemIdentityAndRetry(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := existingFirstGoalFixture(t, f)
	r := mustQueueRegistry(t, f)
	beforeTask, _ := findTask(r, in.TaskID)
	problem := taskContentState(r, in.TaskID).ProblemHead
	before := queueStateBytes(t, f)
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.Problem, problem) || problem.Revision != 2 || queueStateBytes(t, f) != before {
		t.Fatal("preview lost preserved Problem or wrote state")
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	r = mustQueueRegistry(t, f)
	afterTask, _ := findTask(r, in.TaskID)
	if !reflect.DeepEqual(beforeTask, afterTask) || !reflect.DeepEqual(problem, out.Problem) || len(r.TaskProblemRevisions) != 2 || out.Title != in.Title || out.Contract == nil || out.Readiness != "not_checked" || out.Run != nil {
		t.Fatal("registration replaced Task/Problem or lost supplied Goal")
	}
	before = queueStateBytes(t, f)
	f.dependencies.WorkClock = fixedWorkClock{value: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
	again, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil || again.Goal != out.Goal || queueStateBytes(t, f) != before {
		t.Fatalf("retry changed state: %v", err)
	}
}

func TestFirstGoalRegistrationPreviewApplyAndRetry(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := firstGoalFixture(t, f)
	before := queueStateBytes(t, f)
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Applied || p.Goal.Spec.ManifestSHA256 != "" || queueStateBytes(t, f) != before {
		t.Fatal("preview wrote state")
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Applied || out.Goal.Spec.Revision != 1 || out.Contract == nil || len(out.Contract.RequiredInputs) < 3 || len(out.Queue.Pending) != 1 {
		t.Fatalf("incomplete registration: %#v", out)
	}
	before = queueStateBytes(t, f)
	again, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil || again.Goal != out.Goal || queueStateBytes(t, f) != before {
		t.Fatalf("retry changed publication: %v", err)
	}
	in.Objective = "Changed real input"
	if _, err = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); err == nil {
		t.Fatal("changed product accepted")
	}
	r := mustQueueRegistry(t, f)
	if len(r.TaskSpecRevisions) != 1 || len(r.TaskSolutionSelections) != 0 || len(r.TaskPreparations) != 0 {
		t.Fatal("registration created execution or duplicate goal")
	}
}

// The real store performs each immutable publication/CAS before this adapter
// drops its acknowledgement. It changes no native store bytes or locks.
type firstGoalLostAckStore struct {
	WorkItemStore
	failAt, published int
}
type firstGoalLostAckSession struct {
	WorkItemStoreSession
	owner *firstGoalLostAckStore
}

func (s *firstGoalLostAckStore) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(inner WorkItemStoreSession) error { return fn(firstGoalLostAckSession{inner, s}) })
}
func (s firstGoalLostAckSession) Publish(r WorkItemRegistry) error {
	if e := s.WorkItemStoreSession.Publish(r); e != nil {
		return e
	}
	s.owner.published++
	if s.owner.published == s.owner.failAt {
		return fmt.Errorf("synthetic lost acknowledgement after native publication %d", s.owner.failAt)
	}
	return nil
}

func TestFirstGoalRegistrationRecoversNativePublicationBoundaries(t *testing.T) {
	for _, n := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in := firstGoalFixture(t, f)
			old := goalFixture(t, f, "task", "old-goal")
			goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &old}})
			p, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			store := &firstGoalLostAckStore{WorkItemStore: f.dependencies.WorkItems, failAt: n}
			f.dependencies.WorkItems = store
			partial, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e == nil || partial.Phase != "partial" {
				t.Fatalf("lost acknowledgement not surfaced: %#v %v", partial, e)
			}
			if n >= 2 && (partial.Contract == nil || partial.Goal.Spec.ManifestSHA256 == "") {
				t.Fatal("lost actual published Goal ref")
			}
			if n == 3 && len(partial.Queue.Pending) != 2 {
				t.Fatal("lost actual append readback")
			}
			f.dependencies.WorkItems = store.WorkItemStore
			// Legitimate independent queue change and a later observation clock are not
			// new registration inputs and must survive partial recovery.
			if n < 3 {
				goalQueueFixture(t, f, []TaskGoalQueueEntry{})
			}
			f.dependencies.WorkClock = fixedWorkClock{value: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			p2, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			if p2.Confirmation != p.Confirmation {
				t.Fatal("retry identity included own phases/time/unrelated queue drift")
			}
			out, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			r := mustQueueRegistry(t, f)
			count := 0
			for _, s := range r.TaskSpecRevisions {
				if s.TaskID == in.TaskID {
					count++
				}
			}
			if count != 1 || !out.Applied {
				t.Fatal("recovery duplicated Goal")
			}
			if n < 3 && (len(out.Queue.Pending) != 1 || out.Queue.Pending[0].TaskID != in.TaskID) {
				t.Fatal("recovery restored old unrelated queue")
			}
			pub := taskPublication(r, "first-goal/fresh/goal")
			manifest, e := readRegisteredTaskManifest(f.dependencies.TaskContent, filepath.Dir(f.wrapper), r, pub.OutcomeRef.ManifestSHA256)
			if e != nil {
				t.Fatal(e)
			}
			recorder := contentFields(contentFields(manifest)["recorder"])
			create := taskPublication(r, "task-create/fresh")
			if contentString(recorder, "recorded_at_utc") != create.RecordedAtUTC {
				t.Fatal("retry changed original input publication date")
			}
		})
	}
}

func TestFirstGoalRegistrationRejectsRealInputDriftWithoutDependentWrites(t *testing.T) {
	for _, name := range []string{"criteria", "policy", "document-bytes", "document-identity", "target", "executor", "foreign-task", "inactive-task", "invalid-source", "invalid-criterion", "invalid-constraint", "invalid-executor", "wrong-ref"} {
		t.Run(name, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in := firstGoalFixture(t, f)
			p, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			if name == "foreign-task" {
				_, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: in.TaskID, Title: in.Title, Description: "Foreign owner", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"})
				if e != nil {
					t.Fatal(e)
				}
				ensureFixtureTaskProblem(t, f, in.TaskID)
			} else if name == "inactive-task" {
				store := &firstGoalLostAckStore{WorkItemStore: f.dependencies.WorkItems, failAt: 1}
				f.dependencies.WorkItems = store
				if _, e = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); e == nil {
					t.Fatal("missing injected partial registration")
				}
				f.dependencies.WorkItems = store.WorkItemStore
				if _, e = SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: string(in.TaskID), State: LifecycleArchived, ActorClaim: "fixture"}); e != nil {
					t.Fatal(e)
				}
			} else if name != "invalid-source" && name != "invalid-criterion" && name != "invalid-constraint" && name != "invalid-executor" && name != "wrong-ref" {
				if _, e = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); e != nil {
					t.Fatal(e)
				}
			}
			expected := "first_goal_identity_conflict"
			switch name {
			case "criteria":
				in.Criteria[0] = "Different criterion"
			case "policy":
				in.Delivery.Acceptance = "automatic"
				in.Delivery.ResponsibleActor = "fixture"
				in.Delivery.TimeoutSeconds = 30
			case "document-bytes":
				if e = os.WriteFile(in.Documents[0].Locator, []byte("Different supplied bytes\n"), 0600); e != nil {
					t.Fatal(e)
				}
			case "document-identity":
				in.Documents[0].ID = "different"
			case "target":
				in.Target.EpicID = "missing"
				expected = "task_queue_context_required"
			case "executor":
				p := "claude"
				in.Executor.Provider = &p
			case "foreign-task":
				expected = "first_goal_confirmation_conflict"
			case "inactive-task":
				expected = "first_goal_inactive"
			case "invalid-source":
				if e = os.WriteFile(in.Documents[0].Locator, []byte{0xff}, 0600); e != nil {
					t.Fatal(e)
				}
				expected = "document must be"
			case "invalid-criterion":
				in.Criteria[0] = ""
				expected = "invalid observable criterion"
			case "invalid-constraint":
				in.Constraints = []string{""}
				expected = "constraints"
			case "invalid-executor":
				p := "invalid-provider"
				in.Executor.Provider = &p
				expected = "executor"
			case "wrong-ref":
				runLocalGit(t, filepath.Join(f.wrapper, "epic"), "checkout", "-b", "wrong-return")
				expected = "first_goal_target_conflict"
			}
			before := queueStateBytes(t, f)
			apply := name != "invalid-constraint" && name != "invalid-executor"
			if _, e = RegisterFirstGoal(f.dependencies, in, apply, p.Confirmation); e == nil || !strings.Contains(e.Error(), expected) {
				t.Fatalf("did not reach intended validation %s: %v", expected, e)
			}
			if queueStateBytes(t, f) != before {
				t.Fatal("rejected input changed native state")
			}
		})
	}
}

func TestExistingFirstGoalRecoversNativeBoundariesAndPreservesQueue(t *testing.T) {
	for _, n := range []int{1, 2} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in := existingFirstGoalFixture(t, f)
			if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "other", Title: "Other", Description: "Unrelated", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
				t.Fatal(e)
			}
			old := goalFixture(t, f, "other", "other-goal")
			before := goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "other", Goal: &old}})
			p, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			store := &firstGoalLostAckStore{WorkItemStore: f.dependencies.WorkItems, failAt: n}
			f.dependencies.WorkItems = store
			partial, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e == nil || partial.Phase != "partial" || partial.Contract == nil {
				t.Fatalf("missing native partial Goal: %#v %v", partial, e)
			}
			f.dependencies.WorkItems = store.WorkItemStore
			f.dependencies.WorkClock = fixedWorkClock{value: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			out, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if len(out.Queue.Pending) != 2 || !reflect.DeepEqual(out.Queue.Pending[0], before.Pending[0]) || !reflect.DeepEqual(out.Queue.Current, before.Current) || out.Problem.Revision != 2 {
				t.Fatal("recovery changed preserved queue or Problem")
			}
			r := mustQueueRegistry(t, f)
			count := 0
			for _, rev := range r.TaskSpecRevisions {
				if rev.TaskID == in.TaskID {
					count++
				}
			}
			if count != 1 || len(r.TaskProblemRevisions) != 3 {
				t.Fatal("duplicate native history")
			}
			changed := in
			changed.Objective = "Different supplied behavior"
			p2, e := RegisterFirstGoal(f.dependencies, changed, false, "")
			if e == nil {
				_, e = RegisterFirstGoal(f.dependencies, changed, true, p2.Confirmation)
			}
			if e == nil {
				t.Fatal("changed input obtained a new first registration")
			}
		})
	}
}

func TestExistingFirstGoalRejectsDriftAndUsedTasks(t *testing.T) {
	for _, name := range []string{"problem", "spec", "scope", "archived", "worktree", "documents"} {
		t.Run(name, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in := existingFirstGoalFixture(t, f)
			p, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			switch name {
			case "problem":
				_, e = RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "import/problem-3"))
			case "spec":
				goalFixture(t, f, "task", "foreign-goal")
			case "scope":
				in.Target.RepoID = "missing"
			case "archived":
				_, e = SetWorkItemLifecycle(f.dependencies, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: "task", State: LifecycleArchived, ActorClaim: "fixture"})
			case "worktree":
				input, _ := ParseTaskWorktreeCreateInput("task", "used-task", filepath.Join(f.wrapper, "used-task"), f.oid)
				_, e = CreateTaskWorktree(f.dependencies, input)
			case "documents":
				e = os.WriteFile(in.Documents[0].Locator, []byte("changed source bytes\n"), 0600)
			}
			if e != nil {
				t.Fatal(e)
			}
			before := queueStateBytes(t, f)
			if _, e = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); e == nil {
				t.Fatal("drift/used Task accepted")
			}
			if queueStateBytes(t, f) != before {
				t.Fatal("rejected registration wrote native state")
			}
		})
	}
}

func TestExistingFirstGoalAllowsOrdinaryDescriptionPrefix(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "imported", Title: "Imported backlog label", Description: "First Goal registration for the historical backlog item", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	ensureFixtureTaskProblem(t, f, "imported")
	in := firstGoalFixture(t, f)
	in.TaskID, in.SpecID = "imported", "imported-v1"
	p, e := RegisterFirstGoal(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	out, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if e != nil || !out.Applied || !reflect.DeepEqual(out.Problem, p.Problem) {
		t.Fatalf("ordinary descriptive prefix rejected: %v", e)
	}
}
