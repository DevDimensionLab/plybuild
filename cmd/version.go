package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	format := "text"
	c := &cobra.Command{Use: "version", Short: "Identify this executable without reading a workspace", Args: cobra.NoArgs,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		RunE: func(c *cobra.Command, _ []string) error {
			if format != "text" && format != "json" {
				return fmt.Errorf("format must be text or json")
			}
			i, err := buildinfo.Current()
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(c.OutOrStdout()).Encode(i)
			}
			_, err = fmt.Fprintf(c.OutOrStdout(), "Ply %s (%s, dirty=%t) %s/%s\n%s\n", i.Version, i.SourceCommit, i.Dirty, i.OS, i.Arch, i.SHA256)
			return err
		}}
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	return c
}

func init() { RootCmd.AddCommand(newVersionCommand()) }
