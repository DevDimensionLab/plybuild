// Package system binds repoestablish to Git, gh and the Ply workspace. It
// reads no working directory and no environment variable of its own: every
// workspace operation runs with the explicit workspace as its working
// directory.
package system

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/githubrepo"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// DefaultsFile holds the optional repository defaults in a Project root.
const DefaultsFile = "ply-repo-defaults.json"

const (
	gitTimeout    = 30 * time.Second
	githubTimeout = 60 * time.Second
	defaultsLimit = 64 << 10
)

var (
	invalidClasses    = map[string]bool{string(workspace.ErrorProjectInvalidArguments): true, string(workspace.ErrorProjectNotFound): true, string(workspace.ErrorWorkInvalidArguments): true}
	defineConflicts   = map[string]bool{string(workspace.ErrorProjectConflict): true, string(workspace.ErrorProjectStoreConflict): true, string(workspace.ErrorWorkIdentityConflict): true, string(workspace.ErrorWorkProjectConflict): true}
	metadataConflicts = map[string]bool{string(workspace.ErrorProjectConflict): true, string(workspace.ErrorWorkProjectConflict): true}
)

// hooks is a private, per-call test seam; production has none.
type hooks struct{ afterMetadataRead func() }

// New returns the production dependencies of repoestablish.Run.
func New(runner process.Runner) repoestablish.Dependencies {
	return newDependencies(runner, hooks{})
}

func newDependencies(runner process.Runner, h hooks) repoestablish.Dependencies {
	return repoestablish.Dependencies{
		Prepare: prepare,
		Check:   check,
		Observe: func(ctx context.Context, dir string) (gitrepo.Repo, error) {
			ctx, cancel := context.WithTimeout(ctx, gitTimeout)
			defer cancel()
			repo, err := gitrepo.Observe(ctx, runner, dir)
			return repo, capability(err, gitrepo.ErrInvalid, gitrepo.ErrConflict)
		},
		Init: func(ctx context.Context, dir, branch string) (bool, error) {
			ctx, cancel := context.WithTimeout(ctx, gitTimeout)
			defer cancel()
			_, changed, err := gitrepo.EnsureInit(ctx, runner, dir, branch)
			return changed, capability(err, gitrepo.ErrInvalid, gitrepo.ErrConflict)
		},
		GitHub: func(ctx context.Context, owner, name string) (string, bool, error) {
			ctx, cancel := context.WithTimeout(ctx, githubTimeout)
			defer cancel()
			repo, changed, err := githubrepo.EnsurePrivate(ctx, runner, owner, name)
			return repo.SSHURL, changed, capability(err, githubrepo.ErrInvalid, githubrepo.ErrConflict)
		},
		Remote: func(ctx context.Context, dir, name, url string) (bool, error) {
			ctx, cancel := context.WithTimeout(ctx, gitTimeout)
			defer cancel()
			changed, err := gitrepo.EnsureRemote(ctx, runner, dir, name, url)
			return changed, capability(err, gitrepo.ErrInvalid, gitrepo.ErrConflict)
		},
		Define: define,
		SetWrapper: func(ctx context.Context, r repoestablish.Resolved) (bool, error) {
			return setWrapper(ctx, r, h)
		},
	}
}

// classified carries one repoestablish sentinel and the original error.
type classified struct{ kind, err error }

func (c classified) Error() string {
	if message := c.err.Error(); strings.HasPrefix(message, c.kind.Error()+":") {
		return message
	}
	return c.kind.Error() + ": " + c.err.Error()
}
func (c classified) Unwrap() []error { return []error{c.kind, c.err} }

func fail(kind error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", kind, fmt.Sprintf(format, args...))
}

// capability keeps a capability's invalid and conflict sentinels, which it
// returns only before its effect or after a readback showing none; anything
// else is unknown.
func capability(err, invalid, conflict error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, invalid):
		return classified{repoestablish.ErrInvalid, err}
	case errors.Is(err, conflict):
		return classified{repoestablish.ErrConflict, err}
	}
	return classified{repoestablish.ErrUnknown, err}
}

