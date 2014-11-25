package backup

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
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
	backupscript := "./backup/scripts/backup.sh"
	params := []string{"usage"}
	if len(args) > 0 {
		cmd.CommandRunner.Run(backupscript, params...)
		return nil
	}

	ops_manager_host := cmd.Config.OpsManagerHost
	tempest_passwd := cmd.Config.TempestPassword
	ops_manager_admin := cmd.Config.OpsManagerAdminUser
	ops_manager_admin_passwd := cmd.Config.OpsManagerAdminPassword
	backup_location := cmd.Config.BackupFileLocation

	currenttime := time.Now().Local()
	formattedtime := currenttime.Format("2006_01_02")
	backup_dir := backup_location + "/Backup_" + formattedtime

	deployment_dir := backup_location + "/" + formattedtime + "/deployments"
	database_dir := backup_location + "/" + formattedtime + "/database"
	nfs_dir := backup_location + "/" + formattedtime + "/nfs_share"
	jsonfile := backup_dir + "/installation.json"

	createDirectories(backup_dir, deployment_dir, database_dir, nfs_dir)

	backupDeploymentFiles(ops_manager_host, tempest_passwd, deployment_dir)

	extractEncryptionKey(cmd, backupscript, backup_dir, deployment_dir)

	exportInstallationSettings(ops_manager_host, ops_manager_admin, ops_manager_admin_passwd, jsonfile)

	ip, username, password := verifyBoshLogin(jsonfile)

	deploymentName := downloadElasticRuntimeDeploymentFile(ip, username, password, backup_dir)

	ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backup_dir)

	toggleCCJobs(ip, username, password, deploymentName, ccJobs, "stopped")

	exportCCDB(cmd, backupscript, jsonfile, database_dir)

	exportUAADB(cmd, backupscript, jsonfile, database_dir)

	exportConsoleDB(cmd, backupscript, jsonfile, database_dir)

	// arguments = []string{jsonfile, "cf", "nfs_server", "vcap"}
	// password = getPassword(arguments)
	// ip = getIP(arguments)
	//
	// src_url := "vcap@" + ip + ":/var/vcap/store/shared/**/*"
	// dest_url := nfs_dir + "/"
	// options := "-P 22 -rp"
	// ScpCli([]string{options, src_url, dest_url, password})

	toggleCCJobs(ip, username, password, deploymentName, ccJobs, "started")

	return nil
}

func createDirectories(backup_dir string, deployment_dir string, database_dir string, nfs_dir string) {
	os.MkdirAll(backup_dir, 0777)
	os.MkdirAll(deployment_dir, 0777)
	os.MkdirAll(database_dir, 0777)
	os.MkdirAll(nfs_dir, 0777)
}

func backupDeploymentFiles(ops_manager_host string, tempest_passwd string, deployment_dir string) {
	src_url := "tempest@" + ops_manager_host + ":/var/tempest/workspaces/default/deployments/*.yml"
	dest_url := deployment_dir
	options := "-P 22 -r"

	ScpCli([]string{options, src_url, dest_url, tempest_passwd})

	src_url = "tempest@" + ops_manager_host + ":/var/tempest/workspaces/default/deployments/micro/*.yml"
	ScpCli([]string{options, src_url, dest_url, tempest_passwd})
}

func extractEncryptionKey(cmd BackupCommand, backupscript string, backup_dir string, deployment_dir string) {
	params := []string{"export_Encryption_key", backup_dir, deployment_dir}
	cmd.CommandRunner.Run(backupscript, params...)
}

func exportInstallationSettings(ops_manager_host string, ops_manager_admin string, ops_manager_admin_passwd string, jsonfile string) {
	connectionUrl := "https://" + ops_manager_host + "/api/installation_settings"

	resp, err := invoke("GET", connectionUrl, ops_manager_admin, ops_manager_admin_passwd, false)
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

}

func verifyBoshLogin(jsonfile string) (directorIP string, directorUser string, directorPassword string) {
	var username = "director"
	ip, password := getConnectionDetails(jsonfile, "microbosh", "director", username)

	connectionUrl := "https://" + ip + ":25555/info"

	resp, err := invoke("GET", connectionUrl, username, password, false)
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
	password := getPassword(arguments)
	ip := getIP(arguments)

	return ip, password
}

func downloadElasticRuntimeDeploymentFile(ip string, username string, password string, backup_dir string) string {
	connectionUrl := "https://" + ip + ":25555/deployments"

	resp, err := invoke("GET", connectionUrl, username, password, false)
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

	deploymentObjects := getDeploymentObject(contents)

	for _, deploymentObject := range deploymentObjects {
		if strings.Contains(deploymentObject.Name, "cf-") {
			cfDeploymentName = deploymentObject.Name
			break
		}
	}

	fmt.Println(cfDeploymentName)

	connectionUrl = "https://" + ip + ":25555/deployments/" + cfDeploymentName

	resp, err = invoke("GET", connectionUrl, username, password, true)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid Bosh Director Credentials")
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	cfManifestFile := backup_dir + "/" + cfDeploymentName + ".yml"

	err = ioutil.WriteFile(cfManifestFile, contents, 0644)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	return cfDeploymentName

}

func getAllCloudControllerVMs(ip string, username string, password string, deploymentName string, backup_dir string) []string {
	connectionUrl := "https://" + ip + ":25555/deployments/" + deploymentName + "/vms"

	resp, err := invoke("GET", connectionUrl, "director", password, false)
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

	vmObjects := getVMSObject(contents)
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

	fmt.Println(ccjobs)

	return ccjobs
}

func toggleCCJobs(ip string, username string, password string, deploymentName string, ccjobs []string, state string) {
	for i, ccjob := range ccjobs {
		connectionUrl := "https://" + ip + ":25555/deployments/" + deploymentName + "/jobs/" + ccjob + "/" + strconv.Itoa(i) + "?state=" + state

		resp, err := invoke("PUT", connectionUrl, username, password, true)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		if resp.StatusCode == 302 {
			fmt.Println("response %s", resp)
		} else if resp.StatusCode != 200 {
			fmt.Println("Invalid Bosh Director Credentials")
			os.Exit(1)
		}
	}
}

func invoke(method string, connectionUrl string, username string, password string, isYaml bool) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, string(connectionUrl), nil)
	req.SetBasicAuth(username, password)

	if isYaml {
		req.Header.Set("Content-Type", "text/yaml")
	}

	fmt.Println("Request : %s", req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	return resp, err
}

func exportCCDB(cmd BackupCommand, backupscript string, jsonfile string, database_dir string) {
	ip, password := getConnectionDetails(jsonfile, "cf", "ccdb", "admin")

	dbparams := []string{"export_db", ip, "admin", password, "2544", "ccdb", database_dir + "/ccdb.sql"}

	cmd.CommandRunner.Run(backupscript, dbparams...)
}

func exportUAADB(cmd BackupCommand, backupscript string, jsonfile string, database_dir string) {
	ip, password := getConnectionDetails(jsonfile, "cf", "uaadb", "root")

	dbparams := []string{"export_db", ip, "root", password, "2544", "uaa", database_dir + "/uaa.sql"}

	cmd.CommandRunner.Run(backupscript, dbparams...)
}

func exportConsoleDB(cmd BackupCommand, backupscript string, jsonfile string, database_dir string) {
	ip, password := getConnectionDetails(jsonfile, "cf", "consoledb", "root")

	dbparams := []string{"export_db", ip, "root", password, "2544", "console", database_dir + "/console.sql"}

	cmd.CommandRunner.Run(backupscript, dbparams...)
}
