package maven

import (
	"github.com/devdimensionlab/ply/pkg/file"
	"github.com/devdimensionlab/ply/pkg/http"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"strings"
)

func (repository Repository) GetMetaData(groupID string, artifactId string) (metaData RepositoryMetadata, err error) {
	repo := repository

	url := file.Path("%s/%s/%s/maven-metadata.xml",
		repo.Url,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	log.Debugf("using url for metadata: %s", url)

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
