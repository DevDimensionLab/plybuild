package workspace

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMachineReadbackUsesCanonicalVersionedShapes(t *testing.T) {
	epic := EpicRecord{ID: "epic", Title: "Epic", ProjectID: "ply", RepoBindings: []EpicRepoBinding{{RepoID: "ply", GitCommonDir: "/repo/.git", Worktree: EpicWorktreeBinding{ID: "wt_11111111111111111111111111111111", OwnerKind: "epic", OwnerID: "epic", Origin: "adopted", Locator: "/repo/epic", Ref: "refs/heads/epic", OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), GitCommonDir: "/repo/.git"}}}}
	readback := EpicReadbackResult{Workspace: "/workspace", Epic: epic, Observed: []EpicObservedRepo{{RepoID: "ply"}}, Freshness: []EpicFreshnessRepo{{RepoID: "ply", Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}}, Reasons: []string{"epic_observation_unknown"}}
	encoded, err := MarshalEpicReadback(readback)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if !strings.HasPrefix(text, `{"freshness":`) || !strings.Contains(text, `"kind":"WorkspaceEpicReadback@1"`) || strings.HasSuffix(text, "\n") {
		t.Fatalf("canonical readback = %s", text)
	}
}

func TestTaskReadbackSeparatesPersistedReadyFromFreshHandoffReadiness(t *testing.T) {
	result := TaskReadbackResult{Workspace: "/workspace", Task: TaskRecord{ID: "task", Title: "Task", Description: "Description", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: "/repo/.git", WorktreeState: WorkItemReady}, WorktreeReady: true, ReadyForHandoff: false, Reasons: []string{"target_observation_unknown"}, ProjectFreshness: "fresh", ParentFreshness: unknownWorktreeFreshness(), SourceFreshness: TaskFreshnessSource{Ref: "unknown", OID: "unknown", Tree: "unknown", CheckedOutAt: "unknown"}, TargetFreshness: unknownTargetFreshness(), Target: TaskObservedTarget{Kind: "unknown"}}
	encoded, err := MarshalTaskReadback(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if !strings.Contains(text, `"worktree_ready":true`) || !strings.Contains(text, `"ready_for_handoff":false`) || !strings.Contains(text, `"agent_started_by_this_command":false`) {
		t.Fatalf("task readback = %s", text)
	}
}

func TestReadbackMarksRemovedProjectMembershipStale(t *testing.T) {
	fixture := newWorkItemJourneyFixture(t)
	_, repositories, err := fixture.dependencies.Projects.Snapshot(filepath.Dir(fixture.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	fixture.dependencies.Projects = projectSnapshotOverride{ProjectStore: fixture.dependencies.Projects, repositories: repositories}

	epic, err := ShowEpic(fixture.dependencies, "epic")
	if err != nil || epic.Ready || !containsReason(epic.Reasons, "project_binding_stale") {
		t.Fatalf("Epic readback = %#v, %v", epic, err)
	}
	task, err := ShowTask(fixture.dependencies, "task")
	if err != nil || task.ReadyForHandoff || !containsReason(task.Reasons, "project_binding_stale") {
		t.Fatalf("Task readback = %#v, %v", task, err)
	}
}

type projectSnapshotOverride struct {
	ProjectStore
	projects     []ProjectRecord
	repositories []RepoRecord
}

func (store projectSnapshotOverride) Snapshot(string) ([]ProjectRecord, []RepoRecord, error) {
	return store.projects, store.repositories, nil
}

func containsReason(reasons []string, expected string) bool {
	for _, reason := range reasons {
		if reason == expected {
			return true
		}
	}
	return false
}
