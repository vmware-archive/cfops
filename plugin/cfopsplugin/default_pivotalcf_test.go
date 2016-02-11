package cfopsplugin_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/cfops/tileregistry"
)

var _ = Describe("DefaultPivotalCF initialized with valid installationSettings & TileSpec", func() {
	var configParser *cfbackup.ConfigurationParser
	var pivotalCF PivotalCF
	var controlTileSpec = tileregistry.TileSpec{
		OpsManagerHost: "localhost",
	}
	BeforeEach(func() {
		configParser = cfbackup.NewConfigurationParser("./fixtures/installation-settings-1-6-default.json")
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

	XContext("when NewArchiveWriter is called w/ a name", func() {
		controlName := "myarchive"
		It("then it should return a writer based on the cfops configured directory target", func() {
			writer := pivotalCF.NewArchiveWriter(controlName)
			Ω(writer).Should(BeAssignableToTypeOf(bytes.NewBufferString("")))
		})
	})

	XContext("when NewArchiveReader is called w/ a name", func() {
		controlName := "myarchive"
		It("then it should return a reader based on the cfops configured directory target", func() {
			reader := pivotalCF.NewArchiveReader(controlName)
			contents, _ := ioutil.ReadAll(reader)
			Ω(contents).Should(Equal(""))
		})
	})

})
