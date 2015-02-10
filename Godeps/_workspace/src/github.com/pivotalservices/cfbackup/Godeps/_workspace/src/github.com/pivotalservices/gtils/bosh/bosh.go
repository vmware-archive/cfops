package bosh

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/http"
)

type Bosh interface {
	GetDeploymentManifest(deploymentName string) (io.Reader, error)
	//Currently there is a defect on bosh director we need to pass in the manifest file
	ChangeJobState(string, string, string, int, io.Reader) (int, error)
	RetrieveTaskStatus(int) (*Task, error)
}

type BoshDirector struct {
	ip       string
	port     int
	username string
	password string
	gateway  http.HttpGateway
}

func NewBoshDirector(ip, username, password string, port int, gateway http.HttpGateway) *BoshDirector {
	return &BoshDirector{
		ip:       ip,
		port:     port,
		username: username,
		password: password,
		gateway:  gateway,
	}
}

func (director *BoshDirector) GetDeploymentManifest(deploymentName string) (manifest io.Reader, err error) {
	endpoint := fmt.Sprintf("https://%s:%d/deployments/%s", director.ip, director.port, deploymentName)
	httpEntity := director.getEntity(endpoint, "text/yaml")
	request := director.gateway.Get(httpEntity)
	resp, err := request()
	if err != nil {
		return
	}
	manifest, err = retrieveManifest(resp)
	return
}

func (director *BoshDirector) ChangeJobState(deploymentName, jobName, state string, index int, manifest io.Reader) (taskId int, err error) {
	endpoint := fmt.Sprintf("https://%s:%d/deployments/%s/jobs/%s/%d?state=%s", director.ip, director.port, deploymentName, jobName, index, state)
	httpEntity := director.getEntity(endpoint, "text/yaml")
	request := director.gateway.Put(httpEntity, manifest)
	resp, err := request()
	if err != nil {
		return
	}
	return retrieveTaskId(resp)
}

func (director *BoshDirector) RetrieveTaskStatus(taskId int) (task *Task, err error) {
	endpoint := fmt.Sprintf("https://%s:%d/tasks/%d", director.ip, director.port, taskId)
	httpEntity := director.getEntity(endpoint, http.NO_CONTENT_TYPE)
	request := director.gateway.Get(httpEntity)
	resp, err := request()
	if err != nil {
		return
	}
	return retrieveTaskStatus(resp)
}

func (director *BoshDirector) getEntity(endpoint, contentType string) (httpEntity http.HttpRequestEntity) {
	httpEntity = http.HttpRequestEntity{
		Url:         endpoint,
		Username:    director.username,
		Password:    director.password,
		ContentType: contentType,
	}
	return
}
