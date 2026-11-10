package backlog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// threeNoteItem is a legacy-file item with item tags, an attachment and three notes at distinct times.
func threeNoteItem(t *testing.T) (Store, Item) {
	t.Helper()
	store := testStore(t)
	legacy, err := store.AddFile("t", sourceFile(t, "spec.md", []byte("legacy")))
	if err != nil {
		t.Fatal(err)
	}
	for i, author := range []string{"", "Per", "agent"} {
		store.now = func() time.Time { return fixedTime.Add(time.Duration(i) * time.Minute) }
		if _, err := store.AddNoteAs(legacy.ID, fmt.Sprintf(" note %d\n", i), author); err != nil {
			t.Fatal(err)
		}
	}
	store.now = func() time.Time { return fixedTime }
	if _, err := store.SetTags(legacy.ID, []string{"ply"}, []string{"Backlog"}); err != nil {
		t.Fatal(err)
	}
	attach(t, store, legacy.ID, "trace.png", "png")
	item, err := store.Read(legacy.ID)
	if err != nil {
		t.Fatal(err)
	}
	return store, item
}

func TestAddNoteTaggedNormalizesAndReadsBack(t *testing.T) {
	store, item := textItem(t)
	noted, err := store.AddNoteTagged(item.ID, " verbatim\n", "  Per  ", []string{" ply ", "ply", "web"}, []string{"Web app", "web APP"})
	want := Note{At: attachedAt, Text: " verbatim\n", Author: "Per", Repos: []string{"ply", "web"}, Topics: []string{"Web app"}}
	if err != nil || len(noted.Notes) != 1 || !reflect.DeepEqual(noted.Notes[0], want) {
		t.Fatalf("add = %+v, %v", noted, err)
	}
	if read, err := store.Read(item.ID); err != nil || !reflect.DeepEqual(read, noted) {
		t.Fatalf("read = %+v, %v", read, err)
	}
}

func TestAddNoteTaggedWithoutTagsWritesSameBytesAsAddNoteAs(t *testing.T) {
	plain, item := textItem(t)
	tagged, _ := textItem(t)
	if _, err := plain.AddNoteAs(item.ID, "n", "Per"); err != nil {
		t.Fatal(err)
	}
	if _, err := tagged.AddNoteTagged(item.ID, "n", "Per", nil, []string{}); err != nil {
		t.Fatal(err)
	}
	want := mustRead(t, plain.path(item.ID))
	if got := mustRead(t, tagged.path(item.ID)); got != want || strings.Contains(got, "repos") || strings.Contains(got, "topics") {
		t.Fatalf("tagged =\n%s\nwant\n%s", got, want)
	}
}

func TestSetNoteTagsChangesOnlyThatNote(t *testing.T) {
	store, before := threeNoteItem(t)
	tree := backlogTree(t, store)
	item, err := store.SetNoteTags(before.ID, 1, []string{"web", " ply"}, []string{" PWA "})
	want := before
	want.Notes = append([]Note(nil), before.Notes...)
	want.Notes[1].Repos, want.Notes[1].Topics = []string{"web", "ply"}, []string{"PWA"}
	if err != nil || !reflect.DeepEqual(item, want) {
		t.Fatalf("set = %+v, %v", item, err)
	}
	if read, err := store.Read(before.ID); err != nil || !reflect.DeepEqual(read, want) {
		t.Fatalf("read = %+v, %v", read, err)
	}
	contents, err := encode(want)
	if err != nil {
		t.Fatal(err)
	}
	tree[before.ID+".json"] = string(contents)
	requireTree(t, store, tree)
}

