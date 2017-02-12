package jobs

type APIResponse struct {
	Title       string
	City        string
	Country     string
	Apply       string
	Company     string
	Type        string
	Description string
	PublishedOn uint64
}

func ConvertToResponse(jobsList []Job) (response []APIResponse) {
	response = make([]APIResponse, len(jobsList))
	responseIndex := 0

	for _, newJob := range jobsList {
		var newResponseItem APIResponse

		newResponseItem.Apply = newJob.Apply
		newResponseItem.Title = newJob.Title
		newResponseItem.Country = newJob.Country
		newResponseItem.City = newJob.City
		newResponseItem.Company = newJob.Company
		newResponseItem.Type = newJob.Job_Type
		newResponseItem.Description = newJob.Description
		newResponseItem.PublishedOn = newJob.Published_Date

		response[responseIndex] = newResponseItem
		responseIndex = responseIndex + 1
	}

	return response
}
