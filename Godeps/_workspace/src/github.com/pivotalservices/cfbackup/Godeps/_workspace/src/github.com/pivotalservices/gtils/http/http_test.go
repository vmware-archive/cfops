package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
)

var (
	roundTripSuccess bool
	requestCatcher   *http.Request
	handlerSuccess   bool
	successString    string = `{"state":"done"}`
	failureString    string = `{"state":"notdone"}`
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

type MockRoundTripper struct {
}

func (roundTripper *MockRoundTripper) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: 200,
	}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}

	if !roundTripSuccess {
		resp.StatusCode = 500
		resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
		err = errors.New("Mock error")
	}
	*requestCatcher = *request
	return
}

func MockHandlerFunc(resp *http.Response) (val interface{},
	err error) {
	if !handlerSuccess {
		return nil, errors.New("Mock error")
	}
	return "Success", nil
}

var _ = Describe("Http", func() {
	var (
		handler func (resp *http.Response) (val interface{}, err error)
		gateway HttpGateway
	)
	BeforeEach(func() {
		requestCatcher = &http.Request{}
		handler = MockHandlerFunc
		gateway = NewHttpGateway("http://endpoint/test", "username", "password", "contentType", handler)
		NewRoundTripper = func() http.RoundTripper {
			return &MockRoundTripper{}
		}
	})

	Context("The http is request and handled successfully", func() {
		BeforeEach(func() {
			roundTripSuccess = true
			handlerSuccess = true
		})
		It("Should return nil error on success", func() {
			_, err := gateway.Execute("Get")
			Ω(err).Should(BeNil())
		})
		It("Should execute correct request", func() {
			val, _ := gateway.Execute("Get")
			Ω(requestCatcher.URL.Host).Should(Equal("endpoint"))
			Ω(requestCatcher.Method).Should(Equal("Get"))
			Ω(requestCatcher.Header["Content-Type"][0]).Should(Equal("contentType"))
			Ω(requestCatcher.Header["Authorization"][0]).Should(Equal("Basic dXNlcm5hbWU6cGFzc3dvcmQ="))
			Ω(val).Should(Equal("Success"))
		})
	})

	Context("The round trip request failed", func() {
		BeforeEach(func() {
			roundTripSuccess = false
			handlerSuccess = true
		})
		It("Should return error", func() {
			_, err := gateway.Execute("Get")
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("The handler failed", func() {
		BeforeEach(func() {
			roundTripSuccess = true
			handlerSuccess = false
		})
		It("Should return error", func() {
			_, err := gateway.Execute("Get")
			Ω(err).ShouldNot(BeNil())
		})
	})
	Describe("Upload function", func() {
		Context("call to endpoint is successful", func() {
			BeforeEach(func() {
				roundTripSuccess = true
			})

			It("Should return nil error and a valid response", func() {
				fileRef, _ := os.Open("fixtures/installation.json")
				res, err := gateway.Upload("installation[file]", "installation.json", fileRef, nil)
				Ω(err).Should(BeNil())
				Ω(res).ShouldNot(BeNil())
				Ω(res.StatusCode).Should(Equal(200))
			})
		})

		Context("call to endpoint is not successful", func() {
			BeforeEach(func() {
				roundTripSuccess = false
			})

			It("Should return non-nil error and a non 200 statuscode", func() {
				fileRef, _ := os.Open("fixtures/installation.json")
				res, err := gateway.Upload("installation[file]", "installation.json", fileRef, nil)
				Ω(err).ShouldNot(BeNil())
				Ω(res.StatusCode).ShouldNot(Equal(200))
			})
		})
	})
})
