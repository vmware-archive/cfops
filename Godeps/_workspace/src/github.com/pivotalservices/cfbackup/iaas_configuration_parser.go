package cfbackup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

//NewConfigurationParser - constructor for a ConfigurationParser from a json installationsettings file
func NewConfigurationParser(installationFilePath string) *ConfigurationParser {
	is := InstallationSettings{}
	b, _ := ioutil.ReadFile(installationFilePath)
	fmt.Printf("the length of our installation settings file is %d\n", len(b))
	json.Unmarshal(b, &is)
	fmt.Printf("our installation settings is %s\n", is)

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
func (s *ConfigurationParser) GetIaaS() (config IaaSConfiguration, hasSSHKey bool) {
	config = s.installationSettings.Infrastructure.IaaSConfig
	hasSSHKey = false

	if config.SSHPrivateKey != "" {
		hasSSHKey = true
	}
	return
}
