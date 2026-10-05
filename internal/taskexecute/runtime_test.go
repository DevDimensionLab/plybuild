package taskexecute

import (
	"os"
	"path/filepath"
	"strings"
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
	if x.Runtime.Provider != "claude" || x.Policy.PermissionProfile != "auto" || x.Runtime.ConfigProfile != nil || x.Policy.ApprovalPolicy != "" {
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

func TestClaudeRuntimeRequestsAutoWithoutChangingProviderSettings(t *testing.T) {
	root := physicalTemp(t)
	bin := filepath.Join(root, "fake-executable")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	settingsPath := filepath.Join(root, "settings.json")
	settings := []byte(`{"permissions":{"defaultMode":"manual","deny":["Bash(git push *)"]}}`)
	if err := os.WriteFile(settingsPath, settings, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", root)
	assignment := workspace.TaskExecutorAssignment{Provider: ptrString("claude")}
	for _, requested := range []string{"", "auto"} {
		t.Run("profile="+requested, func(t *testing.T) {
			opts := RuntimeOptions{ProviderExecutable: bin, HerdrExecutable: bin, HerdrWorkspace: "w-test", PermissionProfile: requested}
			x, err := PreviewRuntime(assignment, bin, opts)
			if err != nil {
				t.Fatal(err)
			}
			if x.Runtime.Provider != "claude" || x.Runtime.Model != "opus" || x.ReasoningEffort != "medium" || x.Policy.PermissionProfile != "auto" || x.Runtime.PermissionBinding.ProfileID != "auto" {
				t.Fatalf("new Claude launch did not request Opus/medium with auto mode: %+v", x)
			}
			if x.Runtime.PermissionBinding.AuthorityKind != "launch_contract_pending_runtime_acceptance" || x.Runtime.PermissionBinding.EffectivePolicySHA256 != "" || len(x.Runtime.PermissionBinding.Evidence) != 0 {
				t.Fatalf("requested auto mode was presented as observed runtime authority: %+v", x.Runtime.PermissionBinding)
			}
			if x.Runtime.ConfigProfile != nil || x.Policy.ApprovalPolicy != "" || x.Policy.ApprovalReviewer != "" || len(x.Policy.Sources) != 0 {
				t.Fatalf("Codex permission configuration leaked into Claude: %+v", x)
			}
		})
	}
	for _, requested := range []string{"manual", "local-codex-profile", "bypassPermissions"} {
		opts := RuntimeOptions{ProviderExecutable: bin, HerdrExecutable: bin, HerdrWorkspace: "w-test", PermissionProfile: requested}
		if _, err := PreviewRuntime(assignment, bin, opts); err == nil || !strings.Contains(err.Error(), "native auto permission mode") {
			t.Errorf("unsupported Claude launch profile %q was not rejected clearly: %v", requested, err)
		}
	}
	if after, err := os.ReadFile(settingsPath); err != nil || string(after) != string(settings) {
		t.Fatal("runtime preview changed provider permission configuration")
	}
	if _, err := os.Stat(filepath.Join(root, "launch-policy.json")); !os.IsNotExist(err) {
		t.Fatal("runtime preview wrote a launch policy")
	}
}
