package googleDrive

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

type OAuthCredentialsWeb struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
	JavascriptOrigins       []string `json:"javascript_origins"`
}

type OAuthCredentials struct {
	Web OAuthCredentialsWeb `json:"web"`
}

func GetCredentials() OAuthCredentials {
	path, _ := os.Getwd()
	b, err := os.ReadFile(path + "/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	var j OAuthCredentials

	if err := json.Unmarshal(b, &j); err != nil {
		log.Fatalf("Unable to decode json: %v", err)
	}

	return j
}

func GenerateOAuthURL(credentials OAuthCredentials, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = append(scopes, drive.DriveMetadataReadonlyScope)
	}
	return fmt.Sprintf("%s?access_type=offline&client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=state-token", credentials.Web.AuthURI, credentials.Web.ClientID, credentials.Web.RedirectUris[0], strings.Join(scopes, " "))
}
