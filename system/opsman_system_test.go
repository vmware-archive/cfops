package system

import (
	"os"
	"os/exec"

	"code.cloudfoundry.org/lager"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var cfopsExecutablePath string
var cfopsLinuxExecutablePath string
var logger lager.Logger

func checkNoSecretsInSession(session []byte) {
	if cfConfig.OMAdminPassword != "" {
		Expect(session).NotTo(ContainSubstring(cfConfig.OMAdminPassword))
	}
	Expect(session).NotTo(ContainSubstring("RSA PRIVATE KEY"))
}

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

var _ = Describe("CFOps Ops Manager plugin", func() {
	It("backs up and restores successfully", func() {
		if os.Getenv("ONLY_ERT") == "true" {
			return
		}

		vm := createInstance("cfops-test", cfConfig.AmiID, cfConfig.SecurityGroup)

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

		Eventually(backupSession, 1200).Should(gexec.Exit(0))
		checkNoSecretsInSession(backupSession.Out.Contents())
		checkNoSecretsInSession(backupSession.Err.Contents())

		if os.Getenv("OM_VERSION") == "1.6" {
			createAdminUser(newVMIP, cfConfig.OMAdminUser, cfConfig.OMAdminPassword)
		}

		restoreSession, err := gexec.Start(restoreCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(restoreSession, 1800).Should(gexec.Exit(0))
		checkNoSecretsInSession(restoreSession.Out.Contents())
		checkNoSecretsInSession(restoreSession.Err.Contents())

		time.Sleep(2 * time.Minute) // TODO make this better

		checkOpsManagersIdentical(cfConfig.OMHostname, newVMIP)

		vm.Destroy()
	})
})
