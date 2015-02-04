package bosh_test

import (
	"bytes"
	"errors"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/bosh"
	. "github.com/pivotalservices/gtils/http"
	. "github.com/pivotalservices/gtils/http/httptest"
)

var (
	responseMock          http.Response
	httpRequestSuccessful bool
	ip                    string = "10.10.10.10"
	port                  int    = 25555
	username              string = "test"
	password              string = "test"
	capturedEntity        HttpRequestEntity
)

func captureEntity(entity HttpRequestEntity) {
	capturedEntity = entity
}

var fakeAdaptor RequestAdaptor = func() (response *http.Response, err error) {
	if !httpRequestSuccessful {
		err = errors.New("Http Request failed")
		return
	}
	response = &responseMock
	return
}

var _ = Describe("Bosh", func() {
	var (
		gateway *MockGateway = &MockGateway{
			FakeGetAdaptor:  fakeAdaptor,
			FakePutAdaptor:  fakeAdaptor,
			FakePostAdaptor: fakeAdaptor,
			Capture:         captureEntity,
		}
		director = NewBoshDirector(ip, username, password, port, gateway)
	)
	Describe("Bosh Deployment", func() {
		Context("Http Request Failure", func() {
			It("Should return error", func() {
				httpRequestSuccessful = false
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("Http Request Succeed", func() {
			BeforeEach(func() {
				httpRequestSuccessful = true
			})
			It("Should return error when http status code is non 200", func() {
				responseMock = readFakeResponse(500, "fixtures/manifest.json", nil)
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return error when http response is not valid manifest", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest_invalid.json", nil)
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return nil error when http response is valid manifest", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest.json", nil)
				_, err := director.GetDeploymentManifest("name")
				Ω(err).Should(BeNil())
			})
			It("Should compose correct request when the http response is valid", func() {
				responseMock = readFakeResponse(200, "fixtures/manifest.json", nil)
				director.GetDeploymentManifest("name")
				Ω(capturedEntity.Url).Should(Equal("https://10.10.10.10:25555/deployments/name"))
				Ω(capturedEntity.Username).Should(Equal("test"))
				Ω(capturedEntity.Password).Should(Equal("test"))
				Ω(capturedEntity.ContentType).Should(Equal("text/yaml"))
			})
		})
	})
	Describe("Bosh Change job state", func() {
		var (
			reader        bytes.Buffer
			emptyHeader   = http.Header{}
			invalidHeader = http.Header{"Location:": {"xxxxx"}}
			validHeader   = http.Header{"Location": {"http://localhost/tasks/231"}}
		)
		Context("Http Request Failure", func() {
			It("Should return error", func() {
				httpRequestSuccessful = false
				_, err := director.GetDeploymentManifest("name")
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("Http Request Succeed", func() {
			BeforeEach(func() {
				httpRequestSuccessful = true
			})
			It("Should return error when http status code is non 302", func() {
				responseMock = readFakeResponse(500, "", nil)
				_, err := director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return error when header does not have Location", func() {
				responseMock = readFakeResponse(302, "", emptyHeader)
				_, err := director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return error when header does not have valid Location", func() {
				responseMock = readFakeResponse(302, "", invalidHeader)
				_, err := director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(err).ShouldNot(BeNil())
			})
			It("Should return nil error when header has have valid Location", func() {
				responseMock = readFakeResponse(302, "", validHeader)
				_, err := director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(err).Should(BeNil())
			})
			It("Should return correct task id when header has have valid Location", func() {
				responseMock = readFakeResponse(302, "", validHeader)
				id, _ := director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(id).Should(Equal(231))
			})
			It("Should compose correct request when the http response is valid", func() {
				responseMock = readFakeResponse(302, "", validHeader)
				director.ChangeJobState("name", "jobName", "startted", 0, &reader)
				Ω(capturedEntity.Url).Should(Equal("https://10.10.10.10:25555/deployments/name/jobs/jobName/0?state=startted"))
				Ω(capturedEntity.ContentType).Should(Equal("text/yaml"))
			})
		})
		Describe("Bosh Retrieve Task staus", func() {
			Context("Http Request Failure", func() {
				It("Should return error", func() {
					httpRequestSuccessful = false
					_, err := director.RetrieveTaskStatus(100)
					Ω(err).ShouldNot(BeNil())
				})
			})
			Context("Http Request Succeed", func() {
				BeforeEach(func() {
					httpRequestSuccessful = true
				})
				It("Should return error when http status code is non 200", func() {
					responseMock = readFakeResponse(500, "", nil)
					_, err := director.RetrieveTaskStatus(100)
					Ω(err).ShouldNot(BeNil())
				})
				It("Should compose correct request when the http response is valid", func() {
					responseMock = readFakeResponse(200, "", validHeader)
					director.RetrieveTaskStatus(100)
					Ω(capturedEntity.Url).Should(Equal("https://10.10.10.10:25555/tasks/100"))
					Ω(capturedEntity.ContentType).Should(Equal(NO_CONTENT_TYPE))
				})
				It("Should return error when response task is invalid", func() {
					responseMock = readFakeResponse(200, "fixtures/task_invalid.json", validHeader)
					_, err := director.RetrieveTaskStatus(100)
					Ω(err).ShouldNot(BeNil())
				})
				It("Should return nil error when response task is valid", func() {
					responseMock = readFakeResponse(200, "fixtures/task.json", validHeader)
					_, err := director.RetrieveTaskStatus(100)
					Ω(err).Should(BeNil())
				})
				It("Should return correct task data", func() {
					responseMock = readFakeResponse(200, "fixtures/task.json", validHeader)
					task, _ := director.RetrieveTaskStatus(19)
					Ω(task.State).Should(Equal("done"))
					Ω(task.Id).Should(Equal(19))
					Ω(task.Description).Should(Equal("create deployment"))
				})
			})
		})
	})

})

func readFakeResponse(statusCode int, file string, header http.Header) (res http.Response) {
	body, _ := os.Open(file)
	response := http.Response{}
	response.Body = body
	response.Header = header
	response.StatusCode = statusCode
	return response
}
