package routes

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	//go:embed templates/oauth.html
	authTpl embed.FS
	//go:embed templates/index.html
	indexTpl embed.FS
	//go:embed templates/manifest.json
	manifestTpl embed.FS
)

func AuthTpl() (string, error) {
	c, err := authTpl.ReadFile("templates/oauth.html")

	return string(c), err
}

func IndexTpl() (string, error) {
	c, err := indexTpl.ReadFile("templates/index.html")

	return string(c), err
}

func ManifestTpl() (string, error) {
	c, err := manifestTpl.ReadFile("templates/manifest.json")

	return string(c), err
}

var MimeTypes = map[string]string{
	"css": "text/css",
	"js":  "text/javascript",
	"png": "image/png",
}

func Routes() {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage)
	r.HandleFunc("/oauth", oAuthRedirectionPage)
	r.HandleFunc("/manifest.json", getManifest)
	r.HandleFunc("/service-worker.js", getServiceWorker)
	r.HandleFunc("/assets/{type:[a-zA-Z0-9-_]+}/{file:[a-zA-Z0-9-_.]+}", getAsset)

	http.Handle("/", r)
}
