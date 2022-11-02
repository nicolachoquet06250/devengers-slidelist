package main

import (
	"devengers-slidelist/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type RequiredEnv struct {
	IP   string
	PORT string
}

func GetEnv() RequiredEnv {
	IP := os.Getenv("IP")
	PORT := os.Getenv("PORT")

	if strings.Contains(IP, ":") {
		IP = fmt.Sprintf(`[%s]`, IP)
	}

	if IP == "" {
		IP = "localhost"
	}

	if PORT == "" {
		PORT = "80"
	}

	return RequiredEnv{IP, PORT}
}

func main() {
	routes.Routes()

	var env = GetEnv()

	err := http.ListenAndServe(env.IP+":"+env.PORT, nil)
	if err != nil {
		if strings.Contains(err.Error(), "invalid port") {
			println(err.Error())
		} else {
			log.Fatal(err)
		}
	}
}
