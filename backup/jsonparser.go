package backup

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "os"
  "strings"
)

type JsonObject struct {
  Infrastructure Infrastructure `json:"infrastructure"`
  Products []Products `json:"products"`
}

type Infrastructure struct {
  Type string `json:"type"`
}

type Products struct {
  Type string `json:"type"`
  IPS map[string][]string `json:"ips"`
  Jobs []Jobs `json:"jobs"`
}

type IPS struct {
  identifier string
  value []string
}

type Jobs struct {
  Type string `json:"type"`
  Properties []Properties `json:"properties"`
}

type Properties struct {
  Definition string `json:"definition"`
  Value Value `json:"value"`
}

type Value struct {
  Identity string `json:"identity"`
  Password string `json:"password"`
}

func getPassword(args []string) string {

  var password string

  fmt.Println(args)

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
	fmt.Printf("%+v", jsonObject)

  for _, product := range jsonObject.Products {
    if(product.Type == args[1]) {
      for _, job := range product.Jobs {
        fmt.Println(job.Type)
        if(job.Type == args[2]) {
          for _, property := range job.Properties {
            if(property.Value.Identity == args[3]) {
              password = property.Value.Password
              break
            }
          }
        }
      }
    }
  }

  fmt.Println(password)
  return password
}

func  getIP(args []string) string {
  var ip string

  fmt.Println(args)

  jsonfile := args[0]

  file, e := ioutil.ReadFile(jsonfile)

  if e != nil {
      fmt.Printf("File error: %v\n", e)
      os.Exit(1)
  }

  var jsonObject JsonObject
  json.Unmarshal(file, &jsonObject)

  for _, product := range jsonObject.Products {
    if(product.Type == args[1]) {
      for k, vals := range product.IPS {
        if(strings.Contains(k, args[2])) {
          ip = vals[0]
          break
        }
      }
    }
  }

  fmt.Println(ip)
  return ip
}
