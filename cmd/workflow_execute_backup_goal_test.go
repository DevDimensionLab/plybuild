package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"strings"
	"testing"
)

func TestBackupGoalOptInFlags(t *testing.T) {
	c := newWorkflowExecuteBackupGoalCommand(taskrun.Dependencies{})
	if err := c.ParseFlags([]string{"run", "--task", "fresh", "--spec", "fresh-v1", "--provider", "codex", "--start", "--herdr-workspace", "w-fixture"}); err != nil {
		t.Fatal(err)
	}
	if start, _ := c.Flags().GetBool("start"); !start {
		t.Fatal("missing explicit opt-in")
	}
	defaults := newWorkflowExecuteBackupGoalCommand(taskrun.Dependencies{})
	for _, flag := range []string{"start", "apply"} {
		if value, _ := defaults.Flags().GetBool(flag); value {
			t.Fatalf("%s changed registration-only default", flag)
		}
	}
	for _, flag := range []string{"restart", "resume", "headless"} {
		if c.Flags().Lookup(flag) != nil {
			t.Fatalf("unsupported backup recovery option %s", flag)
		}
	}
	if !strings.Contains(c.Long, "Identical retries") || !strings.Contains(c.Long, "Without --start") {
		t.Fatal("missing startup boundary guidance")
	}
}

func TestSelectedBackupRetirementFlags(t *testing.T) {
	c := newWorkflowExecuteBackupGoalCommand(taskrun.Dependencies{})
	if err := c.ParseFlags([]string{"run", "--task", "fresh", "--spec", "fresh-v1", "--provider", "codex", "--retire", "--coordinator", "modeled coordinator", "--reason", "Selected replacement"}); err != nil {
		t.Fatalf("selected pre-work replacement is unavailable: %v", err)
	}
	if value, _ := newWorkflowExecuteBackupGoalCommand(taskrun.Dependencies{}).Flags().GetBool("retire"); value {
		t.Fatal("retirement must remain opt-in")
	}
}
