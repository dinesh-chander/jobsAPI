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

	insertErr = tx.Create(newJob).Error

	if insertErr != nil {
		loggerInstance.Println(insertErr.Error())
	} else {
		newSearchableContent := addSearchableContent(newJob)

		insertErr = tx.Create(newSearchableContent).Error

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

func findFromNormalTable(searchQuerySQLString string) (searchResult []Job) {

	tx := db.Begin()
	rows, err := tx.Raw(searchQuerySQLString).Rows()

	defer rows.Close()
	defer tx.Commit()

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	searchResult = []Job{}

	for rows.Next() {
		var newJob Job
		scanErr := tx.ScanRows(rows, &newJob)

		if scanErr != nil {
			loggerInstance.Println(scanErr.Error())
		} else {
			searchResult = append(searchResult, newJob)
		}
	}

	return
}

func findFromSearchableTable(searchQuerySQLString string) (searchResult []Job) {

	loggerInstance.Println(searchQuerySQLString)

	tx := db.Begin()
	rows, err := tx.Raw(searchQuerySQLString).Rows()

	defer rows.Close()
	defer tx.Commit()

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	var ID string
	searchResult = []Job{}

	for rows.Next() {
		scanErr := rows.Scan(&ID)

		if scanErr != nil {
			loggerInstance.Println(scanErr.Error())
		} else {
			var newJob Job
			tx.First(&newJob, "Source_Id = ?", ID)
			searchResult = append(searchResult, newJob)
		}
	}

	return
}

func FindContent(searchQuery *jobInterface.Query) (searchResult []Job) {

	var searchQuerySQLString string
	searchString := buildSearchString(searchQuery)

	if searchString != "" {
		searchQuerySQLString = fmt.Sprintf("SELECT DISTINCT(ID) FROM %s WHERE %s MATCH %s ORDER BY rank", searchTableName, searchTableName, searchString)
		return findFromSearchableTable(searchQuerySQLString)
	} else {
		searchQuerySQLString = fmt.Sprintf("SELECT * FROM %s", tableName)
		return findFromNormalTable(searchQuerySQLString)
	}
}

func buildSearchString(searchQuery *jobInterface.Query) (searchString string) {

	queryStringList := [5]string{}

	queryStringList[0] = formatSearchQuery(searchQuery.Location, "Location")
	queryStringList[1] = formatSearchQuery(searchQuery.Company, "Company")
	queryStringList[2] = formatSearchQuery(searchQuery.Tags, "Tags")
	queryStringList[3] = formatSearchQuery(searchQuery.Title, "Title")

	for _, value := range queryStringList {
		if value != "" {
			if searchString != "" {
				searchString = searchString + " AND " + value
			} else {
				searchString = value
			}
		}
	}

	if searchString != "" {
		searchString = "'" + searchString + "'"
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
