package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token  string `yaml:"token"`
	Prefix string `yaml:"prefix"`
}

// GetConfig returns unmarshals config.yml file and returns new Config struct.
func GetConfig() (config Config, err error) {
	var file []byte
	if file, err = os.ReadFile("config.yml"); err != nil {
		return
	}
	if err = yaml.Unmarshal(file, &config); err != nil {
		return
	}
	return
}
