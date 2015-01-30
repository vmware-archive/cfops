package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
)

type HttpGateway interface {
	Execute(method string) (interface{}, error)
	Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (*http.Response, error)
}

type DefaultHttpGateway struct {
	endpoint       string
	username       string
	password       string
	contentType    string
	handleResponse HandleRespFunc
}

type HandleRespFunc func(response *http.Response) (interface{}, error)

func NewHttpGateway(endpoint, username, password, contentType string, handler HandleRespFunc) HttpGateway {
	if handler == nil {
		handler = func(response *http.Response) (interface{}, error) {
			return nil, nil
		}
	}
	return &DefaultHttpGateway{
		endpoint:       endpoint,
		username:       username,
		password:       password,
		contentType:    contentType,
		handleResponse: handler,
	}
}

var NewRoundTripper = func() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (gateway *DefaultHttpGateway) Execute(method string) (val interface{}, err error) {
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
	return gateway.handleResponse(resp)
}

func (gateway *DefaultHttpGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	var part io.Writer

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if part, err = writer.CreateFormFile(paramName, filename); err == nil {

		if _, err = io.Copy(part, fileRef); err == nil {

			for key, val := range params {
				_ = writer.WriteField(key, val)
			}
			writer.Close()
			gateway.contentType = writer.FormDataContentType()
			res, err = gateway.makeRequest(body)
		}
	}
	return
}

func (gateway *DefaultHttpGateway) makeRequest(body *bytes.Buffer) (res *http.Response, err error) {
	var req *http.Request
	transport := NewRoundTripper()

	if req, err = http.NewRequest("POST", gateway.endpoint, body); err == nil {
		req.SetBasicAuth(gateway.username, gateway.password)
		req.Header.Add("Content-Type", gateway.contentType)
		res, err = transport.RoundTrip(req)
	}
	return
}
