package job

import (
	"config"
	"strconv"
	"time"
	jobType "types/jobs"
)

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

func migrateJobRowsToSearchableContent() {
	loggerInstance.Println("Migrating old data to searchable_content table")

	var indexEntriesOfLastXDays uint64

	if config.GetConfig("indexEntriesOfLastXDays") != "" {
		var err error

		indexEntriesOfLastXDays, err = strconv.ParseUint(config.GetConfig("indexEntriesOfLastXDays"), 10, 64)

		if err != nil {
			loggerInstance.Panicln("invalid value for 'indexEntriesOfLastXDays' in config. Value is :", config.GetConfig("indexEntriesOfLastXDays"))
		}
	} else {
		indexEntriesOfLastXDays = 0 // index everything
	}

	loggerInstance.Println("Indexing Entries for last :", indexEntriesOfLastXDays, "Days")

	if indexEntriesOfLastXDays != 0 {
		indexEntriesOfLastXDays = (uint64(time.Now().Unix()) - (indexEntriesOfLastXDays * 24 * 3600)) * 1000
	}

	jobsList := GetAll()

	tx := db.Begin()
	defer tx.Commit()

	var newSearchableContent *jobType.SearchableContent
	var insertErr error

	for _, job := range jobsList {
		if indexEntriesOfLastXDays != 0 {
			if job.Published_Date > indexEntriesOfLastXDays {
				newSearchableContent = addSearchableContent(&job)
				insertErr = tx.Table(searchTableName).Create(newSearchableContent).Error

				if insertErr != nil {
					loggerInstance.Println(insertErr.Error())
				}
			}
		}
	}

	loggerInstance.Println("Migration Complete")
}
