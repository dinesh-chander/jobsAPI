package job

import (
	"config"
	"log"
	"logger"
	jobType "types/jobs"

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

	tableCreationErr = db.Table(tableName).AutoMigrate(&jobType.Job{}).Error

	if tableCreationErr != nil {
		loggerInstance.Fatalln(tableCreationErr.Error())
	}
}

func createDatabaseTableIndexes() {
	indexCreationErr := db.Table(tableName).RemoveIndex("idx_source_id").AddUniqueIndex("idx_source_id", "source_id").Error

	if indexCreationErr != nil {
		loggerInstance.Fatalln(indexCreationErr.Error())
	}
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

	tableName = config.GetConfig("mode") + "_" + "jobs"

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

	createDatabaseTableIndexes()
}
