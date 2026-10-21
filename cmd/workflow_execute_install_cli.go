package cmd

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteInstallCLICommand(d taskrun.Dependencies) *cobra.Command {
	var destination, coordinator, format string
	var check, observe, rollback bool
	c := &cobra.Command{Use: "install-cli RUN_ID", Short: "Coordinator installation of an exact delivered self-hosted Ply CLI", Long: "Run from the registered local Epic delivery target after native closeout. Uses the retained qualified snapshot, not a rebuild. --destination explicitly selects the Ply executable. This separate coordinator operation does not widen the owner's installation authority or reopen its Task. --check returns the retained controller without writing; --observe probes existing effects; --rollback restores only the guarded backup. Reuse identical inputs after interruption or reply loss.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		action := "install"
		if check {
			action = "check"
		}
		if observe {
			action = "observe"
		}
		if rollback {
			action = "rollback"
		}
		out, effectErr := taskrun.InstallQualifiedCLI(d, ws.Observation.Root, args[0], destination, coordinator, action)
		if format == "json" {
			if err = json.NewEncoder(c.OutOrStdout()).Encode(out); err != nil {
				return err
			}
		} else {
			_, err = c.OutOrStdout().Write([]byte(out.State + ": " + out.NextAction + "\nController: " + out.Controller.Locator + "\n"))
			if err != nil {
				return err
			}
		}
		return effectErr
	}}
	c.Flags().StringVar(&destination, "destination", "", "explicit physical destination of the self-hosted Ply CLI")
	c.Flags().StringVar(&coordinator, "coordinator", "", "actual authorized coordinator attribution; grants no owner permission")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.Flags().BoolVar(&check, "check", false, "read qualification, delivery, closeout and retained controller without writing")
	c.Flags().BoolVar(&observe, "observe", false, "observe the same operation's installed hash and smoke without replacement")
	c.Flags().BoolVar(&rollback, "rollback", false, "restore the exact guarded backup without overwriting a later operation")
	c.MarkFlagsMutuallyExclusive("check", "observe", "rollback")
	setWorkFlagErrors(c)
	return c
}
