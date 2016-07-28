package system

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var CfAPI string
var CfUser string
var CfPassword string
var OMAdminUser string
var OMAdminPassword string
var OMHostname string
var OMSSHKey string

func TestBrokerintegration(t *testing.T) {
	CfAPI = os.Getenv("CF_API_URL")
	CfUser = os.Getenv("CF_USER")
	CfPassword = os.Getenv("CF_PASSWORD")
	OMAdminUser = os.Getenv("OM_USER")
	OMAdminPassword = os.Getenv("OM_PASSWORD")
	OMHostname = os.Getenv("OM_HOSTNAME")
	OMSSHKey = os.Getenv("OM_SSH_KEY")

	RegisterFailHandler(Fail)
	RunSpecs(t, "cfops System Test Suite")
}
