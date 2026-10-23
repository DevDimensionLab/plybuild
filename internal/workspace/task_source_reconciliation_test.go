package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelConflictRetryPreservesOriginalAbortEvidence(t *testing.T) {
	for _, staged := range []bool{false, true} {
		name := "edited"
		if staged {
			name = "staged"
		}
		t.Run(name, func(t *testing.T) {
			f, prepared, in := sourceReconcileFixture(t, true)
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply, in.Confirmation = true, preview.Confirmation
			conflict, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil || conflict.State != "source_conflict" {
				t.Fatalf("conflict: %+v %v", conflict, e)
			}
			source := prepared.Preparation.Plan.WorktreePath
			resolution := []byte("owner resolution must survive retry and abort\n")
			if e = os.WriteFile(filepath.Join(source, "README.md"), resolution, 0600); e != nil {
				t.Fatal(e)
			}
			if staged {
				runLocalGit(t, source, "add", "README.md")
			}
			indexBefore := runLocalGit(t, source, "diff", "--cached", "--binary")
			in.OperationID = conflict.Operation.ID
			observed, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			if observed.Operation.ConflictSHA256 != conflict.Operation.ConflictSHA256 || !contentTypedEqual(observed.Operation.ConflictPaths, conflict.Operation.ConflictPaths) {
				t.Error("retry replaced original conflict/abort evidence")
			}
			in.Action = "abort"
			result, abortErr := ReconcileTaskSource(f.dependencies, in)
			raw, readErr := os.ReadFile(filepath.Join(source, "README.md"))
			t.Logf("staged=%t initialSHA=%s observedSHA=%s abortState=%s err=%v retainedBytes=%q", staged, conflict.Operation.ConflictSHA256, observed.Operation.ConflictSHA256, result.State, abortErr, raw)
			if abortErr == nil || !strings.Contains(abortErr.Error(), "source_abort_edits_present") || readErr != nil || string(raw) != string(resolution) || runLocalGit(t, source, "diff", "--cached", "--binary") != indexBefore {
				t.Fatal("retry then abort discarded owner resolution or staged content")
			}
		})
	}
}

func TestParallelConflictRecoveryCannotInventAbortBaseline(t *testing.T) {
	f, prepared, in := sourceReconcileFixture(t, true)
	p, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	in.Apply, in.Confirmation = true, p.Confirmation
	f.dependencies.SourceReconciliationFault = func(stage, id string) error {
		if stage == "after_git" {
			return errInjected
		}
		return nil
	}
	if _, e = ReconcileTaskSource(f.dependencies, in); e == nil {
		t.Fatal("lost conflict publication not observed")
	}
	f.dependencies.SourceReconciliationFault = nil
	source := prepared.Preparation.Plan.WorktreePath
	resolution := []byte("owner resolution after lost conflict publication\n")
	if e = os.WriteFile(filepath.Join(source, "README.md"), resolution, 0600); e != nil {
		t.Fatal(e)
	}
	recovered, e := ReconcileTaskSource(f.dependencies, in)
	if e != nil || recovered.State != "source_conflict" {
		t.Fatalf("conflict observation: %+v %v", recovered, e)
	}
	in.Action, in.OperationID = "abort", recovered.Operation.ID
	result, e := ReconcileTaskSource(f.dependencies, in)
	raw, readErr := os.ReadFile(filepath.Join(source, "README.md"))
	t.Logf("recovered baseline=%s abortState=%s err=%v retainedBytes=%q", recovered.Operation.ConflictSHA256, result.State, e, raw)
	if e == nil || readErr != nil || string(raw) != string(resolution) {
		t.Fatal("recovery invented an abort baseline from owner edits and discarded them")
	}
}

