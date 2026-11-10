package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func firstGoalFixture(t *testing.T, f workItemJourneyFixture) FirstGoalInput {
	t.Helper()
	doc := filepath.Join(filepath.Dir(f.wrapper), "supplied.md")
	if err := os.WriteFile(doc, []byte("# Supplied design\nPreserve these exact bytes.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return FirstGoalInput{TaskID: "fresh", SpecID: "fresh-v1", Title: "Supplied Goal", Objective: "Register the supplied behavior", Criteria: []string{"The specified behavior is observable."}, Documents: []GoalDocumentInput{{ID: "spec", Locator: doc}}, Target: queueTargetFixture()}
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
	for _, name := range []string{"criteria", "document-bytes", "document-identity", "target", "executor", "foreign-task", "inactive-task", "invalid-source", "invalid-criterion", "invalid-constraint", "invalid-executor", "wrong-ref"} {
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
			if len(out.Queue.Pending) != 2 || !reflect.DeepEqual(out.Queue.Pending[0], before.Pending[0]) || out.Problem.Revision != 2 {
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
	for _, name := range []string{"problem", "spec", "scope", "archived", "documents", "retired", "creating", "reconciliation", "invalid-receipt"} {
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
			case "documents":
				e = os.WriteFile(in.Documents[0].Locator, []byte("changed source bytes\n"), 0600)
			case "retired":
				nativeFirstGoalWorktree(t, f)
				e = f.dependencies.WorkItems.WithLock(filepath.Dir(f.wrapper), func(s WorkItemStoreSession) error {
					r, err := s.Snapshot()
					if err != nil {
						return err
					}
					r.Tasks[0].WorktreeState = WorkItemRetired
					return s.Publish(r)
				})
			case "creating":
				base := f.dependencies.WorkGit
				f.dependencies.WorkGit = noEffectGit{WorkItemGit: base}
				input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
				if _, err := CreateTaskWorktree(f.dependencies, input); err == nil {
					t.Fatal("creating fixture completed")
				}
				f.dependencies.WorkGit = base
			case "reconciliation":
				persistReconciliationFixture(t, &f)
			case "invalid-receipt":
				digest := nativeFirstGoalWorktree(t, f).Worktree.IntentDigest
				path := workItemsPath(filepath.Dir(f.wrapper))
				b, err := os.ReadFile(path)
				if err != nil || strings.Count(string(b), digest) != 1 {
					t.Fatalf("receipt digest is not unique: %v", err)
				}
				e = os.WriteFile(path, []byte(strings.Replace(string(b), digest, "sha256:"+strings.Repeat("e", 64), 1)), 0600)
			}
			if e != nil {
				t.Fatal(e)
			}
			before := queueStateBytes(t, f)
			_, e = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e == nil {
				t.Fatal("drift/used Task accepted")
			}
			if bound := map[string]bool{"retired": true, "creating": true, "reconciliation": true}; bound[name] && !strings.Contains(e.Error(), "first_goal_task_bound") {
				t.Fatalf("did not reach the bound-work refusal: %v", e)
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

func TestFirstGoalWithoutDocumentsRegistersWholeAssignment(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := firstGoalFixture(t, f)
	in.Documents = nil
	in.Constraints = []string{"Keep the human's limit."}
	before := queueStateBytes(t, f)
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Applied || queueStateBytes(t, f) != before || len(p.Sources) != 0 || strings.Contains(strings.Join(p.Effects, " "), "document") {
		t.Fatalf("preview misreports a document-free registration: %#v", p)
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	c := out.Contract
	if !out.Applied || c == nil || c.Title != in.Title || c.Objective != in.Objective || !reflect.DeepEqual(c.Constraints, in.Constraints) || len(c.Requirements) != 1 || c.Requirements[0].Acceptance != in.Criteria[0] || len(out.Queue.Pending) != 1 {
		t.Fatalf("assignment lost: %#v", out)
	}
	for _, input := range c.RequiredInputs {
		if input.ID == "task-goal-doc/goal/spec" || input.Role == "design" {
			t.Fatalf("invented a design input: %#v", input)
		}
	}
	registered := queueStateBytes(t, f)
	again, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil || again.Goal != out.Goal || queueStateBytes(t, f) != registered {
		t.Fatalf("retry changed publication: %v", err)
	}
	for name, change := range map[string]func(*FirstGoalInput){
		"objective":  func(x *FirstGoalInput) { x.Objective = "Changed" },
		"criterion":  func(x *FirstGoalInput) { x.Criteria = []string{"Changed"} },
		"constraint": func(x *FirstGoalInput) { x.Constraints = nil },
		"document":   func(x *FirstGoalInput) { x.Documents = firstGoalFixture(t, f).Documents },
	} {
		changed := in
		change(&changed)
		if _, err = RegisterFirstGoal(f.dependencies, changed, true, p.Confirmation); err == nil || queueStateBytes(t, f) != registered {
			t.Fatalf("%s drift accepted: %v", name, err)
		}
		preview, err := RegisterFirstGoal(f.dependencies, changed, false, "")
		if err == nil {
			if _, err = RegisterFirstGoal(f.dependencies, changed, true, preview.Confirmation); err == nil {
				t.Fatalf("%s drift replaced the immutable Goal", name)
			}
		}
		if queueStateBytes(t, f) != registered {
			t.Fatalf("%s drift wrote state", name)
		}
	}
}

func TestGoalContentStillRequiresRequirementAndBoundDocuments(t *testing.T) {
	doc := contentObject(map[string]canonicaljson.Value{"id": "spec"})
	req := []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "ac-001", "acceptance": "Observable"})}
	ref := func(id string) []canonicaljson.Value {
		return []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"document_id": id})}
	}
	for name, m := range map[string]map[string]canonicaljson.Value{
		"no requirement":  {"requirements": []canonicaljson.Value{}, "design": []canonicaljson.Value{}, "documents": []canonicaljson.Value{}},
		"unused document": {"requirements": req, "design": []canonicaljson.Value{}, "documents": []canonicaljson.Value{doc}},
		"missing target":  {"requirements": req, "design": ref("spec"), "documents": []canonicaljson.Value{}},
	} {
		if validateTaskGoalSemantics(m) == nil {
			t.Fatalf("%s accepted", name)
		}
	}
	if e := validateTaskGoalSemantics(map[string]canonicaljson.Value{"requirements": req, "design": []canonicaljson.Value{}, "documents": []canonicaljson.Value{}}); e != nil {
		t.Fatal(e)
	}
}

func TestExistingTaskGoalWithoutDocumentsKeepsItsProblem(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := existingFirstGoalFixture(t, f)
	in.Documents = nil
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Problem == nil || p.Problem.Revision != 2 || !reflect.DeepEqual(p.Effects, []string{"publish first Goal", "append exact Goal to queue"}) {
		t.Fatalf("existing Task preview: %#v", p)
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Applied || out.Contract == nil || out.Contract.Problem != *p.Problem {
		t.Fatalf("Goal replaced current Problem: %#v", out)
	}
}

func nativeFirstGoalWorktree(t *testing.T, f workItemJourneyFixture) TaskRecord {
	t.Helper()
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	ready, err := CreateTaskWorktree(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	return ready.Task
}

func firstGoalQueueAppends(r WorkItemRegistry, id TaskID) int {
	count := 0
	for _, ev := range r.TaskQueueEvents {
		if ev.Kind == "set" && contentString(contentFields(queueRequestValue(ev.Request)), "publication_key") == "first-goal/"+string(id)+"/queue" {
			count++
		}
	}
	return count
}

func firstGoalSpecCount(r WorkItemRegistry, id TaskID) int {
	count := 0
	for _, s := range r.TaskSpecRevisions {
		if s.TaskID == id {
			count++
		}
	}
	return count
}

func TestFirstGoalAfterNativeWorktreeKeepsBinding(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := existingFirstGoalFixture(t, f)
	nativeFirstGoalWorktree(t, f)
	r := mustQueueRegistry(t, f)
	beforeTask, _ := findTask(r, in.TaskID)
	problem := taskContentState(r, in.TaskID).ProblemHead
	before := queueStateBytes(t, f)
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Applied || queueStateBytes(t, f) != before {
		t.Fatal("preview wrote state")
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	r = mustQueueRegistry(t, f)
	afterTask, _ := findTask(r, in.TaskID)
	if !reflect.DeepEqual(beforeTask, afterTask) || afterTask.Worktree == nil || afterTask.WorktreeState != WorkItemReady || len(r.WorktreeOperations) != 0 {
		t.Fatalf("registration changed the native binding: %#v", afterTask)
	}
	if !reflect.DeepEqual(out.Problem, problem) || !reflect.DeepEqual(taskContentState(r, in.TaskID).ProblemHead, problem) || out.Contract == nil || out.Contract.Problem != *problem {
		t.Fatal("registration replaced the Problem head")
	}
	if firstGoalSpecCount(r, in.TaskID) != 1 || firstGoalQueueAppends(r, in.TaskID) != 1 || len(out.Queue.Pending) != 1 || out.Queue.Pending[0].TaskID != in.TaskID {
		t.Fatalf("expected exactly one Goal and one queue append: %#v", out.Queue)
	}
	registered := queueStateBytes(t, f)
	again, err := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
	if err != nil || again.Goal != out.Goal || queueStateBytes(t, f) != registered {
		t.Fatalf("retry changed state: %v", err)
	}
}

func TestFirstGoalAfterNativeWorktreeRecoversLostAcknowledgement(t *testing.T) {
	for _, n := range []int{1, 2} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			in := existingFirstGoalFixture(t, f)
			bound := nativeFirstGoalWorktree(t, f)
			p, e := RegisterFirstGoal(f.dependencies, in, false, "")
			if e != nil {
				t.Fatal(e)
			}
			store := &firstGoalLostAckStore{WorkItemStore: f.dependencies.WorkItems, failAt: n}
			f.dependencies.WorkItems = store
			if partial, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); e == nil || partial.Phase != "partial" {
				t.Fatalf("lost acknowledgement not surfaced: %#v %v", partial, e)
			}
			f.dependencies.WorkItems = store.WorkItemStore
			out, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			r := mustQueueRegistry(t, f)
			task, _ := findTask(r, in.TaskID)
			if !reflect.DeepEqual(task.Worktree, bound.Worktree) || task.WorktreeState != WorkItemReady || firstGoalSpecCount(r, in.TaskID) != 1 || firstGoalQueueAppends(r, in.TaskID) != 1 || len(out.Queue.Pending) != 1 {
				t.Fatal("recovery duplicated Goal/append or changed the binding")
			}
		})
	}
}

