package taskrun

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The reviewer is a separate session started through the same Herdr launcher
// as the owner: same provider, model, permission mode and Task cwd, a new tab
// and agent name. Each effect is reserved in a write-once round file before it
// is sent and never replayed; an unknown outcome stays visible as a stop.

// Shared by newly generated briefs and launch prompts. Historical round files
// remain immutable; this text grants no runtime authority or owner callback role.
const deliveryReviewerTestMandate = "You may run relevant ordinary tests within the agreed read-only source, cache and private temporary directory boundaries, subject to actual runtime review of each action. Keep the exact candidate unchanged and Go modules readonly with -mod=readonly. Do not edit the Task worktree or commit. For the assigned product Task/Run and its delivery, `workflow execute verify` and `workflow execute integrate` are owner-only callbacks; the reviewer cannot invoke them for that product Run. Relevant ordinary tests in separate isolated fixtures may exercise synthetic owner callbacks, including verify/integrate. These fixtures do not grant owner authority over the assigned product Task/Run or its delivery. This test mandate does not override an actual refusal or unknown prior effect; preserve and report actual outcomes and evidence gaps."

type deliveryReviewLaunch struct {
	Kind           string   `json:"kind"`
	SchemaVersion  int      `json:"schema_version"`
	RunID          string   `json:"run_id"`
	Round          int      `json:"round"`
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	PermissionMode string   `json:"permission_mode"`
	WorkspaceID    string   `json:"workspace_id"`
	AgentName      string   `json:"agent_name"`
	CWD            string   `json:"cwd"`
	TabArgv        []string `json:"tab_argv"`
	ReservedAtUTC  string   `json:"reserved_at_utc"`
}

type deliveryReviewPane struct {
	WorkspaceID string `json:"workspace_id"`
	TabID       string `json:"tab_id"`
	PaneID      string `json:"pane_id"`
	TerminalID  string `json:"terminal_id"`
}

type deliveryReviewAgent struct {
	Kind          string             `json:"kind"`
	SchemaVersion int                `json:"schema_version"`
	RunID         string             `json:"run_id"`
	Round         int                `json:"round"`
	Pane          deliveryReviewPane `json:"pane"`
	AgentName     string             `json:"agent_name"`
	Provider      string             `json:"provider"`
	SessionID     string             `json:"session_id"`
	OwnerSession  string             `json:"owner_session"`
	BoundAtUTC    string             `json:"bound_at_utc"`
}

type deliveryReviewPromptResult struct {
	Outcome string `json:"outcome"`
	Error   string `json:"error"`
}

// deliveryReviewStart is the observed reviewer start of one round.
type deliveryReviewStart struct {
	State  string
	Pane   *deliveryReviewPane
	Agent  *deliveryReviewAgent
	Detail string
}

func deliveryReviewAgentName(id string, round int) string {
	return fmt.Sprintf("ply-review-%s-%03d", id[4:20], round)
}

func deliveryReviewPlan(d Dependencies, s workflowState, round int, workspaceID string) (deliveryReviewLaunch, error) {
	runtime, e := workflowEffectiveRuntime(s)
	if e != nil {
		return deliveryReviewLaunch{}, e
	}
	if workspaceID == "" {
		workspaceID = s.Result.Transport.WorkspaceID
	}
	if !plain(workspaceID, 1, 128) {
		return deliveryReviewLaunch{}, workflowError(2, "--herdr-workspace is required when the owner's Herdr workspace is unknown")
	}
	argv := workflowTabCreateArgv(s, workspaceID)
	for i := range argv {
		if argv[i] == "--label" && i+1 < len(argv) {
			argv[i+1] = workflowTaskTabLabel(s, fmt.Sprintf("review %d", round))
		}
	}
	argv = append(argv, "--env", "GOFLAGS="+deliveryReviewGoFlags())
	return deliveryReviewLaunch{Kind: "ply.workflow.review-launch", SchemaVersion: 1, RunID: s.Result.RunID, Round: round, Provider: runtime.Provider, Model: runtime.Model, PermissionMode: runtime.PermissionBinding.ProfileID, WorkspaceID: workspaceID, AgentName: deliveryReviewAgentName(s.Result.RunID, round), CWD: s.Observed.Target.WorktreeLocator, TabArgv: argv, ReservedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}, nil
}

