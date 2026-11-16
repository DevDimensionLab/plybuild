package workspace

import (
	"errors"
	"fmt"
	"io/fs"
)

// TaskWorktreeRetireInput selects one Task; without Apply nothing changes.
type TaskWorktreeRetireInput struct {
	TaskID TaskID
	Apply  bool
}

// TaskWorktreeRetireResult reports the observed Git state and each effect:
// Worktree is would_remove, removed or absent; Branch is would_delete,
// deleted or absent. State is the Task worktree_state after the call.
type TaskWorktreeRetireResult struct {
	Workspace string        `json:"workspace"`
	TaskID    TaskID        `json:"task_id"`
	Locator   string        `json:"locator"`
	Ref       string        `json:"ref"`
	OID       string        `json:"oid,omitempty"`
	EpicRef   string        `json:"epic_ref"`
	EpicOID   string        `json:"epic_oid"`
	Applied   bool          `json:"applied"`
	Worktree  string        `json:"worktree"`
	Branch    string        `json:"branch"`
	State     WorkItemState `json:"state"`
}

// RetireTaskWorktree removes a delivered Task's clean worktree and its branch
// once the branch commit is contained in the Epic ref, then marks the binding
// retired. Effects already present are accepted, so a rerun after an
// interruption finishes the remaining steps.
func RetireTaskWorktree(d Dependencies, input TaskWorktreeRetireInput) (TaskWorktreeRetireResult, error) {
	if _, err := ParseTaskID(string(input.TaskID)); err != nil {
		return TaskWorktreeRetireResult{}, err
	}
	if err := requireWorkMutationDependencies(d); err != nil {
		return TaskWorktreeRetireResult{}, err
	}
	root, err := containingRegistrationWorkspace(d)
	if err != nil {
		return TaskWorktreeRetireResult{}, err
	}
	result := TaskWorktreeRetireResult{Workspace: root, TaskID: input.TaskID}
	err = withTaskWork(d, root, input.TaskID, func(projects ProjectSnapshot, work ProjectWork, session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		return retireTaskWorktreeLocked(d, projects, session, registry, input, &result)
	})
	if err != nil {
		return TaskWorktreeRetireResult{}, mapMutationError(err)
	}
	return result, nil
}

func retireTaskWorktreeLocked(d Dependencies, projects ProjectSnapshot, session WorkItemStoreSession, registry WorkItemRegistry, input TaskWorktreeRetireInput, result *TaskWorktreeRetireResult) error {
	task, index := findTask(registry, input.TaskID)
	if task == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Task %s is not registered", input.TaskID), nil)
	}
	if task.Worktree == nil || (task.WorktreeState != WorkItemReady && task.WorktreeState != WorkItemRetired) {
		return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Task %s has no ready worktree binding (state %s)", task.ID, task.WorktreeState), nil)
	}
	w := task.Worktree
	result.Locator, result.Ref, result.State = w.Locator, w.Ref, task.WorktreeState
	if task.WorktreeState == WorkItemRetired {
		result.Worktree, result.Branch = "absent", "absent"
		return nil
	}
	_, repo, err := projectAndRepo(projects, task.ProjectID, task.RepoID)
	if err != nil {
		return err
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", task.ParentEpicID), nil)
	}
	binding, _ := findEpicRepo(*epic, task.RepoID)
	if binding == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s has no binding for repository %s", epic.ID, task.RepoID), nil)
	}
	epicRef, err := d.WorkGit.ObserveRef(repo, binding.Worktree.Ref)
	if err != nil {
		return err
	}
	if !epicRef.Exists {
		return workError(ErrorWorkParentStale, fmt.Sprintf("Epic ref %s does not exist", binding.Worktree.Ref), nil)
	}
	result.EpicRef, result.EpicOID = binding.Worktree.Ref, epicRef.OID
	ref, err := d.WorkGit.ObserveRef(repo, w.Ref)
	if err != nil {
		return err
	}
	worktreePresent := true
	if _, err := d.Files.Lstat(w.Locator); errors.Is(err, fs.ErrNotExist) {
		worktreePresent = false
	} else if err != nil {
		return workError(ErrorWorkPathConflict, fmt.Sprintf("cannot inspect Task worktree %s", w.Locator), err)
	}
	if worktreePresent {
		observed, err := d.WorkGit.ObserveWorktree(w.Locator)
		if err != nil {
			return err
		}
		if observed.Ref != w.Ref || !ref.Exists || observed.OID != ref.OID {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Task worktree %s no longer has its branch %s checked out", w.Locator, w.Ref), nil)
		}
		if !observed.Clean {
			return workError(ErrorWorkParentDirty, fmt.Sprintf("Task worktree %s has uncommitted or untracked files; commit, integrate or remove them first", w.Locator), nil)
		}
	}
	if ref.Exists {
		result.OID = ref.OID
		contained, err := d.WorkGit.IsAncestor(repo, ref.OID, epicRef.OID)
		if err != nil {
			return err
		}
		if !contained {
			return workError(ErrorWorkRefConflict, fmt.Sprintf("Task branch %s at %s is not contained in Epic ref %s; integrate it before retiring", w.Ref, ref.OID, binding.Worktree.Ref), nil)
		}
	}
	result.Worktree, result.Branch = effectLabel(worktreePresent, "would_remove"), effectLabel(ref.Exists, "would_delete")
	if !input.Apply {
		return nil
	}
	if worktreePresent {
		if outcome := d.WorkGit.RemoveWorktree(repo, w.Locator); outcome.Err != nil {
			return workError(ErrorWorkGitEffect, fmt.Sprintf("git worktree remove %s: exit %d; rerun to finish", w.Locator, outcome.Exit), outcome.Err)
		}
		result.Worktree = "removed"
	}
	if ref.Exists {
		if outcome := d.WorkGit.DeleteRef(repo, w.Ref, ref.OID); outcome.Err != nil {
			return workError(ErrorWorkGitEffect, fmt.Sprintf("delete %s at %s: exit %d; rerun to finish", w.Ref, ref.OID, outcome.Exit), outcome.Err)
		}
		result.Branch = "deleted"
	}
	registry.Tasks[index].WorktreeState = WorkItemRetired
	if err := session.Publish(registry); err != nil {
		return err
	}
	result.Applied, result.State = true, WorkItemRetired
	return nil
}

func effectLabel(present bool, planned string) string {
	if present {
		return planned
	}
	return "absent"
}
