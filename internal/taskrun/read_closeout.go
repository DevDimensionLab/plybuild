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
	return workflowTaskCloseoutWithRead(root, s, workspace.ReadTaskCloseoutAt)
}

func workflowTaskCloseoutWithRead(root string, s workflowState, read func(string, workspace.TaskID) (*workspace.TaskCloseoutReceipt, error)) (*workspace.TaskCloseoutReceipt, error) {
	if s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
		return nil, nil
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	r, e := read(root, c.TaskResult.TaskID)
	if e != nil || r == nil {
		return nil, e
	}
	p := r.Plan
	if p.TaskResultID != c.TaskResult.ID || p.ResultOID != c.OID || p.ResultTree != c.Tree || p.Source.Locator != c.TaskResult.SourceLocator || p.Source.Ref != c.TaskResult.SourceRef {
		return nil, workflowError(4, "Task closeout does not bind this workflow's exact candidate and source")
	}
	return r, nil
}

// This is read guidance under the frozen actor model, never authority to run an
// effect. Human-owned delivery retains its explicit owner; ordinary released
// automatic delivery continues through the coordinator outside the source.
func workflowCloseoutAction(s workflowState, r workspace.TaskCloseoutReceipt) WorkflowAction {
	if r.State == "complete" {
		return WorkflowAction{"none", "The exact Task delivery is closed; its source resource is " + r.ResourceState + ". No action is required; retained evidence remains available."}
	}
	_, kind, detail := workspace.TaskCloseoutPendingAction(r)
	actor := "coordinator"
	if a := deliveryEffectiveAgreement(s); a != nil && a.HumanOwnedIntegration() || s.Result.Delivery != nil && s.Result.Delivery.HumanIntegration != nil && s.Result.Delivery.HumanIntegration.Accepted {
		actor = "user"
	} else if headlessRun(s.Request) {
		actor = "recipient"
	}
	if kind == "resume_closeout" && automaticCleanupPending(s) {
		return automaticCloseoutPendingAction(s, detail)
	}
	return WorkflowAction{actor, detail}
}
