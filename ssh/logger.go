package ssh

import (
	"github.com/dcu/go-authy"
	"log"
	"os"
	"os/user"
)

var (
	// Logger is the default logger for this package.
	Logger *log.Logger
)

func init() {
	user, err := user.Current()
        if err != nil {
                panic(err)
        }
        uid := user.Uid
	logFile, err := os.OpenFile(`/tmp/onetouch-` + uid + `.log`, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	Logger = log.New(logFile, "", 0)
	authy.Logger = Logger
}
