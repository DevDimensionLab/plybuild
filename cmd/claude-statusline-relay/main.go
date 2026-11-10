// Command claude-statusline-relay records Claude Code's rate_limits for
// `ply agent` and then runs the original status line command unchanged:
//
//	claude-statusline-relay [--then '<original status line command>']
//
// Recording failures never change the status line output or exit status.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agents/quota"
)

const (
	exitCommandNotStarted = 127
	// exitSignalBase follows the shell convention 128+signal.
	exitSignalBase = 128
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var then string
	switch {
	case len(args) == 0:
	case len(args) == 2 && args[0] == "--then":
		then = args[1]
	default:
		// Exit 2 reports the argument error even if stderr cannot be written.
		_, _ = fmt.Fprintln(stderr, "usage: claude-statusline-relay [--then '<original status line command>']")
		return 2
	}
	input, _ := io.ReadAll(stdin)
	_, _ = quota.RecordClaudeStatusLine(input, quota.DefaultClaudeSnapshotPath(), time.Now())
	if then == "" {
		return 0
	}
	original := exec.Command("/bin/sh", "-c", then)
	original.Stdin, original.Stdout, original.Stderr = bytes.NewReader(input), stdout, stderr
	err := original.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if status, ok := exit.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return exitSignalBase + int(status.Signal())
		}
		return exit.ExitCode()
	}
	if err != nil {
		return exitCommandNotStarted
	}
	return 0
}
