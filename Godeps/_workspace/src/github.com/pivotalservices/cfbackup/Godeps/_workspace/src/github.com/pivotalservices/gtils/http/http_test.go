package http_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
)

var (
	requestCatcher   *http.Request
	roundTripSuccess bool
	httpEntity       HttpRequestEntity = HttpRequestEntity{
		Url:         "http://endpoint/test",
		Username:    "username",
		Password:    "password",
		ContentType: "contentType",
	}
)

type MockRoundTripper struct {
}

func (roundTripper *MockRoundTripper) RoundTrip(request *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: 200,
	}
	if !roundTripSuccess {
		err = errors.New("Mock error")
	}
	*requestCatcher = *request
	return
}

var _ = Describe("Http", func() {
	var _ = Describe("Request Function", func() {
		BeforeEach(func() {
			requestCatcher = &http.Request{}
			NewRoundTripper = func() http.RoundTripper {
				return &MockRoundTripper{}
			}
		})

		Context("The http request successfully", func() {
			BeforeEach(func() {
				roundTripSuccess = true
			})
			It("Should return nil error on success", func() {
				_, err := Request(httpEntity, "Get", nil)
				Ω(err).Should(BeNil())
			})
			It("Should execute correct request", func() {
				resp, _ := Request(httpEntity, "Get", nil)
				Ω(requestCatcher.URL.Host).Should(Equal("endpoint"))
				Ω(requestCatcher.Method).Should(Equal("Get"))
				Ω(requestCatcher.Header["Content-Type"][0]).Should(Equal("contentType"))
				Ω(requestCatcher.Header["Authorization"][0]).Should(Equal("Basic dXNlcm5hbWU6cGFzc3dvcmQ="))
				Ω(resp.StatusCode).Should(Equal(200))

			})
		})

		Context("The round trip request failed", func() {
			BeforeEach(func() {
				roundTripSuccess = false
			})
			It("Should return error", func() {
				_, err := Request(httpEntity, "Get", nil)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	var _ = Describe("Http Gateway", func() {
		var (
			gateway = NewHttpGateway()
		)
		BeforeEach(func() {
			roundTripSuccess = true
			requestCatcher = &http.Request{}
			NewRoundTripper = func() http.RoundTripper {
				return &MockRoundTripper{}
			}
		})

		Context("Calling gateway get method", func() {
			It("Http Request should call Get method", func() {
				request := gateway.Get(httpEntity)
				request()
				Ω(requestCatcher.Method).Should(Equal("GET"))
			})
		})

		Context("Calling gateway post method", func() {
			It("Http Request should call post method", func() {
				request := gateway.Post(httpEntity, nil)
				request()
				Ω(requestCatcher.Method).Should(Equal("POST"))
			})
		})
		Context("Calling gateway put method", func() {
			It("Http Request should call put method", func() {
				request := gateway.Put(httpEntity, nil)
				request()
				Ω(requestCatcher.Method).Should(Equal("PUT"))
			})
		})

	})
})
