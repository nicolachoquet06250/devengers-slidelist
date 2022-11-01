package main

import (
	"devengers-slidelist/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	routes.Routes()

	IP := os.Getenv("IP")
	PORT := os.Getenv("PORT")

	if strings.Contains(IP, ":") {
		IP = fmt.Sprintf(`[%s]`, IP)
	}

	err := http.ListenAndServe(IP+":"+PORT, nil)
	if err != nil {
		if strings.Contains(err.Error(), "invalid port") {
			println(err.Error())
		} else {
			log.Fatal(err)
		}
	}
}
