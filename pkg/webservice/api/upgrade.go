package api

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/webservice/templates"
	"html/template"
	"net/http"
)

func GetUpgrade(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("upgradeTemplate").Parse(templates.Upgrade)
	_ = t.Execute(w, GOptions)
}

func PostUpgrade(w http.ResponseWriter, r *http.Request) {
	go func() { CallbackChannel <- true }()
	_, _ = fmt.Fprintf(w, "OK")
}
