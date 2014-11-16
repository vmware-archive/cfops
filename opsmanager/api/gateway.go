package api

import (
	"net/http"

	"github.com/cloudfoundry-community/gogobosh/net"
)

type Gateway struct {
	baseUrl, username, password string
	net.Gateway
}

func NewOpsManagerGateway(url, username, password string) (gateway Gateway) {
	gateway.baseUrl = url
	gateway.username = username
	gateway.password = password
	gateway.Gateway = net.NewDirectorGateway()
	return
}

func (gateway Gateway) GetAPIVersion(resp interface{}) (http.Header, net.ApiResponse) {
	request, _ := gateway.NewRequest("GET", gateway.baseUrl+"api_version", gateway.username, gateway.password, nil)
	return gateway.PerformRequestForJSONResponse(request, &resp)
}

func (gateway Gateway) GetInstallation(resp interface{}) (http.Header, net.ApiResponse) {
	request, _ := gateway.NewRequest("GET", gateway.baseUrl+"installation_settings", gateway.username, gateway.password, nil)
	return gateway.PerformRequestForJSONResponse(request, &resp)
}