func TestFirstGoalAndNativeWorktreeCreateRunConcurrently(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := existingFirstGoalFixture(t, f)
	p, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	start := make(chan struct{})
	created, registered := make(chan error, 1), make(chan error, 1)
	go func() { <-start; _, e := CreateTaskWorktree(f.dependencies, input); created <- e }()
	go func() { <-start; _, e := RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); registered <- e }()
	close(start)
	var createErr, registerErr error
	for i := 0; i < 2; i++ {
		select {
		case createErr = <-created:
		case registerErr = <-registered:
		case <-time.After(10 * time.Second):
			t.Fatal("concurrent register/create timed out")
		}
	}
	if createErr != nil {
		t.Fatal(createErr)
	}
	r := mustQueueRegistry(t, f)
	if registerErr != nil {
		if !strings.Contains(registerErr.Error(), "first_goal_task_bound") || firstGoalSpecCount(r, in.TaskID) != 0 || firstGoalQueueAppends(r, in.TaskID) != 0 {
			t.Fatalf("register refused with effects or for another reason: %v", registerErr)
		}
		if _, err = RegisterFirstGoal(f.dependencies, in, true, p.Confirmation); err != nil {
			t.Fatal(err)
		}
		r = mustQueueRegistry(t, f)
	}
	task, _ := findTask(r, in.TaskID)
	if task.Worktree == nil || task.WorktreeState != WorkItemReady || len(r.WorktreeOperations) != 0 || firstGoalSpecCount(r, in.TaskID) != 1 || firstGoalQueueAppends(r, in.TaskID) != 1 {
		t.Fatalf("expected one bound Task with exactly one Goal: %#v", task)
	}
}
