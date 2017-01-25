package api

import (
	"config"
	"log"
	"logger"

	"net/http"
	"time"
)

var loggerInstance *log.Logger

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

	listenErr := server.ListenAndServe()

	if listenErr != nil {
		loggerInstance.Fatalln(listenErr.Error())
	}
}

func init() {
	loggerInstance = logger.Logger
}
