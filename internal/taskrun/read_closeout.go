package taskrun

import "github.com/devdimensionlab/plybuild/internal/workspace"

// A completed receipt supplies historical resource state, not live authority.
// Exact candidate and source identity prevent another Task's closeout from
// concealing a missing or changed source in an unrelated workflow.
func workflowCompletedCloseout(root string, s workflowState) (*workspace.TaskCloseoutReceipt, error) {
	r, err := workflowTaskCloseout(root, s)
	if err != nil || r == nil || r.State != "complete" {
		return nil, err
	}
	return r, nil
}

func workflowTaskCloseout(root string, s workflowState) (*workspace.TaskCloseoutReceipt, error) {
	if s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
		return nil, nil
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	r, e := workspace.ReadTaskCloseoutAt(root, c.TaskResult.TaskID)
	if e != nil || r == nil {
		return nil, e
	}
	p := r.Plan
	if p.TaskResultID != c.TaskResult.ID || p.ResultOID != c.OID || p.ResultTree != c.Tree || p.Source.Locator != c.TaskResult.SourceLocator || p.Source.Ref != c.TaskResult.SourceRef {
		return nil, workflowError(4, "Task closeout does not bind this workflow's exact candidate and source")
	}
	return r, nil
}
