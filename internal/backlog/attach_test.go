package backlog

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

const attachedAt = "2026-11-04T21:30:00Z"

func textItem(t *testing.T) (Store, Item) {
	t.Helper()
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	return store, item
}

func attach(t *testing.T, store Store, id, name, contents string) []Attachment {
	t.Helper()
	list, err := store.Attach(id, sourceFile(t, name, []byte(contents)))
	if err != nil {
		t.Fatal(err)
	}
	return list
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func sidecarPath(store Store, id string) string {
	return filepath.Join(store.dir, "attachments", id+".json")
}

func attachmentPath(store Store, id, name string) string {
	return filepath.Join(store.dir, "files", id, name)
}

func requireTree(t *testing.T, store Store, want map[string]string) {
	t.Helper()
	if got := backlogTree(t, store); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("tree =\n%v\nwant\n%v", got, want)
	}
}

func TestAttachCopiesFilesToTextItem(t *testing.T) {
	store, item := textItem(t)
	itemFile := mustRead(t, store.path(item.ID))
	first := sourceFile(t, "a&b.txt", []byte("first"))
	second := sourceFile(t, "trace.png", []byte{0xff, 0x00})
	if _, err := store.Attach(item.ID, first); err != nil {
		t.Fatal(err)
	}
	list, err := store.Attach(item.ID, second)
	want := []Attachment{{Name: "a&b.txt", Size: 5, AddedAt: attachedAt}, {Name: "trace.png", Size: 2, AddedAt: attachedAt}}
	if err != nil || !reflect.DeepEqual(list, want) {
		t.Fatalf("attach = %+v, %v", list, err)
	}
	if read, err := store.Attachments(item.ID); err != nil || !reflect.DeepEqual(read, want) {
		t.Fatalf("attachments = %+v, %v", read, err)
	}
	entry := func(name string, size int) string {
		return fmt.Sprintf("    {\n      \"name\": %q,\n      \"size\": %d,\n      \"added_at\": %q\n    }", name, size, attachedAt)
	}
	wantSidecar := "{\n  \"format\": \"ply-backlog-attachments/v1\",\n  \"id\": \"" + item.ID + "\",\n  \"attachments\": [\n" +
		entry("a&b.txt", 5) + ",\n" + entry("trace.png", 2) + "\n  ]\n}\n"
	if sidecar := mustRead(t, sidecarPath(store, item.ID)); sidecar != wantSidecar {
		t.Fatalf("sidecar =\n%s\nwant\n%s", sidecar, wantSidecar)
	}
	if mustRead(t, attachmentPath(store, item.ID, "a&b.txt")) != "first" || mustRead(t, attachmentPath(store, item.ID, "trace.png")) != "\xff\x00" {
		t.Fatal("copies differ")
	}
	if mustRead(t, first) != "first" || mustRead(t, second) != "\xff\x00" || mustRead(t, store.path(item.ID)) != itemFile {
		t.Fatal("source or item file changed")
	}
}

func TestAttachToLegacyFileItem(t *testing.T) {
	store := testStore(t)
	legacy, err := store.AddFile("t", sourceFile(t, "spec.md", []byte("legacy")))
	if err != nil {
		t.Fatal(err)
	}
	itemFile := mustRead(t, store.path(legacy.ID))
	list := attach(t, store, legacy.ID, "trace.png", "new")
	if mustRead(t, copyPath(store, legacy)) != "legacy" || mustRead(t, store.path(legacy.ID)) != itemFile {
		t.Fatal("legacy copy or item file changed")
	}
	if path, state, err := store.AttachmentState(legacy.ID, list[0]); err != nil || state != FilePresent || path != attachmentPath(store, legacy.ID, "trace.png") {
		t.Fatalf("attachment state = %q %q %v", path, state, err)
	}
	if _, state, err := store.FileState(legacy); err != nil || state != FilePresent {
		t.Fatalf("file state = %q %v", state, err)
	}
}

