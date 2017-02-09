package logger

import (
	"config"
	"io/ioutil"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	if config.GetConfig("disableLog") == "true" {
		Logger.SetOutput(ioutil.Discard)
	}
}
