package api

import (
	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// DirectorRepository is the interface for accessing a BOSH director
type DirectorRepository interface {
	GetInfo() (directorInfo models.DirectorInfo, apiResponse net.ApiResponse)

	GetStemcells() (stemcells models.Stemcells, apiResponse net.ApiResponse)
	DeleteStemcell(name string, version string) (apiResponse net.ApiResponse)

	GetReleases() (releases models.Releases, apiResponse net.ApiResponse)
	DeleteReleases(name string) (apiResponse net.ApiResponse)
	DeleteRelease(name string, version string) (apiResponse net.ApiResponse)

	GetDeployments() (deployments models.Deployments, apiResponse net.ApiResponse)
	GetDeploymentManifest(deploymentName string) (manifest *models.DeploymentManifest, apiResponse net.ApiResponse)
	DeleteDeployment(deploymentName string) (apiResponse net.ApiResponse)
	ListDeploymentVMs(deploymentName string) (deploymentVMs []models.DeploymentVM, apiResponse net.ApiResponse)
	FetchVMsStatus(deploymentName string) (vmsStatus []models.VMStatus, apiResponse net.ApiResponse)

	GetTaskStatuses() (task []models.TaskStatus, apiResponse net.ApiResponse)
	GetTaskStatus(taskID int) (task models.TaskStatus, apiResponse net.ApiResponse)
}

// BoshDirectorRepository represents a Director
type BoshDirectorRepository struct {
	config  *models.Director
	gateway net.Gateway
}

// NewBoshDirectorRepository is a constructor for a BoshDirectorRepository
func NewBoshDirectorRepository(config *models.Director, gateway net.Gateway) (repo BoshDirectorRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}
