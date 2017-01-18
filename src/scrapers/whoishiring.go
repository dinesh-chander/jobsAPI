package scrapers

import (
	"fmt"
	"job"
	"net/http"
	"time"

	"strings"

	"github.com/gorhill/cronexpr"
)

type whoIsHiringJobStruct struct {
	address, country, city        string
	company, company_profile      string
	description, kind, url, title string
	source                        string
	remote                        bool
	time                          time.Duration
	tags, tags_share              []string
	apply, visa                   string
}

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string) {
	whoIsHiringJobsStream := make(chan *whoIsHiringJobStruct, 10)

	go fetchJobs(whoIsHiringJobsStream, scheduleAt)

	for {
		select {
		case newJob := <-whoIsHiringJobsStream:
			//			fmt.Println("Got New Job")
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
		//		fmt.Println("Getting Jobs from whoishiring")

		jobsList := makeRequest()

		for _, jobDetails := range jobsList {
			whoIsHiringJobsStream <- jobDetails
		}

		time.Sleep(time.Duration(nextTime.Unix()))
		nextTime = expr.Next(time.Now())
	}
}

func makeRequest() (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	postDataReader := strings.NewReader("{'query':{'bool':{'must':[],'should':[],'must_not':[],'filter':{'bool':{'must':[{'geo_bounding_box':{'location':{'bottom_left':{'lat':-70.8676081294354,'lon':123.61865624999996},'top_right':{'lat':83.82242395874371,'lon':-66.57665625000004}}}}],'should':[],'must_not':[]}}}},'sort':[{'_geo_distance':{'location':{'lat':'30.993','lon':'-151.479'},'order':'asc','unit':'km','distance_type':'plane'}},'_score',{'time':{'order':'desc','mode':'min'}}],'size':20}")
	response, err := httpClient.Post("https://search.whoishiring.io/item/item/_search?scroll=10m", "application/x-www-form-urlencoded", postDataReader)
	fmt.Println(response)

	if err != nil {
		panic(err)
	}

	jobsList = [](*whoIsHiringJobStruct){&whoIsHiringJobStruct{}}
	return
}
