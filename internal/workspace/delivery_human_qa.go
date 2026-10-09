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
