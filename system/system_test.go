package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
)

var cfopsExecutablePath string
var err error

var cfConfig Config

var _ = BeforeSuite(func() {
	os.Setenv("GOOS", "linux")
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
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("CFOps Elastic Runtime plugin", func() {
	cfopsPath := "/tmp/cfops"
	backupPath := "/tmp/cfops-backup-" + uuid.NewRandom().String()
	BeforeEach(func() {
		pushTestApp(cfConfig)
	})

	AfterEach(func() {
		deleteTestApp(cfConfig)

		if cfopsPath != "" {
			rexecHelper("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "rm -rf "+cfopsPath, os.Stdout)
		}

		if backupPath != "" {
			rexecHelper("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "rm -rf "+backupPath, os.Stdout)
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

		scpHelper("ubuntu", cfConfig.OMHostname, 22, cfopsExecutablePath, cfopsPath, cfConfig.OMSSHKey)
		rexecHelper("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, "chmod +x /tmp/cfops", os.Stdout)

		fmt.Println("Backing up ERT...")
		rexecHelper("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, backupCmd, os.Stdout)

		deleteTestApp(cfConfig)

		fmt.Println("Restoring ERT...")
		rexecHelper("ubuntu", cfConfig.OMHostname, 22, cfConfig.OMSSHKey, restoreCmd, os.Stdout)

		cfDo("target", "-o", cfConfig.OrgName, "-s", cfConfig.SpaceName)
		Eventually(cf.Cf("apps")).Should(gbytes.Say(cfConfig.AppName))
	})
})
