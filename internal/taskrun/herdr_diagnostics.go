package taskrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// HerdrCallError contains only bounded, allowlisted transport facts. A missing
// pane is a failed read observation, not evidence that its provider exited.
// Unwrap preserves the existing workflow error and command exit classification.
type HerdrCallError struct {
	Phase        string
	Class        string
	ProviderCode string
	ReadOnly     bool
	cause        *Error
}

func (e *HerdrCallError) Error() string {
	if e == nil || e.cause == nil {
		return "workflow_run: Herdr transport failure"
	}
	return e.cause.Error()
}

func (e *HerdrCallError) Unwrap() error {
	if e == nil || e.cause == nil {
		return nil
	}
	return e.cause
}

// Child output is untrusted and may contain capabilities or terminal contents.
// Retain only fixed classifications, exit status, sizes and known provider codes.
// Never include argv payloads, raw streams, error messages or arbitrary codes.
func workflowCallFailure(args []string, contextErr, processErr error, stdout, stderr workflowLimitedOutput, class string) error {
	phase := "transport"
	readOnly := false
	if len(args) >= 2 {
		switch args[0] + " " + args[1] {
		case "agent get", "pane get", "pane list", "pane process-info":
			phase, readOnly = args[0]+" "+args[1], true
		case "tab create", "tab rename", "agent start", "agent prompt", "pane run":
			phase = args[0] + " " + args[1]
		}
	}
	switch class {
	case "output_limit", "invalid_response", "unsuccessful_envelope":
	default:
		class = "unknown"
	}
	exit := "0"
	if processErr != nil {
		exit, class = "unknown", "process_error"
		var child *exec.ExitError
		if errors.As(processErr, &child) {
			class = "nonzero_exit"
			if child.ExitCode() >= 0 {
				exit = strconv.Itoa(child.ExitCode())
			}
		}
	}
	if errors.Is(contextErr, context.DeadlineExceeded) {
		class = "timeout"
	}
	code := "unknown"
	if class == "nonzero_exit" || class == "unsuccessful_envelope" {
		code = workflowFailureProviderCode(stdout, stderr)
	}
	consequence := "effect and readiness unknown; no retry was sent"
	if readOnly {
		consequence = "observation unavailable; no retry was sent"
	}
	cause := &Error{Code: "workflow_run", Exit: 5, Detail: fmt.Sprintf("Herdr %s failed: class=%s exit=%s provider_code=%s stdout_bytes=%d stderr_bytes=%d truncated=%t; %s", phase, class, exit, code, stdout.bytes, stderr.bytes, stdout.overflow || stderr.overflow, consequence)}
	return &HerdrCallError{Phase: phase, Class: class, ProviderCode: code, ReadOnly: readOnly, cause: cause}
}

func workflowFailureProviderCode(stdout, stderr workflowLimitedOutput) string {
	if stdout.overflow || stderr.overflow {
		return "unknown"
	}
	code := ""
	for _, stream := range []workflowLimitedOutput{stdout, stderr} {
		data := bytes.TrimSpace(stream.data)
		if len(data) == 0 {
			continue
		}
		// Reuse the strict parser to reject duplicate members, invalid Unicode,
		// or trailing JSON. Numeric metadata is validated but never retained.
		if _, err := canonicaljson.DecodeStrictWithNumberPolicy(data, func(_ []string, number string) (float64, error) {
			return strconv.ParseFloat(number, 64)
		}); err != nil {
			return "unknown"
		}
		var envelope map[string]json.RawMessage
		if err := json.Unmarshal(data, &envelope); err != nil || envelope == nil {
			return "unknown"
		}
		if result, ok := envelope["result"]; ok && !bytes.Equal(bytes.TrimSpace(result), []byte("null")) {
			return "unknown"
		}
		raw, ok := envelope["error"]
		if !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			continue
		}
		var failure struct {
			Code string `json:"code"`
		}
		if err := json.Unmarshal(raw, &failure); err != nil {
			return "unknown"
		}
		switch failure.Code {
		case "agent_not_found", "pane_not_found":
		default:
			return "unknown"
		}
		if code != "" && code != failure.Code {
			return "unknown"
		}
		code = failure.Code
	}
	if code == "" {
		return "unknown"
	}
	return code
}
