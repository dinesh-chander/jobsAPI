package whoishiring

import (
	"models/job"
	"strings"
	jobType "types/jobs"
	"types/scrapers/whoishiring"
	"utils/geoUtils"
	miscellaneousUtils "utils/miscellaneous"
)

func convertToStandardJobStruct(newJob *whoishiring.WhoIsHiringType) (singleJob *jobType.Job) {
	singleJob = job.NewJob()

	singleJob.Company = newJob.Company
	singleJob.Description = newJob.Description
	singleJob.Address = newJob.Address
	singleJob.Is_Remote = newJob.Remote
	singleJob.Published_Date = newJob.Time
	singleJob.Title = newJob.Title
	singleJob.Job_Type = newJob.Kind
	singleJob.Source = newJob.Source
	singleJob.Source_Id = miscellaneousUtils.GenerateSHAChecksum(newJob.Description)
	singleJob.Source_Name = newJob.Source_name

	singleJob.Channel_Name = channelName
	singleJob.Tags = strings.Join(newJob.Tags, " # ")

	locationMap := make(map[string]string)
	geoUtils.GetLocationFromCoordinates(newJob.Location.Lat, newJob.Location.Lon, locationMap)
	singleJob.City = locationMap["locality"]
	singleJob.Country = locationMap["country"]

	return
}
