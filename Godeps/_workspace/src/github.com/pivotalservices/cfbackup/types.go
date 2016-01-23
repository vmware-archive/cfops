package cfbackup

import (
	"io"
	"net/http"

	"github.com/pivotalservices/gtils/command"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/xchapter7x/goutil"
)

type (

	//NFSBackup - this is a nfs backup object
	NFSBackup struct {
		Caller    command.Executer
		RemoteOps remoteOpsInterface
	}

	//BackupContext - stores the base context information for a backup/restore
	BackupContext struct {
		TargetDir string
		IsS3      bool
		StorageProvider
	}

	// StorageProvider is responsible for obtaining/managing a reader/writer to
	// a storage type (eg disk/s3)
	StorageProvider interface {
		Reader(path ...string) (io.ReadCloser, error)
		Writer(path ...string) (io.WriteCloser, error)
	}

	// Tile is a deployable component that can be backed up
	Tile interface {
		Backup() error
		Restore() error
	}

	//InstallationSettings - an object to house installationsettings elements from the json
	InstallationSettings struct {
		Infrastructure Infrastructure `json:"infrastructure"`
		Products       []Products     `json:"products"`
	}

	//Infrastructure - a struct to house Infrastructure block elements from the json
	Infrastructure struct {
		Type       string            `json:"type"`
		IaaSConfig IaaSConfiguration `json:"iaas_configuration"`
	}

	//IaaSConfiguration - a struct to house the IaaSConfiguration block elements from the json
	IaaSConfiguration struct {
		SSHPrivateKey string `json:"ssh_private_key"`
	}

	// Products contains installation settings for a product
	Products struct {
		Identifer      string              `json:"identifier"`
		IPS            map[string][]string `json:"ips"`
		Jobs           []Jobs              `json:"jobs"`
		ProductVersion string              `json:"product_version"`
	}

	// Jobs contains job settings for a product
	Jobs struct {
		Identifier string       `json:"identifier"`
		Properties []Properties `json:"properties"`
		Instances  []Instances  `json:"instances"`
	}

	// Properties contains property settings for a job
	Properties struct {
		Definition string      `json:"definition"`
		Value      interface{} `json:"value"`
	}

	// Instances contains instances for a job
	Instances struct {
		Identifier string `json:"identifier"`
		Value      int    `json:"value"`
	}

	//ConfigurationParser - the parser to handle installation settings file parsing
	ConfigurationParser struct {
		installationSettings InstallationSettings
	}

	//CCJob - a cloud controller job object
	CCJob struct {
		Job   string
		Index int
	}

	//PersistanceBackup - a struct representing a persistence backup
	PersistanceBackup interface {
		Dump(io.Writer) error
		Import(io.Reader) error
	}

	stringGetterSetter interface {
		Get(string) string
		Set(string, string)
	}
	//SystemDump - definition for a SystemDump interface
	SystemDump interface {
		stringGetterSetter
		Error() error
		GetPersistanceBackup() (dumper PersistanceBackup, err error)
	}
	//SystemInfo - a struct representing a base systemdump implementation
	SystemInfo struct {
		goutil.GetSet
		systemInfo    map[string]SystemDump
		Product       string
		Component     string
		Identity      string
		Ip            string
		User          string
		Pass          string
		VcapUser      string
		VcapPass      string
		SSHPrivateKey string
	}
	//PgInfo - a struct representing a pgres systemdump implementation
	PgInfo struct {
		SystemInfo
		Database string
	}
	//MysqlInfo - a struct representing a mysql systemdump implementation
	MysqlInfo struct {
		SystemInfo
		Database string
	}
	//NfsInfo - a struct representing a nfs systemdump implementation
	NfsInfo struct {
		SystemInfo
	}
	//DirectorInfo - a struct representing a director systemdump implementation
	DirectorInfo struct {
		SystemInfo
		Database string
	}
	//SystemsInfo holds the values for all the supported SystemDump used by an installation
	SystemsInfo struct {
		SystemDumps map[string]SystemDump
	}

	connBucketInterface interface {
		Host() string
		AdminUser() string
		AdminPass() string
		OpsManagerUser() string
		OpsManagerPass() string
		Destination() string
	}

	action func() error

	actionAdaptor func(t Tile) action

	httpUploader func(conn ghttp.ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error)

	httpRequestor interface {
		Get(ghttp.HttpRequestEntity) ghttp.RequestAdaptor
		Post(ghttp.HttpRequestEntity, io.Reader) ghttp.RequestAdaptor
		Put(ghttp.HttpRequestEntity, io.Reader) ghttp.RequestAdaptor
	}

	remoteOpsInterface interface {
		UploadFile(lfile io.Reader) (err error)
		Path() string
	}
)
