package workspace

import "time"

// LatestTaskHumanQA selects the most recently completed answer for the exact
// candidate. Random record IDs do not establish chronology. Conflicting answers
// at the latest instant are ambiguous and cannot authorize an effect.
func LatestTaskHumanQA(records []TaskHumanQARecord, resultID TaskResultID, oid, tree string) (*TaskHumanQARecord, error) {
	return latestHumanQA(records, func(q TaskHumanQARecord) bool {
		return q.TaskResultID == resultID && q.ResultOID == oid && q.ResultTree == tree
	})
}

// LatestTaskCandidateHumanQA preserves actual judgment across requalification
// of identical Task bytes. It returns the original record and TaskResultID;
// callers must not reinterpret an old positive answer as QA of a new result.
func LatestTaskCandidateHumanQA(records []TaskHumanQARecord, result TaskResultRecord) (*TaskHumanQARecord, error) {
	return latestHumanQA(records, func(q TaskHumanQARecord) bool {
		return q.TaskID == result.TaskID && q.ResultOID == result.ResultOID && q.ResultTree == result.ResultTree
	})
}

// LatestParallelTaskHumanQA keeps a Task's unresolved human veto through source
// reconciliation. Only a later actual pass on this exact result resolves it;
// automatic qualification and an answer on different bytes cannot do so.
func LatestParallelTaskHumanQA(records []TaskHumanQARecord, result TaskResultRecord) (*TaskHumanQARecord, error) {
	exact, err := LatestTaskHumanQA(records, result.ID, result.ResultOID, result.ResultTree)
	if err != nil {
		return nil, err
	}
	veto, err := latestHumanQA(records, func(q TaskHumanQARecord) bool {
		return q.TaskID == result.TaskID && (q.Outcome == "fail" || q.Outcome == "blocked")
	})
	if err != nil || veto == nil {
		return exact, err
	}
	if exact != nil && exact.Outcome == "pass" {
		a, _ := time.Parse(time.RFC3339Nano, exact.Actor.CompletedAtUTC)
		b, _ := time.Parse(time.RFC3339Nano, veto.Actor.CompletedAtUTC)
		if a.After(b) {
			return exact, nil
		}
	}
	return veto, nil
}

func latestRegistryAcceptanceHumanQA(r WorkItemRegistry, result TaskResultRecord) (*TaskHumanQARecord, error) {
	for _, binding := range r.TaskResultSpecBindings {
		if binding.TaskResultID == result.ID && binding.EffectiveBasis != nil {
			return LatestParallelTaskHumanQA(r.HumanQARecords, result)
		}
	}
	return LatestTaskCandidateHumanQA(r.HumanQARecords, result)
}

func latestHumanQA(records []TaskHumanQARecord, matches func(TaskHumanQARecord) bool) (*TaskHumanQARecord, error) {
	var latest *TaskHumanQARecord
	var completed time.Time
	ambiguous := false
	for i := range records {
		q := records[i]
		if !matches(q) {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, q.Actor.CompletedAtUTC)
		if err != nil {
			return nil, workError(ErrorTaskIntegrationBlocked, "candidate human QA has an invalid completion time", err)
		}
		if latest == nil || at.After(completed) {
			latest, completed, ambiguous = &q, at, false
		} else if at.Equal(completed) {
			if q.Outcome != latest.Outcome {
				ambiguous = true
			}
			// Equivalent outcomes use a stable identity irrespective of input order.
			if q.ID > latest.ID {
				latest = &q
			}
		}
	}
	if ambiguous {
		return nil, workError(ErrorTaskIntegrationBlocked, "latest candidate human QA has conflicting outcomes at the same completion time", nil)
	}
	return latest, nil
}
