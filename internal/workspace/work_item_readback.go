package workspace

import (
	"fmt"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type EpicObservedRepo struct {
	RepoID                                RepoID
	Locator, Ref, OID, Tree, GitCommonDir *string
	Clean, InventoryMatch                 *bool
}
type EpicFreshnessRepo struct {
	RepoID                                                       RepoID
	Locator, Ref, OID, Tree, GitCommonDir, Clean, InventoryMatch string
}
type EpicReadbackResult struct {
	Workspace string
	Epic      EpicRecord
	Observed  []EpicObservedRepo
	Freshness []EpicFreshnessRepo
	Ready     bool
	Reasons   []string
}

type TaskObservedWorktree struct {
	Locator, Ref, OID, Tree, GitCommonDir *string
	Clean, InventoryMatch                 *bool
}
type TaskObservedSource struct{ Ref, OID, Tree, CheckedOutAt *string }
type TaskObservedTarget struct {
	Kind                                  string
	Locator, Ref, OID, Tree, GitCommonDir *string
	Clean, InventoryMatch                 *bool
}
type TaskFreshnessWorktree struct{ Locator, Ref, OID, Tree, GitCommonDir, Clean, InventoryMatch string }
type TaskFreshnessSource struct{ Ref, OID, Tree, CheckedOutAt string }
type TaskFreshnessTarget struct{ Kind, Locator, Ref, OID, Tree, GitCommonDir, Clean, InventoryMatch string }
type TaskReadbackResult struct {
	Workspace                      string
	Task                           TaskRecord
	Epic                           EpicRecord
	Operation                      *WorktreeOperationRecord
	ProjectGitCommonDir            *string
	Parent                         TaskObservedWorktree
	Source                         TaskObservedSource
	Target                         TaskObservedTarget
	ProjectFreshness               string
	ParentFreshness                TaskFreshnessWorktree
	SourceFreshness                TaskFreshnessSource
	TargetFreshness                TaskFreshnessTarget
	WorktreeReady, ReadyForHandoff bool
	Reasons                        []string
	Integration                    *WorkspaceTaskIntegrationReadback
}

func ShowEpic(dependencies Dependencies, id EpicID) (EpicReadbackResult, error) {
	if _, err := ParseEpicID(string(id)); err != nil {
		return EpicReadbackResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return EpicReadbackResult{}, err
	}
	if dependencies.WorkItems == nil {
		return EpicReadbackResult{}, workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	registry, err := dependencies.WorkItems.Snapshot(root)
	if err != nil {
		return EpicReadbackResult{}, err
	}
	epic, _ := findEpic(registry, id)
	if epic == nil {
		return EpicReadbackResult{}, workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", id), nil)
	}
	result := EpicReadbackResult{Workspace: root, Epic: *epic, Ready: true, Reasons: []string{}}
	var projects []ProjectRecord
	var repos []RepoRecord
	if dependencies.Projects != nil {
		projects, repos, err = dependencies.Projects.Snapshot(root)
	} else {
		err = fmt.Errorf("project store unavailable")
	}
	for _, binding := range epic.RepoBindings {
		observed := EpicObservedRepo{RepoID: binding.RepoID}
		fresh := EpicFreshnessRepo{RepoID: binding.RepoID, Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}
		var repository RepoRecord
		bindingFresh := false
		if err == nil {
			var bindingErr error
			_, repository, bindingErr = projectAndRepo(ProjectSnapshot{Projects: projects, Repos: repos}, epic.ProjectID, binding.RepoID)
			bindingFresh = bindingErr == nil && repository.GitCommonDir == binding.GitCommonDir
		}
		if err != nil || !bindingFresh {
			result.Ready = false
			if err != nil {
				result.Reasons = append(result.Reasons, "project_binding_unknown")
			} else {
				result.Reasons = append(result.Reasons, "project_binding_stale")
			}
		} else if dependencies.WorkGit == nil {
			result.Ready = false
			result.Reasons = append(result.Reasons, "epic_observation_unknown")
		} else if current, observeErr := dependencies.WorkGit.ObserveWorktree(binding.Worktree.Locator); observeErr != nil {
			result.Ready = false
			result.Reasons = append(result.Reasons, "epic_observation_unknown")
		} else {
			observed.Locator = stringPointer(current.Locator)
			observed.Ref = stringPointer(current.Ref)
			observed.OID = stringPointer(current.OID)
			observed.Tree = stringPointer(current.Tree)
			observed.GitCommonDir = stringPointer(current.GitCommonDir)
			observed.Clean = boolPointer(current.Clean)
			observed.InventoryMatch = boolPointer(current.InventoryMatch)
			fresh.Locator = freshness(current.Locator == binding.Worktree.Locator)
			fresh.Ref = freshness(current.Ref == binding.Worktree.Ref)
			fresh.OID = freshness(current.OID == binding.Worktree.OID)
			fresh.Tree = freshness(current.Tree == binding.Worktree.Tree)
			fresh.GitCommonDir = freshness(current.GitCommonDir == binding.GitCommonDir)
			fresh.Clean = freshness(current.Clean)
			fresh.InventoryMatch = freshness(current.InventoryMatch)
			if current.Ref != binding.Worktree.Ref || current.OID != binding.Worktree.OID || current.Tree != binding.Worktree.Tree {
				result.Reasons = append(result.Reasons, "epic_ref_stale")
			}
			if current.Locator != binding.Worktree.Locator || current.GitCommonDir != binding.GitCommonDir {
				result.Reasons = append(result.Reasons, "epic_worktree_stale")
			}
			if !current.Clean {
				result.Reasons = append(result.Reasons, "epic_worktree_dirty")
			}
			if !current.InventoryMatch {
				result.Reasons = append(result.Reasons, "worktree_inventory_stale")
			}
		}
		result.Observed = append(result.Observed, observed)
		result.Freshness = append(result.Freshness, fresh)
	}
	result.Reasons = sortedReasons(result.Reasons)
	result.Ready = result.Ready && len(result.Reasons) == 0
	return result, nil
}

func ShowTask(dependencies Dependencies, id TaskID) (TaskReadbackResult, error) {
	if _, err := ParseTaskID(string(id)); err != nil {
		return TaskReadbackResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskReadbackResult{}, err
	}
	if dependencies.WorkItems == nil {
		return TaskReadbackResult{}, workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	registry, err := dependencies.WorkItems.Snapshot(root)
	if err != nil {
		return TaskReadbackResult{}, err
	}
	task, _ := findTask(registry, id)
	if task == nil {
		return TaskReadbackResult{}, workError(ErrorWorkNotFound, fmt.Sprintf("Task %s is not registered", id), nil)
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return TaskReadbackResult{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("Task %s refers to missing Epic", id), nil)
	}
	binding, _ := findEpicRepo(*epic, task.RepoID)
	if binding == nil {
		return TaskReadbackResult{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("Task %s has no parent repository binding", id), nil)
	}
	result := TaskReadbackResult{Workspace: root, Task: *task, Epic: *epic, ProjectFreshness: "unknown", ParentFreshness: unknownWorktreeFreshness(), SourceFreshness: TaskFreshnessSource{Ref: "unknown", OID: "unknown", Tree: "unknown", CheckedOutAt: "unknown"}, TargetFreshness: unknownTargetFreshness(), Reasons: []string{}, WorktreeReady: task.WorktreeState == WorkItemReady && task.Worktree != nil}
	if operation, _ := findOperation(registry, id); operation != nil {
		copyOperation := *operation
		result.Operation = &copyOperation
	}
	var projects []ProjectRecord
	var repositories []RepoRecord
	if dependencies.Projects != nil {
		projects, repositories, err = dependencies.Projects.Snapshot(root)
	} else {
		err = fmt.Errorf("project store unavailable")
	}
	var repository RepoRecord
	found := false
	if err == nil {
		var bindingErr error
		_, repository, bindingErr = projectAndRepo(ProjectSnapshot{Projects: projects, Repos: repositories}, task.ProjectID, task.RepoID)
		found = bindingErr == nil
	}
	if err != nil {
		result.Reasons = append(result.Reasons, "project_binding_unknown")
	} else if !found {
		result.Reasons = append(result.Reasons, "project_binding_stale")
	} else {
		result.ProjectGitCommonDir = stringPointer(repository.GitCommonDir)
		result.ProjectFreshness = freshness(repository.GitCommonDir == task.GitCommonDir)
		if result.ProjectFreshness == "stale" {
			result.Reasons = append(result.Reasons, "project_binding_stale")
		}
	}
	if dependencies.WorkGit == nil {
		result.Reasons = append(result.Reasons, "parent_observation_unknown")
	} else if parent, parentErr := dependencies.WorkGit.ObserveWorktree(binding.Worktree.Locator); parentErr != nil {
		result.Reasons = append(result.Reasons, "parent_observation_unknown")
	} else {
		result.Parent = observedWorktree(parent)
		result.ParentFreshness = freshWorktree(parent, binding.Worktree.Locator, binding.Worktree.Ref, binding.Worktree.OID, binding.Worktree.Tree, binding.GitCommonDir)
		if parent.Ref != binding.Worktree.Ref || parent.OID != binding.Worktree.OID || parent.Tree != binding.Worktree.Tree {
			result.Reasons = append(result.Reasons, "parent_ref_stale")
		}
		if parent.Locator != binding.Worktree.Locator || parent.GitCommonDir != binding.GitCommonDir {
			result.Reasons = append(result.Reasons, "parent_worktree_stale")
		}
		if !parent.Clean {
			result.Reasons = append(result.Reasons, "parent_worktree_dirty")
		}
		if !parent.InventoryMatch {
			result.Reasons = append(result.Reasons, "worktree_inventory_stale")
		}
	}
	if result.Operation == nil {
		result.Target = TaskObservedTarget{Kind: "absent"}
		result.TargetFreshness = TaskFreshnessTarget{Kind: "fresh", Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}
		result.Reasons = append(result.Reasons, "task_worktree_unbound")
	} else if found && dependencies.WorkGit != nil {
		operation := *result.Operation
		source, sourceErr := dependencies.WorkGit.ObserveRef(repository, operation.SourceRef)
		if sourceErr != nil {
			result.Reasons = append(result.Reasons, "source_ref_stale")
		} else if source.Exists {
			result.Source = TaskObservedSource{Ref: stringPointer(source.Ref), OID: stringPointer(source.OID), Tree: stringPointer(source.Tree)}
			if source.CheckedOutAt != "" {
				result.Source.CheckedOutAt = stringPointer(source.CheckedOutAt)
			}
			result.SourceFreshness = TaskFreshnessSource{Ref: "fresh", OID: freshness(source.OID == operation.ParentOID), Tree: freshness(source.Tree == operation.ParentTree), CheckedOutAt: freshness(source.CheckedOutAt == operation.TargetLocator)}
			if source.OID != operation.ParentOID || source.Tree != operation.ParentTree {
				result.Reasons = append(result.Reasons, "source_ref_stale")
			}
			if source.CheckedOutAt != "" && source.CheckedOutAt != operation.TargetLocator {
				result.Reasons = append(result.Reasons, "source_ref_checked_out_elsewhere")
			}
		} else {
			result.Reasons = append(result.Reasons, "source_ref_stale")
		}
		if task.WorktreeState == WorkItemReady {
			if target, targetErr := dependencies.WorkGit.ObserveWorktree(operation.TargetLocator); targetErr != nil {
				result.Target = TaskObservedTarget{Kind: "unknown"}
				result.Reasons = append(result.Reasons, "target_observation_unknown")
			} else {
				result.Target = TaskObservedTarget{Kind: "worktree", Locator: stringPointer(target.Locator), Ref: stringPointer(target.Ref), OID: stringPointer(target.OID), Tree: stringPointer(target.Tree), GitCommonDir: stringPointer(target.GitCommonDir), Clean: boolPointer(target.Clean), InventoryMatch: boolPointer(target.InventoryMatch)}
				result.TargetFreshness = freshTarget(target, operation)
				if target.Locator != operation.TargetLocator || target.Ref != operation.SourceRef || target.OID != operation.ParentOID || target.Tree != operation.ParentTree || target.GitCommonDir != operation.GitCommonDir {
					result.Reasons = append(result.Reasons, "target_worktree_stale")
				}
				if !target.Clean {
					result.Reasons = append(result.Reasons, "target_worktree_dirty")
				}
				if !target.InventoryMatch {
					result.Reasons = append(result.Reasons, "worktree_inventory_stale")
				}
			}
		} else {
			_, current := observeTaskCreate(dependencies, repository, *binding, operation)
			result.Target = TaskObservedTarget{Kind: current.TargetKind, Locator: current.TargetLocator, Ref: current.TargetRef, OID: current.TargetOID, Tree: current.TargetTree, GitCommonDir: current.TargetGitCommonDir, Clean: current.TargetClean, InventoryMatch: current.InventoryMatch}
			result.TargetFreshness = freshnessForObservedTarget(current, operation)
			result.Reasons = append(result.Reasons, current.Reasons...)
			if task.WorktreeState == WorkItemCreating {
				result.Reasons = append(result.Reasons, "task_worktree_creation_in_progress")
			} else {
				result.Reasons = append(result.Reasons, "task_worktree_reconciliation_required")
			}
		}
	} else {
		result.Reasons = append(result.Reasons, "target_observation_unknown")
	}
	result.Reasons = sortedReasons(result.Reasons)
	result.ReadyForHandoff = result.WorktreeReady && len(result.Reasons) == 0 && result.ProjectFreshness == "fresh"
	if registry.FormatVersion == 2 {
		integration := buildTaskShowIntegrationReadback(dependencies, root, ProjectSnapshot{Projects: projects, Repos: repositories}, registry, *task, *epic, *binding)
		result.Integration = &integration
	}
	return result, nil
}

func MarshalEpicReadback(result EpicReadbackResult) ([]byte, error) {
	observed := []canonicaljson.Value{}
	for _, item := range result.Observed {
		observed = append(observed, canonicaljson.Object{{Name: "repo_id", Value: string(item.RepoID)}, {Name: "locator", Value: pointerValue(item.Locator)}, {Name: "ref", Value: pointerValue(item.Ref)}, {Name: "oid", Value: pointerValue(item.OID)}, {Name: "tree", Value: pointerValue(item.Tree)}, {Name: "git_common_dir", Value: pointerValue(item.GitCommonDir)}, {Name: "clean", Value: pointerValue(item.Clean)}, {Name: "inventory_match", Value: pointerValue(item.InventoryMatch)}})
	}
	fresh := []canonicaljson.Value{}
	for _, item := range result.Freshness {
		fresh = append(fresh, canonicaljson.Object{{Name: "repo_id", Value: string(item.RepoID)}, {Name: "locator", Value: item.Locator}, {Name: "ref", Value: item.Ref}, {Name: "oid", Value: item.OID}, {Name: "tree", Value: item.Tree}, {Name: "git_common_dir", Value: item.GitCommonDir}, {Name: "clean", Value: item.Clean}, {Name: "inventory_match", Value: item.InventoryMatch}})
	}
	return canonicaljson.Marshal(canonicaljson.Object{{Name: "kind", Value: "WorkspaceEpicReadback@1"}, {Name: "workspace", Value: result.Workspace}, {Name: "persisted", Value: epicCanonical(result.Epic)}, {Name: "observed", Value: observed}, {Name: "freshness", Value: fresh}, {Name: "ready_for_task_worktree_create", Value: result.Ready}, {Name: "reasons", Value: stringValues(result.Reasons)}})
}

func MarshalTaskReadback(result TaskReadbackResult) ([]byte, error) {
	if result.Integration != nil {
		return MarshalTaskIntegrationReadback(*result.Integration)
	}
	persisted := canonicaljson.Object{{Name: "task_id", Value: string(result.Task.ID)}, {Name: "title", Value: result.Task.Title}, {Name: "description", Value: result.Task.Description}, {Name: "parent_epic_id", Value: string(result.Task.ParentEpicID)}, {Name: "project_id", Value: string(result.Task.ProjectID)}, {Name: "repo_id", Value: string(result.Task.RepoID)}, {Name: "git_common_dir", Value: result.Task.GitCommonDir}, {Name: "worktree_state", Value: string(result.Task.WorktreeState)}, {Name: "worktree", Value: taskWorktreeCanonical(result.Task.Worktree)}, {Name: "operation", Value: operationCanonical(result.Operation)}}
	observed := canonicaljson.Object{{Name: "project_git_common_dir", Value: pointerValue(result.ProjectGitCommonDir)}, {Name: "parent", Value: observedWorktreeCanonical(result.Parent)}, {Name: "source", Value: observedSourceCanonical(result.Source)}, {Name: "target", Value: observedTargetCanonical(result.Target)}}
	fresh := canonicaljson.Object{{Name: "project_git_common_dir", Value: result.ProjectFreshness}, {Name: "parent", Value: freshWorktreeCanonical(result.ParentFreshness)}, {Name: "source", Value: freshSourceCanonical(result.SourceFreshness)}, {Name: "target", Value: freshTargetCanonical(result.TargetFreshness)}}
	return canonicaljson.Marshal(canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskReadback@1"}, {Name: "workspace", Value: result.Workspace}, {Name: "persisted", Value: persisted}, {Name: "observed", Value: observed}, {Name: "freshness", Value: fresh}, {Name: "worktree_ready", Value: result.WorktreeReady}, {Name: "ready_for_handoff", Value: result.ReadyForHandoff}, {Name: "reasons", Value: stringValues(result.Reasons)}, {Name: "agent_started_by_this_command", Value: false}})
}

func repoByID(repositories []RepoRecord, id RepoID) (RepoRecord, bool) {
	for _, repository := range repositories {
		if repository.ID == id {
			return repository, true
		}
	}
	return RepoRecord{}, false
}
func freshness(value bool) string {
	if value {
		return "fresh"
	}
	return "stale"
}
func sortedReasons(values []string) []string {
	sort.Strings(values)
	result := []string{}
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
func pointerValue[T any](pointer *T) canonicaljson.Value {
	if pointer == nil {
		return nil
	}
	switch value := any(*pointer).(type) {
	case string:
		return value
	case bool:
		return value
	default:
		return nil
	}
}
func stringValues(values []string) []canonicaljson.Value {
	result := make([]canonicaljson.Value, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
func observedWorktree(value GitWorktreeObservation) TaskObservedWorktree {
	return TaskObservedWorktree{Locator: stringPointer(value.Locator), Ref: stringPointer(value.Ref), OID: stringPointer(value.OID), Tree: stringPointer(value.Tree), GitCommonDir: stringPointer(value.GitCommonDir), Clean: boolPointer(value.Clean), InventoryMatch: boolPointer(value.InventoryMatch)}
}
func freshWorktree(value GitWorktreeObservation, locator, ref, oid, tree, common string) TaskFreshnessWorktree {
	return TaskFreshnessWorktree{Locator: freshness(value.Locator == locator), Ref: freshness(value.Ref == ref), OID: freshness(value.OID == oid), Tree: freshness(value.Tree == tree), GitCommonDir: freshness(value.GitCommonDir == common), Clean: freshness(value.Clean), InventoryMatch: freshness(value.InventoryMatch)}
}
func freshTarget(value GitWorktreeObservation, operation WorktreeOperationRecord) TaskFreshnessTarget {
	return TaskFreshnessTarget{Kind: "fresh", Locator: freshness(value.Locator == operation.TargetLocator), Ref: freshness(value.Ref == operation.SourceRef), OID: freshness(value.OID == operation.ParentOID), Tree: freshness(value.Tree == operation.ParentTree), GitCommonDir: freshness(value.GitCommonDir == operation.GitCommonDir), Clean: freshness(value.Clean), InventoryMatch: freshness(value.InventoryMatch)}
}
func unknownWorktreeFreshness() TaskFreshnessWorktree {
	return TaskFreshnessWorktree{Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}
}
func unknownTargetFreshness() TaskFreshnessTarget {
	return TaskFreshnessTarget{Kind: "unknown", Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}
}

func freshnessForObservedTarget(observation WorktreeObservationRecord, operation WorktreeOperationRecord) TaskFreshnessTarget {
	result := unknownTargetFreshness()
	result.Kind = freshness(observation.TargetKind == "worktree")
	if observation.TargetLocator != nil {
		result.Locator = freshness(*observation.TargetLocator == operation.TargetLocator)
	}
	if observation.TargetRef != nil {
		result.Ref = freshness(*observation.TargetRef == operation.SourceRef)
	}
	if observation.TargetOID != nil {
		result.OID = freshness(*observation.TargetOID == operation.ParentOID)
	}
	if observation.TargetTree != nil {
		result.Tree = freshness(*observation.TargetTree == operation.ParentTree)
	}
	if observation.TargetGitCommonDir != nil {
		result.GitCommonDir = freshness(*observation.TargetGitCommonDir == operation.GitCommonDir)
	}
	if observation.TargetClean != nil {
		result.Clean = freshness(*observation.TargetClean)
	}
	if observation.InventoryMatch != nil {
		result.InventoryMatch = freshness(*observation.InventoryMatch)
	}
	return result
}

func epicCanonical(epic EpicRecord) canonicaljson.Value {
	bindings := []canonicaljson.Value{}
	for _, binding := range epic.RepoBindings {
		bindings = append(bindings, canonicaljson.Object{{Name: "repo_id", Value: string(binding.RepoID)}, {Name: "git_common_dir", Value: binding.GitCommonDir}, {Name: "worktree", Value: canonicaljson.Object{{Name: "id", Value: string(binding.Worktree.ID)}, {Name: "owner_kind", Value: binding.Worktree.OwnerKind}, {Name: "owner_id", Value: string(binding.Worktree.OwnerID)}, {Name: "origin", Value: binding.Worktree.Origin}, {Name: "locator", Value: binding.Worktree.Locator}, {Name: "ref", Value: binding.Worktree.Ref}, {Name: "oid", Value: binding.Worktree.OID}, {Name: "tree", Value: binding.Worktree.Tree}, {Name: "git_common_dir", Value: binding.Worktree.GitCommonDir}}}})
	}
	return canonicaljson.Object{{Name: "epic_id", Value: string(epic.ID)}, {Name: "title", Value: epic.Title}, {Name: "project_id", Value: string(epic.ProjectID)}, {Name: "repo_bindings", Value: bindings}}
}
func taskWorktreeCanonical(worktree *TaskWorktreeBinding) canonicaljson.Value {
	if worktree == nil {
		return nil
	}
	return canonicaljson.Object{{Name: "id", Value: string(worktree.ID)}, {Name: "owner_kind", Value: worktree.OwnerKind}, {Name: "owner_id", Value: string(worktree.OwnerID)}, {Name: "origin", Value: worktree.Origin}, {Name: "locator", Value: worktree.Locator}, {Name: "ref", Value: worktree.Ref}, {Name: "oid", Value: worktree.OID}, {Name: "tree", Value: worktree.Tree}, {Name: "git_common_dir", Value: worktree.GitCommonDir}, {Name: "parent_epic_id", Value: string(worktree.ParentEpicID)}, {Name: "parent_worktree_id", Value: string(worktree.ParentWorktreeID)}, {Name: "parent_ref", Value: worktree.ParentRef}, {Name: "parent_oid", Value: worktree.ParentOID}, {Name: "parent_tree", Value: worktree.ParentTree}, {Name: "operation_id", Value: string(worktree.OperationID)}, {Name: "intent_digest", Value: worktree.IntentDigest}}
}
func operationCanonical(operation *WorktreeOperationRecord) canonicaljson.Value {
	if operation == nil {
		return nil
	}
	o := operation.LastObservation
	last := canonicaljson.Object{{Name: "classification", Value: o.Classification}, {Name: "parent_ref_oid", Value: pointerValue(o.ParentRefOID)}, {Name: "parent_ref_tree", Value: pointerValue(o.ParentRefTree)}, {Name: "parent_worktree_locator", Value: pointerValue(o.ParentWorktreeLocator)}, {Name: "parent_worktree_ref", Value: pointerValue(o.ParentWorktreeRef)}, {Name: "parent_worktree_oid", Value: pointerValue(o.ParentWorktreeOID)}, {Name: "parent_worktree_tree", Value: pointerValue(o.ParentWorktreeTree)}, {Name: "parent_worktree_git_common_dir", Value: pointerValue(o.ParentWorktreeGitCommonDir)}, {Name: "parent_worktree_clean", Value: pointerValue(o.ParentWorktreeClean)}, {Name: "source_ref_oid", Value: pointerValue(o.SourceRefOID)}, {Name: "source_ref_tree", Value: pointerValue(o.SourceRefTree)}, {Name: "source_ref_checked_out_at", Value: pointerValue(o.SourceRefCheckedOutAt)}, {Name: "target_kind", Value: o.TargetKind}, {Name: "target_locator", Value: pointerValue(o.TargetLocator)}, {Name: "target_ref", Value: pointerValue(o.TargetRef)}, {Name: "target_oid", Value: pointerValue(o.TargetOID)}, {Name: "target_tree", Value: pointerValue(o.TargetTree)}, {Name: "target_git_common_dir", Value: pointerValue(o.TargetGitCommonDir)}, {Name: "target_clean", Value: pointerValue(o.TargetClean)}, {Name: "inventory_match", Value: pointerValue(o.InventoryMatch)}, {Name: "reasons", Value: stringValues(o.Reasons)}}
	return canonicaljson.Object{{Name: "id", Value: string(operation.ID)}, {Name: "task_id", Value: string(operation.TaskID)}, {Name: "worktree_id", Value: string(operation.WorktreeID)}, {Name: "intent_digest", Value: operation.IntentDigest}, {Name: "state", Value: operation.State}, {Name: "project_id", Value: string(operation.ProjectID)}, {Name: "repo_id", Value: string(operation.RepoID)}, {Name: "git_common_dir", Value: operation.GitCommonDir}, {Name: "parent_epic_id", Value: string(operation.ParentEpicID)}, {Name: "parent_worktree_id", Value: string(operation.ParentWorktreeID)}, {Name: "parent_ref", Value: operation.ParentRef}, {Name: "parent_oid", Value: operation.ParentOID}, {Name: "parent_tree", Value: operation.ParentTree}, {Name: "source_ref", Value: operation.SourceRef}, {Name: "target_locator", Value: operation.TargetLocator}, {Name: "last_observation", Value: last}}
}
func observedWorktreeCanonical(value TaskObservedWorktree) canonicaljson.Value {
	return canonicaljson.Object{{Name: "locator", Value: pointerValue(value.Locator)}, {Name: "ref", Value: pointerValue(value.Ref)}, {Name: "oid", Value: pointerValue(value.OID)}, {Name: "tree", Value: pointerValue(value.Tree)}, {Name: "git_common_dir", Value: pointerValue(value.GitCommonDir)}, {Name: "clean", Value: pointerValue(value.Clean)}, {Name: "inventory_match", Value: pointerValue(value.InventoryMatch)}}
}
func observedSourceCanonical(value TaskObservedSource) canonicaljson.Value {
	return canonicaljson.Object{{Name: "ref", Value: pointerValue(value.Ref)}, {Name: "oid", Value: pointerValue(value.OID)}, {Name: "tree", Value: pointerValue(value.Tree)}, {Name: "checked_out_at", Value: pointerValue(value.CheckedOutAt)}}
}
func observedTargetCanonical(value TaskObservedTarget) canonicaljson.Value {
	return canonicaljson.Object{{Name: "kind", Value: value.Kind}, {Name: "locator", Value: pointerValue(value.Locator)}, {Name: "ref", Value: pointerValue(value.Ref)}, {Name: "oid", Value: pointerValue(value.OID)}, {Name: "tree", Value: pointerValue(value.Tree)}, {Name: "git_common_dir", Value: pointerValue(value.GitCommonDir)}, {Name: "clean", Value: pointerValue(value.Clean)}, {Name: "inventory_match", Value: pointerValue(value.InventoryMatch)}}
}
func freshWorktreeCanonical(value TaskFreshnessWorktree) canonicaljson.Value {
	return canonicaljson.Object{{Name: "locator", Value: value.Locator}, {Name: "ref", Value: value.Ref}, {Name: "oid", Value: value.OID}, {Name: "tree", Value: value.Tree}, {Name: "git_common_dir", Value: value.GitCommonDir}, {Name: "clean", Value: value.Clean}, {Name: "inventory_match", Value: value.InventoryMatch}}
}
func freshSourceCanonical(value TaskFreshnessSource) canonicaljson.Value {
	return canonicaljson.Object{{Name: "ref", Value: value.Ref}, {Name: "oid", Value: value.OID}, {Name: "tree", Value: value.Tree}, {Name: "checked_out_at", Value: value.CheckedOutAt}}
}
func freshTargetCanonical(value TaskFreshnessTarget) canonicaljson.Value {
	return canonicaljson.Object{{Name: "kind", Value: value.Kind}, {Name: "locator", Value: value.Locator}, {Name: "ref", Value: value.Ref}, {Name: "oid", Value: value.OID}, {Name: "tree", Value: value.Tree}, {Name: "git_common_dir", Value: value.GitCommonDir}, {Name: "clean", Value: value.Clean}, {Name: "inventory_match", Value: value.InventoryMatch}}
}
