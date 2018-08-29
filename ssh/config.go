package ssh

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dcu/onetouch-ssh/utils"
	"gopkg.in/yaml.v2"
)

var (
	// DataPath is the default path to read the config from.
	DataPath = utils.FindUserHome() + "/.authy-onetouch"
)

// Config contains the configuration of the app.
type Config struct {
	APIKey    string   `yaml:"api_key"`
	ShellPath string   `yaml:"shell"`
	ShellArgs []string `yaml:"shell_args"`
}

// NewConfig builds a new config object.
func NewConfig(apiKey string) *Config {
	config := &Config{
		APIKey: apiKey,
	}

	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "/bin/sh"
	}

	config.ShellPath = shell
	config.ShellArgs = make([]string, 0)

	fmt.Println("Your users will be logged in a '" + shell + "' shell")

	return config
}

// LoadConfig loads the default config
func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(configDbFile())
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// Save stores the current config
func (config *Config) Save() error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configDbFile(), data, 0600)
}

func configDbFile() string {
	return DataPath + "/config.yml"
}

func init() {
	_ = os.MkdirAll(DataPath, 0700)
}
