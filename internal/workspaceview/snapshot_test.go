package workspaceview

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

type viewTestFiles struct {
	workspace.FileSystem
	root string
}

func (f viewTestFiles) Getwd() (string, error) { return f.root, nil }

type viewTestItems struct {
	registry                taskwork.WorkItemRegistry
	reads, fullReads, locks int
}

func (s *viewTestItems) Snapshot(work workspace.ProjectWork) (taskwork.WorkItemRegistry, error) {
	s.fullReads++
	return s.of(work.Project), nil
}
func (s *viewTestItems) SnapshotRegistrations(work workspace.ProjectWork) (taskwork.WorkItemRegistry, error) {
	s.reads++
	return s.of(work.Project), nil
}

// of is the registry of one Project: its Epics, Tasks and Epic bases, and
// for the fixture's Project "project" also every other record.
func (s *viewTestItems) of(project workspace.ProjectID) taskwork.WorkItemRegistry {
	r := taskwork.WorkItemRegistry{}
	if project == "project" {
		r = s.registry
	}
	r.Epics, r.Tasks, r.EpicBaseVersions = nil, nil, nil
	for _, epic := range s.registry.Epics {
		if epic.ProjectID == project {
			r.Epics = append(r.Epics, epic)
		}
	}
	for _, task := range s.registry.Tasks {
		if task.ProjectID == project {
			r.Tasks = append(r.Tasks, task)
		}
	}
	for _, base := range s.registry.EpicBaseVersions {
		if base.ProjectID == project {
			r.EpicBaseVersions = append(r.EpicBaseVersions, base)
		}
	}
	return r
}
func (s *viewTestItems) WithLock(_ workspace.ProjectWork, f func(taskwork.WorkItemStoreSession) error) error {
	s.locks++
	return f(nil)
}

// viewTestLocks lets a lifecycle write take the definitions lock in a
// fixture whose read views must not lock.
type viewTestLocks struct{}

func (viewTestLocks) WithSnapshotLock(_ string, f func(workspace.ProjectSnapshot) error) error {
	return f(workspace.ProjectSnapshot{})
}

func withViewLocks(d taskwork.Dependencies) taskwork.Dependencies {
	d.ProjectLocks = viewTestLocks{}
	return d
}

type viewTestProjects struct {
	projects workspace.ProjectSnapshot
	reads    int
}

func (s *viewTestProjects) Snapshot(string) ([]workspace.ProjectRecord, []workspace.RepoRecord, error) {
	s.reads++
	return s.projects.Projects, s.projects.Repos, nil
}

func viewDiskFixture(t *testing.T) (taskwork.Dependencies, *viewTestItems, *viewTestProjects) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := taskwork.SystemDependencies()
	d.Files = viewTestFiles{d.Files, root}
	if _, err = workspace.Init(d.Dependencies); err != nil {
		t.Fatal(err)
	}
	f := progressFixture()
	items := &viewTestItems{registry: f.Registry}
	projects := &viewTestProjects{projects: workspace.ProjectSnapshot{}}
	defineViewEpic(t, root, items, projects, "project", "repo", "epic", "Epic")
	d.WorkItems, d.Projects = items, projects
	d.WorkGit, d.IntegrationGit, d.Repos, d.ProjectLocks = nil, nil, nil, nil
	return d, items, projects
}

// defineViewEpic adds a synthetic Epic with a worktree directory and a
// matching binding to the registry and Project store fakes; its Project and
// repository are defined in ply-workspace.json when missing.
func defineViewEpic(t *testing.T, root string, items *viewTestItems, projects *viewTestProjects, project workspace.ProjectID, repo workspace.RepoID, epic taskwork.EpicID, title string) {
	t.Helper()
	repoPath := filepath.Join(string(project), string(repo))
	epicPath := filepath.Join(string(project), string(epic))
	for _, dir := range []string{repoPath, epicPath} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	ref := "refs/heads/" + string(epic)
	common := filepath.Join(root, repoPath, ".git")
	_, basis, err := workspacedef.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	mutate := func(defs *workspacedef.Definitions) (bool, error) {
		if _, err := workspacedef.AddRepo(defs, string(project), string(repo), repoPath, ""); err != nil {
			defs.Projects = append(defs.Projects, workspacedef.Project{ID: string(project), Name: string(project), Wrapper: string(project), Repos: []workspacedef.Repo{{ID: string(repo), Path: repoPath}}})
		}
		return true, nil
	}
	if _, err := workspacedef.Update(root, basis, mutate, func(workspacedef.Definitions) error { return nil }); err != nil {
		t.Fatal(err)
	}
	binding := []taskwork.EpicRepoBinding{{RepoID: repo, GitCommonDir: common, Worktree: taskwork.EpicWorktreeBinding{Locator: filepath.Join(root, epicPath), Ref: ref}}}
	found := false
	for i := range items.registry.Epics {
		if items.registry.Epics[i].ID == epic {
			items.registry.Epics[i].RepoBindings, found = binding, true
		}
	}
	if !found {
		items.registry.Epics = append(items.registry.Epics, taskwork.EpicRecord{ID: epic, ProjectID: project, Title: title, RepoBindings: binding})
	}
	// Every registered Epic binding has its adopted base, as in a stored registry.
	hasBase := false
	for _, v := range items.registry.EpicBaseVersions {
		hasBase = hasBase || v.EpicID == epic && v.RepoID == repo
	}
	if !hasBase {
		items.registry.EpicBaseVersions = append(items.registry.EpicBaseVersions, taskwork.EpicBaseVersion{ProjectID: project, RepoID: repo, EpicID: epic, Revision: 1, Locator: binding[0].Worktree.Locator, Ref: ref, GitCommonDir: common})
	}
	snapshot := &projects.projects
	if !containsProject(snapshot.Projects, project) {
		snapshot.Projects = append(snapshot.Projects, workspace.ProjectRecord{ID: project, Name: string(project), Wrapper: filepath.Join(root, string(project)), RepoIDs: []workspace.RepoID{repo}})
		snapshot.Repos = append(snapshot.Repos, workspace.RepoRecord{ID: repo, Locator: filepath.Join(root, repoPath), GitCommonDir: common})
	}
}

