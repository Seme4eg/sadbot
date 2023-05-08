package bot

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

// TODO: maybe use viper in future
// XXX: should config logic b a part of bot package or rather main / other one?
func ReadConfig(config *Config) (err error) {
	var file []byte
	if file, err = ioutil.ReadFile("config.json"); err != nil {
		return
	}
	if err = json.Unmarshal(file, &config); err != nil {
		return
	}
	return
}
