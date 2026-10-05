package workspace

import "path/filepath"

// CurrentDeliveryTaskSpec is for building the initial delivery handoff. Ordinary
// legacy validation is unchanged; only an execution contract qualifies here.
func (s *TaskSpecSession) CurrentDeliveryTaskSpec(id TaskID) (TaskSpecEvaluation, error) {
	task, _ := findTask(s.Registry, id)
	if task == nil || task.Worktree == nil {
		return TaskSpecEvaluation{}, contentError("task_spec_binding_conflict", "delivery Task has no worktree", nil)
	}
	out, e := currentTaskSpec(s.Dependencies, s.Root, s.Registry, s.Projects, *task, true)
	if e != nil {
		return out, e
	}
	if contentString(contentFields(out.Spec), "contract_kind") != "execution" || out.Basis == nil || out.SelectionFreshness != "current" || out.ContentIntegrity != "valid" || out.TargetFreshness != "fresh" {
		return out, contentError("task_spec_selection_stale", "a current, intact execution Spec and clean start are required", nil)
	}
	return out, nil
}

// ValidateDeliveryCandidateTarget validates a later candidate in the same owned
// execution. It never changes the frozen start basis or grants an initial start
// for changed code. The caller must separately bind an already accepted owner.
func (s *TaskSpecSession) ValidateDeliveryCandidateTarget(locator, ref, common string, basis *TaskSpecBasis) (TaskSpecEvaluation, error) {
	var out TaskSpecEvaluation
	if basis == nil {
		return out, contentError("task_spec_required", "exact delivery Spec basis is required", nil)
	}
	task, _ := findTask(s.Registry, basis.TaskID)
	physical, e := filepath.EvalSymlinks(locator)
	if e != nil {
		return out, e
	}
	if task == nil || task.Worktree == nil || locator != physical || task.Worktree.Locator != physical || task.Worktree.Ref != ref || task.GitCommonDir != common {
		return out, contentError("task_spec_binding_conflict", "candidate differs from its registered delivery worktree", nil)
	}
	out, e = currentTaskSpec(s.Dependencies, s.Root, s.Registry, s.Projects, *task, false)
	if e != nil {
		return out, e
	}
	if contentString(contentFields(out.Spec), "contract_kind") != "execution" || !contentTypedEqual(out.Basis, basis) || out.SelectionFreshness != "current" || out.ContentIntegrity != "valid" || out.TargetFreshness != "fresh" {
		return out, contentError("task_spec_selection_stale", "candidate requires the current unchanged delivery contract", nil)
	}
	bm := contentFields(contentFields(out.Spec)["implementation_basis"])
	if e = validateSpecImplementationBasis(s.Dependencies, s.Projects, s.Registry, *task, contentFields(out.Spec)["implementation_basis"], false); e != nil {
		return out, e
	}
	epic, binding := queueTargetBinding(s.Registry, QueueTarget{EpicID: task.ParentEpicID, RepoID: task.RepoID})
	base := currentEpicBase(s.Registry, epic, binding)
	if base.OID != contentString(bm, "parent_oid") || base.Tree != contentString(bm, "parent_tree") || pendingBaseUpdate(s.Registry, task.ParentEpicID, task.RepoID) {
		return out, contentError("epic_base_changed", "delivery parent changed before candidate control", nil)
	}
	target := QueueTarget{ProjectID: task.ProjectID, RepoID: task.RepoID, EpicID: task.ParentEpicID, GitCommonDir: common, ParentWorktreeID: binding.Worktree.ID, ParentLocator: binding.Worktree.Locator, ParentRef: binding.Worktree.Ref}
	parent, e := queueParentObservation(s.Dependencies, target)
	if e != nil {
		return out, e
	}
	if parent.OID != base.OID || parent.Tree != base.Tree {
		return out, contentError("task_spec_basis_stale", "live parent differs from the frozen execution parent", nil)
	}
	candidate, e := s.Dependencies.WorkGit.ObserveWorktree(locator)
	if e != nil {
		return out, e
	}
	if !candidate.Clean || !candidate.InventoryMatch || candidate.Locator != locator || candidate.Ref != ref || candidate.GitCommonDir != common {
		return out, contentError("task_spec_basis_stale", "candidate must be the exact clean owned worktree", nil)
	}
	_, repo, e := projectAndRepo(s.Projects, task.ProjectID, task.RepoID)
	if e != nil {
		return out, e
	}
	ancestor, e := s.Dependencies.IntegrationGit.CheckAncestor(repo, contentString(bm, "start_oid"), candidate.OID)
	if e != nil {
		return out, e
	}
	if !ancestor {
		return out, contentError("task_spec_binding_conflict", "candidate is not descended from the frozen execution start", nil)
	}
	return out, nil
}
