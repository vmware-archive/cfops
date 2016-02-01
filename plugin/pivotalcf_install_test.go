package plugin_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops/plugin"
)

var _ = Describe("DefaultPivotalCF initialized with valid installationSettings", func() {
	var configParser *cfbackup.ConfigurationParser
	var pivotalCF PivotalCF
	BeforeEach(func() {
		configParser = cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-default.json")
		pivotalCF = DefaultPivotalCF(configParser)
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
})
