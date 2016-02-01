package plugin

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
		lo.G.Debug("we have some settings", settingsReader)
		installationSettings := cfbackup.NewConfigurationParserFromReader(settingsReader)
		pcf := DefaultPivotalCF(installationSettings)
		tile = &PluginTile{
			PivotalCF: pcf,
			Meta:      s.Meta,
			FilePath:  s.FilePath,
		}
	}
	return
}
