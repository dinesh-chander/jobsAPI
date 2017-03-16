package job

import (
	jobType "types/jobs"
)

func AddJob(newJob *jobType.Job) {
	tx := db.Begin()
	defer tx.Commit()

	insertErr := tx.Table(tableName).Create(newJob).Error

	if insertErr != nil {
		loggerInstance.Println(insertErr.Error())
	} else {
		loggerInstance.Println("new job added from :", newJob.Channel_Name)
	}
}

func findJobs(searchQuery string, resultListLength int, offset int) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchResult = []jobType.Job{}

	tx := db.Begin()

	defer tx.Commit()

	err := tx.Table(tableName).Limit(resultListLength).Offset(offset).Order("published_date DESC").Where("approved = 1").Find(&searchResult, searchQuery).Error

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	tx.Table(tableName).Where("approved = 1").Where(searchQuery).Count(&numberOfAvailableRecords)

	return
}
