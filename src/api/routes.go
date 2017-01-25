package api

import (
	_ "api/routes/jobs"
	routeRegistry "api/routesRegistry"
	"net/http"
)

func pathHandler(response http.ResponseWriter, request *http.Request) {
	if routeRegistry.RouteRegistry[request.URL.Path] != nil {
		routeRegistry.RouteRegistry[request.URL.Path](response, request)
	} else {
		defaultHandler(response, request)
	}
}

func defaultHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(404)
	response.Write([]byte("Path or Method not found"))
}
