package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

func (repo BoshDirectorRepository) FetchVMsStatus(deploymentName string) (vmsStatuses []models.VMStatus, apiResponse net.ApiResponse) {
	var taskStatus models.TaskStatus

	/*
	* Two API calls
	* 1. GET /deployments/%s/vms?format=full and be redirected to a /tasks/123
	* 2. Streaming GET on /tasks/123/output?type=result - each line is a models.VMStatus
	 */
	path := fmt.Sprintf("/deployments/%s/vms?format=full", deploymentName)
	apiResponse = repo.gateway.GetResource(repo.config.TargetURL+path, repo.config.Username, repo.config.Password, &taskStatus)
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

	path = fmt.Sprintf("/tasks/%d/output?type=result", taskStatus.ID)
	request, apiResponse := repo.gateway.NewRequest("GET", repo.config.TargetURL+path, repo.config.Username, repo.config.Password, nil)
	if apiResponse.IsNotSuccessful() {
		return
	}

	bytes, _, apiResponse := repo.gateway.PerformRequestForResponseBytes(request)
	if apiResponse.IsNotSuccessful() {
		return
	}

	if apiResponse.StatusCode > 203 {
		return
	}

	for _, vmStatusItem := range strings.Split(string(bytes), "\n") {
		resource := vmStatusResponse{}
		err := json.Unmarshal([]byte(vmStatusItem), &resource)
		if err == nil {
			vmsStatuses = append(vmsStatuses, resource.ToModel())
		}
	}

	return
}

type vmStatusResponse struct {
	JobName            string         `json:"job_name"`
	Index              int            `json:"index"`
	JobState           string         `json:"job_state"`
	VMCid              string         `json:"vm_cid"`
	AgentID            string         `json:"agent_id"`
	IPs                []string       `json:"ips"`
	DNSs               []string       `json:"dns"`
	ResourcePool       string         `json:"resource_pool"`
	ResurrectionPaused bool           `json:"resurrection_paused"`
	Vitals             vitalsResponse `json:"vitals"`
}

type vitalsResponse struct {
	Load   []string          `json:"load"`
	CPU    cpuResponse       `json:"cpu"`
	Memory percentKbResponse `json:"mem"`
	Swap   percentKbResponse `json:"swap"`
	Disk   diskResponse      `json:"disk"`
}

type cpuResponse struct {
	User   float64 `json:"user,string"`
	System float64 `json:"sys,string"`
	Wait   float64 `json:"wait,string"`
}

type diskResponse struct {
	Persistent percentKbResponse `json:"persistent"`
}

type percentKbResponse struct {
	Percent float64 `json:"percent,string"`
	Kb      int     `json:"kb,string"`
}

func (resource vmStatusResponse) ToModel() (status models.VMStatus) {
	status = models.VMStatus{}
	status.JobName = resource.JobName
	status.Index = resource.Index
	status.JobState = resource.JobState
	status.VMCid = resource.VMCid
	status.AgentID = resource.AgentID
	status.ResourcePool = resource.ResourcePool
	status.ResurrectionPaused = resource.ResurrectionPaused

	status.IPs = resource.IPs
	status.DNSs = resource.DNSs

	status.CPUUser = resource.Vitals.CPU.User
	status.CPUSys = resource.Vitals.CPU.System
	status.CPUWait = resource.Vitals.CPU.Wait
	status.MemoryPercent = resource.Vitals.Memory.Percent
	status.MemoryKb = resource.Vitals.Memory.Kb
	status.SwapPercent = resource.Vitals.Swap.Percent
	status.SwapKb = resource.Vitals.Swap.Kb
	status.DiskPersistentPercent = resource.Vitals.Disk.Persistent.Percent

	return
}
