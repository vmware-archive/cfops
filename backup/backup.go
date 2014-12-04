package backup

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/ssh"
	"github.com/pivotalservices/cfops/system"
	"github.com/pivotalservices/cfops/utils"
	"github.com/xchapter7x/goutil/unpack"
	"github.com/xchapter7x/toggle"
	"github.com/xchapter7x/toggle/engines/localengine"
)

type BackupCommand struct {
	CommandRunner system.CommandRunner
	Logger        *gosteno.Logger
	Config        *BackupConfig
}

func (cmd BackupCommand) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "backup",
		ShortName:   "b",
		Usage:       "backup an existing deployment",
		Description: "backup an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd BackupCommand) Run(args []string) error {
	toggle.Init("CFOPS_BACKUP", localengine.NewLocalEngine())
	toggle.RegisterFeature("CREATE_DIR")

	backupscript := "./backup/scripts/backup.sh"
	params := []string{"usage"}
	if len(args) > 0 {
		cmd.CommandRunner.Run(backupscript, params...)
		return nil
	}

	opsManagerHost := cmd.Config.OpsManagerHost
	tempestPasswd := cmd.Config.TempestPassword
	opsManagerAdmin := cmd.Config.OpsManagerAdminUser
	opsManagerAdminPasswd := cmd.Config.OpsManagerAdminPassword
	backupLocation := cmd.Config.BackupFileLocation

	currenttime := time.Now().Local()
	formattedtime := currenttime.Format("2006_01_02")
	backupDir := backupLocation + "/Backup_" + formattedtime

	deploymentDir := backupDir + "/deployments"
	databaseDir := backupDir + "/database"
	nfsDir := backupDir + "/nfs_share"
	jsonfile := backupDir + "/installation.json"

	responseArray := toggle.Flip("CREATE_DIR", createDirectories, CreateDirectoriesAdaptor, backupDir, deploymentDir, databaseDir, nfsDir)
	var err error
	unpack.Unpack(responseArray, &err)

	if err != nil {
		fmt.Println("Something went terribly wrong")
	}

	backupDeploymentFiles(opsManagerHost, tempestPasswd, deploymentDir)

	extractEncryptionKey(backupscript, backupDir, deploymentDir)

	exportInstallationSettings(opsManagerHost, opsManagerAdmin, opsManagerAdminPasswd, jsonfile)

	ip, username, password := verifyBoshLogin(jsonfile)

	deploymentName := getElasticRuntimeDeploymentName(ip, username, password, backupDir)

	ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)

	toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "stopped")

	backupCCDB(backupscript, jsonfile, databaseDir)

	backupUAADB(backupscript, jsonfile, databaseDir)

	backupConsoleDB(backupscript, jsonfile, databaseDir)

	backupNfs(jsonfile, nfsDir)

	toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "started")

	backupMySqlDB(backupscript, jsonfile, databaseDir)

	return nil
}

func createDirectories(backupDir string, deploymentDir string, databaseDir string, nfsDir string) (err error) {
	os.MkdirAll(backupDir, 0777)
	os.MkdirAll(deploymentDir, 0777)
	os.MkdirAll(databaseDir, 0777)
	os.MkdirAll(nfsDir, 0777)
	fmt.Println("Created all required directories")
	return
}

func backupDeploymentFiles(opsManagerHost string, tempestPasswd string, deploymentDir string) {

	sshCfg := &ssh.SshConfig{
		Username: "tempest",
		Password: tempestPasswd,
		Host:     opsManagerHost,
		Port:     "22",
	}

	file, _ := os.Create(deploymentDir + "/deployments.tar.gz")
	defer file.Close()
	dump := &ssh.DumpToWriter{
		Writer: file,
	}
	command := "cd /var/tempest/workspaces/default && tar cz deployments"

	ssh.DialSsh(sshCfg, command, dump)
	fmt.Println("Backup of Deployment files completed")
}

func extractEncryptionKey(backupscript string, backupDir string, deploymentDir string) {
	params := []string{"export_Encryption_key", backupDir, deploymentDir}
	executeCommand(backupscript, params...)
	fmt.Println("Backup of cloud controller db encryption key completed")
}

func exportInstallationSettings(opsManagerHost string, opsManagerAdmin string, opsManagerAdminPasswd string, jsonfile string) {
	connectionURL := "https://" + opsManagerHost + "/api/installation_settings"

	resp, err := invoke("GET", connectionURL, opsManagerAdmin, opsManagerAdminPasswd, false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(jsonfile, contents, 0644)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Println("Backup of Installation settings completed")
}

func verifyBoshLogin(jsonfile string) (directorIP string, directorUser string, directorPassword string) {
	var username = "director"
	ip, password := getConnectionDetails(jsonfile, "microbosh", "director", username)

	connectionURL := "https://" + ip + ":25555/info"

	resp, err := invoke("GET", connectionURL, username, password, false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	} else {
		fmt.Println("Verified Bosh Director Credentials")
	}

	return ip, username, password
}

func getConnectionDetails(jsonfile string, product string, component string, username string) (string, string) {
	arguments := []string{jsonfile, product, component, username}
	password := utils.GetPassword(arguments)
	ip := utils.GetIP(arguments)

	return ip, password
}

func getElasticRuntimeDeploymentName(ip string, username string, password string, backupDir string) string {
	connectionURL := "https://" + ip + ":25555/deployments"

	resp, err := invoke("GET", connectionURL, username, password, false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var cfDeploymentName string

	var deploymentObjects []utils.DeploymentObject
	utils.GetJSONObject(contents, &deploymentObjects)

	for _, deploymentObject := range deploymentObjects {
		if strings.Contains(deploymentObject.Name, "cf-") {
			cfDeploymentName = deploymentObject.Name
			break
		}
	}

	fmt.Println(fmt.Sprintf("CF deployment Name : %s", cfDeploymentName))

	connectionURL = "https://" + ip + ":25555/deployments/" + cfDeploymentName

	resp, err = invoke("GET", connectionURL, username, password, true)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	}

	return cfDeploymentName

}

func getAllCloudControllerVMs(ip string, username string, password string, deploymentName string, backupDir string) []string {
	connectionURL := "https://" + ip + ":25555/deployments/" + deploymentName + "/vms"

	resp, err := invoke("GET", connectionURL, "director", password, false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var vmObjects []utils.VMObject
	utils.GetJSONObject(contents, &vmObjects)
	i := 0

	for _, vmObject := range vmObjects {
		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
			i++
		}
	}

	ccjobs := make([]string, i)

	for _, vmObject := range vmObjects {
		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
			ccjobs[vmObject.Index] = vmObject.Job
		}
	}

	fmt.Println(fmt.Sprintf("List of Cloud Controller Jobs: %s ", ccjobs))

	return ccjobs
}

