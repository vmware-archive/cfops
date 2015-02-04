package cfbackup_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/http"

	"testing"
)

func TestBackup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup Suite")
}

var (
	redirectUrl       string = "mysite.com"
	successString     string = `{"state":"done"}`
	failureString     string = `{"state":"notdone"}`
	successWaitCalled int
	failureWaitCalled int
	restSuccessCalled int
	restFailureCalled int
)

type MockHttpGateway struct {
	CheckFailureCondition bool
	StatusCode            int
	State                 string
}

// Implements RequestAdaptor. Used to return a successful response
var restSuccess = func() (*http.Response, error) {
	resp := &http.Response{}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
	restSuccessCalled++
	return resp, nil
}

// Implements RequestAdaptor. Used to return a failed response
var restFailure = func() (*http.Response, error) {
	resp := &http.Response{}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
	restFailureCalled++
	err := fmt.Errorf("")
	return resp, err
}

func makeResponse(entity HttpRequestEntity, method string, statusCode int, checkFailure bool, state string, body io.Reader) (*http.Response, error) {
	header := make(map[string][]string)
	locations := []string{redirectUrl}
	header["Location"] = locations
	if statusCode == 0 {
		statusCode = 200
	}
	if state == "" {
		state = "success"
	}
	response := &http.Response{StatusCode: statusCode,
		Header: header,
	}
	if checkFailure {
		restFailureCalled++
		response.Body = &ClosingBuffer{bytes.NewBufferString(state)}
		return response, nil
	}
	restSuccessCalled++
	response.Body = &ClosingBuffer{bytes.NewBufferString(state)}
	return response, nil
}

func (gateway *MockHttpGateway) Get(entity HttpRequestEntity) RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "GET", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, nil)
	}
}

func (gateway *MockHttpGateway) Post(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "POST", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, body)
	}
}

func (gateway *MockHttpGateway) Put(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "PUT", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, body)
	}
}

var MockMultiPartBodyFunc = func(string, string, io.Reader, map[string]string) (io.Reader, error) {
	return &ClosingBuffer{bytes.NewBufferString("success")}, nil
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

type successExecuter struct{}

func (s *successExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	return
}

type failExecuter struct{}

func (s *failExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	err = fmt.Errorf("error failure")
	return
}

type mockLocalExecute func(name string, arg ...string) *exec.Cmd

func (cmd mockLocalExecute) Execute(destination io.Writer, command string) (err error) {
	return
}

func NewLocalMockExecuter() Executer {
	return mockLocalExecute(exec.Command)
}
