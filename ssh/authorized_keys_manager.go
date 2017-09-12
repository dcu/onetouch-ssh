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
func (manager *AuthorizedKeysManager) WriteToDefaultLocation() error {
	home := utils.FindUserHome()
	// FIXME: create the .ssh dir if it doesn't exist.

	filename := home + "/.ssh/authorized_keys"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	err = manager.Write(file)
	if err != nil {
		return err
	}

	return os.Chmod(filename, 0600)
}

func (manager *AuthorizedKeysManager) Write(f io.Writer) error {
	w := bufio.NewWriter(f)
	usersManager := NewUsersManager()

	authyShell, err := exec.LookPath("onetouch-ssh")
	if err != nil {
		return err
	}

	// FIXME: keep the old contents.
	_, _ = w.WriteString("### onetouch-ssh\n")
	_ = usersManager.EachEntry(func(authyID string, publicKey string) {
		if len(publicKey) == 0 {
			return
		}

		_, _ = w.WriteString("# " + authyID + "\n")
		publicKey = strings.Trim(publicKey, " ")
		cmd := fmt.Sprintf("%s %s %s", authyShell, "shell", authyID)
		_, _ = w.WriteString(`command="` + cmd + `" ` + publicKey + "\n")
	})

	_, _ = w.WriteString("###\n")
	_ = w.Flush()
	return nil
}

// Contains returns true if the given text is present in the authorized keys file.
func (manager *AuthorizedKeysManager) Contains(text string) bool {
	home := utils.FindUserHome()
	file, err := os.Open(home + "/.ssh/authorized_keys")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

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
	_ = manager.Write(os.Stdout)
}
