package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type GitRepoObservation struct{ Locator, GitCommonDir string }
type GitWorktreeObservation struct {
	Locator, Ref, OID, Tree, GitCommonDir string
	Clean, InventoryMatch                 bool
}
type GitRefObservation struct {
	Exists                       bool
	Ref, OID, Tree, CheckedOutAt string
}
type GitWorktreeInventoryEntry struct {
	Locator, Ref, OID string
	Locked, Prunable  bool
	Bare, Detached    bool
}
type GitWorktreeInventory struct{ Entries []GitWorktreeInventoryEntry }
type GitCreateInput struct {
	Repository                                  RepoRecord
	Branch, SourceRef, TargetLocator, ParentOID string
}
type GitCommandOutcome struct {
	Exit           int
	Stdout, Stderr []byte
	Err            error
}

type WorkItemGit interface {
	ValidateBranch(string) error
	ObserveRepo(RepoRecord) (GitRepoObservation, error)
	ObserveWorktree(string) (GitWorktreeObservation, error)
	ObserveRef(RepoRecord, string) (GitRefObservation, error)
	ListWorktrees(RepoRecord) (GitWorktreeInventory, error)
	CreateBranchAndWorktree(GitCreateInput) GitCommandOutcome
	AddWorktreeForExistingBranch(GitCreateInput) GitCommandOutcome
}

type workGitRunner func(arguments []string, environment []string) GitCommandOutcome
type systemWorkItemGit struct {
	files FileSystem
	run   workGitRunner
}

func newSystemWorkItemGit(files FileSystem) WorkItemGit {
	return &systemWorkItemGit{files: files, run: runWorkItemGit}
}

func runWorkItemGit(arguments []string, environment []string) GitCommandOutcome {
	command := exec.Command("git", arguments...)
	command.Env = append(os.Environ(), environment...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	exit := 0
	if err != nil {
		exit = -1
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			exit = ee.ExitCode()
		}
	}
	return GitCommandOutcome{Exit: exit, Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), Err: err}
}

func (git *systemWorkItemGit) read(path string, arguments ...string) GitCommandOutcome {
	argv := append([]string{"-C", path}, arguments...)
	return git.run(argv, []string{"GIT_OPTIONAL_LOCKS=0", "LC_ALL=C", "LANG=C"})
}

func (git *systemWorkItemGit) ValidateBranch(branch string) error {
	result := git.run([]string{"check-ref-format", "--branch", branch}, []string{"GIT_OPTIONAL_LOCKS=0", "LC_ALL=C", "LANG=C"})
	if result.Err != nil {
		return WorkInvalidArguments(fmt.Sprintf("--branch %q is not a valid short local branch name", branch))
	}
	return nil
}

func (git *systemWorkItemGit) ObserveRepo(repository RepoRecord) (GitRepoObservation, error) {
	top, err := git.gitPath(repository.Locator, "--show-toplevel")
	if err != nil {
		return GitRepoObservation{}, err
	}
	common, err := git.gitPath(repository.Locator, "--git-common-dir")
	if err != nil {
		return GitRepoObservation{}, err
	}
	return GitRepoObservation{Locator: top, GitCommonDir: common}, nil
}

func (git *systemWorkItemGit) ObserveWorktree(path string) (GitWorktreeObservation, error) {
	inside := git.read(path, "rev-parse", "--is-inside-work-tree")
	if inside.Err != nil {
		return GitWorktreeObservation{}, gitObservationError("observe worktree", path, inside)
	}
	insideValue, ok := singleGitLine(inside.Stdout)
	if !ok || insideValue != "true" {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("%s is not a non-bare Git worktree", path), nil)
	}
	locator, err := git.gitPath(path, "--show-toplevel")
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	if locator != path {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("%s is not an exact Git worktree root", path), nil)
	}
	common, err := git.gitPath(path, "--git-common-dir")
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	refResult := git.read(path, "symbolic-ref", "-q", "HEAD")
	if refResult.Err != nil {
		return GitWorktreeObservation{}, gitObservationError("observe symbolic HEAD", path, refResult)
	}
	ref, ok := singleGitLine(refResult.Stdout)
	if !ok || !strings.HasPrefix(ref, "refs/heads/") {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("worktree %s does not have a full symbolic local branch", path), nil)
	}
	oid, err := git.gitLine(path, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	tree, err := git.gitLine(path, "rev-parse", "--verify", "HEAD^{tree}")
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	if !validOIDText(oid) || !validOIDText(tree) {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, "Git returned a malformed object ID", nil)
	}
	status := git.read(path, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if status.Err != nil {
		return GitWorktreeObservation{}, gitObservationError("observe status", path, status)
	}
	repository := RepoRecord{Locator: path, GitCommonDir: common}
	inventory, err := git.ListWorktrees(repository)
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	matches := 0
	for _, entry := range inventory.Entries {
		if entry.Locator == locator && entry.Ref == ref && entry.OID == oid && !entry.Locked && !entry.Prunable && !entry.Bare && !entry.Detached {
			matches++
		}
	}
	return GitWorktreeObservation{Locator: locator, Ref: ref, OID: oid, Tree: tree, GitCommonDir: common, Clean: len(status.Stdout) == 0, InventoryMatch: matches == 1}, nil
}

