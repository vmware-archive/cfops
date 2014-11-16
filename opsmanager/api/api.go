package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudfoundry-community/gogobosh/net"
	"github.com/pivotalservices/cfops/opsmanager/models"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type ResponseUserAgent struct {
	Useragent string `json:"user-agent"`
}

func main() {

	username := "admin"
	password := "admin"

	url := "https://opsmgr.cf.nono.com/api/"

	fmt.Println("API Version URL: ", url)

	//
	// Struct to hold response data
	//
	res := struct {
		Version string `json:"version"`
	}{}

	//
	// Send request to server
	//
	gateway := net.NewDirectorGateway()
	request, _ := gateway.NewRequest("GET", url+"api_version", username, password, nil)
	headers, apiResponse := gateway.PerformRequestForJSONResponse(request, &res)

	if apiResponse.IsError() {
		log.Fatal(apiResponse.ErrorCode)
	}

	//
	// Process response
	//
	println("")
	fmt.Println("API Response Status:", apiResponse.StatusCode)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("HTTP Request Header")
	fmt.Println(request.HttpReq.Header)
	fmt.Println("HTTP Response Headers")
	fmt.Println(headers)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Ops Manager API Version:")
	fmt.Println(res.Version)
	println("")

	var jsonObject *models.JsonObject
	request, _ = gateway.NewRequest("GET", url+"installation_settings", username, password, nil)
	headers, apiResponse = gateway.PerformRequestForJSONResponse(request, &jsonObject)
	//
	// Process response
	//
	println("")
	fmt.Println("API Response Status:", apiResponse.StatusCode)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("HTTP Request Header")
	fmt.Println(request.HttpReq.Header)
	fmt.Println("HTTP Response Headers")
	fmt.Println(headers)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Ops Manager Installation:")
	fmt.Println(jsonObject)
	println("")
	fmt.Println("Ops Manager Director:")
	fmt.Println("IP Address:", getDirectorIPAddress(jsonObject))
	fmt.Println("Password:", getPassword(jsonObject, "director", "director"))
	println("")
}

func getDirectorIPAddress(jsonObject *models.JsonObject) (ip string) {
	for _, product := range jsonObject.Products {
		if product.Type == "microbosh" {
			for k, vals := range product.IPS {
				if strings.Contains(k, "director") {
					return vals[0]
				}
			}
		}
	}
	return
}

func getPassword(jsonObject *models.JsonObject, jobType, identity string) (password string) {
	for _, product := range jsonObject.Products {
		if product.Type == "microbosh" {
			for _, job := range product.Jobs {
				if job.Type == jobType {
					for _, property := range job.Properties {
						switch property.Value.(type) {
						case map[string]interface{}:
							propertyValue := property.Value.(map[string]interface{})
							field := propertyValue["identity"]
							value := propertyValue["password"]
							if field == identity {
								return value.(string)
							}
						}
					}
				}
			}
		}
	}
	return
}
