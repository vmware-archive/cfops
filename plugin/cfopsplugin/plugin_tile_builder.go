package cfopsplugin

import (
	"io"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

//New - method to create a plugin tile
func (s *PluginTileBuilder) New(tileSpec tileregistry.TileSpec) (tile tileregistry.Tile, err error) {
	var opsManager *opsmanager.OpsManager
	var settingsReader io.Reader
	opsManager, err = opsmanager.NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory)

	if settingsReader, err = opsManager.GetInstallationSettings(); err == nil {
		var brPlugin BackupRestorer
		installationSettings := cfbackup.NewConfigurationParserFromReader(settingsReader)
		pcf := NewPivotalCF(installationSettings, tileSpec)
		lo.G.Debug("", s.Meta.Name, s.FilePath, pcf)
		brPlugin, _ = Call(s.Meta.Name, s.FilePath)
		brPlugin.Setup(pcf)
		tile = brPlugin
	}
	lo.G.Debug("error from getinstallationsettings: ", err)
	return
}
