package workflowhandoff

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestStoreCreateReadAndIdempotentRetry(t *testing.T) {
	root := physicalTempDir(t)
	store := newSystemStore(systemFileSystem{}, systemClock{})
	handoff := testStoreHandoff(t, root)
	input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: handoff}
	first, err := store.Create(input)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Created || first.Snapshot.Handoff.Identity.HandoffID != handoff.Identity.HandoffID {
		t.Fatalf("first create = %#v", first)
	}
	info, err := os.Stat(first.Snapshot.Handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0o400 {
		t.Fatalf("handoff mode = %o", info.Mode().Perm())
	}
	second, err := store.Create(input)
	if err != nil {
		t.Fatal(err)
	}
	if second.Created || second.Snapshot.Handoff.Locator != first.Snapshot.Handoff.Locator {
		t.Fatalf("retry = %#v", second)
	}
	read, err := store.ReadByID(root, HandoffID(handoff.Identity.HandoffID))
	if err != nil || read.Handoff.SHA256 != handoff.SHA256 {
		t.Fatalf("read = %#v, %v", read, err)
	}

	conflict := input
	conflict.SourceDraftSHA256 = digestBytes([]byte("different"))
	conflict.Handoff.SourceDraftSHA256 = conflict.SourceDraftSHA256
	if _, err := store.Create(conflict); !IsClass(err, ErrorConflict) {
		t.Fatalf("publication conflict error = %v", err)
	}
}

func TestStoreAtomicWriteFaultsLeaveAbsentOrCompleteAuthority(t *testing.T) {
	tests := []struct {
		name   string
		inject func(*systemStore)
	}{
		{name: "temp open", inject: func(store *systemStore) {
			store.ops.createTemp = func(string, string) (*os.File, error) { return nil, errors.New("injected temp failure") }
		}},
		{name: "short write", inject: func(store *systemStore) { store.ops.write = func(*os.File, []byte) (int, error) { return 0, nil } }},
		{name: "write", inject: func(store *systemStore) {
			store.ops.write = func(*os.File, []byte) (int, error) { return 0, errors.New("injected write failure") }
		}},
		{name: "chmod", inject: func(store *systemStore) {
			store.ops.chmodFile = func(*os.File, os.FileMode) error { return errors.New("injected chmod failure") }
		}},
		{name: "file sync", inject: func(store *systemStore) {
			store.ops.syncFile = func(*os.File) error { return errors.New("injected sync failure") }
		}},
		{name: "close", inject: func(store *systemStore) {
			store.ops.closeFile = func(file *os.File) error { _ = file.Close(); return errors.New("injected close failure") }
		}},
		{name: "rename", inject: func(store *systemStore) {
			store.ops.rename = func(string, string) error { return errors.New("injected rename failure") }
		}},
		{name: "directory sync", inject: func(store *systemStore) {
			store.ops.syncDirectory = func(string) error { return errors.New("injected directory sync failure") }
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			test.inject(store)
			destination := filepath.Join(directory, "authority.json")
			contents := []byte(`{"complete":true}`)
			if _, err := store.atomicWriteOnce(destination, contents, 0o400); !IsClass(err, ErrorIO) {
				t.Fatalf("error = %v", err)
			}
			actual, err := os.ReadFile(destination)
			if err == nil && string(actual) != string(contents) {
				t.Fatalf("partial destination = %q", actual)
			}
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Fatal(err)
			}
		})
	}
}

