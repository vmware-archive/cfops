package api_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-community/gogobosh/api"
	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
	"github.com/cloudfoundry-community/gogobosh/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get director info", func() {
	It("GET /info to return Director{}", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/info",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "name": "Bosh Lite Director",
				  "uuid": "bd462a15-213d-448c-aa5b-66624dad3f0e",
				  "version": "1.5.0.pre.1657 (14bc162c)",
				  "user": "admin",
				  "cpi": "warden",
				  "features": {
				    "dns": {
				      "status": false,
				      "extras": {
				        "domain_name": "bosh"
				      }
				    },
				    "compiled_package_cache": {
				      "status": true,
				      "extras": {
				        "provider": "local"
				      }
				    },
				    "snapshots": {
				      "status": false
				    }
				  }
				}`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		info, apiResponse := repo.GetInfo()

		Expect(info.Name).To(Equal("Bosh Lite Director"))
		Expect(info.UUID).To(Equal("bd462a15-213d-448c-aa5b-66624dad3f0e"))
		Expect(info.Version).To(Equal("1.5.0.pre.1657 (14bc162c)"))
		Expect(info.User).To(Equal("admin"))
		Expect(info.CPI).To(Equal("warden"))
		Expect(info.DNSEnabled).To(Equal(false))
		Expect(info.DNSDomainName).To(Equal("bosh"))
		Expect(info.CompiledPackageCacheEnabled).To(Equal(true))
		Expect(info.CompiledPackageCacheProvider).To(Equal("local"))
		Expect(info.SnapshotsEnabled).To(Equal(false))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})

func createDirectorRepo(reqs ...testhelpers.TestRequest) (ts *httptest.Server, handler *testhelpers.TestHandler, repo api.DirectorRepository) {
	ts, handler = testhelpers.NewTLSServer(reqs)
	config := &models.Director{
		TargetURL: ts.URL,
		Username:  "admin",
		Password:  "admin",
	}
	gateway := net.NewDirectorGateway()
	repo = api.NewBoshDirectorRepository(config, gateway)
	return
}
