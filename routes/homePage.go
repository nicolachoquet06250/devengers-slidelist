package routes

import (
	"devengers-slidelist/googleDrive"
	"html/template"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
)

func homePage(w http.ResponseWriter, r *http.Request) {
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
	err := tpl.Execute(w, HomePageTemplateParams{
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
