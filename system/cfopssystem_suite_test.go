package system

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var cfConfig Config

func TestBrokerintegration(t *testing.T) {
	cfConfig.APIEndpoint = os.Getenv("CF_API_URL")
	cfConfig.OMAdminUser = os.Getenv("OM_USER")
	cfConfig.OMAdminPassword = os.Getenv("OM_PASSWORD")
	cfConfig.OMHostname = os.Getenv("OM_HOSTNAME")
	cfConfig.OMSSHKey = os.Getenv("OM_SSH_KEY")
	cfConfig.AmiID = os.Getenv("OPSMAN_AMI")
	cfConfig.SecurityGroup = os.Getenv("AWS_SECURITY_GROUP")

	cfConfig.AppName = uuid.NewRandom().String()
	cfConfig.OrgName = uuid.NewRandom().String()
	cfConfig.SpaceName = uuid.NewRandom().String()
	cfConfig.AppPath = "assets/test-app"

	RegisterFailHandler(Fail)
	RunSpecs(t, "cfops System Test Suite")
}