func TestOrdinaryUpdatesKeepAttachments(t *testing.T) {
	store, item := textItem(t)
	want := attach(t, store, item.ID, "a", "x")
	sidecar := mustRead(t, sidecarPath(store, item.ID))
	if _, err := store.AddNoteAs(item.ID, "note", "me"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetStatus(item.ID, StatusDone); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetTags(item.ID, []string{"ply"}, []string{"topic"}); err != nil {
		t.Fatal(err)
	}
	if mustRead(t, sidecarPath(store, item.ID)) != sidecar {
		t.Fatal("sidecar changed")
	}
	if list, err := store.Attachments(item.ID); err != nil || !reflect.DeepEqual(list, want) {
		t.Fatalf("attachments = %+v, %v", list, err)
	}
	if read, err := store.Read(item.ID); err != nil || read.Status != StatusDone {
		t.Fatalf("read = %+v, %v", read, err)
	}
}

func TestAttachRejectsDuplicateNames(t *testing.T) {
	store := testStore(t)
	legacy, err := store.AddFile("t", sourceFile(t, "spec.md", []byte("legacy")))
	if err != nil {
		t.Fatal(err)
	}
	attach(t, store, legacy.ID, "a.txt", "first")
	if err := os.WriteFile(attachmentPath(store, legacy.ID, "left.txt"), []byte("leftover"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := backlogTree(t, store)
	for _, name := range []string{"a.txt", "A.TXT", "spec.md", "left.txt"} {
		t.Run(name, func(t *testing.T) {
			_, err := store.Attach(legacy.ID, sourceFile(t, name, []byte("other")))
			if err == nil || err.Error() != ClassAttachmentExists+": attachment already exists" {
				t.Fatalf("attach = %v", err)
			}
			requireTree(t, store, before)
		})
	}
}

func TestAttachRejectsSourcesWithoutWriting(t *testing.T) {
	store, item := textItem(t)
	sources := t.TempDir()
	regular := sourceFile(t, "regular", []byte("x"))
	link := filepath.Join(sources, "link")
	if err := os.Symlink(regular, link); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ name, source, detail string }{
		{"relative", "regular", "source must be an absolute path"},
		{"missing", filepath.Join(sources, "missing"), "open file: no such file or directory"},
		{"directory", sources, "file must be a regular file"},
		{"symlink", link, "file must be a regular file"},
	}
	before := backlogTree(t, store)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := store.Attach(item.ID, c.source)
			if err == nil || err.Error() != ClassInvalidArguments+": "+c.detail {
				t.Fatalf("attach = %v, want %s", err, c.detail)
			}
			requireTree(t, store, before)
		})
	}
}

func TestAttachRejectsSymlinkedDestinations(t *testing.T) {
	cases := []struct{ name, class, link string }{
		{"file directory", ClassIOError, filepath.Join("files", statusItemID)},
		{"attachments directory", ClassItemInvalid, "attachments"},
		{"sidecar", ClassItemInvalid, filepath.Join("attachments", statusItemID+".json")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store, item := textItem(t)
			backlog := store.dir
			link := filepath.Join(backlog, c.link)
			for _, directory := range []string{filepath.Join(backlog, "files"), filepath.Dir(link)} {
				if err := os.MkdirAll(directory, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			outside := t.TempDir()
			target := outside
			if c.name == "sidecar" {
				target = sourceFile(t, "outside.json", []byte("{}"))
			}
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}
			before := backlogTree(t, store)
			_, err := store.Attach(item.ID, sourceFile(t, "a", []byte("x")))
			requireClass(t, err, c.class)
			requireTree(t, store, before)
			if entries, _ := os.ReadDir(outside); len(entries) != 0 || c.name == "sidecar" && mustRead(t, target) != "{}" {
				t.Fatalf("wrote outside the backlog: %v", entries)
			}
		})
	}
}

func TestAttachToMissingOrInvalidItem(t *testing.T) {
	store := testStore(t)
	source := sourceFile(t, "a", []byte("x"))
	_, err := store.Attach(statusItemID, source)
	requireClass(t, err, ClassItemNotFound)
	_, err = store.Attachments(statusItemID)
	requireClass(t, err, ClassItemNotFound)
	_, err = store.Attach("../x", source)
	requireClass(t, err, ClassInvalidArguments)
	_, err = store.Attachments("../x")
	requireClass(t, err, ClassInvalidArguments)
	writeItemFile(t, store, statusItemID, "{")
	_, err = store.Attach(statusItemID, source)
	requireClass(t, err, ClassItemInvalid)
	if _, err := os.Lstat(filepath.Join(store.dir, "files")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("files directory = %v", err)
	}
}

func TestInvalidSidecarIsRejectedWithoutWriting(t *testing.T) {
	sidecar := func(format, id string, entries ...string) string {
		list := ""
		for i, entry := range entries {
			if i > 0 {
				list += ","
			}
			list += `{"name":"` + entry + `","size":1,"added_at":"` + attachedAt + `"}`
		}
		return `{"format":"` + format + `","id":"` + id + `","attachments":[` + list + `]}`
	}
	cases := map[string]string{
		"bad JSON":        `{`,
		"unknown format":  sidecar("ply-backlog-attachments/v2", statusItemID, "a"),
		"id mismatch":     sidecar(AttachmentsFormat, "20261104-00000000", "a"),
		"duplicate names": sidecar(AttachmentsFormat, statusItemID, "a", "A"),
		"invalid name":    sidecar(AttachmentsFormat, statusItemID, "a/b"),
		"negative size":   `{"format":"` + AttachmentsFormat + `","id":"` + statusItemID + `","attachments":[{"name":"a","size":-1,"added_at":"x"}]}`,
		"empty added_at":  `{"format":"` + AttachmentsFormat + `","id":"` + statusItemID + `","attachments":[{"name":"a","size":1,"added_at":""}]}`,
	}
	for name, contents := range cases {
		t.Run(name, func(t *testing.T) {
			store, item := textItem(t)
			attach(t, store, item.ID, "a", "x")
			if err := os.WriteFile(sidecarPath(store, item.ID), []byte(contents), 0o644); err != nil {
				t.Fatal(err)
			}
			before := backlogTree(t, store)
			_, err := store.Attachments(item.ID)
			requireClass(t, err, ClassItemInvalid)
			_, err = store.Attach(item.ID, sourceFile(t, "b", []byte("y")))
			requireClass(t, err, ClassItemInvalid)
			requireTree(t, store, before)
		})
	}
}

