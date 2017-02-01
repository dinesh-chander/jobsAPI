package main

import (
	"api"
	"config"
	"log"
	"logger"
	"models/job"
	"net/http"
	_ "net/http/pprof"
	"scrapers"
	"strconv"
	"time"
)

var loggerInstance *log.Logger

func main() {

	go http.ListenAndServe(":8081", nil)

	jobsStream := make(chan *job.Job, 500)

	go updateNewJobs(jobsStream)

	scheduleScrappers(jobsStream)

	api.StartServer()
}

func updateNewJobs(jobsStream chan *job.Job) {
	for {
		select {
		case newJob := <-jobsStream:
			job.AddJob(newJob)
		}
	}
}

func scheduleScrappers(jobsStream chan *job.Job) {
	fetchFrom, err := strconv.ParseInt(config.GetConfig("fetchFrom"), 10, 64)

	if err != nil {
		loggerInstance.Panicln(err.Error())
	}

	fetchFrom = (time.Now().Unix() * 1000) - (fetchFrom * 24 * 3600000)

	go scrapers.GetWhoIsHiringJobs(jobsStream, config.GetConfig("whoishiring"), fetchFrom)
}

func init() {
	loggerInstance = logger.Logger
}
