package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/gtils/bosh"
)

var _ = Describe("API", func() {
	var (
		pathParams  map[string]string = map[string]string{"deployment": "cf-1", "job": "cloudcontroller", "index": "1"}
		queryParams map[string]string = map[string]string{"state": "started"}
		ip          string            = "10.10.10.10"
		port        int               = 25555
		path        string            = "deployments/{deployment}/{job}/{index}"
	)
	Describe("Generate Director url", func() {
		Context("With valid configureation", func() {
			It("should return nil error", func() {
				_, err := ParseUrl(ip, port, path, pathParams, queryParams)
				Expect(err).Should(BeNil())
			})
			It("should compose a valid url", func() {
				url, _ := ParseUrl(ip, port, path, pathParams, queryParams)
				Expect(url).To(Equal("https://10.10.10.10:25555/deployments/cf-1/cloudcontroller/1?state=started"))
			})
		})
	})
})
