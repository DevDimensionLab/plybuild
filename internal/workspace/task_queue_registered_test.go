package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func registeredQueueBasis(t *testing.T, f workItemJourneyFixture) ReadViewBasisSnapshot {
	t.Helper()
	basis, err := ReadViewBasis(f.dependencies)
	if err != nil {
		t.Fatal(err)
	}
	return basis
}

func readRegisteredQueueFixture(f workItemJourneyFixture, basis ReadViewBasisSnapshot) RegisteredTaskQueues {
	d := f.dependencies
	// Panics would expose an accidental call into a detailed reader, Git,
	// locks, filesystem mutation, provider, or a fresh registry snapshot.
	d.WorkGit, d.IntegrationGit, d.Repos = nil, nil, nil
	d.Files, d.ProjectLocks, d.Projects, d.WorkItems = nil, nil, nil, nil
	d.WorkIDs, d.TaskContentIDs, d.WorkClock, d.HandoffEvidence = nil, nil, nil, nil
	return ReadRegisteredTaskQueues(d, basis.Workspace, basis.Projects, basis.Registry)
}

func createRegisteredQueueTask(t *testing.T, f workItemJourneyFixture, id TaskID) {
	t.Helper()
	_, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: id, Title: string(id), Description: "Synthetic queue reader fixture", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
}

func hasRegisteredQueueReason(reasons []QueueReason, code string) bool {
	for _, reason := range reasons {
		if reason.Code == code {
			return true
		}
	}
	return false
}

func TestRegisteredQueuesKeepExactRankAndLegacyReadinessWithoutLiveEffects(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	createRegisteredQueueTask(t, f, "legacy")
	createRegisteredQueueTask(t, f, "unselected")
	goal := goalFixture(t, f, "task", "goal")
	selection := queueSelect(t, f, "legacy", "legacy", nil)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "unselected"}, {TaskID: "task", Goal: &goal}, {TaskID: "legacy", Selection: &selection}})
	basis := registeredQueueBasis(t, f)
	before := queueStateBytes(t, f)
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage != "managed-read" {
			t.Fatalf("queue reader reached effect or nonmanaged content stage %q", stage)
		}
		return nil
	}
	for repeat := 0; repeat < 2; repeat++ {
		out := readRegisteredQueueFixture(f, basis)
		if out.Freshness != "fresh" || len(out.Diagnostics) != 0 || len(out.Queues) != 1 || len(out.Queues[0].Pending) != 3 {
			t.Fatalf("unexpected registered queue: %+v", out)
		}
		q := out.Queues[0]
		if q.QueueID != basis.Registry.TaskQueueEvents[0].QueueID || q.Revision != 1 || q.Target.ProjectID != "ply" || q.Target.RepoID != "ply" || q.Target.EpicID != "epic" || q.Target.ParentWorktreeID == "" || q.Current != nil {
			t.Fatalf("lost exact queue target: %+v", q)
		}
		if q.Pending[0].State != "blocked" || !hasRegisteredQueueReason(q.Pending[0].Reasons, "task_queue_selection_required") {
			t.Fatalf("unselected entry gained readiness: %+v", q.Pending[0])
		}
		ready := q.Pending[1]
		if ready.State != "ready" || ready.Rank != 2 || ready.Spec == nil || *ready.Spec != goal.Spec || ready.Goal == nil || *ready.Goal != goal || ready.SpecID == nil || *ready.SpecID != goal.SpecID || ready.SinceUTC == nil || *ready.SinceUTC != basis.Registry.TaskQueueEvents[0].RecordedAtUTC || len(ready.EvidenceIDs) < 3 {
			t.Fatalf("lost exact registered goal/rank/source: %+v", ready)
		}
		legacy := q.Pending[2]
		if legacy.State != "unknown" || legacy.Freshness != "fresh" || legacy.Spec == nil || legacy.Selection == nil || *legacy.Selection != selection || !hasRegisteredQueueReason(legacy.Reasons, "task_queue_live_preflight_required") {
			t.Fatalf("legacy entry claimed live readiness or lost provenance: %+v", legacy)
		}
	}
	if after := queueStateBytes(t, f); after != before {
		t.Fatal("registered read changed persisted metadata or Git state")
	}
}

