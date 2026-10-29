package taskrun

// Startup recovery has both a native contract and a historical projection.
// These read-only types preserve recorded metadata without running native
// preflight or extending mutation authority. Executable identities are never
// a present runtime probe or permission observation.
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
	var record deliveryStartRecoveryRecord
	if err = decode(raw, 1<<20, &record); err != nil {
		return sources, err
	}
	if record.Kind != "PlyDeliveryStartupRecovery@1" || record.SchemaVersion != 1 || record.Generation != metadata.Generation || record.RunID != s.Result.RunID || record.RequestSHA256 != s.Result.RequestSHA256 || record.ControlExecutable.Path != filepath.Join(dir, "ply-control") || record.Context.Locator != filepath.Join(dir, "context.json") || record.Context.Locator != s.Result.Paths.Context || record.Context.SHA256 != s.ContextSHA256 || !digestPattern.MatchString(record.Confirmation) {
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
	previousTransport := record.PreviousTransport
	if record.AbsentTransport != nil {
		previousTransport = &record.AbsentTransport.Transport
	}
	if previousTransport == nil || !equal(metadata.OriginalTransport, previousTransport) {
		return sources, integrity("startup recovery original transport differs from its record")
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
	previous, current := deliveryTransportIdentity(*previousTransport), deliveryTransportIdentity(s.Result.Transport)
	switch record.TransportMode {
	case deliveryRecoveryReuseTerminal:
		if record.ExitEvidence == nil || record.PreviousTransport == nil || record.AbsentTransport != nil || s.RecoveryTransport != nil || previous != current || previous.WorkspaceID != record.ExitEvidence.WorkspaceID || previous.TabID != record.ExitEvidence.TabID || previous.PaneID != record.ExitEvidence.PaneID || previous.TerminalID != record.ExitEvidence.TerminalID || record.ExitEvidence.CWD != s.Observed.Target.WorktreeLocator {
			return sources, integrity("startup recovery reused transport differs from its record")
		}
	case deliveryRecoveryReplaceTerminal:
		absent := record.AbsentTransport
		if record.ExitEvidence != nil || record.PreviousTransport != nil || absent == nil || absent.PaneState != "not_found" || absent.TerminalState != "absent_from_all_panes" || absent.ProcessLiveness != "unknown" || absent.TaskAuthority != "not_issued" || !digestPattern.MatchString(absent.InventorySHA256) || !plain(record.DestinationWorkspaceID, 1, 128) {
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
		var transport deliveryRecoveryTransport
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
