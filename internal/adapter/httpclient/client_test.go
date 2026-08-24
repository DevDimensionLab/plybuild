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

func TestExecutePassesCompleteBearerJSONRequestToDependency(t *testing.T) {
	wantRequest := Request{
		URL: "https://complete.example.invalid/rest/api/1.0/projects?limit=500",
		BearerJSON: &BearerJSON{
			AccessToken: "complete-bearer-token",
		},
	}
	wantResponse := &http.Response{
		StatusCode: http.StatusAccepted,
		Status:     "202 Complete Bearer Response",
		Body:       io.NopCloser(strings.NewReader("complete bearer body")),
	}
	sentinel := errors.New("complete bearer dependency result")
	client := &recordingClient{response: wantResponse, err: sentinel}

	response, err := Execute(Dependencies{Client: client}, wantRequest)

	if response != wantResponse {
		t.Fatalf("Execute returned response %#v, want %#v", response, wantResponse)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute returned %v, want %v", err, sentinel)
	}
	if !reflect.DeepEqual(client.request, wantRequest) {
		t.Fatalf("dependency received an incomplete bearer JSON request:\n got: %#v\nwant: %#v", client.request, wantRequest)
	}
}

func TestExecutePassesCompletePOSTRequestToDependency(t *testing.T) {
	wantRequest := Request{
		URL: "https://complete.example.invalid/internal/bsearch?compress=false",
		POST: &POST{
			Body: []byte(`{"size":500,"query":{"complete":true}}`),
			Header: http.Header{
				"Accept-Language": []string{"nb-NO"},
				"Authorization":   []string{"Bearer complete-token"},
				"Content-Type":    []string{"application/json"},
				"Kbn-Version":     []string{"8.9.0"},
			},
		},
	}
	wantResponse := &http.Response{
		StatusCode: http.StatusAccepted,
		Status:     "202 Complete POST Response",
		Body:       io.NopCloser(strings.NewReader("complete POST body")),
	}
	sentinel := errors.New("complete POST dependency result")
	client := &recordingClient{response: wantResponse, err: sentinel}

	response, err := Execute(Dependencies{Client: client}, wantRequest)

	if response != wantResponse {
		t.Fatalf("Execute returned response %#v, want %#v", response, wantResponse)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute returned %v, want %v", err, sentinel)
	}
	if !reflect.DeepEqual(client.request, wantRequest) {
		t.Fatalf("dependency received an incomplete POST request:\n got: %#v\nwant: %#v", client.request, wantRequest)
	}
}

func TestExecutePassesCompleteDefaultClientFormPOSTRequestToDependency(t *testing.T) {
	wantRequest := Request{
		URL: "https://complete.example.invalid/form-download?complete=value",
		POST: &POST{
			Body:             []byte("a-first=first%2Fvalue&z-last=last+value"),
			Header:           http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			UseDefaultClient: true,
		},
	}
	wantResponse := &http.Response{
		StatusCode: http.StatusAccepted,
		Status:     "202 Complete Form POST Response",
		Body:       io.NopCloser(strings.NewReader("complete form POST body")),
	}
	sentinel := errors.New("complete form POST dependency result")
	client := &recordingClient{response: wantResponse, err: sentinel}

	response, err := Execute(Dependencies{Client: client}, wantRequest)

	if response != wantResponse {
		t.Fatalf("Execute returned response %#v, want %#v", response, wantResponse)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute returned %v, want %v", err, sentinel)
	}
	if !reflect.DeepEqual(client.request, wantRequest) {
		t.Fatalf("dependency received an incomplete form POST request:\n got: %#v\nwant: %#v", client.request, wantRequest)
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

func TestSystemPreservesBearerJSONGETRequestHeaders(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	wantResponse := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("system bearer response")),
	}
	var recorded *http.Request
	http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		recorded = request
		return wantResponse, nil
	})
	request := Request{
		URL: "https://authenticated.example.invalid/complete/bearer.json?limit=500",
		BearerJSON: &BearerJSON{
			AccessToken: "system-complete-token",
		},
	}

	response, err := Execute(System(), request)

	if err != nil {
		t.Fatalf("system bearer request returned an error: %v", err)
	}
	if response != wantResponse {
		t.Fatalf("system bearer response was %#v, want %#v", response, wantResponse)
	}
	if recorded == nil {
		t.Fatal("system bearer HTTP request population is empty")
	}
	if recorded.Method != http.MethodGet || recorded.URL.String() != request.URL {
		t.Fatalf("system bearer request was %s %s", recorded.Method, recorded.URL.String())
	}
	if got := recorded.Header.Get("Authorization"); got != "Bearer system-complete-token" {
		t.Fatalf("system bearer Authorization header was %q", got)
	}
	if got := recorded.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("system bearer Content-Type header was %q", got)
	}
	if username, password, authenticated := recorded.BasicAuth(); authenticated {
		t.Fatalf("system bearer request unexpectedly used basic auth %q:%q", username, password)
	}
}

