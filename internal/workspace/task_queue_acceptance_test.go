package workspace

import (
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"testing"
)

type queueObservationFault struct {
	WorkItemGit
	mode string
}

func (g queueObservationFault) ObserveWorktree(path string) (GitWorktreeObservation, error) {
	o, e := g.WorkItemGit.ObserveWorktree(path)
	if e != nil {
		return o, e
	}
	switch g.mode {
	case "ref":
		o.Ref = "refs/heads/foreign"
	case "worktree":
		o.Locator += "-foreign"
	case "repo":
		o.GitCommonDir += "-foreign"
	case "unknown":
		return o, errors.New("synthetic unavailable observation")
	}
	return o, nil
}
func TestEpicBaseB1RejectsChangedTargetAndB2ShowsDrift(t *testing.T) {
	for _, mode := range []string{"dirty", "ref", "worktree", "repo", "unknown", "divergent", "rewind", "previous"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "upgrade", []QueueEntry{})
			epic := filepath.Join(f.wrapper, "epic")
			runLocalGit(t, epic, "commit", "--allow-empty", "-m", "base two")
			in := EpicBaseInput{EpicID: "epic", RepoID: "ply"}
			p, e := UpdateEpicBase(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			switch mode {
			case "dirty":
				if e = os.WriteFile(filepath.Join(epic, "dirty"), []byte("dirty"), 0600); e != nil {
					t.Fatal(e)
				}
			case "ref", "worktree", "repo", "unknown":
				f.dependencies.WorkGit = queueObservationFault{f.dependencies.WorkGit, mode}
			case "divergent":
				tree := runLocalGit(t, epic, "rev-parse", "HEAD^{tree}")
				foreign := runLocalGit(t, epic, "commit-tree", tree, "-m", "unrelated root")
				runLocalGit(t, epic, "update-ref", "refs/heads/epic", foreign)
			case "rewind":
				if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *p.Confirmation}); e != nil {
					t.Fatal(e)
				}
				runLocalGit(t, epic, "update-ref", "refs/heads/epic", f.oid)
			case "previous":
				runLocalGit(t, epic, "commit", "--allow-empty", "-m", "base three")
				in.Apply = true
				in.Confirmation = *p.Confirmation
			}
			before := queueStateBytes(t, f)
			if _, e = UpdateEpicBase(f.dependencies, in); e == nil {
				t.Fatal("changed base accepted")
			}
			if before != queueStateBytes(t, f) {
				t.Fatal("rejected base changed bytes")
			}
			if mode == "rewind" {
				show, e := ShowEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
				if e != nil || show.Freshness != "stale" || len(show.Versions) != 2 {
					t.Fatalf("postpublish drift: %#v %v", show, e)
				}
			}
		})
	}
}
func TestQueueQ2SelectionAssessmentAndProblemDrift(t *testing.T) {
	for _, mode := range []string{"selection", "assessment", "problem"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			sel := queueSelect(t, f, "task", "initial", nil)
			queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
			r := mustQueueRegistry(t, f)
			switch mode {
			case "selection":
				next := queueSelect(t, f, "task", "second", taskContentState(r, "task").SpecHead)
				q, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
				if e != nil || q.Pending[0].State != "blocked" {
					t.Fatal("old selection remains ready")
				}
				q = queueSetFixture(t, f, "rebind", []QueueEntry{{"task", &next}})
				if q.Pending[0].State != "ready" {
					t.Fatal("exact new binding not ready")
				}
				return
			case "problem":
				if _, e := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "revised")); e != nil {
					t.Fatal(e)
				}
			case "assessment":
				a := lastTaskAssessment(r, "task", "solution", 1)
				manifest, e := readRegisteredTaskManifest(f.dependencies.TaskContent, filepath.Dir(f.wrapper), r, a.ManifestSHA256)
				if e != nil {
					t.Fatal(e)
				}
				m := contentFields(manifest)
				draft := contentFixtureDraft(t, filepath.Dir(f.wrapper), "WorkspaceTaskSpecAssessmentDraft@1", "negative", map[string]canonicaljson.Value{"spec_id": "solution", "spec": m["spec"], "expected_previous_assessment": contentRefValue(TaskDecisionRef{a.ID, a.ManifestSHA256}), "outcome": "needs_work", "reason": "Synthetic negative assessment", "open_questions": []canonicaljson.Value{}, "checks": m["checks"], "documents": []canonicaljson.Value{}})
				if _, e = AssessTaskSpec(f.dependencies, draft); e != nil {
					t.Fatal(e)
				}
			}
			before := queueStateBytes(t, f)
			q, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false)
			if e != nil || len(q.Pending) != 1 || q.Pending[0].State != "blocked" {
				t.Fatalf("stale selection: %#v %v", q, e)
			}
			out, e := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}})
			if e != nil || out.Confirmation != nil || out.Plan != nil {
				t.Fatal("no-ready queue offered effect")
			}
			if before != queueStateBytes(t, f) {
				t.Fatal("read-only changed bytes")
			}
		})
	}
}
func TestPreparationP2TitleAfterReservationP4NamedConflictP6CAS(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "Task", Description: "Second", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	other := queueSelect(t, f, "task-b", "b", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}, {"task-b", &other}})
	p := queuePrepareFixture(t, f, true)
	oldPlan := queueDigest(p.Plan)
	if _, e := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "retitle")); e != nil {
		t.Fatal(e)
	}
	out, e := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: true, Confirmation: *p.Confirmation})
	if e != nil || queueDigest(out.Plan) != oldPlan {
		t.Fatal("reserved name or basis changed")
	}
	id := TaskID("task-b")
	before := queueStateBytes(t, f)
	if _, e = PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "task", TaskID: &id}}); e == nil {
		t.Fatal("different named task bypassed current")
	}
	if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 1, "wrong revision", "advance"); e == nil {
		t.Fatal("wrong revision accepted")
	}
	wrongID := "pre_" + string(make([]byte, 64))
	if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), wrongID, 2, "wrong current", "advance"); e == nil {
		t.Fatal("wrong current accepted")
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("conflict changed state")
	}
	q, e := CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "synthetic release", "release")
	if e != nil || q.Current != nil || len(q.Pending) != 1 || len(mustQueueRegistry(t, f).TaskPreparations) != 1 {
		t.Fatal("release created next or lost history")
	}
	next := queuePrepareFixture(t, f, false)
	if next.Plan.TaskID != "task-b" {
		t.Fatal("next pending not available")
	}
}

