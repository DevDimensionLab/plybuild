package maven

import (
	"co-pilot/pkg/http"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

func GetMetaData(groupID string, artifactId string) (RepositoryMetadata, error) {
	var metaData RepositoryMetadata
	repos, err := GetRepositories()
	if err != nil {
		return metaData, err
	}
	repo, err := GetRepo(repos)
	if err != nil {
		return metaData, err
	}

	url := fmt.Sprintf("%s/%s/%s/maven-metadata.xml",
		repo,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	log.Debugf("Using url for metadata: %s", url)

	err = http.GetXml(url, &metaData)

	if err != nil {
		return metaData, err
	}

	return metaData, nil
}

func GetRepo(repos Repositories) (string, error) {
	var repo = ""
	if len(repos.Mirror) > 0 {
		repo = repos.Mirror[0]
	} else {
		repo = repos.Fallback
	}

	if repo == "" {
		return "", errors.New("could not find a valid maven repo in repos struct")
	} else {
		return repo, nil
	}
}

func GetBannedModel(url string) (*pom.Model, error) {
	var bannedModel pom.Model
	err := http.GetXml(url, &bannedModel)
	if err != nil {
		return &bannedModel, err
	}

	return &bannedModel, nil
}
