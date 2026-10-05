package taskexecute

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func ptrString(s string) *string { return &s }

func TestRuntimeAssignmentAndMissingBinary(t *testing.T) {
	root := physicalTemp(t)
	bin := filepath.Join(root, "fake-executable")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config.toml"), []byte("model = 'configured-codex'\nmodel_reasoning_effort = 'high'\ndefault_permissions = 'local'\n[permissions.local]\nnetwork = false\n"), 0600); err != nil {
		t.Fatal(err)
	}
	opts := RuntimeOptions{ProviderExecutable: bin, HerdrExecutable: bin, CodexHome: root, HerdrWorkspace: "w-test"}
	x, err := PreviewRuntime(workspace.TaskExecutorAssignment{}, bin, opts)
	if err != nil {
		t.Fatal(err)
	}
	if x.Runtime.Provider != "codex" || x.Runtime.Model != "configured-codex" || x.ReasoningEffort != "high" || x.Policy.PermissionProfile != "local" {
		t.Fatalf("wrong default: %+v", x)
	}
	claude := workspace.TaskExecutorAssignment{Provider: ptrString("claude"), Model: ptrString("opus"), Effort: ptrString("medium")}
	x, err = PreviewRuntime(claude, bin, opts)
	if err != nil {
		t.Fatal(err)
	}
	if x.Runtime.Provider != "claude" || x.Policy.PermissionProfile != "manual" || x.Runtime.ConfigProfile != nil || x.Policy.ApprovalPolicy != "" {
		t.Fatalf("Codex policy leaked: %+v", x)
	}
	opts.ProviderExecutable = filepath.Join(root, "missing")
	if _, err = PreviewRuntime(claude, bin, opts); err == nil {
		t.Fatal("missing Claude accepted")
	}
	claude.Provider = ptrString("unsupported")
	if _, err = PreviewRuntime(claude, bin, opts); err == nil {
		t.Fatal("unknown implementor fell back")
	}
	if _, err = os.Stat(filepath.Join(root, "launch-policy.json")); !os.IsNotExist(err) {
		t.Fatal("preview wrote policy")
	}
}
