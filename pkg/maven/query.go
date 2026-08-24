package maven

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
	"strings"
)

func (repository Repository) GetMetaData(groupID string, artifactId string) (metaData RepositoryMetadata, err error) {
	return repository.getMetaData(metadataHTTPDependencies{Client: packageMetadataHTTP{}}, groupID, artifactId)
}

type metadataHTTPClient interface {
	GetXML(httpclient.Request, interface{}) error
}

type metadataHTTPDependencies struct {
	Client metadataHTTPClient
}

func (dependencies metadataHTTPDependencies) GetXML(request httpclient.Request, parsed interface{}) error {
	if dependencies.Client == nil {
		return httpclient.ErrNoClient
	}
	return dependencies.Client.GetXML(request, parsed)
}

type packageMetadataHTTP struct{}

func (packageMetadataHTTP) GetXML(request httpclient.Request, parsed interface{}) error {
	if request.BasicAuth != nil {
		return http.GetAuthXml(request.URL, request.BasicAuth.Username, request.BasicAuth.Password, parsed)
	}
	return http.GetXml(request.URL, parsed)
}

func (repository Repository) getMetaData(dependencies metadataHTTPDependencies, groupID string, artifactId string) (metaData RepositoryMetadata, err error) {
	repo := repository

	url := file.Path("%s/%s/%s/maven-metadata.xml",
		repo.Url,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	log.Debugf("using url for metadata: %s", url)

	request := httpclient.Request{URL: url}
	if repo.Auth != nil {
		request.BasicAuth = &httpclient.BasicAuth{
			Username: repo.Auth.Username,
			Password: repo.Auth.Password,
		}
	}
	err = dependencies.GetXML(request, &metaData)

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
