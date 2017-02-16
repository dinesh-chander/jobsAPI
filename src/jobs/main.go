package main

import (
	"api"
	"config"
	"log"
	"logger"
	_ "models"
	"models/job"
	"scrapers"
	"strconv"
	"strings"
	"time"
	jobTypes "types/jobs"
)

var loggerInstance *log.Logger

func main() {

	jobsStream := make(chan *jobTypes.Job, 500)

	jobManagersCount, err := strconv.ParseInt(config.GetConfig("jobManagersCount"), 10, 64)

	if err != nil {
		loggerInstance.Panicln(err.Error())
	}

	if jobManagersCount == 0 {
		jobManagersCount = 1
	}

	for jobManagersCount > 0 {
		go updateNewJobs(jobsStream)
		jobManagersCount = jobManagersCount - 1
	}

	scheduleScrappers(jobsStream)

	api.StartServer()
}

func updateNewJobs(jobsStream chan *jobTypes.Job) {
	for {
		select {
		case newJob := <-jobsStream:

			if newJob != nil {
				loggerInstance.Println("new job added")
				job.AddJob(newJob)
			}

		}
	}
}

func scheduleScrappers(jobsStream chan *jobTypes.Job) {
	fetchFrom, err := strconv.ParseInt(config.GetConfig("fetchFrom"), 10, 64)

	if err != nil {
		loggerInstance.Panicln(err.Error())
	}

	fetchFrom = (time.Now().Unix() * 1000) - (fetchFrom * 24 * 3600000)

	searchWordsList := strings.Split(config.GetConfig("searchWords"), ",")

	scrapers.InitiallizeScrappers(jobsStream, fetchFrom, searchWordsList)
}

func init() {
	loggerInstance = logger.Logger
}
