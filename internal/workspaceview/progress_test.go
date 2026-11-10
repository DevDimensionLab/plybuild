package workspaceview

import (
	"encoding/json"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func progressFixture() *Snapshot {
	return &Snapshot{
		Workspace: workspace.WorkspaceObservation{Root: "/work"},
		Projects:  workspace.ProjectSnapshot{Projects: []workspace.ProjectRecord{{ID: "project"}}, Repos: []workspace.RepoRecord{{ID: "repo"}}},
		Registry:  workspace.WorkItemRegistry{Tasks: []workspace.TaskRecord{{ID: "task", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo", Title: "Feature", WorktreeState: workspace.WorkItemUnbound}}, Epics: []workspace.EpicRecord{{ID: "epic", ProjectID: "project", Title: "Epic"}}},
		Titles:    map[workspace.TaskID]workspace.TaskListTitle{"task": {Title: stringValue("Feature"), Source: "registration", Status: "available"}},
		Freshness: "fresh", Reasons: []string{},
	}
}

// fixtureAuthority is recorded integration progress without an outcome yet.
func fixtureAuthority() workspace.IntegrationAuthority {
	return workspace.IntegrationAuthority{ID: "authority", TaskID: "task", CreatedAtUTC: "2026-10-06T08:00:00Z"}
}

func TestProgressDoesNotTurnUntouchedOrPlannedTasksIntoNeedsYou(t *testing.T) {
	for _, planned := range []bool{false, true} {
		s := progressFixture()
		if planned {
			s.Registry.TaskProblemRevisions = []workspace.TaskProblemReference{{TaskID: "task"}}
			s.Registry.TaskSpecRevisions = []workspace.TaskSpecReference{{TaskID: "task"}}
		}
		rows, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
		if err != nil {
			t.Fatal(err)
		}
		progress := rows.Tasks[0].Progress
		if progress.State != "nothing" || progress.HasProgress || progress.IntegrationResultID != nil || len(progress.NextActions) != 0 {
			t.Fatalf("planned=%v produced false execution facts: %+v", planned, progress)
		}
		attention, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
		if err != nil || len(attention.Items) != 0 {
			t.Fatalf("attention=%+v err=%v", attention, err)
		}
	}
}

func TestProgressJSONUsesNullsAndEmptyArrays(t *testing.T) {
	s := progressFixture()
	p, err := BuildTaskProgress(s, "task", nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err = json.Unmarshal(b, &fields); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"last_activity_utc", "integration_result_id"} {
		if v, present := fields[name]; !present || v != nil {
			t.Fatalf("%s not explicit null: %s", name, b)
		}
	}
	if next, ok := fields["next_actions"].([]any); !ok || len(next) != 0 {
		t.Fatalf("next actions not []: %s", b)
	}
}
