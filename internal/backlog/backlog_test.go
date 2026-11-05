package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

var fixedTime = time.Date(2026, 11, 4, 21, 30, 0, 0, time.UTC)

const schedulerTitle = "Scheduler: kjør nattlig opprydding"

func testStore(t *testing.T, random ...byte) Store {
	t.Helper()
	store := Open(t.TempDir())
	store.now = func() time.Time { return fixedTime }
	if len(random) == 0 {
		random = []byte{0x47, 0x53, 0x0e, 0x79}
	}
	store.random = bytes.NewReader(random)
	return store
}

func requireClass(t *testing.T, err error, class string) {
	t.Helper()
	var backlogError *Error
	if !errors.As(err, &backlogError) || backlogError.Class != class {
		t.Fatalf("error = %v, want class %s", err, class)
	}
}

func backlogFiles(t *testing.T, store Store) map[string]string {
	t.Helper()
	files := map[string]string{}
	entries, err := os.ReadDir(filepath.Join(store.root, "backlog"))
	if errors.Is(err, os.ErrNotExist) {
		return files
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(store.root, "backlog", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files[entry.Name()] = string(contents)
	}
	return files
}

func TestListWithoutBacklogIsEmptyAndCreatesNothing(t *testing.T) {
	store := testStore(t)
	result, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(result)
	if err != nil || string(encoded) != `{"items":[],"invalid":[]}` {
		t.Fatalf("list = %s, %v", encoded, err)
	}
	if _, err := os.Stat(filepath.Join(store.root, "backlog")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("backlog directory stat = %v", err)
	}
}

func TestAddPreservesTextAndWritesIndentedItemFile(t *testing.T) {
	store := testStore(t)
	text := "  a\n\n## b\n---\n"
	item, err := store.Add("  "+schedulerTitle+"  ", text)
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "20261104-47530e79" || item.Title != schedulerTitle || item.Text != text {
		t.Fatalf("item = %+v", item)
	}
	contents, err := os.ReadFile(filepath.Join(store.root, "backlog", "20261104-47530e79.json"))
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n" +
		"  \"id\": \"20261104-47530e79\",\n" +
		"  \"title\": \"Scheduler: kjør nattlig opprydding\",\n" +
		"  \"text\": \"  a\\n\\n## b\\n---\\n\",\n" +
		"  \"created_at\": \"2026-11-04T21:30:00Z\",\n" +
		"  \"status\": \"open\",\n" +
		"  \"notes\": []\n" +
		"}\n"
	if string(contents) != want {
		t.Fatalf("file =\n%s\nwant\n%s", contents, want)
	}
	read, err := store.Read(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	readJSON, _ := json.Marshal(read)
	addJSON, _ := json.Marshal(item)
	if !bytes.Equal(readJSON, addJSON) {
		t.Fatalf("read = %s, add = %s", readJSON, addJSON)
	}
}

func TestAddDoesNotEscapeHTMLCharacters(t *testing.T) {
	store := testStore(t)
	item, err := store.Add("a <b> & c", "x <y> & z\n")
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(store.root, "backlog", item.ID+".json"))
	if err != nil || !strings.Contains(string(contents), `"text": "x <y> & z\n"`) {
		t.Fatalf("file = %s, %v", contents, err)
	}
}

func TestNotesAppendInOrderAndKeepExistingFields(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text\n")
	if err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedTime.Add(time.Hour) }
	if _, err := store.AddNote(item.ID, "first\n"); err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedTime.Add(2*time.Hour + 500*time.Millisecond) }
	noted, err := store.AddNote(item.ID, "  second")
	if err != nil {
		t.Fatal(err)
	}
	want := []Note{{At: "2026-11-04T22:30:00Z", Text: "first\n"}, {At: "2026-11-04T23:30:00Z", Text: "  second"}}
	if noted.Text != "text\n" || noted.CreatedAt != "2026-11-04T21:30:00Z" || fmt.Sprint(noted.Notes) != fmt.Sprint(want) {
		t.Fatalf("item = %+v", noted)
	}
	read, err := store.Read(item.ID)
	if err != nil || fmt.Sprint(read) != fmt.Sprint(noted) {
		t.Fatalf("read = %+v, %v", read, err)
	}
	result, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	wantSummary := Summary{ID: item.ID, Title: schedulerTitle, CreatedAt: "2026-11-04T21:30:00Z", UpdatedAt: "2026-11-04T23:30:00Z", NoteCount: 2, Status: StatusOpen}
	if len(result.Items) != 1 || result.Items[0] != wantSummary {
		t.Fatalf("list = %+v", result)
	}
}

