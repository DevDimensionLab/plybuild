package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/backlog"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type backlogFixture struct {
	roots     map[string]string
	dataPaths map[string]map[string]string
}

func newBacklogFixture(t *testing.T, projects ...string) backlogFixture {
	t.Helper()
	f := backlogFixture{roots: map[string]string{}, dataPaths: map[string]map[string]string{}}
	for _, project := range projects {
		f.roots[project] = t.TempDir()
	}
	return f
}

// show stands in for workspace.ShowProject so the real data path resolution runs.
func (f backlogFixture) show(id workspace.ProjectID) (workspace.ProjectResult, error) {
	if root, ok := f.roots[string(id)]; ok {
		return workspace.ProjectResult{Project: workspace.ProjectRecord{ID: id, Wrapper: root, DataPaths: f.dataPaths[string(id)]}}, nil
	}
	return workspace.ProjectResult{}, &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: "project " + string(id) + " is not registered in Ply workspace /ws"}
}

func (f backlogFixture) run(stdin string, args ...string) (string, error) {
	command := newBacklogCommandWithDirectory(backlogDirectory(f.show))
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetIn(strings.NewReader(stdin))
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func (f backlogFixture) success(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	out, err := f.run(stdin, args...)
	if err != nil {
		t.Fatalf("backlog %v: %v", args, err)
	}
	return out
}

func snapshotTree(t *testing.T, roots map[string]string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				out[path] = "dir"
				return nil
			}
			b, err := os.ReadFile(path)
			out[path] = string(b)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

func decodeBacklogItem(t *testing.T, raw string) backlog.Item {
	t.Helper()
	var item backlog.Item
	if err := json.Unmarshal([]byte(raw), &item); err != nil {
		t.Fatalf("invalid item JSON %q: %v", raw, err)
	}
	return item
}

type backlogListOutput struct {
	Project string                `json:"project"`
	Items   []backlog.Summary     `json:"items"`
	Invalid []backlog.InvalidFile `json:"invalid"`
}

func decodeBacklogList(t *testing.T, raw string) backlogListOutput {
	t.Helper()
	var list backlogListOutput
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("invalid list JSON %q: %v", raw, err)
	}
	return list
}

func TestBacklogListWithoutDirectoryIsEmptyAndCreatesNothing(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	out := f.success(t, "", "list", "--project", "ply", "--format", "json")
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw["project"]) != `"ply"` || string(raw["items"]) != "[]" || string(raw["invalid"]) != "[]" || len(raw) != 3 {
		t.Fatalf("list = %s", out)
	}
	if _, err := os.Stat(filepath.Join(f.roots["ply"], "backlog")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("list created backlog directory: %v", err)
	}
}

func TestBacklogAddShowNoteListUseSelectedProjectRoot(t *testing.T) {
	f := newBacklogFixture(t, "ply", "other")
	text := "  a\n\n## b\n---\n"
	added := f.success(t, text, "add", "--project", "ply", "--title", "  Scheduler: kjør nattlig opprydding ", "--format", "json")
	item := decodeBacklogItem(t, added)
	if item.Title != "Scheduler: kjør nattlig opprydding" || item.Text != text || item.Status != backlog.StatusOpen || len(item.Notes) != 0 {
		t.Fatalf("added item = %#v", item)
	}
	file, err := os.ReadFile(filepath.Join(f.roots["ply"], "backlog", item.ID+".json"))
	if err != nil {
		t.Fatal(err)
	}
	if added != string(file) {
		t.Fatalf("add stdout differs from item file:\n%s\n%s", added, file)
	}
	if shown := f.success(t, "", "show", "--project", "ply", item.ID, "--format", "json"); shown != string(file) {
		t.Fatalf("show = %s, want %s", shown, file)
	}
	if _, err := os.Stat(filepath.Join(f.roots["other"], "backlog")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("add wrote outside selected project: %v", err)
	}
	if other := decodeBacklogList(t, f.success(t, "", "list", "--project", "other", "--format", "json")); other.Project != "other" || len(other.Items) != 0 {
		t.Fatalf("other list = %#v", other)
	}

	first := decodeBacklogItem(t, f.success(t, "første\n", "note", "--project", "ply", item.ID, "--format", "json"))
	second := f.success(t, "andre", "note", "--project", "ply", item.ID, "--format", "json")
	noted := decodeBacklogItem(t, second)
	if len(first.Notes) != 1 || len(noted.Notes) != 2 || noted.Notes[0].Text != "første\n" || noted.Notes[1].Text != "andre" ||
		noted.Text != item.Text || noted.CreatedAt != item.CreatedAt || noted.Title != item.Title {
		t.Fatalf("noted item = %#v", noted)
	}
	if file, _ := os.ReadFile(filepath.Join(f.roots["ply"], "backlog", item.ID+".json")); second != string(file) {
		t.Fatalf("note stdout differs from item file:\n%s\n%s", second, file)
	}
	list := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--format", "json"))
	want := []backlog.Summary{{ID: item.ID, Title: item.Title, CreatedAt: item.CreatedAt, UpdatedAt: noted.Notes[1].At, NoteCount: 2, Status: backlog.StatusOpen}}
	if list.Project != "ply" || !reflect.DeepEqual(list.Items, want) || len(list.Invalid) != 0 {
		t.Fatalf("list = %#v", list)
	}
}

