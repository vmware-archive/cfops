package backup

import (
	"os"
	"time"

	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
)

type BackupCommand struct {
	CommandRunner system.CommandRunner
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
	if len(os.Args) < 7 {
		cmd.CommandRunner.Run(backupscript, params...)
		return nil
	}
	ops_manager_host := os.Args[2]
	tempest_passwd := os.Args[3]
	ops_manager_admin := os.Args[4]
	ops_manager_admin_passwd := os.Args[5]
	backup_location := os.Args[6]

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

	params = []string{"copy_deployment_files", ops_manager_host, tempest_passwd, ops_manager_admin, ops_manager_admin_passwd, backup_dir, deployment_dir, database_dir, nfs_dir}
	cmd.CommandRunner.Run(backupscript, params...)

	params[0] = "export_Encryption_key"
	cmd.CommandRunner.Run(backupscript, params...)

	params[0] = "export_installation_settings"
	cmd.CommandRunner.Run(backupscript, params...)

	jsonfile := backup_dir + "/installation.yml"

	arguments := []string{jsonfile, "microbosh", "director", "director"}
	password := getPassword(arguments)
	ip := getIP(arguments)

	boshparams := []string{"bosh_login", ip, "director", password}
	cmd.CommandRunner.Run(backupscript, boshparams...)

	params[0] = "verify_deployment_backedUp"
	cmd.CommandRunner.Run(backupscript, params...)

	params[0] = "bosh_status"
	cmd.CommandRunner.Run(backupscript, params...)

	params[0] = "set_bosh_deployment"
	cmd.CommandRunner.Run(backupscript, params...)

	params[0] = "export_cloud_controller_bosh_vms"
	cmd.CommandRunner.Run(backupscript, params...)

	// params[0] = "stop_cloud_controller"
	// cmd.CommandRunner.Run(backupscript, params...)

	arguments = []string{jsonfile, "cf", "ccdb", "admin"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams := []string{"export_db", ip, "admin", password, "2544", "ccdb", database_dir + "/ccdb.sql"}

	cmd.CommandRunner.Run(backupscript, dbparams...)

	arguments = []string{jsonfile, "cf", "uaadb", "root"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string{"export_db", ip, "root", password, "2544", "uaa", database_dir + "/uaa.sql"}
	cmd.CommandRunner.Run(backupscript, dbparams...)

	arguments = []string{jsonfile, "cf", "consoledb", "root"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string{"export_db", ip, "root", password, "2544", "console", database_dir + "/console.sql"}
	cmd.CommandRunner.Run(backupscript, dbparams...)

	arguments = []string{jsonfile, "cf", "nfs_server", "vcap"}
	password = getPassword(arguments)
	ip = getIP(arguments)

	dbparams = []string{"export_nfs_server", ip, "vcap", nfs_dir}
	cmd.CommandRunner.Run(backupscript, params...)

	// params[0] = "start_cloud_controller"
	// cmd.CommandRunner.Run(backupscript, params...)

	return nil
}
