package cfopsplugin_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)



var _ = Describe("DefaultPivotalCF initialized with valid installationSettings & TileSpec", func() {
	var configParser *cfbackup.ConfigurationParser
	var pivotalCF PivotalCF
	var controlTmpDir, _ = ioutil.TempDir("", "unit-test")
	var controlTileSpec = tileregistry.TileSpec{
		ArchiveDirectory: controlTmpDir,
		OpsManagerHost:   "localhost",
	}
	BeforeEach(func() {
		configParser = cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-aws.json")
		pivotalCF = NewPivotalCF(configParser, controlTileSpec)
	})
	Context("when GetCredentials is called", func() {
		It("then it should return a list of my systems credentials", func() {
			Ω(len(pivotalCF.GetCredentials()["p-bosh"]["director"])).Should(BeNumerically(">", 0))
			Ω(len(pivotalCF.GetCredentials()["cf"])).Should(BeNumerically(">", 0))
		})
	})

	Context("when GetProducts is called", func() {
		It("then it should return a list of my systems products", func() {
			Ω(len(pivotalCF.GetProducts()["p-bosh"].Jobs)).Should(BeNumerically(">", 0))
			Ω(len(pivotalCF.GetProducts()["cf"].Jobs)).Should(BeNumerically(">", 0))
		})
	})

	Context("when GetHostDetails is called", func() {
		It("then it should return its targeted host information", func() {
			Ω(pivotalCF.GetHostDetails()).Should(Equal(controlTileSpec))
		})
	})

	Context("when NewArchiveWriter is called w/ a archive name", func() {
		var controlName string
		BeforeEach(func() {
			u, _ := uuid.NewV4()
			controlName = u.String()
		})
		AfterEach(func() {
			os.Remove(path.Join(controlTmpDir, controlName))
		})
		It("then it should create a writer based on the cfops configured target and the archive name", func() {
			_, errStatBefore := os.Stat(path.Join(controlTmpDir, controlName))
			Ω(os.IsNotExist(errStatBefore)).ShouldNot(BeFalse())
			writer, err := pivotalCF.NewArchiveWriter(controlName)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(writer).ShouldNot(BeNil())
			Ω(func() {
				func(x io.Writer) {
					lo.G.Debug("the returned var should be a valid writer", writer)
				}(writer)
			}).ShouldNot(Panic())
			_, errStatAfter := os.Stat(path.Join(controlTmpDir, controlName))
			Ω(os.IsNotExist(errStatAfter)).Should(BeFalse())
		})
	})

	Context("when NewArchiveReader is called w/ a valid archive name", func() {
		var controlName string
		BeforeEach(func() {
			u, _ := uuid.NewV4()
			controlName := u.String()
			os.Create(path.Join(controlTmpDir, controlName))
		})
		AfterEach(func() {
			os.Remove(path.Join(controlTmpDir, controlName))
		})

		It("then it should return a reader based on the cfops configured target and the archive name", func() {
			reader, err := pivotalCF.NewArchiveReader(controlName)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(reader).ShouldNot(BeNil())
			Ω(func() {
				func(x io.Reader) {
					lo.G.Debug("the returned var should be a valid Reader", reader)
				}(reader)
			}).ShouldNot(Panic())
		})
	})
	products := map[string]string{
		"p-mysql": "mysql",
		"p-redis": "cf-redis-broker",
	}
	for productName, jobName := range products {
		controlProductName := productName
		controlJobName := jobName
		Context("when GetJobProperties is called", func() {

			Context(fmt.Sprintf("when called with %s product name and %s job name", controlProductName, controlJobName), func() {
				It("then it should return a list of properties", func() {
					properties, err := pivotalCF.GetJobProperties(controlProductName, controlJobName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(len(properties)).Should(BeNumerically(">", 0))
				})
			})
			Context("when called with invalid productName", func() {
				It("then it should return an error", func() {
					_, err := pivotalCF.GetJobProperties("missingProgram", "missingJobName")
					Ω(err).Should(HaveOccurred())
					Ω(err.Error()).Should(Equal("job missingJobName not found for product missingProgram"))
				})
			})
			Context("when called with invalid jobName", func() {
				It("then it should return an error", func() {
					_, err := pivotalCF.GetJobProperties("p-mysql", "missingJobName")
					Ω(err).Should(HaveOccurred())
					Ω(err.Error()).Should(Equal("job missingJobName not found for product p-mysql"))

				})
			})
		})
		Context("when GetPropertyValues is called", func() {
			Context(fmt.Sprintf("when called with %s product name and %s job name", controlProductName, controlJobName), func() {
				It("then it should map of properties", func() {
					pMap, err := pivotalCF.GetPropertyValues(controlProductName, controlJobName, "vm_credentials")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(pMap["identity"]).ShouldNot(BeEmpty())
					Ω(pMap["password"]).ShouldNot(BeEmpty())
				})
			})
		})
	}

})
