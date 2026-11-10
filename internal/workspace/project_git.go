package workspace

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
)

type RepoObservation struct {
	Locator      string
	GitCommonDir string
}

type RepoObserver interface {
	Observe(path string) (RepoObservation, error)
}

type systemRepoObserver struct {
	files  FileSystem
	runner process.Runner
}

func newSystemRepoObserver(files FileSystem) RepoObserver {
	return systemRepoObserver{files: files, runner: process.SystemRunner().Runner}
}

func (observer systemRepoObserver) Observe(path string) (RepoObservation, error) {
	identity, err := gitrepo.Identify(context.Background(), observer.runner, path)
	var identifyError *gitrepo.Error
	var exitError *exec.ExitError
	if errors.As(err, &identifyError) && identifyError.Op == gitrepo.OpInsideWorkTree && errors.As(identifyError.Err, &exitError) {
		return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), err)
	}
	if err != nil && (identifyError == nil || identifyError.Op != gitrepo.OpCommonDir) {
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe Git identity for %s: %v", path, err), err)
	}
	if !identity.Worktree {
		return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), nil)
	}
	topLevel, canonicalError := canonicalDirectory(observer.files, "", identity.Toplevel, "Git path")
	if canonicalError != nil {
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe toplevel for %s: %v", path, canonicalError), canonicalError)
	}
	if topLevel != path {
		return RepoObservation{}, projectError(ErrorProjectRepoInvalid, fmt.Sprintf("%s is not a Git worktree root", path), nil)
	}
	if err != nil {
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe Git identity for %s: %v", path, err), err)
	}
	commonDirectory, canonicalError := canonicalDirectory(observer.files, "", identity.CommonDir, "Git path")
	if canonicalError != nil {
		return RepoObservation{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe common dir for %s: %v", path, canonicalError), canonicalError)
	}
	return RepoObservation{Locator: topLevel, GitCommonDir: commonDirectory}, nil
}

func singleGitLine(output []byte) (string, bool) {
	return gitrepo.SingleLine(output)
}
