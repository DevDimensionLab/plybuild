package taskexecute

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type Input struct {
	Target              workspace.QueueTargetInput
	Next                bool
	SpecID              string
	Check               bool
	Runtime             RuntimeOptions
	NotificationContext string
}

type Result struct {
	Kind          string                            `json:"kind"`
	SchemaVersion int                               `json:"schema_version"`
	State         string                            `json:"state"`
	Goal          *workspace.TaskGoalExecutePreview `json:"goal"`
	Runtime       *RuntimePreview                   `json:"runtime"`
	Run           *taskrun.WorkflowRun              `json:"run"`
	RequestPath   string                            `json:"request_path"`
	NextAction    string                            `json:"next_action"`
}

type launchIntent struct {
	Kind          string                           `json:"kind"`
	SchemaVersion int                              `json:"schema_version"`
	Goal          workspace.TaskGoalExecutePreview `json:"goal"`
	Runtime       RuntimePreview                   `json:"runtime"`
	Human         workspace.QueueHumanDecision     `json:"human"`
	Notification  *taskrun.FileBinding             `json:"notification"`
}

func checkOrigin(d taskrun.Dependencies, target workspace.QueueTarget) error {
	cwd, err := d.CWD()
	if err != nil {
		return err
	}
	cwd, err = filepath.EvalSymlinks(cwd)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	b, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("execute must run inside its registered Epic worktree: %w", err)
	}
	top, err := filepath.EvalSymlinks(strings.TrimSpace(string(b)))
	if err != nil || top != target.ParentLocator {
		return fmt.Errorf("execute must start from its return worktree %s (current repository: %s)", target.ParentLocator, top)
	}
	return nil
}

func boundFile(path string) (*taskrun.FileBinding, error) {
	if path == "" {
		return nil, nil
	}
	if err := physicalPath(path); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(b) > 1<<20 {
		return nil, fmt.Errorf("execution input exceeds 1 MiB: %s", path)
	}
	return &taskrun.FileBinding{Locator: path, SHA256: digestBytes(b)}, nil
}

