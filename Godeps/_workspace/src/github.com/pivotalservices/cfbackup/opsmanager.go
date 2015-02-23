package cfbackup

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/pivotalservices/gtils/osutils"
)

const (
	OPSMGR_INSTALLATION_SETTINGS_FILENAME       = "installation.json"
	OPSMGR_INSTALLATION_SETTINGS_POSTFIELD_NAME = "installation[file]"
	OPSMGR_INSTALLATION_ASSETS_FILENAME         = "installation.zip"
	OPSMGR_INSTALLATION_ASSETS_POSTFIELD_NAME   = "installation[file]"
	OPSMGR_DEPLOYMENTS_FILENAME                 = "deployments.tar.gz"
	OPSMGR_ENCRYPTIONKEY_FILENAME               = "cc_db_encryption_key.txt"
	OPSMGR_BACKUP_DIR                           = "opsmanager"
	OPSMGR_DEPLOYMENTS_DIR                      = "deployments"
	OPSMGR_DEFAULT_USER                         = "tempest"
	OPSMGR_INSTALLATION_SETTINGS_URL            = "https://%s/api/installation_settings"
	OPSMGR_INSTALLATION_ASSETS_URL              = "https://%s/api/installation_asset_collection"
	OPSMGR_DEPLOYMENTS_FILE                     = "/var/tempest/workspaces/default/deployments/bosh-deployments.yml"
	OPSMGR_MULTIPART_FORM_CONTENT_TYPE          = "multipart/form-data"
	OPSMGR_URL_FORM_CONTENT_TYPE                = "application/x-www-form-urlencoded"
)

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
	SettingsUploader    UploadFunc
	AssetsUploader      MultiPartBodyFunc
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
		assetsHttpRequestor := NewHttpGateway()

		context = &OpsManager{
			SettingsUploader:  Upload,
			AssetsUploader:    MultiPartBody,
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

	if err = context.importInstallationSettings(OPSMGR_INSTALLATION_SETTINGS_URL, OPSMGR_MULTIPART_FORM_CONTENT_TYPE,
		OPSMGR_INSTALLATION_SETTINGS_FILENAME, OPSMGR_INSTALLATION_SETTINGS_POSTFIELD_NAME, context.SettingsUploader); err == nil {
		// err = context.importInstallationAssets(OPSMGR_INSTALLATION_ASSETS_URL, OPSMGR_URL_FORM_CONTENT_TYPE,
		// 	OPSMGR_INSTALLATION_ASSETS_FILENAME, OPSMGR_INSTALLATION_ASSETS_POSTFIELD_NAME, context.AssetsRequestor)
	}
	return
}

func (context *OpsManager) removeExistingDeploymentFiles() (err error) {
	var w bytes.Buffer
	command := fmt.Sprintf("if [ -f %s ]; then sudo rm %s;fi", OPSMGR_DEPLOYMENTS_FILE, OPSMGR_DEPLOYMENTS_FILE)
	context.Logger.Debug("Removing bosh-deployments.yml")
	err = context.Executer.Execute(&w, command)
	return
}

func (context *OpsManager) importInstallationSettings(urlFormat, contentType, filename, fieldname string, upload UploadFunc) (err error) {

	filePath := path.Join(context.TargetDir, context.OpsmanagerBackupDir, filename)

	url := fmt.Sprintf(urlFormat, context.Hostname)

	context.Logger.Debug("Importing opsmanager installation", log.Data{"url": url, "filePath": filePath, "filename": filename, "fieldname": fieldname})

	if fileRef, err := os.Open(filePath); err == nil {

		resp, err := upload(url, context.Username, context.Password, fieldname, filename, fileRef, nil)
		if err != nil {
			context.Logger.Error("Error uploading settings", err)
			panic(err)
			// return
		}

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("import settings failed with status: %s", resp.Status)
		}
		context.Logger.Debug("Imported installation settings", log.Data{"StatusCode": resp.StatusCode, "Body": resp.Body})
	}
	return
}

func (context *OpsManager) importInstallationAssets(urlFormat, contentType, filename, fieldname string, requestor httpRequestor) (err error) {

	filePath := path.Join(context.TargetDir, context.OpsmanagerBackupDir, filename)

	assetsUrl := fmt.Sprintf(urlFormat, context.Hostname)

	context.Logger.Debug("Importing opsmanager installation", log.Data{"url": assetsUrl})

	body := strings.NewReader(url.QueryEscape(fmt.Sprint(fieldname + "=" + filePath)))

	context.Logger.Debug("Importing opsmanager installation", log.Data{"body": url.QueryEscape(fmt.Sprint(fieldname + "=" + filePath))})

	resp, err := requestor.Post(HttpRequestEntity{
		Url:         assetsUrl,
		Username:    context.Username,
		Password:    context.Password,
		ContentType: contentType,
	}, body)()
	if err != nil {
		context.Logger.Error("Error uploading assets", err)
		panic(err)
		// return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("import assets failed with status: %s", resp.Status)
	}
	context.Logger.Debug("Imported installation settings", log.Data{"StatusCode": resp.StatusCode, "Body": resp.Body})
	return
}

func (context *OpsManager) exportAndExtract() (err error) {
	if err = context.extract(); err == nil {
		err = context.export()
	}
	return
}

func (context *OpsManager) export() (err error) {

	if err = context.exportUrlToFile(OPSMGR_INSTALLATION_SETTINGS_URL, "application/json", OPSMGR_INSTALLATION_SETTINGS_FILENAME); err == nil {
		err = context.exportUrlToFile(OPSMGR_INSTALLATION_ASSETS_URL, "application/zip", OPSMGR_INSTALLATION_ASSETS_FILENAME)
	}
	return
}

func (context *OpsManager) exportUrlToFile(urlFormat, contentType string, filename string) (err error) {
	var settingsFileRef *os.File
	defer settingsFileRef.Close()

	url := fmt.Sprintf(urlFormat, context.Hostname)

	context.Logger.Debug("Exporting url '%s' to file '%s'", log.Data{"url": url, "filename": filename})

	if settingsFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, filename); err == nil {
		err = context.exportUrlToWriter(url, contentType, settingsFileRef, context.SettingsRequestor)
	}
	return
}

func (context *OpsManager) exportUrlToWriter(url, contentType string, dest io.Writer, requestor httpRequestor) (err error) {
	resp, err := requestor.Get(HttpRequestEntity{
		Url:         url,
		Username:    context.Username,
		Password:    context.Password,
		ContentType: contentType,
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
