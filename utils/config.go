package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

func ReadConfig(config *Config) (err error) {
	var file []byte
	if file, err = ioutil.ReadFile("config.json"); err != nil {
		return
	}
	if err = json.Unmarshal(file, config); err != nil {
		return
	}
	return
}
