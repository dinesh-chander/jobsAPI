package job

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var tableName string = "dev_jobs"

type Job struct {
	gorm.Model
	Title         string
	Company       string
	Description   string
	PublishedDate int64
	Compensation  string
	Location      string
	IsRemote      string
	JobFrom       string
	Tags          []string
	Share_Tags    []String
}

func (Job) TableName() string {
	return tableName
}

func New() *Job {
	return &Job{}
}

func AddJob(newJob *Job) {
	fmt.Println("Adding New job in the db")
	db.Create(newJob)
	fmt.Println("Total Job Count :", GetJobsCount())
}

func GetJob() *Job {
	fmt.Println("Returning first job in the db")
	var job Job
	db.First(&job)
	return &job
}

func GetJobsCount() (count int) {
	db.Table(tableName).Count(&count)
	return count
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "job.db")
	if err != nil {
		panic("Unable to open database connection")
	}

	db.AutoMigrate(&Job{})
}
