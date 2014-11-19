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

	var jsonObject JsonObject
	err := json.Unmarshal(file, &jsonObject)
	if err != nil {
		fmt.Println("error:", err)
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

	var jsonObject JsonObject
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