// Execute composes existing APIs. It preserves the launch intent before Task
// creation and reuses an existing request before selecting or starting again.
func Execute(d taskrun.Dependencies, input Input) (Result, error) {
	out := Result{Kind: "ply.workflow.execute", SchemaVersion: 1, State: "unknown"}
	goalInput := workspace.TaskGoalExecuteInput{Target: input.Target, Next: input.Next, SpecID: input.SpecID}
	plan, err := workspace.PreviewTaskGoalExecution(d.Workspace, goalInput)
	if err != nil {
		return out, err
	}
	out.Goal = &plan
	if err = checkOrigin(d, plan.Target); err != nil {
		return out, err
	}
	// A pre-preparation interruption may leave only launcher files. A later
	// explicit execution can bind a fresh queue/base without replacing those
	// files. After Spec publication, Preview recovers the original confirmation.
	directory := filepath.Join(filepath.Dir(plan.AcceptancePath), "attempts", strings.TrimPrefix(plan.Confirmation, "sha256:"))
	if err = physicalPath(directory); err != nil {
		return out, err
	}
	out.RequestPath = filepath.Join(directory, "workflow-request.json")
	if _, err = os.Lstat(out.RequestPath); err == nil {
		// Existing start state is authoritative even after the installed binary,
		// queue or provider configuration changes. No new provider is launched.
		var request taskrun.WorkflowRequest
		if _, err = readPreservedJSON(out.RequestPath, &request); err != nil {
			return out, err
		}
		if request.WorkspaceRoot != plan.Workspace || request.Delivery == nil || request.Delivery.AcceptancePath != plan.AcceptancePath || request.Delivery.Goal.SHA256 != plan.Goal.Goal.Spec.ManifestSHA256 || request.RequestKey != "execute/"+strings.TrimPrefix(plan.Confirmation, "sha256:") {
			return out, fmt.Errorf("preserved execution request differs from the selected goal")
		}
		idBytes, err := taskrun.Canonical([]string{request.WorkspaceRoot, request.RequestKey})
		if err != nil {
			return out, err
		}
		id := "wfr_" + strings.TrimPrefix(digestBytes(idBytes), "sha256:")
		run, err := taskrun.WorkflowShow(d, plan.Workspace, id)
		if os.IsNotExist(err) {
			// The immutable request may have been written before the transport was
			// called. Its native reservation logic owns any subsequent first start.
			return startRequest(d, out, input.Check)
		}
		if err != nil {
			return out, fmt.Errorf("inspect preserved start before retrying (%s): %w", out.RequestPath, err)
		}
		requestBytes, err := taskrun.Canonical(request)
		if err != nil || run.RequestSHA256 != digestBytes(requestBytes) {
			return out, fmt.Errorf("preserved request differs from the existing native reservation")
		}
		out.State, out.Run, out.NextAction = "existing", &run, run.NextAction.Message
		return out, nil
	} else if !os.IsNotExist(err) {
		return out, err
	}

	intentPath := filepath.Join(directory, "intent.json")
	var intent launchIntent
	old, err := readPreservedJSON(intentPath, &intent)
	if err == nil {
		actual, actualErr := taskrun.Canonical(intent.Goal)
		expected, expectedErr := taskrun.Canonical(plan)
		if intent.Kind != "ply.workflow.execute-intent" || intent.SchemaVersion != 1 || actualErr != nil || expectedErr != nil || !bytes.Equal(actual, expected) {
			return out, fmt.Errorf("preserved execution intent has changed inputs; inspect %s before another start", intentPath)
		}
		assignment := plan.Goal.Executor
		if intent.Runtime.Runtime.Provider != stringChoice(assignment.Provider, "codex") ||
			(assignment.Model != nil && intent.Runtime.Runtime.Model != *assignment.Model) ||
			(assignment.Effort != nil && intent.Runtime.ReasoningEffort != *assignment.Effort) {
			return out, fmt.Errorf("preserved runtime differs from the goal's assigned implementor")
		}
	} else if !os.IsNotExist(err) {
		return out, err
	} else {
		control, err := d.Executable()
		if err != nil {
			return out, err
		}
		runtime, err := PreviewRuntime(plan.Goal.Executor, control, input.Runtime)
		if err != nil {
			return out, err
		}
		notice, err := boundFile(input.NotificationContext)
		if err != nil {
			return out, err
		}
		intent = launchIntent{Kind: "ply.workflow.execute-intent", SchemaVersion: 1, Goal: plan, Runtime: runtime, Notification: notice,
			Human: workspace.QueueHumanDecision{ActorClaim: "Local execute caller", DecidedAtUTC: d.Now().UTC().Format(time.RFC3339Nano), Source: "explicit_human_instruction", Statement: "The local caller requested this exact goal from its registered return worktree through ply workflow execute; caller identity is a local attestation."},
		}
	}
	out.Runtime = &intent.Runtime
	if input.Check {
		out.State, out.NextAction = "ready", "Run the same execute command without --check to create the worktree and start the assigned interactive owner."
		return out, nil
	}
	if os.Getenv("HERDR_ENV") != "1" {
		return out, fmt.Errorf("execute start requires a local Herdr terminal (HERDR_ENV=1)")
	}
	if old == nil {
		intent.Runtime, err = preserveRuntime(intent.Runtime, directory)
		if err != nil {
			return out, err
		}
		b, err := taskrun.Canonical(intent)
		if err != nil {
			return out, err
		}
		if err = writeOnce(intentPath, b, 0600); err != nil {
			return out, err
		}
	}
	if d.Fault != nil {
		if err = d.Fault("execute_before_preparation"); err != nil {
			return out, err
		}
	}
	prepared, err := workspace.PrepareTaskGoalExecution(d.Workspace, intent.Goal.Input, intent.Goal.Confirmation, intent.Human)
	if err != nil {
		out.State = "preparation_pending"
		return out, err
	}
	if prepared.Preparation.Preparation == nil {
		return out, fmt.Errorf("goal preparation returned no immutable preparation")
	}
	preparation := *prepared.Preparation.Preparation
	key := "execute/" + strings.TrimPrefix(intent.Goal.Confirmation, "sha256:")
	draft, err := workflowhandoff.BuildDeliveryHandoffDraft(d.Workflow, preparation.ID, intent.Human.ActorClaim, intent.Runtime.Runtime.Provider+" delivery owner", key, plan.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{AllowLocalInstall: true})
	if err != nil {
		return out, err
	}
	var goal taskrun.FileBinding
	for _, in := range plan.Goal.RequiredInputs {
		if in.SHA256 == plan.Goal.Goal.Spec.ManifestSHA256 {
			goal = taskrun.FileBinding{Locator: in.Locator, SHA256: in.SHA256}
			break
		}
	}
	if goal.Locator == "" {
		return out, fmt.Errorf("goal manifest is missing from the preserved execution inputs")
	}
	pbytes, err := taskrun.Canonical(preparation)
	if err != nil {
		return out, err
	}
	r := intent.Runtime
	request, err := taskrun.BuildDeliveryWorkflowRequest(taskrun.DeliveryRequestInput{
		RequestKey: key, WorkspaceRoot: plan.Workspace, PreparationID: preparation.ID, PreparationSHA256: digestBytes(pbytes), HandoffDraft: draft, Runtime: r.Runtime,
		HumanAuthority:  taskrun.HumanAuthority{ActorClaim: intent.Human.ActorClaim, StartSurface: "human_authorized_herdr", Authorized: true},
		HerdrExecutable: r.Herdr, HerdrWorkspaceID: r.HerdrWorkspace, TabLabel: tabLabel(plan.Goal.Title),
		Delivery: taskrun.DeliveryContract{OwnerClaim: r.Runtime.Provider + " delivery owner", Goal: goal, AcceptancePath: plan.AcceptancePath, AllowSubagents: true, AllowLocalInstall: true, LocalIntegration: "after_human_pass", NotificationContext: intent.Notification, ReasoningEffort: r.ReasoningEffort},
	})
	if err != nil {
		return out, err
	}
	// Repository trust remains with the provider's existing configuration and
	// interactive onboarding. A delivery goal does not grant project trust.
	b, err := taskrun.Canonical(request)
	if err != nil {
		return out, err
	}
	if err = writeOnce(out.RequestPath, b, 0600); err != nil {
		return out, err
	}
	return startRequest(d, out, false)
}

func tabLabel(title string) string {
	r := []rune(title)
	if len(r) > 80 {
		return string(r[:77]) + "..."
	}
	return title
}

func startRequest(d taskrun.Dependencies, out Result, check bool) (Result, error) {
	v, err := taskrun.WorkflowPreviewStart(d, out.RequestPath)
	if err != nil {
		out.State = "start_pending"
		return out, err
	}
	if run, ok := v.(taskrun.WorkflowRun); ok {
		out.State, out.Run, out.NextAction = "existing", &run, run.NextAction.Message
		return out, nil
	}
	preview, ok := v.(taskrun.WorkflowPreview)
	if !ok || preview.Confirmation == nil {
		return out, fmt.Errorf("delivery startup did not provide a controlled preview")
	}
	if check {
		out.State = "ready"
		out.NextAction = "Run execute without --check to continue the preserved start."
		return out, nil
	}
	run, err := taskrun.WorkflowStart(d, out.RequestPath, *preview.Confirmation)
	out.State, out.Run, out.NextAction = "started", &run, run.NextAction.Message
	if err != nil {
		out.State = "start_unknown"
	}
	return out, err
}
