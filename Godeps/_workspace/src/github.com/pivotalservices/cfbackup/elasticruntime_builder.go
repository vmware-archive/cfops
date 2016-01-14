package cfbackup

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pivotalservices/cfops/tileregistry"
)

//New -- method to generate an initialized elastic runtime
func (s *ElasticRuntimeBuilder) New(tileSpec tileregistry.TileSpec) (elasticRuntime tileregistry.Tile, err error) {
	var (
		installationSettings io.Reader
		installationTmpFile  *os.File
	)

	if installationSettings, err = GetInstallationSettings(tileSpec); err == nil {
		installationTmpFile, err = ioutil.TempFile("", OpsMgrInstallationSettingsFilename)
		defer installationTmpFile.Close()
		io.Copy(installationTmpFile, installationSettings)
		elasticRuntime = NewElasticRuntime(installationTmpFile.Name(), tileSpec.ArchiveDirectory)
	}
	return
}

//GetInstallationSettings - makes a call to ops manager and returns a io.reader containing the contents of the installation settings file.
var GetInstallationSettings = func(tileSpec tileregistry.TileSpec) (settings io.Reader, err error) {
	var (
		opsManager *OpsManager
	)

	if opsManager, err = NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory); err == nil {
		settings, err = opsManager.GetInstallationSettings()
	}
	return
}
