package cfbackup_test

import (
	"os"

	. "github.com/pivotalservices/cfbackup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get_password_ip", func() {

	Describe("Ops Manager Elastic Runtime v1.3", func() {
		var installationSettingsFilePath = "fixtures/installation-settings-1-3.json"
		testGetPasswordWithVersionSpecificFile(installationSettingsFilePath)
	})

	Describe("Ops Manager Elastic Runtime v1.4 file variant with getpassword IP index error", func() {
		var installationSettingsFilePath = "fixtures/installation-settings-1-4-variant.json"
		testGetPasswordWithVersionSpecificFile(installationSettingsFilePath)
	})

	Describe("Ops Manager Elastic Runtime v1.4", func() {
		var installationSettingsFilePath = "fixtures/installation-settings-1-4.json"
		testGetPasswordWithVersionSpecificFile(installationSettingsFilePath)
	})
})

func testGetPasswordWithVersionSpecificFile(installationSettingsFilePath string) {
	Describe("GetDeploymentName function", func() {
		Context("when given a valid installation.json", func() {
			var (
				jsonObj InstallationCompareObject
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)
			})

			AfterEach(func() {
			})

			It("Should return the proper installation_name and nil error", func() {
				name, err := GetDeploymentName(jsonObj)
				Ω(err).Should(BeNil())
				Ω(name).Should(Equal("cf-f21eea2dbdb8555f89fb"))
			})
		})

		Context("when given a installation.json without a er install", func() {
			var (
				jsonObj InstallationCompareObject
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open("fixtures/config.json")
				jsonObj, _ = ReadAndUnmarshal(fileRef)
			})

			AfterEach(func() {
			})

			It("Should return non nil error", func() {
				_, err := GetDeploymentName(jsonObj)
				Ω(err).ShouldNot(BeNil())
			})
		})

	})

	Describe("GetPasswordAndIP function", func() {
		Context("when given a valid installation.json", func() {
			var (
				jsonObj     InstallationCompareObject
				product     string = "cf"
				component   string = "ccdb"
				username    string = "admin"
				controlIp   string = "172.16.1.46"
				controlPass string = "e3e89a528625d819160d"
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)
			})

			AfterEach(func() {
			})

			It("Should return nil error, correct ip & password", func() {
				ip, password, err := GetPasswordAndIP(jsonObj, product, component, username)
				Ω(err).Should(BeNil())
				Ω(ip).Should(Equal(controlIp))
				Ω(password).Should(Equal(controlPass))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					GetPasswordAndIP(jsonObj, product, component, username)
				}).ShouldNot(Panic())
			})
		})
	})

	Describe("IpPasswordParser struct", func() {
		Context("when given a valid installation.json", func() {
			var (
				parser      *IpPasswordParser
				jsonObj     InstallationCompareObject
				product     string = "cf"
				component   string = "ccdb"
				username    string = "admin"
				controlIp   string = "172.16.1.46"
				controlPass string = "e3e89a528625d819160d"
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)

				parser = &IpPasswordParser{
					Product:   product,
					Component: component,
					Username:  username,
				}
			})

			AfterEach(func() {
			})

			It("Should return nil error, correct ip & password", func() {
				ip, password, err := parser.Parse(jsonObj)
				Ω(err).Should(BeNil())
				Ω(ip).Should(Equal(controlIp))
				Ω(password).Should(Equal(controlPass))
			})

			It("Should not panic on complete real world dataset", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})

		Context("when no valid component found", func() {
			var (
				parser    *IpPasswordParser
				jsonObj   InstallationCompareObject
				product   string = "cf"
				component string = "aaaa"
				username  string = "admin"
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)

				parser = &IpPasswordParser{
					Product:   product,
					Component: component,
					Username:  username,
				}
			})

			AfterEach(func() {
			})

			It("Should return error", func() {
				ip, password, err := parser.Parse(jsonObj)
				Ω(err).ShouldNot(BeNil())
				Ω(ip).Should(BeEmpty())
				Ω(password).Should(BeEmpty())
			})

			It("Should not panic", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})

		Context("when no valid product found", func() {
			var (
				parser    *IpPasswordParser
				jsonObj   InstallationCompareObject
				product   string = "fail"
				component string = "ccdb"
				username  string = "admin"
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)

				parser = &IpPasswordParser{
					Product:   product,
					Component: component,
					Username:  username,
				}
			})

			AfterEach(func() {
			})

			It("Should return error", func() {
				ip, password, err := parser.Parse(jsonObj)
				Ω(err).ShouldNot(BeNil())
				Ω(ip).Should(BeEmpty())
				Ω(password).Should(BeEmpty())
			})

			It("Should not panic", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})

		Context("when no valid username found", func() {
			var (
				parser    *IpPasswordParser
				jsonObj   InstallationCompareObject
				product   string = "cf"
				component string = "ccdb"
				username  string = "fail"
			)

			BeforeEach(func() {
				var fileRef *os.File
				defer fileRef.Close()
				fileRef, _ = os.Open(installationSettingsFilePath)
				jsonObj, _ = ReadAndUnmarshal(fileRef)

				parser = &IpPasswordParser{
					Product:   product,
					Component: component,
					Username:  username,
				}
			})

			AfterEach(func() {
			})

			It("Should return error", func() {
				ip, password, err := parser.Parse(jsonObj)
				Ω(err).ShouldNot(BeNil())
				Ω(ip).Should(BeEmpty())
				Ω(password).Should(BeEmpty())
			})

			It("Should not panic", func() {
				Ω(func() {
					parser.Parse(jsonObj)
				}).ShouldNot(Panic())
			})
		})
	})
}
