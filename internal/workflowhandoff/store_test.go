package workflowhandoff

import (
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
