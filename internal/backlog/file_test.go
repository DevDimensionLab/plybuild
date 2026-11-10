package backlog

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// backlogTree snapshots every entry below backlog/ without following symlinks.
func backlogTree(t *testing.T, store Store) map[string]string {
	t.Helper()
	tree := map[string]string{}
	root := store.dir
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, _ := filepath.Rel(root, path)
		switch {
		case entry.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(path)
			tree[relative] = "link:" + target
			return err
		case entry.IsDir():
			tree[relative] = "dir"
		case entry.Type().IsRegular():
			contents, err := os.ReadFile(path)
			tree[relative] = string(contents)
			return err
		default:
			tree[relative] = "other"
		}
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		t.Fatal(err)
	}
	return tree
}

func sourceFile(t *testing.T, name string, contents []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func copyPath(store Store, item Item) string {
	return filepath.Join(store.dir, "files", item.ID, item.File.Name)
}

func TestItemWithoutFileKeepsBytesThroughNoteAndStatus(t *testing.T) {
	store := testStore(t)
	contents := statusItemFile("  \"status\": \"ready\",\n")
	path := writeItemFile(t, store, statusItemID, contents)
	if _, err := store.AddNote(statusItemID, "n"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetStatus(statusItemID, StatusReady); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetStatus(statusItemID, StatusDone); err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"id\": \"" + statusItemID + "\",\n  \"title\": \"t\",\n  \"text\": \"x\",\n  \"created_at\": \"2026-11-04T21:30:00Z\",\n" +
		"  \"status\": \"done\",\n  \"notes\": [\n    {\n      \"at\": \"2026-11-04T21:30:00Z\",\n      \"text\": \"n\"\n    }\n  ]\n}\n"
	if after, err := os.ReadFile(path); err != nil || string(after) != want {
		t.Fatalf("file =\n%s\nwant\n%s", after, want)
	}
}

func TestReadValidatesFileMetadata(t *testing.T) {
	item := func(text, file string) string {
		return "{\"id\":\"" + statusItemID + "\",\"title\":\"t\",\"text\":\"" + text + "\",\"created_at\":\"2026-11-04T21:30:00Z\",\"status\":\"open\",\"notes\":[]" + file + "}"
	}
	cases := map[string]struct{ contents, detail string }{
		"empty text without file": {item("", ""), "missing text"},
		"empty text null file":    {item("", `,"file":null`), "missing text"},
		"empty name":              {item("", `,"file":{"name":"","size":1}`), "invalid file name"},
		"dot dot name":            {item("", `,"file":{"name":"..","size":1}`), "invalid file name"},
		"slash name":              {item("", `,"file":{"name":"a/b","size":1}`), "invalid file name"},
		"backslash name":          {item("", `,"file":{"name":"a\\\\b","size":1}`), "invalid file name"},
		"newline name":            {item("", `,"file":{"name":"a\nb","size":1}`), "invalid file name"},
		"negative size":           {item("", `,"file":{"name":"a","size":-1}`), "invalid file size"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			writeItemFile(t, store, statusItemID, c.contents)
			_, err := store.Read(statusItemID)
			if err == nil || err.Error() != ClassItemInvalid+": "+c.detail {
				t.Fatalf("read = %v, want %s", err, c.detail)
			}
		})
	}
	store := testStore(t)
	writeItemFile(t, store, statusItemID, item("", `,"file":{"name":"min spec.md","size":0}`))
	read, err := store.Read(statusItemID)
	if err != nil || read.File == nil || *read.File != (File{Name: "min spec.md", Size: 0}) {
		t.Fatalf("read = %+v, %v", read, err)
	}
}

func TestAddFileCopiesEmptyNonUTF8AndMultiBufferFiles(t *testing.T) {
	large := make([]byte, 3*32*1024+17)
	for i := range large {
		large[i] = byte(i * 7)
	}
	cases := map[string][]byte{"empty": {}, "binary.bin": {0xff, 0xfe, 0x00, 0x80}, "large.dat": large}
	for name, contents := range cases {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			source := sourceFile(t, name, contents)
			item, err := store.AddFile("  Spec  ", source)
			if err != nil {
				t.Fatal(err)
			}
			if item.Title != "Spec" || item.Text != "" || item.Status != StatusOpen || item.File == nil || *item.File != (File{Name: name, Size: int64(len(contents))}) {
				t.Fatalf("item = %+v", item)
			}
			if original, err := os.ReadFile(source); err != nil || !bytes.Equal(original, contents) {
				t.Fatalf("original changed: %v", err)
			}
			if err := os.Remove(source); err != nil {
				t.Fatal(err)
			}
			copied, err := os.ReadFile(copyPath(store, item))
			if err != nil || !bytes.Equal(copied, contents) {
				t.Fatalf("copy differs: %v", err)
			}
			if info, err := os.Lstat(copyPath(store, item)); err != nil || info.Mode().Perm() != 0o644 {
				t.Fatalf("copy mode = %v, %v", info, err)
			}
			if _, err := store.AddNote(item.ID, "note"); err != nil {
				t.Fatal(err)
			}
			done, err := store.SetStatus(item.ID, StatusDone)
			if err != nil || done.File == nil || *done.File != *item.File || done.Text != "" {
				t.Fatalf("status = %+v, %v", done, err)
			}
			if read, err := store.Read(item.ID); err != nil || *read.File != *item.File {
				t.Fatalf("read = %+v, %v", read, err)
			}
			path, state, err := store.FileState(item)
			if err != nil || path != copyPath(store, item) || state != FilePresent {
				t.Fatalf("state = %q %q %v", path, state, err)
			}
		})
	}
}

