package system

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	establish "github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/reporemove"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

type cwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f cwdFiles) Getwd() (string, error) { return f.cwd, nil }

// fixture is a workspace whose Project ply (root = workspace root) has an
// anchor repository and the repository map that repo create established,
// with a private local bare repository as its remote.
type fixture struct {
	base, root, wrapper, repo, remote, url, gitLog string
}

func physical(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=Ply tests", "-c", "user.email=tests@example.invalid"}, args...)...)
	command.Env = append(os.Environ(), "LC_ALL=C")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func write(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

// newFixture isolates Git configuration, puts a recording git wrapper first in
// PATH and runs repo create with the GitHub step answering url.
func newFixture(t *testing.T, url func(f fixture) string) fixture {
	t.Helper()
	base := physical(t, t.TempDir())
	config := filepath.Join(base, "gitconfig")
	write(t, config, "")
	t.Setenv("HOME", base)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	real, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(base, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	f := fixture{base: base, root: filepath.Join(base, "ws"), remote: filepath.Join(base, "remote.git"), gitLog: filepath.Join(base, "git.log")}
	f.wrapper = filepath.Join(f.root, "map")
	f.repo = filepath.Join(f.wrapper, "main")
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte("#!/bin/sh\nprintf '%s\\n' \"$*\" >> '"+f.gitLog+"'\nexec '"+real+"' \"$@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := os.MkdirAll(filepath.Join(f.root, "anchor", "main"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Init(workspace.Dependencies{Files: cwdFiles{workspace.SystemDependencies().Files, f.root}}); err != nil {
		t.Fatal(err)
	}
	runGit(t, filepath.Join(f.root, "anchor", "main"), "init", "-q")
	write(t, filepath.Join(f.root, workspacedef.File), `{"format_version": 1, "projects": [
  {"id": "ply", "name": "Ply", "wrapper": ".", "repos": [{"id": "anchor", "path": "anchor/main"}], "epics": []}]}
`)
	runGit(t, base, "init", "-q", "--bare", f.remote)
	f.url = f.remote
	if url != nil {
		f.url = url(f)
	}
	f.create(t)
	return f
}

// create runs repo create for map; only the GitHub step is local.
func (f fixture) create(t *testing.T) repoestablish.Result {
	t.Helper()
	deps := establish.New(process.SystemRunner().Runner)
	deps.Check = func(context.Context, repoestablish.Resolved) error { return nil }
	deps.GitHub = func(context.Context, string, string) (string, bool, error) { return f.url, false, nil }
	result, err := repoestablish.Run(context.Background(), repoestablish.Request{Workspace: f.root, Project: "ply", RepoName: "map", Actor: "tester",
		Overrides: repoestablish.Overrides{GitHubOwner: "tester"}}, deps)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func (f fixture) commit(t *testing.T, name string) string {
	t.Helper()
	write(t, filepath.Join(f.repo, name), name+"\n")
	runGit(t, f.repo, "add", name)
	runGit(t, f.repo, "commit", "-q", "-m", name)
	return runGit(t, f.repo, "rev-parse", "HEAD")
}

func (f fixture) run(t *testing.T, h hooks, req reporemove.Request) (reporemove.Result, error) {
	t.Helper()
	req.Workspace, req.Project, req.RepoName = f.root, "ply", "map"
	return reporemove.Run(context.Background(), req, newDependencies(process.SystemRunner().Runner, h))
}

func (f fixture) preview(t *testing.T) reporemove.Plan {
	t.Helper()
	result, err := f.run(t, hooks{}, reporemove.Request{})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	return result.Plan
}

func (f fixture) apply(t *testing.T, h hooks, expect string) (reporemove.Result, error) {
	t.Helper()
	return f.run(t, h, reporemove.Request{Apply: true, Expect: expect, Actor: "tester"})
}

// state records every entry under the workspace and the remote without
// following links: type, mode, size, modification time and content or target.
func (f fixture) state(t *testing.T) map[string]string {
	t.Helper()
	state := map[string]string{}
	for _, root := range []string{f.root, f.remote} {
		if _, err := os.Lstat(root); errors.Is(err, fs.ErrNotExist) {
			continue
		}
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			value := fmt.Sprintf("%s|%d|%d", info.Mode(), info.Size(), info.ModTime().UnixNano())
			switch {
			case info.Mode()&fs.ModeSymlink != 0:
				target, _ := os.Readlink(path)
				value += "|" + target
			case info.Mode().IsRegular():
				raw, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				value += fmt.Sprintf("|%x", sha256.Sum256(raw))
			}
			state[path] = value
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return state
}

func requireSame(t *testing.T, before, after map[string]string) {
	t.Helper()
	if reflect.DeepEqual(before, after) {
		return
	}
	for path, value := range before {
		if after[path] != value {
			t.Errorf("changed or removed: %s", path)
		}
	}
	for path := range after {
		if _, ok := before[path]; !ok {
			t.Errorf("added: %s", path)
		}
	}
	t.FailNow()
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("%s exists: %v", path, err)
	}
}

func requireStop(t *testing.T, err error, step reporemove.Step, kind error, reason string) {
	t.Helper()
	var stopped *reporemove.Error
	if !errors.As(err, &stopped) || stopped.Step != step || !errors.Is(err, kind) || !strings.Contains(err.Error(), reason) {
		t.Fatalf("error %v, want %v at %s with %q", err, kind, step, reason)
	}
}

func (f fixture) definedRepos(t *testing.T) []workspacedef.Repo {
	t.Helper()
	defs, _, err := workspacedef.Read(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return defs.Projects[0].Repos
}

func (f fixture) metadata(t *testing.T) workspace.ProjectMetadataSnapshot {
	t.Helper()
	d, err := workspace.AtRoot(workspace.SystemDependencies(), f.root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := workspace.ReadProjectMetadata(d, f.root)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func (f fixture) gitCalls(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile(f.gitLog)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// requireRemoved checks the full effect: wrapper gone without quarantine, the
// member removed, one clear appended and the remote unchanged.
func (f fixture) requireRemoved(t *testing.T, result reporemove.Result, err error, remote map[string]string, revision int) {
	t.Helper()
	want := []reporemove.StepResult{{Step: reporemove.StepMetadata, Changed: true}, {Step: reporemove.StepDefine, Changed: true}, {Step: reporemove.StepWrapper, Changed: true}}
	if err != nil || !reflect.DeepEqual(result.Steps, want) {
		t.Fatalf("apply %+v %v", result.Steps, err)
	}
	requireAbsent(t, f.wrapper)
	if matches, _ := filepath.Glob(filepath.Join(f.root, ".ply-repo-remove-*")); len(matches) != 0 {
		t.Fatalf("quarantine left: %v", matches)
	}
	if repos := f.definedRepos(t); len(repos) != 1 || repos[0].ID != "anchor" {
		t.Fatalf("repos %+v", repos)
	}
	snapshot := f.metadata(t)
	last := snapshot.Events[len(snapshot.Events)-1]
	if snapshot.Revision != revision+1 || last.Operation != workspace.ProjectRepoWrapperClear || *last.RepoID != "map" || snapshot.RepositoryWrapper("map") != nil {
		t.Fatalf("metadata %+v", snapshot.Events)
	}
	after := map[string]string{}
	for path, value := range f.state(t) {
		if strings.HasPrefix(path, f.remote) {
			after[path] = value
		}
	}
	requireSame(t, remote, after)
}

func (f fixture) remoteState(t *testing.T) map[string]string {
	t.Helper()
	remote := map[string]string{}
	for path, value := range f.state(t) {
		if strings.HasPrefix(path, f.remote) {
			remote[path] = value
		}
	}
	return remote
}

func TestUnbornWithEmptyRemoteRemovesLocallyWithoutContactingRemoteAndRecreates(t *testing.T) {
	f := newFixture(t, nil)
	before := f.state(t)
	revision := f.metadata(t).Revision

	plan := f.preview(t)

	requireSame(t, before, f.state(t))
	if plan.Verdict != reporemove.VerdictUnborn || plan.RemoteURL != f.url || plan.MetadataRevision != revision || plan.MetadataCleared || !strings.HasPrefix(plan.Digest, "sha256:") {
		t.Fatalf("plan %+v", plan)
	}
	if names := []string{}; len(plan.Checks) != 5 {
		for _, c := range plan.Checks {
			names = append(names, c.Name)
		}
		t.Fatalf("checks %v", names)
	}
	remote := f.remoteState(t)

	result, err := f.apply(t, hooks{}, plan.Digest)

	f.requireRemoved(t, result, err, remote, revision)
	if calls := f.gitCalls(t); strings.Contains(calls, "ls-remote") || strings.Contains(calls, "push") || strings.Contains(calls, "fetch") {
		t.Fatalf("remote contacted:\n%s", calls)
	}
	if again, err := f.run(t, hooks{}, reporemove.Request{}); err != nil || !again.Nothing {
		t.Fatalf("rerun %+v %v", again, err)
	}

	recreated := f.create(t)

	outcomes := map[repoestablish.Step]bool{}
	for _, step := range recreated.Steps {
		outcomes[step.Step] = step.Changed
	}
	if want := map[repoestablish.Step]bool{repoestablish.StepWrapper: true, repoestablish.StepInit: true, repoestablish.StepGitHub: false,
		repoestablish.StepRemote: true, repoestablish.StepDefine: true, repoestablish.StepMetadata: true}; !reflect.DeepEqual(outcomes, want) {
		t.Fatalf("recreate %+v", recreated.Steps)
	}
}

func TestUnbornWithUnreachableRemoteIsRemovable(t *testing.T) {
	f := newFixture(t, func(f fixture) string { return filepath.Join(f.base, "missing.git") })
	remote := f.remoteState(t)
	plan := f.preview(t)

	result, err := f.apply(t, hooks{}, plan.Digest)

	f.requireRemoved(t, result, err, remote, plan.MetadataRevision)
}

func TestPushedHistoryIsProvenAndOutwardLinkTargetSurvives(t *testing.T) {
	f := newFixture(t, nil)
	outside := filepath.Join(f.base, "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(outside, "keep"), "keep\n")
	if err := os.Symlink(outside, filepath.Join(f.repo, "link")); err != nil {
		t.Fatal(err)
	}
	runGit(t, f.repo, "add", "link")
	f.commit(t, "one")
	runGit(t, f.repo, "push", "-q", "origin", "main")
	runGit(t, f.repo, "tag", "-a", "-m", "v1", "v1")
	runGit(t, f.repo, "push", "-q", "origin", "v1")
	remote := f.remoteState(t)
	before := f.state(t)

	plan := f.preview(t)

	requireSame(t, before, f.state(t))
	if plan.Verdict != reporemove.VerdictProven {
		t.Fatalf("plan %+v", plan)
	}

	result, err := f.apply(t, hooks{}, plan.Digest)

	f.requireRemoved(t, result, err, remote, plan.MetadataRevision)
	if raw, err := os.ReadFile(filepath.Join(outside, "keep")); err != nil || string(raw) != "keep\n" {
		t.Fatalf("link target changed: %q %v", raw, err)
	}
}

// requireRejected previews the fixture and requires a reject without any write.
func requireRejected(t *testing.T, f fixture, kind error, reason string) {
	t.Helper()
	before := f.state(t)
	_, err := f.run(t, hooks{}, reporemove.Request{})
	requireStop(t, err, reporemove.StepCheck, kind, reason)
	requireSame(t, before, f.state(t))
}

func TestHistoryThatIsNotOnTheRemoteRejects(t *testing.T) {
	for name, tc := range map[string]struct {
		setup  func(*testing.T, fixture)
		reason string
	}{
		"unpushed commit": {reason: "not on the remote", setup: func(t *testing.T, f fixture) {
			f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			f.commit(t, "two")
		}},
		"commit only in a stale tracking ref": {reason: "not on the remote", setup: func(t *testing.T, f fixture) {
			first := f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			f.commit(t, "two")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			runGit(t, f.remote, "update-ref", "refs/heads/main", first)
			runGit(t, f.repo, "reset", "-q", "--hard", first)
			runGit(t, f.repo, "reflog", "expire", "--expire=now", "--all")
		}},
		"commit only in the reflog": {reason: "not on the remote", setup: func(t *testing.T, f fixture) {
			first := f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			f.commit(t, "two")
			runGit(t, f.repo, "reset", "-q", "--hard", first)
		}},
		"local tag absent remotely": {reason: "refs/tags/local", setup: func(t *testing.T, f fixture) {
			f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			runGit(t, f.repo, "tag", "local")
		}},
		"dangling blob": {reason: "unreferenced local objects", setup: func(t *testing.T, f fixture) {
			write(t, filepath.Join(f.repo, "f"), "lost\n")
			runGit(t, f.repo, "add", "f")
			runGit(t, f.repo, "rm", "-q", "--cached", "f")
			if err := os.Remove(filepath.Join(f.repo, "f")); err != nil {
				t.Fatal(err)
			}
		}},
		"no remote": {reason: "history unproven", setup: func(t *testing.T, f fixture) {
			f.commit(t, "one")
			runGit(t, f.repo, "remote", "remove", "origin")
		}},
		"rewritten URL": {reason: "history unproven", setup: func(t *testing.T, f fixture) {
			f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			runGit(t, f.repo, "config", "url."+filepath.Join(f.base, "other.git")+".insteadOf", f.url)
		}},
		"ls-remote fails": {reason: "history unproven", setup: func(t *testing.T, f fixture) {
			f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			if err := os.Rename(f.remote, f.remote+".moved"); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, nil)
			tc.setup(t, f)
			requireRejected(t, f, reporemove.ErrConflict, tc.reason)
		})
	}
}

func TestLocalStateAndLayoutRejectWithoutWriting(t *testing.T) {
	outside := func(f fixture) string {
		dir := filepath.Join(f.base, "outside")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
		return dir
	}
	pushed := func(t *testing.T, f fixture) {
		f.commit(t, "one")
		runGit(t, f.repo, "push", "-q", "origin", "main")
	}
	for name, tc := range map[string]struct {
		setup  func(*testing.T, fixture)
		kind   error
		reason string
	}{
		"untracked file": {kind: reporemove.ErrConflict, reason: "local changes", setup: func(t *testing.T, f fixture) {
			write(t, filepath.Join(f.repo, "new"), "new\n")
		}},
		"ignored file": {kind: reporemove.ErrConflict, reason: "local changes", setup: func(t *testing.T, f fixture) {
			if err := os.MkdirAll(filepath.Join(f.repo, ".git", "info"), 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, filepath.Join(f.repo, ".git", "info", "exclude"), "ignored\n")
			write(t, filepath.Join(f.repo, "ignored"), "x\n")
		}},
		"staged file": {kind: reporemove.ErrConflict, reason: "local changes", setup: func(t *testing.T, f fixture) {
			write(t, filepath.Join(f.repo, "new"), "new\n")
			runGit(t, f.repo, "add", "new")
		}},
		"stash": {kind: reporemove.ErrConflict, reason: "refs/stash", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			write(t, filepath.Join(f.repo, "one"), "changed\n")
			runGit(t, f.repo, "stash", "-q")
		}},
		"merge in progress": {kind: reporemove.ErrConflict, reason: "MERGE_HEAD", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			write(t, filepath.Join(f.repo, ".git", "MERGE_HEAD"), runGit(t, f.repo, "rev-parse", "HEAD")+"\n")
		}},
		"extra worktree": {kind: reporemove.ErrConflict, reason: "worktree", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			runGit(t, f.repo, "worktree", "add", "-q", "-b", "other", filepath.Join(f.base, "other"))
		}},
		"gitfile": {kind: reporemove.ErrInvalid, reason: ".git", setup: func(t *testing.T, f fixture) {
			store := filepath.Join(f.base, "store.git")
			if err := os.Rename(filepath.Join(f.repo, ".git"), store); err != nil {
				t.Fatal(err)
			}
			write(t, filepath.Join(f.repo, ".git"), "gitdir: "+store+"\n")
		}},
		"symlinked wrapper": {kind: reporemove.ErrInvalid, reason: "symlink", setup: func(t *testing.T, f fixture) {
			real := filepath.Join(f.base, "real-wrapper")
			if err := os.Rename(f.wrapper, real); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(real, f.wrapper); err != nil {
				t.Fatal(err)
			}
		}},
		"extra wrapper entry": {kind: reporemove.ErrConflict, reason: "extra", setup: func(t *testing.T, f fixture) {
			write(t, filepath.Join(f.wrapper, "extra"), "x\n")
		}},
		"tracked nested repository with local history": {kind: reporemove.ErrConflict, reason: "nested", setup: func(t *testing.T, f fixture) {
			child := filepath.Join(f.repo, "sub")
			runGit(t, f.base, "init", "-q", child)
			runGit(t, child, "commit", "-q", "--allow-empty", "-m", "child one")
			runGit(t, f.repo, "add", "sub")
			f.commit(t, "one")
			runGit(t, f.repo, "push", "-q", "origin", "main")
			runGit(t, child, "commit", "-q", "--allow-empty", "-m", "child unpublished")
			runGit(t, f.repo, "config", "diff.ignoreSubmodules", "all")
		}},
		"modified skip-worktree file": {kind: reporemove.ErrConflict, reason: "index entry", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			runGit(t, f.repo, "update-index", "--skip-worktree", "one")
			write(t, filepath.Join(f.repo, "one"), "local override\n")
		}},
		"modified assume-unchanged file": {kind: reporemove.ErrConflict, reason: "index entry", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			runGit(t, f.repo, "update-index", "--assume-unchanged", "one")
			write(t, filepath.Join(f.repo, "one"), "local override\n")
		}},
		"LFS object store": {kind: reporemove.ErrConflict, reason: "lfs", setup: func(t *testing.T, f fixture) {
			pushed(t, f)
			store := filepath.Join(f.repo, ".git", "lfs", "objects", "ab")
			if err := os.MkdirAll(store, 0o755); err != nil {
				t.Fatal(err)
			}
			write(t, filepath.Join(store, "abcd"), "large local content\n")
		}},
		"untracked outward symlink": {kind: reporemove.ErrConflict, reason: "local changes", setup: func(t *testing.T, f fixture) {
			write(t, filepath.Join(outside(f), "keep"), "keep\n")
			if err := os.Symlink(outside(f), filepath.Join(f.repo, "link")); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, nil)
			tc.setup(t, f)
			external := map[string]string{}
			for _, dir := range []string{filepath.Join(f.base, "outside"), filepath.Join(f.base, "real-wrapper"), filepath.Join(f.base, "store.git")} {
				if _, err := os.Lstat(dir); err == nil {
					external[dir] = fmt.Sprint(listTree(t, dir))
				}
			}
			requireRejected(t, f, tc.kind, tc.reason)
			for dir, before := range external {
				if after := fmt.Sprint(listTree(t, dir)); after != before {
					t.Fatalf("%s changed", dir)
				}
			}
		})
	}
}

func listTree(t *testing.T, root string) []string {
	t.Helper()
	var entries []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		entries = append(entries, fmt.Sprintf("%s %s %d %d", path, info.Mode(), info.Size(), info.ModTime().UnixNano()))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

// fsckRunner answers git fsck with a fixed result and runs everything else.
type fsckRunner struct {
	process.Runner
	stdout, stderr string
	code           int
}

func (r fsckRunner) Run(c process.Command) error {
	fsck := false
	for _, arg := range c.Args {
		fsck = fsck || arg == "fsck"
	}
	if !fsck {
		return r.Runner.Run(c)
	}
	_, _ = io.WriteString(c.Stdout, r.stdout)
	_, _ = io.WriteString(c.Stderr, r.stderr)
	if r.code != 0 {
		return exec.Command("/bin/sh", "-c", "exit "+strconv.Itoa(r.code)).Run()
	}
	return nil
}

func TestFsckOutputIsParsedStrictly(t *testing.T) {
	for name, tc := range map[string]struct {
		runner fsckRunner
		ok     bool
	}{
		"notice only":         {runner: fsckRunner{stderr: "notice: No default references\nnotice: HEAD points to an unborn branch (main)\n"}, ok: true},
		"unknown stdout line": {runner: fsckRunner{stdout: "broken link from tree 1\n"}},
		"other stderr line":   {runner: fsckRunner{stderr: "error: object file is empty\n"}},
		"non-zero exit":       {runner: fsckRunner{code: 2}},
		"short object id":     {runner: fsckRunner{stdout: "dangling blob 45b983be\n"}},
		"unterminated output": {runner: fsckRunner{stdout: "dangling blob 45b983be36b73c0788dc9cbcb76cbb80fc7bb057"}},
		"unknown object type": {runner: fsckRunner{stdout: "dangling note 45b983be36b73c0788dc9cbcb76cbb80fc7bb057\n"}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, nil)
			tc.runner.Runner = process.SystemRunner().Runner
			before := f.state(t)

			result, err := reporemove.Run(context.Background(), reporemove.Request{Workspace: f.root, Project: "ply", RepoName: "map"}, newDependencies(tc.runner, hooks{}))

			requireSame(t, before, f.state(t))
			if tc.ok {
				if err != nil || result.Plan.Verdict != reporemove.VerdictUnborn {
					t.Fatalf("result %+v %v", result.Plan, err)
				}
				return
			}
			requireStop(t, err, reporemove.StepCheck, reporemove.ErrInvalid, "fsck")
		})
	}
}

func TestChangesAfterPreviewConflictWithoutWriting(t *testing.T) {
	for name, change := range map[string]func(*testing.T, fixture){
		"metadata revision": func(t *testing.T, f fixture) {
			d, err := workspace.AtRoot(workspace.SystemDependencies(), f.root)
			if err != nil {
				t.Fatal(err)
			}
			input := workspace.ProjectMetadataInput{Operation: workspace.ProjectCompanionAdd, ProjectID: "ply", Role: workspace.CompanionDocs,
				Locator: filepath.Join(f.root, "anchor"), ActorClaim: "tester"}
			if _, err := workspace.UpdateProjectMetadata(d, input); err != nil {
				t.Fatal(err)
			}
		},
		"definitions": func(t *testing.T, f fixture) {
			raw, err := os.ReadFile(filepath.Join(f.root, workspacedef.File))
			if err != nil {
				t.Fatal(err)
			}
			write(t, filepath.Join(f.root, workspacedef.File), string(bytes.TrimSpace(raw))+"\n\n")
		},
		"wrapper file": func(t *testing.T, f fixture) {
			write(t, filepath.Join(f.repo, ".git", "info-extra"), "x\n")
		},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, nil)
			plan := f.preview(t)
			change(t, f)
			before := f.state(t)

			_, err := f.apply(t, hooks{}, plan.Digest)

			requireStop(t, err, reporemove.StepCheck, reporemove.ErrConflict, "")
			requireSame(t, before, f.state(t))
		})
	}
}

