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
	"time"
)

type Request struct {
	Herdr, Executable, Cwd, Workspace, Label, InitialPrompt string
	Env                                                     []string
	Timeout                                                 time.Duration
}
type Result struct {
	State         string `json:"state"`
	Pane          string `json:"pane,omitempty"`
	Tab           string `json:"tab,omitempty"`
	Terminal      string `json:"terminal,omitempty"`
	Session       string `json:"session,omitempty"`
	AgentStatus   string `json:"agent_status,omitempty"`
	EffectUnknown bool   `json:"effect_unknown"`
}
type runner func(context.Context, []string) (json.RawMessage, error)
type pane struct {
	Pane      string `json:"pane_id"`
	Tab       string `json:"tab_id"`
	Terminal  string `json:"terminal_id"`
	Workspace string `json:"workspace_id"`
	Agent     string `json:"agent"`
	Status    string `json:"agent_status"`
	Cwd       string `json:"foreground_cwd"`
	Session   *struct {
		Value string `json:"value"`
	} `json:"agent_session"`
}

func Start(ctx context.Context, in Request) (Result, error) {
	if os.Getenv("HERDR_ENV") != "1" {
		return Result{}, errors.New("launch requires a Herdr-managed caller")
	}
	if !filepath.IsAbs(in.Herdr) {
		return Result{}, errors.New("Herdr executable must be absolute")
	}
	return launch(ctx, in, func(ctx context.Context, args []string) (json.RawMessage, error) { return call(ctx, in.Herdr, args) }, 200*time.Millisecond)
}
func validate(in Request) error {
	if !filepath.IsAbs(in.Executable) || !filepath.IsAbs(in.Cwd) || in.Workspace == "" || in.InitialPrompt == "" || in.Timeout <= 0 {
		return errors.New("absolute executable/cwd, workspace, prompt and positive timeout required")
	}
	for _, s := range append([]string{in.Executable, in.Cwd, in.Workspace, in.Label, in.InitialPrompt}, in.Env...) {
		if strings.ContainsAny(s, "\x00\r\n") {
			return errors.New("launch arguments must be single-line without NUL")
		}
	}
	for _, e := range in.Env {
		if !strings.Contains(e, "=") {
			return errors.New("environment must use KEY=VALUE")
		}
	}
	return nil
}
func launch(ctx context.Context, in Request, run runner, interval time.Duration) (r Result, err error) {
	r.State = "not_started"
	if err = validate(in); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, in.Timeout)
	defer cancel()
	args := []string{"tab", "create", "--workspace", in.Workspace, "--cwd", in.Cwd, "--label", in.Label, "--no-focus"}
	for _, e := range in.Env {
		args = append(args, "--env", e)
	}
	r.State = "tab_create_unknown"
	r.EffectUnknown = true
	raw, e := run(ctx, args)
	if e != nil {
		return r, e
	}
	var created struct {
		Root pane `json:"root_pane"`
	}
	if e = json.Unmarshal(raw, &created); e != nil {
		return r, e
	}
	p := created.Root
	r.Pane, r.Tab, r.Terminal = p.Pane, p.Tab, p.Terminal
	if p.Pane == "" || p.Tab == "" || p.Terminal == "" || p.Workspace != in.Workspace {
		return r, errors.New("invalid created pane identity")
	}
	r.State = "tab_created"
	r.EffectUnknown = false
	if e = waitShell(ctx, p, run, interval); e != nil {
		return r, e
	}
	command := launchCommand(in)
	r.State = "submission_unknown"
	r.EffectUnknown = true
	if _, e = run(ctx, []string{"pane", "run", p.Pane, command}); e != nil {
		return r, e
	}
	r.State = "command_submitted"
	for {
		observed, e := readPane(ctx, p, run)
		if e != nil {
			return r, e
		}
		if observed.Agent != "" {
			if observed.Agent != "claude" || observed.Cwd != in.Cwd {
				return r, errors.New("unexpected agent or foreground directory")
			}
			r.AgentStatus = observed.Status
			if observed.Session != nil {
				r.Session = observed.Session.Value
			}
			r.State = "agent_observed"
			r.EffectUnknown = false
			return r, nil
		}
		if e = pause(ctx, interval); e != nil {
			return r, e
		}
	}
}
func readPane(ctx context.Context, expected pane, run runner) (pane, error) {
	raw, err := run(ctx, []string{"pane", "current", "--pane", expected.Pane})
	if err != nil {
		return pane{}, err
	}
	var result struct {
		Pane pane `json:"pane"`
	}
	if err = json.Unmarshal(raw, &result); err != nil {
		return pane{}, err
	}
	p := result.Pane
	if p.Pane != expected.Pane || p.Terminal != expected.Terminal || p.Tab != expected.Tab || p.Workspace != expected.Workspace {
		return p, errors.New("pane identity changed")
	}
	return p, nil
}
func waitShell(ctx context.Context, p pane, run runner, interval time.Duration) error {
	consecutive := 0
	for {
		current, err := readPane(ctx, p, run)
		if err != nil {
			return err
		}
		if current.Agent != "" {
			return errors.New("created pane already contains an agent")
		}
		raw, err := run(ctx, []string{"pane", "process-info", "--pane", p.Pane})
		if err != nil {
			return err
		}
		ready, err := shellAvailable(raw, p.Pane)
		if err != nil {
			return err
		}
		if ready {
			consecutive++
		} else {
			consecutive = 0
		}
		if consecutive >= 2 {
			return nil
		}
		if err = pause(ctx, interval); err != nil {
			return err
		}
	}
}
func shellAvailable(raw json.RawMessage, id string) (bool, error) {
	var result struct {
		Info struct {
			Pane      string `json:"pane_id"`
			PID       int    `json:"shell_pid"`
			Group     int    `json:"foreground_process_group_id"`
			Processes []struct {
				PID  int    `json:"pid"`
				Name string `json:"name"`
			} `json:"foreground_processes"`
		} `json:"process_info"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return false, err
	}
	p := result.Info
	if p.Pane != id {
		return false, errors.New("process pane identity changed")
	}
	if p.PID == 0 || p.Group != p.PID || len(p.Processes) != 1 || p.Processes[0].PID != p.PID {
		return false, nil
	}
	switch strings.TrimPrefix(filepath.Base(p.Processes[0].Name), "-") {
	case "sh", "bash", "zsh":
		return true, nil
	}
	return false, errors.New("pilot supports sh, bash and zsh shell quoting only")
}
func pause(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func quote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }

type bounded struct {
	data     []byte
	overflow bool
}

func (b *bounded) Write(p []byte) (int, error) {
	n := len(p)
	remaining := (1 << 20) - len(b.data)
	if len(p) > remaining {
		b.overflow = true
		p = p[:remaining]
	}
	b.data = append(b.data, p...)
	return n, nil
}
func call(ctx context.Context, binary string, args []string) (json.RawMessage, error) {
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.WaitDelay = time.Second
	var out, stderr bounded
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("herdr %s %s: %w; code: %s", args[0], args[1], err, diagnostic(stderr.data))
	}
	if out.overflow || stderr.overflow {
		return nil, errors.New("Herdr output exceeded limit")
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if err := json.Unmarshal(out.data, &envelope); err != nil {
		return nil, errors.New("invalid Herdr JSON")
	}
	if len(envelope.Error) > 0 && string(envelope.Error) != "null" {
		return nil, errors.New("Herdr returned an error")
	}
	if len(envelope.Result) == 0 || string(envelope.Result) == "null" {
		return nil, errors.New("missing Herdr result")
	}
	return envelope.Result, nil
}

func diagnostic(raw []byte) string {
	var envelope struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &envelope) != nil {
		return "unavailable"
	}
	code := envelope.Error.Code
	if len(code) == 0 || len(code) > 64 || code[0] < 'a' || code[0] > 'z' {
		return "unavailable"
	}
	for _, c := range code {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_') {
			return "unavailable"
		}
	}
	return code
}

func launchCommand(in Request) string {
	return "cd -- " + quote(in.Cwd) + " && exec " + quote(in.Executable) + " " + quote(in.InitialPrompt)
}
