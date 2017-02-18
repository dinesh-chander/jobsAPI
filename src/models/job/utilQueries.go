package job

func FindLastAddedEntryTimestampForChannel(channelName string) (lastPublishedAt int64) {

	tx := db.Begin()
	defer tx.Commit()

	row := tx.Table(tableName).Select("max(published_date)").Where("channel_name == ?", channelName).Row()

	row.Scan(&lastPublishedAt)

	return
}
