package cfbackup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type (
	//VMObject - a struct representing a vm
	VMObject struct {
		Job   string
		Index int
	}
	//CloudControllerDeploymentParser - a struct which will handle the parsing of deployments
	CloudControllerDeploymentParser struct {
		vms []CCJob
	}
)

//GetCCVMs - a function to get a list of ccjobs
func GetCCVMs(jsonObj []VMObject) ([]CCJob, error) {
	parser := &CloudControllerDeploymentParser{}
	return parser.Parse(jsonObj)
}

//Parse - a method which will parse a given vmobject array
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

//ReadAndUnmarshalVMObjects - read the io.reader and unmarshal its contents into an vmobject array
func ReadAndUnmarshalVMObjects(src io.Reader) (jsonObj []VMObject, err error) {
	var contents []byte

	if contents, err = ioutil.ReadAll(src); err == nil {
		err = json.Unmarshal(contents, &jsonObj)
	}
	return
}
