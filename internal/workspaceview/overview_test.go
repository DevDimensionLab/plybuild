package workspaceview

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestOverviewUsesOneSnapshotAndMatchesItsComponentViews(t *testing.T) {
	d, items, projects := viewDiskFixture(t)
	items.registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
	root, err := d.Files.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	before := viewTree(t, root)
	overview, err := ReadOverview(d, EpicFilters{})
	if err != nil {
		t.Fatal(err)
	}
	if items.reads != 1 || items.fullReads != 0 || items.locks != 0 || projects.reads != 1 {
		t.Fatalf("overview did not share one registry read: items=%+v projects=%+v", items, projects)
	}
	if !reflect.DeepEqual(before, viewTree(t, root)) {
		t.Fatal("overview wrote to the workspace")
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	process, err := ReadProcessFacts(d, s)
	if err != nil {
		t.Fatal(err)
	}
	firstBatch := s.process
	if _, err := ReadProcessFacts(d, s); err != nil || firstBatch == nil || s.process != firstBatch {
		t.Fatal("process batch was not reused")
	}
	s.RefreshFreshness(d)
	tasks, err := BuildTaskList(s, workspace.TaskListFilters{}, process)
	if err != nil {
		t.Fatal(err)
	}
	epics, err := BuildEpics(s, EpicFilters{}, process)
	if err != nil {
		t.Fatal(err)
	}
	attention, err := BuildAttention(s, workspace.TaskListFilters{}, process)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Kind != OverviewKind || !reflect.DeepEqual(overview.Tasks, tasks.Tasks) || !reflect.DeepEqual(overview.Epics, epics.Epics) || !reflect.DeepEqual(overview.Attention, attention.Items) {
		t.Fatalf("overview differs from its component views: %+v", overview)
	}
	if len(overview.Tasks) != 1 || len(overview.Epics) != 1 || len(overview.Attention) != 1 || overview.Epics[0].TaskCount != 1 || overview.Epics[0].ProgressCounts.InProgress != 1 {
		t.Fatalf("inconsistent aggregates: %+v", overview)
	}
}

func TestOverviewFiltersKeepEmptyEpicsAndEmptyCollections(t *testing.T) {
	d, items, projects := viewDiskFixture(t)
	projects.projects.Projects = append(projects.projects.Projects, workspace.ProjectRecord{ID: "other"})
	items.registry.Epics = append(items.registry.Epics, workspace.EpicRecord{ID: "empty", ProjectID: "other", Title: "Empty"})
	id := workspace.ProjectID("other")
	out, err := ReadOverview(d, EpicFilters{ProjectID: &id})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Tasks) != 0 || out.Tasks == nil || len(out.Attention) != 0 || out.Attention == nil || len(out.Epics) != 1 || out.Epics[0].EpicID != "empty" {
		t.Fatalf("filter or empty Epic lost: %+v", out)
	}
	id = "missing"
	if _, err := ReadOverview(d, EpicFilters{ProjectID: &id}); err == nil {
		t.Fatal("unknown project accepted")
	}
}

func viewTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[path] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
