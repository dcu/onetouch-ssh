package ssh

import (
	"errors"
	"fmt"
	"github.com/dcu/go-authy"
	"net/url"
	"strconv"
)

// User is a struct contains the user's info.
type User struct {
	Username    string
	PublicKeys  []string
	Email       string
	CountryCode int
	PhoneNumber string
	AuthyID     string
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

// AuthyIDStr returns the authy id as a string.
func (user *User) AuthyIDStr() string {
	if user.AuthyID == "" {
		return "<not set>"
	}

	return user.AuthyID
}

// ToMap converts the user to a map
func (user *User) ToMap() DatabaseData {
	return DatabaseData{
		"Username":    user.Username,
		"AuthyID":     user.AuthyID,
		"Email":       user.Email,
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
		user.AuthyID = value.(string)
	}
	if value := data["Email"]; value != nil {
		user.Email = value.(string)
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

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	api := authy.NewAuthyAPI(config.APIKey)
	api.BaseURL = "https://api.authy.com"

	authyUser, err := api.RegisterUser(user.Email, user.CountryCode, user.PhoneNumber, url.Values{})
	if err != nil {
		return err
	}

	user.AuthyID = authyUser.ID
	return nil
}

// ValueForColumn returns the value for the given column name.
func (user *User) ValueForColumn(columnName string) string {
	switch columnName {
	case "ID":
		{
			return user.AuthyIDStr()
		}
	case "Username":
		{
			return user.Username
		}
	case "Email":
		{
			return user.Email
		}
	case "Phone Number":
		{
			return fmt.Sprintf("+%s %s", user.CountryCodeStr(), user.PhoneNumber)
		}
	case "Configured":
		{
			if user.IsConfigured() {
				return "YES"
			}
			return "NO"
		}
	case "Protected":
		{
			if user.IsProtected() {
				return "YES"
			}
			return "NO"
		}
	}
	return "<unknown>"
}

// IsConfigured returns true if the user is fully configured
func (user *User) IsConfigured() bool {
	if user.AuthyID != "" && len(user.PublicKeys) > 0 {
		return true
	}
	return false
}

// IsProtected returns true if the user is fully configured
func (user *User) IsProtected() bool {
	manager := NewAuthorizedKeysManager()
	return manager.Contains(fmt.Sprintf("authy-shell %d", user.AuthyID))
}
