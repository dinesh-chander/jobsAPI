package scrapers

import (
	"encoding/json"
	"io/ioutil"
	"logger"
	"models/job"
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
	Time            uint64
	Tags            []string
	Apply           string
	Visa            bool
}

var loggerInstance = logger.Logger

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string) {
	expr := cronexpr.MustParse(scheduleAt)

	for {
		jobsList := makeRequestForNewJobs()

		for _, jobDetails := range jobsList {
			go func(newJob *whoIsHiringJobStruct) {
				jobsStream <- convertToStandardJobStruct(newJob)
			}(jobDetails)
		}

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func convertToStandardJobStruct(newJob *whoIsHiringJobStruct) (singleJob *job.Job) {
	singleJob = job.NewJob()

	singleJob.Company = newJob.Company
	singleJob.Description = newJob.Description
	singleJob.Address = newJob.Address
	singleJob.City = newJob.City
	singleJob.Country = newJob.Country
	singleJob.Is_Remote = newJob.Remote
	singleJob.Published_Date = newJob.Time
	singleJob.Title = newJob.Title
	singleJob.Source = newJob.Source
	singleJob.Source_Id = newJob.Id

	singleJob.Tags = strings.Join(newJob.Tags, " ")

	return
}

func makeRequestForNewJobs() (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	postDataReader := strings.NewReader(`{"query":{"bool":{"must":[],"should":[],"must_not":[],"filter":{"bool":{"must":[{"geo_bounding_box":{"location":{"bottom_left":{"lat":-70.8676081294354,"lon":123.61865624999996},"top_right":{"lat":83.82242395874371,"lon":-66.57665625000004}}}}],"should":[],"must_not":[]}}}},"sort":[{"_geo_distance":{"location":{"lat":"30.993","lon":"-151.479"},"order":"asc","unit":"km","distance_type":"plane"}},"_score",{"time":{"order":"desc","mode":"min"}}],"size":20}`)
	response, err := httpClient.Post("https://search.whoishiring.io/item/item/_search?scroll=10m", "application/x-www-form-urlencoded", postDataReader)

	if (err != nil) || (response.StatusCode != 200) {
		loggerInstance.Println(err)
	} else {
		responseBody, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			loggerInstance.Println("Response Body read error")
		} else {
			hits := gjson.GetBytes(responseBody, "hits.hits")
			jobsList = [](*whoIsHiringJobStruct){}

			hits.ForEach(func(key, value gjson.Result) bool {
				var jobDetails whoIsHiringJobStruct
				jobJSON := []byte(value.Get("_source").String())
				parseErr := json.Unmarshal(jobJSON, &jobDetails)

				if parseErr != nil {
					loggerInstance.Println("Unable to UNMARSHAL")
				} else {
					jobsList = append(jobsList, &jobDetails)
				}

				return true
			})
		}
	}

	return
}

func makeRequestForOldJobs() (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	postDataReader := strings.NewReader(`{"query":{"bool":{"must":[],"should":[],"must_not":[],"filter":{"bool":{"must":[{"geo_bounding_box":{"location":{"bottom_left":{"lat":-70.8676081294354,"lon":123.61865624999996},"top_right":{"lat":83.82242395874371,"lon":-66.57665625000004}}}}],"should":[],"must_not":[]}}}},"sort":[{"_geo_distance":{"location":{"lat":"30.993","lon":"-151.479"},"order":"asc","unit":"km","distance_type":"plane"}},"_score",{"time":{"order":"desc","mode":"min"}}],"size":20}`)
	response, err := httpClient.Post("https://search.whoishiring.io/item/item/_search?scroll=10m", "application/x-www-form-urlencoded", postDataReader)

	if (err != nil) || (response.StatusCode != 200) {
		loggerInstance.Println(err)
	} else {
		responseBody, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			loggerInstance.Println("Response Body read error")
		} else {
			hits := gjson.GetBytes(responseBody, "hits.hits")
			jobsList = [](*whoIsHiringJobStruct){}

			hits.ForEach(func(key, value gjson.Result) bool {
				var jobDetails whoIsHiringJobStruct
				jobJSON := []byte(value.Get("_source").String())
				parseErr := json.Unmarshal(jobJSON, &jobDetails)

				if parseErr != nil {
					loggerInstance.Println("Unable to UNMARSHAL")
				} else {
					jobsList = append(jobsList, &jobDetails)
				}

				return true
			})
		}
	}

	return
}
