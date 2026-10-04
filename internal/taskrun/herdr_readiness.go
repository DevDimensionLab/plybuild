package taskrun

import (
	"path/filepath"
	"time"
)

const workflowBootstrapPrompt = "Transport readiness check. Reply exactly PLY_READY. Do not use tools, read files, or change anything. Wait for the next message."

// Codex can defer SessionStart until its first turn. This exchange carries no
// Task authority. Only the owner of this launch may attempt it, and a lost reply
// is followed by observation of the same attempt, never by repeated input.
func workflowBootstrap(d Dependencies, root, id string, deadline time.Time) error {
	s, e := workflowRead(root, id)
	if e != nil {
		return e
	}
	argv := []string{"agent", "prompt", s.Result.Transport.PaneID, workflowBootstrapPrompt, "--wait", "--until", "working", "--timeout", "10000"}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if s.Phase != "agent_start_attempted" {
			return workflowError(4, "startup readiness exchange already attempted or not authorized")
		}
		if !time.Now().Before(deadline) {
			return workflowError(5, "Herdr startup deadline expired before readiness exchange; no Task prompt sent")
		}
		if e := d.writeValue(filepath.Join(s.Result.Paths.RunRoot, "startup-bootstrap-attempt.json"), map[string]any{"argv": argv}); e != nil {
			return e
		}
		s.Phase = "bootstrap_attempted"
		return nil
	})
	if e != nil {
		return e
	}
	if e = d.fault("workflow_before_bootstrap_send"); e != nil {
		return e
	}
	if _, callErr := workflowCallUntil(d, s.Request, deadline, argv...); callErr != nil {
		return workflowUpdate(d, root, id, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"herdr_bootstrap_response", callErr.Error()})
			return nil
		})
	}
	return nil
}

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
				return workflowError(4, "native Codex session changed during startup; no Task prompt sent")
			}
			observedSession = a.Session.Value
		}
		ready := a.Session != nil && workflowSettled(a)
		bootstrap := a.Session == nil && a.InteractiveReady != nil && *a.InteractiveReady && workflowSettled(a) && s.Phase == "agent_start_attempted"
		if e = workflowUpdate(d, root, id, func(s *workflowState) error {
			if s.Phase != "agent_start_attempted" && s.Phase != "bootstrap_attempted" {
				return workflowError(4, "launch attempt is no longer awaiting readiness")
			}
			workflowObserve(d, s, a)
			if ready {
				if e := workflowIdentity(*s, a, true); e != nil {
					return e
				}
				if !time.Now().Before(deadline) {
					return workflowError(5, "Herdr startup deadline expired before session binding; no Task prompt sent")
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
		if bootstrap {
			if e = workflowBootstrap(d, root, id, deadline); e != nil {
				return e
			}
			continue
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
	return workflowError(5, "Herdr startup deadline expired while awaiting readiness; last observation retained; no Task prompt sent")
}
