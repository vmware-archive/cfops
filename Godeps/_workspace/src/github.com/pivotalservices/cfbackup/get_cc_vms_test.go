package cfbackup_test

import (
	"os"

	. "github.com/pivotalservices/cfbackup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get_cc_vms", func() {
	Describe("GetCCVMs function", func() {
		var (
			jsonObj []VMObject
		)
		Context("when given a valid deployment_vms.json", func() {
			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/deployment_vms.json")
				jsonObj, _ = ReadAndUnmarshalVMObjects(fileRef)
			})

			It("Should return nil error, correct cc jobs", func() {
				vms, err := GetCCVMs(jsonObj)
				Ω(err).Should(BeNil())
				Ω(vms).Should(HaveLen(1))
				Ω(vms[0]).Should(Equal("cloud_controller-partition-7bc61fd2fa9d654696df"))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					GetCCVMs(jsonObj)
				}).ShouldNot(Panic())
			})
		})

		Context("when cc job is not found", func() {
			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/deployment_without_cc.json")
				jsonObj, _ = ReadAndUnmarshalVMObjects(fileRef)
			})
			It("Should return error", func() {
				vms, err := GetCCVMs(jsonObj)
				Ω(err).ShouldNot(BeNil())
				Ω(vms).Should(BeNil())
			})
		})
	})

	Describe("CloudControllerDeploymentParser struct", func() {
		Context("when given a valid deployment_vms.json", func() {
			var (
				parser  *CloudControllerDeploymentParser
				jsonObj []VMObject
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/deployment_vms.json")
				jsonObj, _ = ReadAndUnmarshalVMObjects(fileRef)

				parser = &CloudControllerDeploymentParser{}
			})

			AfterEach(func() {
			})

			It("Should return nil error, correct cc job", func() {
				vms, err := parser.Parse(jsonObj)
				Ω(err).Should(BeNil())
				Ω(vms).Should(HaveLen(1))
				Ω(vms[0]).Should(Equal("cloud_controller-partition-7bc61fd2fa9d654696df"))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})
	})
})
