package job

import (
	"log"
	"logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var loggerInstance *log.Logger
var db *gorm.DB
var tableName string

func init() {
	var err error
	tableName = "dev_jobs"
	loggerInstance = logger.Logger
	db, err = gorm.Open("sqlite3", "job.db")
	db.SingularTable(true)

	if err != nil {
		panic("Unable to open database connection")
	}

	loggerInstance.Println("Trying to create table :", tableName)

	err = db.AutoMigrate(&Job{}).Error

	if err != nil {
		panic("Unable to create table")
	}

	loggerInstance.Println("Trying to create table : searchable_content")

	err = db.Exec("CREATE VIRTUAL TABLE IF NOT EXISTS searchable_content USING fts5(Id, Title, Description, Location, Tags)").Error

	if err != nil {
		panic("Unable to create table")
	}
}
