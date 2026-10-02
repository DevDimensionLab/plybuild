package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func TestTaskRunHelpSeparatesInteractiveAndExec(t *testing.T) {
	for _, tc := range []struct {
		leaf string
		want []string
	}{
		{"", []string{"Request schema 1 (interactive)", "human-observed task inactivity", "Request schema 2 (factory exec) uses pipes", "runner-owned completion evidence", "Existing interactive requests and historical results remain compatible."}},
		{"start", []string{"Request schema 1 (interactive) requires text output and a foreground macOS terminal; its child inherits the terminal.", "In interactive mode, without a human task observation", "Factory exec request schema 2 uses pipes", "runner-owned completion evidence", "Unknown completion exits 5"}},
		{"collect", []string{"For request schema 1 (interactive)", "inactive task observation after the report", "For request schema 2 (factory exec)", "ProviderCompletion@1", "rejects --task-status", "Historical TaskResults remain immutable."}},
		{"report", []string{"mandatory task-requirements artifact", "Report received does not prove provider or process completion."}},
	} {
		t.Run(tc.leaf, func(t *testing.T) {
			root := &cobra.Command{Use: "test"}
			root.AddCommand(newTaskRunCommand(taskrun.Dependencies{}))
			var out, stderr bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&stderr)
			args := []string{"run"}
			if tc.leaf != "" {
				args = append(args, tc.leaf)
			}
			root.SetArgs(append(args, "--help"))
			if err := root.Execute(); err != nil || stderr.Len() != 0 {
				t.Fatalf("help failed: %v stderr=%q", err, stderr.String())
			}
			for _, want := range tc.want {
				if !strings.Contains(out.String(), want) {
					t.Errorf("help missing mode-specific contract %q", want)
				}
			}
		})
	}
}
