package backlog

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

type archiveFixture struct {
	store   Store
	archive Archive
	root    string
}

func newArchiveFixture(t *testing.T) archiveFixture {
	t.Helper()
	store := testStore(t)
	store.random = rand.Reader
	root := filepath.Dir(store.dir)
	return archiveFixture{store: store, archive: OpenArchive(filepath.Join(root, "data", "archive", "backlog"), root), root: root}
}

// done adds a done item with an authored and an anonymous note and tags; a
// non-empty source makes it a file item.
func (f archiveFixture) done(t *testing.T, source string) Item {
	t.Helper()
	item, err := f.store.Add("Tittel", "Tekst <æ>\n")
	if source != "" {
		item, err = f.store.AddFile("Spec", source)
	}
	if err == nil {
		_, err = f.store.AddNoteAs(item.ID, "med", "ply-admin")
	}
	if err == nil {
		_, err = f.store.AddNote(item.ID, "uten")
	}
	if err == nil {
		_, err = f.store.SetTags(item.ID, []string{"ply", "gone-repo"}, []string{"Backlog", "Web app"})
	}
	if err == nil {
		item, err = f.store.SetStatus(item.ID, StatusDone)
	}
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func (f archiveFixture) tree(t *testing.T) map[string]string {
	return backlogTree(t, Store{dir: f.root})
}

// restore writes the given backlogTree entries back, as an interrupted archiving leaves them.
func restore(t *testing.T, dir string, tree map[string]string, only ...string) {
	t.Helper()
	for relative, contents := range tree {
		if len(only) > 0 && !strings.HasPrefix(relative, only[0]) {
			continue
		}
		path := filepath.Join(dir, relative)
		var err error
		if contents == "dir" {
			err = os.MkdirAll(path, 0o755)
		} else {
			err = errors.Join(os.MkdirAll(filepath.Dir(path), 0o755), os.WriteFile(path, []byte(contents), 0o644))
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (f archiveFixture) archiveOK(t *testing.T, id, group string, newGroup bool) ArchivedItem {
	t.Helper()
	archived, err := f.store.Archive(id, f.archive, group, newGroup)
	if err != nil {
		t.Fatalf("archive %s to %s: %v", id, group, err)
	}
	if _, err := os.Lstat(f.store.path(id)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("active item remains: %v", err)
	}
	return archived
}

func (f archiveFixture) groupIDs(t *testing.T, group string) []string {
	t.Helper()
	read, err := f.archive.ReadGroup(group)
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, item := range read.Items {
		ids = append(ids, item.ID)
	}
	return ids
}

func (f archiveFixture) requireUnchanged(t *testing.T, before map[string]string) {
	t.Helper()
	if after := f.tree(t); !reflect.DeepEqual(after, before) {
		t.Fatalf("files changed:\nbefore %v\nafter  %v", before, after)
	}
}

func groupJSON(group string, ids ...string) string {
	items := []string{}
	for _, id := range ids {
		items = append(items, `{"id":"`+id+`","title":"t","text":"x","created_at":"c","notes":[],"archived_at":"a"}`)
	}
	return `{"format":"ply-backlog-archive-group/v1","group":"` + group + `","items":[` + strings.Join(items, ",") + "]}\n"
}

func writeGroupFile(t *testing.T, archive Archive, name, contents string) {
	t.Helper()
	if err := os.MkdirAll(archive.dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(archive.dir, name), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestArchiveKeepsAllFieldsInOrderAndRepeatsWithoutWriting(t *testing.T) {
	f := newArchiveFixture(t)
	first := f.done(t, "")
	archived := f.archiveOK(t, first.ID, "backlog-kjerne", true)
	want := ArchivedItem{ID: first.ID, Title: "Tittel", Text: "Tekst <æ>\n", CreatedAt: first.CreatedAt,
		Notes: []ArchivedNote{{At: first.Notes[0].At, Text: "med", Author: "ply-admin"}, {At: first.Notes[1].At, Text: "uten"}},
		Repos: []string{"ply", "gone-repo"}, Topics: []string{"Backlog", "Web app"}, ArchivedAt: "2026-11-04T21:30:00Z"}
	if !reflect.DeepEqual(archived, want) {
		t.Fatalf("archived = %#v", archived)
	}
	path := filepath.Join(f.archive.dir, "backlog-kjerne.json")
	raw, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(raw), "\n  \"group\": \"backlog-kjerne\",\n") || !strings.Contains(string(raw), `<æ>`) ||
		strings.Count(string(raw), `"author"`) != 1 || strings.Contains(string(raw), `"status"`) {
		t.Fatalf("group file = %s, %v", raw, err)
	}
	second := f.done(t, "")
	f.archiveOK(t, second.ID, "backlog-kjerne", false)
	if ids := f.groupIDs(t, "backlog-kjerne"); !reflect.DeepEqual(ids, []string{first.ID, second.ID}) {
		t.Fatalf("ids = %v", ids)
	}
	contents, _ := os.ReadFile(path)
	modified := ageFile(t, path)
	for _, newGroup := range []bool{false, true} {
		if again := f.archiveOK(t, first.ID, "backlog-kjerne", newGroup); !reflect.DeepEqual(again, want) {
			t.Fatalf("again = %#v", again)
		}
	}
	requireFileUnchanged(t, path, string(contents), modified)
	if _, err := f.store.Read(first.ID); err == nil {
		t.Fatal("archived item is still readable as active")
	}
}

func TestArchiveRejectionsChangeNoFiles(t *testing.T) {
	f := newArchiveFixture(t)
	other := f.done(t, "")
	active := backlogTree(t, f.store)
	f.archiveOK(t, other.ID, "andre", true)
	restore(t, f.store.dir, active, other.ID)
	open, _ := f.store.Add("Åpen", "x")
	ready, _ := f.store.Add("Klar", "x")
	if _, err := f.store.SetStatus(ready.ID, StatusReady); err != nil {
		t.Fatal(err)
	}
	target := f.done(t, "")
	before := f.tree(t)
	for _, test := range []struct {
		id, group string
		newGroup  bool
		class     string
	}{
		{open.ID, "andre", false, ClassItemNotDone},
		{ready.ID, "andre", false, ClassItemNotDone},
		{target.ID, "Ugyldig", true, ClassInvalidArguments},
		{target.ID, "-x", true, ClassInvalidArguments},
		{target.ID, strings.Repeat("a", 64), true, ClassInvalidArguments},
		{"20261104-0000000", "andre", false, ClassInvalidArguments},
		{target.ID, "ukjent", false, ClassGroupNotFound},
		{target.ID, "andre", true, ClassArchiveConflict},
		{other.ID, "ny", true, ClassArchiveConflict},
		{"20261104-00000000", "andre", false, ClassItemNotFound},
	} {
		_, err := f.store.Archive(test.id, f.archive, test.group, test.newGroup)
		requireClass(t, err, test.class)
		f.requireUnchanged(t, before)
	}
	_, err := f.store.Archive(target.ID, f.archive, "ukjent", false)
	if !strings.Contains(err.Error(), "ply backlog groups") {
		t.Fatalf("detail = %v", err)
	}
}

func TestArchiveResumesAfterInterruptionAndRejectsChangedSource(t *testing.T) {
	for _, newGroup := range []bool{false, true} {
		f := newArchiveFixture(t)
		item := f.done(t, "")
		active := backlogTree(t, f.store)
		f.archiveOK(t, item.ID, "g", true)
		restore(t, f.store.dir, active)
		f.archiveOK(t, item.ID, "g", newGroup)
		if ids := f.groupIDs(t, "g"); !reflect.DeepEqual(ids, []string{item.ID}) {
			t.Fatalf("ids = %v", ids)
		}
	}
	f := newArchiveFixture(t)
	item := f.done(t, "")
	active := backlogTree(t, f.store)
	f.archiveOK(t, item.ID, "g", true)
	restore(t, f.store.dir, active)
	if _, err := f.store.AddNote(item.ID, "ny"); err != nil {
		t.Fatal(err)
	}
	before := f.tree(t)
	_, err := f.store.Archive(item.ID, f.archive, "g", false)
	requireClass(t, err, ClassArchiveConflict)
	f.requireUnchanged(t, before)
}

func TestArchiveFailuresAreResumable(t *testing.T) {
	boom := errors.New("boom")
	t.Run("replace fails before rename", func(t *testing.T) {
		f := newArchiveFixture(t)
		item := f.done(t, "")
		f.archive.replace = func(string, []byte, fs.FileMode) error { return &atomicfile.Error{Op: "replace", Err: boom} }
		_, err := f.store.Archive(item.ID, f.archive, "g", true)
		requireClass(t, err, ClassIOError)
		if _, err := f.store.Read(item.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Lstat(filepath.Join(f.archive.dir, "g.json")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("group file published: %v", err)
		}
	})
	t.Run("replace publishes then fails", func(t *testing.T) {
		f := newArchiveFixture(t)
		item := f.done(t, "")
		f.archive.replace = func(path string, contents []byte, mode fs.FileMode) error {
			return errors.Join(atomicfile.Replace(path, contents, mode), &atomicfile.Error{Op: "sync directory", Replaced: true, Err: boom})
		}
		_, err := f.store.Archive(item.ID, f.archive, "g", true)
		requireClass(t, err, ClassIOError)
		if _, err := f.store.Read(item.ID); err != nil {
			t.Fatal(err)
		}
		f.archive.replace = atomicfile.Replace
		f.archiveOK(t, item.ID, "g", true)
		if ids := f.groupIDs(t, "g"); len(ids) != 1 {
			t.Fatalf("ids = %v", ids)
		}
	})
	for _, failing := range []string{"F", "H", "H-file", "after-remove", "R"} {
		failing := failing
		t.Run("sync fails in "+failing, func(t *testing.T) {
			f := newArchiveFixture(t)
			item := f.done(t, sourceFile(t, "a.md", []byte("innhold")))
			target := map[string]string{
				"F": filepath.Join(f.archive.dir, "files", item.ID), "H": filepath.Join(f.archive.dir, "g.json"),
				"H-file": filepath.Join(f.archive.dir, "files", item.ID, "a.md"), "after-remove": f.store.dir,
				"R": filepath.Join(f.store.dir, "files"),
			}[failing]
			f.archive.sync = func(path string) error {
				if path == target {
					return boom
				}
				return syncDirectory(path)
			}
			_, err := f.store.Archive(item.ID, f.archive, "g", true)
			requireClass(t, err, ClassIOError)
			f.archive.sync = syncDirectory
			f.archiveOK(t, item.ID, "g", failing == "F") // F fails before the group exists
			if ids := f.groupIDs(t, "g"); len(ids) != 1 {
				t.Fatalf("ids = %v", ids)
			}
			if copied, err := os.ReadFile(filepath.Join(f.archive.dir, "files", item.ID, "a.md")); err != nil || string(copied) != "innhold" {
				t.Fatalf("archive copy = %q, %v", copied, err)
			}
			if _, err := os.Lstat(filepath.Join(f.store.dir, "files", item.ID)); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("active copy remains: %v", err)
			}
		})
	}
}

func TestParallelNotesAndArchivingLoseNothing(t *testing.T) {
	f := newArchiveFixture(t)
	f.archiveOK(t, f.done(t, "").ID, "g", true)
	var items []Item
	for i := 0; i < 4; i++ {
		items = append(items, f.done(t, ""))
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	noted := map[string][]string{}
	for _, item := range items {
		item := item
		wg.Go(func() {
			if _, err := f.store.Archive(item.ID, f.archive, "g", false); err != nil {
				t.Error(err)
			}
		})
		for n := 0; n < 5; n++ {
			n := n
			wg.Go(func() {
				text := fmt.Sprintf("note %d", n)
				if _, err := f.store.AddNote(item.ID, text); err == nil {
					mu.Lock()
					noted[item.ID] = append(noted[item.ID], text)
					mu.Unlock()
				} else if backlogError := (*Error)(nil); !errors.As(err, &backlogError) || backlogError.Class != ClassItemNotFound {
					t.Error(err)
				}
			})
		}
	}
	wg.Wait()
	group, err := f.archive.ReadGroup("g")
	if err != nil || len(group.Items) != 5 {
		t.Fatalf("group = %+v, %v", group, err)
	}
	for _, archived := range group.Items[1:] {
		kept := map[string]bool{}
		for _, note := range archived.Notes {
			kept[note.Text] = true
		}
		for _, text := range noted[archived.ID] {
			if !kept[text] {
				t.Fatalf("note %q of %s lost", text, archived.ID)
			}
		}
	}
}

func TestInvalidArchiveIsListedAndBlocksArchiving(t *testing.T) {
	const id = "20261101-0a0b0c0d"
	for _, test := range []struct {
		name    string
		files   map[string]string
		invalid []InvalidFile
	}{
		{"unreadable", map[string]string{"g.json": groupJSON("g"), "broken.json": "{"}, []InvalidFile{{File: "broken.json", Error: "invalid JSON at offset 1"}}},
		{"duplicate in group", map[string]string{"g.json": groupJSON("g", id, id)}, []InvalidFile{{File: "g.json", Error: "duplicate item id"}}},
		{"duplicate across groups", map[string]string{"g.json": groupJSON("g", id), "h.json": groupJSON("h", id), "i.json": groupJSON("i")},
			[]InvalidFile{{File: "g.json", Error: "item id archived in another group"}, {File: "h.json", Error: "item id archived in another group"}}},
		{"raw invalid UTF-8", map[string]string{"g.json": strings.Replace(groupJSON("g", id), `"text":"x"`, "\"text\":\"\xff\"", 1)}, []InvalidFile{{File: "g.json", Error: "group file is not valid UTF-8"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newArchiveFixture(t)
			target := f.done(t, "")
			for name, contents := range test.files {
				writeGroupFile(t, f.archive, name, contents)
			}
			result, err := f.archive.ListGroups()
			if err != nil || !reflect.DeepEqual(result.Invalid, test.invalid) {
				t.Fatalf("list = %+v, %v", result, err)
			}
			before := f.tree(t)
			_, err = f.store.Archive(target.ID, f.archive, "g", false)
			requireClass(t, err, ClassArchiveInvalid)
			f.requireUnchanged(t, before)
		})
	}
}

func TestListGroupsAndReadGroupCreateNothing(t *testing.T) {
	f := newArchiveFixture(t)
	result, err := f.archive.ListGroups()
	if err != nil || len(result.Groups) != 0 || result.Invalid == nil {
		t.Fatalf("empty list = %#v, %v", result, err)
	}
	_, err = f.archive.ReadGroup("g")
	requireClass(t, err, ClassGroupNotFound)
	_, err = f.archive.ReadGroup("G")
	requireClass(t, err, ClassInvalidArguments)
	if _, err := os.Lstat(f.archive.dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("archive created: %v", err)
	}
	writeGroupFile(t, f.archive, "b.json", groupJSON("b", "20261101-0a0b0c0d", "20261101-0a0b0c0e"))
	writeGroupFile(t, f.archive, "a.json", groupJSON("a", "20261101-0a0b0c0f"))
	writeGroupFile(t, f.archive, "other.json", groupJSON("c"))
	if err := os.Symlink(filepath.Join(f.archive.dir, "a.json"), filepath.Join(f.archive.dir, "link.json")); err != nil {
		t.Fatal(err)
	}
	result, err = f.archive.ListGroups()
	want := GroupListResult{Groups: []GroupSummary{{"a", 1, "a"}, {"b", 2, "a"}},
		Invalid: []InvalidFile{{File: "link.json", Error: "group file is not a regular file"}, {File: "other.json", Error: "group does not match file name"}}}
	if err != nil || !reflect.DeepEqual(result, want) {
		t.Fatalf("list = %#v, %v", result, err)
	}
}

func TestArchiveFileItemKeepsItsOwnCopy(t *testing.T) {
	f := newArchiveFixture(t)
	contents := bytes.Repeat([]byte("ab\xff"), copyBufferBytes)
	source := sourceFile(t, "my spec.md", contents)
	item := f.done(t, source)
	archived := f.archiveOK(t, item.ID, "filer", true)
	if archived.Text != "" || *archived.File != (ArchivedFile{Name: "my spec.md", Size: int64(len(contents))}) {
		t.Fatalf("archived = %#v", archived)
	}
	if _, err := os.Lstat(filepath.Join(f.store.dir, "files", item.ID)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("active copy remains: %v", err)
	}
	if err := errors.Join(os.RemoveAll(f.store.dir), os.Remove(source)); err != nil {
		t.Fatal(err)
	}
	path, state, err := f.archive.FileState(archived)
	if err != nil || state != FilePresent || path != filepath.Join(f.archive.dir, "files", item.ID, "my spec.md") {
		t.Fatalf("state = %s %s %v", path, state, err)
	}
	if copied, err := os.ReadFile(path); err != nil || !bytes.Equal(copied, contents) {
		t.Fatalf("copy differs: %v", err)
	}
	group, err := f.archive.ReadGroup("filer")
	if err != nil || !reflect.DeepEqual(group.Items, []ArchivedItem{archived}) {
		t.Fatalf("group = %#v, %v", group, err)
	}
	// Without a backlog directory the repeat completes without creating one.
	f.archiveOK(t, item.ID, "filer", false)
	if _, err := os.Lstat(f.store.dir); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("backlog created: %v", err)
	}
}

func TestArchiveFileCopyConflicts(t *testing.T) {
	setup := func(t *testing.T) (archiveFixture, Item, string, string) {
		f := newArchiveFixture(t)
		item := f.done(t, sourceFile(t, "a.md", []byte("innhold")))
		return f, item, filepath.Join(f.store.dir, "files", item.ID, "a.md"), filepath.Join(f.archive.dir, "files", item.ID, "a.md")
	}
	for name, mutate := range map[string]func(t *testing.T, active string) error{
		"missing":   func(t *testing.T, active string) error { return os.Remove(active) },
		"truncated": func(t *testing.T, active string) error { return os.WriteFile(active, []byte("inn"), 0o644) },
		"symlink": func(t *testing.T, active string) error {
			return errors.Join(os.Remove(active), os.Symlink(sourceFile(t, "b", []byte("innhold")), active))
		},
	} {
		t.Run("active copy "+name, func(t *testing.T) {
			f, item, active, _ := setup(t)
			if err := mutate(t, active); err != nil {
				t.Fatal(err)
			}
			before := f.tree(t)
			_, err := f.store.Archive(item.ID, f.archive, "g", true)
			requireClass(t, err, ClassItemInvalid)
			f.requireUnchanged(t, before)
		})
	}
	for name, test := range map[string]struct {
		place func(path string) error
		class string
	}{
		"same bytes":  {func(path string) error { return os.WriteFile(path, []byte("innhold"), 0o644) }, ""},
		"other bytes": {func(path string) error { return os.WriteFile(path, []byte("annet!!"), 0o644) }, ClassArchiveConflict},
		"directory":   {func(path string) error { return os.Mkdir(path, 0o755) }, ClassArchiveConflict},
	} {
		t.Run("leftover archive copy "+name, func(t *testing.T) {
			f, item, _, archived := setup(t)
			if err := errors.Join(os.MkdirAll(filepath.Dir(archived), 0o755), test.place(archived)); err != nil {
				t.Fatal(err)
			}
			if test.class == "" {
				f.archiveOK(t, item.ID, "g", true)
				return
			}
			before := f.tree(t)
			_, err := f.store.Archive(item.ID, f.archive, "g", true)
			requireClass(t, err, test.class)
			f.requireUnchanged(t, before)
		})
	}
	for name, test := range map[string]struct {
		change func(active, archived string) error
		class  string
	}{
		"active copy differs":   {func(active, _ string) error { return os.WriteFile(active, []byte("annet!!"), 0o644) }, ClassArchiveConflict},
		"archive copy missing":  {func(_, archived string) error { return os.Remove(archived) }, ClassArchiveInvalid},
		"archive copy replaced": {func(_, archived string) error { return errors.Join(os.Remove(archived), os.Mkdir(archived, 0o755)) }, ClassArchiveInvalid},
	} {
		t.Run("entry exists and "+name, func(t *testing.T) {
			f, item, active, archived := setup(t)
			snapshot := backlogTree(t, f.store)
			f.archiveOK(t, item.ID, "g", true)
			restore(t, f.store.dir, snapshot)
			if err := test.change(active, archived); err != nil {
				t.Fatal(err)
			}
			before := f.tree(t)
			_, err := f.store.Archive(item.ID, f.archive, "g", false)
			requireClass(t, err, test.class)
			f.requireUnchanged(t, before)
		})
	}
	for name, contents := range map[string]string{"same": "innhold", "other": "annet!!"} {
		t.Run("active copy left after removal with "+name+" bytes", func(t *testing.T) {
			f, item, active, _ := setup(t)
			snapshot := backlogTree(t, f.store)
			f.archiveOK(t, item.ID, "g", true)
			restore(t, f.store.dir, snapshot, "files")
			if err := os.WriteFile(active, []byte(contents), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := f.store.Archive(item.ID, f.archive, "g", false)
			if _, statErr := os.Lstat(active); contents == "innhold" && (err != nil || !errors.Is(statErr, fs.ErrNotExist)) {
				t.Fatalf("cleanup: %v, %v", err, statErr)
			} else if contents != "innhold" && (statErr != nil || err == nil) {
				t.Fatalf("other bytes: %v, %v", err, statErr)
			}
			if contents != "innhold" {
				requireClass(t, err, ClassArchiveConflict)
			}
		})
	}
}

func TestArchiveRejectsPresentEmptyAuthor(t *testing.T) {
	f := newArchiveFixture(t)
	const item = `{"id":"20261101-0a0b0c0d","title":"t","text":"x","created_at":"c","notes":[{"at":"a","text":"x"%s}],"archived_at":"a"}`
	writeGroupFile(t, f.archive, "absent.json", `{"format":"ply-backlog-archive-group/v1","group":"absent","items":[`+fmt.Sprintf(item, "")+`]}`)
	writeGroupFile(t, f.archive, "empty.json", `{"format":"ply-backlog-archive-group/v1","group":"empty","items":[`+fmt.Sprintf(strings.Replace(item, "0a0b0c0d", "0a0b0c0e", 1), `,"author":""`)+`]}`)
	result, err := f.archive.ListGroups()
	want := GroupListResult{Groups: []GroupSummary{{"absent", 1, "a"}}, Invalid: []InvalidFile{{File: "empty.json", Error: "invalid note author"}}}
	if err != nil || !reflect.DeepEqual(result, want) {
		t.Fatalf("list = %#v, %v", result, err)
	}
	_, err = f.archive.ReadGroup("empty")
	requireClass(t, err, ClassArchiveInvalid)
}

func TestArchiveRejectsBacklogDirectoryAsArchive(t *testing.T) {
	f := newArchiveFixture(t)
	item := f.done(t, "")
	alias := filepath.Join(f.root, "alias")
	if err := os.Symlink(f.store.dir, alias); err != nil {
		t.Fatal(err)
	}
	directories := []string{f.store.dir, alias}
	if upper := filepath.Join(f.root, "BACKLOG"); func() bool { _, err := os.Stat(upper); return err == nil }() {
		directories = append(directories, upper) // case-insensitive file system
	}
	before := f.tree(t)
	for _, directory := range directories {
		_, err := f.store.Archive(item.ID, OpenArchive(directory, f.root), "g", true)
		requireClass(t, err, ClassInvalidArguments)
		if !strings.Contains(err.Error(), "must differ") {
			t.Fatalf("%s: %v", directory, err)
		}
		f.requireUnchanged(t, before)
	}
}

func TestArchiveRefusesItemsWithAttachmentsWithoutChangingFiles(t *testing.T) {
	f := newArchiveFixture(t)
	f.archiveOK(t, f.done(t, "").ID, "g", true)
	text := f.done(t, "")
	file := f.done(t, sourceFile(t, "spec.md", []byte("spec")))
	interrupted := f.done(t, "")
	broken := f.done(t, "")
	active := backlogTree(t, f.store)
	f.archiveOK(t, interrupted.ID, "g", false)
	restore(t, f.store.dir, active, interrupted.ID)
	for _, id := range []string{text.ID, file.ID, interrupted.ID} {
		attach(t, f.store, id, "trace.png", "png")
	}
	if err := os.Remove(f.store.path(interrupted.ID)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.store.attachmentsPath(broken.ID), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := f.tree(t)
	for _, test := range []struct {
		id, group string
		newGroup  bool
	}{
		{text.ID, "ny", true},
		{text.ID, "g", false},
		{text.ID, "ukjent", false},
		{file.ID, "g", false},
		{interrupted.ID, "g", false},
		{broken.ID, "g", false},
	} {
		_, err := f.store.Archive(test.id, f.archive, test.group, test.newGroup)
		requireClass(t, err, ClassItemInvalid)
		if !strings.Contains(err.Error(), "item has attachments") {
			t.Fatalf("detail = %v", err)
		}
		f.requireUnchanged(t, before)
	}
}
