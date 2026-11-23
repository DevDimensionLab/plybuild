package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// memberBreak replaces the directory of member gone after registration.
type memberBreak struct {
	name, status, diagnostic string
	apply                    func(path string) error
}

var (
	missingBreak     = memberBreak{"missing", "missing", "the directory does not exist", func(string) error { return nil }}
	unreadableBreaks = []memberBreak{
		{"file", "unreadable", "repository is not a directory", func(path string) error { return os.WriteFile(path, []byte("not a directory"), 0o600) }},
		{"not git", "unreadable", "is not a Git worktree root", func(path string) error { return os.Mkdir(path, 0o755) }},
		{"broken git pointer", "unreadable", "is not a Git worktree root", func(path string) error {
			if err := os.Mkdir(path, 0o755); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(path, ".git"), []byte("gitdir: "+filepath.Join(path, "absent")+"\n"), 0o600)
		}},
	}
)

// missingMemberFixture is a workspace whose Project ply has the members kept
// and gone; the directory of gone is removed after registration.
func missingMemberFixture(t *testing.T) (taskwork.Dependencies, string) {
	t.Helper()
	return brokenMemberFixture(t, missingBreak)
}

func brokenMemberFixture(t *testing.T, broken memberBreak) (taskwork.Dependencies, string) {
	t.Helper()
	d, root := viewFixture(t)
	for _, repo := range []string{"kept", "gone"} {
		path := filepath.Join(root, "ply", repo)
		for _, args := range [][]string{{"init", "-q", "-b", "main", path}, {"-C", path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture"}} {
			if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v %s", args, err, out)
			}
		}
	}
	definitions := `{"format_version": 1, "projects": [{"id": "ply", "name": "Ply", "wrapper": "ply", "repos": [{"id": "gone", "path": "ply/gone"}, {"id": "kept", "path": "ply/kept"}]}]}`
	if err := os.WriteFile(filepath.Join(root, workspace.DefinitionsFile), []byte(definitions), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.ShowProject(d.Dependencies, "ply"); err != nil {
		t.Fatalf("complete Project: %v", err)
	}
	gone := filepath.Join(root, "ply", "gone")
	if err := os.RemoveAll(gone); err != nil {
		t.Fatal(err)
	}
	if err := broken.apply(gone); err != nil {
		t.Fatal(err)
	}
	return d, root
}

func runBacklog(t *testing.T, d taskwork.Dependencies, args ...string) (string, error) {
	t.Helper()
	return runBacklogInput(t, d, "", args...)
}

func runBacklogInput(t *testing.T, d taskwork.Dependencies, stdin string, args ...string) (string, error) {
	t.Helper()
	c := newBacklogCommand(d.Dependencies)
	var out bytes.Buffer
	c.SetIn(strings.NewReader(stdin))
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceErrors = true
	c.SilenceUsage = true
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func runInboxInput(t *testing.T, d taskwork.Dependencies, stdin string, args ...string) (string, error) {
	t.Helper()
	c := newInboxCommand(d.Dependencies)
	var out bytes.Buffer
	c.SetIn(strings.NewReader(stdin))
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceErrors = true
	c.SilenceUsage = true
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

// showRepositories is project show --format json of Project ply, by repo ID.
func showRepositories(t *testing.T, d taskwork.Dependencies) (map[string]map[string]any, string) {
	t.Helper()
	out, err := runWorkspaceView(t, d, "project", "show", "--format", "json", "--", "ply")
	if err != nil {
		t.Fatalf("project show: %v", err)
	}
	var shown struct {
		Repositories []map[string]any `json:"repositories"`
	}
	if err := json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatalf("project show = %s %v", out, err)
	}
	byID := map[string]map[string]any{}
	for _, repo := range shown.Repositories {
		byID[repo["repo_id"].(string)] = repo
	}
	return byID, out
}

func TestProjectWithMissingMemberDirectoryStaysReadable(t *testing.T) {
	d, root := missingMemberFixture(t)

	out, err := runWorkspaceView(t, d, "project", "list", "--format", "json")
	if err != nil {
		t.Fatalf("project list: %v", err)
	}
	var list struct {
		Projects []struct {
			ID        string `json:"project_id"`
			RepoCount int    `json:"repo_count"`
		} `json:"projects"`
	}
	if err := json.Unmarshal([]byte(out), &list); err != nil || len(list.Projects) != 1 || list.Projects[0].ID != "ply" || list.Projects[0].RepoCount != 2 {
		t.Fatalf("project list = %s %v", out, err)
	}

	out, err = runWorkspaceView(t, d, "project", "show", "--format", "json", "--", "ply")
	if err != nil {
		t.Fatalf("project show: %v", err)
	}
	var shown struct {
		Repositories []map[string]any `json:"repositories"`
	}
	if err := json.Unmarshal([]byte(out), &shown); err != nil || len(shown.Repositories) != 2 {
		t.Fatalf("project show = %s %v", out, err)
	}
	gone, kept := shown.Repositories[0], shown.Repositories[1]
	if gone["repo_id"] != "gone" || gone["status"] != "missing" || gone["locator"] != filepath.Join(root, "ply", "gone") {
		t.Fatalf("missing member = %s", out)
	}
	if commonDir, present := gone["git_common_dir"]; !present || commonDir != nil || gone["diagnostic"] != "the directory does not exist" {
		t.Fatalf("missing member git_common_dir and diagnostic = %s", out)
	}
	if kept["repo_id"] != "kept" || kept["status"] != "present" || kept["locator"] != filepath.Join(root, "ply", "kept") || kept["git_common_dir"] != filepath.Join(root, "ply", "kept", ".git") {
		t.Fatalf("present member = %s", out)
	}
	if diagnostic, present := kept["diagnostic"]; !present || diagnostic != nil {
		t.Fatalf("present member diagnostic = %s", out)
	}

	out, err = runWorkspaceView(t, d, "project", "show", "ply")
	if err != nil || !strings.Contains(out, "  gone:\n    locator: "+filepath.Join(root, "ply", "gone")+"\n    missing: the directory does not exist\n") {
		t.Fatalf("project show text = %s %v", out, err)
	}

	if _, err := runBacklog(t, d, "list", "--project", "ply", "--format", "json"); err != nil {
		t.Fatalf("backlog list: %v", err)
	}

	// Missing members have no Git identity, so two of them do not conflict.
	if err := os.RemoveAll(filepath.Join(root, "ply", "kept")); err != nil {
		t.Fatal(err)
	}
	out, err = runWorkspaceView(t, d, "project", "show", "--format", "json", "--", "ply")
	if err != nil || strings.Count(out, `"status":"missing"`) != 2 {
		t.Fatalf("project show with every member missing = %s %v", out, err)
	}
}

func TestProjectWithUnreadableMemberStaysReadable(t *testing.T) {
	for _, broken := range unreadableBreaks {
		t.Run(broken.name, func(t *testing.T) {
			d, root := brokenMemberFixture(t, broken)
			if out, err := runWorkspaceView(t, d, "project", "list", "--format", "json"); err != nil || !strings.Contains(out, `"repo_count":2`) {
				t.Fatalf("project list = %s %v", out, err)
			}
			repos, out := showRepositories(t, d)
			gone, kept := repos["gone"], repos["kept"]
			diagnostic, _ := gone["diagnostic"].(string)
			if gone["status"] != "unreadable" || gone["locator"] != filepath.Join(root, "ply", "gone") || !strings.Contains(diagnostic, broken.diagnostic) {
				t.Fatalf("unreadable member = %s", out)
			}
			if commonDir, present := gone["git_common_dir"]; !present || commonDir != nil {
				t.Fatalf("unreadable member git_common_dir = %s", out)
			}
			if kept["status"] != "present" || kept["git_common_dir"] != filepath.Join(root, "ply", "kept", ".git") {
				t.Fatalf("present member = %s", out)
			}
			text, err := runWorkspaceView(t, d, "project", "show", "ply")
			if err != nil || !strings.Contains(text, "  gone:\n    locator: "+filepath.Join(root, "ply", "gone")+"\n    unreadable: "+diagnostic+"\n") {
				t.Fatalf("project show text = %s %v", text, err)
			}
			if out, err := runBacklog(t, d, "list", "--project", "ply", "--format", "json"); err != nil {
				t.Fatalf("backlog list = %s %v", out, err)
			}

			// Unreadable members have no Git identity, so two of them do not conflict.
			keptPath := filepath.Join(root, "ply", "kept")
			if err := os.RemoveAll(keptPath); err != nil {
				t.Fatal(err)
			}
			if err := broken.apply(keptPath); err != nil {
				t.Fatal(err)
			}
			if out, err := runWorkspaceView(t, d, "project", "show", "--format", "json", "--", "ply"); err != nil || strings.Count(out, `"status":"unreadable"`) != 2 {
				t.Fatalf("project show with every member unreadable = %s %v", out, err)
			}
		})
	}
}

// Project, Epic and Task writes need every member's Git identity, so an
// unavailable member stops them with an error that names it.
func TestProjectWithUnavailableMemberKeepsWritesStrict(t *testing.T) {
	for _, broken := range append([]memberBreak{missingBreak}, unreadableBreaks...) {
		t.Run(broken.name, func(t *testing.T) {
			d, root := brokenMemberFixture(t, broken)
			extra := filepath.Join(root, "ply", "extra")
			for _, args := range [][]string{{"init", "-q", "-b", "main", extra}, {"-C", extra, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture"}} {
				if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
					t.Fatalf("git %v: %v %s", args, err, out)
				}
			}
			definitions := filepath.Join(root, workspace.DefinitionsFile)
			before, err := os.ReadFile(definitions)
			if err != nil {
				t.Fatal(err)
			}
			gone := filepath.Join(root, "ply", "gone")
			for _, args := range [][]string{
				{"project", "add-member", "ply", "--repo", "extra=ply/extra"},
				{"project", "companion", "add", "ply", "--role", "planning", "--path", "ply/extra"},
				{"project", "repo-wrapper", "set", "ply", "--repo", "kept", "--wrapper", "ply"},
			} {
				_, err := runWorkspaceView(t, d, args...)
				var projectErr *workspace.ProjectError
				if !errors.As(err, &projectErr) || !strings.Contains(projectErr.Detail, gone) {
					t.Fatalf("%v with member gone %s: %v", args, broken.name, err)
				}
				if after, _ := os.ReadFile(definitions); !bytes.Equal(after, before) {
					t.Fatalf("%v changed the definitions", args)
				}
			}
			// Removing the member needs no directory and makes writes work again.
			if _, err := runWorkspaceView(t, d, "project", "remove-member", "ply", "--repo", "gone"); err != nil {
				t.Fatalf("remove-member: %v", err)
			}
			if _, err := runWorkspaceView(t, d, "project", "add-member", "ply", "--repo", "extra=ply/extra"); err != nil {
				t.Fatalf("add-member after removal: %v", err)
			}
		})
	}
}

// Backlog and Inbox data live in the Project wrapper and record no member's
// Git identity, so they are written also while a member is unavailable.
func TestBacklogAndInboxWriteWithUnavailableMember(t *testing.T) {
	for _, broken := range []memberBreak{missingBreak, unreadableBreaks[0]} {
		t.Run(broken.name, func(t *testing.T) {
			d, _ := brokenMemberFixture(t, broken)
			if out, err := runBacklogInput(t, d, "Tekst\n", "add", "--project", "ply", "--title", "Sak", "--format", "json"); err != nil || !strings.Contains(out, `"title": "Sak"`) {
				t.Fatalf("backlog add = %s %v", out, err)
			}
			if out, err := runInboxInput(t, d, inboxRequest(1, "Svar"), "put", "--project", "ply", "--owner", "owner-1", "--file", "-", "--format", "json"); err != nil || !strings.Contains(out, `"effect"`) {
				t.Fatalf("inbox put = %s %v", out, err)
			}
		})
	}
}

// Two members with one Git identity are a conflict of the Project itself, so
// reading it still stops.
func TestProjectMemberIdentityConflictStopsReads(t *testing.T) {
	d, root := missingMemberFixture(t)
	definitions := `{"format_version": 1, "projects": [{"id": "ply", "name": "Ply", "wrapper": "ply", "repos": [{"id": "gone", "path": "ply/kept"}, {"id": "kept", "path": "ply/kept"}]}]}`
	if err := os.WriteFile(filepath.Join(root, workspace.DefinitionsFile), []byte(definitions), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"project", "list", "--format", "json"}, {"project", "show", "--format", "json", "ply"}} {
		_, err := runWorkspaceView(t, d, args...)
		var projectErr *workspace.ProjectError
		if !errors.As(err, &projectErr) || projectErr.Class != workspace.ErrorProjectConflict {
			t.Fatalf("%v with an identity conflict: %v", args, err)
		}
	}
	if out, err := runBacklog(t, d, "list", "--project", "ply", "--format", "json"); err == nil {
		t.Fatalf("backlog list succeeded: %s", out)
	}
}
