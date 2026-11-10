// Package system binds reporemove to Git and the Ply workspace. It reads no
// working directory and no environment variable of its own: every workspace
// operation runs at the explicit workspace root. Git only reads, through
// gitrepo.Read; the remote is contacted by one ls-remote only when local
// history exists, and never while projects.lock is held. Where the platform
// reports no device of a file (Windows), the wrapper check fails closed, so
// removal is not supported there. Not detected: an unregistered repository
// elsewhere that borrows this object database through alternates.
package system

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
	"github.com/devdimensionlab/plybuild/internal/gitworktree"
	establish "github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/reporemove"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

const (
	gitTimeout     = 30 * time.Second
	networkTimeout = 60 * time.Second
	quarantineName = ".ply-repo-remove-"
)

var (
	invalidClasses    = map[string]bool{string(workspace.ErrorProjectInvalidArguments): true, string(workspace.ErrorProjectNotFound): true, string(workspace.ErrorWorkInvalidArguments): true}
	checkConflicts    = map[string]bool{string(workspace.ErrorProjectConflict): true}
	defineConflicts   = map[string]bool{string(workspace.ErrorProjectConflict): true, string(workspace.ErrorProjectStoreConflict): true, string(workspace.ErrorWorkIdentityConflict): true, string(workspace.ErrorWorkProjectConflict): true}
	metadataConflicts = map[string]bool{string(workspace.ErrorProjectConflict): true, string(workspace.ErrorWorkProjectConflict): true}
	objectID          = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)
	stateFiles        = []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "BISECT_LOG", "rebase-merge", "rebase-apply", "sequencer"}
	// Object stores outside Git's object database that the history proof does not see.
	foreignStores = []string{"lfs", "annex"}
)

// hooks is a private, per-call test seam; production has none.
type hooks struct {
	beforeWrapper func()
	removeAll     func(string) error
}

type adapter struct {
	runner process.Runner
	hooks  hooks
}

// New returns the production dependencies of reporemove.Run.
func New(runner process.Runner) reporemove.Dependencies {
	return newDependencies(runner, hooks{})
}

func newDependencies(runner process.Runner, h hooks) reporemove.Dependencies {
	a := adapter{runner, h}
	return reporemove.Dependencies{Locate: locate, Inspect: a.inspect, ClearWrapper: clearWrapper, Undefine: undefine, RemoveWrapper: a.removeWrapper}
}

// classified carries one reporemove sentinel and the original error.
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

// classify maps the class of a workspace error: invalid classes and the given
// conflict classes arise only without effect; any other error is otherwise.
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
		return classified{reporemove.ErrInvalid, err}
	case conflicts[class]:
		return classified{reporemove.ErrConflict, err}
	}
	return classified{otherwise, err}
}

func atRoot(root string) (workspace.Dependencies, error) {
	d, err := workspace.AtRoot(workspace.SystemDependencies(), root)
	if err != nil {
		return d, classified{reporemove.ErrInvalid, err}
	}
	return d, nil
}

// locate resolves the Project like repo create and reads the definitions
// member, the wrapper's presence and leftover quarantines.
func locate(ctx context.Context, req reporemove.Request) (reporemove.Target, error) {
	setup, err := establish.PrepareProject(req.Workspace, req.Project)
	if err != nil {
		return reporemove.Target{}, err
	}
	target := reporemove.Target{ProjectRoot: setup.Root, Remote: setup.Remote}
	defs, _, err := workspacedef.Read(req.Workspace)
	if err != nil || defs == nil {
		return target, fail(reporemove.ErrInvalid, "read %s: %v", workspacedef.File, err)
	}
	for _, p := range defs.Projects {
		for _, r := range p.Repos {
			if p.ID == req.Project && r.ID == req.RepoName {
				target.Path = r.Path
			}
		}
	}
	wrapper := filepath.Join(setup.Root, req.RepoName)
	if _, err := os.Lstat(wrapper); err == nil {
		target.WrapperPresent = true
	} else if !errors.Is(err, fs.ErrNotExist) {
		return target, fail(reporemove.ErrInvalid, "inspect wrapper %s: %v", wrapper, err)
	}
	target.Quarantines, err = quarantines(setup.Root, req.RepoName)
	return target, err
}

