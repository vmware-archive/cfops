package backup

import "os"

var defaultFileMode os.FileMode = 0777

func CreateDirectoriesAdaptor(makeDirectoryFunctor func(string, os.FileMode) error) func(directoryList ...string) (err error) {
	return func(directoryList ...string) (err error) {
		err = MultiDirectoryCreate(directoryList, makeDirectoryFunctor)
		return
	}
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
