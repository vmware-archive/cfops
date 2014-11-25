package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getPassword(args []string) string {

	var password string

	jsonfile := args[0]

	file, e := ioutil.ReadFile(jsonfile)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsonObject InstallationObject
	err := json.Unmarshal(file, &jsonObject)
	if err != nil {
		fmt.Println("error:", err)
	}

	if jsonObject.Infrastructure.Type == "vsphere" {

	} else if jsonObject.Infrastructure.Type == "vcloud" {
		if args[1] == "microbosh" {
			args[1] = "microbosh-vcloud"
		}
	}

	for _, product := range jsonObject.Products {
		if product.Type == args[1] {
			for _, job := range product.Jobs {
				if job.Type == args[2] {
					for _, property := range job.Properties {
						switch property.Value.(type) {
						case map[string]interface{}:
							propertyValue := property.Value.(map[string]interface{})
							field := propertyValue["identity"]
							value := propertyValue["password"]
							if field == args[3] {
								password = value.(string)
								break
							}
						default:
							fmt.Println("unknown")
						}
					}
				}
			}
		}
	}

	return password
}

func getIP(args []string) string {
	var ip string

	jsonfile := args[0]

	file, e := ioutil.ReadFile(jsonfile)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsonObject InstallationObject
	json.Unmarshal(file, &jsonObject)

	for _, product := range jsonObject.Products {
		if product.Type == args[1] {
			for k, vals := range product.IPS {
				if strings.Contains(k, args[2]) {
					ip = vals[0]
					break
				}
			}
		}
	}

	return ip
}

func getDeploymentObject(contents []byte) []DeploymentObject {
	var jsonObject []DeploymentObject
	err := json.Unmarshal(contents, &jsonObject)
	if err != nil {
		fmt.Println("error", err)
	}

	return jsonObject
}

func getVMSObject(contents []byte) []VMObject {
	var vmObjects []VMObject
	err := json.Unmarshal(contents, &vmObjects)
	if err != nil {
		fmt.Println("error", err)
	}

	return vmObjects
}

func getEventsObject(contents []byte) EventObject {
	var eventObject EventObject
	err := json.Unmarshal(contents, &eventObject)
	if err != nil {
		fmt.Println("error", err)
	}

	return eventObject
}
