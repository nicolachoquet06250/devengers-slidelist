package routes

import (
	"net/http"
	"os"
)

func getServiceWorker(w http.ResponseWriter, r *http.Request) {
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
