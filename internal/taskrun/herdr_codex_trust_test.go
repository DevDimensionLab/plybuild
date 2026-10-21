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

	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
)

func TestWorkflowCodexTrustGrantPreview(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTrustTestHome(t, f)
	raw, _ := os.ReadFile(f.File)
	var r map[string]any
	if e := json.Unmarshal(raw, &r); e != nil {
		t.Fatal(e)
	}
	r["codex_project_trust"] = map[string]any{"mode": "process-local", "repository_root": filepath.Join(f.R.WorkspaceRoot, "target")}
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), r)
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatalf("explicit process-local trust grant rejected: %v", err)
	}
	if p.(WorkflowPreview).Confirmation == nil {
		t.Fatal("missing trust confirmation")
	}
	if workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("preview contacted provider")
	}
}

func workflowGrant(t *testing.T, f *workflowFixture) string {
	t.Helper()
	home := workflowTrustTestHome(t, *f)
	f.R.CodexProjectTrust = &CodexProjectTrust{"process-local", filepath.Join(f.R.WorkspaceRoot, "target")}
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), f.R)
	return home
}

func workflowTrustSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		s, e := os.Lstat(p)
		if e != nil {
			return e
		}
		v := s.Mode().String()
		if s.Mode().IsRegular() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			v += hash(b)
		}
		if s.Mode()&os.ModeSymlink != 0 {
			l, e := os.Readlink(p)
			if e != nil {
				return e
			}
			v += l
		}
		out[p] = v
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestWorkflowCodexTrustLegacyAndConflict(t *testing.T) {
	f := workflowTestFixture(t)
	before := workflowTrustSnapshot(t, f.R.WorkspaceRoot)
	v, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	p := v.(WorkflowPreview)
	if !reflect.DeepEqual(before, workflowTrustSnapshot(t, f.R.WorkspaceRoot)) {
		t.Fatal("legacy preview wrote files")
	}
	oldConfirmation := digest(struct {
		Request  WorkflowRequest
		Observed Observed
		Paths    WorkflowPaths
		Effects  []string
	}{f.R, p.Observed, p.Paths, p.Effects})
	if *p.Confirmation != oldConfirmation {
		t.Fatal("legacy confirmation changed")
	}
	raw, _ := os.ReadFile(f.File)
	r, e := ReadWorkflowRequest(f.File)
	if e != nil {
		t.Fatal(e)
	}
	canonical, _ := Canonical(r)
	if !bytes.Equal(raw, canonical) || bytes.Contains(canonical, []byte("codex_project_trust")) {
		t.Fatal("legacy canonical bytes changed")
	}
	args := workflowStartArgv(r, "/task", "pane")
	if strings.Contains(strings.Join(args, " "), "trust_level") {
		t.Fatal("legacy argv grants trust")
	}
	o, e := WorkflowStart(f.D, f.File, *p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	workflowGrant(t, &f)
	if workflowID(f.R) != o.RunID || digest(f.R) == o.RequestSHA256 {
		t.Fatal("trust must change bytes, not the request-key reservation")
	}
	if _, e = WorkflowStart(f.D, f.File, *p.Confirmation); e == nil {
		t.Fatal("changed grant reused old reservation")
	}
	for _, c := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatal("conflict replayed " + c)
		}
	}
}

