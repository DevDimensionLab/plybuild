package bitbucket

import (
	"github.com/devdimensionlab/ply/pkg/file"
	"github.com/devdimensionlab/ply/pkg/http"
	"github.com/devdimensionlab/ply/pkg/shell"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Bitbucket struct {
	host        string
	accessToken string
	log         logrus.FieldLogger
}

func With(logger logrus.FieldLogger, host string, accessToken string) Bitbucket {
	return Bitbucket{
		host:        host,
		accessToken: accessToken,
		log:         logger,
	}
}

func (bitbucket Bitbucket) SynchronizeAllRepos(excludeProjects []string) error {
	projects, err := bitbucket.queryProjects()
	if err != nil {
		return err
	}

	for _, bitBucketProject := range projects.Values {
		log.Debugf("Starting to synchronize: %s", bitBucketProject.Key)
		if skipProject(bitBucketProject.Key, excludeProjects) {
			continue
		}

		projectKey := strings.ToLower(bitBucketProject.Key)
		bitbucket.log.Infoln("project: " + projectKey)

		bitBucketProjectReposResponse, err := QueryRepos(bitbucket.host, projectKey, bitbucket.accessToken)
		if err != nil {
			bitbucket.log.Warnln(err)
		}

		for _, bitBucketRepo := range bitBucketProjectReposResponse.BitBucketRepo {
			bitbucket.log.Infoln("  " + bitBucketRepo.Name)

			err := bitbucket.cloneOrPull(".", "/"+projectKey+"/"+bitBucketRepo.Name)
			if err != nil {
				bitbucket.log.Warnln(err)
			}
		}
	}

	return nil
}

func skipProject(key string, excludeProjects []string) bool {
	for _, exclude := range excludeProjects {
		log.Debugf("Checking against excluded project: %s", exclude)
		if strings.ToLower(key) == strings.ToLower(exclude) {
			return true
		}
	}
	return false
}

func (bitbucket Bitbucket) cloneOrPull(workspace string, repository string) error {
	repoDir := workspace + repository

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return bitbucket.clone(workspace, repository)
	} else {
		return bitbucket.pull(workspace, repository)
	}
}

func (bitbucket Bitbucket) clone(workspace string, repository string) error {
	gitUrl := bitbucket.host + "/scm" + repository + ".git"
	toDir := workspace + repository

	bitbucket.log.Debugln("clone [" + gitUrl + "] -> [" + toDir + "]")
	clone := shell.GitClone(gitUrl, toDir)
	if clone.Err != nil {
		return clone.FormatError()
	}

	return nil
}

func (bitbucket Bitbucket) pull(workspace string, repository string) error {
	repoDir := file.Path("%s/%s", workspace, repository)

	bitbucket.log.Debugln(" pull [" + repoDir + "]")
	pull := shell.GitPull(repoDir)
	if pull.Err != nil {
		return pull.FormatError()
	}

	return nil
}

func (bitbucket Bitbucket) queryProjects() (*ProjectList, error) {
	response := ProjectList{}
	err := http.GetJsonWithAccessToken(bitbucket.host, "/rest/api/1.0/projects?limit=500", bitbucket.accessToken, &response)
	return &response, err
}

func QueryRepos(host string, projectKey string, accessToken string) (*ProjectRepos, error) {
	response := ProjectRepos{}
	err := http.GetJsonWithAccessToken(host, "/rest/api/1.0/projects/"+projectKey+"/repos?limit=1000", accessToken, &response)
	return &response, err
}
