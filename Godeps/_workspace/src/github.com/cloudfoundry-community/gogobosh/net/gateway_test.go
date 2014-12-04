package net

import (
	"github.com/cloudfoundry-community/gogobosh/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"runtime"
)

var _ = Describe("Gateway", func() {
	It("NewRequest successfully", func() {
		gateway := NewDirectorGateway()

		request, apiResponse := gateway.NewRequest("GET", "https://example.com/v2/apps", "admin", "admin", nil)

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(request.HttpReq.Header.Get("Authorization")).To(Equal("Basic YWRtaW46YWRtaW4="))
		Expect(request.HttpReq.Header.Get("accept")).To(Equal("application/json"))
		Expect(request.HttpReq.Header.Get("User-Agent")).To(Equal("gogobosh "+constants.Version+" / "+runtime.GOOS))
	})
})
