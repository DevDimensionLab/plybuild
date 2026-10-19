package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
)

type taskWriterState uint8

const (
	taskWritersUnknown taskWriterState = iota
	taskWritersPresent
	taskWritersAbsent
)

// taskWriters composes Herdr identities and local writer evidence through the
// existing production/fixture adapters. Callers choose their action and retain
// its effects and wait policy; they never infer absence from failed inspection.
type taskWriters struct {
	d Dependencies
	s workflowState
}

type taskWriterEvidence struct {
	state    taskWriterState
	absence  HumanOwnerRecoveryObservation
	liveness *TaskWriterLiveness
	reason   string
	err      error
	// herdrPresent distinguishes interactive writer handoff from local
	// writers and unknown inspection. Unknown never proves owner presence.
	herdrPresent bool
}

func (o taskWriterEvidence) unknown(reason string, err error) taskWriterEvidence {
	o.state, o.reason = taskWritersUnknown, reason
	o.err = workflowError(4, "owner_unknown: writer liveness is unknown; inspection unavailable: "+err.Error()+"; restore read access before continuing the same action; no writer absence was established")
	return o
}

func (o taskWriterEvidence) present(reason, detail string) taskWriterEvidence {
	o.state, o.reason, o.err = taskWritersPresent, reason, workflowError(4, detail)
	o.herdrPresent = true
	return o
}

// Historical human recovery deliberately observes only Herdr. Its preserved
// process_liveness stays unknown and cannot authorize a new writer or cleanup.
func observeAbsentDeliveryOwner(d Dependencies, s workflowState) (HumanOwnerRecoveryObservation, error) {
	o := (taskWriters{d, s}).herdr("")
	return o.absence, o.err
}

func (w taskWriters) resume(invoker string) (taskWriterEvidence, string) {
	o := w.herdr(invoker)
	if o.state != taskWritersAbsent {
		return o, ""
	}
	transcript, err := deliveryProviderSession(w.s)
	if err != nil {
		return o.unknown("provider_session_missing", err), ""
	}
	return w.local(o, transcript, invokerShellPID(w.d, w.s, invoker)), transcript
}

func (w taskWriters) retirement() taskWriterEvidence {
	o := w.herdr("")
	if o.state != taskWritersAbsent || w.s.OwnerResume == nil {
		return o
	}
	record, _, err := workflowOwnerResumeRecord(w.s)
	if err != nil {
		return o.unknown("inspection_unavailable", err)
	}
	return w.local(o, record.ProviderSession, 0)
}

func (w taskWriters) cleanup() taskWriterEvidence {
	o := taskWriterEvidence{state: taskWritersAbsent}
	if !headlessRun(w.s.Request) {
		o = w.herdr("")
		if o.state != taskWritersAbsent {
			return o
		}
	}
	transcript, _ := deliveryProviderSession(w.s)
	return w.local(o, transcript, 0)
}

func (w taskWriters) local(o taskWriterEvidence, transcript string, shell int) taskWriterEvidence {
	observed, err := observeWriterLiveness(w.d, w.s, transcript, shell)
	o.liveness = &observed
	if err != nil {
		return o.unknown("liveness_unknown", err)
	}
	if len(observed.RemainingPIDs) > 0 {
		o.state, o.reason = taskWritersPresent, "writer_process_alive"
		o.err = workflowError(4, fmt.Sprintf("owner_active: local Task/session writer processes remain %v; stop those processes (possibly outside Herdr), then continue this same action from the return worktree", observed.RemainingPIDs))
		return o
	}
	o.absence.ProcessLiveness = "no_task_or_session_process_observed"
	return o
}

