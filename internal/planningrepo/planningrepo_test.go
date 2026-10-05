package planningrepo

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/plantemplate"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func fixture(t *testing.T) (Dependencies, Input) {
	t.Helper()
	for key, value := range map[string]string{"GIT_CONFIG_NOSYSTEM": "1", "GIT_AUTHOR_NAME": "Planning Test", "GIT_AUTHOR_EMAIL": "plan@example.invalid", "GIT_COMMITTER_NAME": "Planning Test", "GIT_COMMITTER_EMAIL": "plan@example.invalid"} {
		t.Setenv(key, value)
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(root, "gitconfig")
	if err := os.WriteFile(config, []byte("[user]\nname = Planning Test\nemail = plan@example.invalid\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	repo := filepath.Join(root, "source")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	gitOK(t, repo, "init", "--initial-branch=source", "--template=")
	gitOK(t, repo, "commit", "--allow-empty", "-m", "source")
	return SystemDependencies(), Input{Root: root, Language: "en", Repositories: []string{repo}}
}

func gitOK(t *testing.T, dir string, args ...string) string {
	t.Helper()
	result, err := systemGit(dir, args...)
	if err != nil {
		t.Fatalf("git %v in %s: %v", args, dir, err)
	}
	return result
}

func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if entry.Type()&fs.ModeSymlink != 0 {
			link, err := os.Readlink(p)
			if err != nil {
				return err
			}
			value += ":" + link
		} else if !entry.IsDir() {
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			value += ":" + string(data)
		}
		result[rel] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func absent(t *testing.T, p string) {
	t.Helper()
	if _, err := os.Lstat(p); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected absent %s: %v", p, err)
	}
}

func TestCreateRealCommitAndManifest(t *testing.T) {
	d, input := fixture(t)
	before := snapshot(t, input.Repositories[0])
	input.Goal = "Literal $(touch NEVER) `printf no` {{.RootPath}}"
	r, err := Create(d, input)
	if err != nil {
		t.Fatal(err)
	}
	if r.State != "created" || r.Git == nil || r.Git.Branch != "main" || r.Kind != "PlanningRepoCreation@1" {
		t.Fatalf("result: %+v", r)
	}
	if gitOK(t, r.PlanningPath, "rev-list", "--count", "HEAD") != "1" || gitOK(t, r.PlanningPath, "status", "--porcelain=v1") != "" {
		t.Fatal("not one clean commit")
	}
	if r.Git.OID != gitOK(t, r.PlanningPath, "rev-parse", "HEAD") || r.Git.Tree != gitOK(t, r.PlanningPath, "rev-parse", "HEAD^{tree}") {
		t.Fatal("invented Git result")
	}
	if common := gitOK(t, r.PlanningPath, "rev-parse", "--path-format=absolute", "--git-common-dir"); common != filepath.Join(r.PlanningPath, ".git") {
		t.Fatal(common)
	}
	data, err := os.ReadFile(filepath.Join(r.PlanningPath, ManifestName))
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if !validManifest(manifest, r.PlanningPath) || len(manifest.Files)+1 != len(r.Files) {
		t.Fatalf("manifest: %+v", manifest)
	}
	for _, file := range manifest.Files {
		b, err := os.ReadFile(filepath.Join(r.PlanningPath, file.Path))
		if err != nil {
			t.Fatal(err)
		}
		if want := "sha256:" + fmtHash(b); file.SHA256 != want {
			t.Fatalf("hash %s: %s != %s", file.Path, file.SHA256, want)
		}
	}
	goal, _ := os.ReadFile(filepath.Join(r.PlanningPath, "goal.md"))
	if !strings.Contains(string(goal), input.Goal) {
		t.Fatal("goal not preserved as data")
	}
	if !reflect.DeepEqual(before, snapshot(t, input.Repositories[0])) {
		t.Fatal("source changed")
	}
	planBefore := snapshot(t, r.PlanningPath)
	for _, preview := range []bool{false, true} {
		input.Check = preview
		if _, err := Create(d, input); err == nil || !strings.Contains(err.Error(), "existing planning repository") {
			t.Fatalf("repeat: %v", err)
		}
		if !reflect.DeepEqual(planBefore, snapshot(t, r.PlanningPath)) {
			t.Fatal("existing plan changed")
		}
	}
}

func fmtHash(b []byte) string {
	const hex = "0123456789abcdef"
	h := sha256.Sum256(b)
	out := make([]byte, 64)
	for i, v := range h {
		out[2*i], out[2*i+1] = hex[v>>4], hex[v&15]
	}
	return string(out)
}

func TestPrewriteFailuresAndPreview(t *testing.T) {
	for _, scenario := range []string{"preview", "render", "size", "identity", "unsafe", "duplicate", "inside-source", "inside-git", "other-worktree", "missing-root", "subdirectory", "registered-drift"} {
		t.Run(scenario, func(t *testing.T) {
			d, input := fixture(t)
			source := input.Repositories[0]
			switch scenario {
			case "preview":
				input.Check = true
			case "render":
				d.Render = func(plantemplate.Input) (plantemplate.Bundle, error) {
					return plantemplate.Bundle{}, errors.New("injected rendering error")
				}
			case "size":
				input.Goal = strings.Repeat("large", 10000)
			case "identity":
				real := d.Git
				d.Git = func(dir string, args ...string) (string, error) {
					if strings.Contains(strings.Join(args, " "), "var GIT_") {
						return "", errors.New("missing identity")
					}
					return real(dir, args...)
				}
			case "unsafe":
				d.Render = func(plantemplate.Input) (plantemplate.Bundle, error) {
					return plantemplate.Bundle{Files: []plantemplate.File{{Path: "../escape"}}}, nil
				}
			case "duplicate":
				linked := filepath.Join(input.Root, "linked")
				gitOK(t, source, "worktree", "add", "-b", "other", linked)
				input.Repositories = append(input.Repositories, linked)
			case "inside-source":
				input.Root = source
			case "inside-git":
				input.Root = filepath.Join(source, ".git")
			case "other-worktree":
				other := filepath.Join(input.Root, "other")
				if err := os.Mkdir(other, 0o755); err != nil {
					t.Fatal(err)
				}
				gitOK(t, other, "init", "--template=")
				input.Root = other
			case "missing-root":
				input.Root = filepath.Join(input.Root, "absent")
			case "subdirectory":
				deep := filepath.Join(source, "deep")
				if err := os.Mkdir(deep, 0o755); err != nil {
					t.Fatal(err)
				}
				input.Repositories = []string{deep}
			case "registered-drift":
				input.Project = "example"
				input.Repositories = nil
				d.Project = func(workspace.ProjectID) (workspace.ProjectResult, error) {
					return workspace.ProjectResult{Project: workspace.ProjectRecord{Wrapper: input.Root}, Repos: []workspace.RepoRecord{{Locator: source, GitCommonDir: "/wrong/identity"}}}, nil
				}
			}
			before := snapshot(t, source)
			d.Mkdir = func(string, fs.FileMode) error { t.Fatal("prewrite failure performed mkdir"); return nil }
			d.WriteFile = func(string, []byte) error { t.Fatal("prewrite failure wrote a file"); return nil }
			r, err := Create(d, input)
			if scenario == "preview" {
				if err != nil || r.State != "preview" || r.Git != nil || len(r.Files) != 14 {
					t.Fatalf("preview: %+v %v", r, err)
				}
			} else if err == nil {
				t.Fatal("expected failure")
			}
			absent(t, filepath.Join(input.Root, "planning"))
			if !reflect.DeepEqual(before, snapshot(t, source)) {
				t.Fatal("source changed during validation")
			}
		})
	}
}

func TestPartialFailuresPreserveStateAndNeverReplay(t *testing.T) {
	for _, stage := range []string{"git init", "write goal.md", "mkdir archive", "git add", "git commit", "verify commit", "verify clean worktree"} {
		t.Run(stage, func(t *testing.T) {
			d, input := fixture(t)
			target := filepath.Join(input.Root, "planning")
			realGit, realWrite, realMkdir := d.Git, d.WriteFile, d.Mkdir
			effects, failed := 0, false
			d.Git = func(dir string, args ...string) (string, error) {
				joined := strings.Join(args, " ")
				if dir == target {
					isEffect := args[0] == "init" || args[0] == "add" || strings.Contains(joined, " commit ")
					if isEffect {
						effects++
					}
					match := (stage == "git init" && args[0] == "init") || (stage == "git add" && args[0] == "add") ||
						(stage == "git commit" && strings.Contains(joined, " commit ")) ||
						(stage == "verify commit" && args[0] == "symbolic-ref") || (stage == "verify clean worktree" && args[0] == "status")
					if match {
						failed = true
						return "", errors.New("injected Git fault")
					}
				}
				return realGit(dir, args...)
			}
			d.WriteFile = func(p string, b []byte) error {
				if stage == "write goal.md" && filepath.Base(p) == "goal.md" {
					failed = true
					return errors.New("injected write fault")
				}
				return realWrite(p, b)
			}
			d.Mkdir = func(p string, mode fs.FileMode) error {
				if stage == "mkdir archive" && filepath.Base(p) == "archive" {
					failed = true
					return errors.New("injected mkdir fault")
				}
				return realMkdir(p, mode)
			}
			r, err := Create(d, input)
			if !failed || err == nil || !strings.Contains(err.Error(), target) || !strings.Contains(err.Error(), stage) || !strings.Contains(err.Error(), "partial state preserved") || r.State != "" {
				t.Fatalf("fault: %+v %v", r, err)
			}
			before, count := snapshot(t, target), effects
			if _, err := Create(d, input); err == nil || !strings.Contains(err.Error(), "creation stopped") {
				t.Fatalf("repeat: %v", err)
			}
			if effects != count || !reflect.DeepEqual(before, snapshot(t, target)) {
				t.Fatal("failed operation was retried or partial state changed")
			}
		})
	}
}

func TestConcurrentCreatorWins(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(map[bool]string{false: "conflict", true: "legacy-plan"}[legacy], func(t *testing.T) {
			d, input := fixture(t)
			target := filepath.Join(input.Root, "planning")
			var before map[string]string
			d.Mkdir = func(p string, mode fs.FileMode) error {
				if p != target {
					t.Fatal("wrote after losing destination")
				}
				if err := os.Mkdir(p, mode); err != nil {
					t.Fatal(err)
				}
				if legacy {
					gitOK(t, p, "init", "--template=")
					for _, name := range []string{"README.md", "AGENTS.md", "goal.md", "next-task.md", "status.md"} {
						if err := os.WriteFile(filepath.Join(p, name), []byte("other actor\n"), 0o644); err != nil {
							t.Fatal(err)
						}
					}
				}
				if err := os.WriteFile(filepath.Join(p, "winner"), []byte("preserve me"), 0o644); err != nil {
					t.Fatal(err)
				}
				before = snapshot(t, p)
				return fs.ErrExist
			}
			_, err := Create(d, input)
			want := "path conflict"
			if legacy {
				want = "existing planning repository"
			}
			if err == nil || !strings.Contains(err.Error(), want) || !reflect.DeepEqual(before, snapshot(t, target)) {
				t.Fatalf("race: %v", err)
			}
		})
	}
}

