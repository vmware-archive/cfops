package http

import "os"

//GetFileSize - gets the filesize of a given path
func GetFileSize(filename string) (fileSize int64) {
	var (
		fileInfo os.FileInfo
		err      error
		file     *os.File
	)

	if file, err = os.Open(filename); err == nil {
		fileInfo, err = file.Stat()
		fileSize = fileInfo.Size()

	} else {
		fileSize = -1
	}
	return
}
