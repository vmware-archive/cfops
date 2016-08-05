package system

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
)

var cfopsExecutablePath string
var cfopsLinuxExecutablePath string
var logger lager.Logger

var _ = BeforeSuite(func() {
	cfConfig.APIEndpoint = os.Getenv("CF_API_URL")
	cfConfig.OMAdminUser = os.Getenv("OM_USER")
	cfConfig.OMAdminPassword = os.Getenv("OM_PASSWORD")
	cfConfig.OMHostname = os.Getenv("OM_HOSTNAME")
	cfConfig.AmiID = os.Getenv("OPSMAN_AMI")
	cfConfig.SecurityGroup = os.Getenv("AWS_SECURITY_GROUP")

	cfConfig.AppName = uuid.NewRandom().String()
	cfConfig.OrgName = uuid.NewRandom().String()
	cfConfig.SpaceName = uuid.NewRandom().String()
	cfConfig.AppPath = "assets/test-app"

	Expect(json.Unmarshal([]byte(os.Getenv("OM_PROXY_INFO")), &cfConfig.OMHostInfo)).To(Succeed())
	cfConfig.OMHostInfo.SSHKey = os.Getenv("OM_SSH_KEY")

	var err error

	logger = lager.NewLogger("Test Logs")
	logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

	os.Setenv("GOOS", "linux")
	cfopsLinuxExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	Expect(err).NotTo(HaveOccurred())
	os.Unsetenv("GOOS")

	cfopsExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("CFOps Ops Manager plugin", func() {
	It("backs up and restores successfully", func() {
		if os.Getenv("ONLY_ERT") == "true" {
			return
		}

		vm := createInstance("cfops-test", cfConfig.AmiID, cfConfig.SecurityGroup)
		defer vm.Destroy()

		ips, err := vm.GetIPs()
		newVMIP := ips[0].String()
		Expect(err).NotTo(HaveOccurred())

		backupCommand := exec.Command(
			cfopsExecutablePath,
			"backup",
			"--opsmanagerhost="+cfConfig.OMHostname,
			"--opsmanageruser=ubuntu",
			"--destination=../tmp/",
			"--adminuser="+cfConfig.OMAdminUser,
			"--adminpass="+cfConfig.OMAdminPassword,
			"--tile=ops-manager",
		)

		restoreCommand := exec.Command(
			cfopsExecutablePath,
			"restore",
			"--opsmanagerhost="+newVMIP,
			"--opsmanageruser=ubuntu",
			"--destination=../tmp/",
			"--adminuser="+cfConfig.OMAdminUser,
			"--adminpass="+cfConfig.OMAdminPassword,
			"--opsmanagerpassphrase="+cfConfig.OMAdminPassword,
			"--tile=ops-manager",
		)

		backupSession, err := gexec.Start(backupCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		if cfConfig.OMAdminPassword != "" {
			Consistently(backupSession.Out.Contents()).ShouldNot(ContainSubstring(cfConfig.OMAdminPassword))
			Consistently(backupSession.Err.Contents()).ShouldNot(ContainSubstring(cfConfig.OMAdminPassword))
		}
		Consistently(backupSession.Out.Contents()).ShouldNot(ContainSubstring("RSA PRIVATE KEY"))
		Consistently(backupSession.Err.Contents()).ShouldNot(ContainSubstring("RSA PRIVATE KEY"))

		Eventually(backupSession, 1200).Should(gexec.Exit(0))

		if os.Getenv("OM_VERSION") == "1.6" {
			createAdminUser(newVMIP, cfConfig.OMAdminUser, cfConfig.OMAdminPassword)
		}

		restoreSession, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(restoreSession, 1800).Should(gexec.Exit(0))

		time.Sleep(2 * time.Minute) // TODO make this better

		checkOpsManagersIdentical(cfConfig.OMHostname, newVMIP)
	})
})

func checkOpsManagersIdentical(oldHost, newHost string) {
	opsManager, err := NewOpsManagerClient(oldHost, cfConfig.OMAdminUser, cfConfig.OMAdminPassword, logger)
	Expect(err).NotTo(HaveOccurred())
	opsManagerProducts, _ := opsManager.GetStagedProducts()
	Expect(err).NotTo(HaveOccurred())

	restoredOpsManager, err := NewOpsManagerClient(newHost, cfConfig.OMAdminUser, cfConfig.OMAdminPassword, logger)
	Expect(err).NotTo(HaveOccurred())
	restoredOpsManagerProducts, _ := restoredOpsManager.GetStagedProducts()
	Expect(err).NotTo(HaveOccurred())

	Expect(opsManagerProducts).To(ConsistOf(restoredOpsManagerProducts))
}

var _ = Describe("CFOps Elastic Runtime plugin", func() {
	cfopsPath := "/tmp/cfops"
	backupPath := "/tmp/cfops-backup-" + uuid.NewRandom().String()

	BeforeEach(func() {
		opsManager, _ := NewOpsManagerClient(cfConfig.OMHostname, cfConfig.OMAdminUser, cfConfig.OMAdminPassword, logger)
		adminUser, adminPassword, err := opsManager.GetAdminCredentials()
		Expect(err).NotTo(HaveOccurred())
		cfConfig.AdminUser, cfConfig.AdminPassword = adminUser, adminPassword

		pushTestApp(cfConfig)
	})

	AfterEach(func() {
		deleteTestApp(cfConfig)

		if cfopsPath != "" {
			remoteExecute(cfConfig.OMHostInfo, "rm -rf "+cfopsPath, os.Stdout)
		}

		if backupPath != "" {
			remoteExecute(cfConfig.OMHostInfo, "rm -rf "+backupPath, os.Stdout)
		}
	})

	It("backs up and restores successfully", func() {
		backupCmd := strings.Join([]string{
			"LOG_LEVEL=debug",
			cfopsPath,
			"backup",
			"--opsmanagerhost=" + cfConfig.OMHostname,
			"--opsmanageruser=ubuntu",
			"--destination=" + backupPath,
			"--adminuser=" + cfConfig.OMAdminUser,
			"--adminpass=" + cfConfig.OMAdminPassword,
			"--tile=elastic-runtime",
		}, " ")

		restoreCmd := strings.Join([]string{
			"LOG_LEVEL=debug",
			cfopsPath,
			"restore",
			"--opsmanagerhost=" + cfConfig.OMHostname,
			"--opsmanageruser=ubuntu",
			"--destination=" + backupPath,
			"--adminuser=" + cfConfig.OMAdminUser,
			"--adminpass=" + cfConfig.OMAdminPassword,
			"--tile=elastic-runtime",
		}, " ")

		scpHelper(cfConfig.OMHostInfo, cfopsLinuxExecutablePath, cfopsPath)
		remoteExecute(cfConfig.OMHostInfo, "chmod +x /tmp/cfops", os.Stdout)

		fmt.Println("Backing up ERT...")
		remoteExecute(cfConfig.OMHostInfo, backupCmd, os.Stdout)

		deleteTestApp(cfConfig)

		fmt.Println("Restoring ERT...")
		remoteExecute(cfConfig.OMHostInfo, restoreCmd, os.Stdout)

		cfDo("target", "-o", cfConfig.OrgName, "-s", cfConfig.SpaceName)
		Eventually(cf.Cf("apps")).Should(gbytes.Say(cfConfig.AppName))
	})
})