func TestBacklogTextFormatIsDefault(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	if out := f.success(t, "", "list", "--project", "ply"); out != "No backlog items in project ply.\n" {
		t.Fatalf("empty list text = %q", out)
	}
	item := decodeBacklogItem(t, f.success(t, "Tekst\n", "add", "--project", "ply", "--title", "Tittel", "--format", "json"))
	f.success(t, "Note\n", "note", "--project", "ply", item.ID)
	shown := f.success(t, "", "show", "--project", "ply", item.ID)
	noted := decodeBacklogItem(t, f.success(t, "", "show", "--project", "ply", item.ID, "--format", "json"))
	wantShow := item.ID + "  Tittel\nCreated: " + item.CreatedAt + "\nStatus: open\n\nTekst\n\nNote 1 " + noted.Notes[0].At + ":\nNote\n"
	if shown != wantShow {
		t.Fatalf("show text = %q, want %q", shown, wantShow)
	}
	if listed := f.success(t, "", "list", "--project", "ply"); listed != item.ID+"  "+noted.Notes[0].At+"  1 notes  open  Tittel\n" {
		t.Fatalf("list text = %q", listed)
	}
	wantReady := strings.Replace(wantShow, "Status: open\n", "Status: ready\n", 1)
	if readied := f.success(t, "", "status", "--project", "ply", "--value", "ready", item.ID); readied != wantReady {
		t.Fatalf("status text = %q, want %q", readied, wantReady)
	}
	if listed := f.success(t, "", "list", "--project", "ply"); listed != item.ID+"  "+noted.Notes[0].At+"  1 notes  ready  Tittel\n" {
		t.Fatalf("ready list text = %q", listed)
	}
	if listed := f.success(t, "", "list", "--project", "ply", "--status", "done"); listed != "No done backlog items in project ply.\n" {
		t.Fatalf("empty filtered list text = %q", listed)
	}
}

