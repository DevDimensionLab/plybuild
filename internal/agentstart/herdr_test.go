package agentstart

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func herdrFailing(stream, output string, status int) string {
	redirect := ""
	if stream == "stderr" {
		redirect = " >&2"
	}
	return "#!/bin/sh\ncat <<'EOF'" + redirect + "\n" + output + "\nEOF\nexit " + strconv.Itoa(status) + "\n"
}

const busyReply = `{"id":"cli:agent:start","error":{"code":"agent_pane_busy","message":"agent target pane wT:p1 is not an available shell"}}`

func TestAgentStartHerdrClientKeepsStructuredErrorFromStderr(t *testing.T) {
	_, err := callHerdr(context.Background(), executable(t, herdrFailing("stderr", busyReply, 1)), []string{"agent", "start", "ply-x", "--pane", "wT:p1"})
	want := "herdr agent start: exit status 1 (herdr error agent_pane_busy: agent target pane wT:p1 is not an available shell)"
	if err == nil || err.Error() != want {
		t.Fatalf("got %v", err)
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		t.Fatalf("exit status lost: %v", err)
	}
}

func TestAgentStartHerdrClientKeepsStructuredErrorFromStdoutEnvelope(t *testing.T) {
	_, err := callHerdr(context.Background(), executable(t, herdrFailing("stdout", busyReply, 0)), []string{"agent", "get", "p"})
	if err == nil || !strings.HasSuffix(err.Error(), "unsuccessful response (herdr error agent_pane_busy: agent target pane wT:p1 is not an available shell)") {
		t.Fatalf("got %v", err)
	}
}

func TestAgentStartHerdrClientDropsUnusableStderr(t *testing.T) {
	for name, body := range map[string]string{
		"not json":     herdrFailing("stderr", "secret raw text", 1),
		"bad code":     herdrFailing("stderr", `{"error":{"code":"Bad Code secret","message":"m"}}`, 1),
		"no code":      herdrFailing("stderr", `{"error":{"message":"secret"}}`, 1),
		"trailing":     herdrFailing("stderr", busyReply+"\nsecret raw text", 1),
		"error string": herdrFailing("stderr", `{"error":"secret"}`, 1),
		"oversized":    "#!/bin/sh\nprintf '{\"error\":{\"code\":\"agent_pane_busy\",\"message\":\"' >&2\nhead -c 2000000 /dev/zero | tr '\\0' 'a' >&2\nprintf '\"}}' >&2\nexit 1\n",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := callHerdr(context.Background(), executable(t, body), []string{"agent", "start"})
			if err == nil || err.Error() != "herdr agent start: exit status 1" {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestAgentStartHerdrClientBoundsMessage(t *testing.T) {
	long := strings.Repeat("é", 200)
	reply := `{"error":{"code":"agent_start_input_failed","message":"line\u0001one\nline two ` + long + `"}}`
	_, err := callHerdr(context.Background(), executable(t, herdrFailing("stderr", reply, 1)), []string{"agent", "start"})
	if err == nil {
		t.Fatal("failure accepted")
	}
	msg := err.Error()
	detail := msg[strings.Index(msg, "agent_start_input_failed: ")+len("agent_start_input_failed: "):]
	detail = strings.TrimSuffix(detail, ")")
	if !strings.HasPrefix(detail, "line one line two é") || !strings.HasSuffix(detail, "…") || len(detail) > 160+len("…") {
		t.Fatalf("unbounded or unsanitized message %q", detail)
	}
}

func TestAgentStartHerdrClientWithholdsMessageForSensitivePayloadCalls(t *testing.T) {
	prompt := "Read the Goal and implement the private change described in the Spec"
	for name, tc := range map[string]struct {
		args    []string
		message string
	}{
		"full prompt":         {promptArgs("wT:p1", prompt), "input rejected: " + prompt},
		"partial prompt":      {promptArgs("wT:p1", prompt), "input rejected near private change"},
		"single token prompt": {promptArgs("wT:p1", "hunter2"), "input rejected: hunter2"},
		"flag-like prompt":    {promptArgs("wT:p1", "--secret-token"), "input rejected: --secret-token"},
		"path-like prompt":    {promptArgs("wT:p1", "/private/plan.md"), "input rejected: /private/plan.md"},
		"short env value":     {[]string{"tab", "create", "--cwd", "/w", "--env", "K=ab", "--no-focus"}, "bad env value ab"},
	} {
		t.Run(name, func(t *testing.T) {
			reply := `{"error":{"code":"agent_start_input_failed","message":"` + tc.message + `"}}`
			_, err := callHerdr(context.Background(), executable(t, herdrFailing("stderr", reply, 1)), tc.args)
			if err == nil || !strings.HasSuffix(err.Error(), "(herdr error agent_start_input_failed: message withheld)") {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestAgentStartHerdrClientKeepsMessageForIdentityOnlyCalls(t *testing.T) {
	model := "opus"
	in := StartInput{Key: "k", Provider: ProviderClaude, Model: &model}
	reply := `{"error":{"code":"agent_pane_busy","message":"agent target pane wT:p1 is not an available shell"}}`
	_, err := callHerdr(context.Background(), executable(t, herdrFailing("stderr", reply, 1)), startArgs(in, "wT:p1"))
	if err == nil || !strings.HasSuffix(err.Error(), "(herdr error agent_pane_busy: agent target pane wT:p1 is not an available shell)") {
		t.Fatalf("got %v", err)
	}
}

func TestAgentStartHerdrClientAddsNoDetailAfterTimeout(t *testing.T) {
	body := "#!/bin/sh\necho '" + busyReply + "' >&2\nexec sleep 5\n"
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	_, err := callHerdr(ctx, executable(t, body), []string{"agent", "start"})
	if err == nil || strings.Contains(err.Error(), "herdr error") {
		t.Fatalf("got %v", err)
	}
	if ctx.Err() == nil {
		t.Fatal("call returned before timeout")
	}
}
