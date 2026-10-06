package taskrun

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"sort"
	"time"
)

const (
	deliveryRecoveryReuseTerminal   = "reuse_terminal"
	deliveryRecoveryReplaceTerminal = "replace_missing_terminal"
)

// Missing transport is not evidence that an OS process exited. Replacement is
// authorized solely before Task authority was issued, under the generation lock.
type DeliveryAbsentTransportEvidence struct {
	Transport       WorkflowTransport `json:"transport"`
	PaneState       string            `json:"pane_state"`
	TerminalState   string            `json:"terminal_state"`
	InventorySHA256 string            `json:"inventory_sha256"`
	ProcessLiveness string            `json:"process_liveness"`
	TaskAuthority   string            `json:"task_authority"`
}

type deliveryRecoveryTransport struct {
	WorkspaceID string `json:"workspace_id"`
	TabID       string `json:"tab_id"`
	PaneID      string `json:"pane_id"`
	TerminalID  string `json:"terminal_id"`
	CWD         string `json:"cwd"`
}

func deliveryTransportIdentity(t WorkflowTransport) deliveryRecoveryTransport {
	return deliveryRecoveryTransport{WorkspaceID: t.WorkspaceID, TabID: t.TabID, PaneID: t.PaneID, TerminalID: t.TerminalID}
}

func deliveryRecoveryNoManagedAgent(d Dependencies, s workflowState) error {
	_, err := workflowCall(d, s.Request, "agent", "get", s.Result.Transport.PaneID)
	if err == nil {
		return workflowError(4, "startup transport still has an active or pending managed agent")
	}
	var call *HerdrCallError
	if !errors.As(err, &call) || !call.ReadOnly || call.Phase != "agent get" || call.ProviderCode != "agent_not_found" {
		return err
	}
	return nil
}

func deliveryRecoveryValidateTransport(s, old workflowState, r deliveryStartRecoveryRecord) error {
	previous := deliveryTransportIdentity(old.Result.Transport)
	current := deliveryTransportIdentity(s.Result.Transport)
	switch r.TransportMode {
	case "", deliveryRecoveryReuseTerminal:
		if r.ExitEvidence == nil || r.AbsentTransport != nil || s.RecoveryTransport != nil || previous != current {
			return workflowError(4, "reused startup transport binding differs")
		}
	case deliveryRecoveryReplaceTerminal:
		if r.ExitEvidence != nil || r.AbsentTransport == nil || !equal(r.AbsentTransport.Transport, old.Result.Transport) || r.AbsentTransport.PaneState != "not_found" || r.AbsentTransport.TerminalState != "absent_from_all_panes" || r.AbsentTransport.ProcessLiveness != "unknown" || r.AbsentTransport.TaskAuthority != "not_issued" || !plain(r.DestinationWorkspaceID, 1, 128) {
			return workflowError(4, "replacement transport lacks its bound absence evidence")
		}
		if s.RecoveryTransport == nil {
			if s.Phase != "recovery_tab_create_attempted" || previous != current {
				return workflowError(4, "replacement tab has no confirmed transport identity")
			}
			return nil
		}
		if s.RecoveryTransport.Locator != filepath.Join(deliveryRecoveryDirectory(s), "transport.json") {
			return workflowError(4, "replacement transport path differs")
		}
		b, err := workflowBound(*s.RecoveryTransport, 1<<20)
		if err != nil {
			return err
		}
		var transport deliveryRecoveryTransport
		if err = decode(b, 1<<20, &transport); err != nil {
			return err
		}
		if transport.CWD != s.Observed.Target.WorktreeLocator || transport.WorkspaceID != r.DestinationWorkspaceID || transport.TerminalID == previous.TerminalID || transport.PaneID == previous.PaneID {
			return workflowError(4, "replacement transport identity differs")
		}
		transport.CWD = ""
		if transport != current {
			return workflowError(4, "replacement transport readback differs")
		}
	default:
		return workflowError(4, "unrecognized startup transport replacement mode")
	}
	return nil
}

