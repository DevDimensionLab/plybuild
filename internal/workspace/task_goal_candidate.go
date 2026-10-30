package workspace

import "path/filepath"

func (s *TaskSpecSession) CurrentDeliveryTaskSpec(id TaskID) (TaskSpecEvaluation, error) {
	task, _ := findTask(s.Registry, id)
	if task == nil || task.Worktree == nil {
		return TaskSpecEvaluation{}, contentError("task_spec_binding_conflict", "Task has no worktree", nil)
	}
	return currentTaskSpec(s.Dependencies, s.Root, s.Registry, s.Projects, *task, false)
}

func (s *TaskSpecSession) ValidateDeliveryCandidateTarget(locator, ref, common string, basis *TaskSpecBasis) (TaskSpecEvaluation, error) {
	return s.ValidateDeliveryCandidateTargetBasis(locator, ref, common, basis, nil)
}

func (s *TaskSpecSession) ValidateDeliveryCandidateTargetBasis(locator, ref, common string, basis *TaskSpecBasis, effective *EffectiveDeliveryBasis) (TaskSpecEvaluation, error) {
	out, err := s.ValidateTarget(locator, ref, common, basis)
	if err != nil {
		return out, err
	}
	physical, err := filepath.EvalSymlinks(locator)
	if err != nil {
		return out, err
	}
	var task *TaskRecord
	for _, candidate := range s.Registry.Tasks {
		if candidate.Worktree != nil && candidate.Worktree.Locator == physical && candidate.Worktree.Ref == ref && candidate.GitCommonDir == common {
			copy := candidate
			task = &copy
			break
		}
	}
	if task == nil || locator != physical {
		return out, contentError("task_spec_binding_conflict", "candidate has no exact registered Task source", nil)
	}
	candidate, err := s.Dependencies.WorkGit.ObserveWorktree(locator)
	if err != nil {
		return out, err
	}
	if !candidate.Clean || !candidate.InventoryMatch || candidate.Locator != locator || candidate.Ref != ref || candidate.GitCommonDir != common {
		return out, contentError("task_spec_basis_stale", "candidate must be the exact clean registered worktree", nil)
	}
	_, repo, err := projectAndRepo(s.Projects, task.ProjectID, task.RepoID)
	if err != nil {
		return out, err
	}
	start := task.Worktree.OID
	if operation, _ := findOperation(s.Registry, task.ID); operation != nil {
		start = operationStartOID(*operation)
	}
	ancestor, err := s.Dependencies.IntegrationGit.CheckAncestor(repo, start, candidate.OID)
	if err != nil {
		return out, err
	}
	if !ancestor {
		return out, contentError("task_spec_binding_conflict", "candidate is not descended from its registered source start", nil)
	}
	if effective != nil {
		if err := ValidateEffectiveDeliveryBasis(s.Registry, *effective); err != nil {
			return out, err
		}
		preparation := findPreparation(s.Registry, effective.PreparationID)
		if preparation == nil || preparation.Plan.TaskID != task.ID {
			return out, contentError("task_spec_binding_conflict", "source reconciliation belongs to another Task", nil)
		}
		current, err := EffectiveTaskDeliveryBasis(s.Registry, preparation.ID)
		if err != nil {
			return out, err
		}
		if !contentTypedEqual(current, *effective) {
			return out, contentError("reconciliation_required", "candidate reconciliation generation differs", nil)
		}
		if err := validateReconciledSourceAncestry(s.Dependencies, s.Registry, repo, *effective, candidate.OID); err != nil {
			return out, err
		}
		ancestor, err := s.Dependencies.IntegrationGit.CheckAncestor(repo, effective.Parent.OID, candidate.OID)
		if err != nil {
			return out, err
		}
		if !ancestor {
			return out, contentError("reconciliation_required", "source does not contain its recorded reconciled parent", nil)
		}
	}
	return out, nil
}
