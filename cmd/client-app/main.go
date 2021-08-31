package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	// listenPort number of web server, and peer client server ports
	listenPort int

	// backendPort number of web server, and peer client server ports
	backendPort int

	// server name from SERVER_NAME env var
	serverName string

	// backend endpoint names that client will send requests to
	backendEndpoints []string

	// rate that client sends requests to backend endpoints
	requestRate float64
)

const (
	defaultListenPort       = 8080
	defaultBackendPort      = 80
	defaultServerNameStatic = "Unknown server"
	defaultRequestRate      = 1

	listenPortEnvVarName       = "LISTEN_PORT"
	backendPortEnvVarName      = "BACKEND_PORT"
	serverNameEnvVarName       = "SERVER_NAME"
	backendEndpointsEnvVarName = "BACKEND_ENDPOINTS"
	requestRateEnvVarName      = "REQUEST_RATE"
)

func bootstrap() {
	var err error

	// Listening port
	listenPort, err = strconv.Atoi(os.Getenv(listenPortEnvVarName))
	if err != nil {
		log.Printf("%s env var is not set. Listening port will be set to %d\n", listenPortEnvVarName, defaultListenPort)
		listenPort = defaultListenPort
	} else {
		log.Printf("Listening port: %d\n", listenPort)
	}

	// Backend port
	backendPort, err = strconv.Atoi(os.Getenv(backendPortEnvVarName))
	if err != nil {
		log.Printf("%s env var is not set. Backend port will be set to %d\n", backendPortEnvVarName, defaultBackendPort)
		backendPort = defaultBackendPort
	} else {
		log.Printf("Backend port: %d\n", listenPort)
	}

	// Server name
	serverName = os.Getenv(serverNameEnvVarName)
	if serverName == "" {
		log.Printf("%s env var is not set. Server name is set to %s\n", serverNameEnvVarName, defaultServerName())
		serverName = defaultServerName()
	} else {
		log.Printf("Server name: %s\n", serverName)
	}

	// Backend endpoints
	backendEndpointsStr := os.Getenv(backendEndpointsEnvVarName)
	backendEndpoints = strings.Split(backendEndpointsStr, ",")
	for idx, item := range backendEndpoints {
		backendEndpoints[idx] = strings.Trim(item, " ")
	}
	log.Printf("Backend endpoints: %v\n", backendEndpoints)

	// Request rate
	requestRate, err = strconv.ParseFloat(os.Getenv(requestRateEnvVarName), 32)
	if err != nil {
		log.Printf("%s env var is not set. Request rate will be set to %d\n", requestRateEnvVarName, defaultRequestRate)
		requestRate = defaultRequestRate
	} else {
		log.Printf("Request rate: %f requests per second\n", requestRate)
	}
}

func defaultServerName() string {
	return defaultServerNameStatic
}

// sends requests at a roughly specified request rate
func sendRequests() {
	intervalMs := 1000 / requestRate

	lengEndpoints := len(backendEndpoints)

	var endpoint string
	for i := 0; ; i++ {
		if i == lengEndpoints {
			i = 0
		}

		// randomize endpoints
		endpoint = backendEndpoints[i]

		go func() {
			url := fmt.Sprintf("http://%s:%d/", endpoint, backendPort)
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			} else {
				log.Printf("Sent to %s, received [%d]\n", url, resp.StatusCode)
			}
		}()

		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

func main() {
	bootstrap()

	go sendRequests()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s, from %v\n", serverName, r.RequestURI, r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", serverName)
	})

	log.Printf("[%s] server starting on port %d\n", serverName, listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), r)
}
