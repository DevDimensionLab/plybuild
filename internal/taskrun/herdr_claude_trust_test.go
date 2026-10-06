package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/unicode/norm"
)

func TestWorkflowClaudeTrustExplicitGrantPreview(t *testing.T) {
	f := workflowProviderFixture(t, "claude")
	v, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	target := v.(WorkflowPreview).Observed.Target.WorktreeLocator
	config := filepath.Join(f.R.WorkspaceRoot, "claude-config")
	if err := os.Mkdir(config, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(config, ".claude.json")
	initial := []byte(`{"projects":{},"unrelated":{"ratio":0.25,"large":9007199254740993}}`)
	if err := os.WriteFile(path, initial, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", config)
	doc := workflowProviderDocument(t, f.File)
	doc["claude_project_trust"] = map[string]any{"mode": "configuration", "worktree_root": target}
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), doc)
	before := workflowTrustSnapshot(t, config)
	v, err = WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatalf("explicit Claude folder trust grant rejected: %v", err)
	}
	raw, err := Canonical(v)
	if err != nil {
		t.Fatal(err)
	}
	var preview map[string]json.RawMessage
	if err = json.Unmarshal(raw, &preview); err != nil {
		t.Fatal(err)
	}
	if len(preview["claude_project_trust"]) == 0 || !bytes.Contains(preview["claude_project_trust"], []byte(path)) {
		t.Fatal("preview must expose the exact Claude configuration effect")
	}
	if !reflect.DeepEqual(before, workflowTrustSnapshot(t, config)) || workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("trust preview wrote configuration or opened a tab")
	}
}

func workflowClaudeTrustFixture(t *testing.T) (workflowFixture, WorkflowPreview, string) {
	t.Helper()
	f := workflowProviderFixture(t, "claude")
	v, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(f.R.WorkspaceRoot, "claude-config")
	if err = os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	t.Setenv("CLAUDE_CODE_CUSTOM_OAUTH_URL", "")
	path := filepath.Join(dir, ".claude.json")
	if err = os.WriteFile(path, []byte(`{"projects":{},"opaque":12345678901234567890,"ratio":0.25}`), 0600); err != nil {
		t.Fatal(err)
	}
	f.R.ClaudeProjectTrust = &ClaudeProjectTrust{"configuration", v.(WorkflowPreview).Observed.Target.WorktreeLocator}
	workflowProviderRebind(t, &f)
	v, err = WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	return f, v.(WorkflowPreview), path
}

func TestWorkflowClaudeTrustStartAndReturn(t *testing.T) {
	f, p, path := workflowClaudeTrustFixture(t)
	workflowTestModel(t, f, map[string]any{"claude_trust_config": path, "claude_trust_key": p.ClaudeTrust.ProjectKey})
	o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
	if err != nil || o.ClaudeTrust == nil || o.ClaudeTrust.State != "written" {
		t.Fatalf("trust not written at start: %+v %v", o.ClaudeTrust, err)
	}
	if o.ClaudeTrust.ProjectKey != p.Observed.Target.WorktreeLocator || o.ClaudeTrust.ConfigPath != path || o.Transport.AgentSessionID != "fixture-session" {
		t.Fatalf("wrong effect/session binding: %+v", o)
	}
	calls, err := os.ReadFile(f.Calls)
	if err != nil || !bytes.Contains(calls, []byte(`"claude_trusted_before_start": true`)) {
		t.Fatalf("provider started before trust was published: %s %v", calls, err)
	}
	tabs := workflowProviderCalls(t, f, "tab create")
	if len(tabs) != 1 || !strings.Contains(strings.Join(tabs[0], "\n"), "CLAUDE_CONFIG_DIR="+filepath.Dir(path)) {
		t.Fatal("tab does not bind the configuration directory")
	}
	// Mutable provider configuration is not callback authority. Losing the trust
	// entry after launch must neither replay its write nor invalidate a return.
	replacement := []byte(`{"projects":{},"native_change":true}`)
	if err = os.WriteFile(path, replacement, 0600); err != nil {
		t.Fatal(err)
	}
	o = workflowTestAccept(t, f, o)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || o.FinalReturn.State != "accepted" {
		t.Fatalf("config mutation broke return: %+v %v", o, err)
	}
	workflowAssertNoReplay(t, f, o)
	after, _ := os.ReadFile(path)
	if !bytes.Equal(after, replacement) || workflowTestCalls(t, f, "agent start") != 1 {
		t.Fatal("readback/replay rewrote configuration or relaunched")
	}
}

