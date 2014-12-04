package api_test

import (
	"github.com/cloudfoundry-community/gogobosh/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("models.TaskStatus", func() {
	It("GetTaskStatus returns models.TaskStatus{}", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/tasks/1",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "id": 1,
				  "state": "done",
				  "description": "create release",
				  "timestamp": 1390068518,
				  "result": "Created release cf/153",
				  "user": "admin"
				}`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		task, apiResponse := repo.GetTaskStatus(1)
		
		Expect(task.ID).To(Equal(1))
		Expect(task.State).To(Equal("done"))
		Expect(task.Description).To(Equal("create release"))
		Expect(task.TimeStamp).To(Equal(1390068518))
		Expect(task.Result).To(Equal("Created release cf/153"))
		Expect(task.User).To(Equal("admin"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	It("() returns []models.TaskStatus{}", func() {
		request := testhelpers.NewDirectorTestRequest(testhelpers.TestRequest{
			Method: "GET",
			Path:   "/tasks",
			Response: testhelpers.TestResponse{
				Status: http.StatusOK,
				Body: `[{
				    "id": 2,
				    "state": "done",
				    "description": "create release",
				    "timestamp": 1390068525,
				    "result": "Created release 'etcd/3'",
				    "user": "admin"
				  },
				  {
				    "id": 1,
				    "state": "done",
				    "description": "create release",
				    "timestamp": 1390068518,
				    "result": "Created release 'cf/153'",
				    "user": "admin"
				  }
				]`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		tasks, apiResponse := repo.GetTaskStatuses()

		Expect(len(tasks)).To(Equal(2))

		task := tasks[1]
		Expect(task.ID).To(Equal(1))
		Expect(task.State).To(Equal("done"))
		Expect(task.Description).To(Equal("create release"))
		Expect(task.TimeStamp).To(Equal(1390068518))
		Expect(task.Result).To(Equal("Created release 'cf/153'"))
		Expect(task.User).To(Equal("admin"))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})

	// verbose: true/false - show internal tasks

	// limit: nil or integer limit
	
	// states: all, processing,cancelling,queued ("running"), or specific list
	XIt("GetRunningTaskStatuses", func() {
		// states: processing,cancelling,queued
	})

})
