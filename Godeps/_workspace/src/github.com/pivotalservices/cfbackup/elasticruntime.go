package cfbackup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/xchapter7x/lo"
)

const (
	//ERDefaultSystemUser - default user for system vms
	ERDefaultSystemUser = "vcap"
	//ERDirectorInfoURL - url format for a director info endpoint
	ERDirectorInfoURL = "https://%s:25555/info"
	//ERBackupDir - default er backup dir
	ERBackupDir = "elasticruntime"
	//ERVmsURL - url format for a vms url
	ERVmsURL = "https://%s:25555/deployments/%s/vms"
	//ERDirector -- key
	ERDirector = "DirectorInfo"
	//ERConsole -- key
	ERConsole = "ConsoledbInfo"
	//ERUaa -- key
	ERUaa = "UaadbInfo"
	//ERCc -- key
	ERCc = "CcdbInfo"
	//ERMySQL -- key
	ERMySQL = "MysqldbInfo"
	//ERNfs -- key
	ERNfs = "NfsInfo"
	//ERBackupFileFormat -- format of archive filename
	ERBackupFileFormat = "%s.backup"
	//ERInvalidDirectorCredsMsg -- error message for invalid creds on director
	ERInvalidDirectorCredsMsg = "invalid director credentials"
	//ERNoPersistenceArchives -- error message for persistence stores
	ERNoPersistenceArchives = "there are no persistence stores in the list"
	//ERFileDoesNotExist -- error message for file does not exist
	ERFileDoesNotExist = "file does not exist"
	//ErrERDBBackupFailure -- error message for backup failure
	ErrERDBBackupFailure = "failed to backup database"
	//ERVersionEnvFlag -- env flag from ER version toggle
	ERVersionEnvFlag = "ER_VERSION"
	//ERVersion16 -- value for 1.6 toggle
	ERVersion16 = "1.6"
)

const (
	//ImportArchive --
	ImportArchive = iota
	//ExportArchive --
	ExportArchive
)

var (
	//ErrERDirectorCreds - error for director creds
	ErrERDirectorCreds = errors.New(ERInvalidDirectorCredsMsg)
	//ErrEREmptyDBList - error for db list empty
	ErrEREmptyDBList = errors.New(ERNoPersistenceArchives)
	//ErrERInvalidPath - invalid filepath error
	ErrERInvalidPath = &os.PathError{Err: errors.New(ERFileDoesNotExist)}
	//ErrERDBBackup - error for db backup failures
	ErrERDBBackup = errors.New(ErrERDBBackupFailure)
)

//BoshName - function which returns proper bosh component name for given version
func BoshName() (bosh string) {
	switch os.Getenv(ERVersionEnvFlag) {
	case ERVersion16:
		bosh = "p-bosh"
	default:
		bosh = "microbosh"
	}
	return
}

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	JSONFile          string
	SystemsInfo       map[string]SystemDump
	PersistentSystems []SystemDump
	HTTPGateway       http.HttpGateway
	InstallationName  string
	BackupContext
}

//CCJob - a cloud controller job object
type CCJob struct {
	Job   string
	Index int
}

// NewElasticRuntime initializes an ElasticRuntime intance
var NewElasticRuntime = func(jsonFile string, target string) *ElasticRuntime {
	var (
		uaadbInfo *PgInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:   "cf",
				Component: "uaadb",
				Identity:  "root",
			},
			Database: "uaa",
		}
		consoledbInfo *PgInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:   "cf",
				Component: "consoledb",
				Identity:  "root",
			},
			Database: "console",
		}
		ccdbInfo *PgInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:   "cf",
				Component: "ccdb",
				Identity:  "admin",
			},
			Database: "ccdb",
		}
		mysqldbInfo *MysqlInfo = &MysqlInfo{
			SystemInfo: SystemInfo{
				Product:   "cf",
				Component: "mysql",
				Identity:  "root",
			},
			Database: "mysql",
		}
		directorInfo *SystemInfo = &SystemInfo{
			Product:   BoshName(),
			Component: "director",
			Identity:  "director",
		}
		nfsInfo *NfsInfo = &NfsInfo{
			SystemInfo: SystemInfo{
				Product:   "cf",
				Component: "nfs_server",
				Identity:  "vcap",
			},
		}
	)

	context := &ElasticRuntime{
		JSONFile:      jsonFile,
		BackupContext: NewBackupContext(target, cfenv.CurrentEnv()),
		SystemsInfo: map[string]SystemDump{
			ERDirector: directorInfo,
			ERConsole:  consoledbInfo,
			ERUaa:      uaadbInfo,
			ERCc:       ccdbInfo,
			ERMySQL:    mysqldbInfo,
			ERNfs:      nfsInfo,
		},
		PersistentSystems: []SystemDump{
			consoledbInfo,
			uaadbInfo,
			ccdbInfo,
			nfsInfo,
			mysqldbInfo,
		},
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	return context.backupRestore(ExportArchive)
}

// Restore performs a restore of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Restore() (err error) {
	err = context.backupRestore(ImportArchive)
	return
}

func (context *ElasticRuntime) backupRestore(action int) (err error) {
	var (
		ccJobs []CCJob
	)

	if err = context.ReadAllUserCredentials(); err == nil && context.directorCredentialsValid() {
		lo.G.Debug("Retrieving All CC VMs")
		manifest, erro := context.getManifest()
		if err != nil {
			return erro
		}
		if ccJobs, err = context.getAllCloudControllerVMs(); err == nil {
			directorInfo := context.SystemsInfo[ERDirector]
			cloudController := NewCloudController(directorInfo.Get(SDIP), directorInfo.Get(SDUser), directorInfo.Get(SDPass), context.InstallationName, manifest, ccJobs)
			lo.G.Debug("Setting up CC jobs")
			defer cloudController.Start()
			cloudController.Stop()
		}
		lo.G.Debug("Running db action")
		if len(context.PersistentSystems) > 0 {
			err = context.RunDbAction(context.PersistentSystems, action)
			if err != nil {
				lo.G.Error("Error backing up db", err)
				err = ErrERDBBackup
			}
		} else {
			err = ErrEREmptyDBList
		}
	} else if err == nil {
		err = ErrERDirectorCreds
	}
	return
}

