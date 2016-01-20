package elasticruntime

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
)

//New -- method to generate an initialized elastic runtime
func (s *ElasticRuntimeBuilder) New(tileSpec tileregistry.TileSpec) (elasticRuntime tileregistry.Tile, err error) {
	var (
		installationSettings io.Reader
		installationTmpFile  *os.File
		sshKey               = ""
	)

	if installationSettings, err = GetInstallationSettings(tileSpec); err == nil {
		installationTmpFile, err = ioutil.TempFile("", opsmanager.OpsMgrInstallationSettingsFilename)
		defer installationTmpFile.Close()
		io.Copy(installationTmpFile, installationSettings)
		config := cfbackup.NewConfigurationParser(installationTmpFile.Name())

		if iaas, hasKey := config.GetIaaS(); hasKey {
			sshKey = iaas.SSHPrivateKey
		}
		elasticRuntime = NewElasticRuntime(installationTmpFile.Name(), tileSpec.ArchiveDirectory, sshKey)
	}
	return
}

//GetInstallationSettings - makes a call to ops manager and returns a io.reader containing the contents of the installation settings file.
var GetInstallationSettings = func(tileSpec tileregistry.TileSpec) (settings io.Reader, err error) {
	var (
		opsManager *opsmanager.OpsManager
	)

	if opsManager, err = opsmanager.NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory); err == nil {
		settings, err = opsManager.GetInstallationSettings()
	}
	return
}
