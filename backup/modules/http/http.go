package http

import (
	"crypto/tls"
	"net/http"
)

type HttpGateway struct {
	endpoint    string
	username    string
	password    string
	contentType string
	handler     HttpResponseHandler
}

func NewHttpGateway(endpoint, username, password, contentType string, handler HttpResponseHandler) *HttpGateway {
	return &HttpGateway{
		endpoint:    endpoint,
		username:    username,
		password:    password,
		contentType: contentType,
		handler:     handler,
	}
}

var NewRoundTripper = func() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (gateway *HttpGateway) Execute(method string) (val interface{}, err error) {
	transport := NewRoundTripper()
	req, err := http.NewRequest(method, gateway.endpoint, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(gateway.username, gateway.password)
	req.Header.Set("Content-Type", gateway.contentType)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return
	}
	return gateway.handler.Handle(resp)
}
