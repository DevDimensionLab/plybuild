package workspace

// goalOperationContent narrows external dependency validation, not write
// authority. Queue materialization legitimately changes multiple Tasks' shared
// queue records, so WithTaskContentScope's write fence cannot be used here.
// All observations still decode and structurally validate the whole registry.
func goalOperationContent(d Dependencies, task TaskID) Dependencies {
	return goalOperationContentScope(d, &taskContentScope{task: task})
}

func goalOperationContentScope(d Dependencies, scope *taskContentScope) Dependencies {
	var narrow func(WorkItemStore) WorkItemStore
	narrow = func(store WorkItemStore) WorkItemStore {
		switch s := store.(type) {
		case *systemWorkItemStore:
			copy := *s
			copy.contentScope = scope
			return &copy
		default:
			return store // Alternate stores retain their stronger validation.
		}
	}
	d.WorkItems = narrow(d.WorkItems)
	if d.TaskContent != nil {
		copy := *d.TaskContent
		copy.scope = scope
		d.TaskContent = &copy
	}
	return d
}

// Begin with structural validation under the existing operation lock, then bind
// and validate the Task resolved by the planner before any publication. The
// pointer is private to this synchronous operation; it retains no content bytes.
func beginGoalOperationContent(d Dependencies) (Dependencies, func(string, WorkItemRegistry, TaskID) error) {
	scope := &taskContentScope{}
	d = goalOperationContentScope(d, scope)
	return d, func(root string, r WorkItemRegistry, task TaskID) error {
		if task == "" || (scope.task != "" && scope.task != task) {
			return queueError("task_goal_execution_conflict", "selected Task changed before execution")
		}
		scope.task = task
		return validateGoalOperationContent(d, root, r, task)
	}
}

func validateGoalOperationContent(d Dependencies, root string, r WorkItemRegistry, task TaskID) error {
	d = goalOperationContent(d, task)
	storage := observeTaskContent(d.TaskContent)
	if err := validateTaskContentClosure(storage, root, r, false); err != nil {
		return err
	}
	return validateQueueClosure(Dependencies{TaskContent: storage}, root, r)
}

// Recovery candidates come from authoritative target/queue relationships before
// external reads. All pending entries are retained for interrupted --next
// recovery, even entries after the newly eligible first Goal.
func TaskStartupDependencySHA256(r WorkItemRegistry, task TaskID) string {
	r.RawSHA256 = ""
	// Keep all shared ownership, lifecycle, queue, preparation, integration and
	// base records. Only external content registrations owned by another Task
	// are omitted. This deliberately retains conservative conflicts for shared
	// records; it is not a general task-only view of the registry.
	r.TaskProblemRevisions = outsideScope(r.TaskProblemRevisions, func(v TaskProblemReference) bool { return v.TaskID != task })
	r.TaskSpecRevisions = outsideScope(r.TaskSpecRevisions, func(v TaskSpecReference) bool { return v.TaskID != task })
	r.TaskSpecAssessments = outsideScope(r.TaskSpecAssessments, func(v TaskAssessmentReference) bool { return v.TaskID != task })
	r.TaskSolutionSelections = outsideScope(r.TaskSolutionSelections, func(v TaskSelectionReference) bool { return v.TaskID != task })
	r.TaskContentPublications = outsideScope(r.TaskContentPublications, func(v TaskContentPublication) bool { return v.TaskID != task })
	return queueDigest(struct {
		Kind     string
		Registry WorkItemRegistry
	}{"TaskStartupDependencies@1", r})
}
