package net

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"github.com/cloudfoundry-community/gogobosh/constants"
)

const (
	INVALID_TOKEN_CODE       = "GATEWAY INVALID TOKEN CODE"
	JOB_FINISHED             = "finished"
	JOB_FAILED               = "failed"
	DEFAULT_POLLING_THROTTLE = 5 * time.Second
)

type JobEntity struct {
	Status string
}

type JobResponse struct {
	Entity JobEntity
}

type AsyncMetadata struct {
	Url string
}

type AsyncResponse struct {
	Metadata AsyncMetadata
}

type errorResponse struct {
	Code        string
	Description string
}

type errorHandler func(*http.Response) errorResponse

type Request struct {
	HttpReq      *http.Request
	SeekableBody io.ReadSeeker
}

type Gateway struct {
	errHandler      errorHandler
}

func newGateway(errHandler errorHandler) (gateway Gateway) {
	gateway.errHandler = errHandler
	return
}

func (gateway Gateway) GetResource(url, username string, password string, resource interface{}) (apiResponse ApiResponse) {
	request, apiResponse := gateway.NewRequest("GET", url, username, password, nil)
	if apiResponse.IsNotSuccessful() {
		return
	}

	_, apiResponse = gateway.PerformRequestForJSONResponse(request, resource)
	return
}

func (gateway Gateway) CreateResource(url, username string, password string, body io.ReadSeeker) (apiResponse ApiResponse) {
	return gateway.createUpdateOrDeleteResource("POST", url, username, password, body, nil)
}

func (gateway Gateway) CreateResourceForResponse(url, username string, password string, body io.ReadSeeker, resource interface{}) (apiResponse ApiResponse) {
	return gateway.createUpdateOrDeleteResource("POST", url, username, password, body, resource)
}

func (gateway Gateway) UpdateResource(url, username string, password string, body io.ReadSeeker) (apiResponse ApiResponse) {
	return gateway.createUpdateOrDeleteResource("PUT", url, username, password, body, nil)
}

func (gateway Gateway) UpdateResourceForResponse(url, username string, password string, body io.ReadSeeker, resource interface{}) (apiResponse ApiResponse) {
	return gateway.createUpdateOrDeleteResource("PUT", url, username, password, body, resource)
}

func (gateway Gateway) DeleteResource(url, username string, password string) (apiResponse ApiResponse) {
	return gateway.createUpdateOrDeleteResource("DELETE", url, username, password, nil, &AsyncResponse{})
}

func (gateway Gateway) createUpdateOrDeleteResource(verb, url, username string, password string, body io.ReadSeeker, resource interface{}) (apiResponse ApiResponse) {
	request, apiResponse := gateway.NewRequest(verb, url, username, password, body)
	if apiResponse.IsNotSuccessful() {
		return
	}

	if resource == nil {
		return gateway.PerformRequest(request)
	}

	_, apiResponse = gateway.PerformRequestForJSONResponse(request, resource)
	return
}

func (gateway Gateway) NewRequest(method, path, username string, password string, body io.ReadSeeker) (req *Request, apiResponse ApiResponse) {
	if body != nil {
		body.Seek(0, 0)
	}

	request, err := http.NewRequest(method, path, body)
	if err != nil {
		apiResponse = NewApiResponseWithError("Error building request", err)
		return
	}

	if password != "" {
		data := []byte(username+":"+password)
		auth := base64.StdEncoding.EncodeToString(data)
		request.Header.Set("Authorization", "Basic "+auth)
	}

	request.Header.Set("accept", "application/json")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("User-Agent", "gogobosh "+constants.Version+" / "+runtime.GOOS)

	if body != nil {
		switch v := body.(type) {
		case *os.File:
			fileStats, err := v.Stat()
			if err != nil {
				break
			}
			request.ContentLength = fileStats.Size()
		}
	}

	req = &Request{HttpReq: request, SeekableBody: body}
	return
}

func (gateway Gateway) PerformRequest(request *Request) (apiResponse ApiResponse) {
	_, apiResponse = gateway.doRequestAndHandlerError(request)
	return
}

func (gateway Gateway) PerformRequestForResponseBytes(request *Request) (bytes []byte, headers http.Header, apiResponse ApiResponse) {
	rawResponse, apiResponse := gateway.doRequestAndHandlerError(request)
	if apiResponse.IsNotSuccessful() {
		return
	}

	bytes, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		apiResponse = NewApiResponseWithError("Error reading response", err)
	}

	headers = rawResponse.Header
	return
}

func (gateway Gateway) PerformRequestForTextResponse(request *Request) (response string, headers http.Header, apiResponse ApiResponse) {
	bytes, headers, apiResponse := gateway.PerformRequestForResponseBytes(request)
	response = string(bytes)
	return
}

func (gateway Gateway) PerformRequestForJSONResponse(request *Request, response interface{}) (headers http.Header, apiResponse ApiResponse) {
	bytes, headers, apiResponse := gateway.PerformRequestForResponseBytes(request)
	if apiResponse.IsNotSuccessful() {
		return
	}

	if apiResponse.StatusCode > 203 || strings.TrimSpace(string(bytes)) == "" {
		return
	}

	err := json.Unmarshal(bytes, &response)
	if err != nil {
		apiResponse = NewApiResponseWithError("Invalid JSON response from server", err)
	}
	return
}

func (gateway Gateway) doRequestAndHandlerError(request *Request) (rawResponse *http.Response, apiResponse ApiResponse) {
	rawResponse, err := doRequest(request.HttpReq)
	if err != nil {
		apiResponse = NewApiResponseWithError("Error performing request", err)
		return
	}

	if rawResponse.StatusCode == 302 {
		/* DELETE requests do not automatically redirect; all others should not return 302 */
		apiResponse = NewApiResponseWithRedirect(rawResponse.Header.Get("location"))
	} else if rawResponse.StatusCode > 299 {
		errorResponse := gateway.errHandler(rawResponse)
		message := fmt.Sprintf(
			"Server error, status code: %d, error code: %s, message: %s",
			rawResponse.StatusCode,
			errorResponse.Code,
			errorResponse.Description,
		)
		apiResponse = NewApiResponse(message, errorResponse.Code, rawResponse.StatusCode)
	} else {
		apiResponse = NewApiResponseWithStatusCode(rawResponse.StatusCode)
	}
	return
}
