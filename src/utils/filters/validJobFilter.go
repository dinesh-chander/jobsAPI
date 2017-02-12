package filters

import (
	"config"
	"strings"
)

var searchWordsList []string

func IsValidJob(value string) bool {

	for _, searchWord := range searchWordsList {
		if strings.Contains(strings.ToUpper(value), searchWord) {
			return true
		}
	}

	return false
}

func init() {
	searchWordsList = strings.Split(config.GetConfig("validWords"), ",")

	for index, searchWord := range searchWordsList {
		searchWordsList[index] = strings.ToUpper(searchWord)
	}
}
