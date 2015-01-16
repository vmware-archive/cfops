package backup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	cfhttp "github.com/pivotalservices/cfops/backup/modules/http"
	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

const (
	OPSMGR_INSTALLATION_SETTINGS_FILENAME string = "installation.json"
	OPSMGR_INSTALLATION_ASSETS_FILENAME   string = "installation.zip"
	OPSMGR_DEPLOYMENTS_FILENAME           string = "deployments.tar.gz"
	OPSMGR_ENCRYPTIONKEY_FILENAME         string = "cc_db_encryption_key.txt"
	OPSMGR_BACKUP_DIR                     string = "opsmanager"
	OPSMGR_DEPLOYMENTS_DIR                string = "deployments"
	OPSMGR_DEFAULT_USER                   string = "tempest"
	OPSMGR_INSTALLATION_SETTINGS_URL      string = "https://%s/api/installation_settings"
	OPSMGR_INSTALLATION_ASSETS_URL        string = "https://%s/api/installation_asset_collection"
)

type httpUploader interface {
	Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error)
}

// OpsManager contains the location and credentials of a Pivotal Ops Manager instance
type OpsManager struct {
	BackupContext
	Hostname            string
	Username            string
	Password            string
	TempestPassword     string
	DbEncryptionKey     string
	RestRunner          RestAdapter
	Executer            command.Executer
	SettingsUploader    httpUploader
	AssetsUploader      httpUploader
	DeploymentDir       string
	OpsmanagerBackupDir string
}

// Backup performs a backup of a Pivotal Ops Manager instance
func (context *OpsManager) Backup() (err error) {
	if err = context.copyDeployments(); err == nil {
		err = context.exportAndExtract()
	}
	return
}

// Restore performs a restore of a Pivotal Ops Manager instance
func (context *OpsManager) Restore() (err error) {
	return
}

func (context *OpsManager) exportAndExtract() (err error) {
	if err = context.extract(); err == nil {
		err = context.export()
	}
	return
}

func (context *OpsManager) export() (err error) {

	if err = context.exportUrlToFile(OPSMGR_INSTALLATION_SETTINGS_URL, OPSMGR_INSTALLATION_SETTINGS_FILENAME); err == nil {
		err = context.exportUrlToFile(OPSMGR_INSTALLATION_ASSETS_URL, OPSMGR_INSTALLATION_ASSETS_FILENAME)
	}
	return
}

func (context *OpsManager) exportUrlToFile(urlFormat string, filename string) (err error) {
	var settingsFileRef *os.File
	defer settingsFileRef.Close()

	if settingsFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, filename); err == nil {
		err = context.exportUrlToWriter(urlFormat, settingsFileRef)
	}
	return
}

func (context *OpsManager) exportUrlToWriter(urlFormat string, dest io.Writer) (err error) {
	var body io.Reader
	connectionURL := fmt.Sprintf(urlFormat, context.Hostname)

	if _, body, err = context.RestRunner.Run("GET", connectionURL, context.Username, context.Password, false); err == nil {
		_, err = io.Copy(dest, body)
	}
	return
}

func (context *OpsManager) extract() (err error) {
	var keyFileRef *os.File
	defer keyFileRef.Close()

	if keyFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, OPSMGR_ENCRYPTIONKEY_FILENAME); err == nil {
		err = ExtractEncryptionKey(keyFileRef, context.DeploymentDir)
	}
	return
}

// NewOpsManager initializes an OpsManager instance
func NewOpsManager(hostname string, username string, password string, tempestpassword string, target string) (context *OpsManager, err error) {
	var remoteExecuter command.Executer

	if remoteExecuter, err = createExecuter(hostname, tempestpassword); err == nil {
		settingsGateway, assetsGateway := createInstallationGateways(hostname, tempestpassword)
		context = &OpsManager{
			SettingsUploader: settingsGateway,
			AssetsUploader:   assetsGateway,
			DeploymentDir:    path.Join(target, OPSMGR_BACKUP_DIR, OPSMGR_DEPLOYMENTS_DIR),
			Hostname:         hostname,
			Username:         username,
			Password:         password,
			BackupContext: BackupContext{
				TargetDir: target,
			},
			RestRunner:          RestAdapter(invoke),
			Executer:            remoteExecuter,
			OpsmanagerBackupDir: OPSMGR_BACKUP_DIR,
		}
	}
	return
}

func createExecuter(hostname, tempestpassword string) (remoteExecuter command.Executer, err error) {
	remoteExecuter, err = command.NewRemoteExecutor(command.SshConfig{
		Username: OPSMGR_DEFAULT_USER,
		Password: tempestpassword,
		Host:     hostname,
		Port:     22,
	})
	return
}

func createInstallationGateways(hostname, tempestpassword string) (settingsGateway, assetsGateway *cfhttp.HttpGateway) {
	settingsURL := fmt.Sprintf(OPSMGR_INSTALLATION_SETTINGS_URL, hostname)
	assetsURL := fmt.Sprintf(OPSMGR_INSTALLATION_ASSETS_URL, hostname)
	settingsGateway = cfhttp.NewHttpGateway(settingsURL, OPSMGR_DEFAULT_USER, tempestpassword, "application/octet-stream", nil)
	assetsGateway = cfhttp.NewHttpGateway(assetsURL, OPSMGR_DEFAULT_USER, tempestpassword, "application/octet-stream", nil)
	return
}

func (context *OpsManager) copyDeployments() (err error) {
	var file *os.File
	defer file.Close()

	if file, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, OPSMGR_DEPLOYMENTS_FILENAME); err == nil {
		command := "cd /var/tempest/workspaces/default && tar cz deployments"
		err = context.Executer.Execute(file, command)
	}
	return
}
