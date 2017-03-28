package jobs

type APIResponse struct {
	Id      uint
	Title   string
	City    string
	Country string
	Apply   string
	Company string
	Type    string
	//	Description string
	PublishedOn uint64
}

type APIResponseList struct {
	Count int
	Data  []APIResponse
}

func ConvertToResponse(jobsList []Job, resultCount int) (response *APIResponseList) {
	responseList := make([]APIResponse, len(jobsList))
	responseIndex := 0

	for _, newJob := range jobsList {
		var newResponseItem APIResponse

		newResponseItem.Id = newJob.ID
		newResponseItem.Apply = newJob.Apply
		newResponseItem.Title = newJob.Title
		newResponseItem.Country = newJob.Country
		newResponseItem.City = newJob.City
		newResponseItem.Company = newJob.Company
		newResponseItem.Type = newJob.Job_Type
		newResponseItem.PublishedOn = newJob.Published_Date

		responseList[responseIndex] = newResponseItem
		responseIndex = responseIndex + 1
	}

	response = &APIResponseList{
		Data:  responseList,
		Count: resultCount,
	}

	return
}
