package workflowhandoff

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type systemFileSystem struct{}

func (systemFileSystem) Getwd() (string, error)                       { return os.Getwd() }
func (systemFileSystem) EvalSymlinks(path string) (string, error)     { return filepath.EvalSymlinks(path) }
func (systemFileSystem) Lstat(path string) (fs.FileInfo, error)       { return os.Lstat(path) }
func (systemFileSystem) Stat(path string) (fs.FileInfo, error)        { return os.Stat(path) }
func (systemFileSystem) ReadFile(path string) ([]byte, error)         { return os.ReadFile(path) }
func (systemFileSystem) ReadDir(path string) ([]fs.DirEntry, error)   { return os.ReadDir(path) }
func (systemFileSystem) Mkdir(path string, mode fs.FileMode) error    { return os.Mkdir(path, mode) }
func (systemFileSystem) MkdirAll(path string, mode fs.FileMode) error { return os.MkdirAll(path, mode) }
func (systemFileSystem) MkdirTemp(path, pattern string) (string, error) {
	return os.MkdirTemp(path, pattern)
}
func (systemFileSystem) OpenFile(path string, flag int, mode fs.FileMode) (File, error) {
	return os.OpenFile(path, flag, mode)
}
func (systemFileSystem) Chmod(path string, mode fs.FileMode) error { return os.Chmod(path, mode) }
func (systemFileSystem) Rename(oldPath, newPath string) error      { return os.Rename(oldPath, newPath) }
func (systemFileSystem) Remove(path string) error                  { return os.Remove(path) }

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

type systemWorkspaceObserver struct{ dependencies workspace.Dependencies }

func (observer systemWorkspaceObserver) ObserveContaining() (WorkspaceSnapshot, error) {
	observation, err := workspace.ObserveContaining(observer.dependencies)
	if err != nil {
		return WorkspaceSnapshot{}, translateWorkspaceError(err)
	}
	return observer.snapshot(observation)
}

func (observer systemWorkspaceObserver) ObserveRoot(root string) (WorkspaceSnapshot, error) {
	observation, err := workspace.ObserveRoot(observer.dependencies, root)
	if err != nil {
		return WorkspaceSnapshot{}, translateWorkspaceError(err)
	}
	return observer.snapshot(observation)
}

func (observer systemWorkspaceObserver) snapshot(observation workspace.WorkspaceObservation) (WorkspaceSnapshot, error) {
	projects, repositories, err := observer.dependencies.Projects.Snapshot(observation.Root)
	if err != nil {
		return WorkspaceSnapshot{}, translateWorkspaceError(err)
	}
	return WorkspaceSnapshot{Observation: observation, Projects: projects, Repositories: repositories}, nil
}

func translateWorkspaceError(err error) error {
	var projectErr *workspace.ProjectError
	if errors.As(err, &projectErr) {
		switch projectErr.Class {
		case workspace.ErrorWorkspaceNotFound:
			return classified(ErrorWorkspaceNotFound, projectErr.Detail, err)
		case workspace.ErrorProjectStoreConflict, workspace.ErrorProjectWorkspace:
			return classified(ErrorWorkspaceConflict, projectErr.Detail, err)
		default:
			return classified(ErrorIO, projectErr.Detail, err)
		}
	}
	return classified(ErrorIO, err.Error(), err)
}

type gitResult struct {
	stdout, stderr []byte
	exit           int
	err            error
}
type systemGitObserver struct {
	files FileSystem
	run   func(string, ...string) gitResult
}

func SystemDependencies() Dependencies {
	files := systemFileSystem{}
	workspaceDependencies := workspace.SystemDependencies()
	clock := systemClock{}
	dependencies := Dependencies{
		TaskWorkspace: &workspaceDependencies,
		Files:         files,
		Workspace:     systemWorkspaceObserver{dependencies: workspaceDependencies},
		Git:           systemGitObserver{files: files, run: runGitCommand},
		Clock:         clock,
		Random:        rand.Reader,
	}
	dependencies.Store = newSystemStore(files, clock)
	return dependencies
}

func runGitCommand(directory string, arguments ...string) gitResult {
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	exit := 0
	if err != nil {
		exit = -1
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exit = exitError.ExitCode()
		}
	}
	return gitResult{stdout: stdout.Bytes(), stderr: stderr.Bytes(), exit: exit, err: err}
}

