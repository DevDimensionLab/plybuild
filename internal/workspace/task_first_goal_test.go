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

func TestFirstGoalRegistrationPreservesQueueAndDoesNotRequeuePreparedGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	_, _ = goalPreparedFixture(t, f)
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "other", Title: "Other", Description: "Keep exact Goal", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	old := goalFixture(t, f, "other", "other-goal")
	before := goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "other", Goal: &old}})
	in := firstGoalFixture(t, f)
	provider, model, effort := "codex", "fixture-model", "high"
	in.Executor = TaskExecutorAssignment{Provider: &provider, Model: &model, Effort: &effort}
	p, e := RegisterFirstGoal(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	out, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if len(out.Queue.Pending) != 2 || out.Queue.Pending[0].Goal == nil || *out.Queue.Pending[0].Goal != old || out.Queue.Current == nil || before.Current == nil || !reflect.DeepEqual(before.Current, out.Queue.Current) {
		t.Fatal("append changed unrelated current/pending")
	}
	// Independently register and prepare using native services in a fresh fixture;
	// this exercises historical retry without process/source-release observations.
	f = newWorkItemJourneyFixture(t)
	in = firstGoalFixture(t, f)
	in.Executor = TaskExecutorAssignment{Provider: &provider, Model: &model, Effort: &effort}
	p, e = RegisterFirstGoal(f.dependencies, in, false, "")
	if e != nil {
		t.Fatal(e)
	}
	out, e = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	next, e := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: in.SpecID})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = PrepareTaskGoalExecution(f.dependencies, next.Input, next.Confirmation, goalHumanFixture()); e != nil {
		t.Fatal(e)
	}
	beforeBytes := queueStateBytes(t, f)
	again, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if e != nil || !again.Applied || again.Goal != out.Goal || queueStateBytes(t, f) != beforeBytes {
		t.Fatalf("historical retry rewrote prepared work: %v", e)
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

func TestFirstGoalRegistrationConcurrentAppendsAndDirtyReturn(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	a := firstGoalFixture(t, f)
	b := a
	b.TaskID = "fresh-b"
	b.SpecID = "fresh-b-v1"
	pa, e := RegisterFirstGoal(f.dependencies, a, false, "")
	if e != nil {
		t.Fatal(e)
	}
	pb, e := RegisterFirstGoal(f.dependencies, b, false, "")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(f.wrapper, "epic", "dirty.txt"), []byte("unrelated change"), 0600); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 2)
	for _, item := range []struct {
		in    FirstGoalInput
		token string
	}{{a, pa.Confirmation}, {b, pb.Confirmation}} {
		item := item
		go func() { _, e := RegisterFirstGoal(f.dependencies, item.in, true, item.token); done <- e }()
	}
	for i := 0; i < 2; i++ {
		if e = <-done; e != nil {
			t.Fatal(e)
		}
	}
	q, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
	if e != nil {
		t.Fatal(e)
	}
	if len(q.Pending) != 2 || q.Pending[0].TaskID == q.Pending[1].TaskID {
		t.Fatal("concurrent legitimate append lost a Goal")
	}
	for _, row := range q.Pending {
		if row.Goal == nil {
			t.Fatal("append lost exact reference")
		}
	}
	// Registration remains successful. Existing execute preview owns the separate
	// dirty-return diagnosis and performs no execution.
	before := queueStateBytes(t, f)
	if _, e = PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: a.SpecID}); e == nil {
		t.Fatal("dirty return was reported ready")
	}
	if queueStateBytes(t, f) != before {
		t.Fatal("read-only readiness changed state")
	}
}
