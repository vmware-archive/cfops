package elasticruntime

import (
	"github.com/pivotalservices/cfbackup"
	ghttp "github.com/pivotalservices/gtils/http"
)

type (

	//ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
	ElasticRuntime struct {
		cfbackup.BackupContext
		JSONFile          string
		SystemsInfo       cfbackup.SystemsInfo
		PersistentSystems []cfbackup.SystemDump
		HTTPGateway       ghttp.HttpGateway
		InstallationName  string
		SSHPrivateKey     string
	}

	//ElasticRuntimeBuilder -- an object that can build an elastic runtime pre-initialized
	ElasticRuntimeBuilder struct{}
)
