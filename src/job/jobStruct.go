package job

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Job struct {
	gorm.Model
	Title         string
	Company       string
	Description   string
	PublishedDate int
	Compensation  string
	Location      string
	IsRemote      bool
	Source        string
	Tags          string
	Approved      bool
}

func (Job) TableName() string {
	return tableName
}

func NewJob() *Job {
	return &Job{}
}

func AddJob(newJob *Job) {
	db.Create(newJob)
	AddSearchAbleContent(newJob)
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
