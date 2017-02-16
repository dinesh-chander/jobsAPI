package job

import (
	"fmt"
	"strconv"
	"strings"
	jobType "types/jobs"
)

func FindContent(searchQuery *jobType.Query) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchQuerySQLString := buildSearchString(searchQuery)

	if searchQuery.Limit == 0 {
		searchResult = []jobType.Job{}
		return
	}

	if searchQuerySQLString != "" {
		return findFromSearchableTable(searchQuerySQLString, searchQuery.Limit, searchQuery.Skip)
	} else {
		return findFromNormalTable(searchQuery.Limit, searchQuery.Skip)
	}
}

func buildSearchString(searchQuery *jobType.Query) (searchString string) {

	queryStringList := [3]string{}

	queryStringList[0] = formatSearchQuery(searchQuery.Locations, "Location")
	queryStringList[1] = formatSearchQuery(searchQuery.Titles, "Title")
	queryStringList[2] = formatSearchQuery(searchQuery.Keywords, "Description")

	for _, value := range queryStringList {

		if value != "" {

			if searchString != "" {
				searchString = searchString + " INTERSECT " + value
			} else {
				searchString = value
			}
		}
	}

	return
}

func formatSearchQuery(searchStringsList []string, propertyName string) (searchSQLString string) {

	for index, value := range searchStringsList {

		if index < (len(searchStringsList) - 1) {
			searchStringsList[index] = strconv.Quote(value) + " OR "
		} else {
			searchStringsList[index] = strconv.Quote(value)
		}
	}

	if len(searchStringsList) != 0 {
		searchSQLString = strings.Join(searchStringsList, "")
		searchSQLString = fmt.Sprintf(`SELECT DISTINCT(id) from "%s" WHERE "%s" MATCH '%s'`, searchTableName, propertyName, searchSQLString)
	}

	return
}
