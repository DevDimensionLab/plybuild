package agentlaunch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWorkingAndSingleSubmission(t *testing.T) {
	calls := 0
	polls := 0
	run := func(ctx context.Context, args []string) (json.RawMessage, error) {
		switch strings.Join(args[:2], " ") {
		case "tab create":
			return json.RawMessage(`{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","workspace_id":"w1","terminal_id":"term2"}}`), nil
		case "pane process-info":
			return json.RawMessage(`{"process_info":{"pane_id":"w1:p2","shell_pid":42,"foreground_process_group_id":42,"foreground_processes":[{"pid":42,"name":"zsh"}]}}`), nil
		case "pane run":
			calls++
			if !strings.HasPrefix(args[3], "cd -- '/work' && exec ") {
				t.Fatalf("cwd not established before agent start: %q", args[3])
			}
			if !strings.Contains(args[3], `'hello '\''$()'`) {
				t.Fatalf("unsafe command: %q", args)
			}
			return json.RawMessage(`{}`), nil
		case "pane current":
			polls++
			agent := ""
			if calls > 0 {
				agent = `,"agent":"claude","agent_status":"working","foreground_cwd":"/work"`
			}
			return json.RawMessage(`{"pane":{"pane_id":"w1:p2","tab_id":"w1:t2","workspace_id":"w1","terminal_id":"term2"` + agent + `}}`), nil
		}
		t.Fatalf("unexpected %v", args)
		return nil, nil
	}
	r, err := launch(context.Background(), Request{Agent: "claude", Executable: "/bin/claude", Cwd: "/work", Workspace: "w1", InitialPrompt: "hello '$()", Timeout: time.Second}, run, time.Millisecond)
	if err != nil || r.State != "agent_observed" || r.Agent != "claude" || r.AgentStatus != "working" || calls != 1 || polls < 2 {
		t.Fatalf("%+v %v calls=%d", r, err, calls)
	}
}
func TestUnknownSendNotRetried(t *testing.T) {
	calls := 0
	run := func(ctx context.Context, args []string) (json.RawMessage, error) {
		switch strings.Join(args[:2], " ") {
		case "tab create":
			return json.RawMessage(`{"root_pane":{"pane_id":"p","tab_id":"t","workspace_id":"w","terminal_id":"term"}}`), nil
		case "pane current":
			return json.RawMessage(`{"pane":{"pane_id":"p","tab_id":"t","workspace_id":"w","terminal_id":"term"}}`), nil
		case "pane process-info":
			return json.RawMessage(`{"process_info":{"pane_id":"p","shell_pid":1,"foreground_process_group_id":1,"foreground_processes":[{"pid":1,"name":"bash"}]}}`), nil
		case "pane run":
			calls++
			return nil, errors.New("lost reply")
		}
		return nil, errors.New("unexpected")
	}
	r, err := launch(context.Background(), Request{Agent: "claude", Executable: "/bin/claude", Cwd: "/work", Workspace: "w", InitialPrompt: "hello", Timeout: time.Second}, run, time.Millisecond)
	if err == nil || !r.EffectUnknown || r.Terminal != "term" || calls != 1 {
		t.Fatalf("%+v %v %d", r, err, calls)
	}
}

func TestShellChecks(t *testing.T) {
	for _, tc := range []struct {
		body           string
		ready, wantErr bool
	}{
		{`{"pane_id":"p","shell_pid":1,"foreground_process_group_id":2,"foreground_processes":[{"pid":2,"name":"sleep"}]}`, false, false},
		{`{"pane_id":"p","shell_pid":1,"foreground_process_group_id":1,"foreground_processes":[{"pid":1,"name":"fish"}]}`, false, true},
		{`{"pane_id":"other"}`, false, true},
		{`{"pane_id":"p","shell_pid":1,"foreground_process_group_id":1,"foreground_processes":[{"pid":1,"name":"-zsh"}]}`, true, false},
	} {
		ready, err := shellAvailable(json.RawMessage(`{"process_info":`+tc.body+`}`), "p")
		if ready != tc.ready || (err != nil) != tc.wantErr {
			t.Fatalf("%s: %v %v", tc.body, ready, err)
		}
	}
}
func TestIdentityReplacement(t *testing.T) {
	_, err := readPane(context.Background(), pane{Pane: "p", Tab: "t", Terminal: "original", Workspace: "w"}, func(context.Context, []string) (json.RawMessage, error) {
		return json.RawMessage(`{"pane":{"pane_id":"p","tab_id":"t","terminal_id":"replacement","workspace_id":"w"}}`), nil
	})
	if err == nil {
		t.Fatal("replacement accepted")
	}
}
func TestInputRejectedBeforeEffect(t *testing.T) {
	for _, prompt := range []string{"", "a\nb", "a\x00b"} {
		_, err := launch(context.Background(), Request{Agent: "claude", Executable: "/bin/claude", Cwd: "/work", Workspace: "w", InitialPrompt: prompt, Timeout: time.Second}, func(context.Context, []string) (json.RawMessage, error) { t.Fatal("side effect"); return nil, nil }, time.Millisecond)
		if err == nil {
			t.Fatal("invalid prompt accepted")
		}
	}
}

