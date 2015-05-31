package ssh

import (
	"strconv"
)

// User is a struct contains the user's info.
type User struct {
	Username    string
	PublicKey   string
	CountryCode int
	PhoneNumber string
}

// NewUser returns a new instance of User
func NewUser(username string) *User {
	user := &User{
		Username: username,
	}
	return user
}

// Save saves the user
func Save() bool {
	return false
}

// CountryCodeStr returns the country code as a string.
func (user *User) CountryCodeStr() string {
	if user.CountryCode == 0 {
		return ""
	}

	return strconv.Itoa(user.CountryCode)
}

// ToMap converts the user to a map
func (user *User) ToMap() ConfigData {
	return ConfigData{
		"Username":    user.Username,
		"PublicKey":   user.PublicKey,
		"CountryCode": user.CountryCode,
		"PhoneNumber": user.PhoneNumber,
	}
}

// FromMap loads the user using a map.
func (user *User) FromMap(data ConfigData) {
	user.Username = data["Username"].(string)
}
