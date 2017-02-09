package routeRegistry

import (
	"net/http"
)

var RouteRegistry map[string]func(http.ResponseWriter, *http.Request) (string, int)

func init() {
	RouteRegistry = make(map[string]func(http.ResponseWriter, *http.Request) (string, int))
}
