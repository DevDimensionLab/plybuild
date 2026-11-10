package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/planningrepo"
)

// These journeys execute the shipped main package, not Cobra callbacks or a
// stand-in command tree. Git and workspace registration are real and isolated.
func TestWorkflowEpicPlanAcceptance(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "ply")
	build := exec.Command("go", "build", "-o", binary, "./ply")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build actual CLI: %v\n%s", err, out)
	}
	t.Run("EnglishExplicitAndPreview", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		repo := f.repo(t, "product main")
		before := planSnapshot(t, repo)
		args := []string{"workflow", "epic", "plan", "--repo", repo, "--root", f.root, "--format", "json"}
		preview := f.result(t, f.root, append(args, "--check")...)
		if preview.State != "preview" || preview.Git != nil || len(preview.Files) != 14 {
			t.Fatalf("preview: %+v", preview)
		}
		planAbsent(t, filepath.Join(f.root, "planning"))
		planAbsent(t, filepath.Join(f.root, ".ply"))
		if !reflect.DeepEqual(before, planSnapshot(t, repo)) {
			t.Fatal("preview changed source")
		}
		goal := "Literal $(touch NEVER) `touch NEVER` {{.Goal}}\nSecond line"
		result := f.result(t, f.root, append(args, "--goal", goal)...)
		f.verifyCreated(t, result, "en", []string{repo})
		b, _ := os.ReadFile(filepath.Join(result.PlanningPath, "goal.md"))
		if !strings.Contains(string(b), "# Goal") || !strings.Contains(string(b), "$(touch NEVER)") || !strings.Contains(string(b), "{{.Goal}}") {
			t.Fatalf("goal: %s", b)
		}
		planAbsent(t, filepath.Join(f.root, "NEVER"))
		planAbsent(t, filepath.Join(result.PlanningPath, "NEVER"))
		planAbsent(t, filepath.Join(f.root, ".ply"))
		if !reflect.DeepEqual(before, planSnapshot(t, repo)) {
			t.Fatal("create changed source")
		}
		planBefore := planSnapshot(t, result.PlanningPath)
		for _, check := range []bool{false, true} {
			repeat := append([]string(nil), args...)
			if check {
				repeat = append(repeat, "--check")
			}
			f.failure(t, f.root, nil, "existing planning repository", repeat...)
			if !reflect.DeepEqual(planBefore, planSnapshot(t, result.PlanningPath)) {
				t.Fatal("repeat changed plan")
			}
		}
	})
	t.Run("NorwegianMultiRepoPhysicalRootAndDeepCWD", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		a, b := f.repo(t, "æ service main"), f.repo(t, "api main")
		deep := filepath.Join(a, "some", "deep", "directory")
		mustPlanMkdir(t, deep)
		alias := filepath.Join(f.root, "root alias")
		if err := os.Symlink(f.root, alias); err != nil {
			t.Fatal(err)
		}
		aAlias := filepath.Join(f.root, "source alias")
		if err := os.Symlink(a, aAlias); err != nil {
			t.Fatal(err)
		}
		beforeA, beforeB := planSnapshot(t, a), planSnapshot(t, b)
		r := f.result(t, deep, "workflow", "epic", "plan", "--repo", aAlias, "--repo", b, "--root", alias, "--language", "nb", "--format", "json")
		f.verifyCreated(t, r, "nb", []string{a, b})
		if r.RootPath != f.root {
			t.Fatalf("root depends on cwd/alias: %s", r.RootPath)
		}
		goal, _ := os.ReadFile(filepath.Join(r.PlanningPath, "goal.md"))
		if !strings.Contains(string(goal), "Målet er ikke valgt ennå") {
			t.Fatalf("Norwegian unselected goal: %s", goal)
		}
		if !reflect.DeepEqual(beforeA, planSnapshot(t, a)) || !reflect.DeepEqual(beforeB, planSnapshot(t, b)) {
			t.Fatal("sources changed")
		}
		planAbsent(t, filepath.Join(deep, "planning"))
		planAbsent(t, filepath.Join(f.root, ".ply"))
	})
	t.Run("RegisteredProjectAllMembersAndRootOverride", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		a, b := f.repo(t, "a"), f.repo(t, "b")
		f.success(t, f.root, "workspace", "init")
		f.success(t, f.root, "workspace", "project", "add", "example", "--name", "Example", "--wrapper", f.root, "--repo", "b="+b, "--repo", "a="+a)
		deep := filepath.Join(a, "deep")
		mustPlanMkdir(t, deep)
		registry := planSnapshot(t, filepath.Join(f.root, ".ply"))
		beforeA, beforeB := planSnapshot(t, a), planSnapshot(t, b)
		r := f.result(t, deep, "workflow", "epic", "plan", "--project", "example", "--language", "nb", "--format", "json")
		f.verifyCreated(t, r, "nb", []string{a, b})
		alternate := filepath.Join(f.root, "alternate root")
		mustPlanMkdir(t, alternate)
		r = f.result(t, deep, "workflow", "epic", "plan", "--project", "example", "--root", alternate, "--format", "json")
		if r.RootPath != alternate || len(r.Repositories) != 2 {
			t.Fatalf("override: %+v", r)
		}
		if !reflect.DeepEqual(registry, planSnapshot(t, filepath.Join(f.root, ".ply"))) || !reflect.DeepEqual(beforeA, planSnapshot(t, a)) || !reflect.DeepEqual(beforeB, planSnapshot(t, b)) {
			t.Fatal("project/source state changed")
		}
	})
	t.Run("ExistingTargetsStopUnchanged", func(t *testing.T) {
		for _, kind := range []string{"legacy", "empty", "file", "symlink", "dangling", "unrelated-git", "legacy-symlink"} {
			t.Run(kind, func(t *testing.T) {
				f := newPlanCLI(t, binary)
				repo := f.repo(t, "source")
				target := filepath.Join(f.root, "planning")
				switch kind {
				case "file":
					mustPlanWrite(t, target, "reserved", 0o644)
				case "symlink":
					if err := os.Symlink(repo, target); err != nil {
						t.Fatal(err)
					}
				case "dangling":
					if err := os.Symlink(filepath.Join(f.root, "missing"), target); err != nil {
						t.Fatal(err)
					}
				default:
					mustPlanMkdir(t, target)
					if kind != "empty" {
						f.git(t, target, "init", "--initial-branch=main", "--template=")
						mustPlanWrite(t, filepath.Join(target, "README.md"), "other repo", 0o644)
					}
					if strings.HasPrefix(kind, "legacy") {
						for _, name := range []string{"AGENTS.md", "goal.md", "next-task.md", "status.md"} {
							mustPlanWrite(t, filepath.Join(target, name), "legacy instructions are not executed", 0o644)
						}
					}
					if kind == "legacy-symlink" {
						if err := os.Remove(filepath.Join(target, "AGENTS.md")); err != nil {
							t.Fatal(err)
						}
						if err := os.Symlink(filepath.Join(target, "goal.md"), filepath.Join(target, "AGENTS.md")); err != nil {
							t.Fatal(err)
						}
					}
				}
				before := planSnapshot(t, f.root)
				want := "path conflict"
				if kind == "legacy" {
					want = "existing planning repository"
				}
				for _, preview := range []bool{false, true} {
					args := []string{"workflow", "epic", "plan", "--repo", repo, "--root", f.root, "--format", "json"}
					if preview {
						args = append(args, "--check")
					}
					errText := f.failure(t, f.root, nil, want, args...)
					if !strings.Contains(errText, target) {
						t.Fatalf("missing absolute target: %s", errText)
					}
					if !reflect.DeepEqual(before, planSnapshot(t, f.root)) {
						t.Fatal("existing target or source changed")
					}
				}
			})
		}
	})
	t.Run("EarlyFailuresDoNotCreateTarget", func(t *testing.T) {
		for _, kind := range []string{"language", "no-selection", "no-root", "missing-root", "invalid-repo", "missing-repo", "duplicate-worktree", "duplicate-alias", "both-modes", "unknown-project", "unknown-flag", "positional", "format", "missing-git", "missing-identity", "render-size", "inside-source", "inside-common-dir", "inside-other-worktree"} {
			t.Run(kind, func(t *testing.T) {
				f := newPlanCLI(t, binary)
				repo := f.repo(t, "source")
				root := f.root
				args := []string{"workflow", "epic", "plan", "--repo", repo, "--root", root, "--format", "json"}
				env := map[string]string{}
				want := ""
				setArgs := func(extra ...string) { args = append([]string{"workflow", "epic", "plan"}, extra...) }
				switch kind {
				case "language":
					args = append(args, "--language", "de")
					want = "unsupported language"
				case "no-selection":
					setArgs("--root", root)
					want = "select exactly one"
				case "no-root":
					setArgs("--repo", repo)
					want = "--root"
				case "missing-root":
					root = filepath.Join(root, "absent")
					setArgs("--repo", repo, "--root", root)
					want = "root"
				case "invalid-repo":
					setArgs("--repo", root, "--root", root)
					want = "Git worktree"
				case "missing-repo":
					setArgs("--repo", filepath.Join(root, "absent"), "--root", root)
					want = "repository"
				case "duplicate-worktree":
					linked := filepath.Join(root, "linked")
					f.git(t, repo, "worktree", "add", "-b", "linked", linked)
					args = append(args, "--repo", linked)
					want = "duplicate repository identity"
				case "duplicate-alias":
					alias := filepath.Join(root, "alias")
					if err := os.Symlink(repo, alias); err != nil {
						t.Fatal(err)
					}
					args = append(args, "--repo", alias)
					want = "duplicate repository identity"
				case "both-modes":
					args = append(args, "--project", "example")
					want = "select exactly one"
				case "unknown-project":
					f.success(t, root, "workspace", "init")
					setArgs("--project", "unknown")
					want = "not registered"
				case "unknown-flag":
					args = append(args, "--continue")
					want = "unknown flag"
				case "positional":
					args = append(args, "extra")
					want = "no positional"
				case "format":
					args = append(args, "--format", "xml")
					want = "unsupported format"
				case "missing-git":
					empty := filepath.Join(root, "empty-bin")
					mustPlanMkdir(t, empty)
					env["PATH"] = empty
					want = "git"
				case "missing-identity":
					mustPlanWrite(t, f.config, "[user]\nuseConfigOnly = true\n", 0o600)
					want = "identity"
				case "render-size":
					args = append(args, "--goal", strings.Repeat("a", 45000))
					want = "render"
				case "inside-source":
					root = repo
					setArgs("--repo", repo, "--root", root)
					want = "inside selected"
				case "inside-common-dir":
					root = filepath.Join(repo, ".git")
					setArgs("--repo", repo, "--root", root)
					want = "inside selected"
				case "inside-other-worktree":
					root = f.repo(t, "unselected")
					setArgs("--repo", repo, "--root", root)
					want = "inside an existing"
				}
				before := planSnapshot(t, f.root)
				f.failure(t, f.root, env, want, args...)
				planAbsent(t, filepath.Join(root, "planning"))
				if !reflect.DeepEqual(before, planSnapshot(t, f.root)) {
					t.Fatal("early failure changed files/registry/source")
				}
			})
		}
	})
	t.Run("ActualCLIPartialGitFailureIsPreserved", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		repo := f.repo(t, "source")
		bin := filepath.Join(f.root, "fault-bin")
		mustPlanMkdir(t, bin)
		realGit, err := exec.LookPath("git")
		if err != nil {
			t.Fatal(err)
		}
		// The test adapter fails exactly the add subprocess. Production receives
		// no testing switch and exercises its ordinary failure path.
		script := "#!/bin/sh\nfor argument do\n  if [ \"$argument\" = add ]; then printf 'injected add failure\\n' >&2; exit 73; fi\ndone\nexec '" + strings.ReplaceAll(realGit, "'", "'\\''") + "' \"$@\"\n"
		mustPlanWrite(t, filepath.Join(bin, "git"), script, 0o755)
		env := map[string]string{"PATH": bin + string(os.PathListSeparator) + os.Getenv("PATH")}
		args := []string{"workflow", "epic", "plan", "--repo", repo, "--root", f.root, "--format", "json"}
		f.failure(t, f.root, env, "during git add; partial state preserved", args...)
		target := filepath.Join(f.root, "planning")
		before := planSnapshot(t, target)
		if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
			t.Fatal(err)
		}
		f.failure(t, f.root, env, "existing planning repository", args...)
		if !reflect.DeepEqual(before, planSnapshot(t, target)) {
			t.Fatal("partial state replayed/changed")
		}
	})
	t.Run("ConcurrentCLICreation", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		repo := f.repo(t, "source")
		before := planSnapshot(t, repo)
		var commands [2]*exec.Cmd
		var stdout, stderr [2]bytes.Buffer
		for i := range commands {
			commands[i] = exec.Command(binary, "workflow", "epic", "plan", "--repo", repo, "--root", f.root, "--format", "json")
			commands[i].Dir, commands[i].Env = f.root, f.environment(nil)
			commands[i].Stdout, commands[i].Stderr = &stdout[i], &stderr[i]
			if err := commands[i].Start(); err != nil {
				t.Fatal(err)
			}
		}
		wins := 0
		for i, c := range commands {
			err := c.Wait()
			if err == nil {
				wins++
				var r planningrepo.Result
				if err := json.Unmarshal(stdout[i].Bytes(), &r); err != nil {
					t.Fatal(err)
				}
				f.verifyCreated(t, r, "en", []string{repo})
			} else if stdout[i].Len() != 0 || !strings.Contains(stderr[i].String(), "creation stopped") {
				t.Fatalf("concurrent loser: %v %q %q", err, stdout[i].String(), stderr[i].String())
			}
		}
		if wins != 1 || !reflect.DeepEqual(before, planSnapshot(t, repo)) {
			t.Fatalf("concurrent wins=%d or source changed", wins)
		}
	})
	t.Run("GitEnvironmentCannotRedirectCreationOrRunHooks", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		repo := f.repo(t, "source")
		hooks := filepath.Join(f.root, "hooks")
		mustPlanMkdir(t, hooks)
		sentinel := filepath.Join(f.root, "hook-ran")
		script := "#!/bin/sh\nprintf ran > '" + strings.ReplaceAll(sentinel, "'", "'\\''") + "'\nexit 1\n"
		mustPlanWrite(t, filepath.Join(hooks, "pre-commit"), script, 0o755)
		mustPlanWrite(t, f.config, "[user]\nname = Plan CLI Test\nemail = cli@example.invalid\n[core]\nhooksPath = "+hooks+"\n[commit]\ngpgSign = true\n", 0o600)
		before := planSnapshot(t, repo)
		env := map[string]string{"GIT_DIR": filepath.Join(repo, ".git"), "GIT_WORK_TREE": repo, "GIT_INDEX_FILE": filepath.Join(repo, ".git", "index"), "GIT_CONFIG_COUNT": "1", "GIT_CONFIG_KEY_0": "core.hooksPath", "GIT_CONFIG_VALUE_0": hooks}
		out, stderr, exit := f.run(t, f.root, env, "workflow", "epic", "plan", "--repo", repo, "--root", f.root, "--format", "json")
		if exit != 0 {
			t.Fatalf("environment isolation: %d %s", exit, stderr)
		}
		var r planningrepo.Result
		if err := json.Unmarshal([]byte(out), &r); err != nil {
			t.Fatal(err)
		}
		f.verifyCreated(t, r, "en", []string{repo})
		planAbsent(t, sentinel)
		if !reflect.DeepEqual(before, planSnapshot(t, repo)) {
			t.Fatal("ambient Git routing mutated source")
		}
	})
	t.Run("HelpAndEnglishDocs", func(t *testing.T) {
		f := newPlanCLI(t, binary)
		for _, args := range [][]string{{"--help"}, {"workflow", "--help"}, {"workflow", "epic", "--help"}, {"workflow", "epic", "plan", "--help"}} {
			out := f.success(t, f.root, args...)
			if !strings.Contains(out, "workflow epic plan") {
				t.Fatalf("route missing from %v: %s", args, out)
			}
			if len(args) == 4 {
				for _, text := range []string{"--project", "--repo", "--root", "--goal", "--language", "--check", "--format", "(default \"en\")", "(default \"text\")", "No agent"} {
					if !strings.Contains(out, text) {
						t.Fatalf("help lacks %q: %s", text, out)
					}
				}
			}
		}
		for _, file := range []string{"../README.md", "../docs/workflow-epic-plan.md"} {
			b, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			for _, text := range []string{"ply workflow epic plan", "--repo", "--root", "--check"} {
				if !strings.Contains(string(b), text) {
					t.Fatalf("%s lacks %s", file, text)
				}
			}
		}
	})
}

