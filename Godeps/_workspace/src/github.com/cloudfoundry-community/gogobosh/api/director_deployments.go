package api

import (
	"fmt"
	"net/url"
	"time"

	"launchpad.net/goyaml"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// GetDeployments returns a list of deployments, and the releases/stemcells being used
func (repo BoshDirectorRepository) GetDeployments() (deployments models.Deployments, apiResponse net.ApiResponse) {
	response := []deploymentResponse{}

	path := "/deployments"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &response)
	if apiResponse.IsNotSuccessful() {
		return
	}

	list := []*models.Deployment{}
	for _, resource := range response {
		list = append(list, resource.ToModel())
	}
	deployments = models.Deployments(list)

	return
}

// GetDeploymentManifest returns a deployment manifest
func (repo BoshDirectorRepository) GetDeploymentManifest(deploymentName string) (manifest *models.DeploymentManifest, apiResponse net.ApiResponse) {
	deploymentManifestResponse := deploymentManifestResponse{}

	path := fmt.Sprintf("/deployments/%s", deploymentName)
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &deploymentManifestResponse)
	if apiResponse.IsNotSuccessful() {
		return
	}

	return deploymentManifestResponse.ToModel(), apiResponse
}

// DeleteDeployment asks the director to delete a deployment
func (repo BoshDirectorRepository) DeleteDeployment(deploymentName string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/deployments/%s?force=true", deploymentName)
	apiResponse = repo.gateway.DeleteResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password)
	if apiResponse.IsNotSuccessful() {
		return
	}
	if !apiResponse.IsRedirection() {
		return
	}

	var taskStatus models.TaskStatus
	taskURL, err := url.Parse(apiResponse.RedirectLocation)
	if err != nil {
		return
	}

	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+taskURL.Path, repo.config.Username, repo.config.Password, &taskStatus)
	if apiResponse.IsNotSuccessful() {
		return
	}

	/* Progression should be: queued, progressing, done */
	/* TODO task might fail; end states: done, error, cancelled */
	for taskStatus.State != "done" {
		time.Sleep(1)
		taskStatus, apiResponse = repo.GetTaskStatus(taskStatus.ID)
		if apiResponse.IsNotSuccessful() {
			return
		}
	}

	return
}

type deploymentResponse struct {
	Name      string        `json:"name"`
	Releases  []nameVersion `json:"releases"`
	Stemcells []nameVersion `json:"stemcells"`
}

type deploymentManifestResponse struct {
	RawManifest string `json:"manifest"`
}

type nameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (resource deploymentResponse) ToModel() (deployment *models.Deployment) {
	deployment = &models.Deployment{}
	deployment.Name = resource.Name
	for _, releaseResponse := range resource.Releases {
		release := models.NameVersion{}
		release.Name = releaseResponse.Name
		release.Version = releaseResponse.Version

		deployment.Releases = append(deployment.Releases, release)
	}

	for _, stemcellResponse := range resource.Stemcells {
		stemcell := models.NameVersion{}
		stemcell.Name = stemcellResponse.Name
		stemcell.Version = stemcellResponse.Version

		deployment.Stemcells = append(deployment.Stemcells, stemcell)
	}
	return
}

// ToModel converts a GetDeploymentManifest API response into models.DeploymentManifest
func (resource deploymentManifestResponse) ToModel() (manifest *models.DeploymentManifest) {
	manifest = &models.DeploymentManifest{}
	goyaml.Unmarshal([]byte(resource.RawManifest), manifest)
	return
}
