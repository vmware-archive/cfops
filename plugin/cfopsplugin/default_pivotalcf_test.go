package cfopsplugin_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tileregistry"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
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
		pivotalCF = NewPivotalCF(configParser.InstallationSettings, controlTileSpec)
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

})
