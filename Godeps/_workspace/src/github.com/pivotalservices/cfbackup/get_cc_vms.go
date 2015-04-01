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
		Job   string
		Index int
	}

	CloudControllerDeploymentParser struct {
		vms []CCJob
	}
)

func GetCCVMs(jsonObj []VMObject) ([]CCJob, error) {
	parser := &CloudControllerDeploymentParser{}
	return parser.Parse(jsonObj)
}

func (s *CloudControllerDeploymentParser) Parse(jsonObj []VMObject) ([]CCJob, error) {
	err := s.setupAndRun(jsonObj)
	return s.vms, err
}

func (s *CloudControllerDeploymentParser) setupAndRun(jsonObj []VMObject) (err error) {

	var ccjobs = make([]CCJob, 0)

	for _, vmObject := range jsonObj {
		if strings.Contains(vmObject.Job, "cloud_controller-partition") {
			ccJob := CCJob{
				Job:   vmObject.Job,
				Index: vmObject.Index,
			}
			ccjobs = append(ccjobs, ccJob)
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
