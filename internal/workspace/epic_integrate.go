package workspace

import (
	"fmt"
)

// EpicIntegrateInput names the Epic, its repository binding and the exact
// target and Epic commits the caller saw; without Apply nothing changes.
type EpicIntegrateInput struct {
	EpicID            EpicID
	RepoID            RepoID
	TargetRef         string
	ExpectedTargetOID string
	ExpectedEpicOID   string
	Apply             bool
}

// EpicIntegrateResult reports the observed refs and the effect:
// would_fast_forward, fast_forwarded or already_integrated.
type EpicIntegrateResult struct {
	Workspace      string `json:"workspace"`
	EpicID         EpicID `json:"epic_id"`
	RepoID         RepoID `json:"repo_id"`
	EpicRef        string `json:"epic_ref"`
	EpicOID        string `json:"epic_oid"`
	TargetRef      string `json:"target_ref"`
	TargetOID      string `json:"target_oid"`
	TargetWorktree string `json:"target_worktree"`
	Applied        bool   `json:"applied"`
	Effect         string `json:"effect"`
}

// IntegrateEpic fast-forwards the target ref, usually refs/heads/main, to the
// exact Epic commit. A checked-out target must be a clean worktree and moves
// with git merge --ff-only; otherwise the ref moves only from the expected
// commit. The reflog entry names the Epic. It records nothing in the registry.
func IntegrateEpic(d Dependencies, input EpicIntegrateInput) (EpicIntegrateResult, error) {
	if !validFullBranchRef(input.TargetRef) {
		return EpicIntegrateResult{}, WorkInvalidArguments("--target-ref must be a full local branch ref such as refs/heads/main")
	}
	if !validOIDText(input.ExpectedTargetOID) || !validOIDText(input.ExpectedEpicOID) {
		return EpicIntegrateResult{}, WorkInvalidArguments("--expected-target-oid and --expected-epic-oid must be lowercase full Git object IDs")
	}
	if err := requireWorkMutationDependencies(d); err != nil {
		return EpicIntegrateResult{}, err
	}
	root, err := containingRegistrationWorkspace(d)
	if err != nil {
		return EpicIntegrateResult{}, err
	}
	result := EpicIntegrateResult{Workspace: root, EpicID: input.EpicID, RepoID: input.RepoID, TargetRef: input.TargetRef}
	err = withEpicWork(d, root, input.EpicID, func(projects ProjectSnapshot, work ProjectWork, session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		return integrateEpicLocked(d, projects, registry, input, &result)
	})
	if err != nil {
		return EpicIntegrateResult{}, mapMutationError(err)
	}
	return result, nil
}

func integrateEpicLocked(d Dependencies, projects ProjectSnapshot, registry WorkItemRegistry, input EpicIntegrateInput, result *EpicIntegrateResult) error {
	epic, _ := findEpic(registry, input.EpicID)
	if epic == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", input.EpicID), nil)
	}
	binding, _ := findEpicRepo(*epic, input.RepoID)
	if binding == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s has no binding for repository %s", epic.ID, input.RepoID), nil)
	}
	if binding.Worktree.Ref == input.TargetRef {
		return workError(ErrorWorkRefConflict, "target ref is the Epic ref itself", nil)
	}
	_, repo, err := projectAndRepo(projects, epic.ProjectID, input.RepoID)
	if err != nil {
		return err
	}
	epicRef, err := d.WorkGit.ObserveRef(repo, binding.Worktree.Ref)
	if err != nil {
		return err
	}
	result.EpicRef, result.EpicOID = binding.Worktree.Ref, epicRef.OID
	if !epicRef.Exists || epicRef.OID != input.ExpectedEpicOID {
		return workError(ErrorWorkParentStale, fmt.Sprintf("Epic ref %s is at %s, not the expected %s", binding.Worktree.Ref, epicRef.OID, input.ExpectedEpicOID), nil)
	}
	target, err := d.WorkGit.ObserveRef(repo, input.TargetRef)
	if err != nil {
		return err
	}
	result.TargetOID, result.TargetWorktree = target.OID, target.CheckedOutAt
	if !target.Exists || target.OID != input.ExpectedTargetOID {
		return workError(ErrorWorkParentStale, fmt.Sprintf("target ref %s is at %s, not the expected %s", input.TargetRef, target.OID, input.ExpectedTargetOID), nil)
	}
	if target.OID == epicRef.OID {
		result.Effect = "already_integrated"
		return nil
	}
	forward, err := d.WorkGit.IsAncestor(repo, target.OID, epicRef.OID)
	if err != nil {
		return err
	}
	if !forward {
		return workError(ErrorWorkRefConflict, fmt.Sprintf("%s has commits that are not in Epic %s; integrate them into the Epic first, Ply only fast-forwards", input.TargetRef, epic.ID), nil)
	}
	if target.CheckedOutAt != "" {
		observed, err := d.WorkGit.ObserveWorktree(target.CheckedOutAt)
		if err != nil {
			return err
		}
		if observed.Ref != input.TargetRef || observed.OID != target.OID {
			return workError(ErrorWorkParentStale, fmt.Sprintf("worktree %s no longer has %s at the expected commit", target.CheckedOutAt, input.TargetRef), nil)
		}
		if !observed.Clean {
			return workError(ErrorWorkParentDirty, fmt.Sprintf("target worktree %s has uncommitted or untracked files", target.CheckedOutAt), nil)
		}
	}
	result.Effect = "would_fast_forward"
	if !input.Apply {
		return nil
	}
	outcome := d.WorkGit.FastForwardRef(GitFastForwardInput{Repository: repo, Ref: input.TargetRef, Worktree: target.CheckedOutAt, From: target.OID, To: epicRef.OID, Reason: "ply-workspace-epic-integrate:" + string(epic.ID)})
	after, err := d.WorkGit.ObserveRef(repo, input.TargetRef)
	if err != nil {
		return workError(ErrorWorkGitEffect, fmt.Sprintf("fast-forward of %s ran (exit %d) but its result could not be read; inspect %s before retrying", input.TargetRef, outcome.Exit, input.TargetRef), err)
	}
	result.TargetOID = after.OID
	if after.OID != epicRef.OID {
		return workError(ErrorWorkGitEffect, fmt.Sprintf("fast-forward of %s failed (exit %d); it is still at %s", input.TargetRef, outcome.Exit, after.OID), outcome.Err)
	}
	result.Applied, result.Effect = true, "fast_forwarded"
	return nil
}
