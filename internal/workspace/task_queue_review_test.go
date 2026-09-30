package workspace

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"testing"
)

func TestQueueReviewRejectsHistoryRewriteAtPublish(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "first", []QueueEntry{})
	root := filepath.Dir(f.wrapper)
	before := queueStateBytes(t, f)
	e := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, e := s.Snapshot()
		if e != nil {
			return e
		}
		r.TaskQueueEvents[0].Request["publication_key"] = "rewritten"
		r.TaskQueueEvents[0].RequestSHA256 = queueDigest(r.TaskQueueEvents[0].Request)
		return s.Publish(r)
	})
	if e == nil {
		t.Fatal("published a replacement for immutable queue history")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("rejected history rewrite changed state")
	}
}
func TestQueueReviewRejectsPreparationWithForeignBaseRevision(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	queuePrepareFixture(t, f, true)
	r := mustQueueRegistry(t, f)
	p := &r.TaskPreparations[0]
	oldID := p.ID
	p.Plan.BaseRevision = 999
	p.PlanSHA256 = queueDigest(p.Plan)
	p.ID = "pre_" + p.PlanSHA256[len("sha256:"):]
	for i := range r.TaskQueueEvents {
		ev := &r.TaskQueueEvents[i]
		if ev.Kind == "reserve" && ev.Request["preparation_id"] == oldID {
			ev.Request["preparation_id"] = p.ID
			ev.RequestSHA256 = queueDigest(ev.Request)
		}
	}
	if _, e := encodeWorkItemRegistry(r); e == nil {
		t.Fatal("accepted a plan with a fabricated base revision")
	}
}
func TestQueueReviewReadRejectsCorruptPreparationInput(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	out := queuePrepareFixture(t, f, true)
	path := out.Plan.RequiredInputs[0].Locator
	if e := os.Chmod(path, 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(path, []byte("corrupt"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper)); e == nil {
		t.Fatal("format 4 registry read skipped content closure")
	}
}
func TestQueueReviewNestedTaskCwdDoesNotSelectEpic(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	path := filepath.Join(f.wrapper, "epic", "nested")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "nested", path, f.oid)
	f.dependencies.Files = projectCwdFileSystem{FileSystem: f.dependencies.Files, cwd: path}
	if _, e := ListTaskQueue(f.dependencies, QueueTargetInput{}, false); e == nil {
		t.Fatal("Task cwd inferred its containing Epic")
	}
}

func TestQueuePreparationDriftDoesNotReserveAlternateTask(t *testing.T) {
	for _, mode := range []string{"priority", "selection", "title", "parent", "base", "path", "ref"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			sel := queueSelect(t, f, "task", "a", nil)
			queueSetFixture(t, f, "first", []QueueEntry{{"task", &sel}})
			preview := queuePrepareFixture(t, f, false)
			in := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: true, Confirmation: *preview.Confirmation}
			switch mode {
			case "priority":
				queueSetFixture(t, f, "reordered", []QueueEntry{{"task", &sel}})
			case "selection":
				r := mustQueueRegistry(t, f)
				queueSelect(t, f, "task", "new", taskContentState(r, "task").SpecHead)
			case "title":
				if _, e := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "title-change")); e != nil {
					t.Fatal(e)
				}
			case "parent", "base":
				runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "new base")
				if mode == "base" {
					p, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
					if e != nil {
						t.Fatal(e)
					}
					if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *p.Confirmation}); e != nil {
						t.Fatal(e)
					}
				}
			case "path":
				if e := os.Mkdir(preview.Plan.WorktreePath, 0700); e != nil {
					t.Fatal(e)
				}
			case "ref":
				runLocalGit(t, filepath.Join(f.wrapper, "main"), "branch", preview.Plan.Branch, f.oid)
			}
			before := queueStateBytes(t, f)
			if _, e := PrepareTask(f.dependencies, in); e == nil {
				t.Fatal("stale confirmation accepted")
			}
			if before != queueStateBytes(t, f) || len(mustQueueRegistry(t, f).TaskPreparations) != 0 {
				t.Fatal("stale confirmation reserved or wrote state")
			}
		})
	}
}
func TestQueueContextSubdirectoryAndAmbiguousTargets(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	dir := filepath.Join(f.wrapper, "epic", "subdir")
	if e := os.Mkdir(dir, 0700); e != nil {
		t.Fatal(e)
	}
	f.dependencies.Files = projectCwdFileSystem{FileSystem: f.dependencies.Files, cwd: dir}
	if _, e := ListTaskQueue(f.dependencies, QueueTargetInput{}, false); e != nil {
		t.Fatal(e)
	}
	r := mustQueueRegistry(t, f)
	other := r.Epics[0]
	other.ID = "other"
	r.Epics = append(r.Epics, other)
	ps, rs, e := f.dependencies.Projects.Snapshot(filepath.Dir(f.wrapper))
	if e != nil {
		t.Fatal(e)
	}
	if _, _, e = resolveQueueTarget(f.dependencies, filepath.Dir(f.wrapper), r, ProjectSnapshot{Projects: ps, Repos: rs}, QueueTargetInput{}); e == nil {
		t.Fatal("ambiguous Epic context selected a target")
	}
}

