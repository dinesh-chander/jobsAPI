package main

func getConfig() (config map[string]string) {
    config = make(map[string]string)

	config["whoishiring"] = "60 0 0 * * * *"

	return config
}