func (observer systemGitObserver) ObserveTarget(request TargetRequest) (TargetObservation, error) {
	worktree, err := observer.physicalDirectory(request.Worktree)
	if err != nil || worktree != request.Worktree {
		return TargetObservation{}, targetConflict(request.Worktree, "target must be an absolute physical directory", err)
	}
	inside, err := observer.gitLine(worktree, "rev-parse", "--is-inside-work-tree")
	if err != nil || inside != "true" {
		return TargetObservation{}, targetConflict(worktree, "target is not a Git worktree", err)
	}
	top, err := observer.gitPath(worktree, "--show-toplevel")
	if err != nil || top != worktree {
		return TargetObservation{}, targetConflict(worktree, "target is not the Git worktree root", err)
	}
	common, err := observer.gitPath(worktree, "--git-common-dir")
	if err != nil {
		return TargetObservation{}, err
	}
	ref, err := observer.gitLine(worktree, "symbolic-ref", "--quiet", "HEAD")
	if err != nil || ref != request.Ref {
		return TargetObservation{}, targetConflict(worktree, "checked out ref does not match "+request.Ref, err)
	}
	oid, err := observer.gitLine(worktree, "rev-parse", "--verify", request.Ref)
	if err != nil || !validateOID(oid) {
		return TargetObservation{}, targetConflict(worktree, "cannot resolve expected ref", err)
	}
	tree, err := observer.gitLine(worktree, "rev-parse", "--verify", oid+"^{tree}")
	if err != nil || !validateOID(tree) || len(tree) != len(oid) {
		return TargetObservation{}, targetConflict(worktree, "cannot resolve target tree", err)
	}
	statusResult := observer.run(worktree, "status", "--porcelain=v2", "--untracked-files=all", "-z")
	if statusResult.err != nil {
		return TargetObservation{}, targetConflict(worktree, fmt.Sprintf("git status exited %d", statusResult.exit), statusResult.err)
	}
	status, err := observer.parseStatus(worktree, statusResult.stdout)
	if err != nil {
		return TargetObservation{}, targetConflict(worktree, "malformed Git status: "+err.Error(), err)
	}
	return TargetObservation{Worktree: worktree, GitCommonDir: common, Ref: ref, OID: oid, Tree: tree, Status: status}, nil
}

func (observer systemGitObserver) ObserveInputGit(request InputGitRequest) (InputGitObservation, error) {
	if !filepath.IsAbs(request.Locator) || filepath.Clean(request.Locator) != request.Locator {
		return InputGitObservation{}, targetConflict(request.Locator, "input locator is not absolute and clean", nil)
	}
	physical, err := observer.files.EvalSymlinks(request.Locator)
	if err != nil || physical != request.Locator {
		return InputGitObservation{}, targetConflict(request.Locator, "input locator is not physical", err)
	}
	info, err := observer.files.Lstat(request.Locator)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return InputGitObservation{}, targetConflict(request.Locator, "input is not a regular non-symlink file", err)
	}
	directory := filepath.Dir(request.Locator)
	top, err := observer.gitPath(directory, "--show-toplevel")
	if err != nil {
		return InputGitObservation{}, err
	}
	common, err := observer.gitPath(directory, "--git-common-dir")
	if err != nil {
		return InputGitObservation{}, err
	}
	ref, err := observer.gitLine(directory, "symbolic-ref", "--quiet", "HEAD")
	if err != nil {
		return InputGitObservation{}, targetConflict(request.Locator, "input worktree is detached", err)
	}
	oid, err := observer.gitLine(directory, "rev-parse", "--verify", request.Binding.Ref)
	if err != nil {
		return InputGitObservation{}, targetConflict(request.Locator, "cannot resolve input ref", err)
	}
	relative, err := filepath.Rel(top, request.Locator)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return InputGitObservation{}, targetConflict(request.Locator, "input is outside its worktree", err)
	}
	relative = filepath.ToSlash(relative)
	result := observer.run(directory, "ls-tree", oid, "--", relative)
	if result.err != nil {
		return InputGitObservation{}, targetConflict(request.Locator, "cannot observe input blob", result.err)
	}
	line, ok := singleLine(result.stdout)
	if !ok {
		return InputGitObservation{}, targetConflict(request.Locator, "malformed input ls-tree output", nil)
	}
	metadata, listedPath, ok := strings.Cut(line, "\t")
	parts := strings.Fields(metadata)
	if !ok || listedPath != relative || len(parts) != 3 || parts[1] != "blob" || !validateOID(parts[2]) {
		return InputGitObservation{}, targetConflict(request.Locator, "input path is not a regular Git blob", nil)
	}
	observed := GitBinding{GitCommonDir: common, Ref: ref, OID: oid, Blob: parts[2]}
	match := observed == request.Binding
	return InputGitObservation{Binding: observed, MatchesExpected: match}, nil
}

