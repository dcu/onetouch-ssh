package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// AuthorizedKeysWriter allows to write the authorized_keys file.
type AuthorizedKeysWriter struct {
}

// NewAuthorizedKeysWriter creates a new writer.
func NewAuthorizedKeysWriter() *AuthorizedKeysWriter {
	keysWriter := &AuthorizedKeysWriter{}

	return keysWriter
}

// WriteToDefaultLocation writes the authorized keys to ~/.ssh/authorized_keys
func (writer *AuthorizedKeysWriter) WriteToDefaultLocation() {
	home := findUserHome()
	// FIXME: create the .ssh dir if it doesn't exist.

	file, err := os.Create(home + "/.ssh/authorized_keys")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer.Write(file)
}

func (writer *AuthorizedKeysWriter) Write(f io.Writer) {
	w := bufio.NewWriter(f)
	manager := NewUsersManager()

	authyShell, err := exec.LookPath("authy-shell")
	if err != nil {
		panic(err)
	}

	// FIXME: keep the old contents.
	w.WriteString("### onetouch-ssh\n")
	for _, user := range manager.Users() {
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

// Dump prints the authorized keys file in the stdout.
func (writer *AuthorizedKeysWriter) Dump() {
	writer.Write(os.Stdout)
}