func deliveryRecoveryObserveTransport(d Dependencies, s workflowState, destination string) (*DeliveryStartupExitEvidence, *DeliveryAbsentTransportEvidence, error) {
	exit, err := deliveryRecoveryExit(d, s)
	if err == nil {
		return &exit, nil, nil
	}
	var call *HerdrCallError
	if !errors.As(err, &call) || !call.ReadOnly || call.Phase != "pane get" || call.ProviderCode != "pane_not_found" {
		return nil, nil, err
	}
	b, err := workflowCall(d, s.Request, "pane", "list")
	if err != nil {
		return nil, nil, err
	}
	var result struct {
		Panes *[]deliveryRecoveryTransport `json:"panes"`
	}
	if err = json.Unmarshal(b, &result); err != nil || result.Panes == nil {
		return nil, nil, workflowError(4, "complete Herdr pane inventory is unavailable")
	}
	seen := map[string]bool{}
	for _, pane := range *result.Panes {
		if !plain(pane.WorkspaceID, 1, 128) || !plain(pane.TabID, 1, 128) || !plain(pane.PaneID, 1, 128) || !plain(pane.TerminalID, 1, 128) || seen[pane.PaneID] {
			return nil, nil, workflowError(4, "Herdr pane inventory contains incomplete or ambiguous identities")
		}
		seen[pane.PaneID] = true
		if pane.TerminalID == s.Result.Transport.TerminalID || pane.PaneID == s.Result.Transport.PaneID {
			return nil, nil, workflowError(4, "the previous startup transport is still present or moved; inspect its bound terminal")
		}
	}
	// Keep only identities in the evidence digest. Pane metadata and terminal
	// contents are neither needed nor preserved for this negative lookup.
	panes := *result.Panes
	for i := range panes {
		panes[i].CWD = ""
	}
	sort.Slice(panes, func(i, j int) bool { return panes[i].PaneID < panes[j].PaneID })
	b, err = workflowCall(d, s.Request, "workspace", "get", destination)
	if err != nil {
		return nil, nil, err
	}
	var workspace struct {
		Workspace struct {
			ID string `json:"workspace_id"`
		} `json:"workspace"`
	}
	if err = json.Unmarshal(b, &workspace); err != nil || workspace.Workspace.ID != destination {
		return nil, nil, workflowError(4, "replacement destination workspace identity is unavailable")
	}
	return nil, &DeliveryAbsentTransportEvidence{Transport: s.Result.Transport, PaneState: "not_found", TerminalState: "absent_from_all_panes", InventorySHA256: digest(panes), ProcessLiveness: "unknown", TaskAuthority: "not_issued"}, nil
}

func deliveryRecoveryCreateTab(d Dependencies, owned workflowState, deadline time.Time) (workflowState, error) {
	generation := workflowStartupGeneration(owned)
	recovery, err := workflowRecoveryRecord(owned)
	if err != nil {
		return owned, err
	}
	if err = d.fault("delivery_recovery_before_tab_send"); err != nil {
		return owned, err
	}
	if err = workflowTrustFresh(owned); err != nil {
		return owned, err
	}
	b, err := workflowStartupCall(d, owned, deadline, workflowTabCreateArgv(owned, recovery.DestinationWorkspaceID)...)
	if err != nil {
		return owned, err
	}
	var tab struct {
		RootPane workflowAgent `json:"root_pane"`
	}
	if err = json.Unmarshal(b, &tab); err != nil {
		return owned, workflowError(5, "replacement tab response has no confirmed identity; do not repeat creation")
	}
	a := tab.RootPane
	cwd := a.ForegroundCWD
	if cwd == "" {
		cwd = a.CWD
	}
	if a.WorkspaceID != recovery.DestinationWorkspaceID || !plain(a.TabID, 1, 128) || !plain(a.PaneID, 1, 128) || !plain(a.TerminalID, 1, 128) || a.TerminalID == owned.Result.Transport.TerminalID || a.PaneID == owned.Result.Transport.PaneID || cwd != owned.Observed.Target.WorktreeLocator {
		return owned, workflowError(5, "replacement tab identity or Task directory is unknown; do not repeat creation")
	}
	transport := deliveryRecoveryTransport{WorkspaceID: a.WorkspaceID, TabID: a.TabID, PaneID: a.PaneID, TerminalID: a.TerminalID, CWD: cwd}
	var bound workflowState
	err = workflowUpdate(d, owned.Request.WorkspaceRoot, owned.Result.RunID, func(s *workflowState) error {
		if workflowStartupGeneration(*s) != generation || s.Phase != "recovery_tab_create_attempted" || s.RecoveryTransport != nil {
			return workflowError(4, "startup generation changed before replacement transport binding")
		}
		binding, err := workflowKeep(d, filepath.Join(deliveryRecoveryDirectory(*s), "transport.json"), transport)
		if err != nil {
			return err
		}
		s.RecoveryTransport = &binding
		s.Result.StartupRecovery.ReplacementTransport = &binding
		s.Result.Transport = WorkflowTransport{WorkspaceID: transport.WorkspaceID, TabID: transport.TabID, PaneID: transport.PaneID, TerminalID: transport.TerminalID, State: "unknown", Observation: "fresh", ObservedAt: d.Now().UTC().Format(time.RFC3339Nano)}
		if err = d.writeValue(workflowStartupPath(*s, "agent-start-attempt.json"), map[string]any{"recovery_sha256": s.Recovery.SHA256, "transport_sha256": binding.SHA256, "argv": workflowStartArgv(s.Request, s.Observed.Target.WorktreeLocator, s.Result.Transport.PaneID)}); err != nil {
			return err
		}
		s.Phase = "agent_start_attempted"
		bound = *s
		return nil
	})
	return bound, err
}
