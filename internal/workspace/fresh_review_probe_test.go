package workspace

import (
	"os"
	"testing"
)

func TestFreshReviewFilteredTaskListSurvivesMissingCheckout(t *testing.T) {
	root, repo, d := durableFixture(t)
	if _, err := CreateTask(d, TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "feature", ProjectID: "sample", RepoID: "code"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repo, repo+"-offline"); err != nil {
		t.Fatal(err)
	}
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	p, r, e := ProjectID("sample"), RepoID("code"), EpicID("feature")
	for _, c := range []struct {
		name    string
		filters TaskListFilters
	}{
		{"unfiltered", TaskListFilters{}},
		{"project", TaskListFilters{ProjectID: &p}},
		{"repo", TaskListFilters{RepoID: &r}},
		{"epic", TaskListFilters{EpicID: &e}},
		{"project-and-repo", TaskListFilters{ProjectID: &p, RepoID: &r}},
	} {
		t.Run(c.name, func(t *testing.T) {
			tasks, err := ListTasksWithFilters(d, c.filters)
			if err != nil {
				t.Fatalf("registration listing depends on unavailable checkout: %v", err)
			}
			if len(tasks.Tasks) != 1 {
				t.Fatalf("want one task, got %d", len(tasks.Tasks))
			}
		})
	}
}
func TestFreshReviewRegisteredViewSurvivesMissingCheckout(t *testing.T) {
	root, repo, d := durableFixture(t)
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repo, repo+"-offline"); err != nil {
		t.Fatal(err)
	}
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	if _, err := ReadViewBasis(d); err != nil {
		t.Fatalf("registration view depends on unavailable checkout: %v", err)
	}
}
