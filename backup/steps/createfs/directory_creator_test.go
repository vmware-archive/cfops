package createfs_test

import (
	"fmt"
	"os"

	. "github.com/pivotalservices/cfops/backup/steps/createfs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var callCountCutoff int = 1
var twiceCounter int = 0
var errorCounter int = 0
var successCounter int = 0

func mockMkDirErrorIfTwiceCalled(path string, perm os.FileMode) (err error) {

	if twiceCounter >= callCountCutoff {
		errorCounter++
		err = fmt.Errorf("mock error")
	} else {
		successCounter++
	}
	twiceCounter++
	return
}

func mockMkDirError(path string, perm os.FileMode) (err error) {
	errorCounter++
	err = fmt.Errorf("mock error")
	return
}

func mockMkDirSuccess(path string, perm os.FileMode) (err error) {
	successCounter++
	return
}

var _ = Describe("Backup", func() {
	directoryList := []string{"testdir1", "testdir2", "testdir3", "testdir4"}

	BeforeEach(func() {
		twiceCounter = 0
		errorCounter = 0
		successCounter = 0
	})

	AfterEach(func() {
		twiceCounter = 0
		errorCounter = 0
		successCounter = 0
	})

	Context("MultiDirectoryCreate function", func() {
		It("Should return nil error on success and have called the mkdir functor the proper amount of times", func() {
			controlCallCount := len(directoryList)
			err := MultiDirectoryCreate(directoryList, mockMkDirSuccess)
			Ω(err).Should(BeNil())
			Expect(successCounter).To(Equal(controlCallCount))
		})

		It("Should return not nil error on error", func() {
			controlCallCount := len(directoryList)
			err := MultiDirectoryCreate(directoryList, mockMkDirError)
			Ω(err).ShouldNot(BeNil())
			Expect(errorCounter).NotTo(Equal(controlCallCount))
		})

		It("Should exit the call loop when a error is encountered", func() {
			controlCallCount := callCountCutoff
			err := MultiDirectoryCreate(directoryList, mockMkDirErrorIfTwiceCalled)
			Ω(err).ShouldNot(BeNil())
			Expect(successCounter).To(Equal(controlCallCount))
		})
	})

	Context("DirectoryCreate function", func() {
		It("Should return nil error on success and have called the mkdir functor only once", func() {
			controlCallCount := 1
			err := DirectoryCreate(directoryList[0], mockMkDirSuccess)
			Ω(err).Should(BeNil())
			Expect(successCounter).To(Equal(controlCallCount))
		})

		It("Should return not nil error on error and have called the mkdir functor only once", func() {
			controlCallCount := 1
			err := DirectoryCreate(directoryList[0], mockMkDirError)
			Ω(err).ShouldNot(BeNil())
			Expect(errorCounter).To(Equal(controlCallCount))
		})
	})
})