// quarantines lists <root>/.ply-repo-remove-<repo>-* literally; none is ever
// deleted by pattern.
func quarantines(root, repo string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fail(reporemove.ErrInvalid, "read Project root %s: %v", root, err)
	}
	var found []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), quarantineName+repo+"-") {
			found = append(found, filepath.Join(root, entry.Name()))
		}
	}
	return found, nil
}

// local is the read-only observation of the repository and its wrapper.
type local struct {
	head     string
	refs     map[string]string // ref name -> object ID
	listing  string
	wrapper  os.FileInfo
	gitFiles []string
}

// digest covers what apply must find unchanged in the repository and wrapper.
func (l local) digest() string {
	refs := []canonicaljson.Value{}
	for _, name := range l.refNames() {
		refs = append(refs, canonicaljson.Object{{Name: "ref", Value: name}, {Name: "oid", Value: l.refs[name]}})
	}
	return digestOf(canonicaljson.Object{{Name: "head", Value: l.head}, {Name: "refs", Value: refs}, {Name: "listing", Value: l.listing}})
}

func (l local) refNames() []string {
	names := make([]string, 0, len(l.refs))
	for name := range l.refs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func digestOf(value canonicaljson.Value) string {
	raw, err := canonicaljson.Marshal(value)
	if err != nil {
		panic(err) // only strings, integers, arrays and objects are marshaled
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (a adapter) read(ctx context.Context, timeout time.Duration, dir string, args ...string) ([]byte, string, int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return gitrepo.Read(ctx, a.runner, dir, args...)
}

// readOK runs a read whose only success is exit 0; anything else is invalid.
func (a adapter) readOK(ctx context.Context, dir string, args ...string) ([]byte, error) {
	stdout, stderr, code, err := a.read(ctx, gitTimeout, dir, args...)
	if err != nil {
		return nil, fail(reporemove.ErrInvalid, "git %s in %s: exit %d: %s", args[0], dir, code, strings.TrimSpace(stderr))
	}
	return stdout, nil
}

func (a adapter) inspect(ctx context.Context, r reporemove.Resolved) (reporemove.Plan, error) {
	plan := reporemove.Plan{}
	d, err := atRoot(r.Workspace)
	if err != nil {
		return plan, err
	}
	pass := func(name, detail string) {
		plan.Checks = append(plan.Checks, reporemove.Check{Name: name, Detail: detail})
	}

	if plan.DefinitionsDigest, plan.MetadataRevision, plan.MetadataCleared, err = provenance(ctx, a, d, r); err != nil {
		return plan, err
	}
	pass("provenance", "")
	if err := workspace.CheckRepoRemoval(d, r.Workspace, r.Project, r.RepoName); err != nil {
		return plan, classify(err, checkConflicts, reporemove.ErrInvalid)
	}
	pass("references", "")
	observed, err := a.observe(ctx, r)
	if err != nil {
		return plan, err
	}
	pass("worktree", "")
	plan.Verdict, plan.RemoteURL, err = a.history(ctx, r, observed)
	if err != nil {
		return plan, err
	}
	if plan.Verdict == reporemove.VerdictUnborn {
		pass("history", "unborn, no objects; remote not contacted")
	} else {
		pass("history", "every local commit and ref is on "+plan.RemoteURL)
	}
	pass("wrapper", "")
	plan.GitFiles, plan.LocalDigest = observed.gitFiles, observed.digest()
	refs := []canonicaljson.Value{}
	for _, name := range observed.refNames() {
		refs = append(refs, canonicaljson.Object{{Name: "ref", Value: name}, {Name: "oid", Value: observed.refs[name]}})
	}
	plan.Digest = digestOf(canonicaljson.Object{
		{Name: "workspace", Value: r.Workspace}, {Name: "project", Value: r.Project}, {Name: "repository", Value: r.RepoName},
		{Name: "project_root", Value: r.ProjectRoot}, {Name: "wrapper", Value: r.Wrapper}, {Name: "leaf", Value: r.Leaf},
		{Name: "remote", Value: r.Remote}, {Name: "remote_url", Value: plan.RemoteURL},
		{Name: "definitions", Value: plan.DefinitionsDigest}, {Name: "metadata_revision", Value: int64(plan.MetadataRevision)},
		{Name: "head", Value: observed.head}, {Name: "refs", Value: refs}, {Name: "verdict", Value: plan.Verdict},
		{Name: "listing", Value: observed.listing}})
	return plan, nil
}

// provenance requires the layout repo create produces and returns the
// definitions digest, the metadata revision and whether the wrapper metadata
// is already cleared by an earlier run.
func provenance(ctx context.Context, a adapter, d workspace.Dependencies, r reporemove.Resolved) (string, int, bool, error) {
	if _, err := workspace.ParseRepoID(r.RepoName); err != nil {
		return "", 0, false, classified{reporemove.ErrInvalid, err}
	}
	defs, digest, err := workspacedef.Read(r.Workspace)
	if err != nil || defs == nil {
		return "", 0, false, fail(reporemove.ErrInvalid, "read %s: %v", workspacedef.File, err)
	}
	members := 0
	for _, p := range defs.Projects {
		for _, repo := range p.Repos {
			if repo.ID == r.RepoName {
				members++
				if p.ID != r.Project || repo.Path != r.Path {
					return "", 0, false, fail(reporemove.ErrConflict, "repository %s is also defined in Project %s at %s", r.RepoName, p.ID, repo.Path)
				}
			}
		}
	}
	if members != 1 {
		return "", 0, false, fail(reporemove.ErrConflict, "repository %s is not defined exactly once", r.RepoName)
	}
	if err := workspace.CheckDefinitions(d, r.Workspace); err != nil {
		return "", 0, false, classify(err, checkConflicts, reporemove.ErrInvalid)
	}
	snapshot, err := workspace.ReadProjectMetadata(d, r.Workspace)
	if err != nil {
		return "", 0, false, classified{reporemove.ErrInvalid, err}
	}
	cleared, err := wrapperEvents(snapshot, r)
	if err != nil {
		return "", 0, false, err
	}
	if err := physicalLayout(r); err != nil {
		return "", 0, false, err
	}
	identity, err := gitrepo.Identify(ctx, a.runner, r.RepoDir)
	if err != nil || !identity.Worktree || identity.Toplevel != r.RepoDir || identity.CommonDir != filepath.Join(r.RepoDir, ".git") {
		return "", 0, false, fail(reporemove.ErrInvalid, "%s must be a Git worktree root with common dir %s/.git: %+v %v", r.RepoDir, r.RepoDir, identity, err)
	}
	return digest, snapshot.Revision, cleared, nil
}

// wrapperEvents requires the latest wrapper event of the repository to be the
// set of the wrapper, or a clear directly preceded by it.
func wrapperEvents(snapshot workspace.ProjectMetadataSnapshot, r reporemove.Resolved) (bool, error) {
	var events []workspace.ProjectMetadataEvent
	for _, event := range snapshot.Events {
		if event.RepoID != nil && string(*event.RepoID) == r.RepoName &&
			(event.Operation == workspace.ProjectRepoWrapperSet || event.Operation == workspace.ProjectRepoWrapperClear) {
			events = append(events, event)
		}
	}
	set := func(event workspace.ProjectMetadataEvent) bool {
		return event.Operation == workspace.ProjectRepoWrapperSet && event.Locator != nil && *event.Locator == r.Wrapper
	}
	n := len(events)
	switch {
	case n > 0 && set(events[n-1]):
		return false, nil
	case n > 1 && events[n-1].Operation == workspace.ProjectRepoWrapperClear && set(events[n-2]):
		return true, nil
	}
	return false, fail(reporemove.ErrConflict, "project metadata does not record wrapper %s for repository %s as repo create sets it", r.Wrapper, r.RepoName)
}

// physicalLayout requires a physical wrapper directory directly in the
// Project root that holds neither the workspace root nor .ply, and a physical
// repository directory whose .git is a directory.
func physicalLayout(r reporemove.Resolved) error {
	if filepath.Dir(r.Wrapper) != r.ProjectRoot || r.Wrapper == r.ProjectRoot {
		return fail(reporemove.ErrInvalid, "wrapper %s must be a direct child of Project root %s", r.Wrapper, r.ProjectRoot)
	}
	if rel, err := filepath.Rel(r.Wrapper, r.Workspace); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fail(reporemove.ErrInvalid, "wrapper %s contains the workspace root %s", r.Wrapper, r.Workspace)
	}
	for _, dir := range []string{r.Wrapper, r.RepoDir, filepath.Join(r.RepoDir, ".git")} {
		info, err := os.Lstat(dir)
		switch {
		case err != nil:
			return fail(reporemove.ErrInvalid, "inspect %s: %v", dir, err)
		case info.Mode()&fs.ModeSymlink != 0:
			return fail(reporemove.ErrInvalid, "%s is a symlink", dir)
		case !info.IsDir():
			return fail(reporemove.ErrInvalid, "%s is not a directory", dir)
		}
	}
	if physical, err := filepath.EvalSymlinks(r.Wrapper); err != nil || physical != r.Wrapper {
		return fail(reporemove.ErrInvalid, "wrapper %s is not a physical path", r.Wrapper)
	}
	return nil
}

// observe runs checks 3, 4, 5 and 7, all local: one unlocked worktree, no
// local state, no unreachable object and a wrapper holding only the leaf.
func (a adapter) observe(ctx context.Context, r reporemove.Resolved) (local, error) {
	var l local
	gitDir := filepath.Join(r.RepoDir, ".git")
	raw, err := a.readOK(ctx, r.RepoDir, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return l, err
	}
	worktrees, err := gitworktree.ParseInventory(raw)
	if err != nil {
		return l, fail(reporemove.ErrInvalid, "worktree list of %s: %v", r.RepoDir, err)
	}
	if len(worktrees) != 1 || worktrees[0].Locator != r.RepoDir || worktrees[0].Locked || worktrees[0].Bare {
		return l, fail(reporemove.ErrConflict, "%s must be the one unlocked worktree of its repository, found %d worktrees", r.RepoDir, len(worktrees))
	}
	if entries, err := os.ReadDir(filepath.Join(gitDir, "worktrees")); err == nil && len(entries) > 0 {
		return l, fail(reporemove.ErrConflict, "%s/worktrees records other worktrees", gitDir)
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return l, fail(reporemove.ErrInvalid, "read %s/worktrees: %v", gitDir, err)
	}
	if _, err := os.Lstat(filepath.Join(gitDir, "modules")); !errors.Is(err, fs.ErrNotExist) {
		return l, fail(reporemove.ErrConflict, "%s/modules holds submodule repositories", gitDir)
	}

	raw, err = a.readOK(ctx, r.RepoDir, "status", "--porcelain=v2", "-z", "--branch", "-uall", "--ignored=matching", "--ignore-submodules=none")
	if err != nil {
		return l, err
	}
	for _, record := range bytes.Split(raw, []byte{0}) {
		if len(record) > 0 && !bytes.HasPrefix(record, []byte("# ")) {
			return l, fail(reporemove.ErrConflict, "local changes in %s: tracked, untracked or ignored files exist", r.RepoDir)
		}
	}
	// A gitlink's history lives in another repository that fsck and rev-list
	// of this one do not see.
	raw, err = a.readOK(ctx, r.RepoDir, "ls-files", "--stage", "-z")
	if err != nil {
		return l, err
	}
	for _, record := range bytes.Split(raw, []byte{0}) {
		mode, _, ok := bytes.Cut(record, []byte(" "))
		switch {
		case len(record) == 0:
		case !ok || len(mode) != 6:
			return l, fail(reporemove.ErrInvalid, "ls-files in %s: malformed record %q", r.RepoDir, record)
		case string(mode) == "160000":
			return l, fail(reporemove.ErrConflict, "%s tracks a nested repository (gitlink) whose history is not checked: %q", r.RepoDir, record)
		}
	}
	// Skip-worktree (S) and assume-unchanged (lowercase) entries hide local
	// changes from status; any tag but H is refused.
	raw, err = a.readOK(ctx, r.RepoDir, "ls-files", "-v", "-z")
	if err != nil {
		return l, err
	}
	for _, record := range bytes.Split(raw, []byte{0}) {
		if len(record) > 0 && !bytes.HasPrefix(record, []byte("H ")) {
			return l, fail(reporemove.ErrConflict, "index entry %q in %s is skip-worktree, assume-unchanged or otherwise flagged; its local content is not checked", record, r.RepoDir)
		}
	}
	for _, name := range foreignStores {
		if _, err := os.Lstat(filepath.Join(gitDir, name)); !errors.Is(err, fs.ErrNotExist) {
			return l, fail(reporemove.ErrConflict, "%s/%s holds a %s object store whose local content is not checked", gitDir, name, name)
		}
	}
	for _, name := range stateFiles {
		if _, err := os.Lstat(filepath.Join(gitDir, name)); !errors.Is(err, fs.ErrNotExist) {
			return l, fail(reporemove.ErrConflict, "%s/%s exists: an operation is in progress", gitDir, name)
		}
	}
	if l.refs, err = a.refs(ctx, r.RepoDir); err != nil {
		return l, err
	}
	if _, ok := l.refs["refs/stash"]; ok {
		return l, fail(reporemove.ErrConflict, "refs/stash exists in %s", r.RepoDir)
	}
	if l.head, err = a.head(ctx, r.RepoDir); err != nil {
		return l, err
	}
	if err := a.fsck(ctx, r.RepoDir); err != nil {
		return l, err
	}
	if l.wrapper, l.listing, err = walk(r, deviceOf); err != nil {
		return l, err
	}
	l.gitFiles = gitFiles(gitDir)
	return l, nil
}

func (a adapter) refs(ctx context.Context, dir string) (map[string]string, error) {
	raw, err := a.readOK(ctx, dir, "for-each-ref", "--format=%(objectname) %(refname)")
	if err != nil {
		return nil, err
	}
	refs := map[string]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		if line == "" {
			continue
		}
		oid, name, ok := strings.Cut(line, " ")
		if !ok || !objectID.MatchString(oid) || name == "" {
			return nil, fail(reporemove.ErrInvalid, "for-each-ref in %s: malformed line %q", dir, line)
		}
		refs[name] = oid
	}
	return refs, nil
}

// head is the symbolic ref, or "detached", and the commit, or "unborn".
func (a adapter) head(ctx context.Context, dir string) (string, error) {
	symbolic, stderr, code, err := a.read(ctx, gitTimeout, dir, "symbolic-ref", "-q", "HEAD")
	name := strings.TrimSpace(string(symbolic))
	if err != nil && code == 1 {
		name = "detached"
	} else if err != nil {
		return "", fail(reporemove.ErrInvalid, "symbolic-ref HEAD in %s: exit %d: %s", dir, code, strings.TrimSpace(stderr))
	}
	oid, stderr, code, err := a.read(ctx, gitTimeout, dir, "rev-parse", "-q", "--verify", "HEAD")
	commit := strings.TrimSpace(string(oid))
	if err != nil && code == 1 && commit == "" {
		commit = "unborn"
	} else if err != nil || !objectID.MatchString(commit) {
		return "", fail(reporemove.ErrInvalid, "rev-parse HEAD in %s: exit %d: %s", dir, code, strings.TrimSpace(stderr))
	}
	return name + "@" + commit, nil
}

// fsck rejects every unreachable or dangling object; any output it does not
// know makes the check invalid, never clean.
func (a adapter) fsck(ctx context.Context, dir string) error {
	stdout, stderr, code, err := a.read(ctx, networkTimeout, dir, "fsck", "--connectivity-only", "--unreachable", "--no-progress")
	if err != nil {
		return fail(reporemove.ErrInvalid, "fsck %s: exit %d: %s", dir, code, strings.TrimSpace(stderr))
	}
	for _, line := range strings.SplitAfter(stderr, "\n") {
		if line != "" && (!strings.HasPrefix(line, "notice: ") || !strings.HasSuffix(line, "\n")) {
			return fail(reporemove.ErrInvalid, "fsck %s: unexpected diagnostic %q", dir, strings.TrimSpace(line))
		}
	}
	var unreachable []string
	for _, line := range strings.SplitAfter(string(stdout), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		valid := strings.HasSuffix(line, "\n") && len(fields) == 3 && (fields[0] == "unreachable" || fields[0] == "dangling") &&
			(fields[1] == "commit" || fields[1] == "tree" || fields[1] == "blob" || fields[1] == "tag") && objectID.MatchString(fields[2])
		if !valid {
			return fail(reporemove.ErrInvalid, "fsck %s: unexpected output %q", dir, strings.TrimSpace(line))
		}
		unreachable = append(unreachable, fields[2])
	}
	if len(unreachable) > 0 {
		if len(unreachable) > 5 {
			unreachable = unreachable[:5]
		}
		return fail(reporemove.ErrConflict, "unreferenced local objects in %s (%s); inspect them before removing", dir, strings.Join(unreachable, ", "))
	}
	return nil
}

// walk lists the wrapper without following links: its only entry is the
// leaf, everything is on one device and no .ply exists inside.
func walk(r reporemove.Resolved, deviceOf func(os.FileInfo) (uint64, bool)) (os.FileInfo, string, error) {
	top, err := os.Lstat(r.Wrapper)
	if err != nil || !top.IsDir() || top.Mode()&fs.ModeSymlink != 0 {
		return nil, "", fail(reporemove.ErrInvalid, "wrapper %s must be a physical directory: %v", r.Wrapper, err)
	}
	device, ok := deviceOf(top)
	if !ok {
		return nil, "", fail(reporemove.ErrInvalid, "the device of wrapper %s is unknown on this platform; removal is not supported", r.Wrapper)
	}
	entries := []canonicaljson.Value{}
	err = filepath.WalkDir(r.Wrapper, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fail(reporemove.ErrInvalid, "walk %s: %v", path, err)
		}
		rel, _ := filepath.Rel(r.Wrapper, path)
		info, err := entry.Info()
		var entryDevice uint64
		if err == nil {
			if entryDevice, ok = deviceOf(info); !ok {
				return fail(reporemove.ErrInvalid, "the device of %s is unknown on this platform", path)
			}
		}
		switch {
		case err != nil:
			return fail(reporemove.ErrInvalid, "inspect %s: %v", path, err)
		case filepath.Dir(rel) == "." && rel != "." && rel != r.Leaf:
			return fail(reporemove.ErrConflict, "wrapper %s holds extra entry %s besides %s", r.Wrapper, rel, r.Leaf)
		case entry.Name() == workspace.MarkerDirectory:
			return fail(reporemove.ErrConflict, "wrapper %s holds %s", r.Wrapper, path)
		case entry.Name() == ".git" && rel != filepath.Join(r.Leaf, ".git"):
			return fail(reporemove.ErrConflict, "%s is a nested Git repository whose history is not checked", path)
		case entryDevice != device:
			return fail(reporemove.ErrConflict, "%s is on another device than wrapper %s", path, r.Wrapper)
		}
		entries = append(entries, canonicaljson.Object{{Name: "path", Value: filepath.ToSlash(rel)}, {Name: "type", Value: info.Mode().Type().String()},
			{Name: "size", Value: info.Size()}, {Name: "mtime", Value: info.ModTime().UnixNano()}})
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return top, digestOf(entries), nil
}

// gitFiles names the local configuration deleted with .git.
func gitFiles(gitDir string) []string {
	var files []string
	if _, err := os.Lstat(filepath.Join(gitDir, "config")); err == nil {
		files = append(files, "config")
	}
	for _, dir := range []string{"hooks", "info"} {
		entries, _ := os.ReadDir(filepath.Join(gitDir, dir))
		for _, entry := range entries {
			files = append(files, dir+"/"+entry.Name())
		}
	}
	return files
}

// history proves that no local history is lost: unborn without contacting the
// remote, or every local commit and every ref outside refs/heads and
// refs/remotes on the configured remote URL.
func (a adapter) history(ctx context.Context, r reporemove.Resolved, l local) (string, string, error) {
	raw, stderr, code, err := a.read(ctx, gitTimeout, r.RepoDir, "config", "--local", "--get-all", "remote."+r.Remote+".url")
	if err != nil && (code != 1 || len(raw) > 0) {
		return "", "", fail(reporemove.ErrInvalid, "read remote %s URL in %s: exit %d: %s", r.Remote, r.RepoDir, code, strings.TrimSpace(stderr))
	}
	url := ""
	if urls := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n"); len(raw) > 0 && len(urls) == 1 {
		url = urls[0]
	}
	reflog, err := a.readOK(ctx, r.RepoDir, "rev-list", "--all", "--reflog", "--max-count=1")
	if err != nil {
		return "", "", err
	}
	if len(l.refs) == 0 && len(reflog) == 0 {
		return reporemove.VerdictUnborn, url, nil
	}
	unproven := func(format string, args ...interface{}) error {
		return fail(reporemove.ErrConflict, "history unproven: "+format, args...)
	}
	if url == "" {
		return "", "", unproven("remote %s of %s has no single configured URL", r.Remote, r.RepoDir)
	}
	effective, _, _, err := a.read(ctx, gitTimeout, r.RepoDir, "remote", "get-url", r.Remote)
	if line, ok := gitrepo.SingleLine(effective); err != nil || !ok || line != url {
		return "", "", unproven("remote %s URL %s is rewritten or unreadable", r.Remote, url)
	}
	advertised, stderr, code, err := a.read(ctx, networkTimeout, r.RepoDir, "ls-remote", "--", url)
	if err != nil {
		return "", "", unproven("ls-remote %s: exit %d: %s", url, code, strings.TrimSpace(stderr))
	}
	remote := map[string]string{}
	oids := map[string]bool{}
	for _, line := range strings.Split(string(advertised), "\n") {
		if line == "" {
			continue
		}
		oid, name, ok := strings.Cut(line, "\t")
		if !ok || !objectID.MatchString(oid) || name == "" {
			return "", "", unproven("ls-remote %s: malformed line %q", url, line)
		}
		remote[name], oids[oid] = oid, true
	}
	for _, name := range l.refNames() {
		if !strings.HasPrefix(name, "refs/heads/") && !strings.HasPrefix(name, "refs/remotes/") && remote[name] != l.refs[name] {
			return "", "", fail(reporemove.ErrConflict, "%s is not on the remote %s with the same object", name, url)
		}
	}
	args := []string{"rev-list", "--all", "--reflog", "--max-count=1", "--not"}
	missing := 0
	for oid := range oids {
		_, _, code, err := a.read(ctx, gitTimeout, r.RepoDir, "cat-file", "-e", oid)
		switch {
		case err == nil:
			args = append(args, oid)
		case code == 1:
			missing++
		default:
			return "", "", fail(reporemove.ErrInvalid, "cat-file %s in %s: exit %d", oid, r.RepoDir, code)
		}
	}
	unpublished, err := a.readOK(ctx, r.RepoDir, args...)
	if err != nil {
		return "", "", err
	}
	if len(unpublished) > 0 {
		hint := ""
		if missing > 0 {
			hint = fmt.Sprintf("; %d advertised objects are not local, fetch them and preview again", missing)
		}
		return "", "", fail(reporemove.ErrConflict, "local commit %s is not on the remote %s%s", strings.TrimSpace(string(unpublished)), url, hint)
	}
	return reporemove.VerdictProven, url, nil
}