// classify maps the class of a workspace error: invalid classes and the given
// conflict classes arise only before writing; any other error is otherwise.
func classify(err error, conflicts map[string]bool, otherwise error) error {
	if err == nil {
		return nil
	}
	class := ""
	var project *workspace.ProjectError
	var work *workspace.WorkItemError
	if errors.As(err, &project) {
		class = string(project.Class)
	} else if errors.As(err, &work) {
		class = string(work.Class)
	}
	switch {
	case invalidClasses[class]:
		return classified{repoestablish.ErrInvalid, err}
	case conflicts[class]:
		return classified{repoestablish.ErrConflict, err}
	}
	return classified{otherwise, err}
}

// fixedCwdFiles answers Getwd with the explicit workspace, so that workspace
// operations never see the process working directory.
type fixedCwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f fixedCwdFiles) Getwd() (string, error) { return f.cwd, nil }

func fixedDependencies(root string) workspace.Dependencies {
	d := workspace.SystemDependencies()
	d.Files = fixedCwdFiles{d.Files, root}
	return d
}

// projectRoots resolves each defined Project's wrapper against the workspace
// root like the definitions reader does.
func projectRoots(root string) (map[string]string, error) {
	defs, _, err := workspacedef.Read(root)
	switch {
	case err != nil:
		return nil, fail(repoestablish.ErrInvalid, "read %s: %v", filepath.Join(root, workspacedef.File), err)
	case defs == nil:
		return nil, fail(repoestablish.ErrInvalid, "%s does not exist", filepath.Join(root, workspacedef.File))
	}
	roots := map[string]string{}
	for _, p := range defs.Projects {
		path := filepath.FromSlash(p.Wrapper)
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		roots[p.ID] = filepath.Clean(path)
	}
	return roots, nil
}

// ProjectAtRoot returns the one Project whose root is the physical cwd. No
// match and several matches are errors that ask for --project.
func ProjectAtRoot(workspaceRoot, cwd string) (string, error) {
	roots, err := projectRoots(workspaceRoot)
	if err != nil {
		return "", err
	}
	matches := []string{}
	for id, root := range roots {
		if root == cwd {
			matches = append(matches, id)
		}
	}
	if len(matches) != 1 {
		return "", fail(repoestablish.ErrInvalid, "%d defined Projects have root %s; run from a Project root or use --project", len(matches), cwd)
	}
	return matches[0], nil
}

func prepare(ctx context.Context, req repoestablish.Request) (repoestablish.ProjectSetup, error) {
	roots, err := projectRoots(req.Workspace)
	if err != nil {
		return repoestablish.ProjectSetup{}, err
	}
	root, ok := roots[req.Project]
	if !ok {
		return repoestablish.ProjectSetup{}, fail(repoestablish.ErrInvalid, "Project %s is not defined in %s", req.Project, workspacedef.File)
	}
	if physical, err := filepath.EvalSymlinks(root); err != nil || physical != root {
		return repoestablish.ProjectSetup{}, fail(repoestablish.ErrInvalid, "Project root %s must be an existing physical directory", root)
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		return repoestablish.ProjectSetup{}, fail(repoestablish.ErrInvalid, "Project root %s must be a directory", root)
	}
	return readDefaults(root)
}

// PrepareProject resolves the physical root of a defined Project and reads
// its defaults exactly as repo create does.
func PrepareProject(workspaceRoot, project string) (repoestablish.ProjectSetup, error) {
	return prepare(context.Background(), repoestablish.Request{Workspace: workspaceRoot, Project: project})
}

