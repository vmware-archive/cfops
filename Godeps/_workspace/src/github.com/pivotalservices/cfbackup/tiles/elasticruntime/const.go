package elasticruntime

import (
	"errors"
	"os"
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