func TestAddFileRejectsWithoutWriting(t *testing.T) {
	store := testStore(t)
	if _, err := store.Add(schedulerTitle, "text"); err != nil {
		t.Fatal(err)
	}
	sources := t.TempDir()
	regular := filepath.Join(sources, "regular")
	if err := os.WriteFile(regular, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(sources, "link")
	if err := os.Symlink(regular, link); err != nil {
		t.Fatal(err)
	}
	newline := filepath.Join(sources, "a\nb")
	if err := os.WriteFile(newline, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ name, title, source, detail string }{
		{"relative", "t", "regular", "source must be an absolute path"},
		{"blank title", " ", regular, "title must be one non-empty line of at most 200 characters"},
		{"symlink", "t", link, "file must be a regular file"},
		{"directory", "t", sources, "file must be a regular file"},
		{"missing", "t", filepath.Join(sources, "missing"), "open file: no such file or directory"},
		{"newline name", "t", newline, "invalid file name"},
		{"non utf8 name", "t", filepath.Join(sources, "a\xff"), "invalid file name"},
		{"root", "t", string(filepath.Separator), "invalid file name"},
	}
	before := backlogTree(t, store)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := store.AddFile(c.title, c.source)
			if err == nil || err.Error() != ClassInvalidArguments+": "+c.detail {
				t.Fatalf("add file = %v, want %s", err, c.detail)
			}
			if fmt.Sprint(backlogTree(t, store)) != fmt.Sprint(before) {
				t.Fatal("backlog changed")
			}
		})
	}
}

