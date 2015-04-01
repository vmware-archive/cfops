package http_test

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/mock"
)

const (
	paramName = "installation[file]"
	fileName  = "installation.json"
)

var _ = Describe("Multipart", func() {
	Describe("MultiPartBody", func() {
		var (
			multipartConstructor MultiPartBodyFunc = MultiPartBody
		)
		Context("Construct the multipart body successfully", func() {
			It("Should return nil error", func() {
				fileRef, _ := os.Open("fixtures/installation.json")
				body, contentType, err := multipartConstructor("installation[file]", "installation.json", fileRef, nil)
				Ω(err).Should(BeNil())
				Ω(body).ShouldNot(BeNil())
				Ω(contentType).ShouldNot(BeNil())
			})
		})

		Context("Construct the multipart body failed", func() {
			It("Should return error when file is missing", func() {
				fileRef, _ := os.Open("fixtures/installa.json")
				_, _, err := multipartConstructor("installation[file]", "installation.json", fileRef, nil)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Describe("MultiPartUpload", func() {
		var (
			httpServerMock *mock.HttpServer = &mock.HttpServer{}
			request        *http.Request
		)

		BeforeEach(func() {
			httpServerMock.Setup()
			httpServerMock.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				request = r
			})
		})

		AfterEach(func() {
			httpServerMock.Teardown()
		})

		It("should send the file to the server", func() {
			filePath := fmt.Sprintf("fixtures/%s", fileName)
			conn := ConnAuth{
				Url: httpServerMock.Server.URL,
			}
			fileRef, _ := os.Open(filePath)
			res, err := MultiPartUpload(conn, paramName, fileName, fileRef, nil)
			Ω(err).Should(BeNil())
			Ω(res.StatusCode).Should(Equal(200))
			Ω(request.Method).Should(Equal("POST"))
			Ω(res).ShouldNot(BeNil())
			Ω(request.ContentLength).Should(Equal(res.Request.ContentLength))
		})
	})
})
