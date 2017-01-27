package main

import (
	"api"
	"config"
	"log"
	"logger"
	"models/job"
	"scrapers"
)

var loggerInstance *log.Logger

func main() {
	jobsStream := make(chan *job.Job, 200)

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
	go scrapers.GetWhoIsHiringJobs(jobsStream, config.GetConfig("whoishiring"))
}

func init() {
	loggerInstance = logger.Logger
}