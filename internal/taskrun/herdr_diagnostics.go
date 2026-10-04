package taskrun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

// Child output is untrusted and may contain capabilities or terminal contents.
// Retain only fixed classifications, exit status, sizes and known provider codes.
// Never include argv payloads, raw streams, error messages or arbitrary codes.
func workflowCallFailure(args []string, contextErr, processErr error, stdout, stderr workflowLimitedOutput, class string) error {
	phase := "transport"
	if len(args) >= 2 {
		switch args[0] + " " + args[1] {
		case "tab create", "tab rename", "agent start", "agent get", "agent prompt":
			phase = args[0] + " " + args[1]
		}
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
	var envelope struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if !stdout.overflow && json.Unmarshal(stdout.data, &envelope) == nil {
		switch envelope.Error.Code {
		case "agent_not_found":
			code = envelope.Error.Code
		}
	}
	return workflowError(5, fmt.Sprintf("Herdr %s failed: class=%s exit=%s provider_code=%s stdout_bytes=%d stderr_bytes=%d truncated=%t; effect and readiness unknown; no retry was sent", phase, class, exit, code, stdout.bytes, stderr.bytes, stdout.overflow || stderr.overflow))
}
