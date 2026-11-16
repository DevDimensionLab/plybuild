package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/publication"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func runPublication(t *testing.T, cwd string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(cwd)
	command := newPublicationCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

// producerDir is a finished producer result with a nested file.
func producerDir(t *testing.T, root, name string, files map[string]string) string {
	t.Helper()
	dir := filepath.Join(root, "out", name)
	for file, contents := range files {
		path := filepath.Join(dir, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func publishJSON(t *testing.T, root string, args ...string) publication.View {
	t.Helper()
	out, err := runPublication(t, root, append([]string{"publish", "--project", "a", "--format", "json"}, args...)...)
	if err != nil {
		t.Fatalf("publish %v: %v", args, err)
	}
	var v publication.View
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("publish json %q: %v", out, err)
	}
	return v
}

func listJSON(t *testing.T, cwd string, args ...string) publication.ListResult {
	t.Helper()
	out, err := runPublication(t, cwd, append([]string{"list", "--format", "json"}, args...)...)
	if err != nil {
		t.Fatalf("list %v: %v", args, err)
	}
	var result struct {
		Kind string `json:"kind"`
		publication.ListResult
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil || result.Kind != "PlyPublicationList@1" {
		t.Fatalf("list json %q: %v", out, err)
	}
	return result.ListResult
}

func requirePublicationClass(t *testing.T, err error, class string) {
	t.Helper()
	var e *publication.Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error %v, want %s", err, class)
	}
}

func addTestPlan(t *testing.T, root string) {
	t.Helper()
	if _, err := runPlan(t, root, planV1Snapshot, "add", "--project", "a", "--actor", "ply-epic", "--reason", "test"); err != nil {
		t.Fatal(err)
	}
}

func TestPublicationPublishListShowRoundTrip(t *testing.T) {
	root := newAssignmentWorkspace(t)
	addTestPlan(t, root)
	files := map[string]string{"index.html": "<html>kart</html>", "data/map.json": `{"a":1}`}
	dir := producerDir(t, root, "map", files)
	v := publishJSON(t, root, "--repo", "a-app", "--category", "code-map", "--dir", dir, "--entry", "index.html",
		"--title", "Kodekart", "--source", "d585196", "--plan", "PLAN-1", "--actor", "ply-epic")
	if v.Kind != "PlyPublication@1" || v.Project != "a" || v.Repo != "a-app" || v.Category != "code-map" || v.Plan != "PLAN-1" ||
		v.Title != "Kodekart" || v.Entry != "index.html" || v.Source != "d585196" || v.Actor != "ply-epic" || v.Agent != "" {
		t.Fatalf("metadata %+v", v)
	}
	if want := filepath.Join(root, "a", "data", "publications", "a-app", "code-map", v.Run); v.Path != want {
		t.Fatalf("path %q, want %q", v.Path, want)
	}
	if len(v.Files) != 2 || v.Files[0].Name != "data/map.json" || v.Files[1].Name != "index.html" || v.Files[1].Size != int64(len(files["index.html"])) || len(v.Files[1].SHA256) != 64 {
		t.Fatalf("files %+v", v.Files)
	}
	for name, contents := range files {
		got, err := os.ReadFile(filepath.Join(v.Path, filepath.FromSlash(name)))
		if err != nil || string(got) != contents {
			t.Fatalf("%s: %q %v", name, got, err)
		}
	}
	list := listJSON(t, root, "--project", "a")
	if len(list.Items) != 1 || list.Items[0].Run != v.Run || list.Items[0].Path != v.Path || len(list.Invalid) != 0 {
		t.Fatalf("list %+v", list)
	}
	out, err := runPublication(t, root, "show", "--project", "a", "--repo", "a-app", "--category", "code-map", "--format", "json")
	var shown publication.View
	if err != nil || json.Unmarshal([]byte(out), &shown) != nil || shown.Run != v.Run || shown.Path != v.Path || shown.Files[0].SHA256 != v.Files[0].SHA256 {
		t.Fatalf("show %q %v", out, err)
	}
	if text, err := runPublication(t, root, "show", "--project", "a", "--repo", "a-app", "--category", "code-map", "--run", v.Run); err != nil || !strings.Contains(text, v.Path) {
		t.Fatalf("show text %q %v", text, err)
	}
}

func TestPublicationKeepsProjectRepoAndAgentPathsApart(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "flow", map[string]string{"index.html": "flow"})
	repo := publishJSON(t, root, "--repo", "a-app", "--category", "workflow", "--dir", dir, "--entry", "index.html", "--title", "Repo")
	project := publishJSON(t, root, "--category", "workflow", "--dir", dir, "--entry", "index.html", "--title", "Prosjekt")
	agent := publishJSON(t, root, "--category", "workflow", "--agent", "ply-epic", "--dir", dir, "--entry", "index.html", "--title", "Agent")
	publications := filepath.Join(root, "a", "data", "publications")
	if project.Repo != "" || filepath.Dir(project.Path) != filepath.Join(publications, "_project", "workflow") ||
		filepath.Dir(agent.Path) != filepath.Join(publications, "_project", "workflow", "ply-epic") ||
		filepath.Dir(repo.Path) != filepath.Join(publications, "a-app", "workflow") {
		t.Fatalf("paths %s %s %s", repo.Path, project.Path, agent.Path)
	}
	if list := listJSON(t, root, "--project", "a", "--category", "workflow"); len(list.Items) != 3 {
		t.Fatalf("latest per path %+v", list.Items)
	}
	if list := listJSON(t, root, "--project", "a", "--agent", "ply-epic"); len(list.Items) != 1 || list.Items[0].Run != agent.Run {
		t.Fatalf("agent filter %+v", list.Items)
	}
	out, err := runPublication(t, root, "show", "--project", "a", "--category", "workflow", "--format", "json")
	if err != nil || !strings.Contains(out, project.Run) {
		t.Fatalf("project show %q %v", out, err)
	}
}

func TestPublicationListShowsLatestAndAllRuns(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	first := publishJSON(t, root, "--repo", "a-app", "--category", "code-map", "--dir", dir, "--entry", "index.html", "--title", "1")
	second := publishJSON(t, root, "--repo", "a-app", "--category", "code-map", "--dir", dir, "--entry", "index.html", "--title", "2")
	if first.Run >= second.Run {
		t.Fatalf("runs out of order: %s %s", first.Run, second.Run)
	}
	if list := listJSON(t, root, "--project", "a"); len(list.Items) != 1 || list.Items[0].Run != second.Run {
		t.Fatalf("latest %+v", list.Items)
	}
	if list := listJSON(t, root, "--project", "a", "--all"); len(list.Items) != 2 || list.Items[0].Run != second.Run || list.Items[1].Run != first.Run {
		t.Fatalf("all %+v", list.Items)
	}
}

func TestPublicationRetentionKeepsPlanRuns(t *testing.T) {
	root := newAssignmentWorkspace(t)
	addTestPlan(t, root)
	dir := producerDir(t, root, "flow", map[string]string{"index.html": "x"})
	args := []string{"--category", "workflow", "--agent", "ply-epic", "--dir", dir, "--entry", "index.html", "--title", "t"}
	planned := publishJSON(t, root, append(args, "--plan", "PLAN-1")...)
	var runs []string
	for i := 0; i < 6; i++ {
		runs = append(runs, publishJSON(t, root, args...).Run)
	}
	all := listJSON(t, root, "--project", "a", "--all")
	if len(all.Items) != 6 {
		t.Fatalf("want 5 runs without plan and 1 with, got %+v", all.Items)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(planned.Path), runs[0])); !os.IsNotExist(err) {
		t.Fatalf("oldest run without plan was kept: %v", err)
	}
	if _, err := os.Stat(planned.Path); err != nil {
		t.Fatalf("plan run removed: %v", err)
	}
	if list := listJSON(t, root, "--project", "a", "--plan", "PLAN-1"); len(list.Items) != 1 || list.Items[0].Run != planned.Run {
		t.Fatalf("plan filter %+v", list.Items)
	}
}

