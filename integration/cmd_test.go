package cfopsintegration_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var cfopsExecutablePath string
var err error

var executionTimeout = 5 * time.Second

const (
	help string = `cfops - Cloud Foundry Operations Tool`
)

var _ = BeforeSuite(func() {
	os.Setenv("LOG_LEVEL", "debug")
	cfopsExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	Expect(err).NotTo(HaveOccurred())

	Expect(directoryExists("/var/vcap/store")).To(BeTrue(), "need the /var/vcap/store directory to run tests")
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("cfops cmd", func() {
	BeforeEach(func() {
		os.RemoveAll("/var/vcap/store/shared")
	})

	It("prints the help page", func() {
		command := exec.Command(cfopsExecutablePath)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred(), "couldn't run cfops executable")
		Eventually(session.Out).Should(gbytes.Say(help))
	})
	var currentUser *user.User
	var privateKey string
	BeforeEach(func() {
		currentUser, err = user.Current()
		Expect(err).NotTo(HaveOccurred())
		var publicKey string
		publicKey, privateKey = createSSHKey(currentUser.Name)
		addToAuthorizedKeys(currentUser, publicKey)
	})

	AfterEach(func() {
		removeKeyFromAuthorized(currentUser)
	})

	Describe("backup blobstore", func() {
		var destinationDirectory string
		var opsmanUri *url.URL
		var additionalFlag string
		var boshDirectorServer *ghttp.Server
		var opsmanServer *ghttp.Server
		var session *gexec.Session
		BeforeEach(func() {
			boshDirectorServer = ghttp.NewUnstartedServer()
			boshDirectorServer.HTTPTestServer.Listener, err = net.Listen("tcp", "127.0.0.1:25555")
			Expect(err).NotTo(HaveOccurred())
			boshDirectorServer.HTTPTestServer.StartTLS()

			opsmanServer = ghttp.NewTLSServer()
			opsmanServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/uaa/oauth/token"),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/installation_settings"),
					ghttp.RespondWith(http.StatusOK, readFixture("fixtures/nfs_blobstore_test_installation_settings.json",
						struct {
							DirectorHost   string
							NfsServerIP    string
							NFSSshUser     string
							NFSSshPassword string
							SSHPrivateKey  string
						}{DirectorHost: "127.0.0.1", NfsServerIP: "127.0.0.1", NFSSshUser: currentUser.Name, NFSSshPassword: "SECRET_nfs_ssh_password", SSHPrivateKey: strings.Replace(privateKey, "\n", "\\n", -1)})),
				))
			opsmanUri, err = url.Parse(opsmanServer.URL())
			Expect(err).NotTo(HaveOccurred())

			boshDirectorServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/info"),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/deployments/cf-f21eea2dbdb8555f89fb"),
					ghttp.RespondWith(http.StatusOK, `---`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/deployments/cf-f21eea2dbdb8555f89fb/vms"),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
			)

			createTestFiles("/var/vcap/store", []string{
				"shared/cc-resources/09/4d/094dc37299e8d0c68e8e22e8f72a7b1632d26cc3",
				"shared/cc-buildpacks/07/7a/077a250e-fdc4-40e5-8d0b-2b8e9cbd1aba_886bd2888127429f7f75120d98a187cbf7289c16",
				"shared/cc-droplets/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a/01f1938d589f3e2aa6e3dfa9d4308e215061a5b6",
				"shared/cc-packages/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a",
			})

			destinationDirectory, err = ioutil.TempDir("", "cfops_destination")
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll(destinationDirectory)
			opsmanServer.Close()
			boshDirectorServer.Close()
		})
		JustBeforeEach(func() {
			command := exec.Command(cfopsExecutablePath, "backup", "--opsmanagerhost", opsmanUri.Host, "--adminuser", "<usr>", "--adminpass", "SECRET_admin_password", "--opsmanageruser", "<opsuser>", "-d", destinationDirectory, "--tile", "elastic-runtime", additionalFlag)

			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			session.Wait(executionTimeout)
		})
		Context("dosen't log secrets", func() {
			It("should Succeed", func() {

				Consistently(session.Err.Contents()).ShouldNot(ContainSubstring("SECRET"))
				Consistently(session.Err.Contents()).ShouldNot(ContainSubstring("BEGIN RSA PRIVATE KEY"))
				Consistently(session.Out.Contents()).ShouldNot(ContainSubstring("SECRET"))
				Consistently(session.Out.Contents()).ShouldNot(ContainSubstring("BEGIN RSA PRIVATE KEY"))
			})
		})

		Context("without any additional flags", func() {
			BeforeEach(func() {
				additionalFlag = ""
			})
			It("should Succeed", func() {
				Eventually(session).Should(gexec.Exit(0))
			})
			It("backups all the files", func() {
				nfsBackupPath := filepath.Join(destinationDirectory, "nfs_server.backup")
				Expect(filesInTar(nfsBackupPath)).To(ConsistOf("shared/cc-resources/09/4d/094dc37299e8d0c68e8e22e8f72a7b1632d26cc3",
					"shared/cc-buildpacks/07/7a/077a250e-fdc4-40e5-8d0b-2b8e9cbd1aba_886bd2888127429f7f75120d98a187cbf7289c16",
					"shared/cc-droplets/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a/01f1938d589f3e2aa6e3dfa9d4308e215061a5b6",
					"shared/cc-packages/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a"))
			})
		})

		Context("with nfs flag set to full", func() {
			BeforeEach(func() {
				additionalFlag = "-nfs=full"
			})
			It("should succeed", func() {
				Eventually(session).Should(gexec.Exit(0))
			})
			It("backups all the files", func() {
				nfsBackupPath := filepath.Join(destinationDirectory, "nfs_server.backup")
				Expect(filesInTar(nfsBackupPath)).To(ConsistOf("shared/cc-resources/09/4d/094dc37299e8d0c68e8e22e8f72a7b1632d26cc3",
					"shared/cc-buildpacks/07/7a/077a250e-fdc4-40e5-8d0b-2b8e9cbd1aba_886bd2888127429f7f75120d98a187cbf7289c16",
					"shared/cc-droplets/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a/01f1938d589f3e2aa6e3dfa9d4308e215061a5b6",
					"shared/cc-packages/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a"))
			})
		})

		Context("with nfs flag set to skip", func() {
			BeforeEach(func() {
				additionalFlag = "-nfs=skip"
			})
			It("should succeed", func() {
				Eventually(session).Should(gexec.Exit(0))
			})
			It("does not back up NFS", func() {
				nfsBackupPath := filepath.Join(destinationDirectory, "nfs_server.backup")
				Expect(nfsBackupPath).NotTo(BeAnExistingFile())
			})
		})

		Context("with nfs flag set to a invalid valid", func() {
			BeforeEach(func() {
				additionalFlag = "-nfs=invalid"
			})
			It("does not run the backup", func() {
				Expect(session.Err).Should(gbytes.Say("invalid cli flag args"))
			})
			It("should fail", func() {
				Eventually(session).ShouldNot(gexec.Exit(0))
			})
			It("does not back up NFS", func() {
				nfsBackupPath := filepath.Join(destinationDirectory, "nfs_server.backup")
				Expect(nfsBackupPath).NotTo(BeAnExistingFile())
			})
		})

		Context("with nfs flag set to lite", func() {
			BeforeEach(func() {
				additionalFlag = "-nfs=lite"
			})

			It("should succeed", func() {
				Eventually(session).Should(gexec.Exit(0))
			})

			It("backs up NFS without cc-resources dir", func() {
				nfsBackupPath := filepath.Join(destinationDirectory, "nfs_server.backup")
				Expect(filesInTar(nfsBackupPath)).To(ConsistOf(
					"shared/cc-buildpacks/07/7a/077a250e-fdc4-40e5-8d0b-2b8e9cbd1aba_886bd2888127429f7f75120d98a187cbf7289c16",
					"shared/cc-droplets/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a/01f1938d589f3e2aa6e3dfa9d4308e215061a5b6",
					"shared/cc-packages/63/94/63945b8b-b3f4-4736-bb46-edb5a5dae80a"))
			})
		})
	})
})

