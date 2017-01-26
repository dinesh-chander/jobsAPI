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

type SearchableContent struct {
	Title       string
	Company     string
	Description string
	Location    string
	Tags        string
}

var loggerInstance *log.Logger
var db *gorm.DB
var tableName string
var searchTableName string

func (Job) TableName() string {
	return tableName
}

func NewJob() *Job {
	return &Job{}
}

func init() {
	var databaseCreationErr error
	var tableCreationErr error

	tableName = "dev_jobs"
	searchTableName = "searchable_content"
	loggerInstance = logger.Logger

	db, databaseCreationErr = gorm.Open("sqlite3", "job.db")

	db.Exec("PRAGMA auto_vacuum  = INCREMENTAL")
	db.Exec("PRAGMA cache_size   = 10000")
	db.Exec("PRAGMA synchronous  = OFF")

	db.SingularTable(true)

	if databaseCreationErr != nil {
		loggerInstance.Fatalln(databaseCreationErr.Error())
	}

	tableCreationErr = db.AutoMigrate(&Job{}).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	createSearchTable()
}

func createSearchTable() {
	err := db.DropTableIfExists(searchTableName).Error

	if err != nil {
		loggerInstance.Fatalln(err.Error())
	}

	tableCreationErr := db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts5(Title, Company, Description, Location, Tags)", searchTableName)).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	migrateJobRowsToSearchableContent()
}

func migrateJobRowsToSearchableContent() {
	loggerInstance.Println("Migrating old data to searchable_content table")

	jobsList := GetAll()

	db.Exec("BEGIN TRANSACTION")

	for _, job := range jobsList {
		addSearchableContent(&job)
	}

	db.Exec("COMMIT TRANSACTION")
	loggerInstance.Println("Migration Complete")
}
