package fake

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tileregistry"
)

//PivotalCF --
type PivotalCF struct {
	FakeActivity             string
	FakeReader               io.ReadCloser
	FakeWriter               io.WriteCloser
	FakeHostDetails          tileregistry.TileSpec
	FakeInstallationSettings cfbackup.InstallationSettings
}

//GetHostDetails --
func (s *PivotalCF) GetHostDetails() tileregistry.TileSpec {
	return s.FakeHostDetails
}

//GetInstallationSettings --
func (s *PivotalCF) GetInstallationSettings() cfbackup.InstallationSettings {
	return s.FakeInstallationSettings
}

//SetActivity --
func (s *PivotalCF) SetActivity(activity string) {
	s.FakeActivity = activity
}

//GetActivity --
func (s *PivotalCF) GetActivity() string {
	return s.FakeActivity
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