func TestParallelNotesAreNotLost(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	const writers = 20
	var group sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			_, err := store.AddNote(item.ID, fmt.Sprintf("note %d", i))
			errs <- err
		}(i)
	}
	group.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	read, err := store.Read(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, note := range read.Notes {
		seen[note.Text] = true
	}
	if len(read.Notes) != writers || len(seen) != writers {
		t.Fatalf("notes = %d, distinct = %d", len(read.Notes), len(seen))
	}
}

func TestListSortsByCreatedAtDescendingThenID(t *testing.T) {
	store := testStore(t, 0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa, 0xcc, 0xcc, 0xcc, 0xcc)
	for i := 0; i < 2; i++ {
		if _, err := store.Add(fmt.Sprintf("same %d", i), "x"); err != nil {
			t.Fatal(err)
		}
	}
	store.now = func() time.Time { return fixedTime.Add(time.Second) }
	if _, err := store.Add("newest", "x"); err != nil {
		t.Fatal(err)
	}
	result, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, item := range result.Items {
		ids = append(ids, item.ID)
	}
	if strings.Join(ids, ",") != "20261104-cccccccc,20261104-aaaaaaaa,20261104-bbbbbbbb" {
		t.Fatalf("ids = %v", ids)
	}
}

func TestInvalidFileIsIsolated(t *testing.T) {
	store := testStore(t)
	valid, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(store.root, "backlog")
	writes := map[string]string{
		"20261103-0badf00d.json": "{\"id\": ",
		"20261103-00000001.json": `{"id":"20261103-00000002","title":"t","text":"x"}`,
		"20261103-00000003.json": `{"id":"20261103-00000003","title":"t"}`,
		"notes.txt":              "ignored",
		".atomicfile-1.tmp":      "ignored",
	}
	for name, contents := range writes {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	result, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].ID != valid.ID {
		t.Fatalf("items = %+v", result.Items)
	}
	want := []InvalidFile{
		{File: "20261103-00000001.json", Error: "id does not match file name"},
		{File: "20261103-00000003.json", Error: "missing text"},
		{File: "20261103-0badf00d.json", Error: "invalid JSON at offset 7"},
	}
	if fmt.Sprint(result.Invalid) != fmt.Sprint(want) {
		t.Fatalf("invalid = %+v", result.Invalid)
	}
	for _, id := range []string{"20261103-0badf00d", "20261103-00000001", "20261103-00000003"} {
		_, err := store.Read(id)
		requireClass(t, err, ClassItemInvalid)
		before := backlogFiles(t, store)
		_, err = store.AddNote(id, "note")
		requireClass(t, err, ClassItemInvalid)
		if fmt.Sprint(backlogFiles(t, store)) != fmt.Sprint(before) {
			t.Fatalf("note on invalid %s changed files", id)
		}
	}
}

func TestRejectedInputsChangeNoFiles(t *testing.T) {
	store := testStore(t, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	before := backlogFiles(t, store)
	tooLarge := strings.Repeat("a", 64*1024+1)
	cases := []struct {
		name  string
		class string
		run   func() error
	}{
		{"empty text", ClassInvalidArguments, func() error { _, err := store.Add("t", ""); return err }},
		{"blank text", ClassInvalidArguments, func() error { _, err := store.Add("t", " \n\t "); return err }},
		{"large text", ClassInvalidArguments, func() error { _, err := store.Add("t", tooLarge); return err }},
		{"invalid utf8", ClassInvalidArguments, func() error { _, err := store.Add("t", "a\xff"); return err }},
		{"multiline title", ClassInvalidArguments, func() error { _, err := store.Add("a\nb", "x"); return err }},
		{"carriage return title", ClassInvalidArguments, func() error { _, err := store.Add("a\rb", "x"); return err }},
		{"blank title", ClassInvalidArguments, func() error { _, err := store.Add("  ", "x"); return err }},
		{"long title", ClassInvalidArguments, func() error { _, err := store.Add(strings.Repeat("æ", 201), "x"); return err }},
		{"blank note", ClassInvalidArguments, func() error { _, err := store.AddNote(item.ID, "  "); return err }},
		{"large note", ClassInvalidArguments, func() error { _, err := store.AddNote(item.ID, tooLarge); return err }},
		{"invalid utf8 note", ClassInvalidArguments, func() error { _, err := store.AddNote(item.ID, "\xff"); return err }},
		{"bad id show", ClassInvalidArguments, func() error { _, err := store.Read("../x"); return err }},
		{"bad id note", ClassInvalidArguments, func() error { _, err := store.AddNote("20261104-ABCDEF12", "x"); return err }},
		{"unknown show", ClassItemNotFound, func() error { _, err := store.Read("20261104-00000000"); return err }},
		{"unknown note", ClassItemNotFound, func() error { _, err := store.AddNote("20261104-00000000", "x"); return err }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.run()
			requireClass(t, err, c.class)
			if strings.Contains(err.Error(), store.root) || strings.Contains(err.Error(), "aaaa") {
				t.Fatalf("error leaks path or text: %v", err)
			}
			if fmt.Sprint(backlogFiles(t, store)) != fmt.Sprint(before) {
				t.Fatal("files changed")
			}
		})
	}
	if _, err := store.Add(strings.Repeat("æ", 200), strings.Repeat("a", 64*1024)); err != nil {
		t.Fatalf("limits are inclusive: %v", err)
	}
}

