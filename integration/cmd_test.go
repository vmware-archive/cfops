package cfopsintegration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var cfopsExecutablePath string
var err error

const (
	help string = `cfops - Cloud Foundry Operations Tool`
)

var _ = BeforeSuite(func() {
	cfopsExecutablePath, err = gexec.Build("github.com/pivotalservices/cfops/cmd/cfops")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("cfops cmd", func() {
	It("prints the help page", func() {
		command := exec.Command(cfopsExecutablePath)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred(), "couldn't run cfops executable")
		Eventually(session.Out).Should(gbytes.Say(help))
	})
})
