package webservice

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/webservice/api"
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
	_ = OpenBrowser(endpoint.Uri())
	<-blockingChannel
	//go StopWebServer()
}

func InitAndBlockProject(endpoint Endpoint, blockingChannel chan bool) func(project config.Project) error {
	return func(project config.Project) error {
		api.CurrentProject = project
		InitAndBlockStandalone(endpoint, blockingChannel)
		return nil
	}
}
