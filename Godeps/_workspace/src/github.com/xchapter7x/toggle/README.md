toggle
======

[![wercker status](https://app.wercker.com/status/9c11e691895a9782a234fcc9bb313819/m "wercker status")](https://app.wercker.com/project/bykey/9c11e691895a9782a234fcc9bb313819)


## supports env_variable backed toggling. 
## It can also be updated via a pubsub interface (tested w/ redis)
## 2 engines for toggle backing are included
	- Local Engine (only env variable backed)
	- Local Engine PubSub (env variable backed w/ updates via redis pubsub)

## can be extended to support any 3rd party backing engine













func **Flip**(flg string, defaultFeature, newFeature interface{}, iargs ...interface{}) (responseInterfaceArray []interface{}) {}


func **SetFeatureStatus**(featureSignature, featureStatus string) (err error) {}

func **IsActive**(featureSignature string) (active bool) {}

func **Close**() {}

func **Init**(ns string, engine storageinterface.StorageEngine) {}

func **ShowFeatures**() map[string]*Feature {}

func **RegisterFeature**(featureSignature string) (err error) {}

func **GetFullFeatureSignature**(partialSignature string) (fullSignature string) {}

func **RegisterFeatureWithStatus**(featureSignature, statusValue string) (err error) {}

type Feature struct {
	name     string
	Status   string
	filter   func(...interface{}) bool
	settings map[string]interface{}
}
func (s \*Feature) **UpdateStatus**(newStatus string) {}



type **DefaultEngine** struct

func (s \*DefaultEngine) **GetFeatureStatusValue**(featureSignature string) (status string, err error)


Sample usage:
(./sample/main.go)
```
package main

import (
	"fmt"

	"github.com/xchapter7x/goutil/unpack"
	"github.com/xchapter7x/toggle"
)

func TestA(s string) (r string) {
	r = fmt.Sprintln("testa", s)
	fmt.Println(r)
	return
}

func TestB(s string) (r string) {
	r = fmt.Sprintln("testb", s)
	fmt.Println(r)
	return
}

func main() {
	toggle.Init("MAINTEST", nil)
	toggle.RegisterFeature("test")
	f := toggle.Flip("test", TestA, TestB, "argstring")
	var output string
	unpack.Unpack(f, &output)
	fmt.Println(output)

}
```


```
$ test=true go run sample/main.go
testb argstring

testb argstring



$ go run sample/main.go
testa argstring

testa argstring
```
