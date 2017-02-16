package job

import jobType "types/jobs"

func AddJob(newJob *jobType.Job) {

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

func GetAll() (jobsList []jobType.Job) {

	jobsList = []jobType.Job{}
	db.Table(tableName).Find(&jobsList)

	return
}

func findFromNormalTable(resultListLength int, offset int) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchResult = []jobType.Job{}

	tx := db.Begin()

	defer tx.Commit()

	searchResult = make([]jobType.Job, resultListLength)

	err := tx.Table(tableName).Limit(resultListLength).Offset(offset).Order("published_date DESC").Find(&searchResult).Error

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	tx.Table(tableName).Count(&numberOfAvailableRecords)

	return
}
