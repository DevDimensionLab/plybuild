package workspace

// TaskStartGuard is an in-process guard for one approved initial start. It is
// never decoded from caller JSON. Check runs while the ordinary native Project
// and WorkItem locks protect the snapshot; Committed advances only this caller's
// expected registry after its own native publication.
type TaskStartGuard interface {
	Check(string, WorkItemRegistry) error
	Committed(string, WorkItemRegistry)
}

type guardedTaskStartStore struct {
	WorkItemStore
	guard TaskStartGuard
}

func (s guardedTaskStartStore) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		if err = s.guard.Check(root, r); err != nil {
			return err
		}
		err = fn(session)
		if after, readErr := session.Snapshot(); readErr == nil {
			s.guard.Committed(root, after)
		}
		return err
	})
}

func WithTaskStartGuard(d Dependencies, guard TaskStartGuard) Dependencies {
	d.WorkItems = guardedTaskStartStore{d.WorkItems, guard}
	return d
}

// PreviewInitialTaskGoalExecution excludes existing/prepared executions. Those
// belong to their existing identity and cannot become a fresh approved start.
// The caller already holds the native snapshot locks.
func (s *TaskSpecSession) PreviewInitialTaskGoalExecution(in TaskGoalExecuteInput) (TaskGoalExecutePreview, error) {
	if err := validateGoalExecuteInput(in); err != nil {
		return TaskGoalExecutePreview{}, err
	}
	p, err := buildGoalExecutionPreview(s.Dependencies, s.Root, s.Registry, s.Projects, in)
	if err != nil {
		return p, err
	}
	if err := validateGoalOperationContent(s.Dependencies, s.Root, s.Registry, p.Goal.TaskID); err != nil {
		return p, err
	}
	q := foldQueue(s.Registry, p.QueueID)
	if !admissionAllowed(s.Registry, q, p.AdmissionMode) || preparationForTask(s.Registry, p.Goal.TaskID) != nil {
		return p, queueError("task_goal_initial_start_required", "an existing preparation requires its preserved Run, not a fresh bound start")
	}
	for _, entry := range q.Pending {
		if entry.TaskID == p.Goal.TaskID && entry.Goal != nil && *entry.Goal == p.Goal.Goal {
			old, _, err := storedGoalExecutionInTarget(s.Dependencies, s.Root, s.Registry, p.Input, "", p.Target)
			if err != nil {
				return p, err
			}
			if old != nil {
				return p, queueError("task_goal_initial_start_required", "an execution intent already owns this Goal")
			}
			return p, nil
		}
	}
	return p, queueError("task_goal_initial_start_required", "the explicit Goal is no longer pending for an initial start")
}
