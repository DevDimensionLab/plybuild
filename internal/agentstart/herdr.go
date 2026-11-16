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
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	herdrCallTimeout = 40 * time.Second
	herdrOutputLimit = 1 << 20
	// herdrMessageLimit bounds the Herdr error message kept in Ply's error text.
	herdrMessageLimit = 160
)

var herdrErrorCode = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

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
	out, err := runHerdr(ctx, h, args)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if err = json.Unmarshal(out, &envelope); err != nil {
		return nil, fmt.Errorf("herdr %s %s: invalid response", args[0], subcommand(args))
	}
	if hasValue(envelope.Error) || !hasValue(envelope.Result) {
		return nil, refusal(out, fmt.Errorf("herdr %s %s: unsuccessful response%s", args[0], subcommand(args), herdrErrorDetail(out, args)))
	}
	return envelope.Result, nil
}

// callHerdrText runs a Herdr command that replies with plain text, such as
// pane read and pane send-keys.
func callHerdrText(ctx context.Context, h HerdrBinary, args []string) (string, error) {
	out, err := runHerdr(ctx, h, args)
	return string(out), err
}

func runHerdr(ctx context.Context, h HerdrBinary, args []string) ([]byte, error) {
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
		if ctx.Err() != nil || stderr.overflow {
			return nil, fmt.Errorf("herdr %s %s: %w", args[0], subcommand(args), err)
		}
		return nil, refusal(stderr.data, fmt.Errorf("herdr %s %s: %w%s", args[0], subcommand(args), err, herdrErrorDetail(stderr.data, args)))
	}
	if stdout.overflow || stderr.overflow {
		return nil, fmt.Errorf("herdr %s %s: output limit exceeded", args[0], subcommand(args))
	}
	return stdout.data, nil
}

// herdrRefusal is a failed Herdr call whose reply is exactly a Herdr JSON
// error with a well-formed code.
type herdrRefusal struct {
	code string
	err  error
}

func (e *herdrRefusal) Error() string { return e.err.Error() }
func (e *herdrRefusal) Unwrap() error { return e.err }

// refusal wraps err with the code of the Herdr JSON error in output, if any.
func refusal(output []byte, err error) error {
	if code, _, ok := herdrErrorReply(output); ok {
		return &herdrRefusal{code: code, err: err}
	}
	return err
}

// herdrCode returns the Herdr error code carried by err, or "".
func herdrCode(err error) string {
	var refused *herdrRefusal
	if errors.As(err, &refused) {
		return refused.code
	}
	return ""
}

func herdrErrorReply(output []byte) (code, message string, ok bool) {
	var reply struct {
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(output, &reply) != nil || reply.Error == nil || !herdrErrorCode.MatchString(reply.Error.Code) {
		return "", "", false
	}
	return reply.Error.Code, reply.Error.Message, true
}

// herdrErrorDetail renders the code and message of a Herdr JSON error reply as
// " (herdr error code: message)", or "" when the reply is not exactly such an
// object. The message is bounded and withheld for calls carrying a sensitive
// payload; no other output text is kept.
func herdrErrorDetail(output []byte, args []string) string {
	code, message, ok := herdrErrorReply(output)
	if !ok {
		return ""
	}
	if withholdsMessage(args) {
		message = "message withheld"
	}
	if message = boundedMessage(message); message == "" {
		return " (herdr error " + code + ")"
	}
	return " (herdr error " + code + ": " + message + ")"
}

// withholdsMessage reports whether the call carries a payload that a Herdr
// message might echo: the prompt of agent prompt, or environment values.
func withholdsMessage(args []string) bool {
	return len(args) > 1 && args[0] == "agent" && args[1] == "prompt" || slices.Contains(args, "--env")
}

func boundedMessage(message string) string {
	message = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, strings.ToValidUTF8(message, "\uFFFD"))
	message = strings.TrimSpace(message)
	if len(message) <= herdrMessageLimit {
		return message
	}
	cut := herdrMessageLimit
	for cut > 0 && !utf8.RuneStart(message[cut]) {
		cut--
	}
	return message[:cut] + "…"
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
		if in.Effort != nil {
			flags = append(flags, "--effort", *in.Effort)
		}
		if in.Permission != nil {
			flags = append(flags, "--permission-mode", *in.Permission)
		}
	} else {
		if in.Effort != nil {
			flags = append(flags, "-c", "model_reasoning_effort="+quoted(*in.Effort))
		}
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