func TestPublicationRefusesWithoutStoring(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	linked := producerDir(t, root, "linked", map[string]string{"index.html": "x"})
	if err := os.Symlink(filepath.Join(dir, "index.html"), filepath.Join(linked, "other.html")); err != nil {
		t.Fatal(err)
	}
	large := producerDir(t, root, "large", map[string]string{"index.html": "x"})
	if err := os.Truncate(filepath.Join(large, "index.html"), 20<<20+1); err != nil {
		t.Fatal(err)
	}
	reserved := producerDir(t, root, "reserved", map[string]string{"index.html": "x", "publication.json": "{}"})
	base := func(d string) []string {
		return []string{"publish", "--project", "a", "--category", "code-map", "--dir", d, "--entry", "index.html", "--title", "t"}
	}
	for name, c := range map[string]struct {
		args  []string
		class string
	}{
		"unknown repo":     {append(base(dir), "--repo", "b-app"), publication.ClassInvalidArguments},
		"unknown plan":     {append(base(dir), "--plan", "PLAN-9"), publication.ClassInvalidArguments},
		"invalid category": {[]string{"publish", "--project", "a", "--category", "Code_Map", "--dir", dir, "--entry", "index.html", "--title", "t"}, publication.ClassInvalidArguments},
		"invalid agent":    {append(base(dir), "--agent", "Ply Epic"), publication.ClassInvalidArguments},
		"missing entry":    {[]string{"publish", "--project", "a", "--category", "code-map", "--dir", dir, "--entry", "map.html", "--title", "t"}, publication.ClassInvalidArguments},
		"symlink":          {base(linked), publication.ClassInvalidArguments},
		"too large":        {base(large), publication.ClassTooLarge},
		"reserved name":    {base(reserved), publication.ClassInvalidArguments},
		"two line title":   {[]string{"publish", "--project", "a", "--category", "code-map", "--dir", dir, "--entry", "index.html", "--title", "a\nb"}, publication.ClassInvalidArguments},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := runPublication(t, root, c.args...)
			requirePublicationClass(t, err, c.class)
		})
	}
	if list := listJSON(t, root, "--project", "a", "--all"); len(list.Items) != 0 || len(list.Invalid) != 0 {
		t.Fatalf("refusal stored something: %+v", list)
	}
}

func TestPublicationHidesUnfinishedRunsAndOtherProjects(t *testing.T) {
	root := newAssignmentWorkspace(t)
	dir := producerDir(t, root, "map", map[string]string{"index.html": "x"})
	v := publishJSON(t, root, "--repo", "a-app", "--category", "code-map", "--dir", dir, "--entry", "index.html", "--title", "t")
	unfinished := filepath.Join(filepath.Dir(v.Path), ".20991231T000000000Z-00000000.tmp")
	if err := os.MkdirAll(unfinished, 0o755); err != nil {
		t.Fatal(err)
	}
	if list := listJSON(t, root, "--project", "a", "--all"); len(list.Items) != 1 || list.Items[0].Run != v.Run || len(list.Invalid) != 0 {
		t.Fatalf("unfinished run visible: %+v", list)
	}
	if list := listJSON(t, root, "--project", "b", "--all"); len(list.Items) != 0 {
		t.Fatalf("project b sees a: %+v", list.Items)
	}
	_, err := runPublication(t, root, "show", "--project", "b", "--category", "code-map")
	requirePublicationClass(t, err, publication.ClassNotFound)
}
