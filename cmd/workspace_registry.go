package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceRegistryCommand(d workspace.Dependencies) *cobra.Command {
	c := &cobra.Command{Use: "registry", Short: "Inspect and explicitly upgrade the shared work-item format"}
	var directory, pin, confirmation, format string
	var apply, check bool
	var to int
	u := &cobra.Command{Use: "upgrade", Short: "Confirm a format boundary against the exact registry and delivered reader", Long: "Preview a format-4 to format-5 transition without writes. Apply requires the exact fresh confirmation and a verified prebuilt release, rechecked under the write lock. Preserve all records and a byte-for-byte backup. Compatible ordinary writes keep their format. Retry the same package and confirmation after interruption; an observed transition is recovered without another write. Upgrade only after all workspace consumers have a compatible delivered reader.", Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if apply && check {
				return fmt.Errorf("choose --check or --apply")
			}
			if format != "text" && format != "json" {
				return fmt.Errorf("format must be text or json")
			}
			upgrade := workspace.UpgradeRegistry
			if to == 6 {
				upgrade = workspace.UpgradeEpicCreateRegistry
			} else if to == 7 {
				upgrade = workspace.UpgradeAdoptionRegistry
			} else if to != 5 {
				return fmt.Errorf("to-format must be 5, 6 or 7")
			}
			p, err := upgrade(d, directory, pin, confirmation, apply)
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(c.OutOrStdout()).Encode(p)
			}
			_, err = fmt.Fprintf(c.OutOrStdout(), "Registry %s: format %d -> %d\nBackup: %s\nConfirmation: %s\n", p.State, p.FromVersion, p.ToVersion, p.Backup, p.Confirmation)
			return err
		}}
	u.Flags().IntVar(&to, "to-format", 5, "explicit target format (5, 6 for Epic creation, or 7 for Task adoption)")
	u.Flags().StringVar(&directory, "package", "", "absolute delivered release directory")
	u.Flags().StringVar(&pin, "sha256", "", "expected sha256: manifest digest")
	u.Flags().StringVar(&confirmation, "confirm", "", "exact fresh preview confirmation")
	u.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	u.Flags().BoolVar(&apply, "apply", false, "publish the confirmed format transition")
	u.Flags().BoolVar(&check, "check", false, "read-only preview (default)")
	_ = u.MarkFlagRequired("package")
	_ = u.MarkFlagRequired("sha256")
	c.AddCommand(u)
	var qualified, coordinator, token, output string
	var activate, preview bool
	a := &cobra.Command{Use: "activate-parallel", Short: "Explicitly activate format8 using an exact completed qualified native CLI", Long: "Read-only --check is the default. --apply requires the exact confirmation over the actual format4-7 predecessor, registered target, qualified candidate bytes and finite native remaining-effect compatibility. Requires a completed self-hosted Ply Run with complete native closeout; this never installs a CLI, changes foreign controls, clears stops or probes processes. Incompatible or unknown remaining readers block activation. Legacy upgrades keep their explicit formats5/6/7.", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if activate && preview {
			return fmt.Errorf("choose --check or --apply")
		}
		if output != "json" && output != "text" {
			return fmt.Errorf("format must be text or json")
		}
		observed, err := workspace.ObserveContaining(d)
		if err != nil {
			return err
		}
		result, err := taskrun.ActivateParallelWorkspace(taskrun.SystemDependencies(d), observed.Root, qualified, coordinator, token, activate)
		if output == "json" {
			if e := json.NewEncoder(cmd.OutOrStdout()).Encode(result); e != nil {
				return e
			}
		} else {
			if _, e := fmt.Fprintf(cmd.OutOrStdout(), "Parallel activation: %s\nConfirmation: %s\nBlockers: %v\n", result.State, result.Confirmation, result.Blockers); e != nil {
				return e
			}
		}
		return err
	}}
	a.Flags().StringVar(&qualified, "qualified-run", "", "completed native feature Run supplying the exact tested reader")
	a.Flags().StringVar(&coordinator, "coordinator", "", "actual coordinator invoking the separately authorized activation")
	a.Flags().StringVar(&token, "confirm", "", "fresh exact activation confirmation")
	a.Flags().StringVar(&output, "format", "text", "output format (text or json)")
	a.Flags().BoolVar(&activate, "apply", false, "publish exactly one confirmed format8 boundary")
	a.Flags().BoolVar(&preview, "check", false, "observe compatibility without writes (default)")
	_ = a.MarkFlagRequired("qualified-run")
	_ = a.MarkFlagRequired("coordinator")
	c.AddCommand(a)
	return c
}
