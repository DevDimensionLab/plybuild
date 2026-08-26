package webservice

import (
	"context"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/webservice/api"
	"log"
	"net/http"
	"runtime"
	"time"
)

const port = 7999

var server = &http.Server{Addr: fmt.Sprintf(":%d", port)}

func StartWebServer() {
	http.HandleFunc("/ui/generate", api.GetGenerate)
	http.HandleFunc("/api/generate", api.PostGenerate)

	http.HandleFunc("/ui/upgrade", api.GetUpgrade)
	http.HandleFunc("/api/upgrade", api.PostUpgrade)

	if err := server.ListenAndServe(); err != nil {
		log.Print(err)
	}
}

func StopWebServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
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
