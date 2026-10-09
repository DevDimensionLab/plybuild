package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowWorkerCommand(d taskrun.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "worker", Short: "Observe and control a durable one-Spec headless worker", Long: "Start with workflow execute --headless. Worker show is read-only. Pause and stop request a preserved boundary after the active attempt ends; resume continues the same Task and native Run. A live attempt retains the owner lock when its controller exits. Unknown effects remain blocked."}
	for _, verb := range []string{"show", "resume", "pause", "stop", "_supervise", "_serve", "_attempt"} {
		var root, format, job string
		v := verb
		c := &cobra.Command{Use: v + " <run-id>", Args: cobra.ExactArgs(1), Hidden: v[0] == '_', RunE: func(c *cobra.Command, args []string) error {
			if err := executeFormat(format); err != nil {
				return err
			}
			if root == "" {
				o, err := workspace.ObserveContaining(d.Workspace)
				if err != nil {
					return err
				}
				root = o.Root
			}
			var run taskrun.WorkflowRun
			var err error
			switch v {
			case "show":
				run, err = taskrun.ShowWorker(d, root, args[0])
			case "resume":
				run, err = taskrun.StartWorker(d, root, args[0], true)
			case "pause", "stop":
				run, err = taskrun.ControlWorker(d, root, args[0], v)
			case "_serve":
				err = taskrun.ServeWorker(d, root, args[0])
				if err == nil {
					run, err = taskrun.ShowWorker(d, root, args[0])
				}
			case "_supervise":
				return taskrun.SuperviseWorker(d, root, args[0])
			case "_attempt":
				return taskrun.RunWorkerAttempt(d, root, args[0], job)
			}
			if e := writeExecuteResult(c, format, run); e != nil {
				return e
			}
			return err
		}}
		c.Flags().StringVar(&root, "workspace", "", "physical workspace root (default: containing workspace)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if v == "_attempt" {
			c.Flags().StringVar(&job, "job", "", "immutable native attempt ID")
		}
		setWorkFlagErrors(c)
		parent.AddCommand(c)
	}
	return parent
}
