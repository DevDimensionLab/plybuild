package testutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWouldLeakIntoRepositoryRejectsWorkingTreeAndAcceptsTempDir(t *testing.T) {
	packagePath, err := filepath.Abs(filepath.Join("test-output", "result.txt"))
	if err != nil {
		t.Fatalf("make package test path absolute: %v", err)
	}
	if !wouldLeakIntoRepository(packagePath) {
		t.Fatalf("guard accepted package path %s", packagePath)
	}

	cwd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("make working directory absolute: %v", err)
	}
	repositoryRoot, err := findRepositoryRoot(cwd)
	if err != nil {
		t.Fatalf("find repository root: %v", err)
	}
	repositorySiblingPath := filepath.Join(repositoryRoot, "root-test-output", "result.txt")
	if !wouldLeakIntoRepository(repositorySiblingPath) {
		t.Fatalf("guard accepted repository-root sibling path %s", repositorySiblingPath)
	}

	tempPath := filepath.Join(t.TempDir(), "result.txt")
	if wouldLeakIntoRepository(tempPath) {
		t.Fatalf("guard rejected temp path %s", tempPath)
	}

	symlinkParent := filepath.Join(t.TempDir(), "worktree-link")
	if err := os.Symlink(repositoryRoot, symlinkParent); err != nil {
		t.Fatalf("create repository symlink fixture: %v", err)
	}
	if !wouldLeakIntoRepository(filepath.Join(symlinkParent, "nested", "result.txt")) {
		t.Fatalf("guard accepted path through repository symlink %s", symlinkParent)
	}

	finalSymlink := filepath.Join(t.TempDir(), "go-mod-link")
	if err := os.Symlink(filepath.Join(repositoryRoot, "go.mod"), finalSymlink); err != nil {
		t.Fatalf("create final symlink fixture: %v", err)
	}
	if !wouldLeakIntoRepository(finalSymlink) {
		t.Fatalf("guard accepted final symlink %s", finalSymlink)
	}
}

func TestSafeFixtureWritersRefuseRepositoryPaths(t *testing.T) {
	repositoryRoot, err := findRepositoryRoot(".")
	if err != nil {
		t.Fatalf("find repository root: %v", err)
	}

	repoFile := filepath.Join(repositoryRoot, ".testutil-write-leak-probe")
	assertPathAbsent(t, repoFile)
	t.Cleanup(func() { _ = os.Remove(repoFile) })
	if err := WriteFileOutsideWorkingTree(repoFile, []byte("must not be written"), 0o600); !errors.Is(err, ErrRepositoryWrite) {
		t.Fatalf("repository write returned %v, want ErrRepositoryWrite", err)
	}
	assertPathAbsent(t, repoFile)

	repoDir := filepath.Join(repositoryRoot, ".testutil-copy-leak-probe")
	assertPathAbsent(t, repoDir)
	t.Cleanup(func() { _ = os.RemoveAll(repoDir) })
	if err := CopyFSOutsideWorkingTree(repoDir, os.DirFS(t.TempDir())); !errors.Is(err, ErrRepositoryWrite) {
		t.Fatalf("repository copy returned %v, want ErrRepositoryWrite", err)
	}
	assertPathAbsent(t, repoDir)

	externalFile := filepath.Join(t.TempDir(), "result.txt")
	if err := WriteFileOutsideWorkingTree(externalFile, []byte("written outside"), 0o600); err != nil {
		t.Fatalf("write external fixture: %v", err)
	}
	contents, err := os.ReadFile(externalFile)
	if err != nil {
		t.Fatalf("read external fixture: %v", err)
	}
	if string(contents) != "written outside" {
		t.Fatalf("external fixture contained %q", contents)
	}
}

func assertPathAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("test path must not exist, stat returned %v for %s", err, path)
	}
}

func wouldLeakIntoRepository(path string) bool {
	return WouldLeakIntoWorkingTree(path)
}