func TestBacklogJSONOutputContainsStatusAndNoReady(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	added := f.success(t, "Tekst\n", "add", "--project", "ply", "--title", "Tittel", "--format", "json")
	item := decodeBacklogItem(t, added)
	noted := f.success(t, "Note\n", "note", "--project", "ply", item.ID, "--format", "json")
	shown := f.success(t, "", "show", "--project", "ply", item.ID, "--format", "json")
	status := f.success(t, "", "status", "--project", "ply", "--value", "open", item.ID, "--format", "json")
	for name, out := range map[string]string{"add": added, "note": noted, "show": shown, "status": status} {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &raw); err != nil {
			t.Fatal(err)
		}
		if _, ok := raw["ready"]; ok || string(raw["status"]) != `"open"` {
			t.Fatalf("%s output = %s", name, out)
		}
	}
	var list struct {
		Items []map[string]json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal([]byte(f.success(t, "", "list", "--project", "ply", "--format", "json")), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || string(list.Items[0]["status"]) != `"open"` {
		t.Fatalf("list items = %v", list.Items)
	}
	if _, ok := list.Items[0]["ready"]; ok {
		t.Fatalf("list item has ready: %v", list.Items[0])
	}
}

func TestBacklogStatusSetsEveryValueAndRepeatsWithoutWriting(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	item := decodeBacklogItem(t, f.success(t, "Tekst\n", "add", "--project", "ply", "--title", "Tittel", "--format", "json"))
	f.success(t, "Note\n", "note", "--project", "ply", item.ID)
	path := filepath.Join(f.roots["ply"], "backlog", item.ID+".json")
	before := decodeBacklogItem(t, f.success(t, "", "show", "--project", "ply", item.ID, "--format", "json"))
	for _, step := range []struct {
		value backlog.Status
		write bool
	}{
		{backlog.StatusOpen, false},
		{backlog.StatusReady, true}, {backlog.StatusReady, false},
		{backlog.StatusDone, true}, {backlog.StatusDone, false},
		{backlog.StatusOpen, true},
		{backlog.StatusDone, true}, {backlog.StatusReady, true}, {backlog.StatusOpen, true},
	} {
		file, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		out := f.success(t, "", "status", "--project", "ply", "--value", string(step.value), item.ID, "--format", "json")
		if shown := f.success(t, "", "show", "--project", "ply", item.ID, "--format", "json"); out != shown {
			t.Fatalf("status %s stdout = %s, want show %s", step.value, out, shown)
		}
		got := decodeBacklogItem(t, out)
		want := before
		want.Status = step.value
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("status %s item = %#v, want %#v", step.value, got, want)
		}
		after, _ := os.ReadFile(path)
		afterInfo, _ := os.Stat(path)
		if step.write && out != string(after) {
			t.Fatalf("status %s file = %s, want %s", step.value, after, out)
		}
		if !step.write && (string(after) != string(file) || !afterInfo.ModTime().Equal(info.ModTime())) {
			t.Fatalf("no-op status %s rewrote the item file", step.value)
		}
	}
}

func TestBacklogNoteKeepsEveryStatus(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	item := decodeBacklogItem(t, f.success(t, "Tekst\n", "add", "--project", "ply", "--title", "Tittel", "--format", "json"))
	for _, status := range []backlog.Status{backlog.StatusReady, backlog.StatusDone, backlog.StatusOpen} {
		f.success(t, "", "status", "--project", "ply", "--value", string(status), item.ID)
		noted := decodeBacklogItem(t, f.success(t, "Note\n", "note", "--project", "ply", item.ID, "--format", "json"))
		if noted.Status != status {
			t.Fatalf("note changed status %s to %s", status, noted.Status)
		}
	}
}

func TestBacklogListFiltersByStatusAndKeepsInvalidFiles(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	ids := map[backlog.Status]string{}
	for _, status := range []backlog.Status{backlog.StatusOpen, backlog.StatusReady, backlog.StatusDone} {
		item := decodeBacklogItem(t, f.success(t, "x", "add", "--project", "ply", "--title", string(status), "--format", "json"))
		f.success(t, "", "status", "--project", "ply", "--value", string(status), item.ID)
		ids[status] = item.ID
	}
	writeBacklogItem(t, filepath.Join(f.roots["ply"], "backlog"), "20261103-0badf00d", "{")

	all := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--format", "json"))
	if len(all.Items) != 3 || len(all.Invalid) != 1 {
		t.Fatalf("unfiltered list = %#v", all)
	}
	for status, id := range ids {
		list := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--status", string(status), "--format", "json"))
		if list.Project != "ply" || len(list.Items) != 1 || list.Items[0].ID != id || list.Items[0].Status != status ||
			len(list.Invalid) != 1 || list.Invalid[0].File != "20261103-0badf00d.json" {
			t.Fatalf("list --status %s = %#v", status, list)
		}
		text := f.success(t, "", "list", "--project", "ply", "--status", string(status))
		if !strings.HasPrefix(text, id+"  ") || strings.Count(text, "\n") != 2 || !strings.Contains(text, "invalid: 20261103-0badf00d.json: ") {
			t.Fatalf("list text --status %s = %q", status, text)
		}
	}
}

