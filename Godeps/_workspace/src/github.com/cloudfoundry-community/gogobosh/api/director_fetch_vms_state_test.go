package api_test

import (
	"github.com/cloudfoundry-community/gogobosh/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("parse full vms task output", func() {
	It("GET /deployments/cf-warden/vms?format=full to return Director{}", func() {
		vmsRequest := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/deployments/cf-warden/vms?format=full",
			Response: testhelpers.TestResponse{
				Status: http.StatusFound,
				Header: http.Header{
					"Location":{"https://some.host/tasks/12"},
				},
			},
		})
		taskOutputRequest := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/tasks/12/output?type=result",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `{"vm_cid":"vm-a1a3d634-367d-4b75-940c-ef7742a970d9","ips":["10.244.1.14"],"dns":[],"agent_id":"c0da6161-e66f-4910-a0eb-dc6fc19b4b25","job_name":"hm9000_z1","index":0,"job_state":"running","resource_pool":"medium_z1","vitals":{"load":["0.11","0.21","0.18"],"cpu":{"user":"1.5","sys":"2.8","wait":"0.1"},"mem":{"percent":"46.8","kb":"2864212"},"swap":{"percent":"0.0","kb":"0"},"disk":{"system":{"percent":null},"persistent":{"percent":"1"}}},"resurrection_paused":false}
				{"vm_cid":"vm-affdbbdb-b91e-4838-b068-f1a057242169","ips":["10.244.0.38"],"dns":[],"agent_id":"bec309f8-0e2d-4843-9db3-a419adab4d38","job_name":"etcd_leader_z1","index":0,"job_state":"running","resource_pool":"medium_z1","vitals":{"load":["0.13","0.22","0.18"],"cpu":{"user":"0.4","sys":"2.0","wait":"0.1"},"mem":{"percent":"46.8","kb":"2863012"},"swap":{"percent":"0.0","kb":"0"},"disk":{"system":{"percent":null},"persistent":{"percent":"1"}}},"resurrection_paused":false}
				`}})

		ts, handler, repo := createDirectorRepo(
			vmsRequest,
			taskTestRequest(12, "queued"),
			taskTestRequest(12, "processing"),
			taskTestRequest(12, "done"),
			taskOutputRequest)
		defer ts.Close()

		vmStatuses, apiResponse := repo.FetchVMsStatus("cf-warden")

		/* TODO convert vmStatuses to a chan and pluck first item from chan */
		Expect(len(vmStatuses)).To(Equal(2))
		vmStatus := vmStatuses[0]
		Expect(vmStatus.JobName).To(Equal("hm9000_z1"))
		Expect(vmStatus.Index).To(Equal(0))
		Expect(vmStatus.JobState).To(Equal("running"))
		Expect(vmStatus.VMCid).To(Equal("vm-a1a3d634-367d-4b75-940c-ef7742a970d9"))
		Expect(vmStatus.AgentID).To(Equal("c0da6161-e66f-4910-a0eb-dc6fc19b4b25"))
		Expect(vmStatus.ResourcePool).To(Equal("medium_z1"))
		Expect(vmStatus.ResurrectionPaused).To(Equal(false))

		Expect(len(vmStatus.IPs)).To(Equal(1))
		Expect(vmStatus.IPs[0]).To(Equal("10.244.1.14"))

		Expect(len(vmStatus.DNSs)).To(Equal(0))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})