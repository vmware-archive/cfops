package cfopsplugin

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/pivotalservices/gtils/command"
)

const (
	vmCredentialsName     = "vm_credentials"
	identityName          = "identity"
	passwordName          = "password"
	defaultSSHPort    int = 22
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

func (s *DefaultPivotalCF) getProduct(productName string) (product cfbackup.Products, err error) {
	if _, ok := s.GetProducts()[productName]; ok {
		product = s.GetProducts()[productName]
	} else {
		err = fmt.Errorf("product %s not found", productName)
	}
	return
}
func (s *DefaultPivotalCF) getJob(productName, jobName string) (job cfbackup.Jobs, err error) {
	var product cfbackup.Products
	var jobFound = false

	product, err = s.getProduct(productName)
	if err != nil {
		return
	}
	for _, theJob := range product.Jobs {
		if theJob.Identifier == jobName {
			job = theJob
			jobFound = true
			break
		}
	}
	if !jobFound {
		err = fmt.Errorf("job %s not found for product %s", jobName, productName)
	}
	return
}

func (s *DefaultPivotalCF) getSSLKey(productName, jobName string) (sslKey string, err error) {
	//sslKey = ""
	return
}

//GetJobProperties - returns []cfbackup.Properties for a given product and job
func (s *DefaultPivotalCF) GetJobProperties(productName, jobName string) (properties []cfbackup.Properties, err error) {
	var job cfbackup.Jobs
	job, err = s.getJob(productName, jobName)
	if err != nil {
		return
	}
	properties = job.Properties
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
				propertyMap[key] = fmt.Sprintf("%v", value)
			}
		}
	}
	return
}

//GetSSHConfig - returns command.SshConfig for a given product, job vm
func (s *DefaultPivotalCF) GetSSHConfig(productName, jobName string) (sshConfig command.SshConfig, err error) {
	var userid, password, ip, sslKey string
	var props map[string]string

	props, err = s.GetPropertyValues(productName, jobName, vmCredentialsName)
	if err != nil {
		return
	}

	ip, err = s.GetJobIP(productName, jobName)
	if err != nil {
		return
	}
	
	sslKey, err = s.getSSLKey(productName, jobName)
	if err != nil {
		return
	}

	userid = props[identityName]
	password = props[passwordName]

	sshConfig = command.SshConfig{
		Username: userid,
		Password: password,
		Host:     ip,
		Port:     defaultSSHPort,
		SSLKey:   sslKey,
	}

	return
}

//GetJobIP - returns ip for a given product, job vm
func (s *DefaultPivotalCF) GetJobIP(productName, jobName string) (ip string, err error) {
	var ipFound = false
	var job cfbackup.Jobs
	var product cfbackup.Products

	product, err = s.getProduct(productName)
	if err != nil {
		return
	}
	job, err = s.getJob(productName, jobName)
	if err != nil {
		return
	}

	for vmName, ipList := range product.IPS {
		if strings.HasPrefix(vmName, job.GUID) {
			ip = ipList[0]
			ipFound = true
			break
		}
	}
	if !ipFound {
		err = fmt.Errorf("ip not found for job %s and product %s", jobName, productName)
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
