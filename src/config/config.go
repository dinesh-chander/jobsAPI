package config

var (
	config = make(map[string]string)
)

func GetConfig(configProperty string) (value string) {
	return config[configProperty]
}

func init() {

	config["mode"] = "production"

	if config["mode"] == "development" {

		config["whoishiring"] = "0-59/30 * * * * *"

		config["interface"] = "localhost"
		config["port"] = "9080"
		config["gzip"] = "false"

		config["tableNamePrefix"] = "dev_"

		config["fetchFrom"] = "10" // default 0 means fetch today's data
		config["translateToEnglish"] = "false"

		config["indexEntriesOfLastXDays"] = "30" // default 0 means index all

		config["dbQueryLog"] = "true"
		config["disableLog"] = "true"
		config["removeOlderIndexes"] = "0-59/30 * * * * *"

		config["googleGeoAPIKey"] = "AIzaSyDunhBDEvzMh1Zijn3fcMVzmegDCCa9L1E"
		config["validWords"] = " junior , entry "

	} else if config["mode"] == "production" {

		config["whoishiring"] = "0-23/12 * * * *"

		config["interface"] = "localhost"
		config["port"] = "8080"
		config["gzip"] = "true"

		config["tableNamePrefix"] = ""

		config["fetchFrom"] = "60" // default 0 means fetch today's data
		config["translateToEnglish"] = "false"

		config["indexEntriesOfLastXDays"] = "60" // default 0 means index all

		config["dbQueryLog"] = "false"
		config["disableLog"] = "false"
		config["removeOlderIndexes"] = "0-23/12 * * * *"

		config["googleGeoAPIKey"] = "AIzaSyDunhBDEvzMh1Zijn3fcMVzmegDCCa9L1E"
		config["validWords"] = " junior , entry "
	}
}
