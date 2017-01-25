package main

import (
	"api"
	"config"
	"models/job"
	"scrapers"
)

func main() {
	programConfig := config.GetConfig()
	jobsStream := make(chan *job.Job, 10)

	scheduleScrappers(jobsStream, programConfig)
	go updateNewJobs(jobsStream)
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
