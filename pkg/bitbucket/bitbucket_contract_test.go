package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/sirupsen/logrus"
)

type recordedBitbucketQuery struct {
	Request     httpclient.Request
	Destination interface{}
}

type recordingBitbucketQueryHTTP struct {
	queries   []recordedBitbucketQuery
	responses map[string]string
	errors    map[string]error
}

func (recording *recordingBitbucketQueryHTTP) dependencies() queryDependencies {
	return queryDependencies{Client: recording}
}

func (recording *recordingBitbucketQueryHTTP) GetBearerJSON(request httpclient.Request, parsed interface{}) error {
	if request.BasicAuth != nil {
		copied := *request.BasicAuth
		request.BasicAuth = &copied
	}
	if request.BearerJSON != nil {
		copied := *request.BearerJSON
		request.BearerJSON = &copied
	}
	recording.queries = append(recording.queries, recordedBitbucketQuery{
		Request:     request,
		Destination: parsed,
	})
	if err := recording.errors[request.URL]; err != nil {
		return err
	}
	return json.Unmarshal([]byte(recording.responses[request.URL]), parsed)
}

func (recording *recordingBitbucketQueryHTTP) assertedQueries() ([]recordedBitbucketQuery, error) {
	if len(recording.queries) == 0 {
		return nil, errors.New("recorded Bitbucket query population is empty")
	}
	return recording.queries, nil
}

func TestBitbucketQueriesKeepCompleteURLsBearerValuesAndParsedResponses(t *testing.T) {
	host := "https://bitbucket.example.invalid:8443/complete-base"
	token := "complete-bitbucket-token"
	projectPath := "/rest/api/1.0/projects?limit=500"
	repositoryPath := "/rest/api/1.0/projects/complete-project/repos?limit=1000"
	recording := &recordingBitbucketQueryHTTP{responses: map[string]string{
		host + projectPath: `{
			"size":1,"limit":500,"isLastPage":true,"start":0,"nextPageStart":1,
			"values":[{"key":"COMPLETE-PROJECT","id":41,"name":"Complete Project","public":true,"type":"NORMAL"}]
		}`,
		host + repositoryPath: `{
			"size":1,"limit":1000,"isLastPage":true,"start":0,
			"values":[{"slug":"complete-repo","id":73,"name":"Complete Repo","scmId":"git","state":"AVAILABLE","project":{"key":"COMPLETE-PROJECT"}}]
		}`,
	}}
	client := With(logrus.New(), host, token)
	client.queries = recording.dependencies()

	projects, err := client.queryProjects()
	if err != nil {
		t.Fatalf("Bitbucket project query returned an error: %v", err)
	}
	repositories, err := client.queryRepos("complete-project")
	if err != nil {
		t.Fatalf("Bitbucket repository query returned an error: %v", err)
	}

	if projects.Size != 1 || projects.Limit != 500 || !projects.IsLastPage ||
		len(projects.Values) != 1 || projects.Values[0].Key != "COMPLETE-PROJECT" ||
		projects.Values[0].ID != 41 || projects.Values[0].Name != "Complete Project" ||
		!projects.Values[0].Public || projects.Values[0].Type != "NORMAL" {
		t.Fatalf("parsed Bitbucket projects response was incomplete: %#v", projects)
	}
	if repositories.Size != 1 || repositories.Limit != 1000 || !repositories.IsLastPage ||
		len(repositories.BitBucketRepo) != 1 || repositories.BitBucketRepo[0].Slug != "complete-repo" ||
		repositories.BitBucketRepo[0].ID != 73 || repositories.BitBucketRepo[0].Name != "Complete Repo" ||
		repositories.BitBucketRepo[0].ScmID != "git" || repositories.BitBucketRepo[0].State != "AVAILABLE" ||
		repositories.BitBucketRepo[0].Project.Key != "COMPLETE-PROJECT" {
		t.Fatalf("parsed Bitbucket repositories response was incomplete: %#v", repositories)
	}
	assertRecordedBitbucketQueries(t, recording, []recordedBitbucketQuery{
		{
			Request: httpclient.Request{
				URL: host + projectPath,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			},
			Destination: (*ProjectList)(nil),
		},
		{
			Request: httpclient.Request{
				URL: host + repositoryPath,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			},
			Destination: (*ProjectRepos)(nil),
		},
	})
}

