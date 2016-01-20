package cfbackup

import (
	ghttp "github.com/pivotalservices/gtils/http"
)

//GetUploader - returns an uploader from a given backup context
func GetUploader(backupContext BackupContext) (uploader httpUploader) {
	uploader = ghttp.LargeMultiPartUpload

	if backupContext.IsS3 {
		uploader = ghttp.MultiPartUpload
	}
	return
}
