package cmd

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestWorkspaceProjectCommandShape(t *testing.T) {
	command := newWorkspaceProjectCommand(workspace.Dependencies{})
	root := &cobra.Command{Use: "ply"}
	workspaceParent := &cobra.Command{Use: "workspace"}
	root.AddCommand(workspaceParent)
	workspaceParent.AddCommand(command)

	if command.Use != "project" || command.Short != "Manage projects in a Ply workspace" ||
		command.Long != "Register and inspect projects owned by the containing Ply workspace." ||
		command.Example != "  ply workspace project add ply --name Ply --wrapper ../ply --repo ply=../ply/main\n  ply workspace project list" ||
		command.CommandPath() != "ply workspace project" || command.Runnable() {
		t.Fatalf("project metadata = %#v", command)
	}

	tests := []struct {
		name    string
		use     string
		short   string
		long    string
		example string
	}{
		{name: "add", use: "add <project-id>", short: "Register a project and its explicit repository members", long: "Register one project, its wrapper, and one or more explicit Git repository members in the containing Ply workspace.", example: "  ply workspace project add ply --name Ply --wrapper ../ply --repo ply=../ply/main"},
		{name: "show", use: "show <project-id>", short: "Show a registered project", long: "Show the stored wrapper and repository members for one project in the containing Ply workspace.", example: "  ply workspace project show ply\n  ply workspace project show --format json -- ply"},
		{name: "list", use: "list", short: "List registered projects", long: "List projects registered in the containing Ply workspace.", example: "  ply workspace project list\n  ply workspace project list --format json"},
	}
	for _, test := range tests {
		child, _, err := command.Find([]string{test.name})
		if err != nil {
			t.Fatal(err)
		}
		if child == command || child.Use != test.use || child.Short != test.short || child.Long != test.long || child.Example != test.example || !child.Runnable() || child.Args == nil {
			t.Fatalf("%s metadata = Use %q Short %q Long %q Example %q runnable %t", test.name, child.Use, child.Short, child.Long, child.Example, child.Runnable())
		}
	}

	add, _, _ := command.Find([]string{"add"})
	add.InitDefaultHelpFlag()
	flags := map[string]*pflag.Flag{}
	add.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) { flags[flag.Name] = flag })
	if len(flags) != 4 || flags["name"].Shorthand != "n" || flags["name"].Usage != "project display name" ||
		flags["wrapper"].Shorthand != "w" || flags["wrapper"].Usage != "project wrapper directory" ||
		flags["repo"].Shorthand != "r" || flags["repo"].Value.Type() != "stringArray" || flags["repo"].Usage != "repository member as <repo-id>=<path>" {
		t.Fatalf("add flags = %#v", flags)
	}
	for _, name := range []string{"name", "wrapper", "repo"} {
		if !reflect.DeepEqual(flags[name].Annotations[cobra.BashCompOneRequiredFlag], []string{"true"}) {
			t.Fatalf("flag %s is not required: %#v", name, flags[name].Annotations)
		}
	}
	for _, name := range []string{"show", "list"} {
		child, _, _ := command.Find([]string{name})
		child.InitDefaultHelpFlag()
		count := 0
		child.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) { count++ })
		if count != 2 || child.Flags().Lookup("format").DefValue != "text" {
			t.Fatalf("%s local flag count = %d", name, count)
		}
	}
}

func TestWorkspaceProjectCommandsCallUseCasesAndRenderSortedOutput(t *testing.T) {
	addCalls := 0
	showCalls := 0
	listCalls := 0
	services := workspaceProjectServices{
		add: func(input workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			addCalls++
			if input.ProjectID != "trip" || input.Name != "Trip" || input.Wrapper != "../trip" || !reflect.DeepEqual(input.Repos, []workspace.RepoInput{{RepoID: "service", Path: "../service"}, {RepoID: "api", Path: "../api,mirror"}}) {
				t.Fatalf("add input = %#v", input)
			}
			return workspace.ProjectResult{
				Workspace: "/workspace",
				Project:   workspace.ProjectRecord{ID: "trip", Name: "Trip", Wrapper: "/trip", RepoIDs: []workspace.RepoID{"api", "service"}},
				Repos: []workspace.RepoRecord{
					{ID: "service", Locator: "/service", GitCommonDir: "/service/.git"},
					{ID: "api", Locator: "/api", GitCommonDir: "/api/.git"},
				},
				Created: true,
			}, nil
		},
		show: func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			showCalls++
			if id != "trip" {
				t.Fatalf("show ID = %s", id)
			}
			return workspace.ProjectResult{Workspace: "/workspace", Project: workspace.ProjectRecord{ID: "trip", Name: "Trip", Wrapper: "/trip", RepoIDs: []workspace.RepoID{"api"}}, Repos: []workspace.RepoRecord{{ID: "api", Locator: "/api", GitCommonDir: "/api/.git"}}}, nil
		},
		list: func() (workspace.ProjectListResult, error) {
			listCalls++
			return workspace.ProjectListResult{Workspace: "/workspace", Projects: []workspace.ProjectRecord{
				{ID: "trip", Name: "Trip", Wrapper: "/trip", RepoIDs: []workspace.RepoID{"api", "service"}},
				{ID: "ply", Name: "Ply", Wrapper: "/ply", RepoIDs: []workspace.RepoID{"ply"}},
			}}, nil
		},
	}

	stdout, stderr, err := executeProjectCommand(t, services, "add", "trip", "--name", "Trip", "--wrapper", "../trip", "--repo", "service=../service", "--repo", "api=../api,mirror")
	if err != nil || stderr != "" {
		t.Fatalf("add error = %v, stderr = %q", err, stderr)
	}
	wantAdd := "Added project trip (Trip) to Ply workspace /workspace.\nWrapper: /trip\nRepositories:\n  api:\n    locator: /api\n    git common directory: /api/.git\n  service:\n    locator: /service\n    git common directory: /service/.git\n"
	if stdout != wantAdd {
		t.Fatalf("add stdout = %q, want %q", stdout, wantAdd)
	}

	stdout, stderr, err = executeProjectCommand(t, services, "show", "trip")
	if err != nil || stderr != "" {
		t.Fatalf("show error = %v, stderr = %q", err, stderr)
	}
	wantShow := "Project trip (Trip) in Ply workspace /workspace.\nWrapper: /trip\nRepositories:\n  api:\n    locator: /api\n    git common directory: /api/.git\n"
	if stdout != wantShow {
		t.Fatalf("show stdout = %q, want %q", stdout, wantShow)
	}

	stdout, stderr, err = executeProjectCommand(t, services, "list")
	if err != nil || stderr != "" {
		t.Fatalf("list error = %v, stderr = %q", err, stderr)
	}
	wantList := "Projects in Ply workspace /workspace:\n  ply: Ply (1 repository) /ply\n  trip: Trip (2 repositories) /trip\n"
	if stdout != wantList {
		t.Fatalf("list stdout = %q, want %q", stdout, wantList)
	}
	if addCalls != 1 || showCalls != 1 || listCalls != 1 {
		t.Fatalf("calls = add %d, show %d, list %d", addCalls, showCalls, listCalls)
	}
}

