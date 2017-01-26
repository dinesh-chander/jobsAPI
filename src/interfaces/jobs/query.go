package jobs

import (
	"net/url"
	"strings"
)

type Query struct {
	Location []string
	Title    []string
	Tags     []string
}

func (query *Query) ParseQueryParamsFromURL(url *url.URL) {
	queryParams := url.Query()

	if len(queryParams.Get("location")) > 0 {
		query.Location = strings.Split(queryParams.Get("location"), ",")
	}

	if len(queryParams.Get("title")) > 0 {
		query.Title = strings.Split(queryParams.Get("title"), ",")
	}

	if len(queryParams.Get("tags")) > 0 {
		query.Tags = strings.Split(queryParams.Get("tags"), ",")
	}
}
