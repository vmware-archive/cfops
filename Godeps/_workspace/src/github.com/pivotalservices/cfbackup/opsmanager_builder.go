package cfbackup

import "github.com/pivotalservices/cfops/tileregistry"

//OpsManagerBuilder - an object that can build ops manager objects
type OpsManagerBuilder struct{}

//New -- builds a new ops manager object pre initialized
func (s *OpsManagerBuilder) New(tileSpec tileregistry.TileSpec) (opsManager tileregistry.Tile, err error) {
	opsManager, err = NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory)
	return
}