func TestReRegistrationBeforeRenameKeepsWrapper(t *testing.T) {
	f := newFixture(t, nil)
	plan := f.preview(t)
	readd := func() {
		d, err := workspace.AtRoot(workspace.SystemDependencies(), f.root)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := workspace.AddDefinitionRepo(d, "ply", "map", "map/main"); err != nil {
			t.Fatal(err)
		}
	}

	result, err := f.apply(t, hooks{beforeWrapper: readd}, plan.Digest)

	requireStop(t, err, reporemove.StepWrapper, reporemove.ErrConflict, "")
	if want := []reporemove.StepResult{{Step: reporemove.StepMetadata, Changed: true}, {Step: reporemove.StepDefine, Changed: true}}; !reflect.DeepEqual(result.Steps, want) {
		t.Fatalf("steps %+v", result.Steps)
	}
	if info, err := os.Lstat(filepath.Join(f.repo, ".git")); err != nil || !info.IsDir() {
		t.Fatalf("repository lost: %v", err)
	}
	if matches, _ := filepath.Glob(filepath.Join(f.root, ".ply-repo-remove-*")); len(matches) != 0 {
		t.Fatalf("quarantine created: %v", matches)
	}
}

func TestEarlierQuarantineHaltsWithoutWriting(t *testing.T) {
	f := newFixture(t, nil)
	earlier := filepath.Join(f.root, ".ply-repo-remove-map-x")
	if err := os.Mkdir(earlier, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(earlier, "kept"), "kept\n")

	requireRejected(t, f, reporemove.ErrConflict, earlier)
}

