package backup

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup", func() {
	var (
		dir     string
		context *BackupContext
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
		context = New("localhost", "admin", "admin", "admin", path.Join(dir, "backup"))
	})

	AfterEach(func() {
		//os.RemoveAll(dir)
	})

	Describe("Prepare Filesystem", func() {
		Context("With an empty target", func() {
			It("should create the parent directory", func() {
				Ω(context.Target).NotTo(BeEquivalentTo(""))
				Ω(FileExists(context.Target)).To(BeFalse())
				context.initPaths()
				context.prepareFilesystem()
				Ω(FileExists(context.Target)).To(BeTrue())
			})

			It("should backup a tempest deployment's files", func() {
				copier := &testCopier{}
				err := context.backupDeployment(copier)
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