func containsProject(projects []workspace.ProjectRecord, id workspace.ProjectID) bool {
	for _, p := range projects {
		if p.ID == id {
			return true
		}
	}
	return false
}

func TestSnapshotUsesOneTypedReadAndDetectsLaterRegisterDrift(t *testing.T) {
	d, items, projects := viewDiskFixture(t)
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	if s.Freshness != "fresh" || items.reads != 1 || items.fullReads != 0 || projects.reads != 0 || items.locks != 0 {
		t.Fatalf("freshness=%s items=%+v projects=%+v", s.Freshness, items, projects)
	}
	works, err := workspace.DefinedWork(s.Workspace.Root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(works[0].Registry()), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(works[0].Registry(), []byte("changed after snapshot"), 0600); err != nil {
		t.Fatal(err)
	}
	s.RefreshFreshness(d)
	if s.Freshness != "unknown" || len(s.Reasons) == 0 || items.reads != 1 || projects.reads != 0 || items.locks != 0 {
		t.Fatalf("drift not observed without a reload: %+v", s)
	}
}

func TestAttentionHonorsTaskAndEpicLifecycleWithoutChangingChildren(t *testing.T) {
	for _, subject := range []string{"task", "epic"} {
		for _, state := range []taskwork.LifecycleState{taskwork.LifecycleParked, taskwork.LifecycleFrozen, taskwork.LifecycleArchived} {
			t.Run(subject+"/"+string(state), func(t *testing.T) {
				d, items, _ := viewDiskFixture(t)
				items.registry.IntegrationAuthorities = []taskwork.IntegrationAuthority{fixtureAuthority()}
				beforeTask := items.registry.Tasks[0]
				if _, err := taskwork.SetWorkItemLifecycle(withViewLocks(d), taskwork.WorkItemLifecycleInput{SubjectKind: subject, SubjectID: subject, State: state, ActorClaim: "fixture owner"}); err != nil {
					t.Fatal(err)
				}
				s, err := LoadSnapshot(d)
				if err != nil {
					t.Fatal(err)
				}
				rows, err := BuildTaskList(s, taskwork.TaskListFilters{}, nil)
				if err != nil || len(rows.Tasks) != 1 || rows.Tasks[0].Progress.State != "in_progress" {
					t.Fatalf("task disappeared or progress changed: %+v %v", rows, err)
				}
				a, err := BuildAttention(s, taskwork.TaskListFilters{}, nil)
				if err != nil || len(a.Items) != 0 {
					t.Fatalf("inactive work asks for attention: %+v %v", a, err)
				}
				if subject == "epic" && rows.Tasks[0].Lifecycle != taskwork.LifecycleActive {
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
	p, r, e := workspace.ProjectID("other"), workspace.RepoID("repo"), taskwork.EpicID("epic")
	rows, err := BuildTaskList(s, taskwork.TaskListFilters{ProjectID: &p, RepoID: &r, EpicID: &e}, nil)
	if err != nil || len(rows.Tasks) != 0 {
		t.Fatalf("known empty intersection: %+v %v", rows, err)
	}
	missing := taskwork.EpicID("missing")
	if _, err := BuildTaskList(s, taskwork.TaskListFilters{ProjectID: &p, EpicID: &missing}, nil); err == nil {
		t.Fatal("unknown Epic hidden by empty intersection")
	} else {
		var workError *workspace.WorkItemError
		if !errors.As(err, &workError) || workError.Class != workspace.ErrorWorkNotFound || err.Error() != "workspace_work_not_found: Epic missing is not registered" || errors.Unwrap(err) != nil {
			t.Fatalf("Epic filter error changed: %#v", err)
		}
	}
	empty := workspace.ProjectID("")
	if _, err := ReadTaskList(taskwork.Dependencies{}, taskwork.TaskListFilters{ProjectID: &empty}); err == nil {
		t.Fatal("invalid filter was not rejected before reads")
	}
}