func TestStoreRejectsSymlinkAuthoritiesWithoutFollowingOrRepairingThem(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink authority assertions are Unix-specific")
	}
	t.Run("immutable destination", func(t *testing.T) {
		directory := physicalTempDir(t)
		external := filepath.Join(directory, "external.json")
		contents := []byte(`{"complete":true}`)
		if err := os.WriteFile(external, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		destination := filepath.Join(directory, "authority.json")
		if err := os.Symlink(external, destination); err != nil {
			t.Fatal(err)
		}
		store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		if _, err := store.atomicWriteOnce(destination, contents, 0o400); !IsClass(err, ErrorWorkspaceConflict) {
			t.Fatalf("symlink destination error = %v", err)
		}
		actual, err := os.ReadFile(external)
		if err != nil || !bytes.Equal(actual, contents) {
			t.Fatalf("external bytes changed: %q, %v", actual, err)
		}
	})

	t.Run("lock path", func(t *testing.T) {
		directory := physicalTempDir(t)
		external := filepath.Join(directory, "external.lock")
		if err := os.WriteFile(external, nil, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(external, 0o644); err != nil {
			t.Fatal(err)
		}
		lockPath := filepath.Join(directory, "authority.lock")
		if err := os.Symlink(external, lockPath); err != nil {
			t.Fatal(err)
		}
		store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		if _, err := store.acquire(lockPath); !IsClass(err, ErrorWorkspaceConflict) {
			t.Fatalf("symlink lock error = %v", err)
		}
		info, err := os.Stat(external)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o644 {
			t.Fatalf("external lock mode changed to %o", info.Mode().Perm())
		}
	})
}

func TestAtomicWriteOnceRejectsDestinationSwappedToSymlinkAfterInspection(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink swap assertion is Unix-specific")
	}
	directory := physicalTempDir(t)
	contents := []byte(`{"complete":true}`)
	destination := filepath.Join(directory, "authority.json")
	backup := filepath.Join(directory, "authority.backup")
	external := filepath.Join(directory, "external.json")
	for _, path := range []string{destination, external} {
		if err := os.WriteFile(path, contents, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
	originalLstat := store.ops.lstat
	swapped := false
	store.ops.lstat = func(path string) (os.FileInfo, error) {
		info, err := originalLstat(path)
		if path == destination && err == nil && !swapped {
			swapped = true
			if err := os.Rename(destination, backup); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(external, destination); err != nil {
				t.Fatal(err)
			}
		}
		return info, err
	}
	if _, err := store.atomicWriteOnce(destination, contents, 0o400); !IsClass(err, ErrorWorkspaceConflict) {
		t.Fatalf("swapped destination error = %v", err)
	}
	actual, err := os.ReadFile(external)
	if err != nil || !bytes.Equal(actual, contents) {
		t.Fatalf("external bytes changed: %q, %v", actual, err)
	}
}

func TestCreateRetryDoesNotReopenActivityWhenHeadIsMissingAfterLaterEvent(t *testing.T) {
	root := physicalTempDir(t)
	store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
	input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}
	created, err := store.Create(input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.publishActivityEvent(storeRoot(root), created.Snapshot.Handoff, "cancelled", nil, "A later lifecycle transition.", nil, nil); err != nil {
		t.Fatal(err)
	}
	headPath := filepath.Join(storeRoot(root), "activities", input.Handoff.Identity.ActivityID, "head.ref")
	if err := os.Remove(headPath); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Create(input); !IsClass(err, ErrorWorkspaceConflict) {
		t.Fatalf("retry with later event error = %v", err)
	}
	if _, err := os.Lstat(headPath); !os.IsNotExist(err) {
		t.Fatalf("retry recreated head from ambiguous history: %v", err)
	}
}

func TestStoreExactLayoutModesAndDirectoryAuthorityTypes(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix mode and symlink assertions")
	}
	t.Run("complete tree modes", func(t *testing.T) {
		root := physicalTempDir(t)
		store := newSystemStore(systemFileSystem{}, systemClock{})
		created, err := store.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
		if err != nil {
			t.Fatal(err)
		}
		start := testAcceptedDocument(t, created.Snapshot.Handoff.Identity.StartReceiptID, "started")
		if _, err := store.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: start}); err != nil {
			t.Fatal(err)
		}
		snapshot, err := store.ReadByLocator(created.Snapshot.Handoff.Locator)
		if err != nil {
			t.Fatal(err)
		}
		terminal := testAcceptedDocument(t, snapshot.Handoff.Identity.TerminalResultID, "complete")
		if _, err := store.SubmitResult(SubmitStoreInput{Snapshot: snapshot, Phase: "terminal", Document: terminal}); err != nil {
			t.Fatal(err)
		}
		if err := filepath.Walk(storeRoot(root), func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if info.Mode()&os.ModeSymlink != 0 {
				t.Errorf("store path is a symlink: %s", path)
				return nil
			}
			if info.IsDir() {
				if info.Mode().Perm() != 0o700 {
					t.Errorf("directory %s mode = %o", path, info.Mode().Perm())
				}
				return nil
			}
			want := os.FileMode(0o400)
			if strings.HasSuffix(path, ".lock") || filepath.Base(path) == "head.ref" {
				want = 0o600
			}
			if info.Mode().Perm() != want {
				t.Errorf("file %s mode = %o, want %o", path, info.Mode().Perm(), want)
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	})

	for _, kind := range []string{"symlink", "regular file"} {
		t.Run(kind+" store root", func(t *testing.T) {
			root := physicalTempDir(t)
			storePath := storeRoot(root)
			if err := os.MkdirAll(filepath.Dir(storePath), 0o700); err != nil {
				t.Fatal(err)
			}
			if kind == "symlink" {
				external := filepath.Join(root, "external")
				if err := os.Mkdir(external, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(external, storePath); err != nil {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(storePath, []byte("wrong type"), 0o600); err != nil {
				t.Fatal(err)
			}
			store := newSystemStore(systemFileSystem{}, systemClock{})
			input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}
			if _, err := store.Create(input); !IsClass(err, ErrorWorkspaceConflict) {
				t.Fatalf("wrong store root type error = %v", err)
			}
		})
	}
}

func TestStoreCreateInjectsDirectoryTempAndLockFaults(t *testing.T) {
	tests := []struct {
		name   string
		inject func(*systemStore)
	}{
		{name: "mkdir", inject: func(store *systemStore) {
			store.ops.mkdirAll = func(string, os.FileMode) error { return errors.New("injected mkdir failure") }
		}},
		{name: "run temp", inject: func(store *systemStore) {
			store.ops.mkdirTemp = func(string, string) (string, error) { return "", errors.New("injected run temp failure") }
		}},
		{name: "lock", inject: func(store *systemStore) {
			store.ops.lock = func(*os.File) error { return errors.New("injected lock failure") }
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := physicalTempDir(t)
			store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			test.inject(store)
			input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}
			if _, err := store.Create(input); !IsClass(err, ErrorIO) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestStoreLockOpenIdentityUnlockAndCloseFaults(t *testing.T) {
	t.Run("lock open", func(t *testing.T) {
		root := physicalTempDir(t)
		store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		original := store.ops.openFile
		calls := 0
		store.ops.openFile = func(path string, flag int, mode os.FileMode) (*os.File, error) {
			calls++
			if calls == 2 {
				return nil, errors.New("injected lock reopen failure")
			}
			return original(path, flag, mode)
		}
		if _, err := store.acquire(filepath.Join(root, "lock")); !IsClass(err, ErrorIO) {
			t.Fatalf("lock open error = %v", err)
		}
	})

	t.Run("handle path identity", func(t *testing.T) {
		root := physicalTempDir(t)
		lockPath := filepath.Join(root, "lock")
		otherPath := filepath.Join(root, "other")
		if err := os.WriteFile(otherPath, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		otherInfo, err := os.Stat(otherPath)
		if err != nil {
			t.Fatal(err)
		}
		store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		originalStat := store.ops.statFile
		statCalls := 0
		store.ops.statFile = func(file *os.File) (os.FileInfo, error) {
			statCalls++
			if statCalls == 2 {
				return otherInfo, nil
			}
			return originalStat(file)
		}
		if _, err := store.acquire(lockPath); !IsClass(err, ErrorWorkspaceConflict) {
			t.Fatalf("identity error = %v", err)
		}
	})

	for _, test := range []struct {
		name       string
		inject     func(*systemStore, *bool)
		wantClosed bool
	}{
		{name: "unlock", inject: func(store *systemStore, closed *bool) {
			originalClose := store.ops.closeFile
			store.ops.unlock = func(*os.File) error { return errors.New("injected unlock failure") }
			store.ops.closeFile = func(file *os.File) error { *closed = true; return originalClose(file) }
		}, wantClosed: true},
		{name: "close", inject: func(store *systemStore, closed *bool) {
			store.ops.closeFile = func(file *os.File) error {
				*closed = true
				_ = file.Close()
				return errors.New("injected close failure")
			}
		}, wantClosed: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := physicalTempDir(t)
			store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			lock, err := store.acquire(filepath.Join(root, "lock"))
			if err != nil {
				t.Fatal(err)
			}
			closed := false
			test.inject(store, &closed)
			if err := store.release(lock); !IsClass(err, ErrorIO) {
				t.Fatalf("release error = %v", err)
			}
			if test.wantClosed && !closed {
				t.Fatal("release did not attempt close")
			}
		})
	}
}

func TestStoreReconcilesCompleteRunAfterInterruptedPublication(t *testing.T) {
	for _, stage := range []string{"runs directory sync", "publication ref"} {
		t.Run(stage, func(t *testing.T) {
			root := physicalTempDir(t)
			store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}
			if stage == "runs directory sync" {
				original := store.ops.syncDirectory
				failed := false
				store.ops.syncDirectory = func(path string) error {
					if !failed && path == filepath.Join(storeRoot(root), "runs") {
						failed = true
						return errors.New("injected post-rename sync failure")
					}
					return original(path)
				}
			} else {
				original := store.ops.createTemp
				store.ops.createTemp = func(directory, pattern string) (*os.File, error) {
					if directory == filepath.Join(storeRoot(root), "publications") {
						return nil, errors.New("injected publication failure")
					}
					return original(directory, pattern)
				}
			}
			if _, err := store.Create(input); !IsClass(err, ErrorIO) {
				t.Fatalf("interrupted create error = %v", err)
			}
			store = newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			reconciled, err := store.Create(input)
			if err != nil {
				t.Fatal(err)
			}
			if reconciled.Created || reconciled.Snapshot.Handoff.Identity.HandoffID != input.Handoff.Identity.HandoffID {
				t.Fatalf("reconciled result = %#v", reconciled)
			}
			entries, err := os.ReadDir(filepath.Join(storeRoot(root), "runs"))
			if err != nil {
				t.Fatal(err)
			}
			count := 0
			for _, entry := range entries {
				if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
					count++
				}
			}
			if count != 1 {
				t.Fatalf("complete run count = %d", count)
			}
		})
	}
}

func TestStoreCreateEventPublicationAndHeadFaultsReconcileOnRetry(t *testing.T) {
	for _, stage := range []string{"event temp", "event directory sync", "publication temp", "publication directory sync", "head temp", "head replace", "head directory sync"} {
		t.Run(stage, func(t *testing.T) {
			root := physicalTempDir(t)
			store := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}
			failed := false
			if strings.Contains(stage, "temp") {
				original := store.ops.createTemp
				store.ops.createTemp = func(directory, pattern string) (*os.File, error) {
					base := filepath.Base(directory)
					match := (stage == "event temp" && base == "events") ||
						(stage == "publication temp" && base == "publications") ||
						(stage == "head temp" && base == input.Handoff.Identity.ActivityID)
					if match && !failed {
						failed = true
						return nil, errors.New("injected " + stage)
					}
					return original(directory, pattern)
				}
			} else if stage == "head replace" {
				original := store.ops.replace
				store.ops.replace = func(source, destination string) error {
					if filepath.Base(destination) == "head.ref" && !failed {
						failed = true
						return errors.New("injected head replace")
					}
					return original(source, destination)
				}
			} else {
				original := store.ops.syncDirectory
				store.ops.syncDirectory = func(directory string) error {
					base := filepath.Base(directory)
					match := (stage == "event directory sync" && base == "events") ||
						(stage == "publication directory sync" && base == "publications") ||
						(stage == "head directory sync" && base == input.Handoff.Identity.ActivityID)
					if match && !failed {
						failed = true
						return errors.New("injected " + stage)
					}
					return original(directory)
				}
			}
			if _, err := store.Create(input); !IsClass(err, ErrorIO) {
				t.Fatalf("faulted create error = %v", err)
			}
			fresh := newSystemStore(systemFileSystem{}, systemClock{})
			retried, err := fresh.Create(input)
			if err != nil {
				t.Fatalf("exact retry failed: %v", err)
			}
			if retried.Snapshot.HeadState != "ready" {
				t.Fatalf("retry snapshot = %#v", retried.Snapshot)
			}
			activity := filepath.Join(storeRoot(root), "activities", input.Handoff.Identity.ActivityID)
			if _, err := os.Stat(filepath.Join(activity, "head.ref")); err != nil {
				t.Fatalf("head ref missing after retry: %v", err)
			}
			events, err := os.ReadDir(filepath.Join(activity, "events"))
			if err != nil || len(events) != 1 {
				t.Fatalf("activity events = %d, %v", len(events), err)
			}
		})
	}
}

func TestStoreAcceptsOnlyOneStartDocument(t *testing.T) {
	root := physicalTempDir(t)
	store := newSystemStore(systemFileSystem{}, systemClock{})
	created, err := store.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
	if err != nil {
		t.Fatal(err)
	}
	document := testAcceptedDocument(t, "rcp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "started")
	first, err := store.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: document})
	if err != nil || !first.Created {
		t.Fatalf("first submit = %#v, %v", first, err)
	}
	retry, err := store.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: document})
	if err != nil || retry.Created || retry.Document.Locator != first.Document.Locator {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	other := testAcceptedDocument(t, "rcp_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "conflict")
	if _, err := store.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: other}); !IsClass(err, ErrorConflict) {
		t.Fatalf("competing submit error = %v", err)
	}
}

