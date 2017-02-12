package logger

import (
	"config"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var Logger *log.Logger

func init() {

	if config.GetConfig("logsDir") == "" || config.GetConfig("logsFile") == "" {
		Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
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

		Logger = log.New(logsFile, "logger: ", log.Lshortfile)

		go func() {
			for {
				time.Sleep(30 * time.Second)
				logsFile.Sync()
			}
		}()
	}

	Logger.SetPrefix(config.GetConfig("mode"))
	Logger.SetFlags(log.LstdFlags | log.LUTC)

	if config.GetConfig("disableLog") == "true" {
		Logger.SetOutput(ioutil.Discard)
	}
}
