package backlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// The v1 fixture and its expectations are frozen; they must not follow changes to the active Item rules.
func TestArchiveV1FixtureIsReadAndKept(t *testing.T) {
	var repos, topics []string
	for i := 0; i < 31; i++ {
		repos, topics = append(repos, fmt.Sprintf("r%02d", i)), append(topics, fmt.Sprintf("Topic %02d", i))
	}
	repos, topics = append(repos, "a.b_c-"+strings.Repeat("x", 94)), append(topics, strings.Repeat("Æ", 100))
	want := []ArchivedItem{
		{ID: "20261106-4458fcb9", Title: "Arkivering av ferdig backlog", Text: "Tekst <æ> & \"sitat\"\n", CreatedAt: "2026-11-06T04:45:02Z",
			Notes: []ArchivedNote{{At: "2026-11-06T06:00:00Z", Text: "med forfatter", Author: strings.Repeat("ø", 100)}, {At: "2026-11-06T07:00:00Z", Text: "uten"}},
			Repos: repos, Topics: topics, ArchivedAt: "2026-11-07T08:00:00Z"},
		{ID: "20261106-7c1e0a52", Title: "Spec", Text: "", CreatedAt: "2026-11-06T09:00:00Z", Notes: []ArchivedNote{},
			File: &ArchivedFile{Name: "my spec ✓ æ.md", Size: 10999}, ArchivedAt: "2026-11-07T08:05:00Z"},
	}
	fixture, err := os.ReadFile(filepath.Join("testdata", "archive-v1", "fixture-v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	f := newArchiveFixture(t)
	writeGroupFile(t, f.archive, "fixture-v1.json", string(fixture))
	group, err := f.archive.ReadGroup("fixture-v1")
	if err != nil || group.Format != ArchiveFormat || group.Group != "fixture-v1" || !reflect.DeepEqual(group.Items, want) {
		t.Fatalf("group = %#v, %v", group, err)
	}
	added := f.archiveOK(t, f.done(t, "").ID, "fixture-v1", false)
	group, err = f.archive.ReadGroup("fixture-v1")
	if err != nil || !reflect.DeepEqual(group.Items, append(want, added)) {
		t.Fatalf("after archiving = %#v, %v", group, err)
	}

	invalid, err := os.ReadFile(filepath.Join("testdata", "archive-v1", "invalid.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range bytes.Split(bytes.TrimSuffix(invalid, []byte("\n")), []byte("\n")) {
		name, doc, _ := bytes.Cut(line, []byte(" "))
		t.Run(string(name), func(t *testing.T) {
			f := newArchiveFixture(t)
			writeGroupFile(t, f.archive, string(name)+".json", string(doc))
			result, err := f.archive.ListGroups()
			if err != nil || len(result.Groups) != 0 || len(result.Invalid) != 1 {
				t.Fatalf("list = %+v, %v", result, err)
			}
		})
	}
}
