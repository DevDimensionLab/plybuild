package workspaceview

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type viewTestFiles struct {
	workspace.FileSystem
	root string
}

func (f viewTestFiles) Getwd() (string, error) { return f.root, nil }

type viewTestItems struct {
	registry                workspace.WorkItemRegistry
	reads, fullReads, locks int
}

func (s *viewTestItems) Snapshot(string) (workspace.WorkItemRegistry, error) {
	s.fullReads++
	return s.registry, nil
}
func (s *viewTestItems) SnapshotRegistrations(string) (workspace.WorkItemRegistry, error) {
	s.reads++
	return s.registry, nil
}
func (s *viewTestItems) WithLock(_ string, f func(workspace.WorkItemStoreSession) error) error {
	s.locks++
	return f(nil)
}

type viewTestProjects struct {
	projects workspace.ProjectSnapshot
	reads    int
}

func (s *viewTestProjects) Snapshot(string) ([]workspace.ProjectRecord, []workspace.RepoRecord, error) {
	s.reads++
	return s.projects.Projects, s.projects.Repos, nil
}
func (s *viewTestProjects) Add(string, workspace.ProjectRecord, []workspace.RepoRecord) (workspace.ProjectRecord, []workspace.RepoRecord, bool, error) {
	return workspace.ProjectRecord{}, nil, false, errors.New("unexpected write")
}

func viewDiskFixture(t *testing.T) (workspace.Dependencies, *viewTestItems, *viewTestProjects) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = viewTestFiles{d.Files, root}
	if _, err = workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	f := progressFixture()
	items := &viewTestItems{registry: f.Registry}
	projects := &viewTestProjects{projects: f.Projects}
	d.WorkItems, d.Projects = items, projects
	d.WorkGit, d.IntegrationGit, d.Repos, d.ProjectLocks = nil, nil, nil, nil
	return d, items, projects
}

func TestSnapshotUsesOneTypedReadAndDetectsLaterRegisterDrift(t *testing.T) {
	d, items, projects := viewDiskFixture(t)
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	if s.Freshness != "fresh" || items.reads != 1 || items.fullReads != 0 || projects.reads != 1 || items.locks != 0 {
		t.Fatalf("freshness=%s items=%+v projects=%+v", s.Freshness, items, projects)
	}
	if err := os.WriteFile(filepath.Join(s.Workspace.Root, ".ply", "work-items.yaml"), []byte("changed after snapshot"), 0600); err != nil {
		t.Fatal(err)
	}
	s.RefreshFreshness(d)
	if s.Freshness != "unknown" || len(s.Reasons) == 0 || items.reads != 1 || projects.reads != 1 || items.locks != 0 {
		t.Fatalf("drift not observed without a reload: %+v", s)
	}
}

func TestAttentionHonorsTaskAndEpicLifecycleWithoutChangingChildren(t *testing.T) {
	for _, subject := range []string{"task", "epic"} {
		for _, state := range []workspace.LifecycleState{workspace.LifecycleParked, workspace.LifecycleFrozen, workspace.LifecycleArchived} {
			t.Run(subject+"/"+string(state), func(t *testing.T) {
				d, items, _ := viewDiskFixture(t)
				items.registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
				beforeTask := items.registry.Tasks[0]
				if _, err := workspace.SetWorkItemLifecycle(d, workspace.WorkItemLifecycleInput{SubjectKind: subject, SubjectID: subject, State: state, ActorClaim: "fixture owner"}); err != nil {
					t.Fatal(err)
				}
				s, err := LoadSnapshot(d)
				if err != nil {
					t.Fatal(err)
				}
				rows, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
				if err != nil || len(rows.Tasks) != 1 || rows.Tasks[0].Progress.State != "in_progress" {
					t.Fatalf("task disappeared or progress changed: %+v %v", rows, err)
				}
				a, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
				if err != nil || len(a.Items) != 0 {
					t.Fatalf("inactive work asks for attention: %+v %v", a, err)
				}
				if subject == "epic" && rows.Tasks[0].Lifecycle != workspace.LifecycleActive {
					t.Fatalf("parent lifecycle leaked into stored Task state: %+v", rows.Tasks[0])
				}
				if !reflect.DeepEqual(beforeTask, items.registry.Tasks[0]) {
					t.Fatal("view mutated Task registration")
				}
			})
		}
	}
}

func TestProgressFiltersValidateKnownIDsBeforeIntersection(t *testing.T) {
	s := progressFixture()
	s.Projects.Projects = append(s.Projects.Projects, workspace.ProjectRecord{ID: "other"})
	p, r, e := workspace.ProjectID("other"), workspace.RepoID("repo"), workspace.EpicID("epic")
	rows, err := BuildTaskList(s, workspace.TaskListFilters{ProjectID: &p, RepoID: &r, EpicID: &e}, nil)
	if err != nil || len(rows.Tasks) != 0 {
		t.Fatalf("known empty intersection: %+v %v", rows, err)
	}
	missing := workspace.EpicID("missing")
	if _, err := BuildTaskList(s, workspace.TaskListFilters{ProjectID: &p, EpicID: &missing}, nil); err == nil {
		t.Fatal("unknown Epic hidden by empty intersection")
	}
	empty := workspace.ProjectID("")
	if _, err := ReadTaskList(workspace.Dependencies{}, workspace.TaskListFilters{ProjectID: &empty}); err == nil {
		t.Fatal("invalid filter was not rejected before reads")
	}
}
