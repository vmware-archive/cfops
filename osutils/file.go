package osutils

import (
	"os"
	"path"
	"path/filepath"
)

// SafeCreate creates a file, creating parent directories if needed
func SafeCreate(name ...string) (file *os.File, err error) {
	p, e := ensurePath(path.Join(name...))
	if e != nil {
		return nil, e
	}
	return os.Create(p)
}

func ensurePath(path string) (string, error) {
	base := filepath.Dir(path)
	e, _ := Exists(base)
	if e {
		return path, nil
	}

	// Create missing directory recursively
	err := os.MkdirAll(base, 0777)
	return path, err
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func OpenFile(name ...string) (file *os.File, err error) {
	p, e := ensurePath(path.Join(name...))
	if e != nil {
		return nil, e
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return f, err
}
