package webservice

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/webservice/api"
)

type Endpoint string

const (
	Generate Endpoint = "generate"
	Upgrade           = "upgrade"
)

const host = "localhost"
const scheme = "http"

func (endpoint Endpoint) Uri() string {
	return fmt.Sprintf("%s://%s:%d/ui/%s", scheme, host, port, endpoint)
}

func InitAndBlockStandalone(endpoint Endpoint, blockingChannel chan bool) {
	go StartWebServer()
	OpenBrowser(endpoint.Uri())
	<-blockingChannel
	//go StopWebServer()
}

func InitAndBlockProject(endpoint Endpoint, blockingChannel chan bool) func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		api.CurrentProject = project
		InitAndBlockStandalone(endpoint, blockingChannel)
		return nil
	}
}
