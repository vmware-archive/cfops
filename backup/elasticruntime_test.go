package backup_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/backup/modules/persistence"
	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockDumper struct{}

func (s mockDumper) Dump(i io.Writer) (err error) {
	i.Write([]byte("sometext"))
	return
}

func mockDumperFunc(port int, database, username, password string, sshCfg command.SshConfig) (dpr persistence.Dumper, err error) {
	dpr = &mockDumper{}
	return
}

var _ = Describe("ElasticRuntime", func() {
	Describe("RunPostgresBackup function", func() {
		Context("with a valid product and component for ccdb", func() {
			var (
				product   string = "cf"
				component string = "ccdb"
				username  string = "admin"
				target    string
				er        ElasticRuntime
				info      DbBackupInfo = DbBackupInfo{
					Product:   product,
					Component: component,
					Username:  username,
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					BackupContext: BackupContext{
						TargetDir: target,
					},
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]DbBackupInfo{info})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]DbBackupInfo{info})
				}).ShouldNot(Panic())
				Ω(err).Should(BeNil())
			})
		})

		Context("with a valid product and component for consoledb", func() {
			var (
				product   string = "cf"
				component string = "consoledb"
				username  string = "root"
				target    string
				er        ElasticRuntime
				info      DbBackupInfo = DbBackupInfo{
					Product:   product,
					Component: component,
					Username:  username,
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					BackupContext: BackupContext{
						TargetDir: target,
					},
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]DbBackupInfo{info})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]DbBackupInfo{info})
				}).ShouldNot(Panic())
				Ω(err).Should(BeNil())
			})
		})

		Context("with a valid product and component for uaadb", func() {
			var (
				product   string = "cf"
				component string = "uaadb"
				username  string = "root"
				target    string
				er        ElasticRuntime
				info      DbBackupInfo = DbBackupInfo{
					Product:   product,
					Component: component,
					Username:  username,
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					BackupContext: BackupContext{
						TargetDir: target,
					},
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]DbBackupInfo{info})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]DbBackupInfo{info})
				}).ShouldNot(Panic())
				Ω(err).Should(BeNil())
			})
		})

		Context("with a invalid product, username and component", func() {
			var (
				product   string = "aaaaaaaa"
				component string = "aaaaaaaa"
				username  string = "aaaaaaaa"
				target    string
				er        ElasticRuntime
				info      DbBackupInfo = DbBackupInfo{
					Product:   product,
					Component: component,
					Username:  username,
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					BackupContext: BackupContext{
						TargetDir: target,
					},
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should not write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]DbBackupInfo{info})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).ShouldNot(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]DbBackupInfo{info})
				}).ShouldNot(Panic())
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
