package scrapers

import "config"

import (
	"scrapers/angellist"
	"scrapers/betalist"
	"scrapers/whoishiring"
	jobType "types/jobs"
)

func InitiallizeScrappers(jobsStream chan *jobType.Job, fetchFrom int64, searchWordsList []string) {
	go whoishiring.GetWhoIsHiringJobs(jobsStream, config.GetConfig("whoishiring"), fetchFrom, searchWordsList)
	go betalist.GetBetaListJobs(jobsStream, config.GetConfig("betalist"), searchWordsList)
	go angellist.GetAngelListJobs(jobsStream, config.GetConfig("angellist"), searchWordsList)
}
