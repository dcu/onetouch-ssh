package ssh

import (
	"errors"
	//"github.com/dcu/go-authy"
	"strconv"
)

// User is a struct contains the user's info.
type User struct {
	Username    string
	PublicKeys  []string
	CountryCode int
	PhoneNumber string
	AuthyID     int
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
func (user *User) ToMap() DatabaseData {
	return DatabaseData{
		"Username":    user.Username,
		"AuthyID":     user.AuthyID,
		"PublicKeys":  user.PublicKeys,
		"CountryCode": user.CountryCode,
		"PhoneNumber": user.PhoneNumber,
	}
}

// FromMap loads the user using a map.
func (user *User) FromMap(data DatabaseData) {
	if value := data["Username"]; value != nil {
		user.Username = value.(string)
	}
	if value := data["AuthyID"]; value != nil {
		user.AuthyID = value.(int)
	}
	if value := data["CountryCode"]; value != nil {
		user.CountryCode = value.(int)
	}
	if value := data["PhoneNumber"]; value != nil {
		user.PhoneNumber = value.(string)
	}
	if value := data["PublicKeys"]; value != nil {
		user.PublicKeys = value.([]string)
	}
}

// Register register the user on Authy
func (user *User) Register() error {
	if len(user.PhoneNumber) == 0 || user.CountryCode == 0 {
		return errors.New("Invalid phone number.")
	}

	//authyUser, err := authyApi.RegisterUser(authy.UserOpts{Email: email, PhoneNumber: cellphone, CountryCode: countryCode})

	return nil
}
