package taskwork

import (
	"errors"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/repohistory"
)

// liveFixture has one live and one finished work item of each kind for
// repository app, each finished one recorded at an old Git location.
func liveFixture() (WorkItemRegistry, WorkItemLifecycles) {
	registry := emptyWorkItemRegistry()
	registry.Epics = []EpicRecord{
		{ID: "E-live", RepoBindings: []EpicRepoBinding{{RepoID: "app", GitCommonDir: "/new/app/.git"}}},
		{ID: "E-archived", RepoBindings: []EpicRepoBinding{{RepoID: "app", GitCommonDir: "/old/app/.git"}}},
	}
	for id, state := range map[TaskID]WorkItemState{"T-unbound": WorkItemUnbound, "T-creating": WorkItemCreating, "T-ready": WorkItemReady, "T-retired": WorkItemRetired, "T-reconcile": WorkItemReconciliationRequired} {
		dir := "/new/app/.git"
		if state == WorkItemRetired || state == WorkItemUnbound {
			dir = "/old/app/.git"
		}
		registry.Tasks = append(registry.Tasks, TaskRecord{ID: id, RepoID: "app", GitCommonDir: dir, WorktreeState: state})
	}
	registry.WorktreeOperations = []WorktreeOperationRecord{{ID: "O-live", TaskID: "T-ready", RepoID: "app"}, {ID: "O-finished", TaskID: "T-retired", RepoID: "app", State: "ready"}}
	registry.EpicCreateOperations = []EpicCreateOperation{{ID: "C-live", Request: EpicCreateInput{RepoID: "app"}, State: "pending"}, {ID: "C-completed", Request: EpicCreateInput{RepoID: "app"}, State: "completed"}}
	lifecycles := WorkItemLifecycles{{states: map[string]LifecycleState{"epic/E-archived": LifecycleArchived}}}
	return registry, lifecycles
}

// Only work that still holds a worktree or an open operation is live: an
// archived Epic, a Task without a worktree, an operation of such a Task and a
// completed Epic create are history.
func TestRepoUsesMarkOnlyLiveWorkLive(t *testing.T) {
	t.Parallel()
	registry, lifecycles := liveFixture()
	want := map[string]bool{"E-live": true, "E-archived": false, "T-unbound": false, "T-creating": true, "T-ready": true, "T-retired": false, "T-reconcile": true, "O-live": true, "O-finished": false, "C-live": true, "C-completed": false}
	uses := repoUses(registry, lifecycles)
	if len(uses) != len(want) {
		t.Fatalf("uses = %+v", uses)
	}
	for _, use := range uses {
		if live, ok := want[use.ID]; !ok || use.Live != live {
			t.Errorf("%s %s live = %v, want %v", use.Kind, use.ID, use.Live, live)
		}
	}
}

// The historical identity of a repository comes from live work only, so a
// moved repository keeps its ID while finished work names the old location.
func TestRegistryBindingsGiveLiveIdentityOnly(t *testing.T) {
	t.Parallel()
	registry, lifecycles := liveFixture()
	identities, err := repohistory.CommonDirs(registryBindings(registry, lifecycles))
	if err != nil || identities["app"] != "/new/app/.git" {
		t.Fatalf("identities = %v, %v", identities, err)
	}
	registry.Tasks = append(registry.Tasks, TaskRecord{ID: "T-other", RepoID: "app", GitCommonDir: "/other/app/.git", WorktreeState: WorkItemReady})
	var conflict *repohistory.ConflictError
	if _, err := repohistory.CommonDirs(registryBindings(registry, lifecycles)); !errors.As(err, &conflict) {
		t.Fatalf("two live locations: %v", err)
	}
}
