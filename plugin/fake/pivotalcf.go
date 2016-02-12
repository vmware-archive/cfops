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
	FakeSshConfig     command.SshConfig
	FakeJobProperties []cfbackup.Properties
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

func (s *PivotalCF) GetSshConfig(productName, jobName string) (sshConfig command.SshConfig, err error) {
	sshConfig = s.FakeSshConfig
	return
}

func (s *PivotalCF) GetJobProperties(productName, jobName string) (properties []cfbackup.Properties, err error) {
    properties = s.FakeJobProperties
    return
}