func (context *ElasticRuntime) getAllCloudControllerVMs() (ccvms []CCJob, err error) {

	lo.G.Debug("Entering getAllCloudControllerVMs() function")
	directorInfo := context.SystemsInfo[ERDirector]
	connectionURL := fmt.Sprintf(ERVmsURL, directorInfo.Get(SDIP), context.InstallationName)
	lo.G.Debug("getAllCloudControllerVMs() function", log.Data{"connectionURL": connectionURL, "directorInfo": directorInfo})
	gateway := context.HTTPGateway
	if gateway == nil {
		gateway = http.NewHttpGateway()
	}
	lo.G.Debug("Retrieving CC vms")
	if resp, err := gateway.Get(http.HttpRequestEntity{
		Url:         connectionURL,
		Username:    directorInfo.Get(SDUser),
		Password:    directorInfo.Get(SDPass),
		ContentType: "application/json",
	})(); err == nil {
		var jsonObj []VMObject

		lo.G.Debug("Unmarshalling CC vms")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &jsonObj); err == nil {
			ccvms, err = GetCCVMs(jsonObj)
		}
	}
	return
}

//RunDbAction - run a db action dump/import against a list of systemdump types
func (context *ElasticRuntime) RunDbAction(dbInfoList []SystemDump, action int) (err error) {

	for _, info := range dbInfoList {
		lo.G.Debug(fmt.Sprintf("RunDbAction info: %+v", info))

		if err = info.Error(); err == nil {
			err = context.readWriterArchive(info, context.TargetDir, action)
		} else {
			// Don't error out yet until issue #111461510 is resolved.
			continue
		}
	}
	return
}

func (context *ElasticRuntime) readWriterArchive(dbInfo SystemDump, databaseDir string, action int) (err error) {
	filename := fmt.Sprintf(ERBackupFileFormat, dbInfo.Get(SDComponent))
	filepath := path.Join(databaseDir, filename)

	var pb PersistanceBackup

	if pb, err = dbInfo.GetPersistanceBackup(); err == nil {
		switch action {
		case ImportArchive:
			lo.G.Debug("Restoring %s", dbInfo.Get(SDComponent))
			var backupReader io.ReadCloser
			if backupReader, err = context.Reader(filepath); err == nil {
				defer backupReader.Close()
				err = pb.Import(backupReader)
				lo.G.Debug("Done restoring %s", dbInfo.Get(SDComponent))
			}
		case ExportArchive:
			lo.G.Info("Exporting %s", dbInfo.Get(SDComponent))
			var backupWriter io.WriteCloser
			if backupWriter, err = context.Writer(filepath); err == nil {
				defer backupWriter.Close()
				err = pb.Dump(backupWriter)
				lo.G.Debug("Done backing up %s", dbInfo.Get(SDComponent))
			}
		}
	}
	return
}

//ReadAllUserCredentials - get all user creds from the installation json
func (context *ElasticRuntime) ReadAllUserCredentials() (err error) {
	var (
		fileRef *os.File
		jsonObj InstallationCompareObject
	)
	defer fileRef.Close()

	if fileRef, err = os.Open(context.JSONFile); err == nil {
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
		var (
			ip    string
			pass  string
			vpass string
		)
		sysInfo.Set(SDVcapUser, ERDefaultSystemUser)
		sysInfo.Set(SDUser, sysInfo.Get(SDIdentity))

		if ip, pass, err = GetPasswordAndIP(jsonObj, sysInfo.Get(SDProduct), sysInfo.Get(SDComponent), sysInfo.Get(SDIdentity)); err == nil {
			sysInfo.Set(SDIP, ip)
			sysInfo.Set(SDPass, pass)
			lo.G.Debug("%s credentials for %s from installation.json are %s", name, sysInfo.Get(SDComponent), sysInfo.Get(SDIdentity), pass)
			_, vpass, err = GetPasswordAndIP(jsonObj, sysInfo.Get(SDProduct), sysInfo.Get(SDComponent), sysInfo.Get(SDVcapUser))
			sysInfo.Set(SDVcapPass, vpass)
			context.SystemsInfo[name] = sysInfo
		}
	}
	return
}

func (context *ElasticRuntime) directorCredentialsValid() (ok bool) {
	var directorInfo SystemDump

	if directorInfo, ok = context.SystemsInfo[ERDirector]; ok {
		connectionURL := fmt.Sprintf(ERDirectorInfoURL, directorInfo.Get(SDIP))
		gateway := context.HTTPGateway
		if gateway == nil {
			gateway = http.NewHttpGateway()
		}
		_, err := gateway.Get(http.HttpRequestEntity{
			Url:         connectionURL,
			Username:    directorInfo.Get(SDUser),
			Password:    directorInfo.Get(SDPass),
			ContentType: "application/json",
		})()
		ok = (err == nil)
	}
	return
}

func (context *ElasticRuntime) getManifest() (manifest string, err error) {
	directorInfo, _ := context.SystemsInfo[ERDirector]
	director := NewDirector(directorInfo.Get(SDIP), directorInfo.Get(SDUser), directorInfo.Get(SDPass), 25555)
	mfs, err := director.GetDeploymentManifest(context.InstallationName)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(mfs)
	if err != nil {
		return
	}
	return string(data), nil
}