type planCLIFixture struct{ binary, root, config string }

func newPlanCLI(t *testing.T, binary string) planCLIFixture {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := planCLIFixture{binary: binary, root: root, config: filepath.Join(root, "gitconfig")}
	mustPlanWrite(t, f.config, "[user]\nname = Plan CLI Test\nemail = cli@example.invalid\n", 0o600)
	return f
}

func (f planCLIFixture) environment(overrides map[string]string) []string {
	env := []string{}
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_") || key == "EMAIL" {
			continue
		}
		if _, ok := overrides[key]; !ok {
			env = append(env, entry)
		}
	}
	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+f.config, "GIT_OPTIONAL_LOCKS=0")
	for key, value := range overrides {
		env = append(env, key+"="+value)
	}
	return env
}

func (f planCLIFixture) run(t *testing.T, cwd string, overrides map[string]string, args ...string) (string, string, int) {
	t.Helper()
	c := exec.Command(f.binary, args...)
	c.Dir = cwd
	c.Env = f.environment(overrides)
	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()
	exit := 0
	if err != nil {
		var e *exec.ExitError
		if !errors.As(err, &e) {
			t.Fatal(err)
		}
		exit = e.ExitCode()
	}
	return out.String(), stderr.String(), exit
}

func (f planCLIFixture) success(t *testing.T, cwd string, args ...string) string {
	t.Helper()
	out, stderr, exit := f.run(t, cwd, nil, args...)
	if exit != 0 {
		t.Fatalf("CLI %v exited %d\n%s\n%s", args, exit, out, stderr)
	}
	return out
}

