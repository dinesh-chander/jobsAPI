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

	gq "github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"utils/filters"
)

type jobListDetails struct {
	jobId  string
	jobURL string
}

type companyAndItsJobs struct {
	startupId  string
	jobIdsList [](*jobListDetails)
}

var channelName string
var loggerInstance *log.Logger
var batchSize int

func GetAngelListJobs(jobsStream chan *jobType.Job, scheduleAt string, searchWordsList []string) {

	expr := cronParse.Parse(scheduleAt)

	jobsURLChannel := make(chan *jobListDetails, 2000)

	loggerInstance.Println("Starting", int(batchSize/2), "Angelist Job Fetchers")

	for workerIndex := 0; workerIndex < 1; workerIndex = workerIndex + 1 {

		go func(workerId int) {

			for {
				select {
				case newJobDetailsForProcessing := <-jobsURLChannel:

					loggerInstance.Println(workerId, ":", newJobDetailsForProcessing.jobURL)

					newJob, ok := parseSingleJobPage(newJobDetailsForProcessing)

					if ok {
						jobsStream <- newJob
					}

					time.Sleep(time.Second * 40)
				}
			}
		}(workerIndex)
	}

	loggerInstance.Println("AngelList Scraper Started")

	for {
		companyAndJobIds := findJobIdsList(searchWordsList)

		fetchJobsURL(companyAndJobIds, jobsURLChannel)

		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func fetchGroupedJobsListPage(batchedCompanyAndItsJobs [](*companyAndItsJobs), pageParams string, jobsURLChannel chan *jobListDetails) {

	angelListURL := `https://angel.co/job_listings/browse_startups_table?`

	response, fetchErr := makeRequestToAngelListServer("GET", (angelListURL + pageParams), "", nil, true)

	if fetchErr != nil {
		loggerInstance.Println(fetchErr.Error())
	} else {
		fetchAllJobsURL(batchedCompanyAndItsJobs, response, jobsURLChannel)
		response.Body.Close()
	}
}

func fetchJobsURL(companyAndJobIds [](*companyAndItsJobs), jobsURLChannel chan *jobListDetails) {

	var lastIndex int

	var listingDetails *jobListDetails
	var batchedCompanyAndItsJobs [](*companyAndItsJobs)

	for index := 0; index < len(companyAndJobIds); index = index + batchSize {

		lastIndex = index + batchSize

		if lastIndex >= len(companyAndJobIds) {
			lastIndex = len(companyAndJobIds) - 1
		}

		batchedCompanyAndItsJobs = companyAndJobIds[index:lastIndex]

		urlParams := url.Values{}

		for selectedStartUpIDsIndex, companyDetails := range batchedCompanyAndItsJobs {
			urlParams.Add("startup_ids[]", companyDetails.startupId)

			for _, listingDetails = range companyDetails.jobIdsList {
				urlParams.Add("listing_ids["+strconv.Itoa(selectedStartUpIDsIndex)+"][]", listingDetails.jobId)
			}
		}

		finalUrl, _ := url.QueryUnescape(urlParams.Encode())

		fetchGroupedJobsListPage(batchedCompanyAndItsJobs, finalUrl, jobsURLChannel)

		time.Sleep(time.Minute * 5)
	}
}

func fetchAllJobsURL(batchedCompanyAndItsJobs [](*companyAndItsJobs), pageMarkup *http.Response, jobsURLChannel chan *jobListDetails) {

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

		loggerInstance.Println("Adding for Processing :", href)

		for _, companyDetails := range batchedCompanyAndItsJobs {

			for _, jobDetails := range companyDetails.jobIdsList {

				if strings.Contains(href, jobDetails.jobId) {
					jobDetails.jobURL = href
					jobsURLChannel <- jobDetails
					return
				}
			}
		}

	})

	return
}

func parseSingleJobPage(jobDetails *jobListDetails) (newJob *jobType.Job, ok bool) {

	response, fetchErr := makeRequestToAngelListServer("GET", jobDetails.jobURL, "", nil, true)

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

	titleAndCompany := strings.Split(doc.Find(".company-summary").Find("h1").First().Text(), " at ")

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

		newJob.Apply = jobDetails.jobURL

		newJob.Source = jobDetails.jobURL
		newJob.Source_Id = jobDetails.jobId
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

func findJobIdsList(searchWordsList []string) (companyAndJobIds [](*companyAndItsJobs)) {

	var ifJobsIDFound bool

	startupIDList, jobsIDList, ifJobsIDFound := getJobsPage(searchWordsList)

	companyAndJobIds = [](*companyAndItsJobs){}

	if ifJobsIDFound {

		var allJobsIDs []string
		var jobIdString string
		var companyDetails *companyAndItsJobs
		var newJobListDetails *jobListDetails
		var jobIdsList [](*jobListDetails)

		for startupIndex, startupId := range startupIDList {

			companyDetails = &companyAndItsJobs{}

			companyDetails.startupId = strconv.Itoa(startupId)

			jobIdsList = [](*jobListDetails){}

			for _, jobId := range jobsIDList[startupIndex] {

				jobIdString = strconv.Itoa(jobId)

				newJobListDetails = &jobListDetails{
					jobId: jobIdString,
				}

				jobIdsList = append(jobIdsList, newJobListDetails)
				allJobsIDs = append(allJobsIDs, jobIdString)
			}

			companyDetails.jobIdsList = jobIdsList
			companyAndJobIds = append(companyAndJobIds, companyDetails)
		}

		newJobIDsList := job.FindAlreadyPresentJobsWithGivenSourceIds(channelName, allJobsIDs)

		newJobIdsMap := make(map[string]bool)

		for _, jobIdString := range newJobIDsList {
			newJobIdsMap[jobIdString] = true
		}

		var idIndex int

		loggerInstance.Println("Old Length :", len(companyAndJobIds))

		for companyIndex := 0; companyIndex < len(companyAndJobIds); {

			companyDetails = companyAndJobIds[companyIndex]

			jobIdsList = companyDetails.jobIdsList

			for idIndex, newJobListDetails = range jobIdsList {

				if !newJobIdsMap[newJobListDetails.jobId] {
					jobIdsList[idIndex] = jobIdsList[len(jobIdsList)-1]
					jobIdsList = jobIdsList[:len(jobIdsList)-1]
				}
			}

			if len(jobIdsList) == 0 {
				companyAndJobIds[companyIndex] = companyAndJobIds[len(companyAndJobIds)-1]
				companyAndJobIds = companyAndJobIds[:len(companyAndJobIds)-1]
			} else {
				companyIndex = companyIndex + 1
			}
		}

		loggerInstance.Println("New Length :", len(companyAndJobIds))
	}

	return
}

func init() {
	batchSize = 10
	channelName = "angellist"
	loggerInstance = logger.Logger
}
