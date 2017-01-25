package api

import (
	"config"
	"log"
	"logger"

	"net/http"
	"time"
)

var loggerInstance *log.Logger

func getAllJobs(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)

		query := &Query{}
		query.parseQueryParamsFromURL(request.URL)

		responseData := []byte("Send Actual Data")
		response.Write(responseData)
	} else {
		errMsg := []byte("Wrong Method for the endpoint")
		response.Write(errMsg)
	}
}

func pathHandler(response http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/jobs":
		getAllJobs(response, request)
		break
	default:
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "text/html")
		response.Write([]byte("Wrong Endpoint"))
	}
}

func StartServer() {
	programConfig := config.GetConfig()

	port := "8080"
	networkInterface := "localhost"

	if programConfig["port"] != "" {
		port = programConfig["port"]
	}

	if programConfig["interface"] != "" {
		networkInterface = programConfig["interface"]
	}

	address := networkInterface + ":" + port

	loggerInstance.Println("Starting API Server At :", address)

	server := &http.Server{
		Addr:           address,
		Handler:        http.HandlerFunc(pathHandler),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}

func init() {
	loggerInstance = logger.Logger
}
