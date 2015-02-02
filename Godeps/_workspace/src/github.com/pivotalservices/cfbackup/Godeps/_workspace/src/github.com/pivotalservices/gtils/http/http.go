package http

import (
	"crypto/tls"
	"io"
	"net/http"
)

const NO_CONTENT_TYPE string = ""

type HttpRequestEntity struct {
	Url         string
	Username    string
	Password    string
	ContentType string
}

type RequestFunc func(HttpRequestEntity, string, io.Reader) (*http.Response, error)

func Request(entity HttpRequestEntity, method string, body io.Reader) (response *http.Response, err error) {
	transport := NewRoundTripper()
	req, err := http.NewRequest(method, entity.Url, body)
	req.SetBasicAuth(entity.Username, entity.Password)
	if entity.ContentType != NO_CONTENT_TYPE {
		req.Header.Add("Content-Type", entity.ContentType)
	}

	if err != nil {
		return
	}
	return transport.RoundTrip(req)
}

type RequestAdaptor func() (*http.Response, error)

type HttpGateway interface {
	Get(HttpRequestEntity) RequestAdaptor
	Post(HttpRequestEntity, io.Reader) RequestAdaptor
	Put(HttpRequestEntity, io.Reader) RequestAdaptor
}

func NewHttpGateway() HttpGateway {
	return &DefaultHttpGateway{}
}

type DefaultHttpGateway struct{}

func (gateway *DefaultHttpGateway) Get(entity HttpRequestEntity) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "GET", nil)
	}
}

func (gateway *DefaultHttpGateway) Post(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "POST", body)
	}
}

func (gateway *DefaultHttpGateway) Put(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "PUT", body)
	}
}

var NewRoundTripper = func() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