func TestBacklogItemWithoutValidStatusIsInvalidAndNeverRewritten(t *testing.T) {
	const id = "20261104-47530e79"
	for name, status := range map[string]string{
		"missing":      "",
		"null":         "  \"status\": null,\n",
		"unknown":      "  \"status\": \"closed\",\n",
		"legacy ready": "  \"ready\": true,\n",
	} {
		t.Run(name, func(t *testing.T) {
			f := newBacklogFixture(t, "ply")
			contents := "{\n  \"id\": \"" + id + "\",\n  \"title\": \"Gammel\",\n  \"text\": \"tekst\\n\",\n  \"created_at\": \"2026-11-04T10:00:00Z\",\n" + status + "  \"notes\": []\n}\n"
			path := writeBacklogItem(t, filepath.Join(f.roots["ply"], "backlog"), id, contents)
			for _, args := range [][]string{
				{"list", "--project", "ply", "--format", "json"},
				{"list", "--project", "ply", "--status", "open", "--format", "json"},
			} {
				list := decodeBacklogList(t, f.success(t, "", args...))
				if len(list.Items) != 0 || len(list.Invalid) != 1 || list.Invalid[0].File != id+".json" {
					t.Fatalf("backlog %v = %#v", args, list)
				}
			}
			for _, args := range [][]string{
				{"show", "--project", "ply", id},
				{"note", "--project", "ply", id},
				{"status", "--project", "ply", "--value", "open", id},
				{"status", "--project", "ply", "--value", "ready", id},
			} {
				if out, err := f.run("x", args...); err == nil || out != "" || !strings.HasPrefix(err.Error(), backlog.ClassItemInvalid+": ") {
					t.Fatalf("backlog %v: stdout=%q err=%v", args, out, err)
				}
			}
			if file, _ := os.ReadFile(path); string(file) != contents {
				t.Fatalf("invalid file rewritten: %s", file)
			}
		})
	}
}

func TestBacklogRemovedReadyCommandIsRejected(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	item := decodeBacklogItem(t, f.success(t, "x", "add", "--project", "ply", "--title", "T", "--format", "json"))
	before := snapshotTree(t, f.roots)
	// Nested as under RootCmd, where Cobra does not report unknown subcommands itself.
	parent := &cobra.Command{Use: "ply"}
	parent.AddCommand(newBacklogCommandWithDirectory(backlogDirectory(f.show)))
	parent.SilenceUsage, parent.SilenceErrors = true, true
	for _, args := range [][]string{
		{"backlog", "ready", "--project", "ply", "--value", "true", item.ID},
		{"backlog", "ready", item.ID},
	} {
		var stdout bytes.Buffer
		parent.SetOut(&stdout)
		parent.SetErr(&bytes.Buffer{})
		parent.SetArgs(args)
		if err := parent.Execute(); err == nil || stdout.Len() != 0 || !strings.HasPrefix(err.Error(), backlog.ClassInvalidArguments+": ") {
			t.Fatalf("%v: stdout=%q err=%v", args, stdout.String(), err)
		}
	}
	if after := snapshotTree(t, f.roots); !reflect.DeepEqual(after, before) {
		t.Fatalf("files changed:\nbefore %v\nafter  %v", before, after)
	}
}

