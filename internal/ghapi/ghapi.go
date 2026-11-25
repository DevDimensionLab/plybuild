// Package ghapi runs one read-only gh api call and classifies its failure
// with the error codes every ply github command reports.
package ghapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/githubrepo"
)

// CallTimeout bounds each gh call.
const CallTimeout = 20 * time.Second

var (
	ErrMissing = errors.New("gh_missing")
	ErrAuth    = errors.New("gh_auth")
	ErrScope   = errors.New("gh_scope")
	ErrTimeout = errors.New("timeout")
	ErrFailed  = errors.New("gh_failed")
)

// ErrNotPermitted is an HTTP 403 or 404; the caller decides what it means.
var ErrNotPermitted = errors.New("not permitted")

// Codes are the error codes above, in the order a caller reports them.
var Codes = []error{ErrMissing, ErrAuth, ErrScope, ErrTimeout, ErrFailed}

// Run runs gh with args and returns its stdout, also when gh fails: a GraphQL
// answer with partial errors exits non-zero but still carries data.
func Run(ctx context.Context, runner process.Runner, args []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, CallTimeout)
	defer cancel()
	var stdout bytes.Buffer
	err := runner.Run(process.Command{Name: "gh", Args: args, Context: ctx,
		Env: githubrepo.Env(os.Environ()), Stdout: &stdout, Stderr: io.Discard})
	var exitErr *exec.ExitError
	switch {
	case errors.Is(err, exec.ErrNotFound):
		return nil, fmt.Errorf("%w: gh is not installed", ErrMissing)
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded):
		return nil, fmt.Errorf("%w: gh api did not answer within %s", ErrTimeout, CallTimeout)
	case errors.As(err, &exitErr) && exitErr.ExitCode() == 4:
		return nil, fmt.Errorf("%w: gh is not logged in", ErrAuth)
	case err != nil:
		var status struct {
			Status string `json:"status"`
		}
		_ = json.Unmarshal(stdout.Bytes(), &status)
		switch status.Status {
		case "401":
			return nil, fmt.Errorf("%w: gh credentials were rejected", ErrAuth)
		case "403", "404":
			return nil, ErrNotPermitted
		}
		return stdout.Bytes(), fmt.Errorf("%w: gh api failed", ErrFailed)
	}
	return stdout.Bytes(), nil
}

// Get runs one `gh api --method GET path` with fields and decodes its stdout into out.
func Get(ctx context.Context, runner process.Runner, path string, fields []string, out interface{}) error {
	args := []string{"api", "--method", "GET", path}
	for _, f := range fields {
		args = append(args, "-f", f)
	}
	stdout, err := Run(ctx, runner, args)
	if err != nil {
		return err
	}
	if json.Unmarshal(stdout, out) != nil {
		return fmt.Errorf("%w: unreadable gh api response", ErrFailed)
	}
	return nil
}
