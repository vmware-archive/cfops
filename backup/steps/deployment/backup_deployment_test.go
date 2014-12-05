package deployment_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfops/backup/steps/deployment"
)

var _ = Describe("Backup", func() {
	var (
		dir           string
		deploymentDir string
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
		deploymentDir = path.Join(dir, "backup", "deployments")
	})

	AfterEach(func() {
		//os.RemoveAll(dir)
	})

	Describe("Prepare Filesystem", func() {
		Context("With an empty target", func() {
			It("should backup a tempest deployment's files", func() {
				copier := &testCopier{}
				d := deployment.New(deploymentDir)
				err := d.Backup(copier)
				Ω(err).Should(BeNil())
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
	io.Copy(dest, src)
	copier.dest = dest
	copier.src = src
	return nil
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
