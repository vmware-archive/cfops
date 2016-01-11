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
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/xchapter7x/lo"
)

const (
	ER_DEFAULT_SYSTEM_USER        = "vcap"
	ER_DIRECTOR_INFO_URL          = "https://%s:25555/info"
	ER_BACKUP_DIR                 = "elasticruntime"
	ER_VMS_URL                    = "https://%s:25555/deployments/%s/vms"
	ER_DIRECTOR                   = "DirectorInfo"
	ER_CONSOLE                    = "ConsoledbInfo"
	ER_UAA                        = "UaadbInfo"
	ER_CC                         = "CcdbInfo"
	ER_MYSQL                      = "MysqldbInfo"
	ER_NFS                        = "NfsInfo"
	ER_BACKUP_FILE_FORMAT         = "%s.backup"
	ER_INVALID_DIRECTOR_CREDS_MSG = "invalid director credentials"
	ER_NO_PERSISTENCE_ARCHIVES    = "there are no persistence stores in the list"
	ER_FILE_DOES_NOT_EXIST        = "file does not exist"
	ER_DB_BACKUP_FAILURE          = "failed to backup database"
	ER_VERSION_ENV_FLAG           = "ER_VERSION"
	ER_VERSION_16                 = "1.6"
)

const (
	IMPORT_ARCHIVE = iota
	EXPORT_ARCHIVE
)

var (
	ER_ERROR_DIRECTOR_CREDS = errors.New(ER_INVALID_DIRECTOR_CREDS_MSG)
	ER_ERROR_EMPTY_DB_LIST  = errors.New(ER_NO_PERSISTENCE_ARCHIVES)
	ER_ERROR_INVALID_PATH   = &os.PathError{Err: errors.New(ER_FILE_DOES_NOT_EXIST)}
	ER_DB_BACKUP            = errors.New(ER_DB_BACKUP_FAILURE)
)

func BoshName() (bosh string) {
	switch os.Getenv(ER_VERSION_ENV_FLAG) {
	case ER_VERSION_16:
		bosh = "p-bosh"
	default:
		bosh = "microbosh"
	}
	return
}

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	JsonFile          string
	SystemsInfo       map[string]SystemDump
	PersistentSystems []SystemDump
	HttpGateway       HttpGateway
	InstallationName  string
	BackupContext
}

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
		JsonFile:      jsonFile,
		BackupContext: NewBackupContext(target, cfenv.CurrentEnv()),
		SystemsInfo: map[string]SystemDump{
			ER_DIRECTOR: directorInfo,
			ER_CONSOLE:  consoledbInfo,
			ER_UAA:      uaadbInfo,
			ER_CC:       ccdbInfo,
			ER_MYSQL:    mysqldbInfo,
			ER_NFS:      nfsInfo,
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
	return context.backupRestore(EXPORT_ARCHIVE)
}

// Restore performs a restore of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Restore() (err error) {
	err = context.backupRestore(IMPORT_ARCHIVE)
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
			directorInfo := context.SystemsInfo[ER_DIRECTOR]
			cloudController := NewCloudController(directorInfo.Get(SD_IP), directorInfo.Get(SD_USER), directorInfo.Get(SD_PASS), context.InstallationName, manifest, ccJobs)
			lo.G.Debug("Setting up CC jobs")
			defer cloudController.Start()
			cloudController.Stop()
		}
		lo.G.Debug("Running db action")
		if len(context.PersistentSystems) > 0 {
			err = context.RunDbAction(context.PersistentSystems, action)
			if err != nil {
				lo.G.Error("Error backing up db", err)
				err = ER_DB_BACKUP
			}
		} else {
			err = ER_ERROR_EMPTY_DB_LIST
		}
	} else if err == nil {
		err = ER_ERROR_DIRECTOR_CREDS
	}
	return
}

