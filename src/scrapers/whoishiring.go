package scrapers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"job"
	"net/http"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/tidwall/gjson"
)

type whoIsHiringJobStruct struct {
	Id              string
	Address         string
	Country         string
	City            string
	Company         string
	Company_profile string
	Description     string
	Kind            string
	Url             string
	Title           string
	Source          string
	Source_name     string
	Remote          bool
	Time            int
	Tags            []string
	Tags_share      []string
	Apply           string
	Visa            bool
}

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string) {
	whoIsHiringJobsStream := make(chan *whoIsHiringJobStruct, 10)

	go fetchJobs(whoIsHiringJobsStream, scheduleAt)

	for {
		select {
		case newJob := <-whoIsHiringJobsStream:
			jobsStream <- convertToStandardJobStruct(newJob)
		}
	}
}

func convertToStandardJobStruct(newJob *whoIsHiringJobStruct) (singleJob *job.Job) {
	singleJob = job.New()

	singleJob.Company = newJob.Company
	singleJob.Description = newJob.Description
	singleJob.Location = newJob.Address
	singleJob.IsRemote = newJob.Remote
	singleJob.PublishedDate = newJob.Time
	singleJob.Title = newJob.Title
	//	singleJob.Tags = newJob.Tags
	//	singleJob.Share_Tags = newJob.Tags_share
	singleJob.Source = newJob.Source

	return
}

func fetchJobs(whoIsHiringJobsStream chan *whoIsHiringJobStruct, scheduleAt string) {
	expr := cronexpr.MustParse(scheduleAt)

	for {
		jobsList := makeRequest()

		for _, jobDetails := range jobsList {
			whoIsHiringJobsStream <- jobDetails
		}

		nextTime := expr.Next(time.Now())
		fmt.Println("Sleeping for :", nextTime.Unix()-time.Now().Unix())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func makeRequest() (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	postDataReader := strings.NewReader(`{"query":{"bool":{"must":[],"should":[],"must_not":[],"filter":{"bool":{"must":[{"geo_bounding_box":{"location":{"bottom_left":{"lat":-70.8676081294354,"lon":123.61865624999996},"top_right":{"lat":83.82242395874371,"lon":-66.57665625000004}}}}],"should":[],"must_not":[]}}}},"sort":[{"_geo_distance":{"location":{"lat":"30.993","lon":"-151.479"},"order":"asc","unit":"km","distance_type":"plane"}},"_score",{"time":{"order":"desc","mode":"min"}}],"size":20}`)
	response, err := httpClient.Post("https://search.whoishiring.io/item/item/_search?scroll=10m", "application/x-www-form-urlencoded", postDataReader)

	if (err != nil) || (response.StatusCode != 200) {
		panic(err)
	}

	responseBody, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		fmt.Println("Response Body read error")
	} else {
		hits := gjson.GetBytes(responseBody, "hits.hits")
		jobsList = [](*whoIsHiringJobStruct){}

		hits.ForEach(func(key, value gjson.Result) bool {
			var jobDetails whoIsHiringJobStruct
			jobJSON := []byte(value.Get("_source").String())
			parseErr := json.Unmarshal(jobJSON, &jobDetails)

			if parseErr != nil {
				fmt.Println("Unable to UNMARSHAL")
			} else {
				jobsList = append(jobsList, &jobDetails)
			}

			return true
		})
	}

	return
}
