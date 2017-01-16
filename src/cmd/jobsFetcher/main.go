package main

import (
	"jobsAPI-helpers/scrapers"
	"jobsAPI-helpers/job"
	"jobsAPI-helpers/api"
)

func main() {
	config := getConfig()
	jobsStream := make(chan job.Job, 10)

	scheduleScrappers(jobsStream, config)
	go updateNewJobs(jobsStream)
	api.StartServer()
}

func updateNewJobs(jobsStream chan job.Job) {
	for {
		select {
		case newJob := <-jobsStream:
			job.AddJob(newJob)
		}
	}
}

func scheduleScrappers(jobsStream chan job.Job, config map[string]string) {
	go scrapers.GetWhoIsHiringJobs(jobsStream, "")
}