func (f planCLIFixture) failure(t *testing.T, cwd string, env map[string]string, want string, args ...string) string {
	t.Helper()
	out, stderr, exit := f.run(t, cwd, env, args...)
	if exit == 0 || out != "" || !strings.Contains(stderr, want) {
		t.Fatalf("CLI %v: exit=%d stdout=%q stderr=%q; want error %q", args, exit, out, stderr, want)
	}
	return stderr
}

func (f planCLIFixture) result(t *testing.T, cwd string, args ...string) planningrepo.Result {
	t.Helper()
	out := f.success(t, cwd, args...)
	var r planningrepo.Result
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("invalid success JSON %q: %v", out, err)
	}
	return r
}

func (f planCLIFixture) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	a := append([]string{"-c", "core.hooksPath=" + os.DevNull, "-c", "commit.gpgsign=false", "-C", dir}, args...)
	c := exec.Command("git", a...)
	c.Env = f.environment(nil)
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (f planCLIFixture) repo(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(f.root, name)
	mustPlanMkdir(t, p)
	f.git(t, p, "init", "--initial-branch=product", "--template=")
	mustPlanWrite(t, filepath.Join(p, "source.txt"), "initial\n", 0o644)
	f.git(t, p, "add", "source.txt")
	f.git(t, p, "commit", "-m", "product")
	mustPlanWrite(t, filepath.Join(p, "untracked.txt"), "keep\n", 0o644)
	return p
}

