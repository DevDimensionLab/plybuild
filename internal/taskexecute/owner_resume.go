package taskexecute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type OwnerResumeInput struct {
	SpecID                 string
	Check                  bool
	Timeout                time.Duration
	DestinationWorkspaceID string
}

// OwnerResumeResult keeps preview, the actual start/session receipt, the
// functional callback and the current phase as separate public facts.
type OwnerResumeResult struct {
	Kind            string                              `json:"kind"`
	SchemaVersion   int                                 `json:"schema_version"`
	State           string                              `json:"state"`
	Run             *taskrun.OwnerRunCandidate          `json:"run,omitempty"`
	Choices         []taskrun.OwnerRunCandidate         `json:"choices,omitempty"`
	ProviderSupport string                              `json:"provider_support,omitempty"`
	Preview         *taskrun.DeliveryOwnerResumePreview `json:"preview,omitempty"`
	Workflow        *taskrun.WorkflowRun                `json:"workflow,omitempty"`
	InvokerPaneID   string                              `json:"invoker_pane_id,omitempty"`
	NextAction      string                              `json:"next_action"`
}

// providerSupport states the existing adapter honestly. Both use the
// provider's own resume of the same session id; neither forks a new session.
var providerSupport = map[string]string{
	"claude": "claude --resume SESSION_ID in a new Herdr tab (same session id; no --fork-session)",
	"codex":  "codex resume SESSION_ID in a new Herdr tab (same session id)",
}

// ResumeOwner finds the same existing Task Run from its Task or return
// worktree, previews the resume and, unless checking, applies it once.
func ResumeOwner(d taskrun.Dependencies, input OwnerResumeInput) (OwnerResumeResult, error) {
	out := OwnerResumeResult{Kind: "ply.workflow.execute-owner-resume", SchemaVersion: 1, State: "blocked"}
	if input.Timeout <= 0 {
		return out, fmt.Errorf("resume timeout must be positive")
	}
	ws, err := d.Workflow.Workspace.ObserveContaining()
	if err != nil {
		return out, err
	}
	root := ws.Observation.Root
	cwd, err := d.CWD()
	if err == nil {
		cwd, err = filepath.EvalSymlinks(cwd)
	}
	if err != nil {
		return out, err
	}
	all, err := taskrun.WorkflowOwnerRunCandidates(d, root)
	if err != nil {
		return out, err
	}
	sel, choices, next := selectOwnerRun(all.Runs, cwd, input.SpecID)
	if sel == nil {
		out.NextAction = next
		if len(choices) > 0 {
			out.State, out.Choices = "ambiguous", choices
		}
		return out, workspace.WorkInvalidArguments(next)
	}
	selected := *sel
	out.Run = &selected
	out.ProviderSupport = providerSupport[selected.Provider]
	// A missing provider is reported by the phase-aware preview, which keeps
	// completed, live-owner and blocked phases with their own next action.
	provider, err := lookupExecutable("", selected.Provider)
	if err != nil {
		provider = taskrun.Executable{}
	}
	out.InvokerPaneID = os.Getenv("HERDR_PANE_ID")
	controlPath, err := d.Executable()
	if err != nil {
		return out, err
	}
	control, err := bindExecutable(controlPath)
	if err != nil {
		return out, err
	}
	destination := input.DestinationWorkspaceID
	if destination == "" {
		destination = os.Getenv("HERDR_WORKSPACE_ID")
	}
	in := taskrun.DeliveryOwnerResumeInput{RunID: selected.RunID, ProviderExecutable: provider, ControlExecutable: control, DestinationWorkspaceID: destination, Timeout: input.Timeout, InvokerPaneID: out.InvokerPaneID}
	preview, err := taskrun.WorkflowPreviewOwnerResume(d, root, in)
	out.Preview = &preview
	out.State, out.NextAction = preview.State, preview.NextAction.Message
	if err != nil || preview.State != "ready" && preview.State != "continue_ready" {
		if preview.State == "" {
			out.State = "blocked"
		}
		return out, err
	}
	if input.Check {
		out.NextAction = "Run the same command without --check to resume this owner in one new Herdr tab."
		if preview.State == "continue_ready" {
			out.NextAction = "Run the same command without --check to send the one reserved continuation to the same live session."
		}
		return out, nil
	}
	if os.Getenv("HERDR_ENV") != "1" {
		return out, fmt.Errorf("owner resume requires a local Herdr terminal (HERDR_ENV=1)")
	}
	if preview.State == "ready" {
		preserved, err := PreserveControl(control, preview.Directory)
		if err != nil {
			return out, err
		}
		if preserved != preview.ControlExecutable {
			return out, fmt.Errorf("preserved resume control differs from the preview")
		}
	}
	run, err := taskrun.WorkflowResumeOwner(d, root, in, preview.Confirmation)
	out.Workflow = &run
	out.State, out.NextAction = "resumed", run.NextAction.Message
	if err != nil {
		out.State = "resume_unknown"
		if run.OwnerResume == nil {
			out.State = "blocked"
		}
	}
	return out, err
}

// selectOwnerRun prefers the Task worktree's own Run. From the return
// worktree, Runs that can no longer resume (completed, retired, released)
// only answer an explicit --spec when no active Run of it exists.
func selectOwnerRun(all []taskrun.OwnerRunCandidate, cwd, spec string) (*taskrun.OwnerRunCandidate, []taskrun.OwnerRunCandidate, string) {
	terminal := func(c taskrun.OwnerRunCandidate) bool {
		return c.Phase == "completed" || c.Phase == "retired" || c.Phase == "owner_released"
	}
	var own, active, ended []taskrun.OwnerRunCandidate
	for _, c := range all {
		switch {
		case spec != "" && c.SpecID != spec:
		case c.TaskWorktree == cwd:
			own = append(own, c)
		case c.ReturnWorktree == cwd && terminal(c):
			ended = append(ended, c)
		case c.ReturnWorktree == cwd:
			active = append(active, c)
		}
	}
	matches := own
	if len(matches) == 0 {
		matches = active
	}
	if len(matches) == 0 && spec != "" {
		matches = ended
	}
	switch {
	case len(matches) == 0:
		return nil, nil, "Run from the exact Task worktree, or from its return worktree with --spec SPEC_ID; no Run was found for " + cwd + "."
	case len(matches) > 1:
		choices := []string{}
		for _, c := range matches {
			choices = append(choices, "cd "+c.TaskWorktree+" && ply workflow execute --resume  (Task "+c.TaskID+", Spec "+c.SpecID+", "+c.Phase+")")
		}
		return nil, matches, "Several Runs match; nothing was started. Resume one from its own Task worktree: " + strings.Join(choices, "; ")
	}
	return &matches[0], nil, ""
}
