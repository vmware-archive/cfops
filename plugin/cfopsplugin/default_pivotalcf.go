package cfopsplugin

import (
	"io"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
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

//NewPivotalCF - creates the default pivotacf
var NewPivotalCF = func(installationSettings *cfbackup.ConfigurationParser, ts tileregistry.TileSpec) PivotalCF {

	return &DefaultPivotalCF{
		TileSpec:             ts,
		InstallationSettings: installationSettings,
	}
}
