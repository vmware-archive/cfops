package http_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/http"
)

var _ = Describe("Multipart", func() {
	var (
		multipartConstructor MultiPartBodyFunc = MultiPartBody
	)
	Context("Construct the multipart body successfully", func() {
		It("Should return nil error", func() {
			fileRef, _ := os.Open("fixtures/installation.json")
			body, err := multipartConstructor("installation[file]", "installation.json", fileRef, nil)
			Ω(err).Should(BeNil())
			Ω(body).ShouldNot(BeNil())
		})
	})

	Context("Construct the multipart body failed", func() {
		It("Should return error when file is missing", func() {
			fileRef, _ := os.Open("fixtures/installa.json")
			_, err := multipartConstructor("installation[file]", "installation.json", fileRef, nil)
			Ω(err).ShouldNot(BeNil())
		})
	})

})
