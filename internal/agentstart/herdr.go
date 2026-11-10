package agentstart

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	herdrCallTimeout = 40 * time.Second
	herdrOutputLimit = 1 << 20
)

type limitedOutput struct {
	data     []byte
	overflow bool
}

func (b *limitedOutput) Write(p []byte) (int, error) {
	n := len(p)
	left := herdrOutputLimit - len(b.data)
	if len(p) > left {
		b.overflow = true
		p = p[:left]
	}
	b.data = append(b.data, p...)
	return n, nil
}

func fileSHA256(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func ResolveHerdr() (HerdrBinary, error) {
	found, err := exec.LookPath("herdr")
	if err == nil {
		found, err = filepath.EvalSymlinks(found)
	}
	if err == nil {
		found, err = filepath.Abs(found)
	}
	if err != nil {
		return HerdrBinary{}, fmt.Errorf("herdr executable not found: %w", err)
	}
	sum, err := fileSHA256(found)
	if err != nil {
		return HerdrBinary{}, err
	}
	return HerdrBinary{Path: found, SHA256: sum}, nil
}

func callHerdr(ctx context.Context, h HerdrBinary, args []string) (json.RawMessage, error) {
	sum, err := fileSHA256(h.Path)
	if err != nil || sum != h.SHA256 {
		return nil, &Error{Code: ClassBinaryChanged, Err: errors.New("herdr executable differs from the bound digest")}
	}
	ctx, cancel := context.WithTimeout(ctx, herdrCallTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, h.Path, args...)
	cmd.WaitDelay = time.Second
	var stdout, stderr limitedOutput
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("herdr %s %s: %w", args[0], subcommand(args), err)
	}
	if stdout.overflow || stderr.overflow {
		return nil, fmt.Errorf("herdr %s %s: output limit exceeded", args[0], subcommand(args))
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if err = json.Unmarshal(stdout.data, &envelope); err != nil {
		return nil, fmt.Errorf("herdr %s %s: invalid response", args[0], subcommand(args))
	}
	if hasValue(envelope.Error) || !hasValue(envelope.Result) {
		return nil, fmt.Errorf("herdr %s %s: unsuccessful response", args[0], subcommand(args))
	}
	return envelope.Result, nil
}

func subcommand(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	return ""
}

func hasValue(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && string(trimmed) != "null"
}

func agentName(key string) string {
	sum := sha256.Sum256([]byte(key))
	return "ply-" + hex.EncodeToString(sum[:])[:28]
}

func tabArgs(in StartInput) []string {
	args := []string{"tab", "create", "--workspace", in.Workspace, "--cwd", in.Cwd, "--label", in.Label}
	for _, env := range in.Env {
		args = append(args, "--env", env)
	}
	return append(args, "--no-focus")
}

func startArgs(in StartInput, pane string) []string {
	args := []string{"agent", "start", agentName(in.Key), "--kind", in.Provider, "--pane", pane, "--timeout", "30000"}
	var flags []string
	if in.Model != nil {
		flags = append(flags, "--model", *in.Model)
	}
	if in.Provider == ProviderClaude {
		if in.Permission != nil {
			flags = append(flags, "--permission-mode", *in.Permission)
		}
	} else {
		if in.Permission != nil {
			flags = append(flags, "-c", "default_permissions="+quoted(*in.Permission))
		}
		if in.CodexProfile != nil {
			flags = append(flags, "--profile", *in.CodexProfile)
		}
	}
	if len(flags) > 0 {
		args = append(args, "--")
		args = append(args, flags...)
	}
	return args
}

func promptArgs(pane, prompt string) []string {
	return []string{"agent", "prompt", pane, prompt, "--wait", "--until", "working", "--timeout", "10000"}
}

func quoted(s string) string {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(s)
	return strings.TrimSuffix(out.String(), "\n")
}
