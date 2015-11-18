package cfbackup

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pivotalservices/cfops/tileregistry"
)

//ElasticRuntimeBuilder -- an object that can build an elastic runtime pre-initialized
type ElasticRuntimeBuilder struct{}

//New -- method to generate an initialized elastic runtime
func (s *ElasticRuntimeBuilder) New(tileSpec tileregistry.TileSpec) (elasticRuntime tileregistry.Tile, err error) {
	var (
		installationSettings io.Reader
		installationTmpFile  *os.File
	)

	if installationSettings, err = GetInstallationSettings(tileSpec); err == nil {
		installationTmpFile, err = ioutil.TempFile("", OPSMGR_INSTALLATION_SETTINGS_FILENAME)
		defer installationTmpFile.Close()
		io.Copy(installationTmpFile, installationSettings)
		elasticRuntime = NewElasticRuntime(installationTmpFile.Name(), tileSpec.ArchiveDirectory)
	}
	return
}

var GetInstallationSettings = func(tileSpec tileregistry.TileSpec) (settings io.Reader, err error) {
	var (
		opsManager *OpsManager
	)

	if opsManager, err = NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory); err == nil {
		settings, err = opsManager.GetInstallationSettings()
	}
	return
}