func TestManifestRecognitionDoesNotReadPlanningInstructions(t *testing.T) {
	d, input := fixture(t)
	r, err := Create(d, input)
	if err != nil {
		t.Fatal(err)
	}
	// A versioned plan remains recognizable after its human documents evolve.
	if err := os.Remove(filepath.Join(r.PlanningPath, "AGENTS.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Create(d, input); err == nil || !strings.Contains(err.Error(), "existing planning repository") {
		t.Fatalf("valid manifest: %v", err)
	}
	b, _ := os.ReadFile(filepath.Join(r.PlanningPath, ManifestName))
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	m.PlanningPath = "/wrong/physical/path"
	b, _ = json.Marshal(m)
	if err := os.WriteFile(filepath.Join(r.PlanningPath, ManifestName), b, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Create(d, input); err == nil || !strings.Contains(err.Error(), "path conflict") {
		t.Fatalf("invalid manifest: %v", err)
	}
}

func TestCreationStopsWhenBoundDirectoryIsReplaced(t *testing.T) {
	for _, stage := range []string{"root-during-render", "target-after-init", "target-during-write", "git-directory-during-write"} {
		t.Run(stage, func(t *testing.T) {
			d, input := fixture(t)
			source := input.Repositories[0]
			input.Root = filepath.Join(input.Root, "chosen root")
			if err := os.Mkdir(input.Root, 0o755); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(input.Root, "planning")
			before := snapshot(t, source)
			swapped := false
			replace := func(p string) {
				t.Helper()
				if err := os.Rename(p, p+"-preserved"); err != nil {
					t.Fatal(err)
				}
				destination := source
				if stage == "git-directory-during-write" {
					destination = filepath.Join(source, ".git")
				}
				if err := os.Symlink(destination, p); err != nil {
					t.Fatal(err)
				}
				swapped = true
			}
			switch stage {
			case "root-during-render":
				real := d.Render
				d.Render = func(in plantemplate.Input) (plantemplate.Bundle, error) {
					bundle, err := real(in)
					replace(input.Root)
					return bundle, err
				}
			case "target-after-init":
				real := d.Git
				d.Git = func(dir string, args ...string) (string, error) {
					out, err := real(dir, args...)
					if dir == target && args[0] == "init" && err == nil {
						replace(target)
					}
					return out, err
				}
			case "target-during-write", "git-directory-during-write":
				real := d.WriteFile
				d.WriteFile = func(p string, content []byte) error {
					err := real(p, content)
					if !swapped && err == nil {
						if stage == "git-directory-during-write" {
							replace(filepath.Join(target, ".git"))
						} else {
							replace(target)
						}
					}
					return err
				}
			}
			r, err := Create(d, input)
			if !swapped {
				t.Fatal("test did not reach replacement boundary")
			}
			if !reflect.DeepEqual(before, snapshot(t, source)) {
				t.Fatal("creation followed a replaced directory and changed the source repository")
			}
			if err == nil || r.State != "" || !strings.Contains(err.Error(), "changed") {
				t.Fatalf("directory replacement should stop: %+v %v", r, err)
			}
		})
	}
}
