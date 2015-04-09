package cfbackup_test

import (
	"io/ioutil"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfbackup"
	. "github.com/pivotalservices/gtils/command"
	ghttp "github.com/pivotalservices/gtils/http"
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
				gw := &MockHttpGateway{}

				opsManager = &OpsManager{
					SettingsUploader:  MockMultiPartUploadFunc,
					AssetsUploader:    MockMultiPartUploadFunc,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:          "localhost",
					Username:          "user",
					Password:          "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
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
				gw := &MockHttpGateway{}

				opsManager = &OpsManager{
					SettingsUploader:  MockMultiPartUploadFunc,
					AssetsUploader:    MockMultiPartUploadFunc,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:          "localhost",
					Username:          "user",
					Password:          "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &successExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
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

		Context("calling restore when unable to upload", func() {
			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &MockHttpGateway{StatusCode: 500, State: failureString}

				opsManager = &OpsManager{
					SettingsUploader:  ghttp.MultiPartUpload,
					AssetsUploader:    ghttp.MultiPartUpload,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:          "localhost",
					Username:          "user",
					Password:          "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
				}
				f, _ := osutils.SafeCreate(opsManager.TargetDir, opsManager.OpsmanagerBackupDir, OPSMGR_INSTALLATION_SETTINGS_FILENAME)
				f.Close()
			})

			It("Should yield a non-nil error", func() {
				err := opsManager.Restore()
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("calling restore unsuccessfully", func() {
			BeforeEach(func() {
				tmpDir, _ = ioutil.TempDir("/tmp", "test")
				backupDir = path.Join(tmpDir, "backup", "opsmanager")
				gw := &MockHttpGateway{StatusCode: 500, State: failureString}

				opsManager = &OpsManager{
					SettingsUploader:  MockMultiPartUploadFunc,
					AssetsUploader:    MockMultiPartUploadFunc,
					SettingsRequestor: gw,
					AssetsRequestor:   gw,
					Hostname:          "localhost",
					Username:          "user",
					Password:          "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &failExecuter{},
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
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
					LocalExecuter:       NewLocalMockExecuter(),
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
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
				gw := &MockHttpGateway{StatusCode: 200, State: successString}
				opsManager = &OpsManager{
					SettingsRequestor: gw,
					Hostname:          "localhost",
					Username:          "user",
					Password:          "password",
					BackupContext: BackupContext{
						TargetDir: path.Join(tmpDir, "backup"),
					},
					Executer:            &successExecuter{},
					LocalExecuter:       NewLocalMockExecuter(),
					DeploymentDir:       "fixtures/encryptionkey",
					OpsmanagerBackupDir: "opsmanager",
					Logger:              Logger(),
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
