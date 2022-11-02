package routes

import (
	"devengers-slidelist/googleDrive"
	_ "embed"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/api/drive/v3"
)

//go:embed templates/oauth.html
var authTpl string

//go:embed templates/index.html
var indexTpl string

var MimeTypes = map[string]string{
	"css": "text/css",
	"js":  "text/javascript",
	"png": "image/png",
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	scopes := []string{
		drive.DriveScope, drive.DriveAppdataScope,
		drive.DriveFileScope, drive.DriveMetadataScope,
		drive.DriveMetadataReadonlyScope, drive.DrivePhotosReadonlyScope,
	}

	cred := googleDrive.GetCredentials()
	oauthUrl := googleDrive.GenerateOAuthURL(cred, scopes...)

	referer := os.Getenv("CLIENT_HOSTNAME")
	if os.Getenv("CLIENT_PORT") != "80" && os.Getenv("CLIENT_PORT") != "443" {
		referer = referer + ":" + os.Getenv("CLIENT_PORT")
	}

	tpl, _ := template.New("index").Parse(indexTpl)
	err := tpl.Execute(w, struct {
		LinkUrl      string
		LinkLabel    string
		ApiKey       string
		ClientID     string
		ClientSecret string
		RedirectURI  string
		Referer      string
	}{
		LinkUrl:      oauthUrl,
		LinkLabel:    "Je me connect",
		ApiKey:       "AIzaSyAqH8eFIa42MngvvKTWNCXNY5jBPhDkTIs",
		ClientID:     cred.Web.ClientID,
		ClientSecret: cred.Web.ClientSecret,
		RedirectURI:  cred.Web.RedirectUris[0],
		Referer:      referer,
	})

	if err != nil {
		println(err.Error())
	}
}

func OAuthRedirectionPage(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.New("auth").Parse(authTpl)
	err := tpl.Execute(w, struct {
		Token string
	}{
		Token: r.URL.Query().Get("code"),
	})

	if err != nil {
		println(err.Error())
	}
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["type"] + "/" + mux.Vars(r)["file"]
	path, _ := os.Getwd()

	fileContent, err := os.ReadFile(path + "/routes/templates/assets/" + fileName)

	if err == nil {
		if string(fileContent) == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set("Content-Type", MimeTypes[mux.Vars(r)["type"]])
		}

		_, err = w.Write(fileContent)
	}

	if err != nil {
		panic(err)
	}
}

func GetManifest(w http.ResponseWriter, r *http.Request) {
	path, _ := os.Getwd()

	w.Header().Set("Content-Type", "application/json")

	content, err := os.ReadFile(path + "/routes/templates/manifest.json")
	if err != nil {
		panic(err)
	}

	tpl, _ := template.New("manifest").Parse(string(content))

	startUrl := os.Getenv("CLIENT_HOSTNAME")
	if os.Getenv("CLIENT_PORT") != "80" && os.Getenv("CLIENT_PORT") != "443" {
		startUrl = startUrl + ":" + os.Getenv("CLIENT_PORT")
	}

	err = tpl.Execute(w, struct{ StartUrl string }{StartUrl: startUrl})
	if err != nil {
		panic(err)
	}
}

func GetServiceWorker(w http.ResponseWriter, r *http.Request) {
	path, _ := os.Getwd()
	content, err := os.ReadFile(path + "/routes/templates/assets/js/service-worker.js")

	if err == nil {
		w.Header().Set("Content-Type", "text/javascript")
		_, err = w.Write(content)
	}

	if err != nil {
		panic(err)
	}
}

func Routes() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomePage)
	r.HandleFunc("/oauth", OAuthRedirectionPage)
	r.HandleFunc("/manifest.json", GetManifest)
	r.HandleFunc("/service-worker.js", GetServiceWorker)
	r.HandleFunc("/assets/{type:[a-zA-Z0-9-_]+}/{file:[a-zA-Z0-9-_.]+}", GetAsset)

	http.Handle("/", r)
}
