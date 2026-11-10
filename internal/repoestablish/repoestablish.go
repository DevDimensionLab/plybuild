// Package repoestablish establishes one repository of a Project in fixed
// order: wrapper directory, local Git repository, private GitHub repository,
// remote, workspace definition and wrapper metadata. Every effect goes through
// Dependencies; the package has no output, flags, working directory or
// environment.
package repoestablish

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/gitrepo"
)

// Request is the caller's claim. Empty override fields mean "use defaults".
type Request struct {
	Workspace, Project, RepoName, Actor string
	Overrides                           Overrides
}

type Overrides struct{ GitHubOwner, Leaf, Branch, Remote string }

// ProjectSetup is what Prepare reads from disk; values are empty when not configured.
type ProjectSetup struct {
	Root                              string // physical Project root
	GitHubOwner, Leaf, Branch, Remote string
}

// Resolved is built by Run only, and is the sole input of the effect adapters.
type Resolved struct {
	Workspace, Project, ProjectRoot, RepoName string // RepoName is repo ID and GitHub name
	Wrapper, RepoDir                          string // ProjectRoot/RepoName, Wrapper/Leaf
	Leaf, Branch, Remote, GitHubOwner, Actor  string
}

const DefaultLeaf, DefaultBranch, DefaultRemote = "main", "main", "origin"

type Step string

const (
	StepPrepare  Step = "prepare"
	StepCheck    Step = "check"
	StepWrapper  Step = "wrapper"
	StepInit     Step = "init"
	StepGitHub   Step = "github"
	StepRemote   Step = "remote"
	StepDefine   Step = "define"
	StepMetadata Step = "metadata"
)

type StepResult struct {
	Step    Step
	Changed bool
}

type Result struct {
	Resolved Resolved     // set once prepare has succeeded
	SSHURL   string       // set when the github step is confirmed
	Steps    []StepResult // confirmed completed effect steps only, in fixed order
}

var ErrInvalid, ErrConflict, ErrUnknown = errors.New("invalid"), errors.New("conflict"), errors.New("unknown")

// Error names the step where Run stopped. ErrInvalid and ErrConflict mean that
// step had no effect; ErrUnknown means it may have had one; a bare context
// error means it was not started.
type Error struct {
	Step Step
	Err  error
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %v", e.Step, e.Err) }
func (e *Error) Unwrap() error { return e.Err }

// Dependencies are function fields, all required. Each effect function's error
// must wrap exactly one of the three sentinels.
type Dependencies struct {
	Prepare    func(ctx context.Context, req Request) (ProjectSetup, error)
	Check      func(ctx context.Context, r Resolved) error
	Observe    func(ctx context.Context, dir string) (gitrepo.Repo, error)
	Init       func(ctx context.Context, dir, branch string) (bool, error)
	GitHub     func(ctx context.Context, owner, name string) (sshURL string, changed bool, err error)
	Remote     func(ctx context.Context, dir, name, url string) (bool, error)
	Define     func(ctx context.Context, r Resolved) (bool, error)
	SetWrapper func(ctx context.Context, r Resolved) (bool, error)
}

func (d Dependencies) complete() bool {
	return d.Prepare != nil && d.Check != nil && d.Observe != nil && d.Init != nil && d.GitHub != nil &&
		d.Remote != nil && d.Define != nil && d.SetWrapper != nil
}

func invalid(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, args...))
}

func conflict(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrConflict, fmt.Sprintf(format, args...))
}

func unknown(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrUnknown, fmt.Sprintf(format, args...))
}

