package taskrun

import (
	"os"
	"path/filepath"
	"time"
)

func deliveryStartupPending(s workflowState) bool {
	if !deliveryRun(s.Request) || s.Result.Delivery == nil || s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Round.Number != 0 {
		return false
	}
	if s.Phase != "agent_start_attempted" && s.Phase != "bootstrap_attempted" && s.Phase != "session_bound" {
		return false
	}
	t := s.Result.Transport
	if t.TabID == "" || t.PaneID == "" || t.TerminalID == "" {
		return false
	}
	// Phase alone cannot excuse a preserved possibly-sent Task prompt.
	_, err := os.Lstat(filepathForRound(s, "prompt-attempt.json"))
	return os.IsNotExist(err)
}

func deliveryStartupAction(s workflowState) WorkflowAction {
	if deliveryStartupPending(s) {
		return WorkflowAction{"user", "Complete any native onboarding in the existing tab, then run: ply workflow execute resume " + s.Result.RunID + ". This follows the same start and sends its first Task prompt only when ready; it creates no new tab or agent."}
	}
	return WorkflowAction{"coordinator", "Inspect the preserved startup and last observation in the same tab. Task input may already have been attempted; do not restart or resend it."}
}

// WorkflowResumeDeliveryStart resumes only the existing v2 launch before its
// first Task prompt. It never calls tab create or agent start. The existing
// bootstrap and prompt reservations arbitrate concurrent observers and retain
// unknown sends rather than replaying them.
func WorkflowResumeDeliveryStart(d Dependencies, root, id string, timeout time.Duration) (WorkflowRun, error) {
	if timeout <= 0 {
		return WorkflowRun{}, workflowError(2, "timeout must be positive")
	}
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	deadline := time.Now().Add(timeout)
	var phase string
	err := workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryStartupPending(*s) {
			return workflowError(4, "resume requires an existing delivery startup with no attempted Task prompt; inspect the same run without restarting or resending")
		}
		_, bootstrapErr := os.Lstat(filepath.Join(s.Result.Paths.RunRoot, "startup-bootstrap-attempt.json"))
		if s.Phase == "agent_start_attempted" && !os.IsNotExist(bootstrapErr) {
			return workflowError(4, "a readiness exchange is already reserved or unknown; inspect its preserved state without replaying it")
		}
		if s.Phase == "bootstrap_attempted" && bootstrapErr != nil {
			return workflowError(4, "reserved readiness exchange evidence is missing; inspect without replaying it")
		}
		if err := workflowFresh(d, *s, false); err != nil {
			return err
		}
		cwd, err := d.CWD()
		if err != nil {
			return err
		}
		if cwd != s.Observed.Target.WorktreeLocator && cwd != s.Observed.Epic.WorktreeLocator {
			return workflowError(4, "resume must run from the bound Task or return worktree")
		}
		if err = physical(cwd, false); err != nil {
			return err
		}
		x, err := deliveryTarget(d, *s)
		if err != nil {
			return err
		}
		if x.OID != s.Observed.Target.OID || x.Tree != s.Observed.Target.Tree {
			return workflowError(4, "Task start target changed before its first prompt")
		}
		phase = s.Phase
		return nil
	})
	if err != nil {
		return deliveryReadback(d, root, id, err)
	}
	if phase != "session_bound" {
		err = workflowAwaitReadiness(d, root, id, deadline)
	}
	if err == nil {
		err = workflowPromptUntil(d, root, id, nil, deadline)
	}
	if err != nil {
		cause := err
		if saveErr := workflowUpdate(d, root, id, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"delivery_start_resume_stopped", cause.Error()})
			// A competing successful first prompt or acceptance owns subsequent
			// progress; a late observer must not replace its next action.
			if s.Acceptance == nil && s.Phase != "following" {
				s.Result.NextAction = deliveryStartupAction(*s)
			}
			return nil
		}); saveErr != nil {
			return deliveryReadback(d, root, id, saveErr)
		}
	}
	return deliveryReadback(d, root, id, err)
}
