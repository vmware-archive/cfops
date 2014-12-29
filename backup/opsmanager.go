package backup

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

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

func (context *OpsManager) exportAndExtract() (err error) {
	if err = context.extract(); err == nil {
		err = context.export()
	}
	return
}

func (context *OpsManager) export() (err error) {
	var settingsFileRef *os.File
	defer settingsFileRef.Close()

	if settingsFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, "installation.yml"); err == nil {
		err = context.exportInstallationSettings(settingsFileRef)
	}
	return
}

func (context *OpsManager) extract() (err error) {
	var keyFileRef *os.File
	defer keyFileRef.Close()

	if keyFileRef, err = osutils.SafeCreate(context.TargetDir, context.OpsmanagerBackupDir, "cc_db_encryption_key.txt"); err == nil {
		err = ExtractEncryptionKey(keyFileRef, context.DeploymentDir)
	}
	return
}

// NewOpsManager initializes an OpsManager instance
func NewOpsManager(hostname string, username string, password string, tempestpassword string, target string) (context *OpsManager, err error) {
	var remoteExecuter command.Executer
	remoteExecuter, err = command.NewRemoteExecutor(command.SshConfig{
		Username: "tempest",
		Password: tempestpassword,
		Host:     hostname,
		Port:     22,
	})

	if err == nil {
		context = &OpsManager{
			DeploymentDir: path.Join(target, "opsmanager", "deployments"),
			Hostname:      hostname,
			Username:      username,
			Password:      password,
			BackupContext: BackupContext{
				TargetDir: target,
			},
			RestRunner:          RestAdapter(invoke),
			Executer:            remoteExecuter,
			OpsmanagerBackupDir: "opsmanager",
		}
	}
	return
}

func (context *OpsManager) copyDeployments() (err error) {
	var file *os.File
	defer file.Close()

	if file, err = osutils.SafeCreate(context.TargetDir, "opsmanager", "deployments.tar.gz"); err == nil {
		command := "cd /var/tempest/workspaces/default && tar cz deployments"
		err = context.Executer.Execute(file, command)
	}
	return
}

func (context *OpsManager) exportInstallationSettings(dest io.Writer) (err error) {
	var body io.Reader
	connectionURL := fmt.Sprintf("https://%s/api/installation_settings", context.Hostname)

	if _, body, err = context.RestRunner.Run("GET", connectionURL, context.Username, context.Password, false); err == nil {
		_, err = io.Copy(dest, body)
	}
	return
}

func invoke(method string, connectionURL string, username string, password string, isYaml bool) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest(method, connectionURL, nil)
	req.SetBasicAuth(username, password)

	if isYaml {
		req.Header.Set("Content-Type", "text/yaml")
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	return resp, err
}
