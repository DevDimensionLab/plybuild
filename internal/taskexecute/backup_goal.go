package taskexecute

import (
	"fmt"
	"os"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// BackupStartResult keeps publication and startup as separate observed facts.
// Registration-only callers retain the original AL27 response schema.
type BackupStartResult struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	State         string                   `json:"state"`
	Registration  taskrun.BackupGoalResult `json:"registration"`
	Confirmation  string                   `json:"confirmation"`
	Runtime       *RuntimePreview          `json:"runtime,omitempty"`
	Execution     *Result                  `json:"execution,omitempty"`
	Readiness     string                   `json:"readiness"`
	NextAction    string                   `json:"next_action"`
}

// BackupGoal adds explicitly selected interactive startup to AL27. Native
// execute owns all intents, reservations and uncertain transport effects.
func BackupGoal(d taskrun.Dependencies, root, source string, in taskrun.BackupGoalInput, runtime RuntimeOptions, apply bool, confirmation string) (out BackupStartResult, err error) {
	out = BackupStartResult{Kind: "ply.workflow.backup-goal-start", SchemaVersion: 1, State: "blocked", Readiness: "Not checked", NextAction: "Inspect the diagnosis; preserve the same source and fresh identity."}
	if runtime.Headless {
		return out, workspace.WorkInvalidArguments("backup start supports interactive native execute only")
	}
	// Resolve the destination once. A changed explicit destination/profile is
	// part of the opt-in confirmation and cannot silently select another start.
	if runtime.HerdrWorkspace == "" {
		runtime.HerdrWorkspace = os.Getenv("HERDR_WORKSPACE_ID")
	}
	registration, err := taskrun.WorkflowBackupGoal(d, root, source, in, false, "")
	out.Registration = registration
	if err != nil {
		return out, err
	}
	startConfirmation := boundHash(struct {
		Registration string
		Runtime      RuntimeOptions
	}{registration.Confirmation, runtime})
	if apply && confirmation != startConfirmation {
		return out, workspace.WorkInvalidArguments("backup start confirmation differs from previewed registration/runtime choices")
	}
	if err = checkOrigin(d, registration.Target); err != nil {
		return out, err
	}
	if registration.Queue.Current != nil && registration.Queue.Current.TaskID != in.TaskID && (registration.Retirement == nil || registration.Queue.Current.PreparationID != registration.Retirement.PreparationID) {
		return out, fmt.Errorf("backup_start_current_conflict: complete the current preparation before starting this backup")
	}
	target := workspace.QueueTargetInput{ProjectID: registration.Target.ProjectID, RepoID: registration.Target.RepoID, EpicID: registration.Target.EpicID}
	execute := Input{SpecID: in.SpecID, Target: target, ExpectedTaskID: in.TaskID, Runtime: runtime, Check: !apply}
	// Already registered copies use native execute's preserved request before
	// fresh runtime checks. This is readback, never restart/resume or fallback.
	queued := registration.Queue.Current != nil && registration.Queue.Current.TaskID == in.TaskID
	for _, pending := range registration.Queue.Pending {
		if pending.TaskID == in.TaskID && pending.Goal != nil && *pending.Goal == registration.Goal {
			queued = true
		}
	}
	if queued {
		out.Registration.Applied = true
	}
	var preservedRuntime *RuntimePreview
	if queued && registration.Goal.Spec.ManifestSHA256 != "" {
		execute.ExpectedGoal = &registration.Goal
		checked := execute
		checked.Check = true
		existing, e := Execute(d, checked)
		out.Execution = &existing
		if e != nil {
			out.State = "registered_start_blocked"
			out.NextAction = "Inspect the registered Goal and native execute diagnosis; repeat identical input only after its cause is resolved."
			return out, e
		}
		preservedRuntime = existing.Runtime
		if existing.State == "existing" {
			out.Runtime = preservedRuntime
			out.Confirmation = startConfirmation
			out.State, out.NextAction = "existing", existing.NextAction
			out.Readiness = "Existing native execution observed; no new startup or runtime acceptance implied."
			return out, nil
		}
	}
	var preview RuntimePreview
	if preservedRuntime != nil {
		preview = *preservedRuntime
	} else {
		control, e := d.Executable()
		if e != nil {
			return out, e
		}
		preview, err = PreviewRuntime(registration.Executor, control, runtime)
		if err != nil {
			if registration.Goal.Spec.ManifestSHA256 != "" {
				out.State = "registered_start_blocked"
			}
			return out, err
		}
	}
	out.Runtime = &preview
	out.Readiness = "Local runtime binaries, provider defaults/profile, return cwd and Herdr destination resolved. Native execute rechecks queue/base/resources and live Herdr at apply; no runtime permission or owner receipt observed."
	if preservedRuntime != nil {
		out.Readiness = "Preserved native runtime/destination/profile choices retained; omitted choices keep the existing launch. Native execute owns remaining startup checks; no owner receipt or runtime permission implied."
	}
	out.Confirmation = startConfirmation
	out.State = "ready"
	out.NextAction = "Repeat with --start --apply and this confirmation to register and start through native execute."
	if !apply {
		return out, nil
	}
	registration, err = taskrun.WorkflowBackupGoal(d, root, source, in, true, registration.Confirmation)
	out.Registration = registration
	if err != nil {
		out.State = "registration_unconfirmed"
		out.NextAction = "Inspect the exact Task publications and queue before repeating the identical input."
		return out, err
	}
	execute.ExpectedGoal = &registration.Goal
	execution, err := Execute(d, execute)
	out.Execution = &execution
	out.State = "registered_start_blocked"
	out.NextAction = "Registration is preserved. Inspect native execution and its next action; repeat identical input for readback, never automatically restart or resume an uncertain launch."
	if err == nil {
		out.State = execution.State
		out.NextAction = execution.NextAction
	}
	return out, err
}
