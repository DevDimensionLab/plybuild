package taskrun

// Startup recovery is a known additive historical format. These read-only
// types deliberately do not extend workflowState or any mutation validator.
// Executable identities are preserved metadata, never a present runtime probe
// or permission observation.
import (
	"fmt"
	"path/filepath"
)

type InventoryStartupRecovery struct {
	Generation                 int                `json:"generation"`
	Binding                    FileBinding        `json:"binding"`
	OriginalProviderExecutable Executable         `json:"original_provider_executable"`
	OriginalControlExecutable  Executable         `json:"original_control_executable"`
	ProviderExecutable         Executable         `json:"provider_executable"`
	ControlExecutable          Executable         `json:"control_executable"`
	TransportMode              string             `json:"transport_mode,omitempty"`
	OriginalTransport          *WorkflowTransport `json:"original_transport,omitempty"`
	ReplacementTransport       *FileBinding       `json:"replacement_transport,omitempty"`
}

type inventoryWorkflowResult struct {
	WorkflowRun
	StartupRecovery *InventoryStartupRecovery `json:"startup_recovery,omitempty"`
}

type inventoryWorkflowState struct {
	workflowState
	Result            inventoryWorkflowResult `json:"result"`
	Recovery          *FileBinding            `json:"startup_recovery,omitempty"`
	RecoveryTransport *FileBinding            `json:"startup_recovery_transport,omitempty"`
}

func (s inventoryWorkflowState) state() workflowState {
	out := s.workflowState
	out.Result = s.Result.WorkflowRun
	return out
}

type inventoryRecoveryRecord struct {
	Kind                   string                    `json:"kind"`
	SchemaVersion          int                       `json:"schema_version"`
	RunID                  string                    `json:"run_id"`
	RequestSHA256          string                    `json:"request_sha256"`
	BeforeState            FileBinding               `json:"before_state"`
	SourceControl          Executable                `json:"source_control"`
	ProviderExecutable     Executable                `json:"provider_executable"`
	ControlExecutable      Executable                `json:"control_executable"`
	Context                FileBinding               `json:"context"`
	ExitEvidence           *inventoryStartupExit     `json:"exit_evidence,omitempty"`
	Confirmation           string                    `json:"confirmation"`
	RecordedAtUTC          string                    `json:"recorded_at_utc"`
	Generation             int                       `json:"generation,omitempty"`
	TransportMode          string                    `json:"transport_mode,omitempty"`
	DestinationWorkspaceID string                    `json:"destination_workspace_id,omitempty"`
	AbsentTransport        *inventoryAbsentTransport `json:"absent_transport,omitempty"`
}

type inventoryStartupExit struct {
	WorkspaceID              string `json:"workspace_id"`
	TabID                    string `json:"tab_id"`
	PaneID                   string `json:"pane_id"`
	TerminalID               string `json:"terminal_id"`
	CWD                      string `json:"cwd"`
	ShellPID                 int    `json:"shell_pid"`
	ForegroundProcessGroupID int    `json:"foreground_process_group_id"`
	ShellName                string `json:"shell_name"`
	ProcessTree              []struct {
		PID            int `json:"pid"`
		ParentPID      int `json:"parent_pid"`
		ProcessGroupID int `json:"process_group_id"`
	} `json:"process_tree"`
	ManagedAgentState string `json:"managed_agent_state,omitempty"`
}

type inventoryAbsentTransport struct {
	Transport       WorkflowTransport `json:"transport"`
	PaneState       string            `json:"pane_state"`
	TerminalState   string            `json:"terminal_state"`
	InventorySHA256 string            `json:"inventory_sha256"`
	ProcessLiveness string            `json:"process_liveness"`
	TaskAuthority   string            `json:"task_authority"`
}

type inventoryRecoveryTransport struct {
	WorkspaceID string `json:"workspace_id"`
	TabID       string `json:"tab_id"`
	PaneID      string `json:"pane_id"`
	TerminalID  string `json:"terminal_id"`
	CWD         string `json:"cwd"`
}

type inventoryCapturedSource struct {
	binding FileBinding
	bytes   []byte
	max     int
}