func TestWorkflowClaudeTrustFallbackStartsOnce(t *testing.T) {
	for _, tc := range []struct{ name, state string }{{"missing", "skipped"}, {"invalid", "failed"}, {"write_failure", "failed"}, {"postwrite_unknown", "unknown"}} {
		t.Run(tc.name, func(t *testing.T) {
			f, p, path := workflowClaudeTrustFixture(t)
			switch tc.name {
			case "missing":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "invalid":
				if err := os.WriteFile(path, []byte(`{"broken"`), 0600); err != nil {
					t.Fatal(err)
				}
			default:
				fault := "claude_trust_before_publish"
				if tc.name == "postwrite_unknown" {
					fault = "claude_trust_after_publish"
				}
				f.D.Fault = func(at string) error {
					if at == fault {
						return workflowError(5, "fixture write failure")
					}
					return nil
				}
			}
			o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
			if err != nil || o.ClaudeTrust == nil || o.ClaudeTrust.State != tc.state {
				t.Fatalf("fallback = %+v %v, want %s", o.ClaudeTrust, err, tc.state)
			}
			found := false
			for _, reason := range o.Reasons {
				found = found || reason.Code == "claude_trust_"+tc.state
			}
			if !found || workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
				t.Fatal("fallback warning missing or startup was repeated/stopped")
			}
			before := workflowTrustSnapshot(t, filepath.Dir(path))
			if _, err = WorkflowStart(f.D, f.File, *p.Confirmation); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) || workflowTestCalls(t, f, "agent start") != 1 {
				t.Fatal("repeat start retried trust or provider startup")
			}
		})
	}
}

func TestWorkflowClaudeTrustRejectsWrongGrant(t *testing.T) {
	for _, tc := range []string{"provider", "mode", "target", "null"} {
		t.Run(tc, func(t *testing.T) {
			f, _, path := workflowClaudeTrustFixture(t)
			doc := workflowProviderDocument(t, f.File)
			g := doc["claude_project_trust"].(map[string]any)
			switch tc {
			case "provider":
				doc["runtime"].(map[string]any)["provider"] = "codex"
			case "mode":
				g["mode"] = "bypass"
			case "target":
				g["worktree_root"] = f.R.WorkspaceRoot
			case "null":
				doc["claude_project_trust"] = nil
			}
			writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), doc)
			before := workflowTrustSnapshot(t, filepath.Dir(path))
			if _, err := WorkflowPreviewStart(f.D, f.File); err == nil {
				t.Fatal("invalid grant accepted")
			}
			workflowProviderNoReservation(t, f)
			if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) {
				t.Fatal("invalid grant wrote config")
			}
		})
	}
}

func TestWorkflowClaudeTrustLegacyHasNoEffect(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			f, _, path := workflowClaudeTrustFixture(t)
			f.R.ClaudeProjectTrust = nil
			f.R.Runtime.Provider = provider
			if provider == "codex" {
				f.R.Runtime.PermissionBinding.ProfileID = "local-development"
			}
			workflowProviderRebind(t, &f)
			workflowTestModel(t, f, map[string]any{"provider": provider})
			before := workflowTrustSnapshot(t, filepath.Dir(path))
			v, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			o, err := WorkflowStart(f.D, f.File, *v.(WorkflowPreview).Confirmation)
			if err != nil {
				t.Fatal(err)
			}
			if o.ClaudeTrust != nil || !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) {
				t.Fatal("old request acquired persistent trust")
			}
		})
	}
}

