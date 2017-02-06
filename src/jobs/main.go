package main

import (
	"api"
	"config"
	"log"
	"logger"
	"models/job"
	"scrapers"
	"strconv"
	"strings"
	"time"
)

var loggerInstance *log.Logger

func main() {

	jobsStream := make(chan *job.Job, 500)

	go updateNewJobs(jobsStream)

	scheduleScrappers(jobsStream)

	api.StartServer()
}

func updateNewJobs(jobsStream chan *job.Job) {
	searchWordsList := strings.Split(config.GetConfig("validWords"), ",")

	for {
		select {
		case newJob := <-jobsStream:

			go func(newJob *job.Job) {
				for _, searchWord := range searchWordsList {
					if strings.Contains(strings.ToUpper(newJob.Title), strings.ToUpper(searchWord)) {
						job.AddJob(newJob)
						break
					}
				}
			}(newJob)
		}
	}
}

func scheduleScrappers(jobsStream chan *job.Job) {
	fetchFrom, err := strconv.ParseInt(config.GetConfig("fetchFrom"), 10, 64)

	if err != nil {
		loggerInstance.Panicln(err.Error())
	}

	fetchFrom = (time.Now().Unix() * 1000) - (fetchFrom * 24 * 3600000)

	searchWordsList := strings.Split(config.GetConfig("validWords"), ",")

	go scrapers.GetWhoIsHiringJobs(jobsStream, config.GetConfig("whoishiring"), fetchFrom, searchWordsList)
}

func init() {
	loggerInstance = logger.Logger
}
