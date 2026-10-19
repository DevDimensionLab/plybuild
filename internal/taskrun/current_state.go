package taskrun

import "github.com/devdimensionlab/plybuild/internal/currentstate"

func currentRunAction(d Dependencies, s workflowState, o WorkflowRun) WorkflowAction {
	r := currentstate.Run{SavedAction: currentstate.Action(o.NextAction), HasReasons: len(o.Reasons) != 0, Completed: o.Delivery != nil && o.Delivery.Phase == "completed"}
	if w := o.Worker; w != nil {
		r.Worker = &currentstate.Worker{Phase: w.Phase, Stage: w.Stage, Reason: w.Reason, ReasonCode: w.ReasonCode, RetryNotBefore: w.RetryNotBefore, ControlAction: w.ControlAction}
		r.Answer = productCurrentFacts(d, s, o)
	}
	return WorkflowAction(currentstate.RunAction(r))
}
