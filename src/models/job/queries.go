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

	insertErr = db.Create(newJob).Error

	if insertErr != nil {
		tx.Rollback()
		loggerInstance.Println(insertErr.Error())
	} else {
		insertErr = addSearchableContent(newJob)

		if insertErr != nil {
			tx.Rollback()
			loggerInstance.Println(insertErr.Error())
		}
	}
}

func addSearchableContent(newJob *Job) error {
	newSearchableContent := &SearchableContent{
		ID:          newJob.Source_Id,
		Title:       newJob.Title,
		Company:     newJob.Company,
		Description: newJob.Description,
		Location:    newJob.Location,
		Tags:        newJob.Tags,
	}

	return db.Create(newSearchableContent).Error
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

func SearchContent(searchQuery *jobInterface.Query) (searchResult []Job) {
	loggerInstance.Println("Search query is ", searchQuery)
	searchString := buildSearchString(searchQuery)

	searchQuerySQLString := fmt.Sprintf("SELECT DISTINCT(ID) FROM %s WHERE %s MATCH %s ORDER BY rank", searchTableName, searchTableName, searchString)

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

func buildSearchString(searchQuery *jobInterface.Query) (searchString string) {
	for index, value := range searchQuery.Location {
		if index < len(searchQuery.Location)-1 {
			searchQuery.Location[index] = "'Location:" + strconv.Quote(value) + "' OR "
		} else {
			searchQuery.Location[index] = "'Location:" + strconv.Quote(value) + "'"
		}
	}

	searchString = strings.Join(searchQuery.Location, "")

	return
}