func (observer systemGitObserver) physicalDirectory(input string) (string, error) {
	if !filepath.IsAbs(input) || filepath.Clean(input) != input {
		return "", fmt.Errorf("path is not absolute and clean")
	}
	physical, err := observer.files.EvalSymlinks(input)
	if err != nil {
		return "", err
	}
	physical = filepath.Clean(physical)
	info, err := observer.files.Stat(physical)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}
	return physical, nil
}
func (observer systemGitObserver) gitLine(directory string, arguments ...string) (string, error) {
	result := observer.run(directory, arguments...)
	if result.err != nil {
		return "", fmt.Errorf("git %s exited %d: %w", arguments[0], result.exit, result.err)
	}
	line, ok := singleLine(result.stdout)
	if !ok {
		return "", fmt.Errorf("git %s returned malformed output", arguments[0])
	}
	return line, nil
}
func (observer systemGitObserver) gitPath(directory, argument string) (string, error) {
	value, err := observer.gitLine(directory, "rev-parse", "--path-format=absolute", argument)
	if err != nil {
		return "", targetConflict(directory, err.Error(), err)
	}
	physical, err := observer.physicalDirectory(value)
	if err != nil {
		return "", targetConflict(value, "Git returned a non-physical directory", err)
	}
	return physical, nil
}
func singleLine(output []byte) (string, bool) {
	text := string(output)
	if strings.HasSuffix(text, "\r\n") {
		text = strings.TrimSuffix(text, "\r\n")
	} else if strings.HasSuffix(text, "\n") {
		text = strings.TrimSuffix(text, "\n")
	} else {
		return "", false
	}
	if strings.ContainsAny(text, "\r\n") || text == "" {
		return "", false
	}
	return text, true
}
func targetConflict(path, detail string, err error) error {
	return classified(ErrorTargetConflict, fmt.Sprintf("%s: %s", path, detail), err)
}

