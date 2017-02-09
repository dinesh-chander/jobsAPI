package cache

import (
	"log"
	"logger"
	"time"

	"github.com/boltdb/bolt"
)

var cache *bolt.DB
var loggerInstance *log.Logger

func init() {
	loggerInstance = logger.Logger

	cache, err := bolt.Open("cache.db", 0600, &bolt.Options{Timeout: 30 * time.Second})

	if err != nil {
		loggerInstance.Panicln(err.Error())
	}
}
