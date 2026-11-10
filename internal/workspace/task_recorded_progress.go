package workspace

import "time"

// TaskRecordedProgressFacts contains only validated registry facts. It is not an
// integration plan or a check of the current checkout.
type TaskRecordedProgressFacts struct {
	HasProgress        bool
	IntegrationResult  *IntegrationResult
	IntegrationPending bool
	Ambiguous          bool
	LastActivityUTC    *string
}

// RecordedTaskProgressFacts groups registry history once and follows the Task's
// unambiguous integration leaf. A draft Problem or Spec alone is deliberately
// not execution progress.
func RecordedTaskProgressFacts(r WorkItemRegistry) map[TaskID]TaskRecordedProgressFacts {
	out := make(map[TaskID]TaskRecordedProgressFacts, len(r.Tasks))
	for _, task := range r.Tasks {
		fact := TaskRecordedProgressFacts{HasProgress: task.Worktree != nil || task.WorktreeState == WorkItemCreating || task.WorktreeState == WorkItemReconciliationRequired}
		authority, ambiguous := selectTaskIntegrationLeaf(r, task.ID)
		fact.Ambiguous = ambiguous
		if authority != nil {
			fact.HasProgress = true
			fact.IntegrationResult = findIntegrationResultForAuthority(r, authority.ID)
			fact.IntegrationPending = fact.IntegrationResult == nil
		}
		out[task.ID] = fact
	}
	for _, publication := range r.TaskContentPublications {
		if fact, exists := out[publication.TaskID]; exists {
			// Authoring is activity, but not evidence of execution progress.
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, publication.RecordedAtUTC)
			out[publication.TaskID] = fact
		}
	}
	for _, authority := range r.IntegrationAuthorities {
		if fact, exists := out[authority.TaskID]; exists {
			fact.HasProgress = true
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, authority.CreatedAtUTC)
			if result := findIntegrationResultForAuthority(r, authority.ID); result != nil {
				fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, result.RecordedAtUTC)
			}
			out[authority.TaskID] = fact
		}
	}
	return out
}

func recordedLatestTime(previous *string, candidate string) *string {
	current, err := time.Parse(time.RFC3339Nano, candidate)
	if err != nil {
		return previous
	}
	if previous != nil {
		old, err := time.Parse(time.RFC3339Nano, *previous)
		if err == nil && !current.After(old) {
			return previous
		}
	}
	normalized := current.UTC().Format(time.RFC3339Nano)
	return &normalized
}
