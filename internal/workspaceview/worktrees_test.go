package workspaceview

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func inventoryGit(t *testing.T, path string, args ...string) string {
	t.Helper()
	prefix := []string{"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "commit.gpgSign=false", "-c", "user.name=Inventory fixture", "-c", "user.email=inventory@example.invalid", "-C", path}
	b, err := exec.Command("git", append(prefix, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("fixture Git: %v: %s", err, b)
	}
	return strings.TrimSpace(string(b))
}
func inventoryFingerprint(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(p string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if i.Mode().IsRegular() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			out[p] = fmt.Sprintf("%x:%d:%d", sha256.Sum256(b), i.Mode(), i.ModTime().UnixNano())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func TestWorktreeInventoryKeepsOwnersUnknownsAndNoGitWrites(t *testing.T) {
	d, root := changesFixture(t)
	main := filepath.Join(root, "main")
	if err := os.Mkdir(main, 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, main, "init", "-b", "main")
	if err := os.WriteFile(filepath.Join(main, "README"), []byte("base\n"), 0600); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, main, "add", "README")
	inventoryGit(t, main, "commit", "-m", "base")
	oid := inventoryGit(t, main, "rev-parse", "HEAD")
	epic := filepath.Join(root, "epic")
	inventoryGit(t, main, "worktree", "add", "-b", "epic", epic)
	if _, err := workspace.AddProject(d, workspace.ProjectAddInput{ProjectID: "project", Name: "Inventory", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: main}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AdoptEpic(d, workspace.EpicAdoptInput{EpicID: "epic", Title: "Epic", ProjectID: "project", RepoID: "repo", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.CreateTask(d, workspace.TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
		t.Fatal(err)
	}
	task := filepath.Join(root, "task")
	if _, err := workspace.CreateTaskWorktree(d, workspace.TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: task, ExpectedParentOID: oid}); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(root, "outside Ply")
	inventoryGit(t, main, "worktree", "add", "-b", "outside", other)
	detached := filepath.Join(root, "detached")
	inventoryGit(t, main, "worktree", "add", "--detach", detached)
	missing := filepath.Join(root, "missing")
	inventoryGit(t, main, "worktree", "add", "-b", "missing", missing)
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	} // owned disposable fixture
	if err := os.WriteFile(filepath.Join(task, "README"), []byte("feature\n"), 0600); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, task, "add", "README")
	inventoryGit(t, task, "commit", "-m", "feature")
	if err := os.WriteFile(filepath.Join(task, "untracked"), []byte("dirty\n"), 0600); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "fsmonitor-ran")
	script := filepath.Join(root, "fsmonitor.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, main, "config", "core.fsmonitor", script)
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	before := inventoryFingerprint(t, root)
	out, err := ListWorktrees(d, s, WorktreeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, inventoryFingerprint(t, root)) {
		t.Fatal("inventory changed workspace or Git bytes/metadata")
	}
	if _, err := os.Lstat(marker); !os.IsNotExist(err) {
		t.Fatal("read executed fsmonitor")
	}
	if len(out.Worktrees) != 6 {
		t.Fatalf("lost inventory rows: %+v", out)
	}
	byPath := map[string]WorktreeItem{}
	for _, w := range out.Worktrees {
		byPath[w.Locator] = w
	}
	if byPath[main].Owner != "project_main" || byPath[epic].Owner != "epic" || byPath[task].Owner != "task" || byPath[other].Owner != "none" {
		t.Fatalf("ownership: %+v", byPath)
	}
	w := byPath[task]
	if w.Ahead == nil || *w.Ahead != 1 || w.Behind == nil || *w.Behind != 0 || w.Clean == nil || *w.Clean || w.MergedIntoBase == nil || *w.MergedIntoBase {
		t.Fatalf("task facts: %+v", w)
	}
	if !byPath[detached].Detached || byPath[detached].Ref != nil {
		t.Fatalf("detached: %+v", byPath[detached])
	}
	if byPath[missing].Clean != nil || byPath[missing].Freshness.Identity != "unknown" {
		t.Fatalf("missing: %+v", byPath[missing])
	}
	if byPath[other].Ahead != nil || byPath[other].MergedIntoBase != nil {
		t.Fatalf("invented unmanaged base: %+v", byPath[other])
	}
}
