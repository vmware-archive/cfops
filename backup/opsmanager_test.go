package backup_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/osutils"
)

var _ = Describe("OpsManager object", func() {
	var (
		opsManager *OpsManager
		tmpDir     string
		backupDir  string
	)
	Describe("Restore method", func() {
		var (
			gateway httpGateway
		)
		Context("calling restore successfully", func() {

			BeforeEach(func() {
				gateway = &successGateway{}
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &successGateway{}

				opsManager = &OpsManager{
					SettingsUploader: gw,
					AssetsUploader:   gw,
					Hostname:         "localhost",
					Username:         "user",
					Password:         "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					RestRunner:          RestAdapter(restFailure),
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}
				f, _ := osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_SETTINGS_FILENAME)
				f.Close()
			})

			It("Should yield nil error", func() {
				fmt.Println(opsManager)
				err := opsManager.Restore()
				Ω(err).Should(BeNil())
			})
		})
	})

	Describe("Backup method", func() {

		Context("called yeilding a error in the chain", func() {
			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")

				opsManager = &OpsManager{
					Hostname: "localhost",
					Username: "user",
					Password: "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					RestRunner:          RestAdapter(restFailure),
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}

			})

			It("should return non nil error and not write installation.json", func() {
				err := opsManager.Backup()
				filepath := path.Join(backupDir, "installation.json")
				Ω(err).ShouldNot(BeNil())
				Ω(osutils.Exists(filepath)).Should(BeFalse())
			})

			It("should return non nil error and not write cc_db_encryption_key.txt", func() {
				err := opsManager.Backup()
				filepath := path.Join(backupDir, "cc_db_encryption_key.txt")
				Ω(err).ShouldNot(BeNil())
				Ω(osutils.Exists(filepath)).Should(BeFalse())
			})

			It("should return non nil error and not write deployments.tar.gz", func() {
				err := opsManager.Backup()
				filepath := path.Join(backupDir, "deployments.tar.gz")
				Ω(err).ShouldNot(BeNil())
				Ω(osutils.Exists(filepath)).Should(BeTrue())
			})
		})

		Context("called yeilding a successful rest call", func() {

			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")

				opsManager = &OpsManager{
					Hostname: "localhost",
					Username: "user",
					Password: "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					RestRunner:          RestAdapter(restSuccess),
					Executer:            &successExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}

			})

			It("should return nil error and write the proper information to the installation.json", func() {
				err := opsManager.Backup()
				filepath := path.Join(backupDir, "installation.json")
				b, _ := ioutil.ReadFile(filepath)
				Ω(err).Should(BeNil())
				Ω(b).Should(Equal([]byte(successString)))
			})

			It("should return nil error and write ", func() {
				opsManager.Backup()
				filepath := path.Join(backupDir, "cc_db_encryption_key.txt")
				Ω(osutils.Exists(filepath)).Should(BeTrue())
			})

			It("should return nil error and write ", func() {
				opsManager.Backup()
				filepath := path.Join(backupDir, "deployments.tar.gz")
				Ω(osutils.Exists(filepath)).Should(BeTrue())
			})
		})
	})
})

var (
	successString string = "successString"
	failureString string = "failureString"
)

func restSuccess(method, connectionURL, username, password string, isYaml bool) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: 200,
	}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
	return
}

func restFailure(method, connectionURL, username, password string, isYaml bool) (resp *http.Response, err error) {
	resp = &http.Response{
		StatusCode: 500,
	}
	resp.Body = &ClosingBuffer{bytes.NewBufferString(failureString)}
	return
}

type successExecuter struct{}

func (s *successExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	return
}

type failExecuter struct{}

func (s *failExecuter) Execute(dest io.Writer, src string) (err error) {
	dest.Write([]byte(src))
	err = fmt.Errorf("error failure")
	return
}

type successGateway struct{}

func (s *successGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	res = &http.Response{
		StatusCode: 200,
	}
	res.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
	return
}

func (s *successGateway) Execute(method string) (val interface{}, err error) {
	return
}

type httpUploader interface {
	Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error)
}

type httpExecuter interface {
	Execute(method string) (val interface{}, err error)
}

type httpGateway interface {
	httpUploader
	httpExecuter
}
