package cfbackup

import (
	"github.com/pivotalservices/cfbackup/tiles/elasticruntime"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
)

func init() {
	tileregistry.Register("ops-manager", new(opsmanager.OpsManagerBuilder))
	tileregistry.Register("elastic-runtime", new(elasticruntime.ElasticRuntimeBuilder))
}