func TestQueueLegacyTaskCannotStartIntegrationAfterBaseChanged(t *testing.T) {
	f := newIntegrationJourneyFixtureVersion(t, true)
	queueSetFixture(t, f.workItemJourneyFixture, "upgrade", []QueueEntry{})
	epic := filepath.Join(f.wrapper, "epic")
	runLocalGit(t, epic, "commit", "--allow-empty", "-m", "new base")
	p, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *p.Confirmation}); e != nil {
		t.Fatal(e)
	}
	current := runLocalGit(t, epic, "rev-parse", "HEAD")
	runLocalGit(t, epic, "update-ref", "-m", "synthetic external rewind", "refs/heads/epic", f.oid, current)
	before := queueStateBytes(t, f.workItemJourneyFixture)
	out, e := CheckTaskIntegration(f.dependencies, TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedParentOID: f.oid, ExpectedResultOID: f.resultOID})
	if e == nil && out.Readback.Classification == "ready" {
		t.Fatal("new integration allowed from superseded legacy base")
	}
	if before != queueStateBytes(t, f.workItemJourneyFixture) {
		t.Fatal("check wrote state")
	}
}

func TestQueuePublishFailureReportsObservedCommittedRevision(t *testing.T) {
	for _, action := range []string{"set", "advance", "release"} {
		t.Run(action, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			sel := queueSelect(t, f, "task", "a", nil)
			queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
			var p TaskPreparationReadback
			if action != "set" {
				p = queuePrepareFixture(t, f, true)
			}
			file := filepath.Join(filepath.Dir(f.wrapper), "queue.json")
			if action == "set" {
				r := mustQueueRegistry(t, f)
				draft := WorkspaceTaskQueueDraft{"WorkspaceTaskQueueDraft@1", 1, "second", "ply", "ply", "epic", 1, []QueueEntry{{"task", &sel}}, QueueHumanDecision{"fixture", "2026-09-30T00:00:00Z", "explicit_human_instruction", "synthetic"}, queueUpgrade(r)}
				raw, _ := contentCanonical(draft)
				if e := os.WriteFile(file, raw, 0600); e != nil {
					t.Fatal(e)
				}
			}
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
				if point == "directory-sync" {
					return errInjected
				}
				return nil
			}}}
			var out WorkspaceTaskQueueReadback
			var e error
			want := 2
			if action == "set" {
				out, e = SetTaskQueue(f.dependencies, file)
			} else {
				want = 3
				out, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "synthetic closure", action)
			}
			if e == nil {
				t.Fatal("expected publication error")
			}
			if out.Revision != want || out.QueueID == "" {
				t.Fatalf("lost observed committed revision: %#v", out)
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			before := queueStateBytes(t, f)
			if action == "set" {
				_, e = SetTaskQueue(f.dependencies, file)
			} else {
				_, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "synthetic closure", action)
			}
			if e != nil || before != queueStateBytes(t, f) {
				t.Fatalf("retry changed effect: %v", e)
			}
		})
	}
}

func TestEpicBaseHistoryReadbackSeparatesSelectionAndEffectRelevance(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	queueSetFixture(t, f.workItemJourneyFixture, "upgrade", []QueueEntry{})
	epic := filepath.Join(f.wrapper, "epic")
	runLocalGit(t, epic, "commit", "--allow-empty", "-m", "new current base")
	p, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *p.Confirmation}); e != nil {
		t.Fatal(e)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	task, _ := findTask(r, "task")
	ps, rs, e := f.dependencies.Projects.Snapshot(f.root)
	if e != nil {
		t.Fatal(e)
	}
	eval, e := currentTaskSpec(f.dependencies, f.root, r, ProjectSnapshot{Projects: ps, Repos: rs}, *task, true)
	if e != nil || eval.SelectionFreshness != "current" || eval.TargetFreshness != "stale" || eval.ContentIntegrity != "valid" {
		t.Errorf("selection conflated with target: %#v %v", eval, e)
	}
	show, e := ShowTask(f.dependencies, "task")
	if e != nil {
		t.Fatal(e)
	}
	raw, e := MarshalTaskReadback(show)
	if e != nil {
		t.Fatal(e)
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		t.Fatal(e)
	}
	results := contentArray(contentFields(v), "results")
	if len(results) != 1 || contentString(contentFields(results[0]), "relevance") != "stale" || contentString(contentFields(results[0]), "basis_status") != "bound" {
		t.Fatalf("historical result falsely current: %#v", results)
	}
}

func TestEpicBaseResultRelevanceObservesLiveParent(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	queueSetFixture(t, f.workItemJourneyFixture, "upgrade", []QueueEntry{})
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	task, _ := findTask(r, "task")
	ps, rs, e := f.dependencies.Projects.Snapshot(f.root)
	if e != nil {
		t.Fatal(e)
	}
	projects := ProjectSnapshot{Projects: ps, Repos: rs}
	if got := taskResultSpecRelevance(f.dependencies, f.root, projects, r, *task, f.resultID); got.Relevance != "current" {
		t.Fatalf("result baseline: %#v", got)
	}
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "external advance")
	if got := taskResultSpecRelevance(f.dependencies, f.root, projects, r, *task, f.resultID); got.Relevance != "stale" {
		t.Errorf("live parent move ignored: %#v", got)
	}
	f.dependencies.WorkGit = queueObservationFault{f.dependencies.WorkGit, "unknown"}
	if got := taskResultSpecRelevance(f.dependencies, f.root, projects, r, *task, f.resultID); got.Relevance != "unknown" {
		t.Errorf("missing parent observation ignored: %#v", got)
	}
}
