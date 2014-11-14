package backup

import (
	"os"
  "fmt"
  "time"
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/system"
)

type Backuper struct {
	CommandRunner system.CommandRunner
}

func New(logger *gosteno.Logger) *Backuper {
	commandRunner := new(system.OSCommandRunner)
	commandRunner.Logger = logger
	return &Backuper{
		CommandRunner: commandRunner,
	}
}

func (installer *Backuper) Backup(args []string) error {
	fmt.Println(len(os.Args), os.Args)

	backupscript := "./backup/scripts/backup.sh"
	params := []string {"usage"}
	if(len(os.Args) < 8) {
		installer.CommandRunner.Run(backupscript, params...)
		return nil
	}
  ops_manager_host := os.Args[3]
  tempest_passwd := os.Args[4]
  ops_manager_admin := os.Args[5]
  ops_manager_admin_passwd := os.Args[6]
  backup_location := os.Args[7]

  currenttime := time.Now().Local()
  formattedtime := currenttime.Format("2006_01_02")
  backup_dir := backup_location + "/" + formattedtime

  deployment_dir := backup_location + "/" + formattedtime + "/deployments"
  database_dir := backup_location + "/" + formattedtime + "/database"
  nfs_dir := backup_location + "/" + formattedtime + "/nfs_share"

  os.MkdirAll(backup_dir, 0777)
  os.MkdirAll(deployment_dir, 0777)
  os.MkdirAll(database_dir, 0777)
  os.MkdirAll(nfs_dir, 0777)

  params = []string {"copy_deployment_files", ops_manager_host, tempest_passwd, ops_manager_admin, ops_manager_admin_passwd, backup_dir, deployment_dir, database_dir, nfs_dir}
  installer.CommandRunner.Run(backupscript, params...)

  params[0] = "export_Encryption_key"
  installer.CommandRunner.Run(backupscript, params...)

  params[0] = "export_installation_settings"
  installer.CommandRunner.Run(backupscript, params...)

	jsonfile := backup_dir + "/installation.yml"

	arguments := []string {jsonfile, "microbosh", "director", "director"}
	password := getPassword(arguments)
	ip := getIP(arguments)

  boshparams := []string {"bosh_login", ip, "director", password}
  installer.CommandRunner.Run(backupscript, boshparams...)

  params[0] = "verify_deployment_backedUp"
  installer.CommandRunner.Run(backupscript, params...)

  params[0] = "bosh_status"
  installer.CommandRunner.Run(backupscript, params...)

  params[0] = "set_bosh_deployment"
  installer.CommandRunner.Run(backupscript, params...)

	params[0] = "export_cloud_controller_bosh_vms"
	installer.CommandRunner.Run(backupscript, params...)

	// params[0] = "stop_cloud_controller"
	// installer.CommandRunner.Run(backupscript, params...)

	arguments = []string {jsonfile, "cf", "ccdb", "admin"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams := []string {"export_db", ip, "admin", password, "2544", "ccdb", database_dir + "/ccdb.sql"}

	installer.CommandRunner.Run(backupscript, dbparams...)

	arguments = []string {jsonfile, "cf", "uaadb", "root"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string {"export_db", ip, "root", password, "2544", "uaa", database_dir + "/uaa.sql"}
	installer.CommandRunner.Run(backupscript, dbparams...)


	arguments = []string {jsonfile, "cf", "consoledb", "root"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string {"export_db", ip, "root", password, "2544", "console", database_dir + "/console.sql"}
	installer.CommandRunner.Run(backupscript, dbparams...)

	arguments = []string {jsonfile, "cf", "nfs_server", "vcap"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string {"export_nfs_server", ip, "vcap", nfs_dir}
	installer.CommandRunner.Run(backupscript, params...)

	// params[0] = "start_cloud_controller"
	// installer.CommandRunner.Run(backupscript, params...)

	return nil
}
