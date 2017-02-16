package logger

import (
	"config"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var Logger *log.Logger

func init() {

	var logsPrefix string = config.GetConfig("mode") + " : "

	if config.GetConfig("disableLog") == "true" {
		Logger = log.New(ioutil.Discard, logsPrefix, log.Lshortfile)
	} else {
		if config.GetConfig("logsDir") == "" || config.GetConfig("logsFile") == "" {
			Logger = log.New(os.Stdout, logsPrefix, log.Lshortfile)
		} else {

			var logsFilePath string
			var logsDirectory string

			ex, wrongVersionErr := os.Executable()

			if wrongVersionErr != nil {
				panic(wrongVersionErr.Error())
			}

			exPath := path.Dir(ex)

			logsDirectory = path.Join(path.Dir(exPath), config.GetConfig("logsDir"))

			_, dirErr := os.Stat(logsDirectory)

			if dirErr != nil {
				if os.IsNotExist(dirErr) {

					mkdirErr := os.Mkdir(logsDirectory, 0700)

					if mkdirErr != nil {
						panic(mkdirErr.Error())
					}
				} else {
					panic(dirErr.Error())
				}
			}

			logsFilePath = path.Join(logsDirectory, config.GetConfig("logsFile"))

			logsFile, fileOpenErr := os.OpenFile(logsFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

			if fileOpenErr != nil {
				panic(fileOpenErr.Error())
			}

			Logger = log.New(logsFile, logsPrefix, log.Lshortfile)
		}

		Logger.SetFlags(log.LstdFlags | log.LUTC | log.Llongfile)
	}
}
