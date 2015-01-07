package backup_test

import (
	"os"

	. "github.com/pivotalservices/cfops/backup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get_deployment_name", func() {
	Describe("GetDeploymentName function", func() {
		Context("when given a valid deployments.json", func() {
			var (
				jsonObj []DeploymentObject
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/deployments.json")
				jsonObj, _ = ReadAndUnmarshalDeploymentName(fileRef)
			})

			AfterEach(func() {
			})

			It("Should return nil error, correct deployment name", func() {
				name, err := GetDeploymentName(jsonObj)
				Ω(err).Should(BeNil())
				Ω(name).Should(Equal("cf-c3aad3ad484c438ebc40"))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					GetDeploymentName(jsonObj)
				}).ShouldNot(Panic())
			})
		})
	})

	Describe("DeploymentParser struct", func() {
		Context("when given a valid deployments.json", func() {
			var (
				parser  *DeploymentParser
				jsonObj []DeploymentObject
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/deployments.json")
				jsonObj, _ = ReadAndUnmarshalDeploymentName(fileRef)

				parser = &DeploymentParser{}
			})

			AfterEach(func() {
			})

			It("Should return nil error, correct deployment name", func() {
				name, err := parser.Parse(jsonObj)
				Ω(err).Should(BeNil())
				Ω(name).Should(Equal("cf-c3aad3ad484c438ebc40"))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})
	})
})