func (f planCLIFixture) verifyCreated(t *testing.T, r planningrepo.Result, language string, repos []string) {
	t.Helper()
	sort.Strings(repos)
	if r.Kind != "PlanningRepoCreation@1" || r.SchemaVersion != 1 || r.State != "created" || r.Language != language || r.Git == nil || r.Git.Branch != "main" || r.Template.ID != "agentic/planning" || r.Template.Version != "1" || r.PlanningPath != filepath.Join(r.RootPath, "planning") {
		t.Fatalf("created result: %+v", r)
	}
	if !sort.StringsAreSorted(r.Files) || len(r.Files) != 14 || len(r.Repositories) != len(repos) {
		t.Fatalf("arrays: %+v", r)
	}
	for i, p := range repos {
		if r.Repositories[i].Path != p || r.Repositories[i].GitCommonDir != filepath.Join(p, ".git") {
			t.Fatalf("repository binding: %+v", r.Repositories)
		}
	}
	if f.git(t, r.PlanningPath, "status", "--porcelain=v1") != "" || f.git(t, r.PlanningPath, "rev-list", "--count", "HEAD") != "1" || f.git(t, r.PlanningPath, "rev-parse", "HEAD") != r.Git.OID || f.git(t, r.PlanningPath, "rev-parse", "HEAD^{tree}") != r.Git.Tree || f.git(t, r.PlanningPath, "rev-parse", "--path-format=absolute", "--git-common-dir") != filepath.Join(r.PlanningPath, ".git") {
		t.Fatal("not a clean standalone initial commit")
	}
	b, err := os.ReadFile(filepath.Join(r.PlanningPath, planningrepo.ManifestName))
	if err != nil {
		t.Fatal(err)
	}
	var m planningrepo.Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.SchemaVersion != 1 || m.Language != language || m.RootPath != r.RootPath || m.PlanningPath != r.PlanningPath || m.Template != r.Template || !reflect.DeepEqual(m.Repositories, r.Repositories) || len(m.Files) != 13 {
		t.Fatalf("manifest binding: %+v", m)
	}
	for i, file := range m.Files {
		if file.Path == planningrepo.ManifestName || (i > 0 && m.Files[i-1].Path >= file.Path) {
			t.Fatal("manifest self hash/order")
		}
		b, err := os.ReadFile(filepath.Join(r.PlanningPath, file.Path))
		if err != nil {
			t.Fatal(err)
		}
		if file.SHA256 != fmt.Sprintf("sha256:%x", sha256.Sum256(b)) {
			t.Fatalf("hash differs for %s", file.Path)
		}
	}
}

func mustPlanMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}
func mustPlanWrite(t *testing.T, p, s string, mode fs.FileMode) {
	t.Helper()
	if err := os.WriteFile(p, []byte(s), mode); err != nil {
		t.Fatal(err)
	}
}
func planAbsent(t *testing.T, p string) {
	t.Helper()
	if _, err := os.Lstat(p); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unexpected path %s: %v", p, err)
	}
}
func planSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if info.Mode()&fs.ModeSymlink != 0 {
			s, err := os.Readlink(p)
			if err != nil {
				return err
			}
			value += s
		} else if !d.IsDir() {
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			value += string(b)
		}
		result[rel] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
