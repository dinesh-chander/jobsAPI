package cronParser

import (
	"sync"

	"github.com/gorhill/cronexpr"
)

var mutex *sync.Mutex

func Parse(expression string) (parsedExpr *cronexpr.Expression) {
	mutex.Lock()
	parsedExpr = cronexpr.MustParse(expression)
	mutex.Unlock()

	return
}

func init() {
	mutex = &sync.Mutex{}
}
