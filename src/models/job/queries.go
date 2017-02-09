package job

import (
	"fmt"
	jobInterface "interfaces/jobs"
	"strconv"
	"strings"
)

func AddJob(newJob *Job) {

	var insertErr error

	tx := db.Begin()

	defer tx.Commit()

	insertErr = tx.Table(tableName).Create(newJob).Error

	if insertErr != nil {
		loggerInstance.Println(insertErr.Error())
	} else {

		newSearchableContent := addSearchableContent(newJob)
		insertErr = tx.Table(searchTableName).Create(newSearchableContent).Error

		if insertErr != nil {
			loggerInstance.Println(insertErr.Error())
		}
	}
}

func addSearchableContent(newJob *Job) (newSearchableContent *SearchableContent) {

	location := newJob.Address

	if newJob.City != "" {
		location = location + " " + newJob.City
	}

	if newJob.Country != "" {
		location = location + " " + newJob.Country
	}

	newSearchableContent = &SearchableContent{
		ID:          newJob.Source_Id,
		Title:       newJob.Title,
		Company:     newJob.Company,
		Description: newJob.Description,
		Location:    location,
		Tags:        newJob.Tags,
	}

	return
}

func GetAll() (jobsList []Job) {

	jobsList = []Job{}
	db.Find(&jobsList)
	return
}

func GetJob() *Job {

	var job Job
	db.First(&job)
	return &job
}

func GetJobsCount() (count int) {

	db.Table(tableName).Count(&count)
	return count
}

func FindLastAddedEntryTimestampForSource(channelName string) (lastPublishedAt int64) {

	tx := db.Begin()
	defer tx.Commit()

	row := tx.Table(tableName).Select("max(published_date)").Where("channel_name == ?", channelName).Row()

	row.Scan(&lastPublishedAt)

	return
}

func findFromNormalTable(searchCondition string, resultListLength int, offset int) (searchResult []Job) {

	searchResult = []Job{}

	tx := db.Begin()
	defer tx.Commit()

	searchResult = make([]Job, resultListLength)

	err := tx.Table(tableName).Limit(resultListLength).Offset(offset).Find(&searchResult, searchCondition).Error

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	return
}

func findFromSearchableTable(searchSQLQuery string, resultListLength int, offset int) (searchResult []Job) {

	searchResult = []Job{}

	tx := db.Begin()

	rows, err := tx.Table(searchTableName).Limit(resultListLength).Offset(offset).Order("rank DESC").Raw(searchSQLQuery).Rows()

	defer rows.Close()
	defer tx.Commit()

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	var id string
	var index int = 0

	matchedIDs := make([]string, resultListLength)

	for rows.Next() {

		scanErr := rows.Scan(&id)

		if scanErr != nil {
			loggerInstance.Println(scanErr.Error())
		} else {
			if id != "" {
				matchedIDs[index] = "'" + id + "'"
				index = index + 1
			}
		}
	}

	if index > 0 {
		var separator string
		matchedIDs = matchedIDs[:index]

		if index > 1 {
			separator = ","
		}

		idsList := "(" + strings.Join(matchedIDs, separator) + ")"

		searchResult = make([]Job, index)

		fetchErr := tx.Table(tableName).Find(&searchResult, ("source_id in " + idsList)).Error

		if fetchErr != nil {
			loggerInstance.Println(fetchErr.Error())
			return
		}
	}

	return
}

func FindContent(searchQuery *jobInterface.Query) (searchResult []Job) {

	var searchQuerySQLString string
	searchString := buildSearchString(searchQuery)

	if searchQuery.Limit == 0 {
		searchResult = []Job{}
		return
	}

	if searchString != "" {
		searchQuerySQLString = fmt.Sprintf(`SELECT DISTINCT(id) from "%s" WHERE "%s" MATCH '%s'`, searchTableName, searchTableName, searchString)
		return findFromSearchableTable(searchQuerySQLString, searchQuery.Limit, searchQuery.Skip)
	} else {
		return findFromNormalTable("", searchQuery.Limit, searchQuery.Skip)
	}
}

func buildSearchString(searchQuery *jobInterface.Query) (searchString string) {

	queryStringList := [5]string{}

	queryStringList[0] = formatSearchQuery(searchQuery.Locations, "Location")
	queryStringList[1] = formatSearchQuery(searchQuery.Companies, "Company")
	queryStringList[2] = formatSearchQuery(searchQuery.Tags, "Tags")
	queryStringList[3] = formatSearchQuery(searchQuery.Titles, "Title")
	queryStringList[4] = formatSearchQuery(searchQuery.Keywords, "Description")

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

func formatSearchQuery(searchStringsList []string, propertyName string) string {

	for index, value := range searchStringsList {
		if index < (len(searchStringsList) - 1) {
			searchStringsList[index] = propertyName + ":" + strconv.Quote(value) + " OR "
		} else {
			searchStringsList[index] = propertyName + ":" + strconv.Quote(value)
		}
	}

	return strings.Join(searchStringsList, "")
}
