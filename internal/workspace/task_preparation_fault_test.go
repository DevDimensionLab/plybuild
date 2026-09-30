package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type queueFaultGit struct {
	WorkItemGit
	TaskContentGit
	mode string
}

func (g queueFaultGit) CreateBranchAndWorktree(in GitCreateInput) GitCommandOutcome {
	switch g.mode {
	case "none":
		return GitCommandOutcome{Exit: 1, Err: errors.New("injected no effect")}
	case "branch":
		return branchOnlyGit{g.WorkItemGit}.CreateBranchAndWorktree(in)
	case "partial":
		out := g.WorkItemGit.CreateBranchAndWorktree(in)
		_ = os.WriteFile(filepath.Join(in.TargetLocator, "partial.txt"), []byte("unknown effect"), 0600)
		return out
	}
	return g.WorkItemGit.CreateBranchAndWorktree(in)
}
func TestTaskPreparationP5FaultBoundariesAndRecovery(t *testing.T) {
	for _, mode := range []string{"before-intent", "after-intent", "branch", "after-exact", "partial", "release-no-effect"} {
		t.Run(mode, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			sel := queueSelect(t, f, "task", "a", nil)
			queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
			preview := queuePrepareFixture(t, f, false)
			in := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: true, Confirmation: *preview.Confirmation}
			original := f.dependencies.WorkGit
			before := queueStateBytes(t, f)
			if mode == "before-intent" || mode == "after-exact" {
				count := 0
				f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
					if point == "temp-open" {
						count++
						if mode == "before-intent" && count == 1 || mode == "after-exact" && count == 2 {
							return errInjected
						}
					}
					return nil
				}}}
			} else {
				gm := mode
				if mode == "after-intent" || mode == "release-no-effect" {
					gm = "none"
				}
				f.dependencies.WorkGit = queueFaultGit{original, original.(TaskContentGit), gm}
			}
			got, e := PrepareTask(f.dependencies, in)
			if e == nil {
				t.Fatal("fault not reported")
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			f.dependencies.WorkGit = original
			r := mustQueueRegistry(t, f)
			if mode == "before-intent" {
				if len(r.TaskPreparations) != 0 || queueStateBytes(t, f) != before {
					t.Fatal("pre-intent fault wrote effect")
				}
			} else {
				if len(r.TaskPreparations) != 1 || foldQueue(r, preview.Plan.QueueID).Current == nil {
					t.Fatal("lost reserved intent")
				}
			}
			if mode == "partial" {
				p := r.TaskPreparations[0]
				before = queueStateBytes(t, f)
				if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.ID, 2, "release unknown", "release"); e == nil {
					t.Fatal("released unknown effect")
				}
				if queueStateBytes(t, f) != before {
					t.Fatal("release changed unknown state")
				}
				if got.State != "unknown" {
					t.Fatalf("partial state %s", got.State)
				}
				return
			}
			if mode == "release-no-effect" {
				p := r.TaskPreparations[0]
				epic := filepath.Join(f.wrapper, "epic")
				runLocalGit(t, epic, "commit", "--allow-empty", "-m", "later parent")
				if _, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.ID, 2, "stop this attempt", "release"); e != nil {
					t.Fatal(e)
				}
				before = queueStateBytes(t, f)
				retry, e := PrepareTask(f.dependencies, in)
				if e != nil || retry.Disposition != "released" {
					t.Fatalf("terminal retry: %#v %v", retry, e)
				}
				legacy := TaskWorktreeCreateInput{TaskID: "task", Branch: p.Plan.Branch, Path: p.Plan.WorktreePath, ExpectedParentOID: p.Plan.ParentOID}
				if _, e = CreateTaskWorktree(f.dependencies, legacy); e != nil {
					t.Fatal(e)
				}
				if before != queueStateBytes(t, f) {
					t.Fatal("terminal legacy create wrote Git")
				}
				base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
				if e != nil || base.Confirmation == nil {
					t.Fatalf("released no-effect blocked base: %v", e)
				}
				return
			}
			recovered, e := PrepareTask(f.dependencies, in)
			if e != nil || recovered.State != "prepared" {
				t.Fatalf("recovery: %#v %v", recovered, e)
			}
			r = mustQueueRegistry(t, f)
			if len(r.TaskPreparations) != 1 || foldQueue(r, preview.Plan.QueueID).Revision != 2 {
				t.Fatal("recovery duplicated reservation")
			}
		})
	}
}
func TestEpicBaseB2PublicationFaultRecovery(t *testing.T) {
	for _, boundary := range []int{1, 2, 3} {
		t.Run(string(rune('0'+boundary)), func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			queueSetFixture(t, f, "upgrade", []QueueEntry{})
			epic := filepath.Join(f.wrapper, "epic")
			runLocalGit(t, epic, "commit", "--allow-empty", "-m", "base two")
			in := EpicBaseInput{EpicID: "epic", RepoID: "ply"}
			plan, e := UpdateEpicBase(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply = true
			in.Confirmation = *plan.Confirmation
			count := 0
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
				if point == "temp-open" {
					count++
					if boundary < 3 && count == boundary {
						return errInjected
					}
				}
				if boundary == 3 && point == "directory-sync" && count == 2 {
					return errInjected
				}
				return nil
			}}}
			_, _ = UpdateEpicBase(f.dependencies, in)
			f.dependencies.WorkItems = newSystemWorkItemStore()
			before := runLocalGit(t, epic, "show-ref")
			out, e := UpdateEpicBase(f.dependencies, in)
			if e != nil || out.State != "committed" {
				t.Fatalf("base recovery: %#v %v", out, e)
			}
			r := mustQueueRegistry(t, f)
			if len(r.EpicBaseVersions) != 2 || len(r.EpicBaseUpdates) != 1 || before != runLocalGit(t, epic, "show-ref") {
				t.Fatal("duplicate base or Git write")
			}
		})
	}
}

func TestEpicBaseConflictRetryKeepsFailureAndPreservedOutcome(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "upgrade", []QueueEntry{})
	epic := filepath.Join(f.wrapper, "epic")
	runLocalGit(t, epic, "commit", "--allow-empty", "-m", "base two")
	in := EpicBaseInput{EpicID: "epic", RepoID: "ply"}
	p, e := UpdateEpicBase(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	in.Apply = true
	in.Confirmation = *p.Confirmation
	count := 0
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
		if point == "temp-open" {
			count++
			if count == 2 {
				return errInjected
			}
		}
		return nil
	}}}
	if _, e = UpdateEpicBase(f.dependencies, in); e == nil {
		t.Fatal("expected pending intent")
	}
	f.dependencies.WorkItems = newSystemWorkItemStore()
	runLocalGit(t, epic, "commit", "--allow-empty", "-m", "conflicting observed base")
	out, e := UpdateEpicBase(f.dependencies, in)
	if e == nil || out.State != "conflict" {
		t.Fatalf("conflict missing: %#v %v", out, e)
	}
	before := queueStateBytes(t, f)
	out, e = UpdateEpicBase(f.dependencies, in)
	if e == nil || out.State != "conflict" || out.Operation == nil {
		t.Fatalf("retry lost conflict exit: %#v %v", out, e)
	}
	if before != queueStateBytes(t, f) {
		t.Fatal("terminal conflict retry wrote state")
	}
}
