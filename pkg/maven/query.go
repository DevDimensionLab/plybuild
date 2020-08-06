package maven

import (
	"encoding/xml"
	"fmt"
	"spring-boot-co-pilot/pkg/http"
	"strings"
)

type Metadata struct {
	XMLName      xml.Name `xml:"metadata"`
	Text         string   `xml:",chardata"`
	ModelVersion string   `xml:"modelVersion,attr"`
	GroupId      string   `xml:"groupId"`
	ArtifactId   string   `xml:"artifactId"`
	Versioning   struct {
		Text     string `xml:",chardata"`
		Latest   string `xml:"latest"`
		Release  string `xml:"release"`
		Versions struct {
			Text    string   `xml:",chardata"`
			Version []string `xml:"version"`
		} `xml:"versions"`
		LastUpdated string `xml:"lastUpdated"`
	} `xml:"versioning"`
}

func GetMetaData(groupID string, artifactId string) (Metadata, error) {
	var metaData Metadata
	// uses hardcoded url for now...
	url := fmt.Sprintf("https://repo1.maven.org/maven2/%s/%s/maven-metadata.xml",
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	err := http.GetXml(url, &metaData)
	if err != nil {
		return metaData, err
	}

	return metaData, nil
}
