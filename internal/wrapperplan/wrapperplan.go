// Package wrapperplan creates the local planning repository <repo>-plan in a
// repository wrapper from a template directory: the template's files with
// {{repo}} and {{date}} filled in, gitignore named .gitignore and the
// template's own README.md left out, committed once on a branch without
// remote. The repository is built beside the target and renamed into place,
// so the target is either absent or complete.
package wrapperplan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrConflict = errors.New("conflict")
	ErrUnknown  = errors.New("unknown")
)

const (
	// DefaultBranch is the branch of a planning repository.
	DefaultBranch = "main"
	// templateReadme describes the template itself and is not copied.
	templateReadme = "README.md"
	commitMessage  = "Create planning repository from template"
)

// Request names the template, the wrapper and the values filled in.
type Request struct {
	Template, Wrapper, Repo, Date, Branch string
}

var (
	placeholder = regexp.MustCompile(`\{\{[^{}]*\}\}`)
	datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func fail(kind error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", kind, fmt.Sprintf(format, args...))
}

// Dir is the planning repository of repo in wrapper.
func Dir(wrapper, repo string) string { return filepath.Join(wrapper, repo+"-plan") }

func staging(req Request) string { return filepath.Join(req.Wrapper, "."+req.Repo+"-plan.creating") }

// Check reads the template and the target without any effect: the template
// must be usable, and the target absent or already a repository.
func Check(req Request) error {
	if _, err := files(req); err != nil {
		return err
	}
	if _, err := os.Lstat(staging(req)); err == nil {
		return fail(ErrConflict, "%s is left from an interrupted creation; inspect and remove it", staging(req))
	}
	target := Dir(req.Wrapper, req.Repo)
	info, err := os.Lstat(target)
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return fail(ErrInvalid, "inspect %s: %v", target, err)
	case info.Mode()&os.ModeSymlink != 0 || !info.IsDir():
		return fail(ErrConflict, "%s is a symlink or not a directory", target)
	}
	if _, err := os.Lstat(filepath.Join(target, ".git")); err != nil {
		return fail(ErrConflict, "%s exists and is not a planning repository", target)
	}
	return nil
}

// Ensure creates the planning repository, or verifies that the target is a
// repository root on the branch with a commit.
func Ensure(ctx context.Context, runner process.Runner, req Request) (bool, error) {
	if err := Check(req); err != nil {
		return false, err
	}
	target := Dir(req.Wrapper, req.Repo)
	if _, err := os.Lstat(target); err == nil {
		return false, existing(ctx, runner, target, req.Branch)
	}
	content, _ := files(req)
	dir := staging(req)
	if _, _, err := gitrepo.EnsureInit(ctx, runner, dir, req.Branch); err != nil {
		return false, fail(ErrUnknown, "%v", err)
	}
	names := map[string]bool{}
	for rel, f := range content {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return false, fail(ErrUnknown, "create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, f.data, f.mode); err != nil {
			return false, fail(ErrUnknown, "write %s: %v", path, err)
		}
		names[strings.SplitN(rel, string(filepath.Separator), 2)[0]] = true
	}
	top := make([]string, 0, len(names))
	for name := range names {
		top = append(top, name)
	}
	if _, err := gitrepo.EnsureInitialCommit(ctx, runner, dir, req.Branch, commitMessage, top); err != nil {
		return false, fail(ErrUnknown, "%v", err)
	}
	if err := os.Rename(dir, target); err != nil {
		return false, fail(ErrUnknown, "move %s to %s: %v", dir, target, err)
	}
	if err := existing(ctx, runner, target, req.Branch); err != nil {
		return false, fail(ErrUnknown, "readback %s: %v", target, err)
	}
	return true, nil
}

// existing accepts a repository root on branch with a commit.
func existing(ctx context.Context, runner process.Runner, target, branch string) error {
	repo, err := gitrepo.Observe(ctx, runner, target)
	if err != nil {
		return fail(ErrUnknown, "%v", err)
	}
	if repo.State != gitrepo.Root || repo.Toplevel != target || repo.Branch != branch {
		return fail(ErrConflict, "%s is %s on branch %q, want a planning repository on %q", target, repo.State, repo.Branch, branch)
	}
	committed, err := gitrepo.HasCommit(ctx, runner, target)
	switch {
	case err != nil:
		return fail(ErrUnknown, "%v", err)
	case !committed:
		return fail(ErrConflict, "%s has no commit on %q", target, branch)
	}
	return nil
}

type file struct {
	data []byte
	mode fs.FileMode
}

// files reads the template as the planning repository's files by relative path.
func files(req Request) (map[string]file, error) {
	switch {
	case !filepath.IsAbs(req.Template) || !filepath.IsAbs(req.Wrapper):
		return nil, fail(ErrInvalid, "template and wrapper must be absolute paths")
	case req.Repo == "" || strings.ContainsAny(req.Repo, "/{}") || req.Repo == "." || req.Repo == "..":
		return nil, fail(ErrInvalid, "repository name %q must be one path element", req.Repo)
	case !datePattern.MatchString(req.Date):
		return nil, fail(ErrInvalid, "date %q must be YYYY-MM-DD", req.Date)
	case req.Branch == "":
		return nil, fail(ErrInvalid, "branch is required")
	}
	out := map[string]file{}
	err := filepath.WalkDir(req.Template, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(req.Template, path)
		switch {
		case entry.IsDir():
			return nil
		case !entry.Type().IsRegular():
			return fmt.Errorf("%s is not a regular file", rel)
		case rel == templateReadme:
			return nil
		case rel == ".gitignore" || filepath.Base(rel) == ".git":
			return fmt.Errorf("%s would act inside the template; name it without the dot", rel)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		data = bytes.ReplaceAll(bytes.ReplaceAll(data, []byte("{{repo}}"), []byte(req.Repo)), []byte("{{date}}"), []byte(req.Date))
		if left := placeholder.Find(data); left != nil {
			return fmt.Errorf("%s has the unknown placeholder %s", rel, left)
		}
		if rel == "gitignore" {
			rel = ".gitignore"
		}
		out[rel] = file{data: data, mode: info.Mode().Perm()}
		return nil
	})
	if err != nil {
		return nil, fail(ErrInvalid, "template %s: %v", req.Template, err)
	}
	if len(out) == 0 {
		return nil, fail(ErrInvalid, "template %s has no files", req.Template)
	}
	return out, nil
}
