package cfbackup

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/xchapter7x/lo"
)

//NewConfigurationParser - constructor for a ConfigurationParser from a json installationsettings file
func NewConfigurationParser(installationFilePath string) *ConfigurationParser {
	is := InstallationSettings{}
	b, _ := ioutil.ReadFile(installationFilePath)
	if err := json.Unmarshal(b, &is); err != nil {
		lo.G.Error("Unmarshal installation settings error", err)
	}
    is.SetPGDumpUtilVersions()
	return &ConfigurationParser{
		InstallationSettings: is,
	}
}

//NewConfigurationParserFromReader - constructor for a ConfigurationParser from a json installationsettings file
func NewConfigurationParserFromReader(settings io.Reader) *ConfigurationParser {
	is := InstallationSettings{}
	b, _ := ioutil.ReadAll(settings)
	if err := json.Unmarshal(b, &is); err != nil {
		lo.G.Error("Unmarshal installation settings error", err)
	}
	
    is.SetPGDumpUtilVersions()
	return &ConfigurationParser{
		InstallationSettings: is,
	}
}

//GetIaaS - get the iaas elements from the installation settings
func (s *ConfigurationParser) GetIaaS() (config IaaSConfiguration, hasSSHKey bool) {
	config = s.InstallationSettings.Infrastructure.IaaSConfig
	hasSSHKey = false

	if config.SSHPrivateKey != "" {
		hasSSHKey = true
	}
	return
}

// FindJobsByProductID finds all the jobs in an installation by product id
func (s *ConfigurationParser) FindJobsByProductID(id string) (jobs []Jobs) {
	jobs = s.InstallationSettings.FindJobsByProductID(id)
	return
}

// FindByProductID finds a product by product id
func (s *ConfigurationParser) FindByProductID(id string) (productResponse Products, err error) {
	productResponse, err = s.InstallationSettings.FindByProductID(id)
	return
}

// FindCFPostgresJobs finds all the postgres jobs in the cf product
func (s *ConfigurationParser) FindCFPostgresJobs() (jobs []Jobs) {
	jobs = s.InstallationSettings.FindCFPostgresJobs()
	return
}

//GetProducts - get the products array
func (s *ConfigurationParser) GetProducts() (products []Products) {
	return s.InstallationSettings.Products
}
