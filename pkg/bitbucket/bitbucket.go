package bitbucket

import (
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"os"
	"strings"
)

var log = logger.Context()

func SynchronizeAllRepos(host string, accessToken string) error {
	projects, err := QueryProjects(host, accessToken)
	if err != nil {
		return err
	}

	for _, bitBucketProject := range projects.Values {
		projectKey := strings.ToLower(bitBucketProject.Key)
		log.Infoln("project: " + projectKey)

		bitBucketProjectReposResponse, err := QueryRepos(host, projectKey, accessToken)
		if err != nil {
			log.Warnln(err)
		}

		for _, bitBucketRepo := range bitBucketProjectReposResponse.BitBucketRepo {
			log.Infoln("  " + bitBucketRepo.Name)

			err := cloneOrPull(host, ".", "/"+projectKey+"/"+bitBucketRepo.Name)
			if err != nil {
				log.Warnln(err)
			}
		}
	}

	return nil
}

func cloneOrPull(host string, workspace string, repository string) error {
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
	clone := shell.GitClone(gitUrl, toDir)
	if clone.Err != nil {
		return clone.FormatError()
	}

	return nil
}

func pull(workspace string, repository string) error {
	repoDir := workspace + "/" + repository

	log.Debugln(" pull [" + repoDir + "]")
	pull := shell.GitPull(repoDir)
	if pull.Err != nil {
		return pull.FormatError()
	}

	return nil
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
