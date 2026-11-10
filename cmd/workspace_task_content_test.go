package cmd

import (
	"bytes"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io"
	"strings"
	"testing"
)

func TestWorkspaceTaskContentEnglishNavigationAndRequiredArguments(t *testing.T) {
	for _, args := range [][]string{{"problem"}, {"problem", "show"}, {"problem", "record"}, {"spec"}, {"spec", "list"}, {"spec", "show"}, {"spec", "record"}, {"spec", "assess"}, {"spec", "select"}, {"spec", "withdraw"}, {"publication"}, {"publication", "show"}} {
		c := newWorkspaceTaskCommand(workspace.SystemDependencies())
		var out, stderr bytes.Buffer
		c.SetOut(&out)
		c.SetErr(&stderr)
		c.SetArgs(append(append([]string{}, args...), "--help"))
		if err := c.Execute(); err != nil {
			t.Fatal(err)
		}
		if stderr.Len() != 0 || !strings.Contains(out.String(), "Usage:") {
			t.Fatalf("help %v: %s %s", args, out.String(), stderr.String())
		}
	}
	for _, args := range [][]string{{"spec", "record", "task"}, {"spec", "show", "task"}, {"publication", "show", "task"}} {
		c := newWorkspaceTaskCommand(workspace.SystemDependencies())
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})
		c.SetArgs(args)
		if err := c.Execute(); err == nil {
			t.Fatalf("accepted missing required flags: %v", args)
		}
	}
}

type taskContentFailWriter struct{}

func (taskContentFailWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestTaskCreateJSONLegacyRetryAndPublicationStreamFailure(t *testing.T) {
	args := []string{"create", "task", "--title", "Task", "--description", "Description", "--epic", "epic", "--project", "ply", "--repo", "ply", "--format", "json"}
	services := workspaceTaskServices{create: func(workspace.TaskCreateInput) (workspace.TaskMutationResult, error) {
		return workspace.TaskMutationResult{Workspace: "/fixture", Task: workspace.TaskRecord{ID: "task", Title: "Task"}}, nil
	}, show: func(workspace.TaskID) (workspace.TaskReadbackResult, error) {
		return workspace.TaskReadbackResult{Workspace: "/fixture", Task: workspace.TaskRecord{ID: "task", Title: "Task"}, Reasons: []string{}}, nil
	}}
	out, err := executeTaskCommand(t, services, args...)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = canonicaljson.DecodeStrict([]byte(out)); err != nil {
		t.Fatalf("legacy JSON retry emitted text: %q", out)
	}
	services.create = func(workspace.TaskCreateInput) (workspace.TaskMutationResult, error) {
		return workspace.TaskMutationResult{Workspace: "/fixture", Task: workspace.TaskRecord{ID: "task"}, Created: true, Content: &workspace.TaskContentMutationResult{Workspace: "/fixture", TaskID: "task", PublicationKey: "task-create/task", Classification: "published"}}, nil
	}
	c := newWorkspaceTaskCommandWithServices(services)
	c.SetOut(taskContentFailWriter{})
	c.SetErr(&bytes.Buffer{})
	c.SetArgs(args)
	err = c.Execute()
	if err == nil || !strings.Contains(err.Error(), "ply workspace task publication show task --key task-create/task") {
		t.Fatalf("lost publication recovery after stdout failure: %v", err)
	}
}
