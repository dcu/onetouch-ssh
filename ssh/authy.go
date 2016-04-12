package ssh

import (
	"github.com/dcu/go-authy"
)

func LoadAuthyAPI() *authy.Authy {
	config := NewConfig()
	api := authy.NewAuthyAPI(config.AuthyAPIKey())
	api.BaseURL = "https://api.authy.com"

	return api
}