func TestSystemPreservesPOSTRequestBodyHeadersAndRedirects(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	firstBody := &trackingBody{reader: strings.NewReader("redirect")}
	finalBody := &trackingBody{reader: strings.NewReader("complete redirected POST response")}
	requests := []*http.Request{}
	requestBodies := []string{}
	http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request)
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		requestBodies = append(requestBodies, string(body))
		if len(requests) == 1 {
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Status:     "307 Temporary Redirect",
				Header:     http.Header{"Location": []string{"https://post.example.invalid/final"}},
				Body:       firstBody,
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       finalBody,
		}, nil
	})
	wantBody := `{"size":500,"query":{"redirected":true}}`
	request := Request{
		URL: "https://post.example.invalid/start",
		POST: &POST{
			Body: []byte(wantBody),
			Header: http.Header{
				"Accept-Language": []string{"nb-NO"},
				"Authorization":   []string{"Bearer redirect-token"},
				"Content-Type":    []string{"application/json"},
				"Kbn-Version":     []string{"8.9.0"},
			},
		},
	}

	response, err := Execute(System(), request)

	if err != nil {
		t.Fatalf("system POST request returned an error: %v", err)
	}
	if response == nil || response.Body != finalBody {
		t.Fatalf("system POST response was %#v", response)
	}
	if len(requests) != 2 {
		t.Fatalf("system POST request population was %d, want 2", len(requests))
	}
	wantURLs := []string{"https://post.example.invalid/start", "https://post.example.invalid/final"}
	for index, recorded := range requests {
		if recorded.Method != http.MethodPost || recorded.URL.String() != wantURLs[index] {
			t.Fatalf("system POST request %d was %s %s", index, recorded.Method, recorded.URL)
		}
		if requestBodies[index] != wantBody {
			t.Fatalf("system POST body %d was %q, want %q", index, requestBodies[index], wantBody)
		}
		for name, want := range map[string]string{
			"Accept-Language": "nb-NO",
			"Authorization":   "Bearer redirect-token",
			"Content-Type":    "application/json",
			"Kbn-Version":     "8.9.0",
		} {
			if got := recorded.Header.Get(name); got != want {
				t.Fatalf("system POST request %d header %s was %q, want %q", index, name, got, want)
			}
		}
	}
	if !firstBody.closed {
		t.Fatal("redirect response body was not closed")
	}
}

func TestSystemPreservesDefaultClientFormPOSTBodyHeaderAndRedirects(t *testing.T) {
	previousClient := http.DefaultClient
	previousTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultClient = previousClient
		http.DefaultTransport = previousTransport
	})
	redirects := 0
	http.DefaultClient = &http.Client{CheckRedirect: func(request *http.Request, via []*http.Request) error {
		redirects++
		return nil
	}}
	firstBody := &trackingBody{reader: strings.NewReader("redirect")}
	finalBody := &trackingBody{reader: strings.NewReader("complete redirected form response")}
	requests := []*http.Request{}
	requestBodies := []string{}
	http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request)
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		requestBodies = append(requestBodies, string(body))
		if len(requests) == 1 {
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Status:     "307 Temporary Redirect",
				Header:     http.Header{"Location": []string{"https://form.example.invalid/final"}},
				Body:       firstBody,
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Status:     "502 Must Still Be Returned",
			Body:       finalBody,
		}, nil
	})
	wantBody := "a-first=first%2Fvalue&z-last=last+value"
	request := Request{
		URL: "https://form.example.invalid/start",
		POST: &POST{
			Body:             []byte(wantBody),
			Header:           http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			UseDefaultClient: true,
		},
	}

	response, err := Execute(System(), request)

	if err != nil {
		t.Fatalf("system form POST request returned an error: %v", err)
	}
	if response == nil || response.Body != finalBody || response.StatusCode != http.StatusBadGateway {
		t.Fatalf("system form POST response was %#v", response)
	}
	if redirects != 1 {
		t.Fatalf("default client redirect callback ran %d times, want 1", redirects)
	}
	if len(requests) != 2 {
		t.Fatalf("system form POST request population was %d, want 2", len(requests))
	}
	wantURLs := []string{"https://form.example.invalid/start", "https://form.example.invalid/final"}
	for index, recorded := range requests {
		if recorded.Method != http.MethodPost || recorded.URL.String() != wantURLs[index] {
			t.Fatalf("system form POST request %d was %s %s", index, recorded.Method, recorded.URL)
		}
		if requestBodies[index] != wantBody {
			t.Fatalf("system form POST body %d was %q, want %q", index, requestBodies[index], wantBody)
		}
		if got := recorded.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("system form POST request %d Content-Type was %q", index, got)
		}
	}
	if !firstBody.closed {
		t.Fatal("form redirect response body was not closed")
	}
}

type trackingBody struct {
	reader io.Reader
	closed bool
}

func (body *trackingBody) Read(buffer []byte) (int, error) {
	return body.reader.Read(buffer)
}

func (body *trackingBody) Close() error {
	body.closed = true
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}
