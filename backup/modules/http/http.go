package http

import (
	"crypto/tls"
	"net/http"
)

type HttpExecutor struct {
	endpoint    string
	username    string
	password    string
	contentType string
	handler     Handler
}

func NewHttpExecutor(endpoint, username, password, contentType string, handler Handler) *HttpExecutor {
	return &HttpExecutor{
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

func (executor *HttpExecutor) Execute(method string) (val interface{}, err error) {
	transport := NewRoundTripper()
	req, err := http.NewRequest(method, executor.endpoint, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(executor.username, executor.password)
	req.Header.Set("Content-Type", executor.contentType)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return
	}
	return executor.handler.Handle(resp)
}
