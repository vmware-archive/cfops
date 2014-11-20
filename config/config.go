package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func LoadConfig(config interface{}, configPath string) error {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot load config file [%s]: %s", configPath, err))
	}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot parse config file [%s]: %s", configPath, err))
	}
	return nil
}
