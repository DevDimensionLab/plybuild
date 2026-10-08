package workspace

// taskContentScope bounds schema interpretation, never registry serialization.
// Cross-Task Spec dependencies are currently forbidden by the contract schema.
type taskContentScope struct{ task TaskID }

// WithTaskContentScope validates this Task's complete content and delivery
// dependency closure while preserving all shared registrations. It grants no
// effect: ordinary operation, authority, candidate and Git guards still apply.
// An alternate store without scope support retains its stronger full checks.
func WithTaskContentScope(d Dependencies, task TaskID) Dependencies {
	if task == "" {
		return d
	}
	scope := &taskContentScope{task: task}
	if store, ok := d.WorkItems.(*systemWorkItemStore); ok {
		copy := *store
		copy.scope = scope
		d.WorkItems = &copy
	}
	if d.TaskContent != nil {
		copy := *d.TaskContent
		copy.scope = scope
		d.TaskContent = &copy
	}
	return d
}

func (s *TaskContentStorage) includesTask(task TaskID) bool {
	return s == nil || s.scope == nil || s.scope.task == task
}

func (s *TaskContentStorage) includesTarget(r WorkItemRegistry, epic EpicID, repo RepoID) bool {
	return s == nil || s.scope.includesTarget(r, epic, repo)
}

func (s *taskContentScope) includesTarget(r WorkItemRegistry, epic EpicID, repo RepoID) bool {
	if s == nil {
		return true
	}
	task, _ := findTask(r, s.task)
	return task != nil && task.ParentEpicID == epic && task.RepoID == repo
}

func outsideScope[T any](records []T, owned func(T) bool) []T {
	var out []T
	for _, record := range records {
		if !owned(record) {
			out = append(out, record)
		}
	}
	return out
}

// The projection exists only for comparing what a scoped writer must preserve.
// Publication always encodes the full registry supplied by the locked session.
func outsideTaskContentScope(r WorkItemRegistry, scope *taskContentScope) WorkItemRegistry {
	out := r
	out.RawSHA256 = ""
	id := scope.task
	out.Tasks = outsideScope(r.Tasks, func(v TaskRecord) bool { return v.ID == id })
	out.WorktreeOperations = outsideScope(r.WorktreeOperations, func(v WorktreeOperationRecord) bool { return v.TaskID == id })
	out.TaskResults = outsideScope(r.TaskResults, func(v TaskResultRecord) bool { return v.TaskID == id })
	out.HumanQARecords = outsideScope(r.HumanQARecords, func(v TaskHumanQARecord) bool { return v.TaskID == id })
	out.IntegrationAuthorities = outsideScope(r.IntegrationAuthorities, func(v IntegrationAuthority) bool { return v.TaskID == id })
	out.IntegrationIntents = outsideScope(r.IntegrationIntents, func(v IntegrationIntent) bool { return v.TaskID == id })
	ownedAuthority := func(id IntegrationAuthorityID) bool {
		a := findAuthority(r, id)
		return a != nil && a.TaskID == scope.task
	}
	out.IntegrationAttempts = outsideScope(r.IntegrationAttempts, func(v IntegrationAttempt) bool { return ownedAuthority(v.AuthorityID) })
	out.IntegrationResults = outsideScope(r.IntegrationResults, func(v IntegrationResult) bool { return ownedAuthority(v.AuthorityID) })
	out.TaskSpecPolicies = outsideScope(r.TaskSpecPolicies, func(v TaskSpecPolicy) bool { return v.TaskID == id })
	out.TaskProblemRevisions = outsideScope(r.TaskProblemRevisions, func(v TaskProblemReference) bool { return v.TaskID == id })
	out.TaskSpecRevisions = outsideScope(r.TaskSpecRevisions, func(v TaskSpecReference) bool { return v.TaskID == id })
	out.TaskSpecAssessments = outsideScope(r.TaskSpecAssessments, func(v TaskAssessmentReference) bool { return v.TaskID == id })
	out.TaskSolutionSelections = outsideScope(r.TaskSolutionSelections, func(v TaskSelectionReference) bool { return v.TaskID == id })
	out.TaskResultSpecBindings = outsideScope(r.TaskResultSpecBindings, func(v TaskResultSpecBinding) bool { return v.TaskID == id })
	out.TaskContentPublications = outsideScope(r.TaskContentPublications, func(v TaskContentPublication) bool { return v.TaskID == id })
	out.TaskPreparations = outsideScope(r.TaskPreparations, func(v TaskPreparation) bool { return v.Plan.TaskID == id })
	out.EpicBaseVersions = outsideScope(r.EpicBaseVersions, func(v EpicBaseVersion) bool { return scope.includesTarget(r, v.EpicID, v.RepoID) })
	out.EpicBaseUpdates = outsideScope(r.EpicBaseUpdates, func(v EpicBaseUpdate) bool {
		return scope.includesTarget(r, v.Plan.Target.EpicID, v.Plan.Target.RepoID)
	})
	// Queue set/reordering is a separate multi-Task operation and is never
	// permitted by a Task scope. Only this Task's preparation transitions qualify.
	out.TaskQueueEvents = outsideScope(r.TaskQueueEvents, func(v TaskQueueEvent) bool {
		if v.Kind == "set" {
			return false
		}
		preparation, _ := v.Request["preparation_id"].(string)
		p := findPreparation(r, preparation)
		return p != nil && p.Plan.TaskID == id
	})
	return out
}

func validateTaskScopeTransition(old, next WorkItemRegistry, scope *taskContentScope) error {
	if scope == nil || old.FormatVersion < 3 {
		return nil
	}
	before, _ := findTask(old, scope.task)
	after, _ := findTask(next, scope.task)
	if before != nil && (after == nil || before.ParentEpicID != after.ParentEpicID || before.ProjectID != after.ProjectID || before.RepoID != after.RepoID || before.GitCommonDir != after.GitCommonDir) {
		return contentError("task_content_scope_conflict", "Task operation cannot change its registered ownership", nil)
	}
	if !contentTypedEqual(outsideTaskContentScope(old, scope), outsideTaskContentScope(next, scope)) {
		return contentError("task_content_scope_conflict", "Task operation changed records outside its Task and registered Epic target", nil)
	}
	return nil
}
