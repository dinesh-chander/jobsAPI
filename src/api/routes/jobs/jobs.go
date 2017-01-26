package jobs

import (
	routesRegistry "api/routesRegistry"
	"encoding/json"
	jobInterface "interfaces/jobs"
	jobModel "models/job"
	http "net/http"
)

func getAllJobs(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {

		query := &jobInterface.Query{}
		query.ParseQueryParamsFromURL(request.URL)

		resultList := jobModel.SearchContent(query)

		responseData, marshallingErr := json.Marshal(resultList)

		if marshallingErr != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(marshallingErr.Error()))
		} else {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusOK)
			response.Write(responseData)
		}

	} else {
		errMsg := []byte("Wrong Method for the endpoint")
		response.Write(errMsg)
	}
}

func init() {
	routesRegistry.RouteRegistry["/jobs"] = getAllJobs
}
