package backup

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cfops/backup/modules/persistence"
	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	NewDumper       func(port int, database, username, password string, sshCfg command.SshConfig) (persistence.Dumper, error)
	JsonFile        string
	DeploymentsFile string
	DbEncryptionKey string
	SystemsInfo     map[string]SystemInfo
	BackupContext
}

type SystemInfo struct {
	Product   string
	Component string
	Identity  string
	Ip        string
	User      string
	Pass      string
	VcapUser  string
	VcapPass  string
}

func (s *SystemInfo) Error() (err error) {
	if s.Product == "" ||
		s.Component == "" ||
		s.Identity == "" ||
		s.Ip == "" ||
		s.User == "" ||
		s.Pass == "" ||
		s.VcapUser == "" ||
		s.VcapPass == "" {
		err = fmt.Errorf("invalid or incomplete system info object: %s", s)
	}
	return
}

// NewElasticRuntime initializes an ElasticRuntime intance
func NewElasticRuntime(jsonFile, deploymentsFile, dbEncryptionKey string, target string) *ElasticRuntime {
	context := &ElasticRuntime{
		NewDumper:       persistence.NewPgRemoteDump,
		JsonFile:        jsonFile,
		DeploymentsFile: deploymentsFile,
		DbEncryptionKey: dbEncryptionKey,
		BackupContext: BackupContext{
			TargetDir: target,
		},
		SystemsInfo: map[string]SystemInfo{
			"ConsoledbInfo": SystemInfo{
				Product:   "cf",
				Component: "consoledb",
				Identity:  "root",
			},
			"UaadbInfo": SystemInfo{
				Product:   "cf",
				Component: "uaadb",
				Identity:  "root",
			},
			"CcdbInfo": SystemInfo{
				Product:   "cf",
				Component: "ccdb",
				Identity:  "admin",
			},
			"DirectorInfo": SystemInfo{
				Product:   "microbosh",
				Component: "director",
				Identity:  "director",
			},
		},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	err = context.ReadAllUserCredentials()

	if context.directorCredentialsValid() {
		// deploymentName := getElasticRuntimeDeploymentName(ip, username, password, backupDir)
		// ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)
		// cc := NewCloudController(ip, username, password, deploymentName, "stopped")
		// cc.ToggleJobs(CloudControllerJobs(ccJobs))
		backupDbList := []SystemInfo{
			context.SystemsInfo["ConsoledbInfo"],
			context.SystemsInfo["UaadbInfo"],
			context.SystemsInfo["CcdbInfo"],
		}
		err = context.RunDbBackups(backupDbList)
		//-       arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
		//-       password := utils.GetPassword(arguments)
		//-       ip := utils.GetIP(arguments)
		// BackupNfs(password, ip, outfileref)
		// toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "started")
		// backupMySqlDB(backupscript, jsonfile, databaseDir)
	}
	return
}

func (context *ElasticRuntime) RunDbBackups(dbInfoList []SystemInfo) (err error) {

	for _, info := range dbInfoList {

		if err = info.Error(); err == nil {
			err = context.runPostgresBackup(info, context.TargetDir)
		}

		if err != nil {
			break
		}
	}
	return
}

func (context *ElasticRuntime) runPostgresBackup(dbInfo SystemInfo, databaseDir string) (err error) {
	var (
		outfile        *os.File
		remotePGBackup persistence.Dumper
	)

	sshConfig := command.SshConfig{
		Username: dbInfo.VcapUser,
		Password: dbInfo.VcapPass,
		Host:     dbInfo.Ip,
		Port:     22,
	}

	if remotePGBackup, err = context.NewDumper(2544, dbInfo.Component, dbInfo.User, dbInfo.Pass, sshConfig); err == nil {
		filename := fmt.Sprintf("%s.sql", dbInfo.Component)

		if outfile, err = osutils.SafeCreate(databaseDir, filename); err == nil {
			err = remotePGBackup.Dump(outfile)
		}
	}
	return
}

func (context *ElasticRuntime) ReadAllUserCredentials() (err error) {
	var (
		fileRef *os.File
		jsonObj InstallationCompareObject
	)
	defer fileRef.Close()

	if fileRef, err = os.Open(context.JsonFile); err == nil {

		if jsonObj, err = ReadAndUnmarshal(fileRef); err == nil {
			err = context.assignCredentials(jsonObj)
		}
	}
	return
}

func (context *ElasticRuntime) assignCredentials(jsonObj InstallationCompareObject) (err error) {

	for name, sysInfo := range context.SystemsInfo {
		sysInfo.VcapUser = "vcap"
		sysInfo.User = sysInfo.Identity

		if sysInfo.Ip, sysInfo.Pass, err = GetPasswordAndIP(jsonObj, sysInfo.Product, sysInfo.Component, sysInfo.Identity); err == nil {
			_, sysInfo.VcapPass, err = GetPasswordAndIP(jsonObj, sysInfo.Product, sysInfo.Component, "vcap")
			context.SystemsInfo[name] = sysInfo
		}
	}
	return
}

func (context *ElasticRuntime) directorCredentialsValid() (ok bool) {
	ok = true
	connectionURL := "https://" + context.SystemsInfo["DirectorInfo"].Ip + ":25555/info"

	if resp, err := invoke("GET", connectionURL, context.SystemsInfo["DirectorInfo"].User, context.SystemsInfo["DirectorInfo"].Pass, false); err != nil || resp.StatusCode != 200 {
		ok = false
	}
	return
}
