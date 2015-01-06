package backup_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	Describe("Backup", func() {
		var (
			restSuccessCalled int
			restFailureCalled int
			successString     string = `{"state":"done"}`
			failureString     string = `{"state":"notdone"}`
		)

		restSuccess := func(method, connectionURL, username, password string, isYaml bool) (resp *http.Response, err error) {
			resp = &http.Response{
				StatusCode: 200,
			}
			resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
			restSuccessCalled++
			return
		}

		restFailure := func(method, connectionURL, username, password string, isYaml bool) (resp *http.Response, err error) {
			resp = &http.Response{
				StatusCode: 500,
			}
			resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
			restFailureCalled++
			err = fmt.Errorf("")
			return
		}

		Context("with valid properties", func() {
			var (
				product   string = "cf"
				component string = "ccdb"
				username  string = "admin"
				target    string
				er        ElasticRuntime
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					RestRunner:      RestAdapter(restSuccess),
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should return nil error", func() {
				err := er.Backup()
				Ω(err).Should(BeNil())
			})
		})
		Context("with invalid properties", func() {
			var (
				product   string = "cf"
				component string = "ccdb"
				username  string = "admin"
				target    string
				er        ElasticRuntime
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					NewDumper:       mockDumperFunc,
					JsonFile:        "fixtures/installation.json",
					DeploymentsFile: "",
					DbEncryptionKey: "",
					RestRunner:      RestAdapter(restFailure),
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should not return nil error", func() {
				err := er.Backup()
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Describe("RunPostgresBackup function", func() {
		Context("with a valid product and component for ccdb", func() {
			var (
				product   string = "cf"
				component string = "ccdb"
				username  string = "admin"
				target    string
				er        ElasticRuntime
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
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
					SystemsInfo: info,
				}
				er.ReadAllUserCredentials()
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
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
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
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
					SystemsInfo: info,
				}
				er.ReadAllUserCredentials()
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
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
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
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
					SystemsInfo: info,
				}
				er.ReadAllUserCredentials()
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
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
				info      map[string]SystemInfo = map[string]SystemInfo{
					"ConsoledbInfo": SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
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
					SystemsInfo: info,
				}
				er.ReadAllUserCredentials()
			})

			AfterEach(func() {
				os.Remove(target)
			})

			It("Should not write the dumped output to a file in the databaseDir", func() {
				er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).ShouldNot(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemInfo{info["ConsoledbInfo"]})
				}).ShouldNot(Panic())
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
