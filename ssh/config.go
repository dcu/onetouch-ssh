package ssh

import (
	"io/ioutil"
	"os"

	"github.com/dcu/onetouch-ssh/utils"
	"gopkg.in/yaml.v2"
)

var (
	DataPath = utils.FindUserHome() + "/.authy-onetouch"
)

// Config contains the configuration of the app.
type Config struct {
	APIKey string `yaml:"api_key"`
}

func NewConfig(apiKey string) *Config {
	config := &Config{
		APIKey: apiKey,
	}

	return config
}

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
	os.MkdirAll(DataPath, 0700)
}
