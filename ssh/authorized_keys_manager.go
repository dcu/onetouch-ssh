package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dcu/onetouch-ssh/utils"
)

// AuthorizedKeysManager allows to write the authorized_keys file.
type AuthorizedKeysManager struct {
}

// NewAuthorizedKeysManager creates a new writer.
func NewAuthorizedKeysManager() *AuthorizedKeysManager {
	keysManager := &AuthorizedKeysManager{}

	return keysManager
}

// WriteToDefaultLocation writes the authorized keys to ~/.ssh/authorized_keys
func (manager *AuthorizedKeysManager) WriteToDefaultLocation() {
	home := utils.FindUserHome()
	// FIXME: create the .ssh dir if it doesn't exist.

	file, err := os.Create(home + "/.ssh/authorized_keys")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	manager.Write(file)
}

func (manager *AuthorizedKeysManager) Write(f io.Writer) {
	w := bufio.NewWriter(f)
	usersManager := NewUsersManager()

	authyShell, err := exec.LookPath("authy-shell")
	if err != nil {
		panic(err)
	}

	// FIXME: keep the old contents.
	w.WriteString("### onetouch-ssh\n")
	for _, user := range usersManager.Users() {
		if len(user.PublicKeys) == 0 {
			continue
		}

		w.WriteString("# " + user.Username + "\n")
		for _, pk := range user.PublicKeys {
			if pk != "" {
				pk = strings.Trim(pk, " ")
				cmd := fmt.Sprintf("%s %d", authyShell, user.AuthyID)
				w.WriteString(`command="` + cmd + `" ` + pk + "\n")
			}
		}
	}
	w.WriteString("###\n")
	w.Flush()
}

// Contains returns true if the given text is present in the authorized keys file.
func (manager *AuthorizedKeysManager) Contains(text string) bool {
	home := utils.FindUserHome()
	file, err := os.Open(home + "/.ssh/authorized_keys")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return true
		}
	}

	return false
}

// Dump prints the authorized keys file in the stdout.
func (manager *AuthorizedKeysManager) Dump() {
	manager.Write(os.Stdout)
}
