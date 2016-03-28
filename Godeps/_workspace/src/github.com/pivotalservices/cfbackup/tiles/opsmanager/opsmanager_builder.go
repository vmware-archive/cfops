package opsmanager

import (
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

//New -- builds a new ops manager object pre initialized
func (s *OpsManagerBuilder) New(tileSpec tileregistry.TileSpec) (opsManagerTile tileregistry.Tile, err error) {
	var opsManager *OpsManager
	opsManager, err = NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory, tileSpec.CryptKey)
	opsManager.ClearBoshManifest = tileSpec.ClearBoshManifest

	if installationSettings, err := opsManager.GetInstallationSettings(); err == nil {
		config := cfbackup.NewConfigurationParserFromReader(installationSettings)

		if iaas, hasKey := config.GetIaaS(); hasKey {
			lo.G.Debug("we found a iaas info block")
			opsManager.SetSSHPrivateKey(iaas.SSHPrivateKey)

		} else {
			lo.G.Debug("No IaaS PEM key found. Defaulting to using ssh username and password credentials")
		}
	}
	opsManagerTile = opsManager
	return
}
