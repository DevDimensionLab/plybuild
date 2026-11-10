package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/backlog"
)

func TestBacklogAttachListAndShowAttachments(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	downloads := t.TempDir()
	for name, contents := range map[string]string{"a.txt": "a", "b.png": "bb", "c.md": "ccc", "d.md": "dddd"} {
		mustPlanWrite(t, filepath.Join(downloads, name), contents, 0o600)
	}
	item := decodeBacklogItem(t, f.success(t, "text", "add", "--project", "ply", "--title", "T", "--format", "json"))
	directory := filepath.Join(f.roots["ply"], "backlog")
	itemFile := filepath.Join(directory, item.ID+".json")
	stored, _ := os.ReadFile(itemFile)
	copied := func(name string) string { return filepath.Join(directory, "files", item.ID, name) }

	t.Chdir(downloads)
	text := f.success(t, "", "attach", "--project", "ply", item.ID, "a.txt", filepath.Join(downloads, "b.png"))
	want := fmt.Sprintf("%s  T\nCreated: %s\nStatus: open\nAttachment: a.txt (1 bytes, present)\nPath: %s\nAttachment: b.png (2 bytes, present)\nPath: %s\n\ntext\n", item.ID, item.CreatedAt, copied("a.txt"), copied("b.png"))
	if text != want {
		t.Fatalf("attach text\n%q\nwant\n%q", text, want)
	}
	if shown := f.success(t, "", "show", "--project", "ply", item.ID); shown != want {
		t.Fatalf("show text %q", shown)
	}

	_, err := f.run("", "attach", "--project", "ply", item.ID, "c.md", "a.txt", "d.md")
	var failure *backlog.Error
	if !asBacklogError(err, &failure) || failure.Class != backlog.ClassAttachmentExists || !strings.Contains(failure.Detail, "attached 1 of 3 files") {
		t.Fatalf("partial attach error = %v", err)
	}
	var listed struct {
		ID          string               `json:"id"`
		Attachments []backlog.Attachment `json:"attachments"`
	}
	if err := json.Unmarshal([]byte(f.success(t, "", "attachments", "--project", "ply", "--format", "json", item.ID)), &listed); err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, attachment := range listed.Attachments {
		names = append(names, attachment.Name)
	}
	if listed.ID != item.ID || !reflect.DeepEqual(names, []string{"a.txt", "b.png", "c.md"}) {
		t.Fatalf("attachments = %#v", listed)
	}
	lines := f.success(t, "", "attachments", "--project", "ply", item.ID)
	wantLine := fmt.Sprintf("c.md  3 bytes  %s  present  %s\n", listed.Attachments[2].AddedAt, copied("c.md"))
	if strings.Count(lines, "\n") != 3 || !strings.HasPrefix(lines, "a.txt  1 bytes  ") || !strings.HasSuffix(lines, wantLine) {
		t.Fatalf("attachments text %q", lines)
	}

	shownJSON := f.success(t, "", "show", "--project", "ply", "--format", "json", item.ID)
	var shown struct {
		backlog.Item
		Attachments []backlog.Attachment `json:"attachments"`
	}
	if err := json.Unmarshal([]byte(shownJSON), &shown); err != nil || !reflect.DeepEqual(shown.Item, item) || !reflect.DeepEqual(shown.Attachments, listed.Attachments) {
		t.Fatalf("show json %s: %v", shownJSON, err)
	}
	if attached := f.success(t, "", "attach", "--project", "ply", "--format", "json", item.ID, "d.md"); !strings.HasPrefix(attached, "{\n  \"id\": \""+item.ID+"\",\n  \"attachments\": [") || !strings.Contains(attached, "\"d.md\"") {
		t.Fatalf("attach json %s", attached)
	}
	if after, _ := os.ReadFile(itemFile); string(after) != string(stored) {
		t.Fatal("item file changed")
	}

	plain := decodeBacklogItem(t, f.success(t, "x", "add", "--project", "ply", "--title", "P", "--format", "json"))
	plainFile, _ := os.ReadFile(filepath.Join(directory, plain.ID+".json"))
	if got := f.success(t, "", "show", "--project", "ply", "--format", "json", plain.ID); got != string(plainFile) {
		t.Fatalf("show json without attachments %q", got)
	}
	if got := f.success(t, "", "attachments", "--project", "ply", plain.ID); got != "No attachments on item "+plain.ID+".\n" {
		t.Fatalf("empty attachments %q", got)
	}
	before := snapshotTree(t, f.roots)
	for _, args := range [][]string{
		{"attach", "--project", "ply", item.ID},
		{"attach", item.ID, "a.txt"},
		{"attachments", "--project", "ply"},
	} {
		_, err := f.run("", args...)
		if !asBacklogError(err, &failure) || failure.Class != backlog.ClassInvalidArguments {
			t.Fatalf("%v: %v", args, err)
		}
	}
	if !reflect.DeepEqual(snapshotTree(t, f.roots), before) {
		t.Fatal("rejected arguments changed files")
	}
}

func asBacklogError(err error, target **backlog.Error) bool {
	failure, ok := err.(*backlog.Error)
	*target = failure
	return ok
}
