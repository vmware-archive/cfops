package backup

// Tile is a deployable component that can be backed up
type Tile interface {
	Backup() error
}

type BackupContext struct {
	TargetDir string
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
// cc := NewCloudController(ip, username, password, deploymentName, "stopped")
// cc.ToggleJobs(CloudControllerJobs(ccJobs))
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
//
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
//
//
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
