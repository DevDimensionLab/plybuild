package shell

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"os/exec"
	"strings"
)

func GitClone(url string, target string) Output {
	return run(exec.Command("git", "clone", url, target))
}

func GitPull(targetDir string) Output {
	return run(exec.Command("git", "-C", targetDir, "pull", "origin"))
}

func GitDirty(targetDir string) (bool, error) {
	output := Run("git", "-C", targetDir, "diff", "--stat")
	if output.Err != nil {
		return true, output.Err
	}

	cmd := output.StdOut.String()
	log.Debugf("%s is git-dirty = %t", targetDir, cmd != "")
	return cmd != "", nil
}

func GitIsRepo(targetDir string) (bool, error) {
	output := Run("git", "-C", targetDir, "rev-parse", "--is-inside-work-tree")
	if output.Err != nil {
		return true, output.Err
	}

	cmd := output.StdOut.String()
	isRepo := strings.TrimSpace(cmd) == "true"
	log.Debugf("%s is a git-repo = %t", targetDir, isRepo)
	return isRepo, nil
}

func GitInit(targetDir string) Output {
	return run(exec.Command("git", "-C", targetDir, "init"))
}

func GitAddAndCommit(targetDir string, message string) Output {
	add := Run("git", "-C", targetDir, "add", ".")
	if add.Err != nil {
		return add
	}
	return Run("git", "-C", targetDir, "commit", "-m", fmt.Sprintf("\"%s\"", message))
}

func InstallGitHooks(sourceDir string, sourceFileNames []string, targetDir string) error {
	if sourceFileNames == nil {
		return nil
	}

	for _, sourceFileName := range sourceFileNames {
		hooksFile := file.Path("%s/.git/hooks/%s", targetDir, sourceFileName)
		sourceFilePath := file.Path("%s/%s", sourceDir, sourceFileName)
		log.Debugf("Copying %s into %s", sourceFilePath, hooksFile)
		err := file.CopyFile(sourceFilePath, hooksFile)
		if err != nil {
			return err
		}
	}
	return nil
}
