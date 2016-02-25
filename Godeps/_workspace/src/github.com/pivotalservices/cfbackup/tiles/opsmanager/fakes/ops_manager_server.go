package fakes

import "github.com/onsi/gomega/ghttp"

//NewFakeOpsManagerServer - spins up a fake ops manager server
func NewFakeOpsManagerServer(server *ghttp.Server, authStatusCode int, authResponseBody string, genericStatusCode int, genericResponseBody string) *ghttp.Server {
	authHandler := ghttp.RespondWith(authStatusCode, authResponseBody)
	genericHandler := ghttp.RespondWith(genericStatusCode, genericResponseBody)
	server.RouteToHandler("POST", "/uaa/oauth/token", authHandler)
	server.RouteToHandler("GET", "/api/installation_settings", genericHandler)
	return server
}
