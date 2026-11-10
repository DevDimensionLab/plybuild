package backlog

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestNoteAuthorIsOptionalAndTrimmed(t *testing.T) {
	store := testStore(t)
	path := writeItemFile(t, store, statusItemID, statusItemFile("  \"status\": \"open\",\n"))
	old, err := store.Read(statusItemID)
	if err != nil || old.Repos != nil || old.Topics != nil {
		t.Fatalf("read old item = %+v, %v", old, err)
	}
	if _, err := store.AddNote(statusItemID, "anonymous"); err != nil {
		t.Fatal(err)
	}
	noted, err := store.AddNoteAs(statusItemID, "signed", "  Per Otto  ")
	if err != nil || noted.Notes[0].Author != "" || noted.Notes[1].Author != "Per Otto" {
		t.Fatalf("note = %+v, %v", noted, err)
	}
	contents, err := os.ReadFile(path)
	if err != nil || strings.Count(string(contents), "\"author\"") != 1 || !strings.Contains(string(contents), "\"author\": \"Per Otto\"") {
		t.Fatalf("file = %s, %v", contents, err)
	}
}

func TestAddNoteAsRejectsInvalidAuthorWithoutWriting(t *testing.T) {
	for name, author := range map[string]string{
		"newline":      "a\nb",
		"control":      "a\x00b",
		"invalid utf8": "a\xffb",
		"too long":     strings.Repeat("æ", 101),
	} {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			contents := statusItemFile("  \"status\": \"open\",\n")
			path := writeItemFile(t, store, statusItemID, contents)
			modified := ageFile(t, path)
			_, err := store.AddNoteAs(statusItemID, "note", author)
			requireClass(t, err, ClassInvalidArguments)
			requireFileUnchanged(t, path, contents, modified)
		})
	}
}

func TestSetTagsReplacesAndClearsWithoutTouchingContent(t *testing.T) {
	store := testStore(t)
	contents := "{\n  \"id\": \"" + statusItemID + "\",\n  \"title\": \"t\",\n  \"text\": \"\",\n  \"created_at\": \"2026-11-04T21:30:00Z\",\n" +
		"  \"status\": \"ready\",\n  \"notes\": [\n    {\n      \"at\": \"2026-11-04T21:30:00Z\",\n      \"text\": \"n\",\n      \"author\": \"Per\"\n    }\n  ],\n" +
		"  \"file\": {\n    \"name\": \"spec.md\",\n    \"size\": 3\n  }\n}\n"
	path := writeItemFile(t, store, statusItemID, contents)
	item, err := store.SetTags(statusItemID, []string{"ply", " plybuild ", "ply"}, []string{" Backlog og notater ", "PWA", "backlog OG notater"})
	wantRepos, wantTopics := []string{"ply", "plybuild"}, []string{"Backlog og notater", "PWA"}
	if err != nil || !slices.Equal(item.Repos, wantRepos) || !slices.Equal(item.Topics, wantTopics) || item.Notes[0].Author != "Per" || item.File.Name != "spec.md" {
		t.Fatalf("set tags = %+v, %v", item, err)
	}
	result, err := store.List()
	if err != nil || len(result.Items) != 1 || !slices.Equal(result.Items[0].Repos, wantRepos) || !slices.Equal(result.Items[0].Topics, wantTopics) {
		t.Fatalf("list = %+v, %v", result, err)
	}
	if item, err = store.SetTags(statusItemID, nil, []string{}); err != nil || item.Repos != nil || item.Topics != nil {
		t.Fatalf("clear tags = %+v, %v", item, err)
	}
	if after, err := os.ReadFile(path); err != nil || string(after) != contents {
		t.Fatalf("file after clear =\n%s\nwant\n%s", after, contents)
	}
}

func TestSetTagsToSameValuesWritesNothing(t *testing.T) {
	store := testStore(t)
	item, err := store.Add(schedulerTitle, "text")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SetTags(item.ID, []string{"ply"}, []string{"PWA"}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(store.dir, item.ID+".json")
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	modified := ageFile(t, path)
	readOnlyBacklog(t, store)
	if _, err := store.SetTags(item.ID, []string{" ply"}, []string{"PWA "}); err != nil {
		t.Fatal(err)
	}
	requireFileUnchanged(t, path, string(contents), modified)
	_, err = store.SetTags(item.ID, nil, nil)
	requireClass(t, err, ClassIOError)
	requireFileUnchanged(t, path, string(contents), modified)
}

func TestSetTagsRejectsInvalidInputWithoutWriting(t *testing.T) {
	many := make([]string, 33)
	for i := range many {
		many[i] = fmt.Sprint("t", i)
	}
	cases := map[string]struct{ repos, topics []string }{
		"uppercase repo":   {[]string{"Ply"}, nil},
		"repo with space":  {[]string{"a b"}, nil},
		"repo path":        {[]string{"../x"}, nil},
		"empty repo":       {[]string{" "}, nil},
		"empty topic":      {nil, []string{"  "}},
		"topic newline":    {nil, []string{"a\nb"}},
		"topic control":    {nil, []string{"a\tb"}},
		"topic not utf8":   {nil, []string{"a\xff"}},
		"topic too long":   {nil, []string{strings.Repeat("ø", 101)}},
		"too many topics":  {nil, many},
		"too long repo id": {[]string{strings.Repeat("a", 101)}, nil},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			store := testStore(t)
			contents := statusItemFile("  \"status\": \"open\",\n")
			path := writeItemFile(t, store, statusItemID, contents)
			modified := ageFile(t, path)
			_, err := store.SetTags(statusItemID, c.repos, c.topics)
			requireClass(t, err, ClassInvalidArguments)
			requireFileUnchanged(t, path, contents, modified)
		})
	}
	store := testStore(t)
	_, err := store.SetTags("../x", nil, nil)
	requireClass(t, err, ClassInvalidArguments)
	_, err = store.SetTags(statusItemID, nil, nil)
	requireClass(t, err, ClassItemNotFound)
}

func TestReadRejectsInvalidStoredTagsAndAuthor(t *testing.T) {
	item := func(extra, author string) string {
		return "{\"id\":\"" + statusItemID + "\",\"title\":\"t\",\"text\":\"x\",\"created_at\":\"2026-11-04T21:30:00Z\",\"status\":\"open\"," +
			"\"notes\":[{\"at\":\"2026-11-04T21:30:00Z\",\"text\":\"n\"" + author + "}]" + extra + "}"
	}
	cases := map[string]struct{ contents, detail string }{
		"repo":   {item(`,"repos":["Ply"]`, ""), "invalid repos"},
		"topic":  {item(`,"topics":["a\nb"]`, ""), "invalid topics"},
		"author": {item("", `,"author":"a\nb"`), "invalid note author"},
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
}

func TestParallelNotesAndTagsAreNotLost(t *testing.T) {
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
			_, err := store.AddNoteAs(item.ID, fmt.Sprintf("note %d", i), "Per")
			errs <- err
		}(i)
		go func(i int) {
			defer group.Done()
			_, err := store.SetTags(item.ID, []string{"ply"}, []string{fmt.Sprint("topic ", i)})
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
	if err != nil || len(read.Notes) != writers || len(read.Repos) != 1 || len(read.Topics) != 1 || read.Text != "text" {
		t.Fatalf("read = %+v, %v", read, err)
	}
}