func (observer systemGitObserver) parseStatus(worktree string, output []byte) (StatusPolicy, error) {
	if len(output) == 0 {
		return StatusPolicy{Mode: "clean", Value: canonicaljson.Object{{Name: "mode", Value: "clean"}}}, nil
	}
	records := bytes.Split(output, []byte{0})
	if len(records) == 0 || len(records[len(records)-1]) != 0 {
		return StatusPolicy{}, fmt.Errorf("status is not NUL terminated")
	}
	records = records[:len(records)-1]
	entries := []StatusEntry{}
	for index := 0; index < len(records); index++ {
		record := string(records[index])
		if len(record) < 2 {
			return StatusPolicy{}, fmt.Errorf("short status record")
		}
		switch record[0] {
		case '?':
			if !strings.HasPrefix(record, "? ") {
				return StatusPolicy{}, fmt.Errorf("malformed untracked record")
			}
			entry, err := observer.statusEntry(worktree, record[2:], "unmodified", "untracked")
			if err != nil {
				return StatusPolicy{}, err
			}
			entries = append(entries, entry)
		case '1':
			parts := strings.SplitN(record, " ", 9)
			if len(parts) != 9 || len(parts[1]) != 2 {
				return StatusPolicy{}, fmt.Errorf("malformed ordinary record")
			}
			indexState, err := statusCode(parts[1][0])
			if err != nil {
				return StatusPolicy{}, err
			}
			worktreeState, err := statusCode(parts[1][1])
			if err != nil {
				return StatusPolicy{}, err
			}
			entry, err := observer.statusEntry(worktree, parts[8], indexState, worktreeState)
			if err != nil {
				return StatusPolicy{}, err
			}
			entries = append(entries, entry)
		case '2':
			parts := strings.SplitN(record, " ", 10)
			if len(parts) != 10 || len(parts[1]) != 2 || index+1 >= len(records) {
				return StatusPolicy{}, fmt.Errorf("malformed rename record")
			}
			oldPath := string(records[index+1])
			if !utf8.ValidString(oldPath) || validateRepoPath(filepath.ToSlash(oldPath)) != nil {
				return StatusPolicy{}, fmt.Errorf("invalid rename source path")
			}
			index++
			worktreeState, err := statusCode(parts[1][1])
			if err != nil {
				return StatusPolicy{}, err
			}
			deleted := StatusEntry{Path: filepath.ToSlash(oldPath), IndexState: "deleted", WorktreeState: "unmodified", Kind: "regular"}
			added, err := observer.statusEntry(worktree, parts[9], "added", worktreeState)
			if err != nil {
				return StatusPolicy{}, err
			}
			entries = append(entries, deleted, added)
		case 'u':
			return StatusPolicy{}, fmt.Errorf("unmerged status is unsupported")
		case '!':
			continue
		default:
			return StatusPolicy{}, fmt.Errorf("unknown status record")
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	values := make([]canonicaljson.Value, 0, len(entries))
	for _, entry := range entries {
		var digest canonicaljson.Value = nil
		if entry.SHA256 != "" {
			digest = entry.SHA256
		}
		values = append(values, canonicaljson.Object{{Name: "path", Value: entry.Path}, {Name: "index_state", Value: entry.IndexState}, {Name: "worktree_state", Value: entry.WorktreeState}, {Name: "kind", Value: entry.Kind}, {Name: "sha256", Value: digest}})
	}
	return StatusPolicy{Mode: "exact_manifest", Entries: entries, Value: canonicaljson.Object{{Name: "mode", Value: "exact_manifest"}, {Name: "entries", Value: values}}}, nil
}

func statusCode(value byte) (string, error) {
	switch value {
	case '.':
		return "unmodified", nil
	case 'M':
		return "modified", nil
	case 'A', 'R', 'C':
		return "added", nil
	case 'D':
		return "deleted", nil
	case 'T':
		return "type_changed", nil
	case '?':
		return "untracked", nil
	default:
		return "", fmt.Errorf("unknown status code %q", value)
	}
}
func (observer systemGitObserver) statusEntry(worktree, path, indexState, worktreeState string) (StatusEntry, error) {
	if !utf8.ValidString(path) {
		return StatusEntry{}, fmt.Errorf("status path is not valid UTF-8")
	}
	if err := validateRepoPath(filepath.ToSlash(path)); err != nil {
		return StatusEntry{}, err
	}
	entry := StatusEntry{Path: filepath.ToSlash(path), IndexState: indexState, WorktreeState: worktreeState, Kind: "regular"}
	if indexState == "deleted" || worktreeState == "deleted" {
		return entry, nil
	}
	full := filepath.Join(worktree, filepath.FromSlash(path))
	info, err := observer.files.Lstat(full)
	if err != nil {
		return StatusEntry{}, err
	}
	var contents []byte
	if info.Mode()&fs.ModeSymlink != 0 {
		entry.Kind = "symlink"
		target, err := os.Readlink(full)
		if err != nil {
			return StatusEntry{}, err
		}
		contents = []byte(target)
	} else if info.Mode().IsRegular() {
		contents, err = observer.files.ReadFile(full)
		if err != nil {
			return StatusEntry{}, err
		}
	} else {
		return StatusEntry{}, fmt.Errorf("status path is not regular or symlink")
	}
	entry.SHA256 = digestBytes(contents)
	return entry, nil
}

func hashFile(files FileSystem, path string) (string, int64, error) {
	info, err := files.Lstat(path)
	if err != nil {
		return "", 0, err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return "", 0, fmt.Errorf("%s is not a regular non-symlink file", path)
	}
	contents, err := files.ReadFile(path)
	if err != nil {
		return "", 0, err
	}
	sum := sha256.Sum256(contents)
	return fmt.Sprintf("sha256:%x", sum), int64(len(contents)), nil
}
