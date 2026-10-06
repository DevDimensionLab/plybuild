package cmd

import (
	"bytes"
	"strings"
	"testing"

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
