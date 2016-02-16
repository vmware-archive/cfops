package cfopsplugin

import (
	"fmt"
	"io"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/pivotalservices/gtils/command"
)

//GetHostDetails - return all of the host and archive details in the form of a tile spec object
func (s *DefaultPivotalCF) GetHostDetails() tileregistry.TileSpec {
	return s.TileSpec
}

//GetProducts - gets a products object from the given pivotalcf
func (s *DefaultPivotalCF) GetProducts() (products map[string]cfbackup.Products) {
	products = make(map[string]cfbackup.Products)

	for _, product := range s.InstallationSettings.GetProducts() {
		products[product.Identifier] = product
	}
	return
}

//GetCredentials - gets a credentials object from the given pivotalcf
func (s *DefaultPivotalCF) GetCredentials() (creds map[string]map[string][]cfbackup.Properties) {
	creds = make(map[string]map[string][]cfbackup.Properties)

	for _, product := range s.InstallationSettings.GetProducts() {
		creds[product.Identifier] = make(map[string][]cfbackup.Properties)

		for _, job := range product.Jobs {
			creds[product.Identifier][job.Identifier] = job.Properties
		}
	}
	return
}

//NewArchiveWriter - creates a writer to a named resource using the given name on the cfops defined target (s3, local, etc)
func (s *DefaultPivotalCF) NewArchiveWriter(name string) (writer io.WriteCloser, err error) {
	backupContext := cfbackup.NewBackupContext(s.TileSpec.ArchiveDirectory, cfenv.CurrentEnv())
	writer, err = backupContext.StorageProvider.Writer(path.Join(s.TileSpec.ArchiveDirectory, name))
	return
}

//NewArchiveReader - creates a reader to a named resource using the given name on the cfops defined target (s3, local, etc)
func (s *DefaultPivotalCF) NewArchiveReader(name string) (reader io.ReadCloser, err error) {
	backupContext := cfbackup.NewBackupContext(s.TileSpec.ArchiveDirectory, cfenv.CurrentEnv())
	reader, err = backupContext.StorageProvider.Reader(path.Join(s.TileSpec.ArchiveDirectory, name))
	return
}

//GetSSHConfig - returns the information needed to Ssh into product VM
func (s *DefaultPivotalCF) GetSSHConfig(productName, jobName string) (sshConfig command.SshConfig, err error) {
	//properties := s.GetJobProperties(productName, jobName)

	return
}

//GetJobProperties - returns []cfbackup.Properties for a given product and job
func (s *DefaultPivotalCF) GetJobProperties(productName, jobName string) (properties []cfbackup.Properties, err error) {
	var jobFound = false
	if _, ok := s.GetProducts()[productName]; ok {
		product := s.GetProducts()[productName]
		for _, job := range product.Jobs {
			if job.Identifier == jobName {
				properties = job.Properties
				jobFound = true
			}
		}
	} else {
		err = fmt.Errorf("product %s not found", productName)
	}
	if !jobFound {
	    err = fmt.Errorf("job %s not found for product %s", jobName, productName)
	}
	return
}

//GetPropertyValues - returns map[string]string for a given product, job and property identifier
func (s *DefaultPivotalCF) GetPropertyValues(productName, jobName, identifier string) (propertyMap map[string]string, err error) {
    var properties []cfbackup.Properties
    
	properties, err = s.GetJobProperties(productName, jobName)
	propertyMap = make(map[string]string)
	for _, property := range properties {

		if property.Identifier == identifier {
			pMap := property.Value.(map[string]interface{})
			for key, value := range pMap {
			    propertyMap[key]  = value.(string)
			}
		}
	}
    return
}

//NewPivotalCF - creates the default pivotacf
var NewPivotalCF = func(installationSettings *cfbackup.ConfigurationParser, ts tileregistry.TileSpec) PivotalCF {

	return &DefaultPivotalCF{
		TileSpec:             ts,
		InstallationSettings: installationSettings,
	}
}
