package scrapers

import (
	"fmt"
	"job"
	"time"

	"github.com/gorhill/cronexpr"
)

type whoIsHiringJobStruct struct {
}

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string) {
	whoIsHiringJobsStream := make(chan *whoIsHiringJobStruct, 5)

	go fetchJobs(whoIsHiringJobsStream, scheduleAt)

	for {
		select {
		case newJob := <-whoIsHiringJobsStream:
			fmt.Println("Got New Job")
			jobsStream <- convertToStandardJobStruct(newJob)
		}
	}
}

func convertToStandardJobStruct(newJob *whoIsHiringJobStruct) *job.Job {
	return job.New()
}

func fetchJobs(whoIsHiringJobsStream chan *whoIsHiringJobStruct, scheduleAt string) {
	expr := cronexpr.MustParse(scheduleAt)
	nextTime := expr.Next(time.Now())

	for {
		fmt.Println("Getting Jobs from whoishiring")
		whoIsHiringJobsStream <- &whoIsHiringJobStruct{}
		time.Sleep(time.Duration(nextTime.Unix()))
		nextTime = expr.Next(time.Now())
	}
}
