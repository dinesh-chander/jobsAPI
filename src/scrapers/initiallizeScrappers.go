package scrapers

import "config"
import jobType "types/jobs"

import "scrapers/angellist"

func InitiallizeScrappers(jobsStream chan *jobType.Job, fetchFrom int64, searchWordsList []string) {

	//	go whoishiring.GetWhoIsHiringJobs(jobsStream, config.GetConfig("whoishiring"), fetchFrom, searchWordsList)
	go angellist.GetAngelListJobs(jobsStream, config.GetConfig("whoishiring"), fetchFrom, searchWordsList)
}
