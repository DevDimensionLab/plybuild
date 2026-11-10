package cmd

import (
	"encoding/json"
	"fmt"

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
	u.Flags().IntVar(&to, "to-format", 5, "explicit target format (5, 6 for Epic creation, or historical format 7)")
	u.Flags().StringVar(&directory, "package", "", "absolute delivered release directory")
	u.Flags().StringVar(&pin, "sha256", "", "expected sha256: manifest digest")
	u.Flags().StringVar(&confirmation, "confirm", "", "exact fresh preview confirmation")
	u.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	u.Flags().BoolVar(&apply, "apply", false, "publish the confirmed format transition")
	u.Flags().BoolVar(&check, "check", false, "read-only preview (default)")
	_ = u.MarkFlagRequired("package")
	_ = u.MarkFlagRequired("sha256")
	c.AddCommand(u)
	return c
}
