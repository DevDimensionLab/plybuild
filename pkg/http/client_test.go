package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
)

type recordedHTTPClient struct {
	requests []httpclient.Request
	response *http.Response
	err      error
}

func (client *recordedHTTPClient) dependencies() httpclient.Dependencies {
	return httpclient.Dependencies{Client: client}
}

func (client *recordedHTTPClient) Do(request httpclient.Request) (*http.Response, error) {
	if request.BasicAuth != nil {
		copied := *request.BasicAuth
		request.BasicAuth = &copied
	}
	client.requests = append(client.requests, request)
	return client.response, client.err
}

type trackingReadCloser struct {
	reader io.Reader
	closed bool
	reads  int
}

func (body *trackingReadCloser) Read(buffer []byte) (int, error) {
	body.reads++
	return body.reader.Read(buffer)
}

func (body *trackingReadCloser) Close() error {
	body.closed = true
	return nil
}

func TestXMLHelpersPreserveAnonymousAndBasicAuthResponseBehavior(t *testing.T) {
	tests := []struct {
		name    string
		request httpclient.Request
		invoke  func(httpclient.Dependencies, interface{}) error
	}{
		{
			name:    "anonymous",
			request: httpclient.Request{URL: "https://anonymous.example.invalid/metadata.xml"},
			invoke: func(dependencies httpclient.Dependencies, parsed interface{}) error {
				return getXml(dependencies, "https://anonymous.example.invalid/metadata.xml", parsed)
			},
		},
		{
			name: "basic-auth",
			request: httpclient.Request{
				URL: "https://authenticated.example.invalid/metadata.xml",
				BasicAuth: &httpclient.BasicAuth{
					Username: "response-user-before-password",
					Password: "response-password-after-user",
				},
			},
			invoke: func(dependencies httpclient.Dependencies, parsed interface{}) error {
				return getAuthXml(
					dependencies,
					"https://authenticated.example.invalid/metadata.xml",
					"response-user-before-password",
					"response-password-after-user",
					parsed,
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := &trackingReadCloser{reader: strings.NewReader(`<metadata><release>complete-response</release></metadata>`)}
			client := &recordedHTTPClient{response: &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       body,
			}}
			parsed := struct {
				Release string `xml:"release"`
			}{}

			err := test.invoke(client.dependencies(), &parsed)

			if err != nil {
				t.Fatalf("XML helper returned an error: %v", err)
			}
			if parsed.Release != "complete-response" {
				t.Fatalf("parsed XML release was %q", parsed.Release)
			}
			if !body.closed || body.reads == 0 {
				t.Fatalf("response body lifecycle was closed=%t reads=%d", body.closed, body.reads)
			}
			if !reflect.DeepEqual(client.requests, []httpclient.Request{test.request}) {
				t.Fatalf("HTTP dependency requests differ:\n got: %#v\nwant: %#v", client.requests, []httpclient.Request{test.request})
			}
		})
	}
}

