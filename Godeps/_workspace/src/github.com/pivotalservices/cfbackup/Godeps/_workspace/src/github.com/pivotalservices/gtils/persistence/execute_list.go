package persistence

import (
	"bytes"

	"github.com/pivotalservices/gtils/command"
)

func execute_list(callList []string, caller command.Executer) (err error) {
	var byteWriter bytes.Buffer

	for _, callstring := range callList {

		if err = caller.Execute(&byteWriter, callstring); err != nil {
			break
		}
	}
	return
}
