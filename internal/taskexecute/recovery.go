package taskexecute

import (
	"fmt"
	"os"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// RecoveryInput keeps the original goal and launch choices. Only the installed
// provider and control executable are rebound for an explicitly requested start.
type RecoveryInput struct {
	RunID              string
	Check              bool
	Timeout            time.Duration
	ProviderExecutable string
}

type RecoveryResult struct {
	Kind          string                                `json:"kind"`
	SchemaVersion int                                   `json:"schema_version"`
	State         string                                `json:"state"`
	Preview       *taskrun.DeliveryStartRecoveryPreview `json:"preview,omitempty"`
	Run           *taskrun.WorkflowRun                  `json:"run,omitempty"`
	NextAction    string                                `json:"next_action"`
}

// RecoverStartup composes a read-only preview, control preservation and one
// guarded replacement attempt. The native run owns all reservations and sends.
func RecoverStartup(d taskrun.Dependencies, input RecoveryInput) (RecoveryResult, error) {
	out := RecoveryResult{Kind: "ply.workflow.execute-recovery", SchemaVersion: 1, State: "blocked"}
	if input.Timeout <= 0 {
		return out, fmt.Errorf("recovery timeout must be positive")
	}
	ws, err := d.Workflow.Workspace.ObserveContaining()
	if err != nil {
		return out, err
	}
	root := ws.Observation.Root
	run, err := taskrun.WorkflowShow(d, root, input.RunID)
	if err != nil {
		return out, err
	}
	// Readback of a reserved replacement must not depend on today's PATH or
	// start a second provider after an earlier caller lost its response.
	if run.StartupRecovery != nil {
		out.State, out.Run, out.NextAction = "existing", &run, run.NextAction.Message
		return out, nil
	}
	if !input.Check && os.Getenv("HERDR_ENV") != "1" {
		return out, fmt.Errorf("startup recovery requires a local Herdr terminal (HERDR_ENV=1)")
	}
	provider, err := lookupExecutable(input.ProviderExecutable, "codex")
	if err != nil {
		return out, err
	}
	controlPath, err := d.Executable()
	if err != nil {
		return out, err
	}
	control, err := bindExecutable(controlPath)
	if err != nil {
		return out, err
	}
	in := taskrun.DeliveryStartRecoveryInput{RunID: input.RunID, ProviderExecutable: provider, ControlExecutable: control, Timeout: input.Timeout}
	preview, err := taskrun.WorkflowPreviewDeliveryStartRecovery(d, root, in)
	out.Preview = &preview
	if err != nil {
		return out, err
	}
	if preview.State == "existing" && preview.Run != nil {
		out.State, out.Run, out.NextAction = "existing", preview.Run, preview.Run.NextAction.Message
		return out, nil
	}
	if preview.Confirmation == "" || preview.RecoveryDirectory == "" {
		return out, fmt.Errorf("startup recovery did not provide a controlled preview")
	}
	out.State = "ready"
	out.NextAction = "Run the same recover-start command without --check to make one replacement startup in the existing Task tab."
	if input.Check {
		return out, nil
	}
	preserved, err := PreserveControl(control, preview.RecoveryDirectory)
	if err != nil {
		out.State = "blocked"
		return out, err
	}
	if preserved != preview.ControlExecutable {
		out.State = "blocked"
		return out, fmt.Errorf("preserved recovery control differs from the preview")
	}
	run, err = taskrun.WorkflowRecoverDeliveryStart(d, root, in, preview.Confirmation)
	out.State, out.Run, out.NextAction = "started", &run, run.NextAction.Message
	if err != nil {
		out.State = "blocked"
		if run.StartupRecovery != nil {
			out.State = "start_unknown"
		}
	}
	return out, err
}
