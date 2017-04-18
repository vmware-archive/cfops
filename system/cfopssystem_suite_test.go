package system

import (
	"encoding/json"
	"os"
	"testing"

	"code.cloudfoundry.org/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
)

var cfConfig Config

func TestBrokerintegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cfops System Test Suite")
}

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
