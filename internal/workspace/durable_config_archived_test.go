package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const archivedEpicConfig = `{"format_version":1,"projects":[{"id":"sample","name":"Sample","wrapper":".","repos":[{"id":"code","path":"repo"}],"epics":[{"id":"feature","title":"Feature","repo":"code","path":"repo","ref":"refs/heads/main"},{"id":"old","title":"Old","repo":"code","path":"box/old","ref":"refs/heads/old"}]}]}`

func archivedEpicFixture(t *testing.T, state LifecycleState) (string, string, Dependencies) {
	t.Helper()
	root, repo, d := durableFixture(t)
	d.Files = projectCwdFileSystem{FileSystem: d.Files.(projectCwdFileSystem).FileSystem, cwd: root}
	if err := os.Mkdir(filepath.Join(root, "box"), 0755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "worktree", "add", "-b", "old", filepath.Join(root, "box", "old"))
	if err := os.WriteFile(filepath.Join(root, "ply-workspace.json"), []byte(archivedEpicConfig), 0644); err != nil {
		t.Fatal(err)
	}
	if epics, err := ListEpics(d); err != nil || len(epics.Epics) != 2 {
		t.Fatalf("bootstrap: %#v, %v", epics, err)
	}
	if _, err := CreateTask(d, TaskCreateInput{TaskID: "child", Title: "Child", Description: "Fixture", ParentEpicID: "old", ProjectID: "sample", RepoID: "code"}); err != nil {
		t.Fatal(err)
	}
	if state == LifecycleCompleted {
		event := WorkItemLifecycleEvent{Sequence: 1, SubjectKind: "epic", SubjectID: "old", From: LifecycleActive, To: state, ActorClaim: "fixture", RecordedAtUTC: "2026-11-10T00:00:00Z"}
		event.ID = lifecycleEventID(event)
		bytes, err := encodeWorkItemLifecycle([]WorkItemLifecycleEvent{event})
		if err == nil {
			err = workspaceMetadataPublisher{}.publish(root, WorkItemLifecycleFile, bytes)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, subject := range []struct{ kind, id string }{{"task", "child"}, {"epic", "old"}} {
		if state == LifecycleCompleted {
			break
		}
		if _, err := SetWorkItemLifecycle(d, WorkItemLifecycleInput{SubjectKind: subject.kind, SubjectID: subject.id, State: state, ActorClaim: "fixture"}); err != nil {
			t.Fatal(err)
		}
	}
	runLocalGit(t, repo, "worktree", "remove", filepath.Join(root, "box", "old"))
	return root, repo, d
}

type editedRegistrations struct {
	WorkItemStore
	edit func(*WorkItemRegistry)
}

func (s editedRegistrations) Snapshot(root string) (WorkItemRegistry, error) {
	r, err := s.WorkItemStore.Snapshot(root)
	s.edit(&r)
	return r, err
}

func (s editedRegistrations) SnapshotRegistrations(root string) (WorkItemRegistry, error) {
	r, err := s.WorkItemStore.SnapshotRegistrations(root)
	s.edit(&r)
	return r, err
}

func editOldEpic(edit func(*EpicRecord)) func(*WorkItemRegistry) {
	return func(r *WorkItemRegistry) {
		epics := append([]EpicRecord(nil), r.Epics...)
		for i := range epics {
			if epics[i].ID == "old" {
				epics[i].RepoBindings = append([]EpicRepoBinding(nil), epics[i].RepoBindings...)
				edit(&epics[i])
			}
		}
		r.Epics = epics
	}
}

func TestDurableConfigArchivedEpicWithoutWorktreeStaysHistorical(t *testing.T) {
	root, _, d := archivedEpicFixture(t, LifecycleArchived)
	before := snapshotTree(t, root)
	epics, err := ListEpics(d)
	if err != nil || len(epics.Epics) != 2 {
		t.Fatalf("list after Git cleanup: %#v, %v", epics, err)
	}
	old, err := ShowEpic(d, "old")
	if err != nil {
		t.Fatal(err)
	}
	if old.Lifecycle != LifecycleArchived || old.Ready || !reflect.DeepEqual(old.Reasons, []string{"epic_observation_unknown"}) {
		t.Fatalf("archived Epic readback: %#v", old)
	}
	if old.Epic.RepoBindings[0].Worktree.Locator != filepath.Join(root, "box", "old") {
		t.Fatalf("historical binding changed: %#v", old.Epic.RepoBindings)
	}
	feature, err := ShowEpic(d, "feature")
	if err != nil || !feature.Ready {
		t.Fatalf("other Epic: %#v, %v", feature, err)
	}
	if tasks, err := ListTasks(d, nil); err != nil || len(tasks.Tasks) != 1 {
		t.Fatalf("historical Task: %#v, %v", tasks, err)
	}
	if after := snapshotTree(t, root); !reflect.DeepEqual(after, before) {
		t.Fatal("reading an archived Epic wrote workspace files")
	}
}

func TestDurableConfigArchivedExceptionIsNarrow(t *testing.T) {
	for _, state := range []LifecycleState{LifecycleActive, LifecycleParked, LifecycleFrozen, LifecycleCompleted} {
		t.Run(string(state), func(t *testing.T) {
			_, _, d := archivedEpicFixture(t, state)
			if _, err := ListEpics(d); err == nil {
				t.Fatal("missing non-archived Epic worktree accepted")
			}
		})
	}
	conflict, pathError := string(ErrorWorkIdentityConflict), string(ErrorProjectPath)
	for test, want := range map[string]string{"project": conflict, "repo": conflict, "common-dir": conflict, "title": conflict, "ref": conflict, "locator": conflict, "corrupt-lifecycle": string(ErrorWorkStoreConflict), "unregistered": pathError, "dangling-symlink": pathError, "file": pathError, "parent-symlink": conflict, "permission": pathError} {
		t.Run(test, func(t *testing.T) {
			root, _, d := archivedEpicFixture(t, LifecycleArchived)
			config := archivedEpicConfig
			path := filepath.Join(root, "box", "old")
			var err error
			switch test {
			case "project":
				d.WorkItems = editedRegistrations{d.WorkItems, editOldEpic(func(e *EpicRecord) { e.ProjectID = "elsewhere" })}
			case "repo":
				d.WorkItems = editedRegistrations{d.WorkItems, editOldEpic(func(e *EpicRecord) { e.RepoBindings[0].RepoID = "elsewhere" })}
			case "common-dir":
				d.WorkItems = editedRegistrations{d.WorkItems, editOldEpic(func(e *EpicRecord) { e.RepoBindings[0].GitCommonDir = filepath.Join(root, "elsewhere.git") })}
			case "title":
				config = strings.Replace(config, `"title":"Old"`, `"title":"Renamed"`, 1)
			case "ref":
				config = strings.Replace(config, "refs/heads/old", "refs/heads/main", 1)
			case "locator":
				config = strings.Replace(config, "box/old", "box/moved", 1)
			case "corrupt-lifecycle":
				err = os.WriteFile(lifecycleLocator(root), []byte("{"), 0644)
			case "unregistered":
				err = os.Remove(workItemsPath(root))
			case "dangling-symlink":
				err = os.Symlink(filepath.Join(root, "missing"), path)
			case "file":
				err = os.WriteFile(path, []byte("file"), 0644)
			case "parent-symlink":
				if err = os.Mkdir(filepath.Join(root, "elsewhere"), 0755); err == nil {
					if err = os.Remove(filepath.Join(root, "box")); err == nil {
						err = os.Symlink(filepath.Join(root, "elsewhere"), filepath.Join(root, "box"))
					}
				}
			case "permission":
				if err = os.Chmod(filepath.Join(root, "box"), 0); err == nil {
					t.Cleanup(func() { _ = os.Chmod(filepath.Join(root, "box"), 0755) })
				}
			}
			if err == nil {
				err = os.WriteFile(filepath.Join(root, "ply-workspace.json"), []byte(config), 0644)
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = ListEpics(d); err == nil || !strings.HasPrefix(err.Error(), want+":") {
				t.Fatalf("archived exception accepted an unproven absence: %v", err)
			}
		})
	}
}

func TestDurableConfigArchivedEpicDoesNotHidePrimaryRepo(t *testing.T) {
	root, repo, d := durableFixture(t)
	d.Files = projectCwdFileSystem{FileSystem: d.Files.(projectCwdFileSystem).FileSystem, cwd: root}
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	if _, err := SetWorkItemLifecycle(d, WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: "feature", State: LifecycleArchived, ActorClaim: "fixture"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(repo, filepath.Join(root, "moved")); err != nil {
		t.Fatal(err)
	}
	if _, err := ListEpics(d); err == nil {
		t.Fatal("archived Epic hid a missing primary repository")
	}
}

func TestDurableConfigArchivedRuleDoesNotHidePendingCreate(t *testing.T) {
	for _, state := range []string{"pending", "reconciliation_required"} {
		t.Run(state, func(t *testing.T) {
			_, _, d := archivedEpicFixture(t, LifecycleArchived)
			d.WorkItems = editedRegistrations{d.WorkItems, func(r *WorkItemRegistry) {
				var epics []EpicRecord
				for _, e := range r.Epics {
					if e.ID != "old" {
						epics = append(epics, e)
					}
				}
				r.Epics = epics
				r.EpicCreateOperations = append(append([]EpicCreateOperation(nil), r.EpicCreateOperations...), EpicCreateOperation{ID: "pending", Request: EpicCreateInput{EpicID: "old"}, State: state})
			}}
			epics, err := ListEpics(d)
			if err != nil || len(epics.Epics) != 1 || epics.Epics[0].ID != "feature" {
				t.Fatalf("pending create adopted or hidden: %#v, %v", epics, err)
			}
			shown, err := ShowEpic(d, "old")
			if err != nil || shown.Creation == nil || shown.Creation.State != state || shown.Ready || !reflect.DeepEqual(shown.Reasons, []string{"epic_creation_pending"}) || shown.Lifecycle == LifecycleArchived {
				t.Fatalf("pending create shown as archived Epic: %#v, %v", shown, err)
			}
		})
	}
}
