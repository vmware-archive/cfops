package backup

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

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
		//os.RemoveAll(dir)
	})

	Describe("Prepare Filesystem", func() {
		var (
			context *OpsManager
		)

		BeforeEach(func() {
			context = NewOpsManager("localhost", "admin", "admin", "admin", path.Join(dir, "backup"))
		})

		AfterEach(func() {

		})

		Context("With an empty target", func() {
			It("should backup a tempest deployment's files", func() {
				copier := &testCopier{}
				Ω(context.TargetDir).NotTo(BeEquivalentTo(""))
				e, _ := osutils.Exists(context.TargetDir)
				Ω(e).To(BeFalse())
				err := context.copyDeployments(copier)
				Ω(err).Should(BeNil())
				e, _ = osutils.Exists(context.TargetDir)
				Ω(e).To(BeTrue())
				Ω(copier.src).ToNot(BeNil())
				fmt.Println(copier.src)
				s, _ := ioutil.ReadAll(copier.src)
				Ω(string(s[:])).To(BeEquivalentTo("cd /var/tempest/workspaces/default && tar cz deployments"))
			})
		})
	})
})

type testCopier struct {
	src  io.Reader
	dest io.Writer
}

func (copier *testCopier) Copy(dest io.Writer, src io.Reader) error {
	s, _ := ioutil.ReadAll(src)
	srcvalue := string(s[:])
	copier.dest = dest
	copier.src = strings.NewReader(srcvalue)
	_, err := io.Copy(dest, strings.NewReader(srcvalue))
	return err
}
