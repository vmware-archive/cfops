package fake

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/pivotalservices/gtils/command"
)

//PivotalCF --
type PivotalCF struct {
	FakeProducts      map[string]cfbackup.Products
	FakeCredentials   map[string]map[string][]cfbackup.Properties
	FakeActivity      string
	FakeReader        io.ReadCloser
	FakeWriter        io.WriteCloser
	FakeHostDetails   tileregistry.TileSpec
	FakeJobProperties []cfbackup.Properties
	FakePropertyMap   map[string]string
	FakeSSHConfig     command.SshConfig
	FakeIP            string
}

//GetHostDetails --
func (s *PivotalCF) GetHostDetails() tileregistry.TileSpec {
	return s.FakeHostDetails
}

//SetActivity --
func (s *PivotalCF) SetActivity(activity string) {
	s.FakeActivity = activity
}

//GetActivity --
func (s *PivotalCF) GetActivity() string {
	return s.FakeActivity
}

//GetProducts --
func (s *PivotalCF) GetProducts() (products map[string]cfbackup.Products) {
	return s.FakeProducts
}

//GetCredentials --
func (s *PivotalCF) GetCredentials() (credentials map[string]map[string][]cfbackup.Properties) {
	return s.FakeCredentials
}

//NewArchiveReader -- fake archive reader
func (s *PivotalCF) NewArchiveReader(name string) (reader io.ReadCloser, err error) {
	reader = s.FakeReader
	return
}

//NewArchiveWriter -- fake archive writer
func (s *PivotalCF) NewArchiveWriter(name string) (writer io.WriteCloser, err error) {
	writer = s.FakeWriter
	return
}

//GetJobProperties - returns fake []cfbackup.Properties
func (s *PivotalCF) GetJobProperties(productName, jobName string) (properties []cfbackup.Properties, err error) {
	properties = s.FakeJobProperties
	return
}

//GetPropertyValues - returns fake map[string]string
func (s *PivotalCF) GetPropertyValues(productName, jobName, identifier string) (propertyMap map[string]string, err error) {
	propertyMap = s.FakePropertyMap
	return
}

//GetSSHConfig - returns command.SshConfig for a given product, job vm
func (s *PivotalCF) GetSSHConfig(productName, jobName string) (sshConfig command.SshConfig, err error) {
	sshConfig = s.FakeSSHConfig
	return
}

//GetJobIP - returns ip for a given product, job vm
func (s *PivotalCF) GetJobIP(productName, jobName string) (ip string, err error) {
	ip = s.FakeIP
	return
}
