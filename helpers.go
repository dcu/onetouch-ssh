package ssh

import (
	"os"
	"os/user"
	"regexp"
)

func findUserHome() string {
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
		homeRx := regexp.MustCompile(`^/home/[^/]+`)

		matches := homeRx.FindStringSubmatch(wd)
		homeDir = matches[0]
	}

	return homeDir
}
