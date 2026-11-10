package maven

import (
	"encoding/xml"
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
)

const mavenHTTPMutation = "4. Maven metadata keeps username before password"

type recordedMavenMetadataRequest struct {
	Request     httpclient.Request
	Destination interface{}
}

type recordingMavenMetadataHTTP struct {
	requests []recordedMavenMetadataRequest
	response string
	err      error
}

func (recording *recordingMavenMetadataHTTP) dependencies() metadataHTTPDependencies {
	return metadataHTTPDependencies{Client: recording}
}

func (recording *recordingMavenMetadataHTTP) GetXML(request httpclient.Request, parsed interface{}) error {
	request.BasicAuth = copyBasicAuth(request.BasicAuth)
	recording.requests = append(recording.requests, recordedMavenMetadataRequest{
		Request:     request,
		Destination: parsed,
	})
	if recording.err != nil {
		return recording.err
	}
	return xml.Unmarshal([]byte(recording.response), parsed)
}

func (recording *recordingMavenMetadataHTTP) assertedRequests() ([]recordedMavenMetadataRequest, error) {
	if len(recording.requests) == 0 {
		return nil, errors.New("recorded Maven metadata request population is empty")
	}
	return recording.requests, nil
}

func copyBasicAuth(auth *httpclient.BasicAuth) *httpclient.BasicAuth {
	if auth == nil {
		return nil
	}
	copied := *auth
	return &copied
}

func TestGetRepo(t *testing.T) {
	const fallbackURL = "https://repo.example.test/maven2"
	repos := Repositories{Fallback: Repository{Url: fallbackURL}}
	localRepo, err := repos.GetDefaultRepository()
	if err != nil {
		t.Fatalf("get fallback repository: %v", err)
	}

	if localRepo.Url != fallbackURL {
		t.Errorf("default repository URL was %q, want %q", localRepo.Url, fallbackURL)
	}
}

func TestMavenMetadataKeepsUsernameBeforePasswordAndParsesResponse(t *testing.T) {
	if mavenHTTPMutation == "" {
		t.Fatal("maven-http mutation label is empty")
	}
	recording := &recordingMavenMetadataHTTP{
		response: `<metadata><groupId>com.example.complete</groupId><artifactId>complete.artifact</artifactId><versioning><release>3.2.1-complete</release></versioning></metadata>`,
	}
	repository := Repository{
		Url: "https://repo.example.invalid/complete-root",
		Auth: &RepositoryAuth{
			Username: "complete-user-before-password",
			Password: "complete-password-after-user",
		},
	}

	metadata, err := repository.getMetaData(recording.dependencies(), "com.example.complete", "complete.artifact")

	if err != nil {
		t.Fatalf("authenticated metadata request returned an error: %v", err)
	}
	wantRequest := httpclient.Request{
		URL: "https://repo.example.invalid/complete-root/com/example/complete/complete/artifact/maven-metadata.xml",
		BasicAuth: &httpclient.BasicAuth{
			Username: "complete-user-before-password",
			Password: "complete-password-after-user",
		},
	}
	assertRecordedMavenMetadataRequests(t, recording, []httpclient.Request{wantRequest})
	if metadata.GroupId != "com.example.complete" || metadata.ArtifactId != "complete.artifact" ||
		metadata.Versioning.Release != "3.2.1-complete" {
		t.Fatalf("parsed metadata was incomplete: %#v", metadata)
	}
}

func TestMavenMetadataWithoutAuthSelectsAnonymousRequest(t *testing.T) {
	recording := &recordingMavenMetadataHTTP{
		response: `<metadata><versioning><release>1.2.3-anonymous</release></versioning></metadata>`,
	}
	repository := Repository{Url: "https://anonymous.example.invalid/maven"}

	metadata, err := repository.getMetaData(recording.dependencies(), "org.example", "anonymous-artifact")

	if err != nil {
		t.Fatalf("anonymous metadata request returned an error: %v", err)
	}
	assertRecordedMavenMetadataRequests(t, recording, []httpclient.Request{{
		URL: "https://anonymous.example.invalid/maven/org/example/anonymous-artifact/maven-metadata.xml",
	}})
	if metadata.Versioning.Release != "1.2.3-anonymous" {
		t.Fatalf("anonymous metadata release was %q", metadata.Versioning.Release)
	}
}

func TestMavenMetadataReturnsDependencyErrorAfterCompleteDelivery(t *testing.T) {
	sentinel := errors.New("complete metadata dependency error")
	recording := &recordingMavenMetadataHTTP{err: sentinel}
	repository := Repository{
		Url:  "https://errors.example.invalid/repository",
		Auth: &RepositoryAuth{Username: "error-user", Password: "error-password"},
	}

	_, err := repository.getMetaData(recording.dependencies(), "error.group", "error.artifact")

	if !errors.Is(err, sentinel) {
		t.Fatalf("metadata dependency error was %v, want %v", err, sentinel)
	}
	assertRecordedMavenMetadataRequests(t, recording, []httpclient.Request{{
		URL: "https://errors.example.invalid/repository/error/group/error/artifact/maven-metadata.xml",
		BasicAuth: &httpclient.BasicAuth{
			Username: "error-user",
			Password: "error-password",
		},
	}})
}

func TestMavenMetadataDependenciesDefaultToNoRequest(t *testing.T) {
	repository := Repository{Url: "https://must-not-request.example.invalid"}

	_, err := repository.getMetaData(metadataHTTPDependencies{}, "safe.default", "no-network")

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe metadata default returned %v, want %v", err, httpclient.ErrNoClient)
	}
}

func TestRecordedMavenMetadataRequestsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingMavenMetadataHTTP{}

	if _, err := recording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Maven metadata population passed")
	}
}

func assertRecordedMavenMetadataRequests(t *testing.T, recording *recordingMavenMetadataHTTP, want []httpclient.Request) {
	t.Helper()
	requests, err := recording.assertedRequests()
	if err != nil {
		t.Fatal(err)
	}
	got := make([]httpclient.Request, len(requests))
	for index, request := range requests {
		got[index] = request.Request
		if _, ok := request.Destination.(*RepositoryMetadata); !ok {
			t.Fatalf("metadata dependency destination was %T, want *RepositoryMetadata", request.Destination)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recorded Maven metadata requests differ:\n got: %#v\nwant: %#v", got, want)
	}
}