// herdr also accepts the identified invoking pane
// when Herdr lists it as a plain terminal without any agent or session. An
// agent in that pane, or any other Task pane, still counts as a writer.
func (w taskWriters) herdr(invoker string) taskWriterEvidence {
	d, s := w.d, w.s
	out := taskWriterEvidence{absence: HumanOwnerRecoveryObservation{Transport: s.Result.Transport, ProcessLiveness: "unknown"}}
	raw, err := workflowCall(d, s.Request, "pane", "get", s.Result.Transport.PaneID)
	var call *HerdrCallError
	if !errors.As(err, &call) || !call.ReadOnly || call.Phase != "pane get" || call.ProviderCode != "pane_not_found" {
		if err == nil {
			var response struct {
				Pane workflowAgent `json:"pane"`
			}
			if json.Unmarshal(raw, &response) != nil || !w.boundIdentity(response.Pane) {
				return out.unknown("inspection_unavailable", errors.New("incomplete or different bound pane identity"))
			}
			return out.present("terminal_present", "owner_active: the bound pane still exists; use ordinary release for an existing terminal")
		}
		return out.unknown("inspection_unavailable", err)
	}
	raw, err = workflowCall(d, s.Request, "agent", "get", s.Result.Transport.PaneID)
	if err == nil {
		var response struct {
			Agent workflowAgent `json:"agent"`
		}
		if json.Unmarshal(raw, &response) != nil || !w.boundIdentity(response.Agent) {
			return out.unknown("inspection_unavailable", errors.New("incomplete or different bound agent identity"))
		}
		return out.present("owner_active", "owner_active: the bound pane still has an active or pending managed agent")
	}
	if !errors.As(err, &call) || !call.ReadOnly || call.Phase != "agent get" || call.ProviderCode != "agent_not_found" {
		return out.unknown("inspection_unavailable", err)
	}
	for _, kind := range []string{"pane", "agent"} {
		raw, err := workflowCall(d, s.Request, kind, "list")
		if err != nil {
			return out.unknown("inspection_unavailable", err)
		}
		var inventory struct {
			Panes  *[]workflowAgent `json:"panes"`
			Agents *[]workflowAgent `json:"agents"`
		}
		if err = json.Unmarshal(raw, &inventory); err != nil {
			return out.unknown("inspection_unavailable", workflowError(4, "owner_unknown: cannot read complete Herdr "+kind+" inventory"))
		}
		rows := inventory.Panes
		if kind == "agent" {
			rows = inventory.Agents
		}
		if rows == nil {
			return out.unknown("inspection_unavailable", workflowError(4, "owner_unknown: incomplete Herdr "+kind+" inventory"))
		}
		seen := map[string]bool{}
		for _, a := range *rows {
			if !plain(a.WorkspaceID, 1, 128) || !plain(a.TabID, 1, 128) || !plain(a.PaneID, 1, 128) || !plain(a.TerminalID, 1, 128) || seen[a.PaneID] {
				return out.unknown("inspection_unavailable", workflowError(4, "owner_unknown: incomplete or ambiguous Herdr "+kind+" identity"))
			}
			seen[a.PaneID] = true
			if kind == "pane" && invoker != "" && a.PaneID == invoker && a.Agent == "" && a.Session == nil && a.PaneID != s.Result.Transport.PaneID && a.TerminalID != s.Result.Transport.TerminalID {
				continue
			}
			if a.PaneID == s.Result.Transport.PaneID || a.TerminalID == s.Result.Transport.TerminalID || a.Name == workflowAgentName(s.Result.RunID) || a.Session != nil && a.Session.Value == s.Result.Transport.AgentSessionID && a.Session.Agent == s.Request.Runtime.Provider || sourceContains(s.Observed.Target.WorktreeLocator, a.ForegroundCWD) || sourceContains(s.Observed.Target.WorktreeLocator, a.CWD) {
				return out.present("owner_active", "owner_active: stop the known Task/session writer in Herdr pane "+a.PaneID+" and leave the Task worktree before recovery")
			}
		}
		if kind == "pane" {
			out.absence.PaneInventorySHA = hash(raw)
		} else {
			out.absence.AgentInventorySHA = hash(raw)
		}
	}
	out.state = taskWritersAbsent
	return out
}

func (w taskWriters) boundIdentity(a workflowAgent) bool {
	t := w.s.Result.Transport
	return a.WorkspaceID == t.WorkspaceID && a.TabID == t.TabID && a.PaneID == t.PaneID && a.TerminalID == t.TerminalID && plain(a.PaneID, 1, 128) && plain(a.TerminalID, 1, 128)
}
