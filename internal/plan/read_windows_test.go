//go:build windows

package plan

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadRefusesSymlinkToTheOriginalFile swaps the leaf for a symlink to the
// original file under another name right before the open: the same file,
// reached through a link, is still refused.
func TestReadRefusesSymlinkToTheOriginalFile(t *testing.T) {
	store, dir := testStore(t)
	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	moved := filepath.Join(dir, "moved.data")
	store.beforeOpen = func(path string) {
		store.beforeOpen = nil
		if err := os.Rename(path, moved); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(moved, path); err != nil {
			t.Skipf("symlinks need a privilege here: %v", err)
		}
	}
	_, err := store.Show("PLAN-1", 0)
	requireClass(t, err, ClassInvalid)
	_, err = store.Show("PLAN-1", 0)
	requireClass(t, err, ClassInvalid)
}

func TestReadRefusesDirectory(t *testing.T) {
	store, dir := testStore(t)
	if err := os.MkdirAll(filepath.Join(dir, "PLAN-1.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := store.Show("PLAN-1", 0)
	requireClass(t, err, ClassInvalid)
}
