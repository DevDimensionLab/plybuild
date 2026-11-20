package taskwork

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"testing"
)

func TestReviewTaskRegistrationReadSurvivesMissingCheckout(t *testing.T) {
	t.Parallel()
	root, repo, d := durableFixture(t)
	if _, err := CreateTask(d, TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "feature", ProjectID: "sample", RepoID: "code"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repo, repo+"-offline"); err != nil {
		t.Fatal(err)
	}
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	registry, err := d.WorkItems.SnapshotRegistrations(testWork(t, root))
	if err != nil || len(registry.Tasks) != 1 {
		t.Fatalf("registration missing: %v", err)
	}
	project, repoID := workspace.ProjectID("sample"), workspace.RepoID("code")
	for _, filter := range []TaskListFilters{{ProjectID: &project}, {RepoID: &repoID}, {ProjectID: &project, RepoID: &repoID}} {
		if _, err := ListTasksWithFilters(d, filter); err != nil {
			t.Fatalf("filtered registration read failed: %v", err)
		}
	}
	if _, err := ReadViewBasis(d); err != nil {
		t.Fatalf("registration view failed: %v", err)
	}
	tasks, err := ListTasksWithFilters(d, TaskListFilters{})
	if err != nil {
		t.Fatalf("registered task cannot be listed while checkout unavailable: %v", err)
	}
	if len(tasks.Tasks) != 1 {
		t.Fatal("lost task")
	}
}
