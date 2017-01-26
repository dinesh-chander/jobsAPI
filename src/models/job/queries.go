package job

import (
	"bytes"
	"compress/gzip"
	"fmt"
	jobInterface "interfaces/jobs"
	"strconv"
	"strings"
)

func AddJob(newJob *Job) {
	var compressedDescription bytes.Buffer
	compressedJob := *newJob

	gz, gzErr := gzip.NewWriterLevel(&compressedDescription, 9)

	if gzErr != nil {
		loggerInstance.Println(gzErr.Error())
	} else {
		var writeErr error

		_, writeErr = gz.Write([]byte(compressedJob.Description))

		if writeErr != nil {
			loggerInstance.Println(writeErr.Error())
		} else if writeErr = gz.Flush(); writeErr != nil {
			loggerInstance.Println(writeErr.Error())
		} else if writeErr = gz.Close(); writeErr != nil {
			loggerInstance.Println(writeErr.Error())
		} else {
			compressedJob.Description = compressedDescription.String()

			loggerInstance.Println(len(compressedJob.Description), len(newJob.Description))

			db.Create(&compressedJob)
			addSearchableContent(newJob)
		}
	}
}

func addSearchableContent(newJob *Job) {
	newSearchableContent := &SearchableContent{
		Title:       newJob.Title,
		Company:     newJob.Company,
		Description: newJob.Description,
		Location:    newJob.Location,
		Tags:        newJob.Tags,
	}

	db.Create(newSearchableContent)
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

func SearchContent(searchQuery *jobInterface.Query) (result []SearchableContent) {
	loggerInstance.Println("Search query is ", searchQuery)
	searchString := buildSearchString(searchQuery)

	searchQuerySQLString := fmt.Sprintf("SELECT * FROM %s WHERE %s MATCH %s ORDER BY rank", searchTableName, searchTableName, searchString)

	loggerInstance.Println(searchQuerySQLString)

	rows, err := db.Raw(searchQuerySQLString).Rows()

	defer rows.Close()

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	result = []SearchableContent{}

	for rows.Next() {
		var newSearchableContent SearchableContent
		db.ScanRows(rows, &newSearchableContent)
		loggerInstance.Println(newSearchableContent.Location)
		result = append(result, newSearchableContent)
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