func deliveryReviewRead(path string, v any) (bool, error) {
	raw, e := readFile(path, 1<<20, true)
	if os.IsNotExist(e) {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	return true, json.Unmarshal(raw, v)
}

func deliveryReviewStartState(dir string) deliveryReviewStart {
	var launch deliveryReviewLaunch
	if ok, e := deliveryReviewRead(filepath.Join(dir, "launch.json"), &launch); e != nil || !ok {
		if e != nil {
			return deliveryReviewStart{State: "unknown", Detail: e.Error()}
		}
		return deliveryReviewStart{State: "not_started"}
	}
	var pane deliveryReviewPane
	if ok, e := deliveryReviewRead(filepath.Join(dir, "tab.json"), &pane); e != nil || !ok {
		return deliveryReviewStart{State: "start_unknown", Detail: "the reviewer tab outcome is unknown; inspect Herdr, it is not replayed"}
	}
	out := deliveryReviewStart{State: "start_unknown", Pane: &pane, Detail: "the reviewer tab was created in pane " + pane.PaneID + " but its session is not bound; inspect that tab, it is not replayed"}
	var agent deliveryReviewAgent
	if ok, e := deliveryReviewRead(filepath.Join(dir, "agent.json"), &agent); e != nil || !ok {
		return out
	}
	out.Agent, out.State, out.Detail = &agent, "prompt_unknown", "the review prompt outcome in pane "+pane.PaneID+" is unknown; inspect that tab, it is not resent"
	var prompt deliveryReviewPromptResult
	if ok, e := deliveryReviewRead(filepath.Join(dir, "prompt-result.json"), &prompt); e == nil && ok && prompt.Outcome == "sent" {
		out.State, out.Detail = "started", ""
	}
	return out
}

func deliveryReviewCall(d Dependencies, s workflowState, deadline time.Time, args ...string) (json.RawMessage, error) {
	return workflowCallUntil(d, s.Request, deadline, args...)
}

func deliveryReviewAgentGet(d Dependencies, s workflowState, pane deliveryReviewPane, name, provider string, deadline time.Time) (workflowAgent, error) {
	var response struct {
		Agent workflowAgent `json:"agent"`
	}
	b, e := deliveryReviewCall(d, s, deadline, "agent", "get", pane.PaneID)
	if e != nil {
		return response.Agent, e
	}
	if e = json.Unmarshal(b, &response); e != nil {
		return response.Agent, workflowError(5, "Herdr agent get for the reviewer returned an invalid result")
	}
	a := response.Agent
	if a.WorkspaceID != pane.WorkspaceID || a.TabID != pane.TabID || a.PaneID != pane.PaneID || a.TerminalID != pane.TerminalID || a.Name != name || a.Agent != provider {
		return a, workflowError(4, "fresh Herdr reviewer location, name or provider differs from the started reviewer")
	}
	if a.Session != nil && (a.Session.Agent != provider || a.Session.Kind != "id" || !plain(a.Session.Value, 1, 256)) {
		return a, workflowError(4, "fresh reviewer session is not a native "+provider+" session")
	}
	return a, nil
}

// deliveryLaunchReviewer performs the reserved start of one round's reviewer.
func deliveryLaunchReviewer(d Dependencies, s workflowState, r deliveryReviewRecord, l deliveryReviewLaunch) error {
	limit := 90 * time.Second
	if d.HerdrTimeout > 0 {
		limit = d.HerdrTimeout
	}
	deadline := time.Now().Add(limit)
	b, e := deliveryReviewCall(d, s, deadline, l.TabArgv...)
	if e != nil {
		return workflowError(5, "reviewer tab create outcome is unknown; it is not replayed: "+e.Error())
	}
	var tab struct {
		RootPane workflowAgent `json:"root_pane"`
	}
	if e = json.Unmarshal(b, &tab); e != nil {
		return workflowError(5, "reviewer tab create returned an invalid result; created identity unknown")
	}
	pane := deliveryReviewPane{tab.RootPane.WorkspaceID, tab.RootPane.TabID, tab.RootPane.PaneID, tab.RootPane.TerminalID}
	if pane.WorkspaceID != l.WorkspaceID || !plain(pane.TabID, 1, 128) || !plain(pane.PaneID, 1, 128) || !plain(pane.TerminalID, 1, 128) || pane.PaneID == s.Result.Transport.PaneID {
		return workflowError(5, "created reviewer tab identity is unknown")
	}
	if _, e = workflowKeep(d, filepath.Join(r.Dir, "tab.json"), pane); e != nil {
		return e
	}
	if l.Provider == "claude" && s.ClaudeTrust != nil && s.ClaudeTrust.Reason == "" {
		if e = deliveryReviewClaudeEnvironment(d, s, r.Dir, pane.PaneID); e != nil {
			return e
		}
	}
	argv := workflowStartArgv(workflowEffectiveRequest(s), l.CWD, pane.PaneID)
	argv[2] = l.AgentName
	if _, e = workflowKeep(d, filepath.Join(r.Dir, "start-attempt.json"), map[string]any{"argv": argv}); e != nil {
		return e
	}
	// A lost start reply is followed only by observation of this one pane.
	_, startErr := deliveryReviewCall(d, s, deadline, argv...)
	if facts := deliveryReviewFailureFacts(startErr); facts != nil {
		if _, e = workflowKeep(d, filepath.Join(r.Dir, "start-failure.json"), facts); e != nil {
			return e
		}
	}
	owner := deliveryOwnerSession(s)
	bootstrapped := false
	var last error
	var firstObservationFailure *deliveryReviewFailure
	var lastSession, lastStatus string
	for {
		a, e := deliveryReviewAgentGet(d, s, pane, l.AgentName, l.Provider, deadline)
		last = e
		observation := deliveryReviewObservation{Outcome: "unavailable", TaskPromptStage: "not_attempted", BootstrapPromptStage: "not_attempted", Failure: deliveryReviewFailureFacts(e)}
		if observation.Failure != nil && firstObservationFailure == nil {
			firstObservationFailure = observation.Failure
		}
		observation.FirstFailure = firstObservationFailure
		if bootstrapped {
			observation.BootstrapPromptStage = "attempted"
		}
		if e == nil {
			observation.Outcome = "observed"
			if a.Session != nil {
				observation.SessionID = a.Session.Value
			}
			switch a.Status {
			case "idle", "working", "starting", "stopped":
				observation.AgentStatus = a.Status
			}
		}
		if observation.SessionID != "" {
			lastSession, lastStatus = observation.SessionID, observation.AgentStatus
		}
		observation.SessionID, observation.AgentStatus = lastSession, lastStatus
		if e == nil && a.Session != nil {
			if a.Session.Value == owner {
				if _, e = workflowKeep(d, filepath.Join(r.Dir, "session-observation.json"), observation); e != nil {
					return e
				}
				return workflowError(4, "the started reviewer reports the owner's own session; an independent review needs a distinct session")
			}
			if workflowSettled(a) {
				if _, e = workflowKeep(d, filepath.Join(r.Dir, "session-observation.json"), observation); e != nil {
					return e
				}
				bound := deliveryReviewAgent{"ply.workflow.review-agent", 1, s.Result.RunID, r.Request.Round, pane, l.AgentName, l.Provider, a.Session.Value, owner, d.Now().UTC().Format(time.RFC3339Nano)}
				if e = withStore(s.Request.WorkspaceRoot, func() error {
					if e := parallelSessionAvailable(d, s, l.Provider, a.Session.Value, false, r.Dir); e != nil {
						return e
					}
					_, e := workflowKeep(d, filepath.Join(r.Dir, "agent.json"), bound)
					return e
				}); e != nil {
					return e
				}
				break
			}
		} else if e == nil && !bootstrapped && a.InteractiveReady != nil && *a.InteractiveReady && workflowSettled(a) {
			// The provider defers its session until a first, tool-free turn.
			bootstrap := []string{"agent", "prompt", pane.PaneID, workflowBootstrapPrompt, "--wait", "--until", "working", "--timeout", "10000"}
			if _, e = workflowKeep(d, filepath.Join(r.Dir, "bootstrap-attempt.json"), map[string]any{"argv": bootstrap}); e != nil {
				return e
			}
			bootstrapped = true
			observation.BootstrapPromptStage = "attempted"
			_, _ = deliveryReviewCall(d, s, deadline, bootstrap...)
		}
		if !time.Now().Before(deadline) {
			if _, e = workflowKeep(d, filepath.Join(r.Dir, "session-observation.json"), observation); e != nil {
				return e
			}
			detail := "reviewer session was not observed before the startup deadline; inspect pane " + pane.PaneID + ", the start is not replayed"
			if startErr != nil {
				detail += " (start reply: " + startErr.Error() + ")"
			}
			if last != nil {
				detail += " (last observation: " + last.Error() + ")"
			}
			return workflowError(5, detail)
		}
		time.Sleep(250 * time.Millisecond)
	}
	prompt := fmt.Sprintf("You are the independent reviewer for Ply review round %d of Run %s, started by %s. Read and follow %s. Review only; %s", r.Request.Round, s.Result.RunID, r.Request.Coordinator, r.Request.Brief.Locator, deliveryReviewerTestMandate)
	send := []string{"agent", "prompt", pane.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000"}
	if _, e = workflowKeep(d, filepath.Join(r.Dir, "prompt-attempt.json"), map[string]any{"argv": send}); e != nil {
		return e
	}
	_, sendErr := deliveryReviewCall(d, s, deadline, send...)
	result := deliveryReviewPromptResult{Outcome: "sent"}
	if sendErr != nil {
		result = deliveryReviewPromptResult{Outcome: "unknown", Error: sendErr.Error()}
	}
	if _, e = workflowKeep(d, filepath.Join(r.Dir, "prompt-result.json"), result); e != nil {
		return e
	}
	if sendErr != nil {
		return workflowError(5, "review prompt outcome is unknown; inspect pane "+pane.PaneID+", it is not resent")
	}
	return nil
}

func deliveryReviewClaudeEnvironment(d Dependencies, s workflowState, dir, pane string) error {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	token := hex.EncodeToString(nonce)
	ack := filepath.Join(dir, "claude-environment-"+token+".ack")
	command := workflowClaudeEnvironmentCommand(s.ClaudeTrust, token, ack)
	if _, err := workflowKeep(d, filepath.Join(dir, "claude-environment-attempt.json"), map[string]any{"pane": pane, "command": command, "ack": ack}); err != nil {
		return err
	}
	deadline := time.Now().Add(5 * time.Second)
	_, _ = deliveryReviewCall(d, s, deadline, "pane", "run", pane, command)
	for time.Now().Before(deadline) {
		if b, err := readFile(ack, 128, true); err == nil && string(b) == token {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return workflowError(5, "reviewer Claude shell environment is unconfirmed; no reviewer was started. Inspect pane "+pane+"; its shell input is not replayed")
}

// deliveryReviewGoFlags is the reviewer's process-local GOFLAGS: the ambient
// flags with readonly Go modules last, so an ambient -mod=mod cannot let the
// reviewer's tests rewrite the candidate's module metadata.
func deliveryReviewGoFlags() string {
	return strings.TrimSpace(os.Getenv("GOFLAGS") + " -mod=readonly")
}