func TestXMLHelpersPreserveStatusDependencyAndReadErrors(t *testing.T) {
	t.Run("status", func(t *testing.T) {
		body := &trackingReadCloser{reader: strings.NewReader("must not be read")}
		client := &recordedHTTPClient{response: &http.Response{
			StatusCode: http.StatusTeapot,
			Status:     "418 Complete Status",
			Body:       body,
		}}

		err := getXml(client.dependencies(), "https://status.example.invalid/metadata.xml", &struct{}{})

		want := "https://status.example.invalid/metadata.xml returned status code [418 Complete Status]"
		if err == nil || err.Error() != want {
			t.Fatalf("status error was %v, want %q", err, want)
		}
		if !body.closed || body.reads != 0 {
			t.Fatalf("status response body lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})

	t.Run("dependency", func(t *testing.T) {
		sentinel := errors.New("HTTP dependency failed")
		client := &recordedHTTPClient{err: sentinel}

		err := getAuthXml(client.dependencies(), "https://error.example.invalid/metadata.xml", "user", "password", &struct{}{})

		if !errors.Is(err, sentinel) {
			t.Fatalf("dependency error was %v, want %v", err, sentinel)
		}
	})

	t.Run("read", func(t *testing.T) {
		sentinel := errors.New("response read failed")
		body := &trackingReadCloser{reader: errorReader{err: sentinel}}
		client := &recordedHTTPClient{response: &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       body,
		}}

		err := getXml(client.dependencies(), "https://read.example.invalid/metadata.xml", &struct{}{})

		if !errors.Is(err, sentinel) {
			t.Fatalf("read error was %v, want %v", err, sentinel)
		}
		if !body.closed || body.reads == 0 {
			t.Fatalf("read-error response lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})
}

func TestJSONHelperPassesCompleteAnonymousRequestAndParsesResponse(t *testing.T) {
	body := &trackingReadCloser{reader: strings.NewReader(`{"name":"complete-json","nested":{"value":"complete-value"}}`)}
	client := &recordedHTTPClient{response: &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 Complete JSON",
		Body:       body,
	}}
	parsed := struct {
		Name   string `json:"name"`
		Nested struct {
			Value string `json:"value"`
		} `json:"nested"`
	}{}
	requestURL := "https://anonymous.example.invalid/complete/discovery?first=one&second=two"

	err := getJson(client.dependencies(), requestURL, &parsed)

	if err != nil {
		t.Fatalf("JSON helper returned an error: %v", err)
	}
	if parsed.Name != "complete-json" || parsed.Nested.Value != "complete-value" {
		t.Fatalf("parsed JSON response was incomplete: %#v", parsed)
	}
	if !body.closed || body.reads == 0 {
		t.Fatalf("JSON response body lifecycle was closed=%t reads=%d", body.closed, body.reads)
	}
	want := []httpclient.Request{{URL: requestURL}}
	if !reflect.DeepEqual(client.requests, want) {
		t.Fatalf("JSON dependency requests differ:\n got: %#v\nwant: %#v", client.requests, want)
	}
}

func TestJSONHelperPreservesStatusDependencyReadAndUnmarshalErrors(t *testing.T) {
	t.Run("status", func(t *testing.T) {
		body := &trackingReadCloser{reader: strings.NewReader("must not be read")}
		client := &recordedHTTPClient{response: &http.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Status:     "422 Complete Status",
			Body:       body,
		}}
		requestURL := "https://status.example.invalid/complete.json"

		err := getJson(client.dependencies(), requestURL, &struct{}{})

		want := requestURL + " returned status code [422 Complete Status]"
		if err == nil || err.Error() != want {
			t.Fatalf("JSON status error was %v, want %q", err, want)
		}
		if !body.closed || body.reads != 0 {
			t.Fatalf("JSON status response lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})

	t.Run("dependency after complete request", func(t *testing.T) {
		sentinel := errors.New("JSON HTTP dependency failed")
		client := &recordedHTTPClient{err: sentinel}
		requestURL := "https://dependency.example.invalid/complete.json"

		err := getJson(client.dependencies(), requestURL, &struct{}{})

		if !errors.Is(err, sentinel) {
			t.Fatalf("JSON dependency error was %v, want %v", err, sentinel)
		}
		want := []httpclient.Request{{URL: requestURL}}
		if !reflect.DeepEqual(client.requests, want) {
			t.Fatalf("JSON dependency error lost the complete request:\n got: %#v\nwant: %#v", client.requests, want)
		}
	})

	t.Run("read", func(t *testing.T) {
		sentinel := errors.New("JSON response read failed")
		body := &trackingReadCloser{reader: errorReader{err: sentinel}}
		client := &recordedHTTPClient{response: &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       body,
		}}

		err := getJson(client.dependencies(), "https://read.example.invalid/complete.json", &struct{}{})

		if !errors.Is(err, sentinel) {
			t.Fatalf("JSON read error was %v, want %v", err, sentinel)
		}
		if !body.closed || body.reads == 0 {
			t.Fatalf("JSON read-error response lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		body := &trackingReadCloser{reader: strings.NewReader(`{"incomplete":`)}
		client := &recordedHTTPClient{response: &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       body,
		}}

		err := getJson(client.dependencies(), "https://unmarshal.example.invalid/complete.json", &struct{}{})

		var syntaxError *json.SyntaxError
		if !errors.As(err, &syntaxError) {
			t.Fatalf("JSON unmarshal error was %T %v, want *json.SyntaxError", err, err)
		}
		if !body.closed || body.reads == 0 {
			t.Fatalf("JSON unmarshal response lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})
}

func TestJSONHelperDependenciesDefaultToSafeNoRequest(t *testing.T) {
	err := getJson(httpclient.Dependencies{}, "https://must-not-request.example.invalid/complete.json", &struct{}{})

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe JSON dependency default returned %v, want %v", err, httpclient.ErrNoClient)
	}
}

func TestGetJSONSystemPreservesAnonymousGETRedirects(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	firstBody := &trackingReadCloser{reader: strings.NewReader("redirect")}
	finalBody := &trackingReadCloser{reader: strings.NewReader(`{"redirected":"complete"}`)}
	requests := []*http.Request{}
	http.DefaultTransport = downloadRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request)
		if len(requests) == 1 {
			return &http.Response{
				StatusCode: http.StatusFound,
				Status:     "302 Found",
				Header:     http.Header{"Location": []string{"https://redirect.example.invalid/final.json"}},
				Body:       firstBody,
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       finalBody,
		}, nil
	})
	parsed := struct {
		Redirected string `json:"redirected"`
	}{}

	err := GetJson("https://redirect.example.invalid/start.json", &parsed)

	if err != nil {
		t.Fatalf("system JSON redirect returned an error: %v", err)
	}
	if parsed.Redirected != "complete" {
		t.Fatalf("redirected JSON value was %q", parsed.Redirected)
	}
	if len(requests) != 2 {
		t.Fatalf("redirect request population was %d, want 2", len(requests))
	}
	wantURLs := []string{
		"https://redirect.example.invalid/start.json",
		"https://redirect.example.invalid/final.json",
	}
	for index, request := range requests {
		if request.Method != http.MethodGet || request.URL.String() != wantURLs[index] {
			t.Fatalf("redirect request %d was %s %s, want GET %s", index, request.Method, request.URL, wantURLs[index])
		}
		if _, _, authenticated := request.BasicAuth(); authenticated {
			t.Fatalf("redirect request %d unexpectedly used basic auth", index)
		}
	}
	if !firstBody.closed || !finalBody.closed {
		t.Fatalf("redirect bodies were closed=(%t, %t)", firstBody.closed, finalBody.closed)
	}
}

type errorReader struct {
	err error
}

func (reader errorReader) Read([]byte) (int, error) {
	return 0, reader.err
}
