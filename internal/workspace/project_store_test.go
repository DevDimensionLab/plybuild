package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectStoreRejectsLockSymlink(t *testing.T) {
	root := createProjectWorkspace(t)
	target := filepath.Join(root, "target")
	mustWriteFile(t, target, "target")
	if err := os.Symlink(target, filepath.Join(root, MarkerDirectory, projectsLock)); err != nil {
		t.Fatal(err)
	}
	err := newSystemProjectStore().WithSnapshotLock(root, func(ProjectSnapshot) error { return nil })
	assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
	contents, readErr := os.ReadFile(target)
	if readErr != nil || string(contents) != "target" {
		t.Fatalf("symlink target changed: %q, %v", contents, readErr)
	}
}
