package cfbackup_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	. "github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/gtils/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PgInfoMock struct {
	SystemInfo
}

func (s *PgInfoMock) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	dumper = &mockDumper{}
	return
}

type mockDumper struct{}

func (s mockDumper) Dump(i io.Writer) (err error) {
	i.Write([]byte("sometext"))
	return
}

func (s mockDumper) Import(i io.Reader) (err error) {
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

		Context("with valid properties (DirectorInfo)", func() {
			var (
				product   string = "microbosh"
				component string = "director"
				username  string = "director"
				target    string
				er        ElasticRuntime
				info      map[string]SystemDump = map[string]SystemDump{
					"DirectorInfo": &SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile:   "fixtures/installation.json",
					RestRunner: RestAdapter(restSuccess),
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
				component string = "consoledb"
				username  string = "root"
				target    string
				er        ElasticRuntime
				info      map[string]SystemDump = map[string]SystemDump{
					"ConsoledbInfo": &SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile:   "fixtures/installation.json",
					RestRunner: RestAdapter(restFailure),
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

			It("Should not panic", func() {
				var err error
				Ω(func() {
					err = er.Backup()
				}).ShouldNot(Panic())
			})
		})
	})

	Describe("RunPostgresBackup function", func() {
		Context("with a valid product and component for ccdb", func() {
			var (
				product   string = "cf"
				component string = "consoledb"
				username  string = "root"
				target    string
				er        ElasticRuntime
				info      map[string]SystemDump = map[string]SystemDump{
					"ConsoledbInfo": &PgInfoMock{
						SystemInfo: SystemInfo{
							Product:   product,
							Component: component,
							Identity:  username,
						},
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile: "fixtures/installation.json",
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
				er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.backup", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
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
				info      map[string]SystemDump = map[string]SystemDump{
					"ConsoledbInfo": &PgInfoMock{
						SystemInfo: SystemInfo{
							Product:   product,
							Component: component,
							Identity:  username,
						},
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile: "fixtures/installation.json",
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
				er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.backup", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
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
				info      map[string]SystemDump = map[string]SystemDump{
					"UaadbInfo": &PgInfoMock{
						SystemInfo: SystemInfo{
							Product:   product,
							Component: component,
							Identity:  username,
						},
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile: "fixtures/installation.json",
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
				er.RunDbBackups([]SystemDump{info["UaadbInfo"]})
				filename := fmt.Sprintf("%s.backup", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).Should(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemDump{info["UaadbInfo"]})
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
				info      map[string]SystemDump = map[string]SystemDump{
					"ConsoledbInfo": &SystemInfo{
						Product:   product,
						Component: component,
						Identity:  username,
					},
				}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile: "fixtures/installation.json",
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
				er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
				filename := fmt.Sprintf("%s.sql", component)
				exists, _ := osutils.Exists(path.Join(target, filename))
				Ω(exists).ShouldNot(BeTrue())
			})

			It("Should have a nil error and not panic", func() {
				var err error
				Ω(func() {
					err = er.RunDbBackups([]SystemDump{info["ConsoledbInfo"]})
				}).ShouldNot(Panic())
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
