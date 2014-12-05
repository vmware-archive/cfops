package backup

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	DbEncryptionKey string
	BackupContext
}

// NewElasticRuntime initializes an ElasticRuntime intance
func NewElasticRuntime(dbEncryptionKey string, target string) *ElasticRuntime {
	context := &ElasticRuntime{
		DbEncryptionKey: dbEncryptionKey,
		BackupContext: BackupContext{
			TargetDir: target,
		},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() error {
	return nil
}
