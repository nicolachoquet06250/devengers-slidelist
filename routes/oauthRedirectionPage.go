package routes

import (
	"html/template"
	"net/http"
)

type OAuthRedirectionPageTemplateParams struct {
	Token string
}

func oAuthRedirectionPage(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.New("auth").Parse(authTpl)
	err := tpl.Execute(w, OAuthRedirectionPageTemplateParams{
		Token: r.URL.Query().Get("code"),
	})

	if err != nil {
		println(err.Error())
	}
}
