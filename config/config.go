package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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

// Poached the code below from https://github.com/cloudfoundry/cli/cf/configuration/config_helpers/config_helpers.go

func DefaultFilePath() string {
	var configDir string

	if os.Getenv("CFOPS_HOME") != "" {
		cfHome := os.Getenv("CFOPS_HOME")
		configDir = filepath.Join(cfHome, ".cfops")
	} else {
		configDir = filepath.Join(userHomeDir(), ".cfops")
	}

	return filepath.Join(configDir, "config.json")
}

// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
var userHomeDir = func() string {

	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}
