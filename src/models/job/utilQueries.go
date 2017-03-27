package job

import "types/jobs"

func FindLastAddedEntryTimestampForChannel(channelName string) (lastPublishedAt int64) {

	tx := db.Begin()

	defer tx.Commit()

	row := tx.Table(tableName).Select("max(published_date)").Where("channel_name == ?", channelName).Row()

	row.Scan(&lastPublishedAt)

	return
}

func FindAlreadyPresentJobsWithGivenSourceIds(channelName string, sourceIdsList []string) (unavailableJobsSourceIds []string) {

	tx := db.Begin()

	defer tx.Commit()

	alreadyPresentJobs := []jobs.Job{}

	tx.Table(tableName).Select("source_id").Where("channel_name = ? AND source_id IN (?)", channelName, sourceIdsList).Find(&alreadyPresentJobs)

	unavailableJobsSourceIds = make([]string, len(sourceIdsList))

	copy(unavailableJobsSourceIds, sourceIdsList)

	if len(alreadyPresentJobs) > 0 {

		var jobListId string

		for _, job := range alreadyPresentJobs {

			for index := 0; index < len(unavailableJobsSourceIds); {

				jobListId = unavailableJobsSourceIds[index]

				if job.Source_Id == jobListId {
					unavailableJobsSourceIds[index] = unavailableJobsSourceIds[len(unavailableJobsSourceIds)-1]
					unavailableJobsSourceIds = unavailableJobsSourceIds[:len(unavailableJobsSourceIds)-1]
				} else {
					index = index + 1
				}
			}
		}
	}

	return
}
