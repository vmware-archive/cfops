package cfbackup_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfbackup"
	cfhttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/osutils"
)

var _ = Describe("OpsManager object", func() {
	var (
		opsManager *OpsManager
		tmpDir     string
		backupDir  string
	)
	Describe("Restore method", func() {

		Context("calling restore with failed removal of deployment files", func() {

			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &successGateway{}

				opsManager = &OpsManager{
					SettingsUploader: gw,
					AssetsUploader:   gw,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:         "localhost",
					Username:         "user",
					Password:         "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}
				f, _ := osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_SETTINGS_FILENAME)
				f.Close()
				f, _ = osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_ASSETS_FILENAME)
				f.Close()
			})

			It("Should yield error", func() {
				err := opsManager.Restore()
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("calling restore successfully", func() {

			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &successGateway{}

				opsManager = &OpsManager{
					SettingsUploader: gw,
					AssetsUploader:   gw,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:         "localhost",
					Username:         "user",
					Password:         "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &successExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}
				f, _ := osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_SETTINGS_FILENAME)
				f.Close()
				f, _ = osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_ASSETS_FILENAME)
				f.Close()
			})

			It("Should yield nil error", func() {
				err := opsManager.Restore()
				Ω(err).Should(BeNil())
			})
		})

		Context("calling restore unsuccessfully", func() {
			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &failGateway{}

				opsManager = &OpsManager{
					SettingsUploader: gw,
					AssetsUploader:   gw,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:         "localhost",
					Username:         "user",
					Password:         "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
				}
				f, _ := osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_SETTINGS_FILENAME)
				f.Close()
			})

			It("Should yield a non-nil error", func() {
				err := opsManager.Restore()
				Ω(err).ShouldNot(BeNil())
			})
		})

	})

	Describe("Backup method", func() {

		Context("called yielding an error in the chain", func() {
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

		Context("called yielding a successful rest call", func() {

			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &successGateway{}
				responseHandler := func(resp *http.Response) (interface{}, error) {
					defer resp.Body.Close()
					var settingsFileRef *os.File
					defer settingsFileRef.Close()
					settingsFileRef, _ = osutils.SafeCreate(path.Join(tmpDir, "backup"), "opsmanager", "installation.json");
					return io.Copy(settingsFileRef, resp.Body)
				}
				opsManager = &OpsManager{
					SettingsRequestor:    gw,
					Hostname: "localhost",
					Username: "user",
					Password: "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &successExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					HttpResponseHandler: responseHandler,
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

func (s *successGateway) ExecuteFunc(method string, handler cfhttp.HandleRespFunc) (val interface{}, err error) {
	res := &http.Response{
		StatusCode: 200,
	}
	res.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
	return handler(res)
}

type failGateway struct{}

func (s *failGateway) Upload(paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	res = &http.Response{
		StatusCode: 500,
	}
	res.Body = &ClosingBuffer{bytes.NewBufferString(successString)}
	return
}

func (s *failGateway) Execute(method string) (val interface{}, err error) {
	return
}

func (s *failGateway) ExecuteFunc(method string, handler cfhttp.HandleRespFunc) (val interface{}, err error) {
	return
}
