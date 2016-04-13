package ssh

import (
	"github.com/dcu/go-authy"
)

func LoadAuthyAPI() (*authy.Authy, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	api := authy.NewAuthyAPI(config.APIKey)
	api.BaseURL = "https://api.authy.com"

	return api, nil
}