func toggleCCJobs(backupscript string, ip string, username string, password string, deploymentName string, ccjobs []string, state string) {
	serverURL := "https://" + ip + ":25555/"
	for i, ccjob := range ccjobs {
		connectionURL := serverURL + "deployments/" + deploymentName + "/jobs/" + ccjob + "/" + strconv.Itoa(i) + "?state=" + state

		params := []string{"toggle_cc_job", connectionURL, username, password}

		output, cmderr := executeCommand(backupscript, params...)
		if cmderr != nil {
			fmt.Println(fmt.Sprintf("%s", cmderr))
			os.Exit(1)
		}

		updatedURL := strings.Replace(output, "https://"+ip+"/", serverURL, 1)
		fmt.Println(fmt.Sprintf("Fetching task status from %s", updatedURL))

		contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
		if taskerr != nil {
			fmt.Printf("%s", taskerr)
			os.Exit(1)
		}

		var eventObject utils.EventObject
		utils.GetJSONObject(contents, &eventObject)
		if eventObject.State != "done" {
			fmt.Println(fmt.Sprintf("Attempting to change state to %s for cloud controller %s instance %v", state, ccjob, i))
		}

		for eventObject.State != "done" {
			fmt.Printf(".")
			time.Sleep(2 * time.Second)
			contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
			if taskerr != nil {
				fmt.Printf("%s", taskerr)
				os.Exit(1)
			}

			utils.GetJSONObject(contents, &eventObject)
		}
		fmt.Println(fmt.Sprintf("Changed state to %s for cloud controller %s instance %v", state, ccjob, i))
	}
}

func getTaskEvents(method string, url string, username string, password string, isYaml bool) ([]byte, error) {
	resp, err := invoke(method, url, username, password, false)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)

	return contents, err
}

func invoke(method string, connectionURL string, username string, password string, isYaml bool) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest(method, connectionURL, nil)
	req.SetBasicAuth(username, password)

	if isYaml {
		req.Header.Set("Content-Type", "text/yaml")
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	return resp, err
}

func backupNfs(jsonfile, destDir string) {
	fmt.Println("Backup NFS Server")
	arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
	password := utils.GetPassword(arguments)
	ip := utils.GetIP(arguments)

	sshCfg := &ssh.SshConfig{
		Username: "vcap",
		Password: password,
		Host:     ip,
		Port:     "22",
	}
	file, _ := os.Create(destDir + "/nfs.tar.gz")
	defer file.Close()
	dump := &ssh.DumpToWriter{
		Writer: file,
	}
	command := "cd /var/vcap/store && tar cz shared"

	ssh.DialSsh(sshCfg, command, dump)
	fmt.Println("Completed Backup of NFS Server")
}

func backupCCDB(backupscript string, jsonfile string, databaseDir string) {
	fmt.Println("Backup Cloud Controller Database")
	ip, password := getConnectionDetails(jsonfile, "cf", "ccdb", "admin")

	dbparams := []string{"export_db", ip, "admin", password, "2544", "ccdb", databaseDir + "/ccdb.sql"}

	executeCommand(backupscript, dbparams...)
	fmt.Println("Completed Backup of Cloud Controller Database")
}

func backupUAADB(backupscript string, jsonfile string, databaseDir string) {
	fmt.Println("Backup UAA Database")
	ip, password := getConnectionDetails(jsonfile, "cf", "uaadb", "root")

	dbparams := []string{"export_db", ip, "root", password, "2544", "uaa", databaseDir + "/uaa.sql"}

	executeCommand(backupscript, dbparams...)
	fmt.Println("Completed Backup of UAA Database")
}

func backupConsoleDB(backupscript string, jsonfile string, databaseDir string) {
	fmt.Println("Backup Console Database")
	ip, password := getConnectionDetails(jsonfile, "cf", "consoledb", "root")

	dbparams := []string{"export_db", ip, "root", password, "2544", "console", databaseDir + "/console.sql"}

	executeCommand(backupscript, dbparams...)
	fmt.Println("Completed Backup of Console Database")
}

func backupMySqlDB(backupscript string, jsonfile string, databaseDir string) {
	fmt.Println("Backup MySQL Database")
	ip, password := getConnectionDetails(jsonfile, "cf", "mysql", "root")

	dbparams := []string{"export_mysqldb", ip, "root", password, databaseDir + "/user_databases.sql"}

	executeCommand(backupscript, dbparams...)
	fmt.Println("Completed Backup of MySQL Database")
}

func executeCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return out.String(), err
}
