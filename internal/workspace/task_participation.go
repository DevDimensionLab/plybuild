package workspace

const (
	ParallelRegistryFormat          = 8
	AdmissionParallel               = "parallel_explicit"
	SourceReconciliationMergeParent = "merge_parent_v1"
)

// Active is derived exclusively from the append-only queue CAS chain. Current
// is a legacy read projection, never a test for whether admission is available.
func (q queueFold) contains(id string) bool {
	for _, member := range q.Active {
		if member == id {
			return true
		}
	}
	return false
}

func (q *queueFold) projectCurrent() {
	q.Current = nil
	if len(q.Active) == 1 {
		id := q.Active[0]
		q.Current = &id
	}
}

func admissionAllowed(r WorkItemRegistry, q queueFold, mode string) bool {
	if mode != AdmissionParallel {
		return len(q.Active) == 0
	}
	if r.FormatVersion != ParallelRegistryFormat {
		return false
	}
	for _, id := range q.Active {
		p := findPreparation(r, id)
		if p == nil || p.Plan.AdmissionMode != AdmissionParallel {
			return false
		}
	}
	return true
}

func preparationActive(disposition string) bool {
	return disposition == "current" || disposition == "active"
}

// PreparationIsActive reads authoritative queue history for native consumers.
func PreparationIsActive(r WorkItemRegistry, id string) bool {
	p := findPreparation(r, id)
	return p != nil && foldQueue(r, p.Plan.QueueID).contains(id)
}

func validGoalAdmission(p TaskGoalExecutePreview) bool {
	if p.SchemaVersion == 1 {
		return p.Kind == "WorkspaceTaskGoalExecutePlan@1" && !p.Input.Parallel && p.AdmissionMode == "" && p.SourceReconciliation == ""
	}
	return p.SchemaVersion == 2 && p.Kind == "WorkspaceTaskGoalExecutePlan@2" && p.Input.Parallel && !p.Input.Next && p.Input.SpecID != "" && p.AdmissionMode == AdmissionParallel && p.SourceReconciliation == SourceReconciliationMergeParent && p.SourceMode == ""
}

func validPreparationAdmission(p WorkspaceTaskPreparePlan, format int) bool {
	if p.SchemaVersion == 1 {
		return p.Kind == "WorkspaceTaskPreparePlan@1" && p.AdmissionMode == "" && p.SourceReconciliation == ""
	}
	return format == ParallelRegistryFormat && p.SchemaVersion == 2 && p.Kind == "WorkspaceTaskPreparePlan@2" && p.AdmissionMode == AdmissionParallel && p.SourceReconciliation == SourceReconciliationMergeParent && p.Selector.Kind == "task" && p.SourceMode == ""
}

// Membership is an exact preparation check. It cannot confer another member's
// owner, provider session, candidate or effect authority.
func (r WorkspaceTaskQueueReadback) HasPreparation(id string) bool {
	if r.Current != nil && r.Current.PreparationID == id {
		return true
	}
	for _, member := range r.Active {
		if member.PreparationID == id {
			return true
		}
	}
	return false
}
