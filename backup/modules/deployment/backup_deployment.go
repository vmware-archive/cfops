package deployment

import (
	"os"
	"path"
	"strings"

	"github.com/pivotalservices/cfops/ssh"
)

func New(deploymentDir string) *Deployment {
	return &Deployment{
		deploymentDir: deploymentDir,
	}
}

type Deployment struct {
	deploymentDir string
}

func (context *Deployment) Backup(copier ssh.Copier) error {
	file, _ := os.Create(path.Join(context.deploymentDir, "deployments.tar.gz"))
	defer file.Close()
	command := "cd /var/tempest/workspaces/default && tar cz deployments"

	err := copier.Copy(file, strings.NewReader(command))
	return err
}
