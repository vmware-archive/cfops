package cfbackup

// NewBackupContext initializes a BackupContext
func NewBackupContext(targetDir string, env map[string]string) (backupContext BackupContext) {
	backupContext = BackupContext{
		TargetDir: targetDir,
	}
	if useS3(env) {
		backupContext.StorageProvider = NewS3Provider(env[S3Domain], env[AccessKeyIDVarname], env[SecretAccessKeyVarname], env[BucketNameVarname])
	} else {
		backupContext.StorageProvider = NewDiskProvider()
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
