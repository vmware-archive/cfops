package cfbackup

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/gtils/command"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/xchapter7x/lo"
)

// Ops Manager Backup constants
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

type httpUploader func(conn ghttp.ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error)

type httpRequestor interface {
	Get(ghttp.HttpRequestEntity) ghttp.RequestAdaptor
	Post(ghttp.HttpRequestEntity, io.Reader) ghttp.RequestAdaptor
	Put(ghttp.HttpRequestEntity, io.Reader) ghttp.RequestAdaptor
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
	SSHPrivateKey       string
	SSHUsername         string
	SSHPassword         string
	SSHPort             int
}

// NewOpsManager initializes an OpsManager instance
var NewOpsManager = func(opsManagerHostname string, adminUsername string, adminPassword string, opsManagerUsername string, opsManagerPassword string, target string) (context *OpsManager, err error) {
	backupContext := NewBackupContext(target, cfenv.CurrentEnv())
	settingsHTTPRequestor := ghttp.NewHttpGateway()
	settingsMultiHTTPRequestor := GetUploader(backupContext)
	assetsHTTPRequestor := ghttp.NewHttpGateway()
	assetsMultiHTTPRequestor := GetUploader(backupContext)

	context = &OpsManager{
		SettingsUploader:    settingsMultiHTTPRequestor,
		AssetsUploader:      assetsMultiHTTPRequestor,
		SettingsRequestor:   settingsHTTPRequestor,
		AssetsRequestor:     assetsHTTPRequestor,
		DeploymentDir:       path.Join(target, OpsMgrBackupDir, OpsMgrDeploymentsDir),
		Hostname:            opsManagerHostname,
		Username:            adminUsername,
		Password:            adminPassword,
		BackupContext:       backupContext,
		LocalExecuter:       command.NewLocalExecuter(),
		OpsmanagerBackupDir: OpsMgrBackupDir,
		SSHUsername:         opsManagerUsername,
		SSHPassword:         opsManagerPassword,
		SSHPort:             OpsMgrDefaultSSHPort,
	}
	err = context.createExecuter()
	return
}

//SetSSHPrivateKey - sets the private key in the ops manager object and rebuilds the remote executer associated with the opsmanager
func (context *OpsManager) SetSSHPrivateKey(key string) {
	context.SSHPrivateKey = key
	context.createExecuter()
}

func (context *OpsManager) createExecuter() (err error) {
	context.Executer, err = command.NewRemoteExecutor(command.SshConfig{
		Username: context.SSHUsername,
		Password: context.SSHPassword,
		Host:     context.Hostname,
		Port:     context.SSHPort,
		SSLKey:   context.SSHPrivateKey,
	})
	return
}

// GetInstallationSettings retrieves all the installation settings from OpsMan
// and returns them in a buffered reader
func (context *OpsManager) GetInstallationSettings() (settings io.Reader, err error) {
	var bytesBuffer = new(bytes.Buffer)
	url := fmt.Sprintf(OpsMgrInstallationSettingsURL, context.Hostname)
	lo.G.Debug(fmt.Sprintf("Exporting url '%s'", url))

	if err = context.saveHTTPResponse(url, bytesBuffer); err == nil {
		settings = bytesBuffer
	}
	return
}

//~ Backup Operations

// Backup performs a backup of a Pivotal Ops Manager instance
func (context *OpsManager) Backup() (err error) {
	if err = context.saveDeployments(); err == nil {
		err = context.saveInstallation()
	}
	return
}

func (context *OpsManager) saveDeployments() (err error) {
	var backupWriter io.WriteCloser
	if backupWriter, err = context.Writer(context.TargetDir, context.OpsmanagerBackupDir, OpsMgrDeploymentsFileName); err == nil {
		defer backupWriter.Close()
		command := "cd /var/tempest/workspaces/default && tar cz deployments"
		err = context.Executer.Execute(backupWriter, command)
	}
	return
}

func (context *OpsManager) saveInstallation() error {
	return context.saveInstallationSettingsAndAssets()
}

func (context *OpsManager) saveInstallationSettingsAndAssets() (err error) {
	if err = context.exportFile(OpsMgrInstallationSettingsURL, OpsMgrInstallationSettingsFilename); err == nil {
		err = context.exportFile(OpsMgrInstallationAssetsURL, OpsMgrInstallationAssetsFileName)
	}
	return
}

func (context *OpsManager) exportFile(urlFormat string, filename string) (err error) {
	url := fmt.Sprintf(urlFormat, context.Hostname)

	lo.G.Debug("Exporting file", log.Data{"url": url, "filename": filename})
	var backupWriter io.WriteCloser

	if backupWriter, err = context.Writer(context.TargetDir, context.OpsmanagerBackupDir, filename); err == nil {
		defer backupWriter.Close()
		err = context.saveHTTPResponse(url, backupWriter)
	}
	return
}

func (context *OpsManager) saveHTTPResponse(url string, dest io.Writer) (err error) {
	requestor := context.SettingsRequestor
	resp, err := requestor.Get(ghttp.HttpRequestEntity{
		Url:         url,
		Username:    context.Username,
		Password:    context.Password,
		ContentType: "application/octet-stream",
	})()

	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		_, err = io.Copy(dest, resp.Body)

	} else if resp != nil && resp.StatusCode != http.StatusOK {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(errMsg[:]))
	}

	if err != nil {
		lo.G.Error("error in save http request", err)
	}
	return
}

//~ Restore Operations

// Restore performs a restore of a Pivotal Ops Manager instance
func (context *OpsManager) Restore() (err error) {
	err = context.importInstallation()
	return
}

func (context *OpsManager) importInstallation() (err error) {
	defer func() {
		if err == nil {
			lo.G.Debug("removing deployment files")
			err = context.removeExistingDeploymentFiles()
		}
	}()
	installAssetsURL := fmt.Sprintf(OpsMgrInstallationAssetsURL, context.Hostname)
	lo.G.Debug("uploading installation assets installAssetsURL: %s", installAssetsURL)
	err = context.importInstallationPart(installAssetsURL, OpsMgrInstallationAssetsFileName, OpsMgrInstallationAssetsPostFieldName, context.AssetsUploader)
	return
}

func (context *OpsManager) importInstallationPart(url, filename, fieldname string, upload httpUploader) (err error) {
	var backupReader io.ReadCloser

	if backupReader, err = context.Reader(context.TargetDir, context.OpsmanagerBackupDir, filename); err == nil {
		defer backupReader.Close()
		var res *http.Response
		conn := ghttp.ConnAuth{
			Url:      url,
			Username: context.Username,
			Password: context.Password,
		}

		filePath := path.Join(context.TargetDir, context.OpsmanagerBackupDir, filename)
		bufferedReader := bufio.NewReader(backupReader)

		lo.G.Debug("upload request", log.Data{"fieldname": fieldname, "filePath": filePath, "conn": conn})

		if res, err = upload(conn, fieldname, filePath, -1, bufferedReader, nil); err != nil {
			err = fmt.Errorf(fmt.Sprintf("ERROR:%s", err.Error()))
			lo.G.Debug("upload failed", log.Data{"err": err, "response": res})
		}
	}
	return
}

func (context *OpsManager) removeExistingDeploymentFiles() (err error) {
	var w bytes.Buffer
	command := fmt.Sprintf("if [ -f %s ]; then sudo rm %s;fi", OpsMgrDeploymentsFile, OpsMgrDeploymentsFile)
	err = context.Executer.Execute(&w, command)
	return
}