func TestBacklogInvalidItemFileIsListedAndRejectedByShow(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	item := decodeBacklogItem(t, f.success(t, "gyldig\n", "add", "--project", "ply", "--title", "Gyldig", "--format", "json"))
	if err := os.WriteFile(filepath.Join(f.roots["ply"], "backlog", "20261103-0badf00d.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	list := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--format", "json"))
	if len(list.Items) != 1 || list.Items[0].ID != item.ID || len(list.Invalid) != 1 || list.Invalid[0].File != "20261103-0badf00d.json" {
		t.Fatalf("list = %#v", list)
	}
	if text := f.success(t, "", "list", "--project", "ply"); !strings.Contains(text, "invalid: 20261103-0badf00d.json: ") {
		t.Fatalf("list text = %q", text)
	}
	if _, err := f.run("", "show", "--project", "ply", "20261103-0badf00d"); err == nil || !strings.HasPrefix(err.Error(), backlog.ClassItemInvalid+": ") {
		t.Fatalf("show invalid = %v", err)
	}
}

func TestBacklogRejectedInputWritesNothing(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	item := decodeBacklogItem(t, f.success(t, "eksisterende\n", "add", "--project", "ply", "--title", "Eksisterende", "--format", "json"))
	before := snapshotTree(t, f.roots)
	tests := []struct {
		name  string
		stdin string
		args  []string
		class string
	}{
		{"empty text", "", []string{"add", "--project", "ply", "--title", "T"}, backlog.ClassInvalidArguments},
		{"blank text", " \n\t", []string{"add", "--project", "ply", "--title", "T"}, backlog.ClassInvalidArguments},
		{"too large text", strings.Repeat("x", 64*1024+1), []string{"add", "--project", "ply", "--title", "T"}, backlog.ClassInvalidArguments},
		{"invalid UTF-8", "\xff\n", []string{"add", "--project", "ply", "--title", "T"}, backlog.ClassInvalidArguments},
		{"multiline title", "tekst", []string{"add", "--project", "ply", "--title", "a\nb"}, backlog.ClassInvalidArguments},
		{"missing title", "tekst", []string{"add", "--project", "ply"}, backlog.ClassInvalidArguments},
		{"blank note", "  ", []string{"note", "--project", "ply", item.ID}, backlog.ClassInvalidArguments},
		{"unknown item note", "tekst", []string{"note", "--project", "ply", "20261104-00000000"}, backlog.ClassItemNotFound},
		{"unknown item show", "", []string{"show", "--project", "ply", "20261104-00000000"}, backlog.ClassItemNotFound},
		{"invalid item id", "", []string{"show", "--project", "ply", "../x"}, backlog.ClassInvalidArguments},
		{"unknown project add", "tekst", []string{"add", "--project", "ghost", "--title", "T"}, string(workspace.ErrorProjectNotFound)},
		{"unknown project list", "", []string{"list", "--project", "ghost"}, string(workspace.ErrorProjectNotFound)},
		{"missing project", "", []string{"list"}, backlog.ClassInvalidArguments},
		{"bad format", "", []string{"list", "--project", "ply", "--format", "yaml"}, backlog.ClassInvalidArguments},
		{"unknown flag", "", []string{"list", "--project", "ply", "--ready"}, backlog.ClassInvalidArguments},
		{"unknown status filter", "", []string{"list", "--project", "ply", "--status", "closed"}, backlog.ClassInvalidArguments},
		{"empty status filter", "", []string{"list", "--project", "ply", "--status="}, backlog.ClassInvalidArguments},
		{"uppercase status filter", "", []string{"list", "--project", "ply", "--status", "Ready"}, backlog.ClassInvalidArguments},
		{"missing item id", "tekst", []string{"note", "--project", "ply"}, backlog.ClassInvalidArguments},
		{"extra argument", "", []string{"list", "--project", "ply", "x"}, backlog.ClassInvalidArguments},
		{"status missing value", "", []string{"status", "--project", "ply", item.ID}, backlog.ClassInvalidArguments},
		{"status empty value", "", []string{"status", "--project", "ply", "--value=", item.ID}, backlog.ClassInvalidArguments},
		{"status uppercase value", "", []string{"status", "--project", "ply", "--value=Done", item.ID}, backlog.ClassInvalidArguments},
		{"status unknown value", "", []string{"status", "--project", "ply", "--value=closed", item.ID}, backlog.ClassInvalidArguments},
		{"status old boolean value", "", []string{"status", "--project", "ply", "--value=true", item.ID}, backlog.ClassInvalidArguments},
		{"status extra argument", "", []string{"status", "--project", "ply", "--value", "done", item.ID, "x"}, backlog.ClassInvalidArguments},
		{"status missing item id", "", []string{"status", "--project", "ply", "--value", "done"}, backlog.ClassInvalidArguments},
		{"status invalid item id", "", []string{"status", "--project", "ply", "--value", "done", "../x"}, backlog.ClassInvalidArguments},
		{"status unknown item", "", []string{"status", "--project", "ply", "--value", "done", "20261104-00000000"}, backlog.ClassItemNotFound},
		{"status unknown project", "", []string{"status", "--project", "ghost", "--value", "done", item.ID}, string(workspace.ErrorProjectNotFound)},
		{"status bare value flag", "", []string{"status", "--project", "ply", "--value", item.ID}, backlog.ClassInvalidArguments},
		{"status on show", "", []string{"show", "--project", "ply", "--status", "done", item.ID}, backlog.ClassInvalidArguments},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := f.run(test.stdin, test.args...)
			if err == nil || out != "" || !strings.HasPrefix(err.Error(), test.class+": ") || strings.Contains(err.Error(), "\n") {
				t.Fatalf("stdout=%q err=%v, want class %s", out, err, test.class)
			}
			for _, root := range f.roots {
				if strings.Contains(err.Error(), root) {
					t.Fatalf("error leaks project root: %v", err)
				}
			}
			if after := snapshotTree(t, f.roots); !reflect.DeepEqual(after, before) {
				t.Fatalf("files changed:\nbefore %v\nafter  %v", before, after)
			}
		})
	}
}

