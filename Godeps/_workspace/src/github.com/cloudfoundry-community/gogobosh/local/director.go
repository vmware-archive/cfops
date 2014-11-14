package local

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"launchpad.net/goyaml"
)

// BoshConfig describes a local ~/.bosh_config file
// See testhelpers/fixtures/bosh_config.yml
type BoshConfig struct {
	Target         string
	Name           string `yaml:"target_name"`
	Version        string `yaml:"target_version"`
	UUID           string `yaml:"target_uuid"`
	Aliases        map[string]map[string]string
	Authentication map[string]*authentication `yaml:"auth"`
}

type authentication struct {
	Username string
	Password string
}

// LoadBoshConfig loads and unmarshals ~/.bosh_config
func LoadBoshConfig(configPath string) (config *BoshConfig, err error) {
	config = &BoshConfig{}

	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	goyaml.Unmarshal(contents, config)
	return
}

// DefaultBoshConfigPath returns the path to ~/.bosh_config
func DefaultBoshConfigPath() (configPath string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Abs(usr.HomeDir + "/.bosh_config")
}

// CurrentBoshTarget returns the connection information for local user's current target BOSH
func (config *BoshConfig) CurrentBoshTarget() (target, username, password string, err error) {
	if config.Target == "" {
		return "", "", "", errors.New("Please target a BOSH first. Run 'bosh target DIRECTOR_IP'.")
	}
	auth := config.Authentication[config.Target]
	if auth == nil {
		return "", "", "", errors.New("Current target has not been authenticated yet. Run 'bosh login'.")
	}
	return config.Target, auth.Username, auth.Password, nil
}