// readDefaults reads the optional defaults file strictly: absent is empty
// defaults; anything present but not valid is invalid, never a fallback.
func readDefaults(root string) (repoestablish.ProjectSetup, error) {
	setup := repoestablish.ProjectSetup{Root: root}
	path := filepath.Join(root, DefaultsFile)
	invalid := func(reason string) (repoestablish.ProjectSetup, error) {
		return repoestablish.ProjectSetup{}, fail(repoestablish.ErrInvalid, "Project defaults %s %s", path, reason)
	}
	info, err := os.Lstat(path)
	switch {
	case os.IsNotExist(err):
		return setup, nil
	case err != nil:
		return invalid(fmt.Sprintf("cannot be inspected: %v", err))
	case !info.Mode().IsRegular():
		return invalid("must be a regular file, not a symlink")
	}
	file, err := os.Open(path)
	if err != nil {
		return invalid(fmt.Sprintf("cannot be read: %v", err))
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return invalid("changed while it was opened")
	}
	raw, err := io.ReadAll(io.LimitReader(file, defaultsLimit+1))
	switch {
	case err != nil:
		return invalid(fmt.Sprintf("cannot be read: %v", err))
	case len(raw) > defaultsLimit:
		return invalid("exceeds 64 KiB")
	}
	// Raw string fields keep an explicit null apart from an omitted key.
	var defaults struct {
		FormatVersion *int            `json:"format_version"`
		GitHubOwner   json.RawMessage `json:"github_owner"`
		Leaf          json.RawMessage `json:"leaf"`
		Branch        json.RawMessage `json:"branch"`
		Remote        json.RawMessage `json:"remote"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&defaults); err != nil {
		return invalid(fmt.Sprintf("is not valid: %v", err))
	}
	if err := decoder.Decode(&json.RawMessage{}); err != io.EOF {
		return invalid("must contain exactly one JSON value")
	}
	if defaults.FormatVersion == nil || *defaults.FormatVersion != 1 {
		return invalid("must have format_version 1")
	}
	for _, field := range []struct {
		name   string
		raw    json.RawMessage
		target *string
	}{{"github_owner", defaults.GitHubOwner, &setup.GitHubOwner}, {"leaf", defaults.Leaf, &setup.Leaf},
		{"branch", defaults.Branch, &setup.Branch}, {"remote", defaults.Remote, &setup.Remote}} {
		if len(field.raw) == 0 {
			continue
		}
		var value *string
		if err := json.Unmarshal(field.raw, &value); err != nil || value == nil || *value == "" {
			return invalid(field.name + " must be a nonempty string when present")
		}
		*field.target = *value
	}
	return setup, nil
}

// storedPath is the repository path as the definitions store it.
func storedPath(r repoestablish.Resolved) (string, error) {
	relative, err := filepath.Rel(r.Workspace, r.RepoDir)
	if err != nil {
		return "", fail(repoestablish.ErrInvalid, "repository %s is not inside workspace %s", r.RepoDir, r.Workspace)
	}
	return filepath.ToSlash(relative), nil
}

// check reads tools, identifiers, definitions and metadata without effect.
func check(ctx context.Context, r repoestablish.Resolved) error {
	for _, tool := range []string{"git", "gh"} {
		if _, err := exec.LookPath(tool); err != nil {
			return fail(repoestablish.ErrInvalid, "%s is not available: %v", tool, err)
		}
	}
	d := fixedDependencies(r.Workspace)
	if _, err := workspace.ObserveRoot(d, r.Workspace); err != nil {
		return classify(err, defineConflicts, repoestablish.ErrInvalid)
	}
	if _, err := workspace.ParseProjectID(r.Project); err != nil {
		return classified{repoestablish.ErrInvalid, err}
	}
	if _, err := workspace.ParseRepoID(r.RepoName); err != nil {
		return classified{repoestablish.ErrInvalid, err}
	}
	if !githubrepo.Valid(r.GitHubOwner, r.RepoName) {
		return fail(repoestablish.ErrInvalid, "%s/%s is not a valid GitHub owner and repository name", r.GitHubOwner, r.RepoName)
	}
	if !gitrepo.ValidRemoteName(r.Remote) {
		return fail(repoestablish.ErrInvalid, "%q is not a valid remote name", r.Remote)
	}
	input := workspace.ProjectMetadataInput{Operation: workspace.ProjectRepoWrapperSet, ProjectID: workspace.ProjectID(r.Project),
		RepoID: workspace.RepoID(r.RepoName), Locator: r.Wrapper, ActorClaim: r.Actor}
	if err := input.Validate(); err != nil {
		return classified{repoestablish.ErrInvalid, err}
	}
	stored, err := storedPath(r)
	if err != nil {
		return err
	}
	defs, _, err := workspacedef.Read(r.Workspace)
	if err != nil || defs == nil {
		return fail(repoestablish.ErrInvalid, "read %s: %v", workspacedef.File, err)
	}
	var defined *workspacedef.Error
	if _, err := workspacedef.AddRepo(defs, r.Project, r.RepoName, stored); errors.As(err, &defined) && defined.Kind == workspacedef.KindConflict {
		return classified{repoestablish.ErrConflict, err}
	} else if err != nil {
		return classified{repoestablish.ErrInvalid, err}
	}
	if err := workspace.CheckDefinitions(d, r.Workspace); err != nil {
		return classify(err, defineConflicts, repoestablish.ErrInvalid)
	}
	snapshot, err := workspace.ReadProjectMetadata(d, r.Workspace)
	if err != nil {
		return classified{repoestablish.ErrConflict, err}
	}
	if wrapper := snapshot.RepositoryWrapper(workspace.RepoID(r.RepoName)); wrapper != nil && *wrapper != r.Wrapper {
		return fail(repoestablish.ErrConflict, "repository %s already has wrapper %s, not %s", r.RepoName, *wrapper, r.Wrapper)
	}
	return nil
}

// notStarted reports a context that ended before a workspace call; once
// started, such a call runs to completion.
func notStarted(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return classified{repoestablish.ErrInvalid, err}
	}
	return nil
}

func define(ctx context.Context, r repoestablish.Resolved) (bool, error) {
	if err := notStarted(ctx); err != nil {
		return false, err
	}
	stored, err := storedPath(r)
	if err != nil {
		return false, err
	}
	result, err := workspace.AddDefinitionRepo(fixedDependencies(r.Workspace), r.Project, r.RepoName, stored)
	return result.Changed, classify(err, defineConflicts, repoestablish.ErrUnknown)
}

// setWrapper reads the wrapper of the repository and sets it only when it is
// absent, with the read revision as expected revision.
func setWrapper(ctx context.Context, r repoestablish.Resolved, h hooks) (bool, error) {
	if err := notStarted(ctx); err != nil {
		return false, err
	}
	d := fixedDependencies(r.Workspace)
	snapshot, err := workspace.ReadProjectMetadata(d, r.Workspace)
	if err != nil {
		return false, classified{repoestablish.ErrConflict, err}
	}
	if wrapper := snapshot.RepositoryWrapper(workspace.RepoID(r.RepoName)); wrapper != nil {
		if *wrapper != r.Wrapper {
			return false, fail(repoestablish.ErrConflict, "repository %s already has wrapper %s, not %s", r.RepoName, *wrapper, r.Wrapper)
		}
		return false, nil
	}
	if h.afterMetadataRead != nil {
		h.afterMetadataRead()
	}
	revision := snapshot.Revision
	result, err := workspace.UpdateProjectMetadata(d, workspace.ProjectMetadataInput{Operation: workspace.ProjectRepoWrapperSet,
		ProjectID: workspace.ProjectID(r.Project), RepoID: workspace.RepoID(r.RepoName), Locator: r.Wrapper, ActorClaim: r.Actor,
		ExpectedRevision: &revision})
	return result.Changed, classify(err, metadataConflicts, repoestablish.ErrUnknown)
}
