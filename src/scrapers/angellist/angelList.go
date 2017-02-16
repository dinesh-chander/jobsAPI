package angellist

import (
	"log"
	"logger"
	"models/job"
	"strconv"
	"strings"
	"sync"
	"time"
	jobType "types/jobs"
	cronParse "utils/cronParser"
	"utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"

	gq "github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"utils/filters"
)

var channelName string
var loggerInstance *log.Logger
var batchSize int

func GetAngelListJobs(jobsStream chan *jobType.Job, scheduleAt string, searchWordsList []string) {

	expr := cronParse.Parse(scheduleAt)

	jobsURLChannel := make(chan []string, 300)

	go func() {

		for {
			select {
			case newURLBatch := <-jobsURLChannel:

				jobsList := getJobsFromJobsURL(newURLBatch)

				for _, newJobReference := range jobsList {

					if newJobReference != nil {
						jobsStream <- newJobReference
					}
				}
			}
		}
	}()

	for {

		jobsID := findJobIdsList(searchWordsList)

		fetchJobsURL(jobsID, jobsURLChannel)

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func fetchJobsURL(idsList []int, jobsURLChannel chan []string) {
	var lastIndex int
	var batchedIDs []int

	angelListURL := `https://angel.co/job_listings/browse_startups_table?`

	for index := 0; index < len(idsList); index = index + batchSize {

		var startupsIds string

		lastIndex = index + batchSize

		if lastIndex >= len(idsList) {
			lastIndex = len(idsList) - 1
		}

		batchedIDs = idsList[index:lastIndex]

		for selectedIDsIndex, jobId := range batchedIDs {
			if selectedIDsIndex == len(batchedIDs)-1 {
				startupsIds = startupsIds + "startup_ids[]=" + strconv.Itoa(jobId)
			} else {
				startupsIds = startupsIds + "startup_ids[]=" + strconv.Itoa(jobId) + "&"
			}
		}

		go func(pageParams string) {

			response, fetchErr := makeRequestToAngelListServer("GET", (angelListURL + pageParams), "", nil, true)

			if fetchErr != nil {
				loggerInstance.Println(fetchErr.Error())
				return
			}

			jobsURLList := fetchAllJobsURL(response)

			response.Body.Close()

			if jobsURLList != nil {
				jobsURLChannel <- jobsURLList
			}

		}(startupsIds)
	}
}

func fetchAllJobsURL(pageMarkup *http.Response) (jobsURL []string) {
	doc, err := gq.NewDocumentFromResponse(pageMarkup)

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	doc.Find(".title a").Each(func(_ int, s *gq.Selection) {
		href, found := s.Attr("href")

		if !found {
			loggerInstance.Println("No HREF found for job : ")
			return
		}

		jobsURL = append(jobsURL, href)
	})

	return
}

func getJobsFromJobsURL(jobsURL []string) (jobsList [](*jobType.Job)) {

	var wg sync.WaitGroup

	wg.Add(len(jobsURL))

	var ml sync.Mutex

	jobsList = [](*jobType.Job){}

	for _, jobURL := range jobsURL {

		go func(jobURL string) {

			defer wg.Done()

			time.Sleep(time.Duration(rand.Intn(900)) * time.Second)

			newJob, ok := parseSingleJobPage(jobURL)

			if !ok {
				return
			}

			ml.Lock()

			jobsList = append(jobsList, newJob)

			ml.Unlock()

		}(jobURL)
	}

	wg.Wait()

	return
}

func parseSingleJobPage(jobURL string) (newJob *jobType.Job, ok bool) {

	response, fetchErr := makeRequestToAngelListServer("GET", jobURL, "", nil, true)

	if fetchErr != nil {
		loggerInstance.Println(fetchErr.Error())
		return
	}

	defer response.Body.Close()

	doc, err := gq.NewDocumentFromReader(response.Body)

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	newJob = job.NewJob()

	titleAndCompany := strings.Split(doc.Find(".company-summary").Find("h1").First().Text(), "at")

	if len(titleAndCompany) > 0 {
		newJob.Title = strings.TrimSpace(titleAndCompany[0])
	}

	if filters.IsValidJob(newJob.Title) {

		if len(titleAndCompany) > 1 {
			newJob.Company = strings.TrimSpace(titleAndCompany[1])
		}

		locationAndJobType := strings.Split(doc.Find(".company-summary").Find("div").First().Text(), "·")

		if len(locationAndJobType) > 0 {
			newJob.Address = strings.TrimSpace(locationAndJobType[0])
		}

		if len(locationAndJobType) > 1 {
			newJob.Job_Type = strings.TrimSpace(locationAndJobType[1])
		}

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

		if newJob.Address != "" {
			locationMap := make(map[string]string)
			geoUtils.GetLocationFromPlaceName(newJob.Address, locationMap)
			newJob.City = locationMap["locality"]
			newJob.Country = locationMap["country"]
		}

		ok = true
	} else {
		loggerInstance.Println("Rejecting Job : ", newJob.Title)
	}

	return
}

func findJobIdsList(searchWordsList []string) (idsList []int) {

	var ifJobsIDFound bool

	idsList, ifJobsIDFound = getJobsPage(searchWordsList)

	if !ifJobsIDFound {
		idsList = []int{}
	}
	//else {
	//	loggerInstance.Println(len(idsList))
	//	if len(idsList) > 20 {
	//		idsList = idsList[:20]
	//	}
	//}

	return
}

func init() {
	batchSize = 20
	channelName = "angellist"
	loggerInstance = logger.Logger
}
