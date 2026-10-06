package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRegisteredTaskListTitleFactsAndReadOnlyDependencies(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	revised, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "fixture/read-title"))
	if err != nil {
		t.Fatal(err)
	}
	literal, err := ParseTaskCreateInput("literal", "Problem content unavailable", "Legitimate title", "epic", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTask(f.dependencies, literal); err != nil {
		t.Fatal(err)
	}
	// Seed a validated legacy registration without a problem head alongside modern Tasks.
	err = f.dependencies.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, e := s.Snapshot()
		if e != nil {
			return e
		}
		task := r.Tasks[0]
		task.ID = "legacy"
		task.Title = "Legacy: æøå (registered)"
		r.Tasks = append(r.Tasks, task)
		r.TaskSpecPolicies = append(r.TaskSpecPolicies, TaskSpecPolicy{TaskID: "legacy", Mode: "legacy", LegacyResultIDs: []TaskResultID{}})
		sortWorkRegistry(&r)
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	// Nil effect/observation dependencies fail immediately if reads begin to use them.
	f.dependencies.Repos = nil
	f.dependencies.WorkGit = nil
	f.dependencies.IntegrationGit = nil
	f.dependencies.ProjectLocks = nil
	for _, fault := range []bool{false, true} {
		if fault {
			f.dependencies.TaskContent = &TaskContentStorage{fault: func(point string) error {
				if point == "managed-read" {
					return errors.New("read denied")
				}
				return nil
			}}
		}
		result, err := ListTasksWithFilters(f.dependencies, TaskListFilters{})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Tasks) != 3 {
			t.Fatalf("lost registration: %+v", result)
		}
		legacy := result.Titles["legacy"]
		if legacy.Title == nil || *legacy.Title != "Legacy: æøå (registered)" || legacy.Source != "registration" || legacy.Status != "available" {
			t.Fatalf("legacy=%+v", legacy)
		}
		current := result.Titles["task"]
		sentinel := result.Titles["literal"]
		if !fault {
			if current.Title == nil || *current.Title != "Revised fixture problem" || current.Source != "problem_revision" || current.Status != "available" {
				t.Fatalf("current=%+v", current)
			}
			if sentinel.Title == nil || *sentinel.Title != "Problem content unavailable" || sentinel.Status != "available" {
				t.Fatalf("literal=%+v", sentinel)
			}
		} else {
			for _, title := range []TaskListTitle{current, sentinel} {
				if title.Title != nil || title.Source != "problem_revision" || title.Status != "unavailable" {
					t.Fatalf("unavailable=%+v", title)
				}
			}
			if result.CurrentTitles["task"] != "Problem content unavailable" {
				t.Fatal("legacy presentation changed")
			}
		}
		b, err := MarshalTaskList(result)
		if err != nil {
			t.Fatal(err)
		}
		var value map[string]any
		if err = json.Unmarshal(b, &value); err != nil {
			t.Fatal(err)
		}
		rows := value["tasks"].([]any)
		if rows[0].(map[string]any)["task_id"] != "legacy" || rows[1].(map[string]any)["task_id"] != "literal" || rows[2].(map[string]any)["task_id"] != "task" {
			t.Fatalf("sort=%s", b)
		}
	}
	if revised.OutcomeRef == nil {
		t.Fatal("fixture did not revise current problem")
	}
}

func TestRegisteredTaskListChecksAllFiltersBeforeIntersection(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	if _, _, _, err := f.dependencies.Projects.Add(root, ProjectRecord{ID: "other", Name: "Other", Wrapper: "/missing/ø (wrapper)", RepoIDs: []RepoID{"elsewhere"}}, []RepoRecord{{ID: "elsewhere", Locator: "/missing/ø (checkout)", GitCommonDir: "/missing/common.git"}}); err != nil {
		t.Fatal(err)
	}
	p, other, r, elsewhere, e := ProjectID("ply"), ProjectID("other"), RepoID("ply"), RepoID("elsewhere"), EpicID("epic")
	f.dependencies.Repos = nil
	f.dependencies.WorkGit = nil
	f.dependencies.ProjectLocks = nil
	for _, tc := range []struct {
		f     TaskListFilters
		count int
	}{{TaskListFilters{}, 1}, {TaskListFilters{ProjectID: &p}, 1}, {TaskListFilters{RepoID: &r}, 1}, {TaskListFilters{EpicID: &e}, 1}, {TaskListFilters{ProjectID: &other, RepoID: &r}, 0}, {TaskListFilters{ProjectID: &p, RepoID: &elsewhere, EpicID: &e}, 0}} {
		result, err := ListTasksWithFilters(f.dependencies, tc.f)
		if err != nil || len(result.Tasks) != tc.count {
			t.Fatalf("filters=%+v result=%+v err=%v", tc.f, result, err)
		}
	}
	missingP, missingR, missingE := ProjectID("unknown"), RepoID("unknown"), EpicID("unknown")
	for _, filter := range []TaskListFilters{{ProjectID: &missingP}, {RepoID: &missingR}, {EpicID: &missingE}, {ProjectID: &other, RepoID: &r, EpicID: &missingE}, {ProjectID: &p, RepoID: &missingR, EpicID: &e}} {
		if _, err := ListTasksWithFilters(f.dependencies, filter); !hasWorkErrorClass(err, ErrorWorkNotFound) {
			t.Fatalf("filters=%+v err=%v", filter, err)
		}
	}
	// Invalid IDs must fail even when no filesystem or registry dependencies exist.
	emptyP, badR, emptyE := ProjectID(""), RepoID("BAD"), EpicID("")
	for _, filter := range []TaskListFilters{{ProjectID: &emptyP}, {RepoID: &badR}, {EpicID: &emptyE}} {
		if _, err := ListTasksWithFilters(Dependencies{}, filter); !hasWorkErrorClass(err, ErrorWorkInvalidArguments) {
			t.Fatalf("filter=%+v err=%v", filter, err)
		}
	}
}