func TestBacklogStatusValuesAreValidatedBeforeProjectLookup(t *testing.T) {
	resolved := false
	command := newBacklogCommandWithDirectory(func(string) (backlogPaths, error) {
		resolved = true
		return backlogPaths{}, &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: "project ghost is not registered"}
	})
	command.SilenceUsage, command.SilenceErrors = true, true
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"status", "--project", "ghost", "20261104-47530e79"}, "--value must be open, ready or done"},
		{[]string{"status", "--project", "ghost", "--value", "true", "20261104-47530e79"}, "--value must be open, ready or done"},
		{[]string{"status", "--project", "ghost", "--value", "", "20261104-47530e79"}, "--value must be open, ready or done"},
		{[]string{"list", "--project", "ghost", "--status", "closed"}, "--status must be open, ready or done"},
		{[]string{"list", "--project", "ghost", "--status", ""}, "--status must be open, ready or done"},
	} {
		command.SetArgs(test.args)
		if err := command.Execute(); err == nil || err.Error() != backlog.ClassInvalidArguments+": "+test.want {
			t.Fatalf("backlog %v: err = %v", test.args, err)
		}
	}
	if resolved {
		t.Fatal("project resolved before status was validated")
	}
}

func TestBacklogFlagErrorsDoNotRepeatFlagInput(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	for _, args := range [][]string{
		{"list", "--project", "ply", "--secret=/abs/hemmelig-tekst"},
		{"add", "--project", "ply", "--title"},
		{"list", "--project", "ply", "-/abs/hemmelig-tekst"},
	} {
		if _, err := f.run("", args...); err == nil || err.Error() != backlog.ClassInvalidArguments+": invalid flag" {
			t.Fatalf("backlog %v: err = %v", args, err)
		}
	}
}

