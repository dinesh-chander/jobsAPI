package jobs

import (
	routesRegistry "api/routesRegistry"
	jobModel "models/job"
	http "net/http"
)

func getAllJobs(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)

		query := &Query{}
		query.parseQueryParamsFromURL(request.URL)
		jobModel.SearchContent("")

		responseData := []byte("Send Actual Data")
		response.Write(responseData)
	} else {
		errMsg := []byte("Wrong Method for the endpoint")
		response.Write(errMsg)
	}
}

func init() {
	routesRegistry.RouteRegistry["/jobs"] = getAllJobs
}
