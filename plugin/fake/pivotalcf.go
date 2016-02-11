package fake

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
)

//PivotalCF --
type PivotalCF struct {
	FakeProducts    map[string]cfbackup.Products
	FakeCredentials map[string]map[string][]cfbackup.Properties
	FakeActivity    string
	FakeReader      io.Reader
	FakeWriter      io.Writer
	FakeHostDetails tileregistry.TileSpec
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
func (s *PivotalCF) NewArchiveReader(name string) (reader io.Reader) {
	reader = s.FakeReader
	return
}

//NewArchiveWriter -- fake archive writer
func (s *PivotalCF) NewArchiveWriter(name string) (writer io.Writer) {
	writer = s.FakeWriter
	return
}
