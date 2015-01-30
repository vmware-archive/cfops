package cfbackup_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	. "github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/gtils/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	successControlOuput string = "successful execute"
	failureControlOuput string = "failed to execute"
	redirectUrl         string = "mysite.com"
	successWaitCalled   int
	failureWaitCalled   int
)

type MockSuccessGateway struct{ Handler HandleRespFunc }

func makeResponse(location string, statusCode int) *http.Response {
	header := make(map[string][]string)
	locations := []string{redirectUrl}
	header["Location"] = locations
	response := &http.Response{StatusCode: statusCode,
		Header: header,
	}
	return response
}

func (gateway *MockSuccessGateway) Execute(method string) (val interface{}, err error) {
	return gateway.Handler(makeResponse(redirectUrl, 302))
}

func (gateway *MockSuccessGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	return nil, nil
}

func (gateway *MockSuccessGateway) ExecuteFunc(method string, handler HandleRespFunc) (interface{}, error) {
	return nil, nil
}

type MockFailerGateway struct{ Handler HandleRespFunc }

func (gateway *MockFailerGateway) Execute(method string) (val interface{}, err error) {
	return gateway.Handler(makeResponse(redirectUrl, 200))
}

func (gateway *MockFailerGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	return nil, nil
}

func (gateway *MockFailerGateway) ExecuteFunc(method string, handler HandleRespFunc) (interface{}, error) {
	return nil, nil
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

type SuccessMockEventTasker struct {
}

func (s SuccessMockEventTasker) WaitForEventStateDone(contents bytes.Buffer, eventObject *EventObject) (err error) {
	successWaitCalled++
	return
}

type FailureMockEventTasker struct {
}

func (s FailureMockEventTasker) WaitForEventStateDone(contents bytes.Buffer, eventObject *EventObject) (err error) {
	err = fmt.Errorf("this is an error")
	failureWaitCalled++
	return
}

type DefaultMockHttpGateway struct {
	handler HandleRespFunc
}

func (m DefaultMockHttpGateway) Execute(method string) (interface{}, error) {
	resp := &http.Response{
		StatusCode: 200,
	}
	return m.handler(resp)
}

func (m DefaultMockHttpGateway) ExecuteFunc(method string, handler HandleRespFunc) (interface{}, error) {
	resp := &http.Response{
		StatusCode: 200,
	}
	return handler(resp)
}

func (m DefaultMockHttpGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (*http.Response, error) {
	return nil, nil
}

var _ = Describe("toggle cc job", func() {
	var (
		restSuccessCalled    int
		restFailureCalled    int
		successToggleCalled  int
		failureToggleCalled  int
		successCreaterCalled int
		failureCreaterCalled int
		successString        string = `{"state":"done"}`
		failureString        string = `{"state":"notdone"}`
		failTryExitCount     int    = 5
		endlessLoopFlag      bool   = false
		mockHttpGateway             = DefaultMockHttpGateway{handler: func(resp *http.Response) (interface{}, error) {
			return "success", nil
		}}
	)
	restSuccess := func(resp *http.Response) (interface{}, error) {
		resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
		restSuccessCalled++
		return resp.Body, nil
	}
	restFailure := func(resp *http.Response) (interface{}, error) {
		resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
		restFailureCalled++
		err := fmt.Errorf("")
		return resp.Body, err
	}
	restNotDone := func(resp *http.Response) (interface{}, error) {
		resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
		restFailureCalled++
		_ = failTryExitCount
		if restFailureCalled > failTryExitCount {
			resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
			endlessLoopFlag = true
		}
		return resp.Body, nil
	}

	successJobToggleMock := func(serverUrl, username, password string) (res string, err error) {
		successToggleCalled++
		return
	}

	failureJobToggleMock := func(serverUrl, username, password string) (res string, err error) {
		failureToggleCalled++
		return
	}

	successTaskCreater := func(method string, gateway HttpGateway, handleRespFunc HandleRespFunc) (task EventTasker) {
		task = EventTasker(SuccessMockEventTasker{})
		successCreaterCalled++
		return
	}

	failureTaskCreater := func(method string, gateway HttpGateway, handleRespFunc HandleRespFunc) (task EventTasker) {
		task = &FailureMockEventTasker{}
		failureCreaterCalled++
		return
	}

	Describe("Task", func() {
		Context("successful call", func() {
			var task EventTasker
			BeforeEach(func() {
				task = &Task{
					Method:              "GET",
					HttpGateway:         mockHttpGateway,
					HttpResponseHandler: restSuccess,
				}
			})

			It("Should return nil error on valid arguments", func() {
				eventObject := &EventObject{}
				bbf := bytes.NewBuffer([]byte(successString))
				err := task.WaitForEventStateDone(*bbf, eventObject)
				Ω(err).Should(BeNil())
			})

			It("Should return nil error and return if rest endpoint returns done status", func() {
				eventObject := &EventObject{}
				bbf := bytes.NewBuffer([]byte(failureString))
				err := task.WaitForEventStateDone(*bbf, eventObject)
				Ω(err).Should(BeNil())
			})
		})

		Context("status not done call", func() {
			var task EventTasker
			BeforeEach(func() {
				endlessLoopFlag = false

				task = &Task{
					Method:              "GET",
					HttpGateway:         &DefaultMockHttpGateway{},
					HttpResponseHandler: restNotDone,
				}
			})

			It("Should loop endlessly if done is never returned", func() {
				eventObject := &EventObject{}
				bbf := bytes.NewBuffer([]byte(failureString))
				err := task.WaitForEventStateDone(*bbf, eventObject)
				Ω(err).Should(BeNil())
				Ω(endlessLoopFlag).Should(BeTrue())
			})
		})

		Context("failed call", func() {
			var task EventTasker
			BeforeEach(func() {
				endlessLoopFlag = false

				task = &Task{
					Method:              "GET",
					HttpGateway:         mockHttpGateway,
					HttpResponseHandler: restFailure,
				}
			})

			It("Should return non nil error for bad event object", func() {
				bbf := bytes.NewBuffer([]byte(""))
				err := task.WaitForEventStateDone(*bbf, nil)
				Ω(err).ShouldNot(BeNil())
			})

			It("Should loop endlessly if done is never returned", func() {
				eventObject := &EventObject{}
				bbf := bytes.NewBuffer([]byte(failureString))
				err := task.WaitForEventStateDone(*bbf, eventObject)
				Ω(err).ShouldNot(BeNil())
			})
		})

	})

	Describe("CloudController", func() {
		Context("successful call", func() {
			var cc *CloudController
			BeforeEach(func() {
				resHandler := func(resp *http.Response) (val interface{}, err error) { return }
				cc = NewCloudController("", "", "", "", "", resHandler)
				cc.JobToggler = JobTogglerAdapter(successJobToggleMock)
				cc.NewEventTaskCreater = EvenTaskCreaterAdapter(successTaskCreater)
				successWaitCalled, failureWaitCalled, successToggleCalled, failureToggleCalled, successCreaterCalled, failureCreaterCalled, restSuccessCalled, restFailureCalled = 0, 0, 0, 0, 0, 0, 0, 0
			})
			AfterEach(func() {
				successWaitCalled, failureWaitCalled, successToggleCalled, failureToggleCalled, successCreaterCalled, failureCreaterCalled, restSuccessCalled, restFailureCalled = 0, 0, 0, 0, 0, 0, 0, 0
			})

			Context("ToggleJobs (with an 's') method", func() {
				It("Should call through the entire chain if there is no error", func() {
					cc.ToggleJobs(CloudControllerJobs([]string{"jobA", "someurl.com"}))
					Ω(successToggleCalled).Should(BeNumerically(">", 0))
					Ω(successCreaterCalled).Should(BeNumerically(">", 0))
					Ω(successWaitCalled).Should(BeNumerically(">", 0))
				})
			})

			Context("ToggleJob method", func() {
				Context("when a call to task.WaitForEventStateDone internally returns error", func() {
					BeforeEach(func() {
						cc.NewEventTaskCreater = EvenTaskCreaterAdapter(failureTaskCreater)
					})

					It("should return an error from ToggleJob", func() {
						err := cc.ToggleJob("jobA", "someurl.com", 1)
						Ω(err).ShouldNot(BeNil())
					})
				})

				It("Should call through the entire chain if there is no error", func() {
					cc.ToggleJob("jobA", "someurl.com", 1)
					Ω(successToggleCalled).Should(BeNumerically(">", 0))
					Ω(successCreaterCalled).Should(BeNumerically(">", 0))
					Ω(successWaitCalled).Should(BeNumerically(">", 0))
				})
			})

		})

		Context("failed call", func() {
			var cc *CloudController
			BeforeEach(func() {
				cc = NewCloudController("", "", "", "", "", func(resp *http.Response) (val interface{}, err error) { return })
				cc.JobToggler = JobTogglerAdapter(failureJobToggleMock)
				cc.NewEventTaskCreater = EvenTaskCreaterAdapter(failureTaskCreater)
			})
			Context("ToggleJobs (with an 's') method", func() {
				It("Should not call through the entire chain if there is an error", func() {
					cc.ToggleJobs(CloudControllerJobs([]string{"jobA", "someurl.com"}))
					Ω(successToggleCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successCreaterCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successWaitCalled).ShouldNot(BeNumerically(">", 0))
				})
			})
			Context("ToggleJob method", func() {
				It("Should not call through the entire chain if there is an error", func() {
					cc.ToggleJob("jobA", "someurl.com", 1)
					Ω(successToggleCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successCreaterCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successWaitCalled).ShouldNot(BeNumerically(">", 0))
				})
			})
		})

		Context("partial failed call", func() {
			var cc *CloudController
			BeforeEach(func() {
				resHandler := func(resp *http.Response) (val interface{}, err error) { return }
				cc = NewCloudController("", "", "", "", "", resHandler)
				cc.JobToggler = JobTogglerAdapter(successJobToggleMock)
				cc.NewEventTaskCreater = EvenTaskCreaterAdapter(failureTaskCreater)
			})
			Context("ToggleJobs (with an 's') method", func() {
				It("Should not call through the entire chain if there is an error", func() {
					cc.ToggleJobs(CloudControllerJobs([]string{"jobA", "someurl.com"}))
					Ω(successToggleCalled).Should(BeNumerically(">", 0))
					Ω(successCreaterCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successWaitCalled).ShouldNot(BeNumerically(">", 0))
				})
			})
			Context("ToggleJob method", func() {
				It("Should not call through the entire chain if there is an error", func() {
					cc.ToggleJob("jobA", "someurl.com", 1)
					Ω(successToggleCalled).Should(BeNumerically(">", 0))
					Ω(successCreaterCalled).ShouldNot(BeNumerically(">", 0))
					Ω(successWaitCalled).ShouldNot(BeNumerically(">", 0))
				})
			})
		})
	})

	Describe("ToggleCCJobRunner", func() {
		Context("successful call", func() {
			var (
				username  string = "usertest"
				password  string = "passwrdtest"
				serverUrl string = "someurl.com"
			)
			It("Should return nil error", func() {
				NewToggleGateway = func(serverUrl, username, password string) HttpGateway {
					return &MockSuccessGateway{Handler: ToggleCCHandler}
				}
				_, err := ToggleCCJobRunner(username, password, serverUrl)
				Ω(err).Should(BeNil())
			})

			It("Should return redirectUrl", func() {
				NewToggleGateway = func(serverUrl, username, password string) HttpGateway {
					return &MockSuccessGateway{Handler: ToggleCCHandler}
				}
				msg, _ := ToggleCCJobRunner(username, password, serverUrl)
				Ω(msg).Should(Equal(redirectUrl))
			})
		})

		Context("failure call", func() {
			var (
				username  string = "usertest"
				password  string = "passwrdtest"
				serverUrl string = "someurl.com"
			)
			It("Should return error on non 302 http code", func() {
				NewToggleGateway = func(serverUrl, username, password string) HttpGateway {
					return &MockFailerGateway{Handler: ToggleCCHandler}
				}
				_, err := ToggleCCJobRunner(username, password, serverUrl)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
