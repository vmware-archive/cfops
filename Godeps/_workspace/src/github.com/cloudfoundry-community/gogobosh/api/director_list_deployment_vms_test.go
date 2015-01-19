package api_test

import (
	"github.com/cloudfoundry-community/gogobosh/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("simple list of vms", func() {
	It("GET /deployments/$name/vms to return []models.DeploymentVM{}", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/deployments/cf-warden/vms",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `[
				  {
				    "agent_id": "b11f259c-79dd-4d6d-8aa5-5969d569a2a6",
				    "cid": "vm-8a03a314-6f16-45f6-a377-1a36e763ee45",
				    "job": "ha_proxy_z1",
				    "index": 0
				  },
				  {
				    "agent_id": "5c7708c9-1631-48b3-9833-6b7d0f6c6cd4",
				    "cid": "vm-37926289-487d-4ee9-b556-9684350d1d14",
				    "job": "login_z1",
				    "index": 0
				  }
				]`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		deploymentVMs, apiResponse := repo.ListDeploymentVMs("cf-warden")

		Expect(len(deploymentVMs)).To(Equal(2))

		vm := deploymentVMs[0]
		Expect(vm.JobName).To(Equal("ha_proxy_z1"))
		Expect(vm.Index).To(Equal(0))
		Expect(vm.AgentID).To(Equal("b11f259c-79dd-4d6d-8aa5-5969d569a2a6"))
		Expect(vm.VMCid).To(Equal("vm-8a03a314-6f16-45f6-a377-1a36e763ee45"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})
