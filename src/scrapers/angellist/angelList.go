package angellist

import (
	"log"
	"logger"
	"models/job"
	"strconv"
	"strings"
	"time"
	jobType "types/jobs"
	cronParse "utils/cronParser"
	"utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"

	gq "github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"utils/filters"
)

var channelName string
var loggerInstance *log.Logger
var batchSize int

func GetAngelListJobs(jobsStream chan *jobType.Job, scheduleAt string, searchWordsList []string) {

	expr := cronParse.Parse(scheduleAt)

	jobsURLChannel := make(chan string, 500)

	loggerInstance.Println("Starting", int(batchSize/2), "Angelist Job Fetchers")

	for workerIndex := 0; workerIndex < int(batchSize/2); workerIndex = workerIndex + 1 {

		go func(workerId int) {

			for {
				select {
				case newURLForProcessing := <-jobsURLChannel:

					loggerInstance.Println("New Job Processing Task Fetched By :", workerId)

					newJob, ok := parseSingleJobPage(newURLForProcessing)

					if ok {
						jobsStream <- newJob
					}

				}
			}
		}(workerIndex)
	}

	for {

		loggerInstance.Println("AngelList Loop Starts")

		startupIDList, jobsIDList := findJobIdsList(searchWordsList)

		fetchJobsURL(startupIDList, jobsIDList, jobsURLChannel)

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func fetchGroupedJobsListPage(pageParams string, jobsURLChannel chan string) {

	angelListURL := `https://angel.co/job_listings/browse_startups_table?`

	response, fetchErr := makeRequestToAngelListServer("GET", (angelListURL + pageParams), "", nil, true)

	if fetchErr != nil {
		loggerInstance.Println(fetchErr.Error())
	} else {

		fetchAllJobsURL(response, jobsURLChannel)
		response.Body.Close()
	}
}

func fetchJobsURL(startupIDList []int, jobsIDList [][]int, jobsURLChannel chan string) {
	var lastIndex int

	var batchedStartupIDs []int
	var batchedJobsIDList [][]int

	for index := 0; index < len(startupIDList); index = index + batchSize {

		lastIndex = index + batchSize

		if lastIndex >= len(startupIDList) {
			lastIndex = len(startupIDList) - 1
		}

		batchedStartupIDs = startupIDList[index:lastIndex]
		batchedJobsIDList = jobsIDList[index:lastIndex]

		urlParams := url.Values{}

		for selectedStartUpIDsIndex, startupId := range batchedStartupIDs {
			urlParams.Add("startup_ids[]", strconv.Itoa(startupId))

			for _, listingId := range batchedJobsIDList[selectedStartUpIDsIndex] {
				urlParams.Add("listing_ids["+strconv.Itoa(selectedStartUpIDsIndex)+"][]", strconv.Itoa(listingId))
			}
		}

		finalUrl, _ := url.QueryUnescape(urlParams.Encode())

		go fetchGroupedJobsListPage(finalUrl, jobsURLChannel)
	}
}

func fetchAllJobsURL(pageMarkup *http.Response, jobsURLChannel chan string) {
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

		loggerInstance.Println("Adding New URL for processing")

		jobsURLChannel <- href
	})

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

		if strings.ToLower(doc.Find(".job-listing-metadata").Children().Eq(0).Text()) == "skills" {
			if doc.Find(".job-listing-metadata").Children().Eq(1).Text() != "" {
				newJob.Tags = strings.Join(strings.Split(doc.Find(".job-listing-metadata").Children().Eq(1).Text(), ","), " # ")
			}
		}

		newJob.Source = jobURL
		newJob.Source_Id = miscellaneousUtils.GenerateSHAChecksum(newJob.Description)
		newJob.Source_Name = "al"

		if newJob.Address != "" {
			locationMap := make(map[string]string)
			geoUtils.GetLocationFromPlaceName(newJob.Address, locationMap)
			newJob.City = strings.TrimSpace(locationMap["locality"])
			newJob.Country = strings.TrimSpace(locationMap["country"])
		}

		ok = true
	} else {
		loggerInstance.Println("Rejecting Job : ", newJob.Title)
	}

	return
}

func findJobIdsList(searchWordsList []string) (startupIDList []int, jobsIDList [][]int) {

	var ifJobsIDFound bool

	startupIDList, jobsIDList, ifJobsIDFound = getJobsPage(searchWordsList)

	if !ifJobsIDFound {
		startupIDList = []int{}
		jobsIDList = [][]int{}
	} else {
		startupIDList = startupIDList[:10]
		jobsIDList = jobsIDList[:10]
	}

	return
}

func init() {
	batchSize = 10
	channelName = "angellist"
	loggerInstance = logger.Logger
}
