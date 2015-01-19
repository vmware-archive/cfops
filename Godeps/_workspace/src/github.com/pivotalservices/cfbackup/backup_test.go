package cfbackup_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfbackup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup", func() {
	var (
		dir     string
		context *OpsManager
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
		context = &OpsManager{
			Hostname:        "localhost",
			Username:        "admin",
			Password:        "admin",
			TempestPassword: "admin",
			BackupContext: BackupContext{
				TargetDir: path.Join(dir, "backup"),
			},
		}
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("Something", func() {
		Context("With an empty directory", func() {
			It("should create the parent directory", func() {
				Ω(context.TargetDir).NotTo(BeEquivalentTo(""))
			})
		})
	})
})
