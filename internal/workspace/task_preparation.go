package workspace

import (
	"fmt"
	"strings"
)

type TaskPrepareInput struct {
	Target       QueueTargetInput
	Selector     QueueSelector
	Apply        bool
	Confirmation string
	Upgrade      string
}
type TaskPreparationReadback struct {
	PersistedBefore *QueueRegistryReadback    `json:"persisted_before"`
	PersistedAfter  *QueueRegistryReadback    `json:"persisted_after"`
	Kind            string                    `json:"kind"`
	SchemaVersion   int                       `json:"schema_version"`
	Workspace       string                    `json:"workspace"`
	Target          QueueTarget               `json:"target"`
	State           string                    `json:"state"`
	Disposition     string                    `json:"disposition"`
	Plan            *WorkspaceTaskPreparePlan `json:"plan"`
	Confirmation    *string                   `json:"confirmation"`
	Preparation     *TaskPreparation          `json:"preparation"`
	Observed        *GitWorktreeObservation   `json:"observed"`
	Freshness       string                    `json:"freshness"`
	Reasons         []QueueReason             `json:"reasons"`
	BackupSHA256    *string                   `json:"backup_sha256"`
}

func preparationReadback(d Dependencies, root string, r WorkItemRegistry, p TaskPreparation) TaskPreparationReadback {
	cur := currentPreparation(d, r, p)
	out := TaskPreparationReadback{Kind: "WorkspaceTaskPreparationReadback@1", SchemaVersion: 1, Workspace: root, Target: p.Plan.Target, State: cur.State, Disposition: preparationDisposition(r, p.ID), Plan: &p.Plan, Confirmation: &p.PlanSHA256, Preparation: &p, Freshness: "unknown", Reasons: cur.Reasons}
	if o, e := d.WorkGit.ObserveWorktree(p.Plan.WorktreePath); e == nil {
		out.Observed = &o
		out.Freshness = freshness(o.Locator == p.Plan.WorktreePath && o.Ref == "refs/heads/"+p.Plan.Branch && o.GitCommonDir == p.Plan.Target.GitCommonDir && o.InventoryMatch)
		if !o.Clean || o.OID != p.Plan.ParentOID {
			out.Reasons = append(out.Reasons, QueueReason{"task_worktree_changed", "worktree has changed since preparation; the prepared outcome is historical"})
		}
	}
	if ep, _ := findEpic(r, p.Plan.Target.EpicID); ep != nil {
		b, _ := findEpicRepo(*ep, p.Plan.Target.RepoID)
		if currentEpicBase(r, *ep, *b).OID != p.Plan.ParentOID {
			out.Reasons = append(out.Reasons, QueueReason{"epic_base_changed", "preparation belongs to a historical Epic base"})
		}
	}
	if p.Plan.RegistryUpgrade != nil {
		out.BackupSHA256 = &p.Plan.RegistryUpgrade.RegistrySHA256
	}
	out.Reasons = sortQueueReasons(out.Reasons)
	return out
}
func validatePrepareSelector(s QueueSelector) error {
	if s.Kind == "next" && s.TaskID == nil {
		return nil
	}
	if s.Kind == "task" && s.TaskID != nil {
		_, e := ParseTaskID(string(*s.TaskID))
		return e
	}
	return WorkInvalidArguments("choose exactly one of --next or <task-id>")
}
func buildPreparePlan(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in TaskPrepareInput) (TaskPreparationReadback, error) {
	out := TaskPreparationReadback{Kind: "WorkspaceTaskPreparationReadback@1", SchemaVersion: 1, Workspace: root, State: "empty", Disposition: "none", Freshness: "fresh", Reasons: []QueueReason{}}
	q, e := buildQueueReadback(d, root, r, projects, in.Target, false)
	if e != nil {
		return out, e
	}
	out.Target = q.Target
	if q.Current != nil {
		p := findPreparation(r, q.Current.PreparationID)
		if in.Selector.Kind == "task" && *in.Selector.TaskID != p.Plan.TaskID {
			return out, queueError("task_queue_current_conflict", "handle current preparation "+p.ID+" before selecting another Task")
		}
		return preparationReadback(d, root, r, *p), nil
	}
	fold := foldQueue(r, q.QueueID)
	found := false
	readyGoal := false
	for i, entry := range fold.Pending {
		if in.Selector.Kind == "task" && entry.TaskID != *in.Selector.TaskID {
			continue
		}
		found = true
		row, p := queueMixedEntryEvaluation(d, root, r, projects, q.Target, q.QueueID, entry, i+1)
		if p == nil {
			if entry.Goal != nil && row.State == "ready" {
				readyGoal = true
				row.State = "blocked"
				row.Reasons = queueReasons("task_goal_requires_execute", "choose this base-free goal with ply workflow execute; prepare requires a selected execution solution")
			}
			if in.Selector.Kind == "task" {
				out.State = row.State
				out.Reasons = row.Reasons
				return out, nil
			}
			continue
		}
		p.QueueRevision = q.Revision
		p.Selector = in.Selector
		digest := queueDigest(p)
		out.State = "preview"
		out.Plan = p
		out.Confirmation = &digest
		return out, nil
	}
	if !found && in.Selector.Kind == "task" {
		return out, queueError("task_queue_task_not_pending", "named Task is not an active pending entry")
	}
	if readyGoal {
		out.State = "blocked"
		out.Reasons = queueReasons("task_goal_requires_execute", "choose the next base-free goal with ply workflow execute; prepare requires a selected execution solution")
		return out, nil
	}
	out.Reasons = queueReasons("task_queue_no_ready", "no eligible pending Task; inspect the queue reasons")
	return out, nil
}
func PrepareTask(d Dependencies, in TaskPrepareInput) (TaskPreparationReadback, error) {
	if e := validatePrepareSelector(in.Selector); e != nil {
		return TaskPreparationReadback{}, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return TaskPreparationReadback{}, e
	}
	var out TaskPreparationReadback
	if !in.Apply {
		e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
			var e error
			out, e = buildPreparePlan(d, root, s.Registry, s.Projects, in)
			return e
		})
		return out, e
	}
	if !digestPattern.MatchString(in.Confirmation) {
		return out, queueError("task_prepare_confirmation_required", "apply requires the exact preview digest")
	}
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			before := QueueRegistryReadback{r.FormatVersion, r.RawSHA256, queueUpgrade(r)}
			defer func() {
				out.PersistedBefore = &before
				if after, err := session.Snapshot(); err == nil {
					out.PersistedAfter = &QueueRegistryReadback{after.FormatVersion, after.RawSHA256, queueUpgrade(after)}
				}
			}()
			id := "pre_" + strings.TrimPrefix(in.Confirmation, "sha256:")
			t, _, e := resolveQueueTarget(d, root, r, projects, in.Target)
			if e != nil {
				return e
			}
			p := findPreparation(r, id)
			if p != nil {
				if p.Plan.Target != t || !contentTypedEqual(p.Plan.Selector, in.Selector) {
					return queueError("task_prepare_conflict", "confirmation belongs to another selector or target")
				}
				out = preparationReadback(d, root, r, *p)
				if p.Outcome != nil || out.Disposition != "current" {
					return nil
				}
			} else {
				out, e = buildPreparePlan(d, root, r, projects, in)
				if e != nil {
					return e
				}
				if out.Confirmation == nil || *out.Confirmation != in.Confirmation || out.Preparation != nil {
					return queueError("task_prepare_conflict", "preview changed; no Task was reserved")
				}
				plan := *out.Plan
				if plan.RegistryUpgrade != nil && in.Upgrade != plan.RegistryUpgrade.RegistrySHA256 {
					return queueError("registry_upgrade_required", "confirm the exact registry upgrade digest")
				}
				if plan.RegistryUpgrade == nil && in.Upgrade != "" {
					return queueError("registry_upgrade_conflict", "format 4 needs no upgrade")
				}
				if e = upgradeQueueRegistry(d, root, &r, plan.RegistryUpgrade); e != nil {
					return e
				}
				p = &TaskPreparation{ID: id, PlanSHA256: in.Confirmation, Plan: plan, CreatedAtUTC: queueNow(d)}
			}
			plan := p.Plan
			input := TaskWorktreeCreateInput{TaskID: plan.TaskID, Branch: plan.Branch, Path: plan.WorktreePath, ExpectedParentOID: plan.ParentOID}
			target, exists, e := canonicalTaskTarget(d.Files, root, plan.WorktreePath)
			if e != nil {
				return e
			}
			var result TaskWorktreeMutationResult
			reserve := func(reg *WorkItemRegistry, op WorktreeOperationRecord) error {
				if findPreparation(*reg, id) != nil {
					return nil
				}
				q := foldQueue(*reg, plan.QueueID)
				if q.Current != nil || q.Revision != plan.QueueRevision || q.Revision >= 2147483647 {
					return queueError("task_queue_revision_conflict", "queue changed before reservation")
				}
				p.OperationID = op.ID
				reg.TaskPreparations = append(reg.TaskPreparations, *p)
				request := map[string]any{"preparation_id": id, "expected_revision": q.Revision}
				reg.TaskQueueEvents = append(reg.TaskQueueEvents, TaskQueueEvent{plan.QueueID, q.Revision + 1, "reserve", queueDigest(request), request, queueNow(d)})
				return nil
			}
			effectErr := createTaskWorktreeLocked(d, root, projects, session, r, input, target, exists, &result, reserve)
			after, readErr := session.Snapshot()
			if readErr != nil {
				out.State = "unknown"
				return readErr
			}
			if actual := findPreparation(after, id); actual != nil {
				out = preparationReadback(d, root, after, *actual)
			} else {
				out.Preparation = nil
				out.State = "not_recorded"
			}
			if effectErr != nil {
				return fmt.Errorf("%w; inspect ply workspace task preparation show %s --format json", effectErr, id)
			}
			return nil
		})
	})
	return out, e
}
func ShowTaskPreparation(d Dependencies, id string) (TaskPreparationReadback, error) {
	var out TaskPreparationReadback
	if e := queuePrefixedID("pre_")(id); e != nil {
		return out, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		if e := validateQueueClosure(d, root, s.Registry); e != nil {
			return e
		}
		p := findPreparation(s.Registry, id)
		if p == nil {
			return queueError("task_preparation_not_found", "preparation is not recorded")
		}
		out = preparationReadback(d, root, s.Registry, *p)
		return nil
	})
	return out, e
}

// TaskRunPreparation reads the selected preparation under the caller's existing
// Project/work-item snapshot locks. It never upgrades or repairs a registry.
func (s *TaskSpecSession) TaskRunPreparation(id string) (TaskPreparationReadback, error) {
	if e := queuePrefixedID("pre_")(id); e != nil {
		return TaskPreparationReadback{}, e
	}
	p := findPreparation(s.Registry, id)
	if p == nil {
		return TaskPreparationReadback{}, queueError("task_preparation_not_found", "preparation is not recorded")
	}
	d := WithTaskContentScope(s.Dependencies, p.Plan.TaskID)
	if e := validateTaskContentClosure(d.TaskContent, s.Root, s.Registry, false); e != nil {
		return TaskPreparationReadback{}, e
	}
	if e := validateQueueClosure(d, s.Root, s.Registry); e != nil {
		return TaskPreparationReadback{}, e
	}
	return preparationReadback(d, s.Root, s.Registry, *p), nil
}
