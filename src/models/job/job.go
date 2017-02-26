package job

import (
	"config"
	"log"
	"logger"
	jobType "types/jobs"

	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var loggerInstance *log.Logger
var db *gorm.DB
var databaseName string
var tableName string

func NewJob() *jobType.Job {
	return &jobType.Job{}
}

func getDbPath(dbName string) (dbPath string) {

	dbUser := config.GetConfig("db_user")
	dbPassword := config.GetConfig("db_password")

	dbHost := config.GetConfig("db_host")
	dbPort := config.GetConfig("db_port")

	connectionURI := dbUser + ":" + dbPassword + "@" + "tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	return connectionURI
}

func createDatabase() {

	databaseName = config.GetConfig("mode") + "_" + config.GetConfig("db_name")

	databaseCreationErr := db.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName).Error

	if databaseCreationErr != nil {
		loggerInstance.Fatalln(databaseCreationErr.Error())
	}
}

func createDatabaseTables() {

	var tableCreationErr error

	tableName = config.GetConfig("mode") + "_" + "jobs"

	tableCreationErr = db.Table(tableName).AutoMigrate(&jobType.Job{}).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}
}

func createTableIndexes() {

	//indexCreationErr := db.Table(tableName).Exec("CREATE UNIQUE INDEX %s ON %s (%s);", "idx_source_id", tableName, "source_id").Error
	//
	//if indexCreationErr != nil {
	//	loggerInstance.Fatalln(indexCreationErr.Error())
	//}

	go func() {
		sourceIdIndexCreationErr := db.Table(tableName).Exec(fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s);", "idx_source_id", tableName, "source_id")).Error

		if sourceIdIndexCreationErr != nil {
			loggerInstance.Println(sourceIdIndexCreationErr.Error())
		}
	}()

	go func() {
		channelNameIndexCreationErr := db.Table(tableName).Exec(fmt.Sprintf("CREATE INDEX %s ON %s (%s);", "idx_channel_name", tableName, "channel_name")).Error

		if channelNameIndexCreationErr != nil {
			loggerInstance.Println(channelNameIndexCreationErr.Error())
		}
	}()

	go func() {
		titleIndexCreationErr := db.Table(tableName).Exec(fmt.Sprintf("CREATE FULLTEXT INDEX %s ON %s (%s);", "fts_idx_title", tableName, "title")).Error

		if titleIndexCreationErr != nil {
			loggerInstance.Println(titleIndexCreationErr.Error())
		}
	}()

	go func() {
		descriptionIndexCreationErr := db.Table(tableName).Exec(fmt.Sprintf("CREATE FULLTEXT INDEX %s ON %s (%s, %s);", "fts_idx_keyword", tableName, "description", "tags")).Error

		if descriptionIndexCreationErr != nil {
			loggerInstance.Println(descriptionIndexCreationErr.Error())
		}
	}()

	go func() {
		locationIndexCreationErr := db.Table(tableName).Exec(fmt.Sprintf("CREATE FULLTEXT INDEX %s ON %s (%s, %s, %s);", "fts_idx_location", tableName, "address", "city", "country")).Error

		if locationIndexCreationErr != nil {
			loggerInstance.Println(locationIndexCreationErr.Error())
		}
	}()
}

func setConnectionConfiguration() {

	if config.GetConfig("dbQueryLog") == "true" {
		db.LogMode(true)
		db.SetLogger(loggerInstance)
	} else {
		db.LogMode(false)
	}

	db.SingularTable(true)
	db.DB().SetMaxOpenConns(5)
}

func init() {

	var connectionErr error

	loggerInstance = logger.Logger

	db, connectionErr = gorm.Open("mysql", getDbPath(databaseName))

	if connectionErr != nil {
		loggerInstance.Fatalln(connectionErr.Error())
	}

	createDatabase()

	db.Close()

	db, connectionErr = gorm.Open("mysql", getDbPath(databaseName))

	if connectionErr != nil {
		loggerInstance.Fatalln(connectionErr.Error())
	}

	setConnectionConfiguration()

	createDatabaseTables()

	createTableIndexes()
}

func FindContent(searchQuery *jobType.Query) (searchResult []jobType.Job, numberOfAvailableRecords int) {

	searchQuerySQLString := buildSearchString(searchQuery)

	if searchQuery.Limit == 0 {
		searchResult = []jobType.Job{}
		return
	}

	return findJobs(searchQuerySQLString, searchQuery.Limit, searchQuery.Skip)
}