func TestWorkflowCodexTrustActualAdapterAndPreview(t *testing.T) {
	f := workflowTestFixture(t)
	home := workflowGrant(t, &f)
	f.R.Runtime.ConfigProfile = ptr("chosen-profile")
	if e := os.WriteFile(filepath.Join(home, "chosen-profile.config.toml"), []byte("model = 'another-default'\n"), 0600); e != nil {
		t.Fatal(e)
	}
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), f.R)
	before := workflowTrustSnapshot(t, f.R.WorkspaceRoot)
	v, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	p := v.(WorkflowPreview)
	if !reflect.DeepEqual(before, workflowTrustSnapshot(t, f.R.WorkspaceRoot)) {
		t.Fatal("trust preview wrote files")
	}
	if workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("preview called Herdr")
	}
	if p.CodexTrust.RepositoryRoot != filepath.Join(f.R.WorkspaceRoot, "target") || p.CodexTrust.WorktreeRoot == p.CodexTrust.RepositoryRoot {
		t.Fatal("linked worktree root was guessed from cwd")
	}
	if !strings.Contains(strings.Join(p.Effects, " "), "process-local") || !strings.Contains(strings.Join(p.Effects, " "), "sandbox") {
		t.Fatal("trust effect was hidden")
	}
	o, e := WorkflowStart(f.D, f.File, *p.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(f.Calls)
	var got []string
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		var c struct{ Argv []string }
		if e = json.Unmarshal(line, &c); e != nil {
			t.Fatal(e)
		}
		if strings.Join(c.Argv[:2], " ") == "agent start" {
			got = c.Argv
		}
	}
	want := workflowStartArgv(f.R, p.CodexTrust.CWD, o.Transport.PaneID)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("actual start argv differs: %v", got)
	}
	for _, s := range []string{"--model synthetic-model", "--profile chosen-profile", "-a on-request", `approvals_reviewer="auto_review"`, `default_permissions="synthetic-policy"`} {
		if !strings.Contains(strings.Join(got, " "), s) {
			t.Fatalf("runtime selection changed: %s", s)
		}
	}
	for _, c := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatal("wrong effect count: " + c)
		}
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowCodexTrustRootAndTOMLEscaping(t *testing.T) {
	root, _ := filepath.EvalSymlinks(t.TempDir())
	repo := filepath.Join(root, `repo.with space "quote" back\slash`)
	if e := os.Mkdir(repo, 0700); e != nil {
		t.Fatal(e)
	}
	runGit(t, repo, "init", "-b", "main")
	if e := os.Mkdir(filepath.Join(repo, "sub"), 0700); e != nil {
		t.Fatal(e)
	}
	f, e := workflowTrustLayout(filepath.Join(repo, "sub"))
	if e != nil {
		t.Fatal(e)
	}
	if f.RepositoryRoot != repo || f.WorktreeRoot != repo {
		t.Fatal("ordinary repo root differs")
	}
	r := WorkflowRequest{RequestKey: "test", WorkspaceRoot: root, CodexProjectTrust: &CodexProjectTrust{"process-local", repo}}
	args := workflowStartArgv(r, f.CWD, "pane")
	var parsed map[string]any
	if e = toml.Unmarshal([]byte(args[len(args)-1]), &parsed); e != nil {
		t.Fatal(e)
	}
	projects := workflowTOMLTable(parsed["projects"])
	if len(projects) != 1 || workflowTOMLTable(projects[repo])["trust_level"] != "trusted" {
		t.Fatalf("path became another TOML key: %#v", parsed)
	}
}

func TestWorkflowCodexTrustInvalidShape(t *testing.T) {
	f := workflowTestFixture(t)
	raw, _ := os.ReadFile(f.File)
	var base map[string]any
	json.Unmarshal(raw, &base)
	for name, g := range map[string]any{"null": nil, "extra": map[string]any{"mode": "process-local", "repository_root": "/root", "extra": true}, "missing": map[string]any{"mode": "process-local"}, "wrong_mode": map[string]any{"mode": "persistent", "repository_root": "/root"}, "relative": map[string]any{"mode": "process-local", "repository_root": "relative"}} {
		t.Run(name, func(t *testing.T) {
			base["codex_project_trust"] = g
			writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), base)
			if _, e := WorkflowPreviewStart(f.D, f.File); e == nil {
				t.Fatal("invalid grant accepted")
			}
		})
	}
	if workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("invalid request had effects")
	}
}

