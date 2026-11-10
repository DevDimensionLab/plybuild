package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/assignment"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

const assignmentDefinitions = `{"format_version": 1, "projects": [
  {"id": "a", "name": "A", "wrapper": "a", "repos": [{"id": "a-app", "path": "a/app"}, {"id": "a-lib", "path": "a/lib"}], "epics": []},
  {"id": "b", "name": "B", "wrapper": "b", "repos": [{"id": "b-app", "path": "b/app"}], "epics": []}
]}
`

// newAssignmentWorkspace is a workspace in its live form: .ply/ and
// ply-workspace.json, without the legacy project registry.
func newAssignmentWorkspace(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if _, err := workspace.Init(workspace.SystemDependencies()); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"a/app", "a/lib", "b/app"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(assignmentDefinitions), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func runAssignment(t *testing.T, cwd, stdin string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(cwd)
	command := newAssignmentCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetIn(strings.NewReader(stdin))
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func requireAssignmentError(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil || err.Error() != want {
		t.Fatalf("error %v, want %q", err, want)
	}
}

func decodeAssignmentJSON(t *testing.T, out string, into any) map[string]json.RawMessage {
	t.Helper()
	var keys map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &keys); err != nil {
		t.Fatalf("json %q: %v", out, err)
	}
	if err := json.Unmarshal([]byte(out), into); err != nil {
		t.Fatal(err)
	}
	return keys
}

type assignmentEnvelope struct {
	Project    string                `json:"project"`
	Assignment assignment.Assignment `json:"assignment"`
}

type assignmentListEnvelope struct {
	Project string `json:"project"`
	assignment.ListResult
}

func TestAssignmentAddShowListJourneyInJSON(t *testing.T) {
	root := newAssignmentWorkspace(t)
	app := filepath.Join(root, "a", "app")
	text := "Vis repo-tags.\n<b>&</b>  \r\nBakgrunn: Task x.\n\n"
	out, err := runAssignment(t, app, text, "add", "--title", "Vis repo-tags", "--repo", "a-app", "--repo", " a-lib ", "--topic", "backlog", "--topic", "Web app", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var added assignmentEnvelope
	if keys := decodeAssignmentJSON(t, out, &added); len(keys) != 2 || keys["assignment"] == nil {
		t.Fatalf("add keys %v", keys)
	}
	want := assignment.Assignment{ID: "ASSIGNMENT-1", Title: "Vis repo-tags", Text: text, CreatedAt: added.Assignment.CreatedAt, Repos: []string{"a-app", "a-lib"}, Topics: []string{"backlog", "Web app"}}
	if added.Project != "a" || !reflect.DeepEqual(added.Assignment, want) {
		t.Fatalf("add %+v", added)
	}

	out, err = runAssignment(t, app, "", "show", "ASSIGNMENT-1", "--format", "json")
	var shown assignmentEnvelope
	if err != nil || decodeAssignmentJSON(t, out, &shown) == nil || shown.Project != "a" || !reflect.DeepEqual(shown.Assignment, want) {
		t.Fatalf("show %q %v", out, err)
	}

	out, err = runAssignment(t, root, "", "list", "--project", "a", "--format", "json")
	var listed assignmentListEnvelope
	keys := decodeAssignmentJSON(t, out, &listed)
	if err != nil || len(keys) != 4 || string(keys["view"]) != `"available"` || string(keys["invalid"]) != "[]" || listed.Project != "a" || len(listed.Items) != 1 {
		t.Fatalf("list %q %v", out, err)
	}
	if item := listed.Items[0]; item.ID != "ASSIGNMENT-1" || !reflect.DeepEqual(item.Topics, want.Topics) || !reflect.DeepEqual(item.Repos, want.Repos) {
		t.Fatalf("summary %+v", item)
	}
	if strings.Contains(string(keys["items"]), `"text"`) {
		t.Fatalf("list carries text: %s", keys["items"])
	}
}

func TestAssignmentTextOutput(t *testing.T) {
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	text := "Linje én.\n\nLinje to."
	if out, err := runAssignment(t, a, text, "add", "--title", "Tittel", "--repo", "a-app", "--topic", "backlog"); err != nil || out != "ASSIGNMENT-1\n" {
		t.Fatalf("add %q %v", out, err)
	}
	out, err := runAssignment(t, a, "", "show", "ASSIGNMENT-1")
	if err != nil || !strings.HasPrefix(out, "ASSIGNMENT-1 Tittel\n") || !strings.HasSuffix(out, "\n\n"+text) || !strings.Contains(out, "a-app") || !strings.Contains(out, "backlog") {
		t.Fatalf("show %q %v", out, err)
	}
	out, err = runAssignment(t, a, "", "list")
	if err != nil || strings.Count(out, "\n") != 1 || !strings.HasPrefix(out, "ASSIGNMENT-1 ") || !strings.HasSuffix(out, " a-app Tittel\n") {
		t.Fatalf("list %q %v", out, err)
	}
}

func TestAssignmentAddRequiresProjectRepos(t *testing.T) {
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	_, err := runAssignment(t, a, "x", "add", "--title", "T", "--repo", "a-app", "--repo", "b-app")
	requireAssignmentError(t, err, "assignment_invalid_arguments: --repo must be a repo ID of project a")
	_, err = runAssignment(t, a, "x", "add", "--title", "T")
	requireAssignmentError(t, err, "assignment_invalid_arguments: --repo is required")
	_, err = runAssignment(t, a, "x", "add", "--repo", "a-app")
	requireAssignmentError(t, err, "assignment_invalid_arguments: --title is required")
	if _, err := os.Lstat(filepath.Join(a, assignment.DataPathCategory)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("assignments directory created: %v", err)
	}
}

func TestAssignmentProjectsHaveSeparateNumbering(t *testing.T) {
	root := newAssignmentWorkspace(t)
	for _, project := range []string{"a", "b"} {
		out, err := runAssignment(t, root, "x", "add", "--project", project, "--title", "Title "+project, "--repo", project+"-app")
		if err != nil || out != "ASSIGNMENT-1\n" {
			t.Fatalf("add %s %q %v", project, out, err)
		}
	}
	out, err := runAssignment(t, filepath.Join(root, "a", "app"), "", "show", "--project", "b", "ASSIGNMENT-1", "--format", "json")
	var shown assignmentEnvelope
	if err != nil || decodeAssignmentJSON(t, out, &shown) == nil || shown.Project != "b" || shown.Assignment.Title != "Title b" {
		t.Fatalf("show b %q %v", out, err)
	}
	for _, project := range []string{"a", "b"} {
		if _, err := os.Stat(filepath.Join(root, project, "assignments", "ASSIGNMENT-1.json")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAssignmentErrorForms(t *testing.T) {
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	for name, tc := range map[string]struct {
		cwd  string
		args []string
		want string
	}{
		"unknown project":   {root, []string{"list", "--project", "nope"}, "workspace_project_not_found: project nope is not registered"},
		"invalid project":   {root, []string{"list", "--project", "Nope/../x"}, "workspace_project_invalid_arguments: could not resolve project Nope/../x"},
		"no project at cwd": {root, []string{"list"}, "workspace_project_invalid_arguments: project is required: 0 projects contain the current directory"},
		"unknown flag":      {a, []string{"list", "--secret=/abs/path"}, "assignment_invalid_arguments: invalid flag"},
		"bad format":        {a, []string{"list", "--format", "yaml"}, "assignment_invalid_arguments: --format must be text or json"},
		"show without id":   {a, []string{"show"}, "assignment_invalid_arguments: expected 1 positional arguments, got 0"},
		"list with args":    {a, []string{"list", "x"}, "assignment_invalid_arguments: expected 0 positional arguments, got 1"},
		"unknown command":   {a, []string{"edit"}, "assignment_invalid_arguments: unknown assignment command"},
		"missing":           {a, []string{"show", "ASSIGNMENT-4"}, "assignment_not_found: ASSIGNMENT-4"},
		"invalid id":        {a, []string{"show", "ASSIGNMENT-04"}, "assignment_invalid_arguments: invalid assignment id"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := runAssignment(t, tc.cwd, "", tc.args...)
			requireAssignmentError(t, err, tc.want)
		})
	}
	outside := t.TempDir()
	_, err := runAssignment(t, outside, "", "list", "--project", "a")
	requireAssignmentError(t, err, "workspace_not_found: could not resolve project a")
}

func TestAssignmentListReportsInvalidFilesWithExitZero(t *testing.T) {
	root := newAssignmentWorkspace(t)
	a := filepath.Join(root, "a")
	if _, err := runAssignment(t, a, "x", "add", "--title", "T", "--repo", "a-app"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(a, "assignments", "ASSIGNMENT-2.json"), []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runAssignment(t, a, "", "list", "--format", "json")
	var listed assignmentListEnvelope
	if err != nil || decodeAssignmentJSON(t, out, &listed) == nil || len(listed.Items) != 1 ||
		!reflect.DeepEqual(listed.Invalid, []assignment.InvalidFile{{File: "ASSIGNMENT-2.json", Error: "invalid JSON"}}) {
		t.Fatalf("list %q %v", out, err)
	}
	out, err = runAssignment(t, a, "", "list")
	if err != nil || !strings.HasSuffix(out, "invalid: ASSIGNMENT-2.json: invalid JSON\n") {
		t.Fatalf("list text %q %v", out, err)
	}
	_, err = runAssignment(t, a, "", "show", "ASSIGNMENT-2")
	requireAssignmentError(t, err, "assignment_invalid: invalid JSON")
	if _, err := runAssignment(t, a, "x", "add", "--title", "T", "--repo", "a-app"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(a, "assignments", "ASSIGNMENT-3.json")); err != nil {
		t.Fatal(err)
	}
}

func TestAssignmentCommandIsRegistered(t *testing.T) {
	if command, _, err := RootCmd.Find([]string{"assignment", "show"}); err != nil || command.Name() != "show" {
		t.Fatalf("assignment show not registered: %v", err)
	}
}
