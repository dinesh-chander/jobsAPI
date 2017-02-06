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
	geoUtils "utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"

	"github.com/tidwall/gjson"
)

type Location struct {
	Lat, Lon float64
}

type whoIsHiringJobStruct struct {
	Apply       string
	Address     string
	Company     string
	Description string
	Kind        string
	Location    Location
	Url         string
	Title       string
	Source      string
	Source_name string
	Remote      bool
	Time        uint64
	Tags        []string
}

var channelName string
var loggerInstance *log.Logger

func GetWhoIsHiringJobs(jobsStream chan *job.Job, scheduleAt string, fetchFrom int64, searchWordsList []string) {

	expr := cronParse.Parse(scheduleAt)

	timestampOfLastEntryInDB := job.FindLastAddedEntryTimestampForSource(channelName)

	if timestampOfLastEntryInDB > fetchFrom {
		fetchFrom = timestampOfLastEntryInDB
	}

	for {

		jobsList := makeRequestForNewJobs(fetchFrom, searchWordsList)

		for _, jobDetails := range jobsList {
			jobsStream <- convertToStandardJobStruct(jobDetails)
		}

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
		fetchFrom = job.FindLastAddedEntryTimestampForSource(channelName)
	}
}

func makeRequestForNewJobs(lastFetchedJobTimeInMilliSeconds int64, searchWordsList []string) (jobsList [](*whoIsHiringJobStruct)) {
	httpClient := &http.Client{
		Timeout: time.Second * 300,
	}

	var searchQuery string

	if len(searchWordsList) > 0 {
		searchQuery = `{"bool":{"should":[`

		for index, searchWord := range searchWordsList {
			if index == (len(searchWordsList) - 1) {
				searchQuery = searchQuery + `{"bool":{"should":[{"match":{"title":{"query":"` + searchWord + `","boost":5,"type":"phrase"}}}]}}`
			} else {
				searchQuery = searchQuery + `{"bool":{"should":[{"match":{"title":{"query":"` + searchWord + `","boost":5,"type":"phrase"}}}]}},`
			}
		}

		searchQuery = searchQuery + `],"minimum_should_match":1}}`
	}

	postData := `{"query":{"bool":{"must":[],"should":[` + searchQuery + `],"must_not":[],"filter":{"bool":{"must":[{"geo_bounding_box":{"location":{"bottom_left":{"lat":-70.8676081294354,"lon":86.70474999999999},"top_right":{"lat":83.82242395874371,"lon":-90.48275000000001}}}},{"range":{"time":{"gt":` + strconv.FormatInt(lastFetchedJobTimeInMilliSeconds, 10) + `}}}],"should":[],"must_not":[]}}}},"sort":[{"time":{"order":"desc","mode":"min"}}],"size":20000}`

	postDataReader := strings.NewReader(postData)
	response, err := httpClient.Post("https://search.whoishiring.io/item/item/_search?scroll=10m", "application/x-www-form-urlencoded", postDataReader)

	if (err != nil) || (response.StatusCode != 200) {
		loggerInstance.Println(err)
	} else {
		responseBody, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			loggerInstance.Println(readErr.Error())
		} else {
			hits := gjson.GetBytes(responseBody, "hits.hits")

			jobsList = [](*whoIsHiringJobStruct){}

			hits.ForEach(func(key, value gjson.Result) bool {

				if value.String() != "" {

					if value.Get("_source").String() != "" {

						var jobDetails whoIsHiringJobStruct
						jobJSON := []byte(value.Get("_source").String())

						if len(jobJSON) > 0 {
							parseErr := json.Unmarshal(jobJSON, &jobDetails)

							if parseErr != nil {
								loggerInstance.Println(parseErr.Error(), value.Get("_source").String())
							} else {
								jobsList = append(jobsList, &jobDetails)
							}
						}
					}
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
	singleJob.Is_Remote = newJob.Remote
	singleJob.Published_Date = newJob.Time
	singleJob.Title = newJob.Title
	singleJob.Job_Type = newJob.Kind
	singleJob.Source = newJob.Source
	singleJob.Source_Id = miscellaneousUtils.GenerateSHAChecksum(newJob.Title + newJob.Description)
	singleJob.Source_Name = newJob.Source_name

	singleJob.Channel_Name = channelName

	singleJob.Tags = strings.Join(newJob.Tags, " # ")

	locationMap := make(map[string]string)
	geoUtils.GetLocationFromGeoCode(newJob.Location.Lat, newJob.Location.Lon, locationMap)
	singleJob.City = locationMap["locality"]
	singleJob.Country = locationMap["country"]

	return
}

func init() {
	channelName = "whoishiring"
	loggerInstance = logger.Logger
}