// Run establishes the repository in fixed order. A rerun with the same
// resolved values resumes: completed steps report Changed false.
func Run(ctx context.Context, req Request, d Dependencies) (Result, error) {
	var result Result
	if ctx == nil || !d.complete() {
		return result, &Error{StepPrepare, invalid("context and all dependencies are required")}
	}
	if err := validateRequest(req); err != nil {
		return result, &Error{StepPrepare, err}
	}
	if err := ctx.Err(); err != nil {
		return result, &Error{StepPrepare, err}
	}
	setup, err := d.Prepare(ctx, req)
	if err != nil {
		return result, &Error{StepPrepare, withoutEffect(err)}
	}
	r, err := resolve(req, setup)
	if err != nil {
		return result, &Error{StepPrepare, err}
	}
	result.Resolved = r
	if err := ctx.Err(); err != nil {
		return result, &Error{StepCheck, err}
	}
	if err := checkPhysical(r); err != nil {
		return result, &Error{StepCheck, err}
	}
	if err := d.Check(ctx, r); err != nil {
		return result, &Error{StepCheck, withoutEffect(err)}
	}
	if err := checkNesting(ctx, r, d.Observe); err != nil {
		return result, &Error{StepCheck, err}
	}
	steps := []struct {
		step Step
		run  func() (bool, error)
	}{
		{StepWrapper, func() (bool, error) { return makeWrapper(r) }},
		{StepInit, func() (bool, error) { return d.Init(ctx, r.RepoDir, r.Branch) }},
		{StepGitHub, func() (bool, error) {
			url, changed, err := d.GitHub(ctx, r.GitHubOwner, r.RepoName)
			if err == nil {
				result.SSHURL = url
			}
			return changed, err
		}},
		{StepRemote, func() (bool, error) { return d.Remote(ctx, r.RepoDir, r.Remote, result.SSHURL) }},
		{StepDefine, func() (bool, error) { return d.Define(ctx, r) }},
		{StepMetadata, func() (bool, error) { return d.SetWrapper(ctx, r) }},
	}
	for _, s := range steps {
		if err := ctx.Err(); err != nil {
			return result, &Error{s.step, err}
		}
		changed, err := s.run()
		if err != nil {
			return result, &Error{s.step, withPossibleEffect(err)}
		}
		result.Steps = append(result.Steps, StepResult{s.step, changed})
	}
	return result, nil
}

// withoutEffect keeps an ErrInvalid or ErrConflict from a reading step and
// classifies anything else, including ErrUnknown, as ErrInvalid.
func withoutEffect(err error) error {
	switch {
	case errors.Is(err, ErrUnknown):
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	case errors.Is(err, ErrInvalid), errors.Is(err, ErrConflict):
		return err
	}
	return fmt.Errorf("%w: %w", ErrInvalid, err)
}

// withPossibleEffect keeps a classified effect error; an unclassified one,
// even a context error, may have had an effect.
func withPossibleEffect(err error) error {
	if errors.Is(err, ErrUnknown) || errors.Is(err, ErrInvalid) || errors.Is(err, ErrConflict) {
		return err
	}
	return fmt.Errorf("%w: %w", ErrUnknown, err)
}

// element reports whether value is exactly one relative path element.
func element(value string) bool {
	return value != "" && value != "." && value != ".." && !strings.ContainsRune(value, '/') &&
		!strings.ContainsRune(value, filepath.Separator) && !filepath.IsAbs(value)
}