func TestStoreReplyPublicationFaultsReconcileOnExactRetry(t *testing.T) {
	for _, phase := range []string{"start", "terminal"} {
		for _, stage := range []string{"object temp", "event temp", "event directory sync", "accepted ref temp", "accepted ref directory sync"} {
			t.Run(phase+"/"+stage, func(t *testing.T) {
				root := physicalTempDir(t)
				base := newSystemStore(systemFileSystem{}, systemClock{})
				created, err := base.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
				if err != nil {
					t.Fatal(err)
				}
				snapshot := created.Snapshot
				if phase == "terminal" {
					start := testAcceptedDocument(t, snapshot.Handoff.Identity.StartReceiptID, "started")
					if _, err := base.SubmitStart(SubmitStoreInput{Snapshot: snapshot, Phase: "start", Document: start}); err != nil {
						t.Fatal(err)
					}
					snapshot, err = base.ReadByLocator(snapshot.Handoff.Locator)
					if err != nil {
						t.Fatal(err)
					}
				}
				id := snapshot.Handoff.Identity.StartReceiptID
				outcome := "started"
				if phase == "terminal" {
					id = snapshot.Handoff.Identity.TerminalResultID
					outcome = "complete"
				}
				document := testAcceptedDocument(t, id, outcome)
				faulted := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
				failed := false
				if strings.Contains(stage, "temp") {
					original := faulted.ops.createTemp
					faulted.ops.createTemp = func(directory, pattern string) (*os.File, error) {
						match := (stage == "object temp" && filepath.Base(directory) == "objects") ||
							(stage == "event temp" && filepath.Base(directory) == "events") ||
							(stage == "accepted ref temp" && filepath.Base(directory) == phase)
						if match && !failed {
							failed = true
							return nil, errors.New("injected " + stage)
						}
						return original(directory, pattern)
					}
				} else {
					original := faulted.ops.syncDirectory
					faulted.ops.syncDirectory = func(directory string) error {
						match := (stage == "event directory sync" && filepath.Base(directory) == "events") ||
							(stage == "accepted ref directory sync" && filepath.Base(directory) == phase)
						if match && !failed {
							failed = true
							return errors.New("injected " + stage)
						}
						return original(directory)
					}
				}
				input := SubmitStoreInput{Snapshot: snapshot, Phase: phase, Document: document}
				var submit func(SubmitStoreInput) (SubmitStoreResult, error)
				if phase == "start" {
					submit = faulted.SubmitStart
				} else {
					submit = faulted.SubmitResult
				}
				if _, err := submit(input); !IsClass(err, ErrorIO) {
					t.Fatalf("faulted submit error = %v", err)
				}
				fresh := newSystemStore(systemFileSystem{}, systemClock{})
				if phase == "start" {
					_, err = fresh.SubmitStart(input)
				} else {
					_, err = fresh.SubmitResult(input)
				}
				if err != nil {
					t.Fatalf("exact retry failed: %v", err)
				}
				read, err := fresh.ReadByLocator(snapshot.Handoff.Locator)
				if err != nil {
					t.Fatal(err)
				}
				accepted := read.Start
				if phase == "terminal" {
					accepted = read.Terminal
				}
				if accepted == nil || accepted.DocumentID != id || accepted.SHA256 != document.SHA256 {
					t.Fatalf("accepted %s after retry = %#v", phase, accepted)
				}
			})
		}
	}
}

