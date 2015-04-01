package cfbackup_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/gtils/command"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/pivotalservices/gtils/mock"
	"github.com/pivotalservices/gtils/osutils"

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

var MockMultiPartUploadFunc = func(ConnAuth, string, string, io.Reader, map[string]string) (*http.Response, error) {
	return &http.Response{}, nil
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

var (
	nfsSuccessString string = "success nfs"
	nfsFailureString string = "failed nfs"
)

type SuccessMockNFSExecuter struct{}

func (s *SuccessMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(nfsSuccessString))
	return
}

var (
	mockNfsCommandError error = errors.New("error occurred")
)

type FailureMockNFSExecuter struct{}

func (s *FailureMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(nfsFailureString))
	err = mockNfsCommandError
	return
}

type mockRemoteOps struct {
	Err    error
	Writer io.Writer
}

func (s *mockRemoteOps) Path() string {
	return osutils.REMOTE_IMPORT_PATH
}

func (s *mockRemoteOps) UploadFile(lfile io.Reader) error {

	if s.Writer == nil {
		s.Writer = mock.NewReadWriteCloser(nil, nil, nil)
	}

	if s.Err == nil {
		_, s.Err = io.Copy(s.Writer, lfile)
	}
	return s.Err
}

var getNfs = func(lf io.Writer, cmdexec Executer) *NFSBackup {
	return &NFSBackup{
		Caller: cmdexec,
		RemoteOps: &mockRemoteOps{
			Writer: lf,
		},
	}
}

var logger = getLogger("debug")

func Logger() log.Logger {
	return logger
}

func getLogger(minLogLevel string) log.Logger {
	log.SetLogLevel(minLogLevel)
	return log.LogFactory("TestLogger", log.Lager, os.Stdout)
}

var (
	mockTileBackupError  = errors.New("backup tile error")
	mockTileRestoreError = errors.New("restore tile error")
)

type mockTile struct {
	ErrBackup     error
	ErrRestore    error
	RestoreCalled int
	BackupCalled  int
}

func (s *mockTile) Backup() error {
	s.BackupCalled++
	return s.ErrBackup
}

func (s *mockTile) Restore() error {
	s.RestoreCalled++
	return s.ErrRestore
}
