package backup

import "os"

var defaultFileMode os.FileMode = 0777

func CreateDirectoriesAdaptor(backupDir string, deploymentDir string, databaseDir string, nfsDir string) (err error) {
	directoryList := []string{backupDir, deploymentDir, databaseDir, nfsDir}
	err = MultiDirectoryCreate(directoryList, os.MkdirAll)
	return
}

func MultiDirectoryCreate(directoryList []string, makeDirectoryFunctor func(string, os.FileMode) error) (err error) {

	for _, dirname := range directoryList {
		create_err := DirectoryCreate(dirname, makeDirectoryFunctor)

		if create_err != nil {
			err = create_err
			break
		}
	}
	return
}

func DirectoryCreate(dirname string, makeDirectoryFunctor func(string, os.FileMode) error) (err error) {
	err = makeDirectoryFunctor(dirname, defaultFileMode)
	return
}
