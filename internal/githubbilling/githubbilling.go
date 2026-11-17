// Package githubbilling reads the billing usage rows of the GitHub user that
// gh is logged in as, for one month, through two read-only gh api calls. It
// needs the gh scope user and never changes the login or its scopes.
package githubbilling

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
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

// Item is one billing usage row: one day, product, SKU and repository.
type Item struct {
	Product, SKU, UnitType, Repository string
	Quantity                           float64
}

// Usage is the login and its billing rows of one month.
type Usage struct {
	Login string
	Items []Item
}

var login = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

// Read returns the billing usage of the logged-in gh user in year and month.
// Errors carry one of the sentinels above and never gh's output.
func Read(ctx context.Context, runner process.Runner, year int, month time.Month) (Usage, error) {
	var user struct {
		Login string `json:"login"`
	}
	if err := call(ctx, runner, "/user", &user); err != nil {
		if errors.Is(err, errNotPermitted) {
			return Usage{}, fmt.Errorf("%w: gh could not read the logged-in user", ErrAuth)
		}
		return Usage{}, err
	}
	if !login.MatchString(user.Login) {
		return Usage{}, fmt.Errorf("%w: unreadable login", ErrFailed)
	}
	var body struct {
		UsageItems []struct {
			Product        string  `json:"product"`
			SKU            string  `json:"sku"`
			UnitType       string  `json:"unitType"`
			RepositoryName string  `json:"repositoryName"`
			Quantity       float64 `json:"quantity"`
		} `json:"usageItems"`
	}
	path := fmt.Sprintf("/users/%s/settings/billing/usage?year=%d&month=%d", user.Login, year, int(month))
	if err := call(ctx, runner, path, &body); err != nil {
		if errors.Is(err, errNotPermitted) {
			return Usage{}, fmt.Errorf("%w: billing usage needs the gh scope user", ErrScope)
		}
		return Usage{}, err
	}
	usage := Usage{Login: user.Login, Items: []Item{}}
	for _, i := range body.UsageItems {
		usage.Items = append(usage.Items, Item{Product: i.Product, SKU: i.SKU, UnitType: i.UnitType, Repository: i.RepositoryName, Quantity: i.Quantity})
	}
	return usage, nil
}

// errNotPermitted is an HTTP 403 or 404: ErrScope for billing, ErrAuth for the user.
var errNotPermitted = errors.New("not permitted")

// call runs one `gh api --method GET path` and decodes its stdout into out.
func call(ctx context.Context, runner process.Runner, path string, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, CallTimeout)
	defer cancel()
	var stdout bytes.Buffer
	err := runner.Run(process.Command{Name: "gh", Args: []string{"api", "--method", "GET", path}, Context: ctx,
		Env: githubrepo.Env(os.Environ()), Stdout: &stdout, Stderr: io.Discard})
	var exitErr *exec.ExitError
	switch {
	case errors.Is(err, exec.ErrNotFound):
		return fmt.Errorf("%w: gh is not installed", ErrMissing)
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded):
		return fmt.Errorf("%w: gh api did not answer within %s", ErrTimeout, CallTimeout)
	case errors.As(err, &exitErr) && exitErr.ExitCode() == 4:
		return fmt.Errorf("%w: gh is not logged in", ErrAuth)
	case err != nil:
		var status struct {
			Status string `json:"status"`
		}
		_ = json.Unmarshal(stdout.Bytes(), &status)
		switch status.Status {
		case "401":
			return fmt.Errorf("%w: gh credentials were rejected", ErrAuth)
		case "403", "404":
			return errNotPermitted
		}
		return fmt.Errorf("%w: gh api failed", ErrFailed)
	}
	if json.Unmarshal(stdout.Bytes(), out) != nil {
		return fmt.Errorf("%w: unreadable gh api response", ErrFailed)
	}
	return nil
}
