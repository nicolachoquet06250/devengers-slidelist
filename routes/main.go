package routes

import (
	_ "embed"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	//go:embed templates/oauth.html
	authTpl string
	//go:embed templates/index.html
	indexTpl string
	//go:embed templates/manifest.json
	manifestTpl string
)

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
