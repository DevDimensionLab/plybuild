package workspaceview

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func treeEntries(lines ...string) []taskwork.CommitTreeEntry {
	out := []taskwork.CommitTreeEntry{}
	for _, line := range lines {
		mode, path, _ := strings.Cut(line, " ")
		out = append(out, taskwork.CommitTreeEntry{Mode: mode, Path: path})
	}
	return out
}

func moduleParents(modules []Module) map[string]string {
	out := map[string]string{}
	for _, m := range modules {
		parent := "<null>"
		if m.Parent != nil {
			parent = *m.Parent
		}
		out[m.Module] = parent
	}
	return out
}

func TestModuleTreeFiveLevelsWithAGap(t *testing.T) {
	got := ModuleTree(treeEntries(
		"100644 pom.xml",
		"100644 dom/pom.xml",
		"100644 dom/wrap/pom.xml",
		"100644 dom/wrap/mid/pom.xml",
		"100755 dom/wrap/mid/app/pom.xml",
		"100644 dom/wrap/gap/leaf/pom.xml",
		"100644 dom/wrap/gap/README.md",
		"100644 modules/a/pom.xml",
		"100644 modules/a/src/main/java/A.java",
	))
	want := map[string]string{".": "<null>", "dom": ".", "dom/wrap": "dom", "dom/wrap/mid": "dom/wrap", "dom/wrap/mid/app": "dom/wrap/mid", "dom/wrap/gap/leaf": "dom/wrap", "modules/a": "."}
	if !reflect.DeepEqual(moduleParents(got), want) {
		t.Fatalf("tree %v", moduleParents(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Module >= got[i].Module {
			t.Fatalf("not sorted: %v", got)
		}
	}
}

func TestModuleTreeWithoutRootPom(t *testing.T) {
	got := moduleParents(ModuleTree(treeEntries("100644 a/pom.xml", "100644 a/b/pom.xml", "100644 c/pom.xml", "100644 README.md")))
	if want := map[string]string{"a": "<null>", "a/b": "a", "c": "<null>"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tree %v", got)
	}
	if got := ModuleTree(treeEntries("100644 README.md")); got == nil || len(got) != 0 {
		t.Fatalf("no pom.xml must give an empty, non-null list: %#v", got)
	}
}

func TestModuleTreeExcludesLikePlycq(t *testing.T) {
	got := moduleParents(ModuleTree(treeEntries(
		"100644 pom.xml",
		"100644 .hidden/pom.xml",
		"100644 a/.mvn/pom.xml",
		"100644 vendor/pom.xml",
		"100644 x/node_modules/pom.xml",
		"100644 testdata/pom.xml",
		"100644 target/pom.xml",
		"100644 x/build/pom.xml",
		"100644 x/target/y/pom.xml",
		"100644 src/main/build/pom.xml",
		"120000 link/pom.xml",
		"160000 sub/pom.xml",
		"100644 y/pom.xml.bak",
		"100644 y/Pom.xml",
	)))
	if want := map[string]string{".": "<null>", "src/main/build": "."}; !reflect.DeepEqual(got, want) {
		t.Fatalf("tree %v", got)
	}
}

func moduleRepo(t *testing.T) (taskwork.Dependencies, string, string) {
	t.Helper()
	d, root := changesFixture(t)
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, repo, "init", "-b", "main")
	return d, root, repo
}

func writeFixtureFiles(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, name := range names {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("<project/>\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func listModulesJSON(t *testing.T, d taskwork.Dependencies, o ModuleOptions) ModuleListResult {
	t.Helper()
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ListModules(d, s, o)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	var back ModuleListResult
	if err := json.Unmarshal(b, &back); err != nil || back.Kind != "WorkspaceModuleListReadback@1" || back.SchemaVersion != 1 {
		t.Fatalf("readback %s %v", b, err)
	}
	return out
}

// The tree is read from a commit, never from the work tree.
func TestListModulesReadsACommit(t *testing.T) {
	d, root, repo := moduleRepo(t)
	writeFixtureFiles(t, repo, "pom.xml", "a/pom.xml")
	if err := os.Symlink("../pom.xml", filepath.Join(repo, "a", "linked")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(repo, "l"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../pom.xml", filepath.Join(repo, "l", "pom.xml")); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, repo, "add", ".")
	inventoryGit(t, repo, "update-index", "--add", "--cacheinfo", "160000,"+strings.Repeat("1", 40)+",g/pom.xml")
	inventoryGit(t, repo, "commit", "-m", "first")
	first := inventoryGit(t, repo, "rev-parse", "HEAD")
	writeFixtureFiles(t, repo, "a/b/pom.xml")
	inventoryGit(t, repo, "add", ".")
	inventoryGit(t, repo, "commit", "-m", "second")
	head := inventoryGit(t, repo, "rev-parse", "HEAD")
	writeFixtureFiles(t, repo, "uncommitted/pom.xml")
	if _, err := workspace.AddProject(d.Dependencies, workspace.ProjectAddInput{ProjectID: "project", Name: "Modules", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); err != nil {
		t.Fatal(err)
	}

	out := listModulesJSON(t, d, ModuleOptions{ProjectID: "project"})
	if out.Project != "project" || len(out.Repos) != 1 {
		t.Fatalf("result %+v", out)
	}
	r := out.Repos[0]
	if r.Repo != "repo" || r.Reason != nil || r.Commit == nil || *r.Commit != head {
		t.Fatalf("row %+v", r)
	}
	if got, want := moduleParents(r.Modules), map[string]string{".": "<null>", "a": ".", "a/b": "a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("modules %v", got)
	}

	r = listModulesJSON(t, d, ModuleOptions{ProjectID: "project", RepoID: "repo", Commit: first[:10]}).Repos[0]
	if r.Commit == nil || *r.Commit != first || !reflect.DeepEqual(moduleParents(r.Modules), map[string]string{".": "<null>", "a": "."}) {
		t.Fatalf("older commit %+v", r)
	}

	for _, commit := range []string{"no-such-ref", strings.Repeat("e", 40), "--output=x"} {
		r = listModulesJSON(t, d, ModuleOptions{ProjectID: "project", RepoID: "repo", Commit: commit}).Repos[0]
		if r.Reason == nil || *r.Reason != "commit_unknown" || r.Commit != nil || r.Modules != nil {
			t.Fatalf("unknown commit %q: %+v", commit, r)
		}
	}
}

func TestListModulesKeepsEveryRepoRow(t *testing.T) {
	d, root, repo := moduleRepo(t)
	empty := filepath.Join(root, "empty")
	if err := os.Mkdir(empty, 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, empty, "init", "-b", "main")
	gone := filepath.Join(root, "gone")
	if err := os.Mkdir(gone, 0700); err != nil {
		t.Fatal(err)
	}
	inventoryGit(t, gone, "init", "-b", "main")
	writeFixtureFiles(t, repo, "README.md")
	inventoryGit(t, repo, "add", ".")
	inventoryGit(t, repo, "commit", "-m", "no pom")
	if _, err := workspace.AddProject(d.Dependencies, workspace.ProjectAddInput{ProjectID: "project", Name: "Modules", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}, {RepoID: "empty", Path: empty}, {RepoID: "gone", Path: gone}}}); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(gone); err != nil {
		t.Fatal(err)
	}
	out := listModulesJSON(t, d, ModuleOptions{ProjectID: "project"})
	got := []string{}
	for _, r := range out.Repos {
		reason := "<null>"
		if r.Reason != nil {
			reason = *r.Reason
		}
		got = append(got, fmt.Sprintf("%s:%s:%v", r.Repo, reason, r.Modules == nil))
	}
	if want := []string{"empty:no_commit:true", "gone:repo_unavailable:true", "repo:<null>:false"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rows %v", got)
	}
	if len(out.Repos[2].Modules) != 0 {
		t.Fatalf("modules %v", out.Repos[2].Modules)
	}
}

func TestListModulesRefusesUnknownProjectAndRepo(t *testing.T) {
	d, root, repo := moduleRepo(t)
	writeFixtureFiles(t, repo, "pom.xml")
	inventoryGit(t, repo, "add", ".")
	inventoryGit(t, repo, "commit", "-m", "first")
	if _, err := workspace.AddProject(d.Dependencies, workspace.ProjectAddInput{ProjectID: "project", Name: "Modules", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range []ModuleOptions{{ProjectID: "other"}, {ProjectID: "project", RepoID: "other"}, {}, {ProjectID: "project", Commit: "HEAD"}} {
		if _, err := ListModules(d, s, o); err == nil {
			t.Fatalf("%+v was accepted", o)
		}
	}
}
