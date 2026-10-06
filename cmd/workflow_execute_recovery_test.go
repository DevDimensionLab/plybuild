package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestExecuteRecoveryRejectsInvalidInputsBeforeWorkspaceAccess(t *testing.T) {
	for _, args := range [][]string{
		{}, {"wfr_example", "extra"},
		{"wfr_example", "--timeout", "0"},
		{"wfr_example", "--timeout", "86401"},
		{"wfr_example", "--format", "yaml"},
		{"wfr_example", "--permission-profile", "different"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			// Empty dependencies would panic if validation reached a workspace
			// operation. Invalid input must have no recovery effects.
			c := newWorkflowExecuteRecoverStartCommand(taskrun.Dependencies{})
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			c.SetArgs(args)
			if err := c.Execute(); err == nil {
				t.Fatal("invalid recovery input was accepted")
			}
		})
	}
}

func TestExecuteRestartRequiresExplicitSpecBeforeWorkspaceAccess(t *testing.T) {
	for _, args := range [][]string{
		{"--restart"}, {"--restart", "--next"},
		{"--restart", "--spec", "goal", "--next"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			c := newWorkflowExecuteCommand(taskrun.Dependencies{})
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			c.SetArgs(args)
			if err := c.Execute(); err == nil {
				t.Fatal("ambiguous restart selector was accepted")
			}
		})
	}
}

func TestExecuteRecoveryOutputDistinguishesPlannedAndPerformedEffects(t *testing.T) {
	c := newWorkflowExecuteCommand(taskrun.Dependencies{})
	var output bytes.Buffer
	c.SetOut(&output)
	preview := &taskrun.DeliveryStartRecoveryPreview{
		State: "blocked", RunID: "wfr_fixture", ControlExecutable: taskrun.Executable{Path: "/not-created/ply-control"},
	}
	r := taskexecute.RecoveryResult{State: "blocked", Preview: preview}
	if err := writeExecuteResult(c, "text", r); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "control:") || strings.Contains(output.String(), "/not-created") {
		t.Fatalf("blocked preview claimed a control artifact: %s", output.String())
	}
	output.Reset()
	r.State, preview.State = "ready", "ready"
	preview.TransportMode, preview.DestinationWorkspaceID = "replace_missing_terminal", "destination"
	if err := writeExecuteResult(c, "text", r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "a new tab is planned") || !strings.Contains(output.String(), "Planned control:") {
		t.Fatalf("ready preview omitted its proposed effect: %s", output.String())
	}
	output.Reset()
	r.State, r.Run = "started", &taskrun.WorkflowRun{RunID: "wfr_fixture"}
	if err := writeExecuteResult(c, "text", r); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "planned") || strings.Contains(output.String(), "Planned control") {
		t.Fatalf("performed restart was still shown as planned: %s", output.String())
	}
}
