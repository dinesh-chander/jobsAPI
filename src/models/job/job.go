package job

import (
	"config"
	"fmt"
	"log"
	"logger"
	"os"
	"path"
	jobType "types/jobs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var loggerInstance *log.Logger
var db *gorm.DB
var tableName string
var searchTableName string

func NewJob() *jobType.Job {
	return &jobType.Job{}
}

func getDbPath() (dbPath string) {
	var logsDirectory string

	ex, wrongVersionErr := os.Executable()

	if wrongVersionErr != nil {
		panic(wrongVersionErr.Error())
	}

	exPath := path.Dir(ex)

	logsDirectory = path.Join(path.Dir(exPath), config.GetConfig("dbDir"))

	_, dirErr := os.Stat(logsDirectory)

	if dirErr != nil {
		if os.IsNotExist(dirErr) {

			mkdirErr := os.Mkdir(logsDirectory, 0700)

			if mkdirErr != nil {
				panic(mkdirErr.Error())
			}
		} else {
			panic(dirErr.Error())
		}
	}

	return path.Join(logsDirectory, "job.db")
}

func init() {
	var databaseCreationErr error
	var tableCreationErr error

	tableName = config.GetConfig("tableNamePrefix") + "jobs"
	searchTableName = config.GetConfig("tableNamePrefix") + "searchable_content"
	loggerInstance = logger.Logger

	db, databaseCreationErr = gorm.Open("sqlite3", getDbPath())

	if config.GetConfig("dbQueryLog") == "true" {

		db.LogMode(true)
		db.SetLogger(loggerInstance)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxOpenConns(5)

	db.Exec(`
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

	tableCreationErr = db.Table(tableName).AutoMigrate(&jobType.Job{}).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	indexCreationErr := db.Table(tableName).AddUniqueIndex("idx_source_id", "source_id").Error

	if indexCreationErr != nil {
		loggerInstance.Fatalln(indexCreationErr.Error())
	}

	createSearchTable()
	go removeOlderEntries()
}

func createSearchTable() {
	err := db.DropTableIfExists(searchTableName).Error

	if err != nil {
		loggerInstance.Fatalln(err.Error())
	}

	tableCreationErr := db.Exec(fmt.Sprintf("CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts4(ID, Title, Description, Location, matchinfo=fts3, tokenize=porter 'remove_diacritics=1', notindexed=ID)", searchTableName)).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}

	go migrateJobRowsToSearchableContent()
}
