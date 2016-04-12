package ssh

import (
	"github.com/dcu/onetouch-ssh/utils"
)

// Config contains the configuration of the app.
type Config struct {
	db          *Database
	authyConfig *AuthyConfig
}

// AuthyConfig is the config related to Authy
type AuthyConfig struct {
	APIKey string
}

// ToMap implements the DatabaseData interface
func (authy *AuthyConfig) ToMap() DatabaseData {
	return DatabaseData{
		"APIKey": authy.APIKey,
	}
}

// FromMap implements the DatabaseData interface
func (authy *AuthyConfig) FromMap(data DatabaseData) {
	if value := data["APIKey"]; value != nil {
		authy.APIKey = value.(string)
	}
}

var configInstance *Config

// NewConfig returns a singleton instance of the config.
func NewConfig() *Config {
	if configInstance == nil {
		configInstance = &Config{
			db:          NewDatabase(configDbPath()),
			authyConfig: &AuthyConfig{},
		}

		configInstance.db.Get("authy", configInstance.authyConfig)
	}

	return configInstance
}

// AuthyAPIKey returns the authy's api key.
func (config *Config) AuthyAPIKey() string {
	return config.authyConfig.APIKey
}

// SetAuthyAPIKey sets the authy's api key.
func (config *Config) SetAuthyAPIKey(apiKey string) {
	config.authyConfig.APIKey = apiKey

	config.sync()
}

func (config *Config) sync() {
	config.db.Put("authy", config.authyConfig)
}

func configDbPath() string {
	return utils.FindUserHome() + "/.authy-onetouch/config/"
}