func (git *systemWorkItemGit) ObserveRef(repository RepoRecord, ref string) (GitRefObservation, error) {
	if !validFullBranchRef(ref) {
		return GitRefObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("ref %s is not a full local branch ref", ref), nil)
	}
	check := git.run([]string{"check-ref-format", ref}, []string{"GIT_OPTIONAL_LOCKS=0", "LC_ALL=C", "LANG=C"})
	if check.Err != nil {
		return GitRefObservation{}, workError(ErrorWorkRefConflict, fmt.Sprintf("ref %s is invalid", ref), check.Err)
	}
	result := git.read(repository.Locator, "show-ref", "--verify", "--hash", ref)
	if result.Err != nil {
		if result.Exit == 1 || (result.Exit == 128 && strings.Contains(string(result.Stderr), "not a valid ref")) {
			return GitRefObservation{Ref: ref}, nil
		}
		return GitRefObservation{}, gitObservationError("observe ref", ref, result)
	}
	oid, ok := singleGitLine(result.Stdout)
	if !ok || !validOIDText(oid) {
		return GitRefObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("Git returned malformed OID for %s", ref), nil)
	}
	tree, err := git.gitLine(repository.Locator, "rev-parse", "--verify", ref+"^{tree}")
	if err != nil {
		return GitRefObservation{}, err
	}
	inventory, err := git.ListWorktrees(repository)
	if err != nil {
		return GitRefObservation{}, err
	}
	checked := ""
	for _, entry := range inventory.Entries {
		if entry.Ref == ref {
			if checked != "" && checked != entry.Locator {
				return GitRefObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("ref %s is checked out more than once", ref), nil)
			}
			checked = entry.Locator
		}
	}
	return GitRefObservation{Exists: true, Ref: ref, OID: oid, Tree: tree, CheckedOutAt: checked}, nil
}

func (git *systemWorkItemGit) ListWorktrees(repository RepoRecord) (GitWorktreeInventory, error) {
	result := git.read(repository.Locator, "worktree", "list", "--porcelain", "-z")
	if result.Err != nil {
		return GitWorktreeInventory{}, gitObservationError("list worktrees", repository.Locator, result)
	}
	entries, err := parseWorktreeInventory(result.Stdout, git.files)
	if err != nil {
		return GitWorktreeInventory{}, workError(ErrorWorkGitObservation, "parse Git worktree inventory", err)
	}
	return GitWorktreeInventory{Entries: entries}, nil
}

func (git *systemWorkItemGit) CreateBranchAndWorktree(input GitCreateInput) GitCommandOutcome {
	return git.run([]string{"-c", "core.hooksPath=" + os.DevNull, "-C", input.Repository.Locator, "worktree", "add", "-b", input.Branch, input.TargetLocator, input.ParentOID}, []string{"LC_ALL=C", "LANG=C"})
}
func (git *systemWorkItemGit) AddWorktreeForExistingBranch(input GitCreateInput) GitCommandOutcome {
	return git.run([]string{"-c", "core.hooksPath=" + os.DevNull, "-C", input.Repository.Locator, "worktree", "add", input.TargetLocator, input.Branch}, []string{"LC_ALL=C", "LANG=C"})
}

