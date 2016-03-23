package opsmanager

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	urllib "net/url"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/gtils/command"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/pivotalservices/gtils/uaa"
	"github.com/xchapter7x/lo"
)

// NewOpsManager initializes an OpsManager instance
var NewOpsManager = func(opsManagerHostname string, adminUsername string, adminPassword string, opsManagerUsername string, opsManagerPassword string, target string, cryptKey string) (context *OpsManager, err error) {
	backupContext := cfbackup.NewBackupContext(target, cfenv.CurrentEnv(), cryptKey)
	settingsHTTPRequestor := ghttp.NewHttpGateway()
	settingsMultiHTTPRequestor := httpUploader(cfbackup.GetUploader(backupContext))
	assetsHTTPRequestor := ghttp.NewHttpGateway()
	assetsMultiHTTPRequestor := httpUploader(cfbackup.GetUploader(backupContext))

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
		ClearBoshManifest:   false,
	}
	err = context.createExecuter()
	return
}

//SetSSHPrivateKey - sets the private key in the ops manager object and rebuilds the remote executer associated with the opsmanager
func (context *OpsManager) SetSSHPrivateKey(key string) {
	lo.G.Debug("Setting SSHKey")
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
	var resp *http.Response
	lo.G.Debug("attempting to auth against", url)

	if resp, err = context.oauthHTTPGet(url); err != nil {
		lo.G.Info("falling back to basic auth for legacy system", err)
		resp, err = context.legacyHTTPGet(url)
	}

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

func (context *OpsManager) legacyHTTPGet(url string) (resp *http.Response, err error) {
	requestor := context.SettingsRequestor
	resp, err = requestor.Get(ghttp.HttpRequestEntity{
		Url:         url,
		Username:    context.Username,
		Password:    context.Password,
		ContentType: "application/octet-stream",
	})()
	lo.G.Debug("called basic auth on legacy ops manager", url, err)
	return
}

func (context *OpsManager) oauthHTTPGet(urlString string) (resp *http.Response, err error) {
	var token string
	var uaaURL, _ = urllib.Parse(urlString)
	var opsManagerUsername = context.Username
	var opsManagerPassword = context.Password
	var clientID = "opsman"
	var clientSecret = ""
	lo.G.Debug("aquiring your token from: ", uaaURL, urlString)

	if token, err = uaa.GetToken("https://"+uaaURL.Host+"/uaa", opsManagerUsername, opsManagerPassword, clientID, clientSecret); err == nil {
		lo.G.Debug("your token", token, "https://"+uaaURL.Host+"/uaa")
		requestor := context.SettingsRequestor
		resp, err = requestor.Get(ghttp.HttpRequestEntity{
			Url:           urlString,
			ContentType:   "application/octet-stream",
			Authorization: "Bearer " + token,
		})()
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
		if err == nil && context.ClearBoshManifest {
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
		var resp *http.Response
		var uaaURL, _ = urllib.Parse(url)
		var clientID = "opsman"
		var clientSecret = ""
		var token, _ = uaa.GetToken("https://"+uaaURL.Host+"/uaa", context.Username, context.Password, clientID, clientSecret)
		conn := ghttp.ConnAuth{
			Url:         url,
			Username:    context.Username,
			Password:    context.Password,
			BearerToken: token,
		}
		filePath := path.Join(context.TargetDir, context.OpsmanagerBackupDir, filename)
		bufferedReader := bufio.NewReader(backupReader)
		lo.G.Debug("upload request", log.Data{"fieldname": fieldname, "filePath": filePath, "conn": conn})
		creds := map[string]string{
			"password": context.Password,
		}
		resp, err = upload(conn, fieldname, filePath, -1, bufferedReader, creds)

		if err == nil && resp.StatusCode == http.StatusOK {
			lo.G.Debug("Request for %s succeeded with status: %s", url, resp.Status)

		} else if resp != nil && resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("Request for %s failed with status: %s", url, resp.Status)
		}

		if err != nil {
			lo.G.Error("error uploading installation: %s", err.Error())
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
