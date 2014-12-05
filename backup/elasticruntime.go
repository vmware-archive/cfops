package backup

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	DeploymentsFile string
	DbEncryptionKey string
	BackupContext
}

// NewElasticRuntime initializes an ElasticRuntime intance
func NewElasticRuntime(deploymentsFile, dbEncryptionKey string, target string) *ElasticRuntime {
	context := &ElasticRuntime{
		DeploymentsFile: deploymentsFile,
		DbEncryptionKey: dbEncryptionKey,
		BackupContext: BackupContext{
			TargetDir: target,
		},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() error {
	// Step 1: Back Up the Cloud Controller DB Encryption Credentials (http://docs.pivotal.io/pivotalcf/customizing/backup-settings.html#encrypt-key)
	err := context.extractDbEncryptionKey()
	if err != nil {
		// TODO: Log
		return err
	}

	// Step 2: Find Job Names

	// Step 3: Stop Jobs
	// 4.1: cloud_controller-partition-default_az_guid
	// 4.2: cloud_controller_worker-partition-default_az_guid

	// Step 4: Backup Databases
	// Cloud Controller Database
	// UAA Database
	// Console Database
	// The NFS Server

	// Step 5: Start Jobs

	return nil
}

func (context *ElasticRuntime) extractDbEncryptionKey() error {
	return nil
}

func (context *ElasticRuntime) startJob(job string) error {
	return nil
}

func (context *ElasticRuntime) stopJob(job string) error {
	return nil
}
