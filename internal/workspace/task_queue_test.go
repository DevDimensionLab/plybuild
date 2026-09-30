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
func queueSetFixture(t *testing.T, f workItemJourneyFixture, key string, entries []QueueEntry) WorkspaceTaskQueueReadback {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	ep, b := queueTargetBinding(r, QueueTarget{EpicID: "epic", RepoID: "ply"})
	tgt := QueueTarget{ProjectID: ep.ProjectID, RepoID: b.RepoID, EpicID: ep.ID}
	q := foldQueue(r, queueID(root, tgt))
	draft := WorkspaceTaskQueueDraft{"WorkspaceTaskQueueDraft@1", 1, key, "ply", "ply", "epic", q.Revision, entries, QueueHumanDecision{"synthetic fixture", "2026-09-30T00:00:00Z", "explicit_human_instruction", "Synthetic priority, not actual approval."}, queueUpgrade(r)}
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
func queuePrepareFixture(t *testing.T, f workItemJourneyFixture, apply bool) TaskPreparationReadback {
	t.Helper()
	in := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}}
	out, e := PrepareTask(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	if apply {
		if out.Confirmation == nil {
			t.Fatalf("no preview: %#v", out)
		}
		in.Apply = true
		in.Confirmation = *out.Confirmation
		out, e = PrepareTask(f.dependencies, in)
		if e != nil {
			t.Fatal(e)
		}
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
	selected := queueSelect(t, f, "task", "a", nil)
	if _, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "Task", Description: "Second", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	other := queueSelect(t, f, "task-b", "b", nil)
	q = queueSetFixture(t, f, "priority", []QueueEntry{{"task-b", nil}, {"task", &selected}})
	if q.Pending[0].State != "blocked" || q.Pending[1].State != "ready" {
		t.Fatalf("queue: %#v", q)
	}
	p := queuePrepareFixture(t, f, false)
	if p.Plan.TaskID != "task" {
		t.Fatal("next did not skip blocked first entry")
	}
	ready, e := ListTaskQueue(f.dependencies, queueTargetFixture(), true)
	if e != nil || len(ready.Pending) != 1 || ready.Pending[0].Rank != 2 {
		t.Fatalf("ready rank: %#v %v", ready, e)
	}
	r, _ := f.dependencies.WorkItems.Snapshot(root)
	ref := taskContentState(r, "task").SpecHead
	if _, e = RecordTaskSpec(f.dependencies, queueFixtureDraft(t, f, "task", "new-draft", ref)); e != nil {
		t.Fatal(e)
	}
	p = queuePrepareFixture(t, f, false)
	if p.Plan.Spec.Revision != 1 {
		t.Fatal("unselected draft replaced selection")
	}
	q = queueSetFixture(t, f, "reorder", []QueueEntry{{"task-b", &other}, {"task", &selected}})
	if queuePrepareFixture(t, f, false).Plan.TaskID != "task-b" {
		t.Fatal("priority ignored")
	}
	all, e := ListTasks(f.dependencies, nil)
	if e != nil || len(all.Tasks) != 2 {
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
func TestTaskPreparationP1P2P3P4P6(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	preview := queuePrepareFixture(t, f, false)
	before := queueStateBytes(t, f)
	bad := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: true, Confirmation: "sha256:" + strings.Repeat("f", 64)}
	if _, e := PrepareTask(f.dependencies, bad); e == nil {
		t.Fatal("bad confirmation accepted")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("bad confirmation wrote state")
	}
	branch1, path1 := queueTaskNames(preview.Plan.QueueID, "task", "Task", preview.Target.ParentLocator)
	branch2, path2 := queueTaskNames(preview.Plan.QueueID, "task-b", "Task", preview.Target.ParentLocator)
	if branch1 == branch2 || path1 == path2 {
		t.Fatal("title suffix collision")
	}
	in := bad
	in.Confirmation = *preview.Confirmation
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := PrepareTask(f.dependencies, in); errs <- e }()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	r := mustQueueRegistry(t, f)
	if len(r.TaskPreparations) != 1 || len(r.WorktreeOperations) != 1 || foldQueue(r, preview.Plan.QueueID).Revision != 2 {
		t.Fatal("concurrent apply duplicated intent")
	}
	p := r.TaskPreparations[0]
	if p.Outcome == nil || runLocalGit(t, p.Plan.WorktreePath, "rev-parse", "HEAD") != f.oid || runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD") != f.oid {
		t.Fatal("wrong worktree or moved Epic")
	}
	for _, input := range p.Plan.RequiredInputs {
		b, e := os.ReadFile(input.Locator)
		if e != nil || digestTaskBytes(b) != input.SHA256 || int64(len(b)) != input.SizeBytes {
			t.Fatal("required input differs")
		}
	}
	os.WriteFile(filepath.Join(p.Plan.WorktreePath, "work.txt"), []byte("normal work"), 0600)
	show, e := ShowTaskPreparation(f.dependencies, p.ID)
	if e != nil || show.State != "prepared" {
		t.Fatalf("dirty work erased historical outcome: %#v %v", show, e)
	}
	q, e := CloseTaskQueue(f.dependencies, queueTargetFixture(), p.ID, 2, "continue", "advance")
	if e != nil || q.Current != nil || q.Revision != 3 {
		t.Fatalf("advance: %#v %v", q, e)
	}
	before = queueStateBytes(t, f)
	if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.ID, 2, "continue", "advance"); e != nil {
		t.Fatal(e)
	}
	retry, e := PrepareTask(f.dependencies, in)
	if e != nil || retry.Disposition != "advanced" || before != queueStateBytes(t, f) {
		t.Fatalf("historical retry wrote state: %#v %v", retry, e)
	}
}
func TestEpicBaseB1B2AndHistoricalSpec(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	prepared := queuePrepareFixture(t, f, true)
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
	show, e := ShowTaskPreparation(f.dependencies, prepared.Preparation.ID)
	if e != nil || show.Plan.ParentOID != f.oid || show.State != "prepared" {
		t.Fatalf("historical readback: %#v %v", show, e)
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
	for _, version := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, version < 3)
			root := filepath.Dir(f.wrapper)
			if version == 2 {
				e := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
					r, e := s.Snapshot()
					if e != nil {
						return e
					}
					upgradeRegistryToV2(&r)
					return s.Publish(r)
				})
				if e != nil {
					t.Fatal(e)
				}
			}
			before, _ := os.ReadFile(workItemsPath(root))
			old := mustQueueRegistry(t, f)
			queueSetFixture(t, f, "migrate", []QueueEntry{})
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
