package taskrun

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHerdrDiagnosticsPaneNotFoundOnStderrIsReadFailure(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "herdr-standin")
	// The real installed CLI reported this JSON envelope on stderr with exit 1.
	// No actual Herdr pane, process observer or workflow state is used here.
	script := []byte("#!/bin/sh\nprintf '%s\\n' '{\"error\":{\"code\":\"pane_not_found\",\"message\":\"pane w1:p1C not found\"},\"id\":\"cli:pane:get\"}' >&2\nexit 1\n")
	if err = os.WriteFile(path, script, 0700); err != nil {
		t.Fatal(err)
	}
	var request WorkflowRequest
	request.Herdr.Executable = Executable{Path: path, SHA256: hash(script)}
	_, err = workflowCall(Dependencies{HerdrTimeout: time.Second}, request, "pane", "get", "w1:p1C")
	if err == nil {
		t.Fatal("missing pane was accepted")
	}
	var diagnostic *HerdrCallError
	var failure *Error
	if !errors.As(err, &diagnostic) || diagnostic.Phase != "pane get" || diagnostic.Class != "nonzero_exit" || diagnostic.ProviderCode != "pane_not_found" || !diagnostic.ReadOnly {
		t.Fatalf("known read failure lost typed classification: %#v", diagnostic)
	}
	if !errors.As(err, &failure) || failure.Code != "workflow_run" || failure.Exit != 5 {
		t.Fatalf("typed diagnostic changed the existing CLI error contract: %#v", failure)
	}
	for _, want := range []string{"Herdr pane get failed", "class=nonzero_exit", "exit=1", "provider_code=pane_not_found", "observation unavailable", "no retry was sent"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("missing %q in %q", want, err)
		}
	}
	for _, forbidden := range []string{"effect and readiness unknown", "w1:p1C", "cli:pane:get"} {
		if strings.Contains(err.Error(), forbidden) {
			t.Errorf("read diagnostic includes unsupported or raw content %q: %v", forbidden, err)
		}
	}
}

func TestHerdrDiagnosticsPreserveSafeCodesAndOperationKind(t *testing.T) {
	agent := `{"error":{"code":"agent_not_found","message":"synthetic-secret"}}`
	pane := `{"error":{"code":"pane_not_found","message":"synthetic-secret"}}`
	for _, tc := range []struct {
		name     string
		args     []string
		stdout   string
		stderr   string
		phase    string
		code     string
		readOnly bool
	}{
		{"agent_stdout", []string{"agent", "get", "synthetic-secret"}, agent, "", "agent get", "agent_not_found", true},
		{"pane_stderr", []string{"pane", "get", "synthetic-secret"}, "", pane, "pane get", "pane_not_found", true},
		{"pane_inventory", []string{"pane", "list"}, "", pane, "pane list", "pane_not_found", true},
		{"matching_streams", []string{"pane", "process-info", "--pane", "synthetic-secret"}, pane, pane, "pane process-info", "pane_not_found", true},
		{"bounded_metadata", []string{"pane", "get", "synthetic-secret"}, `{"elapsed_seconds":0.1}`, pane, "pane get", "pane_not_found", true},
		{"start_remains_uncertain", []string{"agent", "start", "synthetic-secret"}, "", pane, "agent start", "pane_not_found", false},
		{"prompt_remains_uncertain", []string{"agent", "prompt", "synthetic-secret"}, agent, "", "agent prompt", "agent_not_found", false},
		{"pane_run_remains_uncertain", []string{"pane", "run", "synthetic-secret"}, "", pane, "pane run", "pane_not_found", false},
		{"tab_create_remains_uncertain", []string{"tab", "create", "synthetic-secret"}, "", pane, "tab create", "pane_not_found", false},
		{"unknown_argv", []string{"synthetic-secret", "synthetic-secret"}, "", pane, "transport", "pane_not_found", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr workflowLimitedOutput
			_, _ = stdout.Write([]byte(tc.stdout))
			_, _ = stderr.Write([]byte(tc.stderr))
			err := workflowCallFailure(tc.args, nil, nil, stdout, stderr, "unsuccessful_envelope")
			var diagnostic *HerdrCallError
			if !errors.As(err, &diagnostic) || diagnostic.Phase != tc.phase || diagnostic.ProviderCode != tc.code || diagnostic.ReadOnly != tc.readOnly {
				t.Fatalf("safe classification lost: %#v", diagnostic)
			}
			if strings.Contains(err.Error(), "synthetic-secret") || len(err.Error()) > 1024 {
				t.Fatalf("diagnostic retained untrusted content: %v", err)
			}
			if strings.Contains(err.Error(), "effect and readiness unknown") == tc.readOnly || !strings.Contains(err.Error(), "no retry was sent") {
				t.Fatalf("call kind has misleading effect wording: %v", err)
			}
		})
	}
}

