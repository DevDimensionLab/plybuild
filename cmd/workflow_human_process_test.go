package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestHumanOwnedWorkflowRejectsAutomaticControllers(t *testing.T) {
	workflow := newWorkflowCommand()
	for _, command := range workflow.Commands() {
		switch command.Name() {
		case "run", "handoff", "execute", "week", "worker", "control", "delivery":
			t.Errorf("process controller command %q is still available", command.Name())
		}
	}
}

func TestRemovedWorkflowCommandsReturnErrorWithoutFlags(t *testing.T) {
	for _, name := range []string{"run", "handoff", "execute", "week", "worker", "control", "delivery", "typo"} {
		command := newWorkflowCommand()
		command.SetOut(&bytes.Buffer{})
		command.SetErr(&bytes.Buffer{})
		root := &cobra.Command{Use: "ply"}
		root.AddCommand(command)
		root.SetArgs([]string{"workflow", name})
		if err := root.Execute(); err == nil {
			t.Errorf("workflow %s reported success", name)
		}
	}
	command := newWorkflowCommand()
	command.SetOut(&bytes.Buffer{})
	command.SetArgs([]string{})
	if err := command.Execute(); err != nil {
		t.Fatal("parent help failed", err)
	}
}
