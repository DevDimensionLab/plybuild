package workspace

import "time"

// TaskRecordedProgressFacts contains only validated registry facts. It is not an
// integration plan, a check of the current checkout, or a new QA attestation.
type TaskRecordedProgressFacts struct {
	HasProgress        bool
	ResultCount        int
	TaskResult         *TaskResultRecord
	HumanQA            *TaskHumanQARecord
	IntegrationResult  *IntegrationResult
	IntegrationPending bool
	Ambiguous          bool
	ResultBasisStale   bool
	LastActivityUTC    *string
}

// RecordedTaskProgressFacts groups registry history once. It selects an
// unambiguous result for the registered Task basis before following an explicit
// integration authority; a historical green result cannot hide competing work.
// A draft Problem or Spec alone is deliberately not execution progress.
func RecordedTaskProgressFacts(r WorkItemRegistry) map[TaskID]TaskRecordedProgressFacts {
	relevant := recordedResultRelevance(r)
	groups := make(map[TaskID]*WorkItemRegistry, len(r.Tasks))
	out := make(map[TaskID]TaskRecordedProgressFacts, len(r.Tasks))
	for _, task := range r.Tasks {
		groups[task.ID] = &WorkItemRegistry{}
		out[task.ID] = TaskRecordedProgressFacts{HasProgress: task.Worktree != nil || task.WorktreeState == WorkItemCreating || task.WorktreeState == WorkItemReconciliationRequired}
	}
	for _, result := range r.TaskResults {
		if group := groups[result.TaskID]; group != nil {
			group.TaskResults = append(group.TaskResults, result)
		}
	}
	for _, qa := range r.HumanQARecords {
		if group := groups[qa.TaskID]; group != nil {
			group.HumanQARecords = append(group.HumanQARecords, qa)
		}
	}
	authorityTasks := map[IntegrationAuthorityID]TaskID{}
	for _, authority := range r.IntegrationAuthorities {
		if group := groups[authority.TaskID]; group != nil {
			group.IntegrationAuthorities = append(group.IntegrationAuthorities, authority)
			authorityTasks[authority.ID] = authority.TaskID
		}
	}
	for _, result := range r.IntegrationResults {
		if group := groups[authorityTasks[result.AuthorityID]]; group != nil {
			group.IntegrationResults = append(group.IntegrationResults, result)
		}
	}
	for _, preparation := range r.TaskPreparations {
		id := preparation.Plan.TaskID
		if fact, exists := out[id]; exists {
			fact.HasProgress = true
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, preparation.CreatedAtUTC)
			if preparation.Outcome != nil {
				fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, preparation.Outcome.RecordedAtUTC)
			}
			out[id] = fact
		}
	}
	for _, publication := range r.TaskContentPublications {
		if fact, exists := out[publication.TaskID]; exists {
			// Authoring is activity, but not evidence of execution progress.
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, publication.RecordedAtUTC)
			out[publication.TaskID] = fact
		}
	}
	for id, group := range groups {
		fact := out[id]
		fact.ResultCount = len(group.TaskResults)
		fact.HasProgress = fact.HasProgress || fact.ResultCount != 0 || len(group.IntegrationAuthorities) != 0
		current := currentRecordedResults(*group, relevant)
		fact.ResultBasisStale = fact.ResultCount != 0 && len(current.TaskResults) == 0
		var authority *IntegrationAuthority
		if len(current.TaskResults) > 1 {
			// Recorder timestamps and a technically green gate do not select the
			// intended delivery among independent results for the same basis.
			fact.Ambiguous = true
		} else if len(current.TaskResults) == 1 {
			fact.TaskResult, fact.HumanQA, authority, fact.Ambiguous = selectTaskIntegrationLeaf(current, id)
			if fact.TaskResult == nil && !fact.Ambiguous {
				// Failed/unknown is progress too, without qualifying for human QA.
				fact.TaskResult = &current.TaskResults[0]
			}
		}
		if fact.TaskResult != nil && authority == nil && !fact.Ambiguous {
			matching := []*TaskHumanQARecord{}
			for i := range group.HumanQARecords {
				q := &group.HumanQARecords[i]
				if q.TaskResultID == fact.TaskResult.ID && q.ResultOID == fact.TaskResult.ResultOID && q.ResultTree == fact.TaskResult.ResultTree {
					matching = append(matching, q)
				}
			}
			if len(matching) == 1 {
				fact.HumanQA = matching[0]
			} else if len(matching) > 1 {
				fact.Ambiguous = true
				fact.HumanQA = nil
			}
		}
		if authority != nil && !fact.Ambiguous {
			fact.IntegrationResult = findIntegrationResultForAuthority(*group, authority.ID)
			fact.IntegrationPending = fact.IntegrationResult == nil
		}
		if fact.Ambiguous {
			fact.TaskResult, fact.HumanQA = nil, nil
		}
		for _, result := range group.TaskResults {
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, result.Recorder.RecordedAtUTC)
		}
		for _, qa := range group.HumanQARecords {
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, qa.Actor.CompletedAtUTC)
		}
		for _, authority := range group.IntegrationAuthorities {
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, authority.CreatedAtUTC)
		}
		for _, result := range group.IntegrationResults {
			fact.LastActivityUTC = recordedLatestTime(fact.LastActivityUTC, result.RecordedAtUTC)
		}
		out[id] = fact
	}
	return out
}

