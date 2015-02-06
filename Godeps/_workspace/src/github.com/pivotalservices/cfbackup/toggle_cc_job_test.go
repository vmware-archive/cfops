package cfbackup_test

import (
	"errors"
	. "github.com/pivotalservices/cfbackup"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/gtils/bosh"
)

var (
	getManifest             bool = true
	getTaskStatus           bool = true
	changeJobState          bool = true
	manifest                io.Reader
	ip                      string              = "10.10.10.10"
	username                string              = "test"
	password                string              = "test"
	deploymentName          string              = "deployment"
	ccjobs                  CloudControllerJobs = CloudControllerJobs{"job1", "job2", "job3"}
	task                    bosh.Task
	doneTask                bosh.Task = bosh.Task{}
	changeJobStateCount     int       = 0
	retrieveTaskStatusCount int       = 0
)

type mockDirector struct{}

func (director *mockDirector) GetDeploymentManifest(deploymentName string) (io.Reader, error) {
	if !getManifest {
		return nil, errors.New("")
	}
	return manifest, nil
}

func (director *mockDirector) ChangeJobState(deploymentName, jobName, state string, index int, manifest io.Reader) (int, error) {
	changeJobStateCount++
	if !changeJobState {
		return 0, errors.New("")
	}
	return 1, nil
}

func (director *mockDirector) RetrieveTaskStatus(int) (*bosh.Task, error) {
	if !getTaskStatus {
		return nil, errors.New("")
	}
	retrieveTaskStatusCount++
	if retrieveTaskStatusCount%2 == 0 {
		return &bosh.Task{State: "processing"}, nil
	}
	return &task, nil
}

var _ = Describe("ToggleCcJob", func() {
	NewDirector = func(ip, username, password string, port int) bosh.Bosh {
		return &mockDirector{}
	}
	TaskPingFreq = time.Millisecond
	var (
		cloudController *CloudController = NewCloudController(ip, username, password, deploymentName, ccjobs)
	)
	Describe("Toggle All jobs", func() {
		Context("Change Job State failed", func() {
			BeforeEach(func() {
				changeJobState = false
			})
			It("Should return error", func() {
				err := cloudController.Start()
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("Toggle successfully", func() {
			BeforeEach(func() {
				changeJobState = true
				changeJobStateCount = 0
				task = bosh.Task{State: "done"}
				retrieveTaskStatusCount = 0
			})
			It("Should return nil error", func() {
				err := cloudController.Start()
				Ω(err).Should(BeNil())
			})
			It("Should Call changeJobState 3 times with 3 jobs", func() {
				cloudController.Start()
				Ω(changeJobStateCount).Should(Equal(3))
			})
			It("Should Call retriveTaskStatus 5 times with retries when task is processing", func() {
				cloudController.Start()
				Ω(retrieveTaskStatusCount).Should(Equal(5))
			})
		})
		Context("Task status is error", func() {
			BeforeEach(func() {
				changeJobState = true
				task = bosh.Task{State: "error"}
			})
			It("Should return error", func() {
				err := cloudController.Start()
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