func TestClaudeTrustPreviewLocationsAndUnicode(t *testing.T) {
	dir := t.TempDir()
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	t.Setenv("CLAUDE_CODE_CUSTOM_OAUTH_URL", "")
	g := &ClaudeProjectTrust{"configuration", filepath.Join(dir, "cafe\u0301")}
	preview, err := PreviewClaudeProjectTrust(g)
	if err != nil || preview.ConfigPath != filepath.Join(dir, ".claude.json") || preview.ProjectKey != norm.NFC.String(g.WorktreeRoot) {
		t.Fatalf("preview %+v %v", preview, err)
	}
	legacy := filepath.Join(dir, ".config.json")
	if err = os.WriteFile(legacy, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	preview, err = PreviewClaudeProjectTrust(g)
	if err != nil || preview.ConfigPath != legacy {
		t.Fatalf("legacy precedence %+v %v", preview, err)
	}
	for _, configDir := range []string{"", "relative/path"} {
		t.Setenv("CLAUDE_CONFIG_DIR", configDir)
		preview, err = PreviewClaudeProjectTrust(g)
		if err != nil || preview.Reason == "" || preview.ConfigPath != "" {
			t.Fatalf("unsafe resolution %+v %v", preview, err)
		}
	}
}

func TestWorkflowClaudeTrustStartupStopsDoNotWrite(t *testing.T) {
	for _, at := range []string{"workflow_before_tab_send", "workflow_before_agent_send"} {
		t.Run(at, func(t *testing.T) {
			f, p, path := workflowClaudeTrustFixture(t)
			before := workflowTrustSnapshot(t, filepath.Dir(path))
			f.D.Fault = func(point string) error {
				if point == at {
					return workflowError(5, "bounded startup stop")
				}
				return nil
			}
			o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
			if err == nil || o.ClaudeTrust != nil {
				t.Fatalf("unexpected start: %+v %v", o, err)
			}
			if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) || workflowTestCalls(t, f, "agent start") != 0 {
				t.Fatal("trust applied before start boundary")
			}
			f.D.Fault = nil
			if _, err = WorkflowStart(f.D, f.File, *p.Confirmation); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) || workflowTestCalls(t, f, "agent start") != 0 {
				t.Fatal("stopped reservation retried trust or start")
			}
		})
	}
}

func TestWorkflowClaudeTrustChangedEnvironmentRejectsConfirmation(t *testing.T) {
	f, p, path := workflowClaudeTrustFixture(t)
	dir := filepath.Join(filepath.Dir(path), "other")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	before := workflowTrustSnapshot(t, filepath.Dir(path))
	if _, err := WorkflowStart(f.D, f.File, *p.Confirmation); err == nil {
		t.Fatal("config drift accepted")
	}
	workflowProviderNoReservation(t, f)
	if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) {
		t.Fatal("rejected confirmation wrote trust")
	}
}

func TestWorkflowClaudeTrustRequiresFreshTarget(t *testing.T) {
	f, p, path := workflowClaudeTrustFixture(t)
	before := workflowTrustSnapshot(t, filepath.Dir(path))
	f.D.Fault = func(point string) error {
		if point == "workflow_before_agent_send" {
			cwd := p.Observed.Target.WorktreeLocator
			if err := os.Rename(cwd, cwd+"-preserved"); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(cwd, 0700); err != nil {
				t.Fatal(err)
			}
		}
		return nil
	}
	_, _ = WorkflowStart(f.D, f.File, *p.Confirmation)
	if !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) || workflowTestCalls(t, f, "agent start") != 0 {
		t.Fatal("changed physical Task directory received trust or a provider start")
	}
}

