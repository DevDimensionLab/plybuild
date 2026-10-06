package taskexecute

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func recordedClaudeTrust(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err = json.Unmarshal(b, &record); err != nil {
		t.Fatal(err)
	}
	trust, _ := record["claude_project_trust"].(map[string]any)
	return trust
}

func TestExecuteClaudeProjectTrustIsExplicitAndPreserved(t *testing.T) {
	d, in, root, epic := launcherFixture(t)
	configDir := filepath.Join(root, "claude-config")
	if err := os.Mkdir(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, ".claude.json")
	configBefore := []byte("{\"projects\":{},\"unrelated\":9007199254740993}\n")
	if err := os.WriteFile(configPath, configBefore, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)
	before := fixtureGit(t, epic, "worktree", "list", "--porcelain")
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return errors.New("fixture stopped with preserved trust intent")
		}
		return nil
	}
	first, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture stopped with preserved trust intent") {
		t.Fatalf("did not preserve the intent before preparation: %v", err)
	}
	intentPath := filepath.Join(filepath.Dir(first.RequestPath), "intent.json")
	grant := recordedClaudeTrust(t, intentPath)
	if grant["mode"] != "configuration" || grant["worktree_root"] != first.Goal.WorktreePath {
		t.Fatalf("new Claude intent did not bind explicit trust to its planned worktree: %+v", grant)
	}
	intentBefore, err := os.ReadFile(intentPath)
	if err != nil {
		t.Fatal(err)
	}
	in.Check = true
	checked, err := Execute(d, in)
	if err != nil || checked.State != "ready" {
		t.Fatalf("preserved intent was not previewable: state=%s error=%v", checked.State, err)
	}
	if got := fixtureGit(t, epic, "worktree", "list", "--porcelain"); got != before {
		t.Fatal("checking a preserved intent created a worktree")
	}
	if after, err := os.ReadFile(intentPath); err != nil || !bytes.Equal(after, intentBefore) {
		t.Fatal("checking a preserved intent rewrote it")
	}
	reservations := 0
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			reservations++
			return errors.New("fixture stopped before native transport")
		}
		return nil
	}
	in.Check = false
	started, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture stopped before native transport") || reservations != 1 {
		t.Fatalf("preserved trust intent did not reach its first reservation: calls=%d error=%v", reservations, err)
	}
	requestGrant := recordedClaudeTrust(t, started.RequestPath)
	if requestGrant["mode"] != grant["mode"] || requestGrant["worktree_root"] != grant["worktree_root"] {
		t.Fatalf("request did not preserve its original explicit trust: intent=%+v request=%+v", grant, requestGrant)
	}
	if after, err := os.ReadFile(configPath); err != nil || !bytes.Equal(after, configBefore) {
		t.Fatal("execute wrote Claude configuration before transport")
	}
	if after, err := os.ReadFile(intentPath); err != nil || !bytes.Equal(after, intentBefore) {
		t.Fatal("starting from a preserved intent rewrote it")
	}
}

func TestExecuteClaudeProjectTrustPreviewIsReadOnly(t *testing.T) {
	d, in, root, epic := launcherFixture(t)
	configDir := filepath.Join(root, "claude-config")
	if err := os.Mkdir(configDir, 0700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, ".claude.json")
	configBefore := []byte("{\"projects\":{}}\n")
	if err := os.WriteFile(configPath, configBefore, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", configDir)
	registryBefore, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	worktreesBefore := fixtureGit(t, epic, "worktree", "list", "--porcelain")
	in.Check = true
	preview, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	trust := preview.ClaudeProjectTrust
	if trust == nil || trust.WorktreeRoot != preview.Goal.WorktreePath || trust.ProjectKey != preview.Goal.WorktreePath || trust.ConfigPath != configPath || trust.ConfigDir != configDir || trust.Reason != "" {
		t.Fatalf("check did not expose the exact worktree trust and caller config effect: %+v", trust)
	}
	registryAfter, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil || registryAfter.RawSHA256 != registryBefore.RawSHA256 {
		t.Fatal("trust preview changed the workspace registry")
	}
	if after := fixtureGit(t, epic, "worktree", "list", "--porcelain"); after != worktreesBefore {
		t.Fatal("trust preview created a worktree")
	}
	if _, err = os.Stat(filepath.Dir(preview.Goal.AcceptancePath)); !os.IsNotExist(err) {
		t.Fatal("trust preview wrote launch artifacts")
	}
	if after, err := os.ReadFile(configPath); err != nil || !bytes.Equal(after, configBefore) {
		t.Fatal("trust preview changed Claude configuration")
	}
}

func TestExecuteRejectsInvalidPreservedClaudeTrustBeforePreparation(t *testing.T) {
	for _, test := range []struct {
		name     string
		provider string
		grant    any
	}{
		{name: "null", provider: "claude", grant: nil},
		{name: "unknown mode", provider: "claude", grant: map[string]any{"mode": "implicit"}},
		{name: "wrong worktree", provider: "claude", grant: map[string]any{"mode": "configuration", "worktree_root": "/wrong/task"}},
		{name: "wrong provider", provider: "codex", grant: map[string]any{"mode": "configuration"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			d, in, _, epic := launcherFixtureForProvider(t, test.provider)
			d.Fault = func(point string) error {
				if point == "execute_before_preparation" {
					return errors.New("fixture stopped before preparation")
				}
				return nil
			}
			first, err := Execute(d, in)
			if err == nil || !strings.Contains(err.Error(), "fixture stopped before preparation") {
				t.Fatalf("did not preserve fixture intent: %v", err)
			}
			path := filepath.Join(filepath.Dir(first.RequestPath), "intent.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var record map[string]json.RawMessage
			if err = json.Unmarshal(raw, &record); err != nil {
				t.Fatal(err)
			}
			if grant, ok := test.grant.(map[string]any); ok {
				if _, exists := grant["worktree_root"]; !exists {
					grant["worktree_root"] = first.Goal.WorktreePath
				}
			}
			record["claude_project_trust"], err = json.Marshal(test.grant)
			if err != nil {
				t.Fatal(err)
			}
			fixtureJSON(t, path, record)
			before := fixtureGit(t, epic, "worktree", "list", "--porcelain")
			in.Check = true
			if _, err = Execute(d, in); err == nil {
				t.Fatal("malformed or foreign trust was accepted")
			}
			in.Check = false
			if _, err = Execute(d, in); err == nil || strings.Contains(err.Error(), "fixture stopped before preparation") {
				t.Fatalf("invalid trust reached preparation: %v", err)
			}
			if after := fixtureGit(t, epic, "worktree", "list", "--porcelain"); after != before {
				t.Fatal("invalid trust created a worktree")
			}
		})
	}
}
