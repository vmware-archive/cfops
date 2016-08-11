package system

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/PuerkitoBio/goquery"
	librssh "github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/virtualmachine/aws"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
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

func scpHelper(hostInfo HostInfo, localpath, remotepath string) {
	f, err := os.Open(localpath)
	Expect(err).ToNot(HaveOccurred())
	var remoteConn *osutils.RemoteOperations
	if hostInfo.Password == "" {
		remoteConn = osutils.NewRemoteOperationsWithPath(command.SshConfig{
			Username: hostInfo.Username,
			Host:     hostInfo.Hostname,
			Port:     22,
			SSLKey:   hostInfo.SSHKey,
		}, remotepath)
	} else {
		remoteConn = osutils.NewRemoteOperationsWithPath(command.SshConfig{
			Username: hostInfo.Username,
			Host:     hostInfo.Hostname,
			Password: hostInfo.Password,
			Port:     22,
		}, remotepath)
	}
	err = remoteConn.UploadFile(f)
	Expect(err).ToNot(HaveOccurred())
}

type wrappedClientToEnableDebugging struct {
	innerClient  command.ClientInterface
	session      *ssh.Session
	outputWriter io.Writer
}

func (wc wrappedClientToEnableDebugging) NewSession() (command.SSHSession, error) {
	session, err := wc.innerClient.NewSession()
	if err != nil {
		return nil, err
	}
	sess := session.(*ssh.Session)
	sess.Setenv("LOG_LEVEL", "debug")
	wc.session = sess
	stderrReader, err := sess.StderrPipe()
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(wc.outputWriter, stderrReader)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func remoteExecute(hostInfo HostInfo, remotecommand string) ([]byte, error) {
	var authMethod []ssh.AuthMethod
	if hostInfo.Password == "" {
		keySigner, err := ssh.ParsePrivateKey([]byte(hostInfo.SSHKey))
		if err != nil {
			return nil, err
		}

		authMethod = []ssh.AuthMethod{
			ssh.PublicKeys(keySigner),
		}

	} else {
		authMethod = []ssh.AuthMethod{
			ssh.Password(hostInfo.Password),
		}
	}
	clientconfig := &ssh.ClientConfig{
		User: hostInfo.Username,
		Auth: authMethod,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostInfo.Hostname, 22), clientconfig)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	session.Setenv("LOG_LEVEL", "debug")

	return session.CombinedOutput(remotecommand)
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
		Name:         "cfops-test",
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
	Eventually(cf.Cf(cmd...), 300).Should(gexec.Exit(0),
		fmt.Sprintf("Command `cf %s` failed", cmd),
	)
}

func createAdminUser(hostname, username, password string) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transport}
	setupResp, err := client.Get(fmt.Sprintf("https://%s/setup", hostname))
	Expect(err).NotTo(HaveOccurred())

	doc, err := goquery.NewDocumentFromReader(setupResp.Body)
	Expect(err).NotTo(HaveOccurred())

	token, exists := doc.Find(`input[name="authenticity_token"]`).First().Attr("value")
	Expect(exists).To(BeTrue())

	data := url.Values{}
	data.Add("setup[user_name]", username)
	data.Add("setup[password]", password)
	data.Add("setup[password_confirmation]", password)
	data.Add("setup[eula_accepted]", "0")
	data.Add("setup[eula_accepted]", "true")
	data.Add("authenticity_token", token)

	makeUserRequest, err := http.NewRequest(http.MethodPost, "https://"+hostname+"/setup", bytes.NewBufferString(data.Encode()))
	Expect(err).NotTo(HaveOccurred())
	for _, cookie := range setupResp.Cookies() {
		makeUserRequest.AddCookie(cookie)
	}

	resp, err := client.Do(makeUserRequest)
	body, _ := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred(), string(body))
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

	AmiID         string
	SecurityGroup string
	OMHostInfo    HostInfo
}

type HostInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"host"`
	SSHKey   string
}
