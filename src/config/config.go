package config

func GetConfig() (config map[string]string) {
	config = make(map[string]string)

	config["whoishiring"] = "0-59/2 * * * * *"

	config["interface"] = "localhost"
	config["port"] = "9080"

	return config
}
