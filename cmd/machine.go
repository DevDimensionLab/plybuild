package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/machine"
	"github.com/spf13/cobra"
)

func newMachineCommand(read func() (machine.Status, error), processes func() ([]machine.Process, error), labelsPath func() (string, error)) *cobra.Command {
	group := &cobra.Command{Use: "machine", Short: "Read how hard this machine is working", Args: cobra.NoArgs,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil }}
	format, withProcesses := "text", false
	status := &cobra.Command{Use: "status", Short: "Show load per logical core and memory pressure without reading a workspace", Args: cobra.NoArgs,
		Long: "Show the load average for 1, 5 and 15 minutes, as is and as load_percent (load / logical cores x 100, can exceed 100), and macOS memory pressure (normal, warn, critical, or null with a reason). load_percent is queued work per core, not CPU utilisation. With --processes, also the ten processes that use most CPU: pid, program name only (no arguments or path), cpu_percent as ps reports it (100 is one full core; not the same measure as load_percent), rss_bytes, and display_name and category from ply machine process-label (null without a label); processes is null with processes_reason when ps cannot be read, and process_labels_reason says why labels are missing when the label table cannot be read. Every call reads the system; nothing is stored. macOS only.",
		RunE: func(c *cobra.Command, _ []string) error {
			if err := checkFormat(format); err != nil {
				return err
			}
			s, err := read()
			if err != nil {
				return err
			}
			var result interface{ Text() string } = s
			if withProcesses {
				list, err := processes()
				labels, labelsErr := readLabels(labelsPath)
				result = machine.WithLabels(machine.WithProcesses(s, list, err), labels, labelsErr)
			}
			return writeMachine(c.OutOrStdout(), format, result, result.Text())
		}}
	status.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	status.Flags().BoolVar(&withProcesses, "processes", false, "also show the ten processes that use most CPU")
	group.AddCommand(status, newProcessLabelCommand(labelsPath))
	return group
}

func newProcessLabelCommand(labelsPath func() (string, error)) *cobra.Command {
	group := &cobra.Command{Use: "process-label", Short: "Give a process name a readable name and category on this machine", Args: cobra.NoArgs,
		Long: "Keep a table, for this machine and user in ~/" + machine.LabelsFile + ", of readable names for the process names ply machine status --processes shows. A name must match the process name exactly."}
	format, displayName, category := "text", "", ""
	set := &cobra.Command{Use: "set <process name> --display-name <name> [--category <category>]", Short: "Add or replace the label for a process name", Args: cobra.ExactArgs(1),
		Long: "Add or replace the label for a process name exactly as ply machine status --processes shows it; quote a name with spaces. --display-name is one line of up to 60 characters. --category is optional, one of " + strings.Join(machine.Categories, ", ") + ".",
		RunE: func(c *cobra.Command, args []string) error {
			path, err := checkedLabelsPath(format, labelsPath)
			if err != nil {
				return err
			}
			l := machine.Label{Name: args[0], DisplayName: displayName}
			if c.Flags().Changed("category") {
				l.Category = &category
			}
			if l, err = machine.SetLabel(path, l); err != nil {
				return err
			}
			return writeMachine(c.OutOrStdout(), format, l, labelText(l))
		}}
	set.Flags().StringVar(&displayName, "display-name", "", "readable name, one line of up to 60 characters")
	set.Flags().StringVar(&category, "category", "", "one of "+strings.Join(machine.Categories, ", "))
	remove := &cobra.Command{Use: "remove <process name>", Short: "Remove the label for a process name", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			path, err := checkedLabelsPath(format, labelsPath)
			if err != nil {
				return err
			}
			removed, err := machine.RemoveLabel(path, args[0])
			if err != nil {
				return err
			}
			text := "Removed the label for " + args[0]
			if !removed {
				text = "No label for " + args[0]
			}
			return writeMachine(c.OutOrStdout(), format, struct {
				Name    string `json:"name"`
				Removed bool   `json:"removed"`
			}{args[0], removed}, text)
		}}
	list := &cobra.Command{Use: "list", Short: "List the labels on this machine", Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			path, err := checkedLabelsPath(format, labelsPath)
			if err != nil {
				return err
			}
			labels, err := machine.ListLabels(path)
			if err != nil {
				return err
			}
			lines := make([]string, len(labels))
			for i, l := range labels {
				lines[i] = labelText(l)
			}
			text := strings.Join(lines, "\n")
			if len(labels) == 0 {
				text = "No process labels"
			}
			return writeMachine(c.OutOrStdout(), format, struct {
				Kind          string          `json:"kind"`
				SchemaVersion int             `json:"schema_version"`
				Labels        []machine.Label `json:"labels"`
			}{machine.LabelsKind, 1, labels}, text)
		}}
	for _, c := range []*cobra.Command{set, remove, list} {
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	}
	group.AddCommand(set, remove, list)
	return group
}

func checkFormat(format string) error {
	if format != "text" && format != "json" {
		return fmt.Errorf("format must be text or json")
	}
	return nil
}

func checkedLabelsPath(format string, labelsPath func() (string, error)) (string, error) {
	if err := checkFormat(format); err != nil {
		return "", err
	}
	path, err := labelsPath()
	if err != nil {
		return "", fmt.Errorf("%s: home directory for %s: %v", machine.ReasonReadFailed, machine.LabelsFile, err)
	}
	return path, nil
}

func readLabels(labelsPath func() (string, error)) (map[string]machine.Label, error) {
	path, err := labelsPath()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", machine.ErrLabelsRead, err)
	}
	return machine.ReadLabels(path)
}

func labelText(l machine.Label) string {
	if l.Category == nil {
		return l.Name + ": " + l.DisplayName
	}
	return l.Name + ": " + l.DisplayName + " (" + *l.Category + ")"
}

func writeMachine(w io.Writer, format string, doc any, text string) error {
	if format == "json" {
		return json.NewEncoder(w).Encode(doc)
	}
	_, err := fmt.Fprintln(w, text)
	return err
}

func init() {
	system := machine.System()
	RootCmd.AddCommand(newMachineCommand(
		func() (machine.Status, error) { return machine.Read(system, time.Now()) },
		system.(machine.ProcessReader).Processes,
		func() (string, error) {
			home, err := os.UserHomeDir()
			return filepath.Join(home, filepath.FromSlash(machine.LabelsFile)), err
		}))
}
