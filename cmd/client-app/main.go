package main

import (
	"fmt"
	"log"
	"math/rand"
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
	serviceName string

	// backend endpoint names that client will send requests to
	backendEndpoints []string

	// rate that client sends requests to backend endpoints
	requestRate float64

	// name of the pod
	podName string

	// lifetime of the pod
	lifetime string
)

const (
	defaultListenPort       = 8080
	defaultBackendPort      = 80
	defaultServiceNameStatic = "Unknown server"
	defaultRequestRate      = 1

	listenPortEnvVarName       = "LISTEN_PORT"
	backendPortEnvVarName      = "BACKEND_PORT"
	serviceNameEnvVarName      = "SERVICE_NAME"
	backendEndpointsEnvVarName = "BACKEND_ENDPOINTS"
	requestRateEnvVarName      = "REQUEST_RATE"
	podNameEnvVarName		   = "POD_NAME"
	lifetimeEnvVarName		   = "LIFETIME"
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
	serviceName = os.Getenv(serviceNameEnvVarName)
	if serviceName == "" {
		log.Printf("%s env var is not set. Server name is set to %s\n", serviceNameEnvVarName, defaultServerName())
		serviceName = defaultServerName()
	} else {
		log.Printf("Server name: %s\n", serviceName)
	}

	// Backend endpoints
	backendEndpointsStr := os.Getenv(backendEndpointsEnvVarName)
	backendEndpointsTokens := strings.Split(backendEndpointsStr, ",")
	for _, token := range backendEndpointsTokens {
		token = strings.Trim(token, " ")
		if token != "" {
			backendEndpoints = append(backendEndpoints, token)
		}
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

	// Pod name
	podName = os.Getenv(podNameEnvVarName)
	if podName == "" {
		log.Printf("%s env var is not set. Pod name is empty\n", podNameEnvVarName)
	} else {
		log.Printf("Pod name: %s\n", podName)
	}

	// Lifetime
	lifetime = os.Getenv(lifetimeEnvVarName)
	if lifetime == "" {
		log.Printf("%s env var is not set. Lifetime is empty\n", lifetimeEnvVarName)
	} else {
		log.Printf("Lifetime: %s\n", lifetime)
	}
}

func defaultServerName() string {
	return defaultServiceNameStatic
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
				log.Println(err)
			} else {
				log.Printf("Sent to %s, received [%d]\n", url, resp.StatusCode)
			}
		}()

		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

func measureLife() {
	if lifetime == "" {
		return
	}

	duration, err := time.ParseDuration(lifetime)
	if err != nil {
		log.Printf("Lifetime value is invalid and ignored. %s.\n", lifetime)
		return
	}

	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(30)
	totalTime := time.Duration(int(duration.Seconds()) + randNum) * time.Second
	log.Printf("Pod will live for a total of %v seconds.\n", totalTime.Seconds())

	time.Sleep(totalTime)
	log.Printf("Pod time is up. Total time is %v seconds. Exit 0.\n", totalTime.Seconds())
	os.Exit(0)
}

func main() {
	bootstrap()

	go measureLife()

	go sendRequests()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s, from %v\n", serviceName, r.RequestURI, r.RemoteAddr)
		fmt.Fprintf(w, "%s - %s\n", serviceName, podName)
	})

	log.Printf("[%s] server starting on port %d\n", serviceName, listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), r)
}
