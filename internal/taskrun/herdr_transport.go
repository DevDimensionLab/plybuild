package taskrun

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"time"
)

type workflowAgent struct {
	WorkspaceID   string `json:"workspace_id"`
	TabID         string `json:"tab_id"`
	PaneID        string `json:"pane_id"`
	TerminalID    string `json:"terminal_id"`
	Name          string `json:"name"`
	Agent         string `json:"agent"`
	Status        string `json:"agent_status"`
	LaunchPending *bool  `json:"launch_pending"`
	Session       *struct {
		Agent string `json:"agent"`
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"agent_session"`
}
type workflowLimitedOutput struct {
	data     []byte
	overflow bool
	bytes    int64
}

func (b *workflowLimitedOutput) Write(p []byte) (int, error) {
	n := len(p)
	b.bytes += int64(n)
	left := (1 << 20) - len(b.data)
	if n > left {
		b.overflow = true
		p = p[:left]
	}
	b.data = append(b.data, p...)
	return n, nil
}
func workflowCall(d Dependencies, r WorkflowRequest, args ...string) (json.RawMessage, error) {
	return workflowCallUntil(d, r, time.Time{}, args...)
}
func workflowCallUntil(d Dependencies, r WorkflowRequest, deadline time.Time, args ...string) (json.RawMessage, error) {
	if e := verifyExecutable(r.Herdr.Executable); e != nil {
		return nil, e
	}
	limit := 15 * time.Second
	if len(args) > 1 && args[0] == "agent" && args[1] == "start" {
		limit = 35 * time.Second
	}
	if d.HerdrTimeout > 0 {
		limit = d.HerdrTimeout
	}
	callDeadline := time.Now().Add(limit)
	if !deadline.IsZero() && deadline.Before(callDeadline) {
		callDeadline = deadline
	}
	ctx, cancel := context.WithDeadline(context.Background(), callDeadline)
	defer cancel()
	c := exec.CommandContext(ctx, r.Herdr.Executable.Path, args...)
	c.WaitDelay = time.Second
	var stdout, stderr workflowLimitedOutput
	c.Stdout = &stdout
	c.Stderr = &stderr
	if e := c.Run(); e != nil {
		return nil, workflowCallFailure(args, ctx.Err(), e, stdout, stderr, "")
	}
	if stdout.overflow || stderr.overflow {
		return nil, workflowCallFailure(args, nil, nil, stdout, stderr, "output_limit")
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if _, e := canonicalTransportJSON(stdout.data); e != nil {
		return nil, workflowCallFailure(args, nil, nil, stdout, stderr, "invalid_response")
	}
	if e := json.Unmarshal(stdout.data, &envelope); e != nil || len(envelope.Result) == 0 || string(envelope.Result) == "null" || len(envelope.Error) > 0 && string(envelope.Error) != "null" {
		return nil, workflowCallFailure(args, nil, nil, stdout, stderr, "unsuccessful_envelope")
	}
	return envelope.Result, nil
}
func canonicalTransportJSON(b []byte) (any, error) {
	var v any
	if !json.Valid(b) {
		return nil, workflowError(5, "invalid Herdr response")
	}
	e := json.Unmarshal(b, &v)
	return v, e
}
func workflowAgentName(id string) string { return "ply-" + id[4:32] }
func workflowIdentity(s workflowState, a workflowAgent, first bool) error {
	if e := workflowLocation(s, a); e != nil {
		return e
	}
	t := s.Result.Transport
	if a.Session == nil || a.Session.Agent != "codex" || a.Session.Kind != "id" || !plain(a.Session.Value, 1, 256) || (!first || t.AgentSessionID != "") && a.Session.Value != t.AgentSessionID {
		return workflowError(4, "fresh Herdr identity or native Codex session differs from the bound run")
	}
	return nil
}
func workflowLocation(s workflowState, a workflowAgent) error {
	t := s.Result.Transport
	if !plain(a.WorkspaceID, 1, 128) || !plain(a.TabID, 1, 128) || !plain(a.PaneID, 1, 128) || !plain(a.TerminalID, 1, 128) || a.WorkspaceID != t.WorkspaceID || a.TabID != t.TabID || a.PaneID != t.PaneID || a.TerminalID != t.TerminalID || a.Name != workflowAgentName(s.Result.RunID) || a.Agent != "codex" {
		return workflowError(4, "fresh Herdr location, name or agent differs from the reserved attempt")
	}
	return nil
}
func workflowAgentGet(d Dependencies, s workflowState, first bool) (workflowAgent, error) {
	return workflowAgentGetUntil(d, s, first, time.Time{})
}
func workflowAgentGetUntil(d Dependencies, s workflowState, first bool, deadline time.Time) (workflowAgent, error) {
	a, e := workflowAgentRead(d, s, deadline)
	if e == nil {
		e = workflowIdentity(s, a, first)
	}
	return a, e
}
func workflowAgentRead(d Dependencies, s workflowState, deadline time.Time) (workflowAgent, error) {
	var response struct {
		Agent workflowAgent `json:"agent"`
	}
	b, e := workflowCallUntil(d, s.Request, deadline, "agent", "get", s.Result.Transport.PaneID)
	if e != nil {
		return response.Agent, e
	}
	if e = json.Unmarshal(b, &response); e != nil {
		return response.Agent, workflowError(5, "Herdr agent get failed: class=invalid_result exit=0; identity unknown; no prompt sent")
	}
	return response.Agent, nil
}
func workflowSettled(a workflowAgent) bool {
	return (a.Status == "idle" || a.Status == "done") && (a.LaunchPending == nil || !*a.LaunchPending)
}
func workflowObserve(d Dependencies, s *workflowState, a workflowAgent) {
	s.Result.Transport.State = "unknown"
	switch a.Status {
	case "idle", "done", "working", "blocked":
		s.Result.Transport.State = a.Status
	}
	s.Result.Transport.LaunchPending = a.LaunchPending
	s.Result.Transport.Observation = "fresh"
	s.Result.Transport.ObservedAt = d.Now().UTC().Format(time.RFC3339Nano)
}
func workflowLaunch(d Dependencies, r WorkflowRequest) error {
	root, id := r.WorkspaceRoot, workflowID(r)
	if e := d.fault("workflow_before_tab_send"); e != nil {
		return e
	}
	s, e := workflowRead(root, id)
	if e != nil {
		return e
	}
	var tab struct {
		RootPane workflowAgent `json:"root_pane"`
	}
	b, e := workflowCall(d, r, "tab", "create", "--workspace", r.Herdr.WorkspaceID, "--cwd", s.Observed.Target.WorktreeLocator, "--label", "run "+r.Herdr.TabLabel, "--env", "PATH="+os.Getenv("PATH"), "--no-focus")
	if e != nil {
		return e
	}
	if e = json.Unmarshal(b, &tab); e != nil {
		return workflowError(5, "Herdr tab create failed: class=invalid_result exit=0; created identity unknown")
	}
	a := tab.RootPane
	if a.WorkspaceID != r.Herdr.WorkspaceID || !plain(a.PaneID, 1, 128) || !plain(a.TabID, 1, 128) || !plain(a.TerminalID, 1, 128) {
		return workflowError(5, "created tab identity is unknown")
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		s.Result.Transport.TabID = a.TabID
		s.Result.Transport.PaneID = a.PaneID
		s.Result.Transport.TerminalID = a.TerminalID
		s.Phase = "agent_start_attempted"
		return nil
	})
	if e != nil {
		return e
	}
	if e = d.fault("workflow_before_agent_send"); e != nil {
		return e
	}
	argv := []string{"agent", "start", workflowAgentName(id), "--kind", "codex", "--pane", a.PaneID, "--timeout", "30000", "--", "-C", s.Observed.Target.WorktreeLocator, "--model", r.Runtime.Model, "-a", "on-request", "-c", `approvals_reviewer="auto_review"`, "-c", "default_permissions=" + workflowJSON(r.Runtime.PermissionBinding.ProfileID)}
	if r.Runtime.ConfigProfile != nil {
		argv = append(argv, "--profile", *r.Runtime.ConfigProfile)
	}
	limit := 35 * time.Second
	if d.HerdrTimeout > 0 {
		limit = d.HerdrTimeout
	}
	deadline := time.Now().Add(limit)
	if _, startErr := workflowCallUntil(d, r, deadline, argv...); startErr != nil {
		// A failed reply does not undo the single reserved start. Only fresh
		// read-only observations of this exact attempt may establish readiness.
		if e = workflowUpdate(d, root, id, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"herdr_start_response", startErr.Error()})
			return nil
		}); e != nil {
			return e
		}
	}
	if e = workflowAwaitReadiness(d, root, id, deadline); e != nil {
		return e
	}
	return workflowPromptUntil(d, root, id, nil, deadline)
}
func workflowPrompt(d Dependencies, root, id string, findings []WorkflowFinding) error {
	return workflowPromptUntil(d, root, id, findings, time.Time{})
}
func workflowPromptUntil(d Dependencies, root, id string, findings []WorkflowFinding, deadline time.Time) error {
	s, e := workflowRead(root, id)
	if e != nil {
		return e
	}
	if e = workflowFresh(d, s, false); e != nil {
		return e
	}
	a, e := workflowAgentGetUntil(d, s, false, deadline)
	if e != nil {
		return e
	}
	if !workflowSettled(a) {
		return workflowError(4, "same session must be freshly idle before input")
	}
	prompt := workflowInstructions(s, findings)
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deadline.IsZero() && !time.Now().Before(deadline) {
			return workflowError(5, "Herdr startup deadline expired before prompt reservation; no prompt sent")
		}
		if s.Phase != "session_bound" && s.Phase != "correction_reserved" {
			return workflowError(4, "prompt already attempted or not authorized")
		}
		if e := d.writeValue(filepathForRound(*s, "prompt-attempt.json"), map[string]any{"argv": []string{"agent", "prompt", s.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000"}, "round": s.Result.Round.Number}); e != nil {
			return e
		}
		if s.Result.Round.Number == 0 && s.Acceptance == nil {
			s.Result.NextAction = WorkflowAction{"recipient", "Accept the bound runtime before target writes, then report this round."}
		}
		s.Phase = "prompt_attempted"
		return nil
	})
	if e != nil {
		return e
	}
	if e = d.fault("workflow_before_prompt_send"); e != nil {
		return e
	}
	_, callErr := workflowCallUntil(d, s.Request, deadline, "agent", "prompt", s.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000")
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if callErr != nil {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"herdr_prompt_response", callErr.Error()})
		}
		if s.Phase != "prompt_attempted" {
			return nil
		}
		if callErr != nil {
			s.Phase = "prompt_unknown"
			s.Result.Transport.State = "unknown"
			if s.Result.Round.ReportSHA256 == nil && s.Result.FinalReturn.State == "pending" {
				s.Result.NextAction = WorkflowAction{"coordinator", "Inspect the preserved prompt attempt in the same session; do not resend input."}
			}
		} else {
			s.Phase = "following"
			s.Result.Transport.State = "working"
		}
		return nil
	})
	if e != nil {
		return e
	}
	return callErr
}

func WorkflowFollow(d Dependencies, root, id string, timeout time.Duration) (WorkflowRun, error) {
	if timeout <= 0 {
		return WorkflowRun{}, workflowError(2, "timeout must be positive")
	}
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	deadline := time.Now().Add(timeout)
	quiet := 0
	lastRound := -1
	lastReport := ""
	var lastQuiet time.Time
	for {
		ready := false
		var observationErr error
		e := workflowUpdate(d, root, id, func(s *workflowState) error {
			if e := workflowFresh(d, *s, true); e != nil {
				return e
			}
			observing := d
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return workflowError(5, "observation timeout")
			}
			if observing.HerdrTimeout <= 0 || observing.HerdrTimeout > remaining {
				observing.HerdrTimeout = remaining
			}
			a, e := workflowAgentGet(observing, *s, false)
			if e != nil {
				observationErr = e
				s.Result.Transport.State = "unknown"
				s.Result.Transport.Observation = "fresh"
				s.Result.Transport.ObservedAt = d.Now().UTC().Format(time.RFC3339Nano)
				if s.Result.Round.State != "sealed" {
					s.Result.Round.State = "unknown"
				}
				s.Result.NextAction = WorkflowAction{"coordinator", "Inspect this run's uncertain transport identity; no input was sent."}
				return nil
			}
			workflowObserve(d, s, a)
			if s.Result.Round.State == "sealed" {
				ready = true
				return nil
			}
			r := s.Result.Round
			rh := ""
			if r.ReportSHA256 != nil {
				rh = *r.ReportSHA256
			}
			if r.Number != lastRound || rh != lastReport {
				quiet = 0
			}
			lastRound, lastReport = r.Number, rh
			now := time.Now()
			if now.Sub(lastQuiet) > time.Second {
				quiet = 0
			}
			if rh != "" && workflowSettled(a) {
				quiet++
				lastQuiet = now
			} else {
				quiet = 0
			}
			if quiet >= 2 {
				s.Result.Round.State = "ready_for_review"
				s.Result.NextAction = WorkflowAction{"coordinator", "Review the immutable report against the frozen Task Spec."}
				ready = true
			}
			return nil
		})
		if e == nil {
			e = observationErr
		}
		if e != nil {
			o, _ := WorkflowShow(d, root, id)
			return o, workflowError(5, "Transport observation is unknown: "+e.Error())
		}
		if ready {
			s, e := workflowRead(root, id)
			if e != nil {
				return WorkflowRun{}, e
			}
			// Tab labels express transport status only. They never grant review.
			_, e = workflowCall(d, s.Request, "tab", "rename", s.Result.Transport.TabID, "finished "+s.Request.Herdr.TabLabel)
			o, se := WorkflowShow(d, root, id)
			o.Transport.Observation = "fresh"
			if e == nil {
				e = se
			}
			return o, e
		}
		if !time.Now().Before(deadline) {
			o, _ := WorkflowShow(d, root, id)
			o.Transport.Observation = "fresh"
			return o, workflowError(5, "Waiting for this round's report and two fresh settled observations; the agent was not stopped")
		}
		pause := 250 * time.Millisecond
		if remaining := time.Until(deadline); remaining < pause {
			pause = remaining
		}
		time.Sleep(pause)
	}
}
