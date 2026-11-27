package taskwork

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func durableFixture(t *testing.T) (string, string, Dependencies) {
	t.Helper()
	root := physicalPath(t, t.TempDir())
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "init", "-b", "main")
	runLocalGit(t, repo, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, repo, "config", "user.name", "Ply tests")
	if err := os.WriteFile(filepath.Join(repo, "README"), []byte("fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "add", "README")
	runLocalGit(t, repo, "commit", "-m", "fixture")
	config := `{"format_version":1,"projects":[{"id":"sample","name":"Sample","wrapper":".","repos":[{"id":"code","path":"repo"}]}]}`
	if err := os.WriteFile(filepath.Join(root, "ply-workspace.json"), []byte(config), 0644); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: repo}
	oid := runLocalGit(t, repo, "rev-parse", "HEAD")
	if _, err := AdoptEpic(d, EpicAdoptInput{EpicID: "feature", Title: "Feature", ProjectID: "sample", RepoID: "code", Worktree: repo, Ref: "refs/heads/main", ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	return root, repo, d
}

func TestDurableConfigGuards(t *testing.T) {
	t.Parallel()
	for _, test := range []string{"duplicate-git", "duplicate-json", "unknown-field", "symlink"} {
		t.Run(test, func(t *testing.T) {
			root, _, d := durableFixture(t)
			path := filepath.Join(root, "ply-workspace.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			switch test {
			case "duplicate-git":
				raw = []byte(strings.Replace(string(raw), `{"id":"code","path":"repo"}`, `{"id":"code","path":"repo"},{"id":"alias","path":"repo"}`, 1))
			case "duplicate-json":
				raw = []byte(strings.Replace(string(raw), `"format_version":1`, `"format_version":1,"format_version":1`, 1))
			case "unknown-field":
				raw = []byte(strings.Replace(string(raw), `"format_version":1`, `"format_version":1,"old_oid":"stale"`, 1))
			case "symlink":
				if err = os.Rename(path, path+".source"); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(path+".source", path)
			}
			if err != nil {
				t.Fatal(err)
			}
			if test != "symlink" {
				if err = os.WriteFile(path, raw, 0644); err != nil {
					t.Fatal(err)
				}
			}
			// The list reads no Git, so a Git identity conflict shows on the
			// Project's own read.
			read := func() error { _, err := workspace.ListProjects(d.Dependencies); return err }
			if test == "duplicate-git" {
				read = func() error { _, err := workspace.ShowProject(d.Dependencies, "sample"); return err }
			}
			if err = read(); err == nil {
				t.Fatal("invalid definitions accepted")
			}
		})
	}
}

func TestDurableConfigAuthorityAndStableReadback(t *testing.T) {
	t.Parallel()
	root, _, d := durableFixture(t)
	first, err := ListEpics(d)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	again, err := ListEpics(d)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(testWork(t, root).Registry())
	if err != nil {
		t.Fatal(err)
	}
	if first.Epics[0].RepoBindings[0].Worktree.ID != again.Epics[0].RepoBindings[0].Worktree.ID || string(before) != string(after) {
		t.Fatal("ordinary read rewrote existing runtime binding")
	}
	if got, err := workspace.ListProjects(d.Dependencies); err != nil || len(got.Projects) != 1 || got.Projects[0].ID != "sample" {
		t.Fatalf("config not authoritative: %#v, %v", got, err)
	}
	input := workspace.ProjectAddInput{ProjectID: "extra", Name: "Extra", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "code", Path: filepath.Join(root, "repo")}}}
	if _, err = workspace.AddProject(d.Dependencies, input); err == nil {
		t.Fatal("Project add overrode config")
	}
	basis, err := ReadViewBasis(d)
	if err != nil {
		t.Fatal(err)
	}
	if basis.Sources[filepath.Join(root, "ply-workspace.json")] == "" {
		t.Fatal("view omitted authoritative source")
	}
	path := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, append(raw, ' '), 0644); err != nil {
		t.Fatal(err)
	}
	if state, _ := CheckReadViewBasis(d, basis); state == "fresh" {
		t.Fatal("config change did not stale view")
	}
}

type definitionRegistrationReader struct{ WorkItemStore }

func (s definitionRegistrationReader) Snapshot(workspace.ProjectWork) (WorkItemRegistry, error) {
	return WorkItemRegistry{}, fmt.Errorf("full content read is unavailable")
}
func TestDurableConfigViewsKeepRegistrationOnlyRead(t *testing.T) {
	t.Parallel()
	_, _, d := durableFixture(t)
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	d.WorkItems = definitionRegistrationReader{d.WorkItems}
	if _, err := ReadViewBasis(d); err != nil {
		t.Fatalf("view needed unrelated full content: %v", err)
	}
}