func (context *ElasticRuntime) getAllCloudControllerVMs() (ccvms []CCJob, err error) {

	lo.G.Debug("Entering getAllCloudControllerVMs() function")
	directorInfo := context.SystemsInfo[ER_DIRECTOR]
	connectionURL := fmt.Sprintf(ER_VMS_URL, directorInfo.Get(SD_IP), context.InstallationName)
	lo.G.Debug("getAllCloudControllerVMs() function", log.Data{"connectionURL": connectionURL, "directorInfo": directorInfo})
	gateway := context.HttpGateway
	if gateway == nil {
		gateway = NewHttpGateway()
	}
	lo.G.Debug("Retrieving CC vms")
	if resp, err := gateway.Get(HttpRequestEntity{
		Url:         connectionURL,
		Username:    directorInfo.Get(SD_USER),
		Password:    directorInfo.Get(SD_PASS),
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

func (context *ElasticRuntime) RunDbAction(dbInfoList []SystemDump, action int) (err error) {

	for _, info := range dbInfoList {
		lo.G.Debug(fmt.Sprintf("RunDbAction info: %+v", info))

		if err = info.Error(); err == nil {
			err = context.readWriterArchive(info, context.TargetDir, action)
			lo.G.Debug("backed up db", log.Data{"info": info})
		} else {
			break
		}
	}
	return
}

func (context *ElasticRuntime) readWriterArchive(dbInfo SystemDump, databaseDir string, action int) (err error) {
	filename := fmt.Sprintf(ER_BACKUP_FILE_FORMAT, dbInfo.Get(SD_COMPONENT))
	filepath := path.Join(databaseDir, filename)

	var pb PersistanceBackup

	if pb, err = dbInfo.GetPersistanceBackup(); err == nil {
		switch action {
		case IMPORT_ARCHIVE:
			lo.G.Debug("we are doing something here now")
			var backupReader io.ReadCloser
			if backupReader, err = context.Reader(filepath); err == nil {
				defer backupReader.Close()
				err = pb.Import(backupReader)
			}
		case EXPORT_ARCHIVE:
			lo.G.Info("Exporting database")
			var backupWriter io.WriteCloser
			if backupWriter, err = context.Writer(filepath); err == nil {
				defer backupWriter.Close()
				err = pb.Dump(backupWriter)
			}
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
		var (
			ip    string
			pass  string
			vpass string
		)
		sysInfo.Set(SD_VCAPUSER, ER_DEFAULT_SYSTEM_USER)
		sysInfo.Set(SD_USER, sysInfo.Get(SD_IDENTITY))

		if ip, pass, err = GetPasswordAndIP(jsonObj, sysInfo.Get(SD_PRODUCT), sysInfo.Get(SD_COMPONENT), sysInfo.Get(SD_IDENTITY)); err == nil {
			sysInfo.Set(SD_IP, ip)
			sysInfo.Set(SD_PASS, pass)
			lo.G.Debug("%s credentials for %s from installation.json are %s", name, sysInfo.Get(SD_COMPONENT), sysInfo.Get(SD_IDENTITY), pass)
			_, vpass, err = GetPasswordAndIP(jsonObj, sysInfo.Get(SD_PRODUCT), sysInfo.Get(SD_COMPONENT), sysInfo.Get(SD_VCAPUSER))
			sysInfo.Set(SD_VCAPPASS, vpass)
			context.SystemsInfo[name] = sysInfo
		}
	}
	return
}

func (context *ElasticRuntime) directorCredentialsValid() (ok bool) {
	var directorInfo SystemDump

	if directorInfo, ok = context.SystemsInfo[ER_DIRECTOR]; ok {
		connectionURL := fmt.Sprintf(ER_DIRECTOR_INFO_URL, directorInfo.Get(SD_IP))
		gateway := context.HttpGateway
		if gateway == nil {
			gateway = NewHttpGateway()
		}
		_, err := gateway.Get(HttpRequestEntity{
			Url:         connectionURL,
			Username:    directorInfo.Get(SD_USER),
			Password:    directorInfo.Get(SD_PASS),
			ContentType: "application/json",
		})()
		ok = (err == nil)
	}
	return
}

func (context *ElasticRuntime) getManifest() (manifest string, err error) {
	directorInfo, _ := context.SystemsInfo[ER_DIRECTOR]
	director := NewDirector(directorInfo.Get(SD_IP), directorInfo.Get(SD_USER), directorInfo.Get(SD_PASS), 25555)
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
