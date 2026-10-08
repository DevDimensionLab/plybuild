package cmd

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestIntegrationRecoveryCommandsExposeExplicitHumanChoices(t *testing.T) {
	command := newIntegrationCommand(taskrun.Dependencies{})
	for _, name := range []string{"reconsider", "retry-notification"} {
		child, _, err := command.Find([]string{name})
		if err != nil || child == command || child.Name() != name {
			t.Errorf("no public controlled recovery path for %s: %v", name, err)
		}
	}
	resume, _, err := command.Find([]string{"resume"})
	if err != nil || resume.Flags().Lookup("retry-merge") == nil {
		t.Fatal("known rejected PR cannot be explicitly retried from the public command")
	}
}
