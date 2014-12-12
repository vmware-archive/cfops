package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/utils"
)

type CloudControllerJobs []string

type RestAdapter func(method, connectionURL, username, password string, isYaml bool) (*http.Response, error)

type JobTogglerAdapter func(serverUrl, username, password string, exec command.Executer) (res string, err error)

type EvenTaskCreaterAdapter func(method, url, username, password string, isYaml bool) (task EventTasker)

type EventTasker interface {
	WaitForEventStateDone(contents bytes.Buffer, eventObject *utils.EventObject) (err error)
}

type CloudController struct {
	ip                  string
	username            string
	password            string
	deploymentName      string
	state               string
	JobToggler          JobTogglerAdapter
	NewEventTaskCreater EvenTaskCreaterAdapter
}

type Task struct {
	Method     string
	Url        string
	Username   string
	Password   string
	IsYaml     bool
	RestRunner RestAdapter
}

func (restAdapter RestAdapter) Run(method, connectionURL, username, password string, isYaml bool) (statusCode int, body io.Reader, err error) {
	res, err := restAdapter(method, connectionURL, username, password, isYaml)
	defer res.Body.Close()
	body = res.Body
	statusCode = res.StatusCode
	return
}

func ToggleCCJobRunner(serverUrl, username, password string, exec command.Executer) (res string, err error) {
	var b bytes.Buffer
	formatString := `curl -v -XPUT -u "%s:%s" %s --insecure -H "Content-Type:text/yaml" -i -s | grep Location: | grep Location: | cut -d ' ' -f 2`
	cmd := fmt.Sprintf(formatString, username, password, serverUrl)
	err = exec.Execute(&b, cmd)
	res = b.String()
	return
}

func NewCloudController(ip, username, password, deploymentName, state string) *CloudController {
	return &CloudController{
		ip:                  ip,
		username:            username,
		password:            password,
		deploymentName:      deploymentName,
		state:               state,
		JobToggler:          JobTogglerAdapter(ToggleCCJobRunner),
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
		eventObject   utils.EventObject
		connectionURL string = newConnectionURL(serverURL, s.deploymentName, ccjob, s.state, ccjobindex)
	)

	if originalUrl, err := s.JobToggler(connectionURL, s.username, s.password, command.NewLocalExecuter()); err == nil {
		task := s.NewEventTaskCreater("GET", modifyUrl(s.ip, serverURL, originalUrl), s.username, s.password, false)
		err = task.WaitForEventStateDone(contents, &eventObject)
	}
	return
}

func NewTask(method, url, username, password string, isYaml bool) (task EventTasker) {
	task = &Task{
		Method:     method,
		Url:        url,
		Username:   username,
		Password:   password,
		IsYaml:     isYaml,
		RestRunner: RestAdapter(invoke),
	}
	return
}

func (s *Task) getEvents(dest io.Writer) (err error) {
	statusCode, body, err := s.RestRunner.Run(s.Method, s.Url, s.Username, s.Password, s.IsYaml)

	if statusCode == 200 {
		io.Copy(dest, body)

	} else {
		err = fmt.Errorf("Invalid Bosh Director Credentials")
	}
	return
}

func (s *Task) WaitForEventStateDone(contents bytes.Buffer, eventObject *utils.EventObject) (err error) {

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
