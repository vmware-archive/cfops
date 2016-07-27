package system

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
)

var cfopsExecutablePath string
var err error

var cfConfig Config

var _ = BeforeSuite(func() {
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
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("CFOps Elastic Runtime plugin", func() {
	BeforeEach(func() {
		pushTestApp(cfConfig)
	})

	AfterEach(func() {
		deleteTestApp(cfConfig)
	})

	It("backs up and restores successfully", func() {
		// TODO this will need to be run on the remote OM machine directly
		command := exec.Command(
			cfopsExecutablePath,
			"backup",
			"--opsmanagerhost="+cfConfig.OMHostname,
			"--opsmanageruser=ubuntu",
			"--destination=.",
			"--adminuser="+cfConfig.OMAdminUser,
			"--adminpass="+cfConfig.OMAdminPassword,
			"--tile=elastic-runtime",
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Eventually(session.Out).Should(gexec.Exit(0))
		Expect(err).NotTo(HaveOccurred())

		deleteTestApp(cfConfig)

		// TODO do restore here
		// TODO assert app exists
	})
})
