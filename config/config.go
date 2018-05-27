package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//RedisConfig defines the config options for a Redis instance
type RedisConfig struct {
	Address string `yaml:"address"`
}

//ServerConfig defines the server config options
type ServerConfig struct {
	Address  string `yaml:"address"`
	AuthFile string `yaml:"auth_file"`
}

//AppConfig holds all app configuration
type AppConfig struct {
	Redis  RedisConfig
	Server ServerConfig
}

// NewConfig loads the config file and returns the Config instance
func NewConfig(file string) (*AppConfig, error) {

	var (
		in  []byte // data structure to hold bytified config when reading from file
		err error  // error handling
	)

	if in, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}

	var cfg = new(AppConfig)

	if err = yaml.Unmarshal(in, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
