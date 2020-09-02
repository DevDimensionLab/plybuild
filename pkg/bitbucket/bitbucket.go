package bitbucket

import (
	"co-pilot/pkg/git"
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"os"
)

var log = logger.Context()

func Pull(host string, workspace string, repository string) error {
	repoDir := workspace + repository

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return clone(host, workspace, repository)
	} else {
		return pull(workspace, repository)
	}
}

func clone(host string, workspace string, repository string) error {
	gitUrl := host + "/scm" + repository + ".git"
	toDir := workspace + repository

	log.Debugln("clone [" + gitUrl + "] -> [" + toDir + "]")
	return git.Clone(gitUrl, toDir)
}

func pull(workspace string, repository string) error {
	repoDir := workspace + "/" + repository

	log.Debugln(" pull [" + repoDir + "]")
	return git.Pull(repoDir)
}

func QueryProjects(host string, accessToken string) (*ProjectList, error) {
	response := ProjectList{}
	err := http.GetJsonWithAccessToken(host, "/rest/api/1.0/projects?limit=500", accessToken, &response)
	return &response, err
}

func QueryRepos(host string, projectKey string, accessToken string) (*ProjectRepos, error) {
	response := ProjectRepos{}
	err := http.GetJsonWithAccessToken(host, "/rest/api/1.0/projects/"+projectKey+"/repos?limit=1000", accessToken, &response)
	return &response, err
}
