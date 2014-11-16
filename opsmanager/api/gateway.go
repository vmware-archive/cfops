package api

import (
	// "fmt"
	// "net/http"

	"github.com/cloudfoundry-community/gogobosh/net"
)

type Gateway interface {
	GetAPIVersion(resp interface{}) net.ApiResponse
	GetInstallationSettings(resp interface{}) net.ApiResponse
}

type gatewayImpl struct {
	baseUrl, username, password string
	net.Gateway
}

func NewOpsManagerGateway(url, username, password string) (gateway gatewayImpl) {
	gateway.baseUrl = url
	gateway.username = username
	gateway.password = password
	gateway.Gateway = net.NewDirectorGateway()
	return
}

func (gateway gatewayImpl) GetAPIVersion(resp interface{}) net.ApiResponse {
	request, _ := gateway.NewRequest("GET", gateway.baseUrl+"api_version", gateway.username, gateway.password, nil)
	// fmt.Println("--------------------------------------------------------------------------------")
	// fmt.Println("HTTP Request Headers")
	// fmt.Println(request.HttpReq.Header)
	_, apiResponse := gateway.PerformRequestForJSONResponse(request, &resp)
	// fmt.Println("HTTP Response Headers")
	// fmt.Println(headers)
	// fmt.Println("--------------------------------------------------------------------------------")
	return apiResponse
}

func (gateway gatewayImpl) GetInstallationSettings(resp interface{}) net.ApiResponse {
	request, _ := gateway.NewRequest("GET", gateway.baseUrl+"installation_settings", gateway.username, gateway.password, nil)
	_, apiResponse := gateway.PerformRequestForJSONResponse(request, &resp)
	return apiResponse
}
