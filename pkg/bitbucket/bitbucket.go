package bitbucket

import (
	"os"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
	"github.com/devdimensionlab/plybuild/pkg/shell"
	"github.com/sirupsen/logrus"
)

type Bitbucket struct {
	host        string
	accessToken string
	log         logrus.FieldLogger
	queries     queryDependencies
}

func With(logger logrus.FieldLogger, host string, accessToken string) Bitbucket {
	return Bitbucket{
		host:        host,
		accessToken: accessToken,
		log:         logger,
		queries:     systemQueryDependencies(),
	}
}

type queryHTTPClient interface {
	GetBearerJSON(httpclient.Request, interface{}) error
}

type queryDependencies struct {
	Client queryHTTPClient
}

func (dependencies queryDependencies) GetBearerJSON(request httpclient.Request, parsed interface{}) error {
	if dependencies.Client == nil {
		return httpclient.ErrNoClient
	}
	return dependencies.Client.GetBearerJSON(request, parsed)
}

func systemQueryDependencies() queryDependencies {
	return queryDependencies{Client: packageQueryHTTP{}}
}

type packageQueryHTTP struct{}

func (packageQueryHTTP) GetBearerJSON(request httpclient.Request, parsed interface{}) error {
	accessToken := ""
	if request.BearerJSON != nil {
		accessToken = request.BearerJSON.AccessToken
	}
	return http.GetJsonWithAccessToken("", request.URL, accessToken, parsed)
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

		bitBucketProjectReposResponse, err := bitbucket.queryRepos(projectKey)
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
		if strings.EqualFold(key, exclude) {
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
	err := bitbucket.queries.GetBearerJSON(httpclient.Request{
		URL: bitbucket.host + "/rest/api/1.0/projects?limit=500",
		BearerJSON: &httpclient.BearerJSON{
			AccessToken: bitbucket.accessToken,
		},
	}, &response)
	return &response, err
}

func QueryRepos(host string, projectKey string, accessToken string) (*ProjectRepos, error) {
	return Bitbucket{
		host:        host,
		accessToken: accessToken,
		queries:     systemQueryDependencies(),
	}.queryRepos(projectKey)
}

func (bitbucket Bitbucket) queryRepos(projectKey string) (*ProjectRepos, error) {
	response := ProjectRepos{}
	err := bitbucket.queries.GetBearerJSON(httpclient.Request{
		URL: bitbucket.host + "/rest/api/1.0/projects/" + projectKey + "/repos?limit=1000",
		BearerJSON: &httpclient.BearerJSON{
			AccessToken: bitbucket.accessToken,
		},
	}, &response)
	return &response, err
}
