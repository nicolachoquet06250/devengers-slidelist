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

func ListenAndServe(addr string, handler http.Handler) (*http.Server, error) {
	server := &http.Server{Addr: addr, Handler: handler}
	return server, server.ListenAndServe()
}

func StartServer() *http.Server {
	routes.Routes()

	var env = GetEnv()

	server, err := ListenAndServe(env.IP+":"+env.PORT, nil)

	if err != nil {
		if strings.Contains(err.Error(), "invalid port") {
			println(err.Error())
		} else {
			log.Fatal(err)
		}
	}

	return server
}

func RestartServer(server *http.Server) {
	server.Close()
	server.ListenAndServe()
}

func main() {
	/*server :*/ _ = StartServer()

	/*signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP)
	for sig := range signals {
		if sig == syscall.SIGHUP {
			RestartServer(server)
		}
	}*/
}
