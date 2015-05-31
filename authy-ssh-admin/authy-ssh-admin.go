package main

import (
	"bufio"
	"fmt"
	"github.com/authy/onetouch-ssh"
	"os"
)

func main() {
	config := ssh.NewConfig()

	if config.AuthyAPIKey() == "" {
		fmt.Printf("Enter the API key: ")
		scanner := bufio.NewScanner(os.Stdin)
		var apiKey string
		if scanner.Scan() {
			apiKey = scanner.Text()
		}

		config.SetAuthyAPIKey(apiKey)
	}

	app := NewApp()

	app.Start()
}
