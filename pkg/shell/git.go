package shell

import (
	"co-pilot/pkg/logger"
	"os/exec"
	"strings"
)

var log = logger.Context()

func GitClone(url string, target string) (string, error) {
	return run(exec.Command("git", "clone", url, target))
}

func GitPull(targetDir string) (string, error) {
	return run(exec.Command("git", "-C", targetDir, "pull", "origin"))
}

func GitDirty(targetDir string) (bool, error) {
	cmd, err := Run("git", "-C", targetDir, "diff", "--stat")
	if err != nil {
		return true, err
	}

	log.Debugf("%s is git-dirty = %t", targetDir, cmd != "")
	return cmd != "", nil
}

func GitIsRepo(targetDir string) (bool, error) {
	cmd, err := Run("git", "-C", targetDir, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return true, err
	}

	isRepo := strings.TrimSpace(cmd) == "true"
	log.Debugf("%s is a git-repo = %t", targetDir, isRepo)
	return isRepo, nil
}
