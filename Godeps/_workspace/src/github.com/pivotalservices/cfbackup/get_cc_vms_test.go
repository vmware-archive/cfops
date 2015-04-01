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
				_, err := GetCCVMs(jsonObj)
				Ω(err).Should(BeNil())
			})

			It("Should have 4 cc correct jobs", func() {
				vms, _ := GetCCVMs(jsonObj)
				Ω(vms).Should(HaveLen(4))
				Ω(vms[0].Index).Should(Equal(0))
				Ω(vms[1].Index).Should(Equal(1))
				Ω(vms[2].Index).Should(Equal(0))
				Ω(vms[3].Index).Should(Equal(1))
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

			It("Should return nil error", func() {
				_, err := parser.Parse(jsonObj)
				Ω(err).Should(BeNil())
			})

			It("Should return 4 jobs", func() {
				vms, _ := parser.Parse(jsonObj)
				Ω(vms).Should(HaveLen(4))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})
	})
})
