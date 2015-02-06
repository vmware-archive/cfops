package persistence

import (
	"io/ioutil"

	"github.com/pivotalservices/gtils/command"
)

func execute_list(callList []string, caller command.Executer) (err error) {
	for _, callstring := range callList {

		if err = caller.Execute(ioutil.Discard, callstring); err != nil {
			break
		}
	}
	return
}
