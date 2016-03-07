package fakes

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/gtils/bosh"
	"github.com/pivotalservices/gtils/command"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/pivotalservices/gtils/mock"
	"github.com/pivotalservices/gtils/osutils"
)

var (
	redirectURL = "mysite.com"
	//SuccessString --
	SuccessString = `{"state":"done"}`
	//FailureString --
	FailureString     = `{"state":"notdone"}`
	successWaitCalled int
	failureWaitCalled int
	restSuccessCalled int
	restFailureCalled int
)

//NewFakeDirector ---
func NewFakeDirector(ip, username, password string, port int) bosh.Bosh {
	return &mockDirector{
		getManifest:             true,
		manifest:                strings.NewReader("manifest"),
		changeJobState:          true,
		changeJobStateCount:     0,
		getTaskStatus:           true,
		retrieveTaskStatusCount: 0,
	}
}

type mockDirector struct {
	getManifest             bool
	manifest                io.Reader
	changeJobStateCount     int
	changeJobState          bool
	getTaskStatus           bool
	retrieveTaskStatusCount int
}

func (director *mockDirector) GetDeploymentManifest(deploymentName string) (io.Reader, error) {
	if !director.getManifest {
		return nil, errors.New("")
	}
	return director.manifest, nil
}

func (director *mockDirector) ChangeJobState(deploymentName, jobName, state string, index int, manifest io.Reader) (int, error) {
	director.changeJobStateCount++
	if !director.changeJobState {
		return 0, errors.New("")
	}
	return 1, nil
}

func (director *mockDirector) RetrieveTaskStatus(int) (task *bosh.Task, err error) {
	if !director.getTaskStatus {
		return nil, errors.New("")
	}
	director.retrieveTaskStatusCount++
	if director.retrieveTaskStatusCount%2 == 0 {
		return &bosh.Task{State: "processing"}, nil
	}
	return task, nil
}

//MockHTTPGateway --
type MockHTTPGateway struct {
	CheckFailureCondition bool
	StatusCode            int
	State                 string
}

// Implements RequestAdaptor. Used to return a successful response
var restSuccess = func() (*http.Response, error) {
	resp := &http.Response{}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(SuccessString)}
	restSuccessCalled++
	return resp, nil
}

// Implements RequestAdaptor. Used to return a failed response
var restFailure = func() (*http.Response, error) {
	resp := &http.Response{}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(FailureString)}
	restFailureCalled++
	err := fmt.Errorf("")
	return resp, err
}

func makeResponse(entity ghttp.HttpRequestEntity, method string, statusCode int, checkFailure bool, state string, body io.Reader) (*http.Response, error) {
	header := make(map[string][]string)
	locations := []string{redirectURL}
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

//Get --
func (gateway *MockHTTPGateway) Get(entity ghttp.HttpRequestEntity) ghttp.RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "GET", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, nil)
	}
}

//Post ---
func (gateway *MockHTTPGateway) Post(entity ghttp.HttpRequestEntity, body io.Reader) ghttp.RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "POST", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, body)
	}
}

//Put --
func (gateway *MockHTTPGateway) Put(entity ghttp.HttpRequestEntity, body io.Reader) ghttp.RequestAdaptor {
	return func() (*http.Response, error) {
		return makeResponse(entity, "PUT", gateway.StatusCode, gateway.CheckFailureCondition, gateway.State, body)
	}
}

//MockMultiPartBodyFunc --
var MockMultiPartBodyFunc = func(string, string, io.Reader, map[string]string) (io.Reader, error) {
	return &ClosingBuffer{bytes.NewBufferString("success")}, nil
}

//MockMultiPartUploadFunc ---
var MockMultiPartUploadFunc = func(ghttp.ConnAuth, string, string, int64, io.Reader, map[string]string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
	}, nil
}

//ClosingBuffer ---
type ClosingBuffer struct {
	*bytes.Buffer
}

//Close --
func (cb *ClosingBuffer) Close() (err error) {
	return
}

//SuccessExecuter ---
type SuccessExecuter struct{}

//Execute --
func (s *SuccessExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	return
}

//FailExecuter ---
type FailExecuter struct{}

//Execute --
func (s *FailExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	err = fmt.Errorf("error failure")
	return
}

type mockLocalExecute func(name string, arg ...string) *exec.Cmd

func (cmd mockLocalExecute) Execute(destination io.Writer, command string) (err error) {
	return
}

//NewLocalMockExecuter ---
func NewLocalMockExecuter() command.Executer {
	return mockLocalExecute(exec.Command)
}

//Exported status strings
var (
	NfsSuccessString = "success nfs"
	NfsFailureString = "failed nfs"
)

//SuccessMockNFSExecuter ---
type SuccessMockNFSExecuter struct{}

//Execute --
func (s *SuccessMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(NfsSuccessString))
	return
}

var (
	//ErrMockNfsCommand ---
	ErrMockNfsCommand = errors.New("error occurred")
)

//FailureMockNFSExecuter --
type FailureMockNFSExecuter struct{}

//Execute --
func (s *FailureMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(NfsFailureString))
	err = ErrMockNfsCommand
	return
}

type mockRemoteOps struct {
	Err    error
	Writer io.Writer
}

func (s *mockRemoteOps) Path() string {
	return osutils.RemoteImportPath
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

func (s *mockRemoteOps) RemoveRemoteFile() (err error) {

	return
}

//GetNfs ---
var GetNfs = func(lf io.Writer, cmdexec command.Executer) *cfbackup.NFSBackup {
	return &cfbackup.NFSBackup{
		Caller: cmdexec,
		RemoteOps: &mockRemoteOps{
			Writer: lf,
		},
	}
}

var logger = getLogger("debug")

//Logger ---
func Logger() log.Logger {
	return logger
}

func getLogger(minLogLevel string) log.Logger {
	log.SetLogLevel(minLogLevel)
	return log.LogFactory("TestLogger", log.Lager, os.Stdout)
}

var (
	errMockTileBackup  = errors.New("backup tile error")
	errMockTileRestore = errors.New("restore tile error")
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

//NewFakeBackupContext --
func NewFakeBackupContext(target string, env map[string]string, storageProvider cfbackup.StorageProvider) (backupContext cfbackup.BackupContext) {
	backupContext = cfbackup.NewBackupContext(target, env, "")
	backupContext.StorageProvider = storageProvider
	return
}

//FakeStorageProvider --
type FakeStorageProvider struct {
}

//Reader ---
func (d *FakeStorageProvider) Reader(path ...string) (io.ReadCloser, error) {
	return &ClosingBuffer{
		bytes.NewBufferString("Fake storage provider doesn't care about your data"),
	}, nil
}

//Writer --
func (d *FakeStorageProvider) Writer(path ...string) (closer io.WriteCloser, err error) {
	return closer, nil
}

//NewMockStringStorageProvider ---
func NewMockStringStorageProvider() *MockStringStorageProvider {
	return &MockStringStorageProvider{
		Buffer: bytes.NewBufferString(""),
	}
}

//MockStringStorageProvider ----
type MockStringStorageProvider struct {
	*bytes.Buffer
	ErrFakeResponse      error
	ErrFakeCloseResponse error
}

//Close ----
func (s *MockStringStorageProvider) Close() (err error) {
	return s.ErrFakeCloseResponse
}

//Reader ---
func (s *MockStringStorageProvider) Reader(path ...string) (read io.ReadCloser, err error) {
	return s, s.ErrFakeResponse
}

//Writer ----
func (s *MockStringStorageProvider) Writer(path ...string) (writer io.WriteCloser, err error) {
	return s, s.ErrFakeResponse
}
