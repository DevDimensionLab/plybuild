package spring

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

type recordedSpringDiscoveryRequest struct {
	Request     httpclient.Request
	Destination interface{}
}

type recordingSpringDiscoveryHTTP struct {
	requests  []recordedSpringDiscoveryRequest
	responses map[string]string
	err       error
}

func (recording *recordingSpringDiscoveryHTTP) dependencies() discoveryDependencies {
	return discoveryDependencies{Client: recording}
}

func (recording *recordingSpringDiscoveryHTTP) GetJSON(request httpclient.Request, parsed interface{}) error {
	if request.BasicAuth != nil {
		copied := *request.BasicAuth
		request.BasicAuth = &copied
	}
	recording.requests = append(recording.requests, recordedSpringDiscoveryRequest{
		Request:     request,
		Destination: parsed,
	})
	if recording.err != nil {
		return recording.err
	}
	return json.Unmarshal([]byte(recording.responses[request.URL]), parsed)
}

func (recording *recordingSpringDiscoveryHTTP) assertedRequests() ([]recordedSpringDiscoveryRequest, error) {
	if len(recording.requests) == 0 {
		return nil, errors.New("recorded Spring discovery request population is empty")
	}
	return recording.requests, nil
}

func TestSpringDiscoveryKeepsCompleteRootAndDependenciesURLsAndParsesResponses(t *testing.T) {
	previousBaseURL := baseUrl
	baseUrl = "https://spring.example.invalid/complete-base"
	t.Cleanup(func() { baseUrl = previousBaseURL })
	recording := &recordingSpringDiscoveryHTTP{responses: map[string]string{
		baseUrl: `{
			"artifactId":{"default":"complete-artifact","text":"complete-text"},
			"bootVersion":{"default":"3.4.5-complete","type":"select","values":[{"id":"3.4.5","name":"3.4.5 complete"}]}
		}`,
		baseUrl + "/dependencies": `{
			"bootVersion":"3.4.5-complete",
			"dependencies":{"complete-valid":{"groupId":"com.example.complete","artifactId":"complete-artifact","scope":"compile"}},
			"repositories":{"complete-repository":{"name":"Complete Repository","url":"https://repo.example.invalid/complete","snapshotEnabled":true}}
		}`,
	}}

	root, err := getRoot(recording.dependencies())
	if err != nil {
		t.Fatalf("Spring root discovery returned an error: %v", err)
	}
	dependencies, err := getDependencies(recording.dependencies())
	if err != nil {
		t.Fatalf("Spring dependency discovery returned an error: %v", err)
	}

	if root.ArtifactId.Default != "complete-artifact" || root.ArtifactId.Type != "complete-text" ||
		root.BootVersion.Default != "3.4.5-complete" || len(root.BootVersion.Values) != 1 ||
		root.BootVersion.Values[0].Name != "3.4.5 complete" {
		t.Fatalf("parsed Spring root response was incomplete: %#v", root)
	}
	dependency := dependencies.Dependencies["complete-valid"]
	repository := dependencies.Repositories["complete-repository"]
	if dependencies.BootVersion != "3.4.5-complete" || dependency.GroupId != "com.example.complete" ||
		dependency.ArtifactId != "complete-artifact" || dependency.Scope != "compile" ||
		repository.Name != "Complete Repository" || repository.URL != "https://repo.example.invalid/complete" ||
		!repository.SnapshotEnabled {
		t.Fatalf("parsed Spring dependencies response was incomplete: %#v", dependencies)
	}
	assertRecordedSpringDiscoveryRequests(t, recording, []recordedSpringDiscoveryRequest{
		{Request: httpclient.Request{URL: baseUrl}, Destination: (*IoRootResponse)(nil)},
		{Request: httpclient.Request{URL: baseUrl + "/dependencies"}, Destination: (*IoDependenciesResponse)(nil)},
	})
}

func TestSpringDiscoveryReturnsDependencyErrorsAfterCompleteRequestDelivery(t *testing.T) {
	previousBaseURL := baseUrl
	baseUrl = "https://errors.spring.example.invalid/complete-base"
	t.Cleanup(func() { baseUrl = previousBaseURL })
	sentinel := errors.New("complete Spring discovery dependency failed")
	tests := []struct {
		name   string
		url    string
		invoke func(discoveryDependencies) error
		target interface{}
	}{
		{
			name: "root",
			url:  baseUrl,
			invoke: func(dependencies discoveryDependencies) error {
				_, err := getRoot(dependencies)
				return err
			},
			target: (*IoRootResponse)(nil),
		},
		{
			name: "dependencies",
			url:  baseUrl + "/dependencies",
			invoke: func(dependencies discoveryDependencies) error {
				_, err := getDependencies(dependencies)
				return err
			},
			target: (*IoDependenciesResponse)(nil),
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestSpringDiscoveryReturnsDependencyErrorsAfterCompleteRequestDelivery test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingSpringDiscoveryHTTP{err: sentinel}

			err := test.invoke(recording.dependencies())

			if !errors.Is(err, sentinel) {
				t.Fatalf("Spring discovery dependency error was %v, want %v", err, sentinel)
			}
			requests, assertErr := recording.assertedRequests()
			if assertErr != nil {
				t.Fatal(assertErr)
			}
			if len(requests) != 1 || !reflect.DeepEqual(requests[0].Request, httpclient.Request{URL: test.url}) {
				t.Fatalf("Spring discovery error request was %#v, want %q", requests, test.url)
			}
			if reflect.TypeOf(requests[0].Destination) != reflect.TypeOf(test.target) {
				t.Fatalf("Spring discovery destination was %T, want %T", requests[0].Destination, test.target)
			}
		})
	}
}

