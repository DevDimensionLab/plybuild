package cmd

import (
	"bytes"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
	"testing"
)

func TestTaskRunArgumentsAndEnglishHelp(t *testing.T) {
	for _, args := range [][]string{{"run", "start"}, {"run", "start", "--file", "/x", "extra"}, {"run", "start", "--file", "/x", "--apply", "--confirm", "sha256:x", "--format", "json"}, {"run", "start", "--file", "/x", "--confirm", "x"}, {"run", "start", "--file", "/x", "--check", "--apply", "--confirm", "x"}, {"run", "start", "--next"}, {"run", "show", "id", "--file", "/x"}, {"run", "accept", "id"}} {
		root := &cobra.Command{Use: "test", SilenceErrors: true, SilenceUsage: true}
		root.AddCommand(newTaskRunCommand(taskrun.Dependencies{}))
		root.SetArgs(args)
		root.SetOut(&bytes.Buffer{})
		root.SetErr(&bytes.Buffer{})
		e := root.Execute()
		var typed *taskrun.Error
		if !errors.As(e, &typed) || typed.Exit != 2 {
			t.Fatalf("%v: %v", args, e)
		}
	}
	for _, leaf := range []string{"", "start", "show", "collect", "accept", "report"} {
		root := &cobra.Command{Use: "test"}
		root.AddCommand(newTaskRunCommand(taskrun.Dependencies{}))
		b := &bytes.Buffer{}
		root.SetOut(b)
		args := []string{"run"}
		if leaf != "" {
			args = append(args, leaf, "--help")
		}
		root.SetArgs(args)
		if e := root.Execute(); e != nil || b.Len() == 0 {
			t.Fatalf("help %s: %v", leaf, e)
		}
	}
}

func TestTaskRunExecutionReturnsAndRestoresRoot(t *testing.T) {
	original := RootCmd
	t.Cleanup(func() { RootCmd = original })
	RootCmd = &cobra.Command{Use: "ply"}
	RootCmd.AddCommand(newTaskRunCommand(taskrun.Dependencies{}))
	RootCmd.SetArgs([]string{"run", "start"})
	var stdout, stderr bytes.Buffer
	RootCmd.SetOut(&stdout)
	RootCmd.SetErr(&stderr)
	deferred := false
	err := func() error {
		defer func() { deferred = true }()
		return ExecuteE()
	}()
	var typed *taskrun.Error
	if !errors.As(err, &typed) || typed.Exit != 2 || !deferred {
		t.Fatalf("execution did not return its error through defers: err=%v deferred=%v", err, deferred)
	}
	if RootCmd.SilenceErrors || RootCmd.SilenceUsage {
		t.Fatal("ExecuteE did not restore diagnostic settings")
	}
	if stdout.Len() != 0 || stderr.String() != "Error: task_run_invalid_arguments: --file is required\n" {
		t.Fatalf("unexpected or duplicate diagnostics: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}
