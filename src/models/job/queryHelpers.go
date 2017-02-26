package job

import (
	"fmt"
	"strconv"
	"strings"
	jobType "types/jobs"
)

func buildSearchString(searchQuery *jobType.Query) (searchString string) {

	queryStringList := [3]string{}

	queryStringList[0] = formatSearchQuery(searchQuery.Locations, []string{"address", "city", "country"})
	queryStringList[1] = formatSearchQuery(searchQuery.Titles, []string{"title"})
	queryStringList[2] = formatSearchQuery(searchQuery.Keywords, []string{"description", "tags"})

	for _, value := range queryStringList {

		if value != "" {

			if searchString != "" {
				searchString = searchString + " AND " + value
			} else {
				searchString = value
			}
		}
	}

	return
}

func formatSearchQuery(searchStringsList []string, propertiesName []string) (searchSQLString string) {

	for index, value := range searchStringsList {
		searchStringsList[index] = strconv.Quote(value)
	}

	if len(searchStringsList) != 0 {
		searchSQLString = fmt.Sprintf(`MATCH(%s) AGAINST('%s' IN BOOLEAN MODE)`, strings.Join(propertiesName, ","), strings.Join(searchStringsList, " "))
	}

	return
}