func TestSetNoteTagsReplacesAndClears(t *testing.T) {
	store, item := threeNoteItem(t)
	original := mustRead(t, store.path(item.ID))
	for _, clear := range [][]string{nil, {}} {
		if _, err := store.SetNoteTags(item.ID, 2, []string{"ply"}, []string{"A"}); err != nil {
			t.Fatal(err)
		}
		replaced, err := store.SetNoteTags(item.ID, 2, []string{"web"}, nil)
		if err != nil || fmt.Sprint(replaced.Notes[2].Repos, replaced.Notes[2].Topics) != "[web] []" {
			t.Fatalf("replace = %+v, %v", replaced.Notes[2], err)
		}
		cleared, err := store.SetNoteTags(item.ID, 2, clear, clear)
		if err != nil || cleared.Notes[2].Repos != nil || cleared.Notes[2].Topics != nil {
			t.Fatalf("clear = %+v, %v", cleared.Notes[2], err)
		}
		if after := mustRead(t, store.path(item.ID)); after != original {
			t.Fatalf("file after clear =\n%s\nwant\n%s", after, original)
		}
	}
}

func TestSetNoteTagsToSameValuesWritesNothing(t *testing.T) {
	store, item := threeNoteItem(t)
	if _, err := store.SetNoteTags(item.ID, 0, []string{"ply"}, []string{"PWA"}); err != nil {
		t.Fatal(err)
	}
	path := store.path(item.ID)
	contents := mustRead(t, path)
	modified := ageFile(t, path)
	readOnlyBacklog(t, store)
	if _, err := store.SetNoteTags(item.ID, 0, []string{" ply", "ply"}, []string{"PWA ", "pwa"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetNoteTags(item.ID, 1, nil, []string{}); err != nil {
		t.Fatal(err)
	}
	requireFileUnchanged(t, path, contents, modified)
	_, err := store.SetNoteTags(item.ID, 0, nil, nil)
	requireClass(t, err, ClassIOError)
	requireFileUnchanged(t, path, contents, modified)
}

func TestNoteTagWritesRejectInvalidInputWithoutWriting(t *testing.T) {
	many := make([]string, 33)
	for i := range many {
		many[i] = fmt.Sprint("r", i)
	}
	cases := map[string]struct {
		id     string
		index  int
		repos  []string
		topics []string
		class  string
		detail string
	}{
		"bad id":           {"../x", 0, nil, nil, ClassInvalidArguments, "invalid item id"},
		"negative index":   {statusItemID, -1, nil, nil, ClassInvalidArguments, "no such note"},
		"index past notes": {statusItemID, 1, nil, nil, ClassInvalidArguments, "no such note"},
		"uppercase repo":   {statusItemID, 0, []string{"Ply"}, nil, ClassInvalidArguments, "repos must be at most 32 lowercase ASCII slugs of at most 100 characters"},
		"33 repos":         {statusItemID, 0, many, nil, ClassInvalidArguments, "repos must be at most 32 lowercase ASCII slugs of at most 100 characters"},
		"topic too long":   {statusItemID, 0, nil, []string{strings.Repeat("ø", 101)}, ClassInvalidArguments, "topics must be at most 32 non-blank lines of at most 100 characters"},
		"topic control":    {statusItemID, 0, nil, []string{"a\x00b"}, ClassInvalidArguments, "topics must be at most 32 non-blank lines of at most 100 characters"},
		"missing item":     {"20261104-00000000", 0, nil, nil, ClassItemNotFound, "20261104-00000000"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			path := writeItemFile(t, store, statusItemID, strings.Replace(statusItemFile("  \"status\": \"done\",\n"), "\"notes\": []", "\"notes\": [{\"at\": \"x\", \"text\": \"n\"}]", 1))
			modified := ageFile(t, path)
			if err := os.WriteFile(filepath.Join(store.dir, ".lock"), nil, 0o644); err != nil { // as after any earlier write
				t.Fatal(err)
			}
			tree := backlogTree(t, store)
			_, err := store.SetNoteTags(c.id, c.index, c.repos, c.topics)
			if err == nil || err.Error() != c.class+": "+c.detail {
				t.Fatalf("set = %v, want %s: %s", err, c.class, c.detail)
			}
			if c.class == ClassInvalidArguments && c.detail != "invalid item id" && c.detail != "no such note" {
				_, err = store.AddNoteTagged(statusItemID, "n", "", c.repos, c.topics)
				requireClass(t, err, c.class)
			}
			requireTree(t, store, tree)
			requireFileUnchanged(t, path, tree[statusItemID+".json"], modified)
		})
	}
	store := testStore(t)
	_, err := store.SetNoteTags(statusItemID, 0, nil, nil)
	requireClass(t, err, ClassItemNotFound)
	_, err = store.AddNoteTagged(statusItemID, "n", "", nil, nil)
	requireClass(t, err, ClassItemNotFound)
	requireTree(t, store, map[string]string{})
}

func TestTaggedNoteSurvivesOrdinaryWriters(t *testing.T) {
	store, item := textItem(t)
	if _, err := store.AddNoteTagged(item.ID, "tagged", "Per", []string{"ply"}, []string{"PWA"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddNoteAs(item.ID, "plain", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetStatus(item.ID, StatusDone); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetTags(item.ID, []string{"web"}, nil); err != nil {
		t.Fatal(err)
	}
	attach(t, store, item.ID, "a", "x")
	read, err := store.Read(item.ID)
	want := Note{At: attachedAt, Text: "tagged", Author: "Per", Repos: []string{"ply"}, Topics: []string{"PWA"}}
	if err != nil || len(read.Notes) != 2 || !reflect.DeepEqual(read.Notes[0], want) || read.Notes[1].Repos != nil || read.Status != StatusDone {
		t.Fatalf("read = %+v, %v", read, err)
	}
}

func TestReadRejectsInvalidStoredNoteTags(t *testing.T) {
	for detail, tags := range map[string]string{
		"invalid note repos":  `"repos":["Ply"]`,
		"invalid note topics": `"topics":[" a"]`,
	} {
		t.Run(detail, func(t *testing.T) {
			store := testStore(t)
			writeItemFile(t, store, statusItemID, "{\"id\":\""+statusItemID+"\",\"title\":\"t\",\"text\":\"x\",\"created_at\":\"2026-11-04T21:30:00Z\","+
				"\"status\":\"open\",\"notes\":[{\"at\":\"2026-11-04T21:30:00Z\",\"text\":\"n\","+tags+"}]}")
			_, err := store.Read(statusItemID)
			if err == nil || err.Error() != ClassItemInvalid+": "+detail {
				t.Fatalf("read = %v, want %s", err, detail)
			}
			result, err := store.List()
			want := []InvalidFile{{File: statusItemID + ".json", Error: detail}}
			if err != nil || len(result.Items) != 0 || !reflect.DeepEqual(result.Invalid, want) {
				t.Fatalf("list = %+v, %v", result, err)
			}
		})
	}
}

func TestParallelNotesAndNoteTagsAreNotLost(t *testing.T) {
	store, item := textItem(t)
	if _, err := store.AddNote(item.ID, "first"); err != nil {
		t.Fatal(err)
	}
	const writers = 20
	var group sync.WaitGroup
	errs := make(chan error, 2*writers)
	for i := 0; i < writers; i++ {
		group.Add(2)
		go func(i int) {
			defer group.Done()
			_, err := store.AddNoteAs(item.ID, fmt.Sprintf("note %d", i), "Per")
			errs <- err
		}(i)
		go func(i int) {
			defer group.Done()
			_, err := store.SetNoteTags(item.ID, 0, []string{"ply"}, []string{fmt.Sprint("topic ", i)})
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
	if err != nil || len(read.Notes) != writers+1 || read.Notes[0].Text != "first" || len(read.Notes[0].Repos) != 1 || len(read.Notes[0].Topics) != 1 {
		t.Fatalf("read = %+v, %v", read, err)
	}
}
