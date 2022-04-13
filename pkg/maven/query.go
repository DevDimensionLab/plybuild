package maven

import (
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"github.com/co-pilot-cli/co-pilot/pkg/http"
	"github.com/co-pilot-cli/mvn-pom-mutator/pkg/pom"
	"strings"
)

func GetMetaData(groupID string, artifactId string) (metaData RepositoryMetadata, err error) {
	settings, _ := NewSettings()
	repos, err := settings.GetRepositories()
	if err != nil {
		return metaData, err
	}
	repo := repos.GetDefaultRepository()

	url := file.Path("%s/%s/%s/maven-metadata.xml",
		repo.Url,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	log.Debugf("Using url for metadata: %s", url)

	if repo.Auth != nil {
		err = http.GetAuthXml(url, repo.Auth.Username, repo.Auth.Password, &metaData)
	} else {
		err = http.GetXml(url, &metaData)
	}

	if err != nil {
		return metaData, err
	}

	return metaData, nil
}

func GetBannedModel(url string) (*pom.Model, error) {
	var bannedModel pom.Model
	err := http.GetXml(url, &bannedModel)
	if err != nil {
		return &bannedModel, err
	}

	return &bannedModel, nil
}
