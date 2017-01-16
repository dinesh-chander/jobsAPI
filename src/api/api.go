package api

import (
	"config"
	"fmt"
	"net/http"
)

func getAllJobs(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)

		responseData := []byte("Send Actual Data")
		response.Write(responseData)
	} else {
		errMsg := []byte("Wrong Method for the endpoint")
		response.Write(errMsg)
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

	fmt.Println("Starting API Server At : ", address)
	http.HandleFunc("/get", getAllJobs)
	http.ListenAndServe(address, nil)
}
