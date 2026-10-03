package taskrun

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"time"
)

type workflowAgent struct {
	WorkspaceID string `json:"workspace_id"`
	TabID       string `json:"tab_id"`
	PaneID      string `json:"pane_id"`
	TerminalID  string `json:"terminal_id"`
	Name        string `json:"name"`
	Agent       string `json:"agent"`
	Status      string `json:"agent_status"`
	Session     *struct {
		Agent string `json:"agent"`
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"agent_session"`
}
type workflowLimitedOutput struct {
	data     []byte
	overflow bool
}

func (b *workflowLimitedOutput) Write(p []byte) (int, error) {
	n := len(p)
	left := (1 << 20) - len(b.data)
	if n > left {
		b.overflow = true
		p = p[:left]
	}
	b.data = append(b.data, p...)
	return n, nil
}
func workflowCall(d Dependencies, r WorkflowRequest, args ...string) (json.RawMessage, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	c := exec.CommandContext(ctx, r.Herdr.Executable.Path, args...)
	c.WaitDelay = time.Second
	var stdout, stderr workflowLimitedOutput
	c.Stdout = &stdout
	c.Stderr = &stderr
	if e := c.Run(); e != nil {
		return nil, workflowError(5, "Herdr response is unknown (transport failure or timeout); no retry was sent")
	}
	if stdout.overflow || stderr.overflow {
		return nil, workflowError(5, "Herdr response exceeded the transport limit")
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if _, e := canonicalTransportJSON(stdout.data); e != nil {
		return nil, e
	}
	if e := json.Unmarshal(stdout.data, &envelope); e != nil || len(envelope.Result) == 0 || string(envelope.Result) == "null" || len(envelope.Error) > 0 && string(envelope.Error) != "null" {
		return nil, workflowError(5, "Herdr did not return a successful result envelope")
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
	t := s.Result.Transport
	if a.WorkspaceID != t.WorkspaceID || a.TabID != t.TabID || a.PaneID != t.PaneID || a.TerminalID != t.TerminalID || a.Name != workflowAgentName(s.Result.RunID) || a.Agent != "codex" || a.Session == nil || a.Session.Agent != "codex" || a.Session.Kind != "id" || !plain(a.Session.Value, 1, 256) || !first && a.Session.Value != t.AgentSessionID {
		return workflowError(4, "fresh Herdr identity or native Codex session differs from the bound run")
	}
	return nil
}
func workflowAgentGet(d Dependencies, s workflowState, first bool) (workflowAgent, error) {
	var response struct {
		Agent workflowAgent `json:"agent"`
	}
	b, e := workflowCall(d, s.Request, "agent", "get", s.Result.Transport.PaneID)
	if e != nil {
		return response.Agent, e
	}
	if e = json.Unmarshal(b, &response); e != nil {
		return response.Agent, e
	}
	e = workflowIdentity(s, response.Agent, first)
	return response.Agent, e
}
func workflowSettled(a workflowAgent) bool { return a.Status == "idle" || a.Status == "done" }
func workflowObserve(d Dependencies, s *workflowState, a workflowAgent) {
	s.Result.Transport.State = a.Status
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
		return e
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
	if _, e = workflowCall(d, r, argv...); e != nil {
		return e
	}
	s, e = workflowRead(root, id)
	if e != nil {
		return e
	}
	a, e = workflowAgentGet(d, s, true)
	if e != nil {
		return e
	}
	if !workflowSettled(a) {
		return workflowError(5, "new session is not freshly idle; no prompt sent")
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		s.Result.Transport.AgentSessionID = a.Session.Value
		workflowObserve(d, s, a)
		s.Phase = "session_bound"
		return nil
	})
	if e != nil {
		return e
	}
	return workflowPrompt(d, root, id, nil)
}
func workflowPrompt(d Dependencies, root, id string, findings []WorkflowFinding) error {
	s, e := workflowRead(root, id)
	if e != nil {
		return e
	}
	if e = workflowFresh(d, s, false); e != nil {
		return e
	}
	a, e := workflowAgentGet(d, s, false)
	if e != nil {
		return e
	}
	if !workflowSettled(a) {
		return workflowError(4, "same session must be freshly idle before input")
	}
	prompt := workflowInstructions(s, findings)
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if s.Phase != "session_bound" && s.Phase != "correction_reserved" {
			return workflowError(4, "prompt already attempted or not authorized")
		}
		if e := d.writeValue(filepathForRound(*s, "prompt-attempt.json"), map[string]any{"argv": []string{"agent", "prompt", s.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000"}, "round": s.Result.Round.Number}); e != nil {
			return e
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
	_, callErr := workflowCall(d, s.Request, "agent", "prompt", s.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000")
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if callErr != nil {
			s.Phase = "prompt_unknown"
			s.Result.Transport.State = "unknown"
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
