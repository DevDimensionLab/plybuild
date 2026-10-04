package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
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
