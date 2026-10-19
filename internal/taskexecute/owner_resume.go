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
	var matches []taskrun.OwnerRunCandidate
	inTask := false
	for _, c := range all.Runs {
		if input.SpecID != "" && c.SpecID != input.SpecID {
			continue
		}
		if c.TaskWorktree == cwd {
			if !inTask {
				matches, inTask = nil, true
			}
			matches = append(matches, c)
		} else if !inTask && c.ReturnWorktree == cwd && (input.SpecID != "" || c.Phase != "completed") {
			matches = append(matches, c)
		}
	}
	switch {
	case len(matches) == 0:
		out.NextAction = "Run from the exact Task worktree, or from its return worktree with --spec SPEC_ID; no Run was found for " + cwd + "."
		return out, workspace.WorkInvalidArguments(out.NextAction)
	case len(matches) > 1:
		out.State, out.Choices = "ambiguous", matches
		ids := []string{}
		for _, c := range matches {
			ids = append(ids, c.SpecID+" ("+c.Phase+")")
		}
		out.NextAction = "Choose one with --spec: " + strings.Join(ids, ", ") + ". Nothing was started."
		return out, workspace.WorkInvalidArguments(out.NextAction)
	}
	selected := matches[0]
	out.Run = &selected
	out.ProviderSupport = providerSupport[selected.Provider]
	provider, err := lookupExecutable("", selected.Provider)
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
	destination := input.DestinationWorkspaceID
	if destination == "" {
		destination = os.Getenv("HERDR_WORKSPACE_ID")
	}
	in := taskrun.DeliveryOwnerResumeInput{RunID: selected.RunID, ProviderExecutable: provider, ControlExecutable: control, DestinationWorkspaceID: destination, Timeout: input.Timeout}
	preview, err := taskrun.WorkflowPreviewOwnerResume(d, root, in)
	out.Preview = &preview
	out.State, out.NextAction = preview.State, preview.NextAction.Message
	if err != nil || preview.State != "ready" {
		if preview.State == "" {
			out.State = "blocked"
		}
		return out, err
	}
	if input.Check {
		out.NextAction = "Run the same command without --check to resume this owner in one new Herdr tab."
		return out, nil
	}
	if os.Getenv("HERDR_ENV") != "1" {
		return out, fmt.Errorf("owner resume requires a local Herdr terminal (HERDR_ENV=1)")
	}
	preserved, err := PreserveControl(control, preview.Directory)
	if err != nil {
		return out, err
	}
	if preserved != preview.ControlExecutable {
		return out, fmt.Errorf("preserved resume control differs from the preview")
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
