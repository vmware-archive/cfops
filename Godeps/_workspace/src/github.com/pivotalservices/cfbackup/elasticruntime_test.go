package cfbackup_test

import (
	"fmt"
	"io"
	"io/ioutil"
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
					HttpGateway: &MockHttpGateway{},
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
					HttpGateway: &MockHttpGateway{true, 500, `{"state":"notdone"}`},
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
					HttpGateway: &MockHttpGateway{},
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
					HttpGateway: &MockHttpGateway{},
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
					HttpGateway: &MockHttpGateway{},
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
					HttpGateway: &MockHttpGateway{},
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
