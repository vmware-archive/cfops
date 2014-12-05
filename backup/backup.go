package backup

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/pivotalservices/cfops/backup/steps/createfs"
	"github.com/pivotalservices/cfops/ssh"
)

// BackupContext provides context for a backup job
type BackupContext struct {
	Hostname      string
	Username      string
	Password      string
	TPassword     string
	Target        string
	backupDir     string
	deploymentDir string
	databaseDir   string
	nfsDir        string
	json          string
}

func New(hostname string, username string, password string, tempestpassword string, target string) *BackupContext {
	context := &BackupContext{
		Hostname:  hostname,
		Username:  username,
		Password:  password,
		TPassword: tempestpassword,
		Target:    target,
	}

	context.initPaths()
	return context
}

// Run performs a backup of a target Cloud Foundry deployment
func (context *BackupContext) Run() (err error) {
	pipeline := context.getStandardPipeline()
	err = context.ExecutePipeline(pipeline)
	return
}

//ExecutePipeline runs through a pipeline of backup tasks
func (context *BackupContext) ExecutePipeline(pipeline []func() error) (err error) {
	for _, functor := range pipeline {
		err = functor()

		if err != nil {
			break
		}
	}
	return
}

func (context *BackupContext) getStandardPipeline() (pipeline []func() error) {
	pipeline = []func() error{
		context.prepareFilesystem,
		context.backupTempestFiles,
	}
	return
}

func (context *BackupContext) initPaths() {
	context.backupDir = path.Join(context.Target, time.Now().Local().Format("2006_01_02"))
	context.deploymentDir = path.Join(context.backupDir, "deployments")
	context.databaseDir = path.Join(context.backupDir, "database")
	context.nfsDir = path.Join(context.backupDir, "nfs")
	context.json = path.Join(context.backupDir, "installation.json")
}

func (context *BackupContext) backupTempestFiles() error {
	copier := ssh.New("tempest", context.TPassword, context.Hostname, 22)
	return context.backupDeployment(copier)
}

func (context *BackupContext) prepareFilesystem() (err error) {
	directoryList := []string{
		context.backupDir,
		context.deploymentDir,
		context.databaseDir,
		context.nfsDir,
	}
	err = createfs.MultiDirectoryCreate(directoryList, os.MkdirAll)
	return
}

func (context *BackupContext) backupDeployment(copier ssh.Copier) error {
	file, _ := os.Create(path.Join(context.deploymentDir, "deployments.tar.gz"))
	defer file.Close()
	command := "cd /var/tempest/workspaces/default && tar cz deployments"

	err := copier.Copy(file, strings.NewReader(command))
	return err
}

// func Run(context *BackupContext) error {

// extractEncryptionKey(backupscript, backupDir, deploymentDir)
//
// exportInstallationSettings(context.Hostname, context.Username, context.Password, jsonfile)
//
// ip, username, password := verifyBoshLogin(jsonfile)
//
// deploymentName := getElasticRuntimeDeploymentName(ip, username, password, backupDir)
//
// ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)
//
// toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "stopped")
//
// backupCCDB(backupscript, jsonfile, databaseDir)
//
// backupUAADB(backupscript, jsonfile, databaseDir)
//
// backupConsoleDB(backupscript, jsonfile, databaseDir)
//
// backupNfs(jsonfile, nfsDir)
//
// toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "started")
//
// backupMySqlDB(backupscript, jsonfile, databaseDir)
//
// 	return nil
// }

