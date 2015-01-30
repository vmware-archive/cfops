package cfbackup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	. "github.com/pivotalservices/gtils/http"
)

type CloudControllerJobs []string

type JobTogglerAdapter func(serverUrl, username, password string) (res string, err error)

type EvenTaskCreaterAdapter func(method string, httpGateway HttpGateway, handleRespFunc HandleRespFunc) (task EventTasker)

type EventTasker interface {
	WaitForEventStateDone(contents bytes.Buffer, eventObject *EventObject) (err error)
}

type EventObject struct {
	Id          int    `json:"id"`
	State       string `json:"state"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

type CloudController struct {
	ip                  string
	username            string
	password            string
	deploymentName      string
	state               string
	JobToggler          JobTogglerAdapter
	NewEventTaskCreater EvenTaskCreaterAdapter
	httpGateway         HttpGateway
	HttpResponseHandler HandleRespFunc
}

type Task struct {
	Method              string
	HttpGateway         HttpGateway
	HttpResponseHandler HandleRespFunc
}

func ToggleCCHandler(response *http.Response) (redirectUrl interface{}, err error) {
	if response.StatusCode != 302 {
		err = errors.New("The response code from toggle request should return 302")
		return
	}
	redirectUrls := response.Header["Location"]
	if redirectUrls == nil || len(redirectUrls) < 1 {
		err = errors.New("Could not find redirect url for bosh tasks")
		return
	}
	return redirectUrls[0], nil
}

var NewToggleGateway = func(serverUrl, username, password string) HttpGateway {
	return NewHttpGateway(serverUrl, username, password, "text/yaml", ToggleCCHandler)
}

func ToggleCCJobRunner(serverUrl, username, password string) (redirectUrl string, err error) {
	httpGateway := NewToggleGateway(serverUrl, username, password)
	ret, err := httpGateway.Execute("PUT")
	if err != nil {
		return
	}
	return ret.(string), err
}

func NewCloudController(ip, username, password, deploymentName, state string, handleRespFunc HandleRespFunc) *CloudController {
	return &CloudController{
		ip:                  ip,
		username:            username,
		password:            password,
		deploymentName:      deploymentName,
		state:               state,
		JobToggler:          ToggleCCJobRunner,
		NewEventTaskCreater: EvenTaskCreaterAdapter(NewTask),
		HttpResponseHandler: handleRespFunc,
	}
}

func (s *CloudController) ToggleJobs(ccjobs CloudControllerJobs) (err error) {
	serverURL := serverUrlFromIp(s.ip)

	for ccjobindex, ccjob := range ccjobs {
		err = s.ToggleJob(ccjob, serverURL, ccjobindex)
	}
	return
}

func (s *CloudController) ToggleJob(ccjob, serverURL string, ccjobindex int) (err error) {
	var (
		contents      bytes.Buffer
		eventObject   EventObject
		originalUrl   string
		connectionURL string = newConnectionURL(serverURL, s.deploymentName, ccjob, s.state, ccjobindex)
	)

	if originalUrl, err = s.JobToggler(connectionURL, s.username, s.password); err == nil {
		gateway := NewHttpGateway(modifyUrl(s.ip, serverURL, originalUrl), s.username, s.password, "application/json", s.HttpResponseHandler)
		task := s.NewEventTaskCreater("GET", gateway, s.HttpResponseHandler)
		err = task.WaitForEventStateDone(contents, &eventObject)
	}
	return
}

func NewTask(method string, httpGateway HttpGateway, handleRespFunc HandleRespFunc) (task EventTasker) {
	task = &Task{
		Method:              method,
		HttpGateway:         httpGateway,
		HttpResponseHandler: handleRespFunc,
	}
	return
}

func (s *Task) getEvents(dest io.Writer) (err error) {
	var responseHandler = s.HttpResponseHandler
	if responseHandler == nil {
		responseHandler = func(resp *http.Response) (interface{}, error) {
			defer resp.Body.Close()
			return io.Copy(dest, resp.Body)
		}
	}
	contents, err := s.HttpGateway.ExecuteFunc("GET", responseHandler)
	if err != nil {
		err = fmt.Errorf("Invalid Bosh Director Credentials")
	}
	if responseHandler != nil {
		io.Copy(dest, contents.(io.Reader))
	}
	return
}

func (s *Task) WaitForEventStateDone(contents bytes.Buffer, eventObject *EventObject) (err error) {

	if err = json.Unmarshal(contents.Bytes(), eventObject); err == nil && eventObject.State != "done" {
		contents.Reset()

		if err = s.getEvents(&contents); err == nil {
			s.WaitForEventStateDone(contents, eventObject)
		}
	}
	return
}

func modifyUrl(ip, serverURL, originalUrl string) (newUrl string) {
	newUrl = strings.Replace(originalUrl, "https://"+ip+"/", serverURL, 1)
	return
}

func serverUrlFromIp(ip string) (serverUrl string) {
	serverUrl = "https://" + ip + ":25555/"
	return
}

func newConnectionURL(serverURL, deploymentName, ccjob, state string, ccjobindex int) (connectionURL string) {
	connectionURL = serverURL + "deployments/" + deploymentName + "/jobs/" + ccjob + "/" + strconv.Itoa(ccjobindex) + "?state=" + state
	return
}