func TestStoreArtifactPublicationFaultsNeverPublishPartialTerminal(t *testing.T) {
	for _, stage := range []string{"temp", "write", "chmod", "file sync", "rename", "directory sync", "terminal object after artifact"} {
		t.Run(stage, func(t *testing.T) {
			root := physicalTempDir(t)
			base := newSystemStore(systemFileSystem{}, systemClock{})
			created, err := base.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
			if err != nil {
				t.Fatal(err)
			}
			start := testAcceptedDocument(t, created.Snapshot.Handoff.Identity.StartReceiptID, "started")
			if _, err := base.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: start}); err != nil {
				t.Fatal(err)
			}
			snapshot, err := base.ReadByLocator(created.Snapshot.Handoff.Locator)
			if err != nil {
				t.Fatal(err)
			}
			artifactBytes := []byte("complete artifact bytes\n")
			artifact := managedArtifact{ID: "verifier-stdout", SHA256: digestBytes(artifactBytes), Bytes: artifactBytes}
			document := testAcceptedDocument(t, snapshot.Handoff.Identity.TerminalResultID, "complete")
			input := SubmitStoreInput{Snapshot: snapshot, Phase: "terminal", Document: document, Artifacts: []managedArtifact{artifact}}
			faulted := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
			failed := false
			isArtifactPath := func(path string) bool { return filepath.Base(filepath.Dir(path)) == "sha256" }
			switch stage {
			case "temp", "terminal object after artifact":
				original := faulted.ops.createTemp
				faulted.ops.createTemp = func(directory, pattern string) (*os.File, error) {
					match := (stage == "temp" && filepath.Base(directory) == "sha256") || (stage == "terminal object after artifact" && filepath.Base(directory) == "objects" && filepath.Base(filepath.Dir(directory)) == "terminal")
					if match && !failed {
						failed = true
						return nil, errors.New("injected " + stage)
					}
					return original(directory, pattern)
				}
			case "write":
				original := faulted.ops.write
				faulted.ops.write = func(file *os.File, contents []byte) (int, error) {
					if isArtifactPath(file.Name()) && !failed {
						failed = true
						return 0, errors.New("injected artifact write")
					}
					return original(file, contents)
				}
			case "chmod":
				original := faulted.ops.chmodFile
				faulted.ops.chmodFile = func(file *os.File, mode os.FileMode) error {
					if isArtifactPath(file.Name()) && !failed {
						failed = true
						return errors.New("injected artifact chmod")
					}
					return original(file, mode)
				}
			case "file sync":
				original := faulted.ops.syncFile
				faulted.ops.syncFile = func(file *os.File) error {
					if isArtifactPath(file.Name()) && !failed {
						failed = true
						return errors.New("injected artifact sync")
					}
					return original(file)
				}
			case "rename":
				original := faulted.ops.rename
				faulted.ops.rename = func(source, destination string) error {
					if filepath.Base(filepath.Dir(destination)) == "sha256" && !failed {
						failed = true
						return errors.New("injected artifact rename")
					}
					return original(source, destination)
				}
			case "directory sync":
				original := faulted.ops.syncDirectory
				faulted.ops.syncDirectory = func(directory string) error {
					if filepath.Base(directory) == "sha256" && !failed {
						failed = true
						return errors.New("injected artifact directory sync")
					}
					return original(directory)
				}
			}
			if _, err := faulted.SubmitResult(input); !IsClass(err, ErrorIO) {
				t.Fatalf("faulted terminal error = %v", err)
			}
			if _, err := os.Lstat(filepath.Join(snapshot.Handoff.ReplyRoot, "terminal", "accepted.ref")); !os.IsNotExist(err) {
				t.Fatalf("terminal authority appeared after artifact fault: %v", err)
			}
			artifactPath := filepath.Join(snapshot.Handoff.ReplyRoot, "artifacts", "sha256", strings.TrimPrefix(artifact.SHA256, "sha256:"))
			if actual, err := os.ReadFile(artifactPath); err == nil && !bytes.Equal(actual, artifactBytes) {
				t.Fatalf("partial artifact = %q", actual)
			} else if err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			fresh := newSystemStore(systemFileSystem{}, systemClock{})
			if _, err := fresh.SubmitResult(input); err != nil {
				t.Fatalf("exact retry failed: %v", err)
			}
			actual, err := os.ReadFile(artifactPath)
			if err != nil || !bytes.Equal(actual, artifactBytes) {
				t.Fatalf("final artifact = %q, %v", actual, err)
			}
			read, err := fresh.ReadByLocator(snapshot.Handoff.Locator)
			if err != nil || read.Terminal == nil {
				t.Fatalf("terminal after retry = %#v, %v", read.Terminal, err)
			}
		})
	}
}

