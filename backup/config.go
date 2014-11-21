package backup

type BackupConfig struct {
	OpsManagerHost          string
	TempestPassword         string
	OpsManagerAdminUser     string
	OpsManagerAdminPassword string
	BackupFileLocation      string
}
