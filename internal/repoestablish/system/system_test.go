package system

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// fakeGH answers gh in process: NOT_FOUND until create, then private and empty.
// Git runs for real.
type fakeGH struct {
	calls      [][]string
	created    bool
	visibility string
}

func (f *fakeGH) Run(c process.Command) error {
	if c.Name != "gh" {
		return process.SystemRunner().Runner.Run(c)
	}
	f.calls = append(f.calls, c.Args)
	if c.Args[0] == "repo" {
		f.created = true
		return nil
	}
	if !f.created {
		fmt.Fprint(c.Stdout, `{"data":{"repository":null},"errors":[{"type":"NOT_FOUND","path":["repository"],"message":"x"}]}`)
		return exec.Command("/bin/sh", "-c", "exit 1").Run()
	}
	owner, name := strings.TrimPrefix(c.Args[7], "owner="), strings.TrimPrefix(c.Args[9], "name=")
	visibility := f.visibility
	if visibility == "" {
		visibility = "PRIVATE"
	}
	fmt.Fprintf(c.Stdout, `{"data":{"repository":{"nameWithOwner":"%s/%s","visibility":"%s","isEmpty":true,"sshUrl":"git@github.com:%[1]s/%[2]s.git"}}}`, owner, name, visibility)
	return nil
}

type fixture struct {
	t                      *testing.T
	root, projects, marker string
	gh                     *fakeGH
	ghMarker               string
}

