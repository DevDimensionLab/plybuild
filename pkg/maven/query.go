package maven

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	stdhttp "net/http"
	"net/url"
	"strings"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
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
	response, err := httpclient.Execute(httpclient.System(), request)
	if err != nil {
		reason := "metadata transport error"
		var networkError net.Error
		if errors.As(err, &networkError) && networkError.Timeout() {
			reason += " (timeout)"
		}
		return metadataDiagnostic{message: reason, cause: err}
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode >= 400 {
		// Use the status code, not server-controlled status text or response body.
		return metadataDiagnostic{message: fmt.Sprintf("HTTP %d %s", response.StatusCode, stdhttp.StatusText(response.StatusCode))}
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return metadataDiagnostic{message: "metadata transport error reading response", cause: err}
	}
	if err := xml.Unmarshal(body, parsed); err != nil {
		// XML errors may quote response content. Retain the cause for callers,
		// but never print that content in an upgrade warning.
		return metadataDiagnostic{message: "invalid XML metadata", cause: err}
	}
	return nil
}

// metadataDiagnostic preserves error identity without exposing credentials,
// request headers, or server-controlled content through Error().
type metadataDiagnostic struct {
	message string
	cause   error
}

func (err metadataDiagnostic) Error() string { return err.message }
func (err metadataDiagnostic) Unwrap() error { return err.cause }

func metadataDiagnosticURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Opaque != "" {
		return "[invalid repository URL]"
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.ForceQuery = false
	parsed.Fragment = ""
	parsed.RawFragment = ""
	return parsed.String()
}

func (repository Repository) metadataError(groupID, artifactID, reason string, cause error) error {
	return metadataDiagnostic{
		message: fmt.Sprintf("could not determine release for %s:%s from %s: %s", groupID, artifactID, metadataDiagnosticURL(repository.Url), reason),
		cause:   cause,
	}
}

func (repository Repository) latestRelease(groupID, artifactID string) (JavaVersion, error) {
	metadata, err := repository.GetMetaData(groupID, artifactID)
	if err != nil {
		return JavaVersion{}, err
	}
	release, err := metadata.LatestRelease()
	if err != nil {
		reason := "invalid version in metadata"
		if errors.Is(err, errNoSuitableRelease) {
			reason = errNoSuitableRelease.Error()
		}
		return JavaVersion{}, repository.metadataError(groupID, artifactID, reason, err)
	}
	return release, nil
}

func (repository Repository) getMetaData(dependencies metadataHTTPDependencies, groupID string, artifactId string) (metaData RepositoryMetadata, err error) {
	repo := repository

	url := file.Path("%s/%s/%s/maven-metadata.xml",
		repo.Url,
		strings.ReplaceAll(groupID, ".", "/"),
		strings.ReplaceAll(artifactId, ".", "/"))
	log.Debugf("using url for metadata: %s", metadataDiagnosticURL(url))

	request := httpclient.Request{URL: url}
	if repo.Auth != nil {
		request.BasicAuth = &httpclient.BasicAuth{
			Username: repo.Auth.Username,
			Password: repo.Auth.Password,
		}
	}
	err = dependencies.GetXML(request, &metaData)

	if err != nil {
		reason := "metadata lookup failed"
		var diagnostic metadataDiagnostic
		if errors.As(err, &diagnostic) {
			reason = diagnostic.Error()
		}
		err = repository.metadataError(groupID, artifactId, reason, err)
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
