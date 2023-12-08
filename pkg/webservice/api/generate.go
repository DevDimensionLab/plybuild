package api

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/plybuild/pkg/webservice/templates"
	"net/http"
	"text/template"
)

type GenerateOptions struct {
	ProjectConfig *config.ProjectConfiguration
	CloudConfig   config.CloudConfig
	IoResponse    spring.IoRootResponse
}

var GOptions GenerateOptions

func GetGenerate(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("generateTemplate").Parse(templates.Generate)
	_ = t.Execute(w, GOptions)
}

func PostGenerate(w http.ResponseWriter, r *http.Request) {
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
