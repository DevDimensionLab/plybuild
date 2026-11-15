package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// EpicRetireInput names one repository binding of an integrated Epic and the
// target ref, such as refs/heads/main, that must contain its branch.
type EpicRetireInput struct {
	EpicID    EpicID
	RepoID    RepoID
	TargetRef string
	Actor     string
	Apply     bool
}

// EpicRetireResult reports the observed refs and each effect: Worktree is
// would_remove, removed or absent; Branch is would_delete, deleted or absent;
// Lifecycle is the Epic lifecycle after the call.
type EpicRetireResult struct {
	Workspace string         `json:"workspace"`
	EpicID    EpicID         `json:"epic_id"`
	RepoID    RepoID         `json:"repo_id"`
	Locator   string         `json:"locator"`
	Ref       string         `json:"ref"`
	OID       string         `json:"oid,omitempty"`
	TargetRef string         `json:"target_ref"`
	TargetOID string         `json:"target_oid"`
	Applied   bool           `json:"applied"`
	Worktree  string         `json:"worktree"`
	Branch    string         `json:"branch"`
	Lifecycle LifecycleState `json:"lifecycle"`
}

// RetireEpic removes an integrated Epic's clean worktree and its branch once
// the branch commit is contained in the target ref and no Task of the Epic
// still has a worktree. It archives the Epic before any Git effect: the
// definitions keep naming the worktree path, and only an archived Epic may
// lack it, so a stop between the steps never leaves the workspace unusable.
// The registry keeps the binding as history. A rerun finishes what an
// interrupted run left.
func RetireEpic(d Dependencies, input EpicRetireInput) (EpicRetireResult, error) {
	if _, err := ParseEpicID(string(input.EpicID)); err != nil {
		return EpicRetireResult{}, err
	}
	if input.TargetRef == "" || input.Actor == "" {
		return EpicRetireResult{}, workError(ErrorWorkInvalidArguments, "target ref and actor are required", nil)
	}
	if err := requireWorkMutationDependencies(d); err != nil {
		return EpicRetireResult{}, err
	}
	observed, err := ObserveContainingRoot(d)
	if err != nil {
		return EpicRetireResult{}, err
	}
	if err := checkEpicPathPresent(d, observed.Root, input); err != nil {
		return EpicRetireResult{}, err
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return EpicRetireResult{}, err
	}
	run := func(apply bool) (EpicRetireResult, error) {
		result := EpicRetireResult{Workspace: root, EpicID: input.EpicID, RepoID: input.RepoID, TargetRef: input.TargetRef}
		step := input
		step.Apply = apply
		err := d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
			return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
				registry, err := session.Snapshot()
				if err != nil {
					return err
				}
				return retireEpicLocked(d, projects, registry, step, &result)
			})
		})
		if err != nil {
			return EpicRetireResult{}, mapMutationError(err)
		}
		lifecycle, err := ShowWorkItemLifecycle(d, "epic", string(input.EpicID))
		if err != nil {
			return result, err
		}
		result.Lifecycle = lifecycle.Lifecycle
		return result, nil
	}
	result, err := run(false)
	if err != nil || !input.Apply {
		return result, err
	}
	if result.Lifecycle != LifecycleArchived {
		if _, err = SetWorkItemLifecycle(d, WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: string(input.EpicID), State: LifecycleArchived, ActorClaim: input.Actor}); err != nil {
			return result, err
		}
	}
	return run(true)
}

