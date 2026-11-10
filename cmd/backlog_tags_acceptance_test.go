//go:build unix

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// The tags journey runs the actual CLI against an isolated workspace whose
// project has two repos and a configured data_paths.backlog.
func TestBacklogAuthorTagsAndFiltersAcceptance(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "ply")
	if out, err := exec.Command("go", "build", "-o", binary, "./ply").CombinedOutput(); err != nil {
		t.Fatalf("build actual CLI: %v\n%s", err, out)
	}
	f := newPlanCLI(t, binary)
	ws := filepath.Join(f.root, "ws")
	for _, repo := range []string{"code", "web"} {
		dir := filepath.Join(ws, repo)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		f.git(t, dir, "init", "-b", "main")
		f.git(t, dir, "commit", "--allow-empty", "-m", "fixture")
	}
	f.success(t, ws, "workspace", "init")
	mustPlanWrite(t, filepath.Join(ws, "ply-workspace.json"), `{"format_version":1,"projects":[{"id":"ply","name":"Ply","wrapper":".","data_paths":{"backlog":"workflow/backlog"},"repos":[{"id":"code","path":"code"},{"id":"web","path":"web"}],"epics":[]}]}`, 0o644)
	backlogDir := filepath.Join(ws, "workflow", "backlog")
	run := func(stdin string, args ...string) string {
		c := exec.Command(binary, append([]string{"backlog", "--workspace", ws}, append(args, "--project", "ply")...)...)
		c.Dir, c.Env, c.Stdin = f.root, f.environment(nil), strings.NewReader(stdin)
		out, err := c.Output()
		if err != nil {
			t.Fatalf("backlog %v: %v", args, err)
		}
		return string(out)
	}
	fail := func(args ...string) {
		f.failure(t, f.root, nil, "backlog_invalid_arguments", append([]string{"backlog", "--workspace", ws}, append(args, "--project", "ply")...)...)
	}

	first := decodeBacklogItem(t, run("Tekst\n", "add", "--title", "Første", "--format", "json"))
	tagged := decodeBacklogItem(t, run("", "tags", first.ID, "--repo", "web", "--repo", "code", "--topic", " Web app ", "--topic", "web APP", "--topic", "API", "--format", "json"))
	if !reflect.DeepEqual(tagged.Repos, []string{"web", "code"}) || !reflect.DeepEqual(tagged.Topics, []string{"Web app", "API"}) || tagged.Text != "Tekst\n" {
		t.Fatalf("tagged = %#v", tagged)
	}
	stored, _ := os.ReadFile(filepath.Join(backlogDir, first.ID+".json"))
	for _, bad := range [][]string{{"--repo", "ghost"}, {"--repo", "code", "--repo", "ghost"}, {"--clear", "--repo", "code"}, {}} {
		fail(append([]string{"tags", first.ID}, bad...)...)
	}
	if after, _ := os.ReadFile(filepath.Join(backlogDir, first.ID+".json")); string(after) != string(stored) {
		t.Fatal("rejected tags changed the item")
	}
	noted := decodeBacklogItem(t, run("Notat\n", "note", first.ID, "--author", " Per ", "--format", "json"))
	if len(noted.Notes) != 1 || noted.Notes[0].Author != "Per" || noted.Notes[0].Text != "Notat\n" {
		t.Fatalf("noted = %#v", noted.Notes)
	}
	show := run("", "show", first.ID)
	if !strings.Contains(show, "Status: open\nRepos: web, code\nTopics: Web app, API\n\nTekst\n") || !strings.Contains(show, " by Per:\nNotat\n") {
		t.Fatalf("show = %q", show)
	}

	second := decodeBacklogItem(t, run("To\n", "add", "--title", "Andre", "--format", "json"))
	run("", "tags", second.ID, "--repo", "code", "--topic", "api")
	writeBacklogItem(t, backlogDir, "20261101-0000000a", `{"id":"20261101-0000000a","title":"Gammel","text":"t","created_at":"2026-11-01T10:00:00Z","status":"done","notes":[],"repos":["retired"]}`)
	for _, test := range []struct {
		filter []string
		want   []string
	}{
		{[]string{"--repo", "code"}, []string{first.ID, second.ID}},
		{[]string{"--repo", "code", "--topic", "WEB APP"}, []string{first.ID}},
		{[]string{"--topic", "Api", "--repo", "web"}, []string{first.ID}},
		{[]string{"--repo", "retired"}, []string{"20261101-0000000a"}},
		{[]string{"--repo", "Code"}, nil},
	} {
		var got []string
		for _, item := range decodeBacklogList(t, run("", append([]string{"list", "--format", "json"}, test.filter...)...)).Items {
			got = append(got, item.ID)
		}
		slices.Sort(got)
		if slices.Sort(test.want); !reflect.DeepEqual(got, test.want) {
			t.Fatalf("list %v = %v, want %v", test.filter, got, test.want)
		}
	}
	if list := run("", "list", "--topic", "api"); !strings.Contains(list, "Andre  repos: code  topics: api\n") {
		t.Fatalf("list text = %q", list)
	}

	cleared := run("", "tags", first.ID, "--clear", "--format", "json")
	if strings.Contains(cleared, "repos") || strings.Contains(cleared, "topics") || !strings.Contains(cleared, `"author": "Per"`) {
		t.Fatalf("cleared = %s", cleared)
	}
}
