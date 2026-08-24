package httpclient

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type recordingClient struct {
	request  Request
	response *http.Response
	err      error
}

func (client *recordingClient) Do(request Request) (*http.Response, error) {
	client.request = request
	return client.response, client.err
}

func TestDependenciesDefaultToSafeNoRequest(t *testing.T) {
	response, err := Execute(Dependencies{}, Request{URL: "https://must-not-request.example.invalid"})

	if response != nil {
		t.Fatalf("zero-value HTTP dependencies returned response %#v", response)
	}
	if !errors.Is(err, ErrNoClient) {
		t.Fatalf("zero-value HTTP dependencies returned %v, want %v", err, ErrNoClient)
	}
}

func TestExecutePassesCompleteRequestToDependency(t *testing.T) {
	wantRequest := Request{
		URL: "https://complete.example.invalid/repository/path?complete=value",
		BasicAuth: &BasicAuth{
			Username: "complete-user-before-password",
			Password: "complete-password-after-user",
		},
	}
	wantResponse := &http.Response{
		StatusCode: http.StatusAccepted,
		Status:     "202 Complete Response",
		Body:       io.NopCloser(strings.NewReader("complete body")),
	}
	sentinel := errors.New("complete dependency result")
	client := &recordingClient{response: wantResponse, err: sentinel}

	response, err := Execute(Dependencies{Client: client}, wantRequest)

	if response != wantResponse {
		t.Fatalf("Execute returned response %#v, want %#v", response, wantResponse)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute returned %v, want %v", err, sentinel)
	}
	if !reflect.DeepEqual(client.request, wantRequest) {
		t.Fatalf("dependency received an incomplete request:\n got: %#v\nwant: %#v", client.request, wantRequest)
	}
}

func TestSystemPreservesAnonymousAndBasicAuthGETRequests(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = previous
	})

	tests := []struct {
		name     string
		request  Request
		username string
		password string
	}{
		{
			name:    "anonymous",
			request: Request{URL: "https://anonymous.example.invalid/complete/path"},
		},
		{
			name: "basic-auth",
			request: Request{
				URL: "https://authenticated.example.invalid/complete/path",
				BasicAuth: &BasicAuth{
					Username: "system-user-before-password",
					Password: "system-password-after-user",
				},
			},
			username: "system-user-before-password",
			password: "system-password-after-user",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wantResponse := &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("system response")),
			}
			var recorded *http.Request
			http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
				recorded = request
				return wantResponse, nil
			})

			response, err := Execute(System(), test.request)

			if err != nil {
				t.Fatalf("system request returned an error: %v", err)
			}
			if response != wantResponse {
				t.Fatalf("system response was %#v, want %#v", response, wantResponse)
			}
			if recorded == nil {
				t.Fatal("system HTTP request population is empty")
			}
			if recorded.Method != http.MethodGet || recorded.URL.String() != test.request.URL {
				t.Fatalf("system request was %s %s", recorded.Method, recorded.URL.String())
			}
			username, password, authenticated := recorded.BasicAuth()
			if test.request.BasicAuth == nil {
				if authenticated {
					t.Fatalf("anonymous request unexpectedly used basic auth %q:%q", username, password)
				}
			} else if !authenticated || username != test.username || password != test.password {
				t.Fatalf("basic auth was (%q, %q, %t), want (%q, %q, true)", username, password, authenticated, test.username, test.password)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}
