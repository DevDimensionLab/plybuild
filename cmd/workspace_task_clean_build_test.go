package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanBuildEmptiesOnlyTheBuildDirectory(t *testing.T) {
	worktree := t.TempDir()
	outside := filepath.Join(t.TempDir(), "keep.txt")
	build := filepath.Join(worktree, ".build")
	for _, path := range []string{outside, filepath.Join(worktree, "src.go"), filepath.Join(build, "a.out"), filepath.Join(build, "cache", "x")} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Dir(outside), filepath.Join(build, "link")); err != nil {
		t.Fatal(err)
	}
	if n, err := cleanBuild(worktree); err != nil || n != 3 {
		t.Fatalf("clean %d %v", n, err)
	}
	if entries, err := os.ReadDir(build); err != nil || len(entries) != 0 {
		t.Fatalf(".build after clean: %v %v", entries, err)
	}
	for _, kept := range []string{outside, filepath.Join(worktree, "src.go")} {
		if _, err := os.Stat(kept); err != nil {
			t.Fatalf("%s removed", kept)
		}
	}
	if n, err := cleanBuild(t.TempDir()); err != nil || n != 0 {
		t.Fatalf("no .build: %d %v", n, err)
	}
	linked := t.TempDir()
	if err := os.Symlink(filepath.Dir(outside), filepath.Join(linked, ".build")); err != nil {
		t.Fatal(err)
	}
	if _, err := cleanBuild(linked); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("symlinked .build: %v", err)
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatal("symlink target removed")
	}
}
