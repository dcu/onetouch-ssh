package main

import (
	"bufio"
	"fmt"
	"github.com/dcu/onetouch-ssh"
	"os"
)

func editSettings() {
	config := ssh.NewConfig()
	fmt.Printf("Enter the API key: ")
	scanner := bufio.NewScanner(os.Stdin)
	var apiKey string
	if scanner.Scan() {
		apiKey = scanner.Text()
	}

	config.SetAuthyAPIKey(apiKey)
}

func writeAuthorizedKeys() {
	writer := ssh.NewAuthorizedKeysWriter()
	writer.WriteToDefaultLocation()
	fmt.Println("Public keys were wrote to: ~/.ssh/authorized_keys")
}

func dumpAuthorizedKeys() {
	writer := ssh.NewAuthorizedKeysWriter()
	writer.Dump()
}

func main() {
	config := ssh.NewConfig()
	if config.AuthyAPIKey() == "" {
		editSettings()
	}

	if len(os.Args) < 2 {
		usage := `Usage: %s <command>

edit-settings
	Edit initial settings like API key.

edit-users
	Opens a UI to admin the users.

write-authorized-keys
	Writes authorized keys to ~/.ssh/authorized_keys

dump-authorized-keys
	Writes authorized keys to stdout.
`
		fmt.Printf(usage, os.Args[0])
		os.Exit(0)
	}

	switch os.Args[1] {
	case "edit-users":
		{
			app := NewApp()
			app.Start()
		}
	case "edit-settings":
		{
			editSettings()
		}
	case "write-authorized-keys":
		{
			writeAuthorizedKeys()
		}
	case "dump-authorized-keys":
		{
			dumpAuthorizedKeys()
		}
	}
}
