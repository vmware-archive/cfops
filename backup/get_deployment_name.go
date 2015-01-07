package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type (
	DeploymentObject struct {
		Name string
	}

	DeploymentParser struct {
		name string
	}
)

func GetDeploymentName(jsonObj []DeploymentObject) (string, error) {
	parser := &DeploymentParser{}
	return parser.Parse(jsonObj)
}

func (s *DeploymentParser) Parse(jsonObj []DeploymentObject) (name string, err error) {
	if err := s.setupAndRun(jsonObj); err == nil {
		name = s.name
	}
	return
}

func (s *DeploymentParser) setupAndRun(jsonObj []DeploymentObject) (err error) {
	for _, deploymentObject := range jsonObj {
		if strings.Contains(deploymentObject.Name, "cf-") {
			fmt.Println(fmt.Sprintf("CF deployment Name : %s", deploymentObject.Name))
			s.name = deploymentObject.Name
			return nil
		}
	}

	return fmt.Errorf("could not find deployment name")
}

func ReadAndUnmarshalDeploymentName(src io.Reader) (jsonObj []DeploymentObject, err error) {
	var contents []byte

	if contents, err = ioutil.ReadAll(src); err == nil {
		err = json.Unmarshal(contents, &jsonObj)
	}
	return
}
