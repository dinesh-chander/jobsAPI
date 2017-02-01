package job

import (
	"config"
	"fmt"
	"log"
	"logger"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Job struct {
	gorm.Model
	Title          string
	Address        string
	Channel_Name   string
	City           string
	Company        string
	Country        string
	Description    string
	Published_Date uint64
	Compensation   string
	Is_Remote      bool
	Source         string
	Source_Name    string
	Source_Id      string `gorm:"unique_index"`
	Tags           string
	Approved       bool
}

type SearchableContent struct {
	ID          string
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
	db.LogMode(false)

	db.Exec(`
        PRAGMA automatic_index = true;
        PRAGMA synchronous = false;
	    PRAGMA cache_size = 32768;
        PRAGMA cache_spill = false;
	    PRAGMA read_uncommitted = true;
	    PRAGMA parser_trace = false;
	    PRAGMA journal_mode = MEMORY;
	    PRAGMA foreign_keys = false;`)

	db.SingularTable(true)

	if databaseCreationErr != nil {
		loggerInstance.Fatalln(databaseCreationErr.Error())
	}

	tableCreationErr = db.AutoMigrate(&Job{}).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	createSearchTable()
	go removeOlderEntries()
}

func createSearchTable() {
	err := db.DropTableIfExists(searchTableName).Error

	if err != nil {
		loggerInstance.Fatalln(err.Error())
	}

	tableCreationErr := db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts5(ID, Title, Company, Description, Location, Tags, tokenize = 'porter unicode61 remove_diacritics 1')", searchTableName)).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	go migrateJobRowsToSearchableContent()
}

func migrateJobRowsToSearchableContent() {
	loggerInstance.Println("Migrating old data to searchable_content table")

	var indexEntriesOfLastXDays uint64

	if config.GetConfig("indexEntriesOfLastXDays") != "" {
		var err error

		indexEntriesOfLastXDays, err = strconv.ParseUint(config.GetConfig("indexEntriesOfLastXDays"), 10, 64)

		if err != nil {
			loggerInstance.Panicln("invalid value for 'indexEntriesOfLastXDays' in config. Value is :", config.GetConfig("indexEntriesOfLastXDays"))
		}
	} else {
		indexEntriesOfLastXDays = 0 // index everything
	}

	loggerInstance.Println("Indexing Entries for last :", indexEntriesOfLastXDays, "Days")

	if indexEntriesOfLastXDays != 0 {
		indexEntriesOfLastXDays = (uint64(time.Now().Unix()) - (indexEntriesOfLastXDays * 24 * 3600)) * 1000
	}

	jobsList := GetAll()

	tx := db.Begin()
	defer tx.Commit()

	var newSearchableContent *SearchableContent
	var insertErr error

	for _, job := range jobsList {
		if indexEntriesOfLastXDays != 0 {
			if job.Published_Date > indexEntriesOfLastXDays {
				newSearchableContent = addSearchableContent(&job)
				insertErr = tx.Create(newSearchableContent).Error

				if insertErr != nil {
					loggerInstance.Println(insertErr.Error())
				}
			}
		}
	}

	loggerInstance.Println("Migration Complete")
}
