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
	cfhttp "github.com/pivotalservices/gtils/http"
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

var (
	restSuccessCalled int
	restFailureCalled int
)

type mockHttpGateway struct {
	CheckFailureCondition bool
}

func (s *mockHttpGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (*http.Response, error) {
	return nil, nil
}

func (s *mockHttpGateway) Execute(method string) (interface{}, error) {
	if s.CheckFailureCondition {
		restFailureCalled++
		return &ClosingBuffer{bytes.NewBufferString(`{"state":"notdone"}`)}, nil
	}
	restSuccessCalled++
	return bytes.NewBufferString(`[{
		"agent_id": "d4131496-4cdf-4309-907b-e2ce327be029",
		"cid": "vm-8dfe3b38-6e31-4d9a-aeef-74cbf2143bd8",
		"job": "cloud_controller-partition-7bc61fd2fa9d654696df",
 		"index": 0
	}]`), nil
}

func (s *mockHttpGateway) ExecuteFunc(method string, handler cfhttp.HandleRespFunc) (interface{}, error) {
	resp := &http.Response{
		StatusCode: 200,
	}
	return handler(resp)
}

var _ = Describe("ElasticRuntime", func() {
	Describe("Backup", func() {
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
					"ConsoledbInfo": &PgInfoMock{
						SystemInfo: SystemInfo{
							Product:   product,
							Component: component,
							Identity:  username,
						},
					},
				}
				ps []SystemDump = []SystemDump{info["ConsoledbInfo"]}
			)

			BeforeEach(func() {
				target, _ = ioutil.TempDir("/tmp", "spec")
				er = ElasticRuntime{
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo:       info,
					PersistentSystems: ps,
					Logger:            Logger(),
				}
			})

			AfterEach(func() {
				os.Remove(target)
			})

			// It("Should return nil error", func() {
			// 	err := er.Backup()
			// 	Ω(err).Should(BeNil())
			// })
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
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{true},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
					Logger:      Logger(),
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
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
					Logger:      Logger(),
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
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
					Logger:      Logger(),
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
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
					Logger:      Logger(),
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
					JsonFile:    "fixtures/installation.json",
					HttpGateway: &mockHttpGateway{},
					BackupContext: BackupContext{
						TargetDir: target,
					},
					SystemsInfo: info,
					Logger:      Logger(),
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
