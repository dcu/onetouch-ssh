package ssh

import (
	"bufio"
	"os"
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

func (writer *AuthorizedKeysWriter) Write() {
	home := findUserHome()
	file, err := os.Create(home + "/.ssh/authorized_keys")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// FIXME: create the .ssh dir if it doesn't exist.
	w := bufio.NewWriter(file)
	manager := NewUsersManager()

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
				cmd := "authy-shell"
				w.WriteString(`command="` + cmd + `" ` + pk + "\n")
			}
		}
	}
	w.WriteString("###\n")
	w.Flush()
}
