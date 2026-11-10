package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/backlog"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestBacklogNoteTagsThroughNoteTagsShowListAndGroup(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "backlog")
	run := func(args ...string) (string, error) {
		command := newBacklogCommandWithDirectory(func(string) (backlogPaths, error) {
			return backlogPaths{directory: dir, archive: filepath.Join(root, "archive"), wrapper: root, members: []workspace.RepoID{"code", "web"}}, nil
		})
		command.SilenceUsage, command.SilenceErrors = true, true
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(&bytes.Buffer{})
		command.SetIn(strings.NewReader("Notat\n"))
		command.SetArgs(append(args, "--project", "ply"))
		err := command.Execute()
		return stdout.String(), err
	}
	ok := func(args ...string) string {
		t.Helper()
		out, err := run(args...)
		if err != nil {
			t.Fatalf("backlog %v: %v", args, err)
		}
		return out
	}
	item := decodeBacklogItem(t, ok("add", "--title", "T", "--format", "json"))
	ok("note", item.ID)
	noted := decodeBacklogItem(t, ok("note", item.ID, "--author", "Per", "--repo", "web", "--topic", " Web app ", "--topic", "web APP", "--format", "json"))
	if n := noted.Notes[1]; noted.Notes[0].Repos != nil || !reflect.DeepEqual([]any{n.Repos, n.Topics, n.Author}, []any{[]string{"web"}, []string{"Web app"}, "Per"}) {
		t.Fatalf("noted = %#v", noted.Notes)
	}
	path := filepath.Join(dir, item.ID+".json")
	stored, _ := os.ReadFile(path)
	for _, bad := range [][]string{
		{"note", item.ID, "--repo", "ghost"},
		{"tags", item.ID, "--note", "0", "--repo", "code"},
		{"tags", item.ID, "--note", "3", "--repo", "code"},
		{"tags", item.ID, "--note", "1", "--repo", "ghost"},
		{"tags", item.ID, "--note", "1", "--clear", "--repo", "code"},
		{"tags", item.ID, "--note", "1"},
	} {
		if _, err := run(bad...); err == nil || !strings.HasPrefix(err.Error(), backlog.ClassInvalidArguments+": ") {
			t.Fatalf("%v: err = %v", bad, err)
		}
	}
	if after, _ := os.ReadFile(path); string(after) != string(stored) {
		t.Fatal("rejected input changed the item")
	}

	retagged := decodeBacklogItem(t, ok("tags", item.ID, "--note", "1", "--repo", "code", "--topic", "API", "--format", "json"))
	if !reflect.DeepEqual(retagged.Notes[0].Repos, []string{"code"}) || !reflect.DeepEqual(retagged.Notes[1].Repos, []string{"web"}) || retagged.Repos != nil {
		t.Fatalf("retagged = %#v", retagged)
	}
	if cleared := decodeBacklogItem(t, ok("tags", item.ID, "--note", "1", "--clear", "--format", "json")); cleared.Notes[0].Repos != nil || cleared.Notes[0].Topics != nil {
		t.Fatalf("cleared = %#v", cleared.Notes[0])
	}
	if list := decodeBacklogList(t, ok("list", "--repo", "web", "--format", "json")); len(list.Items) != 0 {
		t.Fatalf("note tags rolled up into list: %+v", list.Items)
	}
	wantNotes := "\nNote 1 " + noted.Notes[0].At + ":\nNotat\n\nNote 2 " + noted.Notes[1].At + " by Per:\nRepos: web\nTopics: Web app\nNotat\n"
	if show := ok("show", item.ID); !strings.HasSuffix(show, wantNotes) {
		t.Fatalf("show = %q", show)
	}
	if show := decodeBacklogItem(t, ok("show", item.ID, "--format", "json")); !reflect.DeepEqual(show.Notes[1].Topics, []string{"Web app"}) {
		t.Fatalf("show json = %#v", show.Notes)
	}

	ok("status", item.ID, "--value", "done")
	ok("archive", item.ID, "--group", "g", "--new-group")
	if group := ok("group", "g"); !strings.HasSuffix(group, wantNotes) {
		t.Fatalf("group = %q", group)
	}
	if group := ok("group", "g", "--format", "json"); !strings.Contains(group, `"repos": [`) {
		t.Fatalf("group json = %s", group)
	}
}
