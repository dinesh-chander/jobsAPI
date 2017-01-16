package scrapers

import (
	"fmt"
	"github.com/dispareil/jobsAPI/job"
	"github.com/gorhill/cronexpr"
	//"time"
)

type whoIsHiringJobStruct struct {
}

func GetWhoIsHiringJobs(jobsStream chan job.Job, scheduleAt string) {
	whoIsHiringJobsStream := make(chan whoIsHiringJobStruct, 5)

	go fetchJobs(whoIsHiringJobsStream, scheduleAt)

	for {
		select {
		case newJob := <-whoIsHiringJobsStream:
			fmt.Println("Got New Job")
			jobsStream <- convertToStandardJobStruct(newJob)
			//		default:
			//			fmt.Println("Got No New Job")
		}
	}
}

func convertToStandardJobStruct(newJob whoIsHiringJobStruct) job.Job {
	return &job.New()
}

func fetchJobs(whoIsHiringJobsStream chan whoIsHiringJobStruct, scheduleAt string) {

}
