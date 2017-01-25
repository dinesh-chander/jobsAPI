package job

type searchableContent struct {
	Id          string
	Title       string
	Description string
	Location    string
	Tags        string
}

func NewSearchableContent() *searchableContent {
	return &searchableContent{}
}

func AddSearchAbleContent(newJob *Job) {
	searchableContent := &searchableContent{
		Id:          string(newJob.ID),
		Title:       newJob.Title,
		Description: newJob.Description,
		Location:    newJob.Location,
		Tags:        newJob.Tags,
	}

	db.Create(searchableContent)
}

func SearchContent(searchQuery string) (result [](*searchableContent)) {
	rows, err := db.Exec("SELECT * from searchable_content where searchable_content match ?", searchQuery).Rows()
	defer rows.Close()

	if err != nil {
		loggerInstance.Println(err)
		return
	}

	result = [](*searchableContent){}
	for rows.Next() {
		var searchableContent searchableContent
		db.ScanRows(rows, &searchableContent)
		result = append(result, &searchableContent)
	}

	return
}