func filesInTar(filename string) []string {
	nfsBackupFile, err := os.Open(filename)
	defer nfsBackupFile.Close()
	Expect(err).ShouldNot(HaveOccurred())

	var tarReader *tar.Reader
	gzf, err := gzip.NewReader(nfsBackupFile)
	Expect(err).ShouldNot(HaveOccurred())
	if err == nil {
		tarReader = tar.NewReader(gzf)
	} else {
		tarReader = tar.NewReader(nfsBackupFile)
	}

	files := []string{}
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}
		if hdr.Typeflag == tar.TypeReg {
			files = append(files, hdr.Name)
		}
	}
	return files
}

func readFixture(filename string, variables interface{}) string {
	contents, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())

	t := template.Must(template.New("fixture").Parse(string(contents)))

	buffer := bytes.NewBuffer([]byte{})
	Expect(t.Execute(buffer, variables)).NotTo(HaveOccurred())
	return buffer.String()
}

func createSSHKey(sshKeyUsername string) (string, string) {
	sshKeys, err := ioutil.TempDir("", "integration-ssh-keys")
	Expect(err).ToNot(HaveOccurred())
	privateKeyPath := filepath.Join(sshKeys, "id_rsa")
	Expect(exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", sshKeyUsername,
		"-N", "", "-f", privateKeyPath).Run()).To(Succeed())
	privateKeyContents, err := ioutil.ReadFile(privateKeyPath)
	Expect(err).ToNot(HaveOccurred())
	publicKeyContents, err := ioutil.ReadFile(filepath.Join(sshKeys, "id_rsa.pub"))
	Expect(err).ToNot(HaveOccurred())
	os.RemoveAll(sshKeys)
	return string(publicKeyContents), string(privateKeyContents)
}

