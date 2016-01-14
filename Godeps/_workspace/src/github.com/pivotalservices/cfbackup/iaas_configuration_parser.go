package cfbackup

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

//NewConfigurationParser - constructor for a ConfigurationParser from a json installationsettings file
func NewConfigurationParser(installationFilePath string) *ConfigurationParser {
	is := InstallationSettings{}
	b, _ := ioutil.ReadFile(installationFilePath)
	json.Unmarshal(b, &is)

	return &ConfigurationParser{
		installationSettings: is,
	}
}

//NewConfigurationParserFromReader - constructor for a ConfigurationParser from a json installationsettings file
func NewConfigurationParserFromReader(settings io.Reader) *ConfigurationParser {
	is := InstallationSettings{}
	b, _ := ioutil.ReadAll(settings)
	json.Unmarshal(b, &is)

	return &ConfigurationParser{
		installationSettings: is,
	}
}

//GetIaaS - get the iaas elements from the installation settings
func (s *ConfigurationParser) GetIaaS() (config IaaSConfiguration, err error) {
	config = s.installationSettings.Infrastructure.IaaSConfig

	if config.SSHPrivateKey == "" {
		err = ErrNoSSLKeyFound
	}
	return
}

var (
	//ErrNoSSLKeyFound - error if there are no ssl keys found in the iaas config block of installationsettings
	ErrNoSSLKeyFound = errors.New("no ssl key found in iaas config")
)

type (
	//InstallationSettings - an object to house installationsettings elements from the json
	InstallationSettings struct {
		Infrastructure Infrastructure
	}
	//Infrastructure - a struct to house Infrastructure block elements from the json
	Infrastructure struct {
		IaaSConfig IaaSConfiguration `json:"iaas_configuration"`
	}
	//IaaSConfiguration - a struct to house the IaaSConfiguration block elements from the json
	IaaSConfiguration struct {
		SSHPrivateKey string `json:"ssh_private_key"`
	}
	//ConfigurationParser - the parser to handle installation settings file parsing
	ConfigurationParser struct {
		installationSettings InstallationSettings
	}
)
