package filters

import (
	"config"
	"strings"
)

var searchWordsList []string
var filterWordsList []string

func isWordPresent(value string, wordsList []string) bool {

	if len(wordsList) == 0 {
		return true
	}

	for _, word := range wordsList {

		if strings.Contains(strings.ToUpper(value), word) {
			return true
		}
	}

	return false
}

func IsValidJob(value string) bool {
	return isWordPresent(value, searchWordsList) && isWordPresent(value, filterWordsList)
}

func init() {

	filterWordsList = strings.Split(config.GetConfig("filterWords"), ",")

	searchWordsList = strings.Split(config.GetConfig("searchWords"), ",")

	convertListToUpperCase(searchWordsList)

	convertListToUpperCase(filterWordsList)
}

func convertListToUpperCase(wordsList []string) {

	for index, word := range wordsList {
		wordsList[index] = strings.ToUpper(word)
	}
}
