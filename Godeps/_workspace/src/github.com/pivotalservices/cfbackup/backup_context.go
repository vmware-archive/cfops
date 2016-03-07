package cfbackup

import (
	"fmt"

	"github.com/xchapter7x/lo"
)

// NewBackupContext initializes a BackupContext
func NewBackupContext(targetDir string, env map[string]string, cryptKey string) (backupContext BackupContext) {
	backupContext = BackupContext{
		TargetDir: targetDir,
	}
	if useS3(env) {
		backupContext.StorageProvider = NewS3Provider(env[S3Domain], env[AccessKeyIDVarname], env[SecretAccessKeyVarname], env[BucketNameVarname])
		backupContext.IsS3 = true
	} else {
		backupContext.StorageProvider = NewDiskProvider()
	}

	if isValidCryptKey(cryptKey) {
		var err error

		if backupContext.StorageProvider, err = NewEncryptedStorageProvider(backupContext.StorageProvider, cryptKey); err != nil {
			lo.G.Error("something went wrong when applying encryption to storage provider: ", err)
			panic(err)
		}
	}
	return
}

func isValidCryptKey(key string) (valid bool) {
	l := len(key)
	valid = l != 0 && (l == 16 || l == 24 || l == 32)

	if l > 0 && !valid {
		s := fmt.Sprintf("key is not valid failing now. length should be 16,24 or 32: len is %v", l)
		lo.G.Panic(s)
		panic(s)
	}
	return
}

func useS3(env map[string]string) bool {
	_, akid := env[AccessKeyIDVarname]
	_, sak := env[SecretAccessKeyVarname]
	_, bn := env[BucketNameVarname]
	s3val, is := env[IsS3Varname]
	isS3 := (s3val == "true")
	return (akid && sak && bn && is && isS3)
}
