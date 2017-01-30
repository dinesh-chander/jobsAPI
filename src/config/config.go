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

	config["translateToEnglish"] = "false"

	config["indexEntriesOfLastXDays"] = "30" // default 0 means index all

	config["disableLog"] = "false"
	config["removeOlderIndexes"] = "0-59/30 * * * * *"
}