func validateRequest(req Request) error {
	switch {
	case !filepath.IsAbs(req.Workspace) || filepath.Clean(req.Workspace) != req.Workspace:
		return invalid("workspace %q must be an absolute clean path", req.Workspace)
	case req.Project == "":
		return invalid("Project is required")
	case req.Actor == "":
		return invalid("actor is required; set --actor or PLY_ACTOR")
	case !element(req.RepoName):
		return invalid("repository name %q must be one path element", req.RepoName)
	}
	return nil
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func resolve(req Request, setup ProjectSetup) (Resolved, error) {
	o := req.Overrides
	r := Resolved{Workspace: req.Workspace, Project: req.Project, ProjectRoot: setup.Root, RepoName: req.RepoName, Actor: req.Actor,
		GitHubOwner: first(o.GitHubOwner, setup.GitHubOwner),
		Leaf:        first(o.Leaf, setup.Leaf, DefaultLeaf),
		Branch:      first(o.Branch, setup.Branch, DefaultBranch),
		Remote:      first(o.Remote, setup.Remote, DefaultRemote)}
	switch {
	case !filepath.IsAbs(r.ProjectRoot) || filepath.Clean(r.ProjectRoot) != r.ProjectRoot:
		return Resolved{}, invalid("Project root %q must be an absolute clean path", r.ProjectRoot)
	case r.GitHubOwner == "":
		return Resolved{}, invalid("GitHub owner is required; set github_owner in the Project defaults or --github-owner")
	case !element(r.Leaf):
		return Resolved{}, invalid("leaf %q must be one path element", r.Leaf)
	case strings.HasPrefix(r.Branch, "-") || strings.HasPrefix(r.Remote, "-"):
		return Resolved{}, invalid("branch and remote must not begin with -")
	}
	r.Wrapper = filepath.Join(r.ProjectRoot, r.RepoName)
	r.RepoDir = filepath.Join(r.Wrapper, r.Leaf)
	if rel, err := filepath.Rel(r.Workspace, r.Wrapper); err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return Resolved{}, invalid("wrapper %s must be inside workspace %s", r.Wrapper, r.Workspace)
	}
	return r, nil
}

// checkPhysical requires a physical workspace and Project root, and a wrapper
// that is absent or a real directory.
func checkPhysical(r Resolved) error {
	for _, path := range []string{r.Workspace, r.ProjectRoot} {
		if physical, err := filepath.EvalSymlinks(path); err != nil || physical != path {
			return invalid("%s must be an existing physical directory", path)
		}
	}
	if info, err := os.Stat(r.ProjectRoot); err != nil || !info.IsDir() {
		return invalid("Project root %s must be a directory", r.ProjectRoot)
	}
	info, err := os.Lstat(r.Wrapper)
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return invalid("inspect wrapper %s: %v", r.Wrapper, err)
	case info.Mode()&os.ModeSymlink != 0 || !info.IsDir():
		return conflict("wrapper %s is a symlink or not a directory", r.Wrapper)
	}
	return nil
}

// checkNesting requires an existing wrapper to hold no repository or the
// repository on Branch, and a missing wrapper not to be inside a repository.
func checkNesting(ctx context.Context, r Resolved, observe func(context.Context, string) (gitrepo.Repo, error)) error {
	dir := r.Wrapper
	if _, err := os.Lstat(r.Wrapper); err == nil {
		dir = r.RepoDir
	}
	repo, err := observe(ctx, dir)
	if err != nil {
		return withoutEffect(err)
	}
	if repo.State == gitrepo.Absent || dir == r.RepoDir && repo.State == gitrepo.Root && repo.Branch == r.Branch {
		return nil
	}
	return conflict("%s is %s on branch %q in a Git repository, want none or %s on %s", dir, repo.State, repo.Branch, r.RepoDir, r.Branch)
}

// makeWrapper creates the wrapper once, or accepts a real directory. The
// Project root is resolved again afterwards: when it changed, the wrapper may
// have been created elsewhere.
func makeWrapper(r Resolved) (bool, error) {
	changed := true
	if err := os.Mkdir(r.Wrapper, 0o755); errors.Is(err, os.ErrExist) {
		info, err := os.Lstat(r.Wrapper)
		switch {
		case err != nil:
			return false, unknown("inspect wrapper %s: %v", r.Wrapper, err)
		case info.Mode()&os.ModeSymlink != 0 || !info.IsDir():
			return false, conflict("wrapper %s is a symlink or not a directory", r.Wrapper)
		}
		changed = false
	} else if err != nil {
		return false, unknown("create wrapper %s: %v", r.Wrapper, err)
	}
	if physical, err := filepath.EvalSymlinks(r.ProjectRoot); err != nil || physical != r.ProjectRoot {
		return false, unknown("Project root %s changed while creating wrapper %s", r.ProjectRoot, r.Wrapper)
	}
	return changed, nil
}
