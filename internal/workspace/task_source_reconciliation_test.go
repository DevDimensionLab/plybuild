package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelParentMovementInvalidatesUnusedReconciliation(t *testing.T) {
	f, prepared, in := sourceReconcileFixture(t, false)
	p, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	source := prepared.Preparation.Plan.WorktreePath
	before := runLocalGit(t, source, "rev-parse", "HEAD")
	parent := prepared.Preparation.Plan.Target.ParentLocator
	runLocalGit(t, parent, "commit", "--allow-empty", "-m", "C concurrent parent move")
	base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
		t.Fatal(e)
	}
	state := queueStateBytes(t, f)
	in.Apply = true
	in.Confirmation = p.Confirmation
	if _, e = ReconcileTaskSource(f.dependencies, in); e == nil {
		t.Fatal("stale unused source reconciliation preview was applied")
	}
	if queueStateBytes(t, f) != state || runLocalGit(t, source, "rev-parse", "HEAD") != before {
		t.Fatal("stale reconciliation had an effect")
	}
}

func TestParallelSourceCannotDiscardCommittedReconciliationHistory(t *testing.T) {
	f, prepared, in := sourceReconcileFixture(t, false)
	p, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	in.Apply = true
	in.Confirmation = p.Confirmation
	merged, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	source := prepared.Preparation.Plan.WorktreePath
	// Isolated external drift: source loses the native receipt's merge and B
	// history while still descending from the original source start and P1.
	runLocalGit(t, source, "reset", "--hard", merged.Operation.Plan.Parent.OID)
	r := mustQueueRegistry(t, f)
	basis, e := EffectiveTaskDeliveryBasis(r, prepared.Preparation.Preparation.ID)
	if e != nil {
		t.Fatal(e)
	}
	e = WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
		_, e := s.ValidateDeliveryCandidateTargetBasis(source, "refs/heads/"+prepared.Preparation.Plan.Branch, prepared.Preparation.Plan.Target.GitCommonDir, &prepared.Basis, &basis)
		return e
	})
	if e == nil {
		t.Fatal("candidate discarded committed source reconciliation history")
	}
}

