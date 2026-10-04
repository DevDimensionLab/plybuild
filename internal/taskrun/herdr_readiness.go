package taskrun

import "time"

// This loop belongs only to the owner of the one reserved launch. Repeated apply,
// show and follow never enter it and cannot bind a previously unbound session.
func workflowAwaitReadiness(d Dependencies, root, id string, deadline time.Time) error {
	observedSession := ""
	for time.Now().Before(deadline) {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		a, e := workflowAgentRead(d, s, deadline)
		if e != nil {
			return e
		}
		if e = workflowLocation(s, a); e != nil {
			return e
		}
		if a.Session != nil {
			if e = workflowIdentity(s, a, true); e != nil {
				return e
			}
			if observedSession != "" && observedSession != a.Session.Value {
				return workflowError(4, "native Codex session changed during startup; no prompt sent")
			}
			observedSession = a.Session.Value
		}
		ready := a.Session != nil && workflowSettled(a)
		if e = workflowUpdate(d, root, id, func(s *workflowState) error {
			if s.Phase != "agent_start_attempted" {
				return workflowError(4, "launch attempt is no longer awaiting readiness")
			}
			workflowObserve(d, s, a)
			if ready {
				if e := workflowIdentity(*s, a, true); e != nil {
					return e
				}
				if !time.Now().Before(deadline) {
					return workflowError(5, "Herdr startup deadline expired before session binding; no prompt sent")
				}
				s.Result.Transport.AgentSessionID = a.Session.Value
				s.Phase = "session_bound"
			}
			return nil
		}); e != nil {
			return e
		}
		if ready {
			return nil
		}
		// A missing/false launch_pending flag is only an observation. Herdr can
		// report idle before it discovers the Codex session after an early exit.
		// Keep observing this one identity within the original shared deadline.
		pause := 250 * time.Millisecond
		if remaining := time.Until(deadline); remaining < pause {
			pause = remaining
		}
		if pause > 0 {
			time.Sleep(pause)
		}
	}
	return workflowError(5, "Herdr startup deadline expired while awaiting readiness; last observation retained; no prompt sent")
}
