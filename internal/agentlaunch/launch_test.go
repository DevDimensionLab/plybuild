package agentlaunch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
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
	r, err := launch(context.Background(), Request{Executable: "/bin/claude", Cwd: "/work", Workspace: "w1", InitialPrompt: "hello '$()", Timeout: time.Second}, run, time.Millisecond)
	if err != nil || r.State != "agent_observed" || r.AgentStatus != "working" || calls != 1 || polls < 2 {
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
	r, err := launch(context.Background(), Request{Executable: "/bin/claude", Cwd: "/work", Workspace: "w", InitialPrompt: "hello", Timeout: time.Second}, run, time.Millisecond)
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
		_, err := launch(context.Background(), Request{Executable: "/bin/claude", Cwd: "/work", Workspace: "w", InitialPrompt: prompt, Timeout: time.Second}, func(context.Context, []string) (json.RawMessage, error) { t.Fatal("side effect"); return nil, nil }, time.Millisecond)
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
