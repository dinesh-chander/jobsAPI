package angellist

import (
	"encoding/json"
	"io"
	"log"
	"logger"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	jobType "types/jobs"
	//	cronParse "utils/cronParser"
	//	geoUtils "utils/geoUtils"
	"utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"

	gq "github.com/PuerkitoBio/goquery"
)

var channelName string
var loggerInstance *log.Logger
var httpClient *http.Client

func GetAngelListJobs(jobsStream chan *jobType.Job, scheduleAt string, fetchFrom int64, searchWordsList []string) {

	//	expr := cronParse.Parse(scheduleAt)

	//	for {

	//		jobsID := findJobIdsList()
	//		jobsURL := fetchJobsURL(jobsID)
	//		jobsList := getJobsFromJobsURL(jobsURL)

	//		for _, newJob := range jobsList {
	//			jobsStream <- newJob
	//		}

	//		nextTime := expr.Next(time.Now())
	//		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	//	}
}

func fetchJobsURL(idsList []int) (jobsURL []string) {

	var startupsIds string

	for index, jobId := range idsList {
		if index == len(idsList)-1 {
			startupsIds = startupsIds + "startup_ids[]=" + strconv.Itoa(jobId)
		} else {
			startupsIds = startupsIds + "startup_ids[]=" + strconv.Itoa(jobId) + ","
		}
	}

	if startupsIds != "" {

		response, fetchErr := httpClient.Get(`https://angel.co/job_listings/browse_startups_table?` + startupsIds)
		defer response.Body.Close()

		if fetchErr != nil {
			loggerInstance.Println(fetchErr.Error())
			return
		}

		jobsURL = fetchAllJobsURL(response.Body)
	} else {
		loggerInstance.Println("No startup id's found in the id's list")
	}

	return
}

func fetchAllJobsURL(pageMarkup io.ReadCloser) (jobsURL []string) {
	doc, err := gq.NewDocumentFromReader(pageMarkup)

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	doc.Find(".title").Each(func(_ int, s *gq.Selection) {
		href, found := s.Find("a").Attr("href")

		if !found {
			loggerInstance.Println("No HREF found for job")
			return
		}

		jobsURL = append(jobsURL, href)
	})

	return
}

func getJobsFromJobsURL(jobsURL []string) (jobsList []jobType.Job) {

	var wg sync.WaitGroup

	wg.Add(len(jobsURL))

	for _, jobURL := range jobsURL {

		go func(jobURL string) {
			defer wg.Done()

			newJob, ok := parseSingleJobPage(jobURL)

			if !ok {
				return
			}

			jobsList = append(jobsList, *newJob)

		}(jobURL)
	}

	wg.Wait()

	return
}

func parseSingleJobPage(jobURL string) (newJob *jobType.Job, ok bool) {

	newRequestInstance, newRequestInstanceError := http.NewRequest("GET", jobURL, nil)

	if newRequestInstanceError != nil {
		loggerInstance.Println(newRequestInstanceError.Error())
		return
	}

	newRequestInstance.Host = "angel.co"
	newRequestInstance.Header.Add("Accept", "*/*")
	newRequestInstance.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	response, fetchErr := httpClient.Do(newRequestInstance)

	if fetchErr != nil {
		loggerInstance.Println(fetchErr.Error())
		return
	}

	doc, err := gq.NewDocumentFromReader(response.Body)

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	newJob = &jobType.Job{}

	newJob.Description = doc.Find(".job-description").Text()
	newJob.Published_Date = uint64(time.Now().Unix() * 1000)
	newJob.Channel_Name = channelName

	doc.Find(".job-listing-metadata").Children().EachWithBreak(func(index int, s *gq.Selection) bool {

		if index == 1 {
			newJob.Tags = strings.Join(strings.Split(s.Text(), ","), " # ")
			return false
		}

		return true
	})

	newJob.Source = jobURL
	newJob.Source_Id = miscellaneousUtils.GenerateSHAChecksum(newJob.Description)
	newJob.Source_Name = "al"

	titleAndCompany := strings.Split(doc.Find(".company-summary").Find("h1").First().Text(), "at")

	if len(titleAndCompany) > 0 {
		newJob.Title = strings.Trim(titleAndCompany[0], "")
	}

	if len(titleAndCompany) > 1 {
		newJob.Company = strings.Trim(titleAndCompany[1], "")
	}

	locationAndJobType := strings.Split(doc.Find(".company-summary").Find("div").First().Text(), "·")

	if len(locationAndJobType) > 0 {
		newJob.Address = strings.Trim(locationAndJobType[0], "")
	}

	if len(locationAndJobType) > 1 {
		newJob.Job_Type = strings.Trim(locationAndJobType[1], "")
	}

	if newJob.Address != "" {
		locationMap := make(map[string]string)
		geoUtils.GetLocationFromPlaceName(newJob.Address, locationMap)
		newJob.City = locationMap["locality"]
		newJob.Country = locationMap["country"]
	}

	ok = true

	loggerInstance.Println(newJob)

	return
}

func findJobIdsList() (idsList []int) {

	pageResponse, fetchErr := httpClient.Get("https://angel.co/jobs#find/f!%7B%22keywords%22%3A%5B%22junior%22%5D%7D")

	if fetchErr != nil {
		loggerInstance.Println(fetchErr.Error())
		return
	}

	doc, err := gq.NewDocumentFromReader(pageResponse.Body)

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	idsListString, ok := doc.Find(".startup-container").Attr("data-startup_ids")

	if !ok {
		loggerInstance.Println("Unable To find the ID's Node")
		return
	}

	parseError := json.Unmarshal([]byte(idsListString), &idsList)

	if parseError != nil {
		loggerInstance.Println(parseError.Error())
		return
	}

	return
}

func init() {
	channelName = "angellist"
	loggerInstance = logger.Logger

	httpClient = &http.Client{
		Timeout: time.Second * 300,
	}
}
