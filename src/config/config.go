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

		config["whoishiring"] = "0 0 0 1-31/1 * * *"

		config["angellist"] = "0 0 0 * * 0 *"

		config["interface"] = "localhost"
		config["port"] = "9080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "3"

		config["tableNamePrefix"] = "dev_"

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

		config["whoishiring"] = "0-23/12 * * * *"

		config["angellist"] = "0 0 0 * * 0-6/3 *"

		config["interface"] = "localhost"
		config["port"] = "8080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "3"

		config["tableNamePrefix"] = ""

		config["fetchFrom"] = "120" // default 0 means fetch today's data

		config["indexEntriesOfLastXDays"] = "90" // default 0 means index all

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