func clearWrapper(ctx context.Context, r reporemove.Resolved, p reporemove.Plan) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, classified{reporemove.ErrInvalid, err}
	}
	d, err := atRoot(r.Workspace)
	if err != nil {
		return false, err
	}
	result, err := workspace.ClearRepositoryWrapperAt(d, r.Workspace, r.Project, r.RepoName, r.Actor, p.MetadataRevision)
	return result.Changed, classify(err, metadataConflicts, reporemove.ErrUnknown)
}

func undefine(ctx context.Context, r reporemove.Resolved, p reporemove.Plan) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, classified{reporemove.ErrInvalid, err}
	}
	d, err := atRoot(r.Workspace)
	if err != nil {
		return false, err
	}
	result, err := workspace.RemoveDefinitionRepo(d, r.Workspace, r.Project, r.RepoName, r.Path, p.DefinitionsDigest)
	return result.Changed, classify(err, defineConflicts, reporemove.ErrUnknown)
}

// removeWrapper re-proves the repository outside any lock, then, under
// projects.lock and with local reads only, requires the removal to be still
// registered and the repository unchanged, and renames the wrapper to its
// quarantine. Only this invocation deletes the quarantine it renamed.
func (a adapter) removeWrapper(ctx context.Context, r reporemove.Resolved, p reporemove.Plan) (bool, error) {
	if a.hooks.beforeWrapper != nil {
		a.hooks.beforeWrapper()
	}
	d, err := atRoot(r.Workspace)
	if err != nil {
		return false, err
	}
	observed, err := a.observe(ctx, r)
	if err != nil {
		return false, err
	}
	verdict, _, err := a.history(ctx, r, observed)
	if err != nil {
		return false, err
	}
	if verdict != p.Verdict || observed.digest() != p.LocalDigest {
		return false, fail(reporemove.ErrConflict, "repository %s changed since the preview", r.RepoDir)
	}
	quarantine := filepath.Join(r.ProjectRoot, quarantineName+r.RepoName+"-"+strings.TrimPrefix(p.Digest, "sha256:")[:12])
	renamed := false
	err = workspacedef.WithLock(r.Workspace, func() error {
		defs, _, err := workspacedef.Read(r.Workspace)
		if err != nil || defs == nil {
			return fail(reporemove.ErrInvalid, "read %s: %v", workspacedef.File, err)
		}
		for _, project := range defs.Projects {
			for _, repo := range project.Repos {
				if repo.ID == r.RepoName {
					return fail(reporemove.ErrConflict, "repository %s is defined again in Project %s", r.RepoName, project.ID)
				}
			}
		}
		snapshot, err := workspace.ReadProjectMetadata(d, r.Workspace)
		if err != nil {
			return classified{reporemove.ErrInvalid, err}
		}
		if cleared, err := wrapperEvents(snapshot, r); err != nil || !cleared {
			return fail(reporemove.ErrConflict, "wrapper metadata of repository %s is no longer the clear of this removal", r.RepoName)
		}
		if err := workspace.CheckRepoRemoval(d, r.Workspace, r.Project, r.RepoName); err != nil {
			return classify(err, checkConflicts, reporemove.ErrInvalid)
		}
		if info, err := os.Lstat(r.Wrapper); err != nil || !os.SameFile(info, observed.wrapper) {
			return fail(reporemove.ErrConflict, "wrapper %s is no longer the observed directory", r.Wrapper)
		}
		again, err := a.observe(ctx, r)
		if err != nil {
			return err
		}
		if again.digest() != observed.digest() {
			return fail(reporemove.ErrConflict, "repository %s changed while it was checked", r.RepoDir)
		}
		if found, err := quarantines(r.ProjectRoot, r.RepoName); err != nil || len(found) > 0 {
			return fail(reporemove.ErrConflict, "earlier removal left %s; inspect it", strings.Join(found, ", "))
		}
		if err := os.Rename(r.Wrapper, quarantine); err != nil {
			return fail(reporemove.ErrUnknown, "rename %s to %s: %v", r.Wrapper, quarantine, err)
		}
		renamed = true
		if info, err := os.Lstat(quarantine); err != nil || !os.SameFile(info, observed.wrapper) {
			return fail(reporemove.ErrUnknown, "%s is not the renamed wrapper", quarantine)
		}
		return nil
	})
	switch {
	case err != nil && renamed:
		return false, fail(reporemove.ErrUnknown, "%v; %s is preserved for inspection", err, quarantine)
	case err != nil && (errors.Is(err, reporemove.ErrInvalid) || errors.Is(err, reporemove.ErrConflict) || errors.Is(err, reporemove.ErrUnknown)):
		return false, err
	case err != nil:
		return false, classified{reporemove.ErrInvalid, err}
	}
	removeAll := os.RemoveAll
	if a.hooks.removeAll != nil {
		removeAll = a.hooks.removeAll
	}
	if err := removeAll(quarantine); err != nil {
		return false, fail(reporemove.ErrUnknown, "delete %s: %v; it is preserved for inspection", quarantine, err)
	}
	for _, path := range []string{r.Wrapper, quarantine} {
		if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
			return false, fail(reporemove.ErrUnknown, "%s still exists after deletion: %v", path, err)
		}
	}
	return true, nil
}