func TestWorkflowCodexTrustRefusesBeforeReservation(t *testing.T) {
	for _, name := range []string{"untrusted", "wrong_root", "cwd_layer", "main_layer", "parent_layer", "skills_layer", "symlink_layer", "config_symlink", "policy_drift", "binary_drift", "root_drift", "common_dir_drift", "custom_search", "unknown_provider"} {
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			home := workflowGrant(t, &f)
			v, e := WorkflowPreviewStart(f.D, f.File)
			if e != nil {
				t.Fatal(e)
			}
			p := v.(WorkflowPreview)
			cwd := p.CodexTrust.CWD
			repo := p.CodexTrust.RepositoryRoot
			switch name {
			case "untrusted":
				file := filepath.Join(home, "config.toml")
				b, _ := os.ReadFile(file)
				b = append(b, []byte("\n[projects."+tomlQuote(repo)+"]\ntrust_level='untrusted'\n")...)
				if e = os.WriteFile(file, b, 0600); e != nil {
					t.Fatal(e)
				}
			case "wrong_root":
				f.R.CodexProjectTrust.RepositoryRoot = cwd
				writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), f.R)
			case "cwd_layer", "main_layer", "parent_layer", "skills_layer":
				path := filepath.Join(cwd, ".codex")
				if name == "main_layer" {
					path = filepath.Join(repo, ".codex")
				}
				if name == "parent_layer" {
					path = filepath.Join(f.R.WorkspaceRoot, ".codex")
				}
				if name == "skills_layer" {
					path = filepath.Join(cwd, ".agents")
				}
				if e = os.Mkdir(path, 0700); e != nil {
					t.Fatal(e)
				}
			case "symlink_layer":
				e = os.Symlink(filepath.Join(f.R.WorkspaceRoot, "missing"), filepath.Join(cwd, ".codex"))
				if e != nil {
					t.Fatal(e)
				}
			case "config_symlink":
				file := filepath.Join(home, "config.toml")
				if e = os.Rename(file, file+".saved"); e != nil {
					t.Fatal(e)
				}
				if e = os.Symlink(file+".saved", file); e != nil {
					t.Fatal(e)
				}
			case "policy_drift":
				file := filepath.Join(home, "config.toml")
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append(b, []byte("\nmodel='changed'\n")...), 0600)
			case "binary_drift":
				os.WriteFile(f.R.Herdr.Executable.Path, []byte("#!/bin/sh\nexit 99\n"), 0700)
			case "root_drift":
				if e = os.Rename(repo, repo+"-old"); e != nil {
					t.Fatal(e)
				}
				if e = os.Mkdir(repo, 0700); e != nil {
					t.Fatal(e)
				}
			case "common_dir_drift":
				os.WriteFile(filepath.Join(p.CodexTrust.GitDir, "commondir"), []byte("../wrong\n"), 0600)
			case "custom_search":
				file := filepath.Join(home, "config.toml")
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append([]byte("project_root_markers=['.other']\n"), b...), 0600)
			case "unknown_provider":
				f.R.Runtime.Provider = "unknown"
				writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.File), f.R)
			}
			before := workflowTrustSnapshot(t, f.R.WorkspaceRoot)
			if _, e = WorkflowStart(f.D, f.File, *p.Confirmation); e == nil {
				t.Fatal("drift accepted")
			} else if strings.Contains(strings.ToLower(e.Error()), "human") {
				t.Fatalf("drift routed an agent-observable fact to a human: %v", e)
			}
			if _, e = os.Stat(workflowRoot(f.R.WorkspaceRoot)); !os.IsNotExist(e) {
				t.Fatal("rejection reserved a run")
			}
			if workflowTestCalls(t, f, "tab create") != 0 {
				t.Fatal("rejection contacted provider")
			}
			if !reflect.DeepEqual(before, workflowTrustSnapshot(t, f.R.WorkspaceRoot)) {
				t.Fatal("rejection wrote files")
			}
		})
	}
}

func workflowTrustTestHome(t *testing.T, f workflowFixture) string {
	t.Helper()
	home := filepath.Join(f.R.WorkspaceRoot, "user-home")
	config := filepath.Join(home, ".codex")
	if e := os.MkdirAll(config, 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(config, "config.toml"), []byte("[permissions.synthetic-policy]\nfilesystem={\":root\"=\"read\"}\n"), 0600); e != nil {
		t.Fatal(e)
	}
	t.Setenv("HOME", home)
	// This changes only the isolated stand-in process's config lookup.
	t.Setenv("CODEX_HOME", config)
	return config
}