func TestSpringDiscoveryDependenciesDefaultToSafeNoRequest(t *testing.T) {
	tests := []struct {
		name   string
		invoke func() error
	}{
		{
			name: "root",
			invoke: func() error {
				_, err := getRoot(discoveryDependencies{})
				return err
			},
		},
		{
			name: "dependencies",
			invoke: func() error {
				_, err := getDependencies(discoveryDependencies{})
				return err
			},
		},
		{
			name: "validation",
			invoke: func() error {
				return validate(discoveryDependencies{}, config.ProjectConfiguration{Dependencies: []string{"must-not-request"}})
			},
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestSpringDiscoveryDependenciesDefaultToSafeNoRequest test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.invoke(); !errors.Is(err, httpclient.ErrNoClient) {
				t.Fatalf("safe Spring discovery default returned %v, want %v", err, httpclient.ErrNoClient)
			}
		})
	}
}

func TestSpringValidationSelectsDependenciesRequestAndPreservesErrors(t *testing.T) {
	previousBaseURL := baseUrl
	baseUrl = "https://validation.spring.example.invalid/complete-base"
	t.Cleanup(func() { baseUrl = previousBaseURL })

	t.Run("empty dependencies select no request", func(t *testing.T) {
		recording := &recordingSpringDiscoveryHTTP{}

		err := validate(recording.dependencies(), config.ProjectConfiguration{})

		if err != nil {
			t.Fatalf("empty dependency validation returned an error: %v", err)
		}
		if len(recording.requests) != 0 {
			t.Fatalf("empty dependency validation made requests: %#v", recording.requests)
		}
	})

	t.Run("valid dependencies select one anonymous dependencies request", func(t *testing.T) {
		recording := &recordingSpringDiscoveryHTTP{responses: map[string]string{
			baseUrl + "/dependencies": `{"dependencies":{"first-valid":{},"second-valid":{}}}`,
		}}

		err := validate(recording.dependencies(), config.ProjectConfiguration{
			Dependencies: []string{"second-valid", "first-valid"},
		})

		if err != nil {
			t.Fatalf("valid Spring dependencies returned an error: %v", err)
		}
		assertRecordedSpringDiscoveryRequests(t, recording, []recordedSpringDiscoveryRequest{{
			Request: httpclient.Request{URL: baseUrl + "/dependencies"}, Destination: (*IoDependenciesResponse)(nil),
		}})
	})

	t.Run("invalid dependency reporting preserves user order", func(t *testing.T) {
		recording := &recordingSpringDiscoveryHTTP{responses: map[string]string{
			baseUrl + "/dependencies": `{"dependencies":{"only-valid":{}}}`,
		}}

		err := validate(recording.dependencies(), config.ProjectConfiguration{
			Dependencies: []string{"bad-first", "only-valid", "bad-last"},
		})

		want := "[bad-first bad-last] not found in valid list of dependencies [only-valid]"
		if err == nil || err.Error() != want {
			t.Fatalf("invalid dependency error was %v, want %q", err, want)
		}
		assertRecordedSpringDiscoveryRequests(t, recording, []recordedSpringDiscoveryRequest{{
			Request: httpclient.Request{URL: baseUrl + "/dependencies"}, Destination: (*IoDependenciesResponse)(nil),
		}})
	})

	t.Run("dependency error is returned after validation request", func(t *testing.T) {
		sentinel := errors.New("validation dependency failed")
		recording := &recordingSpringDiscoveryHTTP{err: sentinel}

		err := validate(recording.dependencies(), config.ProjectConfiguration{Dependencies: []string{"complete-request"}})

		if !errors.Is(err, sentinel) {
			t.Fatalf("validation dependency error was %v, want %v", err, sentinel)
		}
		assertRecordedSpringDiscoveryRequests(t, recording, []recordedSpringDiscoveryRequest{{
			Request: httpclient.Request{URL: baseUrl + "/dependencies"}, Destination: (*IoDependenciesResponse)(nil),
		}})
	})
}

func TestRecordedSpringDiscoveryRequestsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingSpringDiscoveryHTTP{}

	if _, err := recording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Spring discovery population passed")
	}
}

func assertRecordedSpringDiscoveryRequests(t *testing.T, recording *recordingSpringDiscoveryHTTP, want []recordedSpringDiscoveryRequest) {
	t.Helper()
	requests, err := recording.assertedRequests()
	if err != nil {
		t.Fatal(err)
	}
	got := make([]httpclient.Request, len(requests))
	wantRequests := make([]httpclient.Request, len(want))
	for index, request := range requests {
		got[index] = request.Request
		if index >= len(want) {
			continue
		}
		wantRequests[index] = want[index].Request
		if reflect.TypeOf(request.Destination) != reflect.TypeOf(want[index].Destination) {
			t.Fatalf("Spring discovery destination %d was %T, want %T", index, request.Destination, want[index].Destination)
		}
	}
	if !reflect.DeepEqual(got, wantRequests) {
		t.Fatalf("recorded Spring discovery requests differ:\n got: %#v\nwant: %#v", got, wantRequests)
	}
}
