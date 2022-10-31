package routes

import (
	"devengers-slidelist/googleDrive"
	_ "embed"
	"github.com/gorilla/mux"
	"google.golang.org/api/drive/v3"
	"html/template"
	"net/http"
	"os"
)

//go:embed templates/oauth.html
var authTpl string

//go:embed templates/index.html
var indexTpl string

//go:embed templates/assets/css/index.css
var indexCss string

//go:embed templates/assets/js/index.js
var indexMainJs string

//go:embed templates/assets/js/oauth.js
var oauthMainJs string

func Routes() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		scopes := []string{
			drive.DriveScope, drive.DriveAppdataScope,
			drive.DriveFileScope, drive.DriveMetadataScope,
			drive.DriveMetadataReadonlyScope, drive.DrivePhotosReadonlyScope,
			drive.DriveReadonlyScope,
		}

		cred := googleDrive.GetCredentials()
		oauthUrl := googleDrive.GenerateOAuthURL(cred, scopes...)

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
			Referer:      "http://" + os.Getenv("CLIENT_HOSTNAME") + ":" + os.Getenv("CLIENT_PORT"),
		})

		if err != nil {
			println(err.Error())
		}
	})

	r.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		tpl, _ := template.New("auth").Parse(authTpl)
		err := tpl.Execute(w, struct {
			Token string
		}{
			Token: r.URL.Query().Get("code"),
		})

		if err != nil {
			println(err.Error())
		}
	})

	r.HandleFunc("/assets/{type:[a-zA-Z0-9-_]+}/{file:[a-zA-Z0-9-_.]+}", func(w http.ResponseWriter, r *http.Request) {
		fileName := mux.Vars(r)["type"] + "/" + mux.Vars(r)["file"]
		fileContent := func() string {
			switch fileName {
			case "js/index.js":
				w.Header().Set("Content-Type", "text/javascript")
				return indexMainJs
			case "js/oauth.js":
				w.Header().Set("Content-Type", "text/javascript")
				return oauthMainJs
			case "css/index.css":
				w.Header().Set("Content-Type", "text/css")
				return indexCss
			default:
				w.Header().Set("Content-Type", "text/plain")
				return ""
			}
		}()

		if fileContent == "" {
			w.WriteHeader(http.StatusNotFound)
		}

		_, err := w.Write([]byte(fileContent))
		if err != nil {
			panic(err)
		}
	})

	http.Handle("/", r)
}
