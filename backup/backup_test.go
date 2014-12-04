package backup_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfops/backup"

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
		context = &BackupContext{
			Hostname:  "localhost",
			Username:  "admin",
			Password:  "admin",
			TPassword: "admin",
			Target:    path.Join(dir, "backup"),
		}
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("Something", func() {
		Context("With an empty directory", func() {
			It("should create the parent directory", func() {
				Ω(context.Target).NotTo(BeEquivalentTo(""))

			})
		})
	})
})
