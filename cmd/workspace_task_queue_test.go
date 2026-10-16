package cmd

import (
	"bytes"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"strings"
	"testing"
)

func TestTaskQueueInvalidCLIStopsBeforeDependencies(t *testing.T) {
	cases := [][]string{
		{"prepare"}, {"prepare", "task", "--next"}, {"prepare", "--next", "--apply"},
		{"prepare", "--next", "--check", "--apply", "--confirm", "sha256:x"},
		{"prepare", "--next", "--confirm", "sha256:x"}, {"prepare", "--next", "--upgrade-registry", "sha256:x"},
		{"prepare", "--next", "--project", "ply"}, {"prepare", "--next", "--format", "yaml"},
		{"queue", "list", "--repo", "ply"}, {"queue", "list", "surprise"}, {"queue", "set"},
		{"queue", "advance", "--preparation", "pre_x", "--expected-revision", "1"},
		{"list", "--ready", "--epic", "epic"}, {"list", "--project", ""}, {"list", "--ready", "--format", "yaml"},
	}
	for _, args := range cases {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			c := newWorkspaceTaskCommand(workspace.Dependencies{})
			var out, stderr bytes.Buffer
			c.SetOut(&out)
			c.SetErr(&stderr)
			c.SetArgs(args)
			c.SilenceErrors = true
			c.SilenceUsage = true
			if e := c.Execute(); e == nil {
				t.Fatalf("invalid input reached execution: %v", args)
			}
		})
	}
}
func TestQueueReadyAliasPreservesOrdinaryListService(t *testing.T) {
	calls := 0
	c := newWorkspaceTaskCommandWithServices(workspaceTaskServices{list: func(filters workspace.TaskListFilters) (workspace.TaskListResult, error) {
		calls++
		if filters.EpicID == nil || *filters.EpicID != "epic" {
			t.Fatal("legacy Epic filter changed")
		}
		return workspace.TaskListResult{Workspace: "/fixture", Tasks: []workspace.TaskRecord{}}, nil
	}})
	addWorkspaceTaskQueueCommands(c, workspace.Dependencies{})
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	c.SetArgs([]string{"list", "--epic", "epic"})
	if e := c.Execute(); e != nil || calls != 1 {
		t.Fatalf("legacy list: %d %v", calls, e)
	}
}

func TestRemovedSourceReleaseCLIRequiresExplicitPair(t *testing.T) {
	base := []string{"queue", "release", "--preparation", "pre_" + strings.Repeat("a", 64), "--expected-revision", "2", "--reason", "preserve history"}
	for _, flags := range [][]string{
		{"--removed-source"}, {"--removed-source", "--expected-source-oid", "HEAD"},
		{"--expected-source-oid", strings.Repeat("a", 40)}, {"--expected-source-oid", ""},
		{"--removed-source=false", "--expected-source-oid", strings.Repeat("a", 40)},
	} {
		t.Run(strings.Join(flags, " "), func(t *testing.T) {
			c := newWorkspaceTaskCommand(workspace.Dependencies{})
			c.SetArgs(append(append([]string{}, base...), flags...))
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			if e := c.Execute(); e == nil {
				t.Fatal("invalid opt-in reached dependencies")
			}
		})
	}
	c := newWorkspaceTaskCommand(workspace.Dependencies{})
	advance, _, e := c.Find([]string{"queue", "advance"})
	if e != nil || advance.Flags().Lookup("removed-source") != nil || advance.Flags().Lookup("expected-source-oid") != nil {
		t.Fatal("removed-source flags leaked into advance")
	}
}