func TestNoteWithoutBacklogIsNotFoundAndCreatesNothing(t *testing.T) {
	store := testStore(t)
	_, err := store.AddNote("20261104-00000000", "x")
	requireClass(t, err, ClassItemNotFound)
	if _, err := os.Stat(filepath.Join(store.root, "backlog")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("backlog directory stat = %v", err)
	}
}

func TestProjectRootMustBeAbsoluteExistingDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, root := range []string{"relative", filepath.Join(t.TempDir(), "missing"), file} {
		store := Open(root)
		_, err := store.List()
		requireClass(t, err, ClassInvalidArguments)
		_, err = store.Add("t", "x")
		requireClass(t, err, ClassInvalidArguments)
		if strings.Contains(err.Error(), root) {
			t.Fatalf("error leaks root: %v", err)
		}
	}
}

func TestAddRetriesExistingIDAtMostFiveTimes(t *testing.T) {
	store := testStore(t, 0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22)
	if err := os.MkdirAll(filepath.Join(store.root, "backlog"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.root, "backlog", "20261104-11111111.json"), []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	item, err := store.Add("t", "x")
	if err != nil || item.ID != "20261104-22222222" {
		t.Fatalf("item = %+v, %v", item, err)
	}
	if files := backlogFiles(t, store); files["20261104-11111111.json"] != "keep" {
		t.Fatalf("existing file changed: %q", files["20261104-11111111.json"])
	}

	exhausted := testStore(t, bytes.Repeat([]byte{0x33}, 4*7)...)
	if _, err := exhausted.Add("t", "x"); err != nil {
		t.Fatal(err)
	}
	_, err = exhausted.Add("t", "x")
	requireClass(t, err, ClassIOError)
	if remaining := exhausted.random.(*bytes.Reader).Len(); remaining != 4 {
		t.Fatalf("attempts used = %d, want 5", 6-remaining/4)
	}
}

func TestPackageDependsOnlyOnAllowedPackages(t *testing.T) {
	output, err := exec.Command("go", "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", ".").Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	for _, dependency := range strings.Fields(string(output)) {
		if dependency != "github.com/devdimensionlab/plybuild/internal/backlog" &&
			dependency != "github.com/devdimensionlab/plybuild/internal/atomicfile" &&
			!strings.HasPrefix(dependency, "golang.org/x/sys/") {
			t.Fatalf("unexpected dependency %s", dependency)
		}
	}
}

const statusItemID = "20261104-47530e79"

func writeItemFile(t *testing.T, store Store, id, contents string) string {
	t.Helper()
	path := filepath.Join(store.root, "backlog", id+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readOnlyBacklog(t *testing.T, store Store) {
	t.Helper()
	directory := filepath.Join(store.root, "backlog")
	if err := os.Chmod(directory, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(directory, 0o755) })
}

func requireFileUnchanged(t *testing.T, path, contents string, modified time.Time) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != contents || !info.ModTime().Equal(modified) {
		t.Fatalf("file changed: %q, mtime %v want %v, %v", after, info.ModTime(), modified, err)
	}
}

func ageFile(t *testing.T, path string) time.Time {
	t.Helper()
	old := fixedTime.Add(-24 * time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
	return old
}

func statusItemFile(status string) string {
	return "{\n  \"id\": \"" + statusItemID + "\",\n  \"title\": \"t\",\n  \"text\": \"x\",\n  \"created_at\": \"2026-11-04T21:30:00Z\",\n" + status + "  \"notes\": []\n}\n"
}

func TestMissingNullOrUnknownStatusIsInvalidWithoutRewrite(t *testing.T) {
	cases := map[string]struct{ status, detail string }{
		"missing":      {"", "missing status"},
		"null":         {"  \"status\": null,\n", "missing status"},
		"empty":        {"  \"status\": \"\",\n", "missing status"},
		"legacy ready": {"  \"ready\": true,\n", "missing status"},
		"unknown":      {"  \"status\": \"closed\",\n", "invalid status"},
		"wrong case":   {"  \"status\": \"Ready\",\n", "invalid status"},
		"not a string": {"  \"status\": 1,\n", "invalid item JSON"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			contents := statusItemFile(c.status)
			path := writeItemFile(t, store, statusItemID, contents)
			modified := ageFile(t, path)
			_, err := store.Read(statusItemID)
			requireClass(t, err, ClassItemInvalid)
			if err.Error() != ClassItemInvalid+": "+c.detail {
				t.Fatalf("read = %v", err)
			}
			result, err := store.List()
			want := []InvalidFile{{File: statusItemID + ".json", Error: c.detail}}
			if err != nil || len(result.Items) != 0 || fmt.Sprint(result.Invalid) != fmt.Sprint(want) {
				t.Fatalf("list = %+v, %v", result, err)
			}
			_, err = store.AddNote(statusItemID, "note")
			requireClass(t, err, ClassItemInvalid)
			_, err = store.SetStatus(statusItemID, StatusOpen)
			requireClass(t, err, ClassItemInvalid)
			requireFileUnchanged(t, path, contents, modified)
		})
	}
}

func TestSetStatusAllowsEveryTransitionAndChangesOnlyStatus(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text\n")
	if err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedTime.Add(time.Hour) }
	if _, err := store.AddNote(item.ID, "first\n"); err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedTime.Add(5 * time.Hour) }
	path := filepath.Join(store.root, "backlog", item.ID+".json")
	encoded := func(status Status) string {
		return "{\n" +
			"  \"id\": \"20261104-47530e79\",\n" +
			"  \"title\": \"Scheduler: kjør nattlig opprydding\",\n" +
			"  \"text\": \"text\\n\",\n" +
			"  \"created_at\": \"2026-11-04T21:30:00Z\",\n" +
			"  \"status\": \"" + string(status) + "\",\n" +
			"  \"notes\": [\n" +
			"    {\n" +
			"      \"at\": \"2026-11-04T22:30:00Z\",\n" +
			"      \"text\": \"first\\n\"\n" +
			"    }\n" +
			"  ]\n" +
			"}\n"
	}
	// open→ready→done→open→done→ready→open covers every transition, including reopening.
	for _, status := range []Status{StatusReady, StatusDone, StatusOpen, StatusDone, StatusReady, StatusOpen} {
		set, err := store.SetStatus(item.ID, status)
		if err != nil || set.Status != status || set.Text != "text\n" || len(set.Notes) != 1 {
			t.Fatalf("set status %s = %+v, %v", status, set, err)
		}
		contents, err := os.ReadFile(path)
		if err != nil || string(contents) != encoded(status) {
			t.Fatalf("file =\n%s\nwant\n%s", contents, encoded(status))
		}
		result, err := store.List()
		want := Summary{ID: item.ID, Title: schedulerTitle, CreatedAt: "2026-11-04T21:30:00Z", UpdatedAt: "2026-11-04T22:30:00Z", NoteCount: 1, Status: status}
		if err != nil || len(result.Items) != 1 || result.Items[0] != want {
			t.Fatalf("list = %+v, %v", result, err)
		}
	}
}

