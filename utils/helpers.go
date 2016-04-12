package utils

import (
	"os"
	"os/user"
	"regexp"
)

// FindUserHome returns the home directory of the current user.
func FindUserHome() string {
	var homeDir string

	user, err := user.Current()
	if err == nil {
		homeDir = user.HomeDir
	}

	if homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	if homeDir == "" {
		wd, _ := os.Getwd()
		homeRx := regexp.MustCompile(`\A/home/[^/]+`)

		matches := homeRx.FindStringSubmatch(wd)
		homeDir = matches[0]
	}

	return homeDir
}
