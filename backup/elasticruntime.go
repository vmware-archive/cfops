package backup

import (
	"fmt"
	"io"
	"os"
	"sync"

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
	BackupContext
}

type credentials struct {
	Ip        string
	VcapUser  string
	VcapPass  string
	AdminUser string
	AdminPass string
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
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	// ip, username, password := verifyBoshLogin(jsonfile)
	// deploymentName := getElasticRuntimeDeploymentName(ip, username, password, backupDir)
	// ccJobs := getAllCloudControllerVMs(ip, username, password, deploymentName, backupDir)
	// cc := NewCloudController(ip, username, password, deploymentName, "stopped")
	// cc.ToggleJobs(CloudControllerJobs(ccJobs))
	err = context.backupCCDB()
	// backupUAADB(backupscript, jsonfile, databaseDir)
	// backupConsoleDB(backupscript, jsonfile, databaseDir)
	//-       arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
	//-       password := utils.GetPassword(arguments)
	//-       ip := utils.GetIP(arguments)
	// BackupNfs(password, ip, outfileref)
	// toggleCCJobs(backupscript, ip, username, password, deploymentName, ccJobs, "started")
	// backupMySqlDB(backupscript, jsonfile, databaseDir)
	return nil
}

func (context *ElasticRuntime) backupCCDB() (err error) {
	var (
		product   string = "cf"
		component string = "ccdb"
	)
	err = context.RunPostgresBackup(product, component, context.TargetDir)
	return
}

func (context *ElasticRuntime) RunPostgresBackup(product, component, databaseDir string) (err error) {
	var (
		creds          credentials
		outfile        *os.File
		remotePGBackup persistence.Dumper
	)

	if err = context.getCredentials(product, component, &creds); err == nil {
		sshConfig := command.SshConfig{
			Username: creds.VcapUser,
			Password: creds.VcapPass,
			Host:     creds.Ip,
			Port:     22,
		}

		if remotePGBackup, err = context.NewDumper(2544, component, creds.AdminUser, creds.AdminPass, sshConfig); err == nil {
			filename := fmt.Sprintf("%s.sql", component)

			if outfile, err = osutils.SafeCreate(databaseDir, filename); err == nil {
				err = remotePGBackup.Dump(outfile)
			}
		}
	}
	return
}

func (context *ElasticRuntime) getCredentials(product, component string, creds *credentials) (err error) {
	var (
		ip            string
		adminPassword string
		vcapPassword  string
		wg            sync.WaitGroup
		fileRef       *os.File
		vcapUser      string = "vcap"
		adminUser     string = "admin"
		reader        io.Reader
	)
	defer fileRef.Close()

	if fileRef, err = os.Open(context.JsonFile); err == nil {
		wg.Add(1)
		r, w := io.Pipe()
		reader = io.TeeReader(fileRef, w)

		go func() {
			defer wg.Done()
			defer w.Close()

			if ip, adminPassword, err = GetPasswordAndIP(reader, product, component, adminUser); err == nil {
				(*creds).Ip = ip
				(*creds).VcapUser = vcapUser
				(*creds).AdminUser = adminUser
				(*creds).AdminPass = adminPassword
			}
		}()

		if _, vcapPassword, err = GetPasswordAndIP(r, product, component, vcapUser); err == nil {
			(*creds).VcapPass = vcapPassword
		}
		wg.Wait()
	}
	return
}
