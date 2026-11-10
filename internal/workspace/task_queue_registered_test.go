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
	d.WorkIDs, d.TaskContentIDs, d.WorkClock = nil, nil, nil
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

func TestRegisteredQueuesEpicAdvanceKeepsBaseFreeGoal(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
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
