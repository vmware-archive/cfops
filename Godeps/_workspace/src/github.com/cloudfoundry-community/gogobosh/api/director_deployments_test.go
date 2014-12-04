package api_test

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/gogobosh/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployments", func() {
	It("GetDeployments() - list of deployments", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/deployments",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `[
				  {
				    "name": "cf-warden",
				    "releases": [
				      {
				        "name": "cf",
				        "version": "153"
				      }
				    ],
				    "stemcells": [
				      {
				        "name": "bosh-stemcell",
				        "version": "993"
				      }
				    ]
				  }
				]`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		deployments, apiResponse := repo.GetDeployments()

		deployment := deployments[0]
		Expect(deployment.Name).To(Equal("cf-warden"))

		deploymentRelease := deployment.Releases[0]
		Expect(deploymentRelease.Name).To(Equal("cf"))
		Expect(deploymentRelease.Version).To(Equal("153"))

		deploymentStemcell := deployment.Stemcells[0]
		Expect(deploymentStemcell.Name).To(Equal("bosh-stemcell"))
		Expect(deploymentStemcell.Version).To(Equal("993"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	It("GetDeployment(name) - get deployment, including manifest", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/deployments/cf-warden",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `{
					"manifest": "name: cf-warden"
				}`,
			}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		manifest, apiResponse := repo.GetDeploymentManifest("cf-warden")
		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))

		Expect(manifest.Name).To(Equal("cf-warden"))
	})

	It("DeleteDeployment(name) forcefully", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "DELETE",
			Path:   "/deployments/cf-warden?force=true",
			Response: testhelpers.TestResponse{
				Status: http.StatusFound,
				Header: http.Header{
					"Location": {"https://some.host/tasks/20"},
				},
			}})
		ts, handler, repo := createDirectorRepo(
			request,
			taskTestRequest(20, "queued"),
			taskTestRequest(20, "processing"),
			taskTestRequest(20, "done"),
		)
		defer ts.Close()

		apiResponse := repo.DeleteDeployment("cf-warden")

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})

// Shared helper for asserting that a /tasks/ID is requested and returns a models.TaskStatus response
func taskTestRequest(taskID int, state string) testhelpers.TestRequest {
	baseJSON := `{
	  "id": %d,
	  "state": "%s",
	  "description": "some task",
	  "timestamp": 1390174354,
	  "result": null,
	  "user": "admin"
	}`
	return testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
		Method: "GET",
		Path:   fmt.Sprintf("/tasks/%d", taskID),
		Response: testhelpers.TestResponse{
			Status: http.StatusOK,
			Body:   fmt.Sprintf(baseJSON, taskID, state),
		},
	})
}
