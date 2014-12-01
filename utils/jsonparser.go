package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetPassword(args []string) string {

	var password string

	jsonfile := args[0]

	contents := getFileContents(jsonfile)

	var jsonObject InstallationObject
	GetJSONObject(contents, &jsonObject)

	if jsonObject.Infrastructure.Type == "vcloud" {
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

func GetIP(args []string) string {
	var ip string

	jsonfile := args[0]

	contents := getFileContents(jsonfile)

	var jsonObject InstallationObject
	GetJSONObject(contents, &jsonObject)

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

func getFileContents(jsonfile string) (contents []byte) {
	contents, e := ioutil.ReadFile(jsonfile)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	return contents
}

func GetJSONObject(contents []byte, response interface{}) {
	err := json.Unmarshal(contents, &response)
	if err != nil {
		fmt.Println("error", err)
	}
}
