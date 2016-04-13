package ssh

import (
	//"github.com/dcu/onetouch-ssh/utils"
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/fatih/set.v0"
)

var (
	ErrUserAlreadyPresent = errors.New("user is already present")
	ErrUserDoesNotExist   = errors.New("user does not exist")
)

type EachUserHandler func(authyID string, publicKey string)

// UsersManager is in charge of adding/deleting/updating/listing users
type UsersManager struct {
}

// NewUsersManager returns the current users manager instance
func NewUsersManager() *UsersManager {
	return &UsersManager{}
}

func (manager *UsersManager) EachUser(fn EachUserHandler) error {
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

func (manager *UsersManager) UserIDList() []string {
	s := set.New()
	manager.EachUser(func(authyID string, publicKey string) {
		s.Add(authyID)
	})

	return set.StringSlice(s)
}

func (manager *UsersManager) HasUser(userID string) bool {
	found := false
	manager.EachUser(func(authyID string, publicKey string) {
		if authyID == userID {
			found = true
			return
		}
	})

	return found
}

func (manager *UsersManager) GetUser(userID string) *User {
	user := &User{}
	matchingKeys := make([]string, 0)

	manager.EachUser(func(authyID string, publicKey string) {
		if authyID == userID {
			if len(strings.Trim(publicKey, " ")) != 0 {
				matchingKeys = append(matchingKeys, publicKey)
			}
		}
	})

	if len(matchingKeys) == 0 {
		return nil
	}

	user.PublicKeys = matchingKeys
	user.AuthyID = userID

	return user
}

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
