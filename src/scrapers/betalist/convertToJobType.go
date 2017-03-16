package betalist

import (
	"models/job"
	"strings"
	jobType "types/jobs"
	"types/scrapers/betalist"
	"utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"
)

func convertToStandardJobStruct(newJob *betalist.BetalistType) (singleJob *jobType.Job) {
	singleJob = job.NewJob()

	singleJob.Company = newJob.Company_name
	singleJob.Description = newJob.Description_html
	singleJob.Address = newJob.Location
	singleJob.Is_Remote = newJob.Remote
	singleJob.Published_Date = newJob.Created_at_i
	singleJob.Title = newJob.Title
	singleJob.Job_Type = newJob.Commitment
	singleJob.Apply = "https://betalist.com" + newJob.Path
	singleJob.Source_Id = miscellaneousUtils.GenerateSHAChecksum(newJob.Description_html)
	singleJob.Source_Name = "bl"

	singleJob.Channel_Name = channelName
	singleJob.Tags = strings.Join(newJob.Tags, " # ")

	if (singleJob.City == "" || singleJob.Country == "") && singleJob.Address != "" {
		locationMap := make(map[string]string)
		geoUtils.GetLocationFromPlaceName(singleJob.Address, locationMap)

		if singleJob.City == "" {
			singleJob.City = locationMap["locality"]
		}

		if singleJob.Country == "" {
			singleJob.Country = locationMap["country"]
		}
	}

	return
}
