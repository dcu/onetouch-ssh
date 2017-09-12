package ssh

import (
	"errors"
	"net/url"
	"strings"
)

var (
	errInvalidAuthyID = errors.New("invalid Authy ID")
)

// User is a struct contains the user's info.
type User struct {
	PublicKeys  []string
	CountryCode int
	PhoneNumber string
	AuthyID     string
}

// FindUser finds the user in the local database
func FindUser(userID string) *User {
	user := &User{
		AuthyID: userID,
	}
	matchingKeys := []string{}

	usersManager := NewUsersManager()
	_ = usersManager.EachEntry(func(authyID string, publicKey string) {
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

	return user
}

// LoadFromAuthy loads user data from the Authy API
func (user *User) LoadFromAuthy() error {
	if len(user.AuthyID) == 0 {
		return errInvalidAuthyID
	}

	api, err := LoadAuthyAPI()
	if err != nil {
		return err
	}

	authyUser, err := api.UserStatus(user.AuthyID, url.Values{})
	if err != nil {
		return err
	}

	user.AuthyID = authyUser.ID
	user.CountryCode = authyUser.StatusData.Country
	user.PhoneNumber = authyUser.StatusData.PhoneNumber

	return nil
}
