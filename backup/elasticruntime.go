package backup

import (
	"fmt"
	"os"
)

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	DeploymentsFile string
	DbEncryptionKey string
	BackupContext
}

// NewElasticRuntime initializes an ElasticRuntime intance
func NewElasticRuntime(deploymentsFile, dbEncryptionKey string, target string) *ElasticRuntime {
	context := &ElasticRuntime{
		DeploymentsFile: deploymentsFile,
		DbEncryptionKey: dbEncryptionKey,
		BackupContext: BackupContext{
			TargetDir: target,
		},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	// ip, username, password := verifyBoshLogin(jsonfile)
	//
	// deploymentName := getElasticRuntimeDeploymentName(ip, username, password, backupDir)
	//
	// ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)
	//
	// cc := NewCloudController(ip, username, password, deploymentName, "stopped")
	// cc.ToggleJobs(CloudControllerJobs(ccJobs))
	//
	//backupCCDB(backupscript, jsonfile, databaseDir)
	//
	// backupUAADB(backupscript, jsonfile, databaseDir)
	//
	// backupConsoleDB(backupscript, jsonfile, databaseDir)
	//
	//-       arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
	//-       password := utils.GetPassword(arguments)
	//-       ip := utils.GetIP(arguments)
	// BackupNfs(password, ip, outfileref)
	//
	// toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "started")
	//
	// backupMySqlDB(backupscript, jsonfile, databaseDir)

	return nil
}

func backupCCDB(backupscript string, jsonfile string, databaseDir string) (err error) {
	var (
		ip       string
		password string
		fileRef  *os.File
	)

	if fileRef, err = os.Open(jsonfile); err == nil {
		ip, password, err = GetPasswordAndIP(fileRef, "cf", "ccdb", "admin")
	}
	//dbparams := []string{"export_db", ip, "admin", password, "2544", "ccdb", databaseDir + "/ccdb.sql"}
	//IP=$2
	//USERNAME=$3
	//export PGPASSWORD=$4
	//PORT=$5
	//DB=$6
	//DB_FILE=$7
	//executeCommand(backupscript, dbparams...)
	fmt.Println("Completed Backup of Cloud Controller Database", ip, password)
	return
}
