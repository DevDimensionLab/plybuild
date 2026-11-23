package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devdimensionlab/plybuild/internal/machine"
	"github.com/spf13/cobra"
)

func newMachineCommand(read func() (machine.Status, error), processes func() ([]machine.Process, error)) *cobra.Command {
	group := &cobra.Command{Use: "machine", Short: "Read how hard this machine is working", Args: cobra.NoArgs,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil }}
	format, withProcesses := "text", false
	status := &cobra.Command{Use: "status", Short: "Show load per logical core and memory pressure without reading a workspace", Args: cobra.NoArgs,
		Long: "Show the load average for 1, 5 and 15 minutes, as is and as load_percent (load / logical cores x 100, can exceed 100), and macOS memory pressure (normal, warn, critical, or null with a reason). load_percent is queued work per core, not CPU utilisation. With --processes, also the ten processes that use most CPU: pid, program name only (no arguments or path), cpu_percent as ps reports it (100 is one full core; not the same measure as load_percent) and rss_bytes; processes is null with processes_reason when ps cannot be read. Every call reads the system; nothing is stored. macOS only.",
		RunE: func(c *cobra.Command, _ []string) error {
			if format != "text" && format != "json" {
				return fmt.Errorf("format must be text or json")
			}
			s, err := read()
			if err != nil {
				return err
			}
			var result interface{ Text() string } = s
			if withProcesses {
				list, err := processes()
				result = machine.WithProcesses(s, list, err)
			}
			if format == "json" {
				return json.NewEncoder(c.OutOrStdout()).Encode(result)
			}
			_, err = fmt.Fprintln(c.OutOrStdout(), result.Text())
			return err
		}}
	status.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	status.Flags().BoolVar(&withProcesses, "processes", false, "also show the ten processes that use most CPU")
	group.AddCommand(status)
	return group
}

func init() {
	system := machine.System()
	RootCmd.AddCommand(newMachineCommand(
		func() (machine.Status, error) { return machine.Read(system, time.Now()) },
		system.(machine.ProcessReader).Processes))
}