func TestSetStatusToSameValueWritesNothing(t *testing.T) {
	for _, status := range []Status{StatusOpen, StatusReady, StatusDone} {
		t.Run(string(status), func(t *testing.T) {
			store := testStore(t)
			contents := statusItemFile("  \"status\": \"" + string(status) + "\",\n")
			path := writeItemFile(t, store, statusItemID, contents)
			if err := os.WriteFile(filepath.Join(store.root, "backlog", ".lock"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
			modified := ageFile(t, path)
			// A read-only directory makes any temp file creation fail.
			readOnlyBacklog(t, store)
			item, err := store.SetStatus(statusItemID, status)
			if err != nil || item.ID != statusItemID || item.Status != status {
				t.Fatalf("set status = %+v, %v", item, err)
			}
			requireFileUnchanged(t, path, contents, modified)
			for file := range backlogFiles(t, store) {
				if strings.HasPrefix(file, ".atomicfile-") {
					t.Fatalf("temp file left: %s", file)
				}
			}
		})
	}
}

func TestListShowsEveryStatus(t *testing.T) {
	store := testStore(t, 0x01, 0x01, 0x01, 0x01, 0x02, 0x02, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03)
	statuses := map[string]Status{}
	for _, status := range []Status{StatusOpen, StatusReady, StatusDone} {
		item, err := store.Add(string(status), "x")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.SetStatus(item.ID, status); err != nil {
			t.Fatal(err)
		}
		statuses[item.ID] = status
	}
	result, err := store.List()
	if err != nil || len(result.Items) != 3 {
		t.Fatalf("list = %+v, %v", result, err)
	}
	for _, summary := range result.Items {
		if summary.Status != statuses[summary.ID] {
			t.Fatalf("summary = %+v, want %s", summary, statuses[summary.ID])
		}
	}
}

func TestStatusIsKeptByNotesAndAddStartsOpen(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil || item.Status != StatusOpen {
		t.Fatalf("add = %+v, %v", item, err)
	}
	if _, err := store.SetStatus(item.ID, StatusDone); err != nil {
		t.Fatal(err)
	}
	noted, err := store.AddNote(item.ID, "note")
	if err != nil || noted.Status != StatusDone || len(noted.Notes) != 1 {
		t.Fatalf("note = %+v, %v", noted, err)
	}
	contents, err := os.ReadFile(filepath.Join(store.root, "backlog", item.ID+".json"))
	if err != nil || !strings.Contains(string(contents), "\"status\": \"done\",\n") {
		t.Fatalf("file = %s, %v", contents, err)
	}
}

func TestParallelNotesAndStatusAreNotLost(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	const writers = 20
	var group sync.WaitGroup
	errs := make(chan error, 2*writers)
	for i := 0; i < writers; i++ {
		group.Add(2)
		go func(i int) {
			defer group.Done()
			_, err := store.AddNote(item.ID, fmt.Sprintf("note %d", i))
			errs <- err
		}(i)
		go func() {
			defer group.Done()
			_, err := store.SetStatus(item.ID, StatusDone)
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
	read, err := store.Read(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, note := range read.Notes {
		seen[note.Text] = true
	}
	if len(read.Notes) != writers || len(seen) != writers || read.Status != StatusDone {
		t.Fatalf("notes = %d, distinct = %d, status = %s", len(read.Notes), len(seen), read.Status)
	}
}

func TestSetStatusWaitsForBacklogLock(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	held, err := os.OpenFile(filepath.Join(store.root, "backlog", ".lock"), os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	if err := lockFile(held); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := store.SetStatus(item.ID, StatusReady)
		done <- err
	}()
	select {
	case err := <-done:
		t.Fatalf("set status returned while lock was held: %v", err)
	case <-time.After(200 * time.Millisecond):
	}
	if err := unlockFile(held); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if read, err := store.Read(item.ID); err != nil || read.Status != StatusReady {
		t.Fatalf("read = %+v, %v", read, err)
	}
}

func TestSetStatusErrors(t *testing.T) {
	store := testStore(t)
	_, err := store.SetStatus("../x", StatusReady)
	requireClass(t, err, ClassInvalidArguments)
	_, err = store.SetStatus(statusItemID, StatusReady)
	requireClass(t, err, ClassItemNotFound)
	if _, err := os.Stat(filepath.Join(store.root, "backlog")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("backlog directory stat = %v", err)
	}

	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(store.root, "backlog", item.ID+".json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []Status{"", "closed", "Ready"} {
		_, err := store.SetStatus(item.ID, status)
		requireClass(t, err, ClassInvalidArguments)
		if err.Error() != ClassInvalidArguments+": status must be open, ready or done" {
			t.Fatalf("set status %q = %v", status, err)
		}
	}
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(after, before) {
		t.Fatalf("file changed after rejected status: %s, %v", after, err)
	}
	_, err = store.SetStatus("20261104-00000000", StatusReady)
	requireClass(t, err, ClassItemNotFound)
	invalidPath := writeItemFile(t, store, "20261103-0badf00d", "{\"id\": ")
	modified := ageFile(t, invalidPath)
	_, err = store.SetStatus("20261103-0badf00d", StatusReady)
	requireClass(t, err, ClassItemInvalid)
	requireFileUnchanged(t, invalidPath, "{\"id\": ", modified)

	readOnlyBacklog(t, store)
	_, err = store.SetStatus(item.ID, StatusReady)
	requireClass(t, err, ClassIOError)
	if strings.Contains(err.Error(), store.root) {
		t.Fatalf("error leaks root: %v", err)
	}
	if after, err := os.ReadFile(path); err != nil || !bytes.Equal(after, before) {
		t.Fatalf("file changed after failed write: %s, %v", after, err)
	}
}
