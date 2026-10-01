package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/cmd"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestMainDelegatesToCommand(t *testing.T) {
	original := execute
	t.Cleanup(func() {
		execute = original
	})

	calls := 0
	execute = func() error {
		calls++
		return nil
	}

	main()

	if calls != 1 {
		t.Fatalf("command entrypoint was called %d times, want 1", calls)
	}
}

func TestMainProcessExitCodes(t *testing.T) {
	for _, tc := range []struct {
		name string
		want int
	}{
		{"success", 0}, {"legacy", 1}, {"2", 2}, {"3", 3}, {"4", 4}, {"5", 5},
		{"wrapped", 3}, {"0", 1}, {"1", 1}, {"-1", 1}, {"6", 1}, {"256", 1}, {"typed-nil", 1},
		{"command", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			child := exec.Command(os.Args[0], "-test.run=^TestMainProcessHelper$")
			child.Env = append(os.Environ(), "PLY_ENTRYPOINT_TEST_CASE="+tc.name)
			var stdout, stderr bytes.Buffer
			child.Stdout, child.Stderr = &stdout, &stderr
			err := child.Run()
			got := 0
			if err != nil {
				var exit *exec.ExitError
				if !errors.As(err, &exit) {
					t.Fatal(err)
				}
				got = exit.ExitCode()
			}
			if got != tc.want {
				t.Fatalf("exit = %d, want %d; stdout=%q stderr=%q", got, tc.want, stdout.String(), stderr.String())
			}
			if tc.name == "command" {
				if want := "Error: task_run_invalid_arguments: --file is required\n"; stderr.String() != want {
					t.Fatalf("diagnostic = %q, want exactly %q", stderr.String(), want)
				}
			} else if stderr.Len() != 0 {
				t.Fatalf("entrypoint unexpectedly printed %q", stderr.String())
			}
			if strings.Count(stdout.String(), "execute defer\n") != 1 {
				t.Fatalf("command defer did not run once: %q", stdout.String())
			}
		})
	}
}

// This child runs the real main boundary with synthetic command outcomes.
func TestMainProcessHelper(t *testing.T) {
	name := os.Getenv("PLY_ENTRYPOINT_TEST_CASE")
	if name == "" {
		return
	}
	execute = func() error {
		defer fmt.Fprintln(os.Stdout, "execute defer")
		switch name {
		case "success":
			return nil
		case "legacy":
			return errors.New("legacy command error")
		case "wrapped":
			return fmt.Errorf("wrapped: %w", &taskrun.Error{Code: "synthetic", Exit: 3})
		case "typed-nil":
			return (*taskrun.Error)(nil)
		case "command":
			cmd.RootCmd.SetArgs([]string{"workspace", "task", "run", "start"})
			return cmd.ExecuteE()
		default:
			code, err := strconv.Atoi(name)
			if err != nil {
				t.Fatal(err)
			}
			return &taskrun.Error{Code: "synthetic", Exit: code}
		}
	}
	main()
}
