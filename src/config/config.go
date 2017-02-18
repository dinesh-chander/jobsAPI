package config

//
//import (
//	viper "github.com/spf13/viper"
//)

var config map[string]string

func GetConfig(configProperty string) (value string) {
	return config[configProperty]
}

func init() {
	config = make(map[string]string)

	config["mode"] = "development"

	if config["mode"] == "development" {

		config["db_name"] = "jobs"
		config["db_user"] = "root"
		config["db_password"] = ""
		config["db_host"] = "127.0.0.1"
		config["db_port"] = "3306"

		config["whoishiring"] = "0 0 0-23/12 * * * *"

		config["angellist"] = "0 0 0 * * 0 *"

		config["interface"] = "localhost"
		config["port"] = "9080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "3"

		config["fetchFrom"] = "30" // default 0 means fetch today's data

		config["indexEntriesOfLastXDays"] = "30" // default 0 means index all

		config["dbDir"] = "db"
		config["dbQueryLog"] = "false"
		config["disableLog"] = "false"
		//config["logsDir"] = "logs"
		config["logsFile"] = "logs.log"

		config["removeOlderIndexes"] = "0-59/30 * * * * *"

		config["googleGeoAPIKey"] = "AIzaSyA87A0cCVeQR1yCbeLjitQlWRzg1hYqQyw"

		config["searchWords"] = "junior"
		config["filterWords"] = ""

	} else if config["mode"] == "production" {

		config["db_name"] = "jobs"
		config["db_user"] = "root"
		config["db_password"] = "mysql"
		config["db_host"] = "localhost"
		config["db_port"] = "3306"

		config["whoishiring"] = "0 0 0 * * 0-6/1 *"

		config["angellist"] = "0 0 0 * * 0-6/5 *"

		config["interface"] = "localhost"
		config["port"] = "8080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "5"

		config["fetchFrom"] = "120" // default 0 means fetch today's data

		config["indexEntriesOfLastXDays"] = "120" // default 0 means index all

		config["dbDir"] = "db"
		config["dbQueryLog"] = "false"
		config["disableLog"] = "false"
		config["logsDir"] = "logs"
		config["logsFile"] = "logs.log"

		config["removeOlderIndexes"] = "0-23/12 * * * *"

		config["googleGeoAPIKey"] = "AIzaSyDunhBDEvzMh1Zijn3fcMVzmegDCCa9L1E"
		config["searchWords"] = "junior"
		config["filterWords"] = ""
	}
}
