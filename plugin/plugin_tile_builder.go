package plugin

import (
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
)

//New - method to create a plugin tile
func (s *PluginTileBuilder) New(tileSpec tileregistry.TileSpec) (tile tileregistry.Tile, err error) {
	var opsManager *opsmanager.OpsManager
	opsManager, err = opsmanager.NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory)
	return &PluginTile{
		OpsManager: opsManager,
		Meta:       s.Meta,
		FilePath:   s.FilePath,
	}, err
}
