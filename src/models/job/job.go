package job

import (
	"fmt"
	"log"
	"logger"

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

var loggerInstance *log.Logger
var db *gorm.DB
var tableName string

func (Job) TableName() string {
	return tableName
}

func NewJob() *Job {
	return &Job{}
}

func init() {
	var databaseCreationErr error
	tableName = "dev_jobs"
	loggerInstance = logger.Logger
	db, databaseCreationErr = gorm.Open("sqlite3", "job.db")
	db.SingularTable(true)

	if databaseCreationErr != nil {
		loggerInstance.Fatalln(databaseCreationErr.Error())
	}

	tableCreationErr := db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts5(Id, Title, Description, Location, Tags)", tableName)).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}
}