func TestBacklogDirectoryUsesProjectDataPathAndHidesWorkspacePaths(t *testing.T) {
	wrapper := t.TempDir()
	show := func(id workspace.ProjectID) (workspace.ProjectResult, error) {
		switch id {
		case "ply":
			return workspace.ProjectResult{Workspace: "/ws", Project: workspace.ProjectRecord{ID: "ply", Wrapper: wrapper}}, nil
		case "configured":
			return workspace.ProjectResult{Workspace: "/ws", Project: workspace.ProjectRecord{ID: "configured", Wrapper: wrapper, DataPaths: map[string]string{workspace.DataPathBacklog: "workflow/backlog"}}}, nil
		case "broken":
			return workspace.ProjectResult{}, &workspace.ProjectError{Class: workspace.ErrorWorkspaceNotFound, Detail: "no Ply workspace contains /secret/cwd"}
		}
		return workspace.ProjectResult{}, &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: "project " + string(id) + " is not registered in Ply workspace /ws"}
	}
	resolve := backlogDirectory(show)
	if paths, err := resolve("ply"); err != nil || paths.directory != filepath.Join(wrapper, "backlog") ||
		paths.archive != filepath.Join(wrapper, "backlog-archive") || paths.wrapper != wrapper {
		t.Fatalf("resolve ply = %+v, %v", paths, err)
	}
	if paths, err := resolve("configured"); err != nil || paths.directory != filepath.Join(wrapper, "workflow", "backlog") {
		t.Fatalf("resolve configured = %+v, %v", paths, err)
	}
	if _, err := resolve("ghost"); err == nil || err.Error() != "workspace_project_not_found: project ghost is not registered" {
		t.Fatalf("resolve ghost = %v", err)
	}
	if _, err := resolve("broken"); err == nil || err.Error() != "workspace_not_found: could not resolve project broken" {
		t.Fatalf("resolve broken = %v", err)
	}
	if _, err := resolve("Not Valid"); err == nil || !strings.HasPrefix(err.Error(), backlog.ClassInvalidArguments+": ") {
		t.Fatalf("resolve invalid id = %v", err)
	}
}

// writeBacklogItem places an existing item file the way an earlier Ply wrote it.
func writeBacklogItem(t *testing.T, directory, id, contents string) string {
	t.Helper()
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, id+".json")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestBacklogCommandsUseConfiguredDataPath(t *testing.T) {
	const existingID = "20261101-0a0b0c0d"
	const existing = `{
  "id": "20261101-0a0b0c0d",
  "title": "Eksisterende",
  "text": "Tekst æøå\n",
  "created_at": "2026-11-01T10:00:00Z",
  "status": "ready",
  "notes": [
    {
      "at": "2026-11-02T10:00:00Z",
      "text": "Notat"
    }
  ]
}
`
	for _, relative := range []string{"", "backlog", "workflow/backlog", "data/backlog"} {
		t.Run("data-path="+relative, func(t *testing.T) {
			f := newBacklogFixture(t, "ply", "empty")
			effective := relative
			if relative == "" {
				effective = "backlog"
			} else {
				f.dataPaths["ply"] = map[string]string{workspace.DataPathBacklog: relative}
				f.dataPaths["empty"] = map[string]string{workspace.DataPathBacklog: relative}
			}
			directory := filepath.Join(f.roots["ply"], filepath.FromSlash(effective))
			path := writeBacklogItem(t, directory, existingID, existing)

			if out := f.success(t, "", "show", "--project", "ply", existingID, "--format", "json"); out != existing {
				t.Fatalf("show = %s", out)
			}
			list := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--format", "json"))
			want := []backlog.Summary{{ID: existingID, Title: "Eksisterende", CreatedAt: "2026-11-01T10:00:00Z", UpdatedAt: "2026-11-02T10:00:00Z", NoteCount: 1, Status: backlog.StatusReady}}
			if !reflect.DeepEqual(list.Items, want) || len(list.Invalid) != 0 {
				t.Fatalf("list = %#v", list)
			}
			if file, _ := os.ReadFile(path); string(file) != existing {
				t.Fatalf("read changed item file: %s", file)
			}

			noted := decodeBacklogItem(t, f.success(t, "Nytt", "note", "--project", "ply", existingID, "--format", "json"))
			if noted.ID != existingID || noted.Text != "Tekst æøå\n" || noted.Status != backlog.StatusReady || len(noted.Notes) != 2 || noted.Notes[0].Text != "Notat" || noted.Notes[1].Text != "Nytt" {
				t.Fatalf("noted = %#v", noted)
			}
			done := decodeBacklogItem(t, f.success(t, "", "status", "--project", "ply", "--value", "done", existingID, "--format", "json"))
			if done.Status != backlog.StatusDone || len(done.Notes) != 2 {
				t.Fatalf("done = %#v", done)
			}
			if file, _ := os.ReadFile(path); !strings.Contains(string(file), "\"status\": \"done\"") || strings.Contains(string(file), "ready") {
				t.Fatalf("status file = %s", file)
			}

			emptyRoot := f.roots["empty"]
			if out := f.success(t, "", "list", "--project", "empty"); out != "No backlog items in project empty.\n" {
				t.Fatalf("empty list = %q", out)
			}
			if entries, err := os.ReadDir(emptyRoot); err != nil || len(entries) != 0 {
				t.Fatalf("list created %v, %v", entries, err)
			}
			added := decodeBacklogItem(t, f.success(t, "x", "add", "--project", "empty", "--title", "Ny", "--format", "json"))
			if _, err := os.Stat(filepath.Join(emptyRoot, filepath.FromSlash(effective), added.ID+".json")); err != nil {
				t.Fatalf("add did not create the configured directory: %v", err)
			}
			if effective != "backlog" {
				for _, root := range []string{emptyRoot, f.roots["ply"]} {
					if _, err := os.Stat(filepath.Join(root, "backlog")); !errors.Is(err, fs.ErrNotExist) {
						t.Fatalf("default backlog directory used: %v", err)
					}
				}
			}
		})
	}
}

