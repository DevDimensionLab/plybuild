package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestWorkflowExecuteInstallCLIHelpRetainedSource(t *testing.T) {
	c := newWorkflowExecuteInstallCLICommand(taskrun.Dependencies{})
	if flag := c.Flags().Lookup("allow-retained-source"); flag == nil || flag.DefValue != "false" {
		t.Fatal("retained-source must be an explicit opt-in")
	}
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs([]string{"--help"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"--allow-retained-source", "default still requires completed closeout", "never inspects writers, removes source or activates a registry", "source_cleanup_pending", "--check", "--observe", "--rollback", "does not widen the owner's installation authority"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("TD125516B install help missing %q", want)
		}
	}
}

func TestWorkflowExecuteInstallCLIActionsRemainExclusive(t *testing.T) {
	for _, args := range [][]string{{"--check", "--observe"}, {"--check", "--rollback"}, {"--observe", "--rollback"}} {
		c := newWorkflowExecuteInstallCLICommand(taskrun.Dependencies{})
		c.SilenceUsage = true
		c.SilenceErrors = true
		c.SetArgs(append([]string{"synthetic-run", "--allow-retained-source"}, args...))
		if err := c.Execute(); err == nil || !strings.Contains(err.Error(), "none of the others can be") {
			t.Fatalf("action selection: %v", err)
		}
	}
}
