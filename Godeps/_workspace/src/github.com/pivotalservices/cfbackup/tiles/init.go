package cfbackup

import (
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tiles/elasticruntime"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
)

func init() {
	cfbackup.SetPGDumpUtilVersions()
	tileregistry.Register("ops-manager", new(opsmanager.OpsManagerBuilder))
	tileregistry.Register("elastic-runtime", new(elasticruntime.ElasticRuntimeBuilder))
}
