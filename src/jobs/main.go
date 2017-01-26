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
	programConfig := config.GetConfig()
	jobsStream := make(chan *job.Job, 200)

	go updateNewJobs(jobsStream)

	scheduleScrappers(jobsStream, programConfig)

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

func scheduleScrappers(jobsStream chan *job.Job, programConfig map[string]string) {
	go scrapers.GetWhoIsHiringJobs(jobsStream, programConfig["whoishiring"])
}

func init() {
	loggerInstance = logger.Logger
}
