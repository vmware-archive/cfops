package cfbackup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pivotal-golang/lager"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/osutils"
)

const (
	ER_DEFAULT_SYSTEM_USER        string = "vcap"
	ER_DIRECTOR_INFO_URL          string = "https://%s:25555/info"
	ER_BACKUP_DIR                 string = "elasticruntime"
	ER_VMS_URL                    string = "https://%s:25555/deployments/%s/vms"
	ER_DIRECTOR                   string = "DirectorInfo"
	ER_CONSOLE                    string = "ConsoledbInfo"
	ER_UAA                        string = "UaadbInfo"
	ER_CC                         string = "CcdbInfo"
	ER_MYSQL                      string = "MysqldbInfo"
	ER_NFS                        string = "NfsInfo"
	ER_BACKUP_FILE_FORMAT         string = "%s.backup"
	ER_INVALID_DIRECTOR_CREDS_MSG string = "invalid director credentials"
	ER_NO_PERSISTENCE_ARCHIVES    string = "there are no persistence stores in the list"
	ER_FILE_DOES_NOT_EXIST        string = "file does not exist"
)

const (
	IMPORT_ARCHIVE = iota
	EXPORT_ARCHIVE
)

var (
	ER_ERROR_DIRECTOR_CREDS error         = errors.New(ER_INVALID_DIRECTOR_CREDS_MSG)
	ER_ERROR_EMPTY_DB_LIST  error         = errors.New(ER_NO_PERSISTENCE_ARCHIVES)
	ER_ERROR_INVALID_PATH   *os.PathError = &os.PathError{Err: errors.New(ER_FILE_DOES_NOT_EXIST)}
)

// ElasticRuntime contains information about a Pivotal Elastic Runtime deployment
type ElasticRuntime struct {
	JsonFile          string
	SystemsInfo       map[string]SystemDump
	PersistentSystems []SystemDump
	HttpGateway       HttpGateway
	InstallationName  string
	BackupContext
	Logger lager.Logger
}

// NewElasticRuntime initializes an ElasticRuntime intance
var NewElasticRuntime = func(jsonFile string, target string, logger lager.Logger) *ElasticRuntime {
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
			Product:   "microbosh",
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
		JsonFile: jsonFile,
		BackupContext: BackupContext{
			TargetDir: target,
		},
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
			// nfsInfo,
			mysqldbInfo,
		},
		Logger: logger,
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	return context.backupRestore(EXPORT_ARCHIVE)
}

// Restore performs a restore of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Restore() (err error) {
	return context.backupRestore(IMPORT_ARCHIVE)
}

func (context *ElasticRuntime) backupRestore(action int) (err error) {
	var (
		ccStop  *CloudController
		ccStart *CloudController
		ccJobs  []string
	)

	if err = context.ReadAllUserCredentials(); err == nil && context.directorCredentialsValid() {
		context.Logger.Debug("Retrieving All CC VMs")
		if ccJobs, err = context.getAllCloudControllerVMs(); err == nil {
			context.Logger.Debug("Setting up CC jobs")
			directorInfo := context.SystemsInfo[ER_DIRECTOR]
			ccStop = NewCloudController(directorInfo.Get(SD_IP), directorInfo.Get(SD_USER), directorInfo.Get(SD_PASS), context.InstallationName, "stopped")
			ccStart = NewCloudController(directorInfo.Get(SD_IP), directorInfo.Get(SD_USER), directorInfo.Get(SD_PASS), context.InstallationName, "started")
			defer ccStart.ToggleJobs(CloudControllerJobs(ccJobs))
			ccStop.ToggleJobs(CloudControllerJobs(ccJobs))
		}
		err = context.RunDbAction(context.PersistentSystems, action)

	} else if err == nil {
		err = ER_ERROR_DIRECTOR_CREDS
	}
	return
}

func (context *ElasticRuntime) getAllCloudControllerVMs() (ccvms []string, err error) {

	context.Logger.Debug("Entering getAllCloudControllerVMs() function")
	directorInfo := context.SystemsInfo[ER_DIRECTOR]
	connectionURL := fmt.Sprintf(ER_VMS_URL, directorInfo.Get(SD_IP), context.InstallationName)
	context.Logger.Debug("getAllCloudControllerVMs() function", lager.Data{"connectionURL": connectionURL, "directorInfo": directorInfo})
	gateway := context.HttpGateway
	if gateway == nil {
		gateway = NewHttpGateway()
	}
	context.Logger.Debug("Retrieving CC vms")
	if resp, err := gateway.Get(HttpRequestEntity{
		Url:         connectionURL,
		Username:    directorInfo.Get(SD_USER),
		Password:    directorInfo.Get(SD_PASS),
		ContentType: "application/json",
	})(); err == nil {
		var jsonObj []VMObject

		context.Logger.Debug("Unmarshalling CC vms")
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

		if err = info.Error(); err == nil {
			err = context.readWriterArchive(info, context.TargetDir, action)

		} else {
			break
		}
	}

	if len(dbInfoList) <= 0 {
		err = ER_ERROR_EMPTY_DB_LIST
	}
	return
}

func (context *ElasticRuntime) getReadWriter(fpath string, action int) (rw io.ReadWriter, err error) {
	switch action {
	case IMPORT_ARCHIVE:
		var exists bool

		if exists, err = osutils.Exists(fpath); exists && err == nil {
			rw, err = os.Open(fpath)

		} else {
			var pathError os.PathError
			pathError = *ER_ERROR_INVALID_PATH
			pathError.Path = fpath
			err = &pathError
		}

	case EXPORT_ARCHIVE:
		rw, err = osutils.SafeCreate(fpath)
	}
	return
}

func (context *ElasticRuntime) readWriterArchive(dbInfo SystemDump, databaseDir string, action int) (err error) {
	var (
		archivefile io.ReadWriter
	)
	filename := fmt.Sprintf(ER_BACKUP_FILE_FORMAT, dbInfo.Get(SD_COMPONENT))
	filepath := path.Join(databaseDir, filename)

	if archivefile, err = context.getReadWriter(filepath, action); err == nil {
		err = context.importExport(archivefile, dbInfo, action)
	}
	return
}

func (context *ElasticRuntime) importExport(rw io.ReadWriter, s SystemDump, action int) (err error) {
	var pb PersistanceBackup

	if pb, err = s.GetPersistanceBackup(); err == nil {

		switch action {
		case IMPORT_ARCHIVE:
			fmt.Println("we are doing something here now")
			err = pb.Import(rw)

		case EXPORT_ARCHIVE:
			err = pb.Dump(rw)
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