// func extractEncryptionKey(backupscript string, backupDir string, deploymentDir string) {
// 	params := []string{"export_Encryption_key", backupDir, deploymentDir}
// 	executeCommand(backupscript, params...)
// 	fmt.Println("Backup of cloud controller db encryption key completed")
// }
//
// func (context *BackupContext) exportInstallationSettings(jsonfile string) {
// 	connectionURL := "https://" + context.Hostname + "/api/installation_settings"
//
// 	resp, err := invoke("GET", connectionURL, context.Username, context.Password, false)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	defer resp.Body.Close()
// 	contents, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	err = ioutil.WriteFile(jsonfile, contents, 0644)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println("Backup of Installation settings completed")
// }
//
// func verifyBoshLogin(jsonfile string) (directorIP string, directorUser string, directorPassword string) {
// 	var username = "director"
// 	ip, password := getConnectionDetails(jsonfile, "microbosh", "director", username)
//
// 	connectionURL := "https://" + ip + ":25555/info"
//
// 	resp, err := invoke("GET", connectionURL, username, password, false)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Invalid Bosh Director Credentials")
// 		os.Exit(1)
// 	} else {
// 		fmt.Println("Verified Bosh Director Credentials")
// 	}
//
// 	return ip, username, password
// }
//
// func getConnectionDetails(jsonfile string, product string, component string, username string) (string, string) {
// 	arguments := []string{jsonfile, product, component, username}
// 	password := utils.GetPassword(arguments)
// 	ip := utils.GetIP(arguments)
//
// 	return ip, password
// }
//
// func getElasticRuntimeDeploymentName(ip string, username string, password string, backupDir string) string {
// 	connectionURL := "https://" + ip + ":25555/deployments"
//
// 	resp, err := invoke("GET", connectionURL, username, password, false)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Invalid Bosh Director Credentials")
// 		os.Exit(1)
// 	}
//
// 	defer resp.Body.Close()
// 	contents, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	var cfDeploymentName string
//
// 	var deploymentObjects []utils.DeploymentObject
// 	utils.GetJSONObject(contents, &deploymentObjects)
//
// 	for _, deploymentObject := range deploymentObjects {
// 		if strings.Contains(deploymentObject.Name, "cf-") {
// 			cfDeploymentName = deploymentObject.Name
// 			break
// 		}
// 	}
//
// 	fmt.Println(fmt.Sprintf("CF deployment Name : %s", cfDeploymentName))
//
// 	connectionURL = "https://" + ip + ":25555/deployments/" + cfDeploymentName
//
// 	resp, err = invoke("GET", connectionURL, username, password, true)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Invalid Bosh Director Credentials")
// 		os.Exit(1)
// 	}
//
// 	return cfDeploymentName
//
// }
//
// func getAllCloudControllerVMs(ip string, username string, password string, deploymentName string, backupDir string) []string {
// 	connectionURL := "https://" + ip + ":25555/deployments/" + deploymentName + "/vms"
//
// 	resp, err := invoke("GET", connectionURL, "director", password, false)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Invalid Bosh Director Credentials")
// 		os.Exit(1)
// 	}
//
// 	defer resp.Body.Close()
// 	contents, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	var vmObjects []utils.VMObject
// 	utils.GetJSONObject(contents, &vmObjects)
// 	i := 0
//
// 	for _, vmObject := range vmObjects {
// 		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
// 			i++
// 		}
// 	}
//
// 	ccjobs := make([]string, i)
//
// 	for _, vmObject := range vmObjects {
// 		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
// 			ccjobs[vmObject.Index] = vmObject.Job
// 		}
// 	}
//
// 	fmt.Println(fmt.Sprintf("List of Cloud Controller Jobs: %s ", ccjobs))
//
// 	return ccjobs
// }
//
// func toggleCCJobs(backupscript string, ip string, username string, password string, deploymentName string, ccjobs []string, state string) {
// 	serverURL := "https://" + ip + ":25555/"
// 	for i, ccjob := range ccjobs {
// 		connectionURL := serverURL + "deployments/" + deploymentName + "/jobs/" + ccjob + "/" + strconv.Itoa(i) + "?state=" + state
//
// 		params := []string{"toggle_cc_job", connectionURL, username, password}
//
// 		output, cmderr := executeCommand(backupscript, params...)
// 		if cmderr != nil {
// 			fmt.Println(fmt.Sprintf("%s", cmderr))
// 			os.Exit(1)
// 		}
//
// 		updatedURL := strings.Replace(output, "https://"+ip+"/", serverURL, 1)
// 		fmt.Println(fmt.Sprintf("Fetching task status from %s", updatedURL))
//
// 		contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
// 		if taskerr != nil {
// 			fmt.Printf("%s", taskerr)
// 			os.Exit(1)
// 		}
//
// 		var eventObject utils.EventObject
// 		utils.GetJSONObject(contents, &eventObject)
// 		if eventObject.State != "done" {
// 			fmt.Println(fmt.Sprintf("Attempting to change state to %s for cloud controller %s instance %v", state, ccjob, i))
// 		}
//
// 		for eventObject.State != "done" {
// 			fmt.Printf(".")
// 			time.Sleep(2 * time.Second)
// 			contents, taskerr := getTaskEvents("GET", updatedURL, username, password, false)
// 			if taskerr != nil {
// 				fmt.Printf("%s", taskerr)
// 				os.Exit(1)
// 			}
//
// 			utils.GetJSONObject(contents, &eventObject)
// 		}
// 		fmt.Println(fmt.Sprintf("Changed state to %s for cloud controller %s instance %v", state, ccjob, i))
// 	}
// }
//
// func getTaskEvents(method string, url string, username string, password string, isYaml bool) ([]byte, error) {
// 	resp, err := invoke(method, url, username, password, false)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 		os.Exit(1)
// 	}
//
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Invalid Bosh Director Credentials")
// 		os.Exit(1)
// 	}
//
// 	defer resp.Body.Close()
// 	contents, err := ioutil.ReadAll(resp.Body)
//
// 	return contents, err
// }
//
// func invoke(method string, connectionURL string, username string, password string, isYaml bool) (*http.Response, error) {
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
//
// 	req, err := http.NewRequest(method, connectionURL, nil)
// 	req.SetBasicAuth(username, password)
//
// 	if isYaml {
// 		req.Header.Set("Content-Type", "text/yaml")
// 	}
//
// 	resp, err := tr.RoundTrip(req)
// 	if err != nil {
// 		fmt.Printf("Error : %s", err)
// 	}
//
// 	return resp, err
// }
//
// func backupNfs(jsonfile, destDir string) {
// 	fmt.Println("Backup NFS Server")
// 	arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
// 	password := utils.GetPassword(arguments)
// 	ip := utils.GetIP(arguments)
//
// 	config := &ssh.Config{
// 		Username: "vcap",
// 		Password: password,
// 		Host:     ip,
// 		Port:     22,
// 	}
// 	file, _ := os.Create(destDir + "/nfs.tar.gz")
// 	defer file.Close()
// 	command := "cd /var/vcap/store && tar cz shared"
// 	config.Copy(file, command)
// }
//
// func backupCCDB(backupscript string, jsonfile string, databaseDir string) {
// 	fmt.Println("Backup Cloud Controller Database")
// 	ip, password := getConnectionDetails(jsonfile, "cf", "ccdb", "admin")
//
// 	dbparams := []string{"export_db", ip, "admin", password, "2544", "ccdb", databaseDir + "/ccdb.sql"}
//
// 	executeCommand(backupscript, dbparams...)
// 	fmt.Println("Completed Backup of Cloud Controller Database")
// }
//
// func backupUAADB(backupscript string, jsonfile string, databaseDir string) {
// 	fmt.Println("Backup UAA Database")
// 	ip, password := getConnectionDetails(jsonfile, "cf", "uaadb", "root")
//
// 	dbparams := []string{"export_db", ip, "root", password, "2544", "uaa", databaseDir + "/uaa.sql"}
//
// 	executeCommand(backupscript, dbparams...)
// 	fmt.Println("Completed Backup of UAA Database")
// }
//
// func backupConsoleDB(backupscript string, jsonfile string, databaseDir string) {
// 	fmt.Println("Backup Console Database")
// 	ip, password := getConnectionDetails(jsonfile, "cf", "consoledb", "root")
//
// 	dbparams := []string{"export_db", ip, "root", password, "2544", "console", databaseDir + "/console.sql"}
//
// 	executeCommand(backupscript, dbparams...)
// 	fmt.Println("Completed Backup of Console Database")
// }
//
// func backupMySqlDB(backupscript string, jsonfile string, databaseDir string) {
// 	fmt.Println("Backup MySQL Database")
// 	ip, password := getConnectionDetails(jsonfile, "cf", "mysql", "root")
//
// 	dbparams := []string{"export_mysqldb", ip, "root", password, databaseDir + "/user_databases.sql"}
//
// 	executeCommand(backupscript, dbparams...)
// 	fmt.Println("Completed Backup of MySQL Database")
// }
//
// func executeCommand(name string, args ...string) (string, error) {
// 	cmd := exec.Command(name, args...)
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
//
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	return out.String(), err
// }
