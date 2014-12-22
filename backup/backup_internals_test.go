package backup

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfops/osutils"
)

var _ = Describe("Backup", func() {
	var (
		dir string
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("Prepare Filesystem", func() {
		var (
			context *OpsManager
		)

		BeforeEach(func() {
			context = &OpsManager{
				Hostname: "localhost",
				Username: "user",
				Password: "password",
				BackupContext: BackupContext{
					TargetDir: path.Join(dir, "backup"),
				},
				RestRunner: RestAdapter(invoke),
				Executer:   &testExecuter{},
			}
		})

		AfterEach(func() {

		})

		Context("With an empty target", func() {
			It("should backup a tempest deployment's files", func() {
				Ω(context.TargetDir).NotTo(BeEquivalentTo(""))
				e, _ := osutils.Exists(context.TargetDir)
				Ω(e).To(BeFalse())
				err := context.copyDeployments()
				Ω(err).Should(BeNil())
				e, _ = osutils.Exists(context.TargetDir)
				Ω(e).To(BeTrue())
			})
		})
	})
})

type testExecuter struct{}

func (s *testExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	return
}
