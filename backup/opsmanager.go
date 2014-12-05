package backup

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pivotalservices/cfops/osutils"
	"github.com/pivotalservices/cfops/ssh"
)

// OpsManager contains the location and credentials of a Pivotal Ops Manager instance
type OpsManager struct {
	BackupContext
	Hostname        string
	Username        string
	Password        string
	TempestPassword string
	DbEncryptionKey string
}

// Backup performs a backup of a Pivotal Ops Manager instance
func (context *OpsManager) Backup() error {
	// Step 1: Download Ops Manager Installation Settings
	// :::TODO:::
	// CONNECTION_URL=https://$OPS_MANAGER_HOST/api/installation_settings
	// echo "EXPORT INSTALLATION FILES FROM " $CONNECTION_URL
	// curl "$CONNECTION_URL" -X GET -u $OPS_MGR_ADMIN_USERNAME:$OPS_MGR_ADMIN_PASSWORD --insecure -k -o $WORK_DIR/installation.yml

	// Step 2: Export Ops Manager Deployments (http://docs.pivotal.io/pivotalcf/customizing/backup-settings.html#export)
	copier := ssh.New("tempest", context.TempestPassword, context.Hostname, 22)
	err := context.copyDeployments(copier)
	if err != nil {
		// TODO: Log
	}

	// Step 3 (Optional): Export Ops Manager Installation
	// CONNECTION_URL=https://$OPS_MANAGER_HOST/api/installation_asset_collection
	// echo "EXPORT INSTALLATION FILES FROM " $CONNECTION_URL
	// curl "$CONNECTION_URL" -X GET -u $OPS_MGR_ADMIN_USERNAME:$OPS_MGR_ADMIN_PASSWORD --insecure -k -o $WORK_DIR/installation.zip

	return err
}

// NewOpsManager initializes an OpsManager instance
func NewOpsManager(hostname string, username string, password string, tempestpassword string, target string) *OpsManager {
	context := &OpsManager{
		Hostname:        hostname,
		Username:        username,
		Password:        password,
		TempestPassword: tempestpassword,
		BackupContext: BackupContext{
			TargetDir: target,
		},
	}
	return context
}

func (context *OpsManager) copyDeployments(copier ssh.Copier) error {
	file, err := osutils.SafeCreate(context.TargetDir, "opsmanager", "deployments.tar.gz")
	defer file.Close()
	if err != nil {
		return err
	}

	command := "cd /var/tempest/workspaces/default && tar cz deployments"
	return copier.Copy(file, strings.NewReader(command))
}

func (context *OpsManager) exportInstallationSettings() {
	connectionURL := "https://" + context.Hostname + "/api/installation_settings"

	resp, err := invoke("GET", connectionURL, context.Username, context.Password, false)
	if err != nil {
		// TODO: Log
	}
	defer resp.Body.Close()
	f, err := osutils.SafeCreate(context.TargetDir, "opsmanager", "installation.yml")
	defer f.Close()
	_, e := io.Copy(f, resp.Body)
	if e != nil {
		// TODO: Log
	}
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
