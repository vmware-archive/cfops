package cfopsplugin

import (
	"io"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tileregistry"
)

//GetHostDetails - return all of the host and archive details in the form of a tile spec object
func (s *DefaultPivotalCF) GetHostDetails() tileregistry.TileSpec {
	return s.TileSpec
}

//GetInstallationSettings - return installation settings
func (s *DefaultPivotalCF) GetInstallationSettings() cfbackup.InstallationSettings {
	return s.InstallationSettings
}

//NewArchiveWriter - creates a writer to a named resource using the given name on the cfops defined target (s3, local, etc)
func (s *DefaultPivotalCF) NewArchiveWriter(name string) (writer io.WriteCloser, err error) {
	backupContext := cfbackup.NewBackupContext(s.TileSpec.ArchiveDirectory, cfenv.CurrentEnv(), s.TileSpec.CryptKey)
	writer, err = backupContext.StorageProvider.Writer(path.Join(s.TileSpec.ArchiveDirectory, name))
	return
}

//NewArchiveReader - creates a reader to a named resource using the given name on the cfops defined target (s3, local, etc)
func (s *DefaultPivotalCF) NewArchiveReader(name string) (reader io.ReadCloser, err error) {
	backupContext := cfbackup.NewBackupContext(s.TileSpec.ArchiveDirectory, cfenv.CurrentEnv(), s.TileSpec.CryptKey)
	reader, err = backupContext.StorageProvider.Reader(path.Join(s.TileSpec.ArchiveDirectory, name))
	return
}

//NewPivotalCF - creates the default pivotacf
var NewPivotalCF = func(installationSettings cfbackup.InstallationSettings, ts tileregistry.TileSpec) PivotalCF {

	return &DefaultPivotalCF{
		TileSpec:             ts,
		InstallationSettings: installationSettings,
	}
}