func TestRegisteredQueuesCurrentPreparationBlocksNextGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	preview, prepared := goalPreparedFixture(t, f)
	createRegisteredQueueTask(t, f, "next")
	goal := goalFixture(t, f, "next", "next-goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "next", Goal: &goal}})
	out := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if len(out.Queues) != 1 || out.Queues[0].Current == nil || len(out.Queues[0].Pending) != 1 {
		t.Fatalf("lost current preparation or pending goal: %+v", out)
	}
	q, current := out.Queues[0], out.Queues[0].Current
	if current.TaskID != "task" || current.PreparationID != prepared.Preparation.Preparation.ID || current.State != "prepared" || current.Freshness != "fresh" || current.Goal == nil || *current.Goal != preview.Goal.Goal || current.Spec == nil || *current.Spec != prepared.Basis.Spec {
		t.Fatalf("current execution lost exact goal/Spec provenance: %+v; diagnostics=%+v", current, out.Diagnostics)
	}
	if q.Pending[0].State != "blocked" || !hasRegisteredQueueReason(q.Pending[0].Reasons, "task_queue_current") {
		t.Fatalf("current preparation allowed false readiness: %+v", q.Pending[0])
	}
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	out = readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if len(out.Queues) != 1 || out.Queues[0].Current == nil || len(out.Queues[0].Pending) != 0 {
		t.Fatalf("current-only queue disappeared: %+v", out)
	}
}

func TestRegisteredQueuesPendingExecutionRetainsOriginalGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	preview, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(f.wrapper)
	// Preserve the real execution publication/queue transition, stopping at
	// the boundary before PrepareTask reserves or creates the worktree.
	err = f.dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return f.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, err := session.Snapshot()
			if err != nil {
				return err
			}
			selection, err := materializeGoalExecution(f.dependencies, root, session, r, projects, preview, goalHumanFixture())
			if err != nil {
				return err
			}
			r, err = session.Snapshot()
			if err != nil {
				return err
			}
			return bindGoalExecutionQueue(f.dependencies, root, session, r, preview, selection)
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	out := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if len(out.Queues) != 1 || out.Queues[0].Current != nil || len(out.Queues[0].Pending) != 1 {
		t.Fatalf("interrupted execution did not retain its pending binding: %+v", out)
	}
	row := out.Queues[0].Pending[0]
	if row.Goal == nil || *row.Goal != goal || row.Selection == nil || row.Spec == nil || *row.Spec == goal.Spec || row.State != "unknown" || row.Freshness != "fresh" {
		t.Fatalf("pending execution lost original goal or invented live readiness: %+v", row)
	}
}

func TestRegisteredQueuesUnresolvedEffectsBlockGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	basis := registeredQueueBasis(t, f)
	target := QueueTarget{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
	for _, kind := range []string{"worktree", "integration", "base"} {
		t.Run(kind, func(t *testing.T) {
			b := basis
			switch kind {
			case "worktree":
				b.Registry.WorktreeOperations = []WorktreeOperationRecord{{ID: "preserved-effect", TaskID: "another", State: "reconciliation_required", ProjectID: "ply", RepoID: "ply", ParentEpicID: "epic"}}
			case "integration":
				b.Registry.IntegrationAuthorities = []IntegrationAuthority{{ID: "preserved-effect", Plan: WorkspaceTaskIntegrationPlan{Project: IntegrationPlanProject{ProjectID: "ply"}, Repository: IntegrationPlanRepository{RepoID: "ply"}, Epic: IntegrationPlanEpic{EpicID: "epic"}}}}
			case "base":
				b.Registry.EpicBaseUpdates = []EpicBaseUpdate{{ID: "preserved-effect", Phase: "intent", Plan: WorkspaceEpicBaseUpdatePlan{Target: target}}}
			}
			out := readRegisteredQueueFixture(f, b)
			q := out.Queues[0]
			if q.Pending[0].State != "blocked" || !hasRegisteredQueueReason(q.Pending[0].Reasons, "task_queue_effect_pending") || len(q.UnresolvedOperationIDs) != 1 || q.UnresolvedOperationIDs[0] != "preserved-effect" || len(out.Diagnostics) == 0 {
				t.Fatalf("unresolved %s effect allowed readiness: %+v", kind, out)
			}
		})
	}
}

