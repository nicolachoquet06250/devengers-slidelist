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
	/*if len(os.Args) > 1 {
		var args []string
		for i, a := range os.Args {
			if i > 0 {
				args = append(args, a)
			}
		}

		if args[0] == "oauth" {
			googleDrive.Main()
		}

		os.Exit(0)
	}*/

	routes.Routes()

	IP := os.Getenv("IP")
	PORT := os.Getenv("PORT")

	if strings.Contains(IP, ":") {
		IP = fmt.Sprintf(`[%s]`, IP)
	}

	err := http.ListenAndServe(IP+":"+PORT, nil)
	if err != nil {
		if strings.Contains(err.Error(), "invalid port") == true {
			println(err.Error())
		} else {
			log.Fatal(err)
		}
	}
}