// Content revision/decision ordinals are registered ordering facts. Recorder
// clocks are not. A new unselected Spec revision deliberately leaves the selected
// basis unchanged, as does an Epic base advance after successful integration.
func recordedResultRelevance(r WorkItemRegistry) map[TaskResultID]bool {
	required := map[TaskID]bool{}
	problems := map[TaskID]TaskProblemReference{}
	selections := map[TaskID]TaskSelectionReference{}
	type assessmentKey struct {
		task     TaskID
		spec     string
		revision int
	}
	assessments := map[assessmentKey]TaskAssessmentReference{}
	links := map[TaskResultID]TaskResultSpecBinding{}
	for _, policy := range r.TaskSpecPolicies {
		required[policy.TaskID] = policy.Mode == "spec_required"
	}
	for _, problem := range r.TaskProblemRevisions {
		if problem.Revision > problems[problem.TaskID].Revision {
			problems[problem.TaskID] = problem
		}
	}
	for _, selection := range r.TaskSolutionSelections {
		if selection.Ordinal > selections[selection.TaskID].Ordinal {
			selections[selection.TaskID] = selection
		}
	}
	for _, assessment := range r.TaskSpecAssessments {
		key := assessmentKey{assessment.TaskID, assessment.SpecID, assessment.SpecRevision}
		if assessment.Ordinal > assessments[key].Ordinal {
			assessments[key] = assessment
		}
	}
	for _, link := range r.TaskResultSpecBindings {
		links[link.TaskResultID] = link
	}
	relevant := make(map[TaskResultID]bool, len(r.TaskResults))
	for _, result := range r.TaskResults {
		if !required[result.TaskID] {
			relevant[result.ID] = true
			continue
		}
		link, exists := links[result.ID]
		if !exists {
			// A preserved legacy result does not establish delivery of the
			// subsequently activated selected Spec.
			continue
		}
		basis := link.Basis
		problem := problems[result.TaskID]
		selection := selections[result.TaskID]
		assessment := assessments[assessmentKey{result.TaskID, basis.SpecID, basis.Spec.Revision}]
		relevant[result.ID] = problem.Revision == basis.Problem.Revision && problem.ManifestSHA256 == basis.Problem.ManifestSHA256 &&
			selection.ID == basis.Selection.ID && selection.ManifestSHA256 == basis.Selection.ManifestSHA256 &&
			assessment.ID == basis.Assessment.ID && assessment.ManifestSHA256 == basis.Assessment.ManifestSHA256
	}
	return relevant
}

func currentRecordedResults(group WorkItemRegistry, relevant map[TaskResultID]bool) WorkItemRegistry {
	current := WorkItemRegistry{HumanQARecords: group.HumanQARecords, IntegrationResults: group.IntegrationResults}
	for _, result := range group.TaskResults {
		if relevant[result.ID] {
			current.TaskResults = append(current.TaskResults, result)
		}
	}
	for _, authority := range group.IntegrationAuthorities {
		if relevant[authority.TaskResultID] {
			current.IntegrationAuthorities = append(current.IntegrationAuthorities, authority)
		}
	}
	return current
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
