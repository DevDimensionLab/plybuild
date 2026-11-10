package backlog

import (
	"reflect"
	"strings"
	"testing"
)

func TestArchiveCarriesNoteTagsThroughGroupAppendAndRepeat(t *testing.T) {
	f := newArchiveFixture(t)
	item := f.done(t, "")
	item, err := f.store.AddNoteTagged(item.ID, "merket", "", []string{"ply"}, []string{"Web app"})
	if err != nil {
		t.Fatal(err)
	}
	active := backlogTree(t, f.store)
	f.archiveOK(t, item.ID, "g", true)
	f.archiveOK(t, f.done(t, "").ID, "g", false)
	restore(t, f.store.dir, active)
	archived := f.archiveOK(t, item.ID, "g", false)
	group, err := f.archive.ReadGroup("g")
	if err != nil || len(group.Items) != 2 {
		t.Fatalf("group = %+v, %v", group, err)
	}
	for _, got := range []ArchivedItem{archived, group.Items[0]} {
		if len(got.Notes) != 3 || got.Notes[0].Repos != nil ||
			!reflect.DeepEqual(got.Notes[2], ArchivedNote{At: item.Notes[2].At, Text: "merket", Repos: []string{"ply"}, Topics: []string{"Web app"}}) {
			t.Fatalf("notes = %+v", got.Notes)
		}
	}
}

func TestArchiveNoteTagDifferenceIsConflict(t *testing.T) {
	f := newArchiveFixture(t)
	item := f.done(t, "")
	active := backlogTree(t, f.store)
	f.archiveOK(t, item.ID, "g", true)
	restore(t, f.store.dir, active)
	if _, err := f.store.SetNoteTags(item.ID, 0, []string{"ply"}, nil); err != nil {
		t.Fatal(err)
	}
	before := f.tree(t)
	_, err := f.store.Archive(item.ID, f.archive, "g", false)
	requireClass(t, err, ClassArchiveConflict)
	f.requireUnchanged(t, before)
}

func TestArchiveRejectsInvalidNoteTagsInGroup(t *testing.T) {
	for field, value := range map[string]string{"repos": `["Ply"]`, "topics": `["a\u0000b"]`} {
		t.Run(field, func(t *testing.T) {
			f := newArchiveFixture(t)
			target := f.done(t, "")
			writeGroupFile(t, f.archive, "g.json", strings.Replace(groupJSON("g", "20261101-0a0b0c0d"), `"notes":[]`, `"notes":[{"at":"a","text":"x","`+field+`":`+value+`}]`, 1))
			result, err := f.archive.ListGroups()
			if err != nil || !reflect.DeepEqual(result.Invalid, []InvalidFile{{File: "g.json", Error: "invalid note " + field}}) {
				t.Fatalf("list = %+v, %v", result, err)
			}
			before := f.tree(t)
			_, err = f.store.Archive(target.ID, f.archive, "g", false)
			requireClass(t, err, ClassArchiveInvalid)
			f.requireUnchanged(t, before)
		})
	}
}
