package bosh

import (
	"io"

	"github.com/pivotalservices/gtils/http"
)

type Bosh interface {
	GetDeploymentManifest(deploymentName string) (io.Reader, error)
}

type BoshDirector struct {
	ip       string
	port     int
	username string
	password string
}

func NewBoshDirector(ip, username, password string, port int) *BoshDirector {
	return &BoshDirector{
		ip:       ip,
		port:     port,
		username: username,
		password: password,
	}
}

// We can continute to plugin more apis
var APIs = map[string]API{"manifest": retrieveManifestAPI}

var NewBoshGateway = func(endpoint, username, password, contentType string, handler http.HandleRespFunc, body io.Reader) (gateway http.HttpGateway) {
	return http.NewHttpGateway(endpoint, username, password, contentType, handler, body)
}

func (director *BoshDirector) execute(api API, pathParams, queryParams map[string]string, body io.Reader) (ret interface{}, err error) {
	endpoint, err := ParseUrl(director.ip, director.port, api.Path, pathParams, queryParams)
	if err != nil {
		return
	}
	gateway := NewBoshGateway(endpoint, director.ip, director.password, api.ContentType, api.HandleResponse, body)
	return gateway.Execute(api.Method)
}

func (director *BoshDirector) GetDeploymentManifest(deploymentName string) (manifest io.Reader, err error) {
	api := APIs["manifest"]
	pathParams := map[string]string{"deployment": deploymentName}
	m, err := director.execute(api, pathParams, nil, nil)
	return m.(io.Reader), nil
}
