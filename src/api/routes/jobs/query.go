package jobs

import (
	"net/url"
	"strings"
)

type Query struct {
	Location []string
	Role     []string
	Keywords []string
}

func (query *Query) parseQueryParamsFromURL(url *url.URL) {
	queryParams := url.Query()

	if len(queryParams.Get("location")) > 0 {
		query.Location = strings.Split(queryParams.Get("location"), ",")
	}

	if len(queryParams.Get("role")) > 0 {
		query.Role = strings.Split(queryParams.Get("role"), ",")
	}

	if len(queryParams.Get("keywords")) > 0 {
		query.Keywords = strings.Split(queryParams.Get("keywords"), ",")
	}
}
