package webservice

import (
	"context"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/spring"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"text/template"
	"time"
)

type GenerateOptions struct {
	ProjectConfig *config.ProjectConfiguration
	CloudConfig   config.CloudConfig
	IoResponse    spring.IoRootResponse
}

var CallbackChannel = make(chan bool)
var GOptions GenerateOptions
var server = &http.Server{Addr: ":7999"}

func getHandlerGenerate(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("generateTemplate").Parse(generateTemplate)
	_ = t.Execute(w, GOptions)
}

func postHandlerGenerate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	cfg := GOptions.ProjectConfig
	cfg.GroupId = r.PostFormValue("groupId")
	cfg.ArtifactId = r.PostFormValue("artifactId")
	cfg.Package = r.PostFormValue("package")
	cfg.Name = r.PostFormValue("name")
	cfg.Description = r.PostFormValue("description")
	cfg.Language = r.PostFormValue("language")

	for key, values := range r.PostForm {
		if key == "templates" {
			for _, tmpl := range values {
				cfg.Templates = append(cfg.Templates, tmpl)
			}
		}
		if key == "dependencies" {
			for _, dep := range values {
				cfg.Dependencies = append(cfg.Dependencies, dep)
			}
		}
	}
	go func() { CallbackChannel <- true }()
	_, _ = fmt.Fprintf(w, "OK")
}

func StartService() {
	http.HandleFunc("/ui/generate", getHandlerGenerate)
	http.HandleFunc("/api/generate", postHandlerGenerate)
	log.Fatal(server.ListenAndServe())
}

func StopService() {
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
