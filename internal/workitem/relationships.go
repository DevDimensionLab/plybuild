package workitem

import "fmt"

// FindEpic returns the first matching slice element and its index, or nil, -1.
// The pointer must not be retained across append/reallocation of the slice.
func FindEpic(epics []EpicRecord, id EpicID) (*EpicRecord, int) {
	for i := range epics {
		if epics[i].ID == id {
			return &epics[i], i
		}
	}
	return nil, -1
}

// FindTask returns the first matching slice element and its index, or nil, -1.
// The pointer must not be retained across append/reallocation of the slice.
func FindTask(tasks []TaskRecord, id TaskID) (*TaskRecord, int) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], i
		}
	}
	return nil, -1
}

// FindEpicRepo returns the first matching binding and its index, or nil, -1.
// The pointer aliases epic.RepoBindings and has the same lifetime as that slice.
func FindEpicRepo(epic EpicRecord, repoID RepoID) (*EpicRepoBinding, int) {
	for i := range epic.RepoBindings {
		if epic.RepoBindings[i].RepoID == repoID {
			return &epic.RepoBindings[i], i
		}
	}
	return nil, -1
}

// ValidateTaskParent checks only the registered Epic/Project and repo/common-dir
// relationship. It does not validate syntax, worktree state or historical bases.
func ValidateTaskParent(task TaskRecord, epic *EpicRecord) error {
	if epic == nil || epic.ID != task.ParentEpicID || epic.ProjectID != task.ProjectID {
		return fmt.Errorf("Task %s has invalid Epic/Project reference", task.ID)
	}
	binding, _ := FindEpicRepo(*epic, task.RepoID)
	if binding == nil || binding.GitCommonDir != task.GitCommonDir {
		return fmt.Errorf("Task %s has invalid repository binding", task.ID)
	}
	return nil
}
