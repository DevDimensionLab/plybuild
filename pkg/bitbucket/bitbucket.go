package bitbucket

import (
	"co-pilot/pkg/git"
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"errors"
	"fmt"
	"os"
	"strings"
)

var log = logger.Context()

func Synchronize(host string, accessToken string) error {
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

			err := Pull(host, ".", "/"+projectKey+"/"+bitBucketRepo.Name)
			if err != nil {
				log.Warnln(err)
			}
		}
	}

	return nil
}

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
	out, err := git.Clone(gitUrl, toDir)
	if err != nil {
		return errors.New(fmt.Sprintf("bitbucket clone failed %s, %v", out, err))
	}

	return nil
}

func pull(workspace string, repository string) error {
	repoDir := workspace + "/" + repository

	log.Debugln(" pull [" + repoDir + "]")
	out, err := git.Pull(repoDir)
	if err != nil {
		return errors.New(fmt.Sprintf("bitbucket pull failed %s, %v", out, err))
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