func (git *systemWorkItemGit) gitPath(path, argument string) (string, error) {
	value, err := git.gitLine(path, "rev-parse", "--path-format=absolute", argument)
	if err != nil {
		return "", err
	}
	physical, err := canonicalDirectory(git.files, "", value, "Git path")
	if err != nil {
		return "", workError(ErrorWorkGitObservation, fmt.Sprintf("observe %s for %s", argument, path), err)
	}
	return physical, nil
}
func (git *systemWorkItemGit) gitLine(path string, arguments ...string) (string, error) {
	result := git.read(path, arguments...)
	if result.Err != nil {
		return "", gitObservationError(strings.Join(arguments, " "), path, result)
	}
	value, ok := singleGitLine(result.Stdout)
	if !ok || value == "" {
		return "", workError(ErrorWorkGitObservation, fmt.Sprintf("Git returned malformed output for %s", strings.Join(arguments, " ")), nil)
	}
	return value, nil
}
func gitObservationError(operation, path string, outcome GitCommandOutcome) error {
	return workError(ErrorWorkGitObservation, fmt.Sprintf("%s for %s: git exit %d", operation, path, outcome.Exit), outcome.Err)
}

func parseWorktreeInventory(output []byte, files FileSystem) ([]GitWorktreeInventoryEntry, error) {
	if len(output) == 0 {
		return []GitWorktreeInventoryEntry{}, nil
	}
	parts := bytes.Split(output, []byte{0})
	entries := []GitWorktreeInventoryEntry{}
	var current *GitWorktreeInventoryEntry
	finish := func() error {
		if current == nil {
			return nil
		}
		if current.Locator == "" || current.OID == "" || (!current.Bare && !current.Detached && current.Ref == "") {
			return errors.New("incomplete worktree record")
		}
		physical, err := canonicalDirectory(files, "", current.Locator, "Git worktree inventory path")
		if err != nil {
			return err
		}
		current.Locator = physical
		entries = append(entries, *current)
		current = nil
		return nil
	}
	seen := map[string]bool{}
	for _, raw := range parts {
		line := string(raw)
		if line == "" {
			if err := finish(); err != nil {
				return nil, err
			}
			continue
		}
		key, value, _ := strings.Cut(line, " ")
		switch key {
		case "worktree":
			if current != nil {
				if err := finish(); err != nil {
					return nil, err
				}
			}
			current = &GitWorktreeInventoryEntry{Locator: value}
			seen = map[string]bool{"worktree": true}
		case "HEAD":
			if current == nil {
				return nil, errors.New("HEAD outside worktree record")
			}
			if seen[key] {
				return nil, errors.New("duplicate HEAD in worktree record")
			}
			seen[key] = true
			current.OID = value
		case "branch":
			if current == nil {
				return nil, errors.New("branch outside worktree record")
			}
			if seen[key] {
				return nil, errors.New("duplicate branch in worktree record")
			}
			seen[key] = true
			current.Ref = value
		case "locked":
			if current == nil {
				return nil, errors.New("locked outside worktree record")
			}
			if seen[key] {
				return nil, errors.New("duplicate locked in worktree record")
			}
			seen[key] = true
			current.Locked = true
		case "prunable":
			if current == nil {
				return nil, errors.New("prunable outside worktree record")
			}
			if seen[key] {
				return nil, errors.New("duplicate prunable in worktree record")
			}
			seen[key] = true
			current.Prunable = true
		case "bare", "detached":
			if current == nil {
				return nil, fmt.Errorf("%s outside worktree record", key)
			}
			if seen[key] {
				return nil, fmt.Errorf("duplicate %s in worktree record", key)
			}
			seen[key] = true
			if key == "bare" {
				current.Bare = true
			} else {
				current.Detached = true
			}
		default:
			return nil, fmt.Errorf("unknown worktree record field %s", key)
		}
	}
	if err := finish(); err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Locator < entries[j].Locator })
	return entries, nil
}

func worktreeEntry(inventory GitWorktreeInventory, locator string) (GitWorktreeInventoryEntry, bool) {
	for _, entry := range inventory.Entries {
		if filepath.Clean(entry.Locator) == filepath.Clean(locator) {
			return entry, true
		}
	}
	return GitWorktreeInventoryEntry{}, false
}
