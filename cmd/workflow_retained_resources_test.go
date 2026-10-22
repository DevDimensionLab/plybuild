package cmd

import (
	"bytes"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
	"strings"
	"testing"
)

func TestAL40RetireFlagHelpAndNoManualSelectors(t *testing.T) {
	for _, constructor := range []func(taskrun.Dependencies) *cobra.Command{newWorkflowExecuteCommand, newWorkflowRunCommand} {
		root := constructor(taskrun.Dependencies{})
		c, _, err := root.Find([]string{"retire"})
		if err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		c.SetOut(&out)
		if err = c.Help(); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out.String(), "--retain-effects") {
			t.Fatal("explicit flag missing")
		}
		for _, selector := range []string{"native-session", "stop", "source-sha256", "force"} {
			if c.Flags().Lookup(selector) != nil {
				t.Fatal("added manual selector", selector)
			}
		}
	}
}
