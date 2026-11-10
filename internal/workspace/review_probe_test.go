package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func TestReviewConfigAppearanceMustStaleView(t *testing.T) {
	root, _, d := durableFixture(t)
	if _, err := ListProjects(d); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, DefinitionsFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var defs workspacedef.Definitions
	if err = json.Unmarshal(raw, &defs); err != nil {
		t.Fatal(err)
	}
	defs.Projects[0].Epics = []workspacedef.Epic{}
	raw, err = json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(path); err != nil {
		t.Fatal(err)
	}
	before, err := ReadViewBasis(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(before.Projects.Projects) != 0 {
		t.Fatal("expected empty legacy registry")
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	after, err := ReadViewBasis(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Projects.Projects) != 1 {
		t.Fatal("config was not authoritative")
	}
	state, reasons := CheckReadViewBasis(d, before)
	if state == "fresh" {
		t.Fatalf("stale legacy view reported fresh after config appeared; before=%d projects after=%d state=%s reasons=%v", len(before.Projects.Projects), len(after.Projects.Projects), state, reasons)
	}
}
func TestReviewTaskRegistrationReadSurvivesMissingCheckout(t *testing.T) {
	root, repo, d := durableFixture(t)
	if _, err := CreateTask(d, TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "feature", ProjectID: "sample", RepoID: "code"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repo, repo+"-offline"); err != nil {
		t.Fatal(err)
	}
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	registry, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil || len(registry.Tasks) != 1 {
		t.Fatalf("registration missing: %v", err)
	}
	project, repoID := ProjectID("sample"), RepoID("code")
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
