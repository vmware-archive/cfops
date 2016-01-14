package cfbackup

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/xchapter7x/goutil/itertools"
	"gopkg.in/yaml.v1"
)

type cc struct {
	Db_encryption_key string
}

type property struct {
	Cc cc
}

type job struct {
	Name       string
	Properties property
}

type jobs []job

type yamlkey struct {
	Jobs jobs
}

func (s yamlkey) EncryptionKey() (key string, err error) {
	if job, err := s.Jobs.Contains("cloud_controller"); err == nil {
		key = job.Properties.Cc.Db_encryption_key
	}
	if key == "" {
		err = fmt.Errorf("empty key error")
	}
	return
}

//ExtractEncryptionKey - extract the encryptionKey from a deployment archive
func ExtractEncryptionKey(dest io.Writer, deploymentDir string) (err error) {
	var flist []os.FileInfo

	if flist, err = ioutil.ReadDir(deploymentDir); err == nil {
		yamlfilename := getYamlFilename(flist)
		yamlfilepath := path.Join(deploymentDir, yamlfilename)
		err = writeKey(dest, yamlfilepath)
	}
	return
}

func (jobs jobs) Contains(value string) (job, error) {
	for p, v := range jobs {
		if strings.Contains(v.Name, value) {
			return jobs[p], nil
		}
	}
	return job{}, fmt.Errorf("job not found")
}

func namefilter(i, v interface{}) (ok bool) {
	file := v.(os.FileInfo)
	filename := file.Name()
	ok = (strings.HasPrefix(filename, "cf-") && strings.HasSuffix(filename, ".yml"))
	return
}

func getYamlFilename(filelist []os.FileInfo) (filename string) {
	var (
		file os.FileInfo
		idx  int
	)

	if out := itertools.Filter(filelist, namefilter); len(out) > 0 {
		itertools.PairUnPack(<-out, &idx, &file)
		filename = file.Name()
	}
	return
}

func writeKey(dest io.Writer, yamlfilepath string) (err error) {
	var encryptionKey string

	if encryptionKey, err = getKeyFromFile(yamlfilepath); err == nil {
		_, err = dest.Write([]byte(encryptionKey))
	}
	return
}

func getKeyFromFile(yamlfilepath string) (encryptionKey string, err error) {
	var filebytes []byte
	keyparse := yamlkey{}

	if filebytes, err = ioutil.ReadFile(yamlfilepath); err == nil {
		err = yaml.Unmarshal(filebytes, &keyparse)
		encryptionKey, err = keyparse.EncryptionKey()
	}
	return
}
