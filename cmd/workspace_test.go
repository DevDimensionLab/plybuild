package cmd

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type commandFileSystem struct {
	workspace.FileSystem
	cwd        string
	getwdCalls int
	getwdErr   error
}

func (filesystem *commandFileSystem) Getwd() (string, error) {
	filesystem.getwdCalls++
	if filesystem.getwdErr != nil {
		return "", filesystem.getwdErr
	}
	return filesystem.cwd, nil
}

func TestNewWorkspaceCommandShape(t *testing.T) {
	command := newWorkspaceCommand(workspace.Dependencies{})
	root := &cobra.Command{Use: "ply"}
	root.AddCommand(command)

	if command.Use != "workspace" ||
		command.Short != "Manage Ply workspaces" ||
		command.Long != "Manage explicit local Ply workspaces." ||
		command.Example != "  ply workspace init\n  ply workspace project list" ||
		command.CommandPath() != "ply workspace" ||
		command.Runnable() {
		t.Fatalf("workspace command metadata = Use %q, Short %q, Long %q, Example %q, path %q, runnable %t",
			command.Use, command.Short, command.Long, command.Example, command.CommandPath(), command.Runnable())
	}
	if command.PersistentPreRunE == nil {
		t.Fatal("workspace command does not shadow the root persistent hook")
	}
	initCommand, _, err := command.Find([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}
	if initCommand == command ||
		initCommand.Use != "init" ||
		initCommand.Short != "Initialize a Ply workspace in the current directory" ||
		initCommand.Long != "Initialize a Ply workspace in the current directory by creating .ply/workspace.yaml.\nThe current directory does not need to be a Git repository." ||
		initCommand.Example != "  ply workspace init" ||
		initCommand.CommandPath() != "ply workspace init" ||
		!initCommand.Runnable() {
		t.Fatalf("init command metadata = Use %q, Short %q, Long %q, Example %q, path %q, runnable %t",
			initCommand.Use, initCommand.Short, initCommand.Long, initCommand.Example, initCommand.CommandPath(), initCommand.Runnable())
	}
	if initCommand.Args == nil {
		t.Fatal("init command has no argument validator")
	}
	initCommand.InitDefaultHelpFlag()
	localFlags := initCommand.LocalNonPersistentFlags()
	localFlagNames := make([]string, 0)
	localFlags.VisitAll(func(flag *pflag.Flag) {
		localFlagNames = append(localFlagNames, flag.Name)
	})
	if !reflect.DeepEqual(localFlagNames, []string{"help"}) || localFlags.Lookup("name") != nil {
		t.Fatalf("workspace init local product flags changed: %s", localFlags.FlagUsages())
	}
}

func TestWorkspaceCommandRendersCreatedAndAlreadyInitialized(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	filesystem := &commandFileSystem{FileSystem: workspace.SystemDependencies().Files, cwd: root}

	stdout, stderr, err := executeWorkspaceCommand(t, filesystem, "init")
	if err != nil {
		t.Fatal(err)
	}
	if want := "Initialized Ply workspace at " + root + ".\n"; stdout != want {
		t.Fatalf("created stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("created stderr = %q", stderr)
	}

	stdout, stderr, err = executeWorkspaceCommand(t, filesystem, "init")
	if err != nil {
		t.Fatal(err)
	}
	if want := "Ply workspace already initialized at " + root + ".\n"; stdout != want {
		t.Fatalf("idempotent stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("idempotent stderr = %q", stderr)
	}
	if filesystem.getwdCalls != 2 {
		t.Fatalf("workspace.Init call count = %d, want 2", filesystem.getwdCalls)
	}
}

func TestWorkspaceCommandRejectsArgumentsAndLocalFlags(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"position", []string{"init", "extra"}, "workspace_init_invalid_arguments: expected no positional arguments, got 1"},
		{"name flag", []string{"init", "--name", "value"}, "workspace_init_invalid_arguments: unknown flag: --name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filesystem := &commandFileSystem{FileSystem: workspace.SystemDependencies().Files, cwd: root}
			stdout, _, err := executeWorkspaceCommand(t, filesystem, test.args...)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
			if stdout != "" {
				t.Fatalf("invalid input stdout = %q", stdout)
			}
			if filesystem.getwdCalls != 0 {
				t.Fatalf("use case called %d times for invalid input", filesystem.getwdCalls)
			}
		})
	}
}

func TestWorkspaceCommandDoesNotPrintSuccessOnUseCaseError(t *testing.T) {
	filesystem := &commandFileSystem{
		FileSystem: workspace.SystemDependencies().Files,
		getwdErr:   errors.New("cwd unavailable"),
	}
	stdout, _, err := executeWorkspaceCommand(t, filesystem, "init")
	if err == nil || !strings.HasPrefix(err.Error(), "workspace_init_path_error: getwd:") {
		t.Fatalf("error = %v", err)
	}
	if stdout != "" {
		t.Fatalf("error stdout = %q", stdout)
	}
	if filesystem.getwdCalls != 1 {
		t.Fatalf("use case called %d times, want 1", filesystem.getwdCalls)
	}
}

func TestWorkspaceHookShadowsLegacyInitialization(t *testing.T) {
	for _, arguments := range [][]string{
		{"workspace", "init"},
		{"workspace", "project", "add", "ply", "--name", "Ply", "--wrapper", ".", "--repo", "ply=."},
		{"workspace", "project", "show", "ply"},
		{"workspace", "project", "list"},
	} {
		rootDirectory := t.TempDir()
		filesystem := &commandFileSystem{FileSystem: workspace.SystemDependencies().Files, cwd: rootDirectory}
		dependencies := workspace.SystemDependencies()
		dependencies.Files = filesystem
		legacyCalls := 0
		root := &cobra.Command{
			Use: "ply",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				legacyCalls++
				return errors.New("legacy initialization ran")
			},
		}
		root.AddCommand(newWorkspaceCommand(dependencies))
		root.SetArgs(arguments)
		root.SilenceErrors = true
		root.SilenceUsage = true
		root.SetOut(&bytes.Buffer{})
		root.SetErr(&bytes.Buffer{})
		err := root.Execute()
		if arguments[1] == "init" && err != nil {
			t.Fatalf("workspace hook did not isolate init: %v", err)
		}
		if err != nil && strings.Contains(err.Error(), "legacy initialization ran") {
			t.Fatalf("workspace hook did not isolate %v: %v", arguments, err)
		}
		if legacyCalls != 0 {
			t.Fatalf("legacy initialization called %d times for %v", legacyCalls, arguments)
		}
	}
}

func executeWorkspaceCommand(t *testing.T, filesystem workspace.FileSystem, args ...string) (string, string, error) {
	t.Helper()
	command := newWorkspaceCommand(workspace.Dependencies{Files: filesystem})
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