func retireEpicLocked(d Dependencies, projects ProjectSnapshot, registry WorkItemRegistry, input EpicRetireInput, result *EpicRetireResult) error {
	epic, _ := findEpic(registry, input.EpicID)
	if epic == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", input.EpicID), nil)
	}
	binding, _ := findEpicRepo(*epic, input.RepoID)
	if binding == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s has no binding for repository %s", epic.ID, input.RepoID), nil)
	}
	w := binding.Worktree
	result.Locator, result.Ref = w.Locator, w.Ref
	if w.Ref == input.TargetRef {
		return workError(ErrorWorkInvalidArguments, fmt.Sprintf("Epic %s is on the target %s itself", epic.ID, input.TargetRef), nil)
	}
	for _, task := range registry.Tasks {
		if task.ParentEpicID == epic.ID && task.WorktreeState == WorkItemReady {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Task %s of Epic %s still has a worktree; integrate and retire it first", task.ID, epic.ID), nil)
		}
	}
	_, repo, err := projectAndRepo(projects, epic.ProjectID, input.RepoID)
	if err != nil {
		return err
	}
	if repo.Locator == w.Locator {
		return workError(ErrorWorkInvalidArguments, fmt.Sprintf("Epic %s uses the repository's main worktree %s, which is never removed", epic.ID, w.Locator), nil)
	}
	target, err := d.WorkGit.ObserveRef(repo, input.TargetRef)
	if err != nil {
		return err
	}
	if !target.Exists {
		return workError(ErrorWorkParentStale, fmt.Sprintf("target ref %s does not exist", input.TargetRef), nil)
	}
	result.TargetOID = target.OID
	ref, err := d.WorkGit.ObserveRef(repo, w.Ref)
	if err != nil {
		return err
	}
	worktreePresent := true
	if _, err := d.Files.Lstat(w.Locator); errors.Is(err, fs.ErrNotExist) {
		worktreePresent = false
	} else if err != nil {
		return workError(ErrorWorkPathConflict, fmt.Sprintf("cannot inspect Epic worktree %s", w.Locator), err)
	}
	if worktreePresent {
		observed, err := d.WorkGit.ObserveWorktree(w.Locator)
		if err != nil {
			return err
		}
		if observed.Ref != w.Ref || !ref.Exists || observed.OID != ref.OID {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Epic worktree %s no longer has its branch %s checked out", w.Locator, w.Ref), nil)
		}
		if !observed.Clean {
			return workError(ErrorWorkParentDirty, fmt.Sprintf("Epic worktree %s has uncommitted or untracked files; commit, integrate or remove them first", w.Locator), nil)
		}
	}
	if ref.Exists {
		result.OID = ref.OID
		contained, err := d.WorkGit.IsAncestor(repo, ref.OID, target.OID)
		if err != nil {
			return err
		}
		if !contained {
			return workError(ErrorWorkRefConflict, fmt.Sprintf("Epic branch %s at %s is not contained in %s; integrate the Epic before retiring", w.Ref, ref.OID, input.TargetRef), nil)
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
	result.Applied = true
	return nil
}

// checkEpicPathPresent names the one recovery when an earlier Ply removed the
// worktree of a defined Epic before archiving it: every Project lookup then
// fails until the worktree is back on its path.
func checkEpicPathPresent(d Dependencies, root string, input EpicRetireInput) error {
	defs, _, err := readDefinitions(root)
	if err != nil || defs == nil {
		return err
	}
	lifecycles, err := readWorkItemLifecycle(root)
	if err != nil || lifecycles.Epic(input.EpicID) == LifecycleArchived {
		return err
	}
	for _, p := range defs.Projects {
		for _, e := range p.Epics {
			if e.ID != string(input.EpicID) || e.Repo != string(input.RepoID) {
				continue
			}
			path := e.Path
			if !filepath.IsAbs(path) {
				path = filepath.Join(root, path)
			}
			if _, err := d.Files.Lstat(path); errors.Is(err, fs.ErrNotExist) {
				branch := strings.TrimPrefix(e.Ref, "refs/heads/")
				return workError(ErrorWorkPathConflict, fmt.Sprintf("Epic %s is defined at %s, which is gone, and the Epic is not archived, so Project lookups fail; put the worktree back with git -C <repository> worktree add %s -b %s %s, then rerun retire", input.EpicID, path, path, branch, input.TargetRef), nil)
			}
		}
	}
	return nil
}
