package cfbackup

import (
	"errors"
	"fmt"
	"time"

	"github.com/pivotalservices/gtils/bosh"
	. "github.com/pivotalservices/gtils/http"
)

// Not ping server so frequently and exausted the resources
var TaskPingFreq time.Duration = 1000 * time.Millisecond

type CloudControllerJobs []string

type CloudController struct {
	deploymentName   string
	director         bosh.Bosh
	cloudControllers CloudControllerJobs
}

var NewDirector = func(ip, username, password string, port int) bosh.Bosh {
	return bosh.NewBoshDirector(ip, username, password, port, NewHttpGateway())
}

func NewCloudController(ip, username, password, deploymentName string, cloudControllers CloudControllerJobs) *CloudController {
	director := NewDirector(ip, username, password, 25555)
	return &CloudController{
		deploymentName:   deploymentName,
		director:         director,
		cloudControllers: cloudControllers,
	}
}

func (c *CloudController) Start() error {
	return c.toggleController("started")
}

func (c *CloudController) Stop() error {
	return c.toggleController("stopped")
}

func (c *CloudController) toggleController(state string) error {
	manifest, err := c.director.GetDeploymentManifest(c.deploymentName)
	if err != nil {
		return err
	}
	for ccjobindex, ccjob := range c.cloudControllers {
		taskId, err := c.director.ChangeJobState(c.deploymentName, ccjob, state, ccjobindex, manifest)
		if err != nil {
			return err
		}
		err = c.waitUntilDone(taskId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CloudController) waitUntilDone(taskId int) (err error) {
	time.Sleep(TaskPingFreq)
	result, err := c.director.RetrieveTaskStatus(taskId)
	if err != nil {
		return
	}
	switch bosh.TASKRESULT[result.State] {
	case bosh.ERROR:
		err = errors.New(fmt.Sprintf("Task %d process failed", taskId))
		return
	case bosh.QUEUED:
		err = c.waitUntilDone(taskId)
		return
	case bosh.PROCESSING:
		err = c.waitUntilDone(taskId)
		return
	case bosh.DONE:
		return
	default:
		err = bosh.TaskResultUnknown
		return
	}
}
