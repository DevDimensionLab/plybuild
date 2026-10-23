package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestParallelActivationRetainedSourceHelpAndActionSelection(t *testing.T) {
	c := newWorkspaceRegistryCommand(workspace.Dependencies{})
	a, _, err := c.Find([]string{"activate-parallel"})
	if err != nil {
		t.Fatal(err)
	}
	if flag := a.Flags().Lookup("allow-retained-source"); flag == nil || flag.DefValue != "false" {
		t.Fatal("TD7671441 retained activation must be explicit")
	}
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs([]string{"activate-parallel", "--help"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"--allow-retained-source", "default requires completed native closeout", "automatic ownership release", "qualified_reader_eligible", "source_cleanup_pending", "state=activated", "installed CLI alone", "Incompatible or unknown remaining readers block", "never installs a CLI, removes source", "--check", "--confirm"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("activation help missing %q", want)
		}
	}
	c = newWorkspaceRegistryCommand(workspace.Dependencies{})
	c.SilenceErrors, c.SilenceUsage = true, true
	c.SetArgs([]string{"activate-parallel", "--qualified-run", "simulated-run", "--coordinator", "simulated coordinator", "--allow-retained-source", "--check", "--apply"})
	if err := c.Execute(); err == nil || !strings.Contains(err.Error(), "choose --check or --apply") {
		t.Fatalf("action selection: %v", err)
	}
}
