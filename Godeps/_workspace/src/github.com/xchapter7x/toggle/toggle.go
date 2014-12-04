package toggle

import (
	"fmt"
	"reflect"

	"github.com/xchapter7x/toggle/engines/localengine"
	"github.com/xchapter7x/toggle/engines/storageinterface"
)

func createReflectValueArgsArray(iargs []interface{}) (args []reflect.Value) {
	for _, arg := range iargs {
		args = append(args, reflect.ValueOf(arg))
	}
	return
}

func createInterfaceArrayFromValuesArray(responseValuesArray []reflect.Value) (responseInterfaceArray []interface{}) {
	for _, ri := range responseValuesArray {
		responseInterfaceArray = append(responseInterfaceArray, ri.Interface())
	}
	return
}

func getActivePointer(flg string, defaultFeature, newFeature interface{}) (activePointer interface{}) {
	if IsActive(flg) {
		activePointer = newFeature

	} else {
		activePointer = defaultFeature
	}
	return
}

func Flip(flg string, defaultFeature, newFeature interface{}, iargs ...interface{}) (responseInterfaceArray []interface{}) {
	args := createReflectValueArgsArray(iargs)
	ptr := getActivePointer(flg, defaultFeature, newFeature)
	responseValuesArray := reflect.ValueOf(ptr).Call(args)
	responseInterfaceArray = createInterfaceArrayFromValuesArray(responseValuesArray)
	return
}

func SetFeatureStatus(featureSignature, featureStatus string) (err error) {
	fullSignature := GetFullFeatureSignature(featureSignature)

	if _, exists := featureList[fullSignature]; exists {
		featureList[fullSignature].Status = featureStatus

	} else {
		err = fmt.Errorf("Feature toggle doesnt exist")
	}
	return
}

func IsActive(featureSignature string) (active bool) {
	fullSignature := GetFullFeatureSignature(featureSignature)

	if feature, exists := featureList[fullSignature]; !exists || feature.Status == FEATURE_OFF {
		active = false

	} else {
		active = true
	}
	return
}

type Feature struct {
	name     string
	Status   string
	filter   func(...interface{}) bool
	settings map[string]interface{}
}

func (s *Feature) UpdateStatus(newStatus string) {
	s.Status = newStatus
}

const (
	FEATURE_ON     = "true"
	FEATURE_OFF    = "false"
	FEATURE_FILTER = "filter:x:x"
)

var featureList map[string]*Feature = make(map[string]*Feature)
var namespace string
var toggleEngine storageinterface.StorageEngine

func Close() {
	featureList = make(map[string]*Feature)
	namespace = ""
	toggleEngine = nil
}

func Init(ns string, engine storageinterface.StorageEngine) {
	namespace = ns

	if engine != nil {
		toggleEngine = engine

	} else {
		toggleEngine = localengine.NewLocalEngine()
	}
}

func ShowFeatures() map[string]*Feature {
	return featureList
}

func getFeatureStatusValue(featureSignature, statusValue string) (status string) {
	var err error

	if status, err = toggleEngine.GetFeatureStatusValue(featureSignature); err != nil {
		status = statusValue
	}
	return
}

func RegisterFeature(featureSignature string) (err error) {
	err = RegisterFeatureWithStatus(featureSignature, FEATURE_OFF)
	return
}

func GetFullFeatureSignature(partialSignature string) (fullSignature string) {
	fullSignature = fmt.Sprintf("%s_%s", namespace, partialSignature)
	return
}

func RegisterFeatureWithStatus(featureSignature, statusValue string) (err error) {
	fullSignature := GetFullFeatureSignature(featureSignature)

	if _, exists := featureList[fullSignature]; !exists {
		featureList[fullSignature] = &Feature{
			name:   fullSignature,
			Status: getFeatureStatusValue(fullSignature, statusValue),
		}

	} else {
		err = fmt.Errorf("feature already registered")
	}
	return
}
