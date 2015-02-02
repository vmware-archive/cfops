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

type EvenTaskCreaterAdapter func(method string, requestAdaptor RequestAdaptor) (task EventTasker)

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
}

type Task struct {
	Method         string
	RequestAdaptor RequestAdaptor
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

var NewToggleGateway = func(method, serverUrl, username, password string) func() (interface{}, error) {
	var (
		err  error
		resp *http.Response
	)
	if resp, err = Request(HttpRequestEntity{
		Url:         serverUrl,
		Username:    username,
		Password:    password,
		ContentType: "text/yaml",
	}, method, nil); err == nil {
		return func() (interface{}, error) {
			return ToggleCCHandler(resp)
		}
	}
	return func() (interface{}, error) {
		return nil, err
	}
}

func ToggleCCJobRunner(serverUrl, username, password string) (redirectUrl string, err error) {
	ret, err := NewToggleGateway("PUT", serverUrl, username, password)()
	if err != nil {
		return
	}
	return ret.(string), err
}

func NewCloudController(ip, username, password, deploymentName, state string) *CloudController {
	return &CloudController{
		ip:                  ip,
		username:            username,
		password:            password,
		deploymentName:      deploymentName,
		state:               state,
		JobToggler:          ToggleCCJobRunner,
		NewEventTaskCreater: EvenTaskCreaterAdapter(NewTask),
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
		gateway := s.httpGateway
		if gateway == nil {
			gateway = NewHttpGateway()
		}
		requestAdapter := gateway.Get(HttpRequestEntity{
			Url:         modifyUrl(s.ip, serverURL, originalUrl),
			Username:    s.username,
			Password:    s.password,
			ContentType: "application/json",
		})
		task := s.NewEventTaskCreater("GET", requestAdapter)
		err = task.WaitForEventStateDone(contents, &eventObject)
	}
	return
}

func NewTask(method string, requestAdaptor RequestAdaptor) (task EventTasker) {
	task = &Task{
		Method:         method,
		RequestAdaptor: requestAdaptor,
	}
	return
}

func (s *Task) getEvents(dest io.Writer) (err error) {
	resp, err := s.RequestAdaptor()
	defer resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("Invalid Bosh Director Credentials")
	}
	_, err = io.Copy(dest, resp.Body)
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
