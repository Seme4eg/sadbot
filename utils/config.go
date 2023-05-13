package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token  string `yaml:"token"`
	Prefix string `yaml:"prefix"`
}

func NewConfig() (config *Config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile("config.yml"); err != nil {
		return
	}
	if err = yaml.Unmarshal(file, &config); err != nil {
		return
	}
	return
}
