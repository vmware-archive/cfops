package cfopsplugin_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"

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

	Context("when NewArchiveWriter is called w/ a name", func() {
		controlName := "myarchive"
		It("then it should create a writer based on the cfops configured directory target", func() {
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

	XContext("when NewArchiveReader is called w/ a name", func() {
		controlName := "myarchive"
		It("then it should return a reader based on the cfops configured directory target", func() {
			reader, err := pivotalCF.NewArchiveReader(controlName)
			contents, _ := ioutil.ReadAll(reader)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(contents).Should(Equal(""))
		})
	})

})
