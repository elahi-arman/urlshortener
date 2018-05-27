package main

import (
	"fmt"
	"os"

	"github.com/elahi-arman/urlshortener/server"
)

func main() {

	var (
		shortlyHome string
		isSet       bool
	)

	var env = "SHORTLY_HOME"

	if shortlyHome, isSet = os.LookupEnv(env); !isSet {
		fmt.Fprintf(os.Stderr, "ERROR: Env Var [%s] is not set\n", env)
		os.Exit(1)
	}

	var appConfigFile = shortlyHome + "/config/config.yaml"
	server.StartServer(appConfigFile)
}
