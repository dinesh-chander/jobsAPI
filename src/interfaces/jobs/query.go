package jobs

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Locations []string
	Titles    []string
	Tags      []string
	Companies []string
	Keywords  []string
	Limit     int
	Skip      int
}

func (query *Query) ParseQueryParamsFromURL(url *url.URL) (parseErr error) {
	queryParams := url.Query()

	if len(queryParams.Get("locations")) > 0 {
		query.Locations = strings.Split(queryParams.Get("locations"), ",")
	}

	if len(queryParams.Get("titles")) > 0 {
		query.Titles = strings.Split(queryParams.Get("titles"), ",")
	}

	if len(queryParams.Get("companies")) > 0 {
		query.Companies = strings.Split(queryParams.Get("companies"), ",")
	}

	if len(queryParams.Get("tags")) > 0 {
		query.Tags = strings.Split(queryParams.Get("tags"), ",")
	}

	if len(queryParams.Get("keywords")) > 0 {
		query.Keywords = strings.Split(queryParams.Get("keywords"), ",")
	}

	if len(queryParams.Get("limit")) > 0 {
		var limit int64
		limit, parseErr = strconv.ParseInt(queryParams.Get("limit"), 10, 64)

		if parseErr != nil {
			return
		}

		query.Limit = int(limit)
	} else {
		query.Limit = 20
	}

	if len(queryParams.Get("skip")) > 0 {
		var skipCount int64

		skipCount, parseErr = strconv.ParseInt(queryParams.Get("skip"), 10, 64)

		if parseErr != nil {
			return
		}

		query.Skip = int(skipCount)
	} else {
		query.Skip = 0
	}

	return
}