func TestWorkflowReadinessIdleWithoutPendingMayBecomeReady(t *testing.T) {
	f := workflowTestFixture(t)
	f.D.HerdrTimeout = 2 * time.Second
	workflowTestModel(t, f, map[string]any{"lost_command": "agent start", "observations": []any{
		map[string]any{"agent_status": "idle", "agent_session": nil, "drop_fields": []string{"launch_pending"}},
		map[string]any{"agent_status": "idle", "agent_session": nil, "launch_pending": false}, map[string]any{},
	}})
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatalf("early nonzero/missing pending stopped a live startup: %v", err)
	}
	if o.Transport.AgentSessionID != "fixture-session" {
		t.Fatal("real session not bound")
	}
	for _, c := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatalf("%s must occur once", c)
		}
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowCodexTrustManagedPermissionTable(t *testing.T) {
	r := WorkflowRequest{}
	r.Runtime.PermissionBinding.ProfileID = "synthetic-policy"
	for _, tc := range []struct {
		name, config string
		denied       bool
	}{
		{"allowed", "allowed_approval_policies=['on-request']\nallowed_approvals_reviewers=['auto_review']\n[allowed_permission_profiles]\nsynthetic-policy=true\n", false},
		{"denied", "[allowed_permission_profiles]\nsynthetic-policy=false\n", true},
		{"omitted", "[allowed_permission_profiles]\nother=true\n", true},
		{"wrong_type", "allowed_permission_profiles=['synthetic-policy']\n", true},
		{"approval_denied", "allowed_approval_policies=['never']\n", true},
		{"reviewer_denied", "allowed_approvals_reviewers=['user']\n", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var cfg map[string]any
			if e := toml.Unmarshal([]byte(tc.config), &cfg); e != nil {
				t.Fatal(e)
			}
			e := workflowTrustRequirements(r, "/synthetic/requirements.toml", cfg)
			if (e != nil) != tc.denied {
				t.Fatalf("managed policy denied=%v, want %v: %v", e != nil, tc.denied, e)
			}
		})
	}
}

// workflowTrustJourney prepares the grant, lets setup change native facts
// before the read-only preview, and captures the warn log of both calls.
func workflowTrustJourney(t *testing.T, setup func(f workflowFixture, home, cwd, repo string)) (workflowFixture, WorkflowPreview, *bytes.Buffer) {
	t.Helper()
	f := workflowTestFixture(t)
	home := workflowGrant(t, &f)
	v, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	setup(f, home, v.(WorkflowPreview).Observed.Target.WorktreeLocator, f.R.CodexProjectTrust.RepositoryRoot)
	var log bytes.Buffer
	previous := logrus.StandardLogger().Out
	logrus.SetOutput(&log)
	t.Cleanup(func() { logrus.SetOutput(previous) })
	v, e = WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatalf("optional trust grant locked the ordinary start: %v", e)
	}
	return f, v.(WorkflowPreview), &log
}

func workflowTrustAppendConfig(t *testing.T, home, text string) {
	t.Helper()
	file := filepath.Join(home, "config.toml")
	b, e := os.ReadFile(file)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(file, append(b, []byte("\n"+text+"\n")...), 0600); e != nil {
		t.Fatal(e)
	}
}

// Project instructions in the Task worktree and its main repository. The shared
// exclude keeps the prepared Task clean; Ply's search sees present files the
// same way whether Git tracks them or not.
func workflowTrustAgentsFile(t *testing.T, cwd, repo string) {
	t.Helper()
	for _, file := range []string{filepath.Join(cwd, "AGENTS.md"), filepath.Join(repo, "AGENTS.md")} {
		if e := os.WriteFile(file, []byte("# Instructions\n"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	exclude := filepath.Join(repo, ".git", "info", "exclude")
	if e := os.MkdirAll(filepath.Dir(exclude), 0700); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(exclude)
	if e := os.WriteFile(exclude, append(b, []byte("\n/AGENTS.md\n")...), 0600); e != nil {
		t.Fatal(e)
	}
}

func workflowTrustAgentStart(t *testing.T, f workflowFixture) string {
	t.Helper()
	b, _ := os.ReadFile(f.Calls)
	got := ""
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		var c struct{ Argv []string }
		if e := json.Unmarshal(line, &c); e != nil {
			t.Fatal(e)
		}
		if len(c.Argv) > 1 && strings.Join(c.Argv[:2], " ") == "agent start" {
			got = strings.Join(c.Argv, " ")
		}
	}
	return got
}

func workflowTrustAssertOrdinaryStart(t *testing.T, f workflowFixture, p WorkflowPreview, state string) WorkflowRun {
	t.Helper()
	text := strings.ToLower(strings.Join(p.Effects, "\n"))
	if p.Confirmation == nil || strings.Contains(text, "human") || strings.Contains(text, "inspect") {
		t.Fatalf("preview routed the agent-observable trust fact to a human: %v", p.Effects)
	}
	if p.CodexTrust == nil || p.CodexTrust.Effect != state || p.CodexTrust.Reason == "" || p.CodexTrust.RepositoryRoot != f.R.CodexProjectTrust.RepositoryRoot {
		t.Fatalf("requested and effective trust are not distinguishable: %+v", p.CodexTrust)
	}
	o, e := WorkflowStart(f.D, f.File, *p.Confirmation)
	if e != nil {
		t.Fatalf("ordinary authorized start stopped: %v", e)
	}
	if argv := workflowTrustAgentStart(t, f); argv == "" || strings.Contains(argv, "trust_level") {
		t.Fatalf("unsupported trust effect reached the provider: %q", argv)
	}
	for _, s := range []string{"-a on-request", `approvals_reviewer="auto_review"`, `default_permissions="synthetic-policy"`} {
		if !strings.Contains(workflowTrustAgentStart(t, f), s) {
			t.Fatalf("native runtime control changed: %s", s)
		}
	}
	if o.CodexTrust == nil || o.CodexTrust.State != state || o.CodexTrust.RepositoryRoot != f.R.CodexProjectTrust.RepositoryRoot {
		t.Fatalf("run does not show the effective trust outcome: %+v", o.CodexTrust)
	}
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil || !equal(s.Request, f.R) || s.Result.RequestSHA256 != digest(f.R) || s.Request.CodexProjectTrust == nil {
		t.Fatal("exact requested grant was not preserved")
	}
	for _, c := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatal("wrong effect count: " + c)
		}
	}
	workflowAssertNoReplay(t, f, o)
	return o
}

