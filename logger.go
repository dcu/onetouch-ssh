package ssh

import (
	"github.com/dcu/go-authy"
	"log"
	"os"
)

var (
	// Logger is the default logger for this package.
	Logger *log.Logger
)

func init() {
	logFile, err := os.OpenFile(`/tmp/onetouch.log`, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	Logger = log.New(logFile, "", 0)
	authy.Logger = Logger
}
