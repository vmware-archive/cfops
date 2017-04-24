package system

import (
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pborman/uuid"
)

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
			remoteExecute(cfConfig.OMHostInfo, "rm -rf "+cfopsPath)
		}

		if backupPath != "" {
			remoteExecute(cfConfig.OMHostInfo, "rm -rf "+backupPath)
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
		_, err := remoteExecute(cfConfig.OMHostInfo, "chmod +x /tmp/cfops")

		Expect(err).NotTo(HaveOccurred())

		By("Backing up ERT...")
		output, err := remoteExecute(cfConfig.OMHostInfo, backupCmd)
		GinkgoWriter.Write(output)
		Expect(err).NotTo(HaveOccurred())
		checkNoSecretsInSession(output)
		deleteTestApp(cfConfig)

		By("Restoring ERT...")
		output, err = remoteExecute(cfConfig.OMHostInfo, restoreCmd)
		GinkgoWriter.Write(output)
		Expect(err).NotTo(HaveOccurred())
		checkNoSecretsInSession(output)

		cfDo("target", "-o", cfConfig.OrgName, "-s", cfConfig.SpaceName)
		Eventually(cf.Cf("apps")).Should(gbytes.Say(cfConfig.AppName))
	})
})
