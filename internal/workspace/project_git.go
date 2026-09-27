package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type RepoObservation struct {
	Locator      string
	GitCommonDir string
}

type RepoObserver interface {
	Observe(path string) (RepoObservation, error)
}

type gitCommandResult struct {
	stdout []byte
	stderr []byte
	exit   int
	err    error
}

type systemRepoObserver struct {
	files FileSystem
	run   func(string, ...string) gitCommandResult
}

func newSystemRepoObserver(files FileSystem) RepoObserver {
	return systemRepoObserver{files: files, run: runGitCommand}
}

func (observer systemRepoObserver) Observe(path string) (RepoObservation, error) {
	inside := observer.run(path, "rev-parse", "--is-inside-work-tree")
	if inside.err != nil {
		var exitError *exec.ExitError
		if errors.As(inside.err, &exitError) {
			return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), inside.err)
		}
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("run git for %s: %v", path, inside.err), inside.err)
	}
	insideValue, ok := singleGitLine(inside.stdout)
	if !ok || insideValue != "true" {
		if ok && insideValue == "false" {
			return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), nil)
		}
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe worktree status for %s: malformed Git output", path), nil)
	}

	topLevel, err := observer.gitPath(path, "--show-toplevel")
	if err != nil {
		return RepoObservation{}, err
	}
	if topLevel != path {
		return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), nil)
	}
	commonDirectory, err := observer.gitPath(path, "--git-common-dir")
	if err != nil {
		return RepoObservation{}, err
	}
	return RepoObservation{Locator: topLevel, GitCommonDir: commonDirectory}, nil
}

func (observer systemRepoObserver) gitPath(path, argument string) (string, error) {
	result := observer.run(path, "rev-parse", "--path-format=absolute", argument)
	if result.err != nil {
		return "", projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe %s for %s: git exit %d", argument, path, result.exit), result.err)
	}
	value, ok := singleGitLine(result.stdout)
	if !ok || value == "" {
		return "", projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe %s for %s: malformed Git output", argument, path), nil)
	}
	physical, err := canonicalDirectory(observer.files, "", value, "Git path")
	if err != nil {
		return "", projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe %s for %s: %v", argument, path, err), err)
	}
	return physical, nil
}

func runGitCommand(path string, arguments ...string) gitCommandResult {
	commandArguments := append([]string{"-C", path}, arguments...)
	command := exec.Command("git", commandArguments...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	exit := 0
	if err != nil {
		exit = -1
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exit = exitError.ExitCode()
		}
	}
	return gitCommandResult{stdout: stdout.Bytes(), stderr: stderr.Bytes(), exit: exit, err: err}
}

func singleGitLine(output []byte) (string, bool) {
	value := string(output)
	switch {
	case strings.HasSuffix(value, "\r\n"):
		value = strings.TrimSuffix(value, "\r\n")
	case strings.HasSuffix(value, "\n"):
		value = strings.TrimSuffix(value, "\n")
	default:
		return "", false
	}
	if strings.ContainsAny(value, "\r\n") {
		return "", false
	}
	return value, true
}
