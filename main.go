package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	// port number of web server, and peer client server ports
	port int

	// server name from SERVER_NAME env var
	serverName string
)

const (
	defaultPort             = 8080
	defaultServerNameStatic = "Unknown server"
	portEnvVarname          = "PORT"
	serverNameEnvVarname    = "SERVER_NAME"
)

func bootstrap() {
	var err error
	port, err = strconv.Atoi(os.Getenv(portEnvVarname))
	if err != nil {
		log.Printf("%s env var is not set. Port will be set to %d\n", portEnvVarname, defaultPort)
		port = defaultPort
	}

	serverName = os.Getenv(serverNameEnvVarname)
	if serverName == "" {
		log.Printf("%s env var is not set. Server name is set to %s\n", serverNameEnvVarname, defaultServerName())
		serverName = defaultServerName()
	}
}

func defaultServerName() string {
	return defaultServerNameStatic
}

func main() {
	bootstrap()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s, from %v\n", serverName, r.RequestURI, r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", serverName)
	})

	log.Printf("[%s] server starting on port %d\n", serverName, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