func TestResumeAfterMetadataClearReportsUnchanged(t *testing.T) {
	f := newFixture(t, nil)
	revision := f.metadata(t).Revision
	d, err := workspace.AtRoot(workspace.SystemDependencies(), f.root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.ClearRepositoryWrapperAt(d, f.root, "ply", "map", "tester", revision); err != nil {
		t.Fatal(err)
	}
	plan := f.preview(t)
	if !plan.MetadataCleared || plan.MetadataRevision != revision+1 {
		t.Fatalf("plan %+v", plan)
	}

	result, err := f.apply(t, hooks{}, plan.Digest)

	want := []reporemove.StepResult{{Step: reporemove.StepMetadata}, {Step: reporemove.StepDefine, Changed: true}, {Step: reporemove.StepWrapper, Changed: true}}
	if err != nil || !reflect.DeepEqual(result.Steps, want) {
		t.Fatalf("apply %+v %v", result.Steps, err)
	}
	requireAbsent(t, f.wrapper)
	if f.metadata(t).Revision != revision+1 {
		t.Fatal("metadata written again")
	}
}

func TestFailedDeletionKeepsQuarantineAndHaltsNextRun(t *testing.T) {
	f := newFixture(t, nil)
	plan := f.preview(t)

	result, err := f.apply(t, hooks{removeAll: func(string) error { return errors.New("injected") }}, plan.Digest)

	quarantine := filepath.Join(f.root, ".ply-repo-remove-map-"+strings.TrimPrefix(plan.Digest, "sha256:")[:12])
	requireStop(t, err, reporemove.StepWrapper, reporemove.ErrUnknown, quarantine)
	if len(result.Steps) != 2 {
		t.Fatalf("steps %+v", result.Steps)
	}
	requireAbsent(t, f.wrapper)
	if info, err := os.Lstat(filepath.Join(quarantine, "main", ".git")); err != nil || !info.IsDir() {
		t.Fatalf("quarantine %v", err)
	}

	requireRejected(t, f, reporemove.ErrConflict, quarantine)
}

func TestWalkFailsClosedWithoutDevice(t *testing.T) {
	f := newFixture(t, nil)
	r := reporemove.Resolved{Wrapper: f.wrapper, Leaf: "main"}
	if _, _, err := walk(r, deviceOf); err != nil {
		t.Fatalf("walk with the platform device: %v", err)
	}

	_, _, err := walk(r, func(os.FileInfo) (uint64, bool) { return 0, false })

	if !errors.Is(err, reporemove.ErrInvalid) || !strings.Contains(err.Error(), "device") {
		t.Fatalf("walk without device: %v", err)
	}
}

func TestConfiguredFsmonitorIsNeverRun(t *testing.T) {
	f := newFixture(t, nil)
	f.commit(t, "one")
	runGit(t, f.repo, "push", "-q", "origin", "main")
	called := filepath.Join(f.base, "monitor-called")
	monitor := filepath.Join(f.base, "monitor")
	if err := os.WriteFile(monitor, []byte("#!/bin/sh\n: > '"+called+"'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGit(t, f.repo, "config", "core.fsmonitor", monitor)
	before := f.state(t)

	_, err := f.run(t, hooks{}, reporemove.Request{})

	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	requireAbsent(t, called)
	requireSame(t, before, f.state(t))
}