func TestWorkflowCodexTrustRedundantGrantWithAgentsInstructions(t *testing.T) {
	f, p, _ := workflowTrustJourney(t, func(f workflowFixture, home, cwd, repo string) {
		workflowTrustAgentsFile(t, cwd, repo)
		workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(repo)+"]\ntrust_level='trusted'")
	})
	workflowTrustAssertOrdinaryStart(t, f, p, "redundant")
	if !strings.Contains(strings.Join(p.Effects, "\n"), "already trusts") {
		t.Fatalf("redundant grant is not shown as a no-op: %v", p.Effects)
	}
}

func TestWorkflowCodexTrustUnconfirmableGrantIsWithheld(t *testing.T) {
	for name, setup := range map[string]func(f workflowFixture, home, cwd, repo string){
		"instructions_untrusted_root": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAgentsFile(t, cwd, repo)
		},
		"explicit_untrusted": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(repo)+"]\ntrust_level='untrusted'")
		},
		"unknown_policy": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(repo)+"]\ntrust_level='future-level'")
		},
		"trusted_parent_explicit_untrusted_root": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(filepath.Dir(repo))+"]\ntrust_level='trusted'\n[projects."+tomlQuote(repo)+"]\ntrust_level='untrusted'")
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, p, log := workflowTrustJourney(t, setup)
			workflowTrustAssertOrdinaryStart(t, f, p, "withheld")
			for _, s := range []string{"level=warning", "codex_project_trust", "withheld", "run_id=" + workflowID(f.R)} {
				if !strings.Contains(log.String(), s) {
					t.Fatalf("warn evidence lacks %q: %s", s, log.String())
				}
			}
		})
	}
}

func TestWorkflowCodexTrustDriftAfterReservationWithholdsGrant(t *testing.T) {
	f := workflowTestFixture(t)
	home := workflowGrant(t, &f)
	v, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	p := v.(WorkflowPreview)
	if p.CodexTrust.Effect != "" {
		t.Fatalf("a confirmable grant changed its historical facts: %+v", p.CodexTrust)
	}
	f.D.Fault = func(point string) error {
		if point == "workflow_before_tab_send" {
			workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(p.CodexTrust.RepositoryRoot)+"]\ntrust_level='untrusted'")
		}
		return nil
	}
	o, e := WorkflowStart(f.D, f.File, *p.Confirmation)
	if e != nil {
		t.Fatalf("pre-launch drift locked the reserved start: %v", e)
	}
	f.D.Fault = nil
	if argv := workflowTrustAgentStart(t, f); strings.Contains(argv, "trust_level") {
		t.Fatalf("drifted grant was applied: %s", argv)
	}
	if o.CodexTrust == nil || o.CodexTrust.State != "withheld" || !strings.Contains(o.CodexTrust.Reason, "changed") {
		t.Fatalf("drift outcome is not visible: %+v", o.CodexTrust)
	}
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = workflowTrustFresh(s); e != nil {
		t.Fatalf("a withheld grant still locks later callbacks: %v", e)
	}
	for _, c := range []string{"tab create", "agent start", "agent prompt"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatal("wrong effect count: " + c)
		}
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowCodexTrustLostStartReplyKeepsDecision(t *testing.T) {
	f, p, _ := workflowTrustJourney(t, func(f workflowFixture, home, cwd, repo string) {
		workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(repo)+"]\ntrust_level='untrusted'")
	})
	f.D.HerdrTimeout = 2 * time.Second
	workflowTestModel(t, f, map[string]any{"lost_command": "agent start", "observations": []any{workflowPendingObservation(), map[string]any{}}})
	o := workflowTrustAssertOrdinaryStart(t, f, p, "withheld")
	if o.Transport.AgentSessionID != "fixture-session" {
		t.Fatal("lost reply did not continue the same session")
	}
}

