package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type al14WarningWorkspace struct {
	workflowhandoff.WorkspaceObserver
	err error
}

func (w al14WarningWorkspace) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	logrus.Warn("retirement observation warning")
	return workflowhandoff.WorkspaceSnapshot{}, w.err
}

// The native closed-preparation acceptance reproduces a real successful
// retirement whose warning corrupted JSON. This narrow command regression
// isolates output routing before any workspace mutation or process inspection.
func TestAL14ReportRetirementWarningsStayOnStderr(t *testing.T) {
	for _, surface := range []string{"run", "execute"} {
		t.Run(surface, func(t *testing.T) {
			observed := errors.New("modeled observation ends before any effect")
			d := taskrun.Dependencies{Workflow: workflowhandoff.Dependencies{Workspace: al14WarningWorkspace{err: observed}}}
			var stdout, stderr bytes.Buffer
			previous := logrus.StandardLogger().Out
			logrus.SetOutput(&stdout) // Match the ordinary CLI root logger default.
			t.Cleanup(func() { logrus.SetOutput(previous) })
			root := &cobra.Command{Use: "ply", SilenceErrors: true, SilenceUsage: true}
			workflow := &cobra.Command{Use: "workflow"}
			if surface == "run" {
				workflow.AddCommand(newWorkflowRunCommand(d))
			} else {
				workflow.AddCommand(newWorkflowExecuteCommand(d))
			}
			root.AddCommand(workflow)
			root.SetOut(&stdout)
			root.SetErr(&stderr)
			root.SetArgs([]string{"workflow", surface, "retire", "wfr_fixture", "--observer", "modeled coordinator", "--reason", "Keep warnings separate", "--format", "json"})
			err := root.Execute()
			if !errors.Is(err, observed) {
				t.Fatalf("unexpected observation result: %v", err)
			}
			if stdout.Len() != 0 || strings.Count(stderr.String(), "retirement observation warning") != 1 {
				t.Fatalf("warning must appear once on stderr; stdout=%s stderr=%s", &stdout, &stderr)
			}
			if logrus.StandardLogger().Out != &stdout {
				t.Fatal("command did not restore the prior logger output")
			}
		})
	}
}
