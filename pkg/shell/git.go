package shell

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"strings"
)

func systemGitDependencies() process.Dependencies {
	return process.SystemRunner()
}

func GitClone(url string, target string) Output {
	return gitClone(systemGitDependencies(), url, target)
}

func gitClone(dependencies process.Dependencies, url string, target string) Output {
	return runGit(dependencies, "clone", url, target)
}

func GitPull(targetDir string) Output {
	return gitPull(systemGitDependencies(), targetDir)
}

func gitPull(dependencies process.Dependencies, targetDir string) Output {
	return runGit(dependencies, "-C", targetDir, "pull", "origin")
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
	return gitInit(systemGitDependencies(), targetDir)
}

func gitInit(dependencies process.Dependencies, targetDir string) Output {
	return runGit(dependencies, "-C", targetDir, "init")
}

func GitAddAndCommit(targetDir string, message string) Output {
	return gitAddAndCommit(systemGitDependencies(), targetDir, message)
}

func gitAddAndCommit(dependencies process.Dependencies, targetDir string, message string) Output {
	add := runGit(dependencies, "-C", targetDir, "add", ".")
	if add.Err != nil {
		return add
	}
	return runGit(dependencies, "-C", targetDir, "commit", "-m", fmt.Sprintf("\"%s\"", message))
}

func runGit(dependencies process.Dependencies, args ...string) (output Output) {
	command := process.Command{
		Name:   "git",
		Args:   args,
		Stdout: &output.StdOut,
		Stderr: &output.StdErr,
	}
	log.Debugf("running: %s %s", command.Name, strings.Join(command.Args, " "))
	if err := process.Execute(dependencies, command); err != nil {
		return output
	}
	return output
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
