package wrapperplan

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

func setup(t *testing.T) (string, string, process.Runner) {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(base, "gitconfig")
	if err := os.WriteFile(config, []byte("[user]\n\tname = Ply Test\n\temail = ply@example.invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", base)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	template := filepath.Join(base, "template")
	for name, content := range map[string]string{
		"AGENTS.md":           "# {{repo}}-plan\n\nLaget {{date}}.\n",
		"README.md":           "# Mal: planrepo\n",
		"gitignore":           ".build/\n",
		"claude-dev-start.sh": "#!/bin/sh\nexec python3 scripts/start.py \"$@\"\n",
		"scripts/start.py":    "OWNER = \"{{repo}}-plan\"\n",
	} {
		path := filepath.Join(template, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(filepath.Join(template, "claude-dev-start.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(base, "demo")
	if err := os.Mkdir(wrapper, 0o755); err != nil {
		t.Fatal(err)
	}
	return template, wrapper, process.SystemRunner().Runner
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

func request(template, wrapper string) Request {
	return Request{Template: template, Wrapper: wrapper, Repo: "demo", Date: "2026-11-15", Branch: DefaultBranch}
}

func TestEnsureCreatesTheFilledPlanningRepositoryOnce(t *testing.T) {
	template, wrapper, runner := setup(t)

	changed, err := Ensure(context.Background(), runner, request(template, wrapper))

	target := filepath.Join(wrapper, "demo-plan")
	if err != nil || !changed || Dir(wrapper, "demo") != target {
		t.Fatalf("first Ensure = %v, %v", changed, err)
	}
	for name, want := range map[string]string{
		"AGENTS.md":        "# demo-plan\n\nLaget 2026-11-15.\n",
		".gitignore":       ".build/\n",
		"scripts/start.py": "OWNER = \"demo-plan\"\n",
	} {
		if got, err := os.ReadFile(filepath.Join(target, name)); err != nil || string(got) != want {
			t.Fatalf("%s = %q, %v", name, got, err)
		}
	}
	for _, name := range []string{"gitignore", "README.md"} {
		if _, err := os.Lstat(filepath.Join(target, name)); !os.IsNotExist(err) {
			t.Fatalf("%s copied: %v", name, err)
		}
	}
	if info, err := os.Stat(filepath.Join(target, "claude-dev-start.sh")); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Fatalf("launcher not executable: %v %v", info, err)
	}
	if branch := git(t, target, "branch", "--show-current"); branch != "main" {
		t.Fatalf("branch %q", branch)
	}
	if count := git(t, target, "rev-list", "--count", "HEAD"); count != "1" {
		t.Fatalf("commits %s", count)
	}
	if status, remotes := git(t, target, "status", "--porcelain"), git(t, target, "remote"); status != "" || remotes != "" {
		t.Fatalf("status %q remotes %q", status, remotes)
	}
	if entries, _ := os.ReadDir(wrapper); len(entries) != 1 {
		t.Fatalf("wrapper holds more than the planning repository: %v", entries)
	}
	head := git(t, target, "rev-parse", "HEAD")
	if changed, err = Ensure(context.Background(), runner, request(template, wrapper)); err != nil || changed || git(t, target, "rev-parse", "HEAD") != head {
		t.Fatalf("second Ensure = %v, %v; want unchanged", changed, err)
	}
}

func TestCheckRejectsAnUnusableTemplateBeforeAnyEffect(t *testing.T) {
	for name, change := range map[string]func(template string) error{
		"missing": func(template string) error { return os.RemoveAll(template) },
		"unknown placeholder": func(template string) error {
			return os.WriteFile(filepath.Join(template, "plan.md"), []byte("{{owner}}\n"), 0o644)
		},
		"symlink":       func(template string) error { return os.Symlink("AGENTS.md", filepath.Join(template, "link.md")) },
		"dot gitignore": func(template string) error { return os.WriteFile(filepath.Join(template, ".gitignore"), nil, 0o644) },
	} {
		t.Run(name, func(t *testing.T) {
			template, wrapper, runner := setup(t)
			if err := change(template); err != nil {
				t.Fatal(err)
			}
			req := request(template, wrapper)
			if err := Check(req); !errors.Is(err, ErrInvalid) {
				t.Fatalf("Check = %v, want ErrInvalid", err)
			}
			if changed, err := Ensure(context.Background(), runner, req); !errors.Is(err, ErrInvalid) || changed {
				t.Fatalf("Ensure = %v, %v", changed, err)
			}
			if entries, _ := os.ReadDir(wrapper); len(entries) != 0 {
				t.Fatalf("effect despite unusable template: %v", entries)
			}
		})
	}
}

func TestEnsureNeverTakesOverAnotherPath(t *testing.T) {
	for name, place := range map[string]func(wrapper string) error{
		"plain directory":     func(wrapper string) error { return os.Mkdir(filepath.Join(wrapper, "demo-plan"), 0o755) },
		"file":                func(wrapper string) error { return os.WriteFile(filepath.Join(wrapper, "demo-plan"), nil, 0o644) },
		"interrupted attempt": func(wrapper string) error { return os.Mkdir(filepath.Join(wrapper, ".demo-plan.creating"), 0o755) },
	} {
		t.Run(name, func(t *testing.T) {
			template, wrapper, runner := setup(t)
			if err := place(wrapper); err != nil {
				t.Fatal(err)
			}
			req := request(template, wrapper)
			checkErr := Check(req)
			changed, err := Ensure(context.Background(), runner, req)
			if !errors.Is(checkErr, ErrConflict) || !errors.Is(err, ErrConflict) || changed {
				t.Fatalf("Check = %v, Ensure = %v, %v; want ErrConflict", checkErr, changed, err)
			}
		})
	}
}
