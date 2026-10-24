package workitem

import "testing"

func TestLookupsKeepFirstSliceElement(t *testing.T) {
	epics := []EpicRecord{{ID: "other"}, {ID: "epic"}, {ID: "epic"}}
	tasks := []TaskRecord{{ID: "other"}, {ID: "task"}, {ID: "task"}}
	epic := EpicRecord{RepoBindings: []EpicRepoBinding{{RepoID: "other"}, {RepoID: "repo"}, {RepoID: "repo"}}}
	foundEpic, index := FindEpic(epics, "epic")
	if index != 1 || foundEpic != &epics[1] {
		t.Fatalf("Epic lookup = %p, %d", foundEpic, index)
	}
	foundEpic.Title = "updated"
	if epics[1].Title != "updated" || epics[2].Title != "" {
		t.Fatal("Epic lookup lost slice identity")
	}
	foundTask, index := FindTask(tasks, "task")
	if index != 1 || foundTask != &tasks[1] {
		t.Fatalf("Task lookup = %p, %d", foundTask, index)
	}
	foundTask.Title = "updated"
	if tasks[1].Title != "updated" || tasks[2].Title != "" {
		t.Fatal("Task lookup lost slice identity")
	}
	foundRepo, index := FindEpicRepo(epic, "repo")
	if index != 1 || foundRepo != &epic.RepoBindings[1] {
		t.Fatalf("repo lookup = %p, %d", foundRepo, index)
	}
	foundRepo.GitCommonDir = "updated"
	if epic.RepoBindings[1].GitCommonDir != "updated" || epic.RepoBindings[2].GitCommonDir != "" {
		t.Fatal("repo lookup lost slice identity")
	}
	for _, values := range [][]EpicRecord{nil, {}, epics} {
		if p, i := FindEpic(values, "missing"); p != nil || i != -1 {
			t.Fatalf("missing Epic = %v, %d", p, i)
		}
	}
	for _, values := range [][]TaskRecord{nil, {}, tasks} {
		if p, i := FindTask(values, "missing"); p != nil || i != -1 {
			t.Fatalf("missing Task = %v, %d", p, i)
		}
	}
	for _, values := range [][]EpicRepoBinding{nil, {}, epic.RepoBindings} {
		if p, i := FindEpicRepo(EpicRecord{RepoBindings: values}, "missing"); p != nil || i != -1 {
			t.Fatalf("missing repo = %v, %d", p, i)
		}
	}
}

func TestValidateTaskParent(t *testing.T) {
	task := TaskRecord{ID: "task", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo", GitCommonDir: "/repo/.git"}
	for _, tc := range []struct {
		name string
		epic *EpicRecord
		want string
	}{
		{"valid", &EpicRecord{ID: "epic", ProjectID: "project", RepoBindings: []EpicRepoBinding{{RepoID: "repo", GitCommonDir: "/repo/.git"}}}, ""},
		{"missing", nil, "Task task has invalid Epic/Project reference"},
		{"wrong Epic", &EpicRecord{ID: "other", ProjectID: "project"}, "Task task has invalid Epic/Project reference"},
		{"wrong Project", &EpicRecord{ID: "epic", ProjectID: "other"}, "Task task has invalid Epic/Project reference"},
		{"missing repo", &EpicRecord{ID: "epic", ProjectID: "project"}, "Task task has invalid repository binding"},
		{"wrong common dir", &EpicRecord{ID: "epic", ProjectID: "project", RepoBindings: []EpicRepoBinding{{RepoID: "repo", GitCommonDir: "/other/.git"}}}, "Task task has invalid repository binding"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateTaskParent(task, tc.epic)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || err.Error() != tc.want {
				t.Fatalf("error = %v, want %s", err, tc.want)
			}
		})
	}
}

func TestStateWireValues(t *testing.T) {
	got := []WorkItemState{WorkItemUnbound, WorkItemCreating, WorkItemReady, WorkItemRetired, WorkItemReconciliationRequired}
	want := []string{"unbound", "creating", "worktree_ready", "retired", "reconciliation_required"}
	for i := range want {
		if string(got[i]) != want[i] {
			t.Fatalf("state %d = %q", i, got[i])
		}
	}
}
