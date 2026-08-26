package webservice

import (
	"context"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	serveradapter "github.com/devdimensionlab/plybuild/internal/adapter/server"
	"github.com/devdimensionlab/plybuild/pkg/webservice/api"
	"log"
	"net/http"
	"runtime"
	"time"
)

const port = 7999

var server = &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port)}

type startWebServerDependencies struct {
	ServerOperations serveradapter.Dependencies
	Server           serveradapter.Selector
	HandleFunc       func(string, func(http.ResponseWriter, *http.Request))
	Print            func(...interface{})
}

func systemStartWebServerDependencies() startWebServerDependencies {
	return startWebServerDependencies{
		ServerOperations: serveradapter.System(),
		Server:           selectWebServer,
		HandleFunc:       http.HandleFunc,
		Print:            log.Print,
	}
}

func StartWebServer() {
	startWebServer(systemStartWebServerDependencies())
}

func startWebServer(dependencies startWebServerDependencies) {
	dependencies.handleFunc("/ui/generate", api.GetGenerate)
	dependencies.handleFunc("/api/generate", api.PostGenerate)

	dependencies.handleFunc("/ui/upgrade", api.GetUpgrade)
	dependencies.handleFunc("/api/upgrade", api.PostUpgrade)

	if err := serveradapter.ListenAndServe(dependencies.ServerOperations, dependencies.Server); err != nil {
		dependencies.print(err)
	}
}

func selectWebServer() *http.Server {
	return server
}

func (dependencies startWebServerDependencies) handleFunc(
	path string,
	handler func(http.ResponseWriter, *http.Request),
) {
	if dependencies.HandleFunc != nil {
		dependencies.HandleFunc(path, handler)
	}
}

func (dependencies startWebServerDependencies) print(arguments ...interface{}) {
	if dependencies.Print != nil {
		dependencies.Print(arguments...)
	}
}

type stopWebServerDependencies struct {
	ServerOperations serveradapter.Dependencies
	Server           serveradapter.Selector
	Background       func() context.Context
	WithTimeout      func(context.Context, time.Duration) (context.Context, context.CancelFunc)
}

func systemStopWebServerDependencies() stopWebServerDependencies {
	return stopWebServerDependencies{
		ServerOperations: serveradapter.System(),
		Server:           selectWebServer,
		Background:       context.Background,
		WithTimeout:      context.WithTimeout,
	}
}

func StopWebServer() {
	stopWebServer(systemStopWebServerDependencies())
}

func stopWebServer(dependencies stopWebServerDependencies) {
	if dependencies.Background == nil || dependencies.WithTimeout == nil {
		return
	}
	ctx, cancel := dependencies.WithTimeout(dependencies.Background(), 5*time.Second)
	defer cancel()
	_ = serveradapter.Shutdown(dependencies.ServerOperations, dependencies.Server, ctx)
}

type browserLauncherDependencies struct {
	Process process.Dependencies
	GOOS    string
}

func systemBrowserLauncherDependencies() browserLauncherDependencies {
	return browserLauncherDependencies{
		Process: process.SystemRunner(),
		GOOS:    runtime.GOOS,
	}
}

func OpenBrowser(url string) error {
	return openBrowser(systemBrowserLauncherDependencies(), url)
}

func openBrowser(dependencies browserLauncherDependencies, url string) error {
	switch dependencies.GOOS {
	case "linux":
		return process.Execute(dependencies.Process, process.Command{
			Name:  "xdg-open",
			Args:  []string{url},
			Start: true,
		})
	case "windows":
		return process.Execute(dependencies.Process, process.Command{
			Name:  "rundll32",
			Args:  []string{"url.dll,FileProtocolHandler", url},
			Start: true,
		})
	case "darwin":
		return process.Execute(dependencies.Process, process.Command{
			Name:  "open",
			Args:  []string{url},
			Start: true,
		})
	default:
		return fmt.Errorf("unsupported platform")
	}
}