func TestShellWaitsForForegroundJob(t *testing.T) {
	reads := 0
	run := func(_ context.Context, args []string) (json.RawMessage, error) {
		if args[1] == "current" {
			return json.RawMessage(`{"pane":{"pane_id":"p","tab_id":"t","terminal_id":"term","workspace_id":"w"}}`), nil
		}
		reads++
		group := 2
		if reads > 2 {
			group = 1
		}
		return json.RawMessage(fmt.Sprintf(`{"process_info":{"pane_id":"p","shell_pid":1,"foreground_process_group_id":%d,"foreground_processes":[{"pid":1,"name":"zsh"}]}}`, group)), nil
	}
	err := waitShell(context.Background(), pane{Pane: "p", Tab: "t", Terminal: "term", Workspace: "w"}, run, time.Millisecond)
	if err != nil || reads != 4 {
		t.Fatalf("reads=%d err=%v", reads, err)
	}
}
func TestQuoteRoundTrip(t *testing.T) {
	value := `quote ' double " dollar $(printf wrong) backtick ` + "`printf wrong`" + ` backslash \`
	out, err := exec.Command("/bin/sh", "-c", "printf '%s' "+quote(value)).Output()
	if err != nil || string(out) != value {
		t.Fatalf("%q %v", out, err)
	}
}
func TestShellWaitDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := waitShell(ctx, pane{Pane: "p", Tab: "t", Terminal: "term", Workspace: "w"}, func(_ context.Context, args []string) (json.RawMessage, error) {
		if args[1] == "current" {
			return json.RawMessage(`{"pane":{"pane_id":"p","tab_id":"t","terminal_id":"term","workspace_id":"w"}}`), nil
		}
		return json.RawMessage(`{"process_info":{"pane_id":"p"}}`), nil
	}, time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
}

func TestHerdrErrorHelper(t *testing.T) {
	if os.Getenv("PLY_LAUNCH_ERROR_HELPER") != "1" {
		return
	}
	fmt.Fprint(os.Stderr, os.Getenv("PLY_LAUNCH_ERROR_RESPONSE"))
	os.Exit(1)
}
func TestCallNeverLeaksHerdrErrorMessage(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PLY_LAUNCH_ERROR_HELPER", "1")
	for _, raw := range []string{
		`{"error":{"code":"agent_pane_busy","message":"PRIVATE_PROMPT secret_env=PRIVATE_ENV"}}`,
		`PRIVATE_PROMPT secret_env=PRIVATE_ENV`,
		`{"error":{"code":"PRIVATE_PROMPT secret_env=PRIVATE_ENV"}}`,
	} {
		t.Setenv("PLY_LAUNCH_ERROR_RESPONSE", raw)
		_, err = call(context.Background(), executable, []string{"-test.run=^TestHerdrErrorHelper$", "--"})
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("exit error not preserved: %v", err)
		}
		if strings.Contains(err.Error(), "PRIVATE_") {
			t.Fatalf("sensitive diagnostic: %v", err)
		}
		if strings.Contains(raw, `"agent_pane_busy"`) && !strings.Contains(err.Error(), "agent_pane_busy") {
			t.Fatalf("missing safe code: %v", err)
		}
	}
}

func TestCommandEstablishesCwdBeforeExecutable(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target ' space")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "fake provider")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\npwd -P\n"), 0700); err != nil {
		t.Fatal(err)
	}
	in := Request{Cwd: target, Executable: executable, InitialPrompt: "no task"}
	command := exec.Command("/bin/sh", "-c", launchCommand(in))
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	physical, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(output)) != physical {
		t.Fatalf("wrong cwd: %q", output)
	}
	in.Cwd = filepath.Join(root, "absent")
	output, err = exec.Command("/bin/sh", "-c", launchCommand(in)).Output()
	if err == nil || len(output) != 0 {
		t.Fatalf("executable ran after cd failed: %q %v", output, err)
	}
}

func TestPaneRunAcceptsEmptySuccessOnlyForSubmission(t *testing.T) {
	success, err := exec.LookPath("true")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = call(context.Background(), success, []string{"pane", "run", "p", "command"}); err != nil {
		t.Fatalf("silent successful admission rejected: %v", err)
	}
	if _, err = call(context.Background(), success, []string{"pane", "current", "--pane", "p"}); err == nil {
		t.Fatal("empty read response accepted")
	}
	failure, err := exec.LookPath("false")
	if err != nil {
		t.Fatal(err)
	}
	_, err = call(context.Background(), failure, []string{"pane", "run", "p", "command"})
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("nonzero submission accepted: %v", err)
	}
}

func observedAgentRunner(agent string) runner {
	return func(_ context.Context, args []string) (json.RawMessage, error) {
		switch strings.Join(args[:2], " ") {
		case "tab create":
			return json.RawMessage(`{"root_pane":{"pane_id":"p","tab_id":"t","workspace_id":"w","terminal_id":"term"}}`), nil
		case "pane process-info":
			return json.RawMessage(`{"process_info":{"pane_id":"p","shell_pid":1,"foreground_process_group_id":1,"foreground_processes":[{"pid":1,"name":"zsh"}]}}`), nil
		case "pane run":
			agent = "," + agent
			return json.RawMessage(`{}`), nil
		case "pane current":
			observed := ""
			if strings.HasPrefix(agent, ",") {
				observed = `,"agent":"` + agent[1:] + `","agent_status":"working","foreground_cwd":"/work"`
			}
			return json.RawMessage(`{"pane":{"pane_id":"p","tab_id":"t","workspace_id":"w","terminal_id":"term"` + observed + `}}`), nil
		}
		return nil, errors.New("unexpected")
	}
}

func TestExpectedProviderAgentObserved(t *testing.T) {
	in := Request{Agent: "codex", Executable: "/bin/codex", Cwd: "/work", Workspace: "w", InitialPrompt: "$ply-assignment ASSIGNMENT-7", Timeout: time.Second}
	r, err := launch(context.Background(), in, observedAgentRunner("codex"), time.Millisecond)
	if err != nil || r.State != "agent_observed" || r.Agent != "codex" || r.AgentStatus != "working" || r.EffectUnknown {
		t.Fatalf("%+v %v", r, err)
	}
	r, err = launch(context.Background(), in, observedAgentRunner("claude"), time.Millisecond)
	if err == nil || r.State == "agent_observed" || r.Agent != "" {
		t.Fatalf("other provider accepted: %+v %v", r, err)
	}
}

func TestMissingAgentRejectedBeforeEffect(t *testing.T) {
	_, err := launch(context.Background(), Request{Executable: "/bin/codex", Cwd: "/work", Workspace: "w", InitialPrompt: "hello", Timeout: time.Second}, func(context.Context, []string) (json.RawMessage, error) { t.Fatal("side effect"); return nil, nil }, time.Millisecond)
	if err == nil {
		t.Fatal("missing agent identity accepted")
	}
}

func TestCommandPassesSkillMentionAsOneLiteralArgument(t *testing.T) {
	root := t.TempDir()
	executable := filepath.Join(root, "fake codex")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nprintf '%s|' \"$#\" \"$@\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	prompt := "$ply-assignment ASSIGNMENT-7 --project ply"
	output, err := exec.Command("/bin/sh", "-c", launchCommand(Request{Cwd: root, Executable: executable, InitialPrompt: prompt})).Output()
	if err != nil || string(output) != "1|"+prompt+"|" {
		t.Fatalf("%q %v", output, err)
	}
}
