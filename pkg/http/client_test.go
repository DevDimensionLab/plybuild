package http

import (
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

type errorReader struct {
	err error
}

func (reader errorReader) Read([]byte) (int, error) {
	return 0, reader.err
}
