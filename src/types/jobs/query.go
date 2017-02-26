package jobs

import (
	"net/url"
	"strconv"
	"strings"
)

type Query struct {
	Locations []string
	Titles    []string
	Keywords  []string
	Limit     int
	Skip      int
}

func (query *Query) ParseQueryParamsFromURL(url *url.URL) (parseErr error) {
	queryParams := url.Query()

	locations := strings.TrimSpace(queryParams.Get("locations"))
	titles := strings.TrimSpace(queryParams.Get("titles"))
	keywords := strings.TrimSpace(queryParams.Get("keywords"))
	limit := strings.TrimSpace(queryParams.Get("limit"))
	skip := strings.TrimSpace(queryParams.Get("skip"))

	if len(locations) > 0 {

		locationsList := strings.Split(locations, ",")

		if len(locationsList) > 0 {
			query.Locations = locationsList
		}
	}

	if len(titles) > 0 {

		titlesList := strings.Split(titles, ",")

		if len(titlesList) > 0 {
			query.Titles = titlesList
		}
	}

	//	if len(queryParams.Get("companies")) > 0 {
	//		query.Companies = strings.Split(queryParams.Get("companies"), ",")
	//	}

	//	if len(queryParams.Get("tags")) > 0 {
	//		query.Tags = strings.Split(queryParams.Get("tags"), ",")
	//	}

	if len(keywords) > 0 {

		keywordsList := strings.Split(keywords, ",")

		if len(keywordsList) > 0 {
			query.Keywords = keywordsList
		}
	}

	if len(limit) > 0 {

		var limitVal int64
		limitVal, parseErr = strconv.ParseInt(limit, 10, 64)

		if parseErr != nil {
			return
		}

		query.Limit = int(limitVal)
	} else {
		query.Limit = 20
	}

	if len(skip) > 0 {

		var skipVal int64

		skipVal, parseErr = strconv.ParseInt(skip, 10, 64)

		if parseErr != nil {
			return
		}

		query.Skip = int(skipVal)
	} else {
		query.Skip = 0
	}

	return
}
