package routes

import (
	"html/template"
	"net/http"
	"os"
)

type ManifestTemplateParams struct {
	StartUrl string
}

func getManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	manifestTpl, _ := ManifestTpl()
	tpl, _ := template.New("manifest").Parse(manifestTpl)

	startUrl := os.Getenv("CLIENT_HOSTNAME")
	if os.Getenv("CLIENT_PORT") != "80" && os.Getenv("CLIENT_PORT") != "443" {
		startUrl = startUrl + ":" + os.Getenv("CLIENT_PORT")
	}

	err := tpl.Execute(w, ManifestTemplateParams{StartUrl: startUrl})
	if err != nil {
		panic(err)
	}
}
