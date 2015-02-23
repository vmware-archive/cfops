package persistence

import (
	"fmt"
	// "io/ioutil"
	"os"

	"github.com/pivotalservices/gtils/command"
)

func execute_list(callList []string, caller command.Executer) (err error) {
	for _, callstring := range callList {

		if err = caller.Execute(os.Stdout, callstring); err != nil {
			fmt.Println()
			fmt.Printf("error executing command: %s, %v", callstring, err)
			fmt.Println()
			fmt.Printf("failed to execute::%v", err.Error())
			fmt.Println()
			break
		}
	}
	return
}
