package routes

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func getAsset(w http.ResponseWriter, r *http.Request) {
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