func TestPreparationP3SimulatedNamingCollisionDoesNotChooseFallback(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	preview := queuePrepareFixture(t, f, false)
	r := mustQueueRegistry(t, f)
	// Simulate a preserved foreign intent that claims the exact derived name, even
	// though neither its path nor branch exists yet. Never publish this test snapshot.
	r.WorktreeOperations = append(r.WorktreeOperations, WorktreeOperationRecord{TaskID: "foreign", GitCommonDir: preview.Target.GitCommonDir, SourceRef: "refs/heads/" + preview.Plan.Branch, TargetLocator: preview.Plan.WorktreePath})
	ps, rs, e := f.dependencies.Projects.Snapshot(filepath.Dir(f.wrapper))
	if e != nil {
		t.Fatal(e)
	}
	before := queueStateBytes(t, f)
	row, plan := queueEntryEvaluation(f.dependencies, filepath.Dir(f.wrapper), r, ProjectSnapshot{Projects: ps, Repos: rs}, preview.Target, preview.Plan.QueueID, QueueEntry{"task", &sel}, 1)
	if plan != nil || row.State != "blocked" || row.Reasons[0].Code != "task_queue_resource_conflict" {
		t.Fatalf("collision selected fallback: %#v %#v", row, plan)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("collision wrote effect")
	}
}
func TestQueueH1LegacyReadAndWrongMigrationHashAreReadOnly(t *testing.T) {
	for _, version := range []int{1, 2, 3} {
		t.Run(string(rune('0'+version)), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, version < 3)
			root := filepath.Dir(f.wrapper)
			if version == 2 {
				if e := f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
					r, e := s.Snapshot()
					if e != nil {
						return e
					}
					upgradeRegistryToV2(&r)
					return s.Publish(r)
				}); e != nil {
					t.Fatal(e)
				}
			}
			before := queueStateBytes(t, f)
			if _, e := ListTaskQueue(f.dependencies, queueTargetFixture(), false); e != nil {
				t.Fatal(e)
			}
			if _, e := ShowEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"}); e != nil {
				t.Fatal(e)
			}
			if before != queueStateBytes(t, f) {
				t.Fatal("legacy read migrated registry")
			}
			draft := WorkspaceTaskQueueDraft{"WorkspaceTaskQueueDraft@1", 1, "bad-upgrade", "ply", "ply", "epic", 0, []QueueEntry{}, QueueHumanDecision{"fixture", "2026-09-30T00:00:00Z", "explicit_human_instruction", "synthetic"}, &TaskRegistryUpgrade{FromVersion: version, RegistrySHA256: "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}}
			raw, _ := contentCanonical(draft)
			file := filepath.Join(root, "queue.json")
			if e := os.WriteFile(file, raw, 0600); e != nil {
				t.Fatal(e)
			}
			if _, e := SetTaskQueue(f.dependencies, file); e == nil {
				t.Fatal("wrong migration hash accepted")
			}
			if before != queueStateBytes(t, f) {
				t.Fatal("bad migration wrote state")
			}
		})
	}
}
