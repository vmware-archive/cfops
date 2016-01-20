package opsmanager

//OpsManager constants
const (
	OpsMgrInstallationSettingsFilename    string = "installation.json"
	OpsMgrInstallationAssetsFileName      string = "installation.zip"
	OpsMgrInstallationAssetsPostFieldName string = "installation[file]"
	OpsMgrDeploymentsFileName             string = "deployments.tar.gz"
	OpsMgrEncryptionKeyFileName           string = "cc_db_encryption_key.txt"
	OpsMgrBackupDir                       string = "opsmanager"
	OpsMgrDeploymentsDir                  string = "deployments"
	OpsMgrDefaultSSHPort                  int    = 22
	OpsMgrInstallationSettingsURL         string = "https://%s/api/installation_settings"
	OpsMgrInstallationAssetsURL           string = "https://%s/api/installation_asset_collection"
	OpsMgrDeploymentsFile                 string = "/var/tempest/workspaces/default/deployments/bosh-deployments.yml"
)
