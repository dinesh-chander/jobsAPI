package job

type SearchableContent struct {
	Id          string
	Title       string
	Description string
	Location    string
	Tags        string
}

func NewSearchableContent() *SearchableContent {
	return &SearchableContent{}
}

func AddSearchAbleContent(newJob *Job) {
	searchableContent := &SearchableContent{
		Id:          string(newJob.ID),
		Title:       newJob.Title,
		Description: newJob.Description,
		Location:    newJob.Location,
		Tags:        newJob.Tags,
	}

	db.Create(searchableContent)
}

func SearchContent(searchQuery string) (result [](*SearchableContent)) {
	rows, err := db.Exec("SELECT * from searchable_content where searchable_content match ?", searchQuery).Rows()
	defer rows.Close()

	if err != nil {
		loggerInstance.Println(err)
		return
	}

	result = [](*SearchableContent){}
	for rows.Next() {
		var searchableContent SearchableContent
		db.ScanRows(rows, &searchableContent)
		result = append(result, &searchableContent)
	}

	return
}
