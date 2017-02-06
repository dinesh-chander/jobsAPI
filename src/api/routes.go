package api

import (
	routeRegistry "api/routesRegistry"
	"net/http"
)

func pathHandler(response http.ResponseWriter, request *http.Request) {
	if routeRegistry.RouteRegistry[request.URL.Path] != nil {
		errMsg, errCode := (routeRegistry.RouteRegistry[request.URL.Path](response, request))

		if errMsg != "" {
			http.Error(response, errMsg, errCode)
		}
	} else {
		http.Error(response, "Path or Method not found", 400)
	}
}
