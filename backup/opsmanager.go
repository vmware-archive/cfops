package backup

import (
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
	copier := ssh.New("tempest", context.TempestPassword, context.Hostname, 22)
	err := context.copyDeployments(copier)
	if err != nil {
		// TODO: Log
		return err
	}
	err = context.extractDbEncryptionKey()
	if err != nil {
		// TODO: Log
	}
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

func (context *OpsManager) extractDbEncryptionKey() error {
	return nil
}
