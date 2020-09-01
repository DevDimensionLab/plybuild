package git

import (
	"co-pilot/pkg/logger"
	"os"
	"os/exec"
)

var log = logger.Context()

func PullRepo(host string, workspace string, repository string) error {
	repoDir := workspace + repository

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return GitClone(host, workspace, repository)
	} else {
		return GitPull(workspace, repository)
	}
}

func GitClone(host string, workspace string, repository string) error {
	gitUrl := host + "/scm" + repository + ".git"
	toDir := workspace + repository

	log.Debugln("GitClone [" + gitUrl + "] -> [" + toDir + "]")

	err := exec.Command("git", "clone", gitUrl, toDir).Run()
	return err
}

func GitPull(workspace string, repository string) error {
	repoDir := workspace + "/" + repository
	log.Debugln(" GitPull [" + repoDir + "]")

	err := exec.Command("git", "-C", repoDir, "pull", "origin").Run()
	return err
}
