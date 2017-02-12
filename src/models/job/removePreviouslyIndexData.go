package job

import (
	"config"
	"strconv"
	"time"
	jobType "types/jobs"
	"utils/cronParser"
)

func removeOlderEntries() {
	var indexEntriesOfLastXDays int64
	var removeOlderEntriesThen int64

	if config.GetConfig("indexEntriesOfLastXDays") != "" {
		var err error

		indexEntriesOfLastXDays, err = strconv.ParseInt(config.GetConfig("indexEntriesOfLastXDays"), 10, 64)

		if err != nil {
			loggerInstance.Panicln("invalid value for 'indexEntriesOfLastXDays' in config. Value is:", config.GetConfig("indexEntriesOfLastXDays"))
		}
	} else {
		indexEntriesOfLastXDays = 0 // index everything
	}

	loggerInstance.Println("Indexing Entries for last:", indexEntriesOfLastXDays, "Days")

	scheduleAt := config.GetConfig("removeOlderIndexes")

	expr := cronParser.Parse(scheduleAt)

	if indexEntriesOfLastXDays != 0 {
		for {
			nextTime := expr.Next(time.Now())
			time.Sleep(time.Duration(nextTime.Unix()-time.Now().Unix()) * time.Second)
			removeOlderEntriesThen = (time.Now().Unix() - (indexEntriesOfLastXDays * 24 * 3600)) * 1000
			go makeDeleteFromDBOperation(removeOlderEntriesThen)
		}
	}
}

func makeDeleteFromDBOperation(lastValidTimestamp int64) {
	tx := db.Begin()
	rows, selectErr := tx.Table(tableName).Select("source_id, published_date").Where("published_date < ?", lastValidTimestamp).Rows()

	defer rows.Close()
	defer tx.Commit()

	if selectErr != nil {
		loggerInstance.Println(selectErr.Error())
	} else {
		for rows.Next() {
			var newJob jobType.Job
			scanErr := tx.ScanRows(rows, &newJob)

			if scanErr != nil {
				loggerInstance.Println(scanErr.Error())
			} else {
				err := tx.Table(searchTableName).Unscoped().Where("id = ?", newJob.Source_Id).Delete(&jobType.SearchableContent{}).Error
				if err != nil {
					loggerInstance.Println(err.Error())
				}
			}
		}

		loggerInstance.Println("Indexes older than", config.GetConfig("indexEntriesOfLastXDays"), "Days are successfully removed")
	}
}
