package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/backlog"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestBacklogArchiveGroupsAndGroupCommands(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	f.dataPaths["ply"] = map[string]string{workspace.DataPathBacklogArchive: "data/archive/backlog"}
	root := f.roots["ply"]
	archiveDir := filepath.Join(root, "data", "archive", "backlog")
	if out := f.success(t, "", "groups", "--project", "ply"); out != "No archive groups in project ply.\n" {
		t.Fatalf("empty groups = %q", out)
	}
	if out := f.success(t, "", "groups", "--project", "ply", "--format", "json"); out != "{\n  \"project\": \"ply\",\n  \"groups\": [],\n  \"invalid\": []\n}\n" {
		t.Fatalf("empty groups json = %q", out)
	}
	const textID = "20261101-0a0b0c0d"
	writeBacklogItem(t, filepath.Join(root, "backlog"), textID, `{"id":"`+textID+`","title":"Tekst","text":"Brødtekst\n","created_at":"2026-11-01T10:00:00Z","status":"done",
"notes":[{"at":"2026-11-02T10:00:00Z","text":"med","author":"ply-admin"},{"at":"2026-11-03T10:00:00Z","text":"uten"}],"repos":["ply","gone"],"topics":["Web app"]}`)
	source := filepath.Join(t.TempDir(), "my spec.md")
	if err := os.WriteFile(source, []byte("innhold"), 0o600); err != nil {
		t.Fatal(err)
	}
	fileItem := decodeBacklogItem(t, f.success(t, "", "add", "--project", "ply", "--title", "Spec", "--file", source, "--format", "json"))

	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"archive", "--project", "ply", textID}, backlog.ClassInvalidArguments + ": --group is required"},
		{[]string{"archive", "--project", "ply", "--group", "ukjent", textID}, backlog.ClassGroupNotFound + ": group ukjent does not exist; see ply backlog groups"},
		{[]string{"archive", "--project", "ply", "--group", "g", "--new-group", fileItem.ID}, backlog.ClassItemNotDone + ": item " + fileItem.ID + " is not done"},
		{[]string{"group", "--project", "ply", "g"}, backlog.ClassGroupNotFound + ": group g does not exist; see ply backlog groups"},
	} {
		if _, err := f.run("", test.args...); err == nil || err.Error() != test.want || strings.Contains(err.Error(), root) {
			t.Fatalf("%v: err = %v", test.args, err)
		}
	}
	f.success(t, "", "status", "--project", "ply", "--value", "done", fileItem.ID)
	if out := f.success(t, "", "archive", "--project", "ply", "--group", "g", "--new-group", textID); out != "Archived "+textID+" to g.\n" {
		t.Fatalf("archive = %q", out)
	}
	var archived struct {
		backlog.ArchivedItem
		Group string `json:"group"`
	}
	if err := json.Unmarshal([]byte(f.success(t, "", "archive", "--project", "ply", "--group", "g", fileItem.ID, "--format", "json")), &archived); err != nil ||
		archived.Group != "g" || archived.ID != fileItem.ID || archived.File == nil || archived.File.Name != "my spec.md" {
		t.Fatalf("archive json = %+v, %v", archived, err)
	}

	stored, err := os.ReadFile(filepath.Join(archiveDir, "g.json"))
	if err != nil {
		t.Fatal(err)
	}
	if out := f.success(t, "", "group", "--project", "ply", "g", "--format", "json"); out != string(stored) {
		t.Fatalf("group json = %s\nstored %s", out, stored)
	}
	var group backlog.ArchiveGroup
	if err := json.Unmarshal(stored, &group); err != nil || len(group.Items) != 2 {
		t.Fatalf("group = %+v, %v", group, err)
	}
	first, second := group.Items[0], group.Items[1]
	copyPath := filepath.Join(archiveDir, "files", fileItem.ID, "my spec.md")
	want := fmt.Sprintf("%s  Tekst\nCreated: 2026-11-01T10:00:00Z\nArchived: %s\nRepos: ply, gone\nTopics: Web app\n\nBrødtekst\n"+
		"\nNote 1 2026-11-02T10:00:00Z by ply-admin:\nmed\n\nNote 2 2026-11-03T10:00:00Z:\nuten\n"+
		"\n%s  Spec\nCreated: %s\nArchived: %s\nFile: my spec.md (7 bytes, present)\nPath: %s\n\n\n",
		textID, first.ArchivedAt, fileItem.ID, second.CreatedAt, second.ArchivedAt, copyPath)
	if out := f.success(t, "", "group", "--project", "ply", "g"); out != want {
		t.Fatalf("group text = %q\nwant         %q", out, want)
	}
	if out := f.success(t, "", "groups", "--project", "ply"); out != "g  2 items  "+second.ArchivedAt+"\n" {
		t.Fatalf("groups = %q", out)
	}
	if err := os.WriteFile(filepath.Join(archiveDir, "broken.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out := f.success(t, "", "groups", "--project", "ply"); out != "g  2 items  "+second.ArchivedAt+"\ninvalid: broken.json: invalid JSON at offset 1\n" {
		t.Fatalf("groups with invalid = %q", out)
	}
	// A hand-formatted valid v1 file is shown as stored, not re-encoded.
	const stored2 = `{"format":"ply-backlog-archive-group/v1","group":"hand","items":[{"id":"20261101-0a0b0c0e","title":"t","text":"<x>","created_at":"c","notes":null,"archived_at":"a"}]}`
	if err := os.WriteFile(filepath.Join(archiveDir, "hand.json"), []byte(stored2), 0o644); err != nil {
		t.Fatal(err)
	}
	if out := f.success(t, "", "group", "--project", "ply", "hand", "--format", "json"); out != stored2 {
		t.Fatalf("hand group json = %q", out)
	}
	if out := f.success(t, "", "list", "--project", "ply"); out != "No backlog items in project ply.\n" {
		t.Fatalf("list = %q", out)
	}
	if _, err := f.run("", "show", "--project", "ply", textID); err == nil || err.Error() != backlog.ClassItemNotFound+": "+textID {
		t.Fatalf("show archived = %v", err)
	}
}
