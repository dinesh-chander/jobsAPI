package job

import (
	"strings"
	jobType "types/jobs"
)

func AddJob(newJob *jobType.Job) {
	tx := db.Begin()
	defer tx.Commit()

	insertErr := tx.Table(tableName).Create(newJob).Error

	if insertErr != nil {
		loggerInstance.Println(insertErr.Error())
	}
}

func findFromNormalTable(resultListLength int, offset int) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchResult = []jobType.Job{}

	tx := db.Begin()

	defer tx.Commit()

	searchResult = make([]jobType.Job, resultListLength)

	err := tx.Table(tableName).Limit(resultListLength).Offset(offset).Order("published_date DESC").Find(&searchResult, "approved", 1).Error

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	tx.Table(tableName).Where("approved", 1).Count(&numberOfAvailableRecords)

	return
}

func findFromSearchableTable(searchSQLQuery string, resultListLength int, offset int) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchResult = []jobType.Job{}

	tx := db.Begin()

	defer tx.Commit()

	rows, err := tx.Table(tableName).Limit(resultListLength).Offset(offset).Raw(searchSQLQuery).Rows()

	if err != nil {
		loggerInstance.Println(err.Error())
		return
	}

	if rows.Err() != nil {
		loggerInstance.Println(rows.Err().Error())
		return
	}

	defer rows.Close()

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

		searchResult = make([]jobType.Job, index)

		query := "source_id in " + idsList

		fetchErr := tx.Table(tableName).Order("published_date DESC").Find(&searchResult, query).Error

		if fetchErr != nil {
			loggerInstance.Println(fetchErr.Error())
			return
		}
	}

	totalRows, countErr := tx.Table(tableName).Raw(searchSQLQuery).Rows()

	if countErr != nil {
		loggerInstance.Println(countErr.Error())
		return
	}

	if totalRows.Err() != nil {
		loggerInstance.Println(totalRows.Err().Error())
		return
	}

	numberOfAvailableRecords = 0

	for totalRows.Next() {
		numberOfAvailableRecords = numberOfAvailableRecords + 1
	}

	return
}
