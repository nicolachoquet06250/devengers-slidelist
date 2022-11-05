package main

import "embed"

var (
	//go:embed routes/templates/oauth.html
	authTpl embed.FS
	//go:embed routes/templates/index.html
	indexTpl embed.FS
	//go:embed routes/templates/manifest.json
	manifestTpl embed.FS
)

func AuthTpl() (string, error) {
	c, err := authTpl.ReadFile("routes/templates/oauth.html")

	return string(c), err
}

func IndexTpl() (string, error) {
	c, err := indexTpl.ReadFile("routes/templates/index.html")

	return string(c), err
}

func ManifestTpl() (string, error) {
	c, err := manifestTpl.ReadFile("routes/templates/manifest.json")

	return string(c), err
}