func testStoreHandoff(t *testing.T, root string) handoffDocument {
	t.Helper()
	identityValue := canonicaljson.Object{
		{Name: "activity_id", Value: "act_11111111111111111111111111111111"},
		{Name: "run_id", Value: "run_22222222222222222222222222222222"},
		{Name: "handoff_id", Value: "hnd_33333333333333333333333333333333"},
		{Name: "start_receipt_id", Value: "rcp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{Name: "terminal_result_id", Value: "res_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
		{Name: "publication_key", Value: "wf01/test"},
		{Name: "activity_key", Value: "wf01/test"},
		{Name: "created_at_utc", Value: "2026-09-27T00:00:00Z"},
	}
	replyRoot := filepath.Join(root, ".ply", "workflow", "handoffs", "v1", "runs", "run_22222222222222222222222222222222", "reply")
	value := canonicaljson.Object{
		{Name: "identity", Value: identityValue},
		{Name: "source_draft_sha256", Value: digestBytes([]byte("draft"))},
		{Name: "goal", Value: canonicaljson.Object{{Name: "title", Value: "Test handoff"}}},
		{Name: "workspace_binding", Value: canonicaljson.Object{{Name: "root", Value: root}, {Name: "marker_format_version", Value: int64(1)}, {Name: "marker_sha256", Value: digestBytes([]byte("marker"))}}},
		{Name: "project_binding", Value: canonicaljson.Object{{Name: "project_id", Value: "ply"}, {Name: "repo_id", Value: "ply"}, {Name: "registered_locator", Value: root}, {Name: "registered_git_common_dir", Value: filepath.Join(root, ".git")}}},
		{Name: "target_binding", Value: canonicaljson.Object{
			{Name: "worktree", Value: root},
			{Name: "git_common_dir", Value: filepath.Join(root, ".git")},
			{Name: "ref", Value: "refs/heads/main"},
			{Name: "oid", Value: "1111111111111111111111111111111111111111"},
			{Name: "tree", Value: "2222222222222222222222222222222222222222"},
			{Name: "status_policy", Value: canonicaljson.Object{{Name: "mode", Value: "clean"}}},
		}},
		{Name: "budget", Value: canonicaljson.Object{{Name: "max_rounds", Value: int64(5)}}},
		{Name: "reply_capability", Value: canonicaljson.Object{{Name: "capability_id", Value: "cap_44444444444444444444444444444444"}, {Name: "secret", Value: "5555555555555555555555555555555555555555555555555555555555555555"}, {Name: "reply_root", Value: replyRoot}, {Name: "allowed_operations", Value: []canonicaljson.Value{"submit-result", "submit-start"}}}},
	}
	bytes, err := canonicaljson.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return handoffDocument{Value: value, Bytes: bytes, SHA256: digestBytes(bytes), Identity: identity{ActivityID: "act_11111111111111111111111111111111", RunID: "run_22222222222222222222222222222222", HandoffID: "hnd_33333333333333333333333333333333", StartReceiptID: "rcp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", TerminalResultID: "res_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", PublicationKey: "wf01/test", ActivityKey: "wf01/test", CreatedAtUTC: "2026-09-27T00:00:00Z"}, SourceDraftSHA256: digestBytes([]byte("draft")), GoalTitle: "Test handoff", ReplyCapabilityID: "cap_44444444444444444444444444444444", ReplySecret: "5555555555555555555555555555555555555555555555555555555555555555", ReplyRoot: replyRoot, MaxRounds: 5}
}

func testAcceptedDocument(t *testing.T, id, outcome string) acceptedDocument {
	t.Helper()
	value := canonicaljson.Object{{Name: "id", Value: id}, {Name: "outcome", Value: outcome}}
	bytes, err := canonicaljson.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return acceptedDocument{Value: value, Bytes: bytes, SHA256: digestBytes(bytes), DocumentID: id, Outcome: outcome}
}

func physicalTempDir(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return root
}