func physical(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func write(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

// newFixture builds a synthetic workspace whose Project root is not the
// workspace root, with one live repository, and puts a gh on PATH that
// records any execution; the runner never executes it.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	root := physical(t, t.TempDir())
	if _, err := workspace.Init(workspace.Dependencies{Files: fixedCwdFiles{workspace.SystemDependencies().Files, root}}); err != nil {
		t.Fatal(err)
	}
	f := &fixture{t: t, root: root, projects: filepath.Join(root, "projects"), gh: &fakeGH{}}
	for _, dir := range []string{"projects/ply/main", "apps/app/main"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
		git(t, filepath.Join(root, dir), "init", "-q")
	}
	write(t, filepath.Join(root, workspacedef.File), fmt.Sprintf(`{
  "format_version": 1,
  "projects": [
    {"id": "ply", "name": "Ply", "wrapper": "projects", "repos": [{"id": "ply", "path": "projects/ply/main"}], "epics": []},
    {"id": "app", "name": "App", "wrapper": %q, "repos": [{"id": "app", "path": "apps/app/main"}], "epics": []}
  ]
}
`, filepath.Join(root, "apps")))
	bin := t.TempDir()
	f.ghMarker = filepath.Join(bin, "gh-executed")
	write(t, filepath.Join(bin, "gh"), "#!/bin/sh\ntouch "+f.ghMarker+"\nexit 97\n")
	if err := os.Chmod(filepath.Join(bin, "gh"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Cleanup(func() {
		if _, err := os.Lstat(f.ghMarker); err == nil {
			t.Error("the gh executable on PATH was executed")
		}
	})
	return f
}

func (f *fixture) request() repoestablish.Request {
	return repoestablish.Request{Workspace: f.root, Project: "ply", RepoName: "code-map", Actor: "ply-epic",
		Overrides: repoestablish.Overrides{GitHubOwner: "Org"}}
}

func (f *fixture) deps(h hooks) repoestablish.Dependencies {
	return newDependencies(f.gh, h)
}

func (f *fixture) read(path string) string {
	f.t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		f.t.Fatal(err)
	}
	return string(raw)
}

func (f *fixture) metadata() string {
	return f.read(filepath.Join(f.root, workspace.MarkerDirectory, workspace.ProjectMetadataFile))
}

func (f *fixture) definitions() string { return f.read(filepath.Join(f.root, workspacedef.File)) }

func requireKind(t *testing.T, err, kind error) {
	t.Helper()
	if !errors.Is(err, kind) {
		t.Fatalf("error %v, want %v", err, kind)
	}
}

func requireStop(t *testing.T, err error, step repoestablish.Step, kind error) {
	t.Helper()
	var stopped *repoestablish.Error
	if !errors.As(err, &stopped) || stopped.Step != step || !errors.Is(err, kind) {
		t.Fatalf("error %v, want %v at %s", err, kind, step)
	}
}

func TestPrepareResolvesProjectRootAgainstWorkspace(t *testing.T) {
	f := newFixture(t)
	d := f.deps(hooks{})
	setup, err := d.Prepare(context.Background(), f.request())
	if err != nil || setup != (repoestablish.ProjectSetup{Root: f.projects}) {
		t.Fatalf("relative wrapper: %+v %v", setup, err)
	}
	req := f.request()
	req.Project = "app"
	if setup, err = d.Prepare(context.Background(), req); err != nil || setup.Root != filepath.Join(f.root, "apps") {
		t.Fatalf("absolute wrapper: %+v %v", setup, err)
	}
	req.Project = "other"
	_, err = d.Prepare(context.Background(), req)
	requireKind(t, err, repoestablish.ErrInvalid)

	if err := os.Symlink(f.projects, filepath.Join(f.root, "alias")); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(f.root, workspacedef.File), strings.Replace(f.definitions(), `"wrapper": "projects"`, `"wrapper": "alias"`, 1))
	_, err = d.Prepare(context.Background(), f.request())
	requireKind(t, err, repoestablish.ErrInvalid)
	write(t, filepath.Join(f.root, workspacedef.File), "{")
	_, err = d.Prepare(context.Background(), f.request())
	requireKind(t, err, repoestablish.ErrInvalid)
}

func TestPrepareReadsProjectDefaults(t *testing.T) {
	f := newFixture(t)
	d := f.deps(hooks{})
	path := filepath.Join(f.projects, DefaultsFile)
	write(t, path, `{"format_version": 1, "github_owner": "Org", "leaf": "trunk", "branch": "dev", "remote": "up"}`)
	setup, err := d.Prepare(context.Background(), f.request())
	if want := (repoestablish.ProjectSetup{Root: f.projects, GitHubOwner: "Org", Leaf: "trunk", Branch: "dev", Remote: "up"}); err != nil || setup != want {
		t.Fatalf("setup %+v, error %v", setup, err)
	}
	for name, contents := range map[string]string{
		"invalid JSON":      `{"format_version": 1,`,
		"unknown key":       `{"format_version": 1, "github_owner": "Org", "repo": "x"}`,
		"wrong version":     `{"format_version": 2, "github_owner": "Org"}`,
		"missing version":   `{"github_owner": "Org"}`,
		"empty owner":       `{"format_version": 1, "github_owner": ""}`,
		"empty leaf":        `{"format_version": 1, "leaf": ""}`,
		"null leaf":         `{"format_version": 1, "leaf": null}`,
		"two values":        `{"format_version": 1} {}`,
		"null":              `null`,
		"over 64 KiB":       `{"format_version": 1, "github_owner": "` + strings.Repeat("a", 64<<10) + `"}`,
		"not a JSON object": `[]`,
	} {
		t.Run(name, func(t *testing.T) {
			write(t, path, contents)
			_, err := d.Prepare(context.Background(), f.request()) // the request carries an owner override
			requireKind(t, err, repoestablish.ErrInvalid)
		})
	}
	t.Run("symlink", func(t *testing.T) {
		target := filepath.Join(f.root, "defaults.json")
		write(t, target, `{"format_version": 1, "github_owner": "Org"}`)
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		_, err := d.Prepare(context.Background(), f.request())
		requireKind(t, err, repoestablish.ErrInvalid)
	})
}

func TestRunEstablishesUnderExplicitWorkspaceAndActorOnly(t *testing.T) {
	f := newFixture(t)
	t.Chdir(t.TempDir())
	t.Setenv("PLY_ACTOR", "someone-else")
	result, err := repoestablish.Run(context.Background(), f.request(), f.deps(hooks{}))
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(f.projects, "code-map")
	if result.Resolved.Wrapper != wrapper || result.SSHURL != "git@github.com:Org/code-map.git" || len(result.Steps) != 6 {
		t.Fatalf("result %+v", result)
	}
	for _, step := range result.Steps {
		if !step.Changed {
			t.Fatalf("steps %+v", result.Steps)
		}
	}
	defs, _, err := workspacedef.Read(f.root)
	if err != nil || !reflect.DeepEqual(defs.Projects[0].Repos[1], workspacedef.Repo{ID: "code-map", Path: "projects/code-map/main"}) {
		t.Fatalf("definitions %+v, error %v", defs, err)
	}
	snapshot, err := workspace.ReadProjectMetadata(workspace.SystemDependencies(), f.root)
	if err != nil || len(snapshot.Events) != 1 || snapshot.Events[0].ActorClaim != "ply-epic" || *snapshot.RepositoryWrapper("code-map") != wrapper {
		t.Fatalf("metadata %+v, error %v", snapshot, err)
	}
	remote, err := exec.Command("git", "-C", result.Resolved.RepoDir, "remote", "get-url", "origin").Output()
	if err != nil || string(remote) != result.SSHURL+"\n" {
		t.Fatalf("remote %q, error %v", remote, err)
	}

	definitions, metadata, calls := f.definitions(), f.metadata(), len(f.gh.calls)
	result, err = repoestablish.Run(context.Background(), f.request(), f.deps(hooks{}))
	if err != nil || f.definitions() != definitions || f.metadata() != metadata {
		t.Fatalf("rerun changed files: %+v, error %v", result, err)
	}
	for _, step := range result.Steps {
		if step.Changed {
			t.Fatalf("rerun steps %+v", result.Steps)
		}
	}
	if got := f.gh.calls[calls:]; len(got) != 1 || got[0][0] != "api" {
		t.Fatalf("rerun gh calls %q", got)
	}
}

// established runs Run until the define step so that SetWrapper can be tested alone.
func established(t *testing.T, f *fixture) repoestablish.Resolved {
	t.Helper()
	d := f.deps(hooks{})
	d.SetWrapper = func(context.Context, repoestablish.Resolved) (bool, error) { return false, nil }
	result, err := repoestablish.Run(context.Background(), f.request(), d)
	if err != nil {
		t.Fatal(err)
	}
	return result.Resolved
}

func TestSetWrapperComparesThenSetsWithExpectedRevision(t *testing.T) {
	f := newFixture(t)
	r := established(t, f)
	d := f.deps(hooks{})
	if changed, err := d.SetWrapper(context.Background(), r); err != nil || !changed {
		t.Fatalf("absent: %v %v", changed, err)
	}
	once := f.metadata()
	if changed, err := d.SetWrapper(context.Background(), r); err != nil || changed || f.metadata() != once {
		t.Fatalf("equal: %v %v", changed, err)
	}
	other := r
	other.Wrapper = f.projects
	_, err := d.SetWrapper(context.Background(), other)
	if requireKind(t, err, repoestablish.ErrConflict); f.metadata() != once {
		t.Fatal("conflict changed metadata")
	}
}

func TestSetWrapperLosesRaceWithoutOverwriting(t *testing.T) {
	f := newFixture(t)
	r := established(t, f)
	var raced string
	h := hooks{afterMetadataRead: func() {
		in := workspace.ProjectMetadataInput{Operation: workspace.ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "code-map", Locator: f.projects, ActorClaim: "racer"}
		if _, err := workspace.UpdateProjectMetadata(fixedDependencies(f.root), in); err != nil {
			t.Fatal(err)
		}
		raced = f.metadata()
	}}
	_, err := f.deps(h).SetWrapper(context.Background(), r)
	if requireKind(t, err, repoestablish.ErrConflict); raced == "" || f.metadata() != raced {
		t.Fatalf("race overwrote metadata: %v", err)
	}
}

func TestCheckStopsBeforeWrapperOnMetadataAndTools(t *testing.T) {
	t.Run("corrupt metadata", func(t *testing.T) {
		f := newFixture(t)
		write(t, filepath.Join(f.root, workspace.MarkerDirectory, workspace.ProjectMetadataFile), "{")
		_, err := repoestablish.Run(context.Background(), f.request(), f.deps(hooks{}))
		requireStop(t, err, repoestablish.StepCheck, repoestablish.ErrConflict)
		if _, err := os.Lstat(filepath.Join(f.projects, "code-map")); !os.IsNotExist(err) || len(f.gh.calls) != 0 {
			t.Fatalf("wrapper %v, gh calls %q", err, f.gh.calls)
		}
	})
	t.Run("different wrapper", func(t *testing.T) {
		f := newFixture(t)
		r := established(t, f)
		in := workspace.ProjectMetadataInput{Operation: workspace.ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "code-map", Locator: f.projects, ActorClaim: "other"}
		if _, err := workspace.UpdateProjectMetadata(fixedDependencies(f.root), in); err != nil {
			t.Fatal(err)
		}
		before, calls := f.metadata(), len(f.gh.calls)
		err := f.deps(hooks{}).Check(context.Background(), r)
		if requireKind(t, err, repoestablish.ErrConflict); f.metadata() != before || len(f.gh.calls) != calls {
			t.Fatal("check had an effect")
		}
	})
	t.Run("missing gh", func(t *testing.T) {
		f := newFixture(t)
		t.Setenv("PATH", filepath.Dir(mustLookPath(t, "git")))
		_, err := repoestablish.Run(context.Background(), f.request(), f.deps(hooks{}))
		requireStop(t, err, repoestablish.StepCheck, repoestablish.ErrInvalid)
	})
	t.Run("definition conflict", func(t *testing.T) {
		f := newFixture(t)
		write(t, filepath.Join(f.root, workspacedef.File), strings.Replace(f.definitions(), `"repos": [{"id": "ply", "path": "projects/ply/main"}]`,
			`"repos": [{"id": "ply", "path": "projects/ply/main"}, {"id": "code-map", "path": "projects/ply/main"}]`, 1))
		_, err := repoestablish.Run(context.Background(), f.request(), f.deps(hooks{}))
		requireStop(t, err, repoestablish.StepCheck, repoestablish.ErrConflict)
	})
	t.Run("invalid actor", func(t *testing.T) {
		f := newFixture(t)
		req := f.request()
		req.Actor = " padded "
		_, err := repoestablish.Run(context.Background(), req, f.deps(hooks{}))
		requireStop(t, err, repoestablish.StepCheck, repoestablish.ErrInvalid)
	})
	t.Run("invalid repository ID", func(t *testing.T) {
		f := newFixture(t)
		req := f.request()
		req.RepoName = "Code_Map"
		_, err := repoestablish.Run(context.Background(), req, f.deps(hooks{}))
		requireStop(t, err, repoestablish.StepCheck, repoestablish.ErrInvalid)
	})
}

func mustLookPath(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestClassificationOfWorkspaceErrors(t *testing.T) {
	project := func(class workspace.ErrorClass) error { return &workspace.ProjectError{Class: class} }
	work := func(class workspace.WorkItemErrorClass) error { return &workspace.WorkItemError{Class: class} }
	release := &filelock.ReleaseError{Path: "/w/.ply/projects.lock", Unlock: errors.New("unlock")}
	for name, tc := range map[string]struct {
		err            error
		define, update error
	}{
		"project invalid":   {project(workspace.ErrorProjectInvalidArguments), repoestablish.ErrInvalid, repoestablish.ErrInvalid},
		"project not found": {project(workspace.ErrorProjectNotFound), repoestablish.ErrInvalid, repoestablish.ErrInvalid},
		"work invalid":      {work(workspace.ErrorWorkInvalidArguments), repoestablish.ErrInvalid, repoestablish.ErrInvalid},
		"project conflict":  {project(workspace.ErrorProjectConflict), repoestablish.ErrConflict, repoestablish.ErrConflict},
		"work project":      {work(workspace.ErrorWorkProjectConflict), repoestablish.ErrConflict, repoestablish.ErrConflict},
		"store conflict":    {project(workspace.ErrorProjectStoreConflict), repoestablish.ErrConflict, repoestablish.ErrUnknown},
		"identity conflict": {work(workspace.ErrorWorkIdentityConflict), repoestablish.ErrConflict, repoestablish.ErrUnknown},
		"work store":        {work(workspace.ErrorWorkStoreConflict), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
		"work IO":           {work(workspace.ErrorWorkIO), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
		"reconciliation":    {work(workspace.ErrorWorkReconciliation), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
		"project IO":        {project(workspace.ErrorProjectIO), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
		"release":           {fmt.Errorf("wrapped: %w", release), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
		"untyped":           {errors.New("boom"), repoestablish.ErrUnknown, repoestablish.ErrUnknown},
	} {
		t.Run(name, func(t *testing.T) {
			for _, got := range []struct{ err, want error }{
				{classify(tc.err, defineConflicts, repoestablish.ErrUnknown), tc.define},
				{classify(tc.err, metadataConflicts, repoestablish.ErrUnknown), tc.update},
			} {
				if !errors.Is(got.err, got.want) || !errors.Is(got.err, tc.err) {
					t.Fatalf("classified %v, want %v wrapping the original", got.err, got.want)
				}
				for _, other := range []error{repoestablish.ErrInvalid, repoestablish.ErrConflict, repoestablish.ErrUnknown} {
					if other != got.want && errors.Is(got.err, other) {
						t.Fatalf("classified %v also as %v", got.err, other)
					}
				}
			}
		})
	}
}

func TestProjectAtRootSelectsExactlyOnePhysicalRoot(t *testing.T) {
	f := newFixture(t)
	if project, err := ProjectAtRoot(f.root, f.projects); err != nil || project != "ply" {
		t.Fatalf("project %q, error %v", project, err)
	}
	for name, cwd := range map[string]string{
		"workspace root": f.root,
		"subdirectory":   filepath.Join(f.projects, "ply"),
		"outside":        filepath.Dir(f.root),
	} {
		if _, err := ProjectAtRoot(f.root, cwd); !errors.Is(err, repoestablish.ErrInvalid) || !strings.Contains(err.Error(), "--project") {
			t.Fatalf("%s: %v", name, err)
		}
	}
	write(t, filepath.Join(f.root, workspacedef.File), strings.Replace(f.definitions(), fmt.Sprintf("%q", filepath.Join(f.root, "apps")), `"projects"`, 1))
	if _, err := ProjectAtRoot(f.root, f.projects); !errors.Is(err, repoestablish.ErrInvalid) || !strings.Contains(err.Error(), "--project") {
		t.Fatalf("two matches: %v", err)
	}
	if _, err := ProjectAtRoot(filepath.Join(f.root, "apps"), f.projects); !errors.Is(err, repoestablish.ErrInvalid) {
		t.Fatalf("no definitions: %v", err)
	}
}
