package main

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cfops/opsmanager/api"
	"github.com/pivotalservices/cfops/opsmanager/models"
)

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
	// Send api_version request to server
	//
	gateway := api.NewOpsManagerGateway(url, username, password)
	apiResponse := gateway.GetAPIVersion(res)

	if apiResponse.IsError() {
		println(apiResponse.ErrorCode)
	}

	//
	// Process response
	//
	println("")
	fmt.Println("API Response Status:", apiResponse.StatusCode)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Ops Manager API Version:")
	fmt.Println(res.Version)
	println("")

	//
	// Send installation_settings request to server
	//
	var installation *models.InstallationSettings
	apiResponse = gateway.GetInstallationSettings(&installation)

	if apiResponse.IsError() {
		println(apiResponse.ErrorCode)
	}

	//
	// Process response
	//
	println("")
	fmt.Println("API Response Status:", apiResponse.StatusCode)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Ops Manager Installation:")
	fmt.Println(installation)
	println("")
	fmt.Println("Ops Manager Director:")
	fmt.Println("IP Address:", getIPAddress(installation, "microbosh", "director"))
	fmt.Println("Password:", getPassword(installation, "microbosh", "director", "director"))
	println("")
	fmt.Println("Ops Manager DEA:")
	fmt.Println("IP Address:", getIPAddress(installation, "cf", "cf"))
	println("")
	fmt.Println("Ops Manager CCDB:")
	fmt.Println("Password:", getPassword(installation, "cf", "ccdb", "admin"))
}

func getIPAddress(installation *models.InstallationSettings, productType, productName string) (ip string) {
	for _, product := range installation.Products {
		if product.Type == productType {
			for k, vals := range product.IPS {
				if strings.Contains(k, productName) {
					return vals[0]
				}
			}
		}
	}
	return
}

func getPassword(installation *models.InstallationSettings, productType, jobType, identity string) (password string) {
	for _, product := range installation.Products {
		if product.Type == productType {
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