func TestCoreReadSerializersPreserveStoredPathsSortAndEmptyArrays(t *testing.T) {
	projects := []ProjectRecord{{ID: "z", Name: "Z: æ (two)", Wrapper: "/abs/missing wrapper", RepoIDs: []RepoID{"b", "a"}}, {ID: "a", Name: "A", Wrapper: "/abs/a", RepoIDs: []RepoID{}}}
	original := append([]ProjectRecord(nil), projects...)
	b, err := MarshalProjectList(ProjectListResult{Workspace: "/root", Projects: projects})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(projects, original) {
		t.Fatal("serializer reordered input")
	}
	var value map[string]any
	if err = json.Unmarshal(b, &value); err != nil {
		t.Fatal(err)
	}
	rows := value["projects"].([]any)
	if rows[0].(map[string]any)["project_id"] != "a" || rows[1].(map[string]any)["repo_count"] != float64(2) {
		t.Fatalf("projects=%s", b)
	}
	b, err = MarshalProject(ProjectResult{Workspace: "/root", Project: projects[0], Repos: []RepoRecord{{ID: "b", Locator: "/missing B", GitCommonDir: "/B.git"}, {ID: "a", Locator: "/missing A", GitCommonDir: "/A.git"}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &value); err != nil {
		t.Fatal(err)
	}
	repoRows := value["repositories"].([]any)
	firstRepo := repoRows[0].(map[string]any)
	if firstRepo["repo_id"] != "a" || firstRepo["locator"] != "/missing A" || firstRepo["git_common_dir"] != "/A.git" || firstRepo["wrapper"] != nil {
		t.Fatalf("repos=%s", b)
	}
	for _, tc := range []struct {
		marshal func() ([]byte, error)
		key     string
	}{
		{func() ([]byte, error) { return MarshalProjectList(ProjectListResult{}) }, "projects"},
		{func() ([]byte, error) { return MarshalProject(ProjectResult{}) }, "repositories"},
		{func() ([]byte, error) { return MarshalTaskList(TaskListResult{}) }, "tasks"},
	} {
		b, err := tc.marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"`+tc.key+`":[]`) {
			t.Fatalf("empty=%s", b)
		}
	}

}

func TestCoreReadsDoNotHideCorruptRegistries(t *testing.T) {
	root := createProjectWorkspace(t)
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	if err := os.WriteFile(workItemsPath(root), []byte("not: [valid yaml"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ListTasks(d, nil); err == nil {
		t.Fatal("corrupt work registry hidden")
	}
}

func TestRegisteredV4TaskListKeepsUnavailableTitleWhileReadyRemainsStrict(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "list-v4", []QueueEntry{})
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if r.FormatVersion != 4 {
		t.Fatal("fixture requires format 4")
	}
	head := taskContentState(r, "task").ProblemHead
	// Corrupt the immutable manifest after constructing a fully valid v4 registry.
	path := taskContentPath(root, "manifests", head.ManifestSHA256)
	if err = os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("corrupt content"), 0400); err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(path, 0400); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	list, err := ListTasksWithFilters(f.dependencies, TaskListFilters{})
	if err != nil || len(list.Tasks) != 1 || list.Titles["task"].Title != nil || list.Titles["task"].Status != "unavailable" {
		t.Fatalf("list=%+v err=%v", list, err)
	}
	if _, err = f.dependencies.WorkItems.Snapshot(root); err == nil {
		t.Fatal("strict registry read lost closure check")
	}
	if _, err = ListTaskQueue(f.dependencies, queueTargetFixture(), true); err == nil {
		t.Fatal("ready swallowed content corruption")
	}
	if _, err = ShowTask(f.dependencies, "task"); err == nil {
		t.Fatal("Task show swallowed content corruption")
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil || string(before) != string(after) {
		t.Fatal("read changed registry")
	}
	if err = os.WriteFile(workItemsPath(root), []byte("format_version: 4\ntasks: [broken"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = ListTasksWithFilters(f.dependencies, TaskListFilters{}); err == nil {
		t.Fatal("registered list ignored corrupt registry")
	}
}

func TestRegisteredListRejectsUnverifiableProblemPublication(t *testing.T) {
	for _, v4 := range []bool{false, true} {
		f := newWorkItemJourneyFixture(t)
		if v4 {
			queueSetFixture(t, f, "binding-test", []QueueEntry{})
		}
		root := filepath.Dir(f.wrapper)
		r, err := f.dependencies.WorkItems.Snapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		head := taskContentState(r, "task").ProblemHead
		changedKey := "fixture/changed-registration-key"
		for i := range r.TaskContentPublications {
			if r.TaskContentPublications[i].OutcomeRef.ManifestSHA256 == head.ManifestSHA256 {
				r.TaskContentPublications[i].PublicationKey = changedKey
			}
		}
		r.TaskSpecPolicies[0].ActivationPublicationKey = &changedKey
		sortWorkRegistry(&r)
		data, err := encodeWorkItemRegistry(r)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(workItemsPath(root), data, 0600); err != nil {
			t.Fatal(err)
		}
		result, err := ListTasksWithFilters(f.dependencies, TaskListFilters{})
		if err != nil {
			t.Fatal(err)
		}
		title := result.Titles["task"]
		if title.Title != nil || title.Status != "unavailable" || title.Source != "problem_revision" {
			t.Fatalf("v4=%v unverifiable title exposed as available: %+v", v4, title)
		}
	}
}