func TestAddFileRejectsSymlinkedFilesDirectory(t *testing.T) {
	store := testStore(t)
	if err := os.MkdirAll(store.dir, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(store.dir, "files")); err != nil {
		t.Fatal(err)
	}
	before := backlogTree(t, store)
	_, err := store.AddFile("t", sourceFile(t, "a", []byte("x")))
	if err == nil || err.Error() != ClassIOError+": files directory is not a directory" {
		t.Fatalf("add file = %v", err)
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 || fmt.Sprint(backlogTree(t, store)) != fmt.Sprint(before) {
		t.Fatalf("wrote through symlink: %v", entries)
	}
}

func TestAddFileSkipsTakenAndOrphanedIDs(t *testing.T) {
	store := testStore(t, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22, 0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44)
	writeItemFile(t, store, "20261104-11111111", "keep")
	orphan := filepath.Join(store.dir, "files", "20261104-22222222")
	if err := os.MkdirAll(orphan, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(orphan, "a"), []byte("orphan"), 0o644); err != nil {
		t.Fatal(err)
	}
	source := sourceFile(t, "a", []byte("same"))
	first, err := store.AddFile("t", source)
	if err != nil || first.ID != "20261104-33333333" {
		t.Fatalf("first = %+v, %v", first, err)
	}
	second, err := store.AddFile("t", source)
	if err != nil || second.ID != "20261104-44444444" {
		t.Fatalf("second = %+v, %v", second, err)
	}
	tree := backlogTree(t, store)
	if tree["20261104-11111111.json"] != "keep" || tree[filepath.Join("files", "20261104-22222222", "a")] != "orphan" ||
		tree[filepath.Join("files", first.ID, "a")] != "same" || tree[filepath.Join("files", second.ID, "a")] != "same" {
		t.Fatalf("tree = %v", tree)
	}
	if _, err := os.Lstat(filepath.Join(store.dir, "20261104-22222222.json")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("orphan was published: %v", err)
	}
}

func TestAddFilePublishFailure(t *testing.T) {
	t.Run("replace fails without writing", func(t *testing.T) {
		store := testStore(t)
		existing, err := store.Add(schedulerTitle, "text")
		if err != nil {
			t.Fatal(err)
		}
		store.random = bytes.NewReader([]byte{0x55, 0x55, 0x55, 0x55})
		store.replace = func(string, []byte, fs.FileMode) error { return errors.New("boom") }
		_, err = store.AddFile("t", sourceFile(t, "a", []byte("x")))
		requireClass(t, err, ClassIOError)
		tree := backlogTree(t, store)
		want := map[string]string{".": "dir", ".lock": "", "files": "dir", existing.ID + ".json": tree[existing.ID+".json"]}
		if fmt.Sprint(tree) != fmt.Sprint(want) || !strings.Contains(tree[existing.ID+".json"], `"text": "text"`) {
			t.Fatalf("tree = %v", tree)
		}
	})
	t.Run("replace writes then fails", func(t *testing.T) {
		store := testStore(t)
		store.replace = func(path string, contents []byte, mode fs.FileMode) error {
			if err := atomicfile.Replace(path, contents, mode); err != nil {
				return err
			}
			return errors.New("sync failed")
		}
		_, err := store.AddFile("t", sourceFile(t, "a", []byte("x")))
		requireClass(t, err, ClassIOError)
		result, err := store.List()
		if err != nil || len(result.Items) != 1 || len(result.Invalid) != 0 {
			t.Fatalf("list = %+v, %v", result, err)
		}
		item, err := store.Read(result.Items[0].ID)
		if err != nil {
			t.Fatal(err)
		}
		if copied, err := os.ReadFile(copyPath(store, item)); err != nil || string(copied) != "x" {
			t.Fatalf("copy = %q, %v", copied, err)
		}
	})
	t.Run("orphans are invisible", func(t *testing.T) {
		store := testStore(t)
		existing, err := store.Add(schedulerTitle, "text")
		if err != nil {
			t.Fatal(err)
		}
		files := filepath.Join(store.dir, "files")
		if err := os.MkdirAll(filepath.Join(files, "20261104-66666666"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(files, ".import-x"), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		result, err := store.List()
		if err != nil || len(result.Items) != 1 || result.Items[0].ID != existing.ID || len(result.Invalid) != 0 {
			t.Fatalf("list = %+v, %v", result, err)
		}
	})
}

func TestFileState(t *testing.T) {
	store := testStore(t)
	item, err := store.AddFile("t", sourceFile(t, "a", []byte("x")))
	if err != nil {
		t.Fatal(err)
	}
	if path, state, err := store.FileState(Item{ID: item.ID}); path != "" || state != "" || err != nil {
		t.Fatalf("no file = %q %q %v", path, state, err)
	}
	for _, bad := range []Item{{ID: "../x", File: item.File}, {ID: item.ID, File: &File{Name: ".."}}, {ID: item.ID, File: &File{Name: "a/b"}}} {
		_, _, err := store.FileState(bad)
		requireClass(t, err, ClassInvalidArguments)
	}
	if path, _, err := Open("relative/backlog").FileState(item); path != "" {
		t.Fatalf("relative root path = %q", path)
	} else {
		requireClass(t, err, ClassInvalidArguments)
	}
	files := filepath.Join(store.dir, "files")
	expect := func(want FileState) {
		t.Helper()
		if path, state, err := store.FileState(item); err != nil || state != want || path != copyPath(store, item) {
			t.Fatalf("state = %q %q %v, want %s", path, state, err, want)
		}
	}
	expect(FilePresent)
	if _, _, err := store.FileState(Item{ID: "20261104-00000000", File: item.File}); err != nil {
		t.Fatal(err)
	}

	copied := copyPath(store, item)
	if err := os.Remove(copied); err != nil {
		t.Fatal(err)
	}
	expect(FileMissing)
	if err := os.Mkdir(copied, 0o755); err != nil {
		t.Fatal(err)
	}
	expect(FileUnexpected)
	if err := os.Remove(copied); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(sourceFile(t, "a", []byte("y")), copied); err != nil {
		t.Fatal(err)
	}
	expect(FileUnexpected)

	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, item.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, item.ID, "a"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(files, item.ID)); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, item.ID), filepath.Join(files, item.ID)); err != nil {
		t.Fatal(err)
	}
	expect(FileUnexpected)
	if err := os.Remove(filepath.Join(files, item.ID)); err != nil {
		t.Fatal(err)
	}
	expect(FileMissing)
	if err := os.Remove(files); err != nil {
		t.Fatal(err)
	}
	expect(FileMissing)
	if err := os.Symlink(outside, files); err != nil {
		t.Fatal(err)
	}
	expect(FileUnexpected)
}