func TestWorkspaceProjectEmptyListAndIdempotentAddOutput(t *testing.T) {
	services := workspaceProjectServices{
		add: func(input workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			return workspace.ProjectResult{Workspace: "/workspace", Project: workspace.ProjectRecord{ID: "ply", Name: "Ply", Wrapper: "/ply", RepoIDs: []workspace.RepoID{"ply"}}, Repos: []workspace.RepoRecord{{ID: "ply", Locator: "/ply/main", GitCommonDir: "/ply/.git"}}}, nil
		},
		show: func(workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ProjectResult{}, errors.New("unused")
		},
		list: func() (workspace.ProjectListResult, error) {
			return workspace.ProjectListResult{Workspace: "/workspace"}, nil
		},
	}
	stdout, _, err := executeProjectCommand(t, services, "add", "ply", "-n", "Ply", "-w", "/ply", "-r", "ply=/ply/main")
	if err != nil || !strings.HasPrefix(stdout, "Project ply (Ply) is already registered in Ply workspace /workspace.\n") {
		t.Fatalf("idempotent stdout = %q, error = %v", stdout, err)
	}
	stdout, _, err = executeProjectCommand(t, services, "list")
	if err != nil || stdout != "No projects are registered in Ply workspace /workspace.\n" {
		t.Fatalf("empty stdout = %q, error = %v", stdout, err)
	}
}

func TestWorkspaceProjectRejectsArgumentsAndFlagsBeforeUseCase(t *testing.T) {
	calls := 0
	services := workspaceProjectServices{
		add: func(workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			calls++
			return workspace.ProjectResult{}, nil
		},
		show: func(workspace.ProjectID) (workspace.ProjectResult, error) {
			calls++
			return workspace.ProjectResult{}, nil
		},
		list: func() (workspace.ProjectListResult, error) { calls++; return workspace.ProjectListResult{}, nil },
	}
	tests := [][]string{
		{"add"},
		{"add", "Trip", "--name", "Trip", "--wrapper", "/trip", "--repo", "repo=/repo"},
		{"add", "trip", "--name", "Trip", "--wrapper", "/trip"},
		{"add", "trip", "--name", "Trip", "--wrapper", "/trip", "--repo", "malformed"},
		{"add", "trip", "--unknown"},
		{"show"},
		{"show", "Trip"},
		{"show", "trip", "extra"},
		{"list", "extra"},
		{"list", "--unknown"},
	}
	for _, arguments := range tests {
		stdout, _, err := executeProjectCommand(t, services, arguments...)
		if err == nil || !strings.HasPrefix(err.Error(), string(workspace.ErrorProjectInvalidArguments)+":") {
			t.Fatalf("args %v error = %v", arguments, err)
		}
		if stdout != "" {
			t.Fatalf("args %v stdout = %q", arguments, stdout)
		}
	}
	if calls != 0 {
		t.Fatalf("use cases called %d times", calls)
	}
}

func TestWorkspaceProjectUseCaseErrorWritesNoSuccessOutput(t *testing.T) {
	want := errors.New("workspace_project_conflict: injected")
	services := workspaceProjectServices{
		add: func(workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			return workspace.ProjectResult{}, want
		},
		show: func(workspace.ProjectID) (workspace.ProjectResult, error) { return workspace.ProjectResult{}, want },
		list: func() (workspace.ProjectListResult, error) { return workspace.ProjectListResult{}, want },
	}
	for _, arguments := range [][]string{{"add", "ply", "-n", "Ply", "-w", "/ply", "-r", "ply=/repo"}, {"show", "ply"}, {"list"}} {
		stdout, _, err := executeProjectCommand(t, services, arguments...)
		if !errors.Is(err, want) || stdout != "" {
			t.Fatalf("args %v = stdout %q, error %v", arguments, stdout, err)
		}
	}
}

func executeProjectCommand(t *testing.T, services workspaceProjectServices, args ...string) (string, string, error) {
	t.Helper()
	command := newWorkspaceProjectCommandWithServices(services)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	command.SetOut(stdout)
	command.SetErr(stderr)
	command.SetArgs(args)
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	return stdout.String(), stderr.String(), err
}
