package cfbackup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type (
	VMObject struct {
		Job string
	}

	CloudControllerDeploymentParser struct {
		vms []string
	}
)

func GetCCVMs(jsonObj []VMObject) ([]string, error) {
	parser := &CloudControllerDeploymentParser{}
	return parser.Parse(jsonObj)
}

func (s *CloudControllerDeploymentParser) Parse(jsonObj []VMObject) ([]string, error) {
	err := s.setupAndRun(jsonObj)
	return s.vms, err
}

func (s *CloudControllerDeploymentParser) setupAndRun(jsonObj []VMObject) (err error) {

	var ccjobs = make([]string, 0)

	for _, vmObject := range jsonObj {
		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
			ccjobs = append(ccjobs, vmObject.Job)
		}
	}

	if len(ccjobs) == 0 {
		return fmt.Errorf("no cc jobs found")
	}
	s.vms = ccjobs
	return nil
}

func ReadAndUnmarshalVMObjects(src io.Reader) (jsonObj []VMObject, err error) {
	var contents []byte

	if contents, err = ioutil.ReadAll(src); err == nil {
		err = json.Unmarshal(contents, &jsonObj)
	}
	return
}