func TestRegisteredQueuesProblemAndLegacySelectionDriftBlockReadiness(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(map[bool]string{false: "goal", true: "legacy"}[legacy], func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			if legacy {
				selection := queueSelect(t, f, "task", "selection", nil)
				queueSetFixture(t, f, "queue", []QueueEntry{{TaskID: "task", Selection: &selection}})
			} else {
				goal := goalFixture(t, f, "task", "goal")
				goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			}
			if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "changed-problem")); err != nil {
				t.Fatal(err)
			}
			out := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
			row := out.Queues[0].Pending[0]
			code := "task_goal_problem_changed"
			if legacy {
				code = "task_queue_selection_stale"
			}
			if row.State != "blocked" || !hasRegisteredQueueReason(row.Reasons, code) {
				t.Fatalf("changed Problem did not block %s: %+v", code, row)
			}
		})
	}
}

func TestRegisteredQueuesEpicAdvanceBlocksLegacyButKeepsBaseFreeGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	createRegisteredQueueTask(t, f, "legacy")
	goal := goalFixture(t, f, "task", "goal")
	selection := queueSelect(t, f, "legacy", "legacy", nil)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}, {TaskID: "legacy", Selection: &selection}})
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "advance registered Epic base")
	preview, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *preview.Confirmation}); err != nil {
		t.Fatal(err)
	}
	out := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if out.Queues[0].Pending[0].State != "ready" || out.Queues[0].Pending[0].ParentOID != nil {
		t.Fatalf("base-free goal acquired a stale Git basis: %+v", out.Queues[0].Pending[0])
	}
	legacy := out.Queues[0].Pending[1]
	if legacy.State != "blocked" || !hasRegisteredQueueReason(legacy.Reasons, "epic_base_changed") || legacy.BaseRevision == nil || *legacy.BaseRevision != 1 {
		t.Fatalf("registered stale legacy basis was hidden: %+v", legacy)
	}
}

func TestRegisteredQueuesInvalidGoalReferenceNeverBecomesReady(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	basis := registeredQueueBasis(t, f)
	goal.Spec.Revision++
	index := indexRegisteredQueues(basis.Workspace.Root, basis.Registry)
	row, err := registeredQueuePending(f.dependencies, basis.Workspace.Root, basis.Registry, index, QueueTarget{}, index.tasks["task"], TaskGoalQueueEntry{TaskID: "task", Goal: &goal}, 1)
	if err == nil || row.State != "unknown" || !hasRegisteredQueueReason(row.Reasons, "task_goal_invalid_reference") || row.Goal == nil || row.Goal.Spec != goal.Spec {
		t.Fatalf("unregistered exact goal revision gained readiness or lost evidence: %+v, %v", row, err)
	}
}

func TestRegisteredQueuesCorruptGoalKeepsValidSibling(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	createRegisteredQueueTask(t, f, "sibling")
	broken := goalFixture(t, f, "task", "broken-goal")
	valid := goalFixture(t, f, "sibling", "valid-goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &broken}, {TaskID: "sibling", Goal: &valid}})
	basis := registeredQueueBasis(t, f)
	path := taskContentPath(basis.Workspace.Root, "manifests", broken.Spec.ManifestSHA256)
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("corrupt synthetic goal"), 0400); err != nil {
		t.Fatal(err)
	}
	out := readRegisteredQueueFixture(f, basis)
	if out.Freshness != "unknown" || len(out.Diagnostics) != 1 || out.Diagnostics[0].TaskID == nil || *out.Diagnostics[0].TaskID != "task" || len(out.Queues[0].Pending) != 2 {
		t.Fatalf("corrupt goal was not isolated: %+v", out)
	}
	if out.Queues[0].Pending[0].State != "unknown" || out.Queues[0].Pending[1].State != "ready" || out.Queues[0].Pending[1].Rank != 2 {
		t.Fatalf("corrupt goal hid or altered valid sibling: %+v", out.Queues[0].Pending)
	}
}

