package cfbackup

import (
	"errors"
	"os"

	"github.com/pivotalservices/gtils/command"
)

const (
	//AccessKeyIDVarname - s3 key flag
	AccessKeyIDVarname = "S3_ACCESS_KEY_ID"
	//SecretAccessKeyVarname - s3 secret key
	SecretAccessKeyVarname = "S3_SECRET_ACCESS_KEY"
	//BucketNameVarname - bucket name var flag
	BucketNameVarname = "S3_BUCKET_NAME"
	//S3Domain - s3 domain value
	S3Domain = "S3_DOMAIN"
	//IsS3Varname - s3 persistence true|false
	IsS3Varname = "S3_ACTIVE"

	//NfsDirPath - this is where the nfs store lives
	NfsDirPath string = "/var/vcap/store"
	//NfsArchiveDir - this is the archive dir name
	NfsArchiveDir string = "shared"
	//NfsDefaultSSHUser - this is the default ssh user for nfs
	NfsDefaultSSHUser string = "vcap"

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

	//BackupLoggerName --
	BackupLoggerName = "Backup"
	//RestoreLoggerName --
	RestoreLoggerName = "Restore"

	//SDProduct --
	SDProduct string = "Product"
	//SDComponent --
	SDComponent string = "Component"
	//SDIdentity --
	SDIdentity string = "Identity"
	//SDIP --
	SDIP string = "Ip"
	//SDUser --
	SDUser string = "User"
	//SDPass --
	SDPass string = "Pass"
	//SDVcapUser --
	SDVcapUser string = "VcapUser"
	//SDVcapPass --
	SDVcapPass string = "VcapPass"
	
	//SDIdentifier
	SDIdentifier string = "Identifier"
)

const (
	//ImportArchive --
	ImportArchive = iota
	//ExportArchive --
	ExportArchive
)

var (
	//NfsNewRemoteExecuter - this is a function which is able to execute a remote command against the nfs server
	NfsNewRemoteExecuter = command.NewRemoteExecutor

	//ErrERDirectorCreds - error for director creds
	ErrERDirectorCreds = errors.New(ERInvalidDirectorCredsMsg)
	//ErrEREmptyDBList - error for db list empty
	ErrEREmptyDBList = errors.New(ERNoPersistenceArchives)
	//ErrERInvalidPath - invalid filepath error
	ErrERInvalidPath = &os.PathError{Err: errors.New(ERFileDoesNotExist)}
	//ErrERDBBackup - error for db backup failures
	ErrERDBBackup = errors.New(ErrERDBBackupFailure)

	//TileRestoreAction -- executes a restore action on the given tile
	TileRestoreAction = func(t Tile) func() error {
		return t.Restore
	}
	//TileBackupAction - executes a backup action on a given tile
	TileBackupAction = func(t Tile) func() error {
		return t.Backup
	}
)
