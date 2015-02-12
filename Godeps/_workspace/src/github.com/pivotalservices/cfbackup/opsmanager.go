package cfbackup

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/pivotalservices/gtils/osutils"
)

const (
	OPSMGR_INSTALLATION_SETTINGS_FILENAME       string = "installation.json"
	OPSMGR_INSTALLATION_SETTINGS_POSTFIELD_NAME string = "installation[file]"
	OPSMGR_INSTALLATION_ASSETS_FILENAME         string = "installation.zip"
	OPSMGR_INSTALLATION_ASSETS_POSTFIELD_NAME   string = "installation[file]"
	OPSMGR_DEPLOYMENTS_FILENAME                 string = "deployments.tar.gz"
	OPSMGR_ENCRYPTIONKEY_FILENAME               string = "cc_db_encryption_key.txt"
	OPSMGR_BACKUP_DIR                           string = "opsmanager"
	OPSMGR_DEPLOYMENTS_DIR                      string = "deployments"
	OPSMGR_DEFAULT_USER                         string = "tempest"
	OPSMGR_INSTALLATION_SETTINGS_URL            string = "https://%s/api/installation_settings"
	OPSMGR_INSTALLATION_ASSETS_URL              string = "https://%s/api/installation_asset_collection"
)

type httpUploader func(string, string, io.Reader, map[string]string) (io.Reader, error)

type httpRequestor interface {
	Get(HttpRequestEntity) RequestAdaptor
	Post(HttpRequestEntity, io.Reader) RequestAdaptor
	Put(HttpRequestEntity, io.Reader) RequestAdaptor
}

// OpsManager contains the location and credentials of a Pivotal Ops Manager instance
type OpsManager struct {
	BackupContext
	Hostname            string
	Username            string
	Password            string
	TempestPassword     string
	DbEncryptionKey     string
	Executer            command.Executer
	LocalExecuter       command.Executer
	SettingsUploader    httpUploader
	AssetsUploader      httpUploader
	SettingsRequestor   httpRequestor
	AssetsRequestor     httpRequestor
	DeploymentDir       string
	OpsmanagerBackupDir string
	Logger              log.Logger
}

// NewOpsManager initializes an OpsManager instance
var NewOpsManager = func(hostname string, username string, password string, tempestpassword string, target string, logger log.Logger) (context *OpsManager, err error) {
	var remoteExecuter command.Executer

	if remoteExecuter, err = createExecuter(hostname, tempestpassword); err == nil {
		settingsHttpRequestor := NewHttpGateway()
		settingsMultiHttpRequestor := MultiPartBody
		assetsHttpRequestor := NewHttpGateway()
		assetsMultiHttpRequestor := MultiPartBody

		context = &OpsManager{
			SettingsUploader:  settingsMultiHttpRequestor,
			AssetsUploader:    assetsMultiHttpRequestor,
			SettingsRequestor: settingsHttpRequestor,
			AssetsRequestor:   assetsHttpRequestor,
			DeploymentDir:     path.Join(target, OPSMGR_BACKUP_DIR, OPSMGR_DEPLOYMENTS_DIR),
			Hostname:          hostname,
			Username:          username,
			Password:          password,
			BackupContext: BackupContext{
				TargetDir: target,
			},
			Executer:            remoteExecuter,
			LocalExecuter:       command.NewLocalExecuter(),
			OpsmanagerBackupDir: OPSMGR_BACKUP_DIR,
			Logger:              logger,
		}
	}
	return
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
	err = context.importInstallation()
	return
}

func (context *OpsManager) importInstallation() (err error) {
	defer func() {
		if err == nil {
			err = context.removeExistingDeploymentFiles()
		}
	}()

	if err = context.importInstallationPart(OPSMGR_INSTALLATION_SETTINGS_FILENAME, OPSMGR_INSTALLATION_SETTINGS_POSTFIELD_NAME, context.SettingsUploader); err == nil {
		err = context.importInstallationPart(OPSMGR_INSTALLATION_ASSETS_FILENAME, OPSMGR_INSTALLATION_ASSETS_POSTFIELD_NAME, context.AssetsUploader)
	}
	return
}

func (context *OpsManager) removeExistingDeploymentFiles() (err error) {
	var w bytes.Buffer
	command := "sudo rm /var/tempest/workspaces/default/deployments/bosh-deployments.yml"
	err = context.Executer.Execute(&w, command)
	return
}

func (context *OpsManager) importInstallationPart(filename, fieldname string, upload httpUploader) (err error) {
	filePath := path.Join(context.TargetDir, context.OpsmanagerBackupDir, filename)

	if fileRef, err := os.Open(filePath); err == nil {

		if res, err := upload(fieldname, filename, fileRef, nil); err == nil {
			err = fmt.Errorf(fmt.Sprintf("Bad Response from Gateway: %v", res))
		}
	}
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

	url := fmt.Sprintf(urlFormat, context.Hostname)

	context.Logger.Debug("Exporting url '%s' to file '%s'", log.Data{"url": url, "filename": filename})

	if settingsFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, filename); err == nil {
		err = context.exportUrlToWriter(url, settingsFileRef, context.SettingsRequestor)
	}
	return
}

func (context *OpsManager) exportUrlToWriter(url string, dest io.Writer, requestor httpRequestor) (err error) {
	resp, err := requestor.Get(HttpRequestEntity{
		Url:         url,
		Username:    context.Username,
		Password:    context.Password,
		ContentType: "application/octet-stream",
	})()
	if err == nil {
		defer resp.Body.Close()
		_, err = io.Copy(dest, resp.Body)
	}
	return
}

func (context *OpsManager) extract() (err error) {
	var keyFileRef *os.File
	defer keyFileRef.Close()
	context.Logger.Debug("Extracting Ops Manager")

	if keyFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, OPSMGR_ENCRYPTIONKEY_FILENAME); err == nil {
		context.Logger.Debug("Extracting encryption key")
		backupDir := path.Join(context.TargetDir, context.OpsmanagerBackupDir)
		deployment := path.Join(backupDir, OPSMGR_DEPLOYMENTS_FILENAME)
		cmd := "tar -xf " + deployment + " -C " + backupDir
		context.Logger.Debug("Extracting : %s", log.Data{"command": cmd})
		context.LocalExecuter.Execute(nil, cmd)

		err = ExtractEncryptionKey(keyFileRef, context.DeploymentDir)
	}
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

func createExecuter(hostname, tempestpassword string) (remoteExecuter command.Executer, err error) {
	remoteExecuter, err = command.NewRemoteExecutor(command.SshConfig{
		Username: OPSMGR_DEFAULT_USER,
		Password: tempestpassword,
		Host:     hostname,
		Port:     22,
	})
	return
}