func addToAuthorizedKeys(unixUser *user.User, pubKey string) {
	Expect(os.MkdirAll(filepath.Join(unixUser.HomeDir, ".ssh"), 0700)).To(Succeed())
	authKeys, err := os.OpenFile(filepath.Join(unixUser.HomeDir, ".ssh", "authorized_keys"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	Expect(err).ToNot(HaveOccurred())
	authKeys.WriteString(pubKey)
	authKeys.Close()
}

func removeKeyFromAuthorized(unixUser *user.User) {
	authKeysFilePath := filepath.Join(unixUser.HomeDir, ".ssh", "authorized_keys")
	authKeysContent, err := ioutil.ReadFile(authKeysFilePath)
	Expect(err).NotTo(HaveOccurred())

	trimmedAuthKeysLines := [][]byte{}
	for _, line := range bytes.Split(authKeysContent, []byte("\n")) {
		if !bytes.Contains(line, []byte(unixUser.Name)) {
			trimmedAuthKeysLines = append(trimmedAuthKeysLines, line)
		}
	}

	trimmedAuthKeysContent := bytes.Join(trimmedAuthKeysLines, []byte("\n"))
	err = ioutil.WriteFile(authKeysFilePath, trimmedAuthKeysContent, 0600)
	Expect(err).NotTo(HaveOccurred())
}

func toJson(s interface{}) string {
	b, err := json.Marshal(s)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}

func directoryExists(dirname string) bool {
	_, err := os.Stat(dirname)
	return err == nil
}

func createTestFiles(dirname string, files []string) {

	for _, file := range files {
		fullFileName := filepath.Join(dirname, file)
		subDirname := filepath.Dir(fullFileName)
		Expect(os.MkdirAll(subDirname, 0777)).NotTo(HaveOccurred())

		w, err := os.Create(fullFileName)
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}
		Expect(w.Close()).NotTo(HaveOccurred())
	}
}