func TestAttachPublicationFailure(t *testing.T) {
	failing := func(string, []byte, fs.FileMode) error { return errors.New("boom") }
	t.Run("replace fails without writing", func(t *testing.T) {
		store, item := textItem(t)
		attach(t, store, item.ID, "a", "x")
		before := backlogTree(t, store)
		store.replace = failing
		_, err := store.Attach(item.ID, sourceFile(t, "b", []byte("y")))
		requireClass(t, err, ClassIOError)
		requireTree(t, store, before)
	})
	t.Run("replace fails before any sidecar", func(t *testing.T) {
		store, item := textItem(t)
		store.replace = failing
		_, err := store.Attach(item.ID, sourceFile(t, "b", []byte("y")))
		requireClass(t, err, ClassIOError)
		if _, err := os.Lstat(attachmentPath(store, item.ID, "b")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("copy = %v", err)
		}
		if list, err := store.Attachments(item.ID); err != nil || len(list) != 0 {
			t.Fatalf("attachments = %+v, %v", list, err)
		}
	})
	t.Run("replace writes then fails", func(t *testing.T) {
		store, item := textItem(t)
		attach(t, store, item.ID, "a", "x")
		store.replace = func(path string, contents []byte, mode fs.FileMode) error {
			if err := atomicfile.Replace(path, contents, mode); err != nil {
				return err
			}
			return errors.New("sync failed")
		}
		_, err := store.Attach(item.ID, sourceFile(t, "b", []byte("y")))
		requireClass(t, err, ClassIOError)
		if list, err := store.Attachments(item.ID); err != nil || len(list) != 2 || list[1].Name != "b" {
			t.Fatalf("attachments = %+v, %v", list, err)
		}
		if mustRead(t, attachmentPath(store, item.ID, "b")) != "y" {
			t.Fatal("published copy removed")
		}
	})
}

func TestParallelAttachmentsAreNotLost(t *testing.T) {
	store, item := textItem(t)
	var group sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		source := sourceFile(t, fmt.Sprintf("f%d", i), []byte{byte(i)})
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := store.Attach(item.ID, source)
			errs <- err
		}()
	}
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	list, err := store.Attachments(item.ID)
	names := map[string]bool{}
	for _, attachment := range list {
		names[attachment.Name] = true
	}
	if err != nil || len(list) != 8 || len(names) != 8 {
		t.Fatalf("attachments = %+v, %v", list, err)
	}
	if entries, err := os.ReadDir(filepath.Join(store.dir, "files")); err != nil || len(entries) != 1 || entries[0].Name() != item.ID {
		t.Fatalf("files = %v, %v", entries, err)
	}
}

func TestAttachLimit(t *testing.T) {
	store, item := textItem(t)
	for i := 0; i < 32; i++ {
		attach(t, store, item.ID, fmt.Sprintf("f%02d", i), "x")
	}
	before := backlogTree(t, store)
	_, err := store.Attach(item.ID, sourceFile(t, "f32", []byte("x")))
	requireClass(t, err, ClassInvalidArguments)
	requireTree(t, store, before)
}

func TestAttachmentReadbackAndState(t *testing.T) {
	store, item := textItem(t)
	if list, err := store.Attachments(item.ID); err != nil || list == nil || len(list) != 0 {
		t.Fatalf("attachments = %#v, %v", list, err)
	}
	attachment := attach(t, store, item.ID, "a", "x")[0]
	copied := attachmentPath(store, item.ID, "a")
	expect := func(want FileState) {
		t.Helper()
		if path, state, err := store.AttachmentState(item.ID, attachment); err != nil || state != want || path != copied {
			t.Fatalf("state = %q %q %v, want %s", path, state, err, want)
		}
	}
	expect(FilePresent)
	if err := os.Remove(copied); err != nil {
		t.Fatal(err)
	}
	expect(FileMissing)
	if err := os.Symlink(sourceFile(t, "a", []byte("y")), copied); err != nil {
		t.Fatal(err)
	}
	expect(FileUnexpected)
	_, _, err := store.AttachmentState(item.ID, Attachment{Name: ".."})
	requireClass(t, err, ClassInvalidArguments)
	_, _, err = store.AttachmentState("../x", attachment)
	requireClass(t, err, ClassInvalidArguments)
}
