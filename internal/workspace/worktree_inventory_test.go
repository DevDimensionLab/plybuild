package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type inventoryPermissionFiles struct{ FileSystem }

func (inventoryPermissionFiles) EvalSymlinks(path string) (string, error) {
	return "", &os.PathError{Op: "eval-symlinks", Path: path, Err: os.ErrPermission}
}

func TestInventoryAdapterPolicies(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	oid := strings.Repeat("a", 40)
	tests := []struct {
		name, raw          string
		readOK, physicalOK bool
	}{
		{"normal", "worktree " + root + "\x00HEAD " + oid + "\x00branch refs/heads/main\x00\x00", true, true},
		{"missing", "worktree " + filepath.Join(root, "missing") + "\x00HEAD " + oid + "\x00branch refs/heads/main\x00\x00", true, false},
		{"missing prunable", "worktree " + filepath.Join(root, "missing") + "\x00HEAD " + oid + "\x00branch refs/heads/main\x00prunable missing\x00\x00", true, false},
		{"bare without head", "worktree " + root + "\x00bare\x00\x00", true, false},
		{"opaque oid", "worktree " + root + "\x00HEAD opaque\x00branch refs/heads/main\x00\x00", false, true},
		{"opaque ref", "worktree " + root + "\x00HEAD " + oid + "\x00branch opaque\x00\x00", false, true},
		{"missing ref", "worktree " + root + "\x00HEAD " + oid + "\x00\x00", false, false},
		{"detached", "worktree " + root + "\x00HEAD " + oid + "\x00detached\x00\x00", true, true},
		{"unclean path", "worktree " + root + "/../" + filepath.Base(root) + "\x00HEAD " + oid + "\x00detached\x00\x00", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseReadInventory([]byte(tt.raw))
			if (err == nil) != tt.readOK {
				t.Fatalf("read error=%v", err)
			}
			_, err = parseWorktreeInventory([]byte(tt.raw), systemFileSystem{})
			if (err == nil) != tt.physicalOK {
				t.Fatalf("physical error=%v", err)
			}
		})
	}
	raw := []byte(tests[0].raw)
	if _, err := parseWorktreeInventory(raw, inventoryPermissionFiles{systemFileSystem{}}); !errors.Is(err, os.ErrPermission) {
		t.Fatalf("permission cause lost: %v", err)
	}
	for _, parse := range []func([]byte) ([]GitWorktreeInventoryEntry, error){parseReadInventory, func(raw []byte) ([]GitWorktreeInventoryEntry, error) {
		return parseWorktreeInventory(raw, systemFileSystem{})
	}} {
		entries, err := parse(nil)
		if err != nil || entries == nil || len(entries) != 0 {
			t.Fatalf("empty=%#v, %v", entries, err)
		}
	}
}

func TestInventoryAdaptersSortAfterResolution(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	a, z := filepath.Join(root, "a"), filepath.Join(root, "z")
	for _, p := range []string{a, z} {
		if err := os.Mkdir(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	alias := filepath.Join(root, "0-alias-to-z")
	if err := os.Symlink(z, alias); err != nil {
		t.Fatal(err)
	}
	oid := strings.Repeat("a", 40)
	raw := []byte("worktree " + alias + "\x00HEAD " + oid + "\x00detached\x00\x00worktree " + a + "\x00HEAD " + oid + "\x00detached\x00\x00")
	read, err := parseReadInventory(raw)
	if err != nil || len(read) != 2 || read[0].Locator != alias || read[1].Locator != a {
		t.Fatalf("read=%#v, %v", read, err)
	}
	physical, err := parseWorktreeInventory(raw, systemFileSystem{})
	if err != nil || len(physical) != 2 || physical[0].Locator != a || physical[1].Locator != z {
		t.Fatalf("physical=%#v, %v", physical, err)
	}
}

func TestInventoryConsumersPreservePrunableFacts(t *testing.T) {
	root := initWorkItemGitRepository(t)
	common := physicalPath(t, filepath.Join(root, ".git"))
	side := filepath.Join(physicalPath(t, t.TempDir()), "side")
	runLocalGit(t, root, "worktree", "add", "-b", "side", side)
	if err := os.RemoveAll(side); err != nil {
		t.Fatal(err)
	}
	repo := RepoRecord{ID: "repo", Locator: root, GitCommonDir: common}
	reader := newSystemTaskIntegrationGit(systemFileSystem{}).(*systemTaskIntegrationGit)
	facts, err := reader.ReadWorktreeInventory(repo, nil)
	if err != nil || facts.Freshness != "fresh" || len(facts.Worktrees) != 2 {
		t.Fatalf("facts=%#v, %v", facts, err)
	}
	found := false
	for _, f := range facts.Worktrees {
		if f.Entry.Locator == side {
			found = true
			if !f.Entry.Prunable || f.Freshness.Identity != "unknown" || f.Clean != nil {
				t.Fatalf("missing checkout facts=%#v", f)
			}
		}
	}
	if !found {
		t.Fatal("missing prunable facts")
	}
	plan := TaskCloseoutPlan{RepositoryLocator: root, Source: TaskWorktreeBinding{GitCommonDir: common}}
	entries, err := closeoutInventory(plan)
	if err != nil || len(entries) != 2 {
		t.Fatalf("closeout entries=%#v, %v", entries, err)
	}
	if _, err := newSystemWorkItemGit(systemFileSystem{}).ListWorktrees(repo); err == nil {
		t.Fatal("master ListWorktrees accepted missing checkout")
	}
	if _, err := reader.ObserveIntegrationRepository(repo); err == nil {
		t.Fatal("master integration accepted missing checkout")
	}
}
