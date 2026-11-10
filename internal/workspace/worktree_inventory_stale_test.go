package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEpicAdoptWithUnrelatedPrunableWorktree(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	repository := filepath.Join(f.wrapper, "main")
	missing := filepath.Join(f.wrapper, "obsolete")
	runLocalGit(t, repository, "worktree", "add", "-b", "obsolete", missing)
	if err := os.RemoveAll(missing); err != nil {
		t.Fatal(err)
	}
	ref := runLocalGit(t, repository, "symbolic-ref", "HEAD")
	adopted, err := AdoptEpic(f.dependencies, EpicAdoptInput{EpicID: "second", Title: "Second Epic", ProjectID: "ply", RepoID: "ply", Worktree: repository, Ref: ref, ExpectedOID: f.oid})
	if err != nil || !adopted.Created {
		t.Fatalf("unrelated prunable worktree blocked clean Epic adoption: %+v, %v", adopted, err)
	}
	inventory, err := f.dependencies.WorkGit.ListWorktrees(RepoRecord{Locator: repository})
	if err != nil {
		t.Fatal(err)
	}
	entry, found := worktreeEntry(inventory, missing)
	if !found || !entry.Prunable || entry.Ref != "refs/heads/obsolete" || entry.OID != f.oid {
		t.Fatalf("prunable Git entry was lost or changed: %+v, found=%v", entry, found)
	}
	if _, err := f.dependencies.WorkGit.ObserveWorktree(missing); err == nil {
		t.Fatal("missing prunable worktree was accepted as a live target")
	}
}

type inventoryPermissionDeniedFS struct{ FileSystem }

func (f inventoryPermissionDeniedFS) EvalSymlinks(string) (string, error) {
	return "", os.ErrPermission
}

func TestWorktreeInventoryDoesNotHideUnqualifiedPathErrors(t *testing.T) {
	files := SystemDependencies().Files
	missing := filepath.Join(physicalPath(t, t.TempDir()), "missing")
	for _, tc := range []struct {
		name, path, marker string
		files              FileSystem
	}{
		{name: "missing ordinary worktree", path: missing, files: files},
		{name: "relative prunable path", path: "missing", marker: "prunable reason\x00", files: files},
		{name: "unclean prunable path", path: missing + "/../other", marker: "prunable reason\x00", files: files},
		{name: "permission is not absence", path: missing, marker: "prunable reason\x00", files: inventoryPermissionDeniedFS{files}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := "worktree " + tc.path + "\x00HEAD " + strings.Repeat("1", 40) + "\x00branch refs/heads/test\x00" + tc.marker + "\x00"
			_, err := parseWorktreeInventory([]byte(raw), tc.files)
			if err == nil {
				t.Fatal("unqualified path failure was hidden")
			}
			if tc.name == "permission is not absence" && !errors.Is(err, os.ErrPermission) {
				t.Fatalf("permission error lost: %v", err)
			}
		})
	}
}
