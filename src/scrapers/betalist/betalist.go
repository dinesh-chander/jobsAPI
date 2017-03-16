package betalist

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"logger"
	"net/http"
	"strings"
	"time"
	jobType "types/jobs"
	"types/scrapers/betalist"
	cronParse "utils/cronParser"
	"utils/filters"

	"github.com/tidwall/gjson"
)

var channelName string
var loggerInstance *log.Logger

func GetBetaListJobs(jobsStream chan *jobType.Job, scheduleAt string, searchWordsList []string) {

	expr := cronParse.Parse(scheduleAt)

	loggerInstance.Println("betalist Scraper Started")

	for {
		makeRequestForNewJobs(searchWordsList, jobsStream)
		nextTime := expr.Next(time.Now())
		time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
	}
}

func makeRequestForNewJobs(searchWordsList []string, jobsStream chan *jobType.Job) (jobsList [](*betalist.BetalistType)) {

	httpClient := &http.Client{
		Timeout: time.Duration(300) * time.Second,
	}

	postData := `{"requests":[{"indexName":"Jobs_Post_production","params":"query=developer,` + strings.Join(searchWordsList, ",") + `&hitsPerPage=1000&maxValuesPerFacet=500&attributesToRetrieve=path%2Ccompany_name%2Cdescription_html%2Ccommitment%2Ccity%2Clocation%2Ctitle%2Ccountry%2Csource_Id%2Cremote%2Ccreated_at_i%2C_tags&filters=&facets=%5B%22remote%22%2C%22commitment%22%5D&tagFilters=&facetFilters=%5B%5B%22commitment%3AFull-Time%22%2C%22commitment%3APart-Time%22%2C%22commitment%3AInternship%22%2C%22commitment%3AContractor%22%5D%5D"}]}`

	postDataReader := strings.NewReader(postData)

	response, err := httpClient.Post("https://4cqmtmmk73-dsn.algolia.net/1/indexes/*/queries?x-algolia-agent=Algolia%20for%20vanilla%20JavaScript%20(lite)%203.21.1%3Binstantsearch.js%201.11.2%3BJS%20Helper%202.18.1&x-algolia-application-id=4CQMTMMK73&x-algolia-api-key=5defc918b05014c76cb000f5c9386c9b", "application/x-www-form-urlencoded", postDataReader)

	if (err != nil) || (response.StatusCode != 200) {
		loggerInstance.Println(err)
	} else {
		responseBody, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			loggerInstance.Println(readErr.Error())
		} else {

			hits := gjson.GetBytes(responseBody, "results.0.hits")

			hits.ForEach(func(key, value gjson.Result) bool {

				if value.String() != "" {

					var jobDetails betalist.BetalistType
					jobJSON := []byte(value.String())

					if len(jobJSON) > 0 {
						parseErr := json.Unmarshal(jobJSON, &jobDetails)
						if parseErr != nil {
							loggerInstance.Println(parseErr.Error(), value.String())
						} else if filters.IsValidJob(jobDetails.Title) {
							go func(jobDetails *betalist.BetalistType) {
								jobsStream <- convertToStandardJobStruct(jobDetails)
							}(&jobDetails)
						}
					}
				}

				return true
			})
		}
	}

	return
}

func init() {
	channelName = "betalist"
	loggerInstance = logger.Logger
}
