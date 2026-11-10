//go:build !windows

package quota

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// fakeCodex writes an isolated sh app-server that logs each request, answers
// initialize with noise in between, then prints response and runs tail.
func fakeCodex(t *testing.T, response, tail string) (executable, log string) {
	t.Helper()
	dir := t.TempDir()
	executable, log = filepath.Join(dir, "codex"), filepath.Join(dir, "requests.log")
	script := "#!/bin/sh\n[ \"$1\" = app-server ] || exit 9\n" +
		"read -r l; printf '%s\\n' \"$l\" >> '" + log + "'\n" +
		"printf '%s\\n' 'not json' '{\"method\":\"account/rateLimits/updated\",\"params\":{}}' '{\"id\":1,\"result\":{\"x\":1}}'\n" +
		"read -r l; printf '%s\\n' \"$l\" >> '" + log + "'\n" +
		"read -r l; printf '%s\\n' \"$l\" >> '" + log + "'\n" +
		"printf '%s\\n' '{\"method\":\"noise\"}'\n" + response + "\n" + tail + "\n"
	if err := os.WriteFile(executable, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return executable, log
}

func catFixture(name string) string {
	path, _ := filepath.Abs(filepath.Join("testdata", name))
	return "cat '" + path + "'"
}

func TestReadCodexWeeklyOnlyAndWireOrder(t *testing.T) {
	exe, log := fakeCodex(t, catFixture("codex-weekly-only.json"), "")
	before := time.Now()
	got, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
	if err != nil {
		t.Fatal(err)
	}
	reset := time.Unix(1794673954, 0)
	want := []Window{{ID: Weekly, UsedPercent: 5, ResetsAt: &reset}}
	if got.Source != SourceCodex || got.Plan != "pro" || got.LimitReached || !reflect.DeepEqual(got.Windows, want) || got.ObservedAt.Before(before) {
		t.Fatalf("snapshot = %+v", got)
	}
	if got.Windows[0].Remaining() != 95 {
		t.Fatalf("remaining %v", got.Windows[0].Remaining())
	}
	raw, _ := os.ReadFile(log)
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	methods := []string{`"id":1,"method":"initialize"`, `"method":"initialized"`, `"id":2,"method":"account/rateLimits/read"`}
	if len(lines) != 3 {
		t.Fatalf("requests %q", lines)
	}
	for i, m := range methods {
		if !strings.Contains(lines[i], m) {
			t.Fatalf("request %d = %s, want %s", i, lines[i], m)
		}
	}
}

func TestReadCodexPoolsLimitReachedAndNullReset(t *testing.T) {
	exe, _ := fakeCodex(t, catFixture("codex-pools.json"), "")
	got, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
	if err != nil {
		t.Fatal(err)
	}
	reset := time.Unix(1794673954, 0)
	want := []Window{
		{ID: FiveHour, UsedPercent: 38},
		{ID: Weekly, UsedPercent: 82, ResetsAt: &reset},
		{ID: WindowID{Minutes: 1440}, Pool: "a-other", UsedPercent: 7},
		{ID: Weekly, Pool: "Spark", UsedPercent: 99, ResetsAt: &reset},
	}
	if !got.LimitReached || got.Plan != "plus" || !reflect.DeepEqual(got.Windows, want) {
		t.Fatalf("snapshot = %+v", got.Windows)
	}
}

func TestReadCodexResponseErrors(t *testing.T) {
	cases := map[string]struct {
		response string
		want     error
	}{
		"null rate limits":    {`{"id":2,"result":{"rateLimits":null}}`, ErrNoQuotaReported},
		"missing rate limits": {`{"id":2,"result":{}}`, ErrNoQuotaReported},
		"missing used":        {`{"id":2,"result":{"rateLimits":{"primary":{"windowDurationMins":300}}}}`, ErrInvalidResponse},
		"wrong type":          {`{"id":2,"result":{"rateLimits":{"primary":{"usedPercent":"5","windowDurationMins":300}}}}`, ErrInvalidResponse},
		"wrong pool type":     {`{"id":2,"result":{"rateLimits":{},"rateLimitsByLimitId":{"p":{"primary":{"usedPercent":true}}}}}`, ErrInvalidResponse},
		"auth":                {`{"id":2,"error":{"code":-32001,"message":"Please Log In first CANARY-secret"}}`, ErrNotSignedIn},
		"early eof":           {`exit 0`, ErrTimeout},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			response := c.response
			if strings.HasPrefix(response, "{") {
				response = "printf '%s\\n' '" + response + "'"
			}
			exe, _ := fakeCodex(t, response, "")
			got, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
			if !errors.Is(err, c.want) || got != nil {
				t.Fatalf("got %v, %v; want %v", got, err, c.want)
			}
			if strings.Contains(err.Error(), "CANARY") {
				t.Fatalf("error leaks upstream text: %v", err)
			}
		})
	}
}

