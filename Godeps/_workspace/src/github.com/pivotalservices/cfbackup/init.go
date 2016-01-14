package cfbackup

import "github.com/pivotalservices/cfops/tileregistry"

func init() {
	SetPGDumpUtilVersions()
	tileregistry.Register("ops-manager", new(OpsManagerBuilder))
	tileregistry.Register("elastic-runtime", new(ElasticRuntimeBuilder))
}
