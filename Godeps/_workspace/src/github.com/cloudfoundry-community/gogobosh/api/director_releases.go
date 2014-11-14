package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// GetReleases returns the list of releases, and versions available
func (repo BoshDirectorRepository) GetReleases() (releases models.Releases, apiResponse net.ApiResponse) {
	response := []releaseResponse{}

	path := "/releases"
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &response)
	if apiResponse.IsNotSuccessful() {
		return
	}

	list := []*models.Release{}
	for _, resource := range response {
		list = append(list, resource.ToModel())
	}
	releases = models.Releases(list)

	return
}

// DeleteReleases deletes all versions of a release from the BOSH director
func (repo BoshDirectorRepository) DeleteReleases(name string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/releases/%s?force=true", name)
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

// DeleteRelease deletes a specific version of a release from the BOSH director
func (repo BoshDirectorRepository) DeleteRelease(name string, version string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("/releases/%s?force=true&version=%s", name, version)
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

type releaseResponse struct {
	Name     string                   `json:"name"`
	Versions []releaseVersionResponse `json:"release_versions"`
}

type releaseVersionResponse struct {
	Version            string `json:"version"`
	CommitHash         string `json:"commit_hash"`
	UncommittedChanges bool   `json:"uncommitted_changes"`
	CurrentlyDeployed  bool   `json:"currently_deployed"`
}

func (resource releaseResponse) ToModel() (release *models.Release) {
	release = &models.Release{}
	release.Name = resource.Name
	for _, versionResponse := range resource.Versions {
		version := models.ReleaseVersion{}
		version.Version = versionResponse.Version
		version.CommitHash = versionResponse.CommitHash
		version.UncommittedChanges = versionResponse.UncommittedChanges
		version.CurrentlyDeployed = versionResponse.CurrentlyDeployed

		release.Versions = append(release.Versions, version)
	}
	return
}
