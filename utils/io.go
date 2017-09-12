package utils

import (
	"fmt"
	"os"
)

// IsInteractiveConnection returns true if the terminal is interactive
func IsInteractiveConnection() bool {
	term := os.Getenv("TERM")

	if term != "" && term != "dumb" {
		return true
	}

	return false
}

// PrintMessage prints a message only if the terminal is interactive, otherwise the message is ignored.
func PrintMessage(message string, args ...interface{}) {
	if IsInteractiveConnection() {
		fmt.Printf(message, args...)
	}
}
