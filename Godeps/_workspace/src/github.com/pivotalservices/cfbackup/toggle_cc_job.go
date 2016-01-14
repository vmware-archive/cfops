package cfbackup

import (
	"fmt"
	"strings"
	"time"

	"github.com/pivotalservices/gtils/bosh"
	. "github.com/pivotalservices/gtils/http"
)

// Not ping server so frequently and exausted the resources
var TaskPingFreq = 1000 * time.Millisecond

//CloudControllerJobs - array storing a list of CCJobs
type CloudControllerJobs []CCJob

//CloudController - a struct representing a cloud controller
type CloudController struct {
	deploymentName   string
	director         bosh.Bosh
	cloudControllers CloudControllerJobs
	manifest         string
}

//NewDirector - a function representing a constructor for a director object
var NewDirector = func(ip, username, password string, port int) bosh.Bosh {
	return bosh.NewBoshDirector(ip, username, password, port, NewHttpGateway())
}

//NewCloudController - a function representing a constructor for a cloud controller
func NewCloudController(ip, username, password, deploymentName, manifest string, cloudControllers CloudControllerJobs) *CloudController {
	director := NewDirector(ip, username, password, 25555)
	return &CloudController{
		deploymentName:   deploymentName,
		director:         director,
		cloudControllers: cloudControllers,
		manifest:         manifest,
	}
}

//Start - a method to execute a start event on a cloud controller
func (c *CloudController) Start() error {
	return c.toggleController("started")
}

//Stop - a method which executes a stop against a cloud controller
func (c *CloudController) Stop() error {
	return c.toggleController("stopped")
}

func (c *CloudController) toggleController(state string) error {
	for _, ccjob := range c.cloudControllers {
		taskID, err := c.director.ChangeJobState(c.deploymentName, ccjob.Job, state, ccjob.Index, strings.NewReader(c.manifest))
		if err != nil {
			return err
		}
		err = c.waitUntilDone(taskID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CloudController) waitUntilDone(taskID int) (err error) {
	time.Sleep(TaskPingFreq)
	result, err := c.director.RetrieveTaskStatus(taskID)
	if err != nil {
		return
	}
	switch bosh.TASKRESULT[result.State] {
	case bosh.ERROR:
		err = fmt.Errorf("Task %d process failed", taskID)
		return
	case bosh.QUEUED:
		err = c.waitUntilDone(taskID)
		return
	case bosh.PROCESSING:
		err = c.waitUntilDone(taskID)
		return
	case bosh.DONE:
		return
	default:
		err = bosh.TaskResultUnknown
		return
	}
}
