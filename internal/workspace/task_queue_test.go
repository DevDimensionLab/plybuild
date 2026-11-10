package workspace

import (
	"bytes"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func queueFixtureDraft(t *testing.T, f workItemJourneyFixture, id TaskID, key string, prev *TaskRevisionRef) TaskContentInput {
	t.Helper()
	ensureFixtureTaskProblem(t, f, id)
	in := contentFixtureSpec(t, f, key, prev)
	raw, e := os.ReadFile(in.File)
	if e != nil {
		t.Fatal(e)
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		t.Fatal(e)
	}
	m := contentFields(v)
	r, e := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if e != nil {
		t.Fatal(e)
	}
	task, _ := findTask(r, id)
	ep, b := queueTargetBinding(r, QueueTarget{EpicID: task.ParentEpicID, RepoID: task.RepoID})
	base := currentEpicBase(r, ep, b)
	m["task_id"] = string(id)
	m["problem"] = contentRefValue(taskContentState(r, id).ProblemHead)
	basis := contentFields(m["implementation_basis"])
	basis["parent_oid"] = base.OID
	basis["parent_tree"] = base.Tree
	basis["start_oid"] = base.OID
	basis["start_tree"] = base.Tree
	m["implementation_basis"] = contentObject(basis)
	if prev != nil {
		docs := contentArray(m, "documents")
		for _, doc := range docs {
			dm := contentFields(doc)
			dm["source"] = contentObject(map[string]canonicaljson.Value{"kind": "snapshot", "manifest_sha256": prev.ManifestSHA256, "document_id": contentString(dm, "id")})
			docs[0] = contentObject(dm)
		}
		m["documents"] = docs
	}
	raw, e = canonicaljson.Marshal(contentObject(m))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(in.File, raw, 0600); e != nil {
		t.Fatal(e)
	}
	in.TaskID = id
	return in
}
func queueSelect(t *testing.T, f workItemJourneyFixture, id TaskID, key string, prev *TaskRevisionRef) TaskDecisionRef {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	record, e := RecordTaskSpec(f.dependencies, queueFixtureDraft(t, f, id, key+"-spec", prev))
	if e != nil {
		t.Fatal(e)
	}
	ref := TaskRevisionRef{*record.OutcomeRef.Revision, record.OutcomeRef.ManifestSHA256}
	checks := []canonicaljson.Value{}
	for _, n := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, contentObject(map[string]canonicaljson.Value{"id": n, "outcome": "pass", "reason": "Synthetic fixture only.", "evidence_document_ids": []canonicaljson.Value{}}))
	}
	rewrite := func(in TaskContentInput) TaskContentInput {
		raw, _ := os.ReadFile(in.File)
		v, _ := canonicaljson.DecodeStrict(raw)
		m := contentFields(v)
		m["task_id"] = string(id)
		raw, _ = canonicaljson.Marshal(contentObject(m))
		if e := os.WriteFile(in.File, raw, 0600); e != nil {
			t.Fatal(e)
		}
		in.TaskID = id
		return in
	}
	assessment, e := AssessTaskSpec(f.dependencies, rewrite(contentFixtureDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", key+"-assessment", map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(ref), "expected_previous_assessment": nil, "outcome": "ready", "reason": "Synthetic readiness.", "open_questions": []canonicaljson.Value{}, "checks": checks, "documents": []canonicaljson.Value{}})))
	if e != nil {
		t.Fatal(e)
	}
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	state := taskContentState(r, id)
	selected, e := SelectTaskSolution(f.dependencies, rewrite(contentFixtureDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", key+"-selection", map[string]canonicaljson.Value{"expected_previous_selection": contentRefValue(state.SelectionEvent), "action": "select", "solution": contentObject(map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(ref), "problem": contentRefValue(state.ProblemHead), "assessment": contentRefValue(TaskDecisionRef{*assessment.OutcomeRef.ID, assessment.OutcomeRef.ManifestSHA256})}), "reason": "Synthetic choice.", "human_decision": contentObject(map[string]canonicaljson.Value{"actor_claim": "fixture only", "decided_at_utc": "2026-09-30T00:00:00Z", "source": "explicit_human_instruction", "statement": "Synthetic selection; not actual human QA."})})))
	if e != nil {
		t.Fatal(e)
	}
	return TaskDecisionRef{*selected.OutcomeRef.ID, selected.OutcomeRef.ManifestSHA256}
}
func queueTargetFixture() QueueTargetInput { return QueueTargetInput{"ply", "ply", "epic"} }
func queueSetFixture(t *testing.T, f workItemJourneyFixture, key string, entries []TaskGoalQueueEntry) WorkspaceTaskQueueReadback {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	ep, b := queueTargetBinding(r, QueueTarget{EpicID: "epic", RepoID: "ply"})
	tgt := QueueTarget{ProjectID: ep.ProjectID, RepoID: b.RepoID, EpicID: ep.ID}
	q := foldQueue(r, queueID(root, tgt))
	draft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: key, ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: q.Revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "synthetic fixture", ControlSurface: "synthetic test", RecordedAtUTC: "2026-09-30T00:00:00Z"}, RegistryUpgrade: queueUpgrade(r)}
	raw, _ := contentCanonical(draft)
	file := filepath.Join(root, "queue.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	out, e := SetTaskQueue(f.dependencies, file)
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func queueStateBytes(t *testing.T, f workItemJourneyFixture) string {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	m := map[string]string{}
	e := filepath.Walk(filepath.Join(root, MarkerDirectory), func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.Mode().IsRegular() {
			b, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			m[path] = digestTaskBytes(b)
		}
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	m["refs"] = runLocalGit(t, filepath.Join(f.wrapper, "main"), "show-ref")
	m["inventory"] = runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "list", "--porcelain")
	return queueDigest(m)
}
func TestTaskQueueQ1Q2Q3Q4(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	before := queueStateBytes(t, f)
	q, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
	if e != nil || q.Revision != 0 || len(q.Pending) != 0 {
		t.Fatalf("empty: %#v %v", q, e)
	}
	if _, e = ListTaskQueue(f.dependencies, QueueTargetInput{}, false); e == nil {
		t.Fatal("workspace cwd inferred an Epic")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("read-only view wrote state")
	}
	if _, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "Task", Description: "Second", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	first, second := goalFixture(t, f, "task", "a"), goalFixture(t, f, "task-b", "b")
	q = queueSetFixture(t, f, "priority", []TaskGoalQueueEntry{{TaskID: "task-b", Goal: &second}, {TaskID: "task", Goal: &first}})
	if len(q.Pending) != 2 || q.Pending[0].TaskID != "task-b" || q.Pending[1].TaskID != "task" || q.Pending[0].State != "ready" {
		t.Fatalf("queue: %#v", q)
	}
	q = queueSetFixture(t, f, "reorder", []TaskGoalQueueEntry{{TaskID: "task", Goal: &first}, {TaskID: "task-b", Goal: &second}})
	if q.Pending[0].TaskID != "task" || q.Pending[0].Rank != 1 {
		t.Fatal("priority ignored")
	}
	// A selected solution alone was input to the removed preparation.
	if _, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-c", Title: "Task", Description: "Third", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	selected := queueSelect(t, f, "task-c", "selection", nil)
	before = queueStateBytes(t, f)
	r := mustQueueRegistry(t, f)
	selectionDraft := WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "selection-only", ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: foldQueue(r, q.QueueID).Revision, Entries: []TaskGoalQueueEntry{{TaskID: "task-c", Selection: &selected}}, Recorder: TaskGoalQueueRecorder{ActorClaim: "synthetic fixture", ControlSurface: "synthetic test", RecordedAtUTC: "2026-09-30T00:00:00Z"}}
	rawSelection, _ := contentCanonical(selectionDraft)
	os.WriteFile(filepath.Join(root, "selection-queue.json"), rawSelection, 0600)
	if _, e = SetTaskQueue(f.dependencies, filepath.Join(root, "selection-queue.json")); e == nil || !strings.Contains(e.Error(), "task_queue_goal_required") || before != queueStateBytes(t, f) {
		t.Fatalf("selection-only entry accepted or wrote state: %v", e)
	}
	all, e := ListTasks(f.dependencies, nil)
	if e != nil || len(all.Tasks) != 3 {
		t.Fatal("registered task list changed")
	}
	draft := queueFixtureDraft(t, f, "task", "dependency", taskContentState(mustQueueRegistry(t, f), "task").SpecHead)
	raw, _ := os.ReadFile(draft.File)
	raw = bytes.Replace(raw, []byte(`"dependencies":[]`), []byte(`"dependencies":["task-b"]`), 1)
	os.WriteFile(draft.File, raw, 0600)
	before = queueStateBytes(t, f)
	if _, e = RecordTaskSpec(f.dependencies, draft); e == nil {
		t.Fatal("dependency accepted")
	}
	raw, _ = os.ReadFile(filepath.Join(root, "queue.json"))
	raw = append([]byte(`{"unknown":true,`), raw[1:]...)
	os.WriteFile(filepath.Join(root, "bad-queue.json"), raw, 0600)
	if _, e = SetTaskQueue(f.dependencies, filepath.Join(root, "bad-queue.json")); e == nil {
		t.Fatal("unknown queue field accepted")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("rejected inputs wrote state")
	}
}
func mustQueueRegistry(t *testing.T, f workItemJourneyFixture) WorkItemRegistry {
	t.Helper()
	r, e := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func TestEpicBaseB1B2AndCurrentSpecFreshness(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []TaskGoalQueueEntry{})
	r := mustQueueRegistry(t, f)
	oldSpecs := queueDigest(r.TaskSpecRevisions)
	oldTasks := queueDigest(r.Tasks)
	oldOps := queueDigest(r.WorktreeOperations)
	epic := filepath.Join(f.wrapper, "epic")
	runLocalGit(t, epic, "commit", "--allow-empty", "-m", "advance base")
	in := EpicBaseInput{EpicID: "epic", RepoID: "ply"}
	preview, e := UpdateEpicBase(f.dependencies, in)
	if e != nil || preview.Confirmation == nil {
		t.Fatalf("base preview: %#v %v", preview, e)
	}
	before := runLocalGit(t, epic, "show-ref")
	in.Apply = true
	in.Confirmation = *preview.Confirmation
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := UpdateEpicBase(f.dependencies, in); errs <- e }()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	r = mustQueueRegistry(t, f)
	if len(r.EpicBaseVersions) != 2 || oldSpecs != queueDigest(r.TaskSpecRevisions) || oldTasks != queueDigest(r.Tasks) || oldOps != queueDigest(r.WorktreeOperations) || before != runLocalGit(t, epic, "show-ref") {
		t.Fatal("base update rewrote history or Git")
	}
	in.Apply = false
	in.Confirmation = ""
	same, e := UpdateEpicBase(f.dependencies, in)
	if e != nil || same.State != "unchanged" || same.Confirmation != nil {
		t.Fatal("same OID is not no-op")
	}
	task, _ := findTask(r, "task")
	ps, rs, _ := f.dependencies.Projects.Snapshot(filepath.Dir(f.wrapper))
	eval, e := currentTaskSpec(f.dependencies, filepath.Dir(f.wrapper), r, ProjectSnapshot{Projects: ps, Repos: rs}, *task, false)
	if e != nil || eval.TargetFreshness != "stale" || eval.ContentIntegrity != "valid" {
		t.Fatalf("historical integrity confused with current: %#v %v", eval, e)
	}
}
func TestTaskQueueH1StrictMigration(t *testing.T) {
	for _, version := range []int{1, 3} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, version < 3)
			root := filepath.Dir(f.wrapper)
			before, _ := os.ReadFile(workItemsPath(root))
			old := mustQueueRegistry(t, f)
			queueSetFixture(t, f, "migrate", []TaskGoalQueueEntry{})
			r := mustQueueRegistry(t, f)
			backup, e := f.dependencies.TaskContent.Read(root, "backups", digestTaskBytes(before))
			if e != nil || !bytes.Equal(before, backup) || r.FormatVersion != 4 || !reflect.DeepEqual(old.Tasks, r.Tasks) || !reflect.DeepEqual(old.Epics, r.Epics) {
				t.Fatal("migration lost preserved bytes")
			}
			encoded, e := encodeWorkItemRegistry(r)
			if e != nil {
				t.Fatal(e)
			}
			for _, change := range [][2]string{{"epic_base_versions:", "unknown_base_versions:"}, {"task_preparations: []", "task_preparations: null"}} {
				bad := bytes.Replace(encoded, []byte(change[0]), []byte(change[1]), 1)
				if _, e = decodeWorkItemRegistry(bad); e == nil {
					t.Fatal("invalid format 4 accepted")
				}
			}
			if _, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: "new-task", Title: "New", Description: "New task", ProjectID: "ply", RepoID: "ply", ParentEpicID: "epic"}); e != nil {
				t.Fatal(e)
			}
			if mustQueueRegistry(t, f).FormatVersion != 4 {
				t.Fatal("legacy command downgraded registry")
			}
		})
	}
}
