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

	config["mode"] = "production"

	//config["mode"] = "development"

	if config["mode"] == "development" {

		config["db_name"] = "jobs"
		config["db_user"] = "root"
		config["db_password"] = ""
		config["db_host"] = "127.0.0.1"
		config["db_port"] = "3306"

		config["betalist"] = "0 0-59/1 0 * * * *"

		config["whoishiring"] = "0 0 0-23/12 * * * *"

		config["angellist"] = "0 0 0 * * 0 *"

		config["interface"] = "localhost"
		config["port"] = "9080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "2"

		config["fetchFrom"] = "30" // default 0 means fetch today's data

		config["dbDir"] = "db"
		config["dbQueryLog"] = "false"
		config["disableLog"] = "false"
		//config["logsDir"] = "logs"
		config["logsFile"] = "logs.log"

		config["proxyURL"] = ""

		config["googleGeoAPIKey"] = "AIzaSyA87A0cCVeQR1yCbeLjitQlWRzg1hYqQyw"

		config["searchWords"] = "junior"
		config["filterWords"] = ""

	} else if config["mode"] == "production" {

		config["db_name"] = "jobs"
		config["db_user"] = "root"
		config["db_password"] = "firstjob@123"
		config["db_host"] = "127.0.0.1"
		config["db_port"] = "4079"

		config["betalist"] = "0 0 0 * * * *"

		config["whoishiring"] = "0 0 0 * * * *"

		config["angellist"] = "0 0 0 * * 0-6/2 *"

		config["interface"] = "localhost"
		config["port"] = "8080"
		config["gzip"] = "true"
		config["jobManagersCount"] = "4"

		config["proxyURL"] = ""

		config["fetchFrom"] = "120" // default 0 means fetch today's data

		config["dbDir"] = "db"
		config["dbQueryLog"] = "false"
		config["disableLog"] = "false"
		config["logsDir"] = "logs"
		config["logsFile"] = "logs.log"

		config["googleGeoAPIKey"] = "AIzaSyDunhBDEvzMh1Zijn3fcMVzmegDCCa9L1E"
		config["searchWords"] = "junior"
		config["filterWords"] = ""
	}
}