func TestReadCodexValidEmptyRateLimitsHasNoWindows(t *testing.T) {
	exe, _ := fakeCodex(t, `printf '%s\n' '{"id":2,"result":{"rateLimits":{"primary":null,"secondary":{"usedPercent":3,"windowDurationMins":null}}}}'`, "")
	got, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
	if err != nil || len(got.Windows) != 0 {
		t.Fatalf("got %+v, %v", got, err)
	}
}

func TestReadCodexUpstreamErrorIsSanitized(t *testing.T) {
	message := strings.Repeat("x", 500) + " CANARY-raw"
	exe, _ := fakeCodex(t, `printf '%s\n' '{"id":2,"error":{"code":-32603,"message":"`+message+`"}}'`, "")
	_, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
	var upstream *UpstreamError
	if !errors.As(err, &upstream) || upstream.Code != -32603 {
		t.Fatalf("err = %v", err)
	}
	if strings.Contains(err.Error(), "CANARY") || strings.Contains(err.Error(), "xxx") {
		t.Fatalf("error leaks upstream text: %v", err)
	}
}

func TestReadCodexTimesOutOnSilentOrNoisyServer(t *testing.T) {
	for name, tail := range map[string]string{
		"silent": "exec sleep 30",
		"noisy":  "while :; do printf 'garbage\\n{\"method\":\"n\"}\\n'; done",
	} {
		t.Run(name, func(t *testing.T) {
			exe, _ := fakeCodex(t, "", tail)
			start := time.Now()
			_, err := ReadCodex(context.Background(), CodexOptions{Executable: exe, Timeout: 500 * time.Millisecond})
			if !errors.Is(err, ErrTimeout) || time.Since(start) > 3*time.Second {
				t.Fatalf("err %v after %v", err, time.Since(start))
			}
		})
	}
}

func TestReadCodexParentCancellationIsNotOwnTimeout(t *testing.T) {
	exe, _ := fakeCodex(t, "", "exec sleep 30")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err := ReadCodex(ctx, CodexOptions{Executable: exe, Timeout: 10 * time.Second})
	if !errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTimeout) {
		t.Fatalf("err = %v", err)
	}
}

func TestReadCodexReapsLingeringServerAndChildren(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "child.pid")
	exe, _ := fakeCodex(t, "sleep 30 & echo $! > '"+pidFile+"'; "+catFixture("codex-weekly-only.json"), "wait")
	start := time.Now()
	got, err := ReadCodex(context.Background(), CodexOptions{Executable: exe})
	if err != nil || got.Windows[0].Remaining() != 95 {
		t.Fatalf("got %+v, %v", got, err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("lingering server delayed return by %v", elapsed)
	}
	raw, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(raw)))
	deadline := time.Now().Add(time.Second)
	for syscall.Kill(pid, 0) == nil {
		if time.Now().After(deadline) {
			t.Fatalf("child %d still running", pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestReadCodexInvalidExplicitExecutableDoesNotFallBack(t *testing.T) {
	_, err := ReadCodex(context.Background(), CodexOptions{Executable: filepath.Join(t.TempDir(), "missing")})
	if !errors.Is(err, ErrCLINotFound) {
		t.Fatalf("err = %v", err)
	}
}

func TestLocateCodexSearchOrder(t *testing.T) {
	first, second, empty := t.TempDir(), t.TempDir(), t.TempDir()
	for _, dir := range []string{first, second} {
		if err := os.WriteFile(filepath.Join(dir, "codex"), []byte("#!/bin/sh\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(empty, "codex"), []byte("not executable"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := locateCodex("", empty, []string{empty, second, first})
	if err != nil || got != filepath.Join(second, "codex") {
		t.Fatalf("locate = %q, %v", got, err)
	}
	if _, err := locateCodex("", empty, []string{empty}); !errors.Is(err, ErrCLINotFound) {
		t.Fatalf("err = %v", err)
	}
}