func TestBacklogConfiguredDataPathKeepsErrors(t *testing.T) {
	f := newBacklogFixture(t, "ply")
	f.dataPaths["ply"] = map[string]string{workspace.DataPathBacklog: "workflow/backlog"}
	directory := filepath.Join(f.roots["ply"], "workflow", "backlog")
	writeBacklogItem(t, directory, "20261103-0badf00d", "{")

	if _, err := f.run("", "list", "--project", "ghost"); err == nil || err.Error() != "workspace_project_not_found: project ghost is not registered" {
		t.Fatalf("unknown project = %v", err)
	}
	if _, err := f.run("", "show", "--project", "ply", "not-an-id"); err == nil || err.Error() != backlog.ClassInvalidArguments+": invalid item id" {
		t.Fatalf("invalid id = %v", err)
	}
	if _, err := f.run("", "show", "--project", "ply", "20261103-0badf00d"); err == nil || !strings.HasPrefix(err.Error(), backlog.ClassItemInvalid+": ") {
		t.Fatalf("unreadable item = %v", err)
	}
	list := decodeBacklogList(t, f.success(t, "", "list", "--project", "ply", "--format", "json"))
	if len(list.Invalid) != 1 || list.Invalid[0].File != "20261103-0badf00d.json" {
		t.Fatalf("list = %#v", list)
	}
}

func TestBacklogRejectsDataPathOutsideWrapperBeforeAccess(t *testing.T) {
	outside := t.TempDir()
	for name, relative := range map[string]string{"absolute": outside, "parent": "../backlog", "symlink": "escape/backlog"} {
		t.Run(name, func(t *testing.T) {
			f := newBacklogFixture(t, "ply")
			if err := os.Symlink(outside, filepath.Join(f.roots["ply"], "escape")); err != nil {
				t.Fatal(err)
			}
			f.dataPaths["ply"] = map[string]string{workspace.DataPathBacklog: relative}
			for _, args := range [][]string{
				{"list", "--project", "ply"},
				{"add", "--project", "ply", "--title", "t"},
				{"note", "--project", "ply", "20261104-47530e79"},
				{"status", "--project", "ply", "--value", "done", "20261104-47530e79"},
			} {
				_, err := f.run("x", args...)
				if err == nil || !strings.HasPrefix(err.Error(), string(workspace.ErrorProjectPath)+": ") || strings.Contains(err.Error(), outside) || strings.Contains(err.Error(), f.roots["ply"]) {
					t.Fatalf("backlog %v: err = %v", args, err)
				}
			}
			if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
				t.Fatalf("outside touched: %v, %v", entries, err)
			}
		})
	}
}
