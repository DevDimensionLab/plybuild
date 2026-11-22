package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io"
	"os"
	"strings"
	"testing"
)

func TestWorkspaceTaskContentEnglishNavigationAndRequiredArguments(t *testing.T) {
	for _, args := range [][]string{{"problem"}, {"problem", "show"}, {"problem", "record"}, {"spec"}, {"spec", "list"}, {"spec", "show"}, {"spec", "record"}, {"spec", "assess"}, {"spec", "select"}, {"spec", "withdraw"}, {"publication"}, {"publication", "show"}} {
		c := newWorkspaceTaskCommand(taskwork.SystemDependencies())
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
		c := newWorkspaceTaskCommand(taskwork.SystemDependencies())
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
	services := workspaceTaskServices{create: func(taskwork.TaskCreateInput) (taskwork.TaskMutationResult, error) {
		return taskwork.TaskMutationResult{Workspace: "/fixture", Task: taskwork.TaskRecord{ID: "task", Title: "Task"}}, nil
	}, show: func(taskwork.TaskID) (taskwork.TaskReadbackResult, error) {
		return taskwork.TaskReadbackResult{Workspace: "/fixture", Task: taskwork.TaskRecord{ID: "task", Title: "Task"}, Reasons: []string{}}, nil
	}}
	out, err := executeTaskCommand(t, services, args...)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = canonicaljson.DecodeStrict([]byte(out)); err != nil {
		t.Fatalf("legacy JSON retry emitted text: %q", out)
	}
	services.create = func(taskwork.TaskCreateInput) (taskwork.TaskMutationResult, error) {
		return taskwork.TaskMutationResult{Workspace: "/fixture", Task: taskwork.TaskRecord{ID: "task"}, Created: true, Content: &taskwork.TaskContentMutationResult{Workspace: "/fixture", TaskID: "task", PublicationKey: "task-create/task", Classification: "published"}}, nil
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

func TestWorkspaceTaskSpecShowDocumentWritesBytesNotFormat(t *testing.T) {
	c := newWorkspaceTaskCommand(taskwork.SystemDependencies())
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SetArgs([]string{"spec", "show", "--help"})
	if err := c.Execute(); err != nil || !strings.Contains(out.String(), "--document") {
		t.Fatalf("spec show lacks --document: %v %s", err, out.String())
	}
	c = newWorkspaceTaskCommand(taskwork.SystemDependencies())
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	c.SetArgs([]string{"spec", "show", "task", "--spec", "s", "--revision", "1", "--document", "spec", "--format", "json"})
	if err := c.Execute(); err == nil || !strings.Contains(err.Error(), "--format does not apply") {
		t.Fatalf("document with format: %v", err)
	}
	c = newWorkspaceTaskCommand(taskwork.SystemDependencies())
	c.SetOut(&bytes.Buffer{})
	c.SetErr(&bytes.Buffer{})
	c.SetArgs([]string{"spec", "show", "task", "--spec", "s", "--revision", "1", "--document", ""})
	if err := c.Execute(); err == nil || !strings.Contains(err.Error(), "--document") {
		t.Fatalf("empty document: %v", err)
	}
}

// problem show --document writes a problem document's preserved bytes the way
// spec show --document writes a Spec document's.
func TestWorkspaceTaskProblemShowDocumentWritesPreservedBytes(t *testing.T) {
	d, _ := problemTaskFixture(t)
	help, err := runWorkspaceView(t, d, "task", "problem", "show", "--help")
	if err != nil || !strings.Contains(help, "--document") {
		t.Fatalf("problem show lacks --document: %v %s", err, help)
	}
	out, err := runWorkspaceView(t, d, "task", "problem", "show", "task", "--revision", "1", "--document", "problem")
	if err != nil || out != "# Explicit fixture problem\n" {
		t.Fatalf("problem document bytes: %q %v", out, err)
	}
	notFound := func(args ...string) {
		t.Helper()
		out, err := runWorkspaceView(t, d, append([]string{"task", "problem", "show", "task"}, args...)...)
		var classified *workspace.WorkItemError
		if out != "" || !errors.As(err, &classified) || classified.Class != workspace.ErrorWorkNotFound {
			t.Fatalf("%v: %q %v", args, out, err)
		}
	}
	notFound("--revision", "1", "--document", "absent")
	notFound("--revision", "2", "--document", "problem")
	for _, args := range [][]string{{"--document", "problem"}, {"--revision", "1", "--document", ""}, {"--revision", "1", "--document", "problem", "--format", "json"}} {
		if out, err := runWorkspaceView(t, d, append([]string{"task", "problem", "show", "task"}, args...)...); err == nil || out != "" {
			t.Fatalf("accepted %v: %q", args, out)
		}
	}
	readback, err := runWorkspaceView(t, d, "task", "problem", "show", "task", "--revision", "1", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var shown struct {
		Revision struct {
			Documents []struct {
				Locator string `json:"locator"`
			} `json:"documents"`
		} `json:"revision"`
	}
	if err = json.Unmarshal([]byte(readback), &shown); err != nil || len(shown.Revision.Documents) != 1 {
		t.Fatalf("%v %s", err, readback)
	}
	stored := shown.Revision.Documents[0].Locator
	if err = os.Chmod(stored, 0600); err != nil {
		t.Fatal(err)
	}
	// Same size and immutable mode, so only the hash tells the bytes apart.
	if err = os.WriteFile(stored, []byte("# EXPLICIT FIXTURE PROBLEM\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(stored, 0400); err != nil {
		t.Fatal(err)
	}
	out, err = runWorkspaceView(t, d, "task", "problem", "show", "task", "--revision", "1", "--document", "problem")
	if out != "" || err == nil || !strings.Contains(err.Error(), "managed bytes do not match their digest") {
		t.Fatalf("wrote bytes that differ from their hash: %q %v", out, err)
	}
}
