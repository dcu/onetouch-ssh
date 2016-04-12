package utils

import (
	"fmt"
	"os"
)

func IsInteractiveConnection() bool {
	term := os.Getenv("TERM")

	if term != "" && term != "dumb" {
		return true
	}

	return false
}

func PrintMessage(message string, args ...interface{}) {
	if IsInteractiveConnection() {
		fmt.Printf(message, args...)
	}
}