func TestParallelSourceReconciliationRealMerge(t *testing.T) {
	f := parallelFixture(t)
	g := parallelGoalFixture(t, f, "task", "parallel-source")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	in := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: g.SpecID, Parallel: true}
	p, err := PreviewTaskGoalExecution(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := PrepareTaskGoalExecution(f.dependencies, in, p.Confirmation, goalHumanFixture())
	if err != nil {
		t.Fatal(err)
	}
	source := prepared.Preparation.Plan.WorktreePath
	if err = os.WriteFile(filepath.Join(source, "b.txt"), []byte("B\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, source, "add", "b.txt")
	runLocalGit(t, source, "commit", "-m", "B")
	before := runLocalGit(t, source, "rev-parse", "HEAD")
	parent := filepath.Join(f.wrapper, "epic")
	if err = os.WriteFile(filepath.Join(parent, "a.txt"), []byte("A\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, parent, "add", "a.txt")
	runLocalGit(t, parent, "commit", "-m", "A native parent fixture")
	base, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); err != nil {
		t.Fatal(err)
	}
	input := TaskSourceReconciliationInput{PreparationID: prepared.Preparation.Preparation.ID, RunID: "synthetic-native-run", RequestSHA256: "sha256:" + strings.Repeat("c", 64), MandateSHA256: "sha256:" + strings.Repeat("d", 64), ReviewSHA256: "sha256:" + strings.Repeat("e", 64)}
	preview, err := ReconcileTaskSource(f.dependencies, input)
	if err != nil {
		t.Fatalf("native reconcile preview: %v", err)
	}
	if runLocalGit(t, source, "rev-parse", "HEAD") != before {
		t.Fatal("preview changed source")
	}
	input.Apply = true
	input.Confirmation = preview.Confirmation
	result, err := ReconcileTaskSource(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "committed" || result.Operation == nil || result.Operation.Plan.Generation != 1 {
		t.Fatalf("no committed exact effective basis: %+v", result)
	}
	after := runLocalGit(t, source, "rev-parse", "HEAD")
	runLocalGit(t, source, "merge-base", "--is-ancestor", before, after)
	runLocalGit(t, source, "merge-base", "--is-ancestor", runLocalGit(t, parent, "rev-parse", "HEAD"), after)
	if raw, err := os.ReadFile(filepath.Join(source, "a.txt")); err != nil || string(raw) != "A\n" {
		t.Fatal("parent composition missing")
	}
	again, err := ReconcileTaskSource(f.dependencies, input)
	if err != nil || again.Operation.ID != result.Operation.ID || runLocalGit(t, source, "rev-parse", "HEAD") != after {
		t.Fatalf("reconciliation retry repeated effect: %v", err)
	}
	r := mustQueueRegistry(t, f)
	if r.TaskPreparations[0].Plan.ParentOID != p.ParentOID || len(r.SourceReconciliations) != 1 {
		t.Fatal("reconciliation replaced original history")
	}
}

func sourceReconcileFixture(t *testing.T, conflict bool) (workItemJourneyFixture, TaskGoalPreparedExecution, TaskSourceReconciliationInput) {
	t.Helper()
	f := parallelFixture(t)
	g := parallelGoalFixture(t, f, "task", "source-operation")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	gi := TaskGoalExecuteInput{Target: queueTargetFixture(), SpecID: g.SpecID, Parallel: true}
	preview, e := PreviewTaskGoalExecution(f.dependencies, gi)
	if e != nil {
		t.Fatal(e)
	}
	prepared, e := PrepareTaskGoalExecution(f.dependencies, gi, preview.Confirmation, goalHumanFixture())
	if e != nil {
		t.Fatal(e)
	}
	source, parent := prepared.Preparation.Plan.WorktreePath, filepath.Join(f.wrapper, "epic")
	file := "b.txt"
	if conflict {
		file = "README.md"
	}
	if e = os.WriteFile(filepath.Join(source, file), []byte("B source\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, source, "add", file)
	runLocalGit(t, source, "commit", "-m", "B source")
	file = "a.txt"
	if conflict {
		file = "README.md"
	}
	if e = os.WriteFile(filepath.Join(parent, file), []byte("A parent\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, parent, "add", file)
	runLocalGit(t, parent, "commit", "-m", "A parent")
	base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
		t.Fatal(e)
	}
	return f, prepared, TaskSourceReconciliationInput{PreparationID: prepared.Preparation.Preparation.ID, RunID: "simulated-run", RequestSHA256: "sha256:" + strings.Repeat("c", 64), MandateSHA256: "sha256:" + strings.Repeat("d", 64), ReviewSHA256: "sha256:" + strings.Repeat("e", 64)}
}
func TestParallelSourceConflictContinueAbortAndPreserveEdits(t *testing.T) {
	for _, action := range []string{"continue", "abort", "abort-edited"} {
		t.Run(action, func(t *testing.T) {
			f, p, in := sourceReconcileFixture(t, true)
			source := p.Preparation.Plan.WorktreePath
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply = true
			in.Confirmation = preview.Confirmation
			conflict, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil || conflict.State != "source_conflict" || len(conflict.Operation.ConflictPaths) != 1 || conflict.Operation.ConflictPaths[0] != "README.md" {
				t.Fatalf("native text conflict: %+v %v", conflict, e)
			}
			if _, e = EffectiveTaskDeliveryBasis(mustQueueRegistry(t, f), in.PreparationID); e == nil {
				t.Fatal("pending conflict qualified a basis")
			}
			in.OperationID = conflict.Operation.ID
			in.Confirmation = ""
			in.Action = action
			if action == "continue" || action == "abort-edited" {
				if e = os.WriteFile(filepath.Join(source, "README.md"), []byte("A parent and B source resolved\n"), 0600); e != nil {
					t.Fatal(e)
				}
			}
			if action == "continue" {
				runLocalGit(t, source, "add", "README.md")
			} else if action == "abort-edited" {
				in.Action = "abort"
			}
			result, e := ReconcileTaskSource(f.dependencies, in)
			if action == "abort-edited" {
				if e == nil || !strings.Contains(e.Error(), "source_abort_edits_present") {
					t.Fatalf("edited abort was not refused: %v", e)
				}
				raw, _ := os.ReadFile(filepath.Join(source, "README.md"))
				if string(raw) != "A parent and B source resolved\n" {
					t.Fatal("abort lost owner edits")
				}
				return
			}
			if e != nil {
				t.Fatal(e)
			}
			want := "committed"
			if action == "abort" {
				want = "aborted"
			}
			if result.State != want {
				t.Fatalf("unexpected native resolution: %+v", result)
			}
			if action == "abort" {
				in.Apply = false
				in.Action = ""
				in.OperationID = ""
				fresh, e := ReconcileTaskSource(f.dependencies, in)
				if e != nil || fresh.Confirmation == preview.Confirmation {
					t.Fatalf("new explicit post-abort preview lost attempt identity: %v", e)
				}
			}
		})
	}
}
func TestParallelSourceReconciliationCrashRecovery(t *testing.T) {
	for _, point := range []string{"after_intent", "after_git", "ambiguous_after_git", "conflict_after_git"} {
		t.Run(point, func(t *testing.T) {
			f, p, in := sourceReconcileFixture(t, point == "conflict_after_git")
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply = true
			in.Confirmation = preview.Confirmation
			failPoint := point
			if strings.Contains(point, "after_git") {
				failPoint = "after_git"
			}
			f.dependencies.SourceReconciliationFault = func(stage, id string) error {
				if stage == failPoint {
					return errInjected
				}
				return nil
			}
			if _, e = ReconcileTaskSource(f.dependencies, in); e == nil {
				t.Fatal("crash not observed")
			}
			f.dependencies.SourceReconciliationFault = nil
			source := p.Preparation.Plan.WorktreePath
			before := runLocalGit(t, source, "rev-parse", "HEAD")
			if point == "ambiguous_after_git" {
				r := mustQueueRegistry(t, f)
				marker := runLocalGit(t, source, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+r.SourceReconciliations[0].ID)
				if e = os.Remove(marker); e != nil {
					t.Fatal(e)
				}
			}
			result, e := ReconcileTaskSource(f.dependencies, in)
			if point == "ambiguous_after_git" {
				if e == nil || result.State != "effect_unknown" {
					t.Fatalf("ambiguous effect replayed: %+v %v", result, e)
				}
			} else if point == "conflict_after_git" {
				if e != nil || result.State != "source_conflict" {
					t.Fatalf("conflict recovery: %+v %v", result, e)
				}
			} else if e != nil || result.State != "committed" {
				t.Fatalf("same-operation recovery: %+v %v", result, e)
			}
			if point != "after_intent" && runLocalGit(t, source, "rev-parse", "HEAD") != before {
				t.Fatal("recovery repeated source Git effect")
			}
			if len(mustQueueRegistry(t, f).SourceReconciliations) != 1 {
				t.Fatal("crash created another operation")
			}
		})
	}
}