func TestBitbucketQueriesReturnDependencyErrorsAfterCompleteDelivery(t *testing.T) {
	host := "https://errors.bitbucket.example.invalid/complete-base"
	token := "complete-error-token"
	projectURL := host + "/rest/api/1.0/projects?limit=500"
	repositoryURL := host + "/rest/api/1.0/projects/error-project/repos?limit=1000"
	sentinel := errors.New("Bitbucket query dependency failed")
	tests := []struct {
		name        string
		url         string
		destination interface{}
		invoke      func(Bitbucket) error
	}{
		{
			name:        "projects",
			url:         projectURL,
			destination: (*ProjectList)(nil),
			invoke: func(client Bitbucket) error {
				_, err := client.queryProjects()
				return err
			},
		},
		{
			name:        "repositories",
			url:         repositoryURL,
			destination: (*ProjectRepos)(nil),
			invoke: func(client Bitbucket) error {
				_, err := client.queryRepos("error-project")
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingBitbucketQueryHTTP{errors: map[string]error{test.url: sentinel}}
			client := With(logrus.New(), host, token)
			client.queries = recording.dependencies()

			err := test.invoke(client)

			if !errors.Is(err, sentinel) {
				t.Fatalf("Bitbucket query dependency error was %v, want %v", err, sentinel)
			}
			queries, assertErr := recording.assertedQueries()
			if assertErr != nil {
				t.Fatal(assertErr)
			}
			wantRequest := httpclient.Request{
				URL: test.url,
				BearerJSON: &httpclient.BearerJSON{
					AccessToken: token,
				},
			}
			if len(queries) != 1 || !reflect.DeepEqual(queries[0].Request, wantRequest) {
				t.Fatalf("Bitbucket error query was %#v, want %#v", queries, wantRequest)
			}
			if reflect.TypeOf(queries[0].Destination) != reflect.TypeOf(test.destination) {
				t.Fatalf("Bitbucket error destination was %T, want %T", queries[0].Destination, test.destination)
			}
		})
	}
}

func TestBitbucketQueryDependenciesDefaultToSafeNoRequest(t *testing.T) {
	client := Bitbucket{
		host:        "https://must-not-request.example.invalid",
		accessToken: "must-not-request-token",
		log:         logrus.New(),
	}

	if _, err := client.queryProjects(); !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Bitbucket project-query default returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if _, err := client.queryRepos("must-not-request"); !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Bitbucket repository-query default returned %v, want %v", err, httpclient.ErrNoClient)
	}
}

func TestBitbucketSynchronizationSelectsProjectThenLowercaseRepositoryQueryAndWarns(t *testing.T) {
	host := "https://sync.bitbucket.example.invalid"
	token := "synchronization-token"
	projectURL := host + "/rest/api/1.0/projects?limit=500"
	repositoryURL := host + "/rest/api/1.0/projects/mixed-project/repos?limit=1000"
	sentinel := errors.New("repository query warning sentinel")
	recording := &recordingBitbucketQueryHTTP{
		responses: map[string]string{
			projectURL: `{"values":[{"key":"MiXeD-PrOjEcT"}]}`,
		},
		errors: map[string]error{repositoryURL: sentinel},
	}
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetOutput(logOutput)
	client := With(logger, host, token)
	client.queries = recording.dependencies()

	err := client.SynchronizeAllRepos(nil)

	if err != nil {
		t.Fatalf("Bitbucket synchronization returned an error: %v", err)
	}
	queries, assertErr := recording.assertedQueries()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	if len(queries) != 2 || queries[0].Request.URL != projectURL || queries[1].Request.URL != repositoryURL {
		t.Fatalf("Bitbucket synchronization query order was %#v", queries)
	}
	if !strings.Contains(logOutput.String(), sentinel.Error()) {
		t.Fatalf("Bitbucket synchronization did not warn on repository query error:\n%s", logOutput.String())
	}
}

func TestRecordedBitbucketQueriesRejectEmptyPopulation(t *testing.T) {
	recording := &recordingBitbucketQueryHTTP{}

	if _, err := recording.assertedQueries(); err == nil {
		t.Fatal("empty recorded Bitbucket query population passed")
	}
}

func assertRecordedBitbucketQueries(t *testing.T, recording *recordingBitbucketQueryHTTP, want []recordedBitbucketQuery) {
	t.Helper()
	queries, err := recording.assertedQueries()
	if err != nil {
		t.Fatal(err)
	}
	gotRequests := make([]httpclient.Request, len(queries))
	wantRequests := make([]httpclient.Request, len(want))
	for index, query := range queries {
		gotRequests[index] = query.Request
		if index >= len(want) {
			continue
		}
		wantRequests[index] = want[index].Request
		if reflect.TypeOf(query.Destination) != reflect.TypeOf(want[index].Destination) {
			t.Fatalf("Bitbucket query destination %d was %T, want %T", index, query.Destination, want[index].Destination)
		}
	}
	if !reflect.DeepEqual(gotRequests, wantRequests) {
		t.Fatalf("recorded Bitbucket queries differ:\n got: %#v\nwant: %#v", gotRequests, wantRequests)
	}
}