func TestWorkflowClaudeTrustBindsActualShellEnvironment(t *testing.T) {
	f, p, _ := workflowClaudeTrustFixture(t)
	if _, err := WorkflowStart(f.D, f.File, *p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if workflowTestCalls(t, f, "pane run") != 1 {
		t.Fatal("provider config environment was not bound in the actual shell")
	}
}

func TestClaudeTrustShellEnvironmentClearsInheritedOverrides(t *testing.T) {
	for _, shell := range []string{"/bin/zsh", "/bin/bash"} {
		t.Run(filepath.Base(shell), func(t *testing.T) {
			if _, err := os.Stat(shell); os.IsNotExist(err) {
				t.Skip("shell is not installed")
			}
			for _, dir := range []string{"", filepath.Join(t.TempDir(), "config's spaces $value")} {
				t.Run(dir, func(t *testing.T) {
					root, err := filepath.EvalSymlinks(t.TempDir())
					if err != nil {
						t.Fatal(err)
					}
					p := &ClaudeTrustPreview{UserHome: os.Getenv("HOME"), ConfigDir: dir}
					ack := filepath.Join(root, "ack")
					command := workflowClaudeEnvironmentCommand(p, "proof-token", ack)
					command = "unset() { return 0; }; export() { return 0; }; " + command
					command += ` && /bin/sh -c 'printf "%s\n%s\n%s" "${CLAUDE_CONFIG_DIR-unset}" "${CLAUDE_CODE_CUSTOM_OAUTH_URL-unset}" "$HOME"'`
					flags := []string{"-f", "-c"}
					if shell == "/bin/bash" {
						flags = []string{"--noprofile", "--norc", "-c"}
					}
					run := func() ([]byte, error) {
						c := exec.Command(shell, append(flags, command)...)
						c.Env = append(os.Environ(), "CLAUDE_CONFIG_DIR=/stale/server/config", "CLAUDE_CODE_CUSTOM_OAUTH_URL=stale-server-value")
						return c.CombinedOutput()
					}
					b, err := run()
					wantDir := dir
					if dir == "" {
						wantDir = "unset"
					}
					if err != nil || string(b) != wantDir+"\nunset\n"+p.UserHome {
						t.Fatalf("wrong actual child environment %q %v", b, err)
					}
					b, err = readFile(ack, 128, true)
					if err != nil || string(b) != "proof-token" {
						t.Fatalf("missing private ack %q %v", b, err)
					}
					st, err := os.Stat(ack)
					if err != nil || st.Mode().Perm() != 0600 {
						t.Fatalf("nonprivate ack %v %v", st, err)
					}
					if b, err = run(); err == nil {
						t.Fatalf("same command overwrote its ack: %s", b)
					}
				})
			}
		})
	}
}

func TestWorkflowClaudeTrustUnknownEnvironmentDoesNotStartOrReplay(t *testing.T) {
	f, p, path := workflowClaudeTrustFixture(t)
	f.D.HerdrTimeout = 400 * time.Millisecond
	workflowTestModel(t, f, map[string]any{"claude_environment_no_ack": true})
	before := workflowTrustSnapshot(t, filepath.Dir(path))
	o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
	if err == nil || o.ClaudeTrust != nil {
		t.Fatalf("unconfirmed shell setup accepted: %+v %v", o, err)
	}
	if workflowTestCalls(t, f, "pane run") != 1 || workflowTestCalls(t, f, "agent start") != 0 || !reflect.DeepEqual(before, workflowTrustSnapshot(t, filepath.Dir(path))) {
		t.Fatal("unknown setup caused dependent effects")
	}
	if _, err = WorkflowStart(f.D, f.File, *p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if workflowTestCalls(t, f, "pane run") != 1 || workflowTestCalls(t, f, "agent start") != 0 {
		t.Fatal("unknown setup was replayed")
	}
}

func TestWorkflowClaudeTrustLostEnvironmentReplyUsesAck(t *testing.T) {
	f, p, _ := workflowClaudeTrustFixture(t)
	workflowTestModel(t, f, map[string]any{"failure": map[string]any{"command": "pane run", "exit": 7}})
	o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
	if err != nil || o.ClaudeTrust == nil || o.ClaudeTrust.State != "written" || workflowTestCalls(t, f, "pane run") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
		t.Fatalf("lost response did not use exact ack: %+v %v", o, err)
	}
}

func TestWorkflowClaudeTrustPreviewDescribesSkip(t *testing.T) {
	f, _, _ := workflowClaudeTrustFixture(t)
	t.Setenv("CLAUDE_CONFIG_DIR", "relative")
	v, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	p := v.(WorkflowPreview)
	if !strings.Contains(strings.Join(p.Effects, "\n"), "Skip automatic Claude folder trust: "+p.ClaudeTrust.Reason) {
		t.Fatal("preview promised a config write despite skip")
	}
}

func TestClaudeTrustEnvironmentFixtureExport(t *testing.T) {
	root := os.Getenv("PLY_CLAUDE_ENV_FIXTURE")
	if root == "" {
		t.Skip("native shell fixture export not requested")
	}
	if err := physical(root, false); err != nil {
		t.Fatal(err)
	}
	ack := filepath.Join(root, "native-environment.ack")
	p := &ClaudeTrustPreview{UserHome: os.Getenv("HOME")}
	const token = "ply-claude-environment-212-native-v1"
	value := struct{ Command, AckPath, Token string }{workflowClaudeEnvironmentCommand(p, token, ack), ack, token}
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(filepath.Join(root, "command.json"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err = f.Write(b); err != nil {
		t.Fatal(err)
	}
}
