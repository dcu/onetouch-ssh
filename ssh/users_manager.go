package ssh

import (
	//"github.com/dcu/onetouch-ssh/utils"
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"gopkg.in/fatih/set.v0"
)

var (
	// ErrUserAlreadyPresent is returned when the user already exists in the db.
	ErrUserAlreadyPresent = errors.New("user is already present")

	// ErrUserDoesNotExist is returned when the user doesn't exist in the db.
	ErrUserDoesNotExist = errors.New("user does not exist")
)

// EachEntryHandler is function prototype for the EachEntry callback
type EachEntryHandler func(authyID string, publicKey string)

// EachUserHandler is function prototype for the EachUser callback
type EachUserHandler func(user *User)

// UsersManager is in charge of adding/deleting/updating/listing users
type UsersManager struct {
}

// NewUsersManager returns the current users manager instance
func NewUsersManager() *UsersManager {
	return &UsersManager{}
}

// EachEntry goes through every entry in the users db and calls fn with it.
func (manager *UsersManager) EachEntry(fn EachEntryHandler) error {
	file, err := os.Open(usersDbPath())
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userData := strings.SplitN(scanner.Text(), " ", 2)
		fn(userData[0], userData[1])
	}

	return nil
}

// EachUser goes through every user in the users db and calls fn with it.
func (manager *UsersManager) EachUser(fn EachUserHandler) error {
	userIDs := manager.UserIDList()
	for _, userID := range userIDs {
		user := FindUser(userID)
		if user != nil {
			fn(user)
		}
	}
	return nil
}

// UserIDList returns the list of user ids present in the users db.
func (manager *UsersManager) UserIDList() []string {
	s := set.New()
	manager.EachEntry(func(authyID string, publicKey string) {
		s.Add(authyID)
	})

	return set.StringSlice(s)
}

// HasUser returns true if the user is present in the local db.
func (manager *UsersManager) HasUser(userID string) bool {
	found := false
	manager.EachEntry(func(authyID string, publicKey string) {
		if authyID == userID {
			found = true
			return
		}
	})

	return found
}

// AddUser adds a user to the users db.
func (manager *UsersManager) AddUser(email string, countryCode int, phoneNumber string, publicKey string) error {
	api, err := LoadAuthyAPI()
	if err != nil {
		return err
	}
	publicKey = strings.Replace(publicKey, "\n", "", -1)

	user, err := api.RegisterUser(email, countryCode, phoneNumber, url.Values{})
	if err != nil {
		for field, msg := range user.Errors {
			fmt.Printf("%s=%s\n", field, msg)
		}
		return err
	}

	return manager.AddUserID(user.ID, publicKey)
}

// AddUserID adds a user id to the users db.
func (manager *UsersManager) AddUserID(authyID string, publicKey string) error {
	if manager.HasUser(authyID) {
		return ErrUserAlreadyPresent
	}

	file, err := os.OpenFile(usersDbPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s %s\n", authyID, publicKey))
	return err
}

// RemoveUser removes the user with the given `authyID`
func (manager *UsersManager) RemoveUser(authyID string) error {
	if match, err := regexp.MatchString(`\A[0-9]+\z`, authyID); !match || err != nil {
		return errors.New("invalid authy id")
	}

	tmpFile, err := os.OpenFile(usersDbPath()+".tmp", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	file, err := os.Open(usersDbPath())
	if err != nil {
		return err
	}

	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, authyID+" ") {
			// ignore
			continue
		}
		tmpFile.WriteString(line)
	}

	file.Close()
	tmpFile.Close()

	return os.Rename(tmpFile.Name(), file.Name())
}

// AddKey associates a key to the given user id.
func (manager *UsersManager) AddKey(authyID string, publicKey string) error {
	if manager.HasUser(authyID) != true {
		return ErrUserDoesNotExist
	}

	file, err := os.OpenFile(usersDbPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s %s\n", authyID, publicKey))
	return err
}

func usersDbPath() string {
	return DataPath + "/users.list"
}
