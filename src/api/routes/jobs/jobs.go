package jobs

import (
	routesRegistry "api/routesRegistry"
	"config"
	"encoding/json"
	"log"
	"logger"
	jobModel "models/job"
	http "net/http"
	"strings"
	jobType "types/jobs"
	httpUtils "utils/httpUtils"
)

var loggerInstance *log.Logger
var gzipSupport bool

func getAllJobs(response http.ResponseWriter, request *http.Request) (errMsg string, errCode int) {

	if request.Method == "GET" {

		query := &jobType.Query{}
		parseErr := query.ParseQueryParamsFromURL(request.URL)

		if parseErr != nil {

			errMsg = parseErr.Error()
			errCode = 400
			return
		}

		resultList := jobModel.FindContent(query)
		finalResult := jobType.ConvertToResponse(resultList)
		responseData, marshallingErr := json.Marshal(finalResult)

		if marshallingErr != nil {

			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(marshallingErr.Error()))
		} else if gzipSupport && strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {

			gzr := httpUtils.Pool.Get().(*httpUtils.GzipResponseWriter)
			gzr.ResponseWriter = response
			gzr.GW.Reset(response)

			defer func() {
				// gzr.w.Close will write a footer even if no data has been written.
				// StatusNotModified and StatusNoContent expect an empty body so don't close it.

				if gzr.StatusCode != http.StatusNotModified && gzr.StatusCode != http.StatusNoContent {
					if err := gzr.GW.Close(); err != nil {
						loggerInstance.Println(err.Error())
					}
				}
				httpUtils.Pool.Put(gzr)
			}()

			gzr.Header().Set("Content-Type", "application/json")
			gzr.WriteHeader(http.StatusOK)
			gzr.Write(responseData)
		} else {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusOK)
			response.Write(responseData)
		}
	} else {
		errMsg = "Wrong Method for the endpoint"
		errCode = 400
		return
	}

	return
}

func init() {
	loggerInstance = logger.Logger

	gzipSupport = false

	if config.GetConfig("gzip") == "true" {
		gzipSupport = true
	}

	routesRegistry.RouteRegistry["/jobs"] = getAllJobs
}
