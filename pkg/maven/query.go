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

	defaultRepo := repos[0]
	if defaultRepo == "" {
		return metaData, errors.New("could not find a maven repo")
	}

	url := fmt.Sprintf("%s/%s/%s/maven-metadata.xml",
		defaultRepo,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	err = http.GetXml(url, &metaData)
	if err != nil {
		return metaData, err
	}

	return metaData, nil
}