// The candidate CLI keeps machine-readable stdout while the warning goes to
// stderr, and the same exact request starts once through preview and apply.
func TestWorkflowCodexTrustCandidateCLIJourney(t *testing.T) {
	binary := deliveryCLIBinary(t)
	for name, setup := range map[string]func(f workflowFixture, home, cwd, repo string){
		"redundant": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAgentsFile(t, cwd, repo)
			workflowTrustAppendConfig(t, home, "[projects."+tomlQuote(repo)+"]\ntrust_level='trusted'")
		},
		"withheld": func(f workflowFixture, home, cwd, repo string) {
			workflowTrustAgentsFile(t, cwd, repo)
		},
	} {
		t.Run(name, func(t *testing.T) {
			f, _, _ := workflowTrustJourney(t, setup)
			run := func(env []string, args ...string) (map[string]json.RawMessage, string) {
				t.Helper()
				c := exec.Command(binary, append([]string{"workflow", "run", "start", "--file", f.File, "--format", "json"}, args...)...)
				c.Dir = f.R.WorkspaceRoot
				c.Env = append(os.Environ(), env...)
				var stdout, stderr bytes.Buffer
				c.Stdout, c.Stderr = &stdout, &stderr
				err := c.Run()
				var out map[string]json.RawMessage
				if e := json.Unmarshal(stdout.Bytes(), &out); e != nil {
					t.Fatalf("stdout is not one JSON result: %v (exit %v)\n%s\n%s", e, err, stdout.String(), stderr.String())
				}
				if err != nil {
					t.Fatalf("candidate start failed: %v\n%s\n%s", err, stdout.String(), stderr.String())
				}
				return out, stderr.String()
			}
			preview, log := run(nil)
			var trust workflowTrustFacts
			var confirmation string
			if json.Unmarshal(preview["codex_project_trust"], &trust) != nil || json.Unmarshal(preview["confirmation"], &confirmation) != nil || trust.Effect != name || confirmation == "" {
				t.Fatalf("preview lacks the effective trust outcome: %s", preview["codex_project_trust"])
			}
			if strings.Contains(strings.ToLower(string(preview["effects"])), "human") {
				t.Fatalf("preview addresses a human: %s", preview["effects"])
			}
			if name == "withheld" && (!strings.Contains(log, "level=warning") || !strings.Contains(log, "codex_project_trust=withheld")) {
				t.Fatalf("missing stderr warn evidence: %s", log)
			}
			started, _ := run([]string{"HERDR_ENV=1"}, "--apply", "--confirm", confirmation)
			var effect CodexTrustEffect
			if json.Unmarshal(started["codex_project_trust"], &effect) != nil || effect.State != name {
				t.Fatalf("run lacks the effective trust state: %s", started["codex_project_trust"])
			}
			if argv := workflowTrustAgentStart(t, f); argv == "" || strings.Contains(argv, "trust_level") {
				t.Fatalf("unsupported trust effect reached the provider: %q", argv)
			}
			again, _ := run([]string{"HERDR_ENV=1"}, "--apply", "--confirm", confirmation)
			if string(again["run_id"]) != string(started["run_id"]) {
				t.Fatal("repeated apply changed the reserved run")
			}
			for _, c := range []string{"tab create", "agent start", "agent prompt"} {
				if workflowTestCalls(t, f, c) != 1 {
					t.Fatal("wrong effect count: " + c)
				}
			}
		})
	}
}
