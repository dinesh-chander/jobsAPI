package config

var (
	config = make(map[string]string)
)

func GetConfig(configProperty string) (value string) {
	return config[configProperty]
}

func init() {

	config["whoishiring"] = "0-59/30 * * * * *"

	config["interface"] = "localhost"
	config["port"] = "9080"
	config["gzip"] = "true"

	config["fetchFrom"] = "10" // default 0 means fetch today's data
	config["translateToEnglish"] = "false"

	config["indexEntriesOfLastXDays"] = "30" // default 0 means index all

	config["disableLog"] = "true"
	config["removeOlderIndexes"] = "0-59/30 * * * * *"
}