func TestHerdrDiagnosticsDoNotPromoteAmbiguousProviderOutput(t *testing.T) {
	known := `{"error":{"code":"pane_not_found","message":"synthetic-secret"}}`
	for _, tc := range []struct {
		name   string
		stdout string
		stderr string
	}{
		{"conflicting_streams", known, `{"error":{"code":"agent_not_found"}}`},
		{"invalid_other_stream", known, "synthetic-secret"},
		{"unsupported_other_code", known, `{"error":{"code":"synthetic-secret"}}`},
		{"duplicate_code", `{"error":{"code":"synthetic-secret","code":"pane_not_found"}}`, ""},
		{"duplicate_envelope", `{"error":{"code":"agent_not_found"},"error":{"code":"pane_not_found"}}`, ""},
		{"incomplete_json", "", `{"error":{"code":"pane_not_found"}`},
		{"success_conflicts_with_error", `{"result":{"pane":{}}}`, known},
		{"invalid_utf8", "", "{\"error\":{\"code\":\"pane_not_found\",\"message\":\"\xff\"}}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr workflowLimitedOutput
			_, _ = stdout.Write([]byte(tc.stdout))
			_, _ = stderr.Write([]byte(tc.stderr))
			err := workflowCallFailure([]string{"pane", "get", "synthetic-secret"}, nil, nil, stdout, stderr, "unsuccessful_envelope")
			if !strings.Contains(err.Error(), "provider_code=unknown") || strings.Contains(err.Error(), "synthetic-secret") {
				t.Fatalf("ambiguous or raw provider output was promoted: %v", err)
			}
		})
	}
}

func TestHerdrDiagnosticsIncompleteTransportCannotProveKnownCode(t *testing.T) {
	known := `{"error":{"code":"pane_not_found","message":"synthetic-secret"}}`
	for _, tc := range []struct {
		name       string
		contextErr error
		processErr error
		class      string
		truncate   string
	}{
		{"timeout", context.DeadlineExceeded, nil, "unsuccessful_envelope", ""},
		{"process_failure", nil, errors.New("synthetic-secret"), "", ""},
		{"invalid_response", nil, nil, "invalid_response", ""},
		{"output_limit", nil, nil, "output_limit", ""},
		{"truncated_stdout", nil, nil, "unsuccessful_envelope", "stdout"},
		{"truncated_stderr", nil, nil, "unsuccessful_envelope", "stderr"},
		{"unknown_class", nil, nil, "synthetic-secret", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr workflowLimitedOutput
			_, _ = stdout.Write([]byte(known))
			_, _ = stderr.Write([]byte(known))
			stdout.overflow, stderr.overflow = tc.truncate == "stdout", tc.truncate == "stderr"
			err := workflowCallFailure([]string{"pane", "get", "synthetic-secret"}, tc.contextErr, tc.processErr, stdout, stderr, tc.class)
			var diagnostic *HerdrCallError
			if !errors.As(err, &diagnostic) || diagnostic.ProviderCode != "unknown" || strings.Contains(err.Error(), "synthetic-secret") {
				t.Fatalf("incomplete output became absence evidence: %v", err)
			}
		})
	}
	if err := (&HerdrCallError{}).Unwrap(); err != nil {
		t.Fatalf("empty diagnostic unwraps to a non-nil error: %#v", err)
	}
}
