// Package httpclient isolates HTTP request execution from callers.
package httpclient

import (
	"errors"
	"net/http"
)

// ErrNoClient reports a zero-value dependency without performing a request.
var ErrNoClient = errors.New("HTTP client dependency is not configured")

// BasicAuth contains the ordered credentials for one request.
type BasicAuth struct {
	Username string
	Password string
}

// BearerJSON contains the access token for a JSON bearer request.
type BearerJSON struct {
	AccessToken string
}

// Request is the complete HTTP request value passed to a dependency.
type Request struct {
	URL        string
	BasicAuth  *BasicAuth
	BearerJSON *BearerJSON
}

// Client performs one complete HTTP request.
type Client interface {
	Do(Request) (*http.Response, error)
}

// Dependencies contains the HTTP effect used by a caller. Its zero value
// returns ErrNoClient without performing a request; production callers must
// select System explicitly.
type Dependencies struct {
	Client Client
}

// Execute passes the complete request to the configured dependency.
func Execute(dependencies Dependencies, request Request) (*http.Response, error) {
	if dependencies.Client == nil {
		return nil, ErrNoClient
	}
	return dependencies.Client.Do(request)
}

// System returns the production dependency that executes an HTTP request.
func System() Dependencies {
	return Dependencies{Client: systemClient{}}
}

type systemClient struct{}

func (systemClient) Do(request Request) (*http.Response, error) {
	if request.BasicAuth == nil && request.BearerJSON == nil {
		return http.Get(request.URL)
	}

	httpRequest, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return nil, err
	}
	if request.BasicAuth != nil {
		httpRequest.SetBasicAuth(request.BasicAuth.Username, request.BasicAuth.Password)
	}
	if request.BearerJSON != nil {
		httpRequest.Header.Add("Authorization", "Bearer "+request.BearerJSON.AccessToken)
		httpRequest.Header.Add("Content-Type", "application/json")
	}
	client := &http.Client{}
	return client.Do(httpRequest)
}
