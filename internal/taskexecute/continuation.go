package taskexecute

import (
	"errors"
	"fmt"
	"os"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type ContinuationInput struct {
	RunID               string
	ContextPath         string
	Check               bool
	RuntimeEvidencePath string
}

type ContinuationResult struct {
	Kind          string                               `json:"kind"`
	SchemaVersion int                                  `json:"schema_version"`
	State         string                               `json:"state"`
	Preview       *taskrun.DeliveryContinuationPreview `json:"preview,omitempty"`
	Run           *taskrun.WorkflowRun                 `json:"run,omitempty"`
	NextAction    string                               `json:"next_action"`
	Diagnostic    *ContinuationDiagnostic              `json:"diagnostic,omitempty"`
}

type ContinuationDiagnostic struct {
	Code             string `json:"code"`
	Operation        string `json:"operation"`
	Artifact         string `json:"artifact,omitempty"`
	ArtifactSchema   string `json:"artifact_schema,omitempty"`
	ReaderCapability string `json:"reader_capability"`
	Cause            string `json:"cause"`
	NextAction       string `json:"next_action"`
}

func continuationDiagnostic(err error) *ContinuationDiagnostic {
	d := &ContinuationDiagnostic{Code: "delivery_continuation_observation_unknown", Operation: "workflow.execute.continue", ArtifactSchema: "ply.workflow.herdr-run-request@2", ReaderCapability: taskrun.DeliveryContinuationReaderCapability, Cause: err.Error(), NextAction: "Inspect the preserved run and exact cause with the installed Ply before retrying; keep the same session and frozen artifacts."}
	var content *workspace.TaskContentError
	if errors.As(err, &content) {
		d.Code, d.Artifact, d.ArtifactSchema = content.Code, content.Artifact, content.Schema
		if content.ReaderCapability != "" {
			d.ReaderCapability = content.ReaderCapability
		}
		if content.NextAction != "" {
			d.NextAction = content.NextAction
		}
		return d
	}
	var native *taskrun.Error
	if errors.As(err, &native) {
		d.Code = native.Code
	} else {
		var path *os.PathError
		if errors.As(err, &path) {
			d.Code, d.Artifact = "delivery_required_artifact_unavailable", path.Path
			d.NextAction = "Inspect the required artifact identified in the cause against preserved history; keep frozen inputs intact. Use the bounded incomplete-report path if its identity dependencies remain available."
		}
	}
	switch d.Code {
	case "delivery_continuation_unsupported":
		d.NextAction = "Use a Ply reader with explicit support for the required contract; preserve the same session, original controls and publication."
	case "delivery_continuation_authority":
		d.NextAction = "Resolve the exact missing or changed runtime authority with the original owner; this continuation cannot grant it."
	case "delivery_continuation_session":
		d.NextAction = "Observe the original live provider/session from its exact Task worktree; do not restart or resend input."
	case "delivery_runtime_observation_required", "delivery_runtime_observation":
		d.NextAction = "The original owner must preserve a current runtime observation and pass --runtime-evidence; observe actual policy, do not copy requested policy or infer authority from an installed provider."
	case "delivery_launch_dependency":
		d.NextAction = "The historical provider launcher is required for this uncontinued operation. An accepted live owner can use continue with current --runtime-evidence; an unaccepted or future launch needs its own supported checks."
	case "delivery_session_dependency":
		d.NextAction = "Restore supported live-session observation through the bound Herdr contract; do not restart or substitute a provider. The provider launcher itself is not a session observation."
	case "delivery_runtime_dependency":
		d.NextAction = "Inspect the required control or launch-contract evidence identified in the cause; preserve its bound bytes. Use report --incomplete for a truthful owner report when runtime dependencies prevent continuation."
	case "delivery_effect_unknown":
		d.NextAction = "Inspect the reserved native effect and its receipt before any retry; no command or input may be replayed while its result is unknown."
	case "delivery_continuation_stale":
		d.NextAction = "Inspect the changed run or worktree, then repeat the native continue command to obtain fresh compatibility observations."
	case "delivery_continuation_context":
		d.NextAction = "Use the exact original private context reported by workflow execute show; leave its bytes unchanged."
	case "delivery_continuation_integrity", "task_run_integrity", "task_run_conflict":
		d.NextAction = "Inspect the changed or missing bound artifact against preserved history; do not overwrite frozen bytes or bypass its hash check."
	}
	return d
}

// ContinueDelivery preserves a compatible control after the original native
// session accepted its mandate. It neither starts a process nor replaces any
// frozen control, context, request or accepted permission evidence.
func ContinueDelivery(d taskrun.Dependencies, input ContinuationInput) (out ContinuationResult, err error) {
	out = ContinuationResult{Kind: "ply.workflow.execute-continuation", SchemaVersion: 1, State: "blocked"}
	defer func() {
		if err != nil {
			out.State, out.Diagnostic = "blocked", continuationDiagnostic(err)
			out.NextAction = out.Diagnostic.NextAction
		}
	}()
	ws, err := d.Workflow.Workspace.ObserveContaining()
	if err != nil {
		return out, err
	}
	path, err := d.Executable()
	if err != nil {
		return out, err
	}
	control, err := bindExecutable(path)
	if err != nil {
		return out, err
	}
	in := taskrun.DeliveryContinuationInput{RunID: input.RunID, ContextPath: input.ContextPath, ControlExecutable: control, RuntimeEvidencePath: input.RuntimeEvidencePath}
	p, err := taskrun.WorkflowPreviewDeliveryContinuation(d, ws.Observation.Root, in)
	out.Preview = &p
	if err != nil {
		out.NextAction = err.Error()
		return out, err
	}
	out.State, out.NextAction = p.State, p.NextAction
	if p.State == "existing" || input.Check {
		return out, nil
	}
	preserved, err := PreserveControl(control, p.Directory)
	if err != nil {
		out.State = "blocked"
		return out, err
	}
	if preserved != p.ControlExecutable {
		return out, fmt.Errorf("preserved continuation control differs from the observed preview")
	}
	p, err = taskrun.WorkflowContinueDelivery(d, ws.Observation.Root, in, p.Confirmation)
	out.Preview, out.State, out.NextAction = &p, p.State, p.NextAction
	if err != nil {
		out.State, out.NextAction = "blocked", err.Error()
		return out, err
	}
	run, err := taskrun.WorkflowShow(d, ws.Observation.Root, input.RunID)
	out.Run = &run
	return out, err
}