func TestRegisteredQueuesDetectSourceChangeDuringRead(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	basis := registeredQueueBasis(t, f)
	reads := 0
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "managed-read" {
			reads++
			if reads == 2 {
				// The goal manifest has already been captured. Replace it while
				// the reader is gathering another source in the same request.
				path := taskContentPath(basis.Workspace.Root, "manifests", goal.Spec.ManifestSHA256)
				if err := os.Chmod(path, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("changed during registered read"), 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
		return nil
	}
	out := readRegisteredQueueFixture(f, basis)
	row := out.Queues[0].Pending[0]
	if row.State != "unknown" || row.Freshness != "stale" || !hasRegisteredQueueReason(row.Reasons, "task_queue_source_changed") || out.Freshness != "unknown" || len(out.Diagnostics) != 1 || out.Diagnostics[0].Code != "task_queue_source_changed" {
		t.Fatalf("source replacement claimed ready or fresh: %+v, %+v", row, out.Diagnostics)
	}
}

func TestRegisteredQueuesOrphanedQueueIsDiagnosticWithoutTask(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	basis := registeredQueueBasis(t, f)
	orphan := basis.Registry.TaskQueueEvents[0]
	orphan.QueueID = "que_" + strings.Repeat("a", 64)
	basis.Registry.TaskQueueEvents = append(basis.Registry.TaskQueueEvents, orphan)
	out := readRegisteredQueueFixture(f, basis)
	if len(out.Queues) != 1 || out.Queues[0].Pending[0].State != "ready" || len(out.Diagnostics) != 1 || out.Diagnostics[0].Code != "task_queue_orphaned" || out.Diagnostics[0].TaskID != nil || out.Diagnostics[0].Target != nil || out.Diagnostics[0].QueueID != orphan.QueueID {
		t.Fatalf("orphaned queue became a Task or hid valid queue: %+v", out)
	}
}

func TestRegisteredQueuesAbsenceIsFreshAndDoesNotReadContent(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	basis := registeredQueueBasis(t, f)
	f.dependencies.TaskContent = &TaskContentStorage{fault: func(stage string) error {
		t.Fatalf("queue-free workspace read content at %s", stage)
		return nil
	}}
	out := readRegisteredQueueFixture(f, basis)
	if out.Freshness != "fresh" || out.Queues == nil || len(out.Queues) != 0 || out.Diagnostics == nil || len(out.Diagnostics) != 0 {
		t.Fatalf("absent queue became unavailable: %+v", out)
	}
}

func TestRegisteredQueueContentCaptureBatchesSharedSources(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	createRegisteredQueueTask(t, f, "sibling")
	first := goalFixture(t, f, "task", "first-goal")
	second := goalFixture(t, f, "sibling", "second-goal")
	basis := registeredQueueBasis(t, f)
	physicalReads := 0
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "managed-read" {
			physicalReads++
		}
		return nil
	}
	capture := captureTaskContent(f.dependencies.TaskContent)
	d := f.dependencies
	d.TaskContent = capture.store
	for repeat := 0; repeat < 3; repeat++ {
		for task, goal := range map[TaskID]TaskGoalRef{"task": first, "sibling": second} {
			capture.owner = string(task)
			if _, err := loadTaskGoal(d, filepath.Dir(f.wrapper), basis.Registry, task, goal); err != nil {
				t.Fatal(err)
			}
		}
	}
	if physicalReads != len(capture.values) {
		t.Fatalf("repeated source reads were not captured: %d physical reads for %d sources", physicalReads, len(capture.values))
	}
	if changed := capture.changedOwners(); len(changed) != 0 || physicalReads != 2*len(capture.values) {
		t.Fatalf("expected one final check per unique source: %d reads, %d sources, changed=%v", physicalReads, len(capture.values), changed)
	}
}
