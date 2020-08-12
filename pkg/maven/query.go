package maven

import (
	"co-pilot/pkg/http"
	"errors"
	"fmt"
	"strings"
)

func GetMetaData(groupID string, artifactId string) (Metadata, error) {
	var metaData Metadata
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
