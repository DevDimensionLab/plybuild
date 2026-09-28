package workspace

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
type workGitInputRunner func(arguments []string, environment []string, stdin []byte) GitCommandOutcome
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

type systemTaskIntegrationGit struct {
	files    FileSystem
	run      workGitRunner
	runInput workGitInputRunner
}

func newSystemTaskIntegrationGit(files FileSystem) TaskIntegrationGit {
	return &systemTaskIntegrationGit{files: files, run: runSafeIntegrationGit, runInput: runSafeIntegrationGitInput}
}

func runSafeIntegrationGit(arguments []string, environment []string) GitCommandOutcome {
	return runSafeIntegrationGitInput(arguments, environment, nil)
}

func runSafeIntegrationGitInput(arguments []string, environment []string, stdin []byte) GitCommandOutcome {
	command := exec.Command("git", arguments...)
	command.Env = safeIntegrationEnvironment(environment)
	if stdin != nil {
		command.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
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

func safeIntegrationEnvironment(extra []string) []string {
	allowed := map[string]bool{"PATH": true, "TMPDIR": true, "TMP": true, "TEMP": true, "SystemRoot": true, "WINDIR": true, "PATHEXT": true}
	env := []string{}
	for _, entry := range os.Environ() {
		name, _, ok := strings.Cut(entry, "=")
		if ok && allowed[name] {
			env = append(env, entry)
		}
	}
	env = append(env, "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "GIT_ALLOW_PROTOCOL=file", "GIT_TERMINAL_PROMPT=0", "LC_ALL=C", "LANG=C", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull)
	env = append(env, extra...)
	return env
}

func integrationGitPrefix(path string) []string {
	return []string{"-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", path}
}
func (git *systemTaskIntegrationGit) read(path string, args ...string) GitCommandOutcome {
	argv := append(integrationGitPrefix(path), args...)
	return git.run(argv, append(safeReadExtra(), "GIT_OPTIONAL_LOCKS=0"))
}
func safeReadExtra() []string { return []string{} }
func integrationLine(out GitCommandOutcome, operation string) (string, error) {
	if out.Err != nil {
		return "", fmt.Errorf("%s: %s", operation, safeGitDetail(out))
	}
	line, ok := singleGitLine(out.Stdout)
	if !ok {
		return "", fmt.Errorf("%s returned malformed output", operation)
	}
	return line, nil
}
func safeGitDetail(out GitCommandOutcome) string {
	detail := strings.TrimSpace(string(out.Stderr))
	if detail == "" {
		detail = out.Err.Error()
	}
	if len(detail) > 600 {
		detail = detail[:600]
	}
	return detail
}

func (git *systemTaskIntegrationGit) ObserveIntegrationRepository(repo RepoRecord) (IntegrationRepositoryObservation, error) {
	version, err := integrationLine(git.run([]string{"--version"}, append(safeReadExtra(), "GIT_OPTIONAL_LOCKS=0")), "git --version")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	if !supportedGitVersion(version) {
		return IntegrationRepositoryObservation{}, fmt.Errorf("unsupported Git version %s; version 2.45.0 or newer is required", version)
	}
	top, err := integrationLine(git.read(repo.Locator, "rev-parse", "--path-format=absolute", "--show-toplevel"), "observe repository top-level")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	top, err = canonicalDirectory(git.files, "", top, "integration repository")
	if err != nil || top != repo.Locator {
		return IntegrationRepositoryObservation{}, fmt.Errorf("registered repository locator is not the exact physical top-level")
	}
	common, err := integrationLine(git.read(repo.Locator, "rev-parse", "--path-format=absolute", "--git-common-dir"), "observe repository common directory")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	common, err = canonicalDirectory(git.files, "", common, "Git common directory")
	if err != nil || common != repo.GitCommonDir {
		return IntegrationRepositoryObservation{}, fmt.Errorf("registered Git common directory differs from the repository")
	}
	bare, err := integrationLine(git.read(repo.Locator, "rev-parse", "--is-bare-repository"), "observe bare repository")
	if err != nil || bare != "false" {
		return IntegrationRepositoryObservation{}, fmt.Errorf("repository must be non-bare")
	}
	shallow, err := integrationLine(git.read(repo.Locator, "rev-parse", "--is-shallow-repository"), "observe shallow repository")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	object, err := integrationLine(git.read(repo.Locator, "rev-parse", "--show-object-format=storage"), "observe object format")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	refFormat, err := integrationLine(git.read(repo.Locator, "rev-parse", "--show-ref-format"), "observe ref format")
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	sparse := false
	sparseOut := git.read(repo.Locator, "config", "--get", "--bool", "core.sparseCheckout")
	if sparseOut.Exit == 0 {
		line, ok := singleGitLine(sparseOut.Stdout)
		if !ok {
			return IntegrationRepositoryObservation{}, fmt.Errorf("malformed sparse-checkout setting")
		}
		sparse = line == "true"
	} else if sparseOut.Exit != 1 {
		return IntegrationRepositoryObservation{}, fmt.Errorf("observe sparse-checkout: %s", safeGitDetail(sparseOut))
	}
	partial := false
	partialOut := git.read(repo.Locator, "config", "--local", "--get", "extensions.partialClone")
	if partialOut.Exit == 0 {
		partial = true
	} else if partialOut.Exit != 1 {
		return IntegrationRepositoryObservation{}, fmt.Errorf("observe partial clone: %s", safeGitDetail(partialOut))
	}
	invOut := git.read(repo.Locator, "worktree", "list", "--porcelain", "-z")
	if invOut.Err != nil {
		return IntegrationRepositoryObservation{}, fmt.Errorf("observe worktree inventory: %s", safeGitDetail(invOut))
	}
	old, err := parseWorktreeInventory(invOut.Stdout, git.files)
	if err != nil {
		return IntegrationRepositoryObservation{}, err
	}
	inventory := make([]IntegrationInventoryEntry, len(old))
	for i, e := range old {
		inventory[i] = IntegrationInventoryEntry{Locator: e.Locator, Ref: e.Ref, OID: e.OID, Locked: e.Locked, Prunable: e.Prunable, Bare: e.Bare, Detached: e.Detached}
	}
	sort.Slice(inventory, func(i, j int) bool {
		if inventory[i].Locator != inventory[j].Locator {
			return inventory[i].Locator < inventory[j].Locator
		}
		if inventory[i].Ref != inventory[j].Ref {
			return inventory[i].Ref < inventory[j].Ref
		}
		return inventory[i].OID < inventory[j].OID
	})
	return IntegrationRepositoryObservation{GitVersion: version, ObjectFormat: object, RefFormat: refFormat, Shallow: shallow == "true", PartialClone: partial, SparseCheckout: sparse, Inventory: inventory}, nil
}

var gitVersionPattern = regexp.MustCompile(`^git version ([0-9]+)\.([0-9]+)(?:\.([0-9]+))?`)

func supportedGitVersion(value string) bool {
	m := gitVersionPattern.FindStringSubmatch(value)
	if m == nil {
		return false
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	return major > 2 || (major == 2 && minor >= 45)
}

func (git *systemTaskIntegrationGit) ObserveIntegrationWorktree(path, expectedRef string) (IntegrationWorktreeObservation, error) {
	top, err := integrationLine(git.read(path, "rev-parse", "--path-format=absolute", "--show-toplevel"), "observe worktree top-level")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	physical, err := canonicalDirectory(git.files, "", top, "integration worktree")
	if err != nil || physical != path {
		return IntegrationWorktreeObservation{}, fmt.Errorf("worktree locator is not the exact physical top-level")
	}
	common, err := integrationLine(git.read(path, "rev-parse", "--path-format=absolute", "--git-common-dir"), "observe common directory")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	common, err = canonicalDirectory(git.files, "", common, "Git common directory")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	ref, err := integrationLine(git.read(path, "symbolic-ref", "-q", "HEAD"), "observe symbolic HEAD")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	if ref != expectedRef {
		return IntegrationWorktreeObservation{}, fmt.Errorf("worktree ref %s differs from %s", ref, expectedRef)
	}
	oid, err := integrationLine(git.read(path, "rev-parse", "--verify", "--end-of-options", ref+"^{commit}"), "observe worktree commit")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	tree, err := integrationLine(git.read(path, "rev-parse", "--verify", "--end-of-options", oid+"^{tree}"), "observe worktree tree")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	object, err := integrationLine(git.read(path, "rev-parse", "--show-object-format=storage"), "observe object format")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	rf, err := integrationLine(git.read(path, "rev-parse", "--show-ref-format"), "observe ref format")
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	statusOut := git.read(path, "status", "--porcelain=v2", "-z", "--untracked-files=all", "--ignore-submodules=none")
	if statusOut.Err != nil {
		return IntegrationWorktreeObservation{}, fmt.Errorf("observe status: %s", safeGitDetail(statusOut))
	}
	status, err := parseIntegrationStatus(statusOut.Stdout)
	if err != nil {
		return IntegrationWorktreeObservation{}, err
	}
	inProgress := []string{}
	for _, marker := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "BISECT_LOG", "rebase-merge", "rebase-apply", "sequencer"} {
		p, err := integrationLine(git.read(path, "rev-parse", "--path-format=absolute", "--git-path", marker), "observe operation marker")
		if err != nil {
			return IntegrationWorktreeObservation{}, err
		}
		if _, e := git.files.Lstat(p); e == nil {
			inProgress = append(inProgress, marker)
		} else if !errors.Is(e, os.ErrNotExist) {
			return IntegrationWorktreeObservation{}, fmt.Errorf("inspect operation marker %s: %v", marker, e)
		}
	}
	return IntegrationWorktreeObservation{Locator: path, Ref: ref, OID: oid, Tree: tree, GitCommonDir: common, ObjectFormat: object, RefFormat: rf, Symbolic: true, Clean: len(status) == 0, StatusEntries: status, InProgress: inProgress}, nil
}

func parseIntegrationStatus(data []byte) ([]IntegrationStatusEntry, error) {
	records := bytes.Split(data, []byte{0})
	result := []IntegrationStatusEntry{}
	for i := 0; i < len(records); i++ {
		record := records[i]
		if len(record) == 0 {
			continue
		}
		prefix := record[0]
		entry := IntegrationStatusEntry{}
		switch prefix {
		case '?':
			if len(record) < 3 || record[1] != ' ' {
				return nil, fmt.Errorf("malformed untracked status")
			}
			entry.RecordKind = "untracked"
			entry.PathBase64 = base64.StdEncoding.EncodeToString(record[2:])
		case '1':
			parts := bytes.SplitN(record, []byte{' '}, 9)
			if len(parts) != 9 || len(parts[1]) != 2 {
				return nil, fmt.Errorf("malformed ordinary status")
			}
			entry.RecordKind = "ordinary"
			entry.PathBase64 = base64.StdEncoding.EncodeToString(parts[8])
			entry.IndexState = stringPointer(string(parts[1][0]))
			entry.WorktreeState = stringPointer(string(parts[1][1]))
			entry.SubmoduleState = stringPointer(string(parts[2]))
		case '2':
			parts := bytes.SplitN(record, []byte{' '}, 10)
			if len(parts) != 10 || i+1 >= len(records) || len(records[i+1]) == 0 || len(parts[1]) != 2 {
				return nil, fmt.Errorf("malformed rename/copy status")
			}
			kind := "renamed"
			if len(parts[8]) > 0 && parts[8][0] == 'C' {
				kind = "copied"
			}
			entry.RecordKind = kind
			entry.PathBase64 = base64.StdEncoding.EncodeToString(parts[9])
			orig := base64.StdEncoding.EncodeToString(records[i+1])
			entry.OriginalPathBase64 = &orig
			entry.IndexState = stringPointer(string(parts[1][0]))
			entry.WorktreeState = stringPointer(string(parts[1][1]))
			entry.SubmoduleState = stringPointer(string(parts[2]))
			i++
		case 'u':
			parts := bytes.SplitN(record, []byte{' '}, 11)
			if len(parts) != 11 || len(parts[1]) != 2 {
				return nil, fmt.Errorf("malformed unmerged status")
			}
			entry.RecordKind = "unmerged"
			entry.PathBase64 = base64.StdEncoding.EncodeToString(parts[10])
			entry.IndexState = stringPointer(string(parts[1][0]))
			entry.WorktreeState = stringPointer(string(parts[1][1]))
			entry.SubmoduleState = stringPointer(string(parts[2]))
		default:
			return nil, fmt.Errorf("unknown porcelain v2 status record")
		}
		result = append(result, entry)
	}
	sort.Slice(result, func(i, j int) bool {
		a, _ := base64.StdEncoding.DecodeString(result[i].PathBase64)
		b, _ := base64.StdEncoding.DecodeString(result[j].PathBase64)
		return bytes.Compare(a, b) < 0
	})
	return result, nil
}

func (git *systemTaskIntegrationGit) CheckAncestor(repo RepoRecord, expected, result string) (bool, error) {
	out := git.read(repo.Locator, "merge-base", "--is-ancestor", expected, result)
	if out.Exit == 1 {
		return false, nil
	}
	if out.Err != nil {
		return false, fmt.Errorf("check ancestry: %s", safeGitDetail(out))
	}
	paths := git.read(repo.Locator, "diff-tree", "-r", "-z", "--no-commit-id", "--name-only", expected, result)
	if paths.Err != nil {
		return false, fmt.Errorf("observe changed paths: %s", safeGitDetail(paths))
	}
	if len(paths.Stdout) > 0 {
		if git.runInput == nil {
			return false, fmt.Errorf("check attributes: stdin-capable Git runner is unavailable")
		}
		for _, oid := range []string{expected, result} {
			argv := append(integrationGitPrefix(repo.Locator), "check-attr", "--source="+oid, "-z", "--stdin", "filter")
			attr := git.runInput(argv, append(safeReadExtra(), "GIT_OPTIONAL_LOCKS=0"), paths.Stdout)
			if attr.Err != nil {
				return false, fmt.Errorf("check attributes: %s", safeGitDetail(attr))
			}
			if err := validateFilterAttributeFrames(paths.Stdout, attr.Stdout); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func validateFilterAttributeFrames(paths, output []byte) error {
	if len(paths) == 0 || paths[len(paths)-1] != 0 || len(output) == 0 || output[len(output)-1] != 0 {
		return fmt.Errorf("check attributes returned malformed NUL framing")
	}
	pathFields := bytes.Split(paths[:len(paths)-1], []byte{0})
	fields := bytes.Split(output[:len(output)-1], []byte{0})
	if len(fields) != len(pathFields)*3 {
		return fmt.Errorf("check attributes returned an unexpected field count")
	}
	for i, path := range pathFields {
		if len(path) == 0 || !bytes.Equal(fields[i*3], path) || string(fields[i*3+1]) != "filter" || string(fields[i*3+2]) != "unspecified" {
			return fmt.Errorf("changed path uses an active or malformed Git filter")
		}
	}
	return nil
}

func (git *systemTaskIntegrationGit) ObserveParentReflog(path, ref string, limit int) ([]IntegrationReflogEntry, error) {
	exists := git.read(path, "reflog", "exists", ref)
	if exists.Err != nil {
		return nil, fmt.Errorf("parent reflog is unavailable")
	}
	out := git.read(path, "reflog", "show", "-z", "-n", strconv.Itoa(limit), "--format=%H%x00%gD%x00%gs", ref)
	if out.Err != nil {
		return nil, fmt.Errorf("observe parent reflog: %s", safeGitDetail(out))
	}
	fields := bytes.Split(out.Stdout, []byte{0})
	values := []string{}
	for _, f := range fields {
		v := strings.TrimSpace(string(f))
		if v != "" {
			values = append(values, v)
		}
	}
	if len(values)%3 != 0 {
		return nil, fmt.Errorf("malformed parent reflog")
	}
	entries := []IntegrationReflogEntry{}
	for i := 0; i < len(values); i += 3 {
		if !validOIDText(values[i]) {
			return nil, fmt.Errorf("malformed reflog object ID")
		}
		entries = append(entries, IntegrationReflogEntry{Ordinal: len(entries), OID: values[i], Selector: values[i+1], Action: values[i+2]})
	}
	if len(entries) < 2 {
		return nil, fmt.Errorf("parent reflog requires two parseable entries")
	}
	return entries, nil
}

func (git *systemTaskIntegrationGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	argv := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "merge.autoStash=false", "-c", "gc.auto=0", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", input.ParentWorktree, "merge", "--ff-only", "--no-stat", "--no-autostash", input.ResultOID}
	return git.run(argv, []string{"GIT_REFLOG_ACTION=ply-workspace-task-integrate:" + string(input.AttemptID)})
}
