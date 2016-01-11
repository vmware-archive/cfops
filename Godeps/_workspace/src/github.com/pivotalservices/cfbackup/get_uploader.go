package cfbackup

import (
	ghttp "github.com/pivotalservices/gtils/http"
)

func GetUploader(backupContext BackupContext) (uploader httpUploader) {
	uploader = ghttp.MultiPartUpload

	if backupContext.IsS3 {
		uploader = ghttp.LargeMultiPartUpload
	}
	return
}
