package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// GetStemcells returns the list of stemcells & versions available
func (repo BoshDirectorRepository) GetStemcells() (stemcells models.Stemcells, apiResponse net.ApiResponse) {
	response := []stemcellResponse{}

	path := "/stemcells"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &response)
	if apiResponse.IsNotSuccessful() {
		return
	}

	list := []*models.Stemcell{}
	for _, resource := range response {
		list = append(list, resource.ToModel())
	}
	stemcells = models.Stemcells(list)

	return
}

// DeleteStemcell deletes a specific stemcell version
func (repo BoshDirectorRepository) DeleteStemcell(name string, version string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/stemcells/%s/%s?force=true", name, version)
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

type stemcellResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Cid     string `json:"cid"`
}

func (resource stemcellResponse) ToModel() (stemcell *models.Stemcell) {
	stemcell = &models.Stemcell{}
	stemcell.Name = resource.Name
	stemcell.Version = resource.Version
	stemcell.Cid = resource.Cid

	return
}
