package webservice

import (
	"context"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/webservice/api"
	"log"
	"net/http"
	"os/exec"
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

	log.Fatal(server.ListenAndServe())
}

func StopWebServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
