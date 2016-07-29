package system

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	librssh "github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/virtualmachine/aws"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pborman/uuid"
	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
)

func pushTestApp(config Config) {
	fmt.Println("Pushing test app...")

	cfDo("api", config.APIEndpoint, "--skip-ssl-validation")
	cfDo("auth", config.AdminUser, config.AdminPassword)
	cfDo("create-org", config.OrgName)
	cfDo("target", "-o", config.OrgName)
	cfDo("create-space", config.SpaceName)
	cfDo("target", "-s", config.SpaceName)
	cfDo("push", config.AppName, "-p", config.AppPath)

	fmt.Println("Done pushing test app.")
}

func getAuthMethod(pemkeycontents []byte) (authMethod []ssh.AuthMethod) {
	keySigner, _ := ssh.ParsePrivateKey(pemkeycontents)
	authMethod = []ssh.AuthMethod{
		ssh.PublicKeys(keySigner),
	}
	return
}

func scpHelper(sshuser, host string, port int, localpath, remotepath, pemkeycontents string) {
	f, err := os.Open(localpath)
	Expect(err).ToNot(HaveOccurred())
	remoteConn := osutils.NewRemoteOperationsWithPath(command.SshConfig{
		Username: sshuser,
		Host:     host,
		Port:     port,
		SSLKey:   pemkeycontents,
	}, remotepath)
	err = remoteConn.UploadFile(f)
	Expect(err).ToNot(HaveOccurred())
}

func remoteExecute(sshuser, host string, port int, pemkeycontents, remotecommand string, rstdout io.Writer) (err error) {
	remoteConn, err := command.NewRemoteExecutor(command.SshConfig{
		Username: sshuser,
		Host:     host,
		Port:     port,
		SSLKey:   pemkeycontents,
	})
	Expect(err).ToNot(HaveOccurred())
	return remoteConn.Execute(rstdout, remotecommand)
}

func deleteTestApp(config Config) {
	fmt.Println("Deleting test app...")

	cfDo("api", config.APIEndpoint, "--skip-ssl-validation")
	cfDo("auth", config.AdminUser, config.AdminPassword)
	cfDo("target", "-o", config.OrgName, "-s", config.SpaceName)
	cfDo("delete", "-f", config.AppName)
	cfDo("delete-org", "-f", config.OrgName)

	fmt.Println("Done deleting test app.")
}

func createInstance(amznkeyname string, amiID string, securityGroup string) *aws.VM {
	fmt.Println("Creating AWS VM...")

	vm := &aws.VM{
		Name:         "cfops-test-" + uuid.NewRandom().String(),
		AMI:          amiID,
		InstanceType: "m3.large",
		SSHCreds: librssh.Credentials{
			SSHUser:       "ubuntu",
			SSHPrivateKey: amznkeyname,
		},
		Volumes: []aws.EBSVolume{
			{
				DeviceName: "/dev/sda1",
				VolumeSize: 100,
			},
		},
		Region:        "eu-west-1",
		KeyPair:       amznkeyname,
		SecurityGroup: securityGroup,
	}

	err := aws.ValidCredentials(vm.Region)
	Expect(err).NotTo(HaveOccurred())

	err = vm.Provision()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println("AWS VM created.")

	return vm
}

func cfDo(cmd ...string) {
	Eventually(cf.Cf(cmd...), 90).Should(gexec.Exit(0),
		fmt.Sprintf("Command `cf %s` failed", cmd),
	)
}

//Config ...
type Config struct {
	APIEndpoint     string
	AdminUser       string
	AdminPassword   string
	AppName         string
	OrgName         string
	SpaceName       string
	AppPath         string
	OMAdminUser     string
	OMAdminPassword string
	OMHostname      string
	OMSSHKey        string
	AmiID           string
	SecurityGroup   string
}
