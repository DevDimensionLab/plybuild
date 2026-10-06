package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

type countingViewWorkItems struct {
	WorkItemStore
	snapshots, registrations, locks int
}

func (s *countingViewWorkItems) Snapshot(root string) (WorkItemRegistry, error) {
	s.snapshots++
	return s.WorkItemStore.Snapshot(root)
}
func (s *countingViewWorkItems) SnapshotRegistrations(root string) (WorkItemRegistry, error) {
	s.registrations++
	return s.WorkItemStore.SnapshotRegistrations(root)
}
func (s *countingViewWorkItems) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	s.locks++
	return s.WorkItemStore.WithLock(root, fn)
}

type countingViewProjects struct {
	ProjectStore
	snapshots int
}

func (s *countingViewProjects) Snapshot(root string) ([]ProjectRecord, []RepoRecord, error) {
	s.snapshots++
	return s.ProjectStore.Snapshot(root)
}

func TestReadViewBasisReadsEachRegistryOnceWithoutGitOrLocks(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	items := &countingViewWorkItems{WorkItemStore: f.dependencies.WorkItems}
	projects := &countingViewProjects{ProjectStore: f.dependencies.Projects}
	d := f.dependencies
	d.WorkItems, d.Projects = items, projects
	d.WorkGit, d.IntegrationGit, d.ProjectLocks, d.Repos = nil, nil, nil, nil
	b, err := ReadViewBasis(d)
	if err != nil {
		t.Fatal(err)
	}
	if items.registrations != 1 || items.snapshots != 0 || items.locks != 0 || projects.snapshots != 1 {
		t.Fatalf("unexpected reads or effects: items=%+v projects=%+v", items, projects)
	}
	if len(b.Registry.Tasks) != 1 || len(b.Projects.Projects) != 1 || b.Titles["task"].Title == nil {
		t.Fatalf("lost registry facts: %+v", b)
	}
	if fresh, reasons := CheckReadViewBasis(d, b); fresh != "fresh" || len(reasons) != 0 {
		t.Fatalf("stable source marked %s: %v", fresh, reasons)
	}
}

func TestReadViewProgressDoesNotMistakePlansForWork(t *testing.T) {
	r := WorkItemRegistry{Tasks: []TaskRecord{{ID: "untouched", WorktreeState: WorkItemUnbound}, {ID: "planned", WorktreeState: WorkItemUnbound}}, TaskProblemRevisions: []TaskProblemReference{{TaskID: "planned"}}, TaskSpecRevisions: []TaskSpecReference{{TaskID: "planned"}}}
	for id, fact := range RecordedTaskProgressFacts(r) {
		if fact.HasProgress || fact.ResultCount != 0 || fact.TaskResult != nil || fact.HumanQA != nil || fact.IntegrationResult != nil {
			t.Fatalf("%s planning became implementation progress: %+v", id, fact)
		}
	}
}

func TestReadViewRejectsSymlinkSourcesAndDetectsReplacedFiles(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	b, err := ReadViewBasis(f.dependencies)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(b.Workspace.Root, MarkerDirectory, ProjectsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "same-project-bytes")
	if err = os.WriteFile(target, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	if fresh, _ := CheckReadViewBasis(f.dependencies, b); fresh != "unknown" {
		t.Fatal("symlink replacement with the same bytes was reported fresh")
	}
	if _, err = ReadViewBasis(f.dependencies); err == nil {
		t.Fatal("project source symlink accepted")
	}
}
