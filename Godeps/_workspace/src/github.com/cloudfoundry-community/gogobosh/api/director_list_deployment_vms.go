package api

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

func (repo BoshDirectorRepository) ListDeploymentVMs(deploymentName string) (deploymentVMs []models.DeploymentVM, apiResponse net.ApiResponse) {
	resources := []deploymentVMResponse{}

	path := fmt.Sprintf("/deployments/%s/vms", deploymentName)
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &resources)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range resources {
		deploymentVMs = append(deploymentVMs, resource.ToModel())
	}

	return
}

type deploymentVMResponse struct {
	JobName string `json:"job"`
	Index   int    `json:"index"`
	VMCid   string `json:"cid"`
	AgentID string `json:"agent_id"`
}

func (resource deploymentVMResponse) ToModel() (vm models.DeploymentVM) {
	vm = models.DeploymentVM{}
	vm.JobName = resource.JobName
	vm.Index = resource.Index
	vm.VMCid = resource.VMCid
	vm.AgentID = resource.AgentID

	return
}
