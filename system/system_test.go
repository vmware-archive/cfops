package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
)

var cfopsExecutablePath string
var cfopsLinuxExecutablePath string
var err error

var cfConfig Config

var _ = BeforeSuite(func() {
	os.Setenv("GOOS", "linux")
	cfopsLinuxExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	os.Unsetenv("GOOS")

	cfopsExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	Expect(err).NotTo(HaveOccurred())

	cfConfig.ApiEndpoint = CfAPI
	cfConfig.AdminUser = CfUser
	cfConfig.AdminPassword = CfPassword
	cfConfig.OMAdminUser = OMAdminUser
	cfConfig.OMAdminPassword = OMAdminPassword
	cfConfig.OMHostname = OMHostname
	cfConfig.AppName = uuid.NewRandom().String()
	cfConfig.OrgName = uuid.NewRandom().String()
	cfConfig.SpaceName = uuid.NewRandom().String()
	cfConfig.AppPath = "assets/test-app"
	cfConfig.OMSSHKey = OMSSHKey
	cfConfig.AmiID = AmiID
	cfConfig.SecurityGroup = SecurityGroup
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("CFOps Ops Manager plugin", func() {
	It("backs up and restores successfully", func() {
		vm := createInstance("cfops-test", cfConfig.AmiID, cfConfig.SecurityGroup)
		defer vm.Destroy()

		ips, err := vm.GetIPs()
		newVMIP := ips[0].String()
		Expect(err).NotTo(HaveOccurred())

		backupCmd := exec.Command(
			cfopsExecutablePath,
			"backup",
			"--opsmanagerhost="+cfConfig.OMHostname,
			"--opsmanageruser=ubuntu",
			"--destination=../tmp/",
			"--adminuser="+cfConfig.OMAdminUser,
			"--adminpass="+cfConfig.OMAdminPassword,
			"--tile=ops-manager",
		)

		restoreCmd := exec.Command(
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

		backupSession, err := gexec.Start(backupCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(backupSession, 1200).Should(gexec.Exit(0))

		restoreSession, err := gexec.Start(restoreCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(restoreSession, 1800).Should(gexec.Exit(0))

		// TODO get installation_settings from new machine
		// should broadly match installation.json on disk
	})
})

var _ = XDescribe("CFOps Elastic Runtime plugin", func() {
	cfopsPath := "/tmp/cfops"
	backupPath := "/tmp/cfops-backup-" + uuid.NewRandom().String()
	BeforeEach(func() {
		pushTestApp(cfConfig)
	})

	AfterEach(func() {
		deleteTestApp(cfConfig)

		if cfopsPath != "" {
			remoteExecute("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "rm -rf "+cfopsPath, os.Stdout)
		}

		if backupPath != "" {
			remoteExecute("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "rm -rf "+backupPath, os.Stdout)
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

		scpHelper("ubuntu", cfConfig.OMHostname, 22, cfopsLinuxExecutablePath, cfopsPath, cfConfig.OMSSHKey)
		remoteExecute("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "chmod +x /tmp/cfops", os.Stdout)

		fmt.Println("Backing up ERT...")
		remoteExecute("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, backupCmd, os.Stdout)

		deleteTestApp(cfConfig)

		fmt.Println("Restoring ERT...")
		remoteExecute("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, restoreCmd, os.Stdout)

		cfDo("target", "-o", cfConfig.OrgName, "-s", cfConfig.SpaceName)
		Eventually(cf.Cf("apps")).Should(gbytes.Say(cfConfig.AppName))
	})
})
