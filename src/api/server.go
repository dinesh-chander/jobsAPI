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
	port := "8080"
	networkInterface := "localhost"

	if config.GetConfig("port") != "" {
		port = config.GetConfig("port")
	}

	if config.GetConfig("interface") != "" {
		networkInterface = config.GetConfig("interface")
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
