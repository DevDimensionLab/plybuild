package githubbilling

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type reply struct {
	stdout string
	code   int
	err    error
}

type fakeGH struct {
	replies []reply
	args    [][]string
	env     [][]string
	timeout []time.Duration
}

func (f *fakeGH) Run(c process.Command) error {
	if c.Name != "gh" || c.Context == nil || c.Env == nil {
		return fmt.Errorf("unexpected command %s %q", c.Name, c.Args)
	}
	deadline, _ := c.Context.Deadline()
	f.args = append(f.args, c.Args)
	f.env = append(f.env, c.Env)
	f.timeout = append(f.timeout, time.Until(deadline).Round(time.Second))
	r := f.replies[0]
	f.replies = f.replies[1:]
	fmt.Fprint(c.Stdout, r.stdout)
	fmt.Fprint(c.Stderr, "STDERR-SECRET gho_FAKETOKEN")
	if r.err != nil {
		return r.err
	}
	if r.code != 0 {
		return exec.Command("/bin/sh", "-c", fmt.Sprintf("exit %d", r.code)).Run()
	}
	return nil
}

const usageBody = `{"usageItems": [
 {"date": "2026-11-01T00:00:00Z", "product": "actions", "sku": "Actions Linux", "quantity": 10.0, "unitType": "Minutes", "repositoryName": "app"},
 {"date": "2026-11-02T00:00:00Z", "product": "actions", "sku": "Actions storage", "quantity": 0.5, "unitType": "GigabyteHours", "repositoryName": "app"}]}`

func TestReadReturnsTheLoginAndTheMonthsRows(t *testing.T) {
	gh := &fakeGH{replies: []reply{{stdout: `{"login": "octo"}`}, {stdout: usageBody}}}
	usage, err := Read(context.Background(), gh, 2026, time.November)
	if err != nil {
		t.Fatal(err)
	}
	want := Usage{Login: "octo", Items: []Item{
		{Product: "actions", SKU: "Actions Linux", UnitType: "Minutes", Repository: "app", Quantity: 10},
		{Product: "actions", SKU: "Actions storage", UnitType: "GigabyteHours", Repository: "app", Quantity: 0.5},
	}}
	if !reflect.DeepEqual(usage, want) {
		t.Fatalf("usage %+v, want %+v", usage, want)
	}
	wantArgs := [][]string{
		{"api", "--method", "GET", "/user"},
		{"api", "--method", "GET", "/users/octo/settings/billing/usage?year=2026&month=11"},
	}
	if !reflect.DeepEqual(gh.args, wantArgs) {
		t.Fatalf("args %q, want %q", gh.args, wantArgs)
	}
	for i, d := range gh.timeout {
		if d != CallTimeout {
			t.Fatalf("call %d timeout %v, want %v", i, d, CallTimeout)
		}
	}
	if env := strings.Join(gh.env[0], "\n"); !strings.Contains(env, "GH_HOST=github.com") || !strings.Contains(env, "GH_PROMPT_DISABLED=1") {
		t.Fatalf("gh environment is not pinned: %q", gh.env[0])
	}
}

func TestReadClassifiesFailuresWithoutGhOutput(t *testing.T) {
	notFound := `{"message": "Not Found", "status": "404"}`
	for name, test := range map[string]struct {
		replies []reply
		want    error
	}{
		"gh missing":        {[]reply{{err: &exec.Error{Name: "gh", Err: exec.ErrNotFound}}}, ErrMissing},
		"not logged in":     {[]reply{{code: 4}}, ErrAuth},
		"bad credentials":   {[]reply{{stdout: `{"message": "Bad credentials", "status": "401"}`, code: 1}}, ErrAuth},
		"user not found":    {[]reply{{stdout: notFound, code: 1}}, ErrAuth},
		"billing not found": {[]reply{{stdout: `{"login": "octo"}`}, {stdout: notFound, code: 1}}, ErrScope},
		"billing forbidden": {[]reply{{stdout: `{"login": "octo"}`}, {stdout: `{"message": "x", "status": "403"}`, code: 1}}, ErrScope},
		"timeout":           {[]reply{{err: context.DeadlineExceeded}}, ErrTimeout},
		"server error":      {[]reply{{stdout: `{"login": "octo"}`}, {stdout: `{"status": "500"}`, code: 1}}, ErrFailed},
		"unreadable usage":  {[]reply{{stdout: `{"login": "octo"}`}, {stdout: `<html>`}}, ErrFailed},
		"unsafe login":      {[]reply{{stdout: `{"login": "a/b"}`}}, ErrFailed},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Read(context.Background(), &fakeGH{replies: test.replies}, 2026, time.November)
			if !errors.Is(err, test.want) {
				t.Fatalf("error %v, want %v", err, test.want)
			}
			if strings.Contains(err.Error(), "SECRET") || strings.Contains(err.Error(), "gho_") || strings.Contains(err.Error(), "Not Found") {
				t.Fatalf("error exposes gh output: %v", err)
			}
		})
	}
}

func TestReadCanceledByTheCallerIsNotATimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Read(ctx, &fakeGH{replies: []reply{{err: context.Canceled}}}, 2026, time.November)
	if !errors.Is(err, ErrFailed) {
		t.Fatalf("error %v, want %v", err, ErrFailed)
	}
}