func readInventoryRecovery(s inventoryWorkflowState) (sources []inventoryCapturedSource, err error) {
	readSource := func(binding FileBinding, max int) ([]byte, error) {
		raw, err := workflowBound(binding, max)
		if err == nil {
			sources = append(sources, inventoryCapturedSource{binding, raw, max})
		}
		return raw, err
	}
	metadata := s.Result.StartupRecovery
	if s.Recovery == nil {
		if metadata != nil || s.RecoveryTransport != nil {
			return sources, integrity("startup recovery metadata has no bound record")
		}
		return sources, nil
	}
	if metadata == nil || metadata.Generation < 1 || metadata.Binding != *s.Recovery || !deliveryRun(s.Request) {
		return sources, integrity("startup recovery record and metadata differ")
	}
	dir := filepath.Join(s.Result.Paths.RunRoot, "startup-recovery", fmt.Sprintf("%03d", metadata.Generation))
	if s.Recovery.Locator != filepath.Join(dir, "recovery.json") {
		return sources, integrity("startup recovery source path differs from its generation")
	}
	raw, err := readSource(*s.Recovery, 1<<20)
	if err != nil {
		return sources, err
	}
	var record inventoryRecoveryRecord
	if err = decode(raw, 1<<20, &record); err != nil {
		return sources, err
	}
	generation := record.Generation
	if generation == 0 {
		generation = 1
	}
	if record.Kind != "PlyDeliveryStartupRecovery@1" || record.SchemaVersion != 1 || generation != metadata.Generation || record.RunID != s.Result.RunID || record.RequestSHA256 != s.Result.RequestSHA256 || record.BeforeState.Locator != filepath.Join(dir, "before-state.json") || record.ControlExecutable.Path != filepath.Join(dir, "ply-control") || record.Context.Locator != filepath.Join(dir, "context.json") || record.Context.Locator != s.Result.Paths.Context || record.Context.SHA256 != s.ContextSHA256 || !digestPattern.MatchString(record.Confirmation) {
		return sources, integrity("startup recovery source identity differs from its preserved run")
	}
	if _, err = parseInventoryTime(record.RecordedAtUTC); err != nil {
		return sources, err
	}
	for _, executable := range []Executable{record.SourceControl, record.ProviderExecutable, record.ControlExecutable} {
		if !filepath.IsAbs(executable.Path) || filepath.Clean(executable.Path) != executable.Path || !digestPattern.MatchString(executable.SHA256) {
			return sources, integrity("startup recovery executable metadata is invalid")
		}
	}
	if record.SourceControl.SHA256 != record.ControlExecutable.SHA256 || metadata.OriginalProviderExecutable != s.Request.Runtime.Executable || metadata.OriginalControlExecutable != s.Request.Runtime.PlyExecutable || metadata.ProviderExecutable != record.ProviderExecutable || metadata.ControlExecutable != record.ControlExecutable || metadata.TransportMode != record.TransportMode || !equal(metadata.ReplacementTransport, s.RecoveryTransport) {
		return sources, integrity("startup recovery readback differs from its bound record")
	}
	before, err := readSource(record.BeforeState, 8<<20)
	if err != nil {
		return sources, err
	}
	var old inventoryWorkflowState
	if err = decode(before, 8<<20, &old); err != nil {
		return sources, err
	}
	previousGeneration := 0
	if old.Recovery != nil {
		if old.Result.StartupRecovery == nil || old.Result.StartupRecovery.Binding != *old.Recovery {
			return sources, integrity("previous startup recovery generation has no matching binding")
		}
		previousGeneration = old.Result.StartupRecovery.Generation
	} else if old.Result.StartupRecovery != nil || old.RecoveryTransport != nil {
		return sources, integrity("previous startup recovery metadata has no source")
	}
	if previousGeneration != metadata.Generation-1 || !equal(old.Request, s.Request) || !equal(old.Observed, s.Observed) || old.Result.RunID != s.Result.RunID || old.Result.RequestSHA256 != s.Result.RequestSHA256 || old.Result.SessionID != s.Result.SessionID || old.Result.Handoff != s.Result.Handoff {
		return sources, integrity("startup recovery before-state differs from its preserved run")
	}
	var previousTransport *WorkflowTransport
	if record.TransportMode != "" {
		previousTransport = &old.Result.Transport
	}
	if !equal(metadata.OriginalTransport, previousTransport) {
		return sources, integrity("startup recovery original transport differs from before-state")
	}
	context, err := readSource(record.Context, 1<<20)
	if err != nil {
		return sources, err
	}
	var claim workflowContext
	if err = decode(context, 1<<20, &claim); err != nil {
		return sources, err
	}
	if claim.Envelope != workflowEnv("run-context") || claim.RunID != s.Result.RunID || claim.RequestSHA256 != s.Result.RequestSHA256 || claim.SessionID != s.Result.SessionID || claim.HandoffSHA256 != s.Result.Handoff.SHA256 || claim.PlyExecutable != record.ControlExecutable {
		return sources, integrity("startup recovery context differs from its run binding")
	}
	identity := func(t WorkflowTransport) inventoryRecoveryTransport {
		return inventoryRecoveryTransport{WorkspaceID: t.WorkspaceID, TabID: t.TabID, PaneID: t.PaneID, TerminalID: t.TerminalID}
	}
	previous, current := identity(old.Result.Transport), identity(s.Result.Transport)
	switch record.TransportMode {
	case "", "reuse_terminal":
		if record.ExitEvidence == nil || record.AbsentTransport != nil || s.RecoveryTransport != nil || previous != current {
			return sources, integrity("startup recovery reused transport differs from its record")
		}
	case "replace_missing_terminal":
		absent := record.AbsentTransport
		if record.ExitEvidence != nil || absent == nil || !equal(absent.Transport, old.Result.Transport) || absent.PaneState != "not_found" || absent.TerminalState != "absent_from_all_panes" || absent.ProcessLiveness != "unknown" || absent.TaskAuthority != "not_issued" || !digestPattern.MatchString(absent.InventorySHA256) || !plain(record.DestinationWorkspaceID, 1, 128) {
			return sources, integrity("startup recovery replacement transport has no matching recorded absence")
		}
		if s.RecoveryTransport == nil {
			if s.Phase != "recovery_tab_create_attempted" || previous != current {
				return sources, integrity("startup recovery replacement transport is unbound")
			}
			return sources, nil
		}
		if s.RecoveryTransport.Locator != filepath.Join(dir, "transport.json") {
			return sources, integrity("startup recovery replacement transport source path differs")
		}
		raw, err := readSource(*s.RecoveryTransport, 1<<20)
		if err != nil {
			return sources, err
		}
		var transport inventoryRecoveryTransport
		if err = decode(raw, 1<<20, &transport); err != nil {
			return sources, err
		}
		if transport.CWD != s.Observed.Target.WorktreeLocator || transport.WorkspaceID != record.DestinationWorkspaceID || transport.TerminalID == previous.TerminalID || transport.PaneID == previous.PaneID {
			return sources, integrity("startup recovery replacement transport identity differs")
		}
		transport.CWD = ""
		if transport != current {
			return sources, integrity("startup recovery replacement transport readback differs")
		}
	default:
		return sources, integrity("unrecognized startup recovery transport mode")
	}
	return sources, nil
}
