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
		wg      sync.WaitGroup
		fileRef *os.File
		reader  io.Reader
		ec      chan error
	)
	defer fileRef.Close()
	ec = make(chan error)

	if fileRef, err = os.Open(context.JsonFile); err == nil {
		r, w := io.Pipe()
		reader = io.TeeReader(fileRef, w)
		wg.Add(2)
		go readAdminUserCredentials(&wg, creds, reader, w, product, component, ec)
		go readVcapUserCredentials(&wg, creds, r, product, component, ec)
		wg.Wait()
		err = readErrors(ec)
	}
	return
}

func readErrors(ec chan error) (err error) {

	if len(ec) > 0 {

		for e := range ec {
			err = fmt.Errorf("%v ; %v", err, e)
		}
	}
	return
}

func readAdminUserCredentials(wg *sync.WaitGroup, creds *credentials, reader io.Reader, writer io.WriteCloser, product, component string, ec chan error) {
	var err error
	defer wg.Done()
	defer writer.Close()
	(*creds).AdminUser = "admin"

	if (*creds).Ip, (*creds).AdminPass, err = GetPasswordAndIP(reader, product, component, (*creds).AdminUser); err != nil {
		ec <- err
	}
}

func readVcapUserCredentials(wg *sync.WaitGroup, creds *credentials, reader io.Reader, product, component string, ec chan error) {
	var err error
	defer wg.Done()
	(*creds).VcapUser = "vcap"

	if (*creds).Ip, (*creds).VcapPass, err = GetPasswordAndIP(reader, product, component, (*creds).VcapUser); err != nil {
		ec <- err
	}
}
