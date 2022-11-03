package routes

import (
	"html/template"
	"net/http"
	"os"
)

type ServiceWorkerTemplateParams struct {
	Hostname string
}

func getServiceWorker(w http.ResponseWriter, r *http.Request) {
	path, _ := os.Getwd()
	content, err := os.ReadFile(path + "/routes/templates/assets/js/service-worker.js")
	if err != nil {
		panic(err)
	}

	tpl, _ := template.New("service-worker").Parse(string(content))

	hostname := os.Getenv("CLIENT_HOSTNAME")
	if os.Getenv("CLIENT_PORT") != "80" && os.Getenv("CLIENT_PORT") != "443" {
		hostname = hostname + ":" + os.Getenv("CLIENT_PORT")
	}

	w.Header().Set("Content-Type", "text/javascript")
	err = tpl.Execute(w, ServiceWorkerTemplateParams{Hostname: hostname})
	if err != nil {
		panic(err)
	}
}
