package scrapers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"logger"
	"models/job"
	"net/http"
	"strconv"
	"strings"
	"time"
	cronParse "utils/cronParser"

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

var channelName string
var loggerInstance *log.Logger

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string, fetchFrom int64) {

	expr := cronParse.Parse(scheduleAt)

	timestampOfLastEntryInDB := job.FindLastAddedEntryTimestampForSource(channelName)

	if timestampOfLastEntryInDB > fetchFrom {
		fetchFrom = timestampOfLastEntryInDB
	}

	for {

		jobsList := makeRequestForNewJobs(fetchFrom)

		for _, jobDetails := range jobsList {
			jobsStream <- convertToStandardJobStruct(jobDetails)
		}

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
		fetchFrom = job.FindLastAddedEntryTimestampForSource(channelName)
	}
}

func makeRequestForNewJobs(lastFetchedJobTimeInMilliSeconds int64) (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 120,
	}

	postData := `{"query":{"bool":{"must":[],"should":[],"must_not":[],"filter":{"bool":{"must":[{"geo_bounding_box":{"location":{"bottom_left":{"lat":-70.8676081294354,"lon":123.61865624999996},"top_right":{"lat":83.82242395874371,"lon":-156.57665625000004}}}},{"range":{"time":{"gt":` + strconv.FormatInt(lastFetchedJobTimeInMilliSeconds, 10) + `}}}],"should":[],"must_not":[]}}}},"sort":[{"time":{"order":"desc","mode":"min"}}],"size":20000}`

	postDataReader := strings.NewReader(postData)
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
	singleJob.Source_Name = newJob.Source_name
	singleJob.Source_Id = newJob.Id
	singleJob.Channel_Name = channelName

	singleJob.Tags = strings.Join(newJob.Tags, " ")

	return
}

func init() {
	channelName = "whoishiring"
	loggerInstance = logger.Logger
}