func TestParallelSourceReconciliationPreservesIgnoredLocalContent(t *testing.T) {
	for _, movedParent := range []bool{false, true} {
		label := "same-intent-progresses"
		if movedParent {
			label = "changed-parent-refuses"
		}
		t.Run(label, func(t *testing.T) {
			f, prepared, in := sourceReconcileFixture(t, false)
			source := prepared.Preparation.Plan.WorktreePath
			parent := prepared.Preparation.Plan.Target.ParentLocator
			const name = "ignored-local.txt"
			if e := os.WriteFile(filepath.Join(source, ".gitignore"), []byte(name+"\n"), 0600); e != nil {
				t.Fatal(e)
			}
			runLocalGit(t, source, "add", ".gitignore")
			runLocalGit(t, source, "commit", "-m", "Ignore owned local content")
			local := []byte("exact local bytes\x00not parent content\n")
			if e := os.WriteFile(filepath.Join(source, name), local, 0600); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(parent, name), []byte("new tracked parent content\n"), 0600); e != nil {
				t.Fatal(e)
			}
			runLocalGit(t, parent, "add", name)
			runLocalGit(t, parent, "commit", "-m", "Parent tracks source ignored path")
			base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
			if e != nil {
				t.Fatal(e)
			}
			if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
				t.Fatal(e)
			}
			before := runLocalGit(t, source, "rev-parse", "HEAD")
			observed, e := f.dependencies.IntegrationGit.ObserveIntegrationWorktree(source, "refs/heads/"+prepared.Preparation.Plan.Branch)
			if e != nil || !observed.Clean {
				t.Fatalf("ignored-content fixture must pass ordinary clean status: %+v %v", observed, e)
			}
			git := f.dependencies.IntegrationGit.(*systemTaskIntegrationGit)
			runner, merges := git.run, 0
			git.run = func(args, env []string) GitCommandOutcome {
				for _, arg := range args {
					if arg == "merge" {
						merges++
					}
				}
				return runner(args, env)
			}
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply, in.Confirmation = true, preview.Confirmation
			result, e := ReconcileTaskSource(f.dependencies, in)
			raw, readErr := os.ReadFile(filepath.Join(source, name))
			after := runLocalGit(t, source, "rev-parse", "HEAD")
			t.Logf("isolated ignored-content merge state=%s err=%v merges=%d before=%s after=%s localSHA=%s observedSHA=%s", result.State, e, merges, before, after, digestTaskBytes(local), digestTaskBytes(raw))
			if e == nil || !strings.Contains(e.Error(), "source_ignored_collision") || result.State != "intent" || result.Operation == nil || merges != 0 || readErr != nil || string(raw) != string(local) || after != before {
				t.Fatalf("pre-effect collision lost content or claimed success: %+v %v", result, e)
			}
			marker := runLocalGit(t, source, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+result.Operation.ID)
			if _, e = os.Stat(marker); !os.IsNotExist(e) {
				t.Fatal("known collision published a Git operation marker")
			}
			if _, e = EffectiveTaskDeliveryBasis(mustQueueRegistry(t, f), in.PreparationID); e == nil {
				t.Fatal("blocked merge qualified an effective basis")
			}
			state := queueStateBytes(t, f)
			in.OperationID = result.Operation.ID
			retry, e := ReconcileTaskSource(f.dependencies, in)
			raw, readErr = os.ReadFile(filepath.Join(source, name))
			if e == nil || !strings.Contains(e.Error(), "source_ignored_collision") || retry.State != "intent" || merges != 0 || queueStateBytes(t, f) != state || readErr != nil || string(raw) != string(local) || runLocalGit(t, source, "rev-parse", "HEAD") != before {
				t.Fatalf("unchanged collision retry had an effect: %+v %v", retry, e)
			}
			preserved := filepath.Join(filepath.Dir(f.wrapper), "owner-preserved-local.txt")
			if e = os.Rename(filepath.Join(source, name), preserved); e != nil {
				t.Fatal(e)
			}
			if movedParent {
				runLocalGit(t, parent, "commit", "--allow-empty", "-m", "Parent moved after known refusal")
			}
			state = queueStateBytes(t, f)
			retry, e = ReconcileTaskSource(f.dependencies, in)
			if movedParent {
				if e == nil || !strings.Contains(e.Error(), "source_reconciliation_parent_changed") || retry.State != "intent" || merges != 0 || queueStateBytes(t, f) != state || runLocalGit(t, source, "rev-parse", "HEAD") != before {
					t.Fatalf("changed parent reached merge: %+v %v", retry, e)
				}
				if _, e = os.Stat(marker); !os.IsNotExist(e) {
					t.Fatal("changed parent published a merge marker")
				}
				in.Action = "abort"
				aborted, abortErr := ReconcileTaskSource(f.dependencies, in)
				if abortErr == nil || !strings.Contains(abortErr.Error(), "source_reconciliation_not_pending") || aborted.State != "intent" {
					t.Fatalf("unexpected unstarted-intent abort contract: %+v %v", aborted, abortErr)
				}
				in.Action = "cancel-unstarted"
				cancelled, cancelErr := ReconcileTaskSource(f.dependencies, in)
				if cancelErr != nil || cancelled.State != "aborted" || cancelled.Operation.ID != in.OperationID || merges != 0 || runLocalGit(t, source, "rev-parse", "HEAD") != before {
					t.Fatalf("exact unstarted intent cannot be cancelled without Git: %+v %v", cancelled, cancelErr)
				}
				base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
				if e != nil {
					t.Fatal(e)
				}
				if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
					t.Fatal(e)
				}
				fresh := in
				fresh.Apply, fresh.Action, fresh.OperationID, fresh.Confirmation = false, "", "", ""
				preview, e := ReconcileTaskSource(f.dependencies, fresh)
				if e != nil {
					t.Fatal(e)
				}
				fresh.Apply, fresh.Confirmation = true, preview.Confirmation
				progressed, e := ReconcileTaskSource(f.dependencies, fresh)
				if e != nil || progressed.State != "committed" || progressed.Operation.ID == in.OperationID || progressed.Operation.Plan.AttemptOrdinal != 2 || merges != 1 {
					t.Fatalf("fresh exact current-parent operation did not progress: %+v %v", progressed, e)
				}
				basis, e := EffectiveTaskDeliveryBasis(mustQueueRegistry(t, f), in.PreparationID)
				if e != nil || basis.Generation != 1 || basis.Parent.OID != runLocalGit(t, parent, "rev-parse", "HEAD") {
					t.Fatalf("fresh effective basis: %+v %v", basis, e)
				}
				if !contentTypedEqual(mustQueueRegistry(t, f).SourceReconciliations[0], *cancelled.Operation) {
					t.Fatal("fresh reconciliation replaced cancelled intent history")
				}
				t.Logf("cancelled=%s fresh=%s generation=%d parent=%s", cancelled.Operation.ID, progressed.Operation.ID, basis.Generation, basis.Parent.OID)
			} else if e != nil || retry.State != "committed" || retry.Operation.ID != in.OperationID || merges != 1 || len(mustQueueRegistry(t, f).SourceReconciliations) != 1 {
				t.Fatalf("same operation did not progress after owner preserved content: %+v %v merges=%d", retry, e, merges)
			}
			raw, e = os.ReadFile(preserved)
			if e != nil || string(raw) != string(local) {
				t.Fatal("owner-preserved content changed")
			}
			t.Logf("same-operation recheck movedParent=%t state=%s merges=%d preservedSHA=%s", movedParent, retry.State, merges, digestTaskBytes(raw))
		})
	}
}

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
