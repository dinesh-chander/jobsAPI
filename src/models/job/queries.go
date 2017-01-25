package job

func AddJob(newJob *Job) {
	db.Create(newJob)
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

func SearchContent(searchQuery string) (result [](*Job)) {
	rows, err := db.Exec("SELECT * from ? where searchable_content match ?", tableName, searchQuery).Rows()

	defer rows.Close()

	if err != nil {
		loggerInstance.Println(err)
		return
	}

	result = [](*Job){}
	for rows.Next() {
		var newJob Job
		db.ScanRows(rows, &newJob)
		result = append(result, &newJob)
	}

	return
}
