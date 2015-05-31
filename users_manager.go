package ssh

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"regexp"
)

var ()

// UsersManagerListener is a interface to listen to user's events.
type UsersManagerListener interface {
	OnUserAdded(user *User)
}

// UsersManager is in charge of adding/deleting/updating/listing users
type UsersManager struct {
	listeners []UsersManagerListener
	config    *Database
}

var usersManagerInstance *UsersManager

// NewUsersManager returns the current users manager instance
func NewUsersManager() *UsersManager {
	if usersManagerInstance == nil {
		usersManagerInstance = &UsersManager{
			config: NewDatabase(usersDbPath()),
		}
	}

	return usersManagerInstance
}

// AddListener adds a new listener
func (manager *UsersManager) AddListener(listener UsersManagerListener) {
	manager.listeners = append(manager.listeners, listener)
}

// UpdateUser updates the given user
func (manager *UsersManager) UpdateUser(user *User) error {
	if len(user.Username) == 0 {
		return errors.New("Username can't be empty.")
	}

	err := manager.config.Put(user.Username, user)
	if err == nil {
		return nil
	}

	return err

}

// AddUser adds a new user
func (manager *UsersManager) AddUser(user *User) error {
	if len(user.Username) == 0 {
		return errors.New("Username can't be empty.")
	}

	err := manager.config.Put(user.Username, user)
	if err == nil {
		manager.onUserAdded(user)
		return nil
	}

	return err
}

// HasUsers returns true if the manager has users.
func (manager *UsersManager) HasUsers() bool {
	return manager.config.HasKeys()
}

// Users returns the list of users
func (manager *UsersManager) Users() []*User {
	users := []*User{}
	for _, data := range manager.config.List() {
		fmt.Printf("%#v\n", data)

		user := &User{}
		user.FromMap(data)
		users = append(users, user)
	}

	return users
}

// LoadUser loads a user using the given username
func (manager *UsersManager) LoadUser(username string) *User {
	if username == "" {
		return nil
	}

	user := &User{}
	err := manager.config.Get(username, user)
	if err != nil {
		panic(err)
	}

	return user
}

func usersDbPath() string {
	return findUserHome() + "/.authy-onetouch/users/"
}

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

func (manager *UsersManager) onUserAdded(user *User) {
	for _, listener := range manager.listeners {
		listener.OnUserAdded(user)
	}
}
