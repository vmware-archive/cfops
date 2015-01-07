package backup

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cfops/backup/modules/persistence"
	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

const (
	ER_DEFAULT_SYSTEM_USER string = "vcap"
	ER_DIRECTOR_INFO_URL   string = "https://%s:25555/info"
	ER_BACKUP_DIR          string = "elasticruntime"
	ER_NFS_DIR             string = "nfs_share"
	ER_NFS_FILE            string = "nfs.backup"
)

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	NewDumper        func(port int, database, username, password string, sshCfg command.SshConfig) (persistence.Dumper, error)
	JsonFile         string
	DeploymentsFile  string
	DbEncryptionKey  string
	SystemsInfo      map[string]SystemInfo
	DbSystems        []string
	RestRunner       RestAdapter
	InstallationName string
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
		RestRunner:      RestAdapter(invoke),
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
			"MysqldbInfo": SystemInfo{
				Product:   "cf",
				Component: "mysql",
				Identity:  "root",
			},
			"DirectorInfo": SystemInfo{
				Product:   "microbosh",
				Component: "director",
				Identity:  "director",
			},
			"NfsInfo": SystemInfo{
				Product:   "cf",
				Component: "nfs_server",
				Identity:  "vcap",
			},
		},
		DbSystems: []string{"ConsoledbInfo", "UaadbInfo", "CcdbInfo"},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	var (
		backupDbList []SystemInfo
		ccStop       *CloudController
		ccStart      *CloudController
		ccJobs       []string
	)

	if err = context.ReadAllUserCredentials(); err == nil && context.directorCredentialsValid() {
		// ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)
		directorInfo := context.SystemsInfo["DirectorInfo"]
		ccStop = NewCloudController(directorInfo.Ip, directorInfo.User, directorInfo.Pass, context.InstallationName, "stopped")
		ccStart = NewCloudController(directorInfo.Ip, directorInfo.User, directorInfo.Pass, context.InstallationName, "started")
		defer ccStart.ToggleJobs(CloudControllerJobs(ccJobs))
		ccStop.ToggleJobs(CloudControllerJobs(ccJobs))

		for _, n := range context.DbSystems {
			backupDbList = append(backupDbList, context.SystemsInfo[n])
		}

		if err = context.RunDbBackups(backupDbList); err == nil {
			//var outfile *os.File
			//outfile, err = osutils.SafeCreate(context.TargetDir, ER_BACKUP_DIR, ER_NFS_DIR, ER_NFS_FILE)
			//nfsInfo := context.SystemsInfo["NfsInfo"]
			//err = BackupNfs(nfsInfo.Pass, nfsInfo.Ip, outfile)
		}
		// backupMySqlDB(backupscript, jsonfile, databaseDir)
	} else if err == nil {
		err = fmt.Errorf("invalid director credentials")
	}
	return
}

func (context *ElasticRuntime) RunDbBackups(dbInfoList []SystemInfo) (err error) {

	for _, info := range dbInfoList {

		if err = info.Error(); err == nil {
			err = context.runSqlBackup(info, context.TargetDir)
		}

		if err != nil {
			break
		}
	}
	return
}

func (context *ElasticRuntime) runSqlBackup(dbInfo SystemInfo, databaseDir string) (err error) {
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
			err = context.assignCredentialsAndInstallationName(jsonObj)
		}
	}
	return
}

func (context *ElasticRuntime) assignCredentialsAndInstallationName(jsonObj InstallationCompareObject) (err error) {

	if err = context.assignCredentials(jsonObj); err == nil {
		context.InstallationName, err = GetDeploymentName(jsonObj)
	}
	return
}

func (context *ElasticRuntime) assignCredentials(jsonObj InstallationCompareObject) (err error) {

	for name, sysInfo := range context.SystemsInfo {
		sysInfo.VcapUser = ER_DEFAULT_SYSTEM_USER
		sysInfo.User = sysInfo.Identity

		if sysInfo.Ip, sysInfo.Pass, err = GetPasswordAndIP(jsonObj, sysInfo.Product, sysInfo.Component, sysInfo.Identity); err == nil {
			_, sysInfo.VcapPass, err = GetPasswordAndIP(jsonObj, sysInfo.Product, sysInfo.Component, sysInfo.VcapUser)
			context.SystemsInfo[name] = sysInfo
		}
	}
	return
}

func (context *ElasticRuntime) directorCredentialsValid() (ok bool) {
	connectionURL := fmt.Sprintf(ER_DIRECTOR_INFO_URL, context.SystemsInfo["DirectorInfo"].Ip)
	statusCode, _, err := context.RestRunner.Run("GET", connectionURL, context.SystemsInfo["DirectorInfo"].User, context.SystemsInfo["DirectorInfo"].Pass, false)
	ok = (err == nil && statusCode == 200)
	return
}
